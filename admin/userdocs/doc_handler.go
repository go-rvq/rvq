package userdocs

import (
	stdhtml "html"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/go-rvq/rvq/web"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// The page of a guide (DocHandler): a document open to anyone — the guides
// of the editors of a text, which the ? of an editor shows —, rendered from
// the same files and edits as the documentation page's, so a change of a
// document is a change everywhere. Its address says what it is: the node
// and the language (?node=, ?locale=); a page of its own (a site's dialog
// opens it in a new tab), or — ?portal=1 — the content of a rvq-portal
// (js/portal: the editors' dialog loads it there); for the admin
// (?view=admin) its code the admin's code viewer. The languages it can be
// read in, and its links to the other guides open, are links to this page
// (of its own: in the dialog they load in its portal).

// Public makes the documents of nodes open to anyone (DocHandler).
func (b *Builder) Public(nodes ...string) *Builder {
	for _, n := range nodes {
		if !slices.Contains(b.public, n) {
			b.public = append(b.public, n)
		}
	}
	return b
}

// PageHead sets what the page of a guide has in its head beside its look:
// what draws its code, formulas and diagrams (rvq mdkit's page script and
// its vendor meta, as the application serves them).
func (b *Builder) PageHead(html string) *Builder {
	b.pageHead = html
	return b
}

// DocPagePath is where the page of a guide is served (DocHandler): out of
// the admin's address — whose every page asks for a login —, a visitor of
// the site sees it too.
const DocPagePath = "/_rvq/user-docs/doc"

// DocPath is where the page of a guide is served (DocHandler, DocPagePath).
func (b *Builder) DocPath() string { return DocPagePath }

// DocDataHref is the address of the page of the guide of node for the
// admin (its code the admin's code viewer, as its documentation page's).
func (b *Builder) DocDataHref(node string) string {
	return docQuery{node: node, admin: true}.href(b.DocPath())
}

// DocDataHref is the address of the page of the guide of node served at
// path, for a site.
func DocDataHref(path, node string) string {
	return docQuery{node: node}.href(path)
}

// docQuery is what the address of a page of a guide says.
type docQuery struct {
	node, locale  string
	admin, portal bool
}

func (q docQuery) href(path string) string {
	v := url.Values{"node": {q.node}}
	if q.locale != "" {
		v.Set("locale", q.locale)
	}
	if q.admin {
		v.Set("view", "admin")
	}
	if q.portal {
		v.Set("portal", "1")
	}
	return path + "?" + v.Encode()
}

// matchLocale is the language of the documents of the one asked for —
// the same, else of its base —, "" none.
func matchLocale(locales []string, asked string) string {
	if asked == "" {
		return ""
	}
	for _, l := range locales {
		if strings.EqualFold(l, asked) {
			return l
		}
	}
	for _, l := range locales {
		if sameLanguage(l, asked) {
			return l
		}
	}
	return ""
}

// localeName is the name of a language in itself ("Português (Brasil)").
func localeName(code string) string {
	t, err := language.Parse(code)
	if err != nil {
		return code
	}
	if n := display.Self.Name(t); n != "" {
		return strings.ToUpper(n[:1]) + n[1:]
	}
	return code
}

// DocHandler serves the page of the guide of ?node= — one made open
// (Public) — in the language of ?locale= (one the documents are not in: the
// first of them, and a word that says so above the languages to choose):
// a page of its own, or the content of a portal (?portal=1); for the admin
// (?view=admin) or a site.
func (b *Builder) DocHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query()
		q := docQuery{node: v.Get("node"), admin: v.Get("view") == "admin", portal: v.Get("portal") != ""}
		locales := b.locales()
		if !slices.Contains(b.public, q.node) || len(locales) == 0 {
			http.NotFound(w, r)
			return
		}
		q.locale = matchLocale(locales, v.Get("locale"))
		match := q.locale != ""
		if !match {
			q.locale = locales[0]
		}
		file, src, content, err := b.find(q.locale, q.node)
		if err != nil || src == nil {
			http.NotFound(w, r)
			return
		}
		// in its language: its words
		r = r.WithContext(b.p.I18n().LanguageContext(r.Context(), q.locale))
		ctx := &web.EventContext{R: r}
		path := b.DocPath()
		rd := &renderer{b: b, ctx: ctx, node: q.node, locale: q.locale, plain: !q.admin,
			// a link to another guide open: its page, as this one is
			href: func(node string) string {
				if !slices.Contains(b.public, node) {
					return ""
				}
				o := q
				o.node, o.portal = node, false
				return o.href(path)
			}}
		out, err := rd.Render(src, file, content)
		if err != nil {
			http.Error(w, "the document does not render", http.StatusInternalServerError)
			return
		}
		body := b.docBody(ctx, q, out, locales, match)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		if q.portal {
			_, _ = w.Write([]byte(body))
			return
		}
		_, _ = w.Write([]byte(`<!DOCTYPE html><html lang="` + stdhtml.EscapeString(q.locale) + `"><head><meta charset="utf-8">` +
			`<meta name="viewport" content="width=device-width, initial-scale=1">` +
			`<title>` + stdhtml.EscapeString(b.docTitle(q.locale, q.node)) + `</title>` +
			// the colors the look of a document takes from Vuetify's theme
			`<style>:root{--v-border-color:0,0,0;--v-border-opacity:.12;--v-theme-on-surface:0,0,0;--v-medium-emphasis-opacity:.6;--v-theme-primary:25,118,210}` +
			`body{margin:0 auto;max-width:56rem;padding:1.5rem 1rem;font-family:system-ui,sans-serif;line-height:1.5;color:#1f2328;background:#fff}</style>` +
			b.pageHead + `</head><body>` + body + `</body></html>`))
	})
}

// docStyleLocales is the look of the languages of a page of a guide.
const docStyleLocales = `.user-doc-locales{display:flex;flex-wrap:wrap;align-items:center;gap:.4rem;margin:0 0 1rem;font-size:.875rem}` +
	`.user-doc-locales a{padding:.1rem .6rem;border:1px solid rgba(127,127,127,.4);border-radius:1rem;text-decoration:none;color:inherit}` +
	`.user-doc-locales a[aria-current]{background:rgba(25,118,210,.12);border-color:rgb(25,118,210);font-weight:600}` +
	`.user-doc-notice{margin:0 0 .5rem;padding:.5rem .75rem;border-radius:4px;background:rgba(2,136,209,.1);font-size:.875rem}`

// docBody is the content of the page of a guide: its look (in a portal, a
// <component is="style">, which Vue keeps; else a style), the
// languages it can be read in (links to it), and the document — for a
// site, nothing of it read by Vue (v-pre).
func (b *Builder) docBody(ctx *web.EventContext, q docQuery, html string, locales []string, match bool) string {
	msgs := GetMessages(ctx.Context())
	esc := stdhtml.EscapeString
	var s strings.Builder
	s.WriteString(`<div class="user-docs" data-user-doc="` + esc(q.node) + `" data-user-doc-locale="` + esc(q.locale) + `">`)
	if q.portal {
		s.WriteString(`<component :is="'style'">` + docStyle + docStyleLocales + `</component>`)
	} else {
		s.WriteString(`<style>` + docStyle + docStyleLocales + `</style>`)
	}
	if !match {
		s.WriteString(`<p class="user-doc-notice" data-user-doc-not-yours>` + esc(msgs.DocNotInYours) + `</p>`)
	}
	if len(locales) > 1 {
		s.WriteString(`<nav class="user-doc-locales" aria-label="` + esc(msgs.DocLanguage) + `"><span>` + esc(msgs.DocLanguage) + `:</span>`)
		for _, l := range locales {
			o := q
			o.locale, o.portal = l, false
			cur := ""
			if l == q.locale {
				cur = ` aria-current="true"`
			}
			s.WriteString(`<a href="` + esc(o.href(b.DocPath())) + `" data-locale="` + esc(l) + `"` + cur + `>` + esc(localeName(l)) + `</a>`)
		}
		s.WriteString(`</nav>`)
	}
	pre := " v-pre"
	if q.admin {
		pre = ""
	}
	s.WriteString(`<div class="user-doc" lang="` + esc(q.locale) + `"` + pre + `>` + html + `</div></div>`)
	return s.String()
}
