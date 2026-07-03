package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VNumberInputBuilder struct {
	VTagBuilder[*VNumberInputBuilder]
}

func VNumberInput(children ...h.HTMLComponent) *VNumberInputBuilder {
	return VTag(&VNumberInputBuilder{}, "v-number-input", children...)
}

func (b *VNumberInputBuilder) Flat(v bool) (r *VNumberInputBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Type(v string) (r *VNumberInputBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VNumberInputBuilder) Reverse(v bool) (r *VNumberInputBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Name(v string) (r *VNumberInputBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VNumberInputBuilder) Error(v bool) (r *VNumberInputBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Label(v string) (r *VNumberInputBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VNumberInputBuilder) Theme(v string) (r *VNumberInputBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VNumberInputBuilder) ID(v string) (r *VNumberInputBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VNumberInputBuilder) BaseColor(v string) (r *VNumberInputBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VNumberInputBuilder) BgColor(v string) (r *VNumberInputBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VNumberInputBuilder) Disabled(v bool) (r *VNumberInputBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Density(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) MaxWidth(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) MinWidth(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Width(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Rounded(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Tile(v bool) (r *VNumberInputBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Color(v string) (r *VNumberInputBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VNumberInputBuilder) Variant(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) ModelValue(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Autofocus(v bool) (r *VNumberInputBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Counter(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":counter", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Prefix(v string) (r *VNumberInputBuilder) {
	b.Attr("prefix", v)
	return b
}

func (b *VNumberInputBuilder) Placeholder(v string) (r *VNumberInputBuilder) {
	b.Attr("placeholder", v)
	return b
}

func (b *VNumberInputBuilder) PersistentPlaceholder(v bool) (r *VNumberInputBuilder) {
	b.Attr(":persistent-placeholder", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) PersistentCounter(v bool) (r *VNumberInputBuilder) {
	b.Attr(":persistent-counter", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Suffix(v string) (r *VNumberInputBuilder) {
	b.Attr("suffix", v)
	return b
}

func (b *VNumberInputBuilder) Role(v string) (r *VNumberInputBuilder) {
	b.Attr("role", v)
	return b
}

func (b *VNumberInputBuilder) AppendIcon(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) CenterAffix(v bool) (r *VNumberInputBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Glow(v bool) (r *VNumberInputBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) IconColor(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) PrependIcon(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) HideSpinButtons(v bool) (r *VNumberInputBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Hint(v string) (r *VNumberInputBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VNumberInputBuilder) PersistentHint(v bool) (r *VNumberInputBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Messages(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Direction(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) ErrorMessages(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) MaxErrors(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Readonly(v bool) (r *VNumberInputBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Rules(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) ValidateOn(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Focused(v bool) (r *VNumberInputBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) HideDetails(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) AppendInnerIcon(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":append-inner-icon", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Clearable(v bool) (r *VNumberInputBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) ClearIcon(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Active(v bool) (r *VNumberInputBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Dirty(v bool) (r *VNumberInputBuilder) {
	b.Attr(":dirty", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) PersistentClear(v bool) (r *VNumberInputBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) PrependInnerIcon(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":prepend-inner-icon", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) SingleLine(v bool) (r *VNumberInputBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) Loading(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) CounterValue(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":counter-value", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) ModelModifiers(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":model-modifiers", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Max(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Min(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Step(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":step", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Inset(v bool) (r *VNumberInputBuilder) {
	b.Attr(":inset", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) HideInput(v bool) (r *VNumberInputBuilder) {
	b.Attr(":hide-input", fmt.Sprint(v))
	return b
}

func (b *VNumberInputBuilder) ControlVariant(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":control-variant", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) Precision(v interface{}) (r *VNumberInputBuilder) {
	b.Attr(":precision", h.JSONString(v))
	return b
}

func (b *VNumberInputBuilder) On(name string, value string) (r *VNumberInputBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VNumberInputBuilder) Bind(name string, value string) (r *VNumberInputBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VNumberInputBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VNumberInputBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VNumberInputBuilder) Slot(name string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VNumberInputBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VNumberInputBuilder) SetSlotDetails(child ...h.HTMLComponent) {
	b.SetSlot("details", child...)
}

func (b *VNumberInputBuilder) SetScopedSlotDetails(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("details", scope, child...)
}

func (b *VNumberInputBuilder) SlotDetails(child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetSlotDetails(child...)
	return b
}

func (b *VNumberInputBuilder) ScopedSlotDetails(scope string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetScopedSlotDetails(scope, child...)
	return b
}

func (b *VNumberInputBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VNumberInputBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VNumberInputBuilder) SlotLabel(child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VNumberInputBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}

func (b *VNumberInputBuilder) SetSlotClear(child ...h.HTMLComponent) {
	b.SetSlot("clear", child...)
}

func (b *VNumberInputBuilder) SetScopedSlotClear(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("clear", scope, child...)
}

func (b *VNumberInputBuilder) SlotClear(child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetSlotClear(child...)
	return b
}

func (b *VNumberInputBuilder) ScopedSlotClear(scope string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetScopedSlotClear(scope, child...)
	return b
}

func (b *VNumberInputBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VNumberInputBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VNumberInputBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VNumberInputBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VNumberInputBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VNumberInputBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VNumberInputBuilder) SlotAppend(child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VNumberInputBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VNumberInputBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VNumberInputBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VNumberInputBuilder) SlotMessage(child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VNumberInputBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}

func (b *VNumberInputBuilder) SetSlotPrependInner(child ...h.HTMLComponent) {
	b.SetSlot("prepend-inner", child...)
}

func (b *VNumberInputBuilder) SetScopedSlotPrependInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend-inner", scope, child...)
}

func (b *VNumberInputBuilder) SlotPrependInner(child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetSlotPrependInner(child...)
	return b
}

func (b *VNumberInputBuilder) ScopedSlotPrependInner(scope string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetScopedSlotPrependInner(scope, child...)
	return b
}

func (b *VNumberInputBuilder) SetSlotAppendInner(child ...h.HTMLComponent) {
	b.SetSlot("append-inner", child...)
}

func (b *VNumberInputBuilder) SetScopedSlotAppendInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append-inner", scope, child...)
}

func (b *VNumberInputBuilder) SlotAppendInner(child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetSlotAppendInner(child...)
	return b
}

func (b *VNumberInputBuilder) ScopedSlotAppendInner(scope string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetScopedSlotAppendInner(scope, child...)
	return b
}

func (b *VNumberInputBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VNumberInputBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VNumberInputBuilder) SlotLoader(child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VNumberInputBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VNumberInputBuilder) SetSlotCounter(child ...h.HTMLComponent) {
	b.SetSlot("counter", child...)
}

func (b *VNumberInputBuilder) SetScopedSlotCounter(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("counter", scope, child...)
}

func (b *VNumberInputBuilder) SlotCounter(child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetSlotCounter(child...)
	return b
}

func (b *VNumberInputBuilder) ScopedSlotCounter(scope string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetScopedSlotCounter(scope, child...)
	return b
}

func (b *VNumberInputBuilder) SetSlotIncrement(child ...h.HTMLComponent) {
	b.SetSlot("increment", child...)
}

func (b *VNumberInputBuilder) SetScopedSlotIncrement(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("increment", scope, child...)
}

func (b *VNumberInputBuilder) SlotIncrement(child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetSlotIncrement(child...)
	return b
}

func (b *VNumberInputBuilder) ScopedSlotIncrement(scope string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetScopedSlotIncrement(scope, child...)
	return b
}

func (b *VNumberInputBuilder) SetSlotDecrement(child ...h.HTMLComponent) {
	b.SetSlot("decrement", child...)
}

func (b *VNumberInputBuilder) SetScopedSlotDecrement(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("decrement", scope, child...)
}

func (b *VNumberInputBuilder) SlotDecrement(child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetSlotDecrement(child...)
	return b
}

func (b *VNumberInputBuilder) ScopedSlotDecrement(scope string, child ...h.HTMLComponent) (r *VNumberInputBuilder) {
	b.SetScopedSlotDecrement(scope, child...)
	return b
}
