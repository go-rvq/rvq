package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VAutocompleteBuilder struct {
	VTagBuilder[*VAutocompleteBuilder]
}

func VAutocomplete(children ...h.HTMLComponent) *VAutocompleteBuilder {
	return VTag(&VAutocompleteBuilder{}, "v-autocomplete", children...)
}

func (b *VAutocompleteBuilder) Flat(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Search(v string) (r *VAutocompleteBuilder) {
	b.Attr("search", v)
	return b
}

func (b *VAutocompleteBuilder) Type(v string) (r *VAutocompleteBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VAutocompleteBuilder) Reverse(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Name(v string) (r *VAutocompleteBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VAutocompleteBuilder) Error(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Label(v string) (r *VAutocompleteBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VAutocompleteBuilder) Menu(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":menu", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Theme(v string) (r *VAutocompleteBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VAutocompleteBuilder) Items(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) ID(v string) (r *VAutocompleteBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VAutocompleteBuilder) AutoSelectFirst(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":auto-select-first", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) ClearOnSelect(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":clear-on-select", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) FilterMode(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":filter-mode", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) NoFilter(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":no-filter", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) CustomFilter(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":custom-filter", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) CustomKeyFilter(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":custom-key-filter", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) FilterKeys(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":filter-keys", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) Chips(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":chips", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) ClosableChips(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":closable-chips", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) CloseText(v string) (r *VAutocompleteBuilder) {
	b.Attr("close-text", v)
	return b
}

func (b *VAutocompleteBuilder) OpenText(v string) (r *VAutocompleteBuilder) {
	b.Attr("open-text", v)
	return b
}

func (b *VAutocompleteBuilder) Eager(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) HideNoData(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":hide-no-data", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) HideSelected(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":hide-selected", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) ListProps(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":list-props", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) BaseColor(v string) (r *VAutocompleteBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VAutocompleteBuilder) BgColor(v string) (r *VAutocompleteBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VAutocompleteBuilder) Disabled(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Multiple(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Density(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) MaxWidth(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) MinWidth(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) Width(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) ItemTitle(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":item-title", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) ItemValue(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":item-value", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) ItemChildren(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":item-children", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) ItemProps(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":item-props", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) ReturnObject(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) ValueComparator(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) Rounded(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) Tile(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Color(v string) (r *VAutocompleteBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VAutocompleteBuilder) Variant(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) MenuIcon(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":menu-icon", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) MenuProps(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":menu-props", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) ModelValue(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) Transition(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) NoDataText(v string) (r *VAutocompleteBuilder) {
	b.Attr("no-data-text", v)
	return b
}

func (b *VAutocompleteBuilder) OpenOnClear(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":open-on-clear", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) ItemColor(v string) (r *VAutocompleteBuilder) {
	b.Attr("item-color", v)
	return b
}

func (b *VAutocompleteBuilder) Autofocus(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Counter(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":counter", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) Prefix(v string) (r *VAutocompleteBuilder) {
	b.Attr("prefix", v)
	return b
}

func (b *VAutocompleteBuilder) Placeholder(v string) (r *VAutocompleteBuilder) {
	b.Attr("placeholder", v)
	return b
}

func (b *VAutocompleteBuilder) PersistentPlaceholder(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":persistent-placeholder", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) PersistentCounter(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":persistent-counter", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Suffix(v string) (r *VAutocompleteBuilder) {
	b.Attr("suffix", v)
	return b
}

func (b *VAutocompleteBuilder) Role(v string) (r *VAutocompleteBuilder) {
	b.Attr("role", v)
	return b
}

func (b *VAutocompleteBuilder) AppendIcon(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) CenterAffix(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Glow(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) IconColor(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) PrependIcon(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) HideSpinButtons(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Hint(v string) (r *VAutocompleteBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VAutocompleteBuilder) PersistentHint(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Messages(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) Direction(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) ErrorMessages(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) MaxErrors(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) Readonly(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Rules(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) ValidateOn(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) Focused(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) HideDetails(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) Clearable(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) ClearIcon(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) Active(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) PersistentClear(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) PrependInnerIcon(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":prepend-inner-icon", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) SingleLine(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Loading(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) CounterValue(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":counter-value", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) ModelModifiers(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":model-modifiers", h.JSONString(v))
	return b
}

func (b *VAutocompleteBuilder) On(name string, value string) (r *VAutocompleteBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAutocompleteBuilder) Bind(name string, value string) (r *VAutocompleteBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VAutocompleteBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VAutocompleteBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VAutocompleteBuilder) Slot(name string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotDetails(child ...h.HTMLComponent) {
	b.SetSlot("details", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotDetails(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("details", scope, child...)
}

func (b *VAutocompleteBuilder) SlotDetails(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotDetails(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotDetails(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotDetails(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VAutocompleteBuilder) SlotLabel(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotClear(child ...h.HTMLComponent) {
	b.SetSlot("clear", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotClear(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("clear", scope, child...)
}

func (b *VAutocompleteBuilder) SlotClear(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotClear(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotClear(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotClear(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VAutocompleteBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VAutocompleteBuilder) SlotAppend(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VAutocompleteBuilder) SlotMessage(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotPrependInner(child ...h.HTMLComponent) {
	b.SetSlot("prepend-inner", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotPrependInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend-inner", scope, child...)
}

func (b *VAutocompleteBuilder) SlotPrependInner(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotPrependInner(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotPrependInner(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotPrependInner(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotAppendInner(child ...h.HTMLComponent) {
	b.SetSlot("append-inner", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotAppendInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append-inner", scope, child...)
}

func (b *VAutocompleteBuilder) SlotAppendInner(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotAppendInner(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotAppendInner(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotAppendInner(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VAutocompleteBuilder) SlotLoader(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VAutocompleteBuilder) SlotItem(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotChip(child ...h.HTMLComponent) {
	b.SetSlot("chip", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotChip(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("chip", scope, child...)
}

func (b *VAutocompleteBuilder) SlotChip(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotChip(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotChip(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotChip(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotSelection(child ...h.HTMLComponent) {
	b.SetSlot("selection", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotSelection(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("selection", scope, child...)
}

func (b *VAutocompleteBuilder) SlotSelection(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotSelection(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotSelection(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotSelection(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotPrependItem(child ...h.HTMLComponent) {
	b.SetSlot("prepend-item", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotPrependItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend-item", scope, child...)
}

func (b *VAutocompleteBuilder) SlotPrependItem(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotPrependItem(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotPrependItem(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotPrependItem(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotAppendItem(child ...h.HTMLComponent) {
	b.SetSlot("append-item", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotAppendItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append-item", scope, child...)
}

func (b *VAutocompleteBuilder) SlotAppendItem(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotAppendItem(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotAppendItem(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotAppendItem(scope, child...)
	return b
}

func (b *VAutocompleteBuilder) SetSlotNoData(child ...h.HTMLComponent) {
	b.SetSlot("no-data", child...)
}

func (b *VAutocompleteBuilder) SetScopedSlotNoData(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("no-data", scope, child...)
}

func (b *VAutocompleteBuilder) SlotNoData(child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetSlotNoData(child...)
	return b
}

func (b *VAutocompleteBuilder) ScopedSlotNoData(scope string, child ...h.HTMLComponent) (r *VAutocompleteBuilder) {
	b.SetScopedSlotNoData(scope, child...)
	return b
}
