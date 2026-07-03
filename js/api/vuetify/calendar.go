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

func (b *VCalendarBuilder) Title(v string) (r *VCalendarBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VCalendarBuilder) Text(v string) (r *VCalendarBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VCalendarBuilder) ModelValue(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) Disabled(v bool) (r *VCalendarBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) Max(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) HideHeader(v bool) (r *VCalendarBuilder) {
	b.Attr(":hide-header", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) HideWeekNumber(v bool) (r *VCalendarBuilder) {
	b.Attr(":hide-week-number", fmt.Sprint(v))
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

func (b *VCalendarBuilder) FirstDayOfWeek(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":first-day-of-week", h.JSONString(v))
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

func (b *VCalendarBuilder) Min(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) HideDayHeader(v bool) (r *VCalendarBuilder) {
	b.Attr(":hide-day-header", fmt.Sprint(v))
	return b
}

func (b *VCalendarBuilder) Intervals(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":intervals", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) Day(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":day", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) DayIndex(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":day-index", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) Events(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":events", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) IntervalDivisions(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":interval-divisions", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) IntervalDuration(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":interval-duration", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) IntervalHeight(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":interval-height", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) IntervalFormat(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":interval-format", h.JSONString(v))
	return b
}

func (b *VCalendarBuilder) IntervalStart(v interface{}) (r *VCalendarBuilder) {
	b.Attr(":interval-start", h.JSONString(v))
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

func (b *VCalendarBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCalendarBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCalendarBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCalendarBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCalendarBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCalendarBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCalendarBuilder) SetSlotHeader(child ...h.HTMLComponent) {
	b.SetSlot("header", child...)
}

func (b *VCalendarBuilder) SetScopedSlotHeader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header", scope, child...)
}

func (b *VCalendarBuilder) SlotHeader(child ...h.HTMLComponent) (r *VCalendarBuilder) {
	b.SetSlotHeader(child...)
	return b
}

func (b *VCalendarBuilder) ScopedSlotHeader(scope string, child ...h.HTMLComponent) (r *VCalendarBuilder) {
	b.SetScopedSlotHeader(scope, child...)
	return b
}

func (b *VCalendarBuilder) SetSlotEvent(child ...h.HTMLComponent) {
	b.SetSlot("event", child...)
}

func (b *VCalendarBuilder) SetScopedSlotEvent(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("event", scope, child...)
}

func (b *VCalendarBuilder) SlotEvent(child ...h.HTMLComponent) (r *VCalendarBuilder) {
	b.SetSlotEvent(child...)
	return b
}

func (b *VCalendarBuilder) ScopedSlotEvent(scope string, child ...h.HTMLComponent) (r *VCalendarBuilder) {
	b.SetScopedSlotEvent(scope, child...)
	return b
}
