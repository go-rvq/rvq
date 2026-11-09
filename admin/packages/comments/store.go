package comments

import (
	"bytes"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// StoredFile is an attachment kept on disk (SaveFiles): its id, its name,
// its type, its size and its path under the directory of the files.
type StoredFile struct {
	ID          uuid.UUID
	Name        string
	ContentType string
	Size        int64
	// Path is under the directory: "<owner>/<id><ext>"
	Path string
}

// SaveFiles keeps the files of owner (a comment, a message) in dir: each in
// dir/<owner>/<id><ext> (its extension kept, its name the record's — what is
// written is never named by who sent it). A failure removes what was
// written.
func SaveFiles(dir string, owner uuid.UUID, files ...Attachment) (out []StoredFile, err error) {
	if len(files) == 0 {
		return nil, nil
	}
	sub := filepath.Join(dir, owner.String())
	if err := os.MkdirAll(sub, 0o750); err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			for _, f := range out {
				_ = os.Remove(filepath.Join(dir, f.Path))
			}
			out = nil
		}
	}()
	for _, a := range files {
		if len(a.Data) > FileMaxSize {
			return out, fmt.Errorf("%w: %s", ErrFileTooLarge, a.Name)
		}
		ct := a.ContentType
		if ct == "" {
			ct = mime.TypeByExtension(filepath.Ext(a.Name))
		}
		id := uuid.New()
		ext := strings.ToLower(filepath.Ext(filepath.Base(a.Name)))
		if len(ext) > 16 || strings.ContainsAny(ext, `/\`) {
			ext = ""
		}
		rel := filepath.ToSlash(filepath.Join(owner.String(), id.String()+ext))
		w, err := os.OpenFile(filepath.Join(dir, rel), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o640)
		if err != nil {
			return out, err
		}
		if _, err := w.Write(a.Data); err != nil {
			w.Close()
			return out, err
		}
		if err := w.Close(); err != nil {
			return out, err
		}
		out = append(out, StoredFile{ID: id, Name: filepath.Base(a.Name), ContentType: ct, Size: int64(len(a.Data)), Path: rel})
	}
	return out, nil
}

// FilePath is the file of path under dir — nothing out of it.
func FilePath(dir, path string) (string, bool) {
	p := filepath.Join(dir, filepath.FromSlash(path))
	rel, err := filepath.Rel(dir, p)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false
	}
	return p, true
}

// inlineTypes are the types an attachment is shown in the browser as: what
// runs nothing.
var inlineTypes = map[string]bool{
	"image/png": true, "image/jpeg": true, "image/gif": true, "image/webp": true, "application/pdf": true,
}

// ServeFile sends the attachment f — its content: of path under dir, or data
// (a file of before kept in the database) — with what keeps a user's file
// from running: shown in the browser only an image or a PDF (an SVG is
// downloaded), sniffing off, sandboxed.
func ServeFile(w http.ResponseWriter, r *http.Request, dir string, f StoredFile, data []byte) {
	ct := strings.ToLower(strings.TrimSpace(strings.Split(f.ContentType, ";")[0]))
	disposition := "attachment"
	if inlineTypes[ct] {
		disposition = "inline"
	} else {
		ct = "application/octet-stream"
	}
	h := w.Header()
	h.Set("Content-Type", ct)
	h.Set("Content-Disposition", mime.FormatMediaType(disposition, map[string]string{"filename": f.Name}))
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Content-Security-Policy", "default-src 'none'; img-src 'self'; style-src 'unsafe-inline'; sandbox")
	h.Set("Cache-Control", "private, max-age=3600")
	if f.Path == "" {
		http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(data))
		return
	}
	p, ok := FilePath(dir, f.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}
	file, err := os.Open(p)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()
	st, err := file.Stat()
	if err != nil || st.IsDir() {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, "", st.ModTime(), file)
}
