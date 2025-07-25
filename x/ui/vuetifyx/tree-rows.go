package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

type VXTreeRowsBuilder struct {
	v.VTagBuilder[*VXTreeRowsBuilder]
	items interface{}
	many  bool
}

func VXTreeRows(children ...h.HTMLComponent) *VXTreeRowsBuilder {
	return v.VTag(&VXTreeRowsBuilder{}, "vx-tree-rows", children...)
}

func (b *VXTreeRowsBuilder) Items(items interface{}) *VXTreeRowsBuilder {
	b.Attr(":items")
	return b
}
