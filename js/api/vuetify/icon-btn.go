package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VIconBtnBuilder struct {
	VTagBuilder[*VIconBtnBuilder]
}

func VIconBtn(children ...h.HTMLComponent) *VIconBtnBuilder {
	return VTag(&VIconBtnBuilder{}, "v-icon-btn", children...)
}

func (b *VIconBtnBuilder) Tag(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) Theme(v string) (r *VIconBtnBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VIconBtnBuilder) Size(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) Text(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":text", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) ActiveColor(v string) (r *VIconBtnBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VIconBtnBuilder) Disabled(v bool) (r *VIconBtnBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VIconBtnBuilder) Border(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) Height(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) Width(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) Elevation(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) Rounded(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) Tile(v bool) (r *VIconBtnBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VIconBtnBuilder) Color(v string) (r *VIconBtnBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VIconBtnBuilder) Variant(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) Opacity(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) IconColor(v string) (r *VIconBtnBuilder) {
	b.Attr("icon-color", v)
	return b
}

func (b *VIconBtnBuilder) Readonly(v bool) (r *VIconBtnBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VIconBtnBuilder) Active(v bool) (r *VIconBtnBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VIconBtnBuilder) Loading(v bool) (r *VIconBtnBuilder) {
	b.Attr(":loading", fmt.Sprint(v))
	return b
}

func (b *VIconBtnBuilder) Rotate(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":rotate", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) Icon(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) BaseVariant(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":base-variant", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) HideOverlay(v bool) (r *VIconBtnBuilder) {
	b.Attr(":hide-overlay", fmt.Sprint(v))
	return b
}

func (b *VIconBtnBuilder) IconSize(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":icon-size", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) IconSizes(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":icon-sizes", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) Sizes(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":sizes", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) ActiveIcon(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":active-icon", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) ActiveVariant(v interface{}) (r *VIconBtnBuilder) {
	b.Attr(":active-variant", h.JSONString(v))
	return b
}

func (b *VIconBtnBuilder) On(name string, value string) (r *VIconBtnBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VIconBtnBuilder) Bind(name string, value string) (r *VIconBtnBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VIconBtnBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VIconBtnBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VIconBtnBuilder) Slot(name string, child ...h.HTMLComponent) (r *VIconBtnBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VIconBtnBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VIconBtnBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VIconBtnBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VIconBtnBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VIconBtnBuilder) SlotDefault(child ...h.HTMLComponent) (r *VIconBtnBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VIconBtnBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VIconBtnBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VIconBtnBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VIconBtnBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VIconBtnBuilder) SlotLoader(child ...h.HTMLComponent) (r *VIconBtnBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VIconBtnBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VIconBtnBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}
