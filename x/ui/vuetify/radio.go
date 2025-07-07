package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VRadioBuilder struct {
	VTagBuilder[*VRadioBuilder]
}

func VRadio(children ...h.HTMLComponent) *VRadioBuilder {
	return VTag(&VRadioBuilder{}, "v-radio", children...)
}

func (b *VRadioBuilder) Label(v string) (r *VRadioBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VRadioBuilder) BaseColor(v string) (r *VRadioBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VRadioBuilder) TrueValue(v interface{}) (r *VRadioBuilder) {
	b.Attr(":true-value", h.JSONString(v))
	return b
}

func (b *VRadioBuilder) FalseValue(v interface{}) (r *VRadioBuilder) {
	b.Attr(":false-value", h.JSONString(v))
	return b
}

func (b *VRadioBuilder) Value(v interface{}) (r *VRadioBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VRadioBuilder) Color(v string) (r *VRadioBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VRadioBuilder) Disabled(v bool) (r *VRadioBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VRadioBuilder) DefaultsTarget(v string) (r *VRadioBuilder) {
	b.Attr("defaults-target", v)
	return b
}

func (b *VRadioBuilder) Error(v bool) (r *VRadioBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VRadioBuilder) Id(v string) (r *VRadioBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VRadioBuilder) Inline(v bool) (r *VRadioBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VRadioBuilder) FalseIcon(v interface{}) (r *VRadioBuilder) {
	b.Attr(":false-icon", h.JSONString(v))
	return b
}

func (b *VRadioBuilder) TrueIcon(v interface{}) (r *VRadioBuilder) {
	b.Attr(":true-icon", h.JSONString(v))
	return b
}

func (b *VRadioBuilder) Ripple(v interface{}) (r *VRadioBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VRadioBuilder) Multiple(v bool) (r *VRadioBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VRadioBuilder) Name(v string) (r *VRadioBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VRadioBuilder) Readonly(v bool) (r *VRadioBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VRadioBuilder) ModelValue(v interface{}) (r *VRadioBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VRadioBuilder) Type(v string) (r *VRadioBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VRadioBuilder) ValueComparator(v interface{}) (r *VRadioBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VRadioBuilder) Density(v interface{}) (r *VRadioBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VRadioBuilder) Theme(v string) (r *VRadioBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VRadioBuilder) On(name string, value string) (r *VRadioBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VRadioBuilder) Bind(name string, value string) (r *VRadioBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
