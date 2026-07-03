package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VAvatarBuilder struct {
	VTagBuilder[*VAvatarBuilder]
}

func VAvatar(children ...h.HTMLComponent) *VAvatarBuilder {
	return VTag(&VAvatarBuilder{}, "v-avatar", children...)
}

func (b *VAvatarBuilder) Text(v string) (r *VAvatarBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VAvatarBuilder) Border(v interface{}) (r *VAvatarBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VAvatarBuilder) End(v bool) (r *VAvatarBuilder) {
	b.Attr(":end", fmt.Sprint(v))
	return b
}

func (b *VAvatarBuilder) Start(v bool) (r *VAvatarBuilder) {
	b.Attr(":start", fmt.Sprint(v))
	return b
}

func (b *VAvatarBuilder) Icon(v interface{}) (r *VAvatarBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VAvatarBuilder) Density(v interface{}) (r *VAvatarBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VAvatarBuilder) Rounded(v interface{}) (r *VAvatarBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VAvatarBuilder) Tile(v bool) (r *VAvatarBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VAvatarBuilder) Tag(v interface{}) (r *VAvatarBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VAvatarBuilder) Theme(v string) (r *VAvatarBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VAvatarBuilder) Color(v string) (r *VAvatarBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VAvatarBuilder) Variant(v interface{}) (r *VAvatarBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VAvatarBuilder) Image(v string) (r *VAvatarBuilder) {
	b.Attr("image", v)
	return b
}

func (b *VAvatarBuilder) Size(v interface{}) (r *VAvatarBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VAvatarBuilder) On(name string, value string) (r *VAvatarBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAvatarBuilder) Bind(name string, value string) (r *VAvatarBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VAvatarBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VAvatarBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VAvatarBuilder) Slot(name string, child ...h.HTMLComponent) (r *VAvatarBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VAvatarBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VAvatarBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VAvatarBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VAvatarBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VAvatarBuilder) SlotDefault(child ...h.HTMLComponent) (r *VAvatarBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VAvatarBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VAvatarBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
