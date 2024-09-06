package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSlideXTransitionBuilder struct {
	VTagBuilder[*VSlideXTransitionBuilder]
}

func VSlideXTransition(children ...h.HTMLComponent) *VSlideXTransitionBuilder {
	return VTag(&VSlideXTransitionBuilder{}, "v-slide-x-transition", children...)
}

func (b *VSlideXTransitionBuilder) Disabled(v bool) (r *VSlideXTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSlideXTransitionBuilder) Group(v bool) (r *VSlideXTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VSlideXTransitionBuilder) HideOnLeave(v bool) (r *VSlideXTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VSlideXTransitionBuilder) LeaveAbsolute(v bool) (r *VSlideXTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VSlideXTransitionBuilder) Mode(v string) (r *VSlideXTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VSlideXTransitionBuilder) Origin(v string) (r *VSlideXTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VSlideXTransitionBuilder) On(name string, value string) (r *VSlideXTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSlideXTransitionBuilder) Bind(name string, value string) (r *VSlideXTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
