package admin

import (
	vmessages "github.com/go-rvq/rvq/admin/packages/validators/messages"
	vmodels "github.com/go-rvq/rvq/admin/packages/validators/models"
	"github.com/go-rvq/rvq/admin/presets"
	"gorm.io/gorm"
)

// Builder mounts the validators package and registers validators. It implements
// presets.Plugin, so it is applied with pb.Use(builder). Register additional
// validators with Register before installing; other packages (e.g. people) seed
// their own validators through this builder or directly via the models store.
type Builder struct {
	db         *gorm.DB
	validators []*vmodels.Validator
	exposeCRUD bool
}

var _ presets.Plugin = (*Builder)(nil)

// New returns a validators Builder. By default it exposes the Validator CRUD.
func New(db *gorm.DB) *Builder {
	return &Builder{db: db, exposeCRUD: true}
}

// Register adds validators to be seeded on Install (upsert by name, preserving
// user overrides).
func (b *Builder) Register(v ...*vmodels.Validator) *Builder {
	b.validators = append(b.validators, v...)
	return b
}

// ExposeCRUD toggles whether the Validator admin resource is registered
// (default true). Set false to only seed validators without the CRUD.
func (b *Builder) ExposeCRUD(v bool) *Builder {
	b.exposeCRUD = v
	return b
}

// Install implements presets.Plugin: it migrates the table, registers the i18n
// messages, seeds the registered validators and (optionally) registers the
// Validator CRUD resource.
func (b *Builder) Install(pb *presets.Builder) error {
	if err := vmodels.AutoMigrate(b.db); err != nil {
		return err
	}
	vmessages.Register(pb.I18n())
	for _, v := range b.validators {
		if err := vmodels.Seed(b.db, v); err != nil {
			return err
		}
	}
	if b.exposeCRUD && pb.GetModelByID(ValidatorModelID) == nil {
		configureValidator(pb, b.db)
	}
	return nil
}

// Configure is a convenience that installs a default validators Builder (CRUD
// exposed) on pb, returning the Validator model builder.
func Configure(pb *presets.Builder, db *gorm.DB) *presets.ModelBuilder {
	if err := New(db).Install(pb); err != nil {
		panic(err)
	}
	return pb.GetModelByID(ValidatorModelID)
}
