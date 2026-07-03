package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VExpansionPanelTitleBuilder struct {
	VTagBuilder[*VExpansionPanelTitleBuilder]
}

func VExpansionPanelTitle(children ...h.HTMLComponent) *VExpansionPanelTitleBuilder {
	return VTag(&VExpansionPanelTitleBuilder{}, "v-expansion-panel-title", children...)
}

func (b *VExpansionPanelTitleBuilder) ExpandIcon(v interface{}) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":expand-icon", h.JSONString(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) CollapseIcon(v interface{}) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":collapse-icon", h.JSONString(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) Height(v interface{}) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) MaxHeight(v interface{}) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) MaxWidth(v interface{}) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) MinHeight(v interface{}) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) MinWidth(v interface{}) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) Width(v interface{}) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) Color(v string) (r *VExpansionPanelTitleBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VExpansionPanelTitleBuilder) Static(v bool) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":static", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelTitleBuilder) Readonly(v bool) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
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

func (b *VExpansionPanelTitleBuilder) Ripple(v interface{}) (r *VExpansionPanelTitleBuilder) {
	b.Attr(":ripple", h.JSONString(v))
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

func (b *VExpansionPanelTitleBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VExpansionPanelTitleBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VExpansionPanelTitleBuilder) Slot(name string, child ...h.HTMLComponent) (r *VExpansionPanelTitleBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VExpansionPanelTitleBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VExpansionPanelTitleBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VExpansionPanelTitleBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VExpansionPanelTitleBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VExpansionPanelTitleBuilder) SlotDefault(child ...h.HTMLComponent) (r *VExpansionPanelTitleBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VExpansionPanelTitleBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VExpansionPanelTitleBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VExpansionPanelTitleBuilder) SetSlotActions(child ...h.HTMLComponent) {
	b.SetSlot("actions", child...)
}

func (b *VExpansionPanelTitleBuilder) SetScopedSlotActions(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("actions", scope, child...)
}

func (b *VExpansionPanelTitleBuilder) SlotActions(child ...h.HTMLComponent) (r *VExpansionPanelTitleBuilder) {
	b.SetSlotActions(child...)
	return b
}

func (b *VExpansionPanelTitleBuilder) ScopedSlotActions(scope string, child ...h.HTMLComponent) (r *VExpansionPanelTitleBuilder) {
	b.SetScopedSlotActions(scope, child...)
	return b
}
