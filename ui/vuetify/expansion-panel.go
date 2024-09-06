package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VExpansionPanelBuilder struct {
	VTagBuilder[*VExpansionPanelBuilder]
}

func VExpansionPanel(children ...h.HTMLComponent) *VExpansionPanelBuilder {
	return VTag(&VExpansionPanelBuilder{}, "v-expansion-panel", children...)
}

func (b *VExpansionPanelBuilder) Title(v string) (r *VExpansionPanelBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VExpansionPanelBuilder) Text(v string) (r *VExpansionPanelBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VExpansionPanelBuilder) BgColor(v string) (r *VExpansionPanelBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VExpansionPanelBuilder) Elevation(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) Value(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) Disabled(v bool) (r *VExpansionPanelBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelBuilder) SelectedClass(v string) (r *VExpansionPanelBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VExpansionPanelBuilder) Rounded(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) Tile(v bool) (r *VExpansionPanelBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelBuilder) Tag(v string) (r *VExpansionPanelBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VExpansionPanelBuilder) Color(v string) (r *VExpansionPanelBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VExpansionPanelBuilder) ExpandIcon(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) CollapseIcon(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) HideActions(v bool) (r *VExpansionPanelBuilder) {
	b.Attr(":hide-actions", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelBuilder) Focusable(v bool) (r *VExpansionPanelBuilder) {
	b.Attr(":focusable", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelBuilder) Static(v bool) (r *VExpansionPanelBuilder) {
	b.Attr(":static", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelBuilder) Ripple(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) Readonly(v bool) (r *VExpansionPanelBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelBuilder) Eager(v bool) (r *VExpansionPanelBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelBuilder) On(name string, value string) (r *VExpansionPanelBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VExpansionPanelBuilder) Bind(name string, value string) (r *VExpansionPanelBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
