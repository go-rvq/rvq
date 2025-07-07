package vuetifyx

import (
	h "github.com/theplant/htmlgo"
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
