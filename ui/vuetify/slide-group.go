package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSlideGroupBuilder struct {
	VTagBuilder[*VSlideGroupBuilder]
}

func VSlideGroup(children ...h.HTMLComponent) *VSlideGroupBuilder {
	return VTag(&VSlideGroupBuilder{}, "v-slide-group", children...)
}

func (b *VSlideGroupBuilder) Symbol(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":symbol", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) CenterActive(v bool) (r *VSlideGroupBuilder) {
	b.Attr(":center-active", fmt.Sprint(v))
	return b
}

func (b *VSlideGroupBuilder) Direction(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) NextIcon(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) PrevIcon(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) ShowArrows(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":show-arrows", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) Mobile(v bool) (r *VSlideGroupBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VSlideGroupBuilder) MobileBreakpoint(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) Tag(v string) (r *VSlideGroupBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VSlideGroupBuilder) ModelValue(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) Multiple(v bool) (r *VSlideGroupBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VSlideGroupBuilder) Max(v int) (r *VSlideGroupBuilder) {
	b.Attr(":max", fmt.Sprint(v))
	return b
}

func (b *VSlideGroupBuilder) SelectedClass(v string) (r *VSlideGroupBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VSlideGroupBuilder) Disabled(v bool) (r *VSlideGroupBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSlideGroupBuilder) Mandatory(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) On(name string, value string) (r *VSlideGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSlideGroupBuilder) Bind(name string, value string) (r *VSlideGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
