package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VFieldLabelBuilder struct {
	VTagBuilder[*VFieldLabelBuilder]
}

func VFieldLabel(children ...h.HTMLComponent) *VFieldLabelBuilder {
	return VTag(&VFieldLabelBuilder{}, "v-field-label", children...)
}

func (b *VFieldLabelBuilder) Floating(v bool) (r *VFieldLabelBuilder) {
	b.Attr(":floating", fmt.Sprint(v))
	return b
}

func (b *VFieldLabelBuilder) On(name string, value string) (r *VFieldLabelBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFieldLabelBuilder) Bind(name string, value string) (r *VFieldLabelBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
