package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VWindowBuilder struct {
	VTagBuilder[*VWindowBuilder]
}

func VWindow(children ...h.HTMLComponent) *VWindowBuilder {
	return VTag(&VWindowBuilder{}, "v-window", children...)
}

func (b *VWindowBuilder) Continuous(v bool) (r *VWindowBuilder) {
	b.Attr(":continuous", fmt.Sprint(v))
	return b
}

func (b *VWindowBuilder) NextIcon(v interface{}) (r *VWindowBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) PrevIcon(v interface{}) (r *VWindowBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) Reverse(v bool) (r *VWindowBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VWindowBuilder) ShowArrows(v interface{}) (r *VWindowBuilder) {
	b.Attr(":show-arrows", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) Touch(v interface{}) (r *VWindowBuilder) {
	b.Attr(":touch", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) Direction(v interface{}) (r *VWindowBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) ModelValue(v interface{}) (r *VWindowBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) Disabled(v bool) (r *VWindowBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VWindowBuilder) SelectedClass(v string) (r *VWindowBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VWindowBuilder) Mandatory(v interface{}) (r *VWindowBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) Tag(v string) (r *VWindowBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VWindowBuilder) Theme(v string) (r *VWindowBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VWindowBuilder) On(name string, value string) (r *VWindowBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VWindowBuilder) Bind(name string, value string) (r *VWindowBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
