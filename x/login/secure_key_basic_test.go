package login

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type skUser struct {
	ID        uint
	DeletedAt gorm.DeletedAt
	UserPass
}

var skSeq int64

// The secure key as the password of HTTP Basic, where a handler takes it
// (SecureKeyBasicAuthMiddleware): the administrator — the initial account
// —, whatever user it names, only when the application turned it on, over
// HTTPS (or from the machine itself), the account not locked; never counted
// as a wrong password of an account. Anything else goes by the password.
func TestSecureKeyBasicAuth(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:sk_%d?mode=memory&cache=shared", atomic.AddInt64(&skSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&skUser{}); err != nil {
		t.Fatal(err)
	}
	admin := &skUser{UserPass: UserPass{Account: "admin", Password: "right"}}
	admin.EncryptPassword()
	other := &skUser{UserPass: UserPass{Account: "ana", Password: "hers"}}
	other.EncryptPassword()
	db.Create(admin)
	db.Create(other)

	b := newSecureKeyBuilder(t)
	b.DB(db).UserModel(&skUser{}).InitialUserAccount("admin")
	key, err := b.RenewSecureKey()
	if err != nil {
		t.Fatal(err)
	}
	on := false
	b.SecureKeyBasicAuth(func() bool { return on })

	var seen, auth string
	h := b.SecureKeyBasicAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u, ok := r.Context().Value(UserKey).(*skUser); ok {
			seen = u.Account
		} else {
			seen = fmt.Sprintf("%T", r.Context().Value(UserKey))
		}
		auth = AuthOf(r.Context())
	}))
	do := func(user, pass, remote, proto string) int {
		seen, auth = "", ""
		r := httptest.NewRequest("GET", "/admin/site-files.git/info/refs", nil)
		r.RemoteAddr = remote
		if proto != "" {
			r.Header.Set("X-Forwarded-Proto", proto)
		}
		r.SetBasicAuth(user, pass)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Code
	}
	const local, outside, proxy = "127.0.0.1:5000", "203.0.113.9:5000", "172.18.0.1:5000"

	if code := do("anyone", key, local, ""); code != http.StatusUnauthorized || seen != "" {
		t.Errorf("off: %d %q", code, seen)
	}
	on = true
	if code := do("anyone", key, local, ""); code != http.StatusOK || seen != "admin" || auth != AuthSecureKey {
		t.Errorf("on, from the machine: %d %q %q", code, seen, auth)
	}
	if code := do("ana", key, local, "https"); code != http.StatusOK || seen != "admin" {
		t.Errorf("another user named: the administrator still: %d %q", code, seen)
	}
	if code := do("admin", key, outside, ""); code != http.StatusUnauthorized {
		t.Errorf("plain HTTP from outside: %d", code)
	}
	if code := do("admin", key, outside, "https"); code != http.StatusUnauthorized {
		t.Errorf("X-Forwarded-Proto of no trusted proxy: %d", code)
	}
	if code := do("admin", key, local, "http"); code != http.StatusUnauthorized {
		t.Errorf("a local proxy saying plain HTTP: %d", code)
	}
	if code := do("admin", key, proxy, "https"); code != http.StatusUnauthorized {
		t.Errorf("a proxy not trusted: %d", code)
	}
	TrustedProxies = "172.16.0.0/12"
	t.Cleanup(func() { TrustedProxies = "" })
	if code := do("admin", key, proxy, "https"); code != http.StatusOK || seen != "admin" {
		t.Errorf("a trusted proxy saying HTTPS: %d %q", code, seen)
	}

	// a password, not the key: by the password, as always
	if code := do("ana", "hers", outside, ""); code != http.StatusOK || seen != "ana" || auth != AuthPassword {
		t.Errorf("a password: %d %q %q", code, seen, auth)
	}
	// the key never counts as a wrong password
	var a skUser
	db.First(&a, "account = ?", "admin")
	if a.LoginRetryCount != 0 {
		t.Errorf("the key counted as a wrong password: %d", a.LoginRetryCount)
	}

	// the administrator locked: the key refused
	if err := db.Model(&skUser{}).Where("account = ?", "admin").
		Updates(map[string]any{"locked": true, "locked_at": time.Now()}).Error; err != nil {
		t.Fatal(err)
	}
	if code := do("admin", key, local, ""); code != http.StatusUnauthorized {
		t.Errorf("the administrator locked: %d", code)
	}
	// a wrong key: by the password, refused
	if code := do("admin", key+"x", local, ""); code != http.StatusUnauthorized || strings.Contains(seen, "admin") {
		t.Errorf("a wrong key: %d", code)
	}
}

// A wrong password of HTTP Basic that locks the account tells the
// application (BasicAuthLocked): the request it came by.
func TestBasicAuthLocked(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:sk_%d?mode=memory&cache=shared", atomic.AddInt64(&skSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&skUser{}); err != nil {
		t.Fatal(err)
	}
	ana := &skUser{UserPass: UserPass{Account: "ana", Password: "hers"}}
	ana.EncryptPassword()
	db.Create(ana)

	b := newSecureKeyBuilder(t)
	b.DB(db).UserModel(&skUser{}).MaxRetryCount(2)
	var locked []string
	b.BasicAuthLocked(func(r *http.Request, user any) {
		locked = append(locked, user.(*skUser).Account+"@"+r.RemoteAddr)
	})
	h := b.SecureKeyBasicAuthMiddleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	for i := 0; i < 2; i++ {
		r := httptest.NewRequest("GET", "/admin/site-files.git/info/refs", nil)
		r.RemoteAddr = "203.0.113.9:5000"
		r.SetBasicAuth("ana", "wrong")
		h.ServeHTTP(httptest.NewRecorder(), r)
	}
	if strings.Join(locked, ",") != "ana@203.0.113.9:5000" {
		t.Errorf("told of the lock: %v", locked)
	}
}
