package vuetify

import (
	h "github.com/theplant/htmlgo"
)

func VTable(children ...h.HTMLComponent) *VTableBuilder {
	return VTag(
		&VTableBuilder{},
		"v-table",
		h.Template(children...).
			Attr("#default", true),
	)
}
