package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDatePickerMonthBuilder struct {
	VTagBuilder[*VDatePickerMonthBuilder]
}

func VDatePickerMonth(children ...h.HTMLComponent) *VDatePickerMonthBuilder {
	return VTag(&VDatePickerMonthBuilder{}, "v-date-picker-month", children...)
}

func (b *VDatePickerMonthBuilder) Disabled(v bool) (r *VDatePickerMonthBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VDatePickerMonthBuilder) Multiple(v interface{}) (r *VDatePickerMonthBuilder) {
	b.Attr(":multiple", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthBuilder) Color(v string) (r *VDatePickerMonthBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDatePickerMonthBuilder) ModelValue(v interface{}) (r *VDatePickerMonthBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthBuilder) Transition(v string) (r *VDatePickerMonthBuilder) {
	b.Attr("transition", v)
	return b
}

func (b *VDatePickerMonthBuilder) Max(v interface{}) (r *VDatePickerMonthBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthBuilder) Min(v interface{}) (r *VDatePickerMonthBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthBuilder) Month(v interface{}) (r *VDatePickerMonthBuilder) {
	b.Attr(":month", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthBuilder) Year(v interface{}) (r *VDatePickerMonthBuilder) {
	b.Attr(":year", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthBuilder) HideWeekdays(v bool) (r *VDatePickerMonthBuilder) {
	b.Attr(":hide-weekdays", fmt.Sprint(v))
	return b
}

func (b *VDatePickerMonthBuilder) ShowWeek(v bool) (r *VDatePickerMonthBuilder) {
	b.Attr(":show-week", fmt.Sprint(v))
	return b
}

func (b *VDatePickerMonthBuilder) ReverseTransition(v string) (r *VDatePickerMonthBuilder) {
	b.Attr("reverse-transition", v)
	return b
}

func (b *VDatePickerMonthBuilder) ShowAdjacentMonths(v bool) (r *VDatePickerMonthBuilder) {
	b.Attr(":show-adjacent-months", fmt.Sprint(v))
	return b
}

func (b *VDatePickerMonthBuilder) Weekdays(v interface{}) (r *VDatePickerMonthBuilder) {
	b.Attr(":weekdays", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthBuilder) WeeksInMonth(v interface{}) (r *VDatePickerMonthBuilder) {
	b.Attr(":weeks-in-month", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthBuilder) FirstDayOfWeek(v interface{}) (r *VDatePickerMonthBuilder) {
	b.Attr(":first-day-of-week", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthBuilder) AllowedDates(v interface{}) (r *VDatePickerMonthBuilder) {
	b.Attr(":allowed-dates", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthBuilder) On(name string, value string) (r *VDatePickerMonthBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDatePickerMonthBuilder) Bind(name string, value string) (r *VDatePickerMonthBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDatePickerMonthBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDatePickerMonthBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDatePickerMonthBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDatePickerMonthBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDatePickerMonthBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDatePickerMonthBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDatePickerMonthBuilder) SetSlotDay(child ...h.HTMLComponent) {
	b.SetSlot("day", child...)
}

func (b *VDatePickerMonthBuilder) SetScopedSlotDay(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("day", scope, child...)
}

func (b *VDatePickerMonthBuilder) SlotDay(child ...h.HTMLComponent) (r *VDatePickerMonthBuilder) {
	b.SetSlotDay(child...)
	return b
}

func (b *VDatePickerMonthBuilder) ScopedSlotDay(scope string, child ...h.HTMLComponent) (r *VDatePickerMonthBuilder) {
	b.SetScopedSlotDay(scope, child...)
	return b
}
