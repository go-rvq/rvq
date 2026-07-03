package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSvgIconBuilder struct {
	VTagBuilder[*VSvgIconBuilder]
}

func VSvgIcon(children ...h.HTMLComponent) *VSvgIconBuilder {
	return VTag(&VSvgIconBuilder{}, "v-svg-icon", children...)
}

func (b *VSvgIconBuilder) Icon(v interface{}) (r *VSvgIconBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VSvgIconBuilder) Tag(v interface{}) (r *VSvgIconBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VSvgIconBuilder) On(name string, value string) (r *VSvgIconBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSvgIconBuilder) Bind(name string, value string) (r *VSvgIconBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VSvgIconBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSvgIconBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSvgIconBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSvgIconBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSvgIconBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSvgIconBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}
