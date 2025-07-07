package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDataTableServerBuilder struct {
	VTagBuilder[*VDataTableServerBuilder]
}

func VDataTableServer(children ...h.HTMLComponent) *VDataTableServerBuilder {
	return VTag(&VDataTableServerBuilder{}, "v-data-table-server", children...)
}

func (b *VDataTableServerBuilder) Width(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) HeaderProps(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":header-props", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) CellProps(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":cell-props", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) Mobile(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) Loading(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) Headers(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":headers", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ItemsLength(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":items-length", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) Page(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":page", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ItemsPerPage(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":items-per-page", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) LoadingText(v string) (r *VDataTableServerBuilder) {
	b.Attr("loading-text", v)
	return b
}

func (b *VDataTableServerBuilder) HideNoData(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":hide-no-data", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) Items(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) NoDataText(v string) (r *VDataTableServerBuilder) {
	b.Attr("no-data-text", v)
	return b
}

func (b *VDataTableServerBuilder) MobileBreakpoint(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) RowProps(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":row-props", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) HideDefaultBody(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":hide-default-body", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) HideDefaultFooter(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":hide-default-footer", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) HideDefaultHeader(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":hide-default-header", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) Search(v string) (r *VDataTableServerBuilder) {
	b.Attr("search", v)
	return b
}

func (b *VDataTableServerBuilder) ExpandOnClick(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":expand-on-click", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) ShowExpand(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":show-expand", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) Expanded(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":expanded", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) GroupBy(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":group-by", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ItemValue(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":item-value", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ItemSelectable(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":item-selectable", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ReturnObject(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) ShowSelect(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":show-select", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) SelectStrategy(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":select-strategy", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ModelValue(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ValueComparator(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) SortBy(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":sort-by", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) MultiSort(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":multi-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) MustSort(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":must-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) CustomKeySort(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":custom-key-sort", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) Color(v string) (r *VDataTableServerBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDataTableServerBuilder) Sticky(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) DisableSort(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":disable-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) SortAscIcon(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":sort-asc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) SortDescIcon(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":sort-desc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) FixedHeader(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":fixed-header", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) FixedFooter(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":fixed-footer", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) Height(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) Hover(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":hover", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) Density(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) Tag(v string) (r *VDataTableServerBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VDataTableServerBuilder) Theme(v string) (r *VDataTableServerBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VDataTableServerBuilder) PrevIcon(v string) (r *VDataTableServerBuilder) {
	b.Attr("prev-icon", v)
	return b
}

func (b *VDataTableServerBuilder) NextIcon(v string) (r *VDataTableServerBuilder) {
	b.Attr("next-icon", v)
	return b
}

func (b *VDataTableServerBuilder) FirstIcon(v string) (r *VDataTableServerBuilder) {
	b.Attr("first-icon", v)
	return b
}

func (b *VDataTableServerBuilder) LastIcon(v string) (r *VDataTableServerBuilder) {
	b.Attr("last-icon", v)
	return b
}

func (b *VDataTableServerBuilder) ItemsPerPageText(v string) (r *VDataTableServerBuilder) {
	b.Attr("items-per-page-text", v)
	return b
}

func (b *VDataTableServerBuilder) PageText(v string) (r *VDataTableServerBuilder) {
	b.Attr("page-text", v)
	return b
}

func (b *VDataTableServerBuilder) FirstPageLabel(v string) (r *VDataTableServerBuilder) {
	b.Attr("first-page-label", v)
	return b
}

func (b *VDataTableServerBuilder) PrevPageLabel(v string) (r *VDataTableServerBuilder) {
	b.Attr("prev-page-label", v)
	return b
}

func (b *VDataTableServerBuilder) NextPageLabel(v string) (r *VDataTableServerBuilder) {
	b.Attr("next-page-label", v)
	return b
}

func (b *VDataTableServerBuilder) LastPageLabel(v string) (r *VDataTableServerBuilder) {
	b.Attr("last-page-label", v)
	return b
}

func (b *VDataTableServerBuilder) ItemsPerPageOptions(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":items-per-page-options", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ShowCurrentPage(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":show-current-page", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) On(name string, value string) (r *VDataTableServerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableServerBuilder) Bind(name string, value string) (r *VDataTableServerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
