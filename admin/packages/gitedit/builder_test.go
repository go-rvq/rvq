package gitedit

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

// app is the admin with the editor of the site's files, the policies given.
func app(t *testing.T, invalid *error, policies ...*perm.PolicyBuilder) (http.Handler, *Repo, string) {
	t.Helper()
	repo, served := site(t)
	p := presets.New(i18n.New()).URIPrefix("/admin")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	p.DataOperator(gorm2op.DataOperator(db))
	pb := perm.New().AllowAll()
	pb.CreatePolicies(policies...)
	p.Permission(pb.SubjectsFunc(func(*http.Request) []string { return []string{"editor"} }))
	b := &Builder{Repo: repo,
		Identity: func(*http.Request) (string, Author) { return "u1", Author{"Ana", "ana@x"} },
		Validate: func(context.Context, string) error { return *invalid },
		// the site of the draft: what it was given
		PreviewPath: "/site-preview",
		Preview: func(w http.ResponseWriter, r *http.Request, d *Draft, prefix string) {
			b, _ := os.ReadFile(filepath.Join(d.Dir, "index.gadx"))
			w.Write([]byte(prefix + " " + r.URL.Path + " " + string(b)))
		}}
	if err := b.Install(p); err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	p.Build(mux)
	return mux, repo, served
}

func doAction(h http.Handler, action string, query [][2]string, form [][2]string) *httptest.ResponseRecorder {
	b := multipartestutils.NewMultipartBuilder().PageURL("/admin/site-files").EventFunc(actions.DoAction).
		Query(presets.ParamAction, action)
	for _, q := range query {
		b = b.Query(q[0], q[1])
	}
	for _, f := range form {
		b = b.AddField(f[0], f[1])
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, b.BuildEventFuncRequest())
	return w
}

// The page: the draft edited by the IDE, committed with a message — refused
// when the files do not work —, a change discarded, published: the site
// changes only then.
func TestBuilder(t *testing.T) {
	var invalid error
	h, repo, served := app(t, &invalid)
	draft := filepath.Join(repo.DraftsDir, "u1")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/site-files", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "/admin/site-files/ide/") {
		t.Fatalf("the page: %d %.300s", w.Code, w.Body.String())
	}
	// the site of the draft, under the admin
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/site-preview/en-us/about", nil))
	if got := w.Body.String(); got != "/admin/site-preview /en-us/about old\n" {
		t.Errorf("the preview: %d %q", w.Code, got)
	}
	// the IDE's app, in the frame: its assets relative to the page
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/site-files/ide/index.html", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `src="./assets/`) {
		t.Errorf("the IDE's app: %d %.300s", w.Code, w.Body.String())
	}
	// the IDE, under the page
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("PUT", "/admin/site-files/ide/api/ide/file",
		strings.NewReader(`{"path":"index.gadx","content":"new\n"}`)))
	if w.Code != 200 {
		t.Fatalf("the IDE writing: %d %s", w.Code, w.Body.String())
	}
	os.WriteFile(filepath.Join(draft, "wip.txt"), []byte("wip"), 0o644)
	doAction(h, ActionDiscard, [][2]string{{"path", "wip.txt"}}, nil)
	if _, err := os.Stat(filepath.Join(draft, "wip.txt")); err == nil {
		t.Error("not discarded")
	}

	invalid = errors.New("index.gadx:1: unexpected token")
	doAction(h, ActionCommit, nil, [][2]string{{"Message", "the index"}})
	d, _ := repo.Draft(context.Background(), "u1")
	if log, _ := d.Log(context.Background(), 1); log[0].Subject == "the index" {
		t.Error("committed files that do not work")
	}
	invalid = nil
	doAction(h, ActionCommit, nil, [][2]string{{"Message", "the index"}})
	if log, _ := d.Log(context.Background(), 1); log[0].Subject != "the index" || log[0].Author != "Ana" {
		t.Fatalf("not committed: %+v", log)
	}
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "old\n" {
		t.Errorf("committed, the site changed: %q", b)
	}
	doAction(h, ActionPublish, nil, nil)
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "new\n" {
		t.Errorf("published, the site has %q", b)
	}
}

// Each action asks the page's permission of it; the IDE writes by "!edit".
func TestBuilderPermissions(t *testing.T) {
	var invalid error
	deny := func(res string) *perm.PolicyBuilder {
		return perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(perm.Anything).On(res)
	}
	h, repo, served := app(t, &invalid, deny("admin:/site-files:!publish"), deny("admin:/site-files:!edit"),
		deny("admin:/site-files:!preview"))
	if w := httptest.NewRecorder(); func() int {
		h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/site-preview/", nil))
		return w.Code
	}() != http.StatusForbidden {
		t.Error("the preview with no !preview")
	}
	draft := filepath.Join(repo.DraftsDir, "u1")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("PUT", "/admin/site-files/ide/api/ide/file",
		strings.NewReader(`{"path":"index.gadx","content":"new\n"}`)))
	if w.Code != http.StatusForbidden {
		t.Errorf("the IDE wrote with no !edit: %d", w.Code)
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/site-files/ide/api/ide/file?path=index.gadx", nil))
	if w.Code != 200 {
		t.Errorf("the IDE reading: %d", w.Code)
	}
	os.WriteFile(filepath.Join(draft, "index.gadx"), []byte("new\n"), 0o644)
	doAction(h, ActionCommit, nil, [][2]string{{"Message", "m"}})
	doAction(h, ActionPublish, nil, nil)
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "old\n" {
		t.Errorf("published with no !publish: %q", b)
	}
}
