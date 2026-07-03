package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBtnGroupBuilder struct {
	VTagBuilder[*VBtnGroupBuilder]
}

func VBtnGroup(children ...h.HTMLComponent) *VBtnGroupBuilder {
	return VTag(&VBtnGroupBuilder{}, "v-btn-group", children...)
}

func (b *VBtnGroupBuilder) Border(v interface{}) (r *VBtnGroupBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VBtnGroupBuilder) Density(v interface{}) (r *VBtnGroupBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VBtnGroupBuilder) Elevation(v interface{}) (r *VBtnGroupBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VBtnGroupBuilder) Rounded(v interface{}) (r *VBtnGroupBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VBtnGroupBuilder) Tile(v bool) (r *VBtnGroupBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VBtnGroupBuilder) Tag(v interface{}) (r *VBtnGroupBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VBtnGroupBuilder) Theme(v string) (r *VBtnGroupBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VBtnGroupBuilder) Color(v string) (r *VBtnGroupBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBtnGroupBuilder) Variant(v interface{}) (r *VBtnGroupBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VBtnGroupBuilder) BaseColor(v string) (r *VBtnGroupBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VBtnGroupBuilder) Divided(v bool) (r *VBtnGroupBuilder) {
	b.Attr(":divided", fmt.Sprint(v))
	return b
}

func (b *VBtnGroupBuilder) On(name string, value string) (r *VBtnGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBtnGroupBuilder) Bind(name string, value string) (r *VBtnGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VBtnGroupBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VBtnGroupBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VBtnGroupBuilder) Slot(name string, child ...h.HTMLComponent) (r *VBtnGroupBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VBtnGroupBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VBtnGroupBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VBtnGroupBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VBtnGroupBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VBtnGroupBuilder) SlotDefault(child ...h.HTMLComponent) (r *VBtnGroupBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VBtnGroupBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VBtnGroupBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
