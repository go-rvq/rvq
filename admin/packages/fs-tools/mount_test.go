package fs_tools

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/packages/fs-tools/fs"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
)

// The WebDAV verifies each request: the directories of the page by its
// permissions (reading @get, writing !write); a mount — the draft of the
// user — by its own, its hidden paths not served; nothing moves between them.
func TestWebDavPermissions(t *testing.T) {
	data, drafts := t.TempDir(), t.TempDir()
	os.MkdirAll(filepath.Join(drafts, "ana", ".git"), 0o755)
	os.WriteFile(filepath.Join(drafts, "ana", ".git", "config"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(drafts, "ana", "index.gadx"), []byte("old"), 0o644)
	os.MkdirAll(filepath.Join(drafts, "ana", "sub"), 0o755)

	var asked []string // what the mount was asked: "METHOD path dest"
	p := presets.New(i18n.New()).URIPrefix("/admin")
	pb := perm.New().AllowAll()
	pb.CreatePolicies(perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(perm.Anything).On("admin:/fs-tools:!write"))
	p.Permission(pb.SubjectsFunc(func(*http.Request) []string { return []string{"editor"} }))
	b, err := New(p, p.I18n(), nil)
	if err != nil {
		t.Fatal(err)
	}
	root, err := fs.Parse("w,/data," + data)
	if err != nil {
		t.Fatal(err)
	}
	b.SetFS(root)
	b.AddMount(&Mount{
		Name: "site-files",
		Dir:  func(r *http.Request) (string, error) { return filepath.Join(drafts, r.Header.Get("X-User")), nil },
		Allowed: func(r *http.Request, a MountAccess) bool {
			asked = append(asked, r.Method+" "+a.Path+" "+a.Dest)
			return !a.Write || r.Header.Get("X-Edit") == "1"
		},
		Hidden: func(p string) bool {
			return p == "/.git" || strings.HasPrefix(p, "/.git/")
		},
	})
	if err := b.Install(p); err != nil {
		t.Fatal(err)
	}
	p.Build(http.NewServeMux())
	h := b.davHandler()

	do := func(method, path, body string, header ...string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, b.davPath+path, strings.NewReader(body))
		r.Header.Set("X-User", "ana")
		for i := 0; i+1 < len(header); i += 2 {
			r.Header.Set(header[i], header[i+1])
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}

	// the root: the page's directories and the mounts
	if w := do("PROPFIND", "/", "", "Depth", "1"); !strings.Contains(w.Body.String(), "site-files") || !strings.Contains(w.Body.String(), "data") {
		t.Errorf("the root: %d %s", w.Code, w.Body.String())
	}
	// the draft: read, written with the permission of editing
	if w := do("GET", "/site-files/index.gadx", ""); w.Body.String() != "old" {
		t.Errorf("reading the draft: %d %q", w.Code, w.Body.String())
	}
	if w := do("PUT", "/site-files/index.gadx", "new"); w.Code != http.StatusForbidden {
		t.Errorf("writing the draft with no permission: %d", w.Code)
	}
	if w := do("PUT", "/site-files/index.gadx", "new", "X-Edit", "1"); w.Code >= 300 {
		t.Errorf("writing the draft: %d", w.Code)
	}
	if b, _ := os.ReadFile(filepath.Join(drafts, "ana", "index.gadx")); string(b) != "new" {
		t.Errorf("the draft has %q", b)
	}
	// its .git: not served, not listed
	if w := do("GET", "/site-files/.git/config", ""); w.Code != http.StatusForbidden {
		t.Errorf(".git served: %d", w.Code)
	}
	if w := do("PROPFIND", "/site-files/", "", "Depth", "1"); strings.Contains(w.Body.String(), ".git") {
		t.Errorf(".git listed: %s", w.Body.String())
	}
	// a move in the draft: the mount told both paths, once
	asked = nil
	if w := do("MOVE", "/site-files/index.gadx", "", "Destination", b.davPath+"/site-files/sub/index.gadx", "X-Edit", "1"); w.Code >= 300 {
		t.Errorf("moving in the draft: %d", w.Code)
	}
	if len(asked) != 1 || asked[0] != "MOVE /index.gadx /sub/index.gadx" {
		t.Errorf("the mount was asked %q", asked)
	}
	do("MOVE", "/site-files/sub/index.gadx", "", "Destination", b.davPath+"/site-files/index.gadx", "X-Edit", "1")
	// nothing moves out of the draft
	if w := do("MOVE", "/site-files/index.gadx", "", "Destination", b.davPath+"/data/x", "X-Edit", "1"); w.Code < 300 {
		t.Errorf("moved out of the draft: %d", w.Code)
	}
	// the page's directories: read with @get, written with !write — denied
	os.WriteFile(filepath.Join(data, "a.txt"), []byte("a"), 0o644)
	if w := do("GET", "/data/a.txt", ""); w.Body.String() != "a" {
		t.Errorf("reading the page's: %d %q", w.Code, w.Body.String())
	}
	if w := do("PUT", "/data/a.txt", "b"); w.Code != http.StatusForbidden {
		t.Errorf("writing the page's with no !write: %d", w.Code)
	}
}
