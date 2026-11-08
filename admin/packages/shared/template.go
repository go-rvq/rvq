package shared

import (
	"time"

	"github.com/go-rvq/rvq/admin/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ShareTemplate is a role of the sharing: a named set of verbs (SeedRoles,
// RolesOf), of a resource (empty: any).
//
// Scope, when set, restricts a role to one context (an organization, say); a
// nil Scope is a global role.
type ShareTemplate struct {
	ID      uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Scope   model.UUID     `gorm:"type:uuid;index"`
	Name    string         `gorm:"size:120;not null"`
	Actions pq.StringArray `gorm:"type:text[]"`
	// Resource is the resource the role is offered for (the model of the
	// record shared); empty: any.
	Resource string `gorm:"size:120;index"`
	// None: the role of no access (RoleNone) — of a record under one shared,
	// it takes it out.
	None      bool `gorm:"not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ShareTemplate) TableName() string { return "shared_templates" }

func (t *ShareTemplate) BeforeCreate(*gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// scopeOf is the scope of a template of id.
func scopeOf(id uuid.UUID) model.UUID { return model.UUID(id) }
