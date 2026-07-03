package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSwitchBuilder struct {
	VTagBuilder[*VSwitchBuilder]
}

func VSwitch(children ...h.HTMLComponent) *VSwitchBuilder {
	return VTag(&VSwitchBuilder{}, "v-switch", children...)
}

func (b *VSwitchBuilder) Flat(v bool) (r *VSwitchBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Type(v string) (r *VSwitchBuilder) {
	b.Attr("type", v)
	return b
}

func (b *VSwitchBuilder) Name(v string) (r *VSwitchBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VSwitchBuilder) Error(v bool) (r *VSwitchBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Label(v string) (r *VSwitchBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VSwitchBuilder) Theme(v string) (r *VSwitchBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSwitchBuilder) ID(v string) (r *VSwitchBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VSwitchBuilder) Value(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) BaseColor(v string) (r *VSwitchBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VSwitchBuilder) Disabled(v bool) (r *VSwitchBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Multiple(v bool) (r *VSwitchBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Density(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) MaxWidth(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) MinWidth(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Width(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) ValueComparator(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Color(v string) (r *VSwitchBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSwitchBuilder) ModelValue(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) AppendIcon(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) CenterAffix(v bool) (r *VSwitchBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Glow(v bool) (r *VSwitchBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) IconColor(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) PrependIcon(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) HideSpinButtons(v bool) (r *VSwitchBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Hint(v string) (r *VSwitchBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VSwitchBuilder) PersistentHint(v bool) (r *VSwitchBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Messages(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Direction(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) ErrorMessages(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) MaxErrors(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Readonly(v bool) (r *VSwitchBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Rules(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) ValidateOn(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) ValidationValue(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Focused(v bool) (r *VSwitchBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) HideDetails(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Loading(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Indeterminate(v bool) (r *VSwitchBuilder) {
	b.Attr(":indeterminate", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Inset(v bool) (r *VSwitchBuilder) {
	b.Attr(":inset", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) Ripple(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) Inline(v bool) (r *VSwitchBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VSwitchBuilder) TrueValue(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":true-value", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) FalseValue(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":false-value", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) DefaultsTarget(v string) (r *VSwitchBuilder) {
	b.Attr("defaults-target", v)
	return b
}

func (b *VSwitchBuilder) FalseIcon(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":false-icon", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) TrueIcon(v interface{}) (r *VSwitchBuilder) {
	b.Attr(":true-icon", h.JSONString(v))
	return b
}

func (b *VSwitchBuilder) On(name string, value string) (r *VSwitchBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSwitchBuilder) Bind(name string, value string) (r *VSwitchBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VSwitchBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSwitchBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSwitchBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSwitchBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VSwitchBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VSwitchBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VSwitchBuilder) SlotDefault(child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VSwitchBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VSwitchBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VSwitchBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VSwitchBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VSwitchBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VSwitchBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VSwitchBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VSwitchBuilder) SlotAppend(child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VSwitchBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VSwitchBuilder) SetSlotDetails(child ...h.HTMLComponent) {
	b.SetSlot("details", child...)
}

func (b *VSwitchBuilder) SetScopedSlotDetails(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("details", scope, child...)
}

func (b *VSwitchBuilder) SlotDetails(child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetSlotDetails(child...)
	return b
}

func (b *VSwitchBuilder) ScopedSlotDetails(scope string, child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetScopedSlotDetails(scope, child...)
	return b
}

func (b *VSwitchBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VSwitchBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VSwitchBuilder) SlotMessage(child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VSwitchBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}

func (b *VSwitchBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VSwitchBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VSwitchBuilder) SlotLabel(child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VSwitchBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}

func (b *VSwitchBuilder) SetSlotInput(child ...h.HTMLComponent) {
	b.SetSlot("input", child...)
}

func (b *VSwitchBuilder) SetScopedSlotInput(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("input", scope, child...)
}

func (b *VSwitchBuilder) SlotInput(child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetSlotInput(child...)
	return b
}

func (b *VSwitchBuilder) ScopedSlotInput(scope string, child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetScopedSlotInput(scope, child...)
	return b
}

func (b *VSwitchBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VSwitchBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VSwitchBuilder) SlotLoader(child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VSwitchBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VSwitchBuilder) SetSlotThumb(child ...h.HTMLComponent) {
	b.SetSlot("thumb", child...)
}

func (b *VSwitchBuilder) SetScopedSlotThumb(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("thumb", scope, child...)
}

func (b *VSwitchBuilder) SlotThumb(child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetSlotThumb(child...)
	return b
}

func (b *VSwitchBuilder) ScopedSlotThumb(scope string, child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetScopedSlotThumb(scope, child...)
	return b
}

func (b *VSwitchBuilder) SetSlotTrackFalse(child ...h.HTMLComponent) {
	b.SetSlot("track-false", child...)
}

func (b *VSwitchBuilder) SetScopedSlotTrackFalse(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("track-false", scope, child...)
}

func (b *VSwitchBuilder) SlotTrackFalse(child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetSlotTrackFalse(child...)
	return b
}

func (b *VSwitchBuilder) ScopedSlotTrackFalse(scope string, child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetScopedSlotTrackFalse(scope, child...)
	return b
}

func (b *VSwitchBuilder) SetSlotTrackTrue(child ...h.HTMLComponent) {
	b.SetSlot("track-true", child...)
}

func (b *VSwitchBuilder) SetScopedSlotTrackTrue(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("track-true", scope, child...)
}

func (b *VSwitchBuilder) SlotTrackTrue(child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetSlotTrackTrue(child...)
	return b
}

func (b *VSwitchBuilder) ScopedSlotTrackTrue(scope string, child ...h.HTMLComponent) (r *VSwitchBuilder) {
	b.SetScopedSlotTrackTrue(scope, child...)
	return b
}
