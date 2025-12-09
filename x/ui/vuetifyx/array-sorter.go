package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

type VXArraySorterBuilder struct {
	v.VTagBuilder[*VXArraySorterBuilder]
}

func (b *VXArraySorterBuilder) Label(v string) *VXArraySorterBuilder {
	return b.Attr("label", v)
}

func (b *VXArraySorterBuilder) Density(v string) *VXArraySorterBuilder {
	return b.Attr("density", v)
}

func (b *VXArraySorterBuilder) ItemKey(v string) *VXArraySorterBuilder {
	return b.Attr(":item-key", v)
}

func (b *VXArraySorterBuilder) ItemTitle(v string) *VXArraySorterBuilder {
	return b.Attr(":item-title", v)
}

func VXArraySorter(children ...h.HTMLComponent) *VXArraySorterBuilder {
	return v.VTag(&VXArraySorterBuilder{}, "vx-array-sorter", children...)
}
