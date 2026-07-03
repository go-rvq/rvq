package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSheetBuilder struct {
	VTagBuilder[*VSheetBuilder]
}

func VSheet(children ...h.HTMLComponent) *VSheetBuilder {
	return VTag(&VSheetBuilder{}, "v-sheet", children...)
}

func (b *VSheetBuilder) Tag(v interface{}) (r *VSheetBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Theme(v string) (r *VSheetBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSheetBuilder) Border(v interface{}) (r *VSheetBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Height(v interface{}) (r *VSheetBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) MaxHeight(v interface{}) (r *VSheetBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) MaxWidth(v interface{}) (r *VSheetBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) MinHeight(v interface{}) (r *VSheetBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) MinWidth(v interface{}) (r *VSheetBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Width(v interface{}) (r *VSheetBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Elevation(v interface{}) (r *VSheetBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Rounded(v interface{}) (r *VSheetBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Tile(v bool) (r *VSheetBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VSheetBuilder) Color(v string) (r *VSheetBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSheetBuilder) Location(v interface{}) (r *VSheetBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Position(v interface{}) (r *VSheetBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) On(name string, value string) (r *VSheetBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSheetBuilder) Bind(name string, value string) (r *VSheetBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VSheetBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSheetBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSheetBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSheetBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSheetBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSheetBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VSheetBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VSheetBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VSheetBuilder) SlotDefault(child ...h.HTMLComponent) (r *VSheetBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VSheetBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VSheetBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
