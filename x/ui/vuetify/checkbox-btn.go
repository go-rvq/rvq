package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCheckboxBtnBuilder struct {
	VTagBuilder[*VCheckboxBtnBuilder]
}

func VCheckboxBtn(children ...h.HTMLComponent) *VCheckboxBtnBuilder {
	return VTag(&VCheckboxBtnBuilder{}, "v-checkbox-btn", children...)
}

func (b *VCheckboxBtnBuilder) Label(v string) (r *VCheckboxBtnBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VCheckboxBtnBuilder) Indeterminate(v bool) (r *VCheckboxBtnBuilder) {
	b.Attr(":indeterminate", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBtnBuilder) IndeterminateIcon(v interface{}) (r *VCheckboxBtnBuilder) {
	b.Attr(":indeterminate-icon", h.JSONString(v))
	return b
}

func (b *VCheckboxBtnBuilder) Type(v string) (r *VCheckboxBtnBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VCheckboxBtnBuilder) BaseColor(v string) (r *VCheckboxBtnBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VCheckboxBtnBuilder) TrueValue(v interface{}) (r *VCheckboxBtnBuilder) {
	b.Attr(":true-value", h.JSONString(v))
	return b
}

func (b *VCheckboxBtnBuilder) FalseValue(v interface{}) (r *VCheckboxBtnBuilder) {
	b.Attr(":false-value", h.JSONString(v))
	return b
}

func (b *VCheckboxBtnBuilder) Value(v interface{}) (r *VCheckboxBtnBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VCheckboxBtnBuilder) Color(v string) (r *VCheckboxBtnBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VCheckboxBtnBuilder) Disabled(v bool) (r *VCheckboxBtnBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBtnBuilder) DefaultsTarget(v string) (r *VCheckboxBtnBuilder) {
	b.Attr("defaults-target", v)
	return b
}

func (b *VCheckboxBtnBuilder) Error(v bool) (r *VCheckboxBtnBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBtnBuilder) Id(v string) (r *VCheckboxBtnBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VCheckboxBtnBuilder) Inline(v bool) (r *VCheckboxBtnBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBtnBuilder) FalseIcon(v interface{}) (r *VCheckboxBtnBuilder) {
	b.Attr(":false-icon", h.JSONString(v))
	return b
}

func (b *VCheckboxBtnBuilder) TrueIcon(v interface{}) (r *VCheckboxBtnBuilder) {
	b.Attr(":true-icon", h.JSONString(v))
	return b
}

func (b *VCheckboxBtnBuilder) Ripple(v interface{}) (r *VCheckboxBtnBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VCheckboxBtnBuilder) Multiple(v bool) (r *VCheckboxBtnBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBtnBuilder) Name(v string) (r *VCheckboxBtnBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VCheckboxBtnBuilder) Readonly(v bool) (r *VCheckboxBtnBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBtnBuilder) ModelValue(v interface{}) (r *VCheckboxBtnBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VCheckboxBtnBuilder) ValueComparator(v interface{}) (r *VCheckboxBtnBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VCheckboxBtnBuilder) Density(v interface{}) (r *VCheckboxBtnBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VCheckboxBtnBuilder) Theme(v string) (r *VCheckboxBtnBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VCheckboxBtnBuilder) On(name string, value string) (r *VCheckboxBtnBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCheckboxBtnBuilder) Bind(name string, value string) (r *VCheckboxBtnBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
