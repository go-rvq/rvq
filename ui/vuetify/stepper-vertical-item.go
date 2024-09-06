package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VStepperVerticalItemBuilder struct {
	VTagBuilder[*VStepperVerticalItemBuilder]
}

func VStepperVerticalItem(children ...h.HTMLComponent) *VStepperVerticalItemBuilder {
	return VTag(&VStepperVerticalItemBuilder{}, "v-stepper-vertical-item", children...)
}

func (b *VStepperVerticalItemBuilder) Icon(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("icon", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Subtitle(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("subtitle", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Title(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Text(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VStepperVerticalItemBuilder) HideActions(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":hide-actions", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Color(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Complete(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":complete", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) CompleteIcon(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("complete-icon", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Editable(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":editable", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) EditIcon(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("edit-icon", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Error(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) ErrorIcon(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("error-icon", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Ripple(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Value(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Rules(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) BgColor(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Elevation(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Disabled(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) SelectedClass(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VStepperVerticalItemBuilder) Rounded(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Tile(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Tag(v string) (r *VStepperVerticalItemBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VStepperVerticalItemBuilder) ExpandIcon(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) CollapseIcon(v interface{}) (r *VStepperVerticalItemBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Focusable(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":focusable", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Static(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":static", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Readonly(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) Eager(v bool) (r *VStepperVerticalItemBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VStepperVerticalItemBuilder) On(name string, value string) (r *VStepperVerticalItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperVerticalItemBuilder) Bind(name string, value string) (r *VStepperVerticalItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
