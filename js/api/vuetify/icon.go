package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VIconBuilder struct {
	VTagBuilder[*VIconBuilder]
}

func VIcon(children ...h.HTMLComponent) *VIconBuilder {
	return VTag(&VIconBuilder{}, "v-icon", children...)
}

func (b *VIconBuilder) Tag(v interface{}) (r *VIconBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VIconBuilder) Theme(v string) (r *VIconBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VIconBuilder) Size(v interface{}) (r *VIconBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VIconBuilder) Disabled(v bool) (r *VIconBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VIconBuilder) Color(v string) (r *VIconBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VIconBuilder) Opacity(v interface{}) (r *VIconBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VIconBuilder) Start(v bool) (r *VIconBuilder) {
	b.Attr(":start", fmt.Sprint(v))
	return b
}

func (b *VIconBuilder) End(v bool) (r *VIconBuilder) {
	b.Attr(":end", fmt.Sprint(v))
	return b
}

func (b *VIconBuilder) Icon(v interface{}) (r *VIconBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VIconBuilder) On(name string, value string) (r *VIconBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VIconBuilder) Bind(name string, value string) (r *VIconBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VIconBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VIconBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VIconBuilder) Slot(name string, child ...h.HTMLComponent) (r *VIconBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VIconBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VIconBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VIconBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VIconBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VIconBuilder) SlotDefault(child ...h.HTMLComponent) (r *VIconBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VIconBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VIconBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
