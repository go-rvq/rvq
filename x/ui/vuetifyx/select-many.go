package vuetifyx

import (
	h "github.com/theplant/htmlgo"
)

func VXSelectMany(children ...h.HTMLComponent) (r *VXAdvancedSelectBuilder) {
	return VXAdvancedSelect(children...).Many(true)
}
