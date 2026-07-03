package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VAppBarBuilder struct {
	VTagBuilder[*VAppBarBuilder]
}

func VAppBar(children ...h.HTMLComponent) *VAppBarBuilder {
	return VTag(&VAppBarBuilder{}, "v-app-bar", children...)
}

func (b *VAppBarBuilder) Title(v string) (r *VAppBarBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VAppBarBuilder) Flat(v bool) (r *VAppBarBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) Border(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) ModelValue(v bool) (r *VAppBarBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) Density(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Height(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Elevation(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Location(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Absolute(v bool) (r *VAppBarBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) Rounded(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Tile(v bool) (r *VAppBarBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) Tag(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Theme(v string) (r *VAppBarBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VAppBarBuilder) Color(v string) (r *VAppBarBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VAppBarBuilder) Name(v string) (r *VAppBarBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VAppBarBuilder) Image(v string) (r *VAppBarBuilder) {
	b.Attr("image", v)
	return b
}

func (b *VAppBarBuilder) Collapse(v bool) (r *VAppBarBuilder) {
	b.Attr(":collapse", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) Extended(v bool) (r *VAppBarBuilder) {
	b.Attr(":extended", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) ExtensionHeight(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":extension-height", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Floating(v bool) (r *VAppBarBuilder) {
	b.Attr(":floating", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) Order(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":order", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) ScrollTarget(v string) (r *VAppBarBuilder) {
	b.Attr("scroll-target", v)
	return b
}

func (b *VAppBarBuilder) ScrollThreshold(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":scroll-threshold", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) ScrollBehavior(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":scroll-behavior", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) On(name string, value string) (r *VAppBarBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAppBarBuilder) Bind(name string, value string) (r *VAppBarBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VAppBarBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VAppBarBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VAppBarBuilder) Slot(name string, child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VAppBarBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VAppBarBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VAppBarBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VAppBarBuilder) SlotDefault(child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VAppBarBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VAppBarBuilder) SetSlotImage(child ...h.HTMLComponent) {
	b.SetSlot("image", child...)
}

func (b *VAppBarBuilder) SetScopedSlotImage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("image", scope, child...)
}

func (b *VAppBarBuilder) SlotImage(child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetSlotImage(child...)
	return b
}

func (b *VAppBarBuilder) ScopedSlotImage(scope string, child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetScopedSlotImage(scope, child...)
	return b
}

func (b *VAppBarBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VAppBarBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VAppBarBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VAppBarBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VAppBarBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VAppBarBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VAppBarBuilder) SlotAppend(child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VAppBarBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VAppBarBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VAppBarBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VAppBarBuilder) SlotTitle(child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VAppBarBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VAppBarBuilder) SetSlotExtension(child ...h.HTMLComponent) {
	b.SetSlot("extension", child...)
}

func (b *VAppBarBuilder) SetScopedSlotExtension(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("extension", scope, child...)
}

func (b *VAppBarBuilder) SlotExtension(child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetSlotExtension(child...)
	return b
}

func (b *VAppBarBuilder) ScopedSlotExtension(scope string, child ...h.HTMLComponent) (r *VAppBarBuilder) {
	b.SetScopedSlotExtension(scope, child...)
	return b
}
