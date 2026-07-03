package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDividerBuilder struct {
	VTagBuilder[*VDividerBuilder]
}

func VDivider(children ...h.HTMLComponent) *VDividerBuilder {
	return VTag(&VDividerBuilder{}, "v-divider", children...)
}

func (b *VDividerBuilder) Length(v interface{}) (r *VDividerBuilder) {
	b.Attr(":length", h.JSONString(v))
	return b
}

func (b *VDividerBuilder) Theme(v string) (r *VDividerBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VDividerBuilder) Color(v string) (r *VDividerBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDividerBuilder) Opacity(v interface{}) (r *VDividerBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VDividerBuilder) Vertical(v bool) (r *VDividerBuilder) {
	b.Attr(":vertical", fmt.Sprint(v))
	return b
}

func (b *VDividerBuilder) Inset(v bool) (r *VDividerBuilder) {
	b.Attr(":inset", fmt.Sprint(v))
	return b
}

func (b *VDividerBuilder) Thickness(v interface{}) (r *VDividerBuilder) {
	b.Attr(":thickness", h.JSONString(v))
	return b
}

func (b *VDividerBuilder) On(name string, value string) (r *VDividerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDividerBuilder) Bind(name string, value string) (r *VDividerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDividerBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDividerBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDividerBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDividerBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDividerBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDividerBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDividerBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VDividerBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VDividerBuilder) SlotDefault(child ...h.HTMLComponent) (r *VDividerBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VDividerBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VDividerBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
