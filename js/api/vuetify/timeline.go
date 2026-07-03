package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTimelineBuilder struct {
	VTagBuilder[*VTimelineBuilder]
}

func VTimeline(children ...h.HTMLComponent) *VTimelineBuilder {
	return VTag(&VTimelineBuilder{}, "v-timeline", children...)
}

func (b *VTimelineBuilder) Density(v interface{}) (r *VTimelineBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VTimelineBuilder) Tag(v interface{}) (r *VTimelineBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VTimelineBuilder) Theme(v string) (r *VTimelineBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTimelineBuilder) Size(v interface{}) (r *VTimelineBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VTimelineBuilder) Direction(v interface{}) (r *VTimelineBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VTimelineBuilder) IconColor(v string) (r *VTimelineBuilder) {
	b.Attr("icon-color", v)
	return b
}

func (b *VTimelineBuilder) Align(v interface{}) (r *VTimelineBuilder) {
	b.Attr(":align", h.JSONString(v))
	return b
}

func (b *VTimelineBuilder) Justify(v string) (r *VTimelineBuilder) {
	b.Attr("justify", v)
	return b
}

func (b *VTimelineBuilder) LineThickness(v interface{}) (r *VTimelineBuilder) {
	b.Attr(":line-thickness", h.JSONString(v))
	return b
}

func (b *VTimelineBuilder) LineColor(v string) (r *VTimelineBuilder) {
	b.Attr("line-color", v)
	return b
}

func (b *VTimelineBuilder) DotColor(v string) (r *VTimelineBuilder) {
	b.Attr("dot-color", v)
	return b
}

func (b *VTimelineBuilder) FillDot(v bool) (r *VTimelineBuilder) {
	b.Attr(":fill-dot", fmt.Sprint(v))
	return b
}

func (b *VTimelineBuilder) HideOpposite(v bool) (r *VTimelineBuilder) {
	b.Attr(":hide-opposite", fmt.Sprint(v))
	return b
}

func (b *VTimelineBuilder) LineInset(v interface{}) (r *VTimelineBuilder) {
	b.Attr(":line-inset", h.JSONString(v))
	return b
}

func (b *VTimelineBuilder) Side(v interface{}) (r *VTimelineBuilder) {
	b.Attr(":side", h.JSONString(v))
	return b
}

func (b *VTimelineBuilder) TruncateLine(v interface{}) (r *VTimelineBuilder) {
	b.Attr(":truncate-line", h.JSONString(v))
	return b
}

func (b *VTimelineBuilder) On(name string, value string) (r *VTimelineBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTimelineBuilder) Bind(name string, value string) (r *VTimelineBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VTimelineBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VTimelineBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VTimelineBuilder) Slot(name string, child ...h.HTMLComponent) (r *VTimelineBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VTimelineBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VTimelineBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VTimelineBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VTimelineBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VTimelineBuilder) SlotDefault(child ...h.HTMLComponent) (r *VTimelineBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VTimelineBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VTimelineBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
