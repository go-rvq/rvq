package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCalendarHeaderBuilder struct {
	VTagBuilder[*VCalendarHeaderBuilder]
}

func VCalendarHeader(children ...h.HTMLComponent) *VCalendarHeaderBuilder {
	return VTag(&VCalendarHeaderBuilder{}, "v-calendar-header", children...)
}

func (b *VCalendarHeaderBuilder) Title(v string) (r *VCalendarHeaderBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VCalendarHeaderBuilder) Text(v string) (r *VCalendarHeaderBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VCalendarHeaderBuilder) NextIcon(v string) (r *VCalendarHeaderBuilder) {
	b.Attr("next-icon", v)
	return b
}

func (b *VCalendarHeaderBuilder) PrevIcon(v string) (r *VCalendarHeaderBuilder) {
	b.Attr("prev-icon", v)
	return b
}

func (b *VCalendarHeaderBuilder) ViewMode(v interface{}) (r *VCalendarHeaderBuilder) {
	b.Attr(":view-mode", h.JSONString(v))
	return b
}

func (b *VCalendarHeaderBuilder) On(name string, value string) (r *VCalendarHeaderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCalendarHeaderBuilder) Bind(name string, value string) (r *VCalendarHeaderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VCalendarHeaderBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCalendarHeaderBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCalendarHeaderBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCalendarHeaderBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCalendarHeaderBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCalendarHeaderBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCalendarHeaderBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCalendarHeaderBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCalendarHeaderBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCalendarHeaderBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCalendarHeaderBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCalendarHeaderBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
