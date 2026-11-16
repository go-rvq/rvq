package explorer

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/net/webdav"
)

func call(t *testing.T, x *Explorer, method, op string, q url.Values, body io.Reader, ctype string) (int, map[string]any, []byte) {
	t.Helper()
	r := httptest.NewRequest(method, Prefix+op+"?"+q.Encode(), body)
	if ctype != "" {
		r.Header.Set("Content-Type", ctype)
	}
	w := httptest.NewRecorder()
	x.ServeHTTP(w, r)
	var m map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &m)
	return w.Code, m, w.Body.Bytes()
}

// The operations over a folder (webdav.Dir): write, read, list, mkdir,
// rename, copy, delete, upload (a folder), download (a file, a folder as a
// zip); what changed told.
func TestExplorer(t *testing.T) {
	root := t.TempDir()
	var changes []string
	x := New(webdav.Dir(root), Options{Changed: func(_ *http.Request, op, p string) { changes = append(changes, op+" "+p) }})
	if c, m, _ := call(t, x, "POST", "mkdir", url.Values{"path": {"a/b"}}, nil, ""); c != 200 {
		t.Fatalf("mkdir: %d %v", c, m)
	}
	if c, m, _ := call(t, x, "POST", "write", url.Values{"path": {"a/b/c.txt"}}, strings.NewReader("olá"), ""); c != 200 {
		t.Fatalf("write: %d %v", c, m)
	}
	if _, m, _ := call(t, x, "GET", "read", url.Values{"path": {"a/b/c.txt"}}, nil, ""); m["content"] != "olá" || m["writable"] != true {
		t.Fatalf("read: %v", m)
	}
	_, m, _ := call(t, x, "GET", "list", url.Values{"path": {"a"}}, nil, "")
	if es := m["entries"].([]any); len(es) != 1 || es[0].(map[string]any)["name"] != "b" || es[0].(map[string]any)["dir"] != true {
		t.Fatalf("list: %v", m)
	}
	if c, _, _ := call(t, x, "POST", "copy", url.Values{"path": {"a"}, "to": {"d"}}, nil, ""); c != 200 {
		t.Fatal("copy")
	}
	if b, _ := os.ReadFile(filepath.Join(root, "d", "b", "c.txt")); string(b) != "olá" {
		t.Fatalf("copied: %q", b)
	}
	if c, _, _ := call(t, x, "POST", "copy", url.Values{"path": {"a"}, "to": {"a/x"}}, nil, ""); c != 400 {
		t.Errorf("a folder into itself: %d", c)
	}
	if c, _, _ := call(t, x, "POST", "rename", url.Values{"path": {"d"}, "to": {"a"}}, nil, ""); c != 409 {
		t.Errorf("rename over one there: %d", c)
	}
	if c, _, _ := call(t, x, "POST", "rename", url.Values{"path": {"d"}, "to": {"e"}}, nil, ""); c != 200 {
		t.Fatal("rename")
	}
	// a folder uploaded: its files with their folders
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fw, _ := mw.CreateFormFile("file", "up/f.txt")
	fw.Write([]byte("up"))
	mw.Close()
	if c, m, _ := call(t, x, "POST", "upload", url.Values{"path": {"a"}}, &body, mw.FormDataContentType()); c != 200 {
		t.Fatalf("upload: %d %v", c, m)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "a", "up", "f.txt")); string(b) != "up" {
		t.Fatalf("uploaded: %q", b)
	}
	// a folder downloaded: a zip
	_, _, zb := call(t, x, "GET", "download", url.Values{"path": {"a"}}, nil, "")
	zr, err := zip.NewReader(bytes.NewReader(zb), int64(len(zb)))
	if err != nil || len(zr.File) != 2 {
		t.Fatalf("zip: %v %d", err, len(zr.File))
	}
	if c, _, _ := call(t, x, "POST", "delete", url.Values{"path": {"e", "a/up"}}, nil, ""); c != 200 {
		t.Fatal("delete")
	}
	if _, err := os.Stat(filepath.Join(root, "e")); !os.IsNotExist(err) {
		t.Errorf("deleted: %v", err)
	}
	if c, _, _ := call(t, x, "POST", "delete", url.Values{"path": {"/"}}, nil, ""); c != 400 {
		t.Errorf("the root deleted: %d", c)
	}
	if c, _, _ := call(t, x, "GET", "read", url.Values{"path": {"../../etc/passwd"}}, nil, ""); c == 200 {
		t.Errorf("out of the root read")
	}
	if len(changes) < 6 {
		t.Errorf("changes: %v", changes)
	}
}

// Allow refuses: a folder read only — read, listed; not written, renamed,
// deleted, uploaded to —; a folder not seen at all.
func TestExplorerAllow(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "ro"), 0o755)
	os.MkdirAll(filepath.Join(root, "secret"), 0o755)
	os.WriteFile(filepath.Join(root, "ro", "f.txt"), []byte("x"), 0o644)
	x := New(webdav.Dir(root), Options{Allow: func(_ *http.Request, p string, write bool) bool {
		if strings.HasPrefix(p, "/secret") {
			return false
		}
		return !(write && strings.HasPrefix(p, "/ro"))
	}})
	if _, m, _ := call(t, x, "GET", "read", url.Values{"path": {"ro/f.txt"}}, nil, ""); m["content"] != "x" || m["writable"] != false {
		t.Errorf("read only, read: %v", m)
	}
	for _, c := range []struct {
		method, op string
		q          url.Values
	}{
		{"POST", "write", url.Values{"path": {"ro/f.txt"}}},
		{"POST", "delete", url.Values{"path": {"ro/f.txt"}}},
		{"POST", "rename", url.Values{"path": {"ro/f.txt"}, "to": {"g.txt"}}},
		{"POST", "copy", url.Values{"path": {"g"}, "to": {"ro/g"}}},
		{"POST", "mkdir", url.Values{"path": {"ro/n"}}},
		{"GET", "list", url.Values{"path": {"secret"}}},
		{"GET", "download", url.Values{"path": {"secret"}}},
	} {
		if code, _, _ := call(t, x, c.method, c.op, c.q, strings.NewReader("y"), ""); code != http.StatusForbidden {
			t.Errorf("%s %v: %d", c.op, c.q, code)
		}
	}
	if b, _ := os.ReadFile(filepath.Join(root, "ro", "f.txt")); string(b) != "x" {
		t.Errorf("written: %q", b)
	}
}

// A large file in parts, out of order; its size checked; not over one
// there.
func TestExplorerParts(t *testing.T) {
	root := t.TempDir()
	x := New(webdav.Dir(root), Options{})
	id := "abcdefghijklmnopqrst"
	q := func(extra ...string) url.Values {
		v := url.Values{"path": {"d"}, "name": {"big.bin"}, "id": {id}, "total": {"10"}}
		for i := 0; i+1 < len(extra); i += 2 {
			v.Set(extra[i], extra[i+1])
		}
		return v
	}
	if c, m, _ := call(t, x, "POST", "upload-part", q("offset", "0"), strings.NewReader("abcde"), ""); c != 200 {
		t.Fatalf("part: %d %v", c, m)
	}
	if c, _, _ := call(t, x, "POST", "upload-done", q(), nil, ""); c != 400 {
		t.Errorf("done short: %d", c)
	}
	call(t, x, "POST", "upload-part", q("offset", "5"), strings.NewReader("fghij"), "")
	if c, m, _ := call(t, x, "POST", "upload-done", q(), nil, ""); c != 200 {
		t.Fatalf("done: %d %v", c, m)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "d", "big.bin")); string(b) != "abcdefghij" {
		t.Fatalf("file: %q", b)
	}
	if c, _, _ := call(t, x, "POST", "upload-part", q("offset", "8"), strings.NewReader("xyz"), ""); c != 400 {
		t.Errorf("past the end: %d", c)
	}
	if c, _, _ := call(t, x, "POST", "upload-abort", q(), nil, ""); c != 200 {
		t.Errorf("abort: %d", c)
	}
	call(t, x, "POST", "upload-part", q("offset", "0"), strings.NewReader("0123456789"), "")
	if c, _, _ := call(t, x, "POST", "upload-done", q(), nil, ""); c != 409 {
		t.Errorf("over one there: %d", c)
	}
	if c, _, _ := call(t, x, "POST", "upload-done", q("replace", "1"), nil, ""); c != 200 {
		t.Errorf("replaced: %d", c)
	}
}

// An archive extracted, nothing out of its folder; off unless Extract.
func TestExplorerExtract(t *testing.T) {
	root := t.TempDir()
	var b bytes.Buffer
	zw := zip.NewWriter(&b)
	for _, n := range []string{"a.txt", "sub/b.txt", "../evil.txt"} {
		f, _ := zw.Create(n)
		f.Write([]byte(n))
	}
	zw.Close()
	os.WriteFile(filepath.Join(root, "x.zip"), b.Bytes(), 0o644)
	if c, _, _ := call(t, New(webdav.Dir(root), Options{}), "POST", "extract", url.Values{"path": {"x.zip"}}, nil, ""); c != 404 {
		t.Errorf("off: %d", c)
	}
	x := New(webdav.Dir(root), Options{Extract: true})
	if c, m, _ := call(t, x, "POST", "extract", url.Values{"path": {"x.zip"}}, nil, ""); c != 200 || m["files"] != float64(2) {
		t.Fatalf("extract: %d %v", c, m)
	}
	if bb, _ := os.ReadFile(filepath.Join(root, "x", "sub", "b.txt")); string(bb) != "sub/b.txt" {
		t.Errorf("extracted: %q", bb)
	}
	if _, err := os.Stat(filepath.Join(root, "evil.txt")); !os.IsNotExist(err) {
		t.Errorf("out of the folder: %v", err)
	}
}
