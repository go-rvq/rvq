package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

type VXPortalBuilder struct {
	v.VTagBuilder[*VXPortalBuilder]
}

func VXPortal(children ...h.HTMLComponent) *VXPortalBuilder {
	return v.VTag(&VXPortalBuilder{}, "vx-portal", children...)
}

func (b *VXPortalBuilder) DefaultSlot(children ...h.HTMLComponent) *VXPortalBuilder {
	return b.Children(h.Template(children...).Attr("v-pre", true))
}
