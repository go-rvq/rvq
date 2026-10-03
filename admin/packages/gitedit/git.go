// Package gitedit edits the files of a git repository from the admin — the
// public/ of a site, its templates and static files — in drafts: each editor
// works on a clone of their own (Draft), commits there, and publishes, which
// brings the commits to the branch the site serves and the site's files to
// them (as a push to it with updateInstead would). A template half written
// never reaches the site.
//
// git is the binary (exec): the repository and its worktree point at each
// other by relative paths, which it follows as the site does.
package gitedit

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Repo is the repository edited: GitDir its git directory, WorkTree the
// files the site serves (checked out from Branch).
type Repo struct {
	GitDir   string
	WorkTree string
	Branch   string
	// DraftsDir holds the drafts, one directory each.
	DraftsDir string
}

// Author is who commits.
type Author struct {
	Name, Email string
}

// ErrNotFastForward says the draft is behind the branch: Update it first.
var ErrNotFastForward = errors.New("the branch moved since the draft was updated: update the draft first")

// ErrUncommitted says the draft has changes not committed.
var ErrUncommitted = errors.New("the draft has changes not committed")

// git runs git in dir; its error carries what git said.
func git(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "LC_ALL=C")
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(errb.String()))
	}
	return out.String(), nil
}

func (r *Repo) branch() string {
	if r.Branch == "" {
		return "main"
	}
	return r.Branch
}

// Draft is the draft of key (a user's), cloned from the repository the
// first time.
func (r *Repo) Draft(ctx context.Context, key string) (*Draft, error) {
	if key == "" || strings.ContainsAny(key, `/\.`) {
		return nil, fmt.Errorf("gitedit: invalid draft key %q", key)
	}
	d := &Draft{Dir: filepath.Join(r.DraftsDir, key), repo: r}
	if _, err := os.Stat(filepath.Join(d.Dir, ".git")); err == nil {
		return d, nil
	}
	if err := os.MkdirAll(r.DraftsDir, 0o755); err != nil {
		return nil, err
	}
	gitDir, err := filepath.Abs(r.GitDir)
	if err != nil {
		return nil, err
	}
	// --shared: the objects are the repository's, nothing copied
	if _, err := git(ctx, r.DraftsDir, "clone", "--quiet", "--shared", "--branch", r.branch(), gitDir, key); err != nil {
		return nil, err
	}
	return d, nil
}

// Draft is the clone an editor works on.
type Draft struct {
	Dir  string
	repo *Repo
}

// Change is a file changed in the draft, not committed: its status (git
// porcelain: "M", "A", "D", "R", "??") and path.
type Change struct {
	Status string `json:"status"`
	Path   string `json:"path"`
}

// Status are the changes of the draft not committed.
func (d *Draft) Status(ctx context.Context) (changes []Change, err error) {
	out, err := git(ctx, d.Dir, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	if err != nil {
		return nil, err
	}
	entries := strings.Split(out, "\x00")
	for i := 0; i < len(entries); i++ {
		e := entries[i]
		if len(e) < 4 {
			continue
		}
		st := strings.TrimSpace(e[:2])
		if st == "" {
			st = "M"
		}
		if st[0] == 'R' {
			i++ // the origin of a rename follows
		}
		changes = append(changes, Change{Status: st, Path: e[3:]})
	}
	return
}

// Diff is what changed in the file path of the draft, not committed.
func (d *Draft) Diff(ctx context.Context, path string) (string, error) {
	if err := d.check(path); err != nil {
		return "", err
	}
	out, err := git(ctx, d.Dir, "diff", "HEAD", "--", path)
	if err != nil || out != "" {
		return out, err
	}
	// a new file: all of it
	out, err = git(ctx, d.Dir, "diff", "--no-index", "--", os.DevNull, path)
	if err != nil && out != "" {
		err = nil // exits 1 when the files differ
	}
	return out, err
}

// Commit commits every change of the draft, as author.
func (d *Draft) Commit(ctx context.Context, message string, author Author) (string, error) {
	if strings.TrimSpace(message) == "" {
		return "", errors.New("gitedit: a commit needs a message")
	}
	if _, err := git(ctx, d.Dir, "add", "--all"); err != nil {
		return "", err
	}
	if _, err := git(ctx, d.Dir, "-c", "user.name="+author.Name, "-c", "user.email="+author.Email,
		"commit", "--quiet", "--message", message); err != nil {
		return "", err
	}
	out, err := git(ctx, d.Dir, "rev-parse", "--short", "HEAD")
	return strings.TrimSpace(out), err
}

// Discard drops the changes of the file path of the draft, not committed: a
// file changed or deleted comes back as committed, a new one goes.
func (d *Draft) Discard(ctx context.Context, path string) error {
	if err := d.check(path); err != nil {
		return err
	}
	if _, err := git(ctx, d.Dir, "ls-files", "--error-unmatch", "--", path); err != nil {
		return os.RemoveAll(filepath.Join(d.Dir, filepath.FromSlash(path)))
	}
	_, err := git(ctx, d.Dir, "checkout", "HEAD", "--", path)
	return err
}

// Commit is a commit of the history.
type Commit struct {
	Hash    string    `json:"hash"`
	Short   string    `json:"short"`
	Author  string    `json:"author"`
	Email   string    `json:"email"`
	Date    time.Time `json:"date"`
	Subject string    `json:"subject"`
}

// Log are the last n commits of the draft.
func (d *Draft) Log(ctx context.Context, n int) (commits []Commit, err error) {
	out, err := git(ctx, d.Dir, "log", "-n", strconv.Itoa(n), "--format=%H%x1f%h%x1f%an%x1f%ae%x1f%aI%x1f%s")
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		f := strings.Split(line, "\x1f")
		if len(f) < 6 {
			continue
		}
		date, _ := time.Parse(time.RFC3339, f[4])
		commits = append(commits, Commit{Hash: f[0], Short: f[1], Author: f[2], Email: f[3], Date: date, Subject: f[5]})
	}
	return
}

// Sync says where the draft stands to the branch: commits of its own not
// published (Ahead) and of the branch not in it (Behind).
type Sync struct {
	Ahead  int `json:"ahead"`
	Behind int `json:"behind"`
}

// Sync is where the draft stands to the branch, fetched now.
func (d *Draft) Sync(ctx context.Context) (s Sync, err error) {
	if _, err = git(ctx, d.Dir, "fetch", "--quiet", "origin"); err != nil {
		return
	}
	out, err := git(ctx, d.Dir, "rev-list", "--left-right", "--count", "origin/"+d.repo.branch()+"...HEAD")
	if err != nil {
		return
	}
	f := strings.Fields(out)
	if len(f) == 2 {
		s.Behind, _ = strconv.Atoi(f[0])
		s.Ahead, _ = strconv.Atoi(f[1])
	}
	return
}

// Update brings the commits of the branch into the draft: its own commits
// rebased on them, the changes not committed kept. A conflict aborts it,
// leaving the draft as it was.
func (d *Draft) Update(ctx context.Context) error {
	if _, err := git(ctx, d.Dir, "fetch", "--quiet", "origin"); err != nil {
		return err
	}
	if _, err := git(ctx, d.Dir, "-c", "user.name=gitedit", "-c", "user.email=gitedit@localhost",
		"rebase", "--autostash", "origin/"+d.repo.branch()); err != nil {
		_, _ = git(ctx, d.Dir, "rebase", "--abort")
		return err
	}
	return nil
}

// ErrSiteChanged says files of the site's checkout were changed there, not
// by a commit: nothing is overwritten.
var ErrSiteChanged = errors.New("files of the site were changed out of git: nothing was overwritten")

// Publish brings the commits of the draft to the branch, and the site's files
// to them — as a push to a checked-out branch with updateInstead, done here:
// the repository's push hook may name paths of another machine (the host of
// a container). The draft must have nothing uncommitted (ErrUncommitted) and
// be up to date with the branch (ErrNotFastForward); the site's files must be
// as committed (ErrSiteChanged).
func (d *Draft) Publish(ctx context.Context) error {
	changes, err := d.Status(ctx)
	if err != nil {
		return err
	}
	if len(changes) > 0 {
		return ErrUncommitted
	}
	s, err := d.Sync(ctx)
	if err != nil {
		return err
	}
	if s.Behind > 0 {
		return ErrNotFastForward
	}
	if s.Ahead == 0 {
		return nil
	}
	return d.repo.checkout(ctx, d.Dir)
}

// checkout makes the branch the commit checked out in the draft dir, and the
// site's files those of it: a fast-forward, the files as committed.
func (r *Repo) checkout(ctx context.Context, dir string) error {
	gitDir, err := filepath.Abs(r.GitDir)
	if err != nil {
		return err
	}
	workTree, err := filepath.Abs(r.WorkTree)
	if err != nil {
		return err
	}
	site := func(args ...string) (string, error) {
		return git(ctx, workTree, append([]string{"--git-dir=" + gitDir, "--work-tree=" + workTree}, args...)...)
	}
	// the site's files as committed
	if _, err := site("update-index", "-q", "--ignore-submodules", "--refresh"); err != nil {
		return err
	}
	if _, err := site("diff-files", "--quiet", "--ignore-submodules", "--"); err != nil {
		return ErrSiteChanged
	}
	if _, err := site("diff-index", "--quiet", "--cached", "--ignore-submodules", "HEAD", "--"); err != nil {
		return ErrSiteChanged
	}
	head, err := site("rev-parse", "HEAD")
	if err != nil {
		return err
	}
	head = strings.TrimSpace(head)
	draftDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	if _, err := site("fetch", "--quiet", draftDir, "HEAD"); err != nil {
		return err
	}
	next, err := site("rev-parse", "FETCH_HEAD")
	if err != nil {
		return err
	}
	next = strings.TrimSpace(next)
	if _, err := site("merge-base", "--is-ancestor", head, next); err != nil {
		return ErrNotFastForward
	}
	if _, err := site("read-tree", "-u", "-m", head, next); err != nil {
		return err
	}
	_, err = site("update-ref", "refs/heads/"+r.branch(), next, head)
	return err
}

// Reset drops the draft: it is cloned again, from the branch, when asked.
func (d *Draft) Reset() error { return os.RemoveAll(d.Dir) }

// check refuses a path out of the draft or in its .git.
func (d *Draft) check(path string) error {
	clean := filepath.ToSlash(filepath.Clean("/" + path))
	if clean == "/" || clean == "/.git" || strings.HasPrefix(clean, "/.git/") {
		return fmt.Errorf("gitedit: invalid path %q", path)
	}
	return nil
}
