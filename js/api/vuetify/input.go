package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VInputBuilder struct {
	VTagBuilder[*VInputBuilder]
}

func VInput(children ...h.HTMLComponent) *VInputBuilder {
	return VTag(&VInputBuilder{}, "v-input", children...)
}

func (b *VInputBuilder) Name(v string) (r *VInputBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VInputBuilder) Error(v bool) (r *VInputBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) Label(v string) (r *VInputBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VInputBuilder) Theme(v string) (r *VInputBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VInputBuilder) ID(v string) (r *VInputBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VInputBuilder) BaseColor(v string) (r *VInputBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VInputBuilder) Disabled(v bool) (r *VInputBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) Density(v interface{}) (r *VInputBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VInputBuilder) MaxWidth(v interface{}) (r *VInputBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VInputBuilder) MinWidth(v interface{}) (r *VInputBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VInputBuilder) Width(v interface{}) (r *VInputBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VInputBuilder) Color(v string) (r *VInputBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VInputBuilder) ModelValue(v interface{}) (r *VInputBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VInputBuilder) AppendIcon(v interface{}) (r *VInputBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VInputBuilder) CenterAffix(v bool) (r *VInputBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) Glow(v bool) (r *VInputBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) IconColor(v interface{}) (r *VInputBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VInputBuilder) PrependIcon(v interface{}) (r *VInputBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VInputBuilder) HideSpinButtons(v bool) (r *VInputBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) Hint(v string) (r *VInputBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VInputBuilder) PersistentHint(v bool) (r *VInputBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) Messages(v interface{}) (r *VInputBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VInputBuilder) Direction(v interface{}) (r *VInputBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VInputBuilder) ErrorMessages(v interface{}) (r *VInputBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VInputBuilder) MaxErrors(v interface{}) (r *VInputBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VInputBuilder) Readonly(v bool) (r *VInputBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) Rules(v interface{}) (r *VInputBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VInputBuilder) ValidateOn(v interface{}) (r *VInputBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VInputBuilder) ValidationValue(v interface{}) (r *VInputBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VInputBuilder) Focused(v bool) (r *VInputBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VInputBuilder) HideDetails(v interface{}) (r *VInputBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VInputBuilder) On(name string, value string) (r *VInputBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VInputBuilder) Bind(name string, value string) (r *VInputBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VInputBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VInputBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VInputBuilder) Slot(name string, child ...h.HTMLComponent) (r *VInputBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VInputBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VInputBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VInputBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VInputBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VInputBuilder) SlotDefault(child ...h.HTMLComponent) (r *VInputBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VInputBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VInputBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VInputBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VInputBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VInputBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VInputBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VInputBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VInputBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VInputBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VInputBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VInputBuilder) SlotAppend(child ...h.HTMLComponent) (r *VInputBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VInputBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VInputBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VInputBuilder) SetSlotDetails(child ...h.HTMLComponent) {
	b.SetSlot("details", child...)
}

func (b *VInputBuilder) SetScopedSlotDetails(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("details", scope, child...)
}

func (b *VInputBuilder) SlotDetails(child ...h.HTMLComponent) (r *VInputBuilder) {
	b.SetSlotDetails(child...)
	return b
}

func (b *VInputBuilder) ScopedSlotDetails(scope string, child ...h.HTMLComponent) (r *VInputBuilder) {
	b.SetScopedSlotDetails(scope, child...)
	return b
}

func (b *VInputBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VInputBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VInputBuilder) SlotMessage(child ...h.HTMLComponent) (r *VInputBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VInputBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VInputBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}
