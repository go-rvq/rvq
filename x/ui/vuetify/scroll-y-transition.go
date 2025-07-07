package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VScrollYTransitionBuilder struct {
	VTagBuilder[*VScrollYTransitionBuilder]
}

func VScrollYTransition(children ...h.HTMLComponent) *VScrollYTransitionBuilder {
	return VTag(&VScrollYTransitionBuilder{}, "v-scroll-y-transition", children...)
}

func (b *VScrollYTransitionBuilder) Disabled(v bool) (r *VScrollYTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VScrollYTransitionBuilder) Group(v bool) (r *VScrollYTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VScrollYTransitionBuilder) HideOnLeave(v bool) (r *VScrollYTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VScrollYTransitionBuilder) LeaveAbsolute(v bool) (r *VScrollYTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VScrollYTransitionBuilder) Mode(v string) (r *VScrollYTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VScrollYTransitionBuilder) Origin(v string) (r *VScrollYTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VScrollYTransitionBuilder) On(name string, value string) (r *VScrollYTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VScrollYTransitionBuilder) Bind(name string, value string) (r *VScrollYTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
