package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VDatePickerControlsBuilder struct {
	VTagBuilder[*VDatePickerControlsBuilder]
}

func VDatePickerControls(children ...h.HTMLComponent) *VDatePickerControlsBuilder {
	return VTag(&VDatePickerControlsBuilder{}, "v-date-picker-controls", children...)
}

func (b *VDatePickerControlsBuilder) Active(v interface{}) (r *VDatePickerControlsBuilder) {
	b.Attr(":active", h.JSONString(v))
	return b
}

func (b *VDatePickerControlsBuilder) Disabled(v interface{}) (r *VDatePickerControlsBuilder) {
	b.Attr(":disabled", h.JSONString(v))
	return b
}

func (b *VDatePickerControlsBuilder) NextIcon(v string) (r *VDatePickerControlsBuilder) {
	b.Attr("next-icon", v)
	return b
}

func (b *VDatePickerControlsBuilder) PrevIcon(v string) (r *VDatePickerControlsBuilder) {
	b.Attr("prev-icon", v)
	return b
}

func (b *VDatePickerControlsBuilder) ModeIcon(v string) (r *VDatePickerControlsBuilder) {
	b.Attr("mode-icon", v)
	return b
}

func (b *VDatePickerControlsBuilder) Text(v string) (r *VDatePickerControlsBuilder) {
	b.Attr("text", v)
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
