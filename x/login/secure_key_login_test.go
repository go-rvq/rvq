package login

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-rvq/rvq/x/security"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// The login by the secure key: off unless the app.yml says so (security:
// secure_key_login), a page of the key (GET), and a POST of the right key —
// the form's or the header's —, over HTTPS (or from the machine itself),
// logs in as the administrator: a session, home. Anything else, refused.
func TestSecureKeyLogin(t *testing.T) {
	defer security.Load(security.Get())
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:skl_%d?mode=memory&cache=shared", atomic.AddInt64(&skSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&skUser{}); err != nil {
		t.Fatal(err)
	}
	admin := &skUser{UserPass: UserPass{Account: "admin", Password: "right"}}
	admin.EncryptPassword()
	db.Create(admin)

	b := newSecureKeyBuilder(t)
	b.DB(db).UserModel(&skUser{}).InitialUserAccount("admin").TOTP(false)
	b.HomeURLFunc(func(*http.Request, any) string { return "/admin" })
	key, err := b.RenewSecureKey()
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	b.MountAPI(mux)

	do := func(method, remote, formKey, headerKey string) *httptest.ResponseRecorder {
		var body *strings.Reader
		if formKey != "" {
			body = strings.NewReader(url.Values{"key": {formKey}}.Encode())
		} else {
			body = strings.NewReader("")
		}
		r := httptest.NewRequest(method, "/auth/secure-key", body)
		r.RemoteAddr = remote
		if method == "POST" {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		if headerKey != "" {
			r.Header.Set(SecureKeyHeader, headerKey)
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w
	}
	session := func(w *httptest.ResponseRecorder) bool {
		for _, c := range w.Result().Cookies() {
			if c.Name == b.authCookieName && c.Value != "" {
				return true
			}
		}
		return false
	}
	const local, outside = "127.0.0.1:5000", "203.0.113.9:5000"

	// off: not there
	security.Load(security.Config{})
	if w := do("GET", local, "", ""); w.Code != http.StatusNotFound {
		t.Errorf("off, the page: %d", w.Code)
	}
	if w := do("POST", local, key, ""); w.Code != http.StatusNotFound || session(w) {
		t.Errorf("off, a login: %d %v", w.Code, session(w))
	}

	security.Load(security.Config{SecureKeyLogin: true})
	if w := do("GET", local, "", ""); w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `name="key"`) {
		t.Errorf("the page: %d %.200s", w.Code, w.Body.String())
	}
	// the key in the form, from the machine: the administrator, home
	if w := do("POST", local, key, ""); w.Code != http.StatusFound || w.Header().Get("Location") != "/admin" || !session(w) {
		t.Errorf("the form's key: %d %q %v", w.Code, w.Header().Get("Location"), session(w))
	}
	// the header's
	if w := do("POST", local, "", key); w.Code != http.StatusFound || !session(w) {
		t.Errorf("the header's key: %d %v", w.Code, session(w))
	}
	// a wrong key, none, plain HTTP from outside: refused, no session
	for name, w := range map[string]*httptest.ResponseRecorder{
		"wrong":   do("POST", local, key+"x", ""),
		"none":    do("POST", local, "", ""),
		"outside": do("POST", outside, key, ""),
	} {
		if w.Code != http.StatusUnauthorized || session(w) {
			t.Errorf("%s: %d %v", name, w.Code, session(w))
		}
	}
	// the administrator locked: refused
	if err := db.Model(&skUser{}).Where("account = ?", "admin").
		Updates(map[string]any{"locked": true, "locked_at": time.Now()}).Error; err != nil {
		t.Fatal(err)
	}
	if w := do("POST", local, key, ""); w.Code != http.StatusUnauthorized || session(w) {
		t.Errorf("locked: %d %v", w.Code, session(w))
	}
}

// The secure key as the password on the login page: the administrator's
// account in the form (or what finds it), the login by the key on, from the
// machine — a session; another account, the key off, from outside: as any
// wrong password.
func TestSecureKeyAsPassword(t *testing.T) {
	defer security.Load(security.Get())
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:skp_%d?mode=memory&cache=shared", atomic.AddInt64(&skSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&skUser{}); err != nil {
		t.Fatal(err)
	}
	for _, a := range []string{"admin", "ana"} {
		u := &skUser{UserPass: UserPass{Account: a, Password: "right-" + a}}
		u.EncryptPassword()
		db.Create(u)
	}
	b := newSecureKeyBuilder(t)
	b.DB(db).UserModel(&skUser{}).InitialUserAccount("admin").TOTP(false)
	b.HomeURLFunc(func(*http.Request, any) string { return "/admin" })
	// the form is filled instantly in a test
	b.ChallengeAges(time.Nanosecond, time.Minute)
	key, err := b.RenewSecureKey()
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	b.MountAPI(mux)
	// a login by the page's form (its challenge)
	login := func(remote, account, password string) bool {
		r := httptest.NewRequest("POST", "/auth/userpass/login", strings.NewReader(url.Values{"account": {account}, "password": {password}, ChallengeFormKey: {b.NewChallenge()}}.Encode()))
		r.RemoteAddr = remote
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		for _, c := range w.Result().Cookies() {
			if c.Name == b.authCookieName && c.Value != "" {
				return true
			}
		}
		return false
	}
	const local, outside = "127.0.0.1:5000", "203.0.113.9:5000"

	security.Load(security.Config{SecureKeyLogin: true})
	if !login(local, "admin", key) {
		t.Error("the key as the administrator's password: no session")
	}
	if !login(local, "admin", "right-admin") || !login(local, "ana", "right-ana") {
		t.Error("the passwords: no session")
	}
	for name, ok := range map[string]bool{
		"another account": login(local, "ana", key),
		"no account":      login(local, "", key),
		"a wrong key":     login(local, "admin", key+"x"),
		"from outside":    login(outside, "admin", key),
	} {
		if ok {
			t.Errorf("%s: a session", name)
		}
	}
	security.Load(security.Config{})
	if login(local, "admin", key) {
		t.Error("the key login off: a session")
	}
}
