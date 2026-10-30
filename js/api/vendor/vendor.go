// Package vendor is what the pages of a site take from CDNs — Vue,
// Tailwind's script, Font Awesome, the Inter font; KaTeX and Mermaid, for
// the formulas and diagrams of a text, and Prism, its code — with a fallback to the
// app: each asset's tag asks its CDN first and, when the CDN does not
// answer (an error: down, blocked), the copy the app serves (js.Vendor, at
// Prefix), of the same version.
//
// A script loaded by the fallback comes later than the page's deferred
// scripts: what uses it waits for it — rvqVendor.when("Vue", fn), of the
// bootstrap (Bootstrap) — instead of finding it at DOMContentLoaded.
package vendor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-rvq/rvq/js"
)

// Prefix is where the copies are served.
const Prefix = "/_rvq/vendor/"

// Kind is what an asset is: a script or a stylesheet.
type Kind int

const (
	Script Kind = iota
	Style
)

// Asset is a file a page takes from a CDN, and its copy.
type Asset struct {
	// Name names it (Tags).
	Name string
	Kind Kind
	// CDN is its address there, of a fixed version.
	CDN string
	// Local is its copy, under Prefix.
	Local string
	// Integrity is the hash of the file (SRI), when the CDN's is the same
	// file always: checked, and the copy taken when it does not match.
	Integrity string
	// Global is what a script defines (window.Vue): rvqVendor.when waits for
	// it.
	Global string
	// Defer loads a script deferred.
	Defer bool
	// Version is the version of a copy whose path does not say it (?v=).
	Version string
}

// The assets.
var (
	Vue = &Asset{Name: "vue", Kind: Script, Defer: true, Global: "Vue",
		CDN:       "https://unpkg.com/vue@3.5.44/dist/vue.global.prod.js",
		Local:     "vue/3.5.44/vue.global.prod.js",
		Integrity: "sha384-tORx0bXmMnTC28EYQe1coaGZKSQupXQEYDBNex8xkzJ9iqUoNy9Y+VFBhVxi4kUr"}
	Tailwind = &Asset{Name: "tailwind", Kind: Script, Global: "tailwind",
		CDN:   "https://cdn.tailwindcss.com/3.4.17",
		Local: "tailwind/3.4.17/tailwind.js"}
	FontAwesome = &Asset{Name: "fontawesome", Kind: Style,
		CDN:       "https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css",
		Local:     "fontawesome/6.5.1/css/all.min.css",
		Integrity: "sha384-t1nt8BQoYMLFN5p42tRAtuAAFQaCQODekUVeKKZrEnEyp4H2R0RHFz0KWpmj7i8g"}
	// Inter: Google's stylesheet (made for each browser: no integrity), its
	// fonts from fonts.gstatic.com; the copy, woff2 of the same weights.
	Inter = &Asset{Name: "inter", Kind: Style,
		CDN:   "https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&display=swap",
		Local: "inter/inter.css"}
	// KaTeX: the formulas of a text (rvq mdkit's blocks: $…$, $$…$$); its
	// fonts beside its css (fonts/, woff2 in the copy)
	KaTeX = &Asset{Name: "katex", Kind: Script, Defer: true, Global: "katex",
		CDN:       "https://cdn.jsdelivr.net/npm/katex@0.16.47/dist/katex.min.js",
		Local:     "katex/0.16.47/katex.min.js",
		Integrity: "sha384-CwjPRVHTvLiMBFjEoij+QZViMV5rhTOIp7CJzl24JEqpRDA1sJFHVXXLURktbYYp"}
	KaTeXStyle = &Asset{Name: "katex-css", Kind: Style,
		CDN:       "https://cdn.jsdelivr.net/npm/katex@0.16.47/dist/katex.min.css",
		Local:     "katex/0.16.47/katex.min.css",
		Integrity: "sha384-nH0MfJ44wi1dd7w6jinlyBgljjS8EJAh2JBoRad8a3VDw2K69vfaaqm4WnR+gXtA"}
	// Mermaid: the diagrams of a text (```mermaid)
	Mermaid = &Asset{Name: "mermaid", Kind: Script, Defer: true, Global: "mermaid",
		CDN:       "https://cdn.jsdelivr.net/npm/mermaid@11.17.2/dist/mermaid.min.js",
		Local:     "mermaid/11.17.2/mermaid.min.js",
		Integrity: "sha384-EOXBFmc3gx5mb+vn0vPvvGqACToJD24hhacX5Yx+8NUUQrHIle/Qi5Bg9o3zKwW2"}
	// Material Design Icons: the icons of the buttons and the tabs of a text
	// (rvq mdkit's widgets.js); the copy has the woff2 and woff fonts
	MDIStyle = &Asset{Name: "mdi-css", Kind: Style,
		CDN:       "https://cdn.jsdelivr.net/npm/@mdi/font@7.4.47/css/materialdesignicons.min.css",
		Local:     "mdi/7.4.47/css/materialdesignicons.min.css",
		Integrity: "sha384-HphS8cQyN+eYiJ5PMbzShG6qZdRtvHPVLPkYb8JwMkmNgaIxrFVDhQe3jIbq3EZ2"}
	// Prism: the code of a text highlighted (rvq mdkit's code.js), with the
	// languages and plugins of the admin's code viewer — a bundle of its own
	// (js/mdkit build.sh), no CDN: the app's only, its address the version of
	// its content (?v=, cached a year)
	Prism = &Asset{Name: "prism", Kind: Script, Defer: true, Global: "Prism",
		Local: "prism/1.30.0/prism.min.js"}
	// CodeMirror: the code blocks of the rich text editor edited (rvq mdkit's
	// code-block-tiptap.js), with its languages — a bundle of its own (js/mdkit
	// build.sh), no CDN, its address the version of its content
	CodeMirror = &Asset{Name: "codemirror", Kind: Script, Defer: true, Global: "RvqCodeMirror",
		Local: "codemirror/6/codemirror.min.js"}
)

// Assets are the assets by name.
var Assets = map[string]*Asset{}

// Register adds a to the assets (Tags knows it by its name).
func Register(a ...*Asset) {
	for _, x := range a {
		Assets[x.Name] = x
	}
}

func init() {
	Prism.Version = version(Prism.Local)
	CodeMirror.Version = version(CodeMirror.Local)
	Register(Vue, Tailwind, FontAwesome, Inter, KaTeX, KaTeXStyle, Mermaid, MDIStyle, Prism, CodeMirror)
}

// version is the version of the copy at local: its hash, short.
func version(local string) string {
	b, err := fs.ReadFile(js.Vendor, "vendor/"+local)
	if err != nil {
		panic(err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:6])
}

// Src is where a is loaded from first: its CDN, else (none) its copy.
func (a *Asset) Src() string {
	if a.CDN == "" {
		return a.URL()
	}
	return a.CDN
}

// Bootstrap is the script before the assets' tags (once a page): the
// fallback of a script — its copy loaded, and who waits for it called —, and
// rvqVendor.when(global, fn): fn now if the global is there, else once its
// script loads (from the CDN or the copy).
const Bootstrap = `<script>(function(w){if(w.rvqVendor)return;var q={};var v=w.rvqVendor={` +
	`when:function(g,f){if(w[g])return f();(q[g]=q[g]||[]).push(f)},` +
	`done:function(g){if(!w[g])return;(q[g]||[]).splice(0).forEach(function(f){try{f()}catch(e){console.error(e)}})},` +
	`fallback:function(el,src,g){el.onerror=null;var s=document.createElement('script');s.src=src;` +
	`s.onload=function(){v.done(g)};document.head.appendChild(s)}}})(window)</script>`

// URL is the address of the copy of a.
func (a *Asset) URL() string {
	if a.Version != "" {
		return Prefix + a.Local + "?v=" + a.Version
	}
	return Prefix + a.Local
}

// Tag is the tag of a: its CDN first; on error, its copy (Bootstrap before
// it, in the page).
func (a *Asset) Tag() string {
	attr := func(name, v string) string { return " " + name + `="` + html.EscapeString(v) + `"` }
	var b strings.Builder
	switch a.Kind {
	case Style:
		b.WriteString(`<link rel="stylesheet"`)
		b.WriteString(attr("href", a.Src()))
		if a.Integrity != "" {
			b.WriteString(attr("integrity", a.Integrity))
			b.WriteString(` crossorigin="anonymous"`)
		}
		// the copy: the same link, its address the app's (no integrity: the
		// copy is ours)
		b.WriteString(attr("onerror", "this.onerror=null;this.removeAttribute('integrity');this.href="+strconv.Quote(a.URL())))
		b.WriteString(">")
	default:
		b.WriteString("<script")
		if a.Defer {
			b.WriteString(" defer")
		}
		b.WriteString(attr("src", a.Src()))
		if a.Integrity != "" {
			b.WriteString(attr("integrity", a.Integrity))
			b.WriteString(` crossorigin="anonymous"`)
		}
		g := strconv.Quote(a.Global)
		b.WriteString(attr("onload", "window.rvqVendor&&rvqVendor.done("+g+")"))
		b.WriteString(attr("onerror", "window.rvqVendor&&rvqVendor.fallback(this,"+strconv.Quote(a.URL())+","+g+")"))
		b.WriteString("></script>")
	}
	return b.String()
}

// Tags are the tags of the assets named, in order — an unknown name an
// error.
func Tags(names ...string) (string, error) {
	var b strings.Builder
	for _, n := range names {
		a := Assets[n]
		if a == nil {
			return "", fmt.Errorf("vendor: no asset %q", n)
		}
		b.WriteString(a.Tag())
	}
	return b.String(), nil
}

// Handler serves the copies (js.Vendor) at Prefix, cached a year: a new
// version is another path.
func Handler() http.Handler {
	sub, err := fs.Sub(js.Vendor, "vendor")
	if err != nil {
		panic(err)
	}
	files := http.StripPrefix(strings.TrimSuffix(Prefix, "/"), http.FileServerFS(sub))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") || strings.HasSuffix(r.URL.Path, ".md") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		// the fonts asked by a stylesheet of another origin (a page of
		// another host, a CDN's css): allowed
		w.Header().Set("Access-Control-Allow-Origin", "*")
		files.ServeHTTP(w, r)
	})
}
