package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VComponentIconBuilder struct {
	VTagBuilder[*VComponentIconBuilder]
}

func VComponentIcon(children ...h.HTMLComponent) *VComponentIconBuilder {
	return VTag(&VComponentIconBuilder{}, "v-component-icon", children...)
}

func (b *VComponentIconBuilder) Tag(v interface{}) (r *VComponentIconBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VComponentIconBuilder) Icon(v interface{}) (r *VComponentIconBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VComponentIconBuilder) On(name string, value string) (r *VComponentIconBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VComponentIconBuilder) Bind(name string, value string) (r *VComponentIconBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VComponentIconBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VComponentIconBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VComponentIconBuilder) Slot(name string, child ...h.HTMLComponent) (r *VComponentIconBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VComponentIconBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VComponentIconBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VComponentIconBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VComponentIconBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VComponentIconBuilder) SlotDefault(child ...h.HTMLComponent) (r *VComponentIconBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VComponentIconBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VComponentIconBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
