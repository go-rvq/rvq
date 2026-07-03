package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VAppBarTitleBuilder struct {
	VTagBuilder[*VAppBarTitleBuilder]
}

func VAppBarTitle(children ...h.HTMLComponent) *VAppBarTitleBuilder {
	return VTag(&VAppBarTitleBuilder{}, "v-app-bar-title", children...)
}

func (b *VAppBarTitleBuilder) Tag(v interface{}) (r *VAppBarTitleBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VAppBarTitleBuilder) Text(v string) (r *VAppBarTitleBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VAppBarTitleBuilder) On(name string, value string) (r *VAppBarTitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAppBarTitleBuilder) Bind(name string, value string) (r *VAppBarTitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VAppBarTitleBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VAppBarTitleBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VAppBarTitleBuilder) Slot(name string, child ...h.HTMLComponent) (r *VAppBarTitleBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VAppBarTitleBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VAppBarTitleBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VAppBarTitleBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VAppBarTitleBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VAppBarTitleBuilder) SlotDefault(child ...h.HTMLComponent) (r *VAppBarTitleBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VAppBarTitleBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VAppBarTitleBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VAppBarTitleBuilder) SetSlotText(child ...h.HTMLComponent) {
	b.SetSlot("text", child...)
}

func (b *VAppBarTitleBuilder) SetScopedSlotText(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("text", scope, child...)
}

func (b *VAppBarTitleBuilder) SlotText(child ...h.HTMLComponent) (r *VAppBarTitleBuilder) {
	b.SetSlotText(child...)
	return b
}

func (b *VAppBarTitleBuilder) ScopedSlotText(scope string, child ...h.HTMLComponent) (r *VAppBarTitleBuilder) {
	b.SetScopedSlotText(scope, child...)
	return b
}
