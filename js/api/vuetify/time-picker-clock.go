package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTimePickerClockBuilder struct {
	VTagBuilder[*VTimePickerClockBuilder]
}

func VTimePickerClock(children ...h.HTMLComponent) *VTimePickerClockBuilder {
	return VTag(&VTimePickerClockBuilder{}, "v-time-picker-clock", children...)
}

func (b *VTimePickerClockBuilder) Disabled(v bool) (r *VTimePickerClockBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) Color(v string) (r *VTimePickerClockBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTimePickerClockBuilder) ModelValue(v interface{}) (r *VTimePickerClockBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VTimePickerClockBuilder) Readonly(v bool) (r *VTimePickerClockBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) Max(v interface{}) (r *VTimePickerClockBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VTimePickerClockBuilder) Min(v interface{}) (r *VTimePickerClockBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VTimePickerClockBuilder) Step(v interface{}) (r *VTimePickerClockBuilder) {
	b.Attr(":step", h.JSONString(v))
	return b
}

func (b *VTimePickerClockBuilder) Rotate(v interface{}) (r *VTimePickerClockBuilder) {
	b.Attr(":rotate", h.JSONString(v))
	return b
}

func (b *VTimePickerClockBuilder) Scrollable(v bool) (r *VTimePickerClockBuilder) {
	b.Attr(":scrollable", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) Ampm(v bool) (r *VTimePickerClockBuilder) {
	b.Attr(":ampm", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) DisplayedValue(v interface{}) (r *VTimePickerClockBuilder) {
	b.Attr(":displayed-value", h.JSONString(v))
	return b
}

func (b *VTimePickerClockBuilder) Double(v bool) (r *VTimePickerClockBuilder) {
	b.Attr(":double", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) Format(v interface{}) (r *VTimePickerClockBuilder) {
	b.Attr(":format", h.JSONString(v))
	return b
}

func (b *VTimePickerClockBuilder) AllowedValues(v interface{}) (r *VTimePickerClockBuilder) {
	b.Attr(":allowed-values", h.JSONString(v))
	return b
}

func (b *VTimePickerClockBuilder) On(name string, value string) (r *VTimePickerClockBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTimePickerClockBuilder) Bind(name string, value string) (r *VTimePickerClockBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VTimePickerClockBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VTimePickerClockBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VTimePickerClockBuilder) Slot(name string, child ...h.HTMLComponent) (r *VTimePickerClockBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VTimePickerClockBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VTimePickerClockBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VTimePickerClockBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VTimePickerClockBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VTimePickerClockBuilder) SlotDefault(child ...h.HTMLComponent) (r *VTimePickerClockBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VTimePickerClockBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VTimePickerClockBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
