package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	htmldiff "github.com/documize/html-diff"
	h "github.com/go-rvq/htmlgo"
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// htmlDiffConfig marks changes: green insertions, red struck deletions, amber
// replacements. Each change span also gets a class (hist-ins / hist-del) so the
// merged diff can be projected onto the OLD side (deletions only) and the NEW
// side (insertions only).
var htmlDiffConfig = &htmldiff.Config{
	Granularity:  5,
	InsertedSpan: []htmldiff.Attribute{{Key: "class", Val: "hist-ins"}, {Key: "style", Val: "background-color:#e6ffed;"}},
	DeletedSpan:  []htmldiff.Attribute{{Key: "class", Val: "hist-del"}, {Key: "style", Val: "background-color:#ffeef0;text-decoration:line-through;"}},
	ReplacedSpan: []htmldiff.Attribute{{Key: "class", Val: "hist-ins"}, {Key: "style", Val: "background-color:#fff5b1;"}},
	CleanTags:    []string{},
}

var (
	histInsRe = regexp.MustCompile(`(?s)<span class="hist-ins"[^>]*>.*?</span>`)
	histDelRe = regexp.MustCompile(`(?s)<span class="hist-del"[^>]*>.*?</span>`)
)

// oldSide is the merged diff with insertions dropped: the OLD content with its
// removed parts highlighted (red). newSide drops the deletions: the NEW content
// with its added parts highlighted (green).
func oldSide(merged string) string { return histInsRe.ReplaceAllString(merged, "") }
func newSide(merged string) string { return histDelRe.ReplaceAllString(merged, "") }

const comparePortal = "historyComparePortal"

// compare wraps the comparison in a portal, so the "invert" button can re-render
// it (OLD ↔ NEW swapped) in place.
func (mh *ModelHistory) compare(recordKey string, aHash, bHash histmodels.Hash, ctx *web.EventContext) (h.HTMLComponent, error) {
	inner, err := mh.compareInner(recordKey, aHash, bHash, ctx)
	if err != nil {
		return nil, err
	}
	return web.Portal(inner).Name(comparePortal), nil
}

// compareEventName is the per-model event that re-renders the comparison (used
// by the invert button).
func (mh *ModelHistory) compareEventName() string { return "history_compare_" + mh.table }

// compareEvent re-renders the comparison into its portal with the a/b the button
// passed (swapped, to invert).
func (mh *ModelHistory) compareEvent(ctx *web.EventContext) (r web.EventResponse, err error) {
	recordKey := ctx.R.FormValue("record")
	a, err := decodeHash(ctx.R.FormValue("a"))
	if err != nil {
		return
	}
	b, err := decodeHash(ctx.R.FormValue("b"))
	if err != nil {
		return
	}
	inner, err := mh.compareInner(recordKey, a, b, ctx)
	if err != nil {
		return
	}
	r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{Name: comparePortal, Body: inner})
	return
}

// compareInner shows the differences between revision aHash (OLD) and bHash
// (NEW), one collapsible panel per CHANGED field (unchanged fields are omitted).
// Each panel shows OLD (removals highlighted red) beside NEW (additions
// highlighted green): text/HTML fields through the model's DetailingBuilder with
// documize/html-diff, structured fields (struct, foreign key, slice) as a
// readable JSON summary. A top "invert" button swaps OLD and NEW.
func (mh *ModelHistory) compareInner(recordKey string, aHash, bHash histmodels.Hash, ctx *web.EventContext) (h.HTMLComponent, error) {
	revA, err := mh.Revision(recordKey, aHash)
	if err != nil {
		return nil, err
	}
	revB, err := mh.Revision(recordKey, bHash)
	if err != nil {
		return nil, err
	}
	oldObj, err := mh.applied(revA)
	if err != nil {
		return nil, err
	}
	newObj, err := mh.applied(revB)
	if err != nil {
		return nil, err
	}
	am, err := fieldMap(revA)
	if err != nil {
		return nil, err
	}
	bm, err := fieldMap(revB)
	if err != nil {
		return nil, err
	}

	msgr := getMessages(ctx.Context())

	// ?field= scopes to a single field; otherwise every versioned field. Only the
	// ones that actually changed between the two revisions are shown.
	fields := mh.resolved
	if only := fieldParam(ctx); only != "" {
		fields = []string{only}
	}

	// Invert swaps OLD and NEW: re-render the portal with a/b exchanged.
	invertBtn := v.VBtn("").Icon("mdi-swap-horizontal").
		Variant(v.VariantTonal).Size(v.SizeSmall).
		Attr("title", msgr.Invert).
		Attr("@click", web.Plaid().
			EventFunc(mh.compareEventName()).
			Query("record", recordKey).
			Query("a", bHash.String()).
			Query("b", aHash.String()).
			Go())

	// Merged toggles each field between OLD | NEW and a single merged diff.
	mergedBtn := v.VBtn("").Icon("mdi-vector-difference").
		Variant(v.VariantTonal).Size(v.SizeSmall).Class("me-2").
		Attr("title", msgr.Merged).
		Attr("@click", "locals.showMerged = !locals.showMerged")

	head := h.Div(
		h.Div(
			revInfoHeader(revA, msgr),
			h.Span(" → ").Class("mx-2 text-medium-emphasis"),
			revInfoHeader(revB, msgr),
		).Class("d-flex align-center flex-wrap"),
		h.Div(mergedBtn, invertBtn).Class("ms-auto d-flex align-center"),
	).Class("d-flex align-center justify-space-between mb-3")

	var panels h.HTMLComponents
	var open []string
	for _, f := range fields {
		if string(am[f]) == string(bm[f]) {
			continue // unchanged — omit
		}
		panels = append(panels, mh.fieldPanel(f, oldObj, newObj, am, bm, recordKey, ctx, msgr))
		open = append(open, strconv.Itoa(len(panels)-1))
	}

	if len(panels) == 0 {
		return h.Div(head, v.VAlert(h.Text(msgr.NoChanges)).
			Type("info").Variant(v.VariantTonal).Density(v.DensityComfortable)), nil
	}

	// The scope holds showMerged (the merged-view toggle); open every changed
	// field by default, each collapsible.
	return web.Scope(
		h.Div(
			head,
			v.VExpansionPanels(panels...).
				Attr("multiple", true).
				Attr(":model-value", "["+strings.Join(open, ",")+"]"),
		),
	).LocalsInit("{ showMerged: false }"), nil
}

// fieldPanel is one collapsible field diff: OLD and NEW side by side, plus — for
// a text/HTML field — a merged view below with removals highlighted red (struck)
// and additions green (documize/html-diff returns one merged result, not one per
// version). A structured field shows its before/after JSON.
func (mh *ModelHistory) fieldPanel(f string, oldObj, newObj any, am, bm map[string]json.RawMessage, recordKey string, ctx *web.EventContext, msgr *Messages) h.HTMLComponent {
	side := func(label string, body h.HTMLComponent, cls string) *h.HTMLTagBuilder {
		return h.Div(
			h.Div(h.Strong(label)).Class("text-caption text-medium-emphasis mb-1"),
			body,
		).Class("flex-1-1-0 " + cls)
	}

	var oldC, newC, mergedC h.HTMLComponent
	if mh.fieldIsSimple(f) {
		oldHTML, _ := mh.detailHTML(oldObj, recordKey, []string{f}, ctx)
		newHTML, _ := mh.detailHTML(newObj, recordKey, []string{f}, ctx)
		// The detail render wraps the value in form scaffold (hidden inputs with
		// Vue directives) that carries no content and breaks html-diff's HTML
		// parser — strip it so the diff is computed on the clean content.
		oldHTML, newHTML = stripScaffold(oldHTML), stripScaffold(newHTML)
		oldC, newC = h.RawHTML(oldHTML), h.RawHTML(newHTML)
		// documize/html-diff returns one merged result (common text + <del> +
		// <ins>). The merged view shows it whole (red removals + green additions);
		// OLD keeps deletions (red), NEW keeps insertions (green).
		if diffs, err := htmlDiffConfig.HTMLdiff([]string{oldHTML, newHTML}); err == nil && len(diffs) > 0 {
			oldC, newC = h.RawHTML(oldSide(diffs[0])), h.RawHTML(newSide(diffs[0]))
			mergedC = h.RawHTML(diffs[0])
		}
	} else {
		oldC, newC = jsonBlock(prettyJSON(am[f])), jsonBlock(prettyJSON(bm[f]))
	}

	children := h.HTMLComponents{
		// OLD | NEW, hidden while the merged view is shown.
		h.Div(
			side(msgr.Old, oldC, "pe-2").Style("border-right:1px solid rgba(0,0,0,.12)"),
			side(msgr.New, newC, "ps-2"),
		).Class("d-flex align-start").Attr("v-if", "!locals.showMerged"),
	}
	if mergedC != nil {
		children = append(children, h.Div(mergedC).Class("pa-1").Attr("v-if", "locals.showMerged"))
	}

	return v.VExpansionPanel(
		v.VExpansionPanelTitle().Children(h.Strong(f)),
		v.VExpansionPanelText().Children(children...),
	)
}

// splitFields separates the fields to compare into text/HTML (rendered through
// the DetailingBuilder) and structured (rendered as a JSON summary). With a
// single field (?field=) it classifies just that one.
func (mh *ModelHistory) splitFields(only string) (simple, structured []string) {
	fields := mh.resolved
	if only != "" {
		fields = []string{only}
	}
	for _, f := range fields {
		if mh.fieldIsSimple(f) {
			simple = append(simple, f)
		} else {
			structured = append(structured, f)
		}
	}
	return
}

// fieldIsSimple reports whether a field's value is a plain string (text or HTML)
// — the kind the detail render handles cleanly. Everything else (struct, slice,
// map, foreign key) is "structured".
func (mh *ModelHistory) fieldIsSimple(field string) bool {
	ft, ok := structFieldType(mh.mb.NewModel(), field)
	if !ok {
		return true
	}
	for ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
	}
	return ft.Kind() == reflect.String
}

// structuredRows renders each structured field's value as a readable JSON block.
func structuredRows(m map[string]json.RawMessage, fields []string) h.HTMLComponent {
	if len(fields) == 0 {
		return h.Div()
	}
	var rows h.HTMLComponents
	for _, f := range fields {
		rows = append(rows, h.Div(
			h.Div(h.Strong(f)).Class("text-caption text-medium-emphasis"),
			jsonBlock(prettyJSON(m[f])),
		).Class("mb-2"))
	}
	return h.Div(rows...).Class("mt-2")
}

// scaffoldRe matches the hidden form inputs the detail render emits around a
// field's value (Vue directives, no content).
var scaffoldRe = regexp.MustCompile(`(?s)<input\b[^>]*>`)

// stripScaffold removes those hidden inputs, leaving the field's clean content.
func stripScaffold(s string) string { return scaffoldRe.ReplaceAllString(s, "") }

func jsonBlock(s string) h.HTMLComponent {
	return h.Pre(s).Class("text-caption pa-1").
		Style("white-space:pre-wrap;overflow-x:auto;background-color:rgba(0,0,0,.03);border-radius:4px")
}

// prettyJSON renders a snapshot value for reading: a JSON string is unquoted, an
// object/array is indented, anything else is its raw form.
func prettyJSON(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	var buf bytes.Buffer
	if json.Indent(&buf, raw, "", "  ") == nil {
		return buf.String()
	}
	return string(raw)
}

// revInfoHeader shows a revision's hash, time, author and status badges.
func revInfoHeader(rev *histmodels.Revision, msgr *Messages) h.HTMLComponent {
	meta := h.HTMLComponents{
		h.Code(shortHash(rev.Hash)),
		h.Span(" · " + rev.CreatedAt.Format("2006-01-02 15:04")).Class("text-medium-emphasis"),
		h.Span(" · " + rev.Creator).Class("text-medium-emphasis"),
	}
	if rev.Published {
		meta = append(meta, v.VChip(h.Text(msgr.Published)).Color("success").Size(v.SizeXSmall).Variant(v.VariantTonal).Class("ms-1"))
	}
	if rev.Tag != "" {
		meta = append(meta, v.VChip(h.Text(msgr.Tag+": "+rev.Tag)).Size(v.SizeXSmall).Variant(v.VariantOutlined).Class("ms-1"))
	}
	if rev.AccessCount > 0 {
		meta = append(meta, v.VChip(h.Text(fmt.Sprintf("%s: %d", msgr.Accesses, rev.AccessCount))).Size(v.SizeXSmall).Variant(v.VariantOutlined).Class("ms-1"))
	}
	return h.Div(meta...).Class("text-caption d-flex align-center flex-wrap")
}

// detailComponent renders obj's given fields through the historized model's own
// DetailingBuilder. It renders as if the request were on the parent's detail
// path (/…/posts/{id}) — the builder derives its URLs/events from the request
// path, so rendering from the nested /…/revisions subpath would point them at
// the wrong route.
func (mh *ModelHistory) detailComponent(obj any, recordKey string, fields []string, ctx *web.EventContext) (h.HTMLComponent, *web.EventContext) {
	pctx := detailRequestCtx(ctx, mh.mb.Info().DetailingHref(recordKey))
	// Tell field detail components they are rendering a diff, so e.g. the TipTap
	// field emits plain HTML instead of its interactive editor.
	presets.SetDiff(pctx)
	if len(fields) == 0 {
		return h.Div(), pctx
	}
	fb := *mh.mb.Detailing().FieldsBuilder.Only(anySlice(fields)...)
	comp := fb.ToComponent(
		&presets.ToComponentOptions{}, mh.mb.Info(), obj, presets.FieldModeStack{presets.DETAIL}, pctx)
	return comp, pctx
}

func anySlice(ss []string) []any {
	r := make([]any, len(ss))
	for i, s := range ss {
		r[i] = s
	}
	return r
}

// detailRequestCtx clones ctx with the request URL set to detailPath.
func detailRequestCtx(ctx *web.EventContext, detailPath string) *web.EventContext {
	c := *ctx
	r := ctx.R.Clone(ctx.R.Context())
	u := *r.URL
	u.Path = detailPath
	u.RawQuery = ""
	r.URL = &u
	c.R = r
	return &c
}

// detailHTML renders obj's given (text/HTML) fields to an HTML string — the
// input to the visual diff. Empty when there are no such fields.
func (mh *ModelHistory) detailHTML(obj any, recordKey string, fields []string, ctx *web.EventContext) (string, error) {
	if len(fields) == 0 {
		return "", nil
	}
	comp, pctx := mh.detailComponent(obj, recordKey, fields, ctx)
	var buf bytes.Buffer
	// Render with the EventContext in the context: detail field components (e.g. a
	// media box, QMediaBoxBuilder.Write) look it up via web.MustGetEventContext,
	// which panics ("EventContext required") when it is missing.
	c := web.ContextWithEventContext(pctx.Context(), pctx)
	if err := h.Fprint(&buf, comp, c); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// fieldParam is the ?field= query — the single field the revisions view is
// scoped to (empty = the whole record).
func fieldParam(ctx *web.EventContext) string { return ctx.R.FormValue("field") }
