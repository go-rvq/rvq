package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSelectBuilder struct {
	VTagBuilder[*VSelectBuilder]
}

func (b *VSelectBuilder) Label(v string) (r *VSelectBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VSelectBuilder) Chips(v bool) (r *VSelectBuilder) {
	b.Attr(":chips", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) ClosableChips(v bool) (r *VSelectBuilder) {
	b.Attr(":closable-chips", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) CloseText(v string) (r *VSelectBuilder) {
	b.Attr("close-text", v)
	return b
}

func (b *VSelectBuilder) Type(v string) (r *VSelectBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VSelectBuilder) OpenText(v string) (r *VSelectBuilder) {
	b.Attr("open-text", v)
	return b
}

func (b *VSelectBuilder) Eager(v bool) (r *VSelectBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) HideNoData(v bool) (r *VSelectBuilder) {
	b.Attr(":hide-no-data", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) HideSelected(v bool) (r *VSelectBuilder) {
	b.Attr(":hide-selected", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) ListProps(v interface{}) (r *VSelectBuilder) {
	b.Attr(":list-props", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) BaseColor(v string) (r *VSelectBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VSelectBuilder) BgColor(v string) (r *VSelectBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VSelectBuilder) Disabled(v bool) (r *VSelectBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) Multiple(v bool) (r *VSelectBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) Reverse(v bool) (r *VSelectBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) Flat(v bool) (r *VSelectBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) Density(v interface{}) (r *VSelectBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) MaxWidth(v interface{}) (r *VSelectBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) MinWidth(v interface{}) (r *VSelectBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Width(v interface{}) (r *VSelectBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Items(v interface{}) (r *VSelectBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) ItemTitle(v interface{}) (r *VSelectBuilder) {
	b.Attr(":item-title", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) ItemValue(v interface{}) (r *VSelectBuilder) {
	b.Attr(":item-value", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) ItemChildren(v interface{}) (r *VSelectBuilder) {
	b.Attr(":item-children", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) ItemProps(v interface{}) (r *VSelectBuilder) {
	b.Attr(":item-props", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) ReturnObject(v bool) (r *VSelectBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) ValueComparator(v interface{}) (r *VSelectBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Rounded(v interface{}) (r *VSelectBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Tile(v bool) (r *VSelectBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) Theme(v string) (r *VSelectBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSelectBuilder) Color(v string) (r *VSelectBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSelectBuilder) Variant(v interface{}) (r *VSelectBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Name(v string) (r *VSelectBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VSelectBuilder) Menu(v bool) (r *VSelectBuilder) {
	b.Attr(":menu", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) MenuIcon(v interface{}) (r *VSelectBuilder) {
	b.Attr(":menu-icon", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) MenuProps(v interface{}) (r *VSelectBuilder) {
	b.Attr(":menu-props", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Id(v string) (r *VSelectBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VSelectBuilder) ModelValue(v interface{}) (r *VSelectBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Transition(v interface{}) (r *VSelectBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) NoDataText(v string) (r *VSelectBuilder) {
	b.Attr("no-data-text", v)
	return b
}

func (b *VSelectBuilder) OpenOnClear(v bool) (r *VSelectBuilder) {
	b.Attr(":open-on-clear", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) ItemColor(v string) (r *VSelectBuilder) {
	b.Attr("item-color", v)
	return b
}

func (b *VSelectBuilder) Autofocus(v bool) (r *VSelectBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) Counter(v interface{}) (r *VSelectBuilder) {
	b.Attr(":counter", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Prefix(v string) (r *VSelectBuilder) {
	b.Attr("prefix", v)
	return b
}

func (b *VSelectBuilder) Placeholder(v string) (r *VSelectBuilder) {
	b.Attr("placeholder", v)
	return b
}

func (b *VSelectBuilder) PersistentPlaceholder(v bool) (r *VSelectBuilder) {
	b.Attr(":persistent-placeholder", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) PersistentCounter(v bool) (r *VSelectBuilder) {
	b.Attr(":persistent-counter", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) Suffix(v string) (r *VSelectBuilder) {
	b.Attr("suffix", v)
	return b
}

func (b *VSelectBuilder) Role(v string) (r *VSelectBuilder) {
	b.Attr("role", v)
	return b
}

func (b *VSelectBuilder) AppendIcon(v interface{}) (r *VSelectBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) CenterAffix(v bool) (r *VSelectBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) PrependIcon(v interface{}) (r *VSelectBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) HideSpinButtons(v bool) (r *VSelectBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) Hint(v string) (r *VSelectBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VSelectBuilder) PersistentHint(v bool) (r *VSelectBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) Messages(v interface{}) (r *VSelectBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Direction(v interface{}) (r *VSelectBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Error(v bool) (r *VSelectBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) MaxErrors(v interface{}) (r *VSelectBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Readonly(v bool) (r *VSelectBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) Rules(v interface{}) (r *VSelectBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) ValidateOn(v interface{}) (r *VSelectBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Focused(v bool) (r *VSelectBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) HideDetails(v interface{}) (r *VSelectBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Clearable(v bool) (r *VSelectBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) ClearIcon(v interface{}) (r *VSelectBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) Active(v bool) (r *VSelectBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) PersistentClear(v bool) (r *VSelectBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) PrependInnerIcon(v interface{}) (r *VSelectBuilder) {
	b.Attr(":prepend-inner-icon", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) SingleLine(v bool) (r *VSelectBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VSelectBuilder) Loading(v interface{}) (r *VSelectBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) CounterValue(v interface{}) (r *VSelectBuilder) {
	b.Attr(":counter-value", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) ModelModifiers(v interface{}) (r *VSelectBuilder) {
	b.Attr(":model-modifiers", h.JSONString(v))
	return b
}

func (b *VSelectBuilder) On(name string, value string) (r *VSelectBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSelectBuilder) Bind(name string, value string) (r *VSelectBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
