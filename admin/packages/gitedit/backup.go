package gitedit

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-rvq/rvq/x/packages/backup"
)

// BackupSource is the source of the repository in a backup state
// (x/packages/backup): its git directory, at <worktree>/.git of the state
// (a submodule's or a worktree's made a repository of its own) — the
// commits; the client checks the worktree out of it —, and the files of the
// worktree not committed (modified, not tracked and not ignored), each
// warned of. The branch does not move meanwhile (Lock). Meta: "git" (its
// git directory in the state), "worktree", "branch".
func (r *Repo) BackupSource() backup.SourceFunc {
	return func(ctx context.Context, st *backup.State) error {
		wt, err := filepath.Abs(r.WorkTree)
		if err != nil {
			return err
		}
		gitDir, err := filepath.Abs(r.GitDir)
		if err != nil {
			return err
		}
		wtRel, err := filepath.Rel(st.FS.Root(), wt)
		if err != nil || wtRel == ".." || strings.HasPrefix(wtRel, ".."+string(filepath.Separator)) || wtRel == "." {
			return fmt.Errorf("gitedit: the worktree %s is not under %s", wt, st.FS.Root())
		}
		wtRel = filepath.ToSlash(wtRel)
		gitAt := path.Join(wtRel, ".git")

		r.Lock()
		defer r.Unlock()
		// the locks of a git running now are not the repository's
		if err := st.FS.AddTreeFrom(gitAt, gitDir, "*.lock", "gc.log"); err != nil {
			return err
		}
		out, err := r.site(ctx, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignore-submodules")
		if err != nil {
			return err
		}
		for _, p := range statusPaths(out) {
			abs := filepath.Join(wt, filepath.FromSlash(p))
			fi, err := os.Lstat(abs)
			switch {
			case errors.Is(err, fs.ErrNotExist):
				st.FS.Warn("%s: deleted, not committed (the copy has it)", p)
			case err != nil:
				return err
			case fi.IsDir():
				// a folder of a repository of its own: its files are its own
			default:
				if err := st.FS.AddFrom(path.Join(wtRel, p), abs); err != nil {
					return err
				}
				st.FS.Warn("%s: not committed", p)
			}
		}
		branch := r.branch()
		for k, v := range map[string]string{"git": gitAt, "worktree": wtRel, "branch": branch} {
			if err := st.FS.Meta(k, v); err != nil {
				return err
			}
		}
		return nil
	}
}

// statusPaths are the paths of git status --porcelain=v1 -z: "XY path\0",
// a rename "XY new\0old\0" (the new one).
func statusPaths(out string) (paths []string) {
	fields := strings.Split(out, "\x00")
	for i := 0; i < len(fields); i++ {
		f := fields[i]
		if len(f) < 4 {
			continue
		}
		xy, p := f[:2], f[3:]
		paths = append(paths, p)
		if xy[0] == 'R' || xy[0] == 'C' {
			i++ // the old path
		}
	}
	return
}
