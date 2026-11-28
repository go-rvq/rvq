package vuetifyx

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

type VXAdvancedSelectBuilder struct {
	v.VTagBuilder[*VXAdvancedSelectBuilder]
	items           interface{}
	searchItemsFunc string
	itemsSearcher   *web.VueEventTagBuilder
	many            bool
	chips           bool
	dialog          bool
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

// Dialog makes the selector open in a dialog — the search, the list of a page
// of the records, the pages — in place of the autocomplete. The items, the
// value and the cache of the chosen ones stay the same; with an ItemsSearcher,
// the search asks for the page too.
func (b *VXAdvancedSelectBuilder) Dialog(v bool) *VXAdvancedSelectBuilder {
	b.dialog = v
	return b
}

// IsDialog says the selector opens in a dialog.
func (b *VXAdvancedSelectBuilder) IsDialog() bool { return b.dialog }

// DialogTexts are the texts of the dialog: its title (none: the label), the
// chip that opens it (to choose, to change the one chosen), the search, the
// empty list and the button that closes it.
func (b *VXAdvancedSelectBuilder) DialogTexts(title, sel, change, search, noItems, close string) *VXAdvancedSelectBuilder {
	return b.Attr("dialog-title", title, "select-text", sel, "change-text", change,
		"search-text", search, "no-items-text", noItems, "close-text", close)
}

// ItemHref is the address of the page of a record chosen, "{id}" its value:
// its open, an icon on its right in place of its change (its text, clicked,
// changes it) — clicked, the event "open-item" (the record: its detail in a
// dialog, say); its middle button, the address in a new tab.
func (b *VXAdvancedSelectBuilder) ItemHref(v string) *VXAdvancedSelectBuilder {
	return b.Attr("item-href", v)
}

// ItemTexts are the tooltips of the icons beside a record chosen: its open
// (ItemHref) and its removal (its change: DialogTexts).
func (b *VXAdvancedSelectBuilder) ItemTexts(open, remove string) *VXAdvancedSelectBuilder {
	return b.Attr("open-text", open, "remove-text", remove)
}

// Readonly shows the ones chosen alone: no change, no removal, no choice.
func (b *VXAdvancedSelectBuilder) Readonly(v bool) *VXAdvancedSelectBuilder {
	return b.Attr(":readonly", v)
}

func (b *VXAdvancedSelectBuilder) Chips(v bool) *VXAdvancedSelectBuilder {
	b.chips = v
	return b
}

func (b *VXAdvancedSelectBuilder) ItemText(v string) *VXAdvancedSelectBuilder {
	return b.Attr("item-text", v)
}

func (b *VXAdvancedSelectBuilder) Hint(v string) *VXAdvancedSelectBuilder {
	return b.Attr("hint", v)
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

func (b *VXAdvancedSelectBuilder) Write(ctx *h.Context) (err error) {
	if b.itemsSearcher != nil {
		s := b.itemsSearcher.String()
		b.Attr(":search-items-func", fmt.Sprintf(`function(val){return %s.query("keyword", val).json().then(v => v.Records)}`, s))
		if b.dialog {
			// a page of the search: the records, the total, the size of a page
			b.Attr(":search-page-func", fmt.Sprintf(`function(val, page){return %s.query("keyword", val).query("page", String(page)).json()}`, s))
		}
	} else if b.searchItemsFunc != "" {
		b.Attr(":search-items-func", fmt.Sprintf(`function(val){return $plaid().eventFunc("%s").query("keyword", val).json().then(v => v.Records)}`, b.searchItemsFunc))
	}

	if b.items != nil {
		b.Attr(":items", b.items)
	}

	if b.many {
		b.Attr("many", true)
	}

	if b.chips {
		b.Attr("chips", true)
	}

	if b.dialog {
		b.Attr("dialog", true)
	}

	return b.GetHTMLTagBuilder().Write(ctx)
}
