package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDataTableVirtualBuilder struct {
	VTagBuilder[*VDataTableVirtualBuilder]
}

func VDataTableVirtual(children ...h.HTMLComponent) *VDataTableVirtualBuilder {
	return VTag(&VDataTableVirtualBuilder{}, "v-data-table-virtual", children...)
}

func (b *VDataTableVirtualBuilder) Width(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) HeaderProps(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":header-props", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) CellProps(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":cell-props", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) Mobile(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) Headers(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":headers", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) Loading(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) LoadingText(v string) (r *VDataTableVirtualBuilder) {
	b.Attr("loading-text", v)
	return b
}

func (b *VDataTableVirtualBuilder) HideNoData(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":hide-no-data", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) Items(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) NoDataText(v string) (r *VDataTableVirtualBuilder) {
	b.Attr("no-data-text", v)
	return b
}

func (b *VDataTableVirtualBuilder) MobileBreakpoint(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) RowProps(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":row-props", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) HideDefaultBody(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":hide-default-body", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) HideDefaultFooter(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":hide-default-footer", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) HideDefaultHeader(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":hide-default-header", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) Search(v string) (r *VDataTableVirtualBuilder) {
	b.Attr("search", v)
	return b
}

func (b *VDataTableVirtualBuilder) ExpandOnClick(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":expand-on-click", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) ShowExpand(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":show-expand", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) Expanded(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":expanded", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) GroupBy(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":group-by", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) ItemValue(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":item-value", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) ItemSelectable(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":item-selectable", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) ReturnObject(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) ShowSelect(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":show-select", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) SelectStrategy(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":select-strategy", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) ModelValue(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) ValueComparator(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) SortBy(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":sort-by", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) MultiSort(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":multi-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) MustSort(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":must-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) CustomKeySort(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":custom-key-sort", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) Color(v string) (r *VDataTableVirtualBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDataTableVirtualBuilder) Sticky(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) DisableSort(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":disable-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) SortAscIcon(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":sort-asc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) SortDescIcon(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":sort-desc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) FixedHeader(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":fixed-header", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) FixedFooter(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":fixed-footer", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) Height(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) Hover(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":hover", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) Density(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) Tag(v string) (r *VDataTableVirtualBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VDataTableVirtualBuilder) Theme(v string) (r *VDataTableVirtualBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VDataTableVirtualBuilder) ItemHeight(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":item-height", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) FilterMode(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":filter-mode", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) NoFilter(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":no-filter", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) CustomFilter(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":custom-filter", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) CustomKeyFilter(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":custom-key-filter", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) FilterKeys(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":filter-keys", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) On(name string, value string) (r *VDataTableVirtualBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableVirtualBuilder) Bind(name string, value string) (r *VDataTableVirtualBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
