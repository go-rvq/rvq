package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VLazyBuilder struct {
	VTagBuilder[*VLazyBuilder]
}

func VLazy(children ...h.HTMLComponent) *VLazyBuilder {
	return VTag(&VLazyBuilder{}, "v-lazy", children...)
}

func (b *VLazyBuilder) Tag(v interface{}) (r *VLazyBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) Height(v interface{}) (r *VLazyBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) MaxHeight(v interface{}) (r *VLazyBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) MaxWidth(v interface{}) (r *VLazyBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) MinHeight(v interface{}) (r *VLazyBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) MinWidth(v interface{}) (r *VLazyBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) Width(v interface{}) (r *VLazyBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) ModelValue(v bool) (r *VLazyBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VLazyBuilder) Transition(v interface{}) (r *VLazyBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) Options(v interface{}) (r *VLazyBuilder) {
	b.Attr(":options", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) On(name string, value string) (r *VLazyBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VLazyBuilder) Bind(name string, value string) (r *VLazyBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VLazyBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VLazyBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VLazyBuilder) Slot(name string, child ...h.HTMLComponent) (r *VLazyBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VLazyBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VLazyBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VLazyBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VLazyBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VLazyBuilder) SlotDefault(child ...h.HTMLComponent) (r *VLazyBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VLazyBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VLazyBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
