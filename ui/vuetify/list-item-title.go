package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VListItemTitleBuilder struct {
	VTagBuilder[*VListItemTitleBuilder]
}

func VListItemTitle(children ...h.HTMLComponent) *VListItemTitleBuilder {
	return VTag(&VListItemTitleBuilder{}, "v-list-item-title", children...)
}

func (b *VListItemTitleBuilder) Tag(v string) (r *VListItemTitleBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VListItemTitleBuilder) On(name string, value string) (r *VListItemTitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListItemTitleBuilder) Bind(name string, value string) (r *VListItemTitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
