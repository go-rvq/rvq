package backup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/packages/backup"
	"github.com/google/uuid"
	"github.com/sunfmin/reflectutils"
	"gorm.io/gorm"
)

// PermDownload is the permission of downloading the states (the listing's):
// the manifest and the files — all the application's data and secrets.
const PermDownload = "!download"

// ServiceName is the backup's service (presets.Service): <prefix>/!backup —
// no page of the menu; its paths stable (a client keeps them).
const ServiceName = "backup"

// ClientPath is where the client's routes are: <prefix>/!backup.
func (b *Builder) ClientPath() string {
	if b.clientBase != "" {
		return b.clientBase
	}
	return b.p.ServicePath(ServiceName)
}

// ClientHandler serves the client's routes, out of the admin's login
// (their own: an access key, then a token):
//
//	GET  <ClientPath>/latest                 the id of the newest state ready and whole
//	POST <ClientPath>/<id>/session           an access key (or a token of it: renewed) → a token
//	GET  <ClientPath>/<id>/manifest          the manifest (token)
//	GET  <ClientPath>/<id>/files/<path>      a file (token): Range, If-Match (the ETag is its sha256)
//
// A state not ready: 409 (queued, building), 410 (failed, canceled). A
// request with no credential of the client's (Authorization) is admin's:
// the manifest by the admin's session, the login.
func (b *Builder) ClientHandler(admin http.Handler) http.Handler {
	base := b.ClientPath()
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+base+"/latest", b.serveLatest)
	mux.HandleFunc("POST "+base+"/{id}/session", b.serveSession)
	mux.HandleFunc("GET "+base+"/{id}/manifest", b.withToken(b.serveManifest))
	mux.HandleFunc("GET "+base+"/{id}/files/{path...}", b.withToken(b.serveFile))
	// a download link: the link is its access (no Authorization)
	archive := http.NewServeMux()
	archive.HandleFunc("GET "+base+"/{id}/archive/{file}", b.serveArchive)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, pattern := archive.Handler(r); pattern != "" {
			archive.ServeHTTP(w, r)
			return
		}
		if r.Header.Get("Authorization") == "" && admin != nil {
			admin.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(v)
}

// keyUser is the user of the access key of r (Bearer, or the password of
// HTTP Basic), r with it, if it may download.
func (b *Builder) keyUser(w http.ResponseWriter, r *http.Request) (*http.Request, any, bool) {
	if b.keys == nil {
		jsonError(w, http.StatusServiceUnavailable, "no access keys")
		return nil, nil, false
	}
	code, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok {
		_, code, _ = r.BasicAuth()
	}
	code = strings.TrimSpace(code)
	if code == "" || !b.keys.IsKeyCode(code) {
		w.Header().Set("WWW-Authenticate", `Bearer realm="backup"`)
		jsonError(w, http.StatusUnauthorized, "an access key is needed")
		return nil, nil, false
	}
	user, ctx, err := b.keys.AuthKey(r, code)
	if err != nil || user == nil {
		jsonError(w, http.StatusUnauthorized, "access key refused")
		return nil, nil, false
	}
	// as the login serves a key: its user, restricted as the key says
	r = r.WithContext(context.WithValue(login.WithAuth(ctx, login.AuthAccessKey), login.UserKey, user))
	if !b.mayDownload(r) {
		jsonError(w, http.StatusForbidden, "the key may not download the backup")
		return nil, nil, false
	}
	return r, user, true
}

// readyState is the state id, ready; else the answer said.
func (b *Builder) readyState(w http.ResponseWriter, id string) (*BackupState, bool) {
	var s BackupState
	if err := b.db.First(&s, "id = ?", id).Error; err != nil {
		jsonError(w, http.StatusNotFound, "no such state")
		return nil, false
	}
	switch s.Status {
	case StatusReady:
		return &s, true
	case StatusQueued, StatusBuilding:
		jsonError(w, http.StatusConflict, "the state is "+s.Status)
	default:
		jsonError(w, http.StatusGone, "the state is "+s.Status)
	}
	return nil, false
}

func (b *Builder) serveLatest(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := b.keyUser(w, r); !ok {
		return
	}
	var s BackupState
	err := b.db.Where("status = ? AND partial = ?", StatusReady, false).Order("created_at DESC").First(&s).Error
	if err != nil {
		jsonError(w, http.StatusNotFound, "no state ready")
		return
	}
	writeJSON(w, map[string]any{"id": s.ID, "created": s.CreatedAt, "files": s.Files, "bytes": s.Bytes})
}

// serveSession opens a session (an access key) or renews one (its token):
// a token of the state, valid for TokenTTL.
func (b *Builder) serveSession(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	now := time.Now()
	// renewed: the token of a session of this state
	if tok, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer "); ok {
		if t, err := parseToken(b.secret, strings.TrimSpace(tok), now); err == nil && t.State == id {
			if _, ok := b.readyState(w, id); !ok {
				return
			}
			b.answerSession(w, t.Session, id, now)
			return
		}
	}
	r, user, ok := b.keyUser(w, r)
	if !ok {
		return
	}
	if _, ok := b.readyState(w, id); !ok {
		return
	}
	uid, name := userOf(user)
	s := BackupSession{ID: uuid.New(), StateID: id, ExpiresAt: now.Add(b.tokenTTL), UserID: model.UUID(uid),
		UserName: name, IP: hostOf(login.ClientAddr(r)), UserAgent: r.UserAgent()}
	if err := b.db.Create(&s).Error; err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	b.answerSession(w, s.ID, id, now)
	// the key's history: its use
	b.keys.KeyServed(r, http.StatusOK)
}

func (b *Builder) answerSession(w http.ResponseWriter, session uuid.UUID, id string, now time.Time) {
	exp := now.Add(b.tokenTTL)
	_ = b.db.Model(&BackupSession{}).Where("id = ?", session).Update("expires_at", exp).Error
	base := b.ClientPath() + "/" + id
	writeJSON(w, map[string]any{
		"token":    signToken(b.secret, sessionToken{Session: session, State: id, Expires: exp}),
		"expires":  exp.UTC(),
		"state":    id,
		"manifest": base + "/manifest",
		"files":    base + "/files/",
	})
}

// withToken serves h to a request with a token of its state ({id}), the
// state ready.
func (b *Builder) withToken(h func(w http.ResponseWriter, r *http.Request, s *BackupState, t sessionToken)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		t, err := parseToken(b.secret, strings.TrimSpace(tok), time.Now())
		if !ok || err != nil || t.State != r.PathValue("id") {
			w.Header().Set("WWW-Authenticate", `Bearer realm="backup"`)
			jsonError(w, http.StatusUnauthorized, "a token of this state is needed (session)")
			return
		}
		s, ok := b.readyState(w, t.State)
		if !ok {
			return
		}
		h(w, r, s, t)
	}
}

func (b *Builder) serveManifest(w http.ResponseWriter, r *http.Request, s *BackupState, t sessionToken) {
	b.sendManifest(w, r, s)
}

// sendManifest serves the manifest of s: its ETag its sha256.
func (b *Builder) sendManifest(w http.ResponseWriter, r *http.Request, s *BackupState) {
	f, err := os.Open(b.core.ManifestPath(s.ID))
	if err != nil {
		jsonError(w, http.StatusNotFound, "no manifest")
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/gzip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="manifest-%s.gz"`, s.ID))
	w.Header().Set("ETag", `"`+s.ManifestSum+`"`)
	http.ServeContent(w, r, "", fi.ModTime(), f)
}

func (b *Builder) serveFile(w http.ResponseWriter, r *http.Request, s *BackupState, t sessionToken) {
	m, err := b.manifest(s.ID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	e, ok := m.Find(r.PathValue("path"))
	if !ok || e.Type != backup.TypeFile {
		jsonError(w, http.StatusNotFound, "not a file of the state")
		return
	}
	f, err := b.core.OpenFile(s.ID, e)
	switch {
	case errors.Is(err, backup.ErrChanged):
		jsonError(w, http.StatusPreconditionFailed, "the file changed since the state")
		return
	case errors.Is(err, fs.ErrNotExist):
		jsonError(w, http.StatusNotFound, "the file is gone")
		return
	case err != nil:
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("ETag", `"`+e.Sum.String()+`"`)
	cw := &countingWriter{ResponseWriter: w}
	http.ServeContent(cw, r, "", e.MTime, f)
	if cw.n > 0 {
		b.usage.add(t.Session, cw.n)
	}
}

// manifest is the manifest of the state id, read once.
func (b *Builder) manifest(id string) (*backup.Manifest, error) {
	if m, ok := b.manifests.Load(id); ok {
		return m.(*backup.Manifest), nil
	}
	m, err := b.core.Manifest(id)
	if err != nil {
		return nil, err
	}
	b.manifests.Store(id, m)
	return m, nil
}

type countingWriter struct {
	http.ResponseWriter
	n int64
}

func (c *countingWriter) Write(p []byte) (int, error) {
	n, err := c.ResponseWriter.Write(p)
	c.n += int64(n)
	return n, err
}

// usage counts what the sessions were served, written now and then
// (flush): never a write a request.
type usage struct {
	mu sync.Mutex
	m  map[uuid.UUID]*served
}

type served struct {
	files int
	bytes int64
}

func (u *usage) add(session uuid.UUID, n int64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.m == nil {
		u.m = map[uuid.UUID]*served{}
	}
	s := u.m[session]
	if s == nil {
		s = &served{}
		u.m[session] = s
	}
	s.files++
	s.bytes += n
}

// flush writes what was counted to the sessions.
func (u *usage) flush(db *gorm.DB) error {
	u.mu.Lock()
	m := u.m
	u.m = nil
	u.mu.Unlock()
	var errs []error
	for id, s := range m {
		errs = append(errs, db.Model(&BackupSession{}).Where("id = ?", id).Updates(map[string]any{
			"files": gorm.Expr("files + ?", s.files),
			"bytes": gorm.Expr("bytes + ?", s.bytes),
		}).Error)
	}
	return errors.Join(errs...)
}

// userOf are the id and the name of a user of the admin.
func userOf(u any) (uuid.UUID, string) {
	if u == nil {
		return uuid.Nil, ""
	}
	var id uuid.UUID
	if v, err := reflectutils.Get(u, "ID"); err == nil {
		switch t := v.(type) {
		case uuid.UUID:
			id = t
		case model.UUID:
			id = t.UUID()
		case string:
			id, _ = uuid.Parse(t)
		}
	}
	name := fmt.Sprint(u)
	if v, ok := u.(interface{ GetAccountName() string }); ok && v.GetAccountName() != "" {
		name = v.GetAccountName()
	}
	return id, name
}

// clientBase is the scheme and the host r came by ("https://site.com"): the
// address the client is told.
func clientBase(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

// mayDownload says r's user may download the states (PermDownload).
func (b *Builder) mayDownload(r *http.Request) bool {
	if b.allowDownload != nil {
		return b.allowDownload(r)
	}
	return !b.mb.Permissioner().ReqListDo(r, PermDownload).Denied()
}

// hostOf is addr without its port.
func hostOf(addr string) string {
	if h, _, err := net.SplitHostPort(addr); err == nil {
		return h
	}
	return addr
}
