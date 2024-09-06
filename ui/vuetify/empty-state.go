package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VEmptyStateBuilder struct {
	VTagBuilder[*VEmptyStateBuilder]
}

func VEmptyState(children ...h.HTMLComponent) *VEmptyStateBuilder {
	return VTag(&VEmptyStateBuilder{}, "v-empty-state", children...)
}

func (b *VEmptyStateBuilder) Headline(v string) (r *VEmptyStateBuilder) {
	b.Attr("headline", v)
	return b
}

func (b *VEmptyStateBuilder) Title(v string) (r *VEmptyStateBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VEmptyStateBuilder) Text(v string) (r *VEmptyStateBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VEmptyStateBuilder) ActionText(v string) (r *VEmptyStateBuilder) {
	b.Attr("action-text", v)
	return b
}

func (b *VEmptyStateBuilder) BgColor(v string) (r *VEmptyStateBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VEmptyStateBuilder) Color(v string) (r *VEmptyStateBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VEmptyStateBuilder) Icon(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) Image(v string) (r *VEmptyStateBuilder) {
	b.Attr("image", v)
	return b
}

func (b *VEmptyStateBuilder) Justify(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":justify", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) TextWidth(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":text-width", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) Href(v string) (r *VEmptyStateBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VEmptyStateBuilder) To(v string) (r *VEmptyStateBuilder) {
	b.Attr("to", v)
	return b
}

func (b *VEmptyStateBuilder) Height(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) MaxHeight(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) MaxWidth(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) MinHeight(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) MinWidth(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) Width(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) Size(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) Theme(v string) (r *VEmptyStateBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VEmptyStateBuilder) On(name string, value string) (r *VEmptyStateBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VEmptyStateBuilder) Bind(name string, value string) (r *VEmptyStateBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
