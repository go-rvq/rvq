package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VAlertTitleBuilder struct {
	VTagBuilder[*VAlertTitleBuilder]
}

func VAlertTitle(children ...h.HTMLComponent) *VAlertTitleBuilder {
	return VTag(&VAlertTitleBuilder{}, "v-alert-title", children...)
}

func (b *VAlertTitleBuilder) Tag(v string) (r *VAlertTitleBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VAlertTitleBuilder) On(name string, value string) (r *VAlertTitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAlertTitleBuilder) Bind(name string, value string) (r *VAlertTitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
