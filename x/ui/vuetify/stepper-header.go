package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VStepperHeaderBuilder struct {
	VTagBuilder[*VStepperHeaderBuilder]
}

func VStepperHeader(children ...h.HTMLComponent) *VStepperHeaderBuilder {
	return VTag(&VStepperHeaderBuilder{}, "v-stepper-header", children...)
}

func (b *VStepperHeaderBuilder) Tag(v string) (r *VStepperHeaderBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VStepperHeaderBuilder) On(name string, value string) (r *VStepperHeaderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperHeaderBuilder) Bind(name string, value string) (r *VStepperHeaderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
