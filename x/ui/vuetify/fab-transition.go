package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VFabTransitionBuilder struct {
	VTagBuilder[*VFabTransitionBuilder]
}

func VFabTransition(children ...h.HTMLComponent) *VFabTransitionBuilder {
	return VTag(&VFabTransitionBuilder{}, "v-fab-transition", children...)
}

func (b *VFabTransitionBuilder) Disabled(v bool) (r *VFabTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VFabTransitionBuilder) Group(v bool) (r *VFabTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VFabTransitionBuilder) HideOnLeave(v bool) (r *VFabTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VFabTransitionBuilder) LeaveAbsolute(v bool) (r *VFabTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VFabTransitionBuilder) Mode(v string) (r *VFabTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VFabTransitionBuilder) Origin(v string) (r *VFabTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VFabTransitionBuilder) On(name string, value string) (r *VFabTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFabTransitionBuilder) Bind(name string, value string) (r *VFabTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
