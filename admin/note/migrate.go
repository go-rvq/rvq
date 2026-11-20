package note

import (
	"context"

	"github.com/go-rvq/rvq/x/lifecycle"
	"github.com/go-rvq/rvq/x/migrate"
	"gorm.io/gorm"
)

// Key is the node of the package in the graphs of the commands
// (x/lifecycle).
const Key lifecycle.Key = "rvq/note"

// its migrations, in the graph of every package's (x/migrate): `app migrate`
// and the serve run them; Install still runs them too, for the applications
// that do not run the graph yet
func init() {
	migrate.Graph.Add(lifecycle.Node{
		Key:   Key,
		Start: migrate.Func(func(_ context.Context, db *gorm.DB) error { return Migrate(db) }),
	})
}

// Migrate makes the tables of the notes, off their old names.
func Migrate(db *gorm.DB) error {
	if err := renameLegacy(db); err != nil {
		return err
	}
	return db.AutoMigrate(Note{}, UserNote{})
}
