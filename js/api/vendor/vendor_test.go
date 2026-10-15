package vendor

import (
	"crypto/sha512"
	"encoding/base64"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/js"
)

// Each asset's copy is there, served, and the file its integrity says (the
// CDN's file of that version is checked against it, and the copy taken in
// its place on a mismatch: they must be the same file).
func TestCopies(t *testing.T) {
	h := Handler()
	for _, a := range Assets {
		b, err := fs.ReadFile(js.Vendor, "vendor/"+a.Local)
		if err != nil {
			t.Errorf("%s: no copy: %v", a.Name, err)
			continue
		}
		if a.Integrity != "" {
			sum := sha512.Sum384(b)
			if got := "sha384-" + base64.StdEncoding.EncodeToString(sum[:]); got != a.Integrity {
				t.Errorf("%s: the copy is not the file of its integrity: %s", a.Name, got)
			}
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, a.URL(), nil))
		if w.Code != 200 || w.Body.Len() != len(b) {
			t.Errorf("%s: served %d (%d bytes)", a.Name, w.Code, w.Body.Len())
		}
		if !strings.Contains(w.Header().Get("Cache-Control"), "immutable") {
			t.Errorf("%s: not cached", a.Name)
		}
	}
	// what a stylesheet of a copy asks (its fonts): there too
	for _, css := range []string{FontAwesome.Local, Inter.Local} {
		b, _ := fs.ReadFile(js.Vendor, "vendor/"+css)
		for _, part := range strings.Split(string(b), "url(")[1:] {
			u := strings.Trim(part[:strings.IndexByte(part, ')')], `"'`)
			if strings.HasPrefix(u, "data:") {
				continue
			}
			p := Prefix + css[:strings.LastIndexByte(css, '/')+1] + u
			w := httptest.NewRecorder()
			h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, p, nil))
			if w.Code != 200 {
				t.Errorf("%s: %s not served (%d)", css, p, w.Code)
			}
		}
	}
	// no listing, no readme
	for _, p := range []string{Prefix, Prefix + "vue/", Prefix + "README.md"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, p, nil))
		if w.Code != 404 {
			t.Errorf("%s: %d", p, w.Code)
		}
	}
}

// A tag asks the CDN first, the copy on error: a script by the bootstrap
// (its global waited for), a stylesheet by its href.
func TestTags(t *testing.T) {
	s, err := Tags("vue", "fontawesome", "tailwind", "inter")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`<script defer src="https://unpkg.com/vue@3.5.44/dist/vue.global.prod.js" integrity="sha384-`,
		`onerror="window.rvqVendor&amp;&amp;rvqVendor.fallback(this,&#34;/_rvq/vendor/vue/3.5.44/vue.global.prod.js&#34;,&#34;Vue&#34;)"`,
		`onload="window.rvqVendor&amp;&amp;rvqVendor.done(&#34;Vue&#34;)"`,
		`<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css" integrity="sha384-`,
		`this.href=&#34;/_rvq/vendor/fontawesome/6.5.1/css/all.min.css&#34;`,
		`<script src="https://cdn.tailwindcss.com/3.4.17"`,
		`href="https://fonts.googleapis.com/css2?family=Inter`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("no %s in %s", want, s)
		}
	}
	if strings.Contains(Inter.Tag(), `integrity="`) || strings.Contains(Tailwind.Tag(), `integrity="`) {
		t.Error("an integrity of a file the CDN makes (Google's css, Tailwind's)")
	}
	if _, err := Tags("nothing"); err == nil {
		t.Error("an unknown asset")
	}
	if !strings.HasPrefix(Bootstrap, "<script>") || !strings.Contains(Bootstrap, "when:function") {
		t.Error("the bootstrap")
	}
}
