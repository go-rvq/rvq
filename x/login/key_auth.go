package login

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// KeyAuthenticator authenticates a request by the code of an access key — an
// automation's: git, WebDAV, a script —, given as the password of HTTP Basic
// or as a bearer token. Each request carries it: no session, no cookie.
type KeyAuthenticator interface {
	// IsKeyCode says whether code is of an access key (its format): when not,
	// the other authentications go on.
	IsKeyCode(code string) bool
	// AuthKey is the user of the key of code — enabled, not expired — and the
	// context of its request: restricted to what the key allows
	// (perm.WithRestriction), and the key told for its log. An error refuses
	// the request.
	AuthKey(r *http.Request, code string) (user any, ctx context.Context, err error)
	// KeyServed is told the request served by the key (KeyOf its context),
	// its status: its history.
	KeyServed(r *http.Request, status int)
}

// KeyAuth sets the authentication by access keys (KeyAuthenticator): by HTTP
// Basic (BasichAuthMiddleware) and by bearer token (Middleware,
// BasichAuthMiddleware).
func (b *Builder) KeyAuth(a KeyAuthenticator) *Builder {
	b.keyAuth = a
	return b
}

// ErrKeyRefused is the refusal of an access key (expired, disabled, of a user
// locked): its reason the authenticator's.
var ErrKeyRefused = errors.New("access key refused")

// keyCode is the code of an access key the request carries — the bearer
// token, or the password of HTTP Basic —, "" when none.
func (b *Builder) keyCode(r *http.Request) string {
	if b.keyAuth == nil {
		return ""
	}
	if v, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer "); ok {
		if v = strings.TrimSpace(v); b.keyAuth.IsKeyCode(v) {
			return v
		}
	}
	if _, pass, ok := r.BasicAuth(); ok && b.keyAuth.IsKeyCode(pass) {
		return pass
	}
	return ""
}

// serveKey serves r by its access key's code: its user in the context — as a
// login's —, restricted as the key says; a user locked, refused. The status
// served told to the authenticator (its history).
func (b *Builder) serveKey(w http.ResponseWriter, r *http.Request, code string, next http.Handler) {
	user, ctx, err := b.keyAuth.AuthKey(r, code)
	if err == nil && user != nil {
		if up, ok := user.(UserPasser); ok && up.GetLocked() {
			err = ErrUserLocked
		}
	}
	if err != nil || user == nil {
		if err == nil {
			err = ErrKeyRefused
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="access key"`)
		http.Error(w, "ERROR: "+err.Error(), http.StatusUnauthorized)
		return
	}
	r = r.WithContext(context.WithValue(context.WithValue(WithAuth(ctx, AuthAccessKey), keyAuthenticatedKey{}, true), UserKey, user))
	sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
	next.ServeHTTP(sw, r)
	b.keyAuth.KeyServed(r, sw.status)
}

// statusWriter remembers the status written.
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Flush flushes the writer under it, when it does (git's streaming).
func (w *statusWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Unwrap is the writer under it (http.ResponseController).
func (w *statusWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }

type keyAuthenticatedKey struct{}

// IsKeyAuthenticated says whether the request of ctx was authenticated by an
// access key (KeyAuth): it has no session — nothing of one to check.
func IsKeyAuthenticated(ctx context.Context) bool {
	v, _ := ctx.Value(keyAuthenticatedKey{}).(bool)
	return v
}

// The ways a request is told whose it is (AuthOf).
const (
	// AuthSession is a login's session (its cookie).
	AuthSession = "session"
	// AuthPassword is the account and the password of HTTP Basic.
	AuthPassword = "password"
	// AuthAccessKey is the code of an access key (KeyAuth).
	AuthAccessKey = "access_key"
	// AuthSecureKey is the secure key (SecureKeyBasicAuthMiddleware).
	AuthSecureKey = "secure_key"
)

type authKey struct{}

// WithAuth is ctx with the way its request was told whose it is (AuthOf).
func WithAuth(ctx context.Context, auth string) context.Context {
	return context.WithValue(ctx, authKey{}, auth)
}

// AuthOf is the way the request of ctx was told whose it is — AuthSession,
// AuthPassword, AuthAccessKey, AuthSecureKey —, "" when it was not told.
func AuthOf(ctx context.Context) string {
	v, _ := ctx.Value(authKey{}).(string)
	return v
}

// BasicAuthLocked is told the user whose account a wrong password of HTTP
// Basic (BasichAuthMiddleware: the WebDAV, the git) locked, and its request:
// the place it came from, to record.
func (b *Builder) BasicAuthLocked(f func(r *http.Request, user any)) *Builder {
	b.basicAuthLockedHook = f
	return b
}
