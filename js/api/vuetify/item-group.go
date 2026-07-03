package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VItemGroupBuilder struct {
	VTagBuilder[*VItemGroupBuilder]
}

func VItemGroup(children ...h.HTMLComponent) *VItemGroupBuilder {
	return VTag(&VItemGroupBuilder{}, "v-item-group", children...)
}

func (b *VItemGroupBuilder) Tag(v interface{}) (r *VItemGroupBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VItemGroupBuilder) Theme(v string) (r *VItemGroupBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VItemGroupBuilder) Disabled(v bool) (r *VItemGroupBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VItemGroupBuilder) Multiple(v bool) (r *VItemGroupBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VItemGroupBuilder) Mandatory(v interface{}) (r *VItemGroupBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VItemGroupBuilder) ModelValue(v interface{}) (r *VItemGroupBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VItemGroupBuilder) Max(v interface{}) (r *VItemGroupBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VItemGroupBuilder) SelectedClass(v string) (r *VItemGroupBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VItemGroupBuilder) On(name string, value string) (r *VItemGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VItemGroupBuilder) Bind(name string, value string) (r *VItemGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VItemGroupBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VItemGroupBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VItemGroupBuilder) Slot(name string, child ...h.HTMLComponent) (r *VItemGroupBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VItemGroupBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VItemGroupBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VItemGroupBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VItemGroupBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VItemGroupBuilder) SlotDefault(child ...h.HTMLComponent) (r *VItemGroupBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VItemGroupBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VItemGroupBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
