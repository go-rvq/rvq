package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCarouselBuilder struct {
	VTagBuilder[*VCarouselBuilder]
}

func VCarousel(children ...h.HTMLComponent) *VCarouselBuilder {
	return VTag(&VCarouselBuilder{}, "v-carousel", children...)
}

func (b *VCarouselBuilder) ModelValue(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Reverse(v bool) (r *VCarouselBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VCarouselBuilder) Height(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Tag(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Theme(v string) (r *VCarouselBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VCarouselBuilder) Color(v string) (r *VCarouselBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VCarouselBuilder) Disabled(v bool) (r *VCarouselBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCarouselBuilder) SelectedClass(v string) (r *VCarouselBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VCarouselBuilder) Progress(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":progress", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Mandatory(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) NextIcon(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) PrevIcon(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Interval(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":interval", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Cycle(v bool) (r *VCarouselBuilder) {
	b.Attr(":cycle", fmt.Sprint(v))
	return b
}

func (b *VCarouselBuilder) DelimiterIcon(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":delimiter-icon", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) HideDelimiters(v bool) (r *VCarouselBuilder) {
	b.Attr(":hide-delimiters", fmt.Sprint(v))
	return b
}

func (b *VCarouselBuilder) HideDelimiterBackground(v bool) (r *VCarouselBuilder) {
	b.Attr(":hide-delimiter-background", fmt.Sprint(v))
	return b
}

func (b *VCarouselBuilder) Continuous(v bool) (r *VCarouselBuilder) {
	b.Attr(":continuous", fmt.Sprint(v))
	return b
}

func (b *VCarouselBuilder) ShowArrows(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":show-arrows", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Touch(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":touch", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) Direction(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) VerticalDelimiters(v interface{}) (r *VCarouselBuilder) {
	b.Attr(":vertical-delimiters", h.JSONString(v))
	return b
}

func (b *VCarouselBuilder) On(name string, value string) (r *VCarouselBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCarouselBuilder) Bind(name string, value string) (r *VCarouselBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VCarouselBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCarouselBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCarouselBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCarouselBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCarouselBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCarouselBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCarouselBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCarouselBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCarouselBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCarouselBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCarouselBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCarouselBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VCarouselBuilder) SetSlotAdditional(child ...h.HTMLComponent) {
	b.SetSlot("additional", child...)
}

func (b *VCarouselBuilder) SetScopedSlotAdditional(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("additional", scope, child...)
}

func (b *VCarouselBuilder) SlotAdditional(child ...h.HTMLComponent) (r *VCarouselBuilder) {
	b.SetSlotAdditional(child...)
	return b
}

func (b *VCarouselBuilder) ScopedSlotAdditional(scope string, child ...h.HTMLComponent) (r *VCarouselBuilder) {
	b.SetScopedSlotAdditional(scope, child...)
	return b
}

func (b *VCarouselBuilder) SetSlotPrev(child ...h.HTMLComponent) {
	b.SetSlot("prev", child...)
}

func (b *VCarouselBuilder) SetScopedSlotPrev(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prev", scope, child...)
}

func (b *VCarouselBuilder) SlotPrev(child ...h.HTMLComponent) (r *VCarouselBuilder) {
	b.SetSlotPrev(child...)
	return b
}

func (b *VCarouselBuilder) ScopedSlotPrev(scope string, child ...h.HTMLComponent) (r *VCarouselBuilder) {
	b.SetScopedSlotPrev(scope, child...)
	return b
}

func (b *VCarouselBuilder) SetSlotNext(child ...h.HTMLComponent) {
	b.SetSlot("next", child...)
}

func (b *VCarouselBuilder) SetScopedSlotNext(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("next", scope, child...)
}

func (b *VCarouselBuilder) SlotNext(child ...h.HTMLComponent) (r *VCarouselBuilder) {
	b.SetSlotNext(child...)
	return b
}

func (b *VCarouselBuilder) ScopedSlotNext(scope string, child ...h.HTMLComponent) (r *VCarouselBuilder) {
	b.SetScopedSlotNext(scope, child...)
	return b
}

func (b *VCarouselBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VCarouselBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VCarouselBuilder) SlotItem(child ...h.HTMLComponent) (r *VCarouselBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VCarouselBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VCarouselBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}
