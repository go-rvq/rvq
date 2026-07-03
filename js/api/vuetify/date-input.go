package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDateInputBuilder struct {
	VTagBuilder[*VDateInputBuilder]
}

func VDateInput(children ...h.HTMLComponent) *VDateInputBuilder {
	return VTag(&VDateInputBuilder{}, "v-date-input", children...)
}

func (b *VDateInputBuilder) Title(v string) (r *VDateInputBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VDateInputBuilder) Text(v string) (r *VDateInputBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VDateInputBuilder) Flat(v bool) (r *VDateInputBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Border(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Type(v string) (r *VDateInputBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VDateInputBuilder) ModelValue(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Error(v bool) (r *VDateInputBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Reverse(v bool) (r *VDateInputBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Density(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Height(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) MaxHeight(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) MaxWidth(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) MinHeight(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) MinWidth(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Width(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Elevation(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Location(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Position(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Rounded(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Tile(v bool) (r *VDateInputBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Tag(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Theme(v string) (r *VDateInputBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VDateInputBuilder) Color(v string) (r *VDateInputBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDateInputBuilder) Variant(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Name(v string) (r *VDateInputBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VDateInputBuilder) Header(v string) (r *VDateInputBuilder) {
	b.Attr("header", v)
	return b
}

func (b *VDateInputBuilder) Active(v bool) (r *VDateInputBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) BaseColor(v string) (r *VDateInputBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VDateInputBuilder) PrependIcon(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) AppendIcon(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Readonly(v bool) (r *VDateInputBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Disabled(v bool) (r *VDateInputBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Loading(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Label(v string) (r *VDateInputBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VDateInputBuilder) Max(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Transition(v string) (r *VDateInputBuilder) {
	b.Attr("transition", v)
	return b
}

func (b *VDateInputBuilder) BgColor(v string) (r *VDateInputBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VDateInputBuilder) Mobile(v bool) (r *VDateInputBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) MobileBreakpoint(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Multiple(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":multiple", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) ID(v string) (r *VDateInputBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VDateInputBuilder) Prefix(v string) (r *VDateInputBuilder) {
	b.Attr("prefix", v)
	return b
}

func (b *VDateInputBuilder) Role(v string) (r *VDateInputBuilder) {
	b.Attr("role", v)
	return b
}

func (b *VDateInputBuilder) Divided(v bool) (r *VDateInputBuilder) {
	b.Attr(":divided", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) HideHeader(v bool) (r *VDateInputBuilder) {
	b.Attr(":hide-header", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Month(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":month", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) ShowAdjacentMonths(v bool) (r *VDateInputBuilder) {
	b.Attr(":show-adjacent-months", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Year(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":year", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Weekdays(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":weekdays", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) WeeksInMonth(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":weeks-in-month", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) FirstDayOfWeek(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":first-day-of-week", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) AllowedDates(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":allowed-dates", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Min(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) NextIcon(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) PrevIcon(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) ViewMode(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":view-mode", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Direction(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Placeholder(v string) (r *VDateInputBuilder) {
	b.Attr("placeholder", v)
	return b
}

func (b *VDateInputBuilder) ReverseTransition(v string) (r *VDateInputBuilder) {
	b.Attr("reverse-transition", v)
	return b
}

func (b *VDateInputBuilder) CenterAffix(v bool) (r *VDateInputBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Glow(v bool) (r *VDateInputBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) IconColor(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) HideSpinButtons(v bool) (r *VDateInputBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Hint(v string) (r *VDateInputBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VDateInputBuilder) PersistentHint(v bool) (r *VDateInputBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Messages(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) ErrorMessages(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) MaxErrors(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Rules(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) ValidateOn(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) ValidationValue(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Focused(v bool) (r *VDateInputBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) HideDetails(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Landscape(v bool) (r *VDateInputBuilder) {
	b.Attr(":landscape", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Autofocus(v bool) (r *VDateInputBuilder) {
	b.Attr(":autofocus", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Counter(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":counter", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) PersistentPlaceholder(v bool) (r *VDateInputBuilder) {
	b.Attr(":persistent-placeholder", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) PersistentCounter(v bool) (r *VDateInputBuilder) {
	b.Attr(":persistent-counter", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) Suffix(v string) (r *VDateInputBuilder) {
	b.Attr("suffix", v)
	return b
}

func (b *VDateInputBuilder) AppendInnerIcon(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":append-inner-icon", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Clearable(v bool) (r *VDateInputBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) ClearIcon(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":clear-icon", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) Dirty(v bool) (r *VDateInputBuilder) {
	b.Attr(":dirty", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) PersistentClear(v bool) (r *VDateInputBuilder) {
	b.Attr(":persistent-clear", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) PrependInnerIcon(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":prepend-inner-icon", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) SingleLine(v bool) (r *VDateInputBuilder) {
	b.Attr(":single-line", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) CounterValue(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":counter-value", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) ModelModifiers(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":model-modifiers", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) DisplayFormat(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":display-format", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) CancelText(v string) (r *VDateInputBuilder) {
	b.Attr("cancel-text", v)
	return b
}

func (b *VDateInputBuilder) OkText(v string) (r *VDateInputBuilder) {
	b.Attr("ok-text", v)
	return b
}

func (b *VDateInputBuilder) HideActions(v bool) (r *VDateInputBuilder) {
	b.Attr(":hide-actions", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) HeaderColor(v string) (r *VDateInputBuilder) {
	b.Attr("header-color", v)
	return b
}

func (b *VDateInputBuilder) ControlHeight(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":control-height", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) ModeIcon(v interface{}) (r *VDateInputBuilder) {
	b.Attr(":mode-icon", h.JSONString(v))
	return b
}

func (b *VDateInputBuilder) HideWeekdays(v bool) (r *VDateInputBuilder) {
	b.Attr(":hide-weekdays", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) ShowWeek(v bool) (r *VDateInputBuilder) {
	b.Attr(":show-week", fmt.Sprint(v))
	return b
}

func (b *VDateInputBuilder) On(name string, value string) (r *VDateInputBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDateInputBuilder) Bind(name string, value string) (r *VDateInputBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDateInputBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDateInputBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDateInputBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDateInputBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDateInputBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VDateInputBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VDateInputBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VDateInputBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VDateInputBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VDateInputBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VDateInputBuilder) SlotAppend(child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VDateInputBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VDateInputBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VDateInputBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VDateInputBuilder) SlotLoader(child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VDateInputBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VDateInputBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VDateInputBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VDateInputBuilder) SlotLabel(child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VDateInputBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}

func (b *VDateInputBuilder) SetSlotClear(child ...h.HTMLComponent) {
	b.SetSlot("clear", child...)
}

func (b *VDateInputBuilder) SetScopedSlotClear(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("clear", scope, child...)
}

func (b *VDateInputBuilder) SlotClear(child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetSlotClear(child...)
	return b
}

func (b *VDateInputBuilder) ScopedSlotClear(scope string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetScopedSlotClear(scope, child...)
	return b
}

func (b *VDateInputBuilder) SetSlotDetails(child ...h.HTMLComponent) {
	b.SetSlot("details", child...)
}

func (b *VDateInputBuilder) SetScopedSlotDetails(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("details", scope, child...)
}

func (b *VDateInputBuilder) SlotDetails(child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetSlotDetails(child...)
	return b
}

func (b *VDateInputBuilder) ScopedSlotDetails(scope string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetScopedSlotDetails(scope, child...)
	return b
}

func (b *VDateInputBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VDateInputBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VDateInputBuilder) SlotMessage(child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VDateInputBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}

func (b *VDateInputBuilder) SetSlotPrependInner(child ...h.HTMLComponent) {
	b.SetSlot("prepend-inner", child...)
}

func (b *VDateInputBuilder) SetScopedSlotPrependInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend-inner", scope, child...)
}

func (b *VDateInputBuilder) SlotPrependInner(child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetSlotPrependInner(child...)
	return b
}

func (b *VDateInputBuilder) ScopedSlotPrependInner(scope string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetScopedSlotPrependInner(scope, child...)
	return b
}

func (b *VDateInputBuilder) SetSlotAppendInner(child ...h.HTMLComponent) {
	b.SetSlot("append-inner", child...)
}

func (b *VDateInputBuilder) SetScopedSlotAppendInner(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append-inner", scope, child...)
}

func (b *VDateInputBuilder) SlotAppendInner(child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetSlotAppendInner(child...)
	return b
}

func (b *VDateInputBuilder) ScopedSlotAppendInner(scope string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetScopedSlotAppendInner(scope, child...)
	return b
}

func (b *VDateInputBuilder) SetSlotCounter(child ...h.HTMLComponent) {
	b.SetSlot("counter", child...)
}

func (b *VDateInputBuilder) SetScopedSlotCounter(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("counter", scope, child...)
}

func (b *VDateInputBuilder) SlotCounter(child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetSlotCounter(child...)
	return b
}

func (b *VDateInputBuilder) ScopedSlotCounter(scope string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetScopedSlotCounter(scope, child...)
	return b
}

func (b *VDateInputBuilder) SetSlotActions(child ...h.HTMLComponent) {
	b.SetSlot("actions", child...)
}

func (b *VDateInputBuilder) SetScopedSlotActions(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("actions", scope, child...)
}

func (b *VDateInputBuilder) SlotActions(child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetSlotActions(child...)
	return b
}

func (b *VDateInputBuilder) ScopedSlotActions(scope string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetScopedSlotActions(scope, child...)
	return b
}

func (b *VDateInputBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VDateInputBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VDateInputBuilder) SlotDefault(child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VDateInputBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VDateInputBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
