package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCalendarDayBuilder struct {
	VTagBuilder[*VCalendarDayBuilder]
}

func VCalendarDay(children ...h.HTMLComponent) *VCalendarDayBuilder {
	return VTag(&VCalendarDayBuilder{}, "v-calendar-day", children...)
}

func (b *VCalendarDayBuilder) HideDayHeader(v bool) (r *VCalendarDayBuilder) {
	b.Attr(":hide-day-header", fmt.Sprint(v))
	return b
}

func (b *VCalendarDayBuilder) Intervals(v int) (r *VCalendarDayBuilder) {
	b.Attr(":intervals", fmt.Sprint(v))
	return b
}

func (b *VCalendarDayBuilder) Day(v interface{}) (r *VCalendarDayBuilder) {
	b.Attr(":day", h.JSONString(v))
	return b
}

func (b *VCalendarDayBuilder) DayIndex(v int) (r *VCalendarDayBuilder) {
	b.Attr(":day-index", fmt.Sprint(v))
	return b
}

func (b *VCalendarDayBuilder) Events(v interface{}) (r *VCalendarDayBuilder) {
	b.Attr(":events", h.JSONString(v))
	return b
}

func (b *VCalendarDayBuilder) IntervalDivisions(v int) (r *VCalendarDayBuilder) {
	b.Attr(":interval-divisions", fmt.Sprint(v))
	return b
}

func (b *VCalendarDayBuilder) IntervalDuration(v int) (r *VCalendarDayBuilder) {
	b.Attr(":interval-duration", fmt.Sprint(v))
	return b
}

func (b *VCalendarDayBuilder) IntervalHeight(v int) (r *VCalendarDayBuilder) {
	b.Attr(":interval-height", fmt.Sprint(v))
	return b
}

func (b *VCalendarDayBuilder) IntervalFormat(v interface{}) (r *VCalendarDayBuilder) {
	b.Attr(":interval-format", h.JSONString(v))
	return b
}

func (b *VCalendarDayBuilder) IntervalStart(v int) (r *VCalendarDayBuilder) {
	b.Attr(":interval-start", fmt.Sprint(v))
	return b
}

func (b *VCalendarDayBuilder) On(name string, value string) (r *VCalendarDayBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCalendarDayBuilder) Bind(name string, value string) (r *VCalendarDayBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
