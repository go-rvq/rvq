package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
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

func (b *VCalendarMonthDayBuilder) Color(v string) (r *VCalendarMonthDayBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VCalendarMonthDayBuilder) Active(v bool) (r *VCalendarMonthDayBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VCalendarMonthDayBuilder) Disabled(v bool) (r *VCalendarMonthDayBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCalendarMonthDayBuilder) Day(v interface{}) (r *VCalendarMonthDayBuilder) {
	b.Attr(":day", h.JSONString(v))
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

func (b *VCalendarMonthDayBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCalendarMonthDayBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCalendarMonthDayBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCalendarMonthDayBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCalendarMonthDayBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCalendarMonthDayBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCalendarMonthDayBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCalendarMonthDayBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCalendarMonthDayBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCalendarMonthDayBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCalendarMonthDayBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCalendarMonthDayBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VCalendarMonthDayBuilder) SetSlotContent(child ...h.HTMLComponent) {
	b.SetSlot("content", child...)
}

func (b *VCalendarMonthDayBuilder) SetScopedSlotContent(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("content", scope, child...)
}

func (b *VCalendarMonthDayBuilder) SlotContent(child ...h.HTMLComponent) (r *VCalendarMonthDayBuilder) {
	b.SetSlotContent(child...)
	return b
}

func (b *VCalendarMonthDayBuilder) ScopedSlotContent(scope string, child ...h.HTMLComponent) (r *VCalendarMonthDayBuilder) {
	b.SetScopedSlotContent(scope, child...)
	return b
}

func (b *VCalendarMonthDayBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VCalendarMonthDayBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VCalendarMonthDayBuilder) SlotTitle(child ...h.HTMLComponent) (r *VCalendarMonthDayBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VCalendarMonthDayBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VCalendarMonthDayBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VCalendarMonthDayBuilder) SetSlotEvent(child ...h.HTMLComponent) {
	b.SetSlot("event", child...)
}

func (b *VCalendarMonthDayBuilder) SetScopedSlotEvent(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("event", scope, child...)
}

func (b *VCalendarMonthDayBuilder) SlotEvent(child ...h.HTMLComponent) (r *VCalendarMonthDayBuilder) {
	b.SetSlotEvent(child...)
	return b
}

func (b *VCalendarMonthDayBuilder) ScopedSlotEvent(scope string, child ...h.HTMLComponent) (r *VCalendarMonthDayBuilder) {
	b.SetScopedSlotEvent(scope, child...)
	return b
}
