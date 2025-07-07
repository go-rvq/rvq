package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDatePickerBuilder struct {
	VTagBuilder[*VDatePickerBuilder]
}

func VDatePicker(children ...h.HTMLComponent) *VDatePickerBuilder {
	return VTag(&VDatePickerBuilder{}, "v-date-picker", children...)
}

func (b *VDatePickerBuilder) Header(v string) (r *VDatePickerBuilder) {
	b.Attr("header", v)
	return b
}

func (b *VDatePickerBuilder) Title(v string) (r *VDatePickerBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VDatePickerBuilder) Active(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":active", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) Disabled(v bool) (r *VDatePickerBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VDatePickerBuilder) NextIcon(v string) (r *VDatePickerBuilder) {
	b.Attr("next-icon", v)
	return b
}

func (b *VDatePickerBuilder) PrevIcon(v string) (r *VDatePickerBuilder) {
	b.Attr("prev-icon", v)
	return b
}

func (b *VDatePickerBuilder) ModeIcon(v string) (r *VDatePickerBuilder) {
	b.Attr("mode-icon", v)
	return b
}

func (b *VDatePickerBuilder) Text(v string) (r *VDatePickerBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VDatePickerBuilder) ViewMode(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":view-mode", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) Month(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":month", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) Year(v int) (r *VDatePickerBuilder) {
	b.Attr(":year", fmt.Sprint(v))
	return b
}

func (b *VDatePickerBuilder) Color(v string) (r *VDatePickerBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDatePickerBuilder) HideWeekdays(v bool) (r *VDatePickerBuilder) {
	b.Attr(":hide-weekdays", fmt.Sprint(v))
	return b
}

func (b *VDatePickerBuilder) ShowWeek(v bool) (r *VDatePickerBuilder) {
	b.Attr(":show-week", fmt.Sprint(v))
	return b
}

func (b *VDatePickerBuilder) Transition(v string) (r *VDatePickerBuilder) {
	b.Attr("transition", v)
	return b
}

func (b *VDatePickerBuilder) ReverseTransition(v string) (r *VDatePickerBuilder) {
	b.Attr("reverse-transition", v)
	return b
}

func (b *VDatePickerBuilder) ShowAdjacentMonths(v bool) (r *VDatePickerBuilder) {
	b.Attr(":show-adjacent-months", fmt.Sprint(v))
	return b
}

func (b *VDatePickerBuilder) Weekdays(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":weekdays", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) WeeksInMonth(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":weeks-in-month", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) AllowedDates(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":allowed-dates", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) DisplayValue(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":display-value", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) ModelValue(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) Max(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) Min(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) Multiple(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":multiple", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) Height(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) BgColor(v string) (r *VDatePickerBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VDatePickerBuilder) Landscape(v bool) (r *VDatePickerBuilder) {
	b.Attr(":landscape", fmt.Sprint(v))
	return b
}

func (b *VDatePickerBuilder) HideHeader(v bool) (r *VDatePickerBuilder) {
	b.Attr(":hide-header", fmt.Sprint(v))
	return b
}

func (b *VDatePickerBuilder) Border(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) MaxHeight(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) MaxWidth(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) MinHeight(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) MinWidth(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) Width(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) Elevation(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) Location(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) Position(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) Rounded(v interface{}) (r *VDatePickerBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VDatePickerBuilder) Tile(v bool) (r *VDatePickerBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VDatePickerBuilder) Tag(v string) (r *VDatePickerBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VDatePickerBuilder) Theme(v string) (r *VDatePickerBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VDatePickerBuilder) On(name string, value string) (r *VDatePickerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDatePickerBuilder) Bind(name string, value string) (r *VDatePickerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
