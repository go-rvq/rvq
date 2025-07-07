package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VToolbarItemsBuilder struct {
	VTagBuilder[*VToolbarItemsBuilder]
}

func VToolbarItems(children ...h.HTMLComponent) *VToolbarItemsBuilder {
	return VTag(&VToolbarItemsBuilder{}, "v-toolbar-items", children...)
}

func (b *VToolbarItemsBuilder) Color(v string) (r *VToolbarItemsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VToolbarItemsBuilder) Variant(v interface{}) (r *VToolbarItemsBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VToolbarItemsBuilder) On(name string, value string) (r *VToolbarItemsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VToolbarItemsBuilder) Bind(name string, value string) (r *VToolbarItemsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
