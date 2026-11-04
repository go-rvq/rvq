package presets

import (
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// The selection of several records: the listing opened to choose
// (OpenListingDialogForSelection) with ParamSelectMany=1 shows a column to
// mark the records — the marks kept across the pages and the filters
// (ParamSelectedIds) —, a button that shows only the marked ones (a table
// like the other, where they are unmarked; ParamSelectedView=1) and one that
// gives them to who opened it: its event (SelectedEventParamName), with the
// ids (ParamSelectedIds) and its config (SelectedEventConfigParamName). The
// records already chosen (ParamSelectedDisabledIds) are marked and fixed.
// The listing of the selection has no aggregates; SelectionConditionFunc
// restricts it and SelectionHiddenFilters hides the filters it fixes.
const (
	ParamSelectMany          = "selectMany"
	ParamSelectedView        = "selectedView"
	ParamSelectedDisabledIds = "selectedDisabledIds"
)

// IsListingSelection says the request is of the listing opened to choose
// (it has the event to give the choice to).
func IsListingSelection(ctx *web.EventContext) bool {
	return ctx.R.URL.Query().Get(SelectedEventParamName) != ""
}

// IsListingSelectMany says the request is of the listing opened to choose
// several records.
func IsListingSelectMany(ctx *web.EventContext) bool {
	return IsListingSelection(ctx) && ctx.R.URL.Query().Get(ParamSelectMany) == "1"
}

// isSelectedView says the listing of the selection shows only the marked.
func isSelectedView(ctx *web.EventContext) bool {
	return IsListingSelectMany(ctx) && ctx.R.URL.Query().Get(ParamSelectedView) == "1"
}

func queryIDs(ctx *web.EventContext, name string) (ids []string) {
	for _, id := range strings.Split(ctx.R.URL.Query().Get(name), ",") {
		if id = strings.TrimSpace(id); id != "" {
			ids = append(ids, id)
		}
	}
	return
}

// SelectionConditionFunc restricts the listing opened to choose: its
// conditions, only there (the context of who opened it in
// SelectedEventConfigParamName, say).
func (b *ListingBuilder) SelectionConditionFunc(f func(ctx *web.EventContext) []*SQLCondition) *ListingBuilder {
	b.selectionConditionFunc = f
	return b
}

// SelectionHiddenFilters hides the filter items keys (as the FilterDataFunc
// declares them) in the listing opened to choose: the ones its condition
// fixes.
func (b *ListingBuilder) SelectionHiddenFilters(keys ...string) *ListingBuilder {
	b.selectionHiddenFilters = append(b.selectionHiddenFilters, keys...)
	return b
}

// selectionConditions are the conditions of the listing opened to choose: its
// own (SelectionConditionFunc) and, showing only the marked, the ids of them.
func (b *ListingBuilder) selectionConditions(ctx *web.EventContext) (conds []*SQLCondition) {
	if !IsListingSelection(ctx) {
		return nil
	}
	if b.selectionConditionFunc != nil {
		conds = append(conds, b.selectionConditionFunc(ctx)...)
	}
	if isSelectedView(ctx) {
		ids := queryIDs(ctx, ParamSelectedIds)
		pks := b.mb.Schema().PrimaryFields().QuotedFullDBNames()
		switch {
		case len(ids) == 0 || len(pks) != 1:
			conds = append(conds, &SQLCondition{Query: "1 = 0"})
		default:
			conds = append(conds, &SQLCondition{Query: pks[0] + " IN ?", Args: []any{ids}})
		}
	}
	return
}

// selectionFilterData is fd without the filters the listing opened to choose
// hides (SelectionHiddenFilters).
func (b *ListingBuilder) selectionFilterData(ctx *web.EventContext, fd vx.FilterData) vx.FilterData {
	if len(b.selectionHiddenFilters) == 0 || !IsListingSelection(ctx) {
		return fd
	}
	var r vx.FilterData
	for _, it := range fd {
		hidden := false
		for _, k := range b.selectionHiddenFilters {
			if it.Key == k || it.Key == FilterKeyPrefix+k {
				hidden = true
				break
			}
		}
		if !hidden {
			r = append(r, it)
		}
	}
	return r
}

// selectionActions are the buttons of the listing opened to choose several:
// the marked ones only (or all again) and the one that gives them to who
// opened it.
func (lcb *ListingComponentBuilder) selectionActions(ctx *web.EventContext, msgr *Messages) h.HTMLComponent {
	var (
		qs        = ctx.R.URL.Query()
		selected  = queryIDs(ctx, ParamSelectedIds)
		portalID  = ctx.R.FormValue(ParamPortalID)
		target    = ctx.R.FormValue(ParamTargetPortal)
		toggleTo  = "1"
		toggleTxt = strings.ReplaceAll(msgr.ListingSelectionSelected, "{count}", strconv.Itoa(len(selected)))
		toggleIco = "mdi-checkbox-multiple-marked-outline"
	)
	if isSelectedView(ctx) {
		toggleTo, toggleTxt, toggleIco = "", msgr.ListingSelectionAll, "mdi-format-list-bulleted"
	}
	toggle := PageBarButton(toggleTxt, toggleIco).
		Attr("@click", web.Plaid().
			URL(ctx.R.RequestURI).
			EventFunc(actions.UpdateListingDialog).
			ValidQuery(ParamTargetPortal, target).
			Query(ParamPortalID, portalID).
			Query(ParamSelectedView, toggleTo).
			Query("page", "1").
			MergeQuery(true).
			Go()).
		Class("rvq-listing-selection__toggle")

	add := VBtn(msgr.ListingSelectionAdd).
		Color("primary").Variant(VariantFlat).PrependIcon("mdi-plus").
		Disabled(len(selected) == 0).
		Attr("@click", web.Plaid().
			EventFunc(qs.Get(SelectedEventParamName)).
			Query(ParamSelectedIds, strings.Join(selected, ",")).
			Query(SelectedEventConfigParamName, qs.Get(SelectedEventConfigParamName)).
			Go()+`;closer.show = false`).
		Class("ml-2 rvq-listing-selection__add")
	return h.Components(toggle, add)
}
