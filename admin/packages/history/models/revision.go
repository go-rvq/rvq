// Package models holds the history plugin's persistent types.
package models

import (
	"time"

	"github.com/google/uuid"
)

// Revision is one point in a model record's history — a git-like commit of the
// versioned fields. One table per versioned model, named
// "<model table>_revisions" (e.g. "posts_revisions"), so there is no ModelName
// column: the table already means the model. A single Go type is mapped to each
// such table with db.Table(name).
type Revision struct {
	// Hash is the primary key: sha256 (32 bytes, bytea) of the record key plus
	// the canonical serialization of the versioned fields — like a git blob id,
	// scoped to the record. Folding the record key in keeps it unique per record
	// (so it can be the sole PK, giving a clean /…/revisions/<hash> route) while
	// still collapsing an unchanged save to the same hash (no duplicate revision).
	Hash []byte `gorm:"primaryKey;type:bytea"`

	// RecordKey identifies the record this revision belongs to: its primary key
	// serialized (e.g. "12" or "12|pt-BR" for id + locale). Indexed for the
	// per-record history walk and to scope the nested revisions listing.
	RecordKey string `gorm:"index;type:varchar(255)"`

	// Parent is the hash of the previous revision of the same record — the chain
	// (a git-like DAG). NULL on the first revision.
	Parent []byte `gorm:"type:bytea"`

	// Fields is the snapshot {field: value} of just the versioned fields, JSON.
	Fields []byte `gorm:"type:jsonb"`

	CreatedAt time.Time

	// CreatorID is the author's user id (FK users.id) — never nil: the logged-in
	// user, or the static AnonymousID. Creator is the display name at the time.
	CreatorID uuid.UUID `gorm:"type:uuid;index"`
	Creator   string

	// Tag names the published version this revision became — the git "tag" set
	// on publication; empty on an ordinary (unpublished) revision.
	Tag string
	// Published marks the revision that was published. The record's current
	// published revision is the latest with Published for its RecordKey.
	Published bool `gorm:"index"`

	// AccessCount counts public accesses of this revision (the published one a
	// request resolves to); LastAccess is when it was last hit. Flushed in
	// batches, so they are eventually-consistent.
	AccessCount int64
	LastAccess  time.Time
}
