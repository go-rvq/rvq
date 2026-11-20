package accesskeys

import (
	"context"

	"github.com/go-rvq/rvq/x/lifecycle"
	"github.com/go-rvq/rvq/x/migrate"
	"gorm.io/gorm"

	"github.com/go-rvq/rvq/admin/role"
)

// Key is the node of the package in the graphs of the commands
// (x/lifecycle).
const Key lifecycle.Key = "rvq/accesskeys"

// its migrations, in the graph of every package's (x/migrate): `app migrate`
// and the serve run them; Builder.AutoMigrate still runs them too, for the applications
// that do not run the graph yet
func init() {
	migrate.Graph.Add(lifecycle.Node{
		Key:   Key,
		After: []lifecycle.Key{role.Key},
		Start: migrate.Func(func(_ context.Context, db *gorm.DB) error { return Migrate(db) }),
	})
}

// Migrate makes the tables of the keys; their policies of before made lines
// (MigrateLines), their groups undone (MigrateUngroup).
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&AccessKey{}, &AccessKeyUse{}); err != nil {
		return err
	}
	if err := MigrateLines(db); err != nil {
		return err
	}
	return MigrateUngroup(db)
}
