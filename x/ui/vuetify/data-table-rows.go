package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VDataTableRowsBuilder struct {
	VTagBuilder[*VDataTableRowsBuilder]
}

func VDataTableRows(children ...h.HTMLComponent) *VDataTableRowsBuilder {
	return VTag(&VDataTableRowsBuilder{}, "v-data-table-rows", children...)
}

func (b *VDataTableRowsBuilder) CellProps(v interface{}) (r *VDataTableRowsBuilder) {
	b.Attr(":cell-props", h.JSONString(v))
	return b
}

func (b *VDataTableRowsBuilder) Mobile(v bool) (r *VDataTableRowsBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VDataTableRowsBuilder) Loading(v interface{}) (r *VDataTableRowsBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VDataTableRowsBuilder) LoadingText(v string) (r *VDataTableRowsBuilder) {
	b.Attr("loading-text", v)
	return b
}

func (b *VDataTableRowsBuilder) HideNoData(v bool) (r *VDataTableRowsBuilder) {
	b.Attr(":hide-no-data", fmt.Sprint(v))
	return b
}

func (b *VDataTableRowsBuilder) Items(v interface{}) (r *VDataTableRowsBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VDataTableRowsBuilder) NoDataText(v string) (r *VDataTableRowsBuilder) {
	b.Attr("no-data-text", v)
	return b
}

func (b *VDataTableRowsBuilder) MobileBreakpoint(v interface{}) (r *VDataTableRowsBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VDataTableRowsBuilder) RowProps(v interface{}) (r *VDataTableRowsBuilder) {
	b.Attr(":row-props", h.JSONString(v))
	return b
}

func (b *VDataTableRowsBuilder) On(name string, value string) (r *VDataTableRowsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableRowsBuilder) Bind(name string, value string) (r *VDataTableRowsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
