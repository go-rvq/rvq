package portal

import (
	"crypto/sha256"
	"encoding/base64"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-rvq/rvq/js"
)

// ScriptPath is where the bundle of the site's portals is served: the
// server's path with ".js" (/_rvq/portal.js).
func (s *Server) ScriptPath() string { return s.Path + ".js" }

// The bundle runs once Vue is there — the page's global: from its CDN, or
// the app's copy when the CDN did not answer (js/api/vendor), which comes
// later than the page's deferred scripts —: by rvqVendor.when (the vendor's
// bootstrap), else as soon as it is seen. Run, it says so: the event
// rvq-portal:ready (what waits for RvqPortal: a preview of a form).
const (
	scriptPrefix = `(function(){var run=function(){`
	scriptSuffix = "\n;try{document.dispatchEvent(new Event('rvq-portal:ready'))}catch(e){}};" +
		`if(window.Vue)run();else if(window.rvqVendor)rvqVendor.when('Vue',run);` +
		`else{var n=0,t=setInterval(function(){if(window.Vue||++n>400){clearInterval(t);if(window.Vue)run()}},50)}})();`
)

// script is the bundle of the portals with the server's path written in it,
// and its version (a hash of it); made once.
func (s *Server) script() (body []byte, version string, err error) {
	s.scriptOnce.Do(func() {
		var b []byte
		if b, s.scriptErr = fs.ReadFile(js.Portal, "portal/dist/rvq-portal.umd.js"); s.scriptErr != nil {
			return
		}
		s.scriptBody = []byte(scriptPrefix + strings.ReplaceAll(string(b), strconv.Quote(DefaultPath), strconv.Quote(s.Path)) + scriptSuffix)
		sum := sha256.Sum256(s.scriptBody)
		s.scriptVersion = base64.RawURLEncoding.EncodeToString(sum[:12])
	})
	return s.scriptBody, s.scriptVersion, s.scriptErr
}

// ScriptURL is the address of the bundle for a page: ScriptPath with its
// version (?v=) — a page of a new bundle never runs the old one a browser
// kept.
func (s *Server) ScriptURL() string {
	if _, v, err := s.script(); err == nil {
		return s.ScriptPath() + "?v=" + v
	}
	return s.ScriptPath()
}

// ScriptHandler serves the bundle of the portals of a site (js.Portal:
// rvq-portal.umd.js) with the server's path written in it — the page loads
// it, and names no other address.
func (s *Server) ScriptHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, version, err := s.script()
		if err != nil {
			http.Error(w, "no bundle", http.StatusInternalServerError)
			return
		}
		etag := `"` + version + `"`
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Header().Set("ETag", etag)
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		_, _ = w.Write(body)
	})
}

// Mount puts the server at its path and its bundle at ScriptPath.
func (s *Server) Mount(mux *http.ServeMux) {
	mux.Handle(s.Path, s)
	mux.Handle(s.ScriptPath(), s.ScriptHandler())
}
