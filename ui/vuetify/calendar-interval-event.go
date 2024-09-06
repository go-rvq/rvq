package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VCalendarIntervalEventBuilder struct {
	VTagBuilder[*VCalendarIntervalEventBuilder]
}

func VCalendarIntervalEvent(children ...h.HTMLComponent) *VCalendarIntervalEventBuilder {
	return VTag(&VCalendarIntervalEventBuilder{}, "v-calendar-interval-event", children...)
}

func (b *VCalendarIntervalEventBuilder) AllDay(v bool) (r *VCalendarIntervalEventBuilder) {
	b.Attr(":all-day", fmt.Sprint(v))
	return b
}

func (b *VCalendarIntervalEventBuilder) Interval(v interface{}) (r *VCalendarIntervalEventBuilder) {
	b.Attr(":interval", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalEventBuilder) IntervalDivisions(v int) (r *VCalendarIntervalEventBuilder) {
	b.Attr(":interval-divisions", fmt.Sprint(v))
	return b
}

func (b *VCalendarIntervalEventBuilder) IntervalDuration(v int) (r *VCalendarIntervalEventBuilder) {
	b.Attr(":interval-duration", fmt.Sprint(v))
	return b
}

func (b *VCalendarIntervalEventBuilder) IntervalHeight(v int) (r *VCalendarIntervalEventBuilder) {
	b.Attr(":interval-height", fmt.Sprint(v))
	return b
}

func (b *VCalendarIntervalEventBuilder) Event(v interface{}) (r *VCalendarIntervalEventBuilder) {
	b.Attr(":event", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalEventBuilder) On(name string, value string) (r *VCalendarIntervalEventBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCalendarIntervalEventBuilder) Bind(name string, value string) (r *VCalendarIntervalEventBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
