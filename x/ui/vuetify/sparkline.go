package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSparklineBuilder struct {
	VTagBuilder[*VSparklineBuilder]
}

func VSparkline(children ...h.HTMLComponent) *VSparklineBuilder {
	return VTag(&VSparklineBuilder{}, "v-sparkline", children...)
}

func (b *VSparklineBuilder) Type(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":type", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) AutoLineWidth(v bool) (r *VSparklineBuilder) {
	b.Attr(":auto-line-width", fmt.Sprint(v))
	return b
}

func (b *VSparklineBuilder) AutoDraw(v bool) (r *VSparklineBuilder) {
	b.Attr(":auto-draw", fmt.Sprint(v))
	return b
}

func (b *VSparklineBuilder) AutoDrawDuration(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":auto-draw-duration", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) AutoDrawEasing(v string) (r *VSparklineBuilder) {
	b.Attr("auto-draw-easing", v)
	return b
}

func (b *VSparklineBuilder) Color(v string) (r *VSparklineBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSparklineBuilder) Gradient(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":gradient", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) GradientDirection(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":gradient-direction", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) Height(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) Labels(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":labels", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) LabelSize(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":label-size", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) LineWidth(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":line-width", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) Id(v string) (r *VSparklineBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VSparklineBuilder) ItemValue(v string) (r *VSparklineBuilder) {
	b.Attr("item-value", v)
	return b
}

func (b *VSparklineBuilder) ModelValue(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) Min(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) Max(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) Padding(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":padding", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) ShowLabels(v bool) (r *VSparklineBuilder) {
	b.Attr(":show-labels", fmt.Sprint(v))
	return b
}

func (b *VSparklineBuilder) Smooth(v bool) (r *VSparklineBuilder) {
	b.Attr(":smooth", fmt.Sprint(v))
	return b
}

func (b *VSparklineBuilder) Width(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) Fill(v bool) (r *VSparklineBuilder) {
	b.Attr(":fill", fmt.Sprint(v))
	return b
}

func (b *VSparklineBuilder) On(name string, value string) (r *VSparklineBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSparklineBuilder) Bind(name string, value string) (r *VSparklineBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
