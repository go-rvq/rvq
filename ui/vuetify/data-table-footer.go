package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VDataTableFooterBuilder struct {
	VTagBuilder[*VDataTableFooterBuilder]
}

func VDataTableFooter(children ...h.HTMLComponent) *VDataTableFooterBuilder {
	return VTag(&VDataTableFooterBuilder{}, "v-data-table-footer", children...)
}

func (b *VDataTableFooterBuilder) On(name string, value string) (r *VDataTableFooterBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableFooterBuilder) Bind(name string, value string) (r *VDataTableFooterBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
