package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VStepperItemBuilder struct {
	VTagBuilder[*VStepperItemBuilder]
}

func VStepperItem(children ...h.HTMLComponent) *VStepperItemBuilder {
	return VTag(&VStepperItemBuilder{}, "v-stepper-item", children...)
}

func (b *VStepperItemBuilder) Icon(v string) (r *VStepperItemBuilder) {
	b.Attr("icon", v)
	return b
}

func (b *VStepperItemBuilder) Title(v string) (r *VStepperItemBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VStepperItemBuilder) Subtitle(v string) (r *VStepperItemBuilder) {
	b.Attr("subtitle", v)
	return b
}

func (b *VStepperItemBuilder) Color(v string) (r *VStepperItemBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VStepperItemBuilder) Complete(v bool) (r *VStepperItemBuilder) {
	b.Attr(":complete", fmt.Sprint(v))
	return b
}

func (b *VStepperItemBuilder) CompleteIcon(v string) (r *VStepperItemBuilder) {
	b.Attr("complete-icon", v)
	return b
}

func (b *VStepperItemBuilder) Editable(v bool) (r *VStepperItemBuilder) {
	b.Attr(":editable", fmt.Sprint(v))
	return b
}

func (b *VStepperItemBuilder) EditIcon(v string) (r *VStepperItemBuilder) {
	b.Attr("edit-icon", v)
	return b
}

func (b *VStepperItemBuilder) Error(v bool) (r *VStepperItemBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VStepperItemBuilder) ErrorIcon(v string) (r *VStepperItemBuilder) {
	b.Attr("error-icon", v)
	return b
}

func (b *VStepperItemBuilder) Ripple(v interface{}) (r *VStepperItemBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VStepperItemBuilder) Value(v interface{}) (r *VStepperItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VStepperItemBuilder) Rules(v interface{}) (r *VStepperItemBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VStepperItemBuilder) Disabled(v bool) (r *VStepperItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VStepperItemBuilder) SelectedClass(v string) (r *VStepperItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VStepperItemBuilder) On(name string, value string) (r *VStepperItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperItemBuilder) Bind(name string, value string) (r *VStepperItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
