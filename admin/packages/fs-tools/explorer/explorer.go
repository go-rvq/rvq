// Package explorer is the API of the explorer of files (<vx-files>) over a
// webdav.FileSystem — the folder of a user of the OS (webdav.Dir), the
// file system the WebDAV of fs-tools serves —: of it and nothing out of it,
// each operation as Allow says. The paths in the query (path, to) are of the
// file system ("/" its root):
//
//	GET  list?path[&dims=1]   {path, entries: [{name, dir, size, mode, mtime, …}], repo, can: {read, write}}
//	GET  read?path            {content, size, binary, truncated, writable}: a text, up to MaxRead
//	GET  download?path[&inline=1]  the file (inline: shown, nothing of it running); a folder, a tar.gz;
//	     path repeated (the items of a folder), one tar.gz of them
//	POST write?path           the body, the file's new content (made when not there)
//	POST mkdir?path           a folder (and its parents)
//	POST rename?path&to       renamed or moved
//	POST copy?path&to         copied (a folder, whole)
//	POST delete?path…         deleted (a folder, whole); path repeated, several
//	POST upload?path[&replace=1]  the files of the form ("file"; names with folders) into path
//	POST upload-part?path&name&id&offset&total  a part of a large file (parts.go);
//	     upload-done?…&total[&replace=1] ends it, upload-abort?… gives it up
//	POST extract?path[&to]    an archive extracted (Options.Extract)
//	POST fetch?path&url[&name]  a file downloaded from the internet (Options.Fetch)
//
// The handler answers /api/fs/<op>: the caller strips what is before.
package explorer

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/gif" // the dimensions of the images listed
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
	"golang.org/x/net/webdav"
)

// MaxRead is the most of a file read as a text.
const MaxRead = 2 << 20

// Prefix is where the operations are: /api/fs/<op>.
const Prefix = "/api/fs/"

// Options of an explorer.
type Options struct {
	// Allow says whether r may do its operation on p (write: it changes
	// it); nil, everything.
	Allow func(r *http.Request, p string, write bool) bool
	// Decorate adds to an entry of the folder dir what the file system does
	// not say (a link, a repository); nil, nothing.
	Decorate func(dir string, e *Entry)
	// Writable says whether the file p may be written (read's writable);
	// nil, Allow's.
	Writable func(r *http.Request, p string) bool
	// IsRepo says whether the folder p is a git repository (list's repo).
	IsRepo func(p string) bool
	// Extract allows extracting archives; Fetch, downloading files from the
	// internet (by HTTPClient: one that reaches no address of the server's
	// network).
	Extract    bool
	Fetch      bool
	HTTPClient func() *http.Client
	// Symlink makes the link p to target (a link of a tar extracted, inside
	// its folder); nil, links left out.
	Symlink func(p, target string) error
	// Changed is told of each change (its operation and path): a log.
	Changed func(r *http.Request, op, p string)
	// ReadOnly says the path p is only read whoever asks (a folder mounted
	// read only): writing it is not allowed, and list's can and read's
	// writable say so; nil, none.
	ReadOnly func(p string) bool
	// Archive writes the tar.gz of paths (of the file system), its names
	// those under base (the folder they are in); nil, the files and the
	// folders read of the file system (writeTarGz).
	Archive func(w io.Writer, r *http.Request, base string, paths []string) error
}

// Explorer serves the explorer of FS.
type Explorer struct {
	FS   webdav.FileSystem
	Opts Options
}

// New is the explorer of fsys.
func New(fsys webdav.FileSystem, opts Options) *Explorer { return &Explorer{FS: fsys, Opts: opts} }

// Entry is a file of a folder.
type Entry struct {
	Name  string    `json:"name"`
	Dir   bool      `json:"dir"`
	Size  int64     `json:"size"`
	Mode  string    `json:"mode"`
	MTime time.Time `json:"mtime"`
	// Link: a symbolic link, to it.
	Link string `json:"link,omitempty"`
	// Repo: a folder that is a git repository.
	Repo bool `json:"repo,omitempty"`
	// W, H: an image's dimensions (list?dims=1), read of its header.
	W int `json:"w,omitempty"`
	H int `json:"h,omitempty"`
}

// ErrNotAllowed is an operation Allow refused.
var ErrNotAllowed = errors.New("not allowed")

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// Fail answers err as the explorer reads it ({error}): the paths of the
// server not told.
func Fail(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	if errors.Is(err, ErrNotAllowed) || errors.Is(err, os.ErrPermission) {
		status = http.StatusForbidden
	}
	w.WriteHeader(status)
	msg := err.Error()
	var pe *fs.PathError
	if errors.As(err, &pe) {
		msg = pe.Op + ": " + pe.Err.Error()
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// Rel is p of the query as a path of the file system ("/" its root): no
// ".." out of it.
func Rel(p string) string { return path.Clean("/" + strings.ReplaceAll(p, `\`, "/")) }

func (x *Explorer) allow(r *http.Request, p string, write bool) error {
	if (write && x.Opts.ReadOnly != nil && x.Opts.ReadOnly(Rel(p))) || (x.Opts.Allow != nil && !x.Opts.Allow(r, Rel(p), write)) {
		return fmt.Errorf("%s: %w", Rel(p), ErrNotAllowed)
	}
	return nil
}

func (x *Explorer) changed(r *http.Request, op, p string) {
	if x.Opts.Changed != nil {
		x.Opts.Changed(r, op, Rel(p))
	}
}

func (x *Explorer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	op, ok := strings.CutPrefix(r.URL.Path, Prefix)
	if !ok {
		http.NotFound(w, r)
		return
	}
	get := r.Method == http.MethodGet
	q := r.URL.Query()
	p := q.Get("path")
	ctx := r.Context()
	switch {
	case op == "list" && get:
		x.list(w, r, p, q.Get("dims") == "1")
	case op == "read" && get:
		x.read(w, r, p)
	case op == "download" && get:
		x.download(w, r, q["path"])
	case op == "write" && !get:
		x.write(w, r, p)
	case op == "mkdir" && !get:
		if err := x.allow(r, p, true); err != nil {
			Fail(w, http.StatusForbidden, err)
			return
		}
		err := mkdirAll(ctx, x.FS, Rel(p))
		if err == nil {
			x.changed(r, op, p)
		}
		done(w, err)
	case op == "rename" && !get:
		x.rename(w, r, p, q.Get("to"))
	case op == "copy" && !get:
		x.copyTo(w, r, p, q.Get("to"))
	case op == "delete" && !get:
		x.delete(w, r, q["path"])
	case op == "upload" && !get:
		x.upload(w, r, p)
	case op == "upload-part" && !get:
		x.uploadPart(w, r)
	case op == "upload-done" && !get:
		x.uploadDone(w, r)
	case op == "upload-abort" && !get:
		x.uploadAbort(w, r)
	case op == "extract" && !get && x.Opts.Extract:
		x.extract(w, r, p, q.Get("to"))
	case op == "fetch" && !get && x.Opts.Fetch:
		x.fetch(w, r, p, q.Get("url"), q.Get("name"))
	default:
		Fail(w, http.StatusNotFound, errors.New("not available"))
	}
}

func done(w http.ResponseWriter, err error) {
	if err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, map[string]bool{"ok": true})
}

// mkdirAll makes the folder p and its parents.
func mkdirAll(ctx context.Context, fsys webdav.FileSystem, p string) error {
	p = Rel(p)
	if p == "/" {
		return nil
	}
	if st, err := fsys.Stat(ctx, p); err == nil {
		if st.IsDir() {
			return nil
		}
		return fmt.Errorf("%s: not a folder", p)
	}
	if err := mkdirAll(ctx, fsys, path.Dir(p)); err != nil {
		return err
	}
	if err := fsys.Mkdir(ctx, p, 0o755); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	return nil
}

// readDir are the files of the folder p.
func readDir(ctx context.Context, fsys webdav.FileSystem, p string) ([]fs.FileInfo, error) {
	f, err := fsys.OpenFile(ctx, p, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Readdir(-1)
}

// imageExt are the images whose dimensions are read (their header).
var imageExt = map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true, ".bmp": true}

func (x *Explorer) dimsOf(ctx context.Context, p string) (int, int) {
	f, err := x.FS.OpenFile(ctx, p, os.O_RDONLY, 0)
	if err != nil {
		return 0, 0
	}
	defer f.Close()
	c, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0
	}
	return c.Width, c.Height
}

func (x *Explorer) list(w http.ResponseWriter, r *http.Request, p string, dims bool) {
	if err := x.allow(r, p, false); err != nil {
		Fail(w, http.StatusForbidden, err)
		return
	}
	dir := Rel(p)
	infos, err := readDir(r.Context(), x.FS, dir)
	if err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	out := []Entry{}
	for _, fi := range infos {
		e := Entry{Name: fi.Name(), Dir: fi.IsDir(), Size: fi.Size(), Mode: fi.Mode().String(), MTime: fi.ModTime()}
		if e.Dir {
			e.Size = 0
		}
		if x.Opts.Decorate != nil {
			x.Opts.Decorate(dir, &e)
		}
		if dims && !e.Dir && imageExt[strings.ToLower(path.Ext(e.Name))] {
			e.W, e.H = x.dimsOf(r.Context(), path.Join(dir, e.Name))
		}
		out = append(out, e)
	}
	// the folders first, then by name
	sort.Slice(out, func(i, j int) bool {
		if out[i].Dir != out[j].Dir {
			return out[i].Dir
		}
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	repo := x.Opts.IsRepo != nil && x.Opts.IsRepo(dir)
	// what r may do in it: the explorer shows what it may (reading — the
	// downloads —, writing)
	can := map[string]bool{"read": true, "write": x.allow(r, dir, true) == nil}
	writeJSON(w, map[string]any{"path": dir, "entries": out, "repo": repo, "can": can})
}

func (x *Explorer) read(w http.ResponseWriter, r *http.Request, p string) {
	if err := x.allow(r, p, false); err != nil {
		Fail(w, http.StatusForbidden, err)
		return
	}
	f, err := x.FS.OpenFile(r.Context(), Rel(p), os.O_RDONLY, 0)
	if err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil || st.IsDir() {
		Fail(w, http.StatusBadRequest, errors.New("not a file"))
		return
	}
	b, err := io.ReadAll(io.LimitReader(f, MaxRead))
	if err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	text := b
	if st.Size() > MaxRead {
		text = trimRune(b)
	}
	binary := bytes.IndexByte(b, 0) >= 0 || !utf8.Valid(text)
	writable := x.allow(r, p, true) == nil
	if writable && x.Opts.Writable != nil {
		writable = x.Opts.Writable(r, Rel(p))
	}
	res := map[string]any{"size": st.Size(), "binary": binary, "truncated": st.Size() > MaxRead, "mode": st.Mode().String(),
		"writable": writable}
	if !binary {
		res["content"] = string(b)
	}
	writeJSON(w, res)
}

// trimRune is b without a rune cut at its end: up to 3 bytes off while it
// is not valid.
func trimRune(b []byte) []byte {
	for i := 0; i < 3 && len(b) > 0 && !utf8.Valid(b); i++ {
		b = b[:len(b)-1]
	}
	return b
}

func (x *Explorer) download(w http.ResponseWriter, r *http.Request, ps []string) {
	if len(ps) == 0 {
		ps = []string{"/"}
	}
	var files []string
	for _, p := range ps {
		if err := x.allow(r, p, false); err != nil {
			Fail(w, http.StatusForbidden, err)
			return
		}
		files = append(files, Rel(p))
	}
	ctx := r.Context()
	file := files[0]
	st, err := x.FS.Stat(ctx, file)
	if err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	name := path.Base(file)
	if file == "/" {
		name = "files"
	}
	if len(files) == 1 && !st.IsDir() {
		f, err := x.FS.OpenFile(ctx, file, os.O_RDONLY, 0)
		if err != nil {
			Fail(w, http.StatusBadRequest, err)
			return
		}
		defer f.Close()
		if r.URL.Query().Get("inline") == "1" {
			// shown (an image in the page): nothing of it runs — a script
			// of an SVG, its links —, its type not guessed
			w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src 'self' data:; style-src 'unsafe-inline'; sandbox")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", name))
		} else {
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", name))
		}
		http.ServeContent(w, r, name, st.ModTime(), f)
		return
	}
	// a folder, or the items of one: a tar.gz, made as it is sent — named
	// as the folder; the items' names under the folder they are in
	base := path.Dir(file)
	if file == "/" {
		base = "/"
	}
	for _, f := range files[1:] {
		if path.Dir(f) != base {
			Fail(w, http.StatusBadRequest, fmt.Errorf("%s: not in %s", f, base))
			return
		}
		if _, err := x.FS.Stat(ctx, f); err != nil {
			Fail(w, http.StatusBadRequest, err)
			return
		}
	}
	if len(files) > 1 {
		if name = path.Base(base); base == "/" {
			name = "files"
		}
	}
	w.Header().Set("Content-Type", "application/gzip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", name+".tar.gz"))
	if x.Opts.Archive != nil {
		_ = x.Opts.Archive(w, r, base, files)
		return
	}
	_ = x.writeTarGz(ctx, w, base, files)
}

// writeTarGz writes the tar.gz of paths, read of the file system — the
// folders and the regular files, their permissions and times —, the names
// those under base.
func (x *Explorer) writeTarGz(ctx context.Context, w io.Writer, base string, paths []string) error {
	gz := gzip.NewWriter(w)
	tw := tar.NewWriter(gz)
	for _, p := range paths {
		err := x.walk(ctx, p, func(fp string, fi fs.FileInfo) error {
			name := strings.TrimPrefix(strings.TrimPrefix(fp, base), "/")
			if name == "" || !(fi.IsDir() || fi.Mode().IsRegular()) {
				return nil
			}
			hd := &tar.Header{Name: name, Mode: int64(fi.Mode().Perm()), ModTime: fi.ModTime(), Typeflag: tar.TypeReg, Size: fi.Size(), Format: tar.FormatPAX}
			if fi.IsDir() {
				hd.Name, hd.Typeflag, hd.Size = name+"/", tar.TypeDir, 0
				return tw.WriteHeader(hd)
			}
			f, err := x.FS.OpenFile(ctx, fp, os.O_RDONLY, 0)
			if err != nil {
				return nil
			}
			defer f.Close()
			if err := tw.WriteHeader(hd); err != nil {
				return err
			}
			_, err = io.CopyN(tw, f, fi.Size())
			return err
		})
		if err != nil {
			return err
		}
	}
	if err := tw.Close(); err != nil {
		return err
	}
	return gz.Close()
}

// walk calls fn of p and of every file under it (a folder before its
// files).
func (x *Explorer) walk(ctx context.Context, p string, fn func(p string, fi fs.FileInfo) error) error {
	st, err := x.FS.Stat(ctx, p)
	if err != nil {
		return err
	}
	if err := fn(p, st); err != nil || !st.IsDir() {
		return err
	}
	infos, err := readDir(ctx, x.FS, p)
	if err != nil {
		return err
	}
	for _, fi := range infos {
		if err := x.walk(ctx, path.Join(p, fi.Name()), fn); err != nil {
			return err
		}
	}
	return nil
}

// writeFile writes r as the file p: whole or nothing (a file beside, then
// renamed).
func (x *Explorer) writeFile(ctx context.Context, p string, r io.Reader, mode fs.FileMode) error {
	tmp := path.Join(path.Dir(p), fmt.Sprintf(".rvq-write-%d", time.Now().UnixNano()))
	f, err := x.FS.OpenFile(ctx, tmp, os.O_CREATE|os.O_WRONLY|os.O_EXCL, mode)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = x.FS.Rename(ctx, tmp, p)
	}
	if err != nil {
		_ = x.FS.RemoveAll(ctx, tmp)
	}
	return err
}

func (x *Explorer) write(w http.ResponseWriter, r *http.Request, p string) {
	file := Rel(p)
	if file == "/" {
		Fail(w, http.StatusBadRequest, errors.New("not a file"))
		return
	}
	if err := x.allow(r, file, true); err != nil {
		Fail(w, http.StatusForbidden, err)
		return
	}
	mode := fs.FileMode(0o644)
	if st, err := x.FS.Stat(r.Context(), file); err == nil {
		if st.IsDir() {
			Fail(w, http.StatusBadRequest, errors.New("a folder"))
			return
		}
		mode = st.Mode().Perm()
	}
	err := x.writeFile(r.Context(), file, r.Body, mode)
	if err == nil {
		x.changed(r, "write", file)
	}
	done(w, err)
}

func (x *Explorer) twoPaths(w http.ResponseWriter, r *http.Request, from, to string) (string, string, bool) {
	a, b := Rel(from), Rel(to)
	if a == "/" || b == "/" || to == "" {
		Fail(w, http.StatusBadRequest, errors.New("not the root"))
		return "", "", false
	}
	if err := x.allow(r, a, true); err != nil {
		Fail(w, http.StatusForbidden, err)
		return "", "", false
	}
	if err := x.allow(r, b, true); err != nil {
		Fail(w, http.StatusForbidden, err)
		return "", "", false
	}
	if _, err := x.FS.Stat(r.Context(), b); err == nil {
		Fail(w, http.StatusConflict, fmt.Errorf("%s: already there", b))
		return "", "", false
	}
	return a, b, true
}

func (x *Explorer) rename(w http.ResponseWriter, r *http.Request, from, to string) {
	a, b, ok := x.twoPaths(w, r, from, to)
	if !ok {
		return
	}
	err := x.FS.Rename(r.Context(), a, b)
	if err == nil {
		x.changed(r, "rename", a+" → "+b)
	}
	done(w, err)
}

func (x *Explorer) copyTo(w http.ResponseWriter, r *http.Request, from, to string) {
	a, b, ok := x.twoPaths(w, r, from, to)
	if !ok {
		return
	}
	if b == a || strings.HasPrefix(b, a+"/") {
		Fail(w, http.StatusBadRequest, errors.New("a folder into itself"))
		return
	}
	ctx := r.Context()
	err := x.walk(ctx, a, func(fp string, fi fs.FileInfo) error {
		dst := b + strings.TrimPrefix(fp, a)
		if fi.IsDir() {
			return mkdirAll(ctx, x.FS, dst)
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		src, err := x.FS.OpenFile(ctx, fp, os.O_RDONLY, 0)
		if err != nil {
			return err
		}
		defer src.Close()
		out, err := x.FS.OpenFile(ctx, dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, fi.Mode().Perm())
		if err != nil {
			return err
		}
		if _, err := io.Copy(out, src); err != nil {
			out.Close()
			return err
		}
		return out.Close()
	})
	if err == nil {
		x.changed(r, "copy", a+" → "+b)
	}
	done(w, err)
}

func (x *Explorer) delete(w http.ResponseWriter, r *http.Request, paths []string) {
	if len(paths) == 0 {
		Fail(w, http.StatusBadRequest, errors.New("no path"))
		return
	}
	for _, p := range paths {
		if Rel(p) == "/" {
			Fail(w, http.StatusBadRequest, errors.New("not the root"))
			return
		}
		if err := x.allow(r, p, true); err != nil {
			Fail(w, http.StatusForbidden, err)
			return
		}
	}
	for _, p := range paths {
		if err := x.FS.RemoveAll(r.Context(), Rel(p)); err != nil {
			Fail(w, http.StatusBadRequest, err)
			return
		}
		x.changed(r, "delete", p)
	}
	writeJSON(w, map[string]bool{"ok": true})
}

// upload writes the files of the form into the folder dir: each its name
// (with its folders, a folder uploaded), never one that is there unless
// ?replace=1.
func (x *Explorer) upload(w http.ResponseWriter, r *http.Request, dir string) {
	mr, err := r.MultipartReader()
	if err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()
	replace := r.URL.Query().Get("replace") == "1"
	var wrote []string
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			Fail(w, http.StatusBadRequest, err)
			return
		}
		if part.FormName() != "file" || part.FileName() == "" {
			continue
		}
		// the browser gives a folder's files their path (webkitRelativePath,
		// sent as the name): the raw name — FileName keeps its base only
		raw := part.FileName()
		if _, ps, err := mime.ParseMediaType(part.Header.Get("Content-Disposition")); err == nil && ps["filename"] != "" {
			raw = ps["filename"]
		}
		name := Rel(path.Join(Rel(dir), Rel(raw)))
		if name == "/" {
			continue
		}
		if err := x.allow(r, name, true); err != nil {
			Fail(w, http.StatusForbidden, err)
			return
		}
		if err := mkdirAll(ctx, x.FS, path.Dir(name)); err != nil {
			Fail(w, http.StatusBadRequest, err)
			return
		}
		if !replace {
			if _, err := x.FS.Stat(ctx, name); err == nil {
				Fail(w, http.StatusConflict, fmt.Errorf("%s: already there", name))
				return
			}
		}
		if err := x.writeFile(ctx, name, part, 0o644); err != nil {
			Fail(w, http.StatusBadRequest, err)
			return
		}
		x.changed(r, "upload", name)
		wrote = append(wrote, name)
	}
	writeJSON(w, map[string]any{"ok": true, "files": wrote})
}
