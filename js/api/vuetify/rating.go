package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VRatingBuilder struct {
	VTagBuilder[*VRatingBuilder]
}

func VRating(children ...h.HTMLComponent) *VRatingBuilder {
	return VTag(&VRatingBuilder{}, "v-rating", children...)
}

func (b *VRatingBuilder) Length(v interface{}) (r *VRatingBuilder) {
	b.Attr(":length", h.JSONString(v))
	return b
}

func (b *VRatingBuilder) ModelValue(v interface{}) (r *VRatingBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VRatingBuilder) Density(v interface{}) (r *VRatingBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VRatingBuilder) Tag(v interface{}) (r *VRatingBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VRatingBuilder) Theme(v string) (r *VRatingBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VRatingBuilder) Color(v string) (r *VRatingBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VRatingBuilder) Name(v string) (r *VRatingBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VRatingBuilder) ActiveColor(v string) (r *VRatingBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VRatingBuilder) Readonly(v bool) (r *VRatingBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VRatingBuilder) Ripple(v bool) (r *VRatingBuilder) {
	b.Attr(":ripple", fmt.Sprint(v))
	return b
}

func (b *VRatingBuilder) Disabled(v bool) (r *VRatingBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VRatingBuilder) Size(v interface{}) (r *VRatingBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VRatingBuilder) Hover(v bool) (r *VRatingBuilder) {
	b.Attr(":hover", fmt.Sprint(v))
	return b
}

func (b *VRatingBuilder) Clearable(v bool) (r *VRatingBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VRatingBuilder) ItemAriaLabel(v string) (r *VRatingBuilder) {
	b.Attr("item-aria-label", v)
	return b
}

func (b *VRatingBuilder) EmptyIcon(v interface{}) (r *VRatingBuilder) {
	b.Attr(":empty-icon", h.JSONString(v))
	return b
}

func (b *VRatingBuilder) FullIcon(v interface{}) (r *VRatingBuilder) {
	b.Attr(":full-icon", h.JSONString(v))
	return b
}

func (b *VRatingBuilder) HalfIncrements(v bool) (r *VRatingBuilder) {
	b.Attr(":half-increments", fmt.Sprint(v))
	return b
}

func (b *VRatingBuilder) ItemLabelPosition(v string) (r *VRatingBuilder) {
	b.Attr("item-label-position", v)
	return b
}

func (b *VRatingBuilder) ItemLabels(v interface{}) (r *VRatingBuilder) {
	b.Attr(":item-labels", h.JSONString(v))
	return b
}

func (b *VRatingBuilder) On(name string, value string) (r *VRatingBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VRatingBuilder) Bind(name string, value string) (r *VRatingBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VRatingBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VRatingBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VRatingBuilder) Slot(name string, child ...h.HTMLComponent) (r *VRatingBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VRatingBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VRatingBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VRatingBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VRatingBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VRatingBuilder) SlotItem(child ...h.HTMLComponent) (r *VRatingBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VRatingBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VRatingBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VRatingBuilder) SetSlotItemLabel(child ...h.HTMLComponent) {
	b.SetSlot("item-label", child...)
}

func (b *VRatingBuilder) SetScopedSlotItemLabel(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item-label", scope, child...)
}

func (b *VRatingBuilder) SlotItemLabel(child ...h.HTMLComponent) (r *VRatingBuilder) {
	b.SetSlotItemLabel(child...)
	return b
}

func (b *VRatingBuilder) ScopedSlotItemLabel(scope string, child ...h.HTMLComponent) (r *VRatingBuilder) {
	b.SetScopedSlotItemLabel(scope, child...)
	return b
}
