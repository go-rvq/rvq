package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VLabelBuilder struct {
	VTagBuilder[*VLabelBuilder]
}

func VLabel(children ...h.HTMLComponent) *VLabelBuilder {
	return VTag(&VLabelBuilder{}, "v-label", children...)
}

func (b *VLabelBuilder) Text(v string) (r *VLabelBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VLabelBuilder) Theme(v string) (r *VLabelBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VLabelBuilder) On(name string, value string) (r *VLabelBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VLabelBuilder) Bind(name string, value string) (r *VLabelBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
