package labs

import (
	v "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type VTreeviewBuilder struct {
	v.VTagBuilder[*VTreeviewBuilder]
}

func Treeview(children ...h.HTMLComponent) *VTreeviewBuilder {
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
