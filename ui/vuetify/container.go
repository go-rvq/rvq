package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VContainerBuilder struct {
	VTagBuilder[*VContainerBuilder]
}

func VContainer(children ...h.HTMLComponent) *VContainerBuilder {
	return VTag(&VContainerBuilder{}, "v-container", children...)
}

func (b *VContainerBuilder) Fluid(v bool) (r *VContainerBuilder) {
	b.Attr(":fluid", fmt.Sprint(v))
	return b
}

func (b *VContainerBuilder) Tag(v string) (r *VContainerBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VContainerBuilder) On(name string, value string) (r *VContainerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VContainerBuilder) Bind(name string, value string) (r *VContainerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
