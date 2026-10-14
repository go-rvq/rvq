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
