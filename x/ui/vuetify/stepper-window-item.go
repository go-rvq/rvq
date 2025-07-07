package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VStepperWindowItemBuilder struct {
	VTagBuilder[*VStepperWindowItemBuilder]
}

func VStepperWindowItem(children ...h.HTMLComponent) *VStepperWindowItemBuilder {
	return VTag(&VStepperWindowItemBuilder{}, "v-stepper-window-item", children...)
}

func (b *VStepperWindowItemBuilder) ReverseTransition(v interface{}) (r *VStepperWindowItemBuilder) {
	b.Attr(":reverse-transition", h.JSONString(v))
	return b
}

func (b *VStepperWindowItemBuilder) Transition(v interface{}) (r *VStepperWindowItemBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VStepperWindowItemBuilder) Value(v interface{}) (r *VStepperWindowItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VStepperWindowItemBuilder) Disabled(v bool) (r *VStepperWindowItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VStepperWindowItemBuilder) SelectedClass(v string) (r *VStepperWindowItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VStepperWindowItemBuilder) Eager(v bool) (r *VStepperWindowItemBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VStepperWindowItemBuilder) On(name string, value string) (r *VStepperWindowItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperWindowItemBuilder) Bind(name string, value string) (r *VStepperWindowItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
