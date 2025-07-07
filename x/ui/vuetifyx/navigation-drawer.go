package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
)

var (
	_ VXAdvancedCardTagGetter[*VXNavigationDrawerBuilder]            = (*VXNavigationDrawerBuilder)(nil)
	_ VXAdvancedCloseCardTagGetter[*VXNavigationDrawerBuilder]       = (*VXNavigationDrawerBuilder)(nil)
	_ VXAdvancedExpandCloseCardTagGetter[*VXNavigationDrawerBuilder] = (*VXNavigationDrawerBuilder)(nil)
)

type VXNavigationDrawerBuilder struct {
	VXAdvancedExpandCloseCardTagBuilder[*VXNavigationDrawerBuilder]
}

func VXNavigationDrawer(children ...h.HTMLComponent) *VXNavigationDrawerBuilder {
	return VXAdvancedExpandCloseCardTag(&VXNavigationDrawerBuilder{}, "vx-navigation-drawer", children...)
}

func (b *VXNavigationDrawerBuilder) Location(s string) *VXNavigationDrawerBuilder {
	return b.Attr("location", s)
}

func (b *VXNavigationDrawerBuilder) Density(s string) *VXNavigationDrawerBuilder {
	return b.Attr("density", s)
}

func (b *VXNavigationDrawerBuilder) Temporary(v bool) *VXNavigationDrawerBuilder {
	return b.Attr("temporary", v)
}

func (b *VXNavigationDrawerBuilder) Variant(v string) *VXNavigationDrawerBuilder {
	return b.Attr("variant", v)
}

func (b *VXNavigationDrawerBuilder) VariantMenu() *VXNavigationDrawerBuilder {
	return b.Variant("menu")
}
