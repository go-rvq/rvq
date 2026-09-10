package admin

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	h "github.com/go-rvq/htmlgo"
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

const historyDiffPortal = "historyDiffPortal"

// compareEventName is the per-model event that renders the diff of two chosen
// revisions into the dialog's portal.
func (mh *ModelHistory) compareEventName() string { return "history_compare_" + mh.table }

// installUI adds a "History" detailing action that opens a VXDialog with a table
// of revisions: select two to compare, or use each row's "compare with current".
func (mh *ModelHistory) installUI() {
	if !mh.mb.HasDetailing() {
		return
	}
	mh.mb.RegisterEventFunc(mh.compareEventName(), mh.compareEvent)
	mh.mb.Detailing().Action("HistoryView").
		Icon("mdi-history").
		SetI18nLabel(func(ctx context.Context) string { return getMessages(ctx).History }).
		ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
			return mh.historyDialog(id, ctx)
		})
}

func shortHash(hash []byte) string {
	s := hex.EncodeToString(hash)
	if len(s) > 14 {
		return s[:14]
	}
	return s
}

// historyDialog renders the revision table (newest first) in a VXDialog, with a
// portal for the diff. Selecting two rows enables "Compare"; a row's own button
// compares it with the current (latest) revision.
func (mh *ModelHistory) historyDialog(recordKey string, ctx *web.EventContext) (h.HTMLComponent, error) {
	msgr := getMessages(ctx.Context())
	revs, err := mh.Chain(recordKey)
	if err != nil {
		return nil, err
	}
	if len(revs) == 0 {
		return vx.VXDialog().Title(msgr.History).Width("980").SlotBody(
			v.VAlert(h.Text(msgr.HistoryEmpty)).Type("info").Variant(v.VariantTonal).Density(v.DensityComfortable),
		), nil
	}

	currentHex := hex.EncodeToString(revs[len(revs)-1].Hash)

	rows := h.HTMLComponents{}
	for i := len(revs) - 1; i >= 0; i-- {
		rows = append(rows, mh.revisionRow(recordKey, revs[i], currentHex, msgr))
	}

	table := v.VTable(
		h.Thead(h.Tr(
			h.Th(""),
			h.Th(msgr.Revision),
			h.Th(msgr.Author),
			h.Th(msgr.When),
			h.Th(msgr.Status),
			h.Th(""),
		)),
		h.Tbody(rows...),
	).Density(v.DensityComfortable)

	compareBtn := v.VBtn(msgr.Compare).
		Color(v.ColorPrimary).Variant(v.VariantTonal).Size(v.SizeSmall).
		Attr(":disabled", "locals.hsel.length!=2").
		Attr("@click", web.Plaid().
			EventFunc(mh.compareEventName()).
			Query("record", recordKey).
			Query("a", web.Var("locals.hsel[0]")).
			Query("b", web.Var("locals.hsel[1]")).
			Go())

	body := web.Scope(
		h.Div(
			h.Div(h.Text(msgr.SelectTwoHint)).Class("text-caption text-medium-emphasis mb-2"),
			table,
			h.Div(compareBtn).Class("my-3"),
			web.Portal(
				h.Div(h.Text(msgr.PickToCompare)).Class("text-medium-emphasis"),
			).Name(historyDiffPortal),
		),
	).LocalsInit("{ hsel: [] }")

	return vx.VXDialog().Title(msgr.History).Width("1100").SlotBody(body), nil
}

func (mh *ModelHistory) revisionRow(recordKey string, rev histmodels.Revision, currentHex string, msgr *Messages) h.HTMLComponent {
	hashHex := hex.EncodeToString(rev.Hash)

	status := h.HTMLComponents{}
	if rev.Published {
		status = append(status, v.VChip(h.Text(msgr.Published)).Color("success").Size(v.SizeXSmall).Variant(v.VariantTonal))
	}
	if rev.Tag != "" {
		status = append(status, v.VChip(h.Text(msgr.Tag+": "+rev.Tag)).Size(v.SizeXSmall).Variant(v.VariantOutlined))
	}
	if rev.AccessCount > 0 {
		status = append(status, v.VChip(h.Text(fmt.Sprintf("%s: %d", msgr.Accesses, rev.AccessCount))).Size(v.SizeXSmall).Variant(v.VariantOutlined))
	}

	compareCurrent := v.VBtn(msgr.CompareCurrent).
		Variant(v.VariantText).Size(v.SizeXSmall).
		Attr("@click", web.Plaid().
			EventFunc(mh.compareEventName()).
			Query("record", recordKey).
			Query("a", hashHex).
			Query("b", currentHex).
			Go())

	return h.Tr(
		h.Td(v.VCheckbox().
			Attr("v-model", "locals.hsel").Attr(":value", hashHex).
			HideDetails(true).Density(v.DensityCompact)),
		h.Td(h.Code(shortHash(rev.Hash))),
		h.Td(h.Text(rev.Creator)),
		h.Td(h.Text(rev.CreatedAt.Format("2006-01-02 15:04"))),
		h.Td(h.Div(status...).Class("d-flex flex-wrap ga-1")),
		h.Td(compareCurrent),
	)
}

// compareEvent renders the field-by-field diff of the two chosen revisions into
// the dialog's portal.
func (mh *ModelHistory) compareEvent(ctx *web.EventContext) (r web.EventResponse, err error) {
	msgr := getMessages(ctx.Context())
	recordKey := ctx.R.FormValue("record")
	a, err := hex.DecodeString(ctx.R.FormValue("a"))
	if err != nil {
		return
	}
	b, err := hex.DecodeString(ctx.R.FormValue("b"))
	if err != nil {
		return
	}
	sec, err := mh.diffSection(recordKey, a, b, msgr)
	if err != nil {
		return
	}
	r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{Name: historyDiffPortal, Body: sec})
	return
}

// diffSection renders the versioned fields like the detail form: an unchanged
// field is just its label, struck through; a changed HTML field (or plain text)
// gets an inline diff; anything else shows Before/After values.
func (mh *ModelHistory) diffSection(recordKey string, aHash, bHash []byte, msgr *Messages) (h.HTMLComponent, error) {
	am, bm, err := mh.twoFieldMaps(recordKey, aHash, bHash)
	if err != nil {
		return nil, err
	}

	var rows h.HTMLComponents
	for _, f := range mh.resolved {
		ov, nv := fieldValue(am, f), fieldValue(bm, f)
		if ov == nv {
			rows = append(rows, h.Div(
				h.Span(f).Style("text-decoration: line-through"),
				h.Span(" "+msgr.Unchanged).Class("text-medium-emphasis text-caption ms-1"),
			).Class("mb-1"))
			continue
		}

		var val h.HTMLComponent
		if mh.IsHTML(f) || (isStr(am, f) && isStr(bm, f)) {
			val = h.RawHTML(HTMLDiff(ov, nv))
		} else {
			val = h.Div(
				h.Div(h.Strong(msgr.Old+": "), h.Text(ov)),
				h.Div(h.Strong(msgr.New+": "), h.Text(nv)),
			)
		}
		rows = append(rows, h.Div(
			h.Div(h.Strong(f)).Class("mb-1"),
			h.Div(val).Class("ps-2"),
		).Class("mb-3"))
	}
	return h.Div(rows...).Class("mt-2"), nil
}

// isStr reports whether a field's stored JSON value is a JSON string (so a text
// or HTML diff makes sense, as opposed to a structured/foreign-key value).
func isStr(m map[string]json.RawMessage, f string) bool {
	raw, ok := m[f]
	return ok && len(raw) > 0 && raw[0] == '"'
}
