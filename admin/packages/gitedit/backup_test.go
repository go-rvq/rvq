package gitedit

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/x/packages/backup"
)

// The repository's backup: its git directory at public/.git of the state
// (the site's, a .git file naming another directory), the files not
// committed with a warning each — not the committed ones, nor the
// ignored.
func TestBackupSource(t *testing.T) {
	repo, served := site(t)
	root := filepath.Dir(served)
	os.WriteFile(filepath.Join(served, ".gitignore"), []byte(".transpiled/\n"), 0o644)
	os.MkdirAll(filepath.Join(served, ".transpiled"), 0o755)
	os.WriteFile(filepath.Join(served, ".transpiled/x.gad"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(served, "index.gadx"), []byte("changed here\n"), 0o644)
	os.WriteFile(filepath.Join(served, "new.css"), []byte("a{}"), 0o644)

	b, err := backup.New(backup.Options{Root: root, Dir: filepath.Join(t.TempDir(), "states")})
	if err != nil {
		t.Fatal(err)
	}
	b.Source("site-files", repo.BackupSource())
	st := &backup.State{ID: "s"}
	if err := b.BuildState(context.Background(), st); err != nil {
		t.Fatal(err)
	}
	m, err := b.Manifest("s")
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"public/.git/HEAD", "public/.git/config", "public/index.gadx", "public/new.css", "public/.gitignore"} {
		if _, ok := m.Find(p); !ok {
			t.Errorf("%s not in the state", p)
		}
	}
	if _, ok := m.Find("public/.transpiled/x.gad"); ok {
		t.Error("an ignored file in the state")
	}
	var objects int
	for _, e := range m.Entries {
		if strings.HasPrefix(e.Path, "public/.git/objects/") && e.Type == backup.TypeFile {
			objects++
		}
	}
	if objects == 0 {
		t.Error("no object of the repository")
	}
	if got := m.Header.Meta; got["site-files.git"] != "public/.git" || got["site-files.worktree"] != "public" || got["site-files.branch"] != "main" {
		t.Errorf("meta %v", got)
	}
	w := strings.Join(st.Result.Sources[0].Warnings, "|")
	for _, want := range []string{"index.gadx: not committed", "new.css: not committed", ".gitignore: not committed"} {
		if !strings.Contains(w, want) {
			t.Errorf("%q not warned: %s", want, w)
		}
	}

	// all committed: no warning; a file deleted is said
	os.Remove(filepath.Join(served, "new.css"))
	os.Remove(filepath.Join(served, ".gitignore"))
	os.RemoveAll(filepath.Join(served, ".transpiled"))
	os.Remove(filepath.Join(served, "index.gadx"))
	st = &backup.State{ID: "s2"}
	if err := b.BuildState(context.Background(), st); err != nil {
		t.Fatal(err)
	}
	if w := st.Result.Sources[0].Warnings; len(w) != 1 || !strings.Contains(w[0], "index.gadx: deleted") {
		t.Errorf("warnings %v", w)
	}
}

func TestStatusPaths(t *testing.T) {
	out := " M a.txt\x00?? b c\x00R  new.txt\x00old.txt\x00 D gone\x00"
	if got := strings.Join(statusPaths(out), "|"); got != "a.txt|b c|new.txt|gone" {
		t.Errorf("paths %q", got)
	}
}
