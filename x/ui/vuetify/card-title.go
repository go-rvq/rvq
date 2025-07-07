package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCardTitleBuilder struct {
	VTagBuilder[*VCardTitleBuilder]
}

func VCardTitle(children ...h.HTMLComponent) *VCardTitleBuilder {
	return VTag(&VCardTitleBuilder{}, "v-card-title", children...)
}

func (b *VCardTitleBuilder) Tag(v string) (r *VCardTitleBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VCardTitleBuilder) On(name string, value string) (r *VCardTitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCardTitleBuilder) Bind(name string, value string) (r *VCardTitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
