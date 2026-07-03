package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTreeviewGroupBuilder struct {
	VTagBuilder[*VTreeviewGroupBuilder]
}

func VTreeviewGroup(children ...h.HTMLComponent) *VTreeviewGroupBuilder {
	return VTag(&VTreeviewGroupBuilder{}, "v-treeview-group", children...)
}

func (b *VTreeviewGroupBuilder) Title(v string) (r *VTreeviewGroupBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VTreeviewGroupBuilder) Tag(v interface{}) (r *VTreeviewGroupBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VTreeviewGroupBuilder) Color(v string) (r *VTreeviewGroupBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTreeviewGroupBuilder) ActiveColor(v string) (r *VTreeviewGroupBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VTreeviewGroupBuilder) BaseColor(v string) (r *VTreeviewGroupBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VTreeviewGroupBuilder) PrependIcon(v interface{}) (r *VTreeviewGroupBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewGroupBuilder) AppendIcon(v interface{}) (r *VTreeviewGroupBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewGroupBuilder) Value(v interface{}) (r *VTreeviewGroupBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VTreeviewGroupBuilder) ExpandIcon(v interface{}) (r *VTreeviewGroupBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewGroupBuilder) CollapseIcon(v interface{}) (r *VTreeviewGroupBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewGroupBuilder) Fluid(v bool) (r *VTreeviewGroupBuilder) {
	b.Attr(":fluid", fmt.Sprint(v))
	return b
}

func (b *VTreeviewGroupBuilder) On(name string, value string) (r *VTreeviewGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTreeviewGroupBuilder) Bind(name string, value string) (r *VTreeviewGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VTreeviewGroupBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VTreeviewGroupBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VTreeviewGroupBuilder) Slot(name string, child ...h.HTMLComponent) (r *VTreeviewGroupBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VTreeviewGroupBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VTreeviewGroupBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VTreeviewGroupBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VTreeviewGroupBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VTreeviewGroupBuilder) SlotDefault(child ...h.HTMLComponent) (r *VTreeviewGroupBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VTreeviewGroupBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VTreeviewGroupBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VTreeviewGroupBuilder) SetSlotActivator(child ...h.HTMLComponent) {
	b.SetSlot("activator", child...)
}

func (b *VTreeviewGroupBuilder) SetScopedSlotActivator(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("activator", scope, child...)
}

func (b *VTreeviewGroupBuilder) SlotActivator(child ...h.HTMLComponent) (r *VTreeviewGroupBuilder) {
	b.SetSlotActivator(child...)
	return b
}

func (b *VTreeviewGroupBuilder) ScopedSlotActivator(scope string, child ...h.HTMLComponent) (r *VTreeviewGroupBuilder) {
	b.SetScopedSlotActivator(scope, child...)
	return b
}
