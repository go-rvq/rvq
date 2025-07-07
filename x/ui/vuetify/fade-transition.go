package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VFadeTransitionBuilder struct {
	VTagBuilder[*VFadeTransitionBuilder]
}

func VFadeTransition(children ...h.HTMLComponent) *VFadeTransitionBuilder {
	return VTag(&VFadeTransitionBuilder{}, "v-fade-transition", children...)
}

func (b *VFadeTransitionBuilder) Disabled(v bool) (r *VFadeTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VFadeTransitionBuilder) Group(v bool) (r *VFadeTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VFadeTransitionBuilder) HideOnLeave(v bool) (r *VFadeTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VFadeTransitionBuilder) LeaveAbsolute(v bool) (r *VFadeTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VFadeTransitionBuilder) Mode(v string) (r *VFadeTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VFadeTransitionBuilder) Origin(v string) (r *VFadeTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VFadeTransitionBuilder) On(name string, value string) (r *VFadeTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFadeTransitionBuilder) Bind(name string, value string) (r *VFadeTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
