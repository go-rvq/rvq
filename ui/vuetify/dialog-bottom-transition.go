package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VDialogBottomTransitionBuilder struct {
	VTagBuilder[*VDialogBottomTransitionBuilder]
}

func VDialogBottomTransition(children ...h.HTMLComponent) *VDialogBottomTransitionBuilder {
	return VTag(&VDialogBottomTransitionBuilder{}, "v-dialog-bottom-transition", children...)
}

func (b *VDialogBottomTransitionBuilder) Disabled(v bool) (r *VDialogBottomTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VDialogBottomTransitionBuilder) Group(v bool) (r *VDialogBottomTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VDialogBottomTransitionBuilder) HideOnLeave(v bool) (r *VDialogBottomTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VDialogBottomTransitionBuilder) LeaveAbsolute(v bool) (r *VDialogBottomTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VDialogBottomTransitionBuilder) Mode(v string) (r *VDialogBottomTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VDialogBottomTransitionBuilder) Origin(v string) (r *VDialogBottomTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VDialogBottomTransitionBuilder) On(name string, value string) (r *VDialogBottomTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDialogBottomTransitionBuilder) Bind(name string, value string) (r *VDialogBottomTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
