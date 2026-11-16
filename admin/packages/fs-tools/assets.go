package fs_tools

import (
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

// VendorPath is where xterm is (the vendor of <vx-terminal>).
func VendorPath(p *presets.Builder) string { return p.GetURIPrefix() + vendorPath }

// cacheFor is h of what does not change but with the binary: kept a day.
func cacheFor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		h.ServeHTTP(w, r)
	})
}
