package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VNoSsrBuilder struct {
	VTagBuilder[*VNoSsrBuilder]
}

func VNoSsr(children ...h.HTMLComponent) *VNoSsrBuilder {
	return VTag(&VNoSsrBuilder{}, "v-no-ssr", children...)
}

func (b *VNoSsrBuilder) On(name string, value string) (r *VNoSsrBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VNoSsrBuilder) Bind(name string, value string) (r *VNoSsrBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VNoSsrBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VNoSsrBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VNoSsrBuilder) Slot(name string, child ...h.HTMLComponent) (r *VNoSsrBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VNoSsrBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VNoSsrBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}
