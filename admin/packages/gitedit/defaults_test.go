package gitedit

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	fs_tools "github.com/go-rvq/rvq/admin/packages/fs-tools"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
)

// emptySite is a repository like the site's, with no commit yet (git init).
func emptySite(t *testing.T) (repo *Repo, served string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("no git")
	}
	root := t.TempDir()
	served = filepath.Join(root, "public")
	if _, err := git(context.Background(), root, "init", "--quiet", "--initial-branch=main", served); err != nil {
		t.Fatal(err)
	}
	return &Repo{GitDir: filepath.Join(served, ".git"), WorkTree: served, Branch: "main", DraftsDir: filepath.Join(root, "drafts")}, served
}

// An empty repository makes no draft (ErrEmpty), is clean, has no head;
// InitFiles gives it its first commit — refused by its commit func, nothing
// is there —, the site's files those; then it is not empty (ErrNotEmpty).
func TestInitFiles(t *testing.T) {
	ctx := context.Background()
	repo, served := emptySite(t)
	if empty, err := repo.Empty(ctx); err != nil || !empty {
		t.Fatalf("empty %v %v", empty, err)
	}
	if head, err := repo.Head(ctx); err != nil || head != "" {
		t.Errorf("head of the empty %q %v", head, err)
	}
	if err := repo.SiteClean(ctx); err != nil {
		t.Errorf("the empty is not clean: %v", err)
	}
	if _, err := repo.Draft(ctx, "u1"); !errors.Is(err, ErrEmpty) {
		t.Errorf("a draft of the empty: %v", err)
	}

	files := fstest.MapFS{
		"templates/main.gad": {Data: []byte("main\n")},
		"static/a.css":       {Data: []byte("a\n")},
	}
	refused := errors.New("no")
	if err := repo.InitFiles(ctx, files, func(*Draft) error { return refused }); !errors.Is(err, refused) {
		t.Errorf("refused: %v", err)
	}
	if empty, _ := repo.Empty(ctx); !empty {
		t.Error("refused, a commit")
	}
	if err := repo.InitFiles(ctx, files, func(d *Draft) error {
		_, err := d.Commit(ctx, "the defaults", Author{"Ana", "ana@x"})
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(served, "templates", "main.gad")); string(b) != "main\n" {
		t.Errorf("the site's files: %q", b)
	}
	if err := repo.SiteClean(ctx); err != nil {
		t.Errorf("not clean: %v", err)
	}
	if out, _ := git(ctx, served, "log", "--format=%an %s"); out != "Ana the defaults\n" {
		t.Errorf("log %q", out)
	}
	if entries, _ := os.ReadDir(repo.DraftsDir); len(entries) != 0 {
		t.Errorf("left in the drafts: %v", entries)
	}
	if err := repo.InitFiles(ctx, files, func(*Draft) error { return nil }); !errors.Is(err, ErrNotEmpty) {
		t.Errorf("again: %v", err)
	}
	if _, err := repo.Draft(ctx, "u1"); err != nil {
		t.Errorf("a draft: %v", err)
	}
}

// Overlay writes the defaults over the draft's files, the others kept — but
// in the folders it is told, where they go.
func TestOverlay(t *testing.T) {
	ctx := context.Background()
	repo, _ := site(t)
	d, err := repo.Draft(ctx, "u1")
	if err != nil {
		t.Fatal(err)
	}
	write := func(p, s string) {
		os.MkdirAll(filepath.Dir(filepath.Join(d.Dir, p)), 0o755)
		os.WriteFile(filepath.Join(d.Dir, p), []byte(s), 0o644)
	}
	write("templates/main.gad", "mine\n")
	write("templates/extra.gadx", "extra\n")
	write("static/mine.css", "mine\n")
	if _, err := d.Commit(ctx, "mine", Author{"Ana", "ana@x"}); err != nil {
		t.Fatal(err)
	}
	files := fstest.MapFS{"templates/main.gad": {Data: []byte("main\n")}, "static/a.css": {Data: []byte("a\n")}}
	read := func(p string) string { b, _ := os.ReadFile(filepath.Join(d.Dir, p)); return string(b) }

	if err := d.Overlay(files); err != nil {
		t.Fatal(err)
	}
	if read("templates/main.gad") != "main\n" || read("static/a.css") != "a\n" || read("templates/extra.gadx") != "extra\n" {
		t.Error("overlay, the others kept")
	}
	if err := d.Overlay(files, "templates"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(d.Dir, "templates", "extra.gadx")); err == nil {
		t.Error("an extra of templates kept")
	}
	if read("static/mine.css") != "mine\n" {
		t.Error("an extra of static gone")
	}
	if err := d.Overlay(files, ".git"); err == nil {
		t.Error("the .git of the draft")
	}
}

// The page of the empty repository: no draft, its first files by git or the
// defaults (InitDefaults, its permission); then the draft's page, whose
// ResetDefaults brings the defaults over the draft's files in a commit —
// the extras of DefaultsDirs gone when asked —, the site unchanged until
// published.
func TestDefaultsPage(t *testing.T) {
	var invalid error
	defaults := fstest.MapFS{
		".gitignore":         {Data: []byte(".*\n")},
		"index.gadx":         {Data: []byte("default\n")},
		"templates/main.gad": {Data: []byte("main\n")},
	}
	page := func(h http.Handler) string {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/site-files", nil))
		if w.Code != 200 {
			t.Fatalf("the page: %d %.300s", w.Code, w.Body.String())
		}
		return w.Body.String()
	}

	// denied: no button, the action refused
	repo, _ := emptySite(t)
	h, _ := appOn(t, repo, &invalid, allowAll(policyOn(perm.Denied, "admin:/site-files:!init_defaults")), defaults)
	if body := page(h); !strings.Contains(body, "data-empty-repo") || strings.Contains(body, "data-init-defaults") {
		t.Errorf("the empty page, denied: %.300s", body)
	}
	doAction(h, ActionInitDefaults, nil, nil)
	if empty, _ := repo.Empty(context.Background()); !empty {
		t.Error("denied, a commit")
	}

	repo, served := emptySite(t)
	h, _ = appOn(t, repo, &invalid, allowAll(), defaults)
	if body := page(h); !strings.Contains(body, "data-init-defaults") || !strings.Contains(body, "Empty repository") {
		t.Errorf("the empty page: %.300s", body)
	}
	// invalid: refused, nothing committed
	invalid = errors.New("broken")
	doAction(h, ActionInitDefaults, nil, nil)
	if empty, _ := repo.Empty(context.Background()); !empty {
		t.Error("invalid, a commit")
	}
	invalid = nil
	doAction(h, ActionInitDefaults, nil, nil)
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "default\n" {
		t.Fatalf("the site's files: %q", b)
	}
	if out, _ := git(context.Background(), served, "log", "--format=%an|%B"); !strings.Contains(out, "Ana|The default files of the site") ||
		!strings.Contains(out, TrailerSiteUser+": ") {
		t.Errorf("the commit: %q", out)
	}
	if body := page(h); strings.Contains(body, "data-empty-repo") {
		t.Error("still the empty page")
	}

	// the draft changed and published, then the defaults restored
	d, _ := repo.Draft(context.Background(), "u1")
	os.WriteFile(filepath.Join(d.Dir, "index.gadx"), []byte("mine\n"), 0o644)
	os.WriteFile(filepath.Join(d.Dir, "templates", "extra.gadx"), []byte("extra\n"), 0o644)
	os.WriteFile(filepath.Join(d.Dir, ".gitignore"), []byte("mine\n"), 0o644)
	doAction(h, ActionCommit, nil, [][2]string{{"Message", "mine"}})
	doAction(h, ActionPublish, nil, nil)
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "mine\n" {
		t.Fatalf("published: %q", b)
	}
	// the dialog: the changes browsed before; with the extras, one more
	ev := func(event string, q ...string) string {
		b := multipartestutils.NewMultipartBuilder().PageURL("/admin/site-files").EventFunc(event)
		for i := 0; i+1 < len(q); i += 2 {
			b = b.Query(q[i], q[i+1])
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, b.BuildEventFuncRequest())
		return w.Body.String()
	}
	if body := ev(actions.Action, presets.ParamAction, ActionResetDefaults); !strings.Contains(body, "vx-diff-browser") ||
		!strings.Contains(body, `data-defaults-changes`) || !strings.Contains(body, "index.gadx") {
		t.Errorf("the dialog: %.500s", body)
	}
	// the extras highlighted (warning)
	if body := ev(eventDefaultsPanel, "extras", "1"); !strings.Contains(body, "extra.gadx") || !strings.Contains(body, "warning") {
		t.Errorf("the panel with the extras: %.500s", body)
	}
	if body := ev(eventDefaultsPanel); strings.Contains(body, "extra.gadx") || strings.Contains(body, ".gitignore") ||
		strings.Contains(body, "warning") {
		t.Errorf("the panel without the extras: %.500s", body)
	}
	if body := ev(eventDefaultsDiff, "path", "index.gadx"); !strings.Contains(body, "mine") {
		t.Errorf("the diff: %.500s", body)
	}
	doAction(h, ActionResetDefaults, nil, [][2]string{{"RemoveExtras", "true"}})
	if b, _ := os.ReadFile(filepath.Join(d.Dir, "index.gadx")); string(b) != "default\n" {
		t.Errorf("the draft: %q", b)
	}
	if _, err := os.Stat(filepath.Join(d.Dir, "templates", "extra.gadx")); err == nil {
		t.Error("the extra kept")
	}
	if b, _ := os.ReadFile(filepath.Join(d.Dir, ".gitignore")); string(b) != "mine\n" {
		t.Errorf("the draft's .gitignore (DefaultsKeep): %q", b)
	}
	if changes, _ := d.Status(context.Background()); len(changes) != 0 {
		t.Errorf("not committed: %v", changes)
	}
	if out, _ := git(context.Background(), d.Dir, "log", "--format=%s"); out != "Restore the default files of the site\nmine\nThe default files of the site\n" {
		t.Errorf("the history: %q", out)
	}
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "mine\n" {
		t.Errorf("the site before publishing: %q", b)
	}
	doAction(h, ActionPublish, nil, nil)
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "default\n" {
		t.Errorf("the site: %q", b)
	}
}

// A push into the empty repository by git: its first commit, the site's
// files those.
func TestGitPushEmpty(t *testing.T) {
	var invalid error
	repo, served := emptySite(t)
	_, b := appOn(t, repo, &invalid, allowAll(), nil)
	p := presets.New(i18n.New()).URIPrefix("/admin")
	fsb, err := fs_tools.New(p, p.I18n(), nil)
	if err != nil {
		t.Fatal(err)
	}
	fsb.AddGitRepo(b.GitRepo("site-files"))
	srv := httptest.NewServer(fsb.GitHandler())
	t.Cleanup(srv.Close)
	w := httptest.NewRecorder()
	b.serveHook(w, httptest.NewRequest("GET", srv.URL+"/admin/site-files/commit-msg", nil))
	testHook = w.Body.String()
	t.Cleanup(func() { testHook = "" })

	clone := t.TempDir()
	mustGit(t, clone, "init", "-q", "-b", "main")
	commitFile(t, clone, "index.gadx", "first\n")
	mustGit(t, clone, "push", "-q", srv.URL+fsb.GitPath("site-files"), "HEAD:main")
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "first\n" {
		t.Errorf("the site's files: %q", b)
	}
	if err := repo.SiteClean(context.Background()); err != nil {
		t.Errorf("not clean: %v", err)
	}
}

// The changes of ResetDefaults, before it is done: each default file the
// draft has not (A) or has otherwise (M), the extras of the folders asked
// (D); the diff of each, the draft's and the default.
func TestDefaultsChanges(t *testing.T) {
	ctx := context.Background()
	repo, _ := site(t)
	d, err := repo.Draft(ctx, "u1")
	if err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(filepath.Join(d.Dir, "templates"), 0o755)
	os.WriteFile(filepath.Join(d.Dir, "templates", "extra.gadx"), []byte("extra\n"), 0o644)
	os.WriteFile(filepath.Join(d.Dir, "templates", "same.gadx"), []byte("same\n"), 0o644)
	files := fstest.MapFS{
		"index.gadx":          {Data: []byte("default\n")},
		"templates/same.gadx": {Data: []byte("same\n")},
		"templates/new.gadx":  {Data: []byte("new\n")},
	}
	str := func(changes []Change) (s string) {
		for _, c := range changes {
			s += c.Status + " " + c.Path + ";"
		}
		return
	}
	changes, err := d.DefaultsChanges(files)
	if err != nil || str(changes) != "M index.gadx;A templates/new.gadx;" {
		t.Errorf("changes %q %v", str(changes), err)
	}
	changes, _ = d.DefaultsChanges(files, "templates")
	if str(changes) != "M index.gadx;D templates/extra.gadx;A templates/new.gadx;" {
		t.Errorf("with the extras %q", str(changes))
	}
	if f, err := d.DefaultsDiff(files, "index.gadx"); err != nil || f == nil || f.Old != "old\n" || f.New != "default\n" {
		t.Errorf("the diff %+v %v", f, err)
	}
	if f, _ := d.DefaultsDiff(files, "templates/extra.gadx", "templates"); f == nil || f.Old != "extra\n" || f.New != "" {
		t.Errorf("the diff of an extra %+v", f)
	}
	if f, _ := d.DefaultsDiff(files, "templates/same.gadx"); f != nil {
		t.Errorf("a diff of a file the same %+v", f)
	}
}
