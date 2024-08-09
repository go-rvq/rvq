package vuetifyx

import (
	"context"
	"fmt"

	"github.com/qor5/web/v3"
	h "github.com/theplant/htmlgo"
)

type VXSelectManyBuilder struct {
	tag             *h.HTMLTagBuilder
	selectedItems   interface{}
	items           interface{}
	searchItemsFunc string
	itemsSearcher   *web.VueEventTagBuilder
}

func VXSelectMany(children ...h.HTMLComponent) (r *VXSelectManyBuilder) {
	r = &VXSelectManyBuilder{
		tag: h.Tag("vx-selectmany").Children(children...),
	}
	return
}

func (b *VXSelectManyBuilder) Attr(vs ...interface{}) *VXSelectManyBuilder {
	b.tag.Attr(vs...)
	return b
}

func (b *VXSelectManyBuilder) SetAttr(k string, v interface{}) *VXSelectManyBuilder {
	b.tag.SetAttr(k, v)
	return b
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
	b.tag.Attr("item-text", v)
	return b
}

func (b *VXSelectManyBuilder) ItemValue(v string) (r *VXSelectManyBuilder) {
	b.tag.Attr("item-value", v)
	return b
}

func (b *VXSelectManyBuilder) Label(v string) (r *VXSelectManyBuilder) {
	b.tag.Attr("label", v)
	return b
}

func (b *VXSelectManyBuilder) AddItemLabel(v string) (r *VXSelectManyBuilder) {
	b.tag.Attr("add-item-label", v)
	return b
}

func (b *VXSelectManyBuilder) MarshalHTML(ctx context.Context) (r []byte, err error) {
	if b.itemsSearcher != nil {
		b.tag.Attr(":search-items-func", fmt.Sprintf(`function(val){return %s.query("keyword", val).json()}`, b.itemsSearcher.String()))
	} else if b.searchItemsFunc != "" {
		b.tag.Attr(":search-items-func", fmt.Sprintf(`function(val){return $plaid().eventFunc("%s").query("keyword", val).json()}`, b.searchItemsFunc))
	} else {
		b.tag.Attr(":items", b.items)
	}

	b.tag.Attr(":selected-items", b.selectedItems)
	return b.tag.MarshalHTML(ctx)
}
