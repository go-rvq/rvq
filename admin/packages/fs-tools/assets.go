package fs_tools

import (
	"context"
	"embed"
	"io/fs"
	"net/http"
	"sync"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

// The explorer of files (<vx-files>) and the terminal (<vx-terminal>): their
// scripts and their style, and xterm (vendor/, loaded by the terminal when
// first used).
//
//go:embed assets
var assetsFS embed.FS

var (
	assetsMu   sync.Mutex
	assetsUsed = map[*presets.Builder]bool{}
)

// UseAssets puts <vx-files> and <vx-terminal> on the pages of p (once): the
// scripts, the style, xterm at VendorPath.
func UseAssets(p *presets.Builder) {
	assetsMu.Lock()
	defer assetsMu.Unlock()
	if assetsUsed[p] {
		return
	}
	assetsUsed[p] = true
	// their words, in the user's language: the namespace vxFiles of the
	// page's vue-i18n
	ConfigureFilesMessages(p.I18n())
	p.VueI18nMessages(VueI18nNamespace, func(ctx context.Context) any { return GetFilesMessages(ctx) })
	read := func(name string) web.ComponentsPack {
		b, err := assetsFS.ReadFile("assets/" + name)
		if err != nil {
			panic(err)
		}
		return web.ComponentsPack(b)
	}
	p.ExtraAsset("/vx-files.js", "text/javascript", read("vx-files.js"))
	p.ExtraAsset("/vx-terminal.js", "text/javascript", read("vx-terminal.js"))
	p.ExtraAsset("/vx-files.css", "text/css", read("vx-files.css"))
	vendor, _ := fs.Sub(assetsFS, "assets/vendor")
	p.MuxSetup(func(prefix string, mux *http.ServeMux) {
		mux.Handle("GET "+prefix+vendorPath+"/", http.StripPrefix(prefix+vendorPath+"/", cacheFor(http.FileServer(http.FS(vendor)))))
	})
}

const vendorPath = "/fs-tools/vendor"

// VueI18nNamespace is the namespace of the vue-i18n where <vx-files> and
// <vx-terminal> read their words (FilesMessages).
const VueI18nNamespace = "vxFiles"

// VendorPath is where xterm is (the vendor of <vx-terminal>).
func VendorPath(p *presets.Builder) string { return p.GetURIPrefix() + vendorPath }

// cacheFor is h of what does not change but with the binary: kept a day.
func cacheFor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		h.ServeHTTP(w, r)
	})
}
