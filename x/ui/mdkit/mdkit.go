// Package mdkit is the pack of rvq's mdkit (js/mdkit) for a page: the editors
// of a text — Markdown (EasyMDE), rich (tiptap) — and their builders (tables,
// footnotes), Vue components with the page's Vuetify, window.RvqMdKit.
package mdkit

import (
	"io/fs"

	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/go-rvq/rvq/web"
)

// JSComponentsPack is its script (mdkit.js): after Vue and Vuetify.
func JSComponentsPack() web.ComponentsPack {
	FS, _ := fs.Sub(rvqjs.MdKit, "mdkit/dist")
	return web.ComponentsPackFromFile(FS, "mdkit.js")
}

// LinkCardsPack is the script of the link cards alone (link-cards.js, a few
// KB): a page that shows a text with cards and no editor. It draws them,
// asking where the page's <meta name="mdkit-link-cards"> says.
func LinkCardsPack() web.ComponentsPack {
	FS, _ := fs.Sub(rvqjs.MdKit, "mdkit/dist")
	return web.ComponentsPackFromFile(FS, "link-cards.js")
}
