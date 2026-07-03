package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VExpandXtransitionBuilder struct {
	VTagBuilder[*VExpandXtransitionBuilder]
}

func VExpandXtransition(children ...h.HTMLComponent) *VExpandXtransitionBuilder {
	return VTag(&VExpandXtransitionBuilder{}, "v-expand-xtransition", children...)
}

func (b *VExpandXtransitionBuilder) Mode(v interface{}) (r *VExpandXtransitionBuilder) {
	b.Attr(":mode", h.JSONString(v))
	return b
}

func (b *VExpandXtransitionBuilder) Disabled(v bool) (r *VExpandXtransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VExpandXtransitionBuilder) Group(v bool) (r *VExpandXtransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VExpandXtransitionBuilder) On(name string, value string) (r *VExpandXtransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VExpandXtransitionBuilder) Bind(name string, value string) (r *VExpandXtransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VExpandXtransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VExpandXtransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VExpandXtransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VExpandXtransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VExpandXtransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VExpandXtransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VExpandXtransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VExpandXtransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VExpandXtransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VExpandXtransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VExpandXtransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VExpandXtransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
