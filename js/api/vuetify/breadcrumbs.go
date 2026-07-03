package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBreadcrumbsBuilder struct {
	VTagBuilder[*VBreadcrumbsBuilder]
}

func VBreadcrumbs(children ...h.HTMLComponent) *VBreadcrumbsBuilder {
	return VTag(&VBreadcrumbsBuilder{}, "v-breadcrumbs", children...)
}

func (b *VBreadcrumbsBuilder) Icon(v interface{}) (r *VBreadcrumbsBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsBuilder) Density(v interface{}) (r *VBreadcrumbsBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsBuilder) Rounded(v interface{}) (r *VBreadcrumbsBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsBuilder) Tile(v bool) (r *VBreadcrumbsBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VBreadcrumbsBuilder) Tag(v interface{}) (r *VBreadcrumbsBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsBuilder) Color(v string) (r *VBreadcrumbsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBreadcrumbsBuilder) ActiveColor(v string) (r *VBreadcrumbsBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VBreadcrumbsBuilder) Disabled(v bool) (r *VBreadcrumbsBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VBreadcrumbsBuilder) BgColor(v string) (r *VBreadcrumbsBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VBreadcrumbsBuilder) Divider(v string) (r *VBreadcrumbsBuilder) {
	b.Attr("divider", v)
	return b
}

func (b *VBreadcrumbsBuilder) ActiveClass(v string) (r *VBreadcrumbsBuilder) {
	b.Attr("active-class", v)
	return b
}

func (b *VBreadcrumbsBuilder) Items(v interface{}) (r *VBreadcrumbsBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsBuilder) On(name string, value string) (r *VBreadcrumbsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBreadcrumbsBuilder) Bind(name string, value string) (r *VBreadcrumbsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VBreadcrumbsBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VBreadcrumbsBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VBreadcrumbsBuilder) Slot(name string, child ...h.HTMLComponent) (r *VBreadcrumbsBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VBreadcrumbsBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VBreadcrumbsBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VBreadcrumbsBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VBreadcrumbsBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VBreadcrumbsBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VBreadcrumbsBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VBreadcrumbsBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VBreadcrumbsBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VBreadcrumbsBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VBreadcrumbsBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VBreadcrumbsBuilder) SlotTitle(child ...h.HTMLComponent) (r *VBreadcrumbsBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VBreadcrumbsBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VBreadcrumbsBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VBreadcrumbsBuilder) SetSlotDivider(child ...h.HTMLComponent) {
	b.SetSlot("divider", child...)
}

func (b *VBreadcrumbsBuilder) SetScopedSlotDivider(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("divider", scope, child...)
}

func (b *VBreadcrumbsBuilder) SlotDivider(child ...h.HTMLComponent) (r *VBreadcrumbsBuilder) {
	b.SetSlotDivider(child...)
	return b
}

func (b *VBreadcrumbsBuilder) ScopedSlotDivider(scope string, child ...h.HTMLComponent) (r *VBreadcrumbsBuilder) {
	b.SetScopedSlotDivider(scope, child...)
	return b
}

func (b *VBreadcrumbsBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VBreadcrumbsBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VBreadcrumbsBuilder) SlotItem(child ...h.HTMLComponent) (r *VBreadcrumbsBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VBreadcrumbsBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VBreadcrumbsBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VBreadcrumbsBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VBreadcrumbsBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VBreadcrumbsBuilder) SlotDefault(child ...h.HTMLComponent) (r *VBreadcrumbsBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VBreadcrumbsBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VBreadcrumbsBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
