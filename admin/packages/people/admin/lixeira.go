package admin

import (
	"github.com/go-rvq/rvq/admin/packages/perms"
	"github.com/go-rvq/rvq/admin/presets"
	"gorm.io/gorm"
)

// configureTrash applies the shared, resource-agnostic rvq behaviors to a people
// resource: localized save errors, the permission manager and the trash view
// (gated by the 'trash' listing permission). Deleted records stay out of the
// default queries (gorm soft delete).
func configureTrash(mb *presets.ModelBuilder, db *gorm.DB, sample interface{ TableName() string }) {
	wrapSaveErrors(mb)
	perms.InstallManager(mb, db)
	perms.SetupTrash(mb, db, perms.TrashOptions{TableName: sample.TableName()})
}
