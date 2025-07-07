package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VCalendarMonthDayBuilder struct {
	VTagBuilder[*VCalendarMonthDayBuilder]
}

func VCalendarMonthDay(children ...h.HTMLComponent) *VCalendarMonthDayBuilder {
	return VTag(&VCalendarMonthDayBuilder{}, "v-calendar-month-day", children...)
}

func (b *VCalendarMonthDayBuilder) Title(v interface{}) (r *VCalendarMonthDayBuilder) {
	b.Attr(":title", h.JSONString(v))
	return b
}

func (b *VCalendarMonthDayBuilder) Active(v bool) (r *VCalendarMonthDayBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VCalendarMonthDayBuilder) Color(v string) (r *VCalendarMonthDayBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VCalendarMonthDayBuilder) Day(v interface{}) (r *VCalendarMonthDayBuilder) {
	b.Attr(":day", h.JSONString(v))
	return b
}

func (b *VCalendarMonthDayBuilder) Disabled(v bool) (r *VCalendarMonthDayBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCalendarMonthDayBuilder) Events(v interface{}) (r *VCalendarMonthDayBuilder) {
	b.Attr(":events", h.JSONString(v))
	return b
}

func (b *VCalendarMonthDayBuilder) On(name string, value string) (r *VCalendarMonthDayBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCalendarMonthDayBuilder) Bind(name string, value string) (r *VCalendarMonthDayBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
