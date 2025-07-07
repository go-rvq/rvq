package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VListItemActionBuilder struct {
	VTagBuilder[*VListItemActionBuilder]
}

func VListItemAction(children ...h.HTMLComponent) *VListItemActionBuilder {
	return VTag(&VListItemActionBuilder{}, "v-list-item-action", children...)
}

func (b *VListItemActionBuilder) Start(v bool) (r *VListItemActionBuilder) {
	b.Attr(":start", fmt.Sprint(v))
	return b
}

func (b *VListItemActionBuilder) End(v bool) (r *VListItemActionBuilder) {
	b.Attr(":end", fmt.Sprint(v))
	return b
}

func (b *VListItemActionBuilder) Tag(v string) (r *VListItemActionBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VListItemActionBuilder) On(name string, value string) (r *VListItemActionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListItemActionBuilder) Bind(name string, value string) (r *VListItemActionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
