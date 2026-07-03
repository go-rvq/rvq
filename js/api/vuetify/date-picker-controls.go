package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDatePickerControlsBuilder struct {
	VTagBuilder[*VDatePickerControlsBuilder]
}

func VDatePickerControls(children ...h.HTMLComponent) *VDatePickerControlsBuilder {
	return VTag(&VDatePickerControlsBuilder{}, "v-date-picker-controls", children...)
}

func (b *VDatePickerControlsBuilder) Text(v string) (r *VDatePickerControlsBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VDatePickerControlsBuilder) Disabled(v interface{}) (r *VDatePickerControlsBuilder) {
	b.Attr(":disabled", h.JSONString(v))
	return b
}

func (b *VDatePickerControlsBuilder) Active(v interface{}) (r *VDatePickerControlsBuilder) {
	b.Attr(":active", h.JSONString(v))
	return b
}

func (b *VDatePickerControlsBuilder) PrevIcon(v interface{}) (r *VDatePickerControlsBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VDatePickerControlsBuilder) NextIcon(v interface{}) (r *VDatePickerControlsBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VDatePickerControlsBuilder) ControlHeight(v interface{}) (r *VDatePickerControlsBuilder) {
	b.Attr(":control-height", h.JSONString(v))
	return b
}

func (b *VDatePickerControlsBuilder) ModeIcon(v interface{}) (r *VDatePickerControlsBuilder) {
	b.Attr(":mode-icon", h.JSONString(v))
	return b
}

func (b *VDatePickerControlsBuilder) ViewMode(v interface{}) (r *VDatePickerControlsBuilder) {
	b.Attr(":view-mode", h.JSONString(v))
	return b
}

func (b *VDatePickerControlsBuilder) On(name string, value string) (r *VDatePickerControlsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDatePickerControlsBuilder) Bind(name string, value string) (r *VDatePickerControlsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDatePickerControlsBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDatePickerControlsBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDatePickerControlsBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDatePickerControlsBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDatePickerControlsBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDatePickerControlsBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDatePickerControlsBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VDatePickerControlsBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VDatePickerControlsBuilder) SlotDefault(child ...h.HTMLComponent) (r *VDatePickerControlsBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VDatePickerControlsBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VDatePickerControlsBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
