package admin

import (
	"github.com/go-rvq/rvq/admin/packages/people/models"
	"github.com/go-rvq/rvq/admin/presets"
	"gorm.io/gorm"
)

// configurePerson registers the people registry resource: a Person (identity)
// owned by the current organization. Its financial receiving data lives in the
// finance package and is mounted as a nested sub-resource by finance.
func configurePerson(b *presets.Builder, db *gorm.DB) *presets.ModelBuilder {
	mb := Model(b, &models.Person{}, presets.ModelWithID("persons")).
		MenuIcon("mdi-account-group")

	mb.Listing("Name", "DocumentType", "Document", "Active").
		SearchColumns("name", "document").
		OrderBy("name")

	ed := mb.Editing("Name", "DocumentType", "Document", "Address", "Notes", "Active")
	ed.Field("Name").Required(true)

	mb.Detailing("Name", "DocumentType", "Document", "Address", "Notes", "Active")

	configureTrash(mb, db, models.Person{})
	return mb
}
