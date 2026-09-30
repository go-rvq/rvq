// Package admin mounts the history plugin (git-like revisions for models) into
// an rvq presets admin. See PLAN.md.
package admin

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/activity"
	"github.com/go-rvq/rvq/admin/geomap"
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"gorm.io/gorm"
)

// revisionsByTable maps each versioned model's revisions table to its history,
// so the activity log view can render a revision's diff (RevisionDiffFunc) by
// table name without importing this package's per-model state.
var revisionsByTable = map[string]*ModelHistory{}

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
	// the map of a revision's origin (ActionOrigin)
	geomap.Install(b)
	// Let the activity log render a change from its revision instead of a
	// duplicated diff (no-op when activity is not mounted).
	activity.RevisionDiffFunc = revisionDiff
	hb := NewBuilder(db)
	b.Use(hb)
	return hb
}

// revisionDiff renders a revision's change (its diff against its parent) for the
// activity log view. It is registered as activity.RevisionDiffFunc.
func revisionDiff(table string, hash []byte, ctx *web.EventContext) h.HTMLComponent {
	mh := revisionsByTable[table]
	if mh == nil {
		return nil
	}
	var rev histmodels.Revision
	if err := mh.db.Table(table).Where("hash = ?", hash).First(&rev).Error; err != nil {
		return nil
	}
	if len(rev.Parent) == 0 {
		return nil
	}
	comp, err := mh.compare(rev.RecordKey, rev.Parent, rev.Hash, ctx)
	if err != nil {
		return nil
	}
	return comp
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
