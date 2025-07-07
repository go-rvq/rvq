package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VExpandXTransitionBuilder struct {
	VTagBuilder[*VExpandXTransitionBuilder]
}

func VExpandXTransition(children ...h.HTMLComponent) *VExpandXTransitionBuilder {
	return VTag(&VExpandXTransitionBuilder{}, "v-expand-x-transition", children...)
}

func (b *VExpandXTransitionBuilder) Disabled(v bool) (r *VExpandXTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VExpandXTransitionBuilder) Group(v bool) (r *VExpandXTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VExpandXTransitionBuilder) Mode(v interface{}) (r *VExpandXTransitionBuilder) {
	b.Attr(":mode", h.JSONString(v))
	return b
}

func (b *VExpandXTransitionBuilder) On(name string, value string) (r *VExpandXTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VExpandXTransitionBuilder) Bind(name string, value string) (r *VExpandXTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
