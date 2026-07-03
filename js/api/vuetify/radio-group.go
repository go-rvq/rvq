package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VRadioGroupBuilder struct {
	VTagBuilder[*VRadioGroupBuilder]
}

func VRadioGroup(children ...h.HTMLComponent) *VRadioGroupBuilder {
	return VTag(&VRadioGroupBuilder{}, "v-radio-group", children...)
}

func (b *VRadioGroupBuilder) Type(v string) (r *VRadioGroupBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VRadioGroupBuilder) ModelValue(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Error(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) Density(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Height(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) MaxWidth(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) MinWidth(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Width(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Theme(v string) (r *VRadioGroupBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VRadioGroupBuilder) Color(v string) (r *VRadioGroupBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VRadioGroupBuilder) Name(v string) (r *VRadioGroupBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VRadioGroupBuilder) BaseColor(v string) (r *VRadioGroupBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VRadioGroupBuilder) PrependIcon(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) AppendIcon(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Readonly(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) Ripple(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Disabled(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) Inline(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) Label(v string) (r *VRadioGroupBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VRadioGroupBuilder) ID(v string) (r *VRadioGroupBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VRadioGroupBuilder) Direction(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) CenterAffix(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) Glow(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) IconColor(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) HideSpinButtons(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) Hint(v string) (r *VRadioGroupBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VRadioGroupBuilder) PersistentHint(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) Messages(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) ErrorMessages(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) MaxErrors(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Rules(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) ValidateOn(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) ValidationValue(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) Focused(v bool) (r *VRadioGroupBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VRadioGroupBuilder) HideDetails(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) DefaultsTarget(v string) (r *VRadioGroupBuilder) {
	b.Attr("defaults-target", v)
	return b
}

func (b *VRadioGroupBuilder) FalseIcon(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":false-icon", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) TrueIcon(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":true-icon", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) ValueComparator(v interface{}) (r *VRadioGroupBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VRadioGroupBuilder) On(name string, value string) (r *VRadioGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VRadioGroupBuilder) Bind(name string, value string) (r *VRadioGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VRadioGroupBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VRadioGroupBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VRadioGroupBuilder) Slot(name string, child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VRadioGroupBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VRadioGroupBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VRadioGroupBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VRadioGroupBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VRadioGroupBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VRadioGroupBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VRadioGroupBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VRadioGroupBuilder) SlotAppend(child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VRadioGroupBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VRadioGroupBuilder) SetSlotDetails(child ...h.HTMLComponent) {
	b.SetSlot("details", child...)
}

func (b *VRadioGroupBuilder) SetScopedSlotDetails(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("details", scope, child...)
}

func (b *VRadioGroupBuilder) SlotDetails(child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetSlotDetails(child...)
	return b
}

func (b *VRadioGroupBuilder) ScopedSlotDetails(scope string, child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetScopedSlotDetails(scope, child...)
	return b
}

func (b *VRadioGroupBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VRadioGroupBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VRadioGroupBuilder) SlotMessage(child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VRadioGroupBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}

func (b *VRadioGroupBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VRadioGroupBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VRadioGroupBuilder) SlotDefault(child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VRadioGroupBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VRadioGroupBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VRadioGroupBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VRadioGroupBuilder) SlotLabel(child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VRadioGroupBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VRadioGroupBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}
