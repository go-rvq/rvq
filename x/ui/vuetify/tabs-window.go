package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VTabsWindowBuilder struct {
	VTagBuilder[*VTabsWindowBuilder]
}

func VTabsWindow(children ...h.HTMLComponent) *VTabsWindowBuilder {
	return VTag(&VTabsWindowBuilder{}, "v-tabs-window", children...)
}

func (b *VTabsWindowBuilder) Reverse(v bool) (r *VTabsWindowBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VTabsWindowBuilder) Direction(v interface{}) (r *VTabsWindowBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VTabsWindowBuilder) ModelValue(v interface{}) (r *VTabsWindowBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VTabsWindowBuilder) Disabled(v bool) (r *VTabsWindowBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTabsWindowBuilder) SelectedClass(v string) (r *VTabsWindowBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VTabsWindowBuilder) Tag(v string) (r *VTabsWindowBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VTabsWindowBuilder) Theme(v string) (r *VTabsWindowBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTabsWindowBuilder) On(name string, value string) (r *VTabsWindowBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTabsWindowBuilder) Bind(name string, value string) (r *VTabsWindowBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
