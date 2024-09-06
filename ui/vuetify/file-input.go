package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VFileInputBuilder struct {
	VTagBuilder[*VFileInputBuilder]
}

func VFileInput(children ...h.HTMLComponent) *VFileInputBuilder {
	return VTag(&VFileInputBuilder{}, "v-file-input", children...)
}

func (b *VFileInputBuilder) Label(v string) (r *VFileInputBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VFileInputBuilder) Counter(v bool) (r *VFileInputBuilder) {
	b.Attr(":counter", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Flat(v bool) (r *VFileInputBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Chips(v bool) (r *VFileInputBuilder) {
	b.Attr(":chips", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) CounterSizeString(v string) (r *VFileInputBuilder) {
	b.Attr("counter-size-string", v)
	return b
}

func (b *VFileInputBuilder) CounterString(v string) (r *VFileInputBuilder) {
	b.Attr("counter-string", v)
	return b
}

func (b *VFileInputBuilder) HideInput(v bool) (r *VFileInputBuilder) {
	b.Attr(":hide-input", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Multiple(v bool) (r *VFileInputBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) ShowSize(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":show-size", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Id(v string) (r *VFileInputBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VFileInputBuilder) AppendIcon(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) CenterAffix(v bool) (r *VFileInputBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) PrependIcon(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) HideSpinButtons(v bool) (r *VFileInputBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Hint(v string) (r *VFileInputBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VFileInputBuilder) PersistentHint(v bool) (r *VFileInputBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Messages(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Direction(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Reverse(v bool) (r *VFileInputBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Density(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) MaxWidth(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) MinWidth(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Width(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Theme(v string) (r *VFileInputBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VFileInputBuilder) Disabled(v bool) (r *VFileInputBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Error(v bool) (r *VFileInputBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) MaxErrors(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Name(v string) (r *VFileInputBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VFileInputBuilder) Readonly(v bool) (r *VFileInputBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Rules(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) ModelValue(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) ValidateOn(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) ValidationValue(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Focused(v bool) (r *VFileInputBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) HideDetails(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) AppendInnerIcon(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":append-inner-icon", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) BgColor(v string) (r *VFileInputBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VFileInputBuilder) Clearable(v bool) (r *VFileInputBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) ClearIcon(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Active(v bool) (r *VFileInputBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Color(v string) (r *VFileInputBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VFileInputBuilder) BaseColor(v string) (r *VFileInputBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VFileInputBuilder) Dirty(v bool) (r *VFileInputBuilder) {
	b.Attr(":dirty", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) PersistentClear(v bool) (r *VFileInputBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) PrependInnerIcon(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":prepend-inner-icon", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) SingleLine(v bool) (r *VFileInputBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Variant(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Loading(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Rounded(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Tile(v bool) (r *VFileInputBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) On(name string, value string) (r *VFileInputBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFileInputBuilder) Bind(name string, value string) (r *VFileInputBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
