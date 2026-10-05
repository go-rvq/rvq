package fs_tools

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// GitAccess is what a request asks of a git repository: to read it (clone,
// fetch) or to write in it (push).
type GitAccess int

const (
	GitRead GitAccess = iota
	GitWrite
)

// RefUpdate is a ref a push changes: from Old to New (a zero id: created or
// deleted).
type RefUpdate struct {
	Old, New, Ref string
}

// ZeroID is the id of no commit: the Old of a ref created, the New of one
// deleted.
const ZeroID = "0000000000000000000000000000000000000000"

// GitReceive is a push being received: its request, the repository (Dir) and
// the refs it changes. Env is the environment of git in the hook — with it,
// git sees the objects of the push, not yet in the repository before it is
// accepted (the quarantine).
type GitReceive struct {
	Request *http.Request
	Dir     string
	// GitDir is the git directory of the repository (absolute).
	GitDir  string
	Updates []RefUpdate
	Env     []string
}

// Git runs git on the repository of the push, the objects of the push seen.
func (rc *GitReceive) Git(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = rc.GitDir
	cmd.Env = rc.Env
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(errb.String()))
	}
	return strings.TrimRight(out.String(), "\n"), nil
}

// GitRepo is a git repository served by HTTP (git's smart protocol), at
// <admin>/<Name>.git: cloned and pushed to by the users of the admin, by
// their login.
type GitRepo struct {
	// Name is its name in the URL: "site-files" is <admin>/site-files.git.
	// Ending in "/{key}", it is a repository of each key: "drafts/{key}" is
	// <admin>/drafts/<key>.git, the key told by GitRepoKey.
	Name string
	// Dir is the repository of the request (its .git, or the directory
	// holding it): one of each user, say.
	Dir func(r *http.Request) (string, error)
	// Allowed says whether the request may read the repository, or write in
	// it; nil allows all (to whoever logged in).
	Allowed func(r *http.Request, a GitAccess) bool
	// Middleware wraps the repository's handler, after the login and before
	// the permission: making the user's repository the first time, say.
	Middleware []func(http.Handler) http.Handler
	// Config is the configuration of git for it, "key=value"
	// ("receive.denyCurrentBranch=updateInstead").
	Config []string
	// Env is more of git's environment, "KEY=value" (GIT_WORK_TREE, over a
	// core.worktree that names another machine's path); the hooks' commands
	// (GitReceive.Git) have its GIT_WORK_TREE too.
	Env []string
	// PreReceive is asked before a push is accepted: an error refuses it —
	// nothing written —, its message shown to whoever pushed.
	PreReceive func(ctx context.Context, rc *GitReceive) error
	// PostReceive is told of a push accepted: its refs updated.
	PostReceive func(ctx context.Context, rc *GitReceive)
}

// gitServer serves the repositories (AddGitRepo), and the hooks of their
// pushes: git runs a hook as a program, the application itself
// (RunGitHookIfInvoked), which asks the server, by a loopback address, with
// the token of the push.
type gitServer struct {
	repos []*GitRepo

	mu       sync.Mutex
	pending  map[string]*pendingPush
	hookURL  string
	hooksDir string
	hookErr  error
	once     sync.Once
}

type pendingPush struct {
	repo *GitRepo
	r    *http.Request
	dir  string
}

// AddGitRepo serves a git repository at <admin>/<Name>.git (GitHandler).
func (b *Builder) AddGitRepo(repo *GitRepo) *Builder {
	if b.git == nil {
		b.git = &gitServer{pending: map[string]*pendingPush{}}
	}
	b.git.repos = append(b.git.repos, repo)
	return b
}

// GitRepos are the repositories added (AddGitRepo).
func (b *Builder) GitRepos() []*GitRepo {
	if b.git == nil {
		return nil
	}
	return b.git.repos
}

// GitPath is the path of the repository name: <admin>/<name>.git.
func (b *Builder) GitPath(name string) string {
	return strings.TrimSuffix(b.p.GetURIPrefix(), "/") + "/" + name + ".git"
}

// GitHandler serves the repositories added (AddGitRepo), each at its path
// (GitPath) — a mux of its own for the paths <admin>/<name>.git/…, each
// request by the user of its login (HTTP Basic, as the WebDAV). nil when
// there is none.
func (b *Builder) GitHandler() http.Handler {
	if b.git == nil {
		return nil
	}
	h := b.gitHandler()
	if b.lb != nil {
		h = b.lb.BasichAuthMiddleware(h)
	}
	return h
}

// gitHandler is GitHandler with no login: the user is the request's.
func (b *Builder) gitHandler() http.Handler {
	mux := http.NewServeMux()
	for _, repo := range b.git.repos {
		wrap := func(h http.Handler) http.Handler {
			for i := len(repo.Middleware) - 1; i >= 0; i-- {
				h = repo.Middleware[i](h)
			}
			return h
		}
		if base, ok := strings.CutSuffix(repo.Name, "/"+GitKeyPattern); ok {
			// a repository of each key: <admin>/<base>/<key>.git/…
			prefix := strings.TrimSuffix(b.p.GetURIPrefix(), "/") + "/" + base + "/"
			mux.Handle(prefix, wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				seg, _, _ := strings.Cut(strings.TrimPrefix(r.URL.Path, prefix), "/")
				key, ok := strings.CutSuffix(seg, ".git")
				if !ok || !validGitKey(key) {
					http.NotFound(w, r)
					return
				}
				r = r.WithContext(context.WithValue(r.Context(), gitKeyCtx{}, key))
				b.git.repoHandler(repo, prefix+seg).ServeHTTP(w, r)
			})))
			continue
		}
		mux.Handle(b.GitPath(repo.Name)+"/", wrap(b.git.repoHandler(repo, b.GitPath(repo.Name))))
	}
	return mux
}

// GitMountPaths are the paths GitHandler serves, each ending in "/": a
// repository's (<admin>/<name>.git/), or the folder of the repositories of
// each key (<admin>/<base>/) — to mount it on another mux.
func (b *Builder) GitMountPaths() (paths []string) {
	if b.git == nil {
		return nil
	}
	for _, repo := range b.git.repos {
		if base, ok := strings.CutSuffix(repo.Name, "/"+GitKeyPattern); ok {
			paths = append(paths, strings.TrimSuffix(b.p.GetURIPrefix(), "/")+"/"+base+"/")
			continue
		}
		paths = append(paths, b.GitPath(repo.Name)+"/")
	}
	return
}

// GitKeyPattern ends the Name of a repository of each key (GitRepo).
const GitKeyPattern = "{key}"

type gitKeyCtx struct{}

// GitRepoKey is the key of the repository of r, of a GitRepo of each key
// ("drafts/{key}": the <key> of <admin>/drafts/<key>.git); "" for one of its
// own.
func GitRepoKey(r *http.Request) string {
	k, _ := r.Context().Value(gitKeyCtx{}).(string)
	return k
}

// GitKeyPath is the path of the repository of key of name ("drafts/{key}"):
// <admin>/drafts/<key>.git.
func (b *Builder) GitKeyPath(name, key string) string {
	return b.GitPath(strings.Replace(name, GitKeyPattern, key, 1))
}

// validGitKey says whether key names a repository: letters, digits, "-"
// and "_", nothing that walks a path.
func validGitKey(key string) bool {
	if key == "" || len(key) > 128 {
		return false
	}
	for _, c := range key {
		if !(c == '-' || c == '_' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
			return false
		}
	}
	return true
}

// gitService is what a request of the smart protocol asks, by its path:
// reading (upload-pack) or writing (receive-pack); ok false for a path of no
// service (the "dumb" protocol, not served).
func gitService(r *http.Request, root string) (a GitAccess, ok bool) {
	p := strings.TrimPrefix(r.URL.Path, root)
	switch {
	case p == "/info/refs" && r.Method == http.MethodGet:
		switch r.URL.Query().Get("service") {
		case "git-upload-pack":
			return GitRead, true
		case "git-receive-pack":
			return GitWrite, true
		}
	case p == "/git-upload-pack" && r.Method == http.MethodPost:
		return GitRead, true
	case p == "/git-receive-pack" && r.Method == http.MethodPost:
		return GitWrite, true
	}
	return 0, false
}

func (s *gitServer) repoHandler(repo *GitRepo, root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		access, ok := gitService(r, root)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if repo.Allowed != nil && !repo.Allowed(r, access) {
			http.Error(w, "permission denied", http.StatusForbidden)
			return
		}
		dir, err := repo.Dir(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		config := append([]string{"safe.directory=*", "http.receivepack=true", "http.uploadpack=true"}, repo.Config...)
		env := append([]string{"GIT_PROJECT_ROOT=" + dir, "GIT_HTTP_EXPORT_ALL=1", "REMOTE_USER=" + remoteUser(r)}, repo.Env...)
		if access == GitWrite && r.Method == http.MethodPost && (repo.PreReceive != nil || repo.PostReceive != nil) {
			url, hooks, err := s.hooks()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			token := s.register(&pendingPush{repo: repo, r: r, dir: dir})
			defer s.forget(token)
			config = append(config, "core.hooksPath="+hooks)
			env = append(env, envHookURL+"="+url, envHookToken+"="+token)
		}
		env = append(env, fmt.Sprintf("GIT_CONFIG_COUNT=%d", len(config)))
		for i, c := range config {
			k, v, _ := strings.Cut(c, "=")
			env = append(env, fmt.Sprintf("GIT_CONFIG_KEY_%d=%s", i, k), fmt.Sprintf("GIT_CONFIG_VALUE_%d=%s", i, v))
		}

		gitBin, err := exec.LookPath("git")
		if err != nil {
			http.Error(w, "git: "+err.Error(), http.StatusInternalServerError)
			return
		}
		(&cgi.Handler{
			Path:       gitBin,
			Args:       []string{"http-backend"},
			Root:       root,
			Dir:        dir,
			Env:        env,
			InheritEnv: []string{"PATH", "HOME", "TMPDIR"},
			Stderr:     io.Discard,
		}).ServeHTTP(w, r)
	})
}

// remoteUser is the user of the request, as git records it: its login's
// name, when HTTP Basic told it.
func remoteUser(r *http.Request) string {
	if u, _, ok := r.BasicAuth(); ok {
		return u
	}
	return "user"
}

// The environment that makes the application a hook of a push.
const (
	envHook      = "RVQ_GIT_HOOK"
	envHookURL   = "RVQ_GIT_HOOK_URL"
	envHookToken = "RVQ_GIT_HOOK_TOKEN"
	// envHookCwd is where git ran the hook (the repository's directory).
	envHookCwd = "RVQ_GIT_HOOK_CWD"
)

func (s *gitServer) register(p *pendingPush) string {
	b := make([]byte, 24)
	_, _ = rand.Read(b)
	token := hex.EncodeToString(b)
	s.mu.Lock()
	s.pending[token] = p
	s.mu.Unlock()
	return token
}

func (s *gitServer) forget(token string) {
	s.mu.Lock()
	delete(s.pending, token)
	s.mu.Unlock()
}

// hooks are the URL the hooks ask (a loopback address of this process) and
// the directory of the hooks (scripts running the application as a hook),
// made the first time.
func (s *gitServer) hooks() (string, string, error) {
	s.once.Do(func() {
		exe, err := os.Executable()
		if err != nil {
			s.hookErr = err
			return
		}
		// the application runs where the server does (its .env, its
		// configuration): git runs a hook in the repository's directory, told
		// to it (envHookCwd) for the paths of git's environment
		wd, err := os.Getwd()
		if err != nil {
			s.hookErr = err
			return
		}
		dir, err := os.MkdirTemp("", "rvq-git-hooks-")
		if err != nil {
			s.hookErr = err
			return
		}
		for _, name := range []string{"pre-receive", "post-receive"} {
			script := "#!/bin/sh\n" + envHookCwd + "=\"$PWD\"; export " + envHookCwd + "\ncd " + shellQuote(wd) +
				" || exit 1\n" + envHook + "=" + name + " exec " + shellQuote(exe) + "\n"
			if err := os.WriteFile(filepath.Join(dir, name), []byte(script), 0o755); err != nil {
				s.hookErr = err
				return
			}
		}
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			s.hookErr = err
			return
		}
		go func() { _ = http.Serve(ln, http.HandlerFunc(s.serveHook)) }()
		s.hookURL, s.hooksDir = "http://"+ln.Addr().String()+"/hook", dir
	})
	return s.hookURL, s.hooksDir, s.hookErr
}

func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

// hookRequest is what a hook tells the server: which, the token of its
// push, the refs (its standard input), its environment.
type hookRequest struct {
	Hook  string   `json:"hook"`
	Token string   `json:"token"`
	Input string   `json:"input"`
	Env   []string `json:"env"`
	// Cwd is where git ran the hook: its git directory, which the paths of
	// its environment may be relative to.
	Cwd string `json:"cwd"`
}

// hookResponse is the server's answer: the message shown to whoever
// pushed, and whether the push goes on.
type hookResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (s *gitServer) serveHook(w http.ResponseWriter, r *http.Request) {
	var req hookRequest
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	p := s.pending[req.Token]
	s.mu.Unlock()
	if p == nil {
		_ = json.NewEncoder(w).Encode(hookResponse{Message: "unknown push"})
		return
	}
	env := hookEnv(req.Cwd, req.Env)
	gitDir := req.Cwd
	for _, kv := range env {
		if v, ok := strings.CutPrefix(kv, "GIT_DIR="); ok {
			gitDir = v
		}
	}
	rc := &GitReceive{Request: p.r, Dir: p.dir, GitDir: gitDir, Updates: parseRefUpdates(req.Input), Env: env}
	res := hookResponse{OK: true}
	switch req.Hook {
	case "pre-receive":
		if p.repo.PreReceive != nil {
			if err := p.repo.PreReceive(p.r.Context(), rc); err != nil {
				res = hookResponse{Message: err.Error()}
			}
		}
	case "post-receive":
		if p.repo.PostReceive != nil {
			p.repo.PostReceive(p.r.Context(), rc)
		}
	}
	_ = json.NewEncoder(w).Encode(res)
}

// hookEnv is the environment of git in a hook that a command of the server
// keeps — what makes it see the objects of the push and its repository —,
// its paths made absolute (from cwd, where git ran the hook), over the
// server's own, out of any git's.
func hookEnv(cwd string, env []string) []string {
	abs := func(p string) string {
		if p == "" || filepath.IsAbs(p) {
			return p
		}
		return filepath.Join(cwd, p)
	}
	var out []string
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "GIT_") {
			out = append(out, kv)
		}
	}
	for _, kv := range env {
		k, v, _ := strings.Cut(kv, "=")
		switch k {
		case "GIT_DIR", "GIT_QUARANTINE_PATH", "GIT_OBJECT_DIRECTORY", "GIT_WORK_TREE":
			out = append(out, k+"="+abs(v))
		case "GIT_ALTERNATE_OBJECT_DIRECTORIES":
			parts := strings.Split(v, string(filepath.ListSeparator))
			for i := range parts {
				parts[i] = abs(parts[i])
			}
			out = append(out, k+"="+strings.Join(parts, string(filepath.ListSeparator)))
		}
	}
	return append(out, "GIT_CONFIG_COUNT=1", "GIT_CONFIG_KEY_0=safe.directory", "GIT_CONFIG_VALUE_0=*")
}

// parseRefUpdates reads the refs a hook is told, a line each: "old new ref".
func parseRefUpdates(input string) (out []RefUpdate) {
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		if f := strings.Fields(line); len(f) == 3 {
			out = append(out, RefUpdate{Old: f[0], New: f[1], Ref: f[2]})
		}
	}
	return
}

// ErrNotHook is RunGitHook's error when the application was not run as a
// hook.
var ErrNotHook = errors.New("not run as a git hook")

// RunGitHookIfInvoked runs the application as the hook of a push, when git
// ran it so (a repository of GitHandler): it tells the server the push and
// exits — with the server's answer, its message shown to whoever pushed.
// Otherwise it returns at once. Call it first thing in main.
func RunGitHookIfInvoked() {
	code, err := runGitHook(os.Stdin, os.Stderr)
	if errors.Is(err, ErrNotHook) {
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(code)
}

func runGitHook(stdin io.Reader, stderr io.Writer) (int, error) {
	hook := os.Getenv(envHook)
	if hook == "" {
		return 0, ErrNotHook
	}
	input, err := io.ReadAll(stdin)
	if err != nil {
		return 1, err
	}
	cwd := os.Getenv(envHookCwd)
	if cwd == "" {
		cwd, _ = os.Getwd()
	}
	body, _ := json.Marshal(hookRequest{Hook: hook, Token: os.Getenv(envHookToken), Input: string(input),
		Env: os.Environ(), Cwd: cwd})
	resp, err := http.Post(os.Getenv(envHookURL), "application/json", bytes.NewReader(body))
	if err != nil {
		return 1, fmt.Errorf("the server of the push: %w", err)
	}
	defer resp.Body.Close()
	var res hookResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 1, fmt.Errorf("the server of the push: %w", err)
	}
	if res.Message != "" {
		fmt.Fprintln(stderr, res.Message)
	}
	if !res.OK {
		return 1, nil
	}
	return 0, nil
}
