package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
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
