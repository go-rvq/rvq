package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VTreeviewItemBuilder struct {
	VTagBuilder[*VTreeviewItemBuilder]
}

func VTreeviewItem(children ...h.HTMLComponent) *VTreeviewItemBuilder {
	return VTag(&VTreeviewItemBuilder{}, "v-treeview-item", children...)
}

func (b *VTreeviewItemBuilder) Title(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":title", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Subtitle(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":subtitle", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Loading(v bool) (r *VTreeviewItemBuilder) {
	b.Attr(":loading", fmt.Sprint(v))
	return b
}

func (b *VTreeviewItemBuilder) ToggleIcon(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":toggle-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Active(v bool) (r *VTreeviewItemBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VTreeviewItemBuilder) ActiveClass(v string) (r *VTreeviewItemBuilder) {
	b.Attr("active-class", v)
	return b
}

func (b *VTreeviewItemBuilder) ActiveColor(v string) (r *VTreeviewItemBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VTreeviewItemBuilder) AppendAvatar(v string) (r *VTreeviewItemBuilder) {
	b.Attr("append-avatar", v)
	return b
}

func (b *VTreeviewItemBuilder) AppendIcon(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) BaseColor(v string) (r *VTreeviewItemBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VTreeviewItemBuilder) Disabled(v bool) (r *VTreeviewItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTreeviewItemBuilder) Link(v bool) (r *VTreeviewItemBuilder) {
	b.Attr(":link", fmt.Sprint(v))
	return b
}

func (b *VTreeviewItemBuilder) Nav(v bool) (r *VTreeviewItemBuilder) {
	b.Attr(":nav", fmt.Sprint(v))
	return b
}

func (b *VTreeviewItemBuilder) PrependAvatar(v string) (r *VTreeviewItemBuilder) {
	b.Attr("prepend-avatar", v)
	return b
}

func (b *VTreeviewItemBuilder) PrependIcon(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Ripple(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Value(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Slim(v bool) (r *VTreeviewItemBuilder) {
	b.Attr(":slim", fmt.Sprint(v))
	return b
}

func (b *VTreeviewItemBuilder) Border(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Density(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Height(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) MaxHeight(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) MaxWidth(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) MinHeight(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) MinWidth(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Width(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Elevation(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Rounded(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Tile(v bool) (r *VTreeviewItemBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VTreeviewItemBuilder) Href(v string) (r *VTreeviewItemBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VTreeviewItemBuilder) Replace(v bool) (r *VTreeviewItemBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VTreeviewItemBuilder) Exact(v bool) (r *VTreeviewItemBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VTreeviewItemBuilder) To(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Tag(v string) (r *VTreeviewItemBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VTreeviewItemBuilder) Theme(v string) (r *VTreeviewItemBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTreeviewItemBuilder) Color(v string) (r *VTreeviewItemBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTreeviewItemBuilder) Variant(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) Lines(v interface{}) (r *VTreeviewItemBuilder) {
	b.Attr(":lines", h.JSONString(v))
	return b
}

func (b *VTreeviewItemBuilder) On(name string, value string) (r *VTreeviewItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTreeviewItemBuilder) Bind(name string, value string) (r *VTreeviewItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
