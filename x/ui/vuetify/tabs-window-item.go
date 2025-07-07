package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VTabsWindowItemBuilder struct {
	VTagBuilder[*VTabsWindowItemBuilder]
}

func VTabsWindowItem(children ...h.HTMLComponent) *VTabsWindowItemBuilder {
	return VTag(&VTabsWindowItemBuilder{}, "v-tabs-window-item", children...)
}

func (b *VTabsWindowItemBuilder) ReverseTransition(v interface{}) (r *VTabsWindowItemBuilder) {
	b.Attr(":reverse-transition", h.JSONString(v))
	return b
}

func (b *VTabsWindowItemBuilder) Transition(v interface{}) (r *VTabsWindowItemBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VTabsWindowItemBuilder) Value(v interface{}) (r *VTabsWindowItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VTabsWindowItemBuilder) Disabled(v bool) (r *VTabsWindowItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTabsWindowItemBuilder) SelectedClass(v string) (r *VTabsWindowItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VTabsWindowItemBuilder) Eager(v bool) (r *VTabsWindowItemBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VTabsWindowItemBuilder) On(name string, value string) (r *VTabsWindowItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTabsWindowItemBuilder) Bind(name string, value string) (r *VTabsWindowItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
