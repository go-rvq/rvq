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

// Migrate makes the tables of the media library, of the uses of its files,
// of the usages of the media boxes and of the formats of the images.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&media_library.MediaLibrary{}, &MediaLibraryUse{}, &MediaBoxUsage{},
		&MediaFormat{}, &MediaBoxUsageFormat{}); err != nil {
		return err
	}
	// the formats of a field once a JSON of its own (never released): now
	// MediaBoxUsageFormat
	if m := db.Migrator(); m.HasColumn(&MediaBoxUsage{}, "formats") {
		return m.DropColumn(&MediaBoxUsage{}, "formats")
	}
	return nil
}
