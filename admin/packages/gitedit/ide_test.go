package gitedit

import (
	"encoding/json"
	"net"
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
	h.Allowed = func(_ *http.Request, op IdeOp, _ string) bool { return op == IdeRead || write }

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
	for _, op := range []string{"run", "eval", "inspect", "debug/start"} {
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

// Each change asks its own permission, on the paths it touches: creating a
// file that is not there, editing one that is, renaming in its folder, moving
// to another, deleting, importing (and creating or editing what it writes);
// the workspace tells the IDE which it may.
func TestIDEOps(t *testing.T) {
	repo, _ := site(t)
	h := NewIDE(repo)
	h.DraftKey = func(*http.Request) string { return "u1" }
	var granted map[IdeOp]bool
	var asked []string
	names := map[IdeOp]string{IdeRead: "read", IdeCreate: "create", IdeEdit: "edit", IdeRename: "rename",
		IdeMove: "move", IdeDelete: "delete", IdeImport: "import"}
	h.Allowed = func(_ *http.Request, op IdeOp, p string) bool {
		if op != IdeRead {
			asked = append(asked, names[op]+" "+p)
		}
		return op == IdeRead || granted[op]
	}
	do := func(method, path, body string, ops ...IdeOp) (*httptest.ResponseRecorder, []string) {
		granted, asked = map[IdeOp]bool{}, nil
		for _, o := range ops {
			granted[o] = true
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, path, strings.NewReader(body)))
		return w, asked
	}
	draft := filepath.Join(repo.DraftsDir, "u1")

	for _, c := range []struct {
		name, method, path, body string
		need                     []string // what it asks, in order
		op                       IdeOp    // the one that lets it
	}{
		{"a new file", "PUT", "/api/ide/file", `{"path":"static/new.css","content":"x"}`, []string{"create static/new.css"}, IdeCreate},
		{"a file there", "PUT", "/api/ide/file", `{"path":"index.gadx","content":"y\n"}`, []string{"edit index.gadx"}, IdeEdit},
		{"a folder", "POST", "/api/ide/mkdir", `{"path":"static/img"}`, []string{"create static/img"}, IdeCreate},
		{"a rename", "POST", "/api/ide/rename", `{"path":"static/new.css","to":"static/site.css"}`, []string{"rename static/new.css", "rename static/site.css"}, IdeRename},
		{"a move", "POST", "/api/ide/rename", `{"path":"static/site.css","to":"static/img/site.css"}`, []string{"move static/site.css", "move static/img/site.css"}, IdeMove},
		{"a delete", "POST", "/api/ide/delete", `{"path":"static/img/site.css"}`, []string{"delete static/img/site.css"}, IdeDelete},
	} {
		w, a := do(c.method, c.path, c.body)
		if w.Code != http.StatusForbidden || strings.Join(a, ",") != c.need[0] {
			t.Errorf("%s with no permission: %d, asked %q", c.name, w.Code, a)
		}
		w, a = do(c.method, c.path, c.body, c.op)
		if w.Code != 200 || strings.Join(a, ",") != strings.Join(c.need, ",") {
			t.Errorf("%s: %d %s, asked %q", c.name, w.Code, w.Body.String(), a)
		}
	}

	// an upload: !import, and !create of each new file (!edit of one there)
	upload := `{"files":[{"path":"static/a.css","content":"a"},{"path":"index.gadx","content":"z"}]}`
	if w, _ := do("POST", "/api/ide/upload", upload, IdeImport, IdeCreate); w.Code != http.StatusForbidden {
		t.Errorf("an upload over a file with no !edit: %d", w.Code)
	}
	if w, a := do("POST", "/api/ide/upload", upload, IdeImport, IdeCreate, IdeEdit); w.Code != 200 ||
		strings.Join(a, ",") != "import static/a.css,create static/a.css,import index.gadx,edit index.gadx" {
		t.Errorf("an upload: %d %s, asked %q", w.Code, w.Body.String(), a)
	}
	if w, _ := do("POST", "/api/ide/upload", `{"files":[{"path":"b.css","content":"b"}]}`, IdeCreate, IdeEdit); w.Code != http.StatusForbidden {
		t.Errorf("an upload with no !import: %d", w.Code)
	}
	if b, _ := os.ReadFile(filepath.Join(draft, "static", "a.css")); string(b) != "a" {
		t.Errorf("uploaded: %q", b)
	}

	// a download from a URL of the server's own network: refused
	local := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("secret")) }))
	defer local.Close()
	if w, _ := do("POST", "/api/ide/fetch", `{"url":"`+local.URL+`/x","path":"static/x.txt"}`, IdeImport, IdeCreate); w.Code == 200 ||
		!strings.Contains(w.Body.String(), "not public") {
		t.Errorf("a download of the server's network: %d %s", w.Code, w.Body.String())
	}
	if _, err := os.Stat(filepath.Join(draft, "static", "x.txt")); err == nil {
		t.Error("downloaded from the server's network")
	}

	// the workspace: what the user may do
	w, _ := do("GET", "/api/ide/workspace", "", IdeEdit, IdeImport)
	var ws struct{ Actions map[string]bool }
	if err := json.Unmarshal(w.Body.Bytes(), &ws); err != nil {
		t.Fatal(err)
	}
	if !ws.Actions["edit"] || !ws.Actions["import"] || ws.Actions["create"] || ws.Actions["delete"] || len(ws.Actions) != 6 {
		t.Errorf("the workspace's actions: %v", ws.Actions)
	}
}

// PublicIP: the internet, not the server's own network.
func TestPublicIP(t *testing.T) {
	for ip, want := range map[string]bool{
		"93.184.216.34": true, "2606:4700::1111": true,
		"127.0.0.1": false, "::1": false, "10.1.2.3": false, "192.168.0.1": false, "172.16.5.4": false,
		"169.254.169.254": false, "0.0.0.0": false, "100.64.0.1": false, "fd00::1": false, "fe80::1": false,
	} {
		if got := PublicIP(net.ParseIP(ip)); got != want {
			t.Errorf("%s: %v", ip, got)
		}
	}
}
