package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCardSubtitleBuilder struct {
	VTagBuilder[*VCardSubtitleBuilder]
}

func VCardSubtitle(children ...h.HTMLComponent) *VCardSubtitleBuilder {
	return VTag(&VCardSubtitleBuilder{}, "v-card-subtitle", children...)
}

func (b *VCardSubtitleBuilder) Tag(v interface{}) (r *VCardSubtitleBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VCardSubtitleBuilder) Opacity(v interface{}) (r *VCardSubtitleBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VCardSubtitleBuilder) On(name string, value string) (r *VCardSubtitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCardSubtitleBuilder) Bind(name string, value string) (r *VCardSubtitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VCardSubtitleBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCardSubtitleBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCardSubtitleBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCardSubtitleBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCardSubtitleBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCardSubtitleBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCardSubtitleBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCardSubtitleBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCardSubtitleBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCardSubtitleBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCardSubtitleBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCardSubtitleBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
