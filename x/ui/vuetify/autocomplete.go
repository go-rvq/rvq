package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VAutocompleteBuilder struct {
	VTagBuilder[*VAutocompleteBuilder]
}

func (b *VAutocompleteBuilder) Label(v string) (r *VAutocompleteBuilder) {
	b.Attr("label", v)
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

func (b *VAutocompleteBuilder) Search(v string) (r *VAutocompleteBuilder) {
	b.Attr("search", v)
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

func (b *VAutocompleteBuilder) Reverse(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VAutocompleteBuilder) Flat(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
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

func (b *VAutocompleteBuilder) Type(v string) (r *VAutocompleteBuilder) {
	b.Attr("type", v)
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

func (b *VAutocompleteBuilder) Items(v interface{}) (r *VAutocompleteBuilder) {
	b.Attr(":items", h.JSONString(v))
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

func (b *VAutocompleteBuilder) Theme(v string) (r *VAutocompleteBuilder) {
	b.Attr("theme", v)
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

func (b *VAutocompleteBuilder) Name(v string) (r *VAutocompleteBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VAutocompleteBuilder) Menu(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":menu", fmt.Sprint(v))
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

func (b *VAutocompleteBuilder) Id(v string) (r *VAutocompleteBuilder) {
	b.Attr("id", v)
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

func (b *VAutocompleteBuilder) Error(v bool) (r *VAutocompleteBuilder) {
	b.Attr(":error", fmt.Sprint(v))
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
