package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VMainBuilder struct {
	VTagBuilder[*VMainBuilder]
}

func VMain(children ...h.HTMLComponent) *VMainBuilder {
	return VTag(&VMainBuilder{}, "v-main", children...)
}

func (b *VMainBuilder) Tag(v interface{}) (r *VMainBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VMainBuilder) Height(v interface{}) (r *VMainBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VMainBuilder) MaxHeight(v interface{}) (r *VMainBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VMainBuilder) MaxWidth(v interface{}) (r *VMainBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VMainBuilder) MinHeight(v interface{}) (r *VMainBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VMainBuilder) MinWidth(v interface{}) (r *VMainBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VMainBuilder) Width(v interface{}) (r *VMainBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VMainBuilder) Scrollable(v bool) (r *VMainBuilder) {
	b.Attr(":scrollable", fmt.Sprint(v))
	return b
}

func (b *VMainBuilder) On(name string, value string) (r *VMainBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VMainBuilder) Bind(name string, value string) (r *VMainBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VMainBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VMainBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VMainBuilder) Slot(name string, child ...h.HTMLComponent) (r *VMainBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VMainBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VMainBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VMainBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VMainBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VMainBuilder) SlotDefault(child ...h.HTMLComponent) (r *VMainBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VMainBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VMainBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
