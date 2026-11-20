package role

import (
	"context"

	"github.com/go-rvq/rvq/x/lifecycle"
	"github.com/go-rvq/rvq/x/migrate"
	"gorm.io/gorm"
)

// Key is the node of the package in the graphs of the commands
// (x/lifecycle).
const Key lifecycle.Key = "rvq/role"

// its migrations, in the graph of every package's (x/migrate): `app migrate`
// and the serve run them; Install still runs them too, for the applications
// that do not run the graph yet
func init() {
	migrate.Graph.Add(lifecycle.Node{
		Key:   Key,
		Start: migrate.Func(func(_ context.Context, db *gorm.DB) error { return Migrate(db) }),
	})
}

// Migrate makes the tables of the roles — their policies of before made
// lines (MigrateLines), cleaned (MigrateUngroup) — and of the sets of
// permissions.
func Migrate(db *gorm.DB) error {
	if err := MigrateLines(db); err != nil {
		return err
	}
	if err := MigrateUngroup(db); err != nil {
		return err
	}
	return db.AutoMigrate(&PermSet{})
}
