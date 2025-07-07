package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VExpansionPanelTitleBuilder struct {
	VTagBuilder[*VExpansionPanelTitleBuilder]
}

func VExpansionPanelTitle(children ...h.HTMLComponent) *VExpansionPanelTitleBuilder {
	return VTag(&VExpansionPanelTitleBuilder{}, "v-expansion-panel-title", children...)
}

func (b *VExpansionPanelTitleBuilder) Color(v string) (r *VExpansionPanelTitleBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VExpansionPanelTitleBuilder) ExpandIcon(v interface{}) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) CollapseIcon(v interface{}) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) HideActions(v bool) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":hide-actions", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) Focusable(v bool) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":focusable", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) Static(v bool) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":static", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) Ripple(v interface{}) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) Readonly(v bool) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) On(name string, value string) (r *VExpansionPanelTitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VExpansionPanelTitleBuilder) Bind(name string, value string) (r *VExpansionPanelTitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
