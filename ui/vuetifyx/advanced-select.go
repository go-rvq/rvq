package vuetifyx

import (
	"context"
	"fmt"

	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type VXAdvancedSelectBuilder struct {
	v.VTagBuilder[*VXAdvancedSelectBuilder]
	items           interface{}
	searchItemsFunc string
	itemsSearcher   *web.VueEventTagBuilder
	many            bool
}

func VXAdvancedSelect(children ...h.HTMLComponent) *VXAdvancedSelectBuilder {
	return v.VTag(&VXAdvancedSelectBuilder{}, "vx-advanced-select", children...)
}

func (b *VXAdvancedSelectBuilder) GetVXAdvancedSelect() *VXAdvancedSelectBuilder {
	return b
}

func (b *VXAdvancedSelectBuilder) Items(v interface{}) *VXAdvancedSelectBuilder {
	b.items = v
	return b.Dot()
}

func (b *VXAdvancedSelectBuilder) SearchItemsFunc(v string) *VXAdvancedSelectBuilder {
	b.searchItemsFunc = v
	return b.Dot()
}

func (b *VXAdvancedSelectBuilder) ItemsSearcher(eb *web.VueEventTagBuilder) *VXAdvancedSelectBuilder {
	b.itemsSearcher = eb
	return b.Dot()
}

func (b *VXAdvancedSelectBuilder) Many(v bool) *VXAdvancedSelectBuilder {
	b.many = v
	return b
}

func (b *VXAdvancedSelectBuilder) ItemText(v string) *VXAdvancedSelectBuilder {
	return b.Attr("item-text", v)
}

func (b *VXAdvancedSelectBuilder) ItemValue(v string) *VXAdvancedSelectBuilder {
	return b.Attr("item-value", v)
}

func (b *VXAdvancedSelectBuilder) Label(v string) *VXAdvancedSelectBuilder {
	return b.Attr("label", v)
}

func (b *VXAdvancedSelectBuilder) AddItemLabel(v string) *VXAdvancedSelectBuilder {
	return b.Attr("add-item-label", v)
}

func (b *VXAdvancedSelectBuilder) ItemTextExpr(v string) *VXAdvancedSelectBuilder {
	return b.Attr(":item-text", v)
}

func (b *VXAdvancedSelectBuilder) MarshalHTML(ctx context.Context) (r []byte, err error) {
	if b.itemsSearcher != nil {
		b.Attr(":search-items-func", fmt.Sprintf(`function(val){return %s.query("keyword", val).json().then(v => v.Records)}`, b.itemsSearcher.String()))
	} else if b.searchItemsFunc != "" {
		b.Attr(":search-items-func", fmt.Sprintf(`function(val){return $plaid().eventFunc("%s").query("keyword", val).json().then(v => v.Records)}`, b.searchItemsFunc))
	}

	if b.items != nil {
		b.Attr(":items", b.items)
	}

	if b.many {
		b.Attr("many", true)
	}

	return b.GetHTMLTagBuilder().MarshalHTML(ctx)
}
