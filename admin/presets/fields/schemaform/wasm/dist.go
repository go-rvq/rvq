//go:build !(js && wasm)

package wasm

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Dist is a check built (Build) — <Name>.wasm.gz and wasm_exec.js in FS —,
// served at Path + Version() + "/<Name>.wasm" (and "/wasm_exec.js"): an
// address of its own per build, kept by the browser.
type Dist struct {
	FS   fs.FS
	Path string
	Name string

	once    sync.Once
	version string
	gz      []byte
	plain   []byte
	exec    []byte
	built   time.Time
}

// NewDist is the check built name in fsys, served at path.
func NewDist(fsys fs.FS, path, name string) *Dist {
	return &Dist{FS: fsys, Path: path, Name: name}
}

func (d *Dist) load() {
	d.once.Do(func() {
		d.built = time.Now()
		gz, err := fs.ReadFile(d.FS, d.Name+".wasm.gz")
		if err != nil {
			return
		}
		sum := sha256.Sum256(gz)
		d.gz, d.version = gz, hex.EncodeToString(sum[:6])
		d.exec, _ = fs.ReadFile(d.FS, "wasm_exec.js")
	})
}

// Version is the build's: the start of the hash of the wasm; "" when there
// is none.
func (d *Dist) Version() string {
	d.load()
	return d.version
}

// Available says the check was built into this program.
func (d *Dist) Available() bool { return d.Version() != "" }

// URL is where the page loads it from (its directory, of this build); ""
// when it was not built.
func (d *Dist) URL() string {
	if !d.Available() {
		return ""
	}
	return d.Path + d.Version() + "/"
}

// Handler serves the wasm — gzip to who takes it — and wasm_exec.js: public
// files, to any page (a preview in a sandboxed frame, of no origin, too).
func (d *Dist) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ver, name, _ := strings.Cut(strings.TrimPrefix(r.URL.Path, d.Path), "/")
		if !d.Available() || ver != d.Version() {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		switch name {
		case d.Name + ".wasm":
			w.Header().Set("Content-Type", "application/wasm")
			w.Header().Set("Vary", "Accept-Encoding")
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.Header().Set("Content-Encoding", "gzip")
				http.ServeContent(w, r, "", d.built, bytes.NewReader(d.gz))
				return
			}
			http.ServeContent(w, r, "", d.built, bytes.NewReader(d.unzipped()))
		case "wasm_exec.js":
			if d.exec == nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			http.ServeContent(w, r, "", d.built, bytes.NewReader(d.exec))
		default:
			http.NotFound(w, r)
		}
	})
}

var unzipMu sync.Mutex

func (d *Dist) unzipped() []byte {
	unzipMu.Lock()
	defer unzipMu.Unlock()
	if d.plain == nil {
		if zr, err := gzip.NewReader(bytes.NewReader(d.gz)); err == nil {
			d.plain, _ = io.ReadAll(zr)
		}
	}
	return d.plain
}
