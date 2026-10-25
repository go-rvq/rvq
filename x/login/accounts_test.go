package login

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-rvq/rvq/x/account"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// a browser: its cookies, kept between the requests
type jar map[string]string

func (j jar) take(w *httptest.ResponseRecorder) {
	for _, c := range w.Result().Cookies() {
		if c.MaxAge < 0 || (!c.Expires.IsZero() && c.Expires.Before(time.Now())) || c.Value == "" {
			delete(j, c.Name)
		} else {
			j[c.Name] = c.Value
		}
	}
}

func (j jar) put(r *http.Request) {
	for k, v := range j {
		r.AddCookie(&http.Cookie{Name: k, Value: v})
	}
}

// Several accounts in a browser: the account 0 at /admin, another — logged
// in by "Entrar com outra conta" (login?add=1) — at /admin/u/N, each with its
// cookies; a login of an account already logged in is its own index; the
// logout of one keeps the others; the logout of all, none; the redirects of
// an account keep it.
func TestAccounts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:acc_%d?mode=memory&cache=shared", atomic.AddInt64(&skSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&skUser{}); err != nil {
		t.Fatal(err)
	}
	for _, a := range [][2]string{{"admin", "right"}, {"maria", "pw2"}} {
		u := &skUser{UserPass: UserPass{Account: a[0], Password: a[1]}}
		u.EncryptPassword()
		db.Create(u)
	}

	b := newSecureKeyBuilder(t)
	b.DB(db).UserModel(&skUser{}).InitialUserAccount("admin").TOTP(false)
	b.HomeURLFunc(func(*http.Request, any) string { return "/admin" })
	b.AccountPrefixes("/admin")
	key, err := b.RenewSecureKey()
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	b.MountAPI(mux)
	mux.HandleFunc(b.loginPageURL, func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "login") })
	mux.HandleFunc("/admin/", func(w http.ResponseWriter, r *http.Request) {
		u, _ := CurrentUserFromRequest(r).(*skUser)
		name := "-"
		if u != nil {
			name = u.Account
		}
		fmt.Fprintf(w, "%s %d %s", name, account.IndexOf(r), r.URL.Path)
	})
	h := b.AccountsMiddleware(b.Middleware()(mux))

	br := jar{}
	do := func(method, target string, form url.Values, header ...string) *httptest.ResponseRecorder {
		t.Helper()
		var body *strings.Reader
		if form != nil {
			body = strings.NewReader(form.Encode())
		} else {
			body = strings.NewReader("")
		}
		r := httptest.NewRequest(method, target, body)
		r.RemoteAddr = "127.0.0.1:5000"
		if form != nil {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			// past the form protection (the account and the password asked)
			r.Header.Set(SecureKeyHeader, key)
		}
		for i := 0; i+1 < len(header); i += 2 {
			r.Header.Set(header[i], header[i+1])
		}
		br.put(r)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		br.take(w)
		return w
	}
	login := func(acc, pass string) string {
		t.Helper()
		w := do("POST", b.passwordLoginURL, url.Values{"account": {acc}, "password": {pass}})
		if w.Code != http.StatusFound {
			t.Fatalf("login %s: %d %s", acc, w.Code, w.Body.String())
		}
		return w.Header().Get("Location")
	}
	body := func(target string, header ...string) string {
		t.Helper()
		w := do("GET", target, nil, header...)
		if w.Code != http.StatusOK {
			t.Fatalf("GET %s: %d %q", target, w.Code, w.Header().Get("Location"))
		}
		return w.Body.String()
	}

	// the account 0, as ever
	do("GET", b.loginPageURL, nil)
	if loc := login("admin", "right"); loc != "/admin" {
		t.Fatalf("admin: %q", loc)
	}
	if got := body("/admin/x"); got != "admin 0 /admin/x" {
		t.Fatalf("account 0: %q", got)
	}

	// another account: a free index (1), back to the page asked, in it
	do("GET", b.loginPageURL+"?add=1&redirect_to=/admin/content", nil)
	if loc := login("maria", "pw2"); loc != "/admin/u/1/content" {
		t.Fatalf("maria: %q", loc)
	}
	if _, ok := br["auth_u1"]; !ok {
		t.Fatalf("the cookie of the account 1: %v", br)
	}
	if got := body("/admin/u/1/content"); got != "maria 1 /admin/content" {
		t.Fatalf("account 1: %q", got)
	}
	if got := body("/admin/content"); got != "admin 0 /admin/content" {
		t.Fatalf("account 0 still: %q", got)
	}
	// a request of a page of the account 1 (its events): the header
	if got := body("/admin/content", account.Header, "1"); got != "maria 1 /admin/content" {
		t.Fatalf("the header: %q", got)
	}

	// what a page of the account 1 loads (no navigation): its account
	if got := body("/admin/assets/x.css", "Referer", "http://example.com/admin/u/1/content", "Sec-Fetch-Mode", "no-cors"); got != "maria 1 /admin/assets/x.css" {
		t.Fatalf("a load of the page of the account 1: %q", got)
	}
	// a link followed from it: the account the link says (0)
	if got := body("/admin/x", "Referer", "http://example.com/admin/u/1/content", "Sec-Fetch-Mode", "navigate"); got != "admin 0 /admin/x" {
		t.Fatalf("a navigation from the account 1: %q", got)
	}
	// the explicit /u/0
	if got := body("/admin/u/0/x"); got != "admin 0 /admin/x" {
		t.Fatalf("/u/0: %q", got)
	}

	// the site (out of the admin) follows the account the admin was last
	// opened in; ?rvq_account=N switches it
	mux.HandleFunc("/site", func(w http.ResponseWriter, r *http.Request) {
		u, _ := CurrentUserFromRequest(r).(*skUser)
		name := "-"
		if u != nil {
			name = u.Account
		}
		fmt.Fprintf(w, "%s %d", name, account.IndexOf(r))
	})
	body("/admin/u/1/content")
	if got := body("/site"); got != "maria 1" {
		t.Fatalf("the site after the admin of the account 1: %q", got)
	}
	body("/admin/content")
	if got := body("/site"); got != "admin 0" {
		t.Fatalf("the site after the admin of the account 0: %q", got)
	}
	if got := body(SiteAccountURL("/site", 1)); got != "maria 1" {
		t.Fatalf("the site switched to the account 1: %q", got)
	}
	if got := body("/site"); got != "maria 1" {
		t.Fatalf("the site stays in the account 1: %q", got)
	}
	// a load of the admin (not a page) does not move the site
	body("/admin/content", "Sec-Fetch-Dest", "empty")
	if got := body("/site"); got != "maria 1" {
		t.Fatalf("a load of the admin moved the site: %q", got)
	}
	body("/admin/content")

	// the list
	r := httptest.NewRequest("GET", "/admin/", nil)
	br.put(r)
	r = r.WithContext(account.WithIndex(r.Context(), 1))
	accs := b.Accounts(r)
	if len(accs) != 2 || accs[0].Index != 0 || accs[1].Index != 1 || !accs[0].LoggedIn || !accs[1].LoggedIn || !accs[1].Current || accs[0].Current {
		t.Fatalf("accounts: %+v %+v", accs[0], accs[len(accs)-1])
	}
	if accs[1].User.(*skUser).Account != "maria" {
		t.Fatalf("the user of the account 1: %+v", accs[1].User)
	}

	// the login page of another account: "another", back to a logged one; of
	// an account whose session expired (its cookie gone, the list keeps it):
	// "again", its account in the field
	loginOf := func(target string, n int, drop ...string) *AccountLogin {
		r := httptest.NewRequest("GET", target, nil)
		for k, v := range br {
			if !slices.Contains(drop, k) {
				r.AddCookie(&http.Cookie{Name: k, Value: v})
			}
		}
		return b.AccountLoginOf(r.WithContext(account.WithIndex(r.Context(), n)))
	}
	if al := loginOf(b.loginPageURL+"?add=1", 2); al == nil || !al.Adding || al.Again || al.BackURL != "/admin/" || al.BackName != "admin" {
		t.Fatalf("the login of another account: %+v", al)
	}
	if al := loginOf(b.loginPageURL+"?account=1", 1, "auth_u1"); al == nil || !al.Again || al.Account != "maria" {
		t.Fatalf("the login of an expired account: %+v", al)
	}
	if al := loginOf(b.loginPageURL, 0, "auth", "auth_u1", accountsCookieName); al != nil {
		t.Fatalf("the plain login: %+v", al)
	}

	// an account already logged in: its own index, not a new one
	do("GET", b.loginPageURL+"?add=1", nil)
	if loc := login("admin", "right"); loc != "/admin" {
		t.Fatalf("admin again: %q", loc)
	}
	if _, ok := br["auth_u2"]; ok {
		t.Fatal("admin got a second index")
	}

	// an index with no account: its login
	if w := do("GET", "/admin/u/3/x", nil); w.Code != http.StatusFound || w.Header().Get("Location") != b.loginPageURL+"?account=3" {
		t.Fatalf("no account 3: %d %q", w.Code, w.Header().Get("Location"))
	}
	// the redirects of the account 1 keep it, "//host/…" too
	mux.HandleFunc("/admin/slash/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "//"+r.Host+"/admin/slash", http.StatusMovedPermanently)
	})
	if w := do("GET", "/admin/u/1/slash/", nil); w.Header().Get("Location") != "//example.com/admin/u/1/slash" {
		t.Fatalf("a redirect of the account 1: %q", w.Header().Get("Location"))
	}
	// out of the accounts a browser holds
	if w := do("GET", "/admin/u/9/x", nil); w.Code != http.StatusNotFound {
		t.Fatalf("index 9: %d", w.Code)
	}

	// the logout of the account 1: the 0 stays, the browser goes to it
	if w := do("GET", b.AccountLogoutURL(1, false), nil); w.Code != http.StatusFound || w.Header().Get("Location") != "/admin/" {
		t.Fatalf("logout 1: %d %q", w.Code, w.Header().Get("Location"))
	}
	if _, ok := br["auth_u1"]; ok {
		t.Fatal("the account 1 still has its cookie")
	}
	if got := body("/admin/x"); got != "admin 0 /admin/x" {
		t.Fatalf("account 0 after the logout of 1: %q", got)
	}
	if accs := b.Accounts(func() *http.Request { r := httptest.NewRequest("GET", "/admin/", nil); br.put(r); return r }()); len(accs) != 1 {
		t.Fatalf("accounts after the logout of 1: %d", len(accs))
	}

	// the logout of all
	do("GET", b.loginPageURL+"?add=1", nil)
	login("maria", "pw2")
	if w := do("GET", b.AccountLogoutURL(0, true), nil); w.Code != http.StatusFound || w.Header().Get("Location") != b.loginPageURL {
		t.Fatalf("logout all: %d %q", w.Code, w.Header().Get("Location"))
	}
	for _, c := range []string{"auth", "auth_u1", accountsCookieName} {
		if _, ok := br[c]; ok {
			t.Fatalf("after the logout of all, %s: %v", c, br)
		}
	}
}

// One account (MaxAccounts(1)): the paths and the cookies as ever.
func TestAccountsOff(t *testing.T) {
	b := newSecureKeyBuilder(t)
	b.MaxAccounts(1).AccountPrefixes("/admin")
	seen := ""
	h := b.AccountsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = fmt.Sprintf("%d %s", account.IndexOf(r), r.URL.Path)
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/admin/u/1/x", nil))
	if seen != "0 /admin/u/1/x" {
		t.Fatalf("off: %q", seen)
	}
}
