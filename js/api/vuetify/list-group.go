package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VListGroupBuilder struct {
	VTagBuilder[*VListGroupBuilder]
}

func VListGroup(children ...h.HTMLComponent) *VListGroupBuilder {
	return VTag(&VListGroupBuilder{}, "v-list-group", children...)
}

func (b *VListGroupBuilder) Tag(v interface{}) (r *VListGroupBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VListGroupBuilder) Title(v string) (r *VListGroupBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VListGroupBuilder) Value(v interface{}) (r *VListGroupBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VListGroupBuilder) BaseColor(v string) (r *VListGroupBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VListGroupBuilder) ActiveColor(v string) (r *VListGroupBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VListGroupBuilder) ExpandIcon(v interface{}) (r *VListGroupBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VListGroupBuilder) CollapseIcon(v interface{}) (r *VListGroupBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VListGroupBuilder) Color(v string) (r *VListGroupBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VListGroupBuilder) AppendIcon(v interface{}) (r *VListGroupBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VListGroupBuilder) PrependIcon(v interface{}) (r *VListGroupBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VListGroupBuilder) Fluid(v bool) (r *VListGroupBuilder) {
	b.Attr(":fluid", fmt.Sprint(v))
	return b
}

func (b *VListGroupBuilder) Subgroup(v bool) (r *VListGroupBuilder) {
	b.Attr(":subgroup", fmt.Sprint(v))
	return b
}

func (b *VListGroupBuilder) On(name string, value string) (r *VListGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListGroupBuilder) Bind(name string, value string) (r *VListGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VListGroupBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VListGroupBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VListGroupBuilder) Slot(name string, child ...h.HTMLComponent) (r *VListGroupBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VListGroupBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VListGroupBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VListGroupBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VListGroupBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VListGroupBuilder) SlotDefault(child ...h.HTMLComponent) (r *VListGroupBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VListGroupBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VListGroupBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VListGroupBuilder) SetSlotActivator(child ...h.HTMLComponent) {
	b.SetSlot("activator", child...)
}

func (b *VListGroupBuilder) SetScopedSlotActivator(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("activator", scope, child...)
}

func (b *VListGroupBuilder) SlotActivator(child ...h.HTMLComponent) (r *VListGroupBuilder) {
	b.SetSlotActivator(child...)
	return b
}

func (b *VListGroupBuilder) ScopedSlotActivator(scope string, child ...h.HTMLComponent) (r *VListGroupBuilder) {
	b.SetScopedSlotActivator(scope, child...)
	return b
}
