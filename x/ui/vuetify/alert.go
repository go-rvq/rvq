package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VAlertBuilder struct {
	VTagBuilder[*VAlertBuilder]
}

func VAlert(children ...h.HTMLComponent) *VAlertBuilder {
	return VTag(&VAlertBuilder{}, "v-alert", children...)
}

func (b *VAlertBuilder) Title(v string) (r *VAlertBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VAlertBuilder) Text(v string) (r *VAlertBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VAlertBuilder) Border(v interface{}) (r *VAlertBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) BorderColor(v string) (r *VAlertBuilder) {
	b.Attr("border-color", v)
	return b
}

func (b *VAlertBuilder) Closable(v bool) (r *VAlertBuilder) {
	b.Attr(":closable", fmt.Sprint(v))
	return b
}

func (b *VAlertBuilder) CloseIcon(v interface{}) (r *VAlertBuilder) {
	b.Attr(":close-icon", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Type(v interface{}) (r *VAlertBuilder) {
	b.Attr(":type", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) CloseLabel(v string) (r *VAlertBuilder) {
	b.Attr("close-label", v)
	return b
}

func (b *VAlertBuilder) Icon(v interface{}) (r *VAlertBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) ModelValue(v bool) (r *VAlertBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VAlertBuilder) Prominent(v bool) (r *VAlertBuilder) {
	b.Attr(":prominent", fmt.Sprint(v))
	return b
}

func (b *VAlertBuilder) Density(v interface{}) (r *VAlertBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Height(v interface{}) (r *VAlertBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) MaxHeight(v interface{}) (r *VAlertBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) MaxWidth(v interface{}) (r *VAlertBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) MinHeight(v interface{}) (r *VAlertBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) MinWidth(v interface{}) (r *VAlertBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Width(v interface{}) (r *VAlertBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Elevation(v interface{}) (r *VAlertBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Location(v interface{}) (r *VAlertBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Position(v interface{}) (r *VAlertBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Rounded(v interface{}) (r *VAlertBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) Tile(v bool) (r *VAlertBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VAlertBuilder) Tag(v string) (r *VAlertBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VAlertBuilder) Theme(v string) (r *VAlertBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VAlertBuilder) Color(v string) (r *VAlertBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VAlertBuilder) Variant(v interface{}) (r *VAlertBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VAlertBuilder) On(name string, value string) (r *VAlertBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAlertBuilder) Bind(name string, value string) (r *VAlertBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
