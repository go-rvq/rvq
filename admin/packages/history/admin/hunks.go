package admin

import (
	"fmt"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/sunfmin/reflectutils"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// hunk is one change region turning the field's current value into the target
// revision's value: Old is the current-side text, New the target-side text.
type hunk struct {
	Index int
	Old   string
	New   string
}

// diffOps is the op list turning current into target (deterministic, so render
// and apply agree on hunk ordering). For HTML the diff runs over whole top-level
// blocks (each `<p>…</p>` / text node is one unit), so a hunk is a run of
// complete blocks — `<p>asd,</p>`, never a `</p><p>` shard.
func diffOps(current, target string, htmlMode bool) []diffmatchpatch.Diff {
	d := diffmatchpatch.New()
	if htmlMode {
		return diffTokens(d, htmlBlocks(current), htmlBlocks(target))
	}
	diffs := d.DiffMain(current, target, false)
	return d.DiffCleanupSemantic(diffs)
}

// htmlBlocks parses an HTML fragment and returns its top-level nodes serialized,
// one string per node (a block element with its whole subtree, or a text run).
// Diffing over these keeps a change a whole block; concatenating them back is
// the exact fragment, so applyHunks stays lossless. Falls back to the raw string
// when it does not parse.
func htmlBlocks(s string) []string {
	nodes, err := html.ParseFragment(strings.NewReader(s), &html.Node{
		Type: html.ElementNode, Data: "body", DataAtom: atom.Body,
	})
	if err != nil || len(nodes) == 0 {
		return []string{s}
	}
	toks := make([]string, 0, len(nodes))
	for _, n := range nodes {
		var b strings.Builder
		if html.Render(&b, n) == nil {
			toks = append(toks, b.String())
		}
	}
	if len(toks) == 0 {
		return []string{s}
	}
	return toks
}

// diffTokens diffs two token slices using the rune-mapping trick (each distinct
// token ↦ one rune) so diffmatchpatch's character diff operates on whole tokens;
// the result carries the token text.
func diffTokens(d *diffmatchpatch.DiffMatchPatch, curToks, tgtToks []string) []diffmatchpatch.Diff {
	arr := []string{}
	index := map[string]rune{}
	encode := func(toks []string) []rune {
		rs := make([]rune, len(toks))
		for i, tk := range toks {
			r, ok := index[tk]
			if !ok {
				r = rune(len(arr))
				index[tk] = r
				arr = append(arr, tk)
			}
			rs[i] = r
		}
		return rs
	}
	r1, r2 := encode(curToks), encode(tgtToks)
	diffs := d.DiffCleanupSemantic(d.DiffMainRunes(r1, r2, false))
	out := make([]diffmatchpatch.Diff, len(diffs))
	for i, df := range diffs {
		var b strings.Builder
		for _, r := range df.Text {
			b.WriteString(arr[r])
		}
		out[i] = diffmatchpatch.Diff{Type: df.Type, Text: b.String()}
	}
	return out
}

// fieldHunks groups the diff between current and target into change hunks.
func fieldHunks(current, target string, htmlMode bool) []hunk {
	var hunks []hunk
	inHunk := false
	for _, df := range diffOps(current, target, htmlMode) {
		if df.Type == diffmatchpatch.DiffEqual {
			inHunk = false
			continue
		}
		if !inHunk {
			hunks = append(hunks, hunk{Index: len(hunks)})
			inHunk = true
		}
		i := len(hunks) - 1
		switch df.Type {
		case diffmatchpatch.DiffDelete:
			hunks[i].Old += df.Text
		case diffmatchpatch.DiffInsert:
			hunks[i].New += df.Text
		}
	}
	return hunks
}

// applyHunks rebuilds the field value from current, adopting the target only for
// the selected hunks (git checkout -p): a selected hunk takes the target side,
// an unselected one keeps the current side.
func applyHunks(current, target string, selected map[int]bool, htmlMode bool) string {
	var b strings.Builder
	hunkIdx := -1
	inHunk := false
	for _, df := range diffOps(current, target, htmlMode) {
		if df.Type == diffmatchpatch.DiffEqual {
			b.WriteString(df.Text)
			inHunk = false
			continue
		}
		if !inHunk {
			hunkIdx++
			inHunk = true
		}
		sel := selected[hunkIdx]
		switch df.Type {
		case diffmatchpatch.DiffDelete:
			// text present now, absent in target: keep it unless reverting this hunk
			if !sel {
				b.WriteString(df.Text)
			}
		case diffmatchpatch.DiffInsert:
			// text absent now, present in target: add it only when reverting this hunk
			if sel {
				b.WriteString(df.Text)
			}
		}
	}
	return b.String()
}

// RevertFieldContent sets a field's whole content and records the change (via
// the editing pipeline, so activity + history both log it). Used by the
// hunk-selection partial revert. Only for fields that accept partial revert.
func (mh *ModelHistory) RevertFieldContent(obj any, field, value string, ctx *web.EventContext) error {
	if !mh.AcceptsPartial(field) {
		return fmt.Errorf("history: field %q does not accept partial revert", field)
	}
	if err := reflectutils.Set(obj, field, value); err != nil {
		return err
	}
	return mh.saveAndCapture(obj, ctx)
}

func (mh *ModelHistory) revertHunksEventName() string { return "history_revert_hunks_" + mh.table }

// fieldStringValue is the live field value of the loaded record as a string.
func fieldStringValue(obj any, field string) string {
	v, err := reflectutils.Get(obj, field)
	if err != nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

// hunkSelectPanel renders, for a partial-capable field, the change hunks between
// the current value and the revision's value, each with a checkbox, plus an
// apply button that reverts only the selected hunks.
func (mh *ModelHistory) hunkSelectPanel(recordKey, field string, current, target string, targetHash string, ctx *web.EventContext) h.HTMLComponent {
	msgr := getMessages(ctx.Context())
	htmlMode := mh.IsHTML(field)
	hunks := fieldHunks(current, target, htmlMode)
	if len(hunks) == 0 {
		return v.VAlert(h.Text(msgr.NoHunks)).Type("info").Variant(v.VariantTonal).Density(v.DensityCompact)
	}

	// For HTML fields render the hunk as HTML (so markup shows as formatting, not
	// as literal `</p><p>` text); for plain text escape it.
	cell := func(s string) h.HTMLComponent {
		if htmlMode {
			return h.RawHTML(s)
		}
		return h.Text(s)
	}

	rows := make(h.HTMLComponents, 0, len(hunks))
	for _, hk := range hunks {
		var change h.HTMLComponents
		if hk.Old != "" {
			change = append(change, h.Div(cell(hk.Old)).Style("background-color:#ffeef0;text-decoration:line-through;").Class("pa-1"))
		}
		if hk.New != "" {
			change = append(change, h.Div(cell(hk.New)).Style("background-color:#e6ffed;").Class("pa-1"))
		}
		rows = append(rows, h.Div(
			v.VCheckbox().
				Attr("v-model", "locals.hunks").Attr(":value", hk.Index).
				HideDetails(true).Density(v.DensityCompact),
			h.Div(change...).Class("flex-1-1-0"),
		).Class("d-flex align-start ga-2 mb-2"))
	}

	apply := v.VBtn(msgr.RevertSelected).
		Color("warning").Variant(v.VariantTonal).PrependIcon("mdi-history").
		Attr(":disabled", "!locals.hunks.length").
		Attr("@click", web.Plaid().
			EventFunc(mh.revertHunksEventName()).
			Query("field", field).
			Query("hash", targetHash).
			Query("hunks", web.Var("locals.hunks")).
			Go())

	return web.Scope(
		v.VCard(
			v.VCardTitle(h.Text(msgr.PartialRevert)),
			v.VCardText(rows...),
			v.VCardActions(apply),
		).Variant(v.VariantOutlined).Class("mt-4"),
	).LocalsInit("{ hunks: [] }")
}

// revertHunksEvent reverts only the selected hunks of a field to the target
// revision, recomputing the hunks from the live current value + target so the
// ordering matches what was shown.
func (mh *ModelHistory) revertHunksEvent(ctx *web.EventContext) (r web.EventResponse, err error) {
	recordKey := parentRecordKey(ctx)
	field := ctx.R.FormValue("field")
	hash, err := decodeHash(ctx.R.FormValue("hash"))
	if err != nil {
		return
	}

	rev, err := mh.Revision(recordKey, hash)
	if err != nil {
		return
	}
	m, err := fieldMap(rev)
	if err != nil {
		return
	}
	target := fieldValue(m, field)

	id, err := mh.mb.ParseRecordID(recordKey)
	if err != nil {
		return
	}
	obj := mh.mb.NewModel()
	id.SetTo(obj)
	if err = mh.db.First(obj).Error; err != nil {
		return
	}
	current := fieldStringValue(obj, field)

	selected := map[int]bool{}
	if e := ctx.R.ParseForm(); e == nil {
		for _, s := range ctx.R.Form["hunks"] {
			if n, ne := strconv.Atoi(s); ne == nil {
				selected[n] = true
			}
		}
	}

	value := applyHunks(current, target, selected, mh.IsHTML(field))
	if err = mh.RevertFieldContent(obj, field, value, ctx); err != nil {
		return
	}
	presets.ShowMessage(&r, getMessages(ctx.Context()).Reverted, "success")
	r.PushState = web.Location(nil).URL(mh.mb.Info().DetailingHref(recordKey))
	return
}
