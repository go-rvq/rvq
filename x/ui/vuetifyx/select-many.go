package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
)

func VXSelectMany(children ...h.HTMLComponent) (r *VXAdvancedSelectBuilder) {
	return VXAdvancedSelect(children...).Many(true)
}
