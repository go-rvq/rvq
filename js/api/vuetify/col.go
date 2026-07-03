package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VColBuilder struct {
	VTagBuilder[*VColBuilder]
}

func VCol(children ...h.HTMLComponent) *VColBuilder {
	return VTag(&VColBuilder{}, "v-col", children...)
}

func (b *VColBuilder) Tag(v interface{}) (r *VColBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VColBuilder) Order(v interface{}) (r *VColBuilder) {
	b.Attr(":order", h.JSONString(v))
	return b
}

func (b *VColBuilder) Sm(v interface{}) (r *VColBuilder) {
	b.Attr(":sm", h.JSONString(v))
	return b
}

func (b *VColBuilder) Md(v interface{}) (r *VColBuilder) {
	b.Attr(":md", h.JSONString(v))
	return b
}

func (b *VColBuilder) Lg(v interface{}) (r *VColBuilder) {
	b.Attr(":lg", h.JSONString(v))
	return b
}

func (b *VColBuilder) Xl(v interface{}) (r *VColBuilder) {
	b.Attr(":xl", h.JSONString(v))
	return b
}

func (b *VColBuilder) Xxl(v interface{}) (r *VColBuilder) {
	b.Attr(":xxl", h.JSONString(v))
	return b
}

func (b *VColBuilder) Offset(v interface{}) (r *VColBuilder) {
	b.Attr(":offset", h.JSONString(v))
	return b
}

func (b *VColBuilder) Cols(v interface{}) (r *VColBuilder) {
	b.Attr(":cols", h.JSONString(v))
	return b
}

func (b *VColBuilder) OffsetSm(v interface{}) (r *VColBuilder) {
	b.Attr(":offset-sm", h.JSONString(v))
	return b
}

func (b *VColBuilder) OffsetMd(v interface{}) (r *VColBuilder) {
	b.Attr(":offset-md", h.JSONString(v))
	return b
}

func (b *VColBuilder) OffsetLg(v interface{}) (r *VColBuilder) {
	b.Attr(":offset-lg", h.JSONString(v))
	return b
}

func (b *VColBuilder) OffsetXl(v interface{}) (r *VColBuilder) {
	b.Attr(":offset-xl", h.JSONString(v))
	return b
}

func (b *VColBuilder) OffsetXxl(v interface{}) (r *VColBuilder) {
	b.Attr(":offset-xxl", h.JSONString(v))
	return b
}

func (b *VColBuilder) OrderSm(v interface{}) (r *VColBuilder) {
	b.Attr(":order-sm", h.JSONString(v))
	return b
}

func (b *VColBuilder) OrderMd(v interface{}) (r *VColBuilder) {
	b.Attr(":order-md", h.JSONString(v))
	return b
}

func (b *VColBuilder) OrderLg(v interface{}) (r *VColBuilder) {
	b.Attr(":order-lg", h.JSONString(v))
	return b
}

func (b *VColBuilder) OrderXl(v interface{}) (r *VColBuilder) {
	b.Attr(":order-xl", h.JSONString(v))
	return b
}

func (b *VColBuilder) OrderXxl(v interface{}) (r *VColBuilder) {
	b.Attr(":order-xxl", h.JSONString(v))
	return b
}

func (b *VColBuilder) AlignSelf(v interface{}) (r *VColBuilder) {
	b.Attr(":align-self", h.JSONString(v))
	return b
}

func (b *VColBuilder) On(name string, value string) (r *VColBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VColBuilder) Bind(name string, value string) (r *VColBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VColBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VColBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VColBuilder) Slot(name string, child ...h.HTMLComponent) (r *VColBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VColBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VColBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VColBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VColBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VColBuilder) SlotDefault(child ...h.HTMLComponent) (r *VColBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VColBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VColBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
