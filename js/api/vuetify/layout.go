package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VLayoutBuilder struct {
	VTagBuilder[*VLayoutBuilder]
}

func VLayout(children ...h.HTMLComponent) *VLayoutBuilder {
	return VTag(&VLayoutBuilder{}, "v-layout", children...)
}

func (b *VLayoutBuilder) FullHeight(v bool) (r *VLayoutBuilder) {
	b.Attr(":full-height", fmt.Sprint(v))
	return b
}

func (b *VLayoutBuilder) Overlaps(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":overlaps", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) Height(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) MaxHeight(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) MaxWidth(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) MinHeight(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) MinWidth(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) Width(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) On(name string, value string) (r *VLayoutBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VLayoutBuilder) Bind(name string, value string) (r *VLayoutBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VLayoutBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VLayoutBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VLayoutBuilder) Slot(name string, child ...h.HTMLComponent) (r *VLayoutBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VLayoutBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VLayoutBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VLayoutBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VLayoutBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VLayoutBuilder) SlotDefault(child ...h.HTMLComponent) (r *VLayoutBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VLayoutBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VLayoutBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
