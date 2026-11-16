package explorer

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// extract: POST extract?path[&to] — the archive path (.zip, .tar, .tar.gz,
// .tgz, .gz) extracted into the folder to (by default, beside it, a folder
// of its name). Nothing out of to: a name with ".." or from the root, a link
// — left out; a file that is there, not replaced.
func (x *Explorer) extract(w http.ResponseWriter, r *http.Request, p, to string) {
	file := Rel(p)
	base := path.Base(file)
	kind := archiveKind(base)
	if kind == "" {
		Fail(w, http.StatusBadRequest, errors.New("not an archive (.zip, .tar, .tar.gz, .tgz, .gz)"))
		return
	}
	if to == "" {
		to = path.Join(path.Dir(file), strings.TrimSuffix(strings.TrimSuffix(base, archiveExt(base)), "."))
	}
	dst := Rel(to)
	if err := x.allow(r, file, false); err != nil {
		Fail(w, http.StatusForbidden, err)
		return
	}
	if err := x.allow(r, dst, true); err != nil {
		Fail(w, http.StatusForbidden, err)
		return
	}
	ctx := r.Context()
	if err := mkdirAll(ctx, x.FS, dst); err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	var n int
	var err error
	switch kind {
	case "zip":
		n, err = x.extractZip(ctx, file, dst)
	case "tar", "tgz":
		n, err = x.extractTar(ctx, file, dst, kind == "tgz")
	case "gz":
		n, err = x.extractGz(ctx, file, path.Join(dst, strings.TrimSuffix(base, ".gz")))
	}
	if err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	x.changed(r, "extract", file+" → "+dst)
	writeJSON(w, map[string]any{"ok": true, "to": dst, "files": n})
}

func archiveKind(name string) string {
	n := strings.ToLower(name)
	switch {
	case strings.HasSuffix(n, ".zip"):
		return "zip"
	case strings.HasSuffix(n, ".tar.gz"), strings.HasSuffix(n, ".tgz"):
		return "tgz"
	case strings.HasSuffix(n, ".tar"):
		return "tar"
	case strings.HasSuffix(n, ".gz"):
		return "gz"
	}
	return ""
}

func archiveExt(name string) string {
	n := strings.ToLower(name)
	for _, e := range []string{".tar.gz", ".tgz", ".tar", ".zip", ".gz"} {
		if strings.HasSuffix(n, e) {
			return name[len(name)-len(e):]
		}
	}
	return ""
}

// inside is the file of the archive's name in dst; "" when it would be out
// of it.
func inside(dst, name string) string {
	name = strings.ReplaceAll(name, `\`, "/")
	if path.IsAbs(name) {
		return ""
	}
	c := path.Clean(name)
	if c == "." || c == ".." || strings.HasPrefix(c, "../") {
		return ""
	}
	return path.Join(dst, c)
}

// create makes the file f (its folders), not one that is there.
func (x *Explorer) create(ctx context.Context, f string, mode fs.FileMode, r io.Reader) error {
	if err := mkdirAll(ctx, x.FS, path.Dir(f)); err != nil {
		return err
	}
	out, err := x.FS.OpenFile(ctx, f, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode.Perm()|0o600)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, r)
	if cerr := out.Close(); err == nil {
		err = cerr
	}
	return err
}

func (x *Explorer) extractZip(ctx context.Context, file, dst string) (int, error) {
	// a zip is read at its end first: a copy here, read at will
	src, err := x.FS.OpenFile(ctx, file, os.O_RDONLY, 0)
	if err != nil {
		return 0, err
	}
	tmp, err := os.CreateTemp("", "rvq-extract-*.zip")
	if err != nil {
		src.Close()
		return 0, err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()
	_, err = io.Copy(tmp, src)
	src.Close()
	if err != nil {
		return 0, err
	}
	st, err := tmp.Stat()
	if err != nil {
		return 0, err
	}
	zr, err := zip.NewReader(tmp, st.Size())
	if err != nil {
		return 0, err
	}
	n := 0
	for _, zf := range zr.File {
		f := inside(dst, zf.Name)
		if f == "" {
			continue
		}
		if zf.FileInfo().IsDir() {
			if err := mkdirAll(ctx, x.FS, f); err != nil {
				return n, err
			}
			continue
		}
		if !zf.Mode().IsRegular() {
			continue // links and the rest: left out
		}
		rc, err := zf.Open()
		if err != nil {
			return n, err
		}
		err = x.create(ctx, f, zf.Mode(), rc)
		rc.Close()
		if err != nil {
			return n, fmt.Errorf("%s: %w", zf.Name, err)
		}
		n++
	}
	return n, nil
}

func (x *Explorer) extractTar(ctx context.Context, file, dst string, gz bool) (int, error) {
	fh, err := x.FS.OpenFile(ctx, file, os.O_RDONLY, 0)
	if err != nil {
		return 0, err
	}
	defer fh.Close()
	var r io.Reader = fh
	if gz {
		zr, err := gzip.NewReader(fh)
		if err != nil {
			return 0, err
		}
		defer zr.Close()
		r = zr
	}
	tr := tar.NewReader(r)
	n := 0
	for {
		hd, err := tr.Next()
		if err == io.EOF {
			return n, nil
		}
		if err != nil {
			return n, err
		}
		f := inside(dst, hd.Name)
		if f == "" {
			continue
		}
		switch hd.Typeflag {
		case tar.TypeDir:
			if err := mkdirAll(ctx, x.FS, f); err != nil {
				return n, err
			}
		case tar.TypeReg:
			if err := x.create(ctx, f, fs.FileMode(hd.Mode), tr); err != nil {
				return n, fmt.Errorf("%s: %w", hd.Name, err)
			}
			n++
		case tar.TypeSymlink:
			// a link inside dst only (relative, not leaving it)
			t := hd.Linkname
			if x.Opts.Symlink == nil || path.IsAbs(t) || inside(path.Dir(f), t) == "" {
				continue
			}
			if target := path.Join(path.Dir(f), t); !strings.HasPrefix(target, dst+"/") {
				continue
			}
			_ = mkdirAll(ctx, x.FS, path.Dir(f))
			if err := x.Opts.Symlink(f, t); err == nil {
				n++
			}
		}
		// the rest: left out
	}
}

func (x *Explorer) extractGz(ctx context.Context, file, out string) (int, error) {
	fh, err := x.FS.OpenFile(ctx, file, os.O_RDONLY, 0)
	if err != nil {
		return 0, err
	}
	defer fh.Close()
	zr, err := gzip.NewReader(fh)
	if err != nil {
		return 0, err
	}
	defer zr.Close()
	if err := x.create(ctx, out, 0o644, zr); err != nil {
		return 0, err
	}
	return 1, nil
}

// fetch: POST fetch?path&url[&name] — the URL downloaded into the folder
// path, as name (by default the server's, else the URL's); never over a file
// that is there. By Options.HTTPClient: one that reaches no address of the
// server's own network.
func (x *Explorer) fetch(w http.ResponseWriter, r *http.Request, dir, raw, name string) {
	if x.Opts.HTTPClient == nil || x.Opts.HTTPClient() == nil {
		Fail(w, http.StatusNotFound, errors.New("not available"))
		return
	}
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		Fail(w, http.StatusBadRequest, errors.New("an http(s) URL"))
		return
	}
	if err := x.allow(r, dir, true); err != nil {
		Fail(w, http.StatusForbidden, err)
		return
	}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, u.String(), nil)
	if err != nil {
		Fail(w, http.StatusBadRequest, err)
		return
	}
	res, err := x.Opts.HTTPClient().Do(req)
	if err != nil {
		Fail(w, http.StatusBadGateway, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		Fail(w, http.StatusBadGateway, fmt.Errorf("%s: %s", u.Host, res.Status))
		return
	}
	if name == "" {
		if _, ps, err := mime.ParseMediaType(res.Header.Get("Content-Disposition")); err == nil {
			name = ps["filename"]
		}
	}
	if name == "" {
		name = path.Base(u.Path)
	}
	name = path.Base("/" + strings.ReplaceAll(name, `\`, "/"))
	if name == "/" || name == "." || name == "" {
		name = "download"
	}
	p := path.Join(Rel(dir), name)
	if err := x.allow(r, p, true); err != nil {
		Fail(w, http.StatusForbidden, err)
		return
	}
	if err := x.create(r.Context(), p, 0o644, res.Body); err != nil {
		Fail(w, http.StatusConflict, err)
		return
	}
	x.changed(r, "fetch", p)
	writeJSON(w, map[string]any{"ok": true, "path": p})
}
