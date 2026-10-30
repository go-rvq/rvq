// Package mdkit is the pack of rvq's mdkit (js/mdkit) for a page: the editors
// of a text — Markdown (EasyMDE), rich (tiptap) — and their builders (tables,
// footnotes), Vue components with the page's Vuetify, window.RvqMdKit.
package mdkit

import (
	"encoding/json"
	"html"
	"io/fs"

	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/go-rvq/rvq/js/api/vendor"
	"github.com/go-rvq/rvq/web"
)

// JSComponentsPack is its script (mdkit.js): after Vue and Vuetify.
func JSComponentsPack() web.ComponentsPack {
	FS, _ := fs.Sub(rvqjs.MdKit, "mdkit/dist")
	return web.ComponentsPackFromFile(FS, "mdkit.js")
}

// PagePack is the script of a page that shows a text and has no editor
// (page.js, a few KB): its link cards (asked where the page's <meta
// name="mdkit-link-cards"> says), its formulas and diagrams (KaTeX, Mermaid,
// loaded from where <meta name="mdkit-vendor"> says: VendorMeta).
func PagePack() web.ComponentsPack {
	FS, _ := fs.Sub(rvqjs.MdKit, "mdkit/dist")
	return web.ComponentsPackFromFile(FS, "page.js")
}

// VendorMeta is the meta of the assets the page's formulas, diagrams and code
// load (mdkit-vendor: KaTeX, its styles, Mermaid — their CDN, their copy —,
// Prism and CodeMirror — the app's).
func VendorMeta() string {
	m := map[string]map[string]string{}
	for _, a := range []*vendor.Asset{vendor.KaTeX, vendor.KaTeXStyle, vendor.Mermaid, vendor.Prism, vendor.CodeMirror} {
		m[a.Name] = map[string]string{"src": a.Src(), "local": a.URL(), "integrity": a.Integrity}
	}
	b, _ := json.Marshal(m)
	return `<meta name="mdkit-vendor" content="` + html.EscapeString(string(b)) + `">`
}
