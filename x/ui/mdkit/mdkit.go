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
