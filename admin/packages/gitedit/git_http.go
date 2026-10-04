package gitedit

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	fs_tools "github.com/go-rvq/rvq/admin/packages/fs-tools"
	"github.com/go-rvq/rvq/admin/presets"
)

// ActionGit is reaching the files by git (cloning, pushing): asked before
// anything else of a repository by git (GitRepo, DraftGitRepo).
const ActionGit = "Git"

// emptyTree is git's tree of no file: what a branch created is compared with.
const emptyTree = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"

// maxDenied is how many files a refusal names.
const maxDenied = 20

// allowedAction says whether r may do the page's action ("!git").
func (b *Builder) allowedAction(r *http.Request, action string) bool {
	ver := b.page.Page().ActionVerifier(r, presets.ActionPerm(action))
	return ver == nil || ver.Allowed()
}

// gitAllowed is the permission of a repository by git: !git first; reading
// asks @get of every file of rev (a clone is the whole repository: a path the
// user may not see is not cloned at all); writing asks extra (!publish of the
// one the site serves) — and, each file of a push, what the IDE asks
// (checkPush).
func (b *Builder) gitAllowed(r *http.Request, a fs_tools.GitAccess, gitDir, workTree, rev string, extra ...string) bool {
	if !b.allowedAction(r, ActionGit) {
		return false
	}
	if a == fs_tools.GitWrite {
		for _, action := range extra {
			if !b.allowedAction(r, action) {
				return false
			}
		}
		return true
	}
	if !b.AllowedOp(r, IdeRead, "") {
		return false
	}
	args := []string{"--git-dir=" + gitDir}
	if workTree != "" {
		// over a core.worktree of another machine
		args = append(args, "--work-tree="+workTree)
	}
	out, err := git(r.Context(), gitDir, append(args, "ls-tree", "-r", "-z", "--name-only", rev)...)
	if err != nil {
		return false
	}
	for _, p := range strings.Split(strings.TrimRight(out, "\x00"), "\x00") {
		if p != "" && !b.AllowedOp(r, IdeRead, p) {
			return false
		}
	}
	return true
}

// fileChange is a file a push changes: what it asks (create, edit…) of its
// path, and of where a rename puts it.
type fileChange struct {
	op       IdeOp
	path, to string
}

// pushChanges are the files changed from old to next (a ref created: from
// no file), renames found.
func pushChanges(ctx context.Context, rc *fs_tools.GitReceive, old, next string) ([]fileChange, error) {
	if old == fs_tools.ZeroID {
		old = emptyTree
	}
	out, err := rc.Git(ctx, "diff-tree", "-r", "-z", "-M", "--name-status", "--no-commit-id", old, next)
	if err != nil {
		return nil, err
	}
	f := strings.Split(strings.TrimRight(out, "\x00"), "\x00")
	var changes []fileChange
	for i := 0; i < len(f); i++ {
		status := f[i]
		if status == "" || i+1 >= len(f) {
			continue
		}
		switch status[0] {
		case 'R', 'C':
			if i+2 >= len(f) {
				return nil, fmt.Errorf("diff-tree: %q", out)
			}
			from, to := f[i+1], f[i+2]
			i += 2
			if status[0] == 'C' {
				changes = append(changes, fileChange{op: IdeCreate, path: to})
				continue
			}
			changes = append(changes, fileChange{op: renameOp(from, to), path: from, to: to})
			if !strings.HasPrefix(status, "R100") {
				// renamed and changed
				changes = append(changes, fileChange{op: IdeEdit, path: to})
			}
		case 'A':
			changes = append(changes, fileChange{op: IdeCreate, path: f[i+1]})
			i++
		case 'D':
			changes = append(changes, fileChange{op: IdeDelete, path: f[i+1]})
			i++
		default: // M, T: changed
			changes = append(changes, fileChange{op: IdeEdit, path: f[i+1]})
			i++
		}
	}
	return changes, nil
}

// checkPush refuses a push that is not a fast-forward of branch, alone, or
// changes a file the user may not (AllowedOp, as the IDE asks: create, edit,
// rename, move, delete — of each path); the refusal names the files.
func (b *Builder) checkPush(ctx context.Context, rc *fs_tools.GitReceive, branch string) error {
	var denied []string
	for _, u := range rc.Updates {
		if u.Ref != "refs/heads/"+branch {
			return fmt.Errorf("only the branch %s is pushed here, not %s", branch, u.Ref)
		}
		if u.New == fs_tools.ZeroID {
			return fmt.Errorf("the branch %s is not deleted", branch)
		}
		if u.Old != fs_tools.ZeroID {
			if _, err := rc.Git(ctx, "merge-base", "--is-ancestor", u.Old, u.New); err != nil {
				return errors.New("not a fast-forward: pull (fetch and merge or rebase) first")
			}
		}
		changes, err := pushChanges(ctx, rc, u.Old, u.New)
		if err != nil {
			return err
		}
		for _, c := range changes {
			for _, p := range []string{c.path, c.to} {
				if p != "" && !b.AllowedOp(rc.Request, c.op, p) {
					denied = append(denied, fmt.Sprintf("%s: no permission (%s)", p, opPerms[c.op]))
				}
			}
		}
	}
	if len(denied) > 0 {
		more := ""
		if len(denied) > maxDenied {
			more = fmt.Sprintf("\n… and %d more", len(denied)-maxDenied)
			denied = denied[:maxDenied]
		}
		return errors.New(strings.Join(denied, "\n") + more)
	}
	return nil
}

// validatePush says whether the files of each commit pushed work (Validate:
// the templates compile), checked out in a directory of their own.
func (b *Builder) validatePush(ctx context.Context, rc *fs_tools.GitReceive) error {
	if b.Validate == nil {
		return nil
	}
	for _, u := range rc.Updates {
		tmp, err := os.MkdirTemp("", "gitedit-push-")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmp)
		files := filepath.Join(tmp, "files") + string(filepath.Separator)
		env := append(append([]string{}, rc.Env...), "GIT_INDEX_FILE="+filepath.Join(tmp, "index"))
		for _, args := range [][]string{{"read-tree", u.New}, {"checkout-index", "-a", "-f", "--prefix=" + files}} {
			cmd := exec.CommandContext(ctx, "git", args...)
			cmd.Dir, cmd.Env = rc.GitDir, env
			var errb bytes.Buffer
			cmd.Stderr = &errb
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(errb.String()))
			}
		}
		if err := b.Validate(ctx, files); err != nil {
			return fmt.Errorf("the files pushed do not work: %w", err)
		}
	}
	return nil
}

// GitRepo is the repository the site serves, by git at <admin>/<name>.git: a
// push to it publishes — with what the page's Publish asks: !git and
// !publish, a fast-forward, the site's files as committed, the files working
// (Validate), and of each file changed what the IDE asks of it. One at a time
// with the page's publishing (Repo.Lock).
func (b *Builder) GitRepo(name string) *fs_tools.GitRepo {
	gitDir := func() (string, error) { return filepath.Abs(b.Repo.GitDir) }
	return &fs_tools.GitRepo{
		Name: name,
		Dir:  func(*http.Request) (string, error) { return gitDir() },
		Allowed: func(r *http.Request, a fs_tools.GitAccess) bool {
			dir, err := gitDir()
			if err != nil {
				return false
			}
			work, err := filepath.Abs(b.Repo.WorkTree)
			return err == nil && b.gitAllowed(r, a, dir, work, b.Repo.branch(), ActionPublish)
		},
		Middleware: []func(http.Handler) http.Handler{func(next http.Handler) http.Handler {
			// a push, the whole of it: no publishing meanwhile
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/git-receive-pack") {
					b.Repo.Lock()
					defer b.Repo.Unlock()
				}
				next.ServeHTTP(w, r)
			})
		}},
		// its worktree the site's (over a core.worktree of another machine),
		// moved by PostReceive: git leaves it as it is
		Env:    []string{"GIT_WORK_TREE=" + absPath(b.Repo.WorkTree)},
		Config: []string{"receive.denyCurrentBranch=ignore", "receive.denyNonFastForwards=true", "receive.denyDeletes=true"},
		PreReceive: func(ctx context.Context, rc *fs_tools.GitReceive) error {
			if err := b.checkPush(ctx, rc, b.Repo.branch()); err != nil {
				return err
			}
			if err := b.Repo.SiteClean(ctx); err != nil {
				return err
			}
			return b.validatePush(ctx, rc)
		},
		PostReceive: func(ctx context.Context, rc *fs_tools.GitReceive) {
			for _, u := range rc.Updates {
				if err := b.Repo.Advance(ctx, u.Old, u.New); err != nil {
					slog.Error("gitedit: the site's files after a push", "from", u.Old, "to", u.New, "error", err)
				}
			}
		},
	}
}

// DraftGitRepo is the draft of the user of the request, by git at
// <admin>/<name>.git (made the first time): a push to it changes the draft —
// seen by the editor and the preview at once —, published from the page. It
// asks !git and, of each file changed, what the IDE asks; refused while the
// draft has changes not committed.
func (b *Builder) DraftGitRepo(name string) *fs_tools.GitRepo {
	draftDir := func(r *http.Request) (string, error) {
		d, _, err := b.Draft(r)
		if err != nil {
			return "", err
		}
		return filepath.Abs(d.Dir)
	}
	return &fs_tools.GitRepo{
		Name: name,
		Dir:  draftDir,
		Allowed: func(r *http.Request, a fs_tools.GitAccess) bool {
			dir, err := draftDir(r)
			return err == nil && b.gitAllowed(r, a, filepath.Join(dir, ".git"), "", "HEAD")
		},
		Config: []string{"receive.denyCurrentBranch=updateInstead", "receive.denyNonFastForwards=true",
			"receive.denyDeletes=true"},
		PreReceive: func(ctx context.Context, rc *fs_tools.GitReceive) error {
			d, _, err := b.Draft(rc.Request)
			if err != nil {
				return err
			}
			if changes, err := d.Status(ctx); err != nil {
				return err
			} else if len(changes) > 0 {
				return errors.New("the draft has changes not committed (in the editor): commit or discard them first")
			}
			return b.checkPush(ctx, rc, b.Repo.branch())
		},
	}
}

// absPath is p absolute (p itself when it cannot be made so).
func absPath(p string) string {
	if a, err := filepath.Abs(p); err == nil {
		return a
	}
	return p
}
