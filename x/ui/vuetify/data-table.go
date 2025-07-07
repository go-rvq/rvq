package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDataTableBuilder struct {
	VTagBuilder[*VDataTableBuilder]
}

func VDataTable(children ...h.HTMLComponent) *VDataTableBuilder {
	return VTag(&VDataTableBuilder{}, "v-data-table", children...)
}

func (b *VDataTableBuilder) Width(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) HeaderProps(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":header-props", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) CellProps(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":cell-props", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) Mobile(v bool) (r *VDataTableBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) Loading(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) Headers(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":headers", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) HeadersVar(v string) (r *VDataTableBuilder) {
	b.Attr(":headers", v)
	return b
}

func (b *VDataTableBuilder) Page(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":page", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ItemsPerPage(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":items-per-page", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) LoadingText(v string) (r *VDataTableBuilder) {
	b.Attr("loading-text", v)
	return b
}

func (b *VDataTableBuilder) HideNoData(v bool) (r *VDataTableBuilder) {
	b.Attr(":hide-no-data", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) Items(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ItemsVar(v string) (r *VDataTableBuilder) {
	b.Attr(":items", v)
	return b
}

func (b *VDataTableBuilder) NoDataText(v string) (r *VDataTableBuilder) {
	b.Attr("no-data-text", v)
	return b
}

func (b *VDataTableBuilder) MobileBreakpoint(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) RowProps(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":row-props", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) HideDefaultBody(v bool) (r *VDataTableBuilder) {
	b.Attr(":hide-default-body", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) HideDefaultFooter(v bool) (r *VDataTableBuilder) {
	b.Attr(":hide-default-footer", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) HideDefaultHeader(v bool) (r *VDataTableBuilder) {
	b.Attr(":hide-default-header", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) Search(v string) (r *VDataTableBuilder) {
	b.Attr("search", v)
	return b
}

func (b *VDataTableBuilder) ExpandOnClick(v bool) (r *VDataTableBuilder) {
	b.Attr(":expand-on-click", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) ShowExpand(v bool) (r *VDataTableBuilder) {
	b.Attr(":show-expand", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) Expanded(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":expanded", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) GroupBy(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":group-by", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ItemValue(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":item-value", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ItemSelectable(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":item-selectable", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ReturnObject(v bool) (r *VDataTableBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) ShowSelect(v bool) (r *VDataTableBuilder) {
	b.Attr(":show-select", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) SelectStrategy(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":select-strategy", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ModelValue(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ValueComparator(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) SortBy(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":sort-by", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) MultiSort(v bool) (r *VDataTableBuilder) {
	b.Attr(":multi-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) MustSort(v bool) (r *VDataTableBuilder) {
	b.Attr(":must-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) CustomKeySort(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":custom-key-sort", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) Color(v string) (r *VDataTableBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDataTableBuilder) Sticky(v bool) (r *VDataTableBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) DisableSort(v bool) (r *VDataTableBuilder) {
	b.Attr(":disable-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) SortAscIcon(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":sort-asc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) SortDescIcon(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":sort-desc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) FixedHeader(v bool) (r *VDataTableBuilder) {
	b.Attr(":fixed-header", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) FixedFooter(v bool) (r *VDataTableBuilder) {
	b.Attr(":fixed-footer", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) Height(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) Hover(v bool) (r *VDataTableBuilder) {
	b.Attr(":hover", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) Density(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) Tag(v string) (r *VDataTableBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VDataTableBuilder) Theme(v string) (r *VDataTableBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VDataTableBuilder) FilterMode(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":filter-mode", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) NoFilter(v bool) (r *VDataTableBuilder) {
	b.Attr(":no-filter", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) CustomFilter(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":custom-filter", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) CustomKeyFilter(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":custom-key-filter", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) FilterKeys(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":filter-keys", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) PrevIcon(v string) (r *VDataTableBuilder) {
	b.Attr("prev-icon", v)
	return b
}

func (b *VDataTableBuilder) NextIcon(v string) (r *VDataTableBuilder) {
	b.Attr("next-icon", v)
	return b
}

func (b *VDataTableBuilder) FirstIcon(v string) (r *VDataTableBuilder) {
	b.Attr("first-icon", v)
	return b
}

func (b *VDataTableBuilder) LastIcon(v string) (r *VDataTableBuilder) {
	b.Attr("last-icon", v)
	return b
}

func (b *VDataTableBuilder) ItemsPerPageText(v string) (r *VDataTableBuilder) {
	b.Attr("items-per-page-text", v)
	return b
}

func (b *VDataTableBuilder) PageText(v string) (r *VDataTableBuilder) {
	b.Attr("page-text", v)
	return b
}

func (b *VDataTableBuilder) FirstPageLabel(v string) (r *VDataTableBuilder) {
	b.Attr("first-page-label", v)
	return b
}

func (b *VDataTableBuilder) PrevPageLabel(v string) (r *VDataTableBuilder) {
	b.Attr("prev-page-label", v)
	return b
}

func (b *VDataTableBuilder) NextPageLabel(v string) (r *VDataTableBuilder) {
	b.Attr("next-page-label", v)
	return b
}

func (b *VDataTableBuilder) LastPageLabel(v string) (r *VDataTableBuilder) {
	b.Attr("last-page-label", v)
	return b
}

func (b *VDataTableBuilder) ItemsPerPageOptions(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":items-per-page-options", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ShowCurrentPage(v bool) (r *VDataTableBuilder) {
	b.Attr(":show-current-page", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) On(name string, value string) (r *VDataTableBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableBuilder) Bind(name string, value string) (r *VDataTableBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
