package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCarouselItemBuilder struct {
	VTagBuilder[*VCarouselItemBuilder]
}

func VCarouselItem(children ...h.HTMLComponent) *VCarouselItemBuilder {
	return VTag(&VCarouselItemBuilder{}, "v-carousel-item", children...)
}

func (b *VCarouselItemBuilder) Height(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) MaxHeight(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) MaxWidth(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) MinHeight(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) MinWidth(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Width(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Position(v string) (r *VCarouselItemBuilder) {
	b.Attr("position", v)
	return b
}

func (b *VCarouselItemBuilder) Absolute(v bool) (r *VCarouselItemBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VCarouselItemBuilder) Rounded(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Tile(v bool) (r *VCarouselItemBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VCarouselItemBuilder) Color(v string) (r *VCarouselItemBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VCarouselItemBuilder) Value(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Disabled(v bool) (r *VCarouselItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCarouselItemBuilder) SelectedClass(v string) (r *VCarouselItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VCarouselItemBuilder) Inline(v bool) (r *VCarouselItemBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VCarouselItemBuilder) Transition(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) ContentClass(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":content-class", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Eager(v bool) (r *VCarouselItemBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VCarouselItemBuilder) Alt(v string) (r *VCarouselItemBuilder) {
	b.Attr("alt", v)
	return b
}

func (b *VCarouselItemBuilder) Cover(v bool) (r *VCarouselItemBuilder) {
	b.Attr(":cover", fmt.Sprint(v))
	return b
}

func (b *VCarouselItemBuilder) Draggable(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":draggable", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Gradient(v string) (r *VCarouselItemBuilder) {
	b.Attr("gradient", v)
	return b
}

func (b *VCarouselItemBuilder) LazySrc(v string) (r *VCarouselItemBuilder) {
	b.Attr("lazy-src", v)
	return b
}

func (b *VCarouselItemBuilder) Options(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":options", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Sizes(v string) (r *VCarouselItemBuilder) {
	b.Attr("sizes", v)
	return b
}

func (b *VCarouselItemBuilder) Src(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":src", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Srcset(v string) (r *VCarouselItemBuilder) {
	b.Attr("srcset", v)
	return b
}

func (b *VCarouselItemBuilder) AspectRatio(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":aspect-ratio", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Crossorigin(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":crossorigin", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) Referrerpolicy(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":referrerpolicy", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) ReverseTransition(v interface{}) (r *VCarouselItemBuilder) {
	b.Attr(":reverse-transition", h.JSONString(v))
	return b
}

func (b *VCarouselItemBuilder) On(name string, value string) (r *VCarouselItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCarouselItemBuilder) Bind(name string, value string) (r *VCarouselItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VCarouselItemBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCarouselItemBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCarouselItemBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCarouselItemBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCarouselItemBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCarouselItemBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCarouselItemBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCarouselItemBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCarouselItemBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCarouselItemBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCarouselItemBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCarouselItemBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VCarouselItemBuilder) SetSlotPlaceholder(child ...h.HTMLComponent) {
	b.SetSlot("placeholder", child...)
}

func (b *VCarouselItemBuilder) SetScopedSlotPlaceholder(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("placeholder", scope, child...)
}

func (b *VCarouselItemBuilder) SlotPlaceholder(child ...h.HTMLComponent) (r *VCarouselItemBuilder) {
	b.SetSlotPlaceholder(child...)
	return b
}

func (b *VCarouselItemBuilder) ScopedSlotPlaceholder(scope string, child ...h.HTMLComponent) (r *VCarouselItemBuilder) {
	b.SetScopedSlotPlaceholder(scope, child...)
	return b
}

func (b *VCarouselItemBuilder) SetSlotError(child ...h.HTMLComponent) {
	b.SetSlot("error", child...)
}

func (b *VCarouselItemBuilder) SetScopedSlotError(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("error", scope, child...)
}

func (b *VCarouselItemBuilder) SlotError(child ...h.HTMLComponent) (r *VCarouselItemBuilder) {
	b.SetSlotError(child...)
	return b
}

func (b *VCarouselItemBuilder) ScopedSlotError(scope string, child ...h.HTMLComponent) (r *VCarouselItemBuilder) {
	b.SetScopedSlotError(scope, child...)
	return b
}

func (b *VCarouselItemBuilder) SetSlotSources(child ...h.HTMLComponent) {
	b.SetSlot("sources", child...)
}

func (b *VCarouselItemBuilder) SetScopedSlotSources(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("sources", scope, child...)
}

func (b *VCarouselItemBuilder) SlotSources(child ...h.HTMLComponent) (r *VCarouselItemBuilder) {
	b.SetSlotSources(child...)
	return b
}

func (b *VCarouselItemBuilder) ScopedSlotSources(scope string, child ...h.HTMLComponent) (r *VCarouselItemBuilder) {
	b.SetScopedSlotSources(scope, child...)
	return b
}
