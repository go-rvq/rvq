package gitedit

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gad-lang/gad/web/ide"
)

// IdePerm is what an operation of the IDE asks: to see the files (read) or
// to change them (write).
type IdePerm int

const (
	IdeRead IdePerm = iota
	IdeWrite
)

// ideOps are the operations of the IDE served, each with what it asks; the
// others — running, evaluating, debugging code, fetching URLs into the files
// — run code or reach out from the server, and are not served.
var ideOps = map[string]struct {
	read, write bool // by method: GET reads; PUT/POST of a write op writes
}{
	"workspace": {read: true},
	"tree":      {read: true},
	"file":      {read: true, write: true}, // GET ?path / PUT
	"config":    {read: true, write: true}, // GET / PUT
	"modules":   {read: true},
	"format":    {read: true},
	"diagnose":  {read: true},
	"transpile": {read: true},
	"doc":       {read: true},
	"mkdir":     {write: true},
	"delete":    {write: true},
	"rename":    {write: true},
}

// IDE serves the API of the gad IDE (web/ide) on the draft of each editor:
// the operations on its files and the language's, not the ones that run
// code; no path in .git.
type IDE struct {
	repo *Repo
	// Name is the workspace's name the IDE shows.
	Name string
	// DraftKey is the draft of the request (its user's).
	DraftKey func(r *http.Request) string
	// Allowed says whether the request may do what op asks.
	Allowed func(r *http.Request, p IdePerm) bool

	mu      sync.Mutex
	servers map[string]*ide.Server
}

func NewIDE(repo *Repo) *IDE {
	return &IDE{repo: repo, Name: "public", servers: map[string]*ide.Server{}}
}

func (h *IDE) server(r *http.Request) (*ide.Server, error) {
	key := h.DraftKey(r)
	h.mu.Lock()
	defer h.mu.Unlock()
	if s := h.servers[key]; s != nil {
		return s, nil
	}
	d, err := h.repo.Draft(r.Context(), key)
	if err != nil {
		return nil, err
	}
	s, err := ide.New(d.Dir)
	if err != nil {
		return nil, err
	}
	h.servers[key] = s
	return s, nil
}

// Forget drops the IDE of a draft (reset): made again when asked.
func (h *IDE) Forget(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.servers, key)
}

func ideError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// ServeHTTP serves /api/ide/<op> (its prefix stripped).
func (h *IDE) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	op, ok := strings.CutPrefix(r.URL.Path, "/api/ide/")
	spec, known := ideOps[op]
	if !ok || !known {
		ideError(w, http.StatusNotFound, "not available")
		return
	}
	need := IdeRead
	if r.Method != http.MethodGet {
		if !spec.write {
			if !spec.read {
				ideError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
		} else {
			need = IdeWrite
		}
	} else if !spec.read {
		ideError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if h.Allowed != nil && !h.Allowed(r, need) {
		ideError(w, http.StatusForbidden, "permission denied")
		return
	}
	// no path in .git: the query's, the body's
	if inGit(r.URL.Query().Get("path")) {
		ideError(w, http.StatusForbidden, "path not allowed")
		return
	}
	if r.Body != nil && r.Method != http.MethodGet {
		body, err := io.ReadAll(io.LimitReader(r.Body, 16<<20))
		if err != nil {
			ideError(w, http.StatusBadRequest, err.Error())
			return
		}
		var fields map[string]any
		if json.Unmarshal(body, &fields) == nil {
			for _, k := range []string{"path", "to"} {
				if s, _ := fields[k].(string); inGit(s) {
					ideError(w, http.StatusForbidden, "path not allowed")
					return
				}
			}
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
	}
	s, err := h.server(r)
	if err != nil {
		ideError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if op == "workspace" {
		// not the path on the server
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"root": "/", "name": h.Name, "openFile": "", "compute": "server"})
		return
	}
	s.Handler().ServeHTTP(w, r)
}

// inGit says whether path is .git or in it.
func inGit(path string) bool {
	for _, part := range strings.Split(strings.ReplaceAll(path, `\`, "/"), "/") {
		if part == ".git" {
			return true
		}
	}
	return false
}
