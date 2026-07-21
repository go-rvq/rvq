// Package models holds the persistence of the generic validators package: a
// Validator whose validation algorithm is written in GAD
// (github.com/gad-lang/gad) and whose error messages are translated per user
// language. Records are seeded by the application (not created by end users);
// only the Value, Doc and Messages overrides are user-editable (see the admin
// package). When an override is empty the corresponding Initial* value is used.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AutoMigrate creates/updates the validators table.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Validator{})
}

// Base is the shared primary key and timestamps of the validators models.
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate assigns a random UUID when none was provided.
func (b *Base) BeforeCreate(*gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
