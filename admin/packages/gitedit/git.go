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
	"io/fs"
	"os"
	"os/exec"
	pathpkg "path"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Repo is the repository edited: GitDir its git directory, WorkTree the
// files the site serves (checked out from Branch).
type Repo struct {
	GitDir   string
	WorkTree string
	Branch   string
	// DraftsDir holds the drafts, one directory each.
	DraftsDir string

	// mu is held while the branch the site serves moves: a publishing, a push
	// to it (Lock).
	mu sync.Mutex
}

// Lock is held while the branch the site serves moves — a publishing, a push
// to the repository —: one at a time.
func (r *Repo) Lock() { r.mu.Lock() }

// Unlock releases Lock.
func (r *Repo) Unlock() { r.mu.Unlock() }

// OpenRepo is the repository of the files of workTree — its .git a
// directory, or a file naming it ("gitdir: ../.public.git", a submodule's or
// a worktree's) —, the branch checked out there, the drafts in draftsDir.
func OpenRepo(ctx context.Context, workTree, draftsDir string) (*Repo, error) {
	gitDir := filepath.Join(workTree, ".git")
	st, err := os.Stat(gitDir)
	if err != nil {
		return nil, fmt.Errorf("gitedit: %s is no git worktree: %w", workTree, err)
	}
	if !st.IsDir() {
		b, err := os.ReadFile(gitDir)
		if err != nil {
			return nil, err
		}
		p, ok := strings.CutPrefix(strings.TrimSpace(string(b)), "gitdir:")
		if !ok {
			return nil, fmt.Errorf("gitedit: %s: no gitdir", gitDir)
		}
		if p = strings.TrimSpace(p); !filepath.IsAbs(p) {
			p = filepath.Join(workTree, p)
		}
		gitDir = p
	}
	// the git directory and the worktree said: the repository's
	// core.worktree may name a path of another machine (the host of a
	// container; a deploy tool writes it)
	absGit, err := filepath.Abs(gitDir)
	if err != nil {
		return nil, err
	}
	absWork, err := filepath.Abs(workTree)
	if err != nil {
		return nil, err
	}
	branch, err := git(ctx, workTree, "--git-dir="+absGit, "--work-tree="+absWork, "symbolic-ref", "--short", "HEAD")
	if err != nil {
		return nil, err
	}
	return &Repo{GitDir: gitDir, WorkTree: workTree, Branch: strings.TrimSpace(branch), DraftsDir: draftsDir}, nil
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
	cmd.Env = append(gitEnv(), "GIT_TERMINAL_PROMPT=0", "LC_ALL=C")
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(errb.String()))
	}
	return out.String(), nil
}

// gitEnv is the environment of the process without git's own variables
// (GIT_DIR, GIT_INDEX_FILE, GIT_WORK_TREE…): run from a git hook, they would
// point every command at that repository instead of dir's.
func gitEnv() (env []string) {
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "GIT_") {
			env = append(env, e)
		}
	}
	return
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
	d := &Draft{Dir: filepath.Join(r.DraftsDir, key), Key: key, repo: r}
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
	Dir string
	// Key is its owner's (the user whose draft it is).
	Key  string
	repo *Repo
}

// editorsFile holds the keys of the users who changed the files of the draft
// since its last commit, one a line: its co-authors.
const editorsFile = "rvq-editors"

// NoteEditor notes that the user of key changed the files of the draft (a
// co-author of its next commit).
func (d *Draft) NoteEditor(key string) error {
	if key == "" || slices.Contains(d.Editors(), key) {
		return nil
	}
	f, err := os.OpenFile(filepath.Join(d.Dir, ".git", editorsFile), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(key + "\n")
	return err
}

// Editors are the keys of the users who changed the files of the draft since
// its last commit (NoteEditor), in the order they did.
func (d *Draft) Editors() (keys []string) {
	b, _ := os.ReadFile(filepath.Join(d.Dir, ".git", editorsFile))
	for _, k := range strings.Split(string(b), "\n") {
		if k = strings.TrimSpace(k); k != "" && !slices.Contains(keys, k) {
			keys = append(keys, k)
		}
	}
	return
}

// clearEditors forgets the editors: a commit took their changes.
func (d *Draft) clearEditors() { _ = os.Remove(filepath.Join(d.Dir, ".git", editorsFile)) }

// Change is a file changed in the draft, not committed: its status (git
// porcelain: "M", "A", "D", "R", "??") and path; a renamed one's path before
// (From).
type Change struct {
	Status string `json:"status"`
	Path   string `json:"path"`
	From   string `json:"from,omitempty"`
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
		ch := Change{Status: st, Path: e[3:]}
		if st[0] == 'R' && i+1 < len(entries) {
			i++ // the origin of a rename follows
			ch.From = entries[i]
		}
		changes = append(changes, ch)
	}
	return
}

// renameSimilarity is how alike a file deleted and one added must be to be
// taken for one renamed (moved): git's default for -M, 50%.
const renameSimilarity = 0.5

// Changes are the changes of the draft (Status) with its renames found: a
// file the IDE renames or moves is, for git, one deleted and one new (it
// only says R of a rename staged) — they are one renamed, From its path
// before, when their contents are the same or alike (renameSimilarity). The
// most alike pair first.
func (d *Draft) Changes(ctx context.Context) ([]Change, error) {
	changes, err := d.Status(ctx)
	if err != nil {
		return nil, err
	}
	type side struct {
		i       int
		content string
	}
	var gone, added []side
	for i, ch := range changes {
		switch {
		case strings.HasPrefix(ch.Status, "D"):
			old, _, err := d.Versions(ctx, ch)
			if err != nil {
				return nil, err
			}
			gone = append(gone, side{i, old})
		case strings.HasPrefix(ch.Status, "A"), strings.HasPrefix(ch.Status, "?"):
			_, cur, err := d.Versions(ctx, ch)
			if err != nil {
				return nil, err
			}
			added = append(added, side{i, cur})
		}
	}
	if len(gone) == 0 || len(added) == 0 {
		return changes, nil
	}
	type pair struct {
		g, a  int
		score float64
	}
	var pairs []pair
	for gi, g := range gone {
		for ai, a := range added {
			if s := similarity(g.content, a.content); s >= renameSimilarity {
				pairs = append(pairs, pair{gi, ai, s})
			}
		}
	}
	// the most alike first; alike the same, the closest by name and folder
	near := func(p pair) float64 {
		from, to := changes[gone[p.g].i].Path, changes[added[p.a].i].Path
		n := similarity(pathpkg.Base(from), pathpkg.Base(to))
		if pathpkg.Dir(from) == pathpkg.Dir(to) {
			n++
		}
		return n
	}
	sort.SliceStable(pairs, func(i, j int) bool {
		if pairs[i].score != pairs[j].score {
			return pairs[i].score > pairs[j].score
		}
		return near(pairs[i]) > near(pairs[j])
	})
	usedG, usedA := map[int]bool{}, map[int]bool{}
	drop := map[int]bool{}
	for _, p := range pairs {
		if usedG[p.g] || usedA[p.a] {
			continue
		}
		usedG[p.g], usedA[p.a] = true, true
		g, a := changes[gone[p.g].i], &changes[added[p.a].i]
		a.Status, a.From = "R", g.Path
		drop[gone[p.g].i] = true
	}
	out := changes[:0]
	for i, ch := range changes {
		if !drop[i] {
			out = append(out, ch)
		}
	}
	return out, nil
}

// similarity is how alike two contents are, from 0 to 1: 1 less the
// characters of their lines changed over the longer's. Two empty are alike.
func similarity(a, b string) float64 {
	if a == b {
		return 1
	}
	longer := max(len([]rune(a)), len([]rune(b)))
	if longer == 0 {
		return 1
	}
	// the characters of the lines changed, over the longer's
	return 1 - float64(diffmatchpatch.New().DiffLevenshtein(vx.LineDiff(a, b)))/float64(longer)
}

// WriteFile writes the file path of the draft with content: made, with its
// folders, when it is not there.
func (d *Draft) WriteFile(path, content string) error {
	if err := d.check(path); err != nil {
		return err
	}
	full := filepath.Join(d.Dir, filepath.FromSlash(strings.TrimPrefix(pathpkg.Clean("/"+path), "/")))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	return os.WriteFile(full, []byte(content), 0o644)
}

// Versions are the file of ch as the last commit has it (old; "" for a file
// it has not: added, untracked) and as the draft has it (new; "" for a file
// deleted).
func (d *Draft) Versions(ctx context.Context, ch Change) (old, new string, err error) {
	if err := d.check(ch.Path); err != nil {
		return "", "", err
	}
	from := ch.Path
	if ch.From != "" {
		if err := d.check(ch.From); err != nil {
			return "", "", err
		}
		from = ch.From
	}
	if st := ch.Status; !strings.HasPrefix(st, "A") && !strings.HasPrefix(st, "?") {
		if old, err = git(ctx, d.Dir, "show", "HEAD:"+from); err != nil {
			return "", "", err
		}
	}
	data, err := os.ReadFile(filepath.Join(d.Dir, filepath.FromSlash(ch.Path)))
	switch {
	case err == nil:
		new = string(data)
	case errors.Is(err, fs.ErrNotExist):
		err = nil // deleted
	default:
		return "", "", err
	}
	return old, new, nil
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
	d.clearEditors()
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

// site runs git on the repository the site serves, its git directory and
// worktree said: its core.worktree may name a path of another machine (the
// host of a container; a deploy tool writes it).
func (r *Repo) site(ctx context.Context, args ...string) (string, error) {
	gitDir, err := filepath.Abs(r.GitDir)
	if err != nil {
		return "", err
	}
	workTree, err := filepath.Abs(r.WorkTree)
	if err != nil {
		return "", err
	}
	return git(ctx, workTree, append([]string{"--git-dir=" + gitDir, "--work-tree=" + workTree}, args...)...)
}

// Head is the commit of the branch the site serves.
func (r *Repo) Head(ctx context.Context) (string, error) {
	head, err := r.site(ctx, "rev-parse", "HEAD")
	return strings.TrimSpace(head), err
}

// SiteClean refuses (ErrSiteChanged) when the site's files are not the ones
// committed: changed out of the editor, they would be overwritten.
func (r *Repo) SiteClean(ctx context.Context) error {
	if _, err := r.site(ctx, "update-index", "-q", "--ignore-submodules", "--refresh"); err != nil {
		return err
	}
	if _, err := r.site(ctx, "diff-files", "--quiet", "--ignore-submodules", "--"); err != nil {
		return ErrSiteChanged
	}
	if _, err := r.site(ctx, "diff-index", "--quiet", "--cached", "--ignore-submodules", "HEAD", "--"); err != nil {
		return ErrSiteChanged
	}
	return nil
}

// Advance makes the site's files those of next, from those of head: a
// fast-forward of the worktree (the branch moved by whoever called it).
func (r *Repo) Advance(ctx context.Context, head, next string) error {
	_, err := r.site(ctx, "read-tree", "-u", "-m", head, next)
	return err
}

// checkout makes the branch the commit checked out in the draft dir, and the
// site's files those of it: a fast-forward, the files as committed.
func (r *Repo) checkout(ctx context.Context, dir string) error {
	r.Lock()
	defer r.Unlock()
	if err := r.SiteClean(ctx); err != nil {
		return err
	}
	head, err := r.Head(ctx)
	if err != nil {
		return err
	}
	draftDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	if _, err := r.site(ctx, "fetch", "--quiet", draftDir, "HEAD"); err != nil {
		return err
	}
	next, err := r.site(ctx, "rev-parse", "FETCH_HEAD")
	if err != nil {
		return err
	}
	next = strings.TrimSpace(next)
	if _, err := r.site(ctx, "merge-base", "--is-ancestor", head, next); err != nil {
		return ErrNotFastForward
	}
	if err := r.Advance(ctx, head, next); err != nil {
		return err
	}
	_, err = r.site(ctx, "update-ref", "refs/heads/"+r.branch(), next, head)
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
