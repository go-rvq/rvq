package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VCardTextBuilder struct {
	VTagBuilder[*VCardTextBuilder]
}

func VCardText(children ...h.HTMLComponent) *VCardTextBuilder {
	return VTag(&VCardTextBuilder{}, "v-card-text", children...)
}

func (b *VCardTextBuilder) Opacity(v interface{}) (r *VCardTextBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VCardTextBuilder) Tag(v string) (r *VCardTextBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VCardTextBuilder) On(name string, value string) (r *VCardTextBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCardTextBuilder) Bind(name string, value string) (r *VCardTextBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
