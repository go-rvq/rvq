package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VKbdBuilder struct {
	VTagBuilder[*VKbdBuilder]
}

func VKbd(children ...h.HTMLComponent) *VKbdBuilder {
	return VTag(&VKbdBuilder{}, "v-kbd", children...)
}

func (b *VKbdBuilder) Tag(v string) (r *VKbdBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VKbdBuilder) On(name string, value string) (r *VKbdBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VKbdBuilder) Bind(name string, value string) (r *VKbdBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
