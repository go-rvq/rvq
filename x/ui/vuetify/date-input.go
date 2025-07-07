package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDateInputBuilder struct {
	VTagBuilder[*VDateInputBuilder]
}

func VDateInput(children ...h.HTMLComponent) *VDateInputBuilder {
	return VTag(&VDateInputBuilder{}, "v-date-input", children...)
}

func (b *VDateInputBuilder) Flat(v bool) (r *VDateInputBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) HideActions(v bool) (r *VDateInputBuilder) {
	b.Attr(":hide-actions", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Focused(v bool) (r *VDateInputBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Reverse(v bool) (r *VDateInputBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) ModelValue(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Color(v string) (r *VDateInputBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDateInputBuilder) CancelText(v string) (r *VDateInputBuilder) {
	b.Attr("cancel-text", v)
	return b
}

func (b *VDateInputBuilder) Type(v string) (r *VDateInputBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VDateInputBuilder) OkText(v string) (r *VDateInputBuilder) {
	b.Attr("ok-text", v)
	return b
}

func (b *VDateInputBuilder) Autofocus(v bool) (r *VDateInputBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Counter(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":counter", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Prefix(v string) (r *VDateInputBuilder) {
	b.Attr("prefix", v)
	return b
}

func (b *VDateInputBuilder) Placeholder(v string) (r *VDateInputBuilder) {
	b.Attr("placeholder", v)
	return b
}

func (b *VDateInputBuilder) PersistentPlaceholder(v bool) (r *VDateInputBuilder) {
	b.Attr(":persistent-placeholder", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) PersistentCounter(v bool) (r *VDateInputBuilder) {
	b.Attr(":persistent-counter", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Suffix(v string) (r *VDateInputBuilder) {
	b.Attr("suffix", v)
	return b
}

func (b *VDateInputBuilder) Role(v string) (r *VDateInputBuilder) {
	b.Attr("role", v)
	return b
}

func (b *VDateInputBuilder) Text(v string) (r *VDateInputBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VDateInputBuilder) Id(v string) (r *VDateInputBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VDateInputBuilder) AppendIcon(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) CenterAffix(v bool) (r *VDateInputBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) PrependIcon(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) HideSpinButtons(v bool) (r *VDateInputBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Hint(v string) (r *VDateInputBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VDateInputBuilder) PersistentHint(v bool) (r *VDateInputBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Messages(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Direction(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Density(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Height(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) MaxHeight(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) MaxWidth(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) MinHeight(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) MinWidth(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Width(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Theme(v string) (r *VDateInputBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VDateInputBuilder) Disabled(v bool) (r *VDateInputBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Error(v bool) (r *VDateInputBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) ErrorMessages(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) MaxErrors(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Name(v string) (r *VDateInputBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VDateInputBuilder) Label(v string) (r *VDateInputBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VDateInputBuilder) Readonly(v bool) (r *VDateInputBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Rules(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) ValidateOn(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) ValidationValue(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) HideDetails(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) AppendInnerIcon(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":append-inner-icon", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) BgColor(v string) (r *VDateInputBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VDateInputBuilder) Clearable(v bool) (r *VDateInputBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) ClearIcon(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Active(v bool) (r *VDateInputBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) BaseColor(v string) (r *VDateInputBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VDateInputBuilder) Dirty(v bool) (r *VDateInputBuilder) {
	b.Attr(":dirty", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) PersistentClear(v bool) (r *VDateInputBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) PrependInnerIcon(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":prepend-inner-icon", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) SingleLine(v bool) (r *VDateInputBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Variant(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Loading(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Rounded(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Tile(v bool) (r *VDateInputBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) CounterValue(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":counter-value", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) ModelModifiers(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":model-modifiers", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Header(v string) (r *VDateInputBuilder) {
	b.Attr("header", v)
	return b
}

func (b *VDateInputBuilder) NextIcon(v string) (r *VDateInputBuilder) {
	b.Attr("next-icon", v)
	return b
}

func (b *VDateInputBuilder) PrevIcon(v string) (r *VDateInputBuilder) {
	b.Attr("prev-icon", v)
	return b
}

func (b *VDateInputBuilder) ModeIcon(v string) (r *VDateInputBuilder) {
	b.Attr("mode-icon", v)
	return b
}

func (b *VDateInputBuilder) ViewMode(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":view-mode", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Month(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":month", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Year(v int) (r *VDateInputBuilder) {
	b.Attr(":year", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) HideWeekdays(v bool) (r *VDateInputBuilder) {
	b.Attr(":hide-weekdays", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) ShowWeek(v bool) (r *VDateInputBuilder) {
	b.Attr(":show-week", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Transition(v string) (r *VDateInputBuilder) {
	b.Attr("transition", v)
	return b
}

func (b *VDateInputBuilder) ReverseTransition(v string) (r *VDateInputBuilder) {
	b.Attr("reverse-transition", v)
	return b
}

func (b *VDateInputBuilder) ShowAdjacentMonths(v bool) (r *VDateInputBuilder) {
	b.Attr(":show-adjacent-months", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Weekdays(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":weekdays", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) WeeksInMonth(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":weeks-in-month", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) AllowedDates(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":allowed-dates", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) DisplayValue(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":display-value", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Max(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Min(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Multiple(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":multiple", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Landscape(v bool) (r *VDateInputBuilder) {
	b.Attr(":landscape", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Title(v string) (r *VDateInputBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VDateInputBuilder) HideHeader(v bool) (r *VDateInputBuilder) {
	b.Attr(":hide-header", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Border(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Elevation(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Location(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Position(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Tag(v string) (r *VDateInputBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VDateInputBuilder) On(name string, value string) (r *VDateInputBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDateInputBuilder) Bind(name string, value string) (r *VDateInputBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
