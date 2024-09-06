package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VListItemBuilder struct {
	VTagBuilder[*VListItemBuilder]
}

func VListItem(children ...h.HTMLComponent) *VListItemBuilder {
	return VTag(&VListItemBuilder{}, "v-list-item", children...)
}

func (b *VListItemBuilder) Title(v interface{}) (r *VListItemBuilder) {
	b.Attr(":title", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Subtitle(v interface{}) (r *VListItemBuilder) {
	b.Attr(":subtitle", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Active(v bool) (r *VListItemBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) ActiveClass(v string) (r *VListItemBuilder) {
	b.Attr("active-class", v)
	return b
}

func (b *VListItemBuilder) ActiveColor(v string) (r *VListItemBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VListItemBuilder) AppendAvatar(v string) (r *VListItemBuilder) {
	b.Attr("append-avatar", v)
	return b
}

func (b *VListItemBuilder) AppendIcon(v interface{}) (r *VListItemBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) BaseColor(v string) (r *VListItemBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VListItemBuilder) Disabled(v bool) (r *VListItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) Link(v bool) (r *VListItemBuilder) {
	b.Attr(":link", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) Nav(v bool) (r *VListItemBuilder) {
	b.Attr(":nav", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) PrependAvatar(v string) (r *VListItemBuilder) {
	b.Attr("prepend-avatar", v)
	return b
}

func (b *VListItemBuilder) PrependIcon(v interface{}) (r *VListItemBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Ripple(v interface{}) (r *VListItemBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Value(v interface{}) (r *VListItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Slim(v bool) (r *VListItemBuilder) {
	b.Attr(":slim", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) Border(v interface{}) (r *VListItemBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Density(v interface{}) (r *VListItemBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Height(v interface{}) (r *VListItemBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) MaxHeight(v interface{}) (r *VListItemBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) MaxWidth(v interface{}) (r *VListItemBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) MinHeight(v interface{}) (r *VListItemBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) MinWidth(v interface{}) (r *VListItemBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Width(v interface{}) (r *VListItemBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Elevation(v interface{}) (r *VListItemBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Rounded(v interface{}) (r *VListItemBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Tile(v bool) (r *VListItemBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) Href(v string) (r *VListItemBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VListItemBuilder) Replace(v bool) (r *VListItemBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) Exact(v bool) (r *VListItemBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VListItemBuilder) To(v interface{}) (r *VListItemBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Tag(v string) (r *VListItemBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VListItemBuilder) Theme(v string) (r *VListItemBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VListItemBuilder) Color(v string) (r *VListItemBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VListItemBuilder) Variant(v interface{}) (r *VListItemBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) Lines(v interface{}) (r *VListItemBuilder) {
	b.Attr(":lines", h.JSONString(v))
	return b
}

func (b *VListItemBuilder) On(name string, value string) (r *VListItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListItemBuilder) Bind(name string, value string) (r *VListItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
