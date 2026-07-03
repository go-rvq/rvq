package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VRangeSliderBuilder struct {
	VTagBuilder[*VRangeSliderBuilder]
}

func VRangeSlider(children ...h.HTMLComponent) *VRangeSliderBuilder {
	return VTag(&VRangeSliderBuilder{}, "v-range-slider", children...)
}

func (b *VRangeSliderBuilder) ModelValue(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Error(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Reverse(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Density(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) MaxWidth(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) MinWidth(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Width(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Elevation(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Rounded(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Tile(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Theme(v string) (r *VRangeSliderBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VRangeSliderBuilder) Color(v string) (r *VRangeSliderBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VRangeSliderBuilder) Name(v string) (r *VRangeSliderBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VRangeSliderBuilder) BaseColor(v string) (r *VRangeSliderBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VRangeSliderBuilder) PrependIcon(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) AppendIcon(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Readonly(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Ripple(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":ripple", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Disabled(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Label(v string) (r *VRangeSliderBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VRangeSliderBuilder) Max(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) ID(v string) (r *VRangeSliderBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VRangeSliderBuilder) Min(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Direction(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) CenterAffix(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Glow(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) IconColor(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) HideSpinButtons(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Hint(v string) (r *VRangeSliderBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VRangeSliderBuilder) PersistentHint(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Messages(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) ErrorMessages(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) MaxErrors(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Rules(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) ValidateOn(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) ValidationValue(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Focused(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) HideDetails(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Step(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":step", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) ThumbColor(v string) (r *VRangeSliderBuilder) {
	b.Attr("thumb-color", v)
	return b
}

func (b *VRangeSliderBuilder) ThumbLabel(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":thumb-label", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) ThumbSize(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":thumb-size", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) ShowTicks(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":show-ticks", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Ticks(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":ticks", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) TickSize(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":tick-size", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) TrackColor(v string) (r *VRangeSliderBuilder) {
	b.Attr("track-color", v)
	return b
}

func (b *VRangeSliderBuilder) TrackFillColor(v string) (r *VRangeSliderBuilder) {
	b.Attr("track-fill-color", v)
	return b
}

func (b *VRangeSliderBuilder) TrackSize(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":track-size", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Strict(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":strict", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) On(name string, value string) (r *VRangeSliderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VRangeSliderBuilder) Bind(name string, value string) (r *VRangeSliderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VRangeSliderBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VRangeSliderBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VRangeSliderBuilder) Slot(name string, child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VRangeSliderBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VRangeSliderBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VRangeSliderBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VRangeSliderBuilder) SlotDefault(child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VRangeSliderBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VRangeSliderBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VRangeSliderBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VRangeSliderBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VRangeSliderBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VRangeSliderBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VRangeSliderBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VRangeSliderBuilder) SlotAppend(child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VRangeSliderBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VRangeSliderBuilder) SetSlotDetails(child ...h.HTMLComponent) {
	b.SetSlot("details", child...)
}

func (b *VRangeSliderBuilder) SetScopedSlotDetails(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("details", scope, child...)
}

func (b *VRangeSliderBuilder) SlotDetails(child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetSlotDetails(child...)
	return b
}

func (b *VRangeSliderBuilder) ScopedSlotDetails(scope string, child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetScopedSlotDetails(scope, child...)
	return b
}

func (b *VRangeSliderBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VRangeSliderBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VRangeSliderBuilder) SlotMessage(child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VRangeSliderBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}

func (b *VRangeSliderBuilder) SetSlotThumbLabel(child ...h.HTMLComponent) {
	b.SetSlot("thumb-label", child...)
}

func (b *VRangeSliderBuilder) SetScopedSlotThumbLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("thumb-label", scope, child...)
}

func (b *VRangeSliderBuilder) SlotThumbLabel(child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetSlotThumbLabel(child...)
	return b
}

func (b *VRangeSliderBuilder) ScopedSlotThumbLabel(scope string, child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetScopedSlotThumbLabel(scope, child...)
	return b
}

func (b *VRangeSliderBuilder) SetSlotTickLabel(child ...h.HTMLComponent) {
	b.SetSlot("tick-label", child...)
}

func (b *VRangeSliderBuilder) SetScopedSlotTickLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("tick-label", scope, child...)
}

func (b *VRangeSliderBuilder) SlotTickLabel(child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetSlotTickLabel(child...)
	return b
}

func (b *VRangeSliderBuilder) ScopedSlotTickLabel(scope string, child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetScopedSlotTickLabel(scope, child...)
	return b
}

func (b *VRangeSliderBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VRangeSliderBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VRangeSliderBuilder) SlotLabel(child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VRangeSliderBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VRangeSliderBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}
