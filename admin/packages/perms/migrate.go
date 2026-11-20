package perms

import (
	"context"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/lifecycle"
	"github.com/go-rvq/rvq/x/migrate"
	"gorm.io/gorm"
)

// TrashKey is the node of the migration of the trash of table (the foreign
// key of who deleted: MigrateDeleted) in the graphs of the commands.
func TrashKey(table string) lifecycle.Key {
	return lifecycle.Key("rvq/perms/trash/" + table)
}

// offerTrashMigration offers p the migration of the trash of table — of the
// application's model obj —, which Trash still runs too, for the
// applications that do not run the graph yet.
func offerTrashMigration(p *presets.Builder, table string, obj any) {
	key := TrashKey(table)
	if p == nil || p.Migrations().Has(key) {
		return
	}
	p.Migrations().Chain(lifecycle.Node{
		Key: key,
		Start: migrate.Func(func(_ context.Context, db *gorm.DB) error {
			return MigrateDeleted(db, table, obj)
		}),
	})
}
