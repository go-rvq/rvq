package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VListItemMediaBuilder struct {
	VTagBuilder[*VListItemMediaBuilder]
}

func VListItemMedia(children ...h.HTMLComponent) *VListItemMediaBuilder {
	return VTag(&VListItemMediaBuilder{}, "v-list-item-media", children...)
}

func (b *VListItemMediaBuilder) Start(v bool) (r *VListItemMediaBuilder) {
	b.Attr(":start", fmt.Sprint(v))
	return b
}

func (b *VListItemMediaBuilder) End(v bool) (r *VListItemMediaBuilder) {
	b.Attr(":end", fmt.Sprint(v))
	return b
}

func (b *VListItemMediaBuilder) Tag(v string) (r *VListItemMediaBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VListItemMediaBuilder) On(name string, value string) (r *VListItemMediaBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListItemMediaBuilder) Bind(name string, value string) (r *VListItemMediaBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
