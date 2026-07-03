package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTimePickerControlsBuilder struct {
	VTagBuilder[*VTimePickerControlsBuilder]
}

func VTimePickerControls(children ...h.HTMLComponent) *VTimePickerControlsBuilder {
	return VTag(&VTimePickerControlsBuilder{}, "v-time-picker-controls", children...)
}

func (b *VTimePickerControlsBuilder) Value(v interface{}) (r *VTimePickerControlsBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VTimePickerControlsBuilder) Disabled(v bool) (r *VTimePickerControlsBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) Color(v string) (r *VTimePickerControlsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTimePickerControlsBuilder) Readonly(v bool) (r *VTimePickerControlsBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) ViewMode(v interface{}) (r *VTimePickerControlsBuilder) {
	b.Attr(":view-mode", h.JSONString(v))
	return b
}

func (b *VTimePickerControlsBuilder) Ampm(v bool) (r *VTimePickerControlsBuilder) {
	b.Attr(":ampm", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) AmpmInTitle(v bool) (r *VTimePickerControlsBuilder) {
	b.Attr(":ampm-in-title", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) AmpmReadonly(v bool) (r *VTimePickerControlsBuilder) {
	b.Attr(":ampm-readonly", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) Hour(v interface{}) (r *VTimePickerControlsBuilder) {
	b.Attr(":hour", h.JSONString(v))
	return b
}

func (b *VTimePickerControlsBuilder) Minute(v interface{}) (r *VTimePickerControlsBuilder) {
	b.Attr(":minute", h.JSONString(v))
	return b
}

func (b *VTimePickerControlsBuilder) Second(v interface{}) (r *VTimePickerControlsBuilder) {
	b.Attr(":second", h.JSONString(v))
	return b
}

func (b *VTimePickerControlsBuilder) UseSeconds(v bool) (r *VTimePickerControlsBuilder) {
	b.Attr(":use-seconds", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) Period(v interface{}) (r *VTimePickerControlsBuilder) {
	b.Attr(":period", h.JSONString(v))
	return b
}

func (b *VTimePickerControlsBuilder) On(name string, value string) (r *VTimePickerControlsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTimePickerControlsBuilder) Bind(name string, value string) (r *VTimePickerControlsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VTimePickerControlsBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VTimePickerControlsBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VTimePickerControlsBuilder) Slot(name string, child ...h.HTMLComponent) (r *VTimePickerControlsBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VTimePickerControlsBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VTimePickerControlsBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VTimePickerControlsBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VTimePickerControlsBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VTimePickerControlsBuilder) SlotDefault(child ...h.HTMLComponent) (r *VTimePickerControlsBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VTimePickerControlsBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VTimePickerControlsBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
