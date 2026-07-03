package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VPickerBuilder struct {
	VTagBuilder[*VPickerBuilder]
}

func VPicker(children ...h.HTMLComponent) *VPickerBuilder {
	return VTag(&VPickerBuilder{}, "v-picker", children...)
}

func (b *VPickerBuilder) Title(v string) (r *VPickerBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VPickerBuilder) Border(v interface{}) (r *VPickerBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Height(v interface{}) (r *VPickerBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) MaxHeight(v interface{}) (r *VPickerBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) MaxWidth(v interface{}) (r *VPickerBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) MinHeight(v interface{}) (r *VPickerBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) MinWidth(v interface{}) (r *VPickerBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Width(v interface{}) (r *VPickerBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Elevation(v interface{}) (r *VPickerBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Location(v interface{}) (r *VPickerBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Position(v interface{}) (r *VPickerBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Rounded(v interface{}) (r *VPickerBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Tile(v bool) (r *VPickerBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VPickerBuilder) Tag(v interface{}) (r *VPickerBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Theme(v string) (r *VPickerBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VPickerBuilder) Color(v string) (r *VPickerBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VPickerBuilder) BgColor(v string) (r *VPickerBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VPickerBuilder) Divided(v bool) (r *VPickerBuilder) {
	b.Attr(":divided", fmt.Sprint(v))
	return b
}

func (b *VPickerBuilder) HideHeader(v bool) (r *VPickerBuilder) {
	b.Attr(":hide-header", fmt.Sprint(v))
	return b
}

func (b *VPickerBuilder) Landscape(v bool) (r *VPickerBuilder) {
	b.Attr(":landscape", fmt.Sprint(v))
	return b
}

func (b *VPickerBuilder) On(name string, value string) (r *VPickerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VPickerBuilder) Bind(name string, value string) (r *VPickerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VPickerBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VPickerBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VPickerBuilder) Slot(name string, child ...h.HTMLComponent) (r *VPickerBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VPickerBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VPickerBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VPickerBuilder) SetSlotHeader(child ...h.HTMLComponent) {
	b.SetSlot("header", child...)
}

func (b *VPickerBuilder) SetScopedSlotHeader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header", scope, child...)
}

func (b *VPickerBuilder) SlotHeader(child ...h.HTMLComponent) (r *VPickerBuilder) {
	b.SetSlotHeader(child...)
	return b
}

func (b *VPickerBuilder) ScopedSlotHeader(scope string, child ...h.HTMLComponent) (r *VPickerBuilder) {
	b.SetScopedSlotHeader(scope, child...)
	return b
}

func (b *VPickerBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VPickerBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VPickerBuilder) SlotDefault(child ...h.HTMLComponent) (r *VPickerBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VPickerBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VPickerBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VPickerBuilder) SetSlotActions(child ...h.HTMLComponent) {
	b.SetSlot("actions", child...)
}

func (b *VPickerBuilder) SetScopedSlotActions(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("actions", scope, child...)
}

func (b *VPickerBuilder) SlotActions(child ...h.HTMLComponent) (r *VPickerBuilder) {
	b.SetSlotActions(child...)
	return b
}

func (b *VPickerBuilder) ScopedSlotActions(scope string, child ...h.HTMLComponent) (r *VPickerBuilder) {
	b.SetScopedSlotActions(scope, child...)
	return b
}

func (b *VPickerBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VPickerBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VPickerBuilder) SlotTitle(child ...h.HTMLComponent) (r *VPickerBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VPickerBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VPickerBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}
