package explorer

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// A file of any size, sent in parts: each part a request of its own — short,
// so no proxy's timeout cuts it, retried alone, several at once —, written
// at its place in a file of its own (.rvq-upload-<id>, in the folder of the
// file); the last request, once every part was, says it is done: its size
// checked, it is renamed to the file's name, never over one that is there
// unless replace=1. A file of parts left is removed after PartsMaxAge.
//
//	POST upload-part?path&name&id&offset&total   the body, the part at offset
//	POST upload-done?path&name&id&total[&replace=1]
//	POST upload-abort?path&name&id

// PartsPrefix is the name of the file of the parts of a sending, before its
// id.
const PartsPrefix = ".rvq-upload-"

// PartsMaxAge is how long a file of parts not touched stays.
const PartsMaxAge = 24 * time.Hour

var uploadID = regexp.MustCompile(`^[A-Za-z0-9_-]{16,64}$`)

// ErrPartsSize is a sending done with a size other than its total.
var ErrPartsSize = errors.New("the parts received do not make the file's size")

// partsFile is the file of the parts of the sending id of the file name in
// the folder dir, and the file's.
func partsFile(dir, name, id string) (parts, file string, err error) {
	if !uploadID.MatchString(id) {
		return "", "", errors.New("not an id of a sending")
	}
	file = Rel(path.Join(Rel(dir), Rel(name)))
	if file == "/" || strings.HasSuffix(name, "/") {
		return "", "", errors.New("not a name of a file")
	}
	return path.Join(path.Dir(file), PartsPrefix+id), file, nil
}

func int64Of(q map[string][]string, key string) int64 {
	v, ok := q[key]
	if !ok || len(v) == 0 {
		return -1
	}
	n, err := strconv.ParseInt(v[0], 10, 64)
	if err != nil || n < 0 {
		return -1
	}
	return n
}

func (x *Explorer) uploadPart(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	parts, file, err := partsFile(q.Get("path"), q.Get("name"), q.Get("id"))
	if err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	if err := x.allow(r, file, true); err != nil {
		Fail(w, http.StatusForbidden, err)
		return
	}
	offset, total := int64Of(q, "offset"), int64Of(q, "total")
	if offset < 0 || total < 0 || offset > total {
		Fail(w, http.StatusBadRequest, errors.New("offset and total: 0 ≤ offset ≤ total"))
		return
	}
	ctx := r.Context()
	if err := mkdirAll(ctx, x.FS, path.Dir(parts)); err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	x.cleanParts(r, path.Dir(parts))
	f, err := x.FS.OpenFile(ctx, parts, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	var n int64
	if _, err = f.Seek(offset, io.SeekStart); err == nil {
		// the part, never past the file's end: one byte more is too much
		n, err = io.Copy(f, io.LimitReader(r.Body, total-offset+1))
	}
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	switch {
	case err != nil:
		Fail(w, http.StatusBadRequest, err)
	case offset+n > total:
		Fail(w, http.StatusBadRequest, fmt.Errorf("the part goes past the file's size (%d)", total))
	default:
		writeJSON(w, map[string]any{"ok": true, "wrote": n})
	}
}

func (x *Explorer) uploadDone(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	parts, file, err := partsFile(q.Get("path"), q.Get("name"), q.Get("id"))
	if err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	if err := x.allow(r, file, true); err != nil {
		Fail(w, http.StatusForbidden, err)
		return
	}
	ctx := r.Context()
	total := int64Of(q, "total")
	st, err := x.FS.Stat(ctx, parts)
	if total == 0 && errors.Is(err, os.ErrNotExist) {
		// an empty file: no part sent
		var f interface{ Close() error }
		if f, err = x.FS.OpenFile(ctx, parts, os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
			f.Close()
			st, err = x.FS.Stat(ctx, parts)
		}
	}
	if err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	if st.Size() != total {
		Fail(w, http.StatusBadRequest, fmt.Errorf("%w: %d of %d", ErrPartsSize, st.Size(), total))
		return
	}
	if q.Get("replace") != "1" {
		if _, err := x.FS.Stat(ctx, file); err == nil {
			Fail(w, http.StatusConflict, fmt.Errorf("%s: already there", q.Get("name")))
			return
		}
	}
	if err := mkdirAll(ctx, x.FS, path.Dir(file)); err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	if err := x.FS.Rename(ctx, parts, file); err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	x.changed(r, "upload", file)
	writeJSON(w, map[string]any{"ok": true, "files": []string{file}})
}

func (x *Explorer) uploadAbort(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	parts, file, err := partsFile(q.Get("path"), q.Get("name"), q.Get("id"))
	if err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	if err := x.allow(r, file, true); err != nil {
		Fail(w, http.StatusForbidden, err)
		return
	}
	if err := x.FS.RemoveAll(r.Context(), parts); err != nil && !errors.Is(err, os.ErrNotExist) {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, map[string]bool{"ok": true})
}

// cleanParts removes, of the folder dir, the files of parts not touched for
// PartsMaxAge: sendings given up.
func (x *Explorer) cleanParts(r *http.Request, dir string) {
	infos, err := readDir(r.Context(), x.FS, dir)
	if err != nil {
		return
	}
	for _, fi := range infos {
		if !fi.IsDir() && strings.HasPrefix(fi.Name(), PartsPrefix) && time.Since(fi.ModTime()) > PartsMaxAge {
			_ = x.FS.RemoveAll(r.Context(), path.Join(dir, fi.Name()))
		}
	}
}
