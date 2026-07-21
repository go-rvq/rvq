// Package models holds the persistence models of the generic organizations
// module: an Organizacao owned by a user and shareable with others through the
// perm API (see the admin package).
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base is the shared primary key and timestamps of every org model. It uses a
// UUID primary key assigned application-side, so it works on any database.
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate assigns a random UUID when none was provided.
func (b *Base) BeforeCreate(*gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// AutoMigrate creates/updates the org tables.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Organizacao{},
	)
}
