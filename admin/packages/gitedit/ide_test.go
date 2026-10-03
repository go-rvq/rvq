package gitedit

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The IDE serves the files of the editor's draft and the language's
// operations; not the ones that run code, not .git; reading and writing as
// the request may.
func TestIDE(t *testing.T) {
	repo, served := site(t)
	h := NewIDE(repo)
	h.DraftKey = func(*http.Request) string { return "u1" }
	write := true
	h.Allowed = func(_ *http.Request, p IdePerm) bool { return p == IdeRead || write }

	do := func(method, path, body string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, path, strings.NewReader(body)))
		return w
	}
	if w := do("GET", "/api/ide/file?path=index.gadx", ""); w.Code != 200 || !strings.Contains(w.Body.String(), "old") {
		t.Errorf("reading a file: %d %s", w.Code, w.Body.String())
	}
	if w := do("GET", "/api/ide/workspace", ""); strings.Contains(w.Body.String(), repo.DraftsDir) {
		t.Errorf("the workspace tells the path on the server: %s", w.Body.String())
	}
	if w := do("PUT", "/api/ide/file", `{"path":"index.gadx","content":"edited\n"}`); w.Code != 200 {
		t.Fatalf("writing: %d %s", w.Code, w.Body.String())
	}
	if b, _ := os.ReadFile(filepath.Join(repo.DraftsDir, "u1", "index.gadx")); string(b) != "edited\n" {
		t.Errorf("the draft has %q", b)
	}
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "old\n" {
		t.Errorf("the site changed: %q", b)
	}

	for _, c := range []struct{ method, path, body string }{
		{"GET", "/api/ide/file?path=.git/config", ""},
		{"PUT", "/api/ide/file", `{"path":".git/hooks/post-commit","content":"x"}`},
		{"POST", "/api/ide/rename", `{"path":"index.gadx","to":"sub/.git/x"}`},
	} {
		if w := do(c.method, c.path, c.body); w.Code != http.StatusForbidden {
			t.Errorf("%s %s %s: %d", c.method, c.path, c.body, w.Code)
		}
	}
	for _, op := range []string{"run", "eval", "inspect", "fetch", "debug/start"} {
		if w := do("POST", "/api/ide/"+op, `{}`); w.Code != http.StatusNotFound {
			t.Errorf("%s served: %d", op, w.Code)
		}
	}
	write = false
	if w := do("PUT", "/api/ide/file", `{"path":"index.gadx","content":"x"}`); w.Code != http.StatusForbidden {
		t.Errorf("writing with no permission: %d", w.Code)
	}
	if w := do("GET", "/api/ide/tree", ""); w.Code != 200 {
		t.Errorf("reading with no permission to write: %d", w.Code)
	}
}
