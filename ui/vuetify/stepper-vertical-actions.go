package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VStepperVerticalActionsBuilder struct {
	VTagBuilder[*VStepperVerticalActionsBuilder]
}

func VStepperVerticalActions(children ...h.HTMLComponent) *VStepperVerticalActionsBuilder {
	return VTag(&VStepperVerticalActionsBuilder{}, "v-stepper-vertical-actions", children...)
}

func (b *VStepperVerticalActionsBuilder) Color(v string) (r *VStepperVerticalActionsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VStepperVerticalActionsBuilder) Disabled(v interface{}) (r *VStepperVerticalActionsBuilder) {
	b.Attr(":disabled", h.JSONString(v))
	return b
}

func (b *VStepperVerticalActionsBuilder) PrevText(v string) (r *VStepperVerticalActionsBuilder) {
	b.Attr("prev-text", v)
	return b
}

func (b *VStepperVerticalActionsBuilder) NextText(v string) (r *VStepperVerticalActionsBuilder) {
	b.Attr("next-text", v)
	return b
}

func (b *VStepperVerticalActionsBuilder) On(name string, value string) (r *VStepperVerticalActionsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperVerticalActionsBuilder) Bind(name string, value string) (r *VStepperVerticalActionsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
