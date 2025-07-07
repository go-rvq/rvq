package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VStepperVerticalBuilder struct {
	VTagBuilder[*VStepperVerticalBuilder]
}

func VStepperVertical(children ...h.HTMLComponent) *VStepperVerticalBuilder {
	return VTag(&VStepperVerticalBuilder{}, "v-stepper-vertical", children...)
}

func (b *VStepperVerticalBuilder) Title(v string) (r *VStepperVerticalBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VStepperVerticalBuilder) Flat(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) PrevText(v string) (r *VStepperVerticalBuilder) {
	b.Attr("prev-text", v)
	return b
}

func (b *VStepperVerticalBuilder) NextText(v string) (r *VStepperVerticalBuilder) {
	b.Attr("next-text", v)
	return b
}

func (b *VStepperVerticalBuilder) AltLabels(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":alt-labels", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) BgColor(v string) (r *VStepperVerticalBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VStepperVerticalBuilder) CompleteIcon(v string) (r *VStepperVerticalBuilder) {
	b.Attr("complete-icon", v)
	return b
}

func (b *VStepperVerticalBuilder) EditIcon(v string) (r *VStepperVerticalBuilder) {
	b.Attr("edit-icon", v)
	return b
}

func (b *VStepperVerticalBuilder) Editable(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":editable", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) ErrorIcon(v string) (r *VStepperVerticalBuilder) {
	b.Attr("error-icon", v)
	return b
}

func (b *VStepperVerticalBuilder) HideActions(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":hide-actions", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) Items(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) ItemTitle(v string) (r *VStepperVerticalBuilder) {
	b.Attr("item-title", v)
	return b
}

func (b *VStepperVerticalBuilder) ItemValue(v string) (r *VStepperVerticalBuilder) {
	b.Attr("item-value", v)
	return b
}

func (b *VStepperVerticalBuilder) Value(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) NonLinear(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":non-linear", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) Mobile(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) MobileBreakpoint(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) ModelValue(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Multiple(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) Max(v int) (r *VStepperVerticalBuilder) {
	b.Attr(":max", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) SelectedClass(v string) (r *VStepperVerticalBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VStepperVerticalBuilder) Disabled(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) Mandatory(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Text(v string) (r *VStepperVerticalBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VStepperVerticalBuilder) Elevation(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Rounded(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Tile(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) Tag(v string) (r *VStepperVerticalBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VStepperVerticalBuilder) Color(v string) (r *VStepperVerticalBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VStepperVerticalBuilder) ExpandIcon(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) CollapseIcon(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Focusable(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":focusable", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) Ripple(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) Readonly(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) Eager(v bool) (r *VStepperVerticalBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalBuilder) Theme(v string) (r *VStepperVerticalBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VStepperVerticalBuilder) Variant(v interface{}) (r *VStepperVerticalBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VStepperVerticalBuilder) On(name string, value string) (r *VStepperVerticalBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperVerticalBuilder) Bind(name string, value string) (r *VStepperVerticalBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
