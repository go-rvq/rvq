package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VStepperActionsBuilder struct {
	VTagBuilder[*VStepperActionsBuilder]
}

func VStepperActions(children ...h.HTMLComponent) *VStepperActionsBuilder {
	return VTag(&VStepperActionsBuilder{}, "v-stepper-actions", children...)
}

func (b *VStepperActionsBuilder) Color(v string) (r *VStepperActionsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VStepperActionsBuilder) Disabled(v interface{}) (r *VStepperActionsBuilder) {
	b.Attr(":disabled", h.JSONString(v))
	return b
}

func (b *VStepperActionsBuilder) PrevText(v string) (r *VStepperActionsBuilder) {
	b.Attr("prev-text", v)
	return b
}

func (b *VStepperActionsBuilder) NextText(v string) (r *VStepperActionsBuilder) {
	b.Attr("next-text", v)
	return b
}

func (b *VStepperActionsBuilder) On(name string, value string) (r *VStepperActionsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperActionsBuilder) Bind(name string, value string) (r *VStepperActionsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
