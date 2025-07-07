package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VItemBuilder struct {
	VTagBuilder[*VItemBuilder]
}

func VItem(children ...h.HTMLComponent) *VItemBuilder {
	return VTag(&VItemBuilder{}, "v-item", children...)
}

func (b *VItemBuilder) Value(v interface{}) (r *VItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VItemBuilder) Disabled(v bool) (r *VItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VItemBuilder) SelectedClass(v string) (r *VItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VItemBuilder) On(name string, value string) (r *VItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VItemBuilder) Bind(name string, value string) (r *VItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
