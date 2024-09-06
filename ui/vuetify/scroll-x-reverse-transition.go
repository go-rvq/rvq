package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VScrollXReverseTransitionBuilder struct {
	VTagBuilder[*VScrollXReverseTransitionBuilder]
}

func VScrollXReverseTransition(children ...h.HTMLComponent) *VScrollXReverseTransitionBuilder {
	return VTag(&VScrollXReverseTransitionBuilder{}, "v-scroll-x-reverse-transition", children...)
}

func (b *VScrollXReverseTransitionBuilder) Disabled(v bool) (r *VScrollXReverseTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VScrollXReverseTransitionBuilder) Group(v bool) (r *VScrollXReverseTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VScrollXReverseTransitionBuilder) HideOnLeave(v bool) (r *VScrollXReverseTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VScrollXReverseTransitionBuilder) LeaveAbsolute(v bool) (r *VScrollXReverseTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VScrollXReverseTransitionBuilder) Mode(v string) (r *VScrollXReverseTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VScrollXReverseTransitionBuilder) Origin(v string) (r *VScrollXReverseTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VScrollXReverseTransitionBuilder) On(name string, value string) (r *VScrollXReverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VScrollXReverseTransitionBuilder) Bind(name string, value string) (r *VScrollXReverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
