package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VLigatureIconBuilder struct {
	VTagBuilder[*VLigatureIconBuilder]
}

func VLigatureIcon(children ...h.HTMLComponent) *VLigatureIconBuilder {
	return VTag(&VLigatureIconBuilder{}, "v-ligature-icon", children...)
}

func (b *VLigatureIconBuilder) Tag(v interface{}) (r *VLigatureIconBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VLigatureIconBuilder) Icon(v interface{}) (r *VLigatureIconBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VLigatureIconBuilder) On(name string, value string) (r *VLigatureIconBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VLigatureIconBuilder) Bind(name string, value string) (r *VLigatureIconBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VLigatureIconBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VLigatureIconBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VLigatureIconBuilder) Slot(name string, child ...h.HTMLComponent) (r *VLigatureIconBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VLigatureIconBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VLigatureIconBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}
