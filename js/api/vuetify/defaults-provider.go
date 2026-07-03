package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDefaultsProviderBuilder struct {
	VTagBuilder[*VDefaultsProviderBuilder]
}

func VDefaultsProvider(children ...h.HTMLComponent) *VDefaultsProviderBuilder {
	return VTag(&VDefaultsProviderBuilder{}, "v-defaults-provider", children...)
}

func (b *VDefaultsProviderBuilder) Reset(v interface{}) (r *VDefaultsProviderBuilder) {
	b.Attr(":reset", h.JSONString(v))
	return b
}

func (b *VDefaultsProviderBuilder) Disabled(v bool) (r *VDefaultsProviderBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VDefaultsProviderBuilder) Root(v interface{}) (r *VDefaultsProviderBuilder) {
	b.Attr(":root", h.JSONString(v))
	return b
}

func (b *VDefaultsProviderBuilder) Scoped(v bool) (r *VDefaultsProviderBuilder) {
	b.Attr(":scoped", fmt.Sprint(v))
	return b
}

func (b *VDefaultsProviderBuilder) Defaults(v interface{}) (r *VDefaultsProviderBuilder) {
	b.Attr(":defaults", h.JSONString(v))
	return b
}

func (b *VDefaultsProviderBuilder) On(name string, value string) (r *VDefaultsProviderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDefaultsProviderBuilder) Bind(name string, value string) (r *VDefaultsProviderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDefaultsProviderBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDefaultsProviderBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDefaultsProviderBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDefaultsProviderBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDefaultsProviderBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDefaultsProviderBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDefaultsProviderBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VDefaultsProviderBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VDefaultsProviderBuilder) SlotDefault(child ...h.HTMLComponent) (r *VDefaultsProviderBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VDefaultsProviderBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VDefaultsProviderBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
