package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSlideYTransitionBuilder struct {
	VTagBuilder[*VSlideYTransitionBuilder]
}

func VSlideYTransition(children ...h.HTMLComponent) *VSlideYTransitionBuilder {
	return VTag(&VSlideYTransitionBuilder{}, "v-slide-y-transition", children...)
}

func (b *VSlideYTransitionBuilder) Disabled(v bool) (r *VSlideYTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSlideYTransitionBuilder) Group(v bool) (r *VSlideYTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VSlideYTransitionBuilder) HideOnLeave(v bool) (r *VSlideYTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VSlideYTransitionBuilder) LeaveAbsolute(v bool) (r *VSlideYTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VSlideYTransitionBuilder) Mode(v string) (r *VSlideYTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VSlideYTransitionBuilder) Origin(v string) (r *VSlideYTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VSlideYTransitionBuilder) On(name string, value string) (r *VSlideYTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSlideYTransitionBuilder) Bind(name string, value string) (r *VSlideYTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
