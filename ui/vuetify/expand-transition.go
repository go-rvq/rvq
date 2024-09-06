package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VExpandTransitionBuilder struct {
	VTagBuilder[*VExpandTransitionBuilder]
}

func VExpandTransition(children ...h.HTMLComponent) *VExpandTransitionBuilder {
	return VTag(&VExpandTransitionBuilder{}, "v-expand-transition", children...)
}

func (b *VExpandTransitionBuilder) Disabled(v bool) (r *VExpandTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VExpandTransitionBuilder) Group(v bool) (r *VExpandTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VExpandTransitionBuilder) Mode(v interface{}) (r *VExpandTransitionBuilder) {
	b.Attr(":mode", h.JSONString(v))
	return b
}

func (b *VExpandTransitionBuilder) On(name string, value string) (r *VExpandTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VExpandTransitionBuilder) Bind(name string, value string) (r *VExpandTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
