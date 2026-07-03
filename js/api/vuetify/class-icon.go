package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VClassIconBuilder struct {
	VTagBuilder[*VClassIconBuilder]
}

func VClassIcon(children ...h.HTMLComponent) *VClassIconBuilder {
	return VTag(&VClassIconBuilder{}, "v-class-icon", children...)
}

func (b *VClassIconBuilder) Icon(v interface{}) (r *VClassIconBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VClassIconBuilder) Tag(v interface{}) (r *VClassIconBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VClassIconBuilder) On(name string, value string) (r *VClassIconBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VClassIconBuilder) Bind(name string, value string) (r *VClassIconBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VClassIconBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VClassIconBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VClassIconBuilder) Slot(name string, child ...h.HTMLComponent) (r *VClassIconBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VClassIconBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VClassIconBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}
