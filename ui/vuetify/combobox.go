package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VComboboxBuilder struct {
	VTagBuilder[*VComboboxBuilder]
}

func VCombobox(children ...h.HTMLComponent) *VComboboxBuilder {
	return VTag(&VComboboxBuilder{}, "v-combobox", children...)
}

func (b *VComboboxBuilder) Label(v string) (r *VComboboxBuilder) {
	b.Attr("label", v)
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

func (b *VComboboxBuilder) Type(v string) (r *VComboboxBuilder) {
	b.Attr("type", v)
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

func (b *VComboboxBuilder) Reverse(v bool) (r *VComboboxBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Flat(v bool) (r *VComboboxBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
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

func (b *VComboboxBuilder) Eager(v bool) (r *VComboboxBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
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

func (b *VComboboxBuilder) BaseColor(v string) (r *VComboboxBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VComboboxBuilder) BgColor(v string) (r *VComboboxBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VComboboxBuilder) Disabled(v bool) (r *VComboboxBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) Multiple(v bool) (r *VComboboxBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
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

func (b *VComboboxBuilder) Items(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":items", h.JSONString(v))
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

func (b *VComboboxBuilder) ValueComparator(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
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

func (b *VComboboxBuilder) Menu(v bool) (r *VComboboxBuilder) {
	b.Attr(":menu", fmt.Sprint(v))
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

func (b *VComboboxBuilder) Id(v string) (r *VComboboxBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VComboboxBuilder) ModelValue(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Transition(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":transition", h.JSONString(v))
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

func (b *VComboboxBuilder) Prefix(v string) (r *VComboboxBuilder) {
	b.Attr("prefix", v)
	return b
}

func (b *VComboboxBuilder) Placeholder(v string) (r *VComboboxBuilder) {
	b.Attr("placeholder", v)
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

func (b *VComboboxBuilder) Role(v string) (r *VComboboxBuilder) {
	b.Attr("role", v)
	return b
}

func (b *VComboboxBuilder) AppendIcon(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) CenterAffix(v bool) (r *VComboboxBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) PrependIcon(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
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

func (b *VComboboxBuilder) Direction(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Error(v bool) (r *VComboboxBuilder) {
	b.Attr(":error", fmt.Sprint(v))
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

func (b *VComboboxBuilder) Readonly(v bool) (r *VComboboxBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
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

func (b *VComboboxBuilder) Clearable(v bool) (r *VComboboxBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VComboboxBuilder) ClearIcon(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VComboboxBuilder) Active(v bool) (r *VComboboxBuilder) {
	b.Attr(":active", fmt.Sprint(v))
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

func (b *VComboboxBuilder) Loading(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":loading", h.JSONString(v))
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

func (b *VComboboxBuilder) Delimiters(v interface{}) (r *VComboboxBuilder) {
	b.Attr(":delimiters", h.JSONString(v))
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
