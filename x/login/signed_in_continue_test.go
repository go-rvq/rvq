package login

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/x/i18n"
)

// Signed in already, the login page goes where the login was to go — its
// continue URL (the end of the login dialog of the admin, which closes it) —,
// else to the home.
func TestSignedInLoginPageGoesToContinue(t *testing.T) {
	b := New(i18n.New()).Secret("a secret worth at least this many bytes").SessionMaxAge(3600).TOTP(false).
		HomeURLFunc(func(r *http.Request, user any) string { return "/home" })
	claims := &UserClaims{UserID: "1", RegisteredClaims: b.genBaseSessionClaim("1", false)}
	token := b.mustGetSessionToken(*claims)
	handler := b.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	get := func(continueURL string) *httptest.ResponseRecorder {
		r := httptest.NewRequest("GET", b.loginPageURL, nil)
		r.AddCookie(&http.Cookie{Name: b.authCookieOf(r), Value: token})
		if continueURL != "" {
			r.AddCookie(&http.Cookie{Name: cookieName(b.continueUrlCookieName), Value: continueURL})
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w
	}
	if w := get("/admin/login-done"); w.Code != http.StatusFound || w.Header().Get("Location") != "/admin/login-done" {
		t.Errorf("with a continue URL: %d to %q", w.Code, w.Header().Get("Location"))
	}
	if w := get(""); w.Code != http.StatusFound || w.Header().Get("Location") != "/home" {
		t.Errorf("without: %d to %q", w.Code, w.Header().Get("Location"))
	}
}
