package markup

import (
	"bytes"
	stdhtml "html"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

// The buttons and the tabs of a text, as the editors make them (rvq mdkit's
// widgets-model.js) — in Markdown a block of HTML, a tab's content Markdown
// between blank lines —:
//
//	<div data-type="buttons" data-align="center"><a data-type="button"
//	href="…" data-href="…" title="description" data-color="primary"
//	data-icon-left="mdi-…" data-icon-right="mdi-…">title</a>…</div>
//	<div data-type="tabs" data-align="start"><div data-type="tab"
//	data-title="…" data-color="…" data-icon-left="…">…</div>…</div>
//
// a button's title empty is its target's (the page shows it: the link
// cards'); the tabs drawn by the page (a section each, printed).
const (
	Buttons = "buttons"
	Button  = "button"
	Tabs    = "tabs"
	Tab     = "tab"
)

var (
	widgetColorRe = regexp.MustCompile(`^(primary|secondary|success|info|warning|error)$`)
	widgetIconRe  = regexp.MustCompile(`^mdi-[a-z0-9-]+$`)
	widgetAlignRe = regexp.MustCompile(`^(start|center|end|stretch)$`)
)

// allowWidgets makes p keep the buttons and the tabs: their kinds, their
// place, the look of each (a color of the theme, icons of Material Design),
// a tab's title, a button's target and its new tab.
func allowWidgets(p *bluemonday.Policy) {
	p.AllowAttrs("data-type").Matching(regexp.MustCompile(`^(` + Buttons + `|` + Tabs + `|` + Tab + `)$`)).OnElements("div")
	p.AllowAttrs("data-type").Matching(regexp.MustCompile(`^` + Button + `$`)).OnElements("a")
	p.AllowAttrs("data-align").Matching(widgetAlignRe).OnElements("div")
	p.AllowAttrs("data-color").Matching(widgetColorRe).OnElements("div", "a")
	p.AllowAttrs("data-icon-left", "data-icon-right").Matching(widgetIconRe).OnElements("div", "a")
	p.AllowAttrs("data-title").Matching(regexp.MustCompile(`^[^\x00]*$`)).OnElements("div")
	p.AllowAttrs("data-href").Matching(cardTargetRe).OnElements("a")
	p.AllowAttrs("target").Matching(regexp.MustCompile(`^_blank$`)).OnElements("a")
}

// widgetMarkdown writes n — buttons, tabs — as Markdown (true), or not one.
func (w *mdWriter) widgetMarkdown(n *html.Node, prefix string) bool {
	switch attr(n, "data-type") {
	case Buttons:
		// a block of HTML, one line
		var b bytes.Buffer
		_ = html.Render(&b, n)
		w.b.WriteString(prefix + strings.ReplaceAll(b.String(), "\n", " ") + "\n\n")
		return true
	case Tabs:
		w.b.WriteString(prefix + `<div data-type="tabs" data-align="` + stdhtml.EscapeString(attrOr(n, "data-align", "start")) + `">` + "\n")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode || attr(c, "data-type") != Tab {
				continue
			}
			w.b.WriteString(prefix + `<div data-type="tab" data-title="` + stdhtml.EscapeString(attr(c, "data-title")) + `"`)
			for _, a := range []string{"data-color", "data-icon-left", "data-icon-right"} {
				if v := attr(c, a); v != "" {
					w.b.WriteString(" " + a + `="` + stdhtml.EscapeString(v) + `"`)
				}
			}
			w.b.WriteString(">\n\n")
			for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
				w.block(cc, prefix)
			}
			w.b.WriteString(prefix + "</div>\n")
		}
		w.b.WriteString(prefix + "</div>\n\n")
		return true
	}
	return false
}

// attrOr is the attribute name of n, or def when it has none.
func attrOr(n *html.Node, name, def string) string {
	if v := attr(n, name); v != "" {
		return v
	}
	return def
}
