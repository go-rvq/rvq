package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VWindowItemBuilder struct {
	VTagBuilder[*VWindowItemBuilder]
}

func VWindowItem(children ...h.HTMLComponent) *VWindowItemBuilder {
	return VTag(&VWindowItemBuilder{}, "v-window-item", children...)
}

func (b *VWindowItemBuilder) ReverseTransition(v interface{}) (r *VWindowItemBuilder) {
	b.Attr(":reverse-transition", h.JSONString(v))
	return b
}

func (b *VWindowItemBuilder) Transition(v interface{}) (r *VWindowItemBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VWindowItemBuilder) Value(v interface{}) (r *VWindowItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VWindowItemBuilder) Disabled(v bool) (r *VWindowItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VWindowItemBuilder) SelectedClass(v string) (r *VWindowItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VWindowItemBuilder) Eager(v bool) (r *VWindowItemBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VWindowItemBuilder) On(name string, value string) (r *VWindowItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VWindowItemBuilder) Bind(name string, value string) (r *VWindowItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
