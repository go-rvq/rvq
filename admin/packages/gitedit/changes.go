package gitedit

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// The changes of a draft, browsed and edited (vx-diff-browser): by the page
// of the editor (its events) and by the IDE (its API, git/…), the same.

// errSaveDenied is a save the user may not make: a path not of the draft as
// its summary gives it, or one they may not write.
var errSaveDenied = errors.New("you may not save this file")

// diffSummary are the changes of d as the diff browser's tree has them —
// path, status, origin of a renamed one, language —, each read only when
// allowed says it may not be written (as the IDE writes it: by the
// permission of its path, IdeEdit or IdeCreate).
func diffSummary(ctx context.Context, d *Draft, allowed func(op IdeOp, path string) bool) ([]vx.DiffFile, error) {
	changes, err := d.Changes(ctx)
	if err != nil {
		return nil, err
	}
	files := make([]vx.DiffFile, 0, len(changes))
	for _, ch := range changes {
		f := vx.DiffSummary(ch.Path, ch.Status, ch.From, CodeLanguage(ch.Path))
		f.ReadOnly = !allowed(writeOp(d, ch.Path), ch.Path)
		files = append(files, f)
	}
	return files, nil
}

// diffOf is the diff of the change of d at path (vuetifyx.NewDiffFile); nil
// when path is not one of its changes now (committed or discarded meanwhile).
func diffOf(ctx context.Context, d *Draft, path string) (*vx.DiffFile, error) {
	changes, err := d.Changes(ctx)
	if err != nil {
		return nil, err
	}
	for _, ch := range changes {
		if ch.Path != path {
			continue
		}
		old, cur, err := d.Versions(ctx, ch)
		if err != nil {
			return nil, err
		}
		f := vx.NewDiffFile(ch.Path, ch.Status, CodeLanguage(ch.Path), old, cur)
		f.From = ch.From
		return &f, nil
	}
	return nil, nil
}

// saveChange writes the file path of d with text, as the diff browser saves
// it: path as the summary gives it — clean, in the draft (no ".."), not in
// .git —, the user allowed to write it (IdeEdit, or IdeCreate for one not
// there), its line endings kept (a form sends CRLF: \n, unless the file has
// \r\n). errSaveDenied when not allowed.
func saveChange(d *Draft, path, text string, allowed func(op IdeOp, path string) bool) error {
	if path == "" || relPath(path) != path || d.check(path) != nil || inGit(path) {
		return errSaveDenied
	}
	if !allowed(writeOp(d, path), path) {
		return errSaveDenied
	}
	if was, err := os.ReadFile(filepath.Join(d.Dir, filepath.FromSlash(path))); err != nil || !strings.Contains(string(was), "\r\n") {
		text = strings.ReplaceAll(text, "\r\n", "\n")
	}
	return d.WriteFile(path, text)
}

// writeOp is what writing the file path of d asks: IdeEdit for one that is
// there, IdeCreate for one that is not.
func writeOp(d *Draft, path string) IdeOp {
	if _, err := os.Stat(filepath.Join(d.Dir, filepath.FromSlash(path))); err == nil {
		return IdeEdit
	}
	return IdeCreate
}

// serveGit serves the draft's repository to the IDE, under api/ide/git/. Its
// changes — the Changes panel — are ours, as the editor's page has them
// (renames found, the files the user may not write read only):
//
//	GET  git/changes            {"files": [the summary]}
//	GET  git/diff?path          the diff of a change; 404 when it is no more one
//	POST git/save               {"path", "content"}: the file written ({"saved": true})
//
// its branches, commits and their files — the Git panel, read only — gad's
// IDE's (web/ide git.go), on the draft: git/branches, git/log, git/commit,
// git/commit/diff, git/file, git/patch. Each as the IDE's own operations
// ask: seeing the files (a file's, the one given), writing the file.
func (h *IDE) serveGit(w http.ResponseWriter, r *http.Request, op string) {
	allowed := func(o IdeOp, p string) bool { return h.Allowed == nil || h.Allowed(r, o, p) }
	if !allowed(IdeRead, "") {
		ideError(w, http.StatusForbidden, "permission denied")
		return
	}
	switch op {
	case "git/changes", "git/diff", "git/save":
	case "git/branches", "git/log", "git/commit", "git/commit/diff", "git/file", "git/patch":
		if r.Method != http.MethodGet {
			ideError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		for _, p := range []string{r.URL.Query().Get("path"), r.URL.Query().Get("from")} {
			if p == "" {
				continue
			}
			if inGit(p) || !allowed(IdeRead, relPath(p)) {
				ideError(w, http.StatusForbidden, "permission denied")
				return
			}
		}
		s, err := h.server(r)
		if err != nil {
			ideError(w, http.StatusInternalServerError, err.Error())
			return
		}
		s.Handler().ServeHTTP(w, r)
		return
	default:
		ideError(w, http.StatusNotFound, "not available")
		return
	}
	d, err := h.repo.Draft(r.Context(), h.DraftKey(r))
	if err != nil {
		ideError(w, http.StatusInternalServerError, err.Error())
		return
	}
	reply := func(v any) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
	switch {
	case op == "git/changes" && r.Method == http.MethodGet:
		files, err := diffSummary(r.Context(), d, allowed)
		if err != nil {
			ideError(w, http.StatusInternalServerError, err.Error())
			return
		}
		reply(map[string]any{"files": files})
	case op == "git/diff" && r.Method == http.MethodGet:
		f, err := diffOf(r.Context(), d, r.URL.Query().Get("path"))
		switch {
		case err != nil:
			ideError(w, http.StatusInternalServerError, err.Error())
		case f == nil:
			ideError(w, http.StatusNotFound, "this file is no longer changed")
		default:
			reply(f)
		}
	case op == "git/save" && r.Method == http.MethodPost:
		var body struct{ Path, Content string }
		if err := json.NewDecoder(io.LimitReader(r.Body, maxBody)).Decode(&body); err != nil {
			ideError(w, http.StatusBadRequest, err.Error())
			return
		}
		err := saveChange(d, body.Path, body.Content, allowed)
		switch {
		case errors.Is(err, errSaveDenied):
			ideError(w, http.StatusForbidden, err.Error())
		case err != nil:
			ideError(w, http.StatusInternalServerError, err.Error())
		default:
			if h.OnWrite != nil {
				h.OnWrite(r)
			}
			reply(map[string]any{"saved": true})
		}
	default:
		ideError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
