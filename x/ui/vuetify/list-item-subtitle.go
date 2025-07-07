package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VListItemSubtitleBuilder struct {
	VTagBuilder[*VListItemSubtitleBuilder]
}

func VListItemSubtitle(children ...h.HTMLComponent) *VListItemSubtitleBuilder {
	return VTag(&VListItemSubtitleBuilder{}, "v-list-item-subtitle", children...)
}

func (b *VListItemSubtitleBuilder) Opacity(v interface{}) (r *VListItemSubtitleBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VListItemSubtitleBuilder) Tag(v string) (r *VListItemSubtitleBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VListItemSubtitleBuilder) On(name string, value string) (r *VListItemSubtitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListItemSubtitleBuilder) Bind(name string, value string) (r *VListItemSubtitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
