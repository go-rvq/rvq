package vuetifyx

import (
	"context"
	"fmt"

	v "github.com/qor5/x/v3/ui/vuetify"

	"github.com/qor5/web/v3"

	h "github.com/theplant/htmlgo"
)

type VXSelectOneBuilder struct {
	v.VTagBuilder[*VXSelectOneBuilder]
	selectedItem    interface{}
	items           interface{}
	searchItemsFunc string
	itemsSearcher   *web.VueEventTagBuilder
}

func VXSelectOne(children ...h.HTMLComponent) (r *VXSelectOneBuilder) {
	return v.VTag(&VXSelectOneBuilder{}, "vx-selectone", children...)
}

func (b *VXSelectOneBuilder) Items(v interface{}) (r *VXSelectOneBuilder) {
	b.items = v
	return b
}

func (b *VXSelectOneBuilder) SelectedItem(v interface{}) (r *VXSelectOneBuilder) {
	b.selectedItem = v
	return b
}

func (b *VXSelectOneBuilder) SearchItemsFunc(v string) (r *VXSelectOneBuilder) {
	b.searchItemsFunc = v
	return b
}

func (b *VXSelectOneBuilder) ItemsSearcher(eb *web.VueEventTagBuilder) (r *VXSelectOneBuilder) {
	b.itemsSearcher = eb
	return b
}

func (b *VXSelectOneBuilder) ItemText(v string) (r *VXSelectOneBuilder) {
	b.Attr("item-text", v)
	return b
}

func (b *VXSelectOneBuilder) ItemTextExpr(v string) (r *VXSelectOneBuilder) {
	b.Attr(":item-text", v)
	return b
}

func (b *VXSelectOneBuilder) ItemValue(v string) (r *VXSelectOneBuilder) {
	b.Attr("item-value", v)
	return b
}

func (b *VXSelectOneBuilder) Label(v string) (r *VXSelectOneBuilder) {
	b.Attr("label", v)
	return b
}

func (b *VXSelectOneBuilder) AddItemLabel(v string) (r *VXSelectOneBuilder) {
	b.Attr("add-item-label", v)
	return b
}

func (b *VXSelectOneBuilder) MarshalHTML(ctx context.Context) (r []byte, err error) {
	if b.itemsSearcher != nil {
		b.Attr(":search-items-func", fmt.Sprintf(`function(val){return %s.query("keyword", val).json()}`, b.itemsSearcher.String()))
	} else if b.searchItemsFunc != "" {
		b.Attr(":search-items-func", fmt.Sprintf(`function(val){return $plaid().eventFunc("%s").query("keyword", val).json()}`, b.searchItemsFunc))
	} else {
		b.Attr(":items", b.items)
	}

	var items []any
	if b.selectedItem != nil {
		items = append(items, b.selectedItem)
	}

	b.Attr(":selected-items", items)
	return b.GetHTMLTagBuilder().MarshalHTML(ctx)
}
