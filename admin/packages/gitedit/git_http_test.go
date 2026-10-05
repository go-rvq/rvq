package gitedit

import (
	"context"
	"errors"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	fs_tools "github.com/go-rvq/rvq/admin/packages/fs-tools"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
)

// TestMain: the test binary is the hook of the pushes, as the application is.
func TestMain(m *testing.M) {
	fs_tools.RunGitHookIfInvoked()
	os.Exit(m.Run())
}

// gitClient runs git as a user's would, asking no one (no terminal, no
// askpass, no credential helper): its output and error.
func gitClient(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-c", "credential.helper=", "-c", "core.askPass=",
		"-c", "user.name=Ana", "-c", "user.email=ana@x", "-c", "init.defaultBranch=main"}, args...)...)
	cmd.Dir = dir
	var env []string
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "SSH_ASKPASS=") && !strings.HasPrefix(kv, "GIT_") {
			env = append(env, kv)
		}
	}
	cmd.Env = append(env, "GIT_TERMINAL_PROMPT=0", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// gitSite is the editor of a site's files with its repositories by git (the
// site's, site-files.git; the draft's, site-files-draft.git), served; the
// permissions pb; invalid the error of Validate.
func gitSite(t *testing.T, pb *perm.Builder, invalid *error) (url func(name string) string, repo *Repo, served string, b *Builder) {
	t.Helper()
	_, repo, served, b = appWith(t, invalid, pb)
	p := presets.New(i18n.New()).URIPrefix("/admin")
	fsb, err := fs_tools.New(p, p.I18n(), nil)
	if err != nil {
		t.Fatal(err)
	}
	fsb.AddGitRepo(b.GitRepo("site-files")).AddGitRepo(b.DraftGitRepo("site-files-draft"))
	srv := httptest.NewServer(fsb.GitHandler())
	t.Cleanup(srv.Close)
	// the hook commit-msg of the user, as the page serves it: the commits of
	// the clones say who made them on the site (Site-User)
	w := httptest.NewRecorder()
	b.serveHook(w, httptest.NewRequest("GET", srv.URL+"/admin/site-files/commit-msg", nil))
	testHook = w.Body.String()
	t.Cleanup(func() { testHook = "" })
	return func(name string) string { return srv.URL + fsb.GitPath(name) }, repo, served, b
}

// testHook is the hook commit-msg the clones of the test have (gitSite).
var testHook string

// installHook gives the clone the hook of the test, when it has none.
func installHook(t *testing.T, clone string) {
	t.Helper()
	p := filepath.Join(clone, ".git", "hooks", "commit-msg")
	if _, err := os.Stat(p); err == nil || testHook == "" {
		return
	}
	os.MkdirAll(filepath.Dir(p), 0o755)
	if err := os.WriteFile(p, []byte(testHook), 0o755); err != nil {
		t.Fatal(err)
	}
}

// mustGit runs gitClient, failing on an error.
func mustGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := gitClient(dir, args...)
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return out
}

// commitFile writes path in the clone and commits it.
func commitFile(t *testing.T, clone, path, content string) {
	t.Helper()
	os.MkdirAll(filepath.Dir(filepath.Join(clone, path)), 0o755)
	os.WriteFile(filepath.Join(clone, path), []byte(content), 0o644)
	installHook(t, clone)
	mustGit(t, clone, "add", "-A")
	mustGit(t, clone, "commit", "-q", "-m", "change "+path)
}

func allowAll(policies ...*perm.PolicyBuilder) *perm.Builder {
	pb := perm.New().AllowAll()
	pb.CreatePolicies(policies...)
	return pb
}

func policyOn(effect, res string) *perm.PolicyBuilder {
	return perm.PolicyFor(perm.Anybody).WhoAre(effect).ToDo(perm.Anything).On(res)
}

// The draft by git: cloned; a push changes it — the site does not —; refused
// while the draft has changes not committed, or when not a fast-forward.
func TestGitDraft(t *testing.T) {
	var invalid error
	url, repo, served, _ := gitSite(t, allowAll(), &invalid)
	root := t.TempDir()
	clone := filepath.Join(root, "clone")
	mustGit(t, root, "clone", "-q", url("site-files-draft"), clone)
	if b, _ := os.ReadFile(filepath.Join(clone, "index.gadx")); string(b) != "old\n" {
		t.Fatalf("the clone has %q", b)
	}

	commitFile(t, clone, "index.gadx", "by git\n")
	mustGit(t, clone, "push", "-q", "origin", "main")
	draft := filepath.Join(repo.DraftsDir, "u1")
	if b, _ := os.ReadFile(filepath.Join(draft, "index.gadx")); string(b) != "by git\n" {
		t.Errorf("the draft has %q", b)
	}
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "old\n" {
		t.Errorf("a push to the draft changed the site: %q", b)
	}

	// a change not committed in the draft (the editor's): refused
	os.WriteFile(filepath.Join(draft, "wip.gadx"), []byte("wip"), 0o644)
	commitFile(t, clone, "about.gadx", "about\n")
	if out, err := gitClient(clone, "push", "origin", "main"); err == nil || !strings.Contains(out, "not committed") {
		t.Errorf("a push over changes not committed: %v\n%s", err, out)
	}
	os.Remove(filepath.Join(draft, "wip.gadx"))
	mustGit(t, clone, "push", "-q", "origin", "main")

	// not a fast-forward: refused
	mustGit(t, clone, "reset", "-q", "--hard", "HEAD~2")
	commitFile(t, clone, "other.gadx", "x\n")
	if out, err := gitClient(clone, "push", "-f", "origin", "main"); err == nil {
		t.Errorf("a push not a fast-forward: %s", out)
	}
}

// The site by git: a push publishes — the site's files change — with what the
// page's Publish asks: the files working, the site's files as committed,
// !publish.
func TestGitPublished(t *testing.T) {
	var invalid error
	url, repo, served, _ := gitSite(t, allowAll(), &invalid)
	root := t.TempDir()
	clone := filepath.Join(root, "clone")
	mustGit(t, root, "clone", "-q", url("site-files"), clone)

	commitFile(t, clone, "index.gadx", "published by git\n")
	mustGit(t, clone, "push", "-q", "origin", "main")
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "published by git\n" {
		t.Errorf("the site has %q", b)
	}
	if err := repo.SiteClean(context.Background()); err != nil {
		t.Errorf("the site's files after a push: %v", err)
	}
	head, _ := repo.Head(context.Background())
	if strings.TrimSpace(mustGit(t, clone, "rev-parse", "HEAD")) != head {
		t.Error("the branch did not move")
	}

	// files that do not work: refused, the site as it was
	invalid = errors.New("index.gadx:1: broken template")
	commitFile(t, clone, "index.gadx", "broken\n")
	if out, err := gitClient(clone, "push", "origin", "main"); err == nil || !strings.Contains(out, "broken template") {
		t.Errorf("a push of files that do not work: %v\n%s", err, out)
	}
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "published by git\n" {
		t.Errorf("a push refused changed the site: %q", b)
	}
	if h, _ := repo.Head(context.Background()); h != head {
		t.Error("a push refused moved the branch")
	}
	invalid = nil

	// the site's files changed out of git: refused, nothing overwritten
	os.WriteFile(filepath.Join(served, "index.gadx"), []byte("by hand\n"), 0o644)
	if out, err := gitClient(clone, "push", "origin", "main"); err == nil {
		t.Errorf("a push over the site changed by hand: %s", out)
	}
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "by hand\n" {
		t.Errorf("overwritten: %q", b)
	}

	// no !publish: no push (the clone is still read)
	url, _, _, _ = gitSite(t, allowAll(policyOn(perm.Denied, "admin:/site-files:!publish")), &invalid)
	other := filepath.Join(root, "other")
	mustGit(t, root, "clone", "-q", url("site-files"), other)
	commitFile(t, other, "a.gadx", "a\n")
	if out, err := gitClient(other, "push", "origin", "main"); err == nil {
		t.Errorf("pushed to the site with no !publish: %s", out)
	}
}

// A push asks, of each file it changes, what the IDE asks (AllowedOp): the
// same rules of paths, the refusal naming the files; a clone asks @get of
// every file; nothing at all with no !git.
func TestGitPermissionsAsTheIDE(t *testing.T) {
	var invalid error
	root := t.TempDir()
	const res = "admin:/site-files:"

	t.Run("a path denied", func(t *testing.T) {
		url, _, _, b := gitSite(t, allowAll(policyOn(perm.Denied, res+"<config/*>:!edit")), &invalid)
		r := httptest.NewRequest("GET", "/", nil)
		clone := filepath.Join(t.TempDir(), "c")
		mustGit(t, root, "clone", "-q", url("site-files-draft"), clone)
		commitFile(t, clone, "config/a.gad", "1\n") // created: !create, allowed
		mustGit(t, clone, "push", "-q", "origin", "main")
		commitFile(t, clone, "config/a.gad", "2\n") // edited: !edit of config/, denied
		out, err := gitClient(clone, "push", "origin", "main")
		if err == nil || !strings.Contains(out, "config/a.gad: no permission (!edit)") {
			t.Errorf("an edit denied: %v\n%s", err, out)
		}
		// the IDE says the same
		if b.AllowedOp(r, IdeEdit, "config/a.gad") || !b.AllowedOp(r, IdeCreate, "config/b.gad") {
			t.Error("the IDE says otherwise")
		}
	})

	t.Run("only static/", func(t *testing.T) {
		pb := perm.New()
		pb.CreatePolicies(policyOn(perm.Allowed, res+"@get"), policyOn(perm.Allowed, res+"!git"),
			policyOn(perm.Allowed, res+"<static/*>:!create"), policyOn(perm.Allowed, res+"<static/*>:!rename"))
		url, _, _, _ := gitSite(t, pb, &invalid)
		clone := filepath.Join(t.TempDir(), "c")
		mustGit(t, root, "clone", "-q", url("site-files-draft"), clone)
		commitFile(t, clone, "static/a.css", "a\n")
		mustGit(t, clone, "push", "-q", "origin", "main")
		// renamed in static/: allowed; moved out of it: denied
		mustGit(t, clone, "mv", "static/a.css", "static/b.css")
		mustGit(t, clone, "commit", "-q", "-m", "rename")
		mustGit(t, clone, "push", "-q", "origin", "main")
		mustGit(t, clone, "mv", "static/b.css", "templates.css")
		mustGit(t, clone, "commit", "-q", "-m", "move")
		if out, err := gitClient(clone, "push", "origin", "main"); err == nil || !strings.Contains(out, "!move") {
			t.Errorf("a move out of static/: %v\n%s", err, out)
		}
		mustGit(t, clone, "reset", "-q", "--hard", "HEAD~1")
		commitFile(t, clone, "templates/x.gadx", "x\n")
		if out, err := gitClient(clone, "push", "origin", "main"); err == nil || !strings.Contains(out, "templates/x.gadx: no permission (!create)") {
			t.Errorf("a file created out of static/: %v\n%s", err, out)
		}
	})

	t.Run("a file not to be seen: no clone", func(t *testing.T) {
		url, repo, _, _ := gitSite(t, allowAll(policyOn(perm.Denied, res+"<secret/*>:@get")), &invalid)
		d, _ := repo.Draft(context.Background(), "u1")
		os.MkdirAll(filepath.Join(d.Dir, "secret"), 0o755)
		os.WriteFile(filepath.Join(d.Dir, "secret", "k.txt"), []byte("k"), 0o644)
		d.Commit(context.Background(), "secret", Author{"Ana", "ana@x"})
		if out, err := gitClient(root, "clone", "-q", url("site-files-draft"), filepath.Join(t.TempDir(), "c")); err == nil {
			t.Errorf("cloned a file not to be seen: %s", out)
		}
		// the site's, with no such file: cloned
		mustGit(t, root, "clone", "-q", url("site-files"), filepath.Join(t.TempDir(), "s"))
	})

	t.Run("no !git", func(t *testing.T) {
		url, _, _, _ := gitSite(t, allowAll(policyOn(perm.Denied, res+"!git")), &invalid)
		for _, name := range []string{"site-files", "site-files-draft"} {
			if out, err := gitClient(root, "clone", "-q", url(name), filepath.Join(t.TempDir(), name)); err == nil {
				t.Errorf("%s cloned with no !git: %s", name, out)
			}
		}
	})
}
