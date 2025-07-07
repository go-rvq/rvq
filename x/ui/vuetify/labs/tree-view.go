package labs

import (
	h "github.com/go-rvq/htmlgo"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

type VTreeviewBuilder struct {
	v.VTagBuilder[*VTreeviewBuilder]
}

func VTreeview(children ...h.HTMLComponent) *VTreeviewBuilder {
	return v.VTag(&VTreeviewBuilder{}, "v-treeview", children...)
}

func (b *VTreeviewBuilder) Activatable(v bool) *VTreeviewBuilder {
	return b.Attr("activatable", v)
}

func (b *VTreeviewBuilder) Items(v interface{}) (r *VTreeviewBuilder) {
	return b.Attr(":items", h.JSONString(v))
}

func (b *VTreeviewBuilder) ItemsVar(v string) (r *VTreeviewBuilder) {
	return b.Attr(":items", v)
}
