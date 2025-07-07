package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSlideGroupItemBuilder struct {
	VTagBuilder[*VSlideGroupItemBuilder]
}

func VSlideGroupItem(children ...h.HTMLComponent) *VSlideGroupItemBuilder {
	return VTag(&VSlideGroupItemBuilder{}, "v-slide-group-item", children...)
}

func (b *VSlideGroupItemBuilder) Value(v interface{}) (r *VSlideGroupItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VSlideGroupItemBuilder) Disabled(v bool) (r *VSlideGroupItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSlideGroupItemBuilder) SelectedClass(v string) (r *VSlideGroupItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VSlideGroupItemBuilder) On(name string, value string) (r *VSlideGroupItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSlideGroupItemBuilder) Bind(name string, value string) (r *VSlideGroupItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
