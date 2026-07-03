package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBannerTextBuilder struct {
	VTagBuilder[*VBannerTextBuilder]
}

func VBannerText(children ...h.HTMLComponent) *VBannerTextBuilder {
	return VTag(&VBannerTextBuilder{}, "v-banner-text", children...)
}

func (b *VBannerTextBuilder) Tag(v string) (r *VBannerTextBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VBannerTextBuilder) On(name string, value string) (r *VBannerTextBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBannerTextBuilder) Bind(name string, value string) (r *VBannerTextBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VBannerTextBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VBannerTextBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VBannerTextBuilder) Slot(name string, child ...h.HTMLComponent) (r *VBannerTextBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VBannerTextBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VBannerTextBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VBannerTextBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VBannerTextBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VBannerTextBuilder) SlotDefault(child ...h.HTMLComponent) (r *VBannerTextBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VBannerTextBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VBannerTextBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
