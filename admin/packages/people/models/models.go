// Package models holds the persistence models of the generic people package: a
// Person (individual/company identity) scoped to a go-rvq orgs Organization, so
// each organization keeps its own people registry (see the admin package).
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base is the shared base of every people model: UUID primary key, timestamps,
// logical deletion and the owning organization.
//
// OrganizacaoID scopes every record to an orgs Organization: each organization
// has its own people registry. The field name matches the orgs scoping contract
// (orgs stamps/filters "OrganizacaoID"/"organizacao_id"); it is stamped from the
// org in the request path (/orgs/{id}/people/...) on save and used to scope
// reads (see orgs MountUnderOrg / ScopeToOrg).
type Base struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrganizacaoID uuid.UUID `gorm:"type:uuid;index"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate assigns a random UUID when none was provided, so it works on any
// database (no server-side uuid function required).
func (b *Base) BeforeCreate(*gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// AutoMigrate creates/updates the people tables.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Person{})
}
