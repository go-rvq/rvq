package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSliderBuilder struct {
	VTagBuilder[*VSliderBuilder]
}

func VSlider(children ...h.HTMLComponent) *VSliderBuilder {
	return VTag(&VSliderBuilder{}, "v-slider", children...)
}

func (b *VSliderBuilder) Label(v string) (r *VSliderBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VSliderBuilder) Focused(v bool) (r *VSliderBuilder) {
	b.Attr(":focused", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) Reverse(v bool) (r *VSliderBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) Disabled(v bool) (r *VSliderBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) Error(v bool) (r *VSliderBuilder) {
	b.Attr(":error", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) Readonly(v bool) (r *VSliderBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
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

func (b *VSliderBuilder) Color(v string) (r *VSliderBuilder) {
	b.Attr("color", v)
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

func (b *VSliderBuilder) Direction(v interface{}) (r *VSliderBuilder) {
	b.Attr(":direction", h.JSONString(v))
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

func (b *VSliderBuilder) Elevation(v interface{}) (r *VSliderBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Ripple(v bool) (r *VSliderBuilder) {
	b.Attr(":ripple", fmt.Sprint(v))
	return b
}

func (b *VSliderBuilder) Id(v string) (r *VSliderBuilder) {
	b.Attr("id", v)
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

func (b *VSliderBuilder) Theme(v string) (r *VSliderBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSliderBuilder) MaxErrors(v interface{}) (r *VSliderBuilder) {
	b.Attr(":max-errors", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) Name(v string) (r *VSliderBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VSliderBuilder) Rules(v interface{}) (r *VSliderBuilder) {
	b.Attr(":rules", h.JSONString(v))
	return b
}

func (b *VSliderBuilder) ModelValue(v interface{}) (r *VSliderBuilder) {
	b.Attr(":model-value", h.JSONString(v))
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

func (b *VSliderBuilder) HideDetails(v interface{}) (r *VSliderBuilder) {
	b.Attr(":hide-details", h.JSONString(v))
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
