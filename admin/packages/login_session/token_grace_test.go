package login_session

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
	"github.com/google/uuid"
)

// A session renewed: its new token good; the one it had before still good
// for a moment (PrevTokenGrace) — a request already on its way with it, the
// page refreshing its parts by itself, is not turned away —, and no more
// after it; another token never.
func TestPrevTokenGrace(t *testing.T) {
	db := testDB(t)
	lb := login.New(i18n.New()).AuthCookieName("auth")
	m := NewManager(lb)
	user := uuid.New()
	req := func(token string) *http.Request {
		r := httptest.NewRequest("GET", "/admin", nil)
		r.RemoteAddr = "203.0.113.9:5000"
		r.AddCookie(&http.Cookie{Name: "auth", Value: token})
		return r
	}
	if login.GetSessionToken(lb, req("old")) != "old" {
		t.Skip("the cookie of the session is not named auth here")
	}
	if err := m.AddSessionLogByUserID(db, req("old"), user); err != nil {
		t.Fatal(err)
	}
	valid := func(token string) bool {
		t.Helper()
		ok, err := m.CheckIsTokenValidFromRequest(db, req(token), user)
		if err != nil {
			t.Fatal(err)
		}
		return ok
	}
	if !valid("old") {
		t.Fatal("the token of the session not good")
	}
	// renewed: the request with the new token, the old one given
	if err := m.UpdateCurrentSessionLog(db, req("new"), user, "old"); err != nil {
		t.Fatal(err)
	}
	if !valid("new") || !valid("old") || valid("other") {
		t.Errorf("after the renewal: new %v, old %v, other %v", valid("new"), valid("old"), valid("other"))
	}
	// the moment over
	past := time.Now().Add(-time.Second)
	db.Model(&LoginSession{}).Where("user_id = ?", user).Update("prev_token_until", &past)
	if valid("old") || !valid("new") {
		t.Errorf("after the grace: old %v, new %v", valid("old"), valid("new"))
	}
}
