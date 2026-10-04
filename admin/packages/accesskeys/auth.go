package accesskeys

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type keyCtxKey struct{}

// KeyOf is the access key the request of ctx is authenticated by; nil when
// it is not (a session's, a login's).
func KeyOf(ctx context.Context) *AccessKey {
	k, _ := ctx.Value(keyCtxKey{}).(*AccessKey)
	return k
}

// The refusals of a key.
var (
	ErrKeyNotFound   = errors.New("access key not found")
	ErrKeyUnusable   = errors.New("access key disabled or expired")
	ErrKeyNoUser     = errors.New("the user of the access key is gone")
	ErrTooManyErrors = errors.New("too many wrong access keys: try again later")
)

// Authenticator authenticates requests by access keys (login.KeyAuthenticator):
// the key found by its prefix, its code checked by its hash, usable, its user
// there; the request restricted to what the key allows; each request kept in
// its history.
type Authenticator struct {
	DB *gorm.DB
	// FindUser is the user of id (nil, error when there is none).
	FindUser func(ctx context.Context, id uuid.UUID) (any, error)
	// Locate is where ip is ("" when it is not known).
	Locate func(ip string) string
	// KindOf is what r reaches ("git", "webdav"); "api" when nil or "".
	KindOf func(r *http.Request) string
	// OnUse is told each request by a key, kept: the user's access log.
	OnUse func(r *http.Request, k *AccessKey, use *AccessKeyUse)
	// Now is the time (time.Now when nil).
	Now func() time.Time

	mu       sync.Mutex
	failures map[string][]time.Time
}

var _ login.KeyAuthenticator = (*Authenticator)(nil)

// maxFailures in failureWindow from an address: its keys refused meanwhile.
const (
	maxFailures   = 10
	failureWindow = 10 * time.Minute
)

func (a *Authenticator) now() time.Time {
	if a.Now != nil {
		return a.Now()
	}
	return time.Now()
}

// IsKeyCode says whether code is of the format of the keys.
func (a *Authenticator) IsKeyCode(code string) bool {
	_, err := ParseCode(code)
	return err == nil
}

// AuthKey is the user of the key of code, and the context of the request:
// restricted to the key's permissions, the key told (KeyOf).
func (a *Authenticator) AuthKey(r *http.Request, code string) (any, context.Context, error) {
	ip := ClientIP(r)
	if a.blocked(ip) {
		return nil, nil, ErrTooManyErrors
	}
	prefix, err := ParseCode(code)
	if err != nil {
		return nil, nil, err
	}
	var k AccessKey
	if err := a.DB.WithContext(r.Context()).Where("prefix = ?", prefix).First(&k).Error; err != nil {
		a.fail(ip)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrKeyNotFound
		}
		return nil, nil, err
	}
	if !k.Matches(code) {
		a.fail(ip)
		return nil, nil, ErrKeyNotFound
	}
	if !k.Usable(a.now()) {
		return nil, nil, ErrKeyUnusable
	}
	user, err := a.FindUser(r.Context(), k.UserID)
	if err != nil || user == nil {
		return nil, nil, ErrKeyNoUser
	}
	ctx := perm.WithRestriction(r.Context(), k.Subject())
	return user, context.WithValue(ctx, keyCtxKey{}, &k), nil
}

// KeyServed keeps the request by the key in its history, and its last use.
func (a *Authenticator) KeyServed(r *http.Request, status int) {
	k := KeyOf(r.Context())
	if k == nil {
		return
	}
	ip := ClientIP(r)
	kind := "api"
	if a.KindOf != nil {
		if s := a.KindOf(r); s != "" {
			kind = s
		}
	}
	use := &AccessKeyUse{KeyID: k.ID, UserID: k.UserID, IP: ip, UserAgent: r.UserAgent(), Kind: kind,
		Method: r.Method, Path: r.URL.Path, Status: status}
	if a.Locate != nil {
		use.Place = a.Locate(ip)
	}
	db := a.DB.WithContext(context.WithoutCancel(r.Context()))
	_ = db.Create(use).Error
	now := a.now()
	_ = db.Model(&AccessKey{}).Where("id = ?", k.ID).
		UpdateColumns(map[string]any{"last_used_at": now, "last_used_ip": ip}).Error
	if a.OnUse != nil {
		a.OnUse(r, k, use)
	}
}

func (a *Authenticator) fail(ip string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.failures == nil {
		a.failures = map[string][]time.Time{}
	}
	a.failures[ip] = append(a.recent(ip), a.now())
}

func (a *Authenticator) blocked(ip string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.recent(ip)) >= maxFailures
}

// recent are the failures of ip in the window (the lock held).
func (a *Authenticator) recent(ip string) []time.Time {
	var out []time.Time
	for _, t := range a.failures[ip] {
		if a.now().Sub(t) < failureWindow {
			out = append(out, t)
		}
	}
	return out
}

// ClientIP is the address of the client of r: the first of X-Forwarded-For
// or X-Real-IP (behind a proxy), else the connection's.
func ClientIP(r *http.Request) string {
	if f := r.Header.Get("X-Forwarded-For"); f != "" {
		return strings.TrimSpace(strings.Split(f, ",")[0])
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}
