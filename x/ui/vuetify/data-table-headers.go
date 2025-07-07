package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VDataTableHeadersBuilder struct {
	VTagBuilder[*VDataTableHeadersBuilder]
}

func VDataTableHeaders(children ...h.HTMLComponent) *VDataTableHeadersBuilder {
	return VTag(&VDataTableHeadersBuilder{}, "v-data-table-headers", children...)
}

func (b *VDataTableHeadersBuilder) Color(v string) (r *VDataTableHeadersBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDataTableHeadersBuilder) Sticky(v bool) (r *VDataTableHeadersBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VDataTableHeadersBuilder) DisableSort(v bool) (r *VDataTableHeadersBuilder) {
	b.Attr(":disable-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableHeadersBuilder) MultiSort(v bool) (r *VDataTableHeadersBuilder) {
	b.Attr(":multi-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableHeadersBuilder) SortAscIcon(v interface{}) (r *VDataTableHeadersBuilder) {
	b.Attr(":sort-asc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableHeadersBuilder) SortDescIcon(v interface{}) (r *VDataTableHeadersBuilder) {
	b.Attr(":sort-desc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableHeadersBuilder) HeaderProps(v interface{}) (r *VDataTableHeadersBuilder) {
	b.Attr(":header-props", h.JSONString(v))
	return b
}

func (b *VDataTableHeadersBuilder) Mobile(v bool) (r *VDataTableHeadersBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VDataTableHeadersBuilder) MobileBreakpoint(v interface{}) (r *VDataTableHeadersBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VDataTableHeadersBuilder) Loading(v interface{}) (r *VDataTableHeadersBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VDataTableHeadersBuilder) On(name string, value string) (r *VDataTableHeadersBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableHeadersBuilder) Bind(name string, value string) (r *VDataTableHeadersBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
