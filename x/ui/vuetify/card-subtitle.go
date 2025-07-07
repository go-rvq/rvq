package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VCardSubtitleBuilder struct {
	VTagBuilder[*VCardSubtitleBuilder]
}

func VCardSubtitle(children ...h.HTMLComponent) *VCardSubtitleBuilder {
	return VTag(&VCardSubtitleBuilder{}, "v-card-subtitle", children...)
}

func (b *VCardSubtitleBuilder) Opacity(v interface{}) (r *VCardSubtitleBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VCardSubtitleBuilder) Tag(v string) (r *VCardSubtitleBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VCardSubtitleBuilder) On(name string, value string) (r *VCardSubtitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCardSubtitleBuilder) Bind(name string, value string) (r *VCardSubtitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
