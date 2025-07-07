package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VCarouselBuilder struct {
	VTagBuilder[*VCarouselBuilder]
}

func VCarousel(children ...h.HTMLComponent) *VCarouselBuilder {
	return VTag(&VCarouselBuilder{}, "v-carousel", children...)
}

func (b *VCarouselBuilder) Color(v string) (r *VCarouselBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VCarouselBuilder) Cycle(v bool) (r *VCarouselBuilder) {
	b.Attr(":cycle", fmt.Sprint(v))
	return b
}

func (b *VCarouselBuilder) DelimiterIcon(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":delimiter-icon", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Height(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) HideDelimiters(v bool) (r *VCarouselBuilder) {
	b.Attr(":hide-delimiters", fmt.Sprint(v))
	return b
}

func (b *VCarouselBuilder) HideDelimiterBackground(v bool) (r *VCarouselBuilder) {
	b.Attr(":hide-delimiter-background", fmt.Sprint(v))
	return b
}

func (b *VCarouselBuilder) Interval(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":interval", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Progress(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":progress", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Continuous(v bool) (r *VCarouselBuilder) {
	b.Attr(":continuous", fmt.Sprint(v))
	return b
}

func (b *VCarouselBuilder) NextIcon(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) PrevIcon(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Reverse(v bool) (r *VCarouselBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VCarouselBuilder) ShowArrows(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":show-arrows", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Touch(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":touch", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Direction(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) ModelValue(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Disabled(v bool) (r *VCarouselBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCarouselBuilder) SelectedClass(v string) (r *VCarouselBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VCarouselBuilder) Mandatory(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Tag(v string) (r *VCarouselBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VCarouselBuilder) Theme(v string) (r *VCarouselBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VCarouselBuilder) VerticalDelimiters(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":vertical-delimiters", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) On(name string, value string) (r *VCarouselBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCarouselBuilder) Bind(name string, value string) (r *VCarouselBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
