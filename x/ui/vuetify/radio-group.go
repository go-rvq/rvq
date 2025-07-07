package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VRadioGroupBuilder struct {
	VTagBuilder[*VRadioGroupBuilder]
}

func VRadioGroup(children ...h.HTMLComponent) *VRadioGroupBuilder {
	return VTag(&VRadioGroupBuilder{}, "v-radio-group", children...)
}

func (b *VRadioGroupBuilder) Label(v string) (r *VRadioGroupBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VRadioGroupBuilder) Height(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Type(v string) (r *VRadioGroupBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VRadioGroupBuilder) Id(v string) (r *VRadioGroupBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VRadioGroupBuilder) AppendIcon(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) CenterAffix(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) PrependIcon(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) HideSpinButtons(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) Hint(v string) (r *VRadioGroupBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VRadioGroupBuilder) PersistentHint(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) Messages(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Direction(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Density(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) MaxWidth(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) MinWidth(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Width(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Theme(v string) (r *VRadioGroupBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VRadioGroupBuilder) Disabled(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) Error(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) MaxErrors(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Name(v string) (r *VRadioGroupBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VRadioGroupBuilder) Readonly(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) Rules(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) ModelValue(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) ValidateOn(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) ValidationValue(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Focused(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) HideDetails(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Color(v string) (r *VRadioGroupBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VRadioGroupBuilder) DefaultsTarget(v string) (r *VRadioGroupBuilder) {
	b.Attr("defaults-target", v)
	return b
}

func (b *VRadioGroupBuilder) Inline(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) FalseIcon(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":false-icon", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) TrueIcon(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":true-icon", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Ripple(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) ValueComparator(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) On(name string, value string) (r *VRadioGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VRadioGroupBuilder) Bind(name string, value string) (r *VRadioGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
