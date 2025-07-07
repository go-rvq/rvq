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

func (b *VRangeSliderBuilder) Label(v string) (r *VRangeSliderBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VRangeSliderBuilder) Focused(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Reverse(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Id(v string) (r *VRangeSliderBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VRangeSliderBuilder) AppendIcon(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) CenterAffix(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":center-affix", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) PrependIcon(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
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

func (b *VRangeSliderBuilder) Direction(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":direction", h.JSONString(v))
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

func (b *VRangeSliderBuilder) Theme(v string) (r *VRangeSliderBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VRangeSliderBuilder) Disabled(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Error(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":error", fmt.Sprint(v))
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

func (b *VRangeSliderBuilder) Name(v string) (r *VRangeSliderBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VRangeSliderBuilder) Readonly(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Rules(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) ModelValue(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":model-value", h.JSONString(v))
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

func (b *VRangeSliderBuilder) HideDetails(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Max(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Min(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":min", h.JSONString(v))
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

func (b *VRangeSliderBuilder) Color(v string) (r *VRangeSliderBuilder) {
	b.Attr("color", v)
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

func (b *VRangeSliderBuilder) Rounded(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Tile(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VRangeSliderBuilder) Elevation(v interface{}) (r *VRangeSliderBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VRangeSliderBuilder) Ripple(v bool) (r *VRangeSliderBuilder) {
	b.Attr(":ripple", fmt.Sprint(v))
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
