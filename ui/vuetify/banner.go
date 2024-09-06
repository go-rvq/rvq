package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VBannerBuilder struct {
	VTagBuilder[*VBannerBuilder]
}

func VBanner(children ...h.HTMLComponent) *VBannerBuilder {
	return VTag(&VBannerBuilder{}, "v-banner", children...)
}

func (b *VBannerBuilder) Text(v string) (r *VBannerBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VBannerBuilder) Avatar(v string) (r *VBannerBuilder) {
	b.Attr("avatar", v)
	return b
}

func (b *VBannerBuilder) BgColor(v string) (r *VBannerBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VBannerBuilder) Color(v string) (r *VBannerBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBannerBuilder) Icon(v interface{}) (r *VBannerBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Stacked(v bool) (r *VBannerBuilder) {
	b.Attr(":stacked", fmt.Sprint(v))
	return b
}

func (b *VBannerBuilder) Sticky(v bool) (r *VBannerBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VBannerBuilder) Border(v interface{}) (r *VBannerBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Density(v interface{}) (r *VBannerBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Height(v interface{}) (r *VBannerBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) MaxHeight(v interface{}) (r *VBannerBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) MaxWidth(v interface{}) (r *VBannerBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) MinHeight(v interface{}) (r *VBannerBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) MinWidth(v interface{}) (r *VBannerBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Width(v interface{}) (r *VBannerBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Mobile(v bool) (r *VBannerBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VBannerBuilder) MobileBreakpoint(v interface{}) (r *VBannerBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Elevation(v interface{}) (r *VBannerBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Location(v interface{}) (r *VBannerBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Position(v interface{}) (r *VBannerBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Rounded(v interface{}) (r *VBannerBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) Tile(v bool) (r *VBannerBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VBannerBuilder) Tag(v string) (r *VBannerBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VBannerBuilder) Theme(v string) (r *VBannerBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VBannerBuilder) Lines(v interface{}) (r *VBannerBuilder) {
	b.Attr(":lines", h.JSONString(v))
	return b
}

func (b *VBannerBuilder) On(name string, value string) (r *VBannerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBannerBuilder) Bind(name string, value string) (r *VBannerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
