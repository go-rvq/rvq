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
)

// hunk is one change region turning the field's current value into the target
// revision's value: Old is the current-side text, New the target-side text.
type hunk struct {
	Index int
	Old   string
	New   string
}

// diffOps is the op list turning current into target (deterministic, so render
// and apply agree on hunk ordering).
func diffOps(current, target string) []diffmatchpatch.Diff {
	d := diffmatchpatch.New()
	diffs := d.DiffMain(current, target, false)
	return d.DiffCleanupSemantic(diffs)
}

// fieldHunks groups the diff between current and target into change hunks.
func fieldHunks(current, target string) []hunk {
	var hunks []hunk
	inHunk := false
	for _, df := range diffOps(current, target) {
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
func applyHunks(current, target string, selected map[int]bool) string {
	var b strings.Builder
	hunkIdx := -1
	inHunk := false
	for _, df := range diffOps(current, target) {
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
	hunks := fieldHunks(current, target)
	if len(hunks) == 0 {
		return v.VAlert(h.Text(msgr.NoHunks)).Type("info").Variant(v.VariantTonal).Density(v.DensityCompact)
	}

	rows := make(h.HTMLComponents, 0, len(hunks))
	for _, hk := range hunks {
		var change h.HTMLComponents
		if hk.Old != "" {
			change = append(change, h.Div(h.Text(hk.Old)).Style("background-color:#ffeef0;text-decoration:line-through;").Class("pa-1"))
		}
		if hk.New != "" {
			change = append(change, h.Div(h.Text(hk.New)).Style("background-color:#e6ffed;").Class("pa-1"))
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

	value := applyHunks(current, target, selected)
	if err = mh.RevertFieldContent(obj, field, value, ctx); err != nil {
		return
	}
	presets.ShowMessage(&r, getMessages(ctx.Context()).Reverted, "success")
	r.PushState = web.Location(nil).URL(mh.mb.Info().DetailingHref(recordKey))
	return
}
