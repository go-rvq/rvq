package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSelectionControlGroupBuilder struct {
	VTagBuilder[*VSelectionControlGroupBuilder]
}

func VSelectionControlGroup(children ...h.HTMLComponent) *VSelectionControlGroupBuilder {
	return VTag(&VSelectionControlGroupBuilder{}, "v-selection-control-group", children...)
}

func (b *VSelectionControlGroupBuilder) Color(v string) (r *VSelectionControlGroupBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSelectionControlGroupBuilder) Disabled(v bool) (r *VSelectionControlGroupBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlGroupBuilder) DefaultsTarget(v string) (r *VSelectionControlGroupBuilder) {
	b.Attr("defaults-target", v)
	return b
}

func (b *VSelectionControlGroupBuilder) Error(v bool) (r *VSelectionControlGroupBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlGroupBuilder) Id(v string) (r *VSelectionControlGroupBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VSelectionControlGroupBuilder) Inline(v bool) (r *VSelectionControlGroupBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlGroupBuilder) FalseIcon(v interface{}) (r *VSelectionControlGroupBuilder) {
	b.Attr(":false-icon", h.JSONString(v))
	return b
}

func (b *VSelectionControlGroupBuilder) TrueIcon(v interface{}) (r *VSelectionControlGroupBuilder) {
	b.Attr(":true-icon", h.JSONString(v))
	return b
}

func (b *VSelectionControlGroupBuilder) Ripple(v interface{}) (r *VSelectionControlGroupBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VSelectionControlGroupBuilder) Multiple(v bool) (r *VSelectionControlGroupBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlGroupBuilder) Name(v string) (r *VSelectionControlGroupBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VSelectionControlGroupBuilder) Readonly(v bool) (r *VSelectionControlGroupBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlGroupBuilder) ModelValue(v interface{}) (r *VSelectionControlGroupBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VSelectionControlGroupBuilder) Type(v string) (r *VSelectionControlGroupBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VSelectionControlGroupBuilder) ValueComparator(v interface{}) (r *VSelectionControlGroupBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VSelectionControlGroupBuilder) Density(v interface{}) (r *VSelectionControlGroupBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VSelectionControlGroupBuilder) Theme(v string) (r *VSelectionControlGroupBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSelectionControlGroupBuilder) On(name string, value string) (r *VSelectionControlGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSelectionControlGroupBuilder) Bind(name string, value string) (r *VSelectionControlGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
