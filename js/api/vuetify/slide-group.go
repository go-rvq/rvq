package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSlideGroupBuilder struct {
	VTagBuilder[*VSlideGroupBuilder]
}

func VSlideGroup(children ...h.HTMLComponent) *VSlideGroupBuilder {
	return VTag(&VSlideGroupBuilder{}, "v-slide-group", children...)
}

func (b *VSlideGroupBuilder) Symbol(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":symbol", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) Tag(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) Disabled(v bool) (r *VSlideGroupBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSlideGroupBuilder) Multiple(v bool) (r *VSlideGroupBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VSlideGroupBuilder) Mandatory(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) ModelValue(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) Direction(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) Max(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) Mobile(v bool) (r *VSlideGroupBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VSlideGroupBuilder) MobileBreakpoint(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) PrevIcon(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) NextIcon(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) SelectedClass(v string) (r *VSlideGroupBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VSlideGroupBuilder) CenterActive(v bool) (r *VSlideGroupBuilder) {
	b.Attr(":center-active", fmt.Sprint(v))
	return b
}

func (b *VSlideGroupBuilder) ShowArrows(v interface{}) (r *VSlideGroupBuilder) {
	b.Attr(":show-arrows", h.JSONString(v))
	return b
}

func (b *VSlideGroupBuilder) On(name string, value string) (r *VSlideGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSlideGroupBuilder) Bind(name string, value string) (r *VSlideGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VSlideGroupBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSlideGroupBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSlideGroupBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSlideGroupBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSlideGroupBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSlideGroupBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VSlideGroupBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VSlideGroupBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VSlideGroupBuilder) SlotDefault(child ...h.HTMLComponent) (r *VSlideGroupBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VSlideGroupBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VSlideGroupBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VSlideGroupBuilder) SetSlotPrev(child ...h.HTMLComponent) {
	b.SetSlot("prev", child...)
}

func (b *VSlideGroupBuilder) SetScopedSlotPrev(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prev", scope, child...)
}

func (b *VSlideGroupBuilder) SlotPrev(child ...h.HTMLComponent) (r *VSlideGroupBuilder) {
	b.SetSlotPrev(child...)
	return b
}

func (b *VSlideGroupBuilder) ScopedSlotPrev(scope string, child ...h.HTMLComponent) (r *VSlideGroupBuilder) {
	b.SetScopedSlotPrev(scope, child...)
	return b
}

func (b *VSlideGroupBuilder) SetSlotNext(child ...h.HTMLComponent) {
	b.SetSlot("next", child...)
}

func (b *VSlideGroupBuilder) SetScopedSlotNext(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("next", scope, child...)
}

func (b *VSlideGroupBuilder) SlotNext(child ...h.HTMLComponent) (r *VSlideGroupBuilder) {
	b.SetSlotNext(child...)
	return b
}

func (b *VSlideGroupBuilder) ScopedSlotNext(scope string, child ...h.HTMLComponent) (r *VSlideGroupBuilder) {
	b.SetScopedSlotNext(scope, child...)
	return b
}
