package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VStepperBuilder struct {
	VTagBuilder[*VStepperBuilder]
}

func VStepper(children ...h.HTMLComponent) *VStepperBuilder {
	return VTag(&VStepperBuilder{}, "v-stepper", children...)
}

func (b *VStepperBuilder) Flat(v bool) (r *VStepperBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) AltLabels(v bool) (r *VStepperBuilder) {
	b.Attr(":alt-labels", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) BgColor(v string) (r *VStepperBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VStepperBuilder) CompleteIcon(v string) (r *VStepperBuilder) {
	b.Attr("complete-icon", v)
	return b
}

func (b *VStepperBuilder) EditIcon(v string) (r *VStepperBuilder) {
	b.Attr("edit-icon", v)
	return b
}

func (b *VStepperBuilder) Editable(v bool) (r *VStepperBuilder) {
	b.Attr(":editable", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) ErrorIcon(v string) (r *VStepperBuilder) {
	b.Attr("error-icon", v)
	return b
}

func (b *VStepperBuilder) HideActions(v bool) (r *VStepperBuilder) {
	b.Attr(":hide-actions", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) Items(v interface{}) (r *VStepperBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) ItemTitle(v string) (r *VStepperBuilder) {
	b.Attr("item-title", v)
	return b
}

func (b *VStepperBuilder) ItemValue(v string) (r *VStepperBuilder) {
	b.Attr("item-value", v)
	return b
}

func (b *VStepperBuilder) NonLinear(v bool) (r *VStepperBuilder) {
	b.Attr(":non-linear", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) Mobile(v bool) (r *VStepperBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) MobileBreakpoint(v interface{}) (r *VStepperBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) ModelValue(v interface{}) (r *VStepperBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Multiple(v bool) (r *VStepperBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) Max(v int) (r *VStepperBuilder) {
	b.Attr(":max", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) SelectedClass(v string) (r *VStepperBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VStepperBuilder) Disabled(v bool) (r *VStepperBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) Mandatory(v interface{}) (r *VStepperBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Color(v string) (r *VStepperBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VStepperBuilder) Border(v interface{}) (r *VStepperBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Height(v interface{}) (r *VStepperBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) MaxHeight(v interface{}) (r *VStepperBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) MaxWidth(v interface{}) (r *VStepperBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) MinHeight(v interface{}) (r *VStepperBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) MinWidth(v interface{}) (r *VStepperBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Width(v interface{}) (r *VStepperBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Elevation(v interface{}) (r *VStepperBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Location(v interface{}) (r *VStepperBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Position(v interface{}) (r *VStepperBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Rounded(v interface{}) (r *VStepperBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VStepperBuilder) Tile(v bool) (r *VStepperBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VStepperBuilder) Tag(v string) (r *VStepperBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VStepperBuilder) Theme(v string) (r *VStepperBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VStepperBuilder) PrevText(v string) (r *VStepperBuilder) {
	b.Attr("prev-text", v)
	return b
}

func (b *VStepperBuilder) NextText(v string) (r *VStepperBuilder) {
	b.Attr("next-text", v)
	return b
}

func (b *VStepperBuilder) On(name string, value string) (r *VStepperBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperBuilder) Bind(name string, value string) (r *VStepperBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
