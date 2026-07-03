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

func (b *VSparklineBuilder) ModelValue(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) Fill(v bool) (r *VSparklineBuilder) {
	b.Attr(":fill", fmt.Sprint(v))
	return b
}

func (b *VSparklineBuilder) Height(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) Width(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) Color(v string) (r *VSparklineBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSparklineBuilder) Max(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) ID(v string) (r *VSparklineBuilder) {
	b.Attr("id", v)
	return b
}

func (b *VSparklineBuilder) Min(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) Gradient(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":gradient", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) ItemValue(v string) (r *VSparklineBuilder) {
	b.Attr("item-value", v)
	return b
}

func (b *VSparklineBuilder) Labels(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":labels", h.JSONString(v))
	return b
}

func (b *VSparklineBuilder) Padding(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":padding", h.JSONString(v))
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

func (b *VSparklineBuilder) GradientDirection(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":gradient-direction", h.JSONString(v))
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

func (b *VSparklineBuilder) ShowLabels(v bool) (r *VSparklineBuilder) {
	b.Attr(":show-labels", fmt.Sprint(v))
	return b
}

func (b *VSparklineBuilder) Smooth(v interface{}) (r *VSparklineBuilder) {
	b.Attr(":smooth", h.JSONString(v))
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

func (b *VSparklineBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSparklineBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSparklineBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSparklineBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSparklineBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSparklineBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VSparklineBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VSparklineBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VSparklineBuilder) SlotDefault(child ...h.HTMLComponent) (r *VSparklineBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VSparklineBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VSparklineBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VSparklineBuilder) SetSlotLabel(child ...h.HTMLComponent) {
	b.SetSlot("label", child...)
}

func (b *VSparklineBuilder) SetScopedSlotLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("label", scope, child...)
}

func (b *VSparklineBuilder) SlotLabel(child ...h.HTMLComponent) (r *VSparklineBuilder) {
	b.SetSlotLabel(child...)
	return b
}

func (b *VSparklineBuilder) ScopedSlotLabel(scope string, child ...h.HTMLComponent) (r *VSparklineBuilder) {
	b.SetScopedSlotLabel(scope, child...)
	return b
}
