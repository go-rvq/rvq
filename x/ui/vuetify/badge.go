package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBadgeBuilder struct {
	VTagBuilder[*VBadgeBuilder]
}

func VBadge(children ...h.HTMLComponent) *VBadgeBuilder {
	return VTag(&VBadgeBuilder{}, "v-badge", children...)
}

func (b *VBadgeBuilder) Bordered(v bool) (r *VBadgeBuilder) {
	b.Attr(":bordered", fmt.Sprint(v))
	return b
}

func (b *VBadgeBuilder) Color(v string) (r *VBadgeBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBadgeBuilder) Content(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":content", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) Dot(v bool) (r *VBadgeBuilder) {
	b.Attr(":dot", fmt.Sprint(v))
	return b
}

func (b *VBadgeBuilder) Floating(v bool) (r *VBadgeBuilder) {
	b.Attr(":floating", fmt.Sprint(v))
	return b
}

func (b *VBadgeBuilder) Icon(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) Inline(v bool) (r *VBadgeBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VBadgeBuilder) Label(v string) (r *VBadgeBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VBadgeBuilder) Max(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) ModelValue(v bool) (r *VBadgeBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VBadgeBuilder) OffsetX(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":offset-x", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) OffsetY(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":offset-y", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) TextColor(v string) (r *VBadgeBuilder) {
	b.Attr("text-color", v)
	return b
}

func (b *VBadgeBuilder) Location(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) Rounded(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) Tile(v bool) (r *VBadgeBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VBadgeBuilder) Tag(v string) (r *VBadgeBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VBadgeBuilder) Theme(v string) (r *VBadgeBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VBadgeBuilder) Transition(v interface{}) (r *VBadgeBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VBadgeBuilder) On(name string, value string) (r *VBadgeBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBadgeBuilder) Bind(name string, value string) (r *VBadgeBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
