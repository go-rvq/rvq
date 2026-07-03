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

func (b *VCalendarDayBuilder) Intervals(v interface{}) (r *VCalendarDayBuilder) {
	b.Attr(":intervals", h.JSONString(v))
	return b
}

func (b *VCalendarDayBuilder) Day(v interface{}) (r *VCalendarDayBuilder) {
	b.Attr(":day", h.JSONString(v))
	return b
}

func (b *VCalendarDayBuilder) DayIndex(v interface{}) (r *VCalendarDayBuilder) {
	b.Attr(":day-index", h.JSONString(v))
	return b
}

func (b *VCalendarDayBuilder) Events(v interface{}) (r *VCalendarDayBuilder) {
	b.Attr(":events", h.JSONString(v))
	return b
}

func (b *VCalendarDayBuilder) IntervalDivisions(v interface{}) (r *VCalendarDayBuilder) {
	b.Attr(":interval-divisions", h.JSONString(v))
	return b
}

func (b *VCalendarDayBuilder) IntervalDuration(v interface{}) (r *VCalendarDayBuilder) {
	b.Attr(":interval-duration", h.JSONString(v))
	return b
}

func (b *VCalendarDayBuilder) IntervalHeight(v interface{}) (r *VCalendarDayBuilder) {
	b.Attr(":interval-height", h.JSONString(v))
	return b
}

func (b *VCalendarDayBuilder) IntervalFormat(v interface{}) (r *VCalendarDayBuilder) {
	b.Attr(":interval-format", h.JSONString(v))
	return b
}

func (b *VCalendarDayBuilder) IntervalStart(v interface{}) (r *VCalendarDayBuilder) {
	b.Attr(":interval-start", h.JSONString(v))
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

func (b *VCalendarDayBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCalendarDayBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCalendarDayBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCalendarDayBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCalendarDayBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCalendarDayBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCalendarDayBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCalendarDayBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCalendarDayBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCalendarDayBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCalendarDayBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCalendarDayBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
