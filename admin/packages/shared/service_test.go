package shared_test

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/packages/mail-sender/mailtest"
	"github.com/go-rvq/rvq/admin/packages/shared"
	"github.com/go-rvq/rvq/admin/packages/shared/sharedtest"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	view = "view"
	edit = "edit"
)

// fixture is a service on the models of sharedtest: folders (owned) and
// their documents; ana owns, bia, caio and dani are others; the roles of
// each resource Viewer, Editor and None.
type fixture struct {
	t     *testing.T
	ctx   context.Context
	db    *gorm.DB
	s     *shared.Service
	users *sharedtest.Users
	mail  *mailtest.Server
	now   time.Time
	roles map[string]map[string]uuid.UUID
}

var origin = shared.Origin{IP: "127.0.0.1", UserAgent: "test"}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	// one database for every connection of the test
	db, err := gorm.Open(sqlite.Open("file:"+strings.ReplaceAll(t.Name(), "/", "_")+"?mode=memory&cache=shared"),
		&gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(1)
		t.Cleanup(func() { _ = sqlDB.Close() })
	}
	if err := sharedtest.Migrate(db); err != nil {
		t.Fatal(err)
	}
	// the mails to a mail server of the test's own
	ms, err := mailtest.Start(mailtest.Options{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ms.Close() })
	f := &fixture{t: t, ctx: context.Background(), db: db, users: sharedtest.NewUsers(), mail: ms,
		now: time.Date(2026, 11, 8, 12, 0, 0, 0, time.UTC), roles: map[string]map[string]uuid.UUID{}}
	f.users.Add(&sharedtest.User{ID: "ana", Name: "Ana", Email: "ana@x.test", Language: "pt-BR"})
	f.users.Add(&sharedtest.User{ID: "bia", Name: "Bia", Email: "bia@x.test", Language: "en-US"})
	f.users.Add(&sharedtest.User{ID: "caio", Name: "Caio", Email: "caio@x.test"})
	f.users.Add(&sharedtest.User{ID: "dani", Name: "Dani", Email: "dani@x.test"})
	f.s, err = shared.New(shared.Options{DB: db, Users: f.users, Sender: ms.Sender(), ConfirmURL: sharedtest.ConfirmURL,
		Now: func() time.Time { return f.now }})
	if err != nil {
		t.Fatal(err)
	}
	f.s.Register(shared.Resource{Name: "pastas", Table: "sharedtest_pastas", OwnerColumn: "owner_id",
		Label: func(_ context.Context, db *gorm.DB, id string) string {
			var p sharedtest.Pasta
			db.First(&p, "id = ?", id)
			return p.Nome
		}})
	f.s.Register(shared.Resource{Name: "documentos", Table: "sharedtest_documentos", Parent: "pastas", ParentColumn: "pasta_id"})
	for _, res := range []string{"pastas", "documentos"} {
		if err := shared.SeedRoles(db, uuid.Nil, res,
			shared.RoleSeed{Name: "Viewer", Actions: []string{view}},
			shared.RoleSeed{Name: "Editor", Actions: []string{view, edit}}); err != nil {
			t.Fatal(err)
		}
		roles, err := shared.RolesOf(db, uuid.Nil, res)
		if err != nil {
			t.Fatal(err)
		}
		f.roles[res] = map[string]uuid.UUID{}
		for _, r := range roles {
			f.roles[res][r.Name] = r.ID
		}
	}
	return f
}

func (f *fixture) role(res, name string) uuid.UUID { return f.roles[res][name] }

// share asks the sharing of the record with the subject as the role (by
// by), confirmed by the owner, the invite (a user's) accepted.
func (f *fixture) share(res, id, subject, role, by string) {
	f.t.Helper()
	if _, err := f.s.Share(f.ctx, res, id, subject, f.role(res, role), by, origin); err != nil {
		f.t.Fatal(err)
	}
	f.confirmLast()
	if user, ok := strings.CutPrefix(subject, shared.SubjectUser); ok {
		f.acceptAll(user)
	}
}

// confirmLast confirms, as the owner, the request of the last mail.
func (f *fixture) confirmLast() {
	f.t.Helper()
	tok := f.token()
	req, err := f.s.RequestOf(f.ctx, tok)
	if err != nil {
		f.t.Fatal(err)
	}
	if err := f.s.Confirm(f.ctx, tok, req.Owner, origin); err != nil {
		f.t.Fatal(err)
	}
}

// last is the last mail sent.
func (f *fixture) last() *mailtest.Message {
	f.t.Helper()
	m := f.mail.Last()
	if m == nil {
		f.t.Fatal("no mail")
	}
	return m
}

// token is the token of the link of the last mail.
func (f *fixture) token() string { return sharedtest.Token(f.last().Body()) }

func (f *fixture) acceptAll(user string) {
	f.t.Helper()
	invs, err := f.s.Invites(f.ctx, user)
	if err != nil {
		f.t.Fatal(err)
	}
	for _, inv := range invs {
		if err := f.s.Accept(f.ctx, inv.ID, user, origin); err != nil {
			f.t.Fatal(err)
		}
	}
}

func (f *fixture) roleName(user, res, id string) string {
	f.t.Helper()
	e, err := f.s.RoleOf(f.ctx, user, res, id)
	if err != nil {
		f.t.Fatal(err)
	}
	switch {
	case e.Owner:
		return "owner"
	case e.Role == nil:
		return ""
	}
	return e.Role.Name
}

func (f *fixture) events(res, id string) (out []string) {
	f.t.Helper()
	h, err := f.s.History(f.ctx, res, id)
	if err != nil {
		f.t.Fatal(err)
	}
	for _, e := range h {
		out = append(out, e.Event)
	}
	return
}

// A sharing: the owner asked by e-mail, in their language, the link to
// confirm; nothing until they confirm (only they), then the invite, the
// access once the user accepts (only them); the token once; all of it in
// the history.
func TestShareConfirmAccept(t *testing.T) {
	f := newFixture(t)
	p := sharedtest.NewPasta(f.db, "Contratos", "ana")
	req, err := f.s.Share(f.ctx, "pastas", p.ID, shared.UserSubject("bia"), f.role("pastas", "Viewer"), "ana", origin)
	if err != nil {
		t.Fatal(err)
	}
	m := f.last()
	if m.To[0] != "ana@x.test" || m.Subject != "Confirme: compartilhamento de Contratos" ||
		!strings.Contains(m.Body(), "Bia") {
		t.Fatalf("the mail to the owner, in Portuguese: %q %q %q", m.To, m.Subject, m.Body())
	}
	if req.State != shared.StatePending || f.roleName("bia", "pastas", p.ID) != "" {
		t.Error("access before the owner confirms")
	}
	tok := sharedtest.Token(m.Body())
	if err := f.s.Confirm(f.ctx, tok, "bia", origin); !errors.Is(err, shared.ErrNotOwner) {
		t.Errorf("confirmed by another: %v", err)
	}
	if err := f.s.Confirm(f.ctx, tok, "ana", origin); err != nil {
		t.Fatal(err)
	}
	if err := f.s.Confirm(f.ctx, tok, "ana", origin); !errors.Is(err, shared.ErrNotPending) {
		t.Errorf("the token twice: %v", err)
	}
	if f.roleName("bia", "pastas", p.ID) != "" {
		t.Error("access before the user accepts")
	}
	invs, _ := f.s.Invites(f.ctx, "bia")
	if len(invs) != 1 {
		t.Fatalf("invites: %d", len(invs))
	}
	if err := f.s.Accept(f.ctx, invs[0].ID, "caio", origin); !errors.Is(err, shared.ErrNotYours) {
		t.Errorf("accepted by another: %v", err)
	}
	if err := f.s.Accept(f.ctx, invs[0].ID, "bia", origin); err != nil {
		t.Fatal(err)
	}
	if f.roleName("bia", "pastas", p.ID) != "Viewer" || !f.s.Can(f.ctx, "bia", "pastas", p.ID, view) ||
		f.s.Can(f.ctx, "bia", "pastas", p.ID, edit) || !f.s.Can(f.ctx, "ana", "pastas", p.ID, edit) {
		t.Error("the role after accepting")
	}
	want := []string{shared.EventConfirmationRequested, shared.EventOwnerConfirmed, shared.EventInvited,
		shared.EventAccepted, shared.EventGranted}
	if got := f.events("pastas", p.ID); !slices.Equal(got, want) {
		t.Errorf("history %v", got)
	}
	if _, err := f.s.Share(f.ctx, "pastas", p.ID, shared.UserSubject("bia"), f.role("pastas", "Editor"), "ana", origin); !errors.Is(err, shared.ErrAlreadyShared) {
		t.Errorf("twice: %v", err)
	}
	if _, err := f.s.Share(f.ctx, "pastas", p.ID, shared.UserSubject("ninguem"), f.role("pastas", "Editor"), "ana", origin); !errors.Is(err, shared.ErrNoSuchSubject) {
		t.Errorf("a user there is not: %v", err)
	}
}

// The mail in the owner's language: English.
func TestMailLanguage(t *testing.T) {
	f := newFixture(t)
	p := sharedtest.NewPasta(f.db, "Docs", "bia")
	if _, err := f.s.Share(f.ctx, "pastas", p.ID, shared.UserSubject("caio"), f.role("pastas", "Editor"), "bia", origin); err != nil {
		t.Fatal(err)
	}
	if s := f.last().Subject; s != "Confirm: sharing Docs" {
		t.Errorf("subject %q", s)
	}
}

// A confirmation late (past 48h) expires; one refused is closed.
func TestExpiredAndRejected(t *testing.T) {
	f := newFixture(t)
	p := sharedtest.NewPasta(f.db, "P", "ana")
	if _, err := f.s.Share(f.ctx, "pastas", p.ID, shared.UserSubject("bia"), f.role("pastas", "Viewer"), "ana", origin); err != nil {
		t.Fatal(err)
	}
	tok := f.token()
	f.now = f.now.Add(49 * time.Hour)
	if err := f.s.Confirm(f.ctx, tok, "ana", origin); !errors.Is(err, shared.ErrExpired) {
		t.Errorf("late: %v", err)
	}
	if err := f.s.Confirm(f.ctx, tok, "ana", origin); !errors.Is(err, shared.ErrNotPending) {
		t.Errorf("expired, again: %v", err)
	}
	if _, err := f.s.Share(f.ctx, "pastas", p.ID, shared.UserSubject("caio"), f.role("pastas", "Viewer"), "ana", origin); err != nil {
		t.Fatal(err)
	}
	if err := f.s.Reject(f.ctx, f.token(), "ana", origin); err != nil {
		t.Fatal(err)
	}
	if inv, _ := f.s.Invites(f.ctx, "caio"); len(inv) != 0 {
		t.Error("an invite of a request refused")
	}
	got := f.events("pastas", p.ID)
	if !slices.Contains(got, shared.EventConfirmationExpired) || !slices.Contains(got, shared.EventOwnerRejected) {
		t.Errorf("history %v", got)
	}
}

// A role with more actions waits for the owner (the old one until then);
// with fewer, at once; revoked, at once.
func TestRoleChangeAndRevoke(t *testing.T) {
	f := newFixture(t)
	p := sharedtest.NewPasta(f.db, "P", "ana")
	f.share("pastas", p.ID, shared.UserSubject("bia"), "Viewer", "ana")
	req, err := f.s.ChangeRole(f.ctx, "pastas", p.ID, shared.UserSubject("bia"), f.role("pastas", "Editor"), "ana", origin)
	if err != nil || req == nil {
		t.Fatalf("more: a request: %v %v", req, err)
	}
	if f.roleName("bia", "pastas", p.ID) != "Viewer" {
		t.Error("more, before the owner confirms")
	}
	f.confirmLast()
	if f.roleName("bia", "pastas", p.ID) != "Editor" {
		t.Error("more, confirmed")
	}
	req, err = f.s.ChangeRole(f.ctx, "pastas", p.ID, shared.UserSubject("bia"), f.role("pastas", "Viewer"), "ana", origin)
	if err != nil || req != nil || f.roleName("bia", "pastas", p.ID) != "Viewer" {
		t.Errorf("fewer, at once: %v %v", req, err)
	}
	if err := f.s.Revoke(f.ctx, "pastas", p.ID, shared.UserSubject("bia"), "ana", origin); err != nil {
		t.Fatal(err)
	}
	if f.roleName("bia", "pastas", p.ID) != "" {
		t.Error("revoked")
	}
	got := f.events("pastas", p.ID)
	if !slices.Contains(got, shared.EventRoleChangeRequested) || !slices.Contains(got, shared.EventRoleChanged) ||
		got[len(got)-1] != shared.EventRevoked {
		t.Errorf("history %v", got)
	}
	h, _ := f.s.History(f.ctx, "pastas", p.ID)
	if last := h[len(h)-1]; last.IP != "127.0.0.1" || len(last.ActionsBefore) != 1 || last.ActionsBefore[0] != view {
		t.Errorf("the revoke: %+v", last)
	}
}

// The records under one shared inherit its role; the most specific wins,
// for more or for less; None takes one out; the accessible records.
func TestInheritance(t *testing.T) {
	f := newFixture(t)
	p := sharedtest.NewPasta(f.db, "P", "ana")
	d1, d2, d3 := sharedtest.NewDocumento(f.db, p, "1"), sharedtest.NewDocumento(f.db, p, "2"), sharedtest.NewDocumento(f.db, p, "3")
	f.share("pastas", p.ID, shared.UserSubject("bia"), "Editor", "ana")
	if f.roleName("bia", "documentos", d1.ID) != "Editor" {
		t.Error("inherited")
	}
	if f.roleName("ana", "documentos", d1.ID) != "owner" {
		t.Error("the owner of the folder owns its documents")
	}
	f.share("documentos", d1.ID, shared.UserSubject("bia"), "Viewer", "ana")
	f.share("documentos", d2.ID, shared.UserSubject("bia"), shared.RoleNone, "ana")
	if f.roleName("bia", "documentos", d1.ID) != "Viewer" || f.s.Can(f.ctx, "bia", "documentos", d1.ID, edit) {
		t.Error("the most specific, for less")
	}
	if e, _ := f.s.RoleOf(f.ctx, "bia", "documentos", d2.ID); e.Visible() {
		t.Error("None: out")
	}
	ids, err := f.s.AccessibleIDs(f.ctx, "bia", "documentos")
	if err != nil {
		t.Fatal(err)
	}
	slices.Sort(ids)
	want := []string{d1.ID, d3.ID}
	slices.Sort(want)
	if !slices.Equal(ids, want) {
		t.Errorf("accessible documents %v, want %v", ids, want)
	}
	if ids, _ := f.s.AccessibleIDs(f.ctx, "ana", "documentos"); len(ids) != 3 {
		t.Errorf("the owner's: %d", len(ids))
	}
	if ids, _ := f.s.AccessibleIDs(f.ctx, "caio", "documentos"); len(ids) != 0 {
		t.Errorf("caio's: %d", len(ids))
	}
	var n int64
	q, _ := f.s.Scope(f.ctx, f.db.Model(&sharedtest.Documento{}), "bia", "documentos")
	q.Count(&n)
	if n != 2 {
		t.Errorf("the scope: %d", n)
	}
	subs, _ := f.s.Subjects(f.ctx, "documentos", d1.ID)
	if len(subs) != 2 || subs[0].RecordID != d1.ID || subs[1].RecordID != p.ID {
		t.Errorf("who reaches the document: %+v", subs)
	}
}

// A group: the access at once after the owner confirms (no invite); its
// members — new ones too — have it, one who leaves no more; the user's own
// access first, then their groups in their order.
func TestGroups(t *testing.T) {
	f := newFixture(t)
	f.users.AddGroup("contab", "Contabilidade")
	f.users.AddGroup("estag", "Estagiários")
	f.users.Join("caio", "contab")
	p := sharedtest.NewPasta(f.db, "P", "ana")
	if _, err := f.s.Share(f.ctx, "pastas", p.ID, shared.GroupSubject("contab"), f.role("pastas", "Viewer"), "ana", origin); err != nil {
		t.Fatal(err)
	}
	f.confirmLast()
	if inv, _ := f.s.Invites(f.ctx, "caio"); len(inv) != 0 || f.roleName("caio", "pastas", p.ID) != "Viewer" {
		t.Error("the group's at once")
	}
	f.users.Join("dani", "contab")
	_ = f.s.MemberChanged(f.ctx, "contab", "dani", true, "ana", origin)
	if f.roleName("dani", "pastas", p.ID) != "Viewer" {
		t.Error("a new member")
	}
	f.users.Leave("caio", "contab")
	_ = f.s.MemberChanged(f.ctx, "contab", "caio", false, "ana", origin)
	if f.roleName("caio", "pastas", p.ID) != "" {
		t.Error("one who left")
	}
	// dani: estag (None) before contab (Viewer) in their order → None
	f.users.Leave("dani", "contab")
	f.users.Join("dani", "estag")
	f.users.Join("dani", "contab")
	f.share("pastas", p.ID, shared.GroupSubject("estag"), shared.RoleNone, "ana")
	if e, _ := f.s.RoleOf(f.ctx, "dani", "pastas", p.ID); e.Visible() {
		t.Error("the first group in the user's order")
	}
	// the user's own access before the groups'
	f.share("pastas", p.ID, shared.UserSubject("dani"), "Editor", "ana")
	if f.roleName("dani", "pastas", p.ID) != "Editor" {
		t.Error("the user's own first")
	}
	got := f.events("pastas", p.ID)
	if !slices.Contains(got, shared.EventMemberAdded) || !slices.Contains(got, shared.EventMemberRemoved) {
		t.Errorf("history %v", got)
	}
}

// A new owner: asked to the owner by e-mail, accepted by the new one; the
// former keeps the role given.
func TestOwnerTransfer(t *testing.T) {
	f := newFixture(t)
	p := sharedtest.NewPasta(f.db, "P", "ana")
	if _, err := f.s.TransferOwner(f.ctx, "pastas", p.ID, "ana", f.role("pastas", "Editor"), "ana", origin); !errors.Is(err, shared.ErrOwnerIsUser) {
		t.Errorf("to themselves: %v", err)
	}
	if _, err := f.s.TransferOwner(f.ctx, "pastas", p.ID, "bia", f.role("pastas", "Editor"), "ana", origin); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(f.last().Subject, "propriet") {
		t.Errorf("the mail: %q", f.last().Subject)
	}
	f.confirmLast()
	if o, _ := f.s.Owner(f.ctx, "pastas", p.ID); o != "ana" {
		t.Error("before the new owner accepts")
	}
	f.acceptAll("bia")
	if o, _ := f.s.Owner(f.ctx, "pastas", p.ID); o != "bia" {
		t.Errorf("owner %q", o)
	}
	if f.roleName("ana", "pastas", p.ID) != "Editor" || f.roleName("bia", "pastas", p.ID) != "owner" {
		t.Error("the former owner's role")
	}
	if got := f.events("pastas", p.ID); !slices.Contains(got, shared.EventOwnerTransferred) {
		t.Errorf("history %v", got)
	}
}
