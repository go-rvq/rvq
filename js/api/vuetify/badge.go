package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBadgeBuilder struct {
	VTagBuilder[*VBadgeBuilder]
}

func VBadge(children ...h.HTMLComponent) *VBadgeBuilder {
	return VTag(&VBadgeBuilder{}, "v-badge", children...)
}

func (b *VBadgeBuilder) Icon(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) ModelValue(v bool) (r *VBadgeBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VBadgeBuilder) Location(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) Rounded(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) Tile(v bool) (r *VBadgeBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VBadgeBuilder) Tag(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) Theme(v string) (r *VBadgeBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VBadgeBuilder) Color(v string) (r *VBadgeBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBadgeBuilder) Floating(v bool) (r *VBadgeBuilder) {
	b.Attr(":floating", fmt.Sprint(v))
	return b
}

func (b *VBadgeBuilder) Bordered(v bool) (r *VBadgeBuilder) {
	b.Attr(":bordered", fmt.Sprint(v))
	return b
}

func (b *VBadgeBuilder) Content(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":content", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) Dot(v bool) (r *VBadgeBuilder) {
	b.Attr(":dot", fmt.Sprint(v))
	return b
}

func (b *VBadgeBuilder) Inline(v bool) (r *VBadgeBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VBadgeBuilder) Label(v string) (r *VBadgeBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VBadgeBuilder) Max(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) OffsetX(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":offset-x", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) OffsetY(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":offset-y", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) TextColor(v string) (r *VBadgeBuilder) {
	b.Attr("text-color", v)
	return b
}

func (b *VBadgeBuilder) Transition(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) On(name string, value string) (r *VBadgeBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBadgeBuilder) Bind(name string, value string) (r *VBadgeBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VBadgeBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VBadgeBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VBadgeBuilder) Slot(name string, child ...h.HTMLComponent) (r *VBadgeBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VBadgeBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VBadgeBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VBadgeBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VBadgeBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VBadgeBuilder) SlotDefault(child ...h.HTMLComponent) (r *VBadgeBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VBadgeBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VBadgeBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VBadgeBuilder) SetSlotBadge(child ...h.HTMLComponent) {
	b.SetSlot("badge", child...)
}

func (b *VBadgeBuilder) SetScopedSlotBadge(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("badge", scope, child...)
}

func (b *VBadgeBuilder) SlotBadge(child ...h.HTMLComponent) (r *VBadgeBuilder) {
	b.SetSlotBadge(child...)
	return b
}

func (b *VBadgeBuilder) ScopedSlotBadge(scope string, child ...h.HTMLComponent) (r *VBadgeBuilder) {
	b.SetScopedSlotBadge(scope, child...)
	return b
}
