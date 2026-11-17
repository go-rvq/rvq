package backup

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/x/packages/backup"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type fakeUser struct {
	ID   uuid.UUID
	Name string
}

func (u *fakeUser) GetAccountName() string { return u.Name }

// fakeKeys: a key is "rk_<anything>"; "rk_bad" refused.
type fakeKeys struct{ served int }

func (fakeKeys) IsKeyCode(code string) bool { return strings.HasPrefix(code, "rk_") }
func (fakeKeys) AuthKey(r *http.Request, code string) (any, context.Context, error) {
	if code == "rk_bad" {
		return nil, nil, errors.New("refused")
	}
	return &fakeUser{ID: uuid.New(), Name: "ana"}, r.Context(), nil
}
func (k *fakeKeys) KeyServed(*http.Request, int) { k.served++ }

// setup is a Builder over a database and an application's directory with
// a source of two files; block, when given, makes the source wait for it.
func setup(t *testing.T, block chan struct{}) (*Builder, string) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "db.sqlite")), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "public"), 0o755)
	os.WriteFile(filepath.Join(root, "public/a.txt"), []byte("hello world"), 0o644)
	os.WriteFile(filepath.Join(root, "public/b.txt"), []byte("bee"), 0o644)
	core, err := backup.New(backup.Options{Root: root, Dir: filepath.Join(root, "data/backup")})
	if err != nil {
		t.Fatal(err)
	}
	core.Source("site", func(ctx context.Context, st *backup.State) error {
		if block != nil {
			select {
			case <-block:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		st.FS.Warn("one warning")
		return st.FS.AddTree("public")
	})
	core.Source("other", func(ctx context.Context, st *backup.State) error { return nil })
	b := New(db, core, nil).Secret([]byte("secret"))
	b.keys = &fakeKeys{}
	b.allowDownload = func(*http.Request) bool { return true }
	if err := b.AutoMigrate(); err != nil {
		t.Fatal(err)
	}
	return b, root
}

func queue(t *testing.T, b *Builder, s *BackupState) *BackupState {
	t.Helper()
	if s.ID == "" {
		s.ID = NewID(time.Now())
	}
	s.Status = StatusQueued
	if err := b.db.Create(s).Error; err != nil {
		t.Fatal(err)
	}
	return s
}

func reload(t *testing.T, b *Builder, id string) *BackupState {
	t.Helper()
	var s BackupState
	if err := b.db.First(&s, "id = ?", id).Error; err != nil {
		t.Fatal(err)
	}
	return &s
}

func TestBuild(t *testing.T) {
	b, _ := setup(t, nil)
	s := queue(t, b, &BackupState{Description: "first"})
	if err := b.Build(context.Background(), s.ID); err != nil {
		t.Fatal(err)
	}
	s = reload(t, b, s.ID)
	if s.Status != StatusReady || s.Files != 2 || s.Bytes != int64(len("hello worldbee")) || s.Partial ||
		s.ManifestSum == "" || s.StartedAt == nil || s.FinishedAt == nil {
		t.Errorf("state %+v", s)
	}
	if len(s.Results) != 2 || s.Results[0].Name != "site" || len(s.Results[0].Warnings) != 1 {
		t.Errorf("results %+v", s.Results)
	}
	if _, err := os.Stat(b.core.ManifestPath(s.ID)); err != nil {
		t.Error("no manifest")
	}
	// made: not again
	if err := b.Build(context.Background(), s.ID); err != nil {
		t.Error(err)
	}

	// partial: a source asked alone
	p := queue(t, b, &BackupState{Sources: []string{"site"}})
	if err := b.Build(context.Background(), p.ID); err != nil {
		t.Fatal(err)
	}
	if p = reload(t, b, p.ID); !p.Partial || p.Status != StatusReady {
		t.Errorf("partial %+v", p)
	}
}

func TestBuildFailsAndCancels(t *testing.T) {
	b, _ := setup(t, nil)
	b.core.Source("bad", func(ctx context.Context, st *backup.State) error { return errors.New("broken") })
	f := queue(t, b, &BackupState{})
	if err := b.Build(context.Background(), f.ID); err == nil {
		t.Error("no error")
	}
	if f = reload(t, b, f.ID); f.Status != StatusFailed || !strings.Contains(f.Error, "broken") {
		t.Errorf("failed %+v", f)
	}

	// canceled while queued: not made
	q := queue(t, b, &BackupState{})
	if err := b.Cancel(q.ID); err != nil {
		t.Fatal(err)
	}
	b.Build(context.Background(), q.ID)
	if q = reload(t, b, q.ID); q.Status != StatusCanceled {
		t.Errorf("queued %+v", q)
	}
	if err := b.Cancel(q.ID); err == nil {
		t.Error("canceled twice")
	}
}

func TestCancelBuilding(t *testing.T) {
	block := make(chan struct{})
	b, _ := setup(t, block)
	s := queue(t, b, &BackupState{})
	done := make(chan error, 1)
	go func() { done <- b.Build(context.Background(), s.ID) }()
	for i := 0; ; i++ {
		if reload(t, b, s.ID).Status == StatusBuilding {
			break
		}
		if i > 200 {
			t.Fatal("not building")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := b.Remove(s.ID); !errors.Is(err, ErrInProgress) {
		t.Errorf("removed while building: %v", err)
	}
	if err := b.Cancel(s.ID); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Errorf("canceled: %v", err)
	}
	if s = reload(t, b, s.ID); s.Status != StatusCanceled {
		t.Errorf("state %+v", s)
	}
	if left, _ := filepath.Glob(filepath.Join(b.core.Dir(), "*")); len(left) != 0 {
		t.Errorf("left %v", left)
	}
}

func TestRetainAndRecover(t *testing.T) {
	b, _ := setup(t, nil)
	b.Keep(1)
	var ids []string
	for i, keep := range []bool{true, false, false} {
		s := queue(t, b, &BackupState{Keep: keep})
		b.db.Model(s).Update("created_at", time.Now().Add(time.Duration(i)*time.Minute))
		if err := b.Build(context.Background(), s.ID); err != nil {
			t.Fatal(err)
		}
		ids = append(ids, s.ID)
	}
	// kept, and the newest: the middle one removed
	var n int64
	b.db.Model(&BackupState{}).Count(&n)
	if n != 2 {
		t.Errorf("%d states", n)
	}
	if err := b.db.First(&BackupState{}, "id = ?", ids[1]).Error; err == nil {
		t.Error("the old one kept")
	}
	if _, err := os.Stat(filepath.Join(b.core.Dir(), ids[1])); err == nil {
		t.Error("its folder kept")
	}

	// a state being made when the application stopped: failed
	s := queue(t, b, &BackupState{})
	b.db.Model(s).Update("status", StatusBuilding)
	os.MkdirAll(filepath.Join(b.core.Dir(), ".building-x-1"), 0o700)
	if err := b.recover(); err != nil {
		t.Fatal(err)
	}
	if s = reload(t, b, s.ID); s.Status != StatusFailed || s.Error == "" {
		t.Errorf("interrupted %+v", s)
	}
	if left, _ := filepath.Glob(filepath.Join(b.core.Dir(), ".building-*")); len(left) != 0 {
		t.Errorf("left %v", left)
	}

	// ready, its folder gone: failed
	os.RemoveAll(filepath.Join(b.core.Dir(), ids[2]))
	if err := b.recover(); err != nil {
		t.Fatal(err)
	}
	if g := reload(t, b, ids[2]); g.Status != StatusFailed || g.Error != Messages_en_US.Gone {
		t.Errorf("gone %+v", g)
	}
	if k := reload(t, b, ids[0]); k.Status != StatusReady {
		t.Errorf("kept %+v", k)
	}
}

func TestToken(t *testing.T) {
	sec := []byte("k")
	now := time.Now()
	tok := signToken(sec, sessionToken{Session: uuid.New(), State: "s1", Expires: now.Add(time.Minute)})
	if got, err := parseToken(sec, tok, now); err != nil || got.State != "s1" {
		t.Fatalf("%+v %v", got, err)
	}
	if _, err := parseToken([]byte("other"), tok, now); err == nil {
		t.Error("another secret")
	}
	if _, err := parseToken(sec, tok, now.Add(2*time.Minute)); err == nil {
		t.Error("expired")
	}
	if _, err := parseToken(sec, tok[:len(tok)-2]+"xx", now); err == nil {
		t.Error("tampered")
	}
	if _, err := parseToken(sec, "x", now); err == nil {
		t.Error("garbage")
	}
}

// The client's routes: the newest state, a session by a key, the manifest
// and the files by its token — Range, If-Match —, the usage counted.
func TestClient(t *testing.T) {
	b, _ := setup(t, nil)
	b.p = nil
	ready := queue(t, b, &BackupState{})
	if err := b.Build(context.Background(), ready.ID); err != nil {
		t.Fatal(err)
	}
	queued := queue(t, b, &BackupState{})
	base := "/admin/!backup"
	b.clientBase = base
	srv := httptest.NewServer(b.ClientHandler(nil))
	defer srv.Close()

	do := func(method, path, auth string, hdr ...string) *http.Response {
		t.Helper()
		req, _ := http.NewRequest(method, srv.URL+path, nil)
		if auth != "" {
			req.Header.Set("Authorization", "Bearer "+auth)
		}
		for i := 0; i+1 < len(hdr); i += 2 {
			req.Header.Set(hdr[i], hdr[i+1])
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		return res
	}
	body := func(res *http.Response) string {
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		return string(b)
	}

	if res := do("GET", base+"/latest", ""); res.StatusCode != 401 {
		t.Errorf("latest with no key: %d", res.StatusCode)
	}
	if res := do("GET", base+"/latest", "rk_bad"); res.StatusCode != 401 {
		t.Errorf("latest with a bad key: %d", res.StatusCode)
	}
	if got := body(do("GET", base+"/latest", "rk_1")); !strings.Contains(got, ready.ID) {
		t.Errorf("latest %s", got)
	}
	if res := do("POST", base+"/"+queued.ID+"/session", "rk_1"); res.StatusCode != http.StatusConflict {
		t.Errorf("session of a queued: %d", res.StatusCode)
	}
	b.allowDownload = func(*http.Request) bool { return false }
	if res := do("POST", base+"/"+ready.ID+"/session", "rk_1"); res.StatusCode != http.StatusForbidden {
		t.Errorf("session not allowed: %d", res.StatusCode)
	}
	b.allowDownload = func(*http.Request) bool { return true }

	var sess struct{ Token, Manifest, Files string }
	if err := json.Unmarshal([]byte(body(do("POST", base+"/"+ready.ID+"/session", "rk_1"))), &sess); err != nil || sess.Token == "" {
		t.Fatalf("session %+v %v", sess, err)
	}
	if b.keys.(*fakeKeys).served != 1 {
		t.Error("the key's use not told")
	}
	// the manifest, by the token; not by a key
	if res := do("GET", sess.Manifest, "rk_1"); res.StatusCode != 401 {
		t.Errorf("manifest by a key: %d", res.StatusCode)
	}
	res := do("GET", sess.Manifest, sess.Token)
	m, err := backup.ReadManifest(res.Body)
	res.Body.Close()
	if err != nil || m.Header.ID != ready.ID || res.Header.Get("ETag") == "" {
		t.Fatalf("manifest %v %v", m, err)
	}
	e, _ := m.Find("public/a.txt")

	// a file: whole, a range, If-Match
	if got := body(do("GET", sess.Files+"public/a.txt", sess.Token)); got != "hello world" {
		t.Errorf("file %q", got)
	}
	res = do("GET", sess.Files+"public/a.txt", sess.Token, "Range", "bytes=6-", "If-Range", `"`+e.Sum.String()+`"`)
	if got := body(res); res.StatusCode != http.StatusPartialContent || got != "world" {
		t.Errorf("range %d %q", res.StatusCode, got)
	}
	if res := do("GET", sess.Files+"public/a.txt", sess.Token, "If-Match", `"other"`); res.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("if-match: %d", res.StatusCode)
	}
	for _, p := range []string{"public/none", "public", "../x"} {
		if res := do("GET", sess.Files+p, sess.Token); res.StatusCode != http.StatusNotFound {
			t.Errorf("%s: %d", p, res.StatusCode)
		}
	}
	// a token of another state
	other := signToken(b.secret, sessionToken{Session: uuid.New(), State: "zzz", Expires: time.Now().Add(time.Hour)})
	if res := do("GET", sess.Files+"public/a.txt", other); res.StatusCode != 401 {
		t.Errorf("another state's token: %d", res.StatusCode)
	}
	// renewed by its token
	if res := do("POST", base+"/"+ready.ID+"/session", sess.Token); res.StatusCode != 200 {
		t.Errorf("renew: %d", res.StatusCode)
	}
	// the usage: written by flush
	if err := b.usage.flush(b.db); err != nil {
		t.Fatal(err)
	}
	var s BackupSession
	b.db.First(&s, "state_id = ?", ready.ID)
	if s.Files != 2 || s.Bytes != int64(len("hello world")+len("world")) || s.UserName != "ana" || strings.Contains(s.IP, ":") {
		t.Errorf("session %+v", s)
	}
}

// A request with no credential of the client's is the admin's (its login,
// the manifest to its session).
func TestClientFallsBackToAdmin(t *testing.T) {
	b, _ := setup(t, nil)
	b.clientBase = "/admin/!backup"
	admin := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("admin")) })
	h := b.ClientHandler(admin)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/!backup/x/manifest", nil))
	if w.Body.String() != "admin" {
		t.Errorf("no credential: %q", w.Body.String())
	}
	r := httptest.NewRequest("GET", "/admin/!backup/x/manifest", nil)
	r.Header.Set("Authorization", "Bearer nope")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("a bad token: %d", w.Code)
	}
}
