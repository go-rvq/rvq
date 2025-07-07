package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSwitchBuilder struct {
	VTagBuilder[*VSwitchBuilder]
}

func VSwitch(children ...h.HTMLComponent) *VSwitchBuilder {
	return VTag(&VSwitchBuilder{}, "v-switch", children...)
}

func (b *VSwitchBuilder) Label(v string) (r *VSwitchBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VSwitchBuilder) Indeterminate(v bool) (r *VSwitchBuilder) {
	b.Attr(":indeterminate", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Inset(v bool) (r *VSwitchBuilder) {
	b.Attr(":inset", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Flat(v bool) (r *VSwitchBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Loading(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Type(v string) (r *VSwitchBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VSwitchBuilder) Id(v string) (r *VSwitchBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VSwitchBuilder) AppendIcon(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) CenterAffix(v bool) (r *VSwitchBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) PrependIcon(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) HideSpinButtons(v bool) (r *VSwitchBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Hint(v string) (r *VSwitchBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VSwitchBuilder) PersistentHint(v bool) (r *VSwitchBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Messages(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Direction(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Density(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) MaxWidth(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) MinWidth(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Width(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Theme(v string) (r *VSwitchBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSwitchBuilder) Disabled(v bool) (r *VSwitchBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Error(v bool) (r *VSwitchBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) MaxErrors(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Name(v string) (r *VSwitchBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VSwitchBuilder) Readonly(v bool) (r *VSwitchBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Rules(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) ModelValue(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) ValidateOn(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) ValidationValue(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Focused(v bool) (r *VSwitchBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) HideDetails(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) BaseColor(v string) (r *VSwitchBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VSwitchBuilder) TrueValue(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":true-value", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) FalseValue(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":false-value", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Value(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Color(v string) (r *VSwitchBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSwitchBuilder) DefaultsTarget(v string) (r *VSwitchBuilder) {
	b.Attr("defaults-target", v)
	return b
}

func (b *VSwitchBuilder) Inline(v bool) (r *VSwitchBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) FalseIcon(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":false-icon", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) TrueIcon(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":true-icon", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Ripple(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Multiple(v bool) (r *VSwitchBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) ValueComparator(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) On(name string, value string) (r *VSwitchBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSwitchBuilder) Bind(name string, value string) (r *VSwitchBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
