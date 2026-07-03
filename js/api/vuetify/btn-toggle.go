package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBtnToggleBuilder struct {
	VTagBuilder[*VBtnToggleBuilder]
}

func VBtnToggle(children ...h.HTMLComponent) *VBtnToggleBuilder {
	return VTag(&VBtnToggleBuilder{}, "v-btn-toggle", children...)
}

func (b *VBtnToggleBuilder) Border(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) ModelValue(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) Density(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) Elevation(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) Rounded(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) Tile(v bool) (r *VBtnToggleBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VBtnToggleBuilder) Tag(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) Theme(v string) (r *VBtnToggleBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VBtnToggleBuilder) Color(v string) (r *VBtnToggleBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBtnToggleBuilder) Variant(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) BaseColor(v string) (r *VBtnToggleBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VBtnToggleBuilder) Disabled(v bool) (r *VBtnToggleBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VBtnToggleBuilder) SelectedClass(v string) (r *VBtnToggleBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VBtnToggleBuilder) Max(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) Multiple(v bool) (r *VBtnToggleBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VBtnToggleBuilder) Mandatory(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) Divided(v bool) (r *VBtnToggleBuilder) {
	b.Attr(":divided", fmt.Sprint(v))
	return b
}

func (b *VBtnToggleBuilder) On(name string, value string) (r *VBtnToggleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBtnToggleBuilder) Bind(name string, value string) (r *VBtnToggleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VBtnToggleBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VBtnToggleBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VBtnToggleBuilder) Slot(name string, child ...h.HTMLComponent) (r *VBtnToggleBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VBtnToggleBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VBtnToggleBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VBtnToggleBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VBtnToggleBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VBtnToggleBuilder) SlotDefault(child ...h.HTMLComponent) (r *VBtnToggleBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VBtnToggleBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VBtnToggleBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
