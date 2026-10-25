package login

import (
	"bufio"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-rvq/rvq/x/account"
)

// Several accounts in a browser, as Gmail's: the account N has the cookies
// of its own — auth / rvq_auth_secure for the account 0 (the names of a
// single account, as ever), auth_uN / rvq_auth_secure_uN for another —, its
// session, its second factor; the request says which one (x/account): the
// prefix /u/N of its path under an account prefix (AccountPrefixes), its
// header X-Rvq-Account, or — the pages of the login — the account being
// logged in. The cookie rvq_accounts lists, signed, the accounts the browser
// logged in (index and user): the menu of the user shows them, those whose
// session expired too.

// DefaultMaxAccounts is how many accounts a browser holds by default.
const DefaultMaxAccounts = 5

const (
	accountsCookieName     = "rvq_accounts"
	loginAccountCookieName = "rvq_login_account"
	accountsCookieMaxAge   = 60 * 60 * 24 * 30
	loginAccountMaxAge     = 60 * 15
)

// Account is an account the browser logged in.
type Account struct {
	Index  int
	UserID string
	// User is the user of the account (nil when it is gone).
	User any
	// LoggedIn: its session is valid; false — expired —: its login asks it
	// again.
	LoggedIn bool
	// Current: the account of the request.
	Current bool
}

// MaxAccounts is how many accounts a browser holds (DefaultMaxAccounts);
// 1 turns the feature off: one account, the cookies as ever.
func (b *Builder) MaxAccounts(n int) *Builder {
	if n < 1 {
		n = 1
	}
	b.maxAccounts = n
	return b
}

func (b *Builder) GetMaxAccounts() int {
	if b.maxAccounts == 0 {
		return DefaultMaxAccounts
	}
	return b.maxAccounts
}

// AccountPrefixes are the paths whose /u/N says the account (the admin's,
// "/admin"); the first one is where the home of an account is.
func (b *Builder) AccountPrefixes(p ...string) *Builder {
	b.accountPrefixes = p
	return b
}

func (b *Builder) GetAccountPrefixes() []string { return b.accountPrefixes }

func (b *Builder) authCookieOf(r *http.Request) string {
	return account.CookieName(b.authCookieName, account.IndexOf(r))
}

func (b *Builder) authSecureCookieOf(r *http.Request) string {
	return account.CookieName(b.authSecureCookieName, account.IndexOf(r))
}

// withAccount is r of the account n.
func withAccount(r *http.Request, n int) *http.Request {
	if st, _ := r.Context().Value(accountStateKey{}).(*accountState); st != nil {
		st.n = n
	}
	return r.WithContext(account.WithIndex(r.Context(), n))
}

type accountStateKey struct{}

// accountState is the account the response is for: a login may change it
// (the user already logged in another index).
type accountState struct{ n int }

// isLoginFlow: the pages of a login — where the account comes from the login
// being done, not from a path.
func (b *Builder) isLoginFlow(p string) bool {
	switch p {
	case b.loginPageURL, b.passwordLoginURL, b.secureKeyLoginURL,
		b.totpSetupPageURL, b.totpValidatePageURL, b.validateTOTPURL,
		b.oauthBeginURL, b.oauthCallbackURL, b.oauthCallbackCompleteURL:
		return true
	}
	return false
}

// AccountsMiddleware says the account of the request (x/account), before the
// login's Middleware: the /u/N of a path under an account prefix (taken out
// of it: the routes see the path without it), its header X-Rvq-Account, the
// query account=N (the login, the logout), or the login being done
// (rvq_login_account). The login page with add=1 logs in another account: a
// free index. The redirects of an account N keep its /u/N.
func (b *Builder) AccountsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		max := b.GetMaxAccounts()
		if max <= 1 {
			next.ServeHTTP(w, r)
			return
		}
		var (
			n      int
			byPath bool
			p      = strings.TrimRight(r.URL.Path, "/")
		)
		for _, prefix := range b.accountPrefixes {
			if stripped, i, ok := account.Split(prefix, r.URL.Path); ok {
				if i >= max {
					http.NotFound(w, r)
					return
				}
				u := *r.URL
				u.Path, u.RawPath = stripped, ""
				r2 := r.Clone(r.Context())
				r2.URL = &u
				r, n, byPath = r2, i, true
				break
			}
		}
		if !byPath {
			q := r.URL.Query()
			login := b.isLoginFlow(p)
			switch {
			case (p == b.loginPageURL || p == b.secureKeyLoginURL) && q.Get("add") == "1":
				n = b.freeAccount(r)
				b.setLoginAccount(w, n)
				b.setAddAccountContinue(w, q.Get("redirect_to"), n)
			case q.Get("account") != "" && (login || p == b.logoutURL):
				n = parseAccount(q.Get("account"), max)
				if login {
					b.setLoginAccount(w, n)
				}
			case r.Header.Get(account.Header) != "":
				n = parseAccount(r.Header.Get(account.Header), max)
				if login {
					b.setLoginAccount(w, n)
				}
			case (p == b.loginPageURL || p == b.secureKeyLoginURL) && r.Method == http.MethodGet:
				// the plain login page: the account 0
				if _, err := r.Cookie(loginAccountCookieName); err == nil {
					b.clearLoginAccount(w)
				}
			case login:
				// the steps of a login: the account it began with
				if c, err := r.Cookie(loginAccountCookieName); err == nil {
					n = parseAccount(c.Value, max)
				}
			case p != b.logoutURL:
				// what a page of the account N loads — its styles, scripts,
				// pictures, frames' forms —: the account of the page
				var ok bool
				if n, ok = b.refererAccount(r, max); !ok && !b.underAccountPrefix(r.URL.Path) {
					// the site (out of the admin): the account the admin was
					// last opened in (rvq_site_account), or the one it
					// switches to (?rvq_account=N)
					n = b.siteAccount(w, r, max)
				}
			}
		}
		if (byPath || b.underAccountPrefix(r.URL.Path)) && isDocument(r) {
			// a page of the admin opened: the site follows its account
			b.setSiteAccount(w, r, n)
		}
		st := &accountState{n: n}
		r = r.WithContext(account.WithIndex(context.WithValue(r.Context(), accountStateKey{}, st), n))
		next.ServeHTTP(&accountWriter{ResponseWriter: w, b: b, st: st, host: r.Host}, r)
	})
}

// refererAccount is the account of the page of the same site r was asked
// from (its Referer, /admin/u/N/…), for what is not a navigation — the
// browser's own loads, not a link followed —; 0 otherwise.
func (b *Builder) refererAccount(r *http.Request, max int) (int, bool) {
	mode := r.Header.Get("Sec-Fetch-Mode")
	if mode == "" || mode == "navigate" {
		return 0, false
	}
	ref := r.Header.Get("Referer")
	if ref == "" {
		return 0, false
	}
	u, err := url.Parse(ref)
	if err != nil || (u.Host != "" && u.Host != r.Host) {
		return 0, false
	}
	for _, prefix := range b.accountPrefixes {
		if _, i, ok := account.Split(prefix, u.Path); ok && i < max {
			return i, true
		}
	}
	return 0, false
}

// SiteAccountQuery is the query that switches the account of the site
// (?rvq_account=N): its pages, its bar of the user.
const SiteAccountQuery = "rvq_account"

const siteAccountCookieName = "rvq_site_account"

func (b *Builder) underAccountPrefix(p string) bool {
	for _, prefix := range b.accountPrefixes {
		prefix = strings.TrimSuffix(prefix, "/")
		if prefix != "" && (p == prefix || strings.HasPrefix(p, prefix+"/")) {
			return true
		}
	}
	return false
}

// isDocument: r opens a page (not a load of one, nor its events).
func isDocument(r *http.Request) bool {
	if r.Method != http.MethodGet || r.Header.Get(account.Header) != "" || r.URL.Query().Has("__execute_event__") {
		return false
	}
	if d := r.Header.Get("Sec-Fetch-Dest"); d != "" && d != "document" {
		return false
	}
	return true
}

// siteAccount is the account of a page of the site.
func (b *Builder) siteAccount(w http.ResponseWriter, r *http.Request, max int) int {
	if v := r.URL.Query().Get(SiteAccountQuery); v != "" {
		n := parseAccount(v, max)
		b.setSiteAccount(w, r, n)
		return n
	}
	if c, err := r.Cookie(siteAccountCookieName); err == nil {
		return parseAccount(c.Value, max)
	}
	return 0
}

func (b *Builder) setSiteAccount(w http.ResponseWriter, r *http.Request, n int) {
	v := strconv.Itoa(n)
	if c, err := r.Cookie(siteAccountCookieName); err == nil && c.Value == v {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     siteAccountCookieName,
		Value:    v,
		Path:     b.cookieConfig.Path,
		Domain:   b.cookieConfig.Domain,
		MaxAge:   accountsCookieMaxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

// SiteAccountURL is p (a page of the site) in the account n.
func SiteAccountURL(p string, n int) string {
	u, err := url.Parse(p)
	if err != nil {
		return p
	}
	q := u.Query()
	q.Set(SiteAccountQuery, strconv.Itoa(n))
	u.RawQuery = q.Encode()
	return u.String()
}

func parseAccount(s string, max int) int {
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 || n >= max {
		return 0
	}
	return n
}

func (b *Builder) setLoginAccount(w http.ResponseWriter, n int) {
	http.SetCookie(w, &http.Cookie{
		Name:     loginAccountCookieName,
		Value:    strconv.Itoa(n),
		Path:     b.cookieConfig.Path,
		Domain:   b.cookieConfig.Domain,
		MaxAge:   loginAccountMaxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

// setAddAccountContinue: the login of another account goes back to the page
// it was asked from (a path of this site), in that account.
func (b *Builder) setAddAccountContinue(w http.ResponseWriter, to string, n int) {
	if to == "" || !strings.HasPrefix(to, "/") || strings.HasPrefix(to, "//") {
		return
	}
	u, err := url.Parse(to)
	if err != nil {
		return
	}
	stripped := u.Path
	for _, prefix := range b.accountPrefixes {
		if s, _, ok := account.Split(prefix, u.Path); ok {
			stripped = s
			break
		}
	}
	to = u.String()
	if b.underAccountPrefix(stripped) {
		u.Path = b.AccountURL(stripped, n)
		to = u.String()
	} else {
		// a page of the site: it, in the new account
		to = SiteAccountURL(to, n)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     b.continueUrlCookieName,
		Value:    to,
		Path:     b.cookieConfig.Path,
		Domain:   b.cookieConfig.Domain,
		HttpOnly: true,
	})
}

func (b *Builder) clearLoginAccount(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     loginAccountCookieName,
		Value:    "",
		Path:     b.cookieConfig.Path,
		Domain:   b.cookieConfig.Domain,
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		HttpOnly: true,
		Secure:   true,
	})
}

// accountWriter keeps the account of the response in its redirects: a path
// under an account prefix gets its /u/N, a page of the login its account=N.
type accountWriter struct {
	http.ResponseWriter
	b           *Builder
	st          *accountState
	host        string
	wroteHeader bool
}

func (aw *accountWriter) WriteHeader(code int) {
	if !aw.wroteHeader {
		aw.wroteHeader = true
		if loc := aw.Header().Get("Location"); loc != "" {
			aw.Header().Set("Location", aw.b.accountLocation(loc, aw.st.n, aw.host))
		}
	}
	aw.ResponseWriter.WriteHeader(code)
}

func (aw *accountWriter) Write(p []byte) (int, error) {
	if !aw.wroteHeader {
		aw.WriteHeader(http.StatusOK)
	}
	return aw.ResponseWriter.Write(p)
}

func (aw *accountWriter) Unwrap() http.ResponseWriter { return aw.ResponseWriter }

func (aw *accountWriter) Flush() {
	if f, ok := aw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (aw *accountWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := aw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, errors.New("login: the response cannot be hijacked")
}

// accountLocation is loc for the account n: a path of this site (host —
// "//host/…" too —, or none).
func (b *Builder) accountLocation(loc string, n int, host string) string {
	if n <= 0 {
		return loc
	}
	u, err := url.Parse(loc)
	if err != nil || (u.Host != "" && u.Host != host) {
		return loc
	}
	if b.isLoginFlow(strings.TrimRight(u.Path, "/")) || u.Path == b.logoutURL {
		q := u.Query()
		if q.Get("account") == "" && q.Get("add") == "" {
			q.Set("account", strconv.Itoa(n))
			u.RawQuery = q.Encode()
		}
		return u.String()
	}
	for _, prefix := range b.accountPrefixes {
		if p := account.URL(prefix, u.Path, n); p != u.Path {
			u.Path = p
			return u.String()
		}
	}
	return loc
}

// AccountURL is the path p for the account n (its /u/N, under an account
// prefix).
func (b *Builder) AccountURL(p string, n int) string {
	for _, prefix := range b.accountPrefixes {
		if q := account.URL(prefix, p, n); q != p {
			return q
		}
	}
	return p
}

// AccountHomeURL is the home of the account n (the first account prefix).
func (b *Builder) AccountHomeURL(n int) string {
	if len(b.accountPrefixes) == 0 {
		return "/"
	}
	return account.Prefix(b.accountPrefixes[0], n) + "/"
}

// AddAccountURL is the login of another account, back to redirectTo (the
// login page; the page of the secure key takes add=1 too).
func (b *Builder) AddAccountURL(redirectTo string) string {
	q := url.Values{"add": {"1"}}
	if redirectTo != "" {
		q.Set("redirect_to", redirectTo)
	}
	return b.loginPageURL + "?" + q.Encode()
}

// AccountLogoutURL is the logout of the account n; all: of every account.
func (b *Builder) AccountLogoutURL(n int, all bool) string {
	if all {
		return b.logoutURL + "?all=1"
	}
	return b.logoutURL + "?account=" + strconv.Itoa(n)
}

// AccountLoginURL is the login of the account n (its session expired).
func (b *Builder) AccountLoginURL(n int) string {
	return b.loginPageURL + "?account=" + strconv.Itoa(n)
}

// --- the list of the accounts (rvq_accounts) ---

func (b *Builder) accountsSig(payload string) string {
	m := hmac.New(sha256.New, []byte(b.secret+"\x00accounts"))
	m.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

// readAccounts is the list of the browser: index → user id.
func (b *Builder) readAccounts(r *http.Request) map[int]string {
	out := map[int]string{}
	c, err := r.Cookie(accountsCookieName)
	if err != nil || c.Value == "" {
		return out
	}
	v, err := url.QueryUnescape(c.Value)
	if err != nil {
		return out
	}
	i := strings.LastIndexByte(v, '.')
	if i < 0 || !hmac.Equal([]byte(b.accountsSig(v[:i])), []byte(v[i+1:])) {
		return out
	}
	max := b.GetMaxAccounts()
	for _, e := range strings.Split(v[:i], ",") {
		n, id, ok := strings.Cut(e, ":")
		if !ok || id == "" {
			continue
		}
		if k, err := strconv.Atoi(n); err == nil && k >= 0 && k < max {
			out[k] = id
		}
	}
	return out
}

func (b *Builder) writeAccounts(w http.ResponseWriter, r *http.Request, list map[int]string) {
	keys := make([]int, 0, len(list))
	for k := range list {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, strconv.Itoa(k)+":"+list[k])
	}
	payload := strings.Join(parts, ",")
	c := &http.Cookie{
		Name:     accountsCookieName,
		Path:     b.cookieConfig.Path,
		Domain:   b.cookieConfig.Domain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	if payload == "" {
		c.MaxAge, c.Expires = -1, time.Unix(1, 0)
	} else {
		c.Value = url.QueryEscape(payload + "." + b.accountsSig(payload))
		c.MaxAge = accountsCookieMaxAge
	}
	http.SetCookie(w, c)
	// the rest of the request sees it
	setCookieForRequest(r, &http.Cookie{Name: accountsCookieName, Value: c.Value})
}

// rememberAccount puts the account of r — userID — in the list.
func (b *Builder) rememberAccount(w http.ResponseWriter, r *http.Request, userID string) {
	if b.GetMaxAccounts() <= 1 || userID == "" {
		return
	}
	list := b.readAccounts(r)
	n := account.IndexOf(r)
	if list[n] == userID {
		return
	}
	list[n] = userID
	b.writeAccounts(w, r, list)
}

// forgetAccount takes the account n out of the list.
func (b *Builder) forgetAccount(w http.ResponseWriter, r *http.Request, n int) {
	list := b.readAccounts(r)
	if _, ok := list[n]; !ok {
		return
	}
	delete(list, n)
	b.writeAccounts(w, r, list)
}

// accountOfUser is r of the index the browser already has userID in — a
// login of an account already logged in is that account, not a new one —; r
// as it is otherwise.
func (b *Builder) accountOfUser(w http.ResponseWriter, r *http.Request, userID string) *http.Request {
	if b.GetMaxAccounts() <= 1 {
		return r
	}
	cur := account.IndexOf(r)
	for n, id := range b.readAccounts(r) {
		if id == userID && n != cur {
			if b.isLoginFlow(strings.TrimRight(r.URL.Path, "/")) {
				b.setLoginAccount(w, n)
			}
			return withAccount(r, n)
		}
	}
	return r
}

// freeAccount is the index of another account: the first one the browser
// has neither in its list nor logged in; all taken, the first whose session
// expired; none: the last one.
func (b *Builder) freeAccount(r *http.Request) int {
	max := b.GetMaxAccounts()
	list := b.readAccounts(r)
	for n := 0; n < max; n++ {
		if _, ok := list[n]; ok {
			continue
		}
		if _, err := r.Cookie(account.CookieName(b.authCookieName, n)); err == nil {
			continue
		}
		return n
	}
	for n := 0; n < max; n++ {
		if !b.loggedIn(withAccount(r, n)) {
			return n
		}
	}
	return max - 1
}

// loggedIn: r has a valid session (the claims of its account).
func (b *Builder) loggedIn(r *http.Request) bool {
	_, err := b.ParseClaims(r)
	return err == nil
}

// Accounts are the accounts the browser logged in, by index: the current
// one among them.
func (b *Builder) Accounts(r *http.Request) []*Account {
	if b.GetMaxAccounts() <= 1 {
		return nil
	}
	cur := account.IndexOf(r)
	list := b.readAccounts(r)
	// an account logged in that the list lost (a login whose request did not
	// carry it): by its cookies
	for n := 0; n < b.GetMaxAccounts(); n++ {
		if _, ok := list[n]; ok {
			continue
		}
		if _, err := r.Cookie(account.CookieName(b.authCookieName, n)); err != nil {
			continue
		}
		if c, err := b.ParseClaims(r.WithContext(account.WithIndex(r.Context(), n))); err == nil {
			list[n] = c.UserID
		}
	}
	keys := make([]int, 0, len(list))
	for k := range list {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	out := make([]*Account, 0, len(keys))
	for _, n := range keys {
		rn := r.WithContext(account.WithIndex(r.Context(), n))
		a := &Account{Index: n, UserID: list[n], Current: n == cur}
		if c, err := b.ParseClaims(rn); err == nil && c.UserID == a.UserID {
			a.LoggedIn = true
		}
		if b.userModel != nil {
			if u, err := b.findUserByID(a.UserID); err == nil {
				a.User = u
			}
		}
		out = append(out, a)
	}
	return out
}

// OtherLoggedInAccount is the logged-in account of the smallest index other
// than n: -1 when none.
func (b *Builder) OtherLoggedInAccount(r *http.Request, n int) int {
	for k := 0; k < b.GetMaxAccounts(); k++ {
		if k != n && b.loggedIn(r.WithContext(account.WithIndex(r.Context(), k))) {
			return k
		}
	}
	return -1
}

// AccountLogin is what the login page says of the account it logs in, in a
// browser holding several: another one (Adding), or one whose session expired
// (Again: its account already in the field); and the way back to the one
// logged in (BackURL, BackName). nil: the plain login.
type AccountLogin struct {
	Adding   bool
	Again    bool
	Account  string
	Name     string
	BackURL  string
	BackName string
}

// AccountLoginOf is the AccountLogin of the login page of r.
func (b *Builder) AccountLoginOf(r *http.Request) *AccountLogin {
	if b.GetMaxAccounts() <= 1 {
		return nil
	}
	n := account.IndexOf(r)
	list := b.readAccounts(r)
	id, known := list[n]
	adding := r.URL.Query().Get("add") == "1" || n > 0
	if !known && !adding {
		return nil
	}
	al := &AccountLogin{Adding: !known}
	if known && b.userModel != nil {
		if u, err := b.findUserByID(id); err == nil {
			al.Again = true
			al.Adding = false
			if up, ok := u.(UserPasser); ok {
				al.Account = up.GetAccountName()
			}
			al.Name = accountName(u, al.Account)
		}
	}
	if m := b.OtherLoggedInAccount(r, n); m >= 0 {
		al.BackURL = b.AccountHomeURL(m)
		if c, err := b.ParseClaims(r.WithContext(account.WithIndex(r.Context(), m))); err == nil && b.userModel != nil {
			if u, err := b.findUserByID(c.UserID); err == nil {
				var acc string
				if up, ok := u.(UserPasser); ok {
					acc = up.GetAccountName()
				}
				al.BackName = accountName(u, acc)
			}
		}
	}
	if !al.Adding && !al.Again && al.BackURL == "" {
		return nil
	}
	return al
}

// accountName is the name of u, else its account.
func accountName(u any, acc string) string {
	if n, ok := u.(interface{ GetName() string }); ok && n.GetName() != "" {
		return n.GetName()
	}
	return acc
}
