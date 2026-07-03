package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VLocaleProviderBuilder struct {
	VTagBuilder[*VLocaleProviderBuilder]
}

func VLocaleProvider(children ...h.HTMLComponent) *VLocaleProviderBuilder {
	return VTag(&VLocaleProviderBuilder{}, "v-locale-provider", children...)
}

func (b *VLocaleProviderBuilder) Messages(v interface{}) (r *VLocaleProviderBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VLocaleProviderBuilder) Locale(v string) (r *VLocaleProviderBuilder) {
	b.Attr("locale", v)
	return b
}

func (b *VLocaleProviderBuilder) FallbackLocale(v string) (r *VLocaleProviderBuilder) {
	b.Attr("fallback-locale", v)
	return b
}

func (b *VLocaleProviderBuilder) Rtl(v bool) (r *VLocaleProviderBuilder) {
	b.Attr(":rtl", fmt.Sprint(v))
	return b
}

func (b *VLocaleProviderBuilder) On(name string, value string) (r *VLocaleProviderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VLocaleProviderBuilder) Bind(name string, value string) (r *VLocaleProviderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VLocaleProviderBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VLocaleProviderBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VLocaleProviderBuilder) Slot(name string, child ...h.HTMLComponent) (r *VLocaleProviderBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VLocaleProviderBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VLocaleProviderBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VLocaleProviderBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VLocaleProviderBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VLocaleProviderBuilder) SlotDefault(child ...h.HTMLComponent) (r *VLocaleProviderBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VLocaleProviderBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VLocaleProviderBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
