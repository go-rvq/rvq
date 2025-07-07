package vuetify

import (
	h "github.com/go-rvq/htmlgo"
)

func VAutocomplete(children ...h.HTMLComponent) *VAutocompleteBuilder {
	return VTag(&VAutocompleteBuilder{}, "v-autocomplete", children...)
}
