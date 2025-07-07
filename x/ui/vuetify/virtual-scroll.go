package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VVirtualScrollBuilder struct {
	VTagBuilder[*VVirtualScrollBuilder]
}

func VVirtualScroll(children ...h.HTMLComponent) *VVirtualScrollBuilder {
	return VTag(&VVirtualScrollBuilder{}, "v-virtual-scroll", children...)
}

func (b *VVirtualScrollBuilder) Items(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) Renderless(v bool) (r *VVirtualScrollBuilder) {
	b.Attr(":renderless", fmt.Sprint(v))
	return b
}

func (b *VVirtualScrollBuilder) ItemHeight(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":item-height", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) Height(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) MaxHeight(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) MaxWidth(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) MinHeight(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) MinWidth(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) Width(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) On(name string, value string) (r *VVirtualScrollBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VVirtualScrollBuilder) Bind(name string, value string) (r *VVirtualScrollBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
