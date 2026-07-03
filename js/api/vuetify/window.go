package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VWindowBuilder struct {
	VTagBuilder[*VWindowBuilder]
}

func VWindow(children ...h.HTMLComponent) *VWindowBuilder {
	return VTag(&VWindowBuilder{}, "v-window", children...)
}

func (b *VWindowBuilder) ModelValue(v interface{}) (r *VWindowBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) Reverse(v bool) (r *VWindowBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VWindowBuilder) Tag(v interface{}) (r *VWindowBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) Theme(v string) (r *VWindowBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VWindowBuilder) Disabled(v bool) (r *VWindowBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VWindowBuilder) SelectedClass(v string) (r *VWindowBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VWindowBuilder) Mandatory(v interface{}) (r *VWindowBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) NextIcon(v interface{}) (r *VWindowBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) PrevIcon(v interface{}) (r *VWindowBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) Continuous(v bool) (r *VWindowBuilder) {
	b.Attr(":continuous", fmt.Sprint(v))
	return b
}

func (b *VWindowBuilder) ShowArrows(v interface{}) (r *VWindowBuilder) {
	b.Attr(":show-arrows", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) Touch(v interface{}) (r *VWindowBuilder) {
	b.Attr(":touch", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) Direction(v interface{}) (r *VWindowBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VWindowBuilder) On(name string, value string) (r *VWindowBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VWindowBuilder) Bind(name string, value string) (r *VWindowBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VWindowBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VWindowBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VWindowBuilder) Slot(name string, child ...h.HTMLComponent) (r *VWindowBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VWindowBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VWindowBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VWindowBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VWindowBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VWindowBuilder) SlotDefault(child ...h.HTMLComponent) (r *VWindowBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VWindowBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VWindowBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VWindowBuilder) SetSlotAdditional(child ...h.HTMLComponent) {
	b.SetSlot("additional", child...)
}

func (b *VWindowBuilder) SetScopedSlotAdditional(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("additional", scope, child...)
}

func (b *VWindowBuilder) SlotAdditional(child ...h.HTMLComponent) (r *VWindowBuilder) {
	b.SetSlotAdditional(child...)
	return b
}

func (b *VWindowBuilder) ScopedSlotAdditional(scope string, child ...h.HTMLComponent) (r *VWindowBuilder) {
	b.SetScopedSlotAdditional(scope, child...)
	return b
}

func (b *VWindowBuilder) SetSlotPrev(child ...h.HTMLComponent) {
	b.SetSlot("prev", child...)
}

func (b *VWindowBuilder) SetScopedSlotPrev(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prev", scope, child...)
}

func (b *VWindowBuilder) SlotPrev(child ...h.HTMLComponent) (r *VWindowBuilder) {
	b.SetSlotPrev(child...)
	return b
}

func (b *VWindowBuilder) ScopedSlotPrev(scope string, child ...h.HTMLComponent) (r *VWindowBuilder) {
	b.SetScopedSlotPrev(scope, child...)
	return b
}

func (b *VWindowBuilder) SetSlotNext(child ...h.HTMLComponent) {
	b.SetSlot("next", child...)
}

func (b *VWindowBuilder) SetScopedSlotNext(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("next", scope, child...)
}

func (b *VWindowBuilder) SlotNext(child ...h.HTMLComponent) (r *VWindowBuilder) {
	b.SetSlotNext(child...)
	return b
}

func (b *VWindowBuilder) ScopedSlotNext(scope string, child ...h.HTMLComponent) (r *VWindowBuilder) {
	b.SetScopedSlotNext(scope, child...)
	return b
}
