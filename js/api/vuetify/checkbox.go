package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCheckboxBuilder struct {
	VTagBuilder[*VCheckboxBuilder]
}

func VCheckbox(children ...h.HTMLComponent) *VCheckboxBuilder {
	return VTag(&VCheckboxBuilder{}, "v-checkbox", children...)
}

func (b *VCheckboxBuilder) Type(v string) (r *VCheckboxBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VCheckboxBuilder) ModelValue(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) Error(v bool) (r *VCheckboxBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBuilder) Density(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) MaxWidth(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) MinWidth(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) Width(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) Theme(v string) (r *VCheckboxBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VCheckboxBuilder) Color(v string) (r *VCheckboxBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VCheckboxBuilder) Name(v string) (r *VCheckboxBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VCheckboxBuilder) BaseColor(v string) (r *VCheckboxBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VCheckboxBuilder) PrependIcon(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) AppendIcon(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) Readonly(v bool) (r *VCheckboxBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBuilder) Ripple(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) Value(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) Disabled(v bool) (r *VCheckboxBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBuilder) Label(v string) (r *VCheckboxBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VCheckboxBuilder) Multiple(v bool) (r *VCheckboxBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBuilder) ID(v string) (r *VCheckboxBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VCheckboxBuilder) Direction(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) CenterAffix(v bool) (r *VCheckboxBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBuilder) Glow(v bool) (r *VCheckboxBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBuilder) IconColor(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) HideSpinButtons(v bool) (r *VCheckboxBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBuilder) Hint(v string) (r *VCheckboxBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VCheckboxBuilder) PersistentHint(v bool) (r *VCheckboxBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBuilder) Messages(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) ErrorMessages(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) MaxErrors(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) Rules(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) ValidateOn(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) ValidationValue(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) Focused(v bool) (r *VCheckboxBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBuilder) HideDetails(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) Indeterminate(v bool) (r *VCheckboxBuilder) {
	b.Attr(":indeterminate", fmt.Sprint(v))
	return b
}

func (b *VCheckboxBuilder) IndeterminateIcon(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":indeterminate-icon", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) TrueValue(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":true-value", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) FalseValue(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":false-value", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) DefaultsTarget(v string) (r *VCheckboxBuilder) {
	b.Attr("defaults-target", v)
	return b
}

func (b *VCheckboxBuilder) FalseIcon(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":false-icon", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) TrueIcon(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":true-icon", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) ValueComparator(v interface{}) (r *VCheckboxBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VCheckboxBuilder) On(name string, value string) (r *VCheckboxBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCheckboxBuilder) Bind(name string, value string) (r *VCheckboxBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VCheckboxBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCheckboxBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCheckboxBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCheckboxBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCheckboxBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VCheckboxBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VCheckboxBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VCheckboxBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VCheckboxBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VCheckboxBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VCheckboxBuilder) SlotAppend(child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VCheckboxBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VCheckboxBuilder) SetSlotDetails(child ...h.HTMLComponent) {
	b.SetSlot("details", child...)
}

func (b *VCheckboxBuilder) SetScopedSlotDetails(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("details", scope, child...)
}

func (b *VCheckboxBuilder) SlotDetails(child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetSlotDetails(child...)
	return b
}

func (b *VCheckboxBuilder) ScopedSlotDetails(scope string, child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetScopedSlotDetails(scope, child...)
	return b
}

func (b *VCheckboxBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VCheckboxBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VCheckboxBuilder) SlotMessage(child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VCheckboxBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}

func (b *VCheckboxBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCheckboxBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCheckboxBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCheckboxBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VCheckboxBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VCheckboxBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VCheckboxBuilder) SlotLabel(child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VCheckboxBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}

func (b *VCheckboxBuilder) SetSlotInput(child ...h.HTMLComponent) {
	b.SetSlot("input", child...)
}

func (b *VCheckboxBuilder) SetScopedSlotInput(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("input", scope, child...)
}

func (b *VCheckboxBuilder) SlotInput(child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetSlotInput(child...)
	return b
}

func (b *VCheckboxBuilder) ScopedSlotInput(scope string, child ...h.HTMLComponent) (r *VCheckboxBuilder) {
	b.SetScopedSlotInput(scope, child...)
	return b
}
