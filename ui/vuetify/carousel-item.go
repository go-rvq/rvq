package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VCarouselItemBuilder struct {
	VTagBuilder[*VCarouselItemBuilder]
}

func VCarouselItem(children ...h.HTMLComponent) (r *VCarouselItemBuilder) {
	return VTag(&VCarouselItemBuilder{}, "v-carousel-item", children...)
}

func (b *VCarouselItemBuilder) Alt(v string) (r *VCarouselItemBuilder) {
	b.Attr("alt", v)
	return b
}

func (b *VCarouselItemBuilder) Cover(v bool) (r *VCarouselItemBuilder) {
	b.Attr(":cover", fmt.Sprint(v))
	return b
}

func (b *VCarouselItemBuilder) Color(v string) (r *VCarouselItemBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VCarouselItemBuilder) Draggable(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":draggable", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Eager(v bool) (r *VCarouselItemBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VCarouselItemBuilder) Gradient(v string) (r *VCarouselItemBuilder) {
	b.Attr("gradient", v)
	return b
}

func (b *VCarouselItemBuilder) LazySrc(v string) (r *VCarouselItemBuilder) {
	b.Attr("lazy-src", v)
	return b
}

func (b *VCarouselItemBuilder) Options(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":options", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Sizes(v string) (r *VCarouselItemBuilder) {
	b.Attr("sizes", v)
	return b
}

func (b *VCarouselItemBuilder) Src(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":src", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Srcset(v string) (r *VCarouselItemBuilder) {
	b.Attr("srcset", v)
	return b
}

func (b *VCarouselItemBuilder) Position(v string) (r *VCarouselItemBuilder) {
	b.Attr("position", v)
	return b
}

func (b *VCarouselItemBuilder) AspectRatio(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":aspect-ratio", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) ContentClass(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":content-class", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Inline(v bool) (r *VCarouselItemBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VCarouselItemBuilder) Height(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) MaxHeight(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) MaxWidth(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) MinHeight(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) MinWidth(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Width(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Rounded(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Tile(v bool) (r *VCarouselItemBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VCarouselItemBuilder) Transition(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Crossorigin(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":crossorigin", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Referrerpolicy(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":referrerpolicy", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) ReverseTransition(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":reverse-transition", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Value(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Disabled(v bool) (r *VCarouselItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCarouselItemBuilder) SelectedClass(v string) (r *VCarouselItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VCarouselItemBuilder) On(name string, value string) (r *VCarouselItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCarouselItemBuilder) Bind(name string, value string) (r *VCarouselItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
