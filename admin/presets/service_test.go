package presets

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServices(t *testing.T) {
	b := &Builder{prefix: "/admin"}
	for _, bad := range []string{"", "Bad", "a/b", "!x"} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%q accepted", bad)
				}
			}()
			b.Service(bad)
		}()
	}
	s := b.Service("files")
	if b.Service("files") != s || s.Path() != "/admin/!files" {
		t.Errorf("service %q %s", s.Name(), s.Path())
	}
	if !b.IsServicePath("/admin/!files/x") || b.IsServicePath("/admin/files") || b.IsServicePath("/admin/!") {
		t.Error("IsServicePath")
	}

	// under the admin's login: in its mux
	s.Handle("GET /{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("file " + r.PathValue("id"))) }))
	mux := http.NewServeMux()
	b.setupServices(mux)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("GET", "/admin/!files/7", nil))
	if w.Body.String() != "file 7" {
		t.Errorf("handle: %q", w.Body.String())
	}

	// its own login: in front of the admin; a protocol's path mounted too;
	// the rest the admin's
	admin := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("admin")) })
	b.Service("git").Mount("/admin/site.git/").Outer(func(admin http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("git")) })
	})
	h := b.OuterServices(admin)
	for path, want := range map[string]string{"/admin/!git/x": "git", "/admin/site.git/info/refs": "git", "/admin/users": "admin", "/admin/!files/7": "admin"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Body.String() != want {
			t.Errorf("%s: %q", path, w.Body.String())
		}
	}
	if names := b.Services(); len(names) != 2 || names[0].Name() != "files" || !names[1].IsOuter() {
		t.Errorf("services %v", names)
	}
}

func TestServiceRoutes(t *testing.T) {
	b := &Builder{prefix: "/admin/"}
	s := b.Service("git").
		Route("GET", "/{id}", ServiceAuthToken+", "+ServiceAuthSession, nil).
		Route("GET", "<prefix>/{name}.git/…", ServiceAuthBasic, nil)
	rs := s.Routes()
	if got := s.RoutePath(rs[0]); got != "/admin/!git/{id}" {
		t.Errorf("route path %q", got)
	}
	if got := s.RoutePath(rs[1]); got != "/admin/{name}.git/…" {
		t.Errorf("protocol route path %q", got)
	}
	r := &ServiceRoute{Auth: ServiceAuthSession + ", !terminal"}
	if got := r.GetAuth(context.Background()); got != "session, !terminal" {
		t.Errorf("auth %q", got)
	}
}

// TestServiceMessages: every text of the services is translated.
func TestServiceMessages(t *testing.T) {
	for _, m := range []*Messages{Messages_en_US, Messages_pt_BR} {
		for _, v := range []string{m.ServiceAssetsTitle, m.ServiceLoginTitle, m.ServiceGadcodeTitle, m.ServiceUserDocsTitle,
			m.RouteAssetsExtra, m.RouteLoginDone, m.RouteGadcodeDoc, m.RouteUserDocsAsset, m.ServiceAuthBasic, m.ServiceAuthSignedLink} {
			if v == "" {
				t.Fatal("a text of the services is empty")
			}
		}
	}
}

// TestServiceMountOnly: a service of a protocol's paths alone is not at its
// /!<name>.
func TestServiceMountOnly(t *testing.T) {
	b := &Builder{prefix: "/admin"}
	b.Service("git").MountOnly().Mount("/admin/site.git/").
		Outer(func(admin http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("git")) })
		})
	h := b.OuterServices(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("admin")) }))
	for path, want := range map[string]string{"/admin/site.git/info/refs": "git", "/admin/!git/x": "admin"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Body.String() != want {
			t.Errorf("%s: %q, want %q", path, w.Body.String(), want)
		}
	}
}
