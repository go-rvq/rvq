package vuetify

import (
	h "github.com/theplant/htmlgo"
)

func VAutocomplete(children ...h.HTMLComponent) *VAutocompleteBuilder {
	return VTag(&VAutocompleteBuilder{}, "v-autocomplete", children...)
}
