package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBreadcrumbsDividerBuilder struct {
	VTagBuilder[*VBreadcrumbsDividerBuilder]
}

func VBreadcrumbsDivider(children ...h.HTMLComponent) *VBreadcrumbsDividerBuilder {
	return VTag(&VBreadcrumbsDividerBuilder{}, "v-breadcrumbs-divider", children...)
}

func (b *VBreadcrumbsDividerBuilder) Divider(v interface{}) (r *VBreadcrumbsDividerBuilder) {
	b.Attr(":divider", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsDividerBuilder) On(name string, value string) (r *VBreadcrumbsDividerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBreadcrumbsDividerBuilder) Bind(name string, value string) (r *VBreadcrumbsDividerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VBreadcrumbsDividerBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VBreadcrumbsDividerBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VBreadcrumbsDividerBuilder) Slot(name string, child ...h.HTMLComponent) (r *VBreadcrumbsDividerBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VBreadcrumbsDividerBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VBreadcrumbsDividerBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VBreadcrumbsDividerBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VBreadcrumbsDividerBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VBreadcrumbsDividerBuilder) SlotDefault(child ...h.HTMLComponent) (r *VBreadcrumbsDividerBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VBreadcrumbsDividerBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VBreadcrumbsDividerBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
