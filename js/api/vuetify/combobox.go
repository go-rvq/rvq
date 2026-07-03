package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VComboboxBuilder struct {
	VTagBuilder[*VComboboxBuilder]
}

func VCombobox(children ...h.HTMLComponent) *VComboboxBuilder {
	return VTag(&VComboboxBuilder{}, "v-combobox", children...)
}

func (b *VComboboxBuilder) Flat(v bool) (r *VComboboxBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Type(v string) (r *VComboboxBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VComboboxBuilder) ModelValue(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Error(v bool) (r *VComboboxBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Reverse(v bool) (r *VComboboxBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Density(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) MaxWidth(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) MinWidth(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Width(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Rounded(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Tile(v bool) (r *VComboboxBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Theme(v string) (r *VComboboxBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VComboboxBuilder) Color(v string) (r *VComboboxBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VComboboxBuilder) Variant(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Name(v string) (r *VComboboxBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VComboboxBuilder) Delimiters(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":delimiters", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Active(v bool) (r *VComboboxBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) BaseColor(v string) (r *VComboboxBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VComboboxBuilder) PrependIcon(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) AppendIcon(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Readonly(v bool) (r *VComboboxBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Disabled(v bool) (r *VComboboxBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Loading(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Label(v string) (r *VComboboxBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VComboboxBuilder) Transition(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) BgColor(v string) (r *VComboboxBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VComboboxBuilder) Menu(v bool) (r *VComboboxBuilder) {
	b.Attr(":menu", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Multiple(v bool) (r *VComboboxBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Eager(v bool) (r *VComboboxBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) ID(v string) (r *VComboboxBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VComboboxBuilder) Prefix(v string) (r *VComboboxBuilder) {
	b.Attr("prefix", v)
	return b
}

func (b *VComboboxBuilder) Role(v string) (r *VComboboxBuilder) {
	b.Attr("role", v)
	return b
}

func (b *VComboboxBuilder) Items(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Direction(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Placeholder(v string) (r *VComboboxBuilder) {
	b.Attr("placeholder", v)
	return b
}

func (b *VComboboxBuilder) CenterAffix(v bool) (r *VComboboxBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Glow(v bool) (r *VComboboxBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) IconColor(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) HideSpinButtons(v bool) (r *VComboboxBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Hint(v string) (r *VComboboxBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VComboboxBuilder) PersistentHint(v bool) (r *VComboboxBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Messages(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) ErrorMessages(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) MaxErrors(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Rules(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) ValidateOn(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Focused(v bool) (r *VComboboxBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) HideDetails(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) ValueComparator(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) AutoSelectFirst(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":auto-select-first", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) ClearOnSelect(v bool) (r *VComboboxBuilder) {
	b.Attr(":clear-on-select", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) FilterMode(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":filter-mode", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) NoFilter(v bool) (r *VComboboxBuilder) {
	b.Attr(":no-filter", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) CustomFilter(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":custom-filter", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) CustomKeyFilter(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":custom-key-filter", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) FilterKeys(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":filter-keys", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Chips(v bool) (r *VComboboxBuilder) {
	b.Attr(":chips", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) ClosableChips(v bool) (r *VComboboxBuilder) {
	b.Attr(":closable-chips", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) CloseText(v string) (r *VComboboxBuilder) {
	b.Attr("close-text", v)
	return b
}

func (b *VComboboxBuilder) OpenText(v string) (r *VComboboxBuilder) {
	b.Attr("open-text", v)
	return b
}

func (b *VComboboxBuilder) HideNoData(v bool) (r *VComboboxBuilder) {
	b.Attr(":hide-no-data", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) HideSelected(v bool) (r *VComboboxBuilder) {
	b.Attr(":hide-selected", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) ListProps(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":list-props", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) ItemTitle(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":item-title", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) ItemValue(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":item-value", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) ItemChildren(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":item-children", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) ItemProps(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":item-props", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) ReturnObject(v bool) (r *VComboboxBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) MenuIcon(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":menu-icon", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) MenuProps(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":menu-props", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) NoDataText(v string) (r *VComboboxBuilder) {
	b.Attr("no-data-text", v)
	return b
}

func (b *VComboboxBuilder) OpenOnClear(v bool) (r *VComboboxBuilder) {
	b.Attr(":open-on-clear", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) ItemColor(v string) (r *VComboboxBuilder) {
	b.Attr("item-color", v)
	return b
}

func (b *VComboboxBuilder) Autofocus(v bool) (r *VComboboxBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Counter(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":counter", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) PersistentPlaceholder(v bool) (r *VComboboxBuilder) {
	b.Attr(":persistent-placeholder", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) PersistentCounter(v bool) (r *VComboboxBuilder) {
	b.Attr(":persistent-counter", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Suffix(v string) (r *VComboboxBuilder) {
	b.Attr("suffix", v)
	return b
}

func (b *VComboboxBuilder) Clearable(v bool) (r *VComboboxBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) ClearIcon(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) PersistentClear(v bool) (r *VComboboxBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) PrependInnerIcon(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":prepend-inner-icon", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) SingleLine(v bool) (r *VComboboxBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) CounterValue(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":counter-value", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) ModelModifiers(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":model-modifiers", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) On(name string, value string) (r *VComboboxBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VComboboxBuilder) Bind(name string, value string) (r *VComboboxBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VComboboxBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VComboboxBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VComboboxBuilder) Slot(name string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VComboboxBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VComboboxBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VComboboxBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VComboboxBuilder) SlotAppend(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VComboboxBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VComboboxBuilder) SlotLoader(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VComboboxBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VComboboxBuilder) SlotLabel(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotClear(child ...h.HTMLComponent) {
	b.SetSlot("clear", child...)
}

func (b *VComboboxBuilder) SetScopedSlotClear(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("clear", scope, child...)
}

func (b *VComboboxBuilder) SlotClear(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotClear(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotClear(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotClear(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotDetails(child ...h.HTMLComponent) {
	b.SetSlot("details", child...)
}

func (b *VComboboxBuilder) SetScopedSlotDetails(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("details", scope, child...)
}

func (b *VComboboxBuilder) SlotDetails(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotDetails(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotDetails(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotDetails(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VComboboxBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VComboboxBuilder) SlotMessage(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotPrependInner(child ...h.HTMLComponent) {
	b.SetSlot("prepend-inner", child...)
}

func (b *VComboboxBuilder) SetScopedSlotPrependInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend-inner", scope, child...)
}

func (b *VComboboxBuilder) SlotPrependInner(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotPrependInner(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotPrependInner(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotPrependInner(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotAppendInner(child ...h.HTMLComponent) {
	b.SetSlot("append-inner", child...)
}

func (b *VComboboxBuilder) SetScopedSlotAppendInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append-inner", scope, child...)
}

func (b *VComboboxBuilder) SlotAppendInner(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotAppendInner(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotAppendInner(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotAppendInner(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VComboboxBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VComboboxBuilder) SlotItem(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotChip(child ...h.HTMLComponent) {
	b.SetSlot("chip", child...)
}

func (b *VComboboxBuilder) SetScopedSlotChip(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("chip", scope, child...)
}

func (b *VComboboxBuilder) SlotChip(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotChip(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotChip(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotChip(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotSelection(child ...h.HTMLComponent) {
	b.SetSlot("selection", child...)
}

func (b *VComboboxBuilder) SetScopedSlotSelection(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("selection", scope, child...)
}

func (b *VComboboxBuilder) SlotSelection(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotSelection(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotSelection(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotSelection(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotPrependItem(child ...h.HTMLComponent) {
	b.SetSlot("prepend-item", child...)
}

func (b *VComboboxBuilder) SetScopedSlotPrependItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend-item", scope, child...)
}

func (b *VComboboxBuilder) SlotPrependItem(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotPrependItem(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotPrependItem(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotPrependItem(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotAppendItem(child ...h.HTMLComponent) {
	b.SetSlot("append-item", child...)
}

func (b *VComboboxBuilder) SetScopedSlotAppendItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append-item", scope, child...)
}

func (b *VComboboxBuilder) SlotAppendItem(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotAppendItem(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotAppendItem(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotAppendItem(scope, child...)
	return b
}

func (b *VComboboxBuilder) SetSlotNoData(child ...h.HTMLComponent) {
	b.SetSlot("no-data", child...)
}

func (b *VComboboxBuilder) SetScopedSlotNoData(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("no-data", scope, child...)
}

func (b *VComboboxBuilder) SlotNoData(child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetSlotNoData(child...)
	return b
}

func (b *VComboboxBuilder) ScopedSlotNoData(scope string, child ...h.HTMLComponent) (r *VComboboxBuilder) {
	b.SetScopedSlotNoData(scope, child...)
	return b
}
