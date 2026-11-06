package vuetifyx

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

var (
	_ VXAdvancedCardTagGetter[*VXDialogBuilder]            = (*VXDialogBuilder)(nil)
	_ VXAdvancedCloseCardTagGetter[*VXDialogBuilder]       = (*VXDialogBuilder)(nil)
	_ VXAdvancedExpandCloseCardTagGetter[*VXDialogBuilder] = (*VXDialogBuilder)(nil)
)

type VXDialogBuilder struct {
	VXAdvancedExpandCloseCardTagBuilder[*VXDialogBuilder]
}

func VXDialog(children ...h.HTMLComponent) *VXDialogBuilder {
	return VXAdvancedExpandCloseCardTag(&VXDialogBuilder{}, "vx-dialog", children...)
}

// Expanded opens the dialog expanded (fullscreen); its expand button
// (Expandable) still turns it.
func (b *VXDialogBuilder) Expanded(v bool) *VXDialogBuilder {
	b.Attr(":expanded", fmt.Sprint(v))
	return b
}

// FullscreenOnMobile: on a phone (narrower than md, 960 px) the
// dialog is always the whole screen, with no expand button (true, the
// default); false, as on a wide screen.
func (b *VXDialogBuilder) FullscreenOnMobile(v bool) *VXDialogBuilder {
	b.Attr(":fullscreen-on-mobile", fmt.Sprint(v))
	return b
}
