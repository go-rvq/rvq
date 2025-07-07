package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VNumberInputBuilder struct {
	VTagBuilder[*VNumberInputBuilder]
}

func VNumberInput(children ...h.HTMLComponent) *VNumberInputBuilder {
	return VTag(&VNumberInputBuilder{}, "v-number-input", children...)
}

func (b *VNumberInputBuilder) Label(v string) (r *VNumberInputBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VNumberInputBuilder) Counter(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":counter", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Flat(v bool) (r *VNumberInputBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) ControlVariant(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":control-variant", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Inset(v bool) (r *VNumberInputBuilder) {
	b.Attr(":inset", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) HideInput(v bool) (r *VNumberInputBuilder) {
	b.Attr(":hide-input", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Min(v int) (r *VNumberInputBuilder) {
	b.Attr(":min", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Type(v string) (r *VNumberInputBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VNumberInputBuilder) Max(v int) (r *VNumberInputBuilder) {
	b.Attr(":max", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Step(v int) (r *VNumberInputBuilder) {
	b.Attr(":step", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Autofocus(v bool) (r *VNumberInputBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Prefix(v string) (r *VNumberInputBuilder) {
	b.Attr("prefix", v)
	return b
}

func (b *VNumberInputBuilder) Placeholder(v string) (r *VNumberInputBuilder) {
	b.Attr("placeholder", v)
	return b
}

func (b *VNumberInputBuilder) PersistentPlaceholder(v bool) (r *VNumberInputBuilder) {
	b.Attr(":persistent-placeholder", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) PersistentCounter(v bool) (r *VNumberInputBuilder) {
	b.Attr(":persistent-counter", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Suffix(v string) (r *VNumberInputBuilder) {
	b.Attr("suffix", v)
	return b
}

func (b *VNumberInputBuilder) Role(v string) (r *VNumberInputBuilder) {
	b.Attr("role", v)
	return b
}

func (b *VNumberInputBuilder) Id(v string) (r *VNumberInputBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VNumberInputBuilder) AppendIcon(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) CenterAffix(v bool) (r *VNumberInputBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) PrependIcon(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) HideSpinButtons(v bool) (r *VNumberInputBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Hint(v string) (r *VNumberInputBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VNumberInputBuilder) PersistentHint(v bool) (r *VNumberInputBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Messages(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Direction(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Reverse(v bool) (r *VNumberInputBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Density(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) MaxWidth(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) MinWidth(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Width(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Theme(v string) (r *VNumberInputBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VNumberInputBuilder) Disabled(v bool) (r *VNumberInputBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Error(v bool) (r *VNumberInputBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) ErrorMessages(v ...string) (r *VNumberInputBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) MaxErrors(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Name(v string) (r *VNumberInputBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VNumberInputBuilder) Readonly(v bool) (r *VNumberInputBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Rules(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) ModelValue(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) ValidateOn(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) ValidationValue(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Focused(v bool) (r *VNumberInputBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) HideDetails(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) BgColor(v string) (r *VNumberInputBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VNumberInputBuilder) Clearable(v bool) (r *VNumberInputBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) ClearIcon(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Active(v bool) (r *VNumberInputBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Color(v string) (r *VNumberInputBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VNumberInputBuilder) BaseColor(v string) (r *VNumberInputBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VNumberInputBuilder) Dirty(v bool) (r *VNumberInputBuilder) {
	b.Attr(":dirty", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) PersistentClear(v bool) (r *VNumberInputBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) SingleLine(v bool) (r *VNumberInputBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Variant(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Loading(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Rounded(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Tile(v bool) (r *VNumberInputBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) CounterValue(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":counter-value", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) ModelModifiers(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":model-modifiers", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) On(name string, value string) (r *VNumberInputBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VNumberInputBuilder) Bind(name string, value string) (r *VNumberInputBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
