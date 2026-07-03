package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBreadcrumbsItemBuilder struct {
	VTagBuilder[*VBreadcrumbsItemBuilder]
}

func VBreadcrumbsItem(children ...h.HTMLComponent) *VBreadcrumbsItemBuilder {
	return VTag(&VBreadcrumbsItemBuilder{}, "v-breadcrumbs-item", children...)
}

func (b *VBreadcrumbsItemBuilder) Title(v string) (r *VBreadcrumbsItemBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VBreadcrumbsItemBuilder) Replace(v bool) (r *VBreadcrumbsItemBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VBreadcrumbsItemBuilder) Tag(v interface{}) (r *VBreadcrumbsItemBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsItemBuilder) Color(v string) (r *VBreadcrumbsItemBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBreadcrumbsItemBuilder) Active(v bool) (r *VBreadcrumbsItemBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VBreadcrumbsItemBuilder) ActiveColor(v string) (r *VBreadcrumbsItemBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VBreadcrumbsItemBuilder) Disabled(v bool) (r *VBreadcrumbsItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VBreadcrumbsItemBuilder) Href(v string) (r *VBreadcrumbsItemBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VBreadcrumbsItemBuilder) Exact(v bool) (r *VBreadcrumbsItemBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VBreadcrumbsItemBuilder) To(v interface{}) (r *VBreadcrumbsItemBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsItemBuilder) ActiveClass(v string) (r *VBreadcrumbsItemBuilder) {
	b.Attr("active-class", v)
	return b
}

func (b *VBreadcrumbsItemBuilder) On(name string, value string) (r *VBreadcrumbsItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBreadcrumbsItemBuilder) Bind(name string, value string) (r *VBreadcrumbsItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VBreadcrumbsItemBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VBreadcrumbsItemBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VBreadcrumbsItemBuilder) Slot(name string, child ...h.HTMLComponent) (r *VBreadcrumbsItemBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VBreadcrumbsItemBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VBreadcrumbsItemBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VBreadcrumbsItemBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VBreadcrumbsItemBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VBreadcrumbsItemBuilder) SlotDefault(child ...h.HTMLComponent) (r *VBreadcrumbsItemBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VBreadcrumbsItemBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VBreadcrumbsItemBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
