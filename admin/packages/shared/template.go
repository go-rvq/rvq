package shared

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ShareTemplate is a named, reusable set of verbs a share can grant (SPEC:
// "the share permissions are templatized"). Templates are the choices the
// sharing dialog offers instead of raw verb checkboxes.
//
// Scope, when set, restricts a template to one context (typically an
// organization id) — its "configuration". A nil Scope is a global template
// available everywhere. This keeps the shared package generic while letting
// each organization define its own templates.
type ShareTemplate struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Scope     *uuid.UUID     `gorm:"type:uuid;index"`
	Name      string         `gorm:"size:120;not null"`
	Actions   pq.StringArray `gorm:"type:text[]"`
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

// AutoMigrateTemplates creates the share-template table.
func AutoMigrateTemplates(db *gorm.DB) error { return db.AutoMigrate(&ShareTemplate{}) }

// Templates lists the templates available for scope: the scope's own templates
// plus the global ones (Scope IS NULL), ordered by name.
func Templates(db *gorm.DB, scope *uuid.UUID) ([]ShareTemplate, error) {
	q := db.Model(&ShareTemplate{})
	if scope != nil {
		q = q.Where("scope = ? OR scope IS NULL", *scope)
	} else {
		q = q.Where("scope IS NULL")
	}
	var out []ShareTemplate
	return out, q.Order("name").Find(&out).Error
}

// TemplateActions returns the verbs of a template by id.
func TemplateActions(db *gorm.DB, id uuid.UUID) ([]string, error) {
	var t ShareTemplate
	if err := db.Select("actions").First(&t, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return []string(t.Actions), nil
}

// CreateTemplate adds a template for scope (nil = global).
func CreateTemplate(db *gorm.DB, scope *uuid.UUID, name string, actions []string) (*ShareTemplate, error) {
	t := &ShareTemplate{Scope: scope, Name: name, Actions: actions}
	return t, db.Create(t).Error
}

// SeedDefaultTemplates creates the standard Viewer/Editor/Manager templates for
// scope if it has none yet (idempotent).
func SeedDefaultTemplates(db *gorm.DB, scope *uuid.UUID) error {
	var n int64
	q := db.Model(&ShareTemplate{})
	if scope != nil {
		q = q.Where("scope = ?", *scope)
	} else {
		q = q.Where("scope IS NULL")
	}
	if err := q.Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	defaults := []struct {
		name    string
		actions []string
	}{
		{"Viewer", VerbView},
		{"Editor", append(append([]string{}, VerbView...), VerbEdit...)},
		{"Manager", append(append(append([]string{}, VerbView...), VerbEdit...), VerbDelete...)},
	}
	for _, d := range defaults {
		if _, err := CreateTemplate(db, scope, d.name, d.actions); err != nil {
			return err
		}
	}
	return nil
}
