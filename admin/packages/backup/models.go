package backup

import (
	"time"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/x/packages/backup"
	"github.com/google/uuid"
)

// The states of a BackupState.
const (
	StatusQueued   = "queued"
	StatusBuilding = "building"
	StatusReady    = "ready"
	StatusFailed   = "failed"
	StatusCanceled = "canceled"
)

// The origins of a BackupState.
const (
	OriginManual    = "manual"
	OriginScheduled = "scheduled"
)

// BackupState is a state of the backup (x/packages/backup): what was asked —
// its description, its sources, kept or not —, who asked, and what it is:
// its status, as it is made and once made.
type BackupState struct {
	// ID is the state's: its folder (x/packages/backup).
	ID          string `gorm:"primaryKey;size:64"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string `gorm:"type:text"`
	// Sources are those asked; empty: all.
	Sources []string `gorm:"type:text;serializer:json"`
	// Keep: the retention never removes it.
	Keep   bool
	Status string `gorm:"size:16;index;not null;default:'queued'"`
	Origin string `gorm:"size:16;not null;default:'manual'"`
	// CreatedByID and CreatedByName are who asked it (none: the schedule).
	CreatedByID   model.UUID `gorm:"type:uuid"`
	CreatedByName string
	StartedAt     *time.Time
	FinishedAt    *time.Time
	Error         string `gorm:"type:text"`
	// Progress, while it is made: the source, its files and bytes so far.
	ProgressSource string
	ProgressFiles  int
	ProgressBytes  int64
	// What it is, once made.
	Files, Dirs, Links int
	Bytes              int64
	Partial            bool
	ManifestSize       int64
	ManifestSum        string
	Results            []backup.SourceResult `gorm:"type:text;serializer:json"`
	Meta               map[string]string     `gorm:"type:text;serializer:json"`
}

// TableName is backup_states.
func (BackupState) TableName() string { return "backup_states" }

// InProgress says it is queued or being made.
func (s *BackupState) InProgress() bool {
	return s.Status == StatusQueued || s.Status == StatusBuilding
}

// BackupSession is a session of a client downloading a state: opened by an
// access key (!session), its token valid until ExpiresAt; what it was served.
type BackupSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	StateID   string `gorm:"size:64;index;not null"`
	ExpiresAt time.Time
	// UserID and UserName are the key's user's.
	UserID    model.UUID `gorm:"type:uuid"`
	UserName  string
	IP        string
	UserAgent string `gorm:"type:text"`
	// Files and Bytes served, counted in memory and written now and then.
	Files int
	Bytes int64
}

// TableName is backup_sessions.
func (BackupSession) TableName() string { return "backup_sessions" }
