package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VToolbarBuilder struct {
	VTagBuilder[*VToolbarBuilder]
}

func VToolbar(children ...h.HTMLComponent) *VToolbarBuilder {
	return VTag(&VToolbarBuilder{}, "v-toolbar", children...)
}

func (b *VToolbarBuilder) Title(v string) (r *VToolbarBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VToolbarBuilder) Flat(v bool) (r *VToolbarBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VToolbarBuilder) Border(v interface{}) (r *VToolbarBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VToolbarBuilder) Density(v interface{}) (r *VToolbarBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VToolbarBuilder) Height(v interface{}) (r *VToolbarBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VToolbarBuilder) Elevation(v interface{}) (r *VToolbarBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VToolbarBuilder) Absolute(v bool) (r *VToolbarBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VToolbarBuilder) Rounded(v interface{}) (r *VToolbarBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VToolbarBuilder) Tile(v bool) (r *VToolbarBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VToolbarBuilder) Tag(v interface{}) (r *VToolbarBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VToolbarBuilder) Theme(v string) (r *VToolbarBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VToolbarBuilder) Color(v string) (r *VToolbarBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VToolbarBuilder) Image(v string) (r *VToolbarBuilder) {
	b.Attr("image", v)
	return b
}

func (b *VToolbarBuilder) Collapse(v bool) (r *VToolbarBuilder) {
	b.Attr(":collapse", fmt.Sprint(v))
	return b
}

func (b *VToolbarBuilder) Extended(v bool) (r *VToolbarBuilder) {
	b.Attr(":extended", fmt.Sprint(v))
	return b
}

func (b *VToolbarBuilder) ExtensionHeight(v interface{}) (r *VToolbarBuilder) {
	b.Attr(":extension-height", h.JSONString(v))
	return b
}

func (b *VToolbarBuilder) Floating(v bool) (r *VToolbarBuilder) {
	b.Attr(":floating", fmt.Sprint(v))
	return b
}

func (b *VToolbarBuilder) On(name string, value string) (r *VToolbarBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VToolbarBuilder) Bind(name string, value string) (r *VToolbarBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VToolbarBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VToolbarBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VToolbarBuilder) Slot(name string, child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VToolbarBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VToolbarBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VToolbarBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VToolbarBuilder) SlotDefault(child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VToolbarBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VToolbarBuilder) SetSlotImage(child ...h.HTMLComponent) {
	b.SetSlot("image", child...)
}

func (b *VToolbarBuilder) SetScopedSlotImage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("image", scope, child...)
}

func (b *VToolbarBuilder) SlotImage(child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetSlotImage(child...)
	return b
}

func (b *VToolbarBuilder) ScopedSlotImage(scope string, child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetScopedSlotImage(scope, child...)
	return b
}

func (b *VToolbarBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VToolbarBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VToolbarBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VToolbarBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VToolbarBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VToolbarBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VToolbarBuilder) SlotAppend(child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VToolbarBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VToolbarBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VToolbarBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VToolbarBuilder) SlotTitle(child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VToolbarBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VToolbarBuilder) SetSlotExtension(child ...h.HTMLComponent) {
	b.SetSlot("extension", child...)
}

func (b *VToolbarBuilder) SetScopedSlotExtension(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("extension", scope, child...)
}

func (b *VToolbarBuilder) SlotExtension(child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetSlotExtension(child...)
	return b
}

func (b *VToolbarBuilder) ScopedSlotExtension(scope string, child ...h.HTMLComponent) (r *VToolbarBuilder) {
	b.SetScopedSlotExtension(scope, child...)
	return b
}
