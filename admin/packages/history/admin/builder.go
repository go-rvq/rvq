// Package admin mounts the history plugin (git-like revisions for models) into
// an rvq presets admin. See PLAN.md.
package admin

import (
	"github.com/go-rvq/rvq/admin/presets"
	"gorm.io/gorm"
)

// RevisionsSuffix is appended to a model's table name to name its per-model
// revisions table (e.g. "posts" -> "posts_revisions").
const RevisionsSuffix = "_revisions"

// Builder is the history plugin (a presets.Plugin). It owns the shared access
// Recorder that counts per-revision public accesses.
type Builder struct {
	db       *gorm.DB
	recorder *Recorder
}

// NewBuilder returns a history plugin builder.
func NewBuilder(db *gorm.DB) *Builder {
	return &Builder{db: db, recorder: NewRecorder(db)}
}

// Configure mounts the history plugin on b (starts the access recorder) and
// returns it. Activate per model with New(db).Model(mb).Build().
func Configure(b *presets.Builder, db *gorm.DB) *Builder {
	ConfigureMessages(b.I18n())
	hb := NewBuilder(db)
	b.Use(hb)
	return hb
}

// Install satisfies presets.Plugin: it publishes the shared recorder and starts
// its background flusher.
func (b *Builder) Install(pb *presets.Builder) error {
	defaultRecorder = b.recorder
	b.recorder.Start()
	return nil
}

// Recorder is the shared access recorder — the site calls Recorder().Hit(table,
// recordKey) on a public render.
func (b *Builder) Recorder() *Recorder { return b.recorder }

// DefaultRecorder returns the process-wide access recorder set by Configure (nil
// until then). The public site uses it to count per-revision accesses:
//
//	history.DefaultRecorder().Hit("posts_revisions", recordKey)
func DefaultRecorder() *Recorder { return defaultRecorder }

// RevisionTableFor returns the revisions table name for a model value
// ("<model table>_revisions"), so callers outside the plugin (e.g. the site)
// can name it without hardcoding.
func RevisionTableFor(db *gorm.DB, m any) string { return revisionTable(db, m) }
