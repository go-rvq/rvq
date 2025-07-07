package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCalendarIntervalBuilder struct {
	VTagBuilder[*VCalendarIntervalBuilder]
}

func VCalendarInterval(children ...h.HTMLComponent) *VCalendarIntervalBuilder {
	return VTag(&VCalendarIntervalBuilder{}, "v-calendar-interval", children...)
}

func (b *VCalendarIntervalBuilder) Index(v int) (r *VCalendarIntervalBuilder) {
	b.Attr(":index", fmt.Sprint(v))
	return b
}

func (b *VCalendarIntervalBuilder) Day(v interface{}) (r *VCalendarIntervalBuilder) {
	b.Attr(":day", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalBuilder) DayIndex(v int) (r *VCalendarIntervalBuilder) {
	b.Attr(":day-index", fmt.Sprint(v))
	return b
}

func (b *VCalendarIntervalBuilder) Events(v interface{}) (r *VCalendarIntervalBuilder) {
	b.Attr(":events", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalBuilder) IntervalDivisions(v int) (r *VCalendarIntervalBuilder) {
	b.Attr(":interval-divisions", fmt.Sprint(v))
	return b
}

func (b *VCalendarIntervalBuilder) IntervalDuration(v int) (r *VCalendarIntervalBuilder) {
	b.Attr(":interval-duration", fmt.Sprint(v))
	return b
}

func (b *VCalendarIntervalBuilder) IntervalHeight(v int) (r *VCalendarIntervalBuilder) {
	b.Attr(":interval-height", fmt.Sprint(v))
	return b
}

func (b *VCalendarIntervalBuilder) IntervalFormat(v interface{}) (r *VCalendarIntervalBuilder) {
	b.Attr(":interval-format", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalBuilder) IntervalStart(v int) (r *VCalendarIntervalBuilder) {
	b.Attr(":interval-start", fmt.Sprint(v))
	return b
}

func (b *VCalendarIntervalBuilder) On(name string, value string) (r *VCalendarIntervalBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCalendarIntervalBuilder) Bind(name string, value string) (r *VCalendarIntervalBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
