package gitedit

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	fs_tools "github.com/go-rvq/rvq/admin/packages/fs-tools"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// testUsers are the users of the tests, by key.
type testUsers map[string]*User

func (u testUsers) User(_ *http.Request, key string) (*User, error) {
	if x, ok := u[key]; ok {
		return x, nil
	}
	return nil, errors.New("no such user")
}

func (u testUsers) Search(*http.Request, string, int) (out []*User, _ error) {
	for _, x := range u {
		out = append(out, x)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return
}

var people = testUsers{
	"ana":  {Key: "ana", Login: "ana", Name: "Ana Lima", Author: Author{"Ana Lima", "ana@x"}},
	"bia":  {Key: "bia", Login: "bia", Name: "Bia Souza", Author: Author{"Bia Souza", "bia@x"}},
	"caio": {Key: "caio", Login: "caio", Name: "Caio", Author: Author{"Caio", "caio@x"}},
	"root": {Key: "root", Login: "root", Name: "Root", Author: Author{"Root", "root@x"}},
}

// shareApp is the editor with users (X-User, ana by default), sharings and
// git: whoever is "root" may see every draft (!drafts), nobody else.
type shareApp struct {
	h   http.Handler
	b   *Builder
	git func(name string) string
	srv string
}

func newShareApp(t *testing.T) *shareApp {
	t.Helper()
	repo, _ := site(t)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	userOf := func(r *http.Request) string {
		if u := r.Header.Get("X-User"); u != "" {
			return u
		}
		return "ana"
	}
	pb := perm.New().AllowAll().SubjectsFunc(func(r *http.Request) []string { return []string{userOf(r)} })
	for _, u := range []string{"ana", "bia", "caio"} {
		pb.CreatePolicies(perm.PolicyFor(u).WhoAre(perm.Denied).ToDo(perm.Anything).On("admin:/site-files:!drafts"))
	}
	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.DataOperator(gorm2op.DataOperator(db))
	p.Permission(pb)
	b := &Builder{Repo: repo, Users: people, DB: db, PreviewPath: "/site-preview",
		Identity: func(r *http.Request) (string, Author) { u := people[userOf(r)]; return u.Key, u.Author },
		Preview: func(w http.ResponseWriter, r *http.Request, d *Draft, prefix string) {
			w.Write([]byte("the draft of " + d.Key + " at " + prefix))
		},
	}
	if err := b.Install(p); err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	p.Build(mux)

	gp := presets.New(i18n.New()).URIPrefix("/admin")
	fsb, err := fs_tools.New(gp, gp.I18n(), nil)
	if err != nil {
		t.Fatal(err)
	}
	fsb.AddGitRepo(b.GitRepo("site-files")).AddGitRepo(b.DraftGitRepo("site-files-draft")).
		AddGitRepo(b.DraftsGitRepo("site-files-drafts/" + fs_tools.GitKeyPattern))
	srv := httptest.NewServer(fsb.GitHandler())
	t.Cleanup(srv.Close)
	return &shareApp{h: mux, b: b, srv: srv.URL, git: func(name string) string { return srv.URL + fsb.GitPath(name) }}
}

// get is path served to user: its status and body.
func (a *shareApp) get(user, path string) (int, string) {
	r := httptest.NewRequest("GET", path, nil)
	r.Header.Set("X-User", user)
	w := httptest.NewRecorder()
	a.h.ServeHTTP(w, r)
	return w.Code, w.Body.String()
}

// reaches says whether user reaches the draft of owner: its page, its IDE,
// its editor, its preview.
func (a *shareApp) reaches(t *testing.T, user, owner string) map[string]bool {
	t.Helper()
	got := map[string]bool{}
	if code, body := a.get(user, "/admin/site-files/u/"+owner); code == 200 {
		got["page"] = strings.Contains(body, "data-draft-owner")
	}
	code, body := a.get(user, "/admin/site-files/u/"+owner+"/ide/api/ide/tree")
	got["ide"] = code == 200 && strings.Contains(body, "index.gadx")
	code, body = a.get(user, "/admin/site-files/u/"+owner+"/editor")
	got["editor"] = code == 200 && strings.Contains(body, "vx-gad-ide src=")
	code, body = a.get(user, "/admin/site-preview/u/"+owner+"/x")
	got["preview"] = code == 200 && body == "the draft of "+owner+" at /admin/site-preview/u/"+owner
	return got
}

func allTrue(m map[string]bool) bool {
	for _, v := range m {
		if !v {
			return false
		}
	}
	return len(m) == 4
}

func noneTrue(m map[string]bool) bool {
	for _, v := range m {
		if v {
			return false
		}
	}
	return true
}

// A draft shared: the user reaches it whole — the page, the IDE, the editor,
// the preview, git — until revoked or expired; someone else never; whoever
// may see every draft (!drafts) always; its owner's page lists the sharings,
// the user's the drafts shared with them.
func TestShareDraft(t *testing.T) {
	a := newShareApp(t)
	ctx := context.Background()
	a.get("ana", "/admin/site-files") // ana's draft made

	if got := a.reaches(t, "bia", "ana"); !noneTrue(got) {
		t.Errorf("bia reaches ana's draft not shared: %v", got)
	}
	if got := a.reaches(t, "root", "ana"); !allTrue(got) {
		t.Errorf("root (!drafts) does not reach ana's: %v", got)
	}
	s, err := a.b.Share(ctx, "ana", "bia", "ana", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := a.reaches(t, "bia", "ana"); !allTrue(got) {
		t.Errorf("bia does not reach ana's shared: %v", got)
	}
	if got := a.reaches(t, "caio", "ana"); !noneTrue(got) {
		t.Errorf("caio reaches ana's: %v", got)
	}
	// a key of no user: nothing, not even to root (no draft made of it)
	if code, _ := a.get("root", "/admin/site-files/u/nobody/ide/api/ide/tree"); code == 200 {
		t.Error("the draft of no user reached")
	}
	if _, err := os.Stat(filepath.Join(a.b.Repo.DraftsDir, "nobody")); err == nil {
		t.Error("a draft made for no user")
	}

	// the pages: ana's lists bia; bia's lists ana's draft; root's every draft
	if _, body := a.get("ana", "/admin/site-files"); !strings.Contains(body, "data-share=") || !strings.Contains(body, "Bia Souza (bia)") {
		t.Error("ana's page does not list the sharing")
	}
	if _, body := a.get("bia", "/admin/site-files"); !strings.Contains(body, "data-shared-draft=&#39;ana&#39;") {
		t.Error("bia's page does not list ana's draft")
	}
	if _, body := a.get("root", "/admin/site-files"); !strings.Contains(body, "data-draft=&#39;ana&#39;") {
		t.Error("root's page does not list every draft")
	}
	if _, body := a.get("bia", "/admin/site-files"); strings.Contains(body, "data-draft=&#39;ana&#39;") {
		t.Error("bia sees every draft")
	}

	// by git: bia clones ana's draft and pushes to it
	root := t.TempDir()
	clone := filepath.Join(root, "clone")
	hdr := []string{"-c", "http.extraHeader=X-User: bia"}
	mustGit(t, root, append(hdr, "clone", "-q", a.git("site-files-drafts/ana"), clone)...)
	w := httptest.NewRecorder()
	hr := httptest.NewRequest("GET", a.srv+"/admin/site-files/commit-msg", nil)
	hr.Header.Set("X-User", "bia")
	a.b.serveHook(w, hr)
	os.WriteFile(filepath.Join(clone, ".git", "hooks", "commit-msg"), w.Body.Bytes(), 0o755)
	os.WriteFile(filepath.Join(clone, "index.gadx"), []byte("by bia\n"), 0o644)
	mustGit(t, clone, "commit", "-q", "-am", "by bia")
	mustGit(t, clone, append(hdr, "push", "-q", "origin", "main")...)
	if got, _ := os.ReadFile(filepath.Join(a.b.Repo.DraftsDir, "ana", "index.gadx")); string(got) != "by bia\n" {
		t.Errorf("ana's draft has %q", got)
	}
	// caio: nothing by git
	if out, err := gitClient(root, "-c", "http.extraHeader=X-User: caio", "clone", "-q",
		a.git("site-files-drafts/ana"), filepath.Join(root, "caio")); err == nil {
		t.Errorf("caio cloned ana's: %s", out)
	}

	// revoked: bia reaches nothing, by git neither; its history kept
	if err := a.b.Revoke(ctx, s.ID, "ana"); err != nil {
		t.Fatal(err)
	}
	if got := a.reaches(t, "bia", "ana"); !noneTrue(got) {
		t.Errorf("bia reaches ana's revoked: %v", got)
	}
	os.WriteFile(filepath.Join(clone, "index.gadx"), []byte("again\n"), 0o644)
	mustGit(t, clone, "commit", "-q", "-am", "again")
	if out, err := gitClient(clone, append(hdr, "push", "-q", "origin", "main")...); err == nil {
		t.Errorf("bia pushed revoked: %s", out)
	}
	if shares, _ := a.b.SharesOf(ctx, "ana"); len(shares) != 1 || shares[0].RevokedAt == nil || shares[0].RevokedBy != "ana" {
		t.Errorf("the history: %+v", shares)
	}

	// expired: nothing
	past := time.Now().Add(-time.Minute)
	if _, err := a.b.Share(ctx, "ana", "bia", "ana", &past); err != nil {
		t.Fatal(err)
	}
	if got := a.reaches(t, "bia", "ana"); !noneTrue(got) {
		t.Errorf("bia reaches ana's expired: %v", got)
	}
	future := time.Now().Add(time.Hour)
	if _, err := a.b.Share(ctx, "ana", "bia", "ana", &future); err != nil {
		t.Fatal(err)
	}
	if got := a.reaches(t, "bia", "ana"); !allTrue(got) {
		t.Errorf("bia does not reach ana's until later: %v", got)
	}
	if _, err := a.b.Share(ctx, "ana", "ana", "ana", nil); !errors.Is(err, ErrShareSelf) {
		t.Errorf("shared with its owner: %v", err)
	}
}

// A commit of the panel in another's draft: by whoever commits, its message
// ending in the other users who changed its files (Co-authored-by), who made
// it on the site (Site-User) and where (Committed-Via).
func TestPanelCommitSigned(t *testing.T) {
	a := newShareApp(t)
	ctx := context.Background()
	if _, err := a.b.Share(ctx, "ana", "bia", "ana", nil); err != nil {
		t.Fatal(err)
	}
	// ana and bia change files by the IDE, in ana's draft
	write := func(user, base, path, content string) {
		r := httptest.NewRequest("PUT", base+"/ide/api/ide/file",
			strings.NewReader(`{"path":"`+path+`","content":"`+content+`"}`))
		r.Header.Set("X-User", user)
		w := httptest.NewRecorder()
		a.h.ServeHTTP(w, r)
		if w.Code != 200 {
			t.Fatalf("%s writing %s: %d %s", user, path, w.Code, w.Body.String())
		}
	}
	write("ana", "/admin/site-files", "a.gadx", "by ana")
	write("bia", "/admin/site-files/u/ana", "b.gadx", "by bia")

	// bia commits, on ana's draft
	r := httptest.NewRequest("POST", "http://site.example/admin/site-files/u/ana", nil)
	r.Header.Set("X-User", "bia")
	r.SetPathValue(userParam, "ana")
	d, author, err := a.b.Draft(r)
	if err != nil {
		t.Fatal(err)
	}
	if author.Email != "bia@x" {
		t.Errorf("the author: %+v", author)
	}
	if _, err := d.Commit(ctx, a.b.panelMessage(r, d, "two files"), author); err != nil {
		t.Fatal(err)
	}
	out, _ := git(ctx, d.Dir, "log", "-1", "--format=%an <%ae>%n%B")
	for _, want := range []string{"Bia Souza <bia@x>", "two files", "Co-authored-by: Ana Lima <ana@x>",
		"Site-User: bia <bia@site.example>", "Committed-Via: admin panel (site.example)"} {
		if !strings.Contains(out, want) {
			t.Errorf("the commit: no %q\n%s", want, out)
		}
	}
	if strings.Contains(out, "Co-authored-by: Bia") {
		t.Errorf("the committer as a co-author:\n%s", out)
	}
	if len(d.Editors()) != 0 {
		t.Errorf("the editors after the commit: %v", d.Editors())
	}
}

// A push of commits that do not say who made them on the site — no
// Site-User, one of another host, of a user the site does not know, a login
// not the key's — is refused, the line to add and the hook said.
func TestPushSigned(t *testing.T) {
	a := newShareApp(t)
	root := t.TempDir()
	clone := filepath.Join(root, "clone")
	mustGit(t, root, "clone", "-q", a.git("site-files-draft"), clone)
	host := strings.Split(strings.TrimPrefix(a.srv, "http://"), ":")[0]
	for _, c := range []struct{ name, trailer string }{
		{"none", ""},
		{"another host", "Site-User: ana <ana@elsewhere.example>"},
		{"no such user", "Site-User: zoe <zoe@" + host + ">"},
		{"not the key's login", "Site-User: bia <ana@" + host + ">"},
	} {
		os.WriteFile(filepath.Join(clone, "index.gadx"), []byte(c.name+"\n"), 0o644)
		msg := c.name
		if c.trailer != "" {
			msg += "\n\n" + c.trailer
		}
		mustGit(t, clone, "commit", "-q", "-am", msg)
		out, err := gitClient(clone, "push", "origin", "main")
		if err == nil || !strings.Contains(out, "does not say who made it on "+host) ||
			!strings.Contains(out, "Site-User: ana <ana@"+host+">") || !strings.Contains(out, "/admin/site-files/commit-msg") {
			t.Errorf("%s: pushed or not said: %v\n%s", c.name, err, out)
		}
		mustGit(t, clone, "reset", "-q", "--hard", "HEAD~1")
	}
	os.WriteFile(filepath.Join(clone, "index.gadx"), []byte("signed\n"), 0o644)
	mustGit(t, clone, "commit", "-q", "-am", "signed\n\nSite-User: ana <ana@"+host+">")
	mustGit(t, clone, "push", "-q", "origin", "main")
}

// The actions of the page: ana shares her draft (with whom, until when) and
// revokes it; bia cannot revoke it, nor share ana's; root (!drafts) revokes
// any.
func TestShareActions(t *testing.T) {
	a := newShareApp(t)
	ctx := context.Background()
	act := func(user, page, action string, query, form [][2]string) *httptest.ResponseRecorder {
		mb := multipartestutils.NewMultipartBuilder().PageURL(page).EventFunc(actions.DoAction).
			Query(presets.ParamAction, action)
		for _, q := range query {
			mb = mb.Query(q[0], q[1])
		}
		for _, f := range form {
			mb = mb.AddField(f[0], f[1])
		}
		r := mb.BuildEventFuncRequest()
		r.Header.Set("X-User", user)
		w := httptest.NewRecorder()
		a.h.ServeHTTP(w, r)
		return w
	}
	active := func() []*DraftShare {
		shares, _ := a.b.SharesOf(ctx, "ana")
		var out []*DraftShare
		for _, s := range shares {
			if s.Active(time.Now()) {
				out = append(out, s)
			}
		}
		return out
	}

	until := time.Now().AddDate(0, 0, 3).Format("2006-01-02")
	act("ana", "/admin/site-files", ActionShare, nil, [][2]string{{"User", "bia"}, {"ExpiresAt", until}})
	got := active()
	if len(got) != 1 || got[0].User != "bia" || got[0].CreatedBy != "ana" || got[0].ExpiresAt == nil ||
		got[0].ExpiresAt.Format("2006-01-02") != time.Now().AddDate(0, 0, 4).Format("2006-01-02") {
		t.Fatalf("shared: %+v", got)
	}
	// no user, a past date, herself: refused
	for _, form := range [][][2]string{{{"User", ""}}, {{"User", "caio"}, {"ExpiresAt", "2001-01-01"}}, {{"User", "ana"}}} {
		act("ana", "/admin/site-files", ActionShare, nil, form)
	}
	if n := len(active()); n != 1 {
		t.Errorf("a sharing refused was made: %d", n)
	}
	// bia, on ana's draft: cannot share it, nor revoke
	act("bia", "/admin/site-files/u/ana", ActionShare, nil, [][2]string{{"User", "caio"}})
	act("bia", "/admin/site-files/u/ana", ActionRevoke, [][2]string{{"id", got[0].ID.String()}}, nil)
	if g := active(); len(g) != 1 || g[0].User != "bia" {
		t.Errorf("bia shared or revoked ana's: %+v", g)
	}
	// ana revokes
	act("ana", "/admin/site-files", ActionRevoke, [][2]string{{"id", got[0].ID.String()}}, nil)
	if len(active()) != 0 {
		t.Error("not revoked")
	}
	// root revokes any
	s, _ := a.b.Share(ctx, "ana", "caio", "ana", nil)
	act("root", "/admin/site-files/u/ana", ActionRevoke, [][2]string{{"id", s.ID.String()}}, nil)
	if len(active()) != 0 {
		t.Error("root did not revoke")
	}
}
