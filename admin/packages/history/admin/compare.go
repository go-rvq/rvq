package admin

import (
	"bytes"
	"fmt"

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
// right), each through the historized model's own DetailingBuilder, with a
// revision-info header on each. A top button reveals a third MERGED column (to
// the right of NEW) with the visual diff of the two, via documize/html-diff.
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

	// When the revisions view is scoped to one field (?field=), compare only it.
	field := fieldParam(ctx)

	oldHTML, err := mh.detailHTML(oldObj, recordKey, field, ctx)
	if err != nil {
		return nil, err
	}
	newHTML, err := mh.detailHTML(newObj, recordKey, field, ctx)
	if err != nil {
		return nil, err
	}

	merged := newHTML
	if diffs, derr := htmlDiffConfig.HTMLdiff([]string{oldHTML, newHTML}); derr == nil && len(diffs) > 0 {
		merged = diffs[len(diffs)-1]
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
				diffColumn(msgr.Old, revA, oldHTML, msgr),
				diffColumn(msgr.New, revB, newHTML, msgr),
				h.Div(diffColumn(msgr.Merged, revB, merged, msgr)).
					Attr("v-if", "locals.showMerged").
					Class("flex-1-1-0 ps-3").
					Style("border-left:1px solid rgba(0,0,0,.12)"),
			).Class("d-flex align-start"),
		),
	).LocalsInit("{ showMerged: false }"), nil
}

// diffColumn is one side of the comparison: a title, the revision-info header
// (hash, time, author, status) and the rendered detail HTML.
func diffColumn(title string, rev *histmodels.Revision, html string, msgr *Messages) h.HTMLComponent {
	return h.Div(
		h.Div(h.Strong(title)).Class("text-subtitle-2 mb-1"),
		revInfoHeader(rev, msgr),
		h.Div(h.RawHTML(html)).Class("pa-2 mt-2"),
	).Class("flex-1-1-0 px-2")
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

// detailComponent renders obj through the historized model's own
// DetailingBuilder. It renders as if the request were on the parent's detail
// path (/…/posts/{id}) — the builder derives its URLs/events from the request
// path, so rendering from the nested /…/revisions subpath would point them at
// the wrong route. When field is non-empty, only that field is rendered
// (the per-field history view).
func (mh *ModelHistory) detailComponent(obj any, recordKey, field string, ctx *web.EventContext) (h.HTMLComponent, *web.EventContext) {
	pctx := detailRequestCtx(ctx, mh.mb.Info().DetailingHref(recordKey))
	fb := mh.mb.Detailing().FieldsBuilder
	if field != "" {
		fb = *fb.Only(field)
	}
	comp := fb.ToComponent(
		&presets.ToComponentOptions{}, mh.mb.Info(), obj, presets.FieldModeStack{presets.DETAIL}, pctx)
	return comp, pctx
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

// detailHTML renders obj's detail to an HTML string (the input to the visual
// diff). When field is non-empty, only that field is rendered.
func (mh *ModelHistory) detailHTML(obj any, recordKey, field string, ctx *web.EventContext) (string, error) {
	comp, pctx := mh.detailComponent(obj, recordKey, field, ctx)
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
