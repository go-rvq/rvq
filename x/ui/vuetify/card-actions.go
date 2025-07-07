package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCardActionsBuilder struct {
	VTagBuilder[*VCardActionsBuilder]
}

func VCardActions(children ...h.HTMLComponent) *VCardActionsBuilder {
	return VTag(&VCardActionsBuilder{}, "v-card-actions", children...)
}

func (b *VCardActionsBuilder) On(name string, value string) (r *VCardActionsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCardActionsBuilder) Bind(name string, value string) (r *VCardActionsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
