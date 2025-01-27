package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VTextFieldBuilder struct {
	VTagBuilder[*VTextFieldBuilder]
}

func VTextField(children ...h.HTMLComponent) *VTextFieldBuilder {
	return VTag(&VTextFieldBuilder{}, "v-text-field", children...)
}

func (b *VTextFieldBuilder) Label(v string) (r *VTextFieldBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VTextFieldBuilder) Counter(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":counter", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Flat(v bool) (r *VTextFieldBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Autofocus(v bool) (r *VTextFieldBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) MaxLength(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":maxlength", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Prefix(v string) (r *VTextFieldBuilder) {
	b.Attr("prefix", v)
	return b
}

func (b *VTextFieldBuilder) Placeholder(v string) (r *VTextFieldBuilder) {
	b.Attr("placeholder", v)
	return b
}

func (b *VTextFieldBuilder) PersistentPlaceholder(v bool) (r *VTextFieldBuilder) {
	b.Attr(":persistent-placeholder", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) PersistentCounter(v bool) (r *VTextFieldBuilder) {
	b.Attr(":persistent-counter", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Suffix(v string) (r *VTextFieldBuilder) {
	b.Attr("suffix", v)
	return b
}

func (b *VTextFieldBuilder) Role(v string) (r *VTextFieldBuilder) {
	b.Attr("role", v)
	return b
}

func (b *VTextFieldBuilder) Type(v string) (r *VTextFieldBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VTextFieldBuilder) Id(v string) (r *VTextFieldBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VTextFieldBuilder) AppendIcon(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) CenterAffix(v bool) (r *VTextFieldBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) PrependIcon(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) HideSpinButtons(v bool) (r *VTextFieldBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Hint(v string) (r *VTextFieldBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VTextFieldBuilder) PersistentHint(v bool) (r *VTextFieldBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Messages(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Direction(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Reverse(v bool) (r *VTextFieldBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Density(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) MaxWidth(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) MinWidth(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Width(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Theme(v string) (r *VTextFieldBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTextFieldBuilder) Disabled(v bool) (r *VTextFieldBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Error(v bool) (r *VTextFieldBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) MaxErrors(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Name(v string) (r *VTextFieldBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VTextFieldBuilder) Readonly(v bool) (r *VTextFieldBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Rules(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) ModelValue(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) ValidateOn(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) ValidationValue(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Focused(v bool) (r *VTextFieldBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) HideDetails(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) AppendInnerIcon(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":append-inner-icon", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) BgColor(v string) (r *VTextFieldBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VTextFieldBuilder) Clearable(v bool) (r *VTextFieldBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) ClearIcon(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Active(v bool) (r *VTextFieldBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Color(v string) (r *VTextFieldBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTextFieldBuilder) BaseColor(v string) (r *VTextFieldBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VTextFieldBuilder) Dirty(v bool) (r *VTextFieldBuilder) {
	b.Attr(":dirty", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) PersistentClear(v bool) (r *VTextFieldBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) PrependInnerIcon(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":prepend-inner-icon", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) SingleLine(v bool) (r *VTextFieldBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Variant(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Loading(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Rounded(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Tile(v bool) (r *VTextFieldBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) CounterValue(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":counter-value", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) ModelModifiers(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":model-modifiers", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) On(name string, value string) (r *VTextFieldBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTextFieldBuilder) Bind(name string, value string) (r *VTextFieldBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VTextFieldBuilder) Event(name string, value string) (r *VTextFieldBuilder) {
	b.Attr("@"+name, value)
	return b
}
