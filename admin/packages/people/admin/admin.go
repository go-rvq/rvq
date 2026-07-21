// Package admin mounts the people package into an rvq presets admin: the Person
// registry (identity), scoped to an organization so each organization keeps its
// own people. Other packages (e.g. finance) depend on it to attach per-person
// data as nested sub-resources.
package admin

import (
	orgsadmin "github.com/go-rvq/rvq/admin/packages/orgs/admin"
	"github.com/go-rvq/rvq/admin/packages/people/messages"
	"github.com/go-rvq/rvq/admin/packages/people/models"
	"github.com/go-rvq/rvq/admin/presets"
	"gorm.io/gorm"
)

// PersonModelID is the presets model id of the Person resource.
const PersonModelID = "persons"

// PersonModel returns the registered Person model builder, or nil when the
// people package has not been configured on b. Packages that attach per-person
// data (e.g. finance) use it to add nested child resources.
func PersonModel(b *presets.Builder) *presets.ModelBuilder {
	return b.GetModelByID(PersonModelID)
}

// Configure mounts the people package on b. It ensures the orgs package is
// configured (idempotent) so the Person resource can be nested under the
// organization.
func Configure(b *presets.Builder, db *gorm.DB) {
	if err := models.AutoMigrate(db); err != nil {
		panic(err)
	}
	if orgsadmin.OrganizacaoModel(b) == nil {
		orgsadmin.Configure(b, db)
	}
	messages.Register(b.I18n())
	registerEnumSelects(b)
	configurePerson(b, db)
}
