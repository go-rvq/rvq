package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VInfiniteScrollBuilder struct {
	VTagBuilder[*VInfiniteScrollBuilder]
}

func VInfiniteScroll(children ...h.HTMLComponent) *VInfiniteScrollBuilder {
	return VTag(&VInfiniteScrollBuilder{}, "v-infinite-scroll", children...)
}

func (b *VInfiniteScrollBuilder) Color(v string) (r *VInfiniteScrollBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VInfiniteScrollBuilder) Direction(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) Side(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":side", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) Mode(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":mode", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) Margin(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":margin", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) LoadMoreText(v string) (r *VInfiniteScrollBuilder) {
	b.Attr("load-more-text", v)
	return b
}

func (b *VInfiniteScrollBuilder) EmptyText(v string) (r *VInfiniteScrollBuilder) {
	b.Attr("empty-text", v)
	return b
}

func (b *VInfiniteScrollBuilder) Height(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) MaxHeight(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) MaxWidth(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) MinHeight(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) MinWidth(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) Width(v interface{}) (r *VInfiniteScrollBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VInfiniteScrollBuilder) Tag(v string) (r *VInfiniteScrollBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VInfiniteScrollBuilder) On(name string, value string) (r *VInfiniteScrollBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VInfiniteScrollBuilder) Bind(name string, value string) (r *VInfiniteScrollBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
