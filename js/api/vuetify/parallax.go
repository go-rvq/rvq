package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VParallaxBuilder struct {
	VTagBuilder[*VParallaxBuilder]
}

func VParallax(children ...h.HTMLComponent) *VParallaxBuilder {
	return VTag(&VParallaxBuilder{}, "v-parallax", children...)
}

func (b *VParallaxBuilder) Scale(v interface{}) (r *VParallaxBuilder) {
	b.Attr(":scale", h.JSONString(v))
	return b
}

func (b *VParallaxBuilder) On(name string, value string) (r *VParallaxBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VParallaxBuilder) Bind(name string, value string) (r *VParallaxBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VParallaxBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VParallaxBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VParallaxBuilder) Slot(name string, child ...h.HTMLComponent) (r *VParallaxBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VParallaxBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VParallaxBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VParallaxBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VParallaxBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VParallaxBuilder) SlotDefault(child ...h.HTMLComponent) (r *VParallaxBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VParallaxBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VParallaxBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VParallaxBuilder) SetSlotPlaceholder(child ...h.HTMLComponent) {
	b.SetSlot("placeholder", child...)
}

func (b *VParallaxBuilder) SetScopedSlotPlaceholder(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("placeholder", scope, child...)
}

func (b *VParallaxBuilder) SlotPlaceholder(child ...h.HTMLComponent) (r *VParallaxBuilder) {
	b.SetSlotPlaceholder(child...)
	return b
}

func (b *VParallaxBuilder) ScopedSlotPlaceholder(scope string, child ...h.HTMLComponent) (r *VParallaxBuilder) {
	b.SetScopedSlotPlaceholder(scope, child...)
	return b
}

func (b *VParallaxBuilder) SetSlotError(child ...h.HTMLComponent) {
	b.SetSlot("error", child...)
}

func (b *VParallaxBuilder) SetScopedSlotError(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("error", scope, child...)
}

func (b *VParallaxBuilder) SlotError(child ...h.HTMLComponent) (r *VParallaxBuilder) {
	b.SetSlotError(child...)
	return b
}

func (b *VParallaxBuilder) ScopedSlotError(scope string, child ...h.HTMLComponent) (r *VParallaxBuilder) {
	b.SetScopedSlotError(scope, child...)
	return b
}

func (b *VParallaxBuilder) SetSlotSources(child ...h.HTMLComponent) {
	b.SetSlot("sources", child...)
}

func (b *VParallaxBuilder) SetScopedSlotSources(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("sources", scope, child...)
}

func (b *VParallaxBuilder) SlotSources(child ...h.HTMLComponent) (r *VParallaxBuilder) {
	b.SetSlotSources(child...)
	return b
}

func (b *VParallaxBuilder) ScopedSlotSources(scope string, child ...h.HTMLComponent) (r *VParallaxBuilder) {
	b.SetScopedSlotSources(scope, child...)
	return b
}
