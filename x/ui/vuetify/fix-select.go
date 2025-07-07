package vuetify

import (
	h "github.com/go-rvq/htmlgo"
)

func VSelect(children ...h.HTMLComponent) *VSelectBuilder {
	return VTag(&VSelectBuilder{}, "v-select", children...)
}

func (b *VSelectBuilder) Value(v interface{}) *VSelectBuilder {
	return b.ModelValue(v)
}
