package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VExpansionPanelsBuilder struct {
	VTagBuilder[*VExpansionPanelsBuilder]
}

func VExpansionPanels(children ...h.HTMLComponent) *VExpansionPanelsBuilder {
	return VTag(&VExpansionPanelsBuilder{}, "v-expansion-panels", children...)
}

func (b *VExpansionPanelsBuilder) Flat(v bool) (r *VExpansionPanelsBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelsBuilder) ModelValue(v interface{}) (r *VExpansionPanelsBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VExpansionPanelsBuilder) Multiple(v bool) (r *VExpansionPanelsBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelsBuilder) Max(v int) (r *VExpansionPanelsBuilder) {
	b.Attr(":max", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelsBuilder) SelectedClass(v string) (r *VExpansionPanelsBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VExpansionPanelsBuilder) Disabled(v bool) (r *VExpansionPanelsBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelsBuilder) Mandatory(v interface{}) (r *VExpansionPanelsBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VExpansionPanelsBuilder) Title(v string) (r *VExpansionPanelsBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VExpansionPanelsBuilder) Text(v string) (r *VExpansionPanelsBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VExpansionPanelsBuilder) BgColor(v string) (r *VExpansionPanelsBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VExpansionPanelsBuilder) Elevation(v interface{}) (r *VExpansionPanelsBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VExpansionPanelsBuilder) Value(v interface{}) (r *VExpansionPanelsBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VExpansionPanelsBuilder) Rounded(v interface{}) (r *VExpansionPanelsBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VExpansionPanelsBuilder) Tile(v bool) (r *VExpansionPanelsBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelsBuilder) Tag(v string) (r *VExpansionPanelsBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VExpansionPanelsBuilder) Color(v string) (r *VExpansionPanelsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VExpansionPanelsBuilder) ExpandIcon(v interface{}) (r *VExpansionPanelsBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VExpansionPanelsBuilder) CollapseIcon(v interface{}) (r *VExpansionPanelsBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VExpansionPanelsBuilder) HideActions(v bool) (r *VExpansionPanelsBuilder) {
	b.Attr(":hide-actions", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelsBuilder) Focusable(v bool) (r *VExpansionPanelsBuilder) {
	b.Attr(":focusable", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelsBuilder) Static(v bool) (r *VExpansionPanelsBuilder) {
	b.Attr(":static", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelsBuilder) Ripple(v interface{}) (r *VExpansionPanelsBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VExpansionPanelsBuilder) Readonly(v bool) (r *VExpansionPanelsBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelsBuilder) Eager(v bool) (r *VExpansionPanelsBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelsBuilder) Theme(v string) (r *VExpansionPanelsBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VExpansionPanelsBuilder) Variant(v interface{}) (r *VExpansionPanelsBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VExpansionPanelsBuilder) On(name string, value string) (r *VExpansionPanelsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VExpansionPanelsBuilder) Bind(name string, value string) (r *VExpansionPanelsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
