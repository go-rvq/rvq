package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VStepperWindowBuilder struct {
	VTagBuilder[*VStepperWindowBuilder]
}

func VStepperWindow(children ...h.HTMLComponent) *VStepperWindowBuilder {
	return VTag(&VStepperWindowBuilder{}, "v-stepper-window", children...)
}

func (b *VStepperWindowBuilder) Reverse(v bool) (r *VStepperWindowBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VStepperWindowBuilder) Direction(v interface{}) (r *VStepperWindowBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VStepperWindowBuilder) ModelValue(v interface{}) (r *VStepperWindowBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VStepperWindowBuilder) Disabled(v bool) (r *VStepperWindowBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VStepperWindowBuilder) SelectedClass(v string) (r *VStepperWindowBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VStepperWindowBuilder) Tag(v string) (r *VStepperWindowBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VStepperWindowBuilder) Theme(v string) (r *VStepperWindowBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VStepperWindowBuilder) On(name string, value string) (r *VStepperWindowBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperWindowBuilder) Bind(name string, value string) (r *VStepperWindowBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
