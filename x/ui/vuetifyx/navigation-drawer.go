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

// Resizable lets the inner edge of the drawer drag its width — between
// MinWidth and MaxWidth —, kept in the browser's localStorage under
// storageKey ("": not kept); a double click on the edge is the width given
// again.
func (b *VXNavigationDrawerBuilder) Resizable(storageKey string) *VXNavigationDrawerBuilder {
	if storageKey == "" {
		return b.Attr("resizable", true)
	}
	return b.Attr("resizable", storageKey)
}

// MinWidth is the narrowest a resizable drawer is dragged to (px).
func (b *VXNavigationDrawerBuilder) MinWidth(px int) *VXNavigationDrawerBuilder {
	return b.Attr(":min-width", px)
}

// MaxWidth is the widest a resizable drawer is dragged to (px).
func (b *VXNavigationDrawerBuilder) MaxWidth(px int) *VXNavigationDrawerBuilder {
	return b.Attr(":max-width", px)
}
