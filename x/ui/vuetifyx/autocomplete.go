package vuetifyx

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

func (b *VXAutocompleteBuilder) AllowOverflow(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":allow-overflow", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) AppendIcon(v string) (r *VXAutocompleteBuilder) {
	b.Attr("append-icon", v)
	return b
}

func (b *VXAutocompleteBuilder) AppendOuterIcon(v string) (r *VXAutocompleteBuilder) {
	b.Attr("append-outer-icon", v)
	return b
}

func (b *VXAutocompleteBuilder) Attach(v interface{}) (r *VXAutocompleteBuilder) {
	b.Attr(":attach", h.JSONString(v))
	return b
}

func (b *VXAutocompleteBuilder) AutoSelectFirst(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":auto-select-first", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Autofocus(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) BackgroundColor(v string) (r *VXAutocompleteBuilder) {
	b.Attr("background-color", v)
	return b
}

func (b *VXAutocompleteBuilder) CacheItems(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":cache-items", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Chips(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":chips", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) ClearIcon(v string) (r *VXAutocompleteBuilder) {
	b.Attr("clear-icon", v)
	return b
}

func (b *VXAutocompleteBuilder) Clearable(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Color(v string) (r *VXAutocompleteBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VXAutocompleteBuilder) Counter(v int) (r *VXAutocompleteBuilder) {
	b.Attr(":counter", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) CounterValue(v interface{}) (r *VXAutocompleteBuilder) {
	b.Attr(":counter-value", h.JSONString(v))
	return b
}

func (b *VXAutocompleteBuilder) Dark(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":dark", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) DeletableChips(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":deletable-chips", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Dense(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":dense", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) DisableLookup(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":disable-lookup", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Disabled(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Eager(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Error(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) ErrorCount(v int) (r *VXAutocompleteBuilder) {
	b.Attr(":error-count", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Filled(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":filled", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Filter(v interface{}) (r *VXAutocompleteBuilder) {
	b.Attr(":filter", h.JSONString(v))
	return b
}

func (b *VXAutocompleteBuilder) Flat(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) FullWidth(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":full-width", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Height(v int) (r *VXAutocompleteBuilder) {
	b.Attr(":height", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) HideDetails(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":hide-details", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) HideNoData(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":hide-no-data", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) HideSelected(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":hide-selected", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Hint(v string) (r *VXAutocompleteBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VXAutocompleteBuilder) Id(v string) (r *VXAutocompleteBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VXAutocompleteBuilder) ItemColor(v string) (r *VXAutocompleteBuilder) {
	b.Attr("item-color", v)
	return b
}

func (b *VXAutocompleteBuilder) ItemDisabled(v string) (r *VXAutocompleteBuilder) {
	b.Attr("item-disabled", v)
	return b
}

func (b *VXAutocompleteBuilder) ItemText(v string) (r *VXAutocompleteBuilder) {
	b.Attr("item-text", v)
	return b
}

func (b *VXAutocompleteBuilder) ItemValue(v string) (r *VXAutocompleteBuilder) {
	b.Attr("item-value", v)
	return b
}

func (b *VXAutocompleteBuilder) Label(v string) (r *VXAutocompleteBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VXAutocompleteBuilder) Light(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":light", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) LoaderHeight(v int) (r *VXAutocompleteBuilder) {
	b.Attr(":loader-height", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Loading(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":loading", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) MenuProps(v interface{}) (r *VXAutocompleteBuilder) {
	b.Attr(":menu-props", h.JSONString(v))
	return b
}

func (b *VXAutocompleteBuilder) Messages(v string) (r *VXAutocompleteBuilder) {
	b.Attr("messages", v)
	return b
}

func (b *VXAutocompleteBuilder) Multiple(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) NoDataText(v string) (r *VXAutocompleteBuilder) {
	b.Attr("no-data-text", v)
	return b
}

func (b *VXAutocompleteBuilder) NoFilter(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":no-filter", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) OpenOnClear(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":open-on-clear", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Outlined(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":outlined", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) PersistentHint(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) PersistentPlaceholder(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":persistent-placeholder", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Placeholder(v string) (r *VXAutocompleteBuilder) {
	b.Attr("placeholder", v)
	return b
}

func (b *VXAutocompleteBuilder) Prefix(v string) (r *VXAutocompleteBuilder) {
	b.Attr("prefix", v)
	return b
}

func (b *VXAutocompleteBuilder) PrependIcon(v string) (r *VXAutocompleteBuilder) {
	b.Attr("prepend-icon", v)
	return b
}

func (b *VXAutocompleteBuilder) PrependInnerIcon(v string) (r *VXAutocompleteBuilder) {
	b.Attr("prepend-inner-icon", v)
	return b
}

func (b *VXAutocompleteBuilder) Readonly(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) ReturnObject(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Reverse(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Rounded(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":rounded", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Rules(v interface{}) (r *VXAutocompleteBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VXAutocompleteBuilder) SearchInput(v string) (r *VXAutocompleteBuilder) {
	b.Attr("search-input", v)
	return b
}

func (b *VXAutocompleteBuilder) Shaped(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":shaped", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) SingleLine(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) SmallChips(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":small-chips", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Solo(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":solo", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) SoloInverted(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":solo-inverted", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Success(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":success", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) SuccessMessages(v string) (r *VXAutocompleteBuilder) {
	b.Attr("success-messages", v)
	return b
}

func (b *VXAutocompleteBuilder) Suffix(v string) (r *VXAutocompleteBuilder) {
	b.Attr("suffix", v)
	return b
}

func (b *VXAutocompleteBuilder) Type(v string) (r *VXAutocompleteBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VXAutocompleteBuilder) ValidateOnBlur(v bool) (r *VXAutocompleteBuilder) {
	b.Attr(":validate-on-blur", fmt.Sprint(v))
	return b
}

func (b *VXAutocompleteBuilder) Value(v interface{}) (r *VXAutocompleteBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VXAutocompleteBuilder) ValueComparator(v interface{}) (r *VXAutocompleteBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VXAutocompleteBuilder) On(name string, value string) (r *VXAutocompleteBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VXAutocompleteBuilder) Bind(name string, value string) (r *VXAutocompleteBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
