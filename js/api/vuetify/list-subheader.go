package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VListSubheaderBuilder struct {
	VTagBuilder[*VListSubheaderBuilder]
}

func VListSubheader(children ...h.HTMLComponent) *VListSubheaderBuilder {
	return VTag(&VListSubheaderBuilder{}, "v-list-subheader", children...)
}

func (b *VListSubheaderBuilder) Tag(v interface{}) (r *VListSubheaderBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VListSubheaderBuilder) Title(v string) (r *VListSubheaderBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VListSubheaderBuilder) Color(v string) (r *VListSubheaderBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VListSubheaderBuilder) Inset(v bool) (r *VListSubheaderBuilder) {
	b.Attr(":inset", fmt.Sprint(v))
	return b
}

func (b *VListSubheaderBuilder) Sticky(v bool) (r *VListSubheaderBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VListSubheaderBuilder) On(name string, value string) (r *VListSubheaderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListSubheaderBuilder) Bind(name string, value string) (r *VListSubheaderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VListSubheaderBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VListSubheaderBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VListSubheaderBuilder) Slot(name string, child ...h.HTMLComponent) (r *VListSubheaderBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VListSubheaderBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VListSubheaderBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VListSubheaderBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VListSubheaderBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VListSubheaderBuilder) SlotDefault(child ...h.HTMLComponent) (r *VListSubheaderBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VListSubheaderBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VListSubheaderBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
