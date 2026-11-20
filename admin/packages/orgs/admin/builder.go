package admin

import (
	"github.com/go-rvq/rvq/admin/packages/orgs/messages"
	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/packages/shared"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/migrate"
	"gorm.io/gorm"
)

// Builder mounts the organizations package. It implements presets.Plugin, so it
// is applied with pb.Use(builder). It is idempotent: installing when the
// Organizacao resource is already registered is a no-op, so other packages
// (people, finance) can safely ensure it. After Install it also exposes the
// org model and the MountUnderOrg helper as methods.
type Builder struct {
	db      *gorm.DB
	pb      *presets.Builder
	sharing *shared.Service
	// roles are the roles of the organizations (OrgRoles when none), the
	// ones of before renamed (Roles)
	roles       []shared.RoleSeed
	roleRenames map[string]string
}

var _ presets.Plugin = (*Builder)(nil)

// New returns an orgs Builder.
func New(db *gorm.DB) *Builder {
	return &Builder{db: db}
}

// Install implements presets.Plugin: it migrates the org tables, registers the
// i18n messages and the Organizacao resource, the selector, the nested menu
// and — shared by a service (Sharing) — the sharing; without one, an
// organization is its owner's alone.
func (b *Builder) Install(pb *presets.Builder) error {
	b.pb = pb
	if OrganizacaoModel(pb) != nil {
		return nil
	}
	if err := migrate.Auto(b.db, models.AutoMigrate); err != nil {
		return err
	}
	messages.Register(pb.I18n())

	configureOrganizacao(pb, b.db, b.sharing == nil)
	// the organization in the side menu: a root item, its chip opens the
	// choice; the one chosen, its resources under it
	configureMenuRoot(pb, b.db)
	configureNestedMenu(pb, b.db)
	if b.sharing != nil {
		// shared by the service: the actions of sharing, the page of the
		// confirmations, the invites
		return installSharing(pb, b.sharing, b.roles, b.roleRenames)
	}
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
