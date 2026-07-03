package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSelectionControlBuilder struct {
	VTagBuilder[*VSelectionControlBuilder]
}

func VSelectionControl(children ...h.HTMLComponent) *VSelectionControlBuilder {
	return VTag(&VSelectionControlBuilder{}, "v-selection-control", children...)
}

func (b *VSelectionControlBuilder) Type(v string) (r *VSelectionControlBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VSelectionControlBuilder) Name(v string) (r *VSelectionControlBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VSelectionControlBuilder) Error(v bool) (r *VSelectionControlBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlBuilder) Label(v string) (r *VSelectionControlBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VSelectionControlBuilder) Theme(v string) (r *VSelectionControlBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSelectionControlBuilder) ID(v string) (r *VSelectionControlBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VSelectionControlBuilder) Value(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) BaseColor(v string) (r *VSelectionControlBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VSelectionControlBuilder) Disabled(v bool) (r *VSelectionControlBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlBuilder) Multiple(v bool) (r *VSelectionControlBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlBuilder) Density(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) ValueComparator(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) Color(v string) (r *VSelectionControlBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSelectionControlBuilder) ModelValue(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) Readonly(v bool) (r *VSelectionControlBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlBuilder) Ripple(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) Inline(v bool) (r *VSelectionControlBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VSelectionControlBuilder) TrueValue(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":true-value", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) FalseValue(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":false-value", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) DefaultsTarget(v string) (r *VSelectionControlBuilder) {
	b.Attr("defaults-target", v)
	return b
}

func (b *VSelectionControlBuilder) FalseIcon(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":false-icon", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) TrueIcon(v interface{}) (r *VSelectionControlBuilder) {
	b.Attr(":true-icon", h.JSONString(v))
	return b
}

func (b *VSelectionControlBuilder) On(name string, value string) (r *VSelectionControlBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSelectionControlBuilder) Bind(name string, value string) (r *VSelectionControlBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VSelectionControlBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSelectionControlBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSelectionControlBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSelectionControlBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSelectionControlBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSelectionControlBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VSelectionControlBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VSelectionControlBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VSelectionControlBuilder) SlotDefault(child ...h.HTMLComponent) (r *VSelectionControlBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VSelectionControlBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VSelectionControlBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VSelectionControlBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VSelectionControlBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VSelectionControlBuilder) SlotLabel(child ...h.HTMLComponent) (r *VSelectionControlBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VSelectionControlBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VSelectionControlBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}

func (b *VSelectionControlBuilder) SetSlotInput(child ...h.HTMLComponent) {
	b.SetSlot("input", child...)
}

func (b *VSelectionControlBuilder) SetScopedSlotInput(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("input", scope, child...)
}

func (b *VSelectionControlBuilder) SlotInput(child ...h.HTMLComponent) (r *VSelectionControlBuilder) {
	b.SetSlotInput(child...)
	return b
}

func (b *VSelectionControlBuilder) ScopedSlotInput(scope string, child ...h.HTMLComponent) (r *VSelectionControlBuilder) {
	b.SetScopedSlotInput(scope, child...)
	return b
}
