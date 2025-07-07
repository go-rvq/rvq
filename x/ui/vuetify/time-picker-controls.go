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

func (b *VTimePickerControlsBuilder) Ampm(v bool) (r *VTimePickerControlsBuilder) {
	b.Attr(":ampm", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) AmpmReadonly(v bool) (r *VTimePickerControlsBuilder) {
	b.Attr(":ampm-readonly", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) Color(v string) (r *VTimePickerControlsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTimePickerControlsBuilder) Disabled(v bool) (r *VTimePickerControlsBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) Hour(v int) (r *VTimePickerControlsBuilder) {
	b.Attr(":hour", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) Minute(v int) (r *VTimePickerControlsBuilder) {
	b.Attr(":minute", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) Second(v int) (r *VTimePickerControlsBuilder) {
	b.Attr(":second", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) Period(v string) (r *VTimePickerControlsBuilder) {
	b.Attr("period", v)
	return b
}

func (b *VTimePickerControlsBuilder) Readonly(v bool) (r *VTimePickerControlsBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) UseSeconds(v bool) (r *VTimePickerControlsBuilder) {
	b.Attr(":use-seconds", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) Selecting(v int) (r *VTimePickerControlsBuilder) {
	b.Attr(":selecting", fmt.Sprint(v))
	return b
}

func (b *VTimePickerControlsBuilder) Value(v int) (r *VTimePickerControlsBuilder) {
	b.Attr(":value", fmt.Sprint(v))
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
