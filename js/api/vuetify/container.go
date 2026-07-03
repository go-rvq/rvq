package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VContainerBuilder struct {
	VTagBuilder[*VContainerBuilder]
}

func VContainer(children ...h.HTMLComponent) *VContainerBuilder {
	return VTag(&VContainerBuilder{}, "v-container", children...)
}

func (b *VContainerBuilder) Tag(v interface{}) (r *VContainerBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VContainerBuilder) Height(v interface{}) (r *VContainerBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VContainerBuilder) MaxHeight(v interface{}) (r *VContainerBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VContainerBuilder) MaxWidth(v interface{}) (r *VContainerBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VContainerBuilder) MinHeight(v interface{}) (r *VContainerBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VContainerBuilder) MinWidth(v interface{}) (r *VContainerBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VContainerBuilder) Width(v interface{}) (r *VContainerBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VContainerBuilder) Fluid(v bool) (r *VContainerBuilder) {
	b.Attr(":fluid", fmt.Sprint(v))
	return b
}

func (b *VContainerBuilder) On(name string, value string) (r *VContainerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VContainerBuilder) Bind(name string, value string) (r *VContainerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VContainerBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VContainerBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VContainerBuilder) Slot(name string, child ...h.HTMLComponent) (r *VContainerBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VContainerBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VContainerBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VContainerBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VContainerBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VContainerBuilder) SlotDefault(child ...h.HTMLComponent) (r *VContainerBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VContainerBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VContainerBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
