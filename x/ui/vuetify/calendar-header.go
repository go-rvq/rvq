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

func (b *VCalendarHeaderBuilder) NextIcon(v string) (r *VCalendarHeaderBuilder) {
	b.Attr("next-icon", v)
	return b
}

func (b *VCalendarHeaderBuilder) PrevIcon(v string) (r *VCalendarHeaderBuilder) {
	b.Attr("prev-icon", v)
	return b
}

func (b *VCalendarHeaderBuilder) Title(v string) (r *VCalendarHeaderBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VCalendarHeaderBuilder) Text(v string) (r *VCalendarHeaderBuilder) {
	b.Attr("text", v)
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
