package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VProgressCircularBuilder struct {
	VTagBuilder[*VProgressCircularBuilder]
}

func VProgressCircular(children ...h.HTMLComponent) *VProgressCircularBuilder {
	return VTag(&VProgressCircularBuilder{}, "v-progress-circular", children...)
}

func (b *VProgressCircularBuilder) ModelValue(v interface{}) (r *VProgressCircularBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VProgressCircularBuilder) Width(v interface{}) (r *VProgressCircularBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VProgressCircularBuilder) Tag(v interface{}) (r *VProgressCircularBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VProgressCircularBuilder) Theme(v string) (r *VProgressCircularBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VProgressCircularBuilder) Color(v string) (r *VProgressCircularBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VProgressCircularBuilder) Size(v interface{}) (r *VProgressCircularBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VProgressCircularBuilder) BgColor(v string) (r *VProgressCircularBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VProgressCircularBuilder) Indeterminate(v interface{}) (r *VProgressCircularBuilder) {
	b.Attr(":indeterminate", h.JSONString(v))
	return b
}

func (b *VProgressCircularBuilder) Rotate(v interface{}) (r *VProgressCircularBuilder) {
	b.Attr(":rotate", h.JSONString(v))
	return b
}

func (b *VProgressCircularBuilder) On(name string, value string) (r *VProgressCircularBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VProgressCircularBuilder) Bind(name string, value string) (r *VProgressCircularBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VProgressCircularBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VProgressCircularBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VProgressCircularBuilder) Slot(name string, child ...h.HTMLComponent) (r *VProgressCircularBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VProgressCircularBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VProgressCircularBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VProgressCircularBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VProgressCircularBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VProgressCircularBuilder) SlotDefault(child ...h.HTMLComponent) (r *VProgressCircularBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VProgressCircularBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VProgressCircularBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
