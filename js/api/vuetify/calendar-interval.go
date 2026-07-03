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

func (b *VCalendarIntervalBuilder) Day(v interface{}) (r *VCalendarIntervalBuilder) {
	b.Attr(":day", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalBuilder) DayIndex(v interface{}) (r *VCalendarIntervalBuilder) {
	b.Attr(":day-index", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalBuilder) Events(v interface{}) (r *VCalendarIntervalBuilder) {
	b.Attr(":events", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalBuilder) IntervalDivisions(v interface{}) (r *VCalendarIntervalBuilder) {
	b.Attr(":interval-divisions", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalBuilder) IntervalDuration(v interface{}) (r *VCalendarIntervalBuilder) {
	b.Attr(":interval-duration", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalBuilder) IntervalHeight(v interface{}) (r *VCalendarIntervalBuilder) {
	b.Attr(":interval-height", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalBuilder) IntervalFormat(v interface{}) (r *VCalendarIntervalBuilder) {
	b.Attr(":interval-format", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalBuilder) IntervalStart(v interface{}) (r *VCalendarIntervalBuilder) {
	b.Attr(":interval-start", h.JSONString(v))
	return b
}

func (b *VCalendarIntervalBuilder) Index(v interface{}) (r *VCalendarIntervalBuilder) {
	b.Attr(":index", h.JSONString(v))
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

func (b *VCalendarIntervalBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCalendarIntervalBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCalendarIntervalBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCalendarIntervalBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCalendarIntervalBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCalendarIntervalBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCalendarIntervalBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCalendarIntervalBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCalendarIntervalBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCalendarIntervalBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCalendarIntervalBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCalendarIntervalBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
