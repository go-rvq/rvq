package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VThemeProviderBuilder struct {
	VTagBuilder[*VThemeProviderBuilder]
}

func VThemeProvider(children ...h.HTMLComponent) *VThemeProviderBuilder {
	return VTag(&VThemeProviderBuilder{}, "v-theme-provider", children...)
}

func (b *VThemeProviderBuilder) Tag(v interface{}) (r *VThemeProviderBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VThemeProviderBuilder) Theme(v string) (r *VThemeProviderBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VThemeProviderBuilder) WithBackground(v bool) (r *VThemeProviderBuilder) {
	b.Attr(":with-background", fmt.Sprint(v))
	return b
}

func (b *VThemeProviderBuilder) On(name string, value string) (r *VThemeProviderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VThemeProviderBuilder) Bind(name string, value string) (r *VThemeProviderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VThemeProviderBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VThemeProviderBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VThemeProviderBuilder) Slot(name string, child ...h.HTMLComponent) (r *VThemeProviderBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VThemeProviderBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VThemeProviderBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VThemeProviderBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VThemeProviderBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VThemeProviderBuilder) SlotDefault(child ...h.HTMLComponent) (r *VThemeProviderBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VThemeProviderBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VThemeProviderBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
