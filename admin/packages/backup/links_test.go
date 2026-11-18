package backup

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/packages/fs-tools/explorer"
)

// tarFiles are the files of the tar.gz of r: their contents by name.
func tarFiles(t *testing.T, r io.Reader) map[string]string {
	t.Helper()
	gz, err := gzip.NewReader(r)
	if err != nil {
		t.Fatal(err)
	}
	tr := tar.NewReader(gz)
	out := map[string]string{}
	for {
		hd, err := tr.Next()
		if err == io.EOF {
			return out
		}
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(tr)
		out[hd.Name] = string(b)
	}
}

// A download link: the state whole as tar.gz, with no login, as many times
// as wanted; its uses counted; expired, revoked, of another state or of a
// state removed, not found.
func TestLink(t *testing.T) {
	b, _ := setup(t, nil)
	b.p = nil
	s := queue(t, b, &BackupState{})
	if err := b.Build(context.Background(), s.ID); err != nil {
		t.Fatal(err)
	}
	queued := queue(t, b, &BackupState{})
	b.clientBase = "/admin/!backup"
	admin := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(299) })
	srv := httptest.NewServer(b.ClientHandler(admin))
	defer srv.Close()
	get := func(p string) *http.Response {
		t.Helper()
		res, err := http.Get(srv.URL + p)
		if err != nil {
			t.Fatal(err)
		}
		return res
	}

	if _, _, err := b.NewLink(httptest.NewRequest("GET", "/", nil), queued.ID); err == nil {
		t.Error("a link of a state not ready")
	}
	tok, l, err := b.NewLink(httptest.NewRequest("GET", "/", nil), s.ID)
	if err != nil || time.Until(l.ExpiresAt) < LinkTTL-time.Minute {
		t.Fatal(err, l)
	}
	var kept BackupLink
	b.db.First(&kept, "id = ?", l.ID)
	if kept.TokenSum == tok || kept.TokenSum != tokenSum(tok) {
		t.Error("the token kept, not its sum")
	}

	// twice: the archive each time, no Authorization asked
	for i := 0; i < 2; i++ {
		res := get(b.LinkPath(s.ID, tok))
		if res.StatusCode != 200 || res.Header.Get("Content-Type") != "application/gzip" {
			t.Fatalf("archive: %d", res.StatusCode)
		}
		files := tarFiles(t, res.Body)
		res.Body.Close()
		if files["public/a.txt"] != "hello world" || files["public/b.txt"] != "bee" {
			t.Fatalf("files: %v", files)
		}
	}
	b.db.First(&kept, "id = ?", l.ID)
	if kept.Uses != 2 || kept.LastUsedAt == nil || kept.LastIP == "" {
		t.Errorf("uses: %+v", kept)
	}

	// a wrong token, the token of another state: not found
	for _, p := range []string{b.LinkPath(s.ID, "nope"), b.LinkPath(queued.ID, tok), b.ClientPath() + "/" + s.ID + "/archive/" + tok} {
		if res := get(p); res.StatusCode != 404 {
			t.Errorf("%s: %d", p, res.StatusCode)
		}
	}
	// the rest of the client's paths with no Authorization: the admin's
	if res := get(b.ClientPath() + "/latest"); res.StatusCode != 299 {
		t.Errorf("latest, no key: %d", res.StatusCode)
	}

	// expired, revoked: not found
	b.db.Model(&BackupLink{}).Where("id = ?", l.ID).Update("expires_at", time.Now().Add(-time.Second))
	if res := get(b.LinkPath(s.ID, tok)); res.StatusCode != 404 {
		t.Errorf("expired: %d", res.StatusCode)
	}
	tok2, l2, _ := b.NewLink(httptest.NewRequest("GET", "/", nil), s.ID)
	now := time.Now()
	b.db.Model(&BackupLink{}).Where("id = ?", l2.ID).Update("revoked_at", &now)
	if res := get(b.LinkPath(s.ID, tok2)); res.StatusCode != 404 {
		t.Errorf("revoked: %d", res.StatusCode)
	}

	// the state removed: its links with it
	if err := b.Remove(s.ID); err != nil {
		t.Fatal(err)
	}
	var n int64
	b.db.Model(&BackupLink{}).Where("state_id = ?", s.ID).Count(&n)
	if n != 0 {
		t.Errorf("links of a state removed: %d", n)
	}
}

// The files of a state in the explorer: the manifest's tree, read; no
// change; a folder and several items downloaded as tar.gz (the core's).
func TestStateFS(t *testing.T) {
	b, _ := setup(t, nil)
	s := queue(t, b, &BackupState{})
	if err := b.Build(context.Background(), s.ID); err != nil {
		t.Fatal(err)
	}
	sfs, err := b.stateFSOf(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	x := explorer.New(sfs, explorer.Options{
		Allow: func(_ *http.Request, _ string, write bool) bool { return !write },
		Archive: func(w io.Writer, _ *http.Request, base string, paths []string) error {
			_, err := b.core.WriteArchive(w, s.ID, archiveOptions(base, paths))
			return err
		},
	})
	call := func(method, op string, q url.Values) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		x.ServeHTTP(w, httptest.NewRequest(method, explorer.Prefix+op+"?"+q.Encode(), strings.NewReader("x")))
		return w
	}

	w := call("GET", "list", url.Values{"path": {"/"}})
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"name":"public"`) || !strings.Contains(w.Body.String(), `"write":false`) {
		t.Fatalf("list /: %d %s", w.Code, w.Body)
	}
	w = call("GET", "list", url.Values{"path": {"/public"}})
	var names []string
	for _, n := range []string{"a.txt", "b.txt"} {
		if strings.Contains(w.Body.String(), `"name":"`+n+`"`) {
			names = append(names, n)
		}
	}
	if len(names) != 2 {
		t.Errorf("list /public: %s", w.Body)
	}
	if w = call("GET", "read", url.Values{"path": {"/public/a.txt"}}); !strings.Contains(w.Body.String(), "hello world") {
		t.Errorf("read: %s", w.Body)
	}
	if w = call("GET", "download", url.Values{"path": {"/public/b.txt"}}); w.Body.String() != "bee" {
		t.Errorf("download: %q", w.Body)
	}
	for _, op := range []string{"write", "mkdir", "delete"} {
		if w = call("POST", op, url.Values{"path": {"/public/a.txt"}}); w.Code != http.StatusForbidden {
			t.Errorf("%s: %d", op, w.Code)
		}
	}
	if _, err := sfs.OpenFile(context.Background(), "/public/a.txt", os.O_RDWR, 0); err == nil {
		t.Error("opened to write")
	}

	// a folder, several items: tar.gz, the names under the folder they are in
	w = call("GET", "download", url.Values{"path": {"/public"}})
	files := tarFiles(t, w.Body)
	if files["public/a.txt"] != "hello world" || files["public/b.txt"] != "bee" {
		t.Errorf("folder: %v", keysOf(files))
	}
	w = call("GET", "download", url.Values{"path": {"/public/a.txt", "/public/b.txt"}})
	if files = tarFiles(t, w.Body); len(files) != 2 || files["a.txt"] != "hello world" || files["b.txt"] != "bee" {
		t.Errorf("items: %v", keysOf(files))
	}
}

func keysOf(m map[string]string) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
