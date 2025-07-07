package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSlideYReverseTransitionBuilder struct {
	VTagBuilder[*VSlideYReverseTransitionBuilder]
}

func VSlideYReverseTransition(children ...h.HTMLComponent) *VSlideYReverseTransitionBuilder {
	return VTag(&VSlideYReverseTransitionBuilder{}, "v-slide-y-reverse-transition", children...)
}

func (b *VSlideYReverseTransitionBuilder) Disabled(v bool) (r *VSlideYReverseTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSlideYReverseTransitionBuilder) Group(v bool) (r *VSlideYReverseTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VSlideYReverseTransitionBuilder) HideOnLeave(v bool) (r *VSlideYReverseTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VSlideYReverseTransitionBuilder) LeaveAbsolute(v bool) (r *VSlideYReverseTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VSlideYReverseTransitionBuilder) Mode(v string) (r *VSlideYReverseTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VSlideYReverseTransitionBuilder) Origin(v string) (r *VSlideYReverseTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VSlideYReverseTransitionBuilder) On(name string, value string) (r *VSlideYReverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSlideYReverseTransitionBuilder) Bind(name string, value string) (r *VSlideYReverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
