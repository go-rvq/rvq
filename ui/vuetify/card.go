package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VCardBuilder struct {
	VTagBuilder[*VCardBuilder]
}

func VCard(children ...h.HTMLComponent) *VCardBuilder {
	return VTag(&VCardBuilder{}, "v-card", children...)
}

func (b *VCardBuilder) Title(v interface{}) (r *VCardBuilder) {
	b.Attr(":title", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Subtitle(v interface{}) (r *VCardBuilder) {
	b.Attr(":subtitle", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Text(v interface{}) (r *VCardBuilder) {
	b.Attr(":text", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Image(v string) (r *VCardBuilder) {
	b.Attr("image", v)
	return b
}

func (b *VCardBuilder) Flat(v bool) (r *VCardBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) AppendAvatar(v string) (r *VCardBuilder) {
	b.Attr("append-avatar", v)
	return b
}

func (b *VCardBuilder) AppendIcon(v interface{}) (r *VCardBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Disabled(v bool) (r *VCardBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) Hover(v bool) (r *VCardBuilder) {
	b.Attr(":hover", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) Link(v bool) (r *VCardBuilder) {
	b.Attr(":link", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) PrependAvatar(v string) (r *VCardBuilder) {
	b.Attr("prepend-avatar", v)
	return b
}

func (b *VCardBuilder) PrependIcon(v interface{}) (r *VCardBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Ripple(v interface{}) (r *VCardBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Border(v interface{}) (r *VCardBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Density(v interface{}) (r *VCardBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Height(v interface{}) (r *VCardBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VCardBuilder) MaxHeight(v interface{}) (r *VCardBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VCardBuilder) MaxWidth(v interface{}) (r *VCardBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VCardBuilder) MinHeight(v interface{}) (r *VCardBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VCardBuilder) MinWidth(v interface{}) (r *VCardBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Width(v interface{}) (r *VCardBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Elevation(v interface{}) (r *VCardBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Loading(v interface{}) (r *VCardBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Location(v interface{}) (r *VCardBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Position(v interface{}) (r *VCardBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Rounded(v interface{}) (r *VCardBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Tile(v bool) (r *VCardBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) Href(v string) (r *VCardBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VCardBuilder) Replace(v bool) (r *VCardBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) Exact(v bool) (r *VCardBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VCardBuilder) To(v interface{}) (r *VCardBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VCardBuilder) Tag(v string) (r *VCardBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VCardBuilder) Theme(v string) (r *VCardBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VCardBuilder) Color(v string) (r *VCardBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VCardBuilder) Variant(v interface{}) (r *VCardBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VCardBuilder) On(name string, value string) (r *VCardBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCardBuilder) Bind(name string, value string) (r *VCardBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
