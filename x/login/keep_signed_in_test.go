package login

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-rvq/rvq/x/i18n"
)

// A session asked to stay signed in on the device lasts the long max age —
// its claims, its cookie, and the max age the request says (the log of the
// sessions) —; one not asked, the SessionMaxAge; the option off, always it.
func TestKeepSignedIn(t *testing.T) {
	b := New(i18n.New()).Secret("s3cret").SessionMaxAge(3600)
	if b.KeepSignedInEnabled() {
		t.Fatal("on by default")
	}
	if got := b.genBaseSessionClaim("u", false, true); !near(got.ExpiresAt.Time, time.Hour) {
		t.Fatalf("off: keep lasts %v, want the session's hour", time.Until(got.ExpiresAt.Time))
	}

	const month = 30 * 24 * 3600
	b.KeepSignedIn(month)
	keep := UserClaims{UserID: "u", KeepSignedIn: true, RegisteredClaims: b.genBaseSessionClaim("u", false, true)}
	if !near(keep.ExpiresAt.Time, month*time.Second) {
		t.Fatalf("keep lasts %v, want 30 days", time.Until(keep.ExpiresAt.Time))
	}
	plain := UserClaims{UserID: "u", RegisteredClaims: b.genBaseSessionClaim("u", false)}
	if !near(plain.ExpiresAt.Time, time.Hour) {
		t.Fatalf("not kept lasts %v, want an hour", time.Until(plain.ExpiresAt.Time))
	}
	if a := cookieMaxAge(&keep); a < month-5 || a > month {
		t.Fatalf("the cookie of keep: %d", a)
	}

	req := func(c UserClaims) *http.Request {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: b.authCookieOf(r), Value: b.mustGetSessionToken(c)})
		return r
	}
	if got := b.SessionMaxAgeOf(req(keep)); got != month {
		t.Fatalf("SessionMaxAgeOf keep %d", got)
	}
	if got := b.SessionMaxAgeOf(req(plain)); got != 3600 {
		t.Fatalf("SessionMaxAgeOf plain %d", got)
	}
	if got := b.SessionMaxAgeOf(httptest.NewRequest(http.MethodGet, "/", nil)); got != 3600 {
		t.Fatalf("SessionMaxAgeOf no session %d", got)
	}
	for _, v := range []string{"1", "true", "on"} {
		if !isTrue(v) {
			t.Fatalf("%q not checked", v)
		}
	}
	if isTrue("") || isTrue("0") {
		t.Fatal("unchecked read as checked")
	}
}

func near(t time.Time, d time.Duration) bool {
	left := time.Until(t)
	return left > d-5*time.Second && left <= d
}

// With the initial user logged (SetInitialUserLoggerd), the session of the
// initial user lasts until the end of the day — at the login and renewed
// alike —; the one of another user, the SessionMaxAge.
func TestInitialUserSessionRenewed(t *testing.T) {
	b := New(i18n.New()).Secret("s3cret").SessionMaxAge(3600).
		InitialUserAccount("admin").SetInitialUserLoggerd(true)
	endOfDay := func(tm time.Time) bool {
		n := time.Now()
		return tm.Equal(time.Date(n.Year(), n.Month(), n.Day()+1, 0, 0, 0, 0, n.Location())) ||
			tm.Sub(time.Date(n.Year(), n.Month(), n.Day()+1, 0, 0, 0, 0, n.Location())).Abs() < 2*time.Second
	}
	if n := time.Now(); n.Hour() == 23 {
		t.Skip("an hour and the end of the day are too near")
	}

	login := b.genBaseSessionClaim("1", true)
	if !endOfDay(login.ExpiresAt.Time) {
		t.Fatalf("the login of the initial user ends %v, want the end of the day", login.ExpiresAt.Time)
	}
	renewed := b.renewedSessionClaims(&UserClaims{UserID: "1"}, "admin")
	if !endOfDay(renewed.ExpiresAt.Time) {
		t.Fatalf("the renewal of the initial user ends %v, want the end of the day", renewed.ExpiresAt.Time)
	}
	other := b.renewedSessionClaims(&UserClaims{UserID: "2"}, "ana")
	if !near(other.ExpiresAt.Time, time.Hour) {
		t.Fatalf("the renewal of another user ends %v, want an hour", other.ExpiresAt.Time)
	}
}
