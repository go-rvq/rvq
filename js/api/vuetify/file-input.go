package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VFileInputBuilder struct {
	VTagBuilder[*VFileInputBuilder]
}

func VFileInput(children ...h.HTMLComponent) *VFileInputBuilder {
	return VTag(&VFileInputBuilder{}, "v-file-input", children...)
}

func (b *VFileInputBuilder) Flat(v bool) (r *VFileInputBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Reverse(v bool) (r *VFileInputBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Name(v string) (r *VFileInputBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VFileInputBuilder) Error(v bool) (r *VFileInputBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Label(v string) (r *VFileInputBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VFileInputBuilder) Theme(v string) (r *VFileInputBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VFileInputBuilder) ID(v string) (r *VFileInputBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VFileInputBuilder) Chips(v bool) (r *VFileInputBuilder) {
	b.Attr(":chips", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) BaseColor(v string) (r *VFileInputBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VFileInputBuilder) BgColor(v string) (r *VFileInputBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VFileInputBuilder) Disabled(v bool) (r *VFileInputBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Multiple(v bool) (r *VFileInputBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Density(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) MaxWidth(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) MinWidth(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Width(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Rounded(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Tile(v bool) (r *VFileInputBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Color(v string) (r *VFileInputBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VFileInputBuilder) Variant(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) ModelValue(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Counter(v bool) (r *VFileInputBuilder) {
	b.Attr(":counter", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) AppendIcon(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) CenterAffix(v bool) (r *VFileInputBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Glow(v bool) (r *VFileInputBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) IconColor(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) PrependIcon(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) HideSpinButtons(v bool) (r *VFileInputBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Hint(v string) (r *VFileInputBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VFileInputBuilder) PersistentHint(v bool) (r *VFileInputBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Messages(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Direction(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) ErrorMessages(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) MaxErrors(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Readonly(v bool) (r *VFileInputBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Rules(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) ValidateOn(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) ValidationValue(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Focused(v bool) (r *VFileInputBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) HideDetails(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) AppendInnerIcon(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":append-inner-icon", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Clearable(v bool) (r *VFileInputBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) ClearIcon(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) Active(v bool) (r *VFileInputBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Dirty(v bool) (r *VFileInputBuilder) {
	b.Attr(":dirty", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) PersistentClear(v bool) (r *VFileInputBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) PrependInnerIcon(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":prepend-inner-icon", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) SingleLine(v bool) (r *VFileInputBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) Loading(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) CounterSizeString(v string) (r *VFileInputBuilder) {
	b.Attr("counter-size-string", v)
	return b
}

func (b *VFileInputBuilder) CounterString(v string) (r *VFileInputBuilder) {
	b.Attr("counter-string", v)
	return b
}

func (b *VFileInputBuilder) HideInput(v bool) (r *VFileInputBuilder) {
	b.Attr(":hide-input", fmt.Sprint(v))
	return b
}

func (b *VFileInputBuilder) ShowSize(v interface{}) (r *VFileInputBuilder) {
	b.Attr(":show-size", h.JSONString(v))
	return b
}

func (b *VFileInputBuilder) On(name string, value string) (r *VFileInputBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFileInputBuilder) Bind(name string, value string) (r *VFileInputBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VFileInputBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VFileInputBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VFileInputBuilder) Slot(name string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VFileInputBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VFileInputBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VFileInputBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VFileInputBuilder) SlotDefault(child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VFileInputBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VFileInputBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VFileInputBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VFileInputBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VFileInputBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VFileInputBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VFileInputBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VFileInputBuilder) SlotAppend(child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VFileInputBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VFileInputBuilder) SetSlotDetails(child ...h.HTMLComponent) {
	b.SetSlot("details", child...)
}

func (b *VFileInputBuilder) SetScopedSlotDetails(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("details", scope, child...)
}

func (b *VFileInputBuilder) SlotDetails(child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetSlotDetails(child...)
	return b
}

func (b *VFileInputBuilder) ScopedSlotDetails(scope string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetScopedSlotDetails(scope, child...)
	return b
}

func (b *VFileInputBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VFileInputBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VFileInputBuilder) SlotMessage(child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VFileInputBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}

func (b *VFileInputBuilder) SetSlotClear(child ...h.HTMLComponent) {
	b.SetSlot("clear", child...)
}

func (b *VFileInputBuilder) SetScopedSlotClear(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("clear", scope, child...)
}

func (b *VFileInputBuilder) SlotClear(child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetSlotClear(child...)
	return b
}

func (b *VFileInputBuilder) ScopedSlotClear(scope string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetScopedSlotClear(scope, child...)
	return b
}

func (b *VFileInputBuilder) SetSlotPrependInner(child ...h.HTMLComponent) {
	b.SetSlot("prepend-inner", child...)
}

func (b *VFileInputBuilder) SetScopedSlotPrependInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend-inner", scope, child...)
}

func (b *VFileInputBuilder) SlotPrependInner(child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetSlotPrependInner(child...)
	return b
}

func (b *VFileInputBuilder) ScopedSlotPrependInner(scope string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetScopedSlotPrependInner(scope, child...)
	return b
}

func (b *VFileInputBuilder) SetSlotAppendInner(child ...h.HTMLComponent) {
	b.SetSlot("append-inner", child...)
}

func (b *VFileInputBuilder) SetScopedSlotAppendInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append-inner", scope, child...)
}

func (b *VFileInputBuilder) SlotAppendInner(child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetSlotAppendInner(child...)
	return b
}

func (b *VFileInputBuilder) ScopedSlotAppendInner(scope string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetScopedSlotAppendInner(scope, child...)
	return b
}

func (b *VFileInputBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VFileInputBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VFileInputBuilder) SlotLabel(child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VFileInputBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}

func (b *VFileInputBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VFileInputBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VFileInputBuilder) SlotLoader(child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VFileInputBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VFileInputBuilder) SetSlotCounter(child ...h.HTMLComponent) {
	b.SetSlot("counter", child...)
}

func (b *VFileInputBuilder) SetScopedSlotCounter(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("counter", scope, child...)
}

func (b *VFileInputBuilder) SlotCounter(child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetSlotCounter(child...)
	return b
}

func (b *VFileInputBuilder) ScopedSlotCounter(scope string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetScopedSlotCounter(scope, child...)
	return b
}

func (b *VFileInputBuilder) SetSlotSelection(child ...h.HTMLComponent) {
	b.SetSlot("selection", child...)
}

func (b *VFileInputBuilder) SetScopedSlotSelection(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("selection", scope, child...)
}

func (b *VFileInputBuilder) SlotSelection(child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetSlotSelection(child...)
	return b
}

func (b *VFileInputBuilder) ScopedSlotSelection(scope string, child ...h.HTMLComponent) (r *VFileInputBuilder) {
	b.SetScopedSlotSelection(scope, child...)
	return b
}
