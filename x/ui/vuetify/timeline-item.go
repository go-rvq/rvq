package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTimelineItemBuilder struct {
	VTagBuilder[*VTimelineItemBuilder]
}

func VTimelineItem(children ...h.HTMLComponent) *VTimelineItemBuilder {
	return VTag(&VTimelineItemBuilder{}, "v-timeline-item", children...)
}

func (b *VTimelineItemBuilder) Icon(v interface{}) (r *VTimelineItemBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VTimelineItemBuilder) DotColor(v string) (r *VTimelineItemBuilder) {
	b.Attr("dot-color", v)
	return b
}

func (b *VTimelineItemBuilder) FillDot(v bool) (r *VTimelineItemBuilder) {
	b.Attr(":fill-dot", fmt.Sprint(v))
	return b
}

func (b *VTimelineItemBuilder) HideDot(v bool) (r *VTimelineItemBuilder) {
	b.Attr(":hide-dot", fmt.Sprint(v))
	return b
}

func (b *VTimelineItemBuilder) HideOpposite(v bool) (r *VTimelineItemBuilder) {
	b.Attr(":hide-opposite", fmt.Sprint(v))
	return b
}

func (b *VTimelineItemBuilder) IconColor(v string) (r *VTimelineItemBuilder) {
	b.Attr("icon-color", v)
	return b
}

func (b *VTimelineItemBuilder) LineInset(v interface{}) (r *VTimelineItemBuilder) {
	b.Attr(":line-inset", h.JSONString(v))
	return b
}

func (b *VTimelineItemBuilder) Height(v interface{}) (r *VTimelineItemBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VTimelineItemBuilder) MaxHeight(v interface{}) (r *VTimelineItemBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VTimelineItemBuilder) MaxWidth(v interface{}) (r *VTimelineItemBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VTimelineItemBuilder) MinHeight(v interface{}) (r *VTimelineItemBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VTimelineItemBuilder) MinWidth(v interface{}) (r *VTimelineItemBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VTimelineItemBuilder) Width(v interface{}) (r *VTimelineItemBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VTimelineItemBuilder) Elevation(v interface{}) (r *VTimelineItemBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VTimelineItemBuilder) Rounded(v interface{}) (r *VTimelineItemBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VTimelineItemBuilder) Tile(v bool) (r *VTimelineItemBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VTimelineItemBuilder) Size(v interface{}) (r *VTimelineItemBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VTimelineItemBuilder) Tag(v string) (r *VTimelineItemBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VTimelineItemBuilder) Density(v interface{}) (r *VTimelineItemBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VTimelineItemBuilder) On(name string, value string) (r *VTimelineItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTimelineItemBuilder) Bind(name string, value string) (r *VTimelineItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
