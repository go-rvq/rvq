package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSelectionControlBuilder struct {
	VTagBuilder[*VSelectionControlBuilder]
}

func VSelectionControl(children ...h.HTMLComponent) *VSelectionControlBuilder {
	return VTag(&VSelectionControlBuilder{}, "v-selection-control", children...)
}

func (b *VSelectionControlBuilder) Label(v string) (r *VSelectionControlBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VSelectionControlBuilder) BaseColor(v string) (r *VSelectionControlBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VSelectionControlBuilder) TrueValue(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":true-value", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) FalseValue(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":false-value", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) Value(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) Color(v string) (r *VSelectionControlBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSelectionControlBuilder) Disabled(v bool) (r *VSelectionControlBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlBuilder) DefaultsTarget(v string) (r *VSelectionControlBuilder) {
	b.Attr("defaults-target", v)
	return b
}

func (b *VSelectionControlBuilder) Error(v bool) (r *VSelectionControlBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlBuilder) Id(v string) (r *VSelectionControlBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VSelectionControlBuilder) Inline(v bool) (r *VSelectionControlBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlBuilder) FalseIcon(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":false-icon", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) TrueIcon(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":true-icon", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) Ripple(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) Multiple(v bool) (r *VSelectionControlBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlBuilder) Name(v string) (r *VSelectionControlBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VSelectionControlBuilder) Readonly(v bool) (r *VSelectionControlBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlBuilder) ModelValue(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) Type(v string) (r *VSelectionControlBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VSelectionControlBuilder) ValueComparator(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) Density(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) Theme(v string) (r *VSelectionControlBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSelectionControlBuilder) On(name string, value string) (r *VSelectionControlBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSelectionControlBuilder) Bind(name string, value string) (r *VSelectionControlBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
