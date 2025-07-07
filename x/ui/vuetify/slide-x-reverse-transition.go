package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSlideXReverseTransitionBuilder struct {
	VTagBuilder[*VSlideXReverseTransitionBuilder]
}

func VSlideXReverseTransition(children ...h.HTMLComponent) *VSlideXReverseTransitionBuilder {
	return VTag(&VSlideXReverseTransitionBuilder{}, "v-slide-x-reverse-transition", children...)
}

func (b *VSlideXReverseTransitionBuilder) Disabled(v bool) (r *VSlideXReverseTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSlideXReverseTransitionBuilder) Group(v bool) (r *VSlideXReverseTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VSlideXReverseTransitionBuilder) HideOnLeave(v bool) (r *VSlideXReverseTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VSlideXReverseTransitionBuilder) LeaveAbsolute(v bool) (r *VSlideXReverseTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VSlideXReverseTransitionBuilder) Mode(v string) (r *VSlideXReverseTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VSlideXReverseTransitionBuilder) Origin(v string) (r *VSlideXReverseTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VSlideXReverseTransitionBuilder) On(name string, value string) (r *VSlideXReverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSlideXReverseTransitionBuilder) Bind(name string, value string) (r *VSlideXReverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
