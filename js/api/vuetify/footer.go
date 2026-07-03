package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VFooterBuilder struct {
	VTagBuilder[*VFooterBuilder]
}

func VFooter(children ...h.HTMLComponent) *VFooterBuilder {
	return VTag(&VFooterBuilder{}, "v-footer", children...)
}

func (b *VFooterBuilder) Tag(v interface{}) (r *VFooterBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VFooterBuilder) Name(v string) (r *VFooterBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VFooterBuilder) Theme(v string) (r *VFooterBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VFooterBuilder) Border(v interface{}) (r *VFooterBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VFooterBuilder) Height(v interface{}) (r *VFooterBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VFooterBuilder) Elevation(v interface{}) (r *VFooterBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VFooterBuilder) Rounded(v interface{}) (r *VFooterBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VFooterBuilder) Tile(v bool) (r *VFooterBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VFooterBuilder) Color(v string) (r *VFooterBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VFooterBuilder) Absolute(v bool) (r *VFooterBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VFooterBuilder) Order(v interface{}) (r *VFooterBuilder) {
	b.Attr(":order", h.JSONString(v))
	return b
}

func (b *VFooterBuilder) App(v bool) (r *VFooterBuilder) {
	b.Attr(":app", fmt.Sprint(v))
	return b
}

func (b *VFooterBuilder) On(name string, value string) (r *VFooterBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFooterBuilder) Bind(name string, value string) (r *VFooterBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VFooterBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VFooterBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VFooterBuilder) Slot(name string, child ...h.HTMLComponent) (r *VFooterBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VFooterBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VFooterBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VFooterBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VFooterBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VFooterBuilder) SlotDefault(child ...h.HTMLComponent) (r *VFooterBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VFooterBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VFooterBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
