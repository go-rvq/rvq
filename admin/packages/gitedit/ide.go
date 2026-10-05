package gitedit

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gad-lang/gad/web/ide"
)

// IdeOp is what an operation of the IDE asks: to see the files, or one of
// the ways of changing them — each its own permission (Builder.AllowedOp).
type IdeOp int

const (
	// IdeRead is seeing the files.
	IdeRead IdeOp = iota
	// IdeCreate is creating a file or a folder that is not there.
	IdeCreate
	// IdeEdit is changing a file that is there.
	IdeEdit
	// IdeRename is renaming a file or a folder in its folder.
	IdeRename
	// IdeMove is moving a file or a folder to another folder.
	IdeMove
	// IdeDelete is deleting a file or a folder.
	IdeDelete
	// IdeImport is bringing files in — uploaded, or downloaded from a URL by
	// the server —: what it writes asks IdeCreate (a new file) or IdeEdit
	// (one that is there) too.
	IdeImport
)

// IdeOps are the ways of changing the files, by the name the IDE is told
// them (the workspace's "actions").
var IdeOps = map[string]IdeOp{
	"create": IdeCreate, "edit": IdeEdit, "rename": IdeRename, "move": IdeMove, "delete": IdeDelete,
	"import": IdeImport,
}

// ideOps are the operations of the IDE served — the files' and the language's
// —, whether a GET reads them, and whether another method changes the files
// (what it asks is said by the request: ServeHTTP). The others — running,
// evaluating, debugging code — run code, and are not served.
var ideOps = map[string]struct{ read, write bool }{
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
	"upload":    {write: true},
	"fetch":     {write: true},
}

// maxBody is the most a write of the IDE may hold: an upload's bytes, in
// base64 (the IDE's own limit, web/ide MaxUploadBody).
const maxBody = 32 << 20

// ideBody is what the IDE's writes say of the paths: the file, where a
// rename puts it, the files of an upload.
type ideBody struct {
	Path  string `json:"path"`
	To    string `json:"to"`
	Files []struct {
		Path string `json:"path"`
	} `json:"files"`
}

// IDE serves the API of the gad IDE (web/ide) on the draft of each editor:
// the operations on its files and the language's, not the ones that run
// code; no path in .git; each change asking its permission.
type IDE struct {
	repo *Repo
	// Name is the workspace's name the IDE shows.
	Name string
	// DraftKey is the draft of the request (its user's).
	DraftKey func(r *http.Request) string
	// Allowed says whether the request may do op on the file at path ("static/a.css";
	// "" the files as a whole).
	Allowed func(r *http.Request, op IdeOp, path string) bool
	// OnWrite, when set, is told of a change of the files allowed (before it
	// is made): who changes them.
	OnWrite func(r *http.Request)
	// HTTPClient downloads the URLs imported (fetch); one that reaches no
	// address of the server's own network (SafeHTTPClient) when not set.
	HTTPClient *http.Client

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
	s.HTTPClient = h.HTTPClient
	if s.HTTPClient == nil {
		s.HTTPClient = SafeHTTPClient()
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

// cleanPath is p as a path of the workspace: from its root, no `..` out of it.
func cleanPath(p string) string {
	return path.Clean("/" + strings.ReplaceAll(p, `\`, "/"))
}

// renameOp is what a rename of from to to asks: a move when it changes the
// folder.
func renameOp(from, to string) IdeOp {
	if path.Dir(cleanPath(from)) != path.Dir(cleanPath(to)) {
		return IdeMove
	}
	return IdeRename
}

// ideNeed is an operation asked on the file at Path (of the draft, "static/a.css").
type ideNeed struct {
	Op   IdeOp
	Path string
}

// relPath is p as a path of the draft: "static/a.css".
func relPath(p string) string { return strings.TrimPrefix(cleanPath(p), "/") }

// writeNeeds are what writing the files at paths asks, in root: IdeCreate for
// one that is not there, IdeEdit for one that is.
func writeNeeds(root string, paths ...string) (needs []ideNeed) {
	for _, p := range paths {
		op := IdeCreate
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(cleanPath(p)))); err == nil {
			op = IdeEdit
		}
		needs = append(needs, ideNeed{op, relPath(p)})
	}
	return
}

// asks are what the write of op asks, its body read (b), in root: each
// operation on the files it touches — a rename on both its paths, an import
// what it writes too.
func asks(op string, b *ideBody, root string) []ideNeed {
	switch op {
	case "file", "mkdir":
		return writeNeeds(root, b.Path)
	case "config":
		return []ideNeed{{IdeEdit, ""}}
	case "delete":
		return []ideNeed{{IdeDelete, relPath(b.Path)}}
	case "rename":
		o := renameOp(b.Path, b.To)
		return []ideNeed{{o, relPath(b.Path)}, {o, relPath(b.To)}}
	case "upload":
		var needs []ideNeed
		for _, f := range b.Files {
			needs = append(needs, ideNeed{IdeImport, relPath(f.Path)})
			needs = append(needs, writeNeeds(root, f.Path)...)
		}
		return needs
	case "fetch":
		return append([]ideNeed{{IdeImport, relPath(b.Path)}}, writeNeeds(root, b.Path)...)
	}
	return nil
}

// ServeHTTP serves /api/ide/<op> (its prefix stripped).
func (h *IDE) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	op, ok := strings.CutPrefix(r.URL.Path, "/api/ide/")
	spec, known := ideOps[op]
	if !ok || !known {
		ideError(w, http.StatusNotFound, "not available")
		return
	}
	write := r.Method != http.MethodGet
	if (write && !spec.write && !spec.read) || (!write && !spec.read) {
		ideError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	allowed := func(op IdeOp, p string) bool { return h.Allowed == nil || h.Allowed(r, op, p) }
	// a file read is asked of it; the rest, of the files as a whole
	readPath := ""
	if op == "file" && !write {
		readPath = relPath(r.URL.Query().Get("path"))
	}
	if !allowed(IdeRead, readPath) {
		ideError(w, http.StatusForbidden, "permission denied")
		return
	}
	// no path in .git: the query's, the body's
	if inGit(r.URL.Query().Get("path")) {
		ideError(w, http.StatusForbidden, "path not allowed")
		return
	}
	var b ideBody
	if r.Body != nil && write {
		body, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
		if err != nil {
			ideError(w, http.StatusBadRequest, err.Error())
			return
		}
		_ = json.Unmarshal(body, &b)
		paths := []string{b.Path, b.To}
		for _, f := range b.Files {
			paths = append(paths, f.Path)
		}
		for _, p := range paths {
			if inGit(p) {
				ideError(w, http.StatusForbidden, "path not allowed")
				return
			}
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
	}
	s, err := h.server(r)
	if err != nil {
		ideError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if write && spec.write {
		for _, need := range asks(op, &b, s.Root) {
			if !allowed(need.Op, need.Path) {
				ideError(w, http.StatusForbidden, "permission denied")
				return
			}
		}
	}
	if write && spec.write && h.OnWrite != nil {
		h.OnWrite(r)
	}
	if op == "workspace" {
		// not the path on the server; what the user may do with the files
		actions := map[string]bool{}
		for name, o := range IdeOps {
			actions[name] = allowed(o, "")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"root": "/", "name": h.Name, "openFile": "", "compute": "server",
			"actions": actions})
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
