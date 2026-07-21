// Package admin mounts the people package into an rvq presets admin: the Person
// registry (identity), scoped to an organization so each organization keeps its
// own people. Other packages (e.g. finance) depend on it to attach per-person
// data as nested sub-resources.
package admin

import (
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

// Configure mounts the people package on b via the people Builder (a
// presets.Plugin). The optional customizers can register document validators
// for other document types (see Builder.RegisterDocumentValidator).
func Configure(b *presets.Builder, db *gorm.DB, customize ...func(*Builder)) {
	pb := NewBuilder(db)
	for _, c := range customize {
		c(pb)
	}
	b.Use(pb)
}
