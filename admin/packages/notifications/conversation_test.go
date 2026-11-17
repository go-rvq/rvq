package notifications

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/packages/comments"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/x/login"
	"github.com/google/uuid"
)

// fakeUser is a minimal user.User.
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
func (f *fakeUser) SetID(v uuid.UUID)             { f.id = v }
func (f *fakeUser) Anonymous() bool               { return false }
func (f *fakeUser) GetAccountName() string        { return "test" }

var _ user.User = (*fakeUser)(nil)

// A message of a user to others: a copy each, one message, its files; the
// answer of one of them goes to everyone else of the conversation, under
// "Re: "; an answer by mail kept once; its files seen by who is of it alone.
func TestDeliverAnswer(t *testing.T) {
	db := testDB(t)
	b := &Builder{db: db, FilesDir: t.TempDir()}
	ana, bia, caio, out := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	var mailed []uuid.UUID
	b.Mailer = func(_ context.Context, m *Mail) error { mailed = append(mailed, m.To); return nil }
	if err := b.Deliver(nil, ana, []uuid.UUID{bia, caio}, Notification{Title: "Reunião", Body: "<p>às <b>15h</b></p>", HTML: true},
		comments.Attachment{Name: "pauta.txt", Data: []byte("1. abertura")}); err != nil {
		t.Fatal(err)
	}
	var copies []Notification
	db.Order("user_id").Find(&copies)
	if len(copies) != 2 || copies[0].MessageID != copies[1].MessageID || copies[0].ThreadID != copies[0].MessageID ||
		copies[0].System || copies[0].Source != SourceUser(ana) || !copies[0].HTML {
		t.Fatalf("the copies: %+v", copies)
	}
	if len(mailed) != 2 {
		t.Fatalf("mailed: %v", mailed)
	}
	files, _ := FilesOf(db, copies[0].MessageID)
	if fs := files[copies[0].MessageID]; len(fs) != 1 || fs[0].Name != "pauta.txt" || fs[0].Size != 11 {
		t.Fatalf("the files: %+v", files)
	}
	// on disk, under the directory of the files; not in the database
	f0 := files[copies[0].MessageID][0]
	if data, err := os.ReadFile(filepath.Join(b.FilesDir, f0.Path)); err != nil || string(data) != "1. abertura" {
		t.Fatalf("the file on disk %q: %q %v", f0.Path, data, err)
	}
	var stored File
	db.First(&stored, "id = ?", f0.ID)
	if len(stored.Data) != 0 {
		t.Error("the data in the database too")
	}
	thread := copies[0].ThreadID

	// bia answers: to ana and caio
	if err := b.Answer(nil, bia, thread, Reply{Body: "ok", ExternalID: "<x@mail>"}); err != nil {
		t.Fatal(err)
	}
	var answers []Notification
	db.Where("source = ?", SourceUser(bia)).Find(&answers)
	got := map[uuid.UUID]bool{}
	for _, a := range answers {
		got[a.UserID] = true
		// none given: under the newest bia received, the first
		if a.Title != "Re: Reunião" || a.ThreadID != thread || a.ParentID != copies[0].MessageID {
			t.Fatalf("the answer: %+v", a)
		}
	}
	if len(answers) != 2 || !got[ana] || !got[caio] {
		t.Fatalf("answered to: %v", got)
	}
	if !b.Answered("<x@mail>") || b.Answered("<y@mail>") || b.Answered("") {
		t.Fatal("Answered")
	}
	// who is of none of it answers nothing
	if err := b.Answer(nil, out, thread, Reply{Body: "eu"}); err == nil {
		t.Fatal("an outsider answered")
	}

	// the files: to ana (sent), bia (received); not to out
	fid := files[copies[0].MessageID][0].ID
	h := FileHandler(db, b.FilesDir)
	get := func(u uuid.UUID) int {
		r := httptest.NewRequest("GET", "/admin/!notifications/files/"+fid.String(), nil)
		r.SetPathValue("id", fid.String())
		r = r.WithContext(login.WithTrustedUser(r.Context(), &fakeUser{id: u}))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Code
	}
	if get(ana) != http.StatusOK || get(bia) != http.StatusOK || get(out) != http.StatusNotFound {
		t.Fatalf("the file: ana %d bia %d out %d", get(ana), get(bia), get(out))
	}
}

// The attachments kept in the database before a directory was given: moved
// there, their data emptied, served from it.
func TestMoveFilesToDir(t *testing.T) {
	db := testDB(t)
	msg := uuid.New()
	if err := AttachFiles(db, "", msg, comments.Attachment{Name: "a.txt", Data: []byte("antes")}); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if n, err := MoveFilesToDir(db, dir); err != nil || n != 1 {
		t.Fatalf("moved %d: %v", n, err)
	}
	var f File
	db.First(&f, "message_id = ?", msg)
	if f.Path == "" || len(f.Data) != 0 {
		t.Fatalf("the record: %+v", f)
	}
	if data, _ := os.ReadFile(filepath.Join(dir, f.Path)); string(data) != "antes" {
		t.Errorf("on disk: %q", data)
	}
	if n, _ := MoveFilesToDir(db, dir); n != 0 {
		t.Errorf("again: %d", n)
	}
}
