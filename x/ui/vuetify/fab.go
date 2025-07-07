package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VFabBuilder struct {
	VTagBuilder[*VFabBuilder]
}

func VFab(children ...h.HTMLComponent) *VFabBuilder {
	return VTag(&VFabBuilder{}, "v-fab", children...)
}

func (b *VFabBuilder) Symbol(v interface{}) (r *VFabBuilder) {
	b.Attr(":symbol", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Flat(v bool) (r *VFabBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) App(v bool) (r *VFabBuilder) {
	b.Attr(":app", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) Appear(v bool) (r *VFabBuilder) {
	b.Attr(":appear", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) Extended(v bool) (r *VFabBuilder) {
	b.Attr(":extended", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) Layout(v bool) (r *VFabBuilder) {
	b.Attr(":layout", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) Location(v interface{}) (r *VFabBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Offset(v bool) (r *VFabBuilder) {
	b.Attr(":offset", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) ModelValue(v bool) (r *VFabBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) Active(v bool) (r *VFabBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) BaseColor(v string) (r *VFabBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VFabBuilder) PrependIcon(v interface{}) (r *VFabBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VFabBuilder) AppendIcon(v interface{}) (r *VFabBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Block(v bool) (r *VFabBuilder) {
	b.Attr(":block", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) Readonly(v bool) (r *VFabBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) Slim(v bool) (r *VFabBuilder) {
	b.Attr(":slim", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) Stacked(v bool) (r *VFabBuilder) {
	b.Attr(":stacked", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) Ripple(v interface{}) (r *VFabBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Value(v interface{}) (r *VFabBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Text(v string) (r *VFabBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VFabBuilder) Border(v interface{}) (r *VFabBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Density(v interface{}) (r *VFabBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Height(v interface{}) (r *VFabBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VFabBuilder) MaxHeight(v interface{}) (r *VFabBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VFabBuilder) MaxWidth(v interface{}) (r *VFabBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VFabBuilder) MinHeight(v interface{}) (r *VFabBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VFabBuilder) MinWidth(v interface{}) (r *VFabBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Width(v interface{}) (r *VFabBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Elevation(v interface{}) (r *VFabBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Disabled(v bool) (r *VFabBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) SelectedClass(v string) (r *VFabBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VFabBuilder) Loading(v interface{}) (r *VFabBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Position(v interface{}) (r *VFabBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Absolute(v bool) (r *VFabBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) Rounded(v interface{}) (r *VFabBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Tile(v bool) (r *VFabBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) Href(v string) (r *VFabBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VFabBuilder) Replace(v bool) (r *VFabBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) Exact(v bool) (r *VFabBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VFabBuilder) To(v interface{}) (r *VFabBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Size(v interface{}) (r *VFabBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Tag(v string) (r *VFabBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VFabBuilder) Theme(v string) (r *VFabBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VFabBuilder) Color(v string) (r *VFabBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VFabBuilder) Variant(v interface{}) (r *VFabBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Icon(v interface{}) (r *VFabBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Name(v string) (r *VFabBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VFabBuilder) Order(v interface{}) (r *VFabBuilder) {
	b.Attr(":order", h.JSONString(v))
	return b
}

func (b *VFabBuilder) Transition(v interface{}) (r *VFabBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VFabBuilder) On(name string, value string) (r *VFabBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFabBuilder) Bind(name string, value string) (r *VFabBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
