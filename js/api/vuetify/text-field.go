package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTextFieldBuilder struct {
	VTagBuilder[*VTextFieldBuilder]
}

func VTextField(children ...h.HTMLComponent) *VTextFieldBuilder {
	return VTag(&VTextFieldBuilder{}, "v-text-field", children...)
}

func (b *VTextFieldBuilder) Flat(v bool) (r *VTextFieldBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Type(v string) (r *VTextFieldBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VTextFieldBuilder) ModelValue(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Error(v bool) (r *VTextFieldBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Reverse(v bool) (r *VTextFieldBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Density(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) MaxWidth(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) MinWidth(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Width(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Rounded(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Tile(v bool) (r *VTextFieldBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Theme(v string) (r *VTextFieldBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTextFieldBuilder) Color(v string) (r *VTextFieldBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTextFieldBuilder) Variant(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Name(v string) (r *VTextFieldBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VTextFieldBuilder) Active(v bool) (r *VTextFieldBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) BaseColor(v string) (r *VTextFieldBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VTextFieldBuilder) PrependIcon(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) AppendIcon(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Readonly(v bool) (r *VTextFieldBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Disabled(v bool) (r *VTextFieldBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Loading(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Label(v string) (r *VTextFieldBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VTextFieldBuilder) BgColor(v string) (r *VTextFieldBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VTextFieldBuilder) ID(v string) (r *VTextFieldBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VTextFieldBuilder) Prefix(v string) (r *VTextFieldBuilder) {
	b.Attr("prefix", v)
	return b
}

func (b *VTextFieldBuilder) Role(v string) (r *VTextFieldBuilder) {
	b.Attr("role", v)
	return b
}

func (b *VTextFieldBuilder) Direction(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Placeholder(v string) (r *VTextFieldBuilder) {
	b.Attr("placeholder", v)
	return b
}

func (b *VTextFieldBuilder) CenterAffix(v bool) (r *VTextFieldBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Glow(v bool) (r *VTextFieldBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) IconColor(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) HideSpinButtons(v bool) (r *VTextFieldBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Hint(v string) (r *VTextFieldBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VTextFieldBuilder) PersistentHint(v bool) (r *VTextFieldBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Messages(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) ErrorMessages(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) MaxErrors(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Rules(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) ValidateOn(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) ValidationValue(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Focused(v bool) (r *VTextFieldBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) HideDetails(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Autofocus(v bool) (r *VTextFieldBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Counter(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":counter", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) PersistentPlaceholder(v bool) (r *VTextFieldBuilder) {
	b.Attr(":persistent-placeholder", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) PersistentCounter(v bool) (r *VTextFieldBuilder) {
	b.Attr(":persistent-counter", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) Suffix(v string) (r *VTextFieldBuilder) {
	b.Attr("suffix", v)
	return b
}

func (b *VTextFieldBuilder) AppendInnerIcon(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":append-inner-icon", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Clearable(v bool) (r *VTextFieldBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) ClearIcon(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) Dirty(v bool) (r *VTextFieldBuilder) {
	b.Attr(":dirty", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) PersistentClear(v bool) (r *VTextFieldBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) PrependInnerIcon(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":prepend-inner-icon", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) SingleLine(v bool) (r *VTextFieldBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VTextFieldBuilder) CounterValue(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":counter-value", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) ModelModifiers(v interface{}) (r *VTextFieldBuilder) {
	b.Attr(":model-modifiers", h.JSONString(v))
	return b
}

func (b *VTextFieldBuilder) On(name string, value string) (r *VTextFieldBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTextFieldBuilder) Bind(name string, value string) (r *VTextFieldBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VTextFieldBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VTextFieldBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VTextFieldBuilder) Slot(name string, child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VTextFieldBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VTextFieldBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VTextFieldBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VTextFieldBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VTextFieldBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VTextFieldBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VTextFieldBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VTextFieldBuilder) SlotAppend(child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VTextFieldBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VTextFieldBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VTextFieldBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VTextFieldBuilder) SlotLoader(child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VTextFieldBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VTextFieldBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VTextFieldBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VTextFieldBuilder) SlotLabel(child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VTextFieldBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}

func (b *VTextFieldBuilder) SetSlotClear(child ...h.HTMLComponent) {
	b.SetSlot("clear", child...)
}

func (b *VTextFieldBuilder) SetScopedSlotClear(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("clear", scope, child...)
}

func (b *VTextFieldBuilder) SlotClear(child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetSlotClear(child...)
	return b
}

func (b *VTextFieldBuilder) ScopedSlotClear(scope string, child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetScopedSlotClear(scope, child...)
	return b
}

func (b *VTextFieldBuilder) SetSlotDetails(child ...h.HTMLComponent) {
	b.SetSlot("details", child...)
}

func (b *VTextFieldBuilder) SetScopedSlotDetails(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("details", scope, child...)
}

func (b *VTextFieldBuilder) SlotDetails(child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetSlotDetails(child...)
	return b
}

func (b *VTextFieldBuilder) ScopedSlotDetails(scope string, child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetScopedSlotDetails(scope, child...)
	return b
}

func (b *VTextFieldBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VTextFieldBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VTextFieldBuilder) SlotMessage(child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VTextFieldBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}

func (b *VTextFieldBuilder) SetSlotPrependInner(child ...h.HTMLComponent) {
	b.SetSlot("prepend-inner", child...)
}

func (b *VTextFieldBuilder) SetScopedSlotPrependInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend-inner", scope, child...)
}

func (b *VTextFieldBuilder) SlotPrependInner(child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetSlotPrependInner(child...)
	return b
}

func (b *VTextFieldBuilder) ScopedSlotPrependInner(scope string, child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetScopedSlotPrependInner(scope, child...)
	return b
}

func (b *VTextFieldBuilder) SetSlotAppendInner(child ...h.HTMLComponent) {
	b.SetSlot("append-inner", child...)
}

func (b *VTextFieldBuilder) SetScopedSlotAppendInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append-inner", scope, child...)
}

func (b *VTextFieldBuilder) SlotAppendInner(child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetSlotAppendInner(child...)
	return b
}

func (b *VTextFieldBuilder) ScopedSlotAppendInner(scope string, child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetScopedSlotAppendInner(scope, child...)
	return b
}

func (b *VTextFieldBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VTextFieldBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VTextFieldBuilder) SlotDefault(child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VTextFieldBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VTextFieldBuilder) SetSlotCounter(child ...h.HTMLComponent) {
	b.SetSlot("counter", child...)
}

func (b *VTextFieldBuilder) SetScopedSlotCounter(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("counter", scope, child...)
}

func (b *VTextFieldBuilder) SlotCounter(child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetSlotCounter(child...)
	return b
}

func (b *VTextFieldBuilder) ScopedSlotCounter(scope string, child ...h.HTMLComponent) (r *VTextFieldBuilder) {
	b.SetScopedSlotCounter(scope, child...)
	return b
}
