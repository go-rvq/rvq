package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VVirtualScrollBuilder struct {
	VTagBuilder[*VVirtualScrollBuilder]
}

func VVirtualScroll(children ...h.HTMLComponent) *VVirtualScrollBuilder {
	return VTag(&VVirtualScrollBuilder{}, "v-virtual-scroll", children...)
}

func (b *VVirtualScrollBuilder) Items(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) Height(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) MaxHeight(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) MaxWidth(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) MinHeight(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) MinWidth(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) Width(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) ItemHeight(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":item-height", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) ItemKey(v interface{}) (r *VVirtualScrollBuilder) {
	b.Attr(":item-key", h.JSONString(v))
	return b
}

func (b *VVirtualScrollBuilder) Renderless(v bool) (r *VVirtualScrollBuilder) {
	b.Attr(":renderless", fmt.Sprint(v))
	return b
}

func (b *VVirtualScrollBuilder) On(name string, value string) (r *VVirtualScrollBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VVirtualScrollBuilder) Bind(name string, value string) (r *VVirtualScrollBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VVirtualScrollBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VVirtualScrollBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VVirtualScrollBuilder) Slot(name string, child ...h.HTMLComponent) (r *VVirtualScrollBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VVirtualScrollBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VVirtualScrollBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VVirtualScrollBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VVirtualScrollBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VVirtualScrollBuilder) SlotDefault(child ...h.HTMLComponent) (r *VVirtualScrollBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VVirtualScrollBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VVirtualScrollBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
