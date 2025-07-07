package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VScaleTransitionBuilder struct {
	VTagBuilder[*VScaleTransitionBuilder]
}

func VScaleTransition(children ...h.HTMLComponent) *VScaleTransitionBuilder {
	return VTag(&VScaleTransitionBuilder{}, "v-scale-transition", children...)
}

func (b *VScaleTransitionBuilder) Disabled(v bool) (r *VScaleTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VScaleTransitionBuilder) Group(v bool) (r *VScaleTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VScaleTransitionBuilder) HideOnLeave(v bool) (r *VScaleTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VScaleTransitionBuilder) LeaveAbsolute(v bool) (r *VScaleTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VScaleTransitionBuilder) Mode(v string) (r *VScaleTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VScaleTransitionBuilder) Origin(v string) (r *VScaleTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VScaleTransitionBuilder) On(name string, value string) (r *VScaleTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VScaleTransitionBuilder) Bind(name string, value string) (r *VScaleTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
