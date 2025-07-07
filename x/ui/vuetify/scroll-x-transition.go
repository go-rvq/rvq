package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VScrollXTransitionBuilder struct {
	VTagBuilder[*VScrollXTransitionBuilder]
}

func VScrollXTransition(children ...h.HTMLComponent) *VScrollXTransitionBuilder {
	return VTag(&VScrollXTransitionBuilder{}, "v-scroll-x-transition", children...)
}

func (b *VScrollXTransitionBuilder) Disabled(v bool) (r *VScrollXTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VScrollXTransitionBuilder) Group(v bool) (r *VScrollXTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VScrollXTransitionBuilder) HideOnLeave(v bool) (r *VScrollXTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VScrollXTransitionBuilder) LeaveAbsolute(v bool) (r *VScrollXTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VScrollXTransitionBuilder) Mode(v string) (r *VScrollXTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VScrollXTransitionBuilder) Origin(v string) (r *VScrollXTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VScrollXTransitionBuilder) On(name string, value string) (r *VScrollXTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VScrollXTransitionBuilder) Bind(name string, value string) (r *VScrollXTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
