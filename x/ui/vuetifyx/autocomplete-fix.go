package vuetifyx

import (
	v "github.com/go-rvq/rvq/x/ui/vuetify"

	h "github.com/go-rvq/htmlgo"
)

type VXAutocompleteBuilder struct {
	v.VTagBuilder[*VXAutocompleteBuilder]
	selectedItems interface{}
	items         interface{}
}

func VXAutocomplete(children ...h.HTMLComponent) *VXAutocompleteBuilder {
	return v.VTag(&VXAutocompleteBuilder{}, "vx-autocomplete", children...).
		Multiple(true).
		Chips(true).
		DeletableChips(true).
		Clearable(true)
}

func (b *VXAutocompleteBuilder) SelectedItems(v interface{}) (r *VXAutocompleteBuilder) {
	b.selectedItems = v
	return b
}

func (b *VXAutocompleteBuilder) HasIcon(v bool) (r *VXAutocompleteBuilder) {
	b.Attr("has-icon", v)
	return b
}

func (b *VXAutocompleteBuilder) Sorting(v bool) (r *VXAutocompleteBuilder) {
	b.Attr("sorting", v)
	return b
}

func (b *VXAutocompleteBuilder) Variant(v string) (r *VXAutocompleteBuilder) {
	b.Attr("variant", v)
	return b
}

func (b *VXAutocompleteBuilder) Density(v string) (r *VXAutocompleteBuilder) {
	b.Attr("density", v)
	return b
}

func (b *VXAutocompleteBuilder) Items(v interface{}) (r *VXAutocompleteBuilder) {
	b.items = v
	return b
}

func (b *VXAutocompleteBuilder) ChipColor(v string) (r *VXAutocompleteBuilder) {
	b.Attr("chip-color", v)
	return b
}

func (b *VXAutocompleteBuilder) ChipTextColor(v string) (r *VXAutocompleteBuilder) {
	b.Attr("chip-text-color", v)
	return b
}

func (b *VXAutocompleteBuilder) SetDataSource(ds *AutocompleteDataSource) (r *VXAutocompleteBuilder) {
	b.Attr("remote-url", ds.RemoteURL)
	b.Attr("event-name", ds.EventName)
	b.Attr("is-paging", ds.IsPaging)
	b.Attr("has-icon", ds.HasIcon)
	return b
}

func (b *VXAutocompleteBuilder) Write(ctx *h.Context) (err error) {
	if b.items == nil {
		b.items = b.selectedItems
	}
	b.Attr(":items", b.items)
	b.Attr(":selected-items", b.selectedItems)
	return b.GetHTMLTagBuilder().Write(ctx)
}

type AutocompleteDataSource struct {
	RemoteURL string `json:"remote-url"`
	EventName string `json:"event-name"`
	IsPaging  bool   `json:"is-paging"`
	HasIcon   bool   `json:"has-icon"`
}
