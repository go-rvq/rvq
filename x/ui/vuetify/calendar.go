package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCalendarBuilder struct {
	VTagBuilder[*VCalendarBuilder]
}

func VCalendar(children ...h.HTMLComponent) *VCalendarBuilder {
	return VTag(&VCalendarBuilder{}, "v-calendar", children...)
}

func (b *VCalendarBuilder) HideHeader(v bool) (r *VCalendarBuilder) {
	b.Attr(":hide-header", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) HideWeekNumber(v bool) (r *VCalendarBuilder) {
	b.Attr(":hide-week-number", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) Disabled(v bool) (r *VCalendarBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) Month(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":month", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) ShowAdjacentMonths(v bool) (r *VCalendarBuilder) {
	b.Attr(":show-adjacent-months", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) Year(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":year", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) Weekdays(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":weekdays", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) WeeksInMonth(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":weeks-in-month", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) AllowedDates(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":allowed-dates", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) DisplayValue(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":display-value", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) ModelValue(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) Max(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) Min(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) HideDayHeader(v bool) (r *VCalendarBuilder) {
	b.Attr(":hide-day-header", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) Intervals(v int) (r *VCalendarBuilder) {
	b.Attr(":intervals", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) Day(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":day", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) DayIndex(v int) (r *VCalendarBuilder) {
	b.Attr(":day-index", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) Events(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":events", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) IntervalDivisions(v int) (r *VCalendarBuilder) {
	b.Attr(":interval-divisions", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) IntervalDuration(v int) (r *VCalendarBuilder) {
	b.Attr(":interval-duration", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) IntervalHeight(v int) (r *VCalendarBuilder) {
	b.Attr(":interval-height", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) IntervalFormat(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":interval-format", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) IntervalStart(v int) (r *VCalendarBuilder) {
	b.Attr(":interval-start", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) NextIcon(v string) (r *VCalendarBuilder) {
	b.Attr("next-icon", v)
	return b
}

func (b *VCalendarBuilder) PrevIcon(v string) (r *VCalendarBuilder) {
	b.Attr("prev-icon", v)
	return b
}

func (b *VCalendarBuilder) Title(v string) (r *VCalendarBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VCalendarBuilder) Text(v string) (r *VCalendarBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VCalendarBuilder) ViewMode(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":view-mode", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) On(name string, value string) (r *VCalendarBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCalendarBuilder) Bind(name string, value string) (r *VCalendarBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
