package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDataTableRowBuilder struct {
	VTagBuilder[*VDataTableRowBuilder]
}

func VDataTableRow(children ...h.HTMLComponent) *VDataTableRowBuilder {
	return VTag(&VDataTableRowBuilder{}, "v-data-table-row", children...)
}

func (b *VDataTableRowBuilder) CellProps(v interface{}) (r *VDataTableRowBuilder) {
	b.Attr(":cell-props", h.JSONString(v))
	return b
}

func (b *VDataTableRowBuilder) Mobile(v bool) (r *VDataTableRowBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VDataTableRowBuilder) Index(v int) (r *VDataTableRowBuilder) {
	b.Attr(":index", fmt.Sprint(v))
	return b
}

func (b *VDataTableRowBuilder) MobileBreakpoint(v interface{}) (r *VDataTableRowBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VDataTableRowBuilder) Item(v interface{}) (r *VDataTableRowBuilder) {
	b.Attr(":item", h.JSONString(v))
	return b
}

func (b *VDataTableRowBuilder) On(name string, value string) (r *VDataTableRowBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableRowBuilder) Bind(name string, value string) (r *VDataTableRowBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
