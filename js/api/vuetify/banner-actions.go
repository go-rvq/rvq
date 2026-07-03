package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBannerActionsBuilder struct {
	VTagBuilder[*VBannerActionsBuilder]
}

func VBannerActions(children ...h.HTMLComponent) *VBannerActionsBuilder {
	return VTag(&VBannerActionsBuilder{}, "v-banner-actions", children...)
}

func (b *VBannerActionsBuilder) Density(v string) (r *VBannerActionsBuilder) {
	b.Attr("density", v)
	return b
}

func (b *VBannerActionsBuilder) Color(v string) (r *VBannerActionsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBannerActionsBuilder) On(name string, value string) (r *VBannerActionsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBannerActionsBuilder) Bind(name string, value string) (r *VBannerActionsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VBannerActionsBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VBannerActionsBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VBannerActionsBuilder) Slot(name string, child ...h.HTMLComponent) (r *VBannerActionsBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VBannerActionsBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VBannerActionsBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VBannerActionsBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VBannerActionsBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VBannerActionsBuilder) SlotDefault(child ...h.HTMLComponent) (r *VBannerActionsBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VBannerActionsBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VBannerActionsBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
