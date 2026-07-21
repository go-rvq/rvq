package admin

import (
	"github.com/go-rvq/rvq/admin/packages/orgs/messages"
	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/packages/shared"
	"github.com/go-rvq/rvq/admin/presets"
	"gorm.io/gorm"
)

// Builder mounts the organizations package. It implements presets.Plugin, so it
// is applied with pb.Use(builder). It is idempotent: installing when the
// Organizacao resource is already registered is a no-op, so other packages
// (people, finance) can safely ensure it. After Install it also exposes the
// org model and the MountUnderOrg helper as methods.
type Builder struct {
	db *gorm.DB
	pb *presets.Builder
}

var _ presets.Plugin = (*Builder)(nil)

// New returns an orgs Builder.
func New(db *gorm.DB) *Builder {
	return &Builder{db: db}
}

// Install implements presets.Plugin: it migrates the org tables, seeds the
// sharing permission templates, registers the i18n messages and the
// Organizacao resource, the selector, the nested menu and the share invites.
func (b *Builder) Install(pb *presets.Builder) error {
	b.pb = pb
	if OrganizacaoModel(pb) != nil {
		return nil
	}
	if err := models.AutoMigrate(b.db); err != nil {
		return err
	}
	if err := shared.AutoMigrateTemplates(b.db); err != nil {
		return err
	}
	if err := shared.SeedDefaultTemplates(b.db, nil); err != nil {
		return err
	}
	messages.Register(pb.I18n())

	configureOrganizacao(pb, b.db)
	configureSelector(pb, b.db)
	configureNestedMenu(pb, b.db)
	shared.InstallInvites(pb, b.db)
	return nil
}

// OrganizacaoModel returns the registered Organizacao model builder (nil before
// Install).
func (b *Builder) OrganizacaoModel() *presets.ModelBuilder {
	if b.pb == nil {
		return nil
	}
	return OrganizacaoModel(b.pb)
}

// MountUnderOrg registers child as a resource nested under this builder's
// organization (see the package-level MountUnderOrg). The builder must have been
// installed (pb.Use) so the org model is resolvable.
func (b *Builder) MountUnderOrg(child *presets.ModelBuilder, group string) {
	MountUnderOrg(b.OrganizacaoModel(), child, group)
}
