package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VChipBuilder struct {
	VTagBuilder[*VChipBuilder]
}

func VChip(children ...h.HTMLComponent) *VChipBuilder {
	return VTag(&VChipBuilder{}, "v-chip", children...)
}

func (b *VChipBuilder) Label(v bool) (r *VChipBuilder) {
	b.Attr(":label", fmt.Sprint(v))
	return b
}

func (b *VChipBuilder) Filter(v bool) (r *VChipBuilder) {
	b.Attr(":filter", fmt.Sprint(v))
	return b
}

func (b *VChipBuilder) ActiveClass(v string) (r *VChipBuilder) {
	b.Attr("active-class", v)
	return b
}

func (b *VChipBuilder) AppendAvatar(v string) (r *VChipBuilder) {
	b.Attr("append-avatar", v)
	return b
}

func (b *VChipBuilder) AppendIcon(v interface{}) (r *VChipBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VChipBuilder) Closable(v bool) (r *VChipBuilder) {
	b.Attr(":closable", fmt.Sprint(v))
	return b
}

func (b *VChipBuilder) CloseIcon(v interface{}) (r *VChipBuilder) {
	b.Attr(":close-icon", h.JSONString(v))
	return b
}

func (b *VChipBuilder) CloseLabel(v string) (r *VChipBuilder) {
	b.Attr("close-label", v)
	return b
}

func (b *VChipBuilder) Draggable(v bool) (r *VChipBuilder) {
	b.Attr(":draggable", fmt.Sprint(v))
	return b
}

func (b *VChipBuilder) FilterIcon(v string) (r *VChipBuilder) {
	b.Attr("filter-icon", v)
	return b
}

func (b *VChipBuilder) Link(v bool) (r *VChipBuilder) {
	b.Attr(":link", fmt.Sprint(v))
	return b
}

func (b *VChipBuilder) Pill(v bool) (r *VChipBuilder) {
	b.Attr(":pill", fmt.Sprint(v))
	return b
}

func (b *VChipBuilder) PrependAvatar(v string) (r *VChipBuilder) {
	b.Attr("prepend-avatar", v)
	return b
}

func (b *VChipBuilder) PrependIcon(v interface{}) (r *VChipBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VChipBuilder) Ripple(v interface{}) (r *VChipBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VChipBuilder) Value(v interface{}) (r *VChipBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VChipBuilder) Text(v string) (r *VChipBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VChipBuilder) ModelValue(v bool) (r *VChipBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VChipBuilder) Border(v interface{}) (r *VChipBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VChipBuilder) Density(v interface{}) (r *VChipBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VChipBuilder) Elevation(v interface{}) (r *VChipBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VChipBuilder) Disabled(v bool) (r *VChipBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VChipBuilder) SelectedClass(v string) (r *VChipBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VChipBuilder) Rounded(v interface{}) (r *VChipBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VChipBuilder) Tile(v bool) (r *VChipBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VChipBuilder) Href(v string) (r *VChipBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VChipBuilder) Replace(v bool) (r *VChipBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VChipBuilder) Exact(v bool) (r *VChipBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VChipBuilder) To(v interface{}) (r *VChipBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VChipBuilder) Size(v interface{}) (r *VChipBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VChipBuilder) Tag(v string) (r *VChipBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VChipBuilder) Theme(v string) (r *VChipBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VChipBuilder) Color(v string) (r *VChipBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VChipBuilder) Variant(v interface{}) (r *VChipBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VChipBuilder) On(name string, value string) (r *VChipBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VChipBuilder) Bind(name string, value string) (r *VChipBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
