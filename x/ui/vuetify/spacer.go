package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSpacerBuilder struct {
	VTagBuilder[*VSpacerBuilder]
}

func VSpacer(children ...h.HTMLComponent) *VSpacerBuilder {
	return VTag(&VSpacerBuilder{}, "v-spacer", children...)
}

func (b *VSpacerBuilder) Tag(v string) (r *VSpacerBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VSpacerBuilder) On(name string, value string) (r *VSpacerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSpacerBuilder) Bind(name string, value string) (r *VSpacerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
