package fs_tools

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"golang.org/x/net/webdav"
)

// Mount is a directory of each request, shown in the WebDAV at /Name: the
// draft of the user's files, say. Dir is the directory of the request;
// Allowed says whether the request may do what it asks (MountAccess: the
// path, the destination of a COPY or a MOVE, its method in r). Hidden paths
// (".git") are not served.
type Mount struct {
	Name    string
	Dir     func(r *http.Request) (string, error)
	Allowed func(r *http.Request, a MountAccess) bool
	// Hidden says whether the path, in the mount ("/.git/config"), is not
	// served.
	Hidden func(p string) bool
}

// MountAccess is what a request asks of a mount: its path in it ("/a.css"),
// where a COPY or a MOVE puts it, in the same mount ("" otherwise), and whether
// it changes the files (Write; r.Method says how).
type MountAccess struct {
	Path  string
	Dest  string
	Write bool
}

// AddMount adds a directory of each request to the WebDAV.
func (b *Builder) AddMount(m *Mount) *Builder {
	b.mounts = append(b.mounts, m)
	return b
}

// writeMethods are the methods of WebDAV that change the files.
var writeMethods = map[string]bool{
	http.MethodPut: true, http.MethodDelete: true, http.MethodPost: true, http.MethodPatch: true,
	"MKCOL": true, "COPY": true, "MOVE": true, "PROPPATCH": true, "LOCK": true, "UNLOCK": true,
}

type davRequestKey struct{}

// split is the mount of p (a path in the WebDAV, "/site-files/a.css") and the
// path in it ("/a.css"); nil when p is not in a mount.
func (b *Builder) split(p string) (*Mount, string) {
	p = path.Clean("/" + p)
	first, rest, _ := strings.Cut(strings.TrimPrefix(p, "/"), "/")
	for _, m := range b.mounts {
		if m.Name == first {
			return m, "/" + rest
		}
	}
	return nil, p
}

// allowed says whether r may do its method on p (to dest, of a COPY or a
// MOVE: "" otherwise): in a mount, as it says — once, both paths told it, when
// they are in the same one —; else by the page's permissions — reading @get,
// writing !write.
func (b *Builder) allowed(r *http.Request, p, dest string, write bool) bool {
	if m, rest := b.split(p); m != nil {
		if m.Hidden != nil && m.Hidden(rest) {
			return false
		}
		a := MountAccess{Path: rest, Write: write}
		if dest != "" {
			if dm, drest := b.split(dest); dm == m {
				if m.Hidden != nil && m.Hidden(drest) {
					return false
				}
				a.Dest = drest
			} else if !b.allowed(r, dest, "", write) {
				return false
			}
		}
		return m.Allowed == nil || m.Allowed(r, a)
	}
	if dest != "" && !b.allowed(r, dest, "", write) {
		return false
	}
	if b.page == nil {
		return true
	}
	perm := "@get"
	if write {
		perm = "!write"
	}
	ver := b.page.Page().ActionVerifier(r, perm)
	return ver == nil || ver.Allowed()
}

// permissions verifies each request of the WebDAV — its path, and the
// destination of a COPY or a MOVE — before it is served.
func (b *Builder) permissions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		write := writeMethods[r.Method]
		p, dest := strings.TrimPrefix(r.URL.Path, b.davPath), ""
		if d := r.Header.Get("Destination"); d != "" {
			if u, err := url.Parse(d); err == nil {
				dest = strings.TrimPrefix(u.Path, b.davPath)
				write = true
			}
		}
		if !b.allowed(r, p, dest, write) {
			http.Error(w, "permission denied", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), davRequestKey{}, r)))
	})
}

// mountedFS is the WebDAV's files: the base's, and the mounts of the request
// at the root.
type mountedFS struct {
	b    *Builder
	base webdav.FileSystem
}

var errCrossMount = errors.New("between different mounts")

func (f *mountedFS) route(ctx context.Context, name string) (webdav.FileSystem, string, *Mount, error) {
	m, rest := f.b.split(name)
	if m == nil {
		return f.base, name, nil, nil
	}
	r, _ := ctx.Value(davRequestKey{}).(*http.Request)
	if r == nil {
		return nil, "", nil, os.ErrPermission
	}
	dir, err := m.Dir(r)
	if err != nil {
		return nil, "", nil, err
	}
	return webdav.Dir(dir), rest, m, nil
}

func (f *mountedFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	fsys, rest, _, err := f.route(ctx, name)
	if err != nil {
		return err
	}
	return fsys.Mkdir(ctx, rest, perm)
}

func (f *mountedFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	fsys, rest, m, err := f.route(ctx, name)
	if err != nil {
		return nil, err
	}
	file, err := fsys.OpenFile(ctx, rest, flag, perm)
	if err != nil {
		return nil, err
	}
	if m == nil && path.Clean("/"+name) == "/" {
		// the root: the mounts beside the base's
		return &rootFile{File: file, mounts: f.b.mounts}, nil
	}
	if m != nil {
		return &hidingFile{File: file, mount: m, dir: rest}, nil
	}
	return file, nil
}

func (f *mountedFS) RemoveAll(ctx context.Context, name string) error {
	fsys, rest, m, err := f.route(ctx, name)
	if err != nil {
		return err
	}
	if m != nil && rest == "/" {
		return os.ErrPermission // the mount itself
	}
	return fsys.RemoveAll(ctx, rest)
}

func (f *mountedFS) Rename(ctx context.Context, oldName, newName string) error {
	fromFS, from, fromM, err := f.route(ctx, oldName)
	if err != nil {
		return err
	}
	_, to, toM, err := f.route(ctx, newName)
	if err != nil {
		return err
	}
	if fromM != toM || (fromM != nil && (from == "/" || to == "/")) {
		return errCrossMount
	}
	return fromFS.Rename(ctx, from, to)
}

func (f *mountedFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	fsys, rest, m, err := f.route(ctx, name)
	if err != nil {
		return nil, err
	}
	fi, err := fsys.Stat(ctx, rest)
	if err == nil && m != nil && rest == "/" {
		fi = namedInfo{FileInfo: fi, name: m.Name}
	}
	return fi, err
}

// rootFile is the root of the WebDAV: the base's entries and the mounts.
type rootFile struct {
	webdav.File
	mounts []*Mount
	done   bool
}

func (f *rootFile) Readdir(count int) ([]fs.FileInfo, error) {
	infos, err := f.File.Readdir(count)
	if f.done || (err != nil && len(infos) == 0 && count > 0) {
		return infos, err
	}
	f.done = true
	for _, m := range f.mounts {
		infos = append(infos, dirInfo(m.Name))
	}
	return infos, nil
}

// hidingFile is a directory of a mount, its hidden entries left out.
type hidingFile struct {
	webdav.File
	mount *Mount
	dir   string
}

func (f *hidingFile) Readdir(count int) ([]fs.FileInfo, error) {
	infos, err := f.File.Readdir(count)
	if f.mount.Hidden == nil {
		return infos, err
	}
	out := infos[:0]
	for _, fi := range infos {
		if !f.mount.Hidden(path.Join(f.dir, fi.Name())) {
			out = append(out, fi)
		}
	}
	return out, err
}

type namedInfo struct {
	fs.FileInfo
	name string
}

func (i namedInfo) Name() string { return i.name }

// dirInfo is a directory named name.
type dirInfo string

func (d dirInfo) Name() string       { return string(d) }
func (d dirInfo) Size() int64        { return 0 }
func (d dirInfo) Mode() fs.FileMode  { return fs.ModeDir | 0o755 }
func (d dirInfo) ModTime() time.Time { return time.Time{} }
func (d dirInfo) IsDir() bool        { return true }
func (d dirInfo) Sys() any           { return nil }
