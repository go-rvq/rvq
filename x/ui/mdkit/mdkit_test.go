package mdkit

import (
	"strings"
	"testing"
)

// TestJSComponentsPack: the bundle is in the binary, its window.RvqMdKit.
func TestJSComponentsPack(t *testing.T) {
	if js := string(JSComponentsPack()); !strings.Contains(js, "RvqMdKit") || len(js) < 100000 {
		t.Errorf("the pack of mdkit: %d bytes, RvqMdKit %v", len(js), strings.Contains(js, "RvqMdKit"))
	}
}

// TestPagePack: the script of a page that shows a text is in the binary,
// small, reading the page's metas; the meta of the assets says KaTeX and
// Mermaid, their CDN and their copies.
func TestPagePack(t *testing.T) {
	if js := string(PagePack()); !strings.Contains(js, "mdkit-link-cards") || !strings.Contains(js, "mdkit-vendor") || len(js) > 20000 {
		t.Errorf("the pack of the page: %d bytes", len(js))
	}
	m := VendorMeta()
	for _, want := range []string{`name="mdkit-vendor"`, "katex-css", "cdn.jsdelivr.net/npm/mermaid", "/_rvq/vendor/katex/0.16.47/katex.min.js", "sha384-", "/_rvq/vendor/prism/1.30.0/prism.min.js?v=", "/_rvq/vendor/codemirror/6/codemirror.min.js?v="} {
		if !strings.Contains(m, want) {
			t.Errorf("the meta has no %s: %s", want, m)
		}
	}
}
