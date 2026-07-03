package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSystemBarBuilder struct {
	VTagBuilder[*VSystemBarBuilder]
}

func VSystemBar(children ...h.HTMLComponent) *VSystemBarBuilder {
	return VTag(&VSystemBarBuilder{}, "v-system-bar", children...)
}

func (b *VSystemBarBuilder) Height(v interface{}) (r *VSystemBarBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VSystemBarBuilder) Elevation(v interface{}) (r *VSystemBarBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VSystemBarBuilder) Absolute(v bool) (r *VSystemBarBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VSystemBarBuilder) Rounded(v interface{}) (r *VSystemBarBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VSystemBarBuilder) Tile(v bool) (r *VSystemBarBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VSystemBarBuilder) Tag(v interface{}) (r *VSystemBarBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VSystemBarBuilder) Theme(v string) (r *VSystemBarBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSystemBarBuilder) Color(v string) (r *VSystemBarBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSystemBarBuilder) Name(v string) (r *VSystemBarBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VSystemBarBuilder) Order(v interface{}) (r *VSystemBarBuilder) {
	b.Attr(":order", h.JSONString(v))
	return b
}

func (b *VSystemBarBuilder) Window(v bool) (r *VSystemBarBuilder) {
	b.Attr(":window", fmt.Sprint(v))
	return b
}

func (b *VSystemBarBuilder) On(name string, value string) (r *VSystemBarBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSystemBarBuilder) Bind(name string, value string) (r *VSystemBarBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VSystemBarBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSystemBarBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSystemBarBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSystemBarBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSystemBarBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSystemBarBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VSystemBarBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VSystemBarBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VSystemBarBuilder) SlotDefault(child ...h.HTMLComponent) (r *VSystemBarBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VSystemBarBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VSystemBarBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
