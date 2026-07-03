package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VAppBuilder struct {
	VTagBuilder[*VAppBuilder]
}

func VApp(children ...h.HTMLComponent) *VAppBuilder {
	return VTag(&VAppBuilder{}, "v-app", children...)
}

func (b *VAppBuilder) FullHeight(v bool) (r *VAppBuilder) {
	b.Attr(":full-height", fmt.Sprint(v))
	return b
}

func (b *VAppBuilder) Overlaps(v interface{}) (r *VAppBuilder) {
	b.Attr(":overlaps", h.JSONString(v))
	return b
}

func (b *VAppBuilder) Theme(v string) (r *VAppBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VAppBuilder) On(name string, value string) (r *VAppBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAppBuilder) Bind(name string, value string) (r *VAppBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VAppBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VAppBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VAppBuilder) Slot(name string, child ...h.HTMLComponent) (r *VAppBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VAppBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VAppBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VAppBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VAppBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VAppBuilder) SlotDefault(child ...h.HTMLComponent) (r *VAppBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VAppBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VAppBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
