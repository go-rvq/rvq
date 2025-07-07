package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDatePickerHeaderBuilder struct {
	VTagBuilder[*VDatePickerHeaderBuilder]
}

func VDatePickerHeader(children ...h.HTMLComponent) *VDatePickerHeaderBuilder {
	return VTag(&VDatePickerHeaderBuilder{}, "v-date-picker-header", children...)
}

func (b *VDatePickerHeaderBuilder) AppendIcon(v string) (r *VDatePickerHeaderBuilder) {
	b.Attr("append-icon", v)
	return b
}

func (b *VDatePickerHeaderBuilder) Color(v string) (r *VDatePickerHeaderBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDatePickerHeaderBuilder) Header(v string) (r *VDatePickerHeaderBuilder) {
	b.Attr("header", v)
	return b
}

func (b *VDatePickerHeaderBuilder) Transition(v string) (r *VDatePickerHeaderBuilder) {
	b.Attr("transition", v)
	return b
}

func (b *VDatePickerHeaderBuilder) On(name string, value string) (r *VDatePickerHeaderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDatePickerHeaderBuilder) Bind(name string, value string) (r *VDatePickerHeaderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
