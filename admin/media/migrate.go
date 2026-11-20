package media

import (
	"context"

	"github.com/go-rvq/rvq/x/lifecycle"
	"github.com/go-rvq/rvq/x/migrate"
	"gorm.io/gorm"

	"github.com/go-rvq/rvq/admin/media/media_library"
)

// Key is the node of the package in the graphs of the commands
// (x/lifecycle).
const Key lifecycle.Key = "rvq/media"

// its migrations, in the graph of every package's (x/migrate): `app migrate`
// and the serve run them; its configure still runs them too, for the applications
// that do not run the graph yet
func init() {
	migrate.Graph.Add(lifecycle.Node{
		Key:   Key,
		Start: migrate.Func(func(_ context.Context, db *gorm.DB) error { return Migrate(db) }),
	})
}

// Migrate makes the tables of the media library, of the uses of its files
// and of the usages of the media boxes.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&media_library.MediaLibrary{}, &MediaLibraryUse{}, &MediaBoxUsage{})
}
