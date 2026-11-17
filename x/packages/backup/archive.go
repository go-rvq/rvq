package backup

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

// SkippedFile is the file an archive ends with when some of its files were
// left out (WriteArchive): their paths, and why.
const SkippedFile = "BACKUP-SKIPPED.txt"

// ArchiveOptions are what WriteArchive puts in the archive.
type ArchiveOptions struct {
	// Paths are the entries (and, of a folder, everything under it): paths
	// of the manifest ("app/data", no leading "/"); none, the whole state.
	Paths []string
	// Base is taken out of the names in the archive: the folder the paths
	// are in ("app": "app/data/x" is "data/x"); "", the names as they are.
	Base string
}

// WriteArchive writes the state id as a tar.gz to w, as it is read — in
// the order of the manifest, the modes, the times and the links kept.
// Each file is opened (OpenFile) before its header: one changed since the
// state (ErrChanged) or gone is left out, and the archive ends with
// SkippedFile, the list of them. skipped are their paths.
func (b *Builder) WriteArchive(w io.Writer, id string, opts ArchiveOptions) (skipped []string, err error) {
	m, err := b.Manifest(id)
	if err != nil {
		return nil, err
	}
	base := strings.Trim(opts.Base, "/")
	var roots []string
	for _, p := range opts.Paths {
		p = strings.Trim(path.Clean("/"+p), "/")
		if p == "" {
			roots = nil
			break
		}
		roots = append(roots, p)
	}
	in := func(p string) bool {
		if len(roots) == 0 {
			return true
		}
		for _, r := range roots {
			if p == r || strings.HasPrefix(p, r+"/") {
				return true
			}
		}
		return false
	}
	name := func(p string) string {
		if base == "" {
			return p
		}
		return strings.TrimPrefix(p, base+"/")
	}

	gz := gzip.NewWriter(w)
	tw := tar.NewWriter(gz)
	var why []string
	for _, e := range m.Entries {
		if !in(e.Path) || (base != "" && !strings.HasPrefix(e.Path, base+"/")) {
			continue
		}
		hd := &tar.Header{Name: name(e.Path), Mode: int64(unixMode(e.Mode)), ModTime: e.MTime, Format: tar.FormatPAX}
		switch e.Type {
		case TypeDir:
			hd.Typeflag, hd.Name = tar.TypeDir, hd.Name+"/"
			err = tw.WriteHeader(hd)
		case TypeLink:
			hd.Typeflag, hd.Linkname = tar.TypeSymlink, e.Target
			err = tw.WriteHeader(hd)
		case TypeFile:
			f, oerr := b.OpenFile(id, e)
			if oerr != nil {
				skipped = append(skipped, e.Path)
				why = append(why, e.Path+"\t"+reason(oerr))
				continue
			}
			hd.Typeflag, hd.Size = tar.TypeReg, e.Size
			if err = tw.WriteHeader(hd); err == nil {
				_, err = io.CopyN(tw, f, e.Size)
			}
			f.Close()
		}
		if err != nil {
			return skipped, err
		}
	}
	if len(why) > 0 {
		body := "These files changed (or were removed) after the state was made and are not in the archive:\n\n" +
			strings.Join(why, "\n") + "\n"
		if err = tw.WriteHeader(&tar.Header{Name: SkippedFile, Mode: 0o644, Size: int64(len(body)), Typeflag: tar.TypeReg}); err == nil {
			_, err = io.WriteString(tw, body)
		}
		if err != nil {
			return skipped, err
		}
	}
	if err = tw.Close(); err != nil {
		return skipped, err
	}
	return skipped, gz.Close()
}

// reason is why a file was left out of an archive.
func reason(err error) string {
	switch {
	case errors.Is(err, ErrChanged):
		return "changed"
	case errors.Is(err, fs.ErrNotExist):
		return "removed"
	}
	return fmt.Sprint(err)
}
