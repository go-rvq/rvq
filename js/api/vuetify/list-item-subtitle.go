package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VListItemSubtitleBuilder struct {
	VTagBuilder[*VListItemSubtitleBuilder]
}

func VListItemSubtitle(children ...h.HTMLComponent) *VListItemSubtitleBuilder {
	return VTag(&VListItemSubtitleBuilder{}, "v-list-item-subtitle", children...)
}

func (b *VListItemSubtitleBuilder) Tag(v interface{}) (r *VListItemSubtitleBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VListItemSubtitleBuilder) Opacity(v interface{}) (r *VListItemSubtitleBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VListItemSubtitleBuilder) On(name string, value string) (r *VListItemSubtitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListItemSubtitleBuilder) Bind(name string, value string) (r *VListItemSubtitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VListItemSubtitleBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VListItemSubtitleBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VListItemSubtitleBuilder) Slot(name string, child ...h.HTMLComponent) (r *VListItemSubtitleBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VListItemSubtitleBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VListItemSubtitleBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VListItemSubtitleBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VListItemSubtitleBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VListItemSubtitleBuilder) SlotDefault(child ...h.HTMLComponent) (r *VListItemSubtitleBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VListItemSubtitleBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VListItemSubtitleBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
