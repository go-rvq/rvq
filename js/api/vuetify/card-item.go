package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCardItemBuilder struct {
	VTagBuilder[*VCardItemBuilder]
}

func VCardItem(children ...h.HTMLComponent) *VCardItemBuilder {
	return VTag(&VCardItemBuilder{}, "v-card-item", children...)
}

func (b *VCardItemBuilder) Title(v interface{}) (r *VCardItemBuilder) {
	b.Attr(":title", h.JSONString(v))
	return b
}

func (b *VCardItemBuilder) Density(v interface{}) (r *VCardItemBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VCardItemBuilder) PrependIcon(v interface{}) (r *VCardItemBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VCardItemBuilder) AppendIcon(v interface{}) (r *VCardItemBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VCardItemBuilder) Subtitle(v interface{}) (r *VCardItemBuilder) {
	b.Attr(":subtitle", h.JSONString(v))
	return b
}

func (b *VCardItemBuilder) AppendAvatar(v string) (r *VCardItemBuilder) {
	b.Attr("append-avatar", v)
	return b
}

func (b *VCardItemBuilder) PrependAvatar(v string) (r *VCardItemBuilder) {
	b.Attr("prepend-avatar", v)
	return b
}

func (b *VCardItemBuilder) On(name string, value string) (r *VCardItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCardItemBuilder) Bind(name string, value string) (r *VCardItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VCardItemBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCardItemBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCardItemBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCardItemBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCardItemBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCardItemBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCardItemBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCardItemBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCardItemBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCardItemBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCardItemBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCardItemBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VCardItemBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VCardItemBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VCardItemBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VCardItemBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VCardItemBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VCardItemBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VCardItemBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VCardItemBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VCardItemBuilder) SlotAppend(child ...h.HTMLComponent) (r *VCardItemBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VCardItemBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VCardItemBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VCardItemBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VCardItemBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VCardItemBuilder) SlotTitle(child ...h.HTMLComponent) (r *VCardItemBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VCardItemBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VCardItemBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VCardItemBuilder) SetSlotSubtitle(child ...h.HTMLComponent) {
	b.SetSlot("subtitle", child...)
}

func (b *VCardItemBuilder) SetScopedSlotSubtitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("subtitle", scope, child...)
}

func (b *VCardItemBuilder) SlotSubtitle(child ...h.HTMLComponent) (r *VCardItemBuilder) {
	b.SetSlotSubtitle(child...)
	return b
}

func (b *VCardItemBuilder) ScopedSlotSubtitle(scope string, child ...h.HTMLComponent) (r *VCardItemBuilder) {
	b.SetScopedSlotSubtitle(scope, child...)
	return b
}
