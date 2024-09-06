package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VDialogTransitionBuilder struct {
	VTagBuilder[*VDialogTransitionBuilder]
}

func VDialogTransition(children ...h.HTMLComponent) *VDialogTransitionBuilder {
	return VTag(&VDialogTransitionBuilder{}, "v-dialog-transition", children...)
}

func (b *VDialogTransitionBuilder) Target(v interface{}) (r *VDialogTransitionBuilder) {
	b.Attr(":target", h.JSONString(v))
	return b
}

func (b *VDialogTransitionBuilder) On(name string, value string) (r *VDialogTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDialogTransitionBuilder) Bind(name string, value string) (r *VDialogTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
