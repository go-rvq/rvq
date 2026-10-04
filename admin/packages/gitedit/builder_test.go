package gitedit

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	rvqjs "github.com/go-rvq/rvq/js"

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
	pb := perm.New().AllowAll()
	pb.CreatePolicies(policies...)
	h, repo, served, _ := appWith(t, invalid, pb)
	return h, repo, served
}

// appWith is the admin with the editor of the site's files, by the
// permissions pb; its builder too.
func appWith(t *testing.T, invalid *error, pb *perm.Builder, group ...string) (http.Handler, *Repo, string, *Builder) {
	t.Helper()
	repo, served := site(t)
	p := presets.New(i18n.New()).URIPrefix("/admin")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	p.DataOperator(gorm2op.DataOperator(db))
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
	for _, g := range group {
		b.Page().Page().SetMenuGroup(p.MenuGroup(g))
	}
	mux := http.NewServeMux()
	p.Build(mux)
	return mux, repo, served, b
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
	if w.Code != 200 || !strings.Contains(w.Body.String(), "/admin/site-files/ide/index.html") || !strings.Contains(w.Body.String(), "Open the editor") {
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
	// an image of the draft, as it is: shown by the IDE
	os.WriteFile(filepath.Join(draft, "logo.png"), []byte("\x89PNG\r\n\x1a\n"), 0o644)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/site-files/ide/api/ide/file?raw=1&path=logo.png", nil))
	if w.Code != 200 || w.Header().Get("Content-Type") != "image/png" {
		t.Errorf("the IDE's image: %d %q", w.Code, w.Header().Get("Content-Type"))
	}
	os.Remove(filepath.Join(draft, "logo.png"))
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
	// the app's assets: served with no permission (the login lets static
	// files by with no session); its page asks it
	if w := httptest.NewRecorder(); func() int {
		h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/site-files/ide/index.html", nil))
		i := strings.Index(w.Body.String(), "./assets/")
		asset := w.Body.String()[i+2 : i+strings.IndexByte(w.Body.String()[i:], '"')]
		w2 := httptest.NewRecorder()
		h.ServeHTTP(w2, httptest.NewRequest("GET", "/admin/site-files/ide/"+asset, nil))
		return w2.Code
	}() != http.StatusOK {
		t.Error("an asset of the IDE's app not served")
	}
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

// The app's assets are served with no permission — the login lets static
// files by with no session, no user to ask one of —; its page asks @get.
func TestBuilderIDEAssets(t *testing.T) {
	var invalid error
	h, _, _ := app(t, &invalid, perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(perm.Anything).On("admin:/site-files:@get"))
	index, err := fs.ReadFile(rvqjs.GadIDE(), "index.html")
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`src="\./(assets/[^"]+)"`).FindSubmatch(index)
	if m == nil {
		t.Fatal("no asset in the index")
	}
	for path, want := range map[string]int{
		"/admin/site-files/ide/index.html":      http.StatusForbidden,
		"/admin/site-files/ide/" + string(m[1]): http.StatusOK,
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != want {
			t.Errorf("%s: %d, want %d", path, w.Code, want)
		}
	}
}

// The permissions of a path, recursive ("<static/*>"), with the page's: a deny
// of the path denies; an allow of the page allows; a deny of the page denies,
// whatever the path; with nothing of the page, an allow of the path allows.
func TestAllowedOpPaths(t *testing.T) {
	var invalid error
	r := httptest.NewRequest("GET", "/admin/site-files", nil)
	policy := func(effect, res string) *perm.PolicyBuilder {
		return perm.PolicyFor(perm.Anybody).WhoAre(effect).ToDo(perm.Anything).On(res)
	}

	// everything allowed, config/ denied for editing; creating denied on the
	// page, allowed under static/ — the page's deny wins
	pb := perm.New().AllowAll()
	pb.CreatePolicies(policy(perm.Denied, "admin:/site-files:<config/*>:!edit"),
		policy(perm.Denied, "admin:/site-files:!create"), policy(perm.Allowed, "admin:/site-files:<static/*>:!create"))
	_, _, _, b := appWith(t, &invalid, pb)
	for _, c := range []struct {
		op   IdeOp
		path string
		want bool
	}{
		{IdeEdit, "config/layout_config.gad", false},
		{IdeEdit, "config/sub/deep.gad", false}, // recursive
		{IdeEdit, "templates/main.gad", true},
		{IdeEdit, "", true},
		{IdeCreate, "static/new.css", false},
		{IdeCreate, "templates/new.gadx", false},
	} {
		if got := b.AllowedOp(r, c.op, c.path); got != c.want {
			t.Errorf("%s %q: %v", opPerms[c.op], c.path, got)
		}
	}

	// only under static/: nothing of the page, editing allowed of static/
	pb = perm.New()
	pb.CreatePolicies(policy(perm.Allowed, "admin:/site-files:@get"),
		policy(perm.Allowed, "admin:/site-files:<static/*>:!edit"))
	_, _, _, b = appWith(t, &invalid, pb)
	for p, want := range map[string]bool{"static/css/a.css": true, "static/x": true, "templates/main.gad": false, "": false} {
		if got := b.AllowedOp(r, IdeEdit, p); got != want {
			t.Errorf("only static/: %q: %v", p, got)
		}
	}
	if b.AllowedOp(r, IdeDelete, "static/css/a.css") {
		t.Error("only static/: deleting allowed with editing")
	}
	if !b.AllowedOp(r, IdeRead, "templates/main.gad") {
		t.Error("only static/: reading denied")
	}
}

// The WebDAV on the draft asks what the IDE asks, by its method.
func TestDavAllowed(t *testing.T) {
	var invalid error
	pb := perm.New()
	allow := func(res string) *perm.PolicyBuilder {
		return perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On(res)
	}
	pb.CreatePolicies(allow("admin:/site-files:@get"), allow("admin:/site-files:!edit"), allow("admin:/site-files:!rename"))
	_, repo, _, b := appWith(t, &invalid, pb)
	if _, err := repo.Draft(context.Background(), "u1"); err != nil {
		t.Fatal(err)
	}
	req := func(method string) *http.Request { return httptest.NewRequest(method, "/x", nil) }
	for _, c := range []struct {
		method, p, dest string
		write, want     bool
	}{
		{"GET", "/index.gadx", "", false, true},
		{"PUT", "/index.gadx", "", true, true}, // there: !edit
		{"PUT", "/new.gadx", "", true, false},  // not there: !create
		{"MKCOL", "/dir", "", true, false},     // !create
		{"DELETE", "/index.gadx", "", true, false},
		{"MOVE", "/index.gadx", "/main.gadx", true, true},       // in its folder: !rename
		{"MOVE", "/index.gadx", "/sub/index.gadx", true, false}, // to another: !move
		{"COPY", "/index.gadx", "/copy.gadx", true, false},      // !create
		{"PROPPATCH", "/index.gadx", "", true, true},            // !edit
	} {
		if got := b.DavAllowed(req(c.method), c.p, c.dest, c.write); got != c.want {
			t.Errorf("%s %s %s: %v", c.method, c.p, c.dest, got)
		}
	}
}

// The examples of the documentation (user_docs …/05-permissions, "Examples by
// path"), each its policies on the page's resource out of its group, do what
// it says — the page in a group of the menu.
func TestPermissionExamples(t *testing.T) {
	var invalid error
	r := httptest.NewRequest("GET", "/admin/site-files", nil)
	// the page in a group, the policies out of it — by its unique name, what
	// admin.page("/site-files").uniqueResource is
	const res = "admin:/site-files:"
	type rule struct{ effect, perm string }
	type check struct {
		op   IdeOp
		path string
		want bool
	}
	for _, ex := range []struct {
		name   string
		rules  []rule
		checks []check
	}{
		{"only the styles", []rule{{perm.Allowed, "@get"}, {perm.Allowed, "<static/css/*>:!edit"}, {perm.Allowed, "<static/css/*>:!create"}},
			[]check{{IdeEdit, "static/css/site.css", true}, {IdeCreate, "static/css/new.css", true}, {IdeEdit, "static/js/a.js", false},
				{IdeEdit, "templates/main.gad", false}, {IdeDelete, "static/css/site.css", false}, {IdeRead, "templates/main.gad", true}}},
		{"all but config, read only", []rule{{perm.Allowed, "*"}, {perm.Denied, "<config/*>:!*"}},
			[]check{{IdeEdit, "templates/main.gad", true}, {IdeDelete, "static/a.css", true}, {IdeEdit, "config/layout_config.gad", false},
				{IdeCreate, "config/new.gad", false}, {IdeDelete, "config/layout_config.gad", false}, {IdeRead, "config/layout_config.gad", true}}},
		{"images into static/img", []rule{{perm.Allowed, "@get"}, {perm.Allowed, "<static/img/*>:!import"},
			{perm.Allowed, "<static/img/*>:!create"}, {perm.Allowed, "<static/img/*>:!edit"}},
			[]check{{IdeImport, "static/img/logo.png", true}, {IdeCreate, "static/img/logo.png", true}, {IdeImport, "static/css/a.css", false},
				{IdeImport, "", false}}},
		{"layouts never deleted nor moved", []rule{{perm.Allowed, "*"}, {perm.Denied, "<templates/layouts/*>:!delete"},
			{perm.Denied, "<templates/layouts/*>:!move"}, {perm.Denied, "<templates/layouts/*>:!rename"}},
			[]check{{IdeDelete, "templates/layouts/default.gadx", false}, {IdeMove, "templates/layouts/default.gadx", false},
				{IdeRename, "templates/layouts/default.gadx", false}, {IdeEdit, "templates/layouts/default.gadx", true},
				{IdeDelete, "templates/main.gad", true}}},
		{"rename and move inside static", []rule{{perm.Allowed, "@get"}, {perm.Allowed, "<static/*>:!rename"}, {perm.Allowed, "<static/*>:!move"}},
			[]check{{IdeMove, "static/a.css", true}, {IdeMove, "static/css/a.css", true}, {IdeMove, "templates/a.css", false},
				{IdeRename, "static/a.css", true}, {IdeRename, "config/x.gad", false}}},
		{"a single file", []rule{{perm.Allowed, "@get"}, {perm.Allowed, "<config/layout_config.gad>:!edit"}},
			[]check{{IdeEdit, "config/layout_config.gad", true}, {IdeEdit, "config/other.gad", false}}},
	} {
		pb := perm.New()
		for _, ru := range ex.rules {
			pb.CreatePolicies(perm.PolicyFor(perm.Anybody).WhoAre(ru.effect).ToDo(perm.Anything).On(res + ru.perm))
		}
		_, _, _, b := appWith(t, &invalid, pb, "site")
		for _, c := range ex.checks {
			if got := b.AllowedOp(r, c.op, c.path); got != c.want {
				t.Errorf("%s: %s %q: %v", ex.name, opPerms[c.op], c.path, got)
			}
		}
	}
	// a move asks both paths: out of static/ is denied
	pb := perm.New()
	pb.CreatePolicies(perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On(res+"<static/*>:!move"),
		perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On(res+"@get"))
	_, _, _, b := appWith(t, &invalid, pb, "site")
	needs := asks("rename", &ideBody{Path: "static/a.css", To: "templates/a.css"}, "")
	allowed := true
	for _, n := range needs {
		allowed = allowed && b.AllowedOp(r, n.Op, n.Path)
	}
	if allowed {
		t.Error("a move out of static/ allowed")
	}
}

// In a group of the menu the page has two resources: by the group
// ("admin:site/:/site-files:") and by its unique name, out of the group
// ("admin:/site-files:"), which decides first — for the page's actions and
// for the paths alike.
func TestPermissionsByUniqueName(t *testing.T) {
	var invalid error
	r := httptest.NewRequest("GET", "/admin/site/site-files", nil)
	const group, unique = "admin:site/:/site-files:", "admin:/site-files:"
	policy := func(effect, res string) *perm.PolicyBuilder {
		return perm.PolicyFor(perm.Anybody).WhoAre(effect).ToDo(perm.Anything).On(res)
	}
	_, _, _, b := appWith(t, &invalid, perm.New(), "site")
	v := b.Page().Page().ActionVerifier(r, "!edit")
	if v.Resource() != group || v.PreferredResource() != unique {
		t.Fatalf("the resources: %q, %q", v.Resource(), v.PreferredResource())
	}
	for _, c := range []struct {
		name     string
		policies []*perm.PolicyBuilder
		op       IdeOp
		path     string
		want     bool
	}{
		{"the group allows", []*perm.PolicyBuilder{policy(perm.Allowed, group+"!edit")}, IdeEdit, "a.css", true},
		{"the unique name denies over the group", []*perm.PolicyBuilder{policy(perm.Allowed, group+"!edit"),
			policy(perm.Denied, unique+"!edit")}, IdeEdit, "a.css", false},
		{"the unique name allows over the group", []*perm.PolicyBuilder{policy(perm.Denied, group+"!edit"),
			policy(perm.Allowed, unique+"!edit")}, IdeEdit, "a.css", true},
		{"a path denied by the unique name, the group allowing", []*perm.PolicyBuilder{policy(perm.Allowed, group+"*"),
			policy(perm.Denied, unique+"<config/*>:!edit")}, IdeEdit, "config/x.gad", false},
		{"a path allowed by the unique name, nothing of the page", []*perm.PolicyBuilder{
			policy(perm.Allowed, unique+"<static/*>:!edit")}, IdeEdit, "static/a.css", true},
		{"a path allowed by the group, the unique name denying the page", []*perm.PolicyBuilder{
			policy(perm.Allowed, group+"<static/*>:!edit"), policy(perm.Denied, unique+"!edit")}, IdeEdit, "static/a.css", false},
	} {
		pb := perm.New()
		pb.CreatePolicies(c.policies...)
		_, _, _, b := appWith(t, &invalid, pb, "site")
		if got := b.AllowedOp(r, c.op, c.path); got != c.want {
			t.Errorf("%s: %v", c.name, got)
		}
	}
}
