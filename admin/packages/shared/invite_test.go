package shared

import (
	"errors"
	"testing"

	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// recordingNotifier captures the notifier callbacks.
type recordingNotifier struct{ invited, accepted int }

func (r *recordingNotifier) ShareInvited(*gorm.DB, *ShareInvite) error  { r.invited++; return nil }
func (r *recordingNotifier) ShareAccepted(*gorm.DB, *ShareInvite) error { r.accepted++; return nil }

func inviteTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&perm.DefaultDBPolicy{}); err != nil {
		t.Fatal(err)
	}
	if err := AutoMigrateInvites(db); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestInviteAcceptGrants(t *testing.T) {
	db := inviteTestDB(t)
	const resource = "presets:orgs:o1:*"
	n := &recordingNotifier{}

	invites, err := Invite(db, resource, []string{"ana", "bruno"}, VerbView, "owner", n)
	if err != nil {
		t.Fatal(err)
	}
	if len(invites) != 2 || n.invited != 2 {
		t.Fatalf("want 2 invites + 2 notifications, got %d/%d", len(invites), n.invited)
	}

	// no access is granted before acceptance
	var pols int64
	db.Model(&perm.DefaultDBPolicy{}).Count(&pols)
	if pols != 0 {
		t.Fatalf("no permission should exist before acceptance, got %d", pols)
	}

	// ana appears in her pending list
	pend, _ := PendingFor(db, "ana")
	if len(pend) != 1 || pend[0].ID != invites[0].ID {
		t.Fatalf("ana should have 1 pending invite, got %v", pend)
	}

	// ana accepts -> a shared policy is created and the inviter is notified
	if err := Accept(db, invites[0].ID, n); err != nil {
		t.Fatal(err)
	}
	if n.accepted != 1 {
		t.Errorf("acceptance should notify the inviter")
	}
	subs, _ := Subjects(db, invites[0].SharedID)
	if len(subs) != 1 || subs[0].Subject != "ana" {
		t.Fatalf("accepting should grant ana a shared policy, got %v", subs)
	}
	if pend, _ = PendingFor(db, "ana"); len(pend) != 0 {
		t.Errorf("accepted invite should leave the pending list")
	}

	// accepting twice is rejected
	if err := Accept(db, invites[0].ID, n); !errors.Is(err, ErrInviteNotPending) {
		t.Errorf("re-accepting should fail with ErrInviteNotPending, got %v", err)
	}

	// bruno rejects -> no policy, invite gone from pending
	if err := Reject(db, invites[1].ID); err != nil {
		t.Fatal(err)
	}
	if pend, _ = PendingFor(db, "bruno"); len(pend) != 0 {
		t.Errorf("rejected invite should leave the pending list")
	}
	if subs, _ := Subjects(db, invites[1].SharedID); len(subs) != 1 {
		// only ana's grant exists under the shared id (bruno rejected)
		t.Errorf("only the accepted subject should have a policy, got %v", subs)
	}
	// rejecting a non-pending invite fails
	if err := Reject(db, uuid.New()); !errors.Is(err, ErrInviteNotPending) {
		t.Errorf("rejecting a missing invite should fail, got %v", err)
	}
}

func TestAcceptForVerifiesSubject(t *testing.T) {
	db := inviteTestDB(t)
	invites, err := Invite(db, "presets:orgs:o1:*", []string{"ana"}, VerbView, "owner", nil)
	if err != nil {
		t.Fatal(err)
	}
	id := invites[0].ID

	// another user cannot accept ana's invite
	if err := AcceptFor(db, id, "mallory", nil); !errors.Is(err, ErrInviteNotYours) {
		t.Fatalf("a foreign user must not accept the invite, got %v", err)
	}
	if err := RejectFor(db, id, "mallory"); !errors.Is(err, ErrInviteNotYours) {
		t.Fatalf("a foreign user must not reject the invite, got %v", err)
	}
	// no permission was granted by the failed attempts
	var pols int64
	db.Model(&perm.DefaultDBPolicy{}).Count(&pols)
	if pols != 0 {
		t.Fatalf("no permission should exist after foreign attempts, got %d", pols)
	}

	// the addressed user accepts -> granted
	if err := AcceptFor(db, id, "ana", nil); err != nil {
		t.Fatal(err)
	}
	if subs, _ := Subjects(db, invites[0].SharedID); len(subs) != 1 {
		t.Errorf("ana's acceptance should grant one policy, got %d", len(subs))
	}
}
