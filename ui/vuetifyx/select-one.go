package vuetifyx

import (
	"context"
	"fmt"

	"github.com/qor5/web/v3"
	h "github.com/theplant/htmlgo"
)

type VXSelectOneBuilder struct {
	tag             *h.HTMLTagBuilder
	selectedItem    interface{}
	items           interface{}
	searchItemsFunc string
	itemsSearcher   *web.VueEventTagBuilder
}

func VXSelectOne(children ...h.HTMLComponent) (r *VXSelectOneBuilder) {
	r = &VXSelectOneBuilder{
		tag: h.Tag("vx-selectone").Children(children...),
	}
	return
}

func (b *VXSelectOneBuilder) Attr(vs ...interface{}) *VXSelectOneBuilder {
	b.tag.Attr(vs...)
	return b
}

func (b *VXSelectOneBuilder) SetAttr(k string, v interface{}) *VXSelectOneBuilder {
	b.tag.SetAttr(k, v)
	return b
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
	b.tag.Attr("item-text", v)
	return b
}

func (b *VXSelectOneBuilder) ItemValue(v string) (r *VXSelectOneBuilder) {
	b.tag.Attr("item-value", v)
	return b
}

func (b *VXSelectOneBuilder) Label(v string) (r *VXSelectOneBuilder) {
	b.tag.Attr("label", v)
	return b
}

func (b *VXSelectOneBuilder) AddItemLabel(v string) (r *VXSelectOneBuilder) {
	b.tag.Attr("add-item-label", v)
	return b
}

func (b *VXSelectOneBuilder) MarshalHTML(ctx context.Context) (r []byte, err error) {
	if b.itemsSearcher != nil {
		b.tag.Attr(":search-items-func", fmt.Sprintf(`function(val){return %s.query("keyword", val).json()}`, b.itemsSearcher.String()))
	} else if b.searchItemsFunc != "" {
		b.tag.Attr(":search-items-func", fmt.Sprintf(`function(val){return $plaid().eventFunc("%s").query("keyword", val).json()}`, b.searchItemsFunc))
	} else {
		b.tag.Attr(":items", b.items)
	}

	var items []any
	if b.selectedItem != nil {
		items = append(items, b.selectedItem)
	}

	b.tag.Attr(":selected-items", items)
	return b.tag.MarshalHTML(ctx)
}
