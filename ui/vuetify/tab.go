package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VTabBuilder struct {
	VTagBuilder[*VTabBuilder]
}

func VTab(children ...h.HTMLComponent) *VTabBuilder {
	return VTag(&VTabBuilder{}, "v-tab", children...)
}

func (b *VTabBuilder) Fixed(v bool) (r *VTabBuilder) {
	b.Attr(":fixed", fmt.Sprint(v))
	return b
}

func (b *VTabBuilder) SliderColor(v string) (r *VTabBuilder) {
	b.Attr("slider-color", v)
	return b
}

func (b *VTabBuilder) HideSlider(v bool) (r *VTabBuilder) {
	b.Attr(":hide-slider", fmt.Sprint(v))
	return b
}

func (b *VTabBuilder) Direction(v interface{}) (r *VTabBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VTabBuilder) BaseColor(v string) (r *VTabBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VTabBuilder) PrependIcon(v interface{}) (r *VTabBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VTabBuilder) AppendIcon(v interface{}) (r *VTabBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VTabBuilder) Readonly(v bool) (r *VTabBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VTabBuilder) Slim(v bool) (r *VTabBuilder) {
	b.Attr(":slim", fmt.Sprint(v))
	return b
}

func (b *VTabBuilder) Stacked(v bool) (r *VTabBuilder) {
	b.Attr(":stacked", fmt.Sprint(v))
	return b
}

func (b *VTabBuilder) Ripple(v interface{}) (r *VTabBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VTabBuilder) Value(v interface{}) (r *VTabBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VTabBuilder) Text(v string) (r *VTabBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VTabBuilder) Border(v interface{}) (r *VTabBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VTabBuilder) Density(v interface{}) (r *VTabBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VTabBuilder) Height(v interface{}) (r *VTabBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VTabBuilder) MaxHeight(v interface{}) (r *VTabBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VTabBuilder) MaxWidth(v interface{}) (r *VTabBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VTabBuilder) MinHeight(v interface{}) (r *VTabBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VTabBuilder) MinWidth(v interface{}) (r *VTabBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VTabBuilder) Width(v interface{}) (r *VTabBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VTabBuilder) Elevation(v interface{}) (r *VTabBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VTabBuilder) Disabled(v bool) (r *VTabBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTabBuilder) SelectedClass(v string) (r *VTabBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VTabBuilder) Loading(v interface{}) (r *VTabBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VTabBuilder) Rounded(v interface{}) (r *VTabBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VTabBuilder) Tile(v bool) (r *VTabBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VTabBuilder) Href(v string) (r *VTabBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VTabBuilder) Replace(v bool) (r *VTabBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VTabBuilder) Exact(v bool) (r *VTabBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VTabBuilder) To(v interface{}) (r *VTabBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VTabBuilder) Size(v interface{}) (r *VTabBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VTabBuilder) Tag(v string) (r *VTabBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VTabBuilder) Theme(v string) (r *VTabBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTabBuilder) Color(v string) (r *VTabBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTabBuilder) Variant(v interface{}) (r *VTabBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VTabBuilder) Icon(v interface{}) (r *VTabBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VTabBuilder) On(name string, value string) (r *VTabBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTabBuilder) Bind(name string, value string) (r *VTabBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
