package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VExpansionPanelBuilder struct {
	VTagBuilder[*VExpansionPanelBuilder]
}

func VExpansionPanel(children ...h.HTMLComponent) *VExpansionPanelBuilder {
	return VTag(&VExpansionPanelBuilder{}, "v-expansion-panel", children...)
}

func (b *VExpansionPanelBuilder) Tag(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) Title(v string) (r *VExpansionPanelBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VExpansionPanelBuilder) Value(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) Text(v string) (r *VExpansionPanelBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VExpansionPanelBuilder) Eager(v bool) (r *VExpansionPanelBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelBuilder) BgColor(v string) (r *VExpansionPanelBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VExpansionPanelBuilder) Disabled(v bool) (r *VExpansionPanelBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
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

func (b *VExpansionPanelBuilder) Height(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) MaxHeight(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) MaxWidth(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) MinHeight(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) MinWidth(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) Width(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VExpansionPanelBuilder) Elevation(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":elevation", h.JSONString(v))
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

func (b *VExpansionPanelBuilder) Color(v string) (r *VExpansionPanelBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VExpansionPanelBuilder) Static(v bool) (r *VExpansionPanelBuilder) {
	b.Attr(":static", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelBuilder) Readonly(v bool) (r *VExpansionPanelBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelBuilder) HideActions(v bool) (r *VExpansionPanelBuilder) {
	b.Attr(":hide-actions", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelBuilder) SelectedClass(v string) (r *VExpansionPanelBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VExpansionPanelBuilder) Focusable(v bool) (r *VExpansionPanelBuilder) {
	b.Attr(":focusable", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelBuilder) Ripple(v interface{}) (r *VExpansionPanelBuilder) {
	b.Attr(":ripple", h.JSONString(v))
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

func (b *VExpansionPanelBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VExpansionPanelBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VExpansionPanelBuilder) Slot(name string, child ...h.HTMLComponent) (r *VExpansionPanelBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VExpansionPanelBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VExpansionPanelBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VExpansionPanelBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VExpansionPanelBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VExpansionPanelBuilder) SlotDefault(child ...h.HTMLComponent) (r *VExpansionPanelBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VExpansionPanelBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VExpansionPanelBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VExpansionPanelBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VExpansionPanelBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VExpansionPanelBuilder) SlotTitle(child ...h.HTMLComponent) (r *VExpansionPanelBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VExpansionPanelBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VExpansionPanelBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VExpansionPanelBuilder) SetSlotText(child ...h.HTMLComponent) {
	b.SetSlot("text", child...)
}

func (b *VExpansionPanelBuilder) SetScopedSlotText(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("text", scope, child...)
}

func (b *VExpansionPanelBuilder) SlotText(child ...h.HTMLComponent) (r *VExpansionPanelBuilder) {
	b.SetSlotText(child...)
	return b
}

func (b *VExpansionPanelBuilder) ScopedSlotText(scope string, child ...h.HTMLComponent) (r *VExpansionPanelBuilder) {
	b.SetScopedSlotText(scope, child...)
	return b
}
