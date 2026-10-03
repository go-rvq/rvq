package gitedit

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// site is a repository like the site's: its git directory apart, its worktree
// the files served, pushes into it updating them.
func site(t *testing.T) (repo *Repo, served string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("no git")
	}
	root := t.TempDir()
	served = filepath.Join(root, "public")
	run := func(dir string, args ...string) {
		t.Helper()
		if _, err := git(context.Background(), dir, args...); err != nil {
			t.Fatal(err)
		}
	}
	os.MkdirAll(served, 0o755)
	run(root, "init", "--quiet", "--initial-branch=main", "--separate-git-dir="+filepath.Join(root, ".public.git"), served)
	os.WriteFile(filepath.Join(served, "index.gadx"), []byte("old\n"), 0o644)
	run(served, "add", ".")
	run(served, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "--quiet", "-m", "first")
	// as on the server: the repository and its worktree by relative paths
	run(served, "config", "core.worktree", "../public")
	os.WriteFile(filepath.Join(served, ".git"), []byte("gitdir: ../.public.git\n"), 0o644)
	return &Repo{GitDir: filepath.Join(root, ".public.git"), WorkTree: served, Branch: "main", DraftsDir: filepath.Join(root, "drafts")}, served
}

// A draft is edited, committed by its author and published: the site's
// files change only then.
func TestDraftCommitPublish(t *testing.T) {
	ctx := context.Background()
	repo, served := site(t)
	d, err := repo.Draft(ctx, "u1")
	if err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(d.Dir, "index.gadx"), []byte("new\n"), 0o644)
	os.WriteFile(filepath.Join(d.Dir, "about.gadx"), []byte("about\n"), 0o644)

	changes, err := d.Status(ctx)
	if err != nil || len(changes) != 2 {
		t.Fatalf("status %v %v", changes, err)
	}
	if diff, _ := d.Diff(ctx, "index.gadx"); !strings.Contains(diff, "+new") {
		t.Errorf("diff %q", diff)
	}
	if diff, _ := d.Diff(ctx, "about.gadx"); !strings.Contains(diff, "+about") {
		t.Errorf("diff of a new file %q", diff)
	}
	if _, err := d.Diff(ctx, ".git/config"); err == nil {
		t.Error("a diff of .git")
	}
	if err := d.Publish(ctx); !errors.Is(err, ErrUncommitted) {
		t.Errorf("publishing with changes: %v", err)
	}
	if _, err := d.Commit(ctx, "", Author{"Ana", "ana@x"}); err == nil {
		t.Error("a commit with no message")
	}
	if _, err := d.Commit(ctx, "the index", Author{"Ana", "ana@x"}); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "old\n" {
		t.Errorf("committed, the site changed: %q", b)
	}
	log, _ := d.Log(ctx, 1)
	if len(log) != 1 || log[0].Author != "Ana" || log[0].Email != "ana@x" || log[0].Subject != "the index" {
		t.Errorf("log %+v", log)
	}
	if s, _ := d.Sync(ctx); s.Ahead != 1 || s.Behind != 0 {
		t.Errorf("sync %+v", s)
	}
	if err := d.Publish(ctx); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(served, "index.gadx")); string(b) != "new\n" {
		t.Errorf("published, the site has %q", b)
	}
}

// Another editor published first: publishing asks for an update, which
// rebases the draft's commits, its changes kept; a change discarded goes.
func TestDraftUpdateDiscard(t *testing.T) {
	ctx := context.Background()
	repo, served := site(t)
	a, _ := repo.Draft(ctx, "a")
	b, _ := repo.Draft(ctx, "b")

	os.WriteFile(filepath.Join(a.Dir, "a.txt"), []byte("a\n"), 0o644)
	a.Commit(ctx, "a", Author{"A", "a@x"})
	if err := a.Publish(ctx); err != nil {
		t.Fatal(err)
	}

	os.WriteFile(filepath.Join(b.Dir, "b.txt"), []byte("b\n"), 0o644)
	b.Commit(ctx, "b", Author{"B", "b@x"})
	os.WriteFile(filepath.Join(b.Dir, "wip.txt"), []byte("wip\n"), 0o644)
	os.WriteFile(filepath.Join(b.Dir, "index.gadx"), []byte("changed\n"), 0o644)
	if err := b.Discard(ctx, "wip.txt"); err != nil {
		t.Fatal(err)
	}
	if err := b.Discard(ctx, "index.gadx"); err != nil {
		t.Fatal(err)
	}
	if ch, _ := b.Status(ctx); len(ch) != 0 {
		t.Errorf("after discarding: %v", ch)
	}
	if err := b.Publish(ctx); !errors.Is(err, ErrNotFastForward) {
		t.Fatalf("publishing behind: %v", err)
	}
	if err := b.Update(ctx); err != nil {
		t.Fatal(err)
	}
	if err := b.Publish(ctx); err != nil {
		t.Fatal(err)
	}
	// a file of the site changed there: publishing overwrites nothing
	os.WriteFile(filepath.Join(served, "index.gadx"), []byte("by hand\n"), 0o644)
	os.WriteFile(filepath.Join(b.Dir, "c.txt"), []byte("c\n"), 0o644)
	b.Commit(ctx, "c", Author{"B", "b@x"})
	if err := b.Publish(ctx); !errors.Is(err, ErrSiteChanged) {
		t.Errorf("publishing over a change of the site: %v", err)
	}
	if _, err := os.Stat(filepath.Join(served, "c.txt")); err == nil {
		t.Error("published over a change of the site")
	}
	for _, f := range []string{"a.txt", "b.txt"} {
		if _, err := os.Stat(filepath.Join(served, f)); err != nil {
			t.Errorf("the site has no %s", f)
		}
	}
}
