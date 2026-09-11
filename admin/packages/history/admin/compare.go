package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	htmldiff "github.com/documize/html-diff"
	h "github.com/go-rvq/htmlgo"
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// htmlDiffConfig marks changes in the merged view: green insertions, red struck
// deletions, amber replacements.
var htmlDiffConfig = &htmldiff.Config{
	Granularity:  5,
	InsertedSpan: []htmldiff.Attribute{{Key: "style", Val: "background-color:#e6ffed;"}},
	DeletedSpan:  []htmldiff.Attribute{{Key: "style", Val: "background-color:#ffeef0;text-decoration:line-through;"}},
	ReplacedSpan: []htmldiff.Attribute{{Key: "style", Val: "background-color:#fff5b1;"}},
	CleanTags:    []string{},
}

// compare renders the record at revision aHash (OLD, left) and bHash (NEW,
// right). Text/HTML fields go through the historized model's DetailingBuilder;
// structured fields (a struct, foreign key, slice, …) render a readable JSON
// summary instead — their detail components are interactive Vue widgets (media
// box, nested sub-resource) that do not survive being injected as diff HTML. A
// top toggle reveals a MERGED column with the visual diff (documize/html-diff)
// of the text/HTML side.
func (mh *ModelHistory) compare(recordKey string, aHash, bHash histmodels.Hash, ctx *web.EventContext) (h.HTMLComponent, error) {
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

	// ?field= scopes to a single field; otherwise all versioned fields, split
	// into text/HTML (DetailBuilder) and structured (JSON summary).
	simple, structured := mh.splitFields(fieldParam(ctx))

	oldHTML, err := mh.detailHTML(oldObj, recordKey, simple, ctx)
	if err != nil {
		return nil, err
	}
	newHTML, err := mh.detailHTML(newObj, recordKey, simple, ctx)
	if err != nil {
		return nil, err
	}

	merged := newHTML
	if len(simple) > 0 {
		if diffs, derr := htmlDiffConfig.HTMLdiff([]string{oldHTML, newHTML}); derr == nil && len(diffs) > 0 {
			merged = diffs[len(diffs)-1]
		}
	}

	msgr := getMessages(ctx.Context())

	return web.Scope(
		h.Div(
			h.Div(
				v.VBtn(msgr.Merged).
					PrependIcon("mdi-vector-difference").
					Variant(v.VariantTonal).Size(v.SizeSmall).
					Attr("@click", "locals.showMerged = !locals.showMerged"),
			).Class("d-flex justify-end mb-2"),
			h.Div(
				diffColumn(msgr.Old, revA, oldHTML, structuredRows(am, structured), msgr),
				diffColumn(msgr.New, revB, newHTML, structuredRows(bm, structured), msgr),
				h.Div(diffColumn(msgr.Merged, revB, merged, structuredDiffRows(am, bm, structured, msgr), msgr)).
					Attr("v-if", "locals.showMerged").
					Class("flex-1-1-0 ps-3").
					Style("border-left:1px solid rgba(0,0,0,.12)"),
			).Class("d-flex align-start"),
		),
	).LocalsInit("{ showMerged: false }"), nil
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

// diffColumn is one side of the comparison: a title, the revision-info header,
// the rendered text/HTML detail and the structured-field summary.
func diffColumn(title string, rev *histmodels.Revision, html string, structured h.HTMLComponent, msgr *Messages) h.HTMLComponent {
	return h.Div(
		h.Div(h.Strong(title)).Class("text-subtitle-2 mb-1"),
		revInfoHeader(rev, msgr),
		h.Div(h.RawHTML(html)).Class("pa-2 mt-2"),
		structured,
	).Class("flex-1-1-0 px-2")
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

// structuredDiffRows renders the structured fields for the merged view: an
// unchanged field is just its label, struck through; a changed one shows the
// before and after JSON.
func structuredDiffRows(am, bm map[string]json.RawMessage, fields []string, msgr *Messages) h.HTMLComponent {
	if len(fields) == 0 {
		return h.Div()
	}
	var rows h.HTMLComponents
	for _, f := range fields {
		ov, nv := string(am[f]), string(bm[f])
		if ov == nv {
			rows = append(rows, h.Div(
				h.Span(f).Style("text-decoration: line-through"),
				h.Span(" "+msgr.Unchanged).Class("text-medium-emphasis text-caption ms-1"),
			).Class("mb-1"))
			continue
		}
		rows = append(rows, h.Div(
			h.Div(h.Strong(f)).Class("text-caption text-medium-emphasis"),
			h.Div(h.Strong(msgr.Old+": "), jsonBlock(prettyJSON(am[f]))),
			h.Div(h.Strong(msgr.New+": "), jsonBlock(prettyJSON(bm[f]))),
		).Class("mb-2"))
	}
	return h.Div(rows...).Class("mt-2")
}

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
