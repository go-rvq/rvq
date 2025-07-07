package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTextareaBuilder struct {
	VTagBuilder[*VTextareaBuilder]
}

func VTextarea(children ...h.HTMLComponent) *VTextareaBuilder {
	return VTag(&VTextareaBuilder{}, "v-textarea", children...)
}

func (b *VTextareaBuilder) Label(v string) (r *VTextareaBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VTextareaBuilder) Counter(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":counter", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Flat(v bool) (r *VTextareaBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) AutoGrow(v bool) (r *VTextareaBuilder) {
	b.Attr(":auto-grow", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Autofocus(v bool) (r *VTextareaBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Prefix(v string) (r *VTextareaBuilder) {
	b.Attr("prefix", v)
	return b
}

func (b *VTextareaBuilder) Placeholder(v string) (r *VTextareaBuilder) {
	b.Attr("placeholder", v)
	return b
}

func (b *VTextareaBuilder) PersistentPlaceholder(v bool) (r *VTextareaBuilder) {
	b.Attr(":persistent-placeholder", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) PersistentCounter(v bool) (r *VTextareaBuilder) {
	b.Attr(":persistent-counter", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) NoResize(v bool) (r *VTextareaBuilder) {
	b.Attr(":no-resize", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Rows(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":rows", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) MaxRows(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":max-rows", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Suffix(v string) (r *VTextareaBuilder) {
	b.Attr("suffix", v)
	return b
}

func (b *VTextareaBuilder) Id(v string) (r *VTextareaBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VTextareaBuilder) AppendIcon(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) CenterAffix(v bool) (r *VTextareaBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) PrependIcon(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) HideSpinButtons(v bool) (r *VTextareaBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Hint(v string) (r *VTextareaBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VTextareaBuilder) PersistentHint(v bool) (r *VTextareaBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Messages(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Direction(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Reverse(v bool) (r *VTextareaBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Density(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) MaxWidth(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) MinWidth(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Width(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Theme(v string) (r *VTextareaBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTextareaBuilder) Disabled(v bool) (r *VTextareaBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Error(v bool) (r *VTextareaBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) MaxErrors(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Name(v string) (r *VTextareaBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VTextareaBuilder) Readonly(v bool) (r *VTextareaBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Rules(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) ModelValue(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) ValidateOn(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) ValidationValue(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Focused(v bool) (r *VTextareaBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) HideDetails(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) AppendInnerIcon(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":append-inner-icon", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) BgColor(v string) (r *VTextareaBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VTextareaBuilder) Clearable(v bool) (r *VTextareaBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) ClearIcon(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Active(v bool) (r *VTextareaBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Color(v string) (r *VTextareaBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTextareaBuilder) BaseColor(v string) (r *VTextareaBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VTextareaBuilder) Dirty(v bool) (r *VTextareaBuilder) {
	b.Attr(":dirty", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) PersistentClear(v bool) (r *VTextareaBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) PrependInnerIcon(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":prepend-inner-icon", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) SingleLine(v bool) (r *VTextareaBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Variant(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Loading(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Rounded(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Tile(v bool) (r *VTextareaBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) CounterValue(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":counter-value", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) ModelModifiers(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":model-modifiers", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) On(name string, value string) (r *VTextareaBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTextareaBuilder) Bind(name string, value string) (r *VTextareaBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
