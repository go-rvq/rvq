package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VScrollYReverseTransitionBuilder struct {
	VTagBuilder[*VScrollYReverseTransitionBuilder]
}

func VScrollYReverseTransition(children ...h.HTMLComponent) *VScrollYReverseTransitionBuilder {
	return VTag(&VScrollYReverseTransitionBuilder{}, "v-scroll-y-reverse-transition", children...)
}

func (b *VScrollYReverseTransitionBuilder) Disabled(v bool) (r *VScrollYReverseTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VScrollYReverseTransitionBuilder) Group(v bool) (r *VScrollYReverseTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VScrollYReverseTransitionBuilder) HideOnLeave(v bool) (r *VScrollYReverseTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VScrollYReverseTransitionBuilder) LeaveAbsolute(v bool) (r *VScrollYReverseTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VScrollYReverseTransitionBuilder) Mode(v string) (r *VScrollYReverseTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VScrollYReverseTransitionBuilder) Origin(v string) (r *VScrollYReverseTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VScrollYReverseTransitionBuilder) On(name string, value string) (r *VScrollYReverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VScrollYReverseTransitionBuilder) Bind(name string, value string) (r *VScrollYReverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
