package vuetifyx

import (
	"context"
	"fmt"

	v "github.com/qor5/x/v3/ui/vuetify"

	"github.com/qor5/web/v3"

	h "github.com/theplant/htmlgo"
)

type VXSelectManyBuilder struct {
	v.VTagBuilder[*VXSelectManyBuilder]
	selectedItems   interface{}
	items           interface{}
	searchItemsFunc string
	itemsSearcher   *web.VueEventTagBuilder
}

func VXSelectMany(children ...h.HTMLComponent) (r *VXSelectManyBuilder) {
	return v.VTag(&VXSelectManyBuilder{}, "vx-selectmany", children...)
}

func (b *VXSelectManyBuilder) Items(v interface{}) (r *VXSelectManyBuilder) {
	b.items = v
	return b
}

func (b *VXSelectManyBuilder) SelectedItems(v interface{}) (r *VXSelectManyBuilder) {
	b.selectedItems = v
	return b
}

func (b *VXSelectManyBuilder) SearchItemsFunc(v string) (r *VXSelectManyBuilder) {
	b.searchItemsFunc = v
	return b
}

func (b *VXSelectManyBuilder) ItemsSearcher(eb *web.VueEventTagBuilder) (r *VXSelectManyBuilder) {
	b.itemsSearcher = eb
	return b
}

func (b *VXSelectManyBuilder) ItemText(v string) (r *VXSelectManyBuilder) {
	b.Attr("item-text", v)
	return b
}

func (b *VXSelectManyBuilder) ItemValue(v string) (r *VXSelectManyBuilder) {
	b.Attr("item-value", v)
	return b
}

func (b *VXSelectManyBuilder) Label(v string) (r *VXSelectManyBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VXSelectManyBuilder) AddItemLabel(v string) (r *VXSelectManyBuilder) {
	b.Attr("add-item-label", v)
	return b
}

func (b *VXSelectManyBuilder) MarshalHTML(ctx context.Context) (r []byte, err error) {
	if b.itemsSearcher != nil {
		b.Attr(":search-items-func", fmt.Sprintf(`function(val){return %s.query("keyword", val).json()}`, b.itemsSearcher.String()))
	} else if b.searchItemsFunc != "" {
		b.Attr(":search-items-func", fmt.Sprintf(`function(val){return $plaid().eventFunc("%s").query("keyword", val).json()}`, b.searchItemsFunc))
	} else {
		b.Attr(":items", b.items)
	}

	b.Attr(":selected-items", b.selectedItems)
	return b.GetHTMLTagBuilder().MarshalHTML(ctx)
}
