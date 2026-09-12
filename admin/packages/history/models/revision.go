// Package models holds the history plugin's persistent types.
package models

import (
	"encoding/json"
	"time"

	"github.com/go-rvq/rvq/thirdpart/gorm/datatypes"
	"github.com/google/uuid"
)

// Fields is the per-revision snapshot: each versioned field's value as raw JSON,
// stored as a JSONB column (datatypes.NullJSONType keeps the map and round-trips
// it), so a diff can tell a plain string from a structured/foreign-key value.
type Fields = datatypes.NullJSONType[map[string]json.RawMessage]

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
	Hash Hash `gorm:"primaryKey;type:bytea"`

	// RecordKey identifies the record this revision belongs to: its primary key
	// serialized (e.g. "12" or "12|pt-BR" for id + locale). Indexed for the
	// per-record history walk and to scope the nested revisions listing.
	RecordKey string `gorm:"index;type:varchar(255)"`

	// Parent is the hash of the previous revision of the same record — the chain
	// (a git-like DAG). NULL on the first revision.
	Parent Hash `gorm:"type:bytea"`

	// Fields is the snapshot {field: raw JSON value} of just the versioned fields.
	Fields Fields `gorm:"type:jsonb"`

	// ChangedFields lists the versioned fields whose value differs from the parent
	// revision (all versioned fields on the first revision). The snapshot above is
	// still complete; this is what makes the history queryable by field — a
	// field's timeline is exactly the revisions that list it.
	ChangedFields datatypes.NullJSONType[[]string] `gorm:"type:jsonb"`

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

// PrimarySlug is the revision's id in a URL and in the listing's row selection:
// the hash as hex. The primary key is Hash (there is no "ID" field), so without
// this the admin data table would read a non-existent "ID" and give every row an
// empty id — breaking selection and the compare bulk action. It round-trips
// through Hash.Parse in ParseRecordID.
func (r Revision) PrimarySlug() string {
	return r.Hash.String()
}
