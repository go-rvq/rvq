package shared

import (
	"bytes"
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/login"
	"github.com/google/uuid"
)

// fakeUser is a minimal user.User (login.UserPass provides the UserPasser part).
type fakeUser struct {
	login.UserPass
	id uuid.UUID
}

func (f *fakeUser) GetID() uuid.UUID              { return f.id }
func (f *fakeUser) GetName() string               { return "Test" }
func (f *fakeUser) SetName(string)                {}
func (f *fakeUser) SetEmail(string)               {}
func (f *fakeUser) SetRegistrationDate(time.Time) {}
func (f *fakeUser) GetStatus() string             { return "active" }
func (f *fakeUser) GetRoles() role.Roles          { return nil }
func (f *fakeUser) SetRoles(role.Roles)           {}

func ctxWithUser(id uuid.UUID) *web.EventContext {
	r := httptest.NewRequest("GET", "/admin", nil)
	r = r.WithContext(context.WithValue(r.Context(), login.UserKey, &fakeUser{id: id}))
	return &web.EventContext{R: r}
}

func TestInvitesBody(t *testing.T) {
	db := inviteTestDB(t)
	me := uuid.New()

	// an invite addressed to me and one to someone else
	const resource = "presets:orgs:o1:*"
	if _, err := Invite(db, resource, []string{me.String()}, VerbView, "owner-x", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := Invite(db, "presets:orgs:o2:*", []string{uuid.New().String()}, VerbView, "owner-y", nil); err != nil {
		t.Fatal(err)
	}

	ctx := ctxWithUser(me)
	var buf bytes.Buffer
	if err := h.Fprint(&buf, invitesBody(db, ctx), ctx.Context()); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	// shows my invite (resource + inviter) with accept/reject events
	if !strings.Contains(html, resource) || !strings.Contains(html, "owner-x") {
		t.Errorf("my pending invite should be listed; got:\n%s", html)
	}
	if !strings.Contains(html, invitesAcceptEvent) || !strings.Contains(html, invitesRejectEvent) {
		t.Errorf("invite should offer accept/decline actions; got:\n%s", html)
	}
	// does not show the other user's invite
	if strings.Contains(html, "owner-y") {
		t.Errorf("must not list invites addressed to other users")
	}

	// a user with no invites renders nothing
	var empty bytes.Buffer
	if err := h.Fprint(&empty, invitesBody(db, ctxWithUser(uuid.New())), ctx.Context()); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(empty.String(), resource) {
		t.Errorf("a user without invites should see none")
	}

	// the notification count reflects the current user's pending invites
	if got := NotificationCount(db)(ctx); got != 1 {
		t.Errorf("NotificationCount = %d, want 1", got)
	}
	if got := NotificationCount(db)(ctxWithUser(uuid.New())); got != 0 {
		t.Errorf("NotificationCount for a user without invites = %d, want 0", got)
	}
}
