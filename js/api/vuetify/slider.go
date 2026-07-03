package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSliderBuilder struct {
	VTagBuilder[*VSliderBuilder]
}

func VSlider(children ...h.HTMLComponent) *VSliderBuilder {
	return VTag(&VSliderBuilder{}, "v-slider", children...)
}

func (b *VSliderBuilder) Reverse(v bool) (r *VSliderBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) Name(v string) (r *VSliderBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VSliderBuilder) Error(v bool) (r *VSliderBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) Label(v string) (r *VSliderBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VSliderBuilder) Theme(v string) (r *VSliderBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSliderBuilder) ID(v string) (r *VSliderBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VSliderBuilder) BaseColor(v string) (r *VSliderBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VSliderBuilder) Disabled(v bool) (r *VSliderBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) Density(v interface{}) (r *VSliderBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) MaxWidth(v interface{}) (r *VSliderBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) MinWidth(v interface{}) (r *VSliderBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Width(v interface{}) (r *VSliderBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Elevation(v interface{}) (r *VSliderBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Rounded(v interface{}) (r *VSliderBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Tile(v bool) (r *VSliderBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) Color(v string) (r *VSliderBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSliderBuilder) ModelValue(v interface{}) (r *VSliderBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) AppendIcon(v interface{}) (r *VSliderBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) CenterAffix(v bool) (r *VSliderBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) Glow(v bool) (r *VSliderBuilder) {
	b.Attr(":glow", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) IconColor(v interface{}) (r *VSliderBuilder) {
	b.Attr(":icon-color", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) PrependIcon(v interface{}) (r *VSliderBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) HideSpinButtons(v bool) (r *VSliderBuilder) {
	b.Attr(":hide-spin-buttons", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) Hint(v string) (r *VSliderBuilder) {
	b.Attr("hint", v)
	return b
}

func (b *VSliderBuilder) PersistentHint(v bool) (r *VSliderBuilder) {
	b.Attr(":persistent-hint", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) Messages(v interface{}) (r *VSliderBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Direction(v interface{}) (r *VSliderBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) ErrorMessages(v interface{}) (r *VSliderBuilder) {
	b.Attr(":error-messages", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) MaxErrors(v interface{}) (r *VSliderBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Readonly(v bool) (r *VSliderBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) Rules(v interface{}) (r *VSliderBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) ValidateOn(v interface{}) (r *VSliderBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) ValidationValue(v interface{}) (r *VSliderBuilder) {
	b.Attr(":validation-value", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Focused(v bool) (r *VSliderBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) HideDetails(v interface{}) (r *VSliderBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Max(v interface{}) (r *VSliderBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Min(v interface{}) (r *VSliderBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Step(v interface{}) (r *VSliderBuilder) {
	b.Attr(":step", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Ripple(v bool) (r *VSliderBuilder) {
	b.Attr(":ripple", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) ThumbColor(v string) (r *VSliderBuilder) {
	b.Attr("thumb-color", v)
	return b
}

func (b *VSliderBuilder) ThumbLabel(v interface{}) (r *VSliderBuilder) {
	b.Attr(":thumb-label", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) ThumbSize(v interface{}) (r *VSliderBuilder) {
	b.Attr(":thumb-size", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) ShowTicks(v interface{}) (r *VSliderBuilder) {
	b.Attr(":show-ticks", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Ticks(v interface{}) (r *VSliderBuilder) {
	b.Attr(":ticks", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) TickSize(v interface{}) (r *VSliderBuilder) {
	b.Attr(":tick-size", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) TrackColor(v string) (r *VSliderBuilder) {
	b.Attr("track-color", v)
	return b
}

func (b *VSliderBuilder) TrackFillColor(v string) (r *VSliderBuilder) {
	b.Attr("track-fill-color", v)
	return b
}

func (b *VSliderBuilder) TrackSize(v interface{}) (r *VSliderBuilder) {
	b.Attr(":track-size", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) On(name string, value string) (r *VSliderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSliderBuilder) Bind(name string, value string) (r *VSliderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VSliderBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSliderBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSliderBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSliderBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VSliderBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VSliderBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VSliderBuilder) SlotDefault(child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VSliderBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VSliderBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VSliderBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VSliderBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VSliderBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VSliderBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VSliderBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VSliderBuilder) SlotAppend(child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VSliderBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VSliderBuilder) SetSlotDetails(child ...h.HTMLComponent) {
	b.SetSlot("details", child...)
}

func (b *VSliderBuilder) SetScopedSlotDetails(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("details", scope, child...)
}

func (b *VSliderBuilder) SlotDetails(child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetSlotDetails(child...)
	return b
}

func (b *VSliderBuilder) ScopedSlotDetails(scope string, child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetScopedSlotDetails(scope, child...)
	return b
}

func (b *VSliderBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VSliderBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VSliderBuilder) SlotMessage(child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VSliderBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}

func (b *VSliderBuilder) SetSlotThumbLabel(child ...h.HTMLComponent) {
	b.SetSlot("thumb-label", child...)
}

func (b *VSliderBuilder) SetScopedSlotThumbLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("thumb-label", scope, child...)
}

func (b *VSliderBuilder) SlotThumbLabel(child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetSlotThumbLabel(child...)
	return b
}

func (b *VSliderBuilder) ScopedSlotThumbLabel(scope string, child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetScopedSlotThumbLabel(scope, child...)
	return b
}

func (b *VSliderBuilder) SetSlotTickLabel(child ...h.HTMLComponent) {
	b.SetSlot("tick-label", child...)
}

func (b *VSliderBuilder) SetScopedSlotTickLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("tick-label", scope, child...)
}

func (b *VSliderBuilder) SlotTickLabel(child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetSlotTickLabel(child...)
	return b
}

func (b *VSliderBuilder) ScopedSlotTickLabel(scope string, child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetScopedSlotTickLabel(scope, child...)
	return b
}

func (b *VSliderBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VSliderBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VSliderBuilder) SlotLabel(child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VSliderBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VSliderBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}
