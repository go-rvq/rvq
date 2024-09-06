package vuetifyx

import (
	"context"

	v "github.com/qor5/x/v3/ui/vuetify"

	h "github.com/theplant/htmlgo"
)

type VXSelectBuilder struct {
	v.VTagBuilder[*VXSelectBuilder]
	selectedItems interface{}
	items         interface{}
}

func VXSelect(children ...h.HTMLComponent) *VXSelectBuilder {
	return v.VTag(&VXSelectBuilder{}, "vx-select", children...)
}

func (b *VXSelectBuilder) Items(v interface{}) (r *VXSelectBuilder) {
	b.items = v
	return b
}

func (b *VXSelectBuilder) SelectedItems(v interface{}) (r *VXSelectBuilder) {
	b.selectedItems = v
	return b
}

func (b *VXSelectBuilder) FieldName(v string) (r *VXSelectBuilder) {
	b.Attr("field-name", v)
	return b
}

func (b *VXSelectBuilder) MarshalHTML(ctx context.Context) (r []byte, err error) {
	if b.items == nil {
		b.items = b.selectedItems
	}
	b.Attr(":items", b.items)
	b.Attr(":selected-items", b.selectedItems)

	return b.GetHTMLTagBuilder().MarshalHTML(ctx)
}
