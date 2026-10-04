package fs_tools

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/i18n"
)

// TestMain: the test binary is the hook of the pushes, as the application is
// (RunGitHookIfInvoked).
func TestMain(m *testing.M) {
	// the hook runs where the server does (its .env, say), not in the
	// repository's directory git runs it in
	if marker := os.Getenv("RVQ_TEST_SERVER_MARKER"); marker != "" && os.Getenv(envHook) != "" {
		if _, err := os.Stat(marker); err != nil {
			println("the hook is not where the server is:", err.Error())
			os.Exit(1)
		}
	}
	RunGitHookIfInvoked()
	os.Exit(m.Run())
}

// gitRun runs git in dir, failing the test on an error; its output.
func gitRun(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := gitTry(dir, args...)
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return out
}

func gitTry(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-c", "credential.helper=", "-c", "core.askPass=",
		"-c", "user.name=Ana", "-c", "user.email=ana@x",
		"-c", "protocol.version=2", "-c", "init.defaultBranch=main"}, args...)...)
	cmd.Dir = dir
	// never asks anyone: no terminal, no askpass, no credential helper
	var env []string
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "SSH_ASKPASS=") && !strings.HasPrefix(kv, "GIT_ASKPASS=") &&
			!strings.HasPrefix(kv, "GIT_CONFIG") {
			env = append(env, kv)
		}
	}
	cmd.Env = append(env, "GIT_TERMINAL_PROMPT=0", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// A repository by HTTP: cloned and pushed to; a push refused by its
// PreReceive — nothing written, its message shown —; one accepted told to
// its PostReceive; its middleware first, then its permission; two of them,
// each its own; only git's smart protocol.
func TestGitRepos(t *testing.T) {
	// the server in a directory of its own, a file of it there (the hook must
	// run there: TestMain)
	serverDir := t.TempDir()
	os.WriteFile(filepath.Join(serverDir, "server.marker"), []byte("x"), 0o644)
	t.Chdir(serverDir)
	t.Setenv("RVQ_TEST_SERVER_MARKER", "server.marker")
	root := t.TempDir()
	// two repositories: a bare one, and a working one (a .git in it)
	bare := filepath.Join(root, "bare.git")
	gitRun(t, root, "init", "-q", "--bare", bare)
	work := filepath.Join(root, "work")
	gitRun(t, root, "init", "-q", work)
	os.WriteFile(filepath.Join(work, "a.txt"), []byte("a\n"), 0o644)
	gitRun(t, work, "add", "-A")
	gitRun(t, work, "commit", "-q", "-m", "first")

	p := presets.New(i18n.New()).URIPrefix("/admin")
	b, err := New(p, p.I18n(), nil)
	if err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	var told []string
	b.AddGitRepo(&GitRepo{
		Name: "bare",
		Dir:  func(*http.Request) (string, error) { return bare, nil },
		Allowed: func(r *http.Request, a GitAccess) bool {
			return a == GitRead || r.Header.Get("X-Write") == "1"
		},
		Middleware: []func(http.Handler) http.Handler{func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("X-User") == "" {
					http.Error(w, "who?", http.StatusUnauthorized)
					return
				}
				next.ServeHTTP(w, r)
			})
		}},
		PreReceive: func(ctx context.Context, rc *GitReceive) error {
			for _, u := range rc.Updates {
				// the files of the push, seen before it is accepted
				names, err := rc.Git(ctx, "diff-tree", "-r", "--name-only", "--no-commit-id", "--root", u.New)
				if err != nil {
					return err
				}
				if strings.Contains(names, "bad.txt") {
					return errors.New("bad.txt: not allowed here")
				}
			}
			return nil
		},
		PostReceive: func(_ context.Context, rc *GitReceive) {
			mu.Lock()
			defer mu.Unlock()
			for _, u := range rc.Updates {
				told = append(told, u.Ref)
			}
		},
	})
	b.AddGitRepo(&GitRepo{
		Name:   "work",
		Dir:    func(*http.Request) (string, error) { return work, nil },
		Config: []string{"receive.denyCurrentBranch=updateInstead"},
	})
	srv := httptest.NewServer(b.gitHandler())
	defer srv.Close()
	url := func(name string) string { return srv.URL + b.GitPath(name) }
	headers := func(h ...string) []string {
		var out []string
		for _, x := range h {
			out = append(out, "-c", "http.extraHeader="+x)
		}
		return out
	}

	// the middleware first: no user, nothing
	if out, err := gitTry(root, "clone", "-q", url("bare"), filepath.Join(root, "nobody")); err == nil {
		t.Errorf("cloned with no user: %s", out)
	}
	// a clone of the working repository, a push to it: its files changed
	clone := filepath.Join(root, "clone")
	gitRun(t, root, append(headers("X-User: ana"), "clone", "-q", url("work"), clone)...)
	if b, _ := os.ReadFile(filepath.Join(clone, "a.txt")); string(b) != "a\n" {
		t.Fatalf("the clone has %q", b)
	}
	os.WriteFile(filepath.Join(clone, "a.txt"), []byte("b\n"), 0o644)
	gitRun(t, clone, "commit", "-q", "-am", "second")
	gitRun(t, clone, append(headers("X-User: ana"), "push", "-q", "origin", "main")...)
	if b, _ := os.ReadFile(filepath.Join(work, "a.txt")); string(b) != "b\n" {
		t.Errorf("the push did not reach the working repository: %q", b)
	}

	// the bare one: no permission to write, refused
	gitRun(t, clone, "remote", "add", "bare", url("bare"))
	if out, err := gitTry(clone, append(headers("X-User: ana"), "push", "-q", "bare", "main")...); err == nil {
		t.Errorf("pushed with no permission: %s", out)
	}
	// with it: its PreReceive refuses a file, its message shown, nothing written
	os.WriteFile(filepath.Join(clone, "bad.txt"), []byte("x"), 0o644)
	gitRun(t, clone, "add", "-A")
	gitRun(t, clone, "commit", "-q", "-m", "bad")
	out, err := gitTry(clone, append(headers("X-User: ana", "X-Write: 1"), "push", "bare", "main")...)
	if err == nil || !strings.Contains(out, "bad.txt: not allowed here") {
		t.Errorf("a push refused: %v\n%s", err, out)
	}
	if refs := gitRun(t, bare, "for-each-ref"); refs != "" {
		t.Errorf("a push refused wrote: %s", refs)
	}
	mu.Lock()
	if len(told) != 0 {
		t.Errorf("a push refused was told: %v", told)
	}
	mu.Unlock()
	// without the file: accepted, told
	gitRun(t, clone, "reset", "-q", "--hard", "HEAD~1")
	gitRun(t, clone, append(headers("X-User: ana", "X-Write: 1"), "push", "-q", "bare", "main")...)
	if head := gitRun(t, bare, "rev-parse", "main"); strings.TrimSpace(head) != strings.TrimSpace(gitRun(t, clone, "rev-parse", "HEAD")) {
		t.Errorf("the bare has %s", head)
	}
	mu.Lock()
	if len(told) != 1 || told[0] != "refs/heads/main" {
		t.Errorf("told %v", told)
	}
	mu.Unlock()

	// only the smart protocol; an unknown repository is not found
	for _, path := range []string{b.GitPath("bare") + "/HEAD", b.GitPath("bare") + "/objects/info/packs",
		b.GitPath("nope") + "/info/refs?service=git-upload-pack"} {
		r := httptest.NewRequest("GET", path, nil)
		r.Header.Set("X-User", "ana")
		w := httptest.NewRecorder()
		b.gitHandler().ServeHTTP(w, r)
		if w.Code != http.StatusNotFound {
			t.Errorf("%s: %d", path, w.Code)
		}
	}
}

// The refs a hook is told, a line each.
func TestParseRefUpdates(t *testing.T) {
	got := parseRefUpdates(ZeroID + " abc refs/heads/main\nold new refs/tags/v1\n\n")
	if len(got) != 2 || got[0] != (RefUpdate{ZeroID, "abc", "refs/heads/main"}) || got[1].Ref != "refs/tags/v1" {
		t.Errorf("%+v", got)
	}
}
