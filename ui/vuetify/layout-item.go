package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VLayoutItemBuilder struct {
	VTagBuilder[*VLayoutItemBuilder]
}

func VLayoutItem(children ...h.HTMLComponent) *VLayoutItemBuilder {
	return VTag(&VLayoutItemBuilder{}, "v-layout-item", children...)
}

func (b *VLayoutItemBuilder) Position(v interface{}) (r *VLayoutItemBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VLayoutItemBuilder) Size(v interface{}) (r *VLayoutItemBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VLayoutItemBuilder) ModelValue(v bool) (r *VLayoutItemBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VLayoutItemBuilder) Name(v string) (r *VLayoutItemBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VLayoutItemBuilder) Order(v interface{}) (r *VLayoutItemBuilder) {
	b.Attr(":order", h.JSONString(v))
	return b
}

func (b *VLayoutItemBuilder) Absolute(v bool) (r *VLayoutItemBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VLayoutItemBuilder) On(name string, value string) (r *VLayoutItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VLayoutItemBuilder) Bind(name string, value string) (r *VLayoutItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
