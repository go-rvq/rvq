package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VListImgBuilder struct {
	VTagBuilder[*VListImgBuilder]
}

func VListImg(children ...h.HTMLComponent) *VListImgBuilder {
	return VTag(&VListImgBuilder{}, "v-list-img", children...)
}

func (b *VListImgBuilder) Tag(v string) (r *VListImgBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VListImgBuilder) On(name string, value string) (r *VListImgBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListImgBuilder) Bind(name string, value string) (r *VListImgBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VListImgBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VListImgBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VListImgBuilder) Slot(name string, child ...h.HTMLComponent) (r *VListImgBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VListImgBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VListImgBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VListImgBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VListImgBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VListImgBuilder) SlotDefault(child ...h.HTMLComponent) (r *VListImgBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VListImgBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VListImgBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
