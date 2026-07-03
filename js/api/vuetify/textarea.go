package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTextareaBuilder struct {
	VTagBuilder[*VTextareaBuilder]
}

func VTextarea(children ...h.HTMLComponent) *VTextareaBuilder {
	return VTag(&VTextareaBuilder{}, "v-textarea", children...)
}

func (b *VTextareaBuilder) Flat(v bool) (r *VTextareaBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Reverse(v bool) (r *VTextareaBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Name(v string) (r *VTextareaBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VTextareaBuilder) Error(v bool) (r *VTextareaBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Label(v string) (r *VTextareaBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VTextareaBuilder) Theme(v string) (r *VTextareaBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTextareaBuilder) ID(v string) (r *VTextareaBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VTextareaBuilder) BaseColor(v string) (r *VTextareaBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VTextareaBuilder) BgColor(v string) (r *VTextareaBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VTextareaBuilder) Disabled(v bool) (r *VTextareaBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Density(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) MaxWidth(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) MinWidth(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Width(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Rounded(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Tile(v bool) (r *VTextareaBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Color(v string) (r *VTextareaBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTextareaBuilder) Variant(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) ModelValue(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Autofocus(v bool) (r *VTextareaBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Counter(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":counter", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Prefix(v string) (r *VTextareaBuilder) {
	b.Attr("prefix", v)
	return b
}

func (b *VTextareaBuilder) Placeholder(v string) (r *VTextareaBuilder) {
	b.Attr("placeholder", v)
	return b
}

func (b *VTextareaBuilder) PersistentPlaceholder(v bool) (r *VTextareaBuilder) {
	b.Attr(":persistent-placeholder", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) PersistentCounter(v bool) (r *VTextareaBuilder) {
	b.Attr(":persistent-counter", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Suffix(v string) (r *VTextareaBuilder) {
	b.Attr("suffix", v)
	return b
}

func (b *VTextareaBuilder) AppendIcon(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) CenterAffix(v bool) (r *VTextareaBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Glow(v bool) (r *VTextareaBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) IconColor(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) PrependIcon(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) HideSpinButtons(v bool) (r *VTextareaBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Hint(v string) (r *VTextareaBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VTextareaBuilder) PersistentHint(v bool) (r *VTextareaBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Messages(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Direction(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) ErrorMessages(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) MaxErrors(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Readonly(v bool) (r *VTextareaBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Rules(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) ValidateOn(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) ValidationValue(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Focused(v bool) (r *VTextareaBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) HideDetails(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) AppendInnerIcon(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":append-inner-icon", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Clearable(v bool) (r *VTextareaBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) ClearIcon(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) Active(v bool) (r *VTextareaBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Dirty(v bool) (r *VTextareaBuilder) {
	b.Attr(":dirty", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) PersistentClear(v bool) (r *VTextareaBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) PrependInnerIcon(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":prepend-inner-icon", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) SingleLine(v bool) (r *VTextareaBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Loading(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) CounterValue(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":counter-value", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) ModelModifiers(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":model-modifiers", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) AutoGrow(v bool) (r *VTextareaBuilder) {
	b.Attr(":auto-grow", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) NoResize(v bool) (r *VTextareaBuilder) {
	b.Attr(":no-resize", fmt.Sprint(v))
	return b
}

func (b *VTextareaBuilder) Rows(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":rows", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) MaxRows(v interface{}) (r *VTextareaBuilder) {
	b.Attr(":max-rows", h.JSONString(v))
	return b
}

func (b *VTextareaBuilder) On(name string, value string) (r *VTextareaBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTextareaBuilder) Bind(name string, value string) (r *VTextareaBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VTextareaBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VTextareaBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VTextareaBuilder) Slot(name string, child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VTextareaBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VTextareaBuilder) SetSlotDetails(child ...h.HTMLComponent) {
	b.SetSlot("details", child...)
}

func (b *VTextareaBuilder) SetScopedSlotDetails(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("details", scope, child...)
}

func (b *VTextareaBuilder) SlotDetails(child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetSlotDetails(child...)
	return b
}

func (b *VTextareaBuilder) ScopedSlotDetails(scope string, child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetScopedSlotDetails(scope, child...)
	return b
}

func (b *VTextareaBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VTextareaBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VTextareaBuilder) SlotLabel(child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VTextareaBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}

func (b *VTextareaBuilder) SetSlotClear(child ...h.HTMLComponent) {
	b.SetSlot("clear", child...)
}

func (b *VTextareaBuilder) SetScopedSlotClear(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("clear", scope, child...)
}

func (b *VTextareaBuilder) SlotClear(child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetSlotClear(child...)
	return b
}

func (b *VTextareaBuilder) ScopedSlotClear(scope string, child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetScopedSlotClear(scope, child...)
	return b
}

func (b *VTextareaBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VTextareaBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VTextareaBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VTextareaBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VTextareaBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VTextareaBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VTextareaBuilder) SlotAppend(child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VTextareaBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VTextareaBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VTextareaBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VTextareaBuilder) SlotMessage(child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VTextareaBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}

func (b *VTextareaBuilder) SetSlotPrependInner(child ...h.HTMLComponent) {
	b.SetSlot("prepend-inner", child...)
}

func (b *VTextareaBuilder) SetScopedSlotPrependInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend-inner", scope, child...)
}

func (b *VTextareaBuilder) SlotPrependInner(child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetSlotPrependInner(child...)
	return b
}

func (b *VTextareaBuilder) ScopedSlotPrependInner(scope string, child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetScopedSlotPrependInner(scope, child...)
	return b
}

func (b *VTextareaBuilder) SetSlotAppendInner(child ...h.HTMLComponent) {
	b.SetSlot("append-inner", child...)
}

func (b *VTextareaBuilder) SetScopedSlotAppendInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append-inner", scope, child...)
}

func (b *VTextareaBuilder) SlotAppendInner(child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetSlotAppendInner(child...)
	return b
}

func (b *VTextareaBuilder) ScopedSlotAppendInner(scope string, child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetScopedSlotAppendInner(scope, child...)
	return b
}

func (b *VTextareaBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VTextareaBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VTextareaBuilder) SlotLoader(child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VTextareaBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VTextareaBuilder) SetSlotCounter(child ...h.HTMLComponent) {
	b.SetSlot("counter", child...)
}

func (b *VTextareaBuilder) SetScopedSlotCounter(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("counter", scope, child...)
}

func (b *VTextareaBuilder) SlotCounter(child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetSlotCounter(child...)
	return b
}

func (b *VTextareaBuilder) ScopedSlotCounter(scope string, child ...h.HTMLComponent) (r *VTextareaBuilder) {
	b.SetScopedSlotCounter(scope, child...)
	return b
}
