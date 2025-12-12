package script

import (
	h "github.com/go-rvq/htmlgo"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

func init() {
	Messages_pt_BR.EditorUsage = h.HTMLComponents{
		vx.VXPortal(),
	}
}
