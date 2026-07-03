package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCalendarIntervalEventBuilder struct {
	VTagBuilder[*VCalendarIntervalEventBuilder]
}

func VCalendarIntervalEvent(children ...h.HTMLComponent) *VCalendarIntervalEventBuilder {
	return VTag(&VCalendarIntervalEventBuilder{}, "v-calendar-interval-event", children...)
}

func (b *VCalendarIntervalEventBuilder) Event(v interface{}) (r *VCalendarIntervalEventBuilder) {
	b.Attr(":event", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalEventBuilder) IntervalDivisions(v interface{}) (r *VCalendarIntervalEventBuilder) {
	b.Attr(":interval-divisions", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalEventBuilder) IntervalDuration(v interface{}) (r *VCalendarIntervalEventBuilder) {
	b.Attr(":interval-duration", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalEventBuilder) IntervalHeight(v interface{}) (r *VCalendarIntervalEventBuilder) {
	b.Attr(":interval-height", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalEventBuilder) Interval(v interface{}) (r *VCalendarIntervalEventBuilder) {
	b.Attr(":interval", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalEventBuilder) AllDay(v bool) (r *VCalendarIntervalEventBuilder) {
	b.Attr(":all-day", fmt.Sprint(v))
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

func (b *VCalendarIntervalEventBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCalendarIntervalEventBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCalendarIntervalEventBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCalendarIntervalEventBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCalendarIntervalEventBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCalendarIntervalEventBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCalendarIntervalEventBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCalendarIntervalEventBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCalendarIntervalEventBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCalendarIntervalEventBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCalendarIntervalEventBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCalendarIntervalEventBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
