// Package uploadparts receives a file of any size in parts — each a short
// request no proxy's limit of size or time cuts —, as the explorer of the
// files of senir does: the parts written at their place (offset) in a file of
// their own, its size checked and its sha256 calculated, once, when they all
// came (Done); the file then taken by a token (Open). Before, a hash may be
// asked (Exists): a content the application has needs no sending.
//
// Its events are of a model (Install): each asks what Allow says (the
// model's @create, by default).
//
//	uploadparts_Exists  hash                       → the application's data of it, or null
//	uploadparts_Part    id, offset, total; "part"  → the part written
//	uploadparts_Done    id, total, name[, hash]    → {id, name, size, hash}
//	uploadparts_Abort   id                         → its parts removed
//
// A sending left (its parts not touched for MaxAge) is removed.
package uploadparts

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
)

const (
	EventExists = "uploadparts_Exists"
	EventPart   = "uploadparts_Part"
	EventDone   = "uploadparts_Done"
	EventAbort  = "uploadparts_Abort"

	// PartSize is the size of a part the client sends (8MB).
	PartSize int64 = 8 << 20
	// MaxAge is how long a sending not touched stays.
	MaxAge = 24 * time.Hour

	partsSuffix = ".parts"
	doneSuffix  = ".done"
)

var (
	// ErrID is an id of a sending that is none (a name of a file, not a path).
	ErrID = errors.New("uploadparts: the id of the sending: 8 to 64 letters, digits or -")
	// ErrSize is a sending done with a size other than its total.
	ErrSize = errors.New("uploadparts: the parts received do not make the file's size")
	// ErrPart is a part past the file's size.
	ErrPart = errors.New("uploadparts: the part goes past the file's size")
	// ErrNotDone is a token of a sending not done (or gone).
	ErrNotDone = errors.New("uploadparts: no file of this sending")

	idRe = regexp.MustCompile(`^[A-Za-z0-9-]{8,64}$`)
)

// Token is a file received, by its sending: its id, its name, its size and
// the sha256 of its content (the server's).
type Token struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Hash string `json:"hash"`
	// MTime is when the file was last modified, as the client said
	// (milliseconds since 1970: the browser's lastModified); 0, unknown.
	MTime int64 `json:"mtime,omitempty"`
}

// ExistsFunc is the application's data of a content of hash it has (the
// client shows it, links it), nil of none.
type ExistsFunc func(ctx *web.EventContext, hash string) (any, error)

// Builder receives the files sent in parts to dir.
type Builder struct {
	dir    string
	exists ExistsFunc
	allow  func(ctx *web.EventContext, mb *presets.ModelBuilder) bool
	now    func() time.Time
}

// New receives the files in dir (created when needed).
func New(dir string) *Builder {
	return &Builder{dir: dir, now: time.Now, allow: func(ctx *web.EventContext, mb *presets.ModelBuilder) bool {
		return mb.Info().Permissioner().ReqCreator(ctx.R).Allowed()
	}}
}

// Exists is the application's data of a content by its hash.
func (b *Builder) Exists(f ExistsFunc) *Builder {
	b.exists = f
	return b
}

// Allow says who may send (the model's @create, by default).
func (b *Builder) Allow(f func(ctx *web.EventContext, mb *presets.ModelBuilder) bool) *Builder {
	b.allow = f
	return b
}

// Install registers the events of the sending in mb.
func (b *Builder) Install(mb *presets.ModelBuilder) {
	guard := func(f func(ctx *web.EventContext) (any, error)) web.EventFunc {
		return func(ctx *web.EventContext) (r web.EventResponse, err error) {
			if !b.allow(ctx, mb) {
				return r, perm.PermissionDenied
			}
			r.Data, err = f(ctx)
			return
		}
	}
	mb.RegisterEventFunc(EventExists, guard(func(ctx *web.EventContext) (any, error) {
		hash := strings.ToLower(strings.TrimSpace(ctx.R.FormValue("hash")))
		if b.exists == nil || !validHash(hash) {
			return nil, nil
		}
		return b.exists(ctx, hash)
	}))
	mb.RegisterEventFunc(EventPart, guard(b.part))
	mb.RegisterEventFunc(EventDone, guard(func(ctx *web.EventContext) (any, error) {
		t, err := b.Done(ctx.R.FormValue("id"), int64Of(ctx.R.FormValue("total")), ctx.R.FormValue("name"), ctx.R.FormValue("hash"))
		if t != nil {
			if m := int64Of(ctx.R.FormValue("mtime")); m > 0 {
				t.MTime = m
			}
		}
		return t, err
	}))
	mb.RegisterEventFunc(EventAbort, guard(func(ctx *web.EventContext) (any, error) {
		b.Remove(ctx.R.FormValue("id"))
		return nil, nil
	}))
}

func validHash(h string) bool {
	if len(h) != 64 {
		return false
	}
	_, err := hex.DecodeString(h)
	return err == nil
}

func int64Of(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1
	}
	return n
}

// paths are the file of the parts of the sending id and the file done.
func (b *Builder) paths(id string) (parts, done string, err error) {
	if !idRe.MatchString(id) {
		return "", "", ErrID
	}
	return filepath.Join(b.dir, id+partsSuffix), filepath.Join(b.dir, id+doneSuffix), nil
}

// part writes the part posted at its offset: never past the file's size.
func (b *Builder) part(ctx *web.EventContext) (any, error) {
	if err := ctx.R.ParseMultipartForm(PartSize + 1<<20); err != nil {
		return nil, err
	}
	f, _, err := ctx.R.FormFile("part")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return nil, b.WritePart(ctx.R.FormValue("id"), int64Of(ctx.R.FormValue("offset")), int64Of(ctx.R.FormValue("total")), f)
}

// WritePart writes r at offset of the file of the parts of the sending id, of
// the size total.
func (b *Builder) WritePart(id string, offset, total int64, r io.Reader) error {
	parts, _, err := b.paths(id)
	if err != nil {
		return err
	}
	if offset < 0 || total < 0 || offset > total {
		return fmt.Errorf("uploadparts: offset and total: 0 ≤ offset ≤ total")
	}
	if err = os.MkdirAll(b.dir, 0o755); err != nil {
		return err
	}
	b.clean()
	f, err := os.OpenFile(parts, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	// one byte more than what is left: too much
	n, err := io.Copy(io.NewOffsetWriter(f, offset), io.LimitReader(r, total-offset+1))
	if err != nil {
		return err
	}
	if offset+n > total {
		return ErrPart
	}
	return nil
}

// Done ends the sending id: the parts, of the size total, its file — its
// sha256 calculated, once —: the token of it. hash, the client's, is not
// trusted: the token's is the content's.
func (b *Builder) Done(id string, total int64, name, hash string) (*Token, error) {
	parts, done, err := b.paths(id)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		// an empty file: no part sent
		if err = os.MkdirAll(b.dir, 0o755); err != nil {
			return nil, err
		}
		if f, err := os.OpenFile(parts, os.O_CREATE|os.O_WRONLY, 0o600); err == nil {
			f.Close()
		}
	}
	f, err := os.Open(parts)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	n, err := io.Copy(h, f)
	f.Close()
	if err != nil {
		return nil, err
	}
	if n != total {
		return nil, fmt.Errorf("%w: %d of %d", ErrSize, n, total)
	}
	if err = os.Rename(parts, done); err != nil {
		return nil, err
	}
	return &Token{ID: id, Name: filepath.Base(strings.ReplaceAll(name, "\\", "/")), Size: n, Hash: hex.EncodeToString(h.Sum(nil))}, nil
}

// Open is the file received of the sending id (done).
func (b *Builder) Open(id string) (*os.File, error) {
	_, done, err := b.paths(id)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(done)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotDone
	}
	return f, err
}

// Remove removes the sending id: its parts, its file.
func (b *Builder) Remove(id string) {
	if parts, done, err := b.paths(id); err == nil {
		os.Remove(parts)
		os.Remove(done)
	}
}

// clean removes the sendings not touched for MaxAge.
func (b *Builder) clean() {
	entries, err := os.ReadDir(b.dir)
	if err != nil {
		return
	}
	limit := b.now().Add(-MaxAge)
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, partsSuffix) && !strings.HasSuffix(name, doneSuffix) {
			continue
		}
		if info, err := e.Info(); err == nil && info.ModTime().Before(limit) {
			os.Remove(filepath.Join(b.dir, name))
		}
	}
}
