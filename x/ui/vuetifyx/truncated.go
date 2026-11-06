package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

// TruncatedText is a text of a cell of a table (a description, a name) of at
// most maxWidth (a CSS width; "", 280px), cut with "…" on one line, the whole
// text its title. A click on it, when cut, shows it whole — broken across
// lines if need be — and another cuts it again; that click does not reach
// the row (its detail). Not cut, its click is the row's.
func TruncatedText(text, maxWidth string) h.HTMLComponent {
	if maxWidth == "" {
		maxWidth = "280px"
	}
	return web.Scope(
		h.Div(h.Text(text)).
			Attr("title", text).
			Attr("data-truncated-text", true).
			Style("max-width: "+maxWidth).
			Attr(":style", "locals.open ? {whiteSpace: 'normal', overflowWrap: 'anywhere'} : "+
				"{whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}").
			Attr("@click", "if (locals.open || $event.currentTarget.scrollWidth > $event.currentTarget.clientWidth) "+
				"{ $event.stopPropagation(); locals.open = !locals.open }"),
	).LocalsInit("{open: false}").Slot("{ locals }")
}

// TruncatedCell is a cell of a table that takes the room left in its row
// (the others as wide as their contents), its text cut at the width of the
// column (TruncatedText) — the description of a listing, say.
func TruncatedCell(text string) *h.HTMLTagBuilder {
	// max-width 0: the cell does not grow with its text; width 100%: the
	// room left is its
	return h.Td(TruncatedText(text, "100%")).Style("width: 100%; max-width: 0").Attr("data-truncated-cell", true)
}
