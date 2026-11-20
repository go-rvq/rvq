package user

import (
	"context"

	"github.com/go-rvq/rvq/x/lifecycle"
	"github.com/go-rvq/rvq/x/migrate"
	"gorm.io/gorm"
)

// Key is the node of the package in the graphs of the commands
// (x/lifecycle).
const Key lifecycle.Key = "rvq/user"

// its migrations, in the graph of every package's (x/migrate): `app migrate`
// and the serve run them; the admin of the users still runs them too, for the applications
// that do not run the graph yet
func init() {
	migrate.Graph.Add(lifecycle.Node{
		Key:   Key,
		Start: migrate.Func(func(_ context.Context, db *gorm.DB) error { return AutoMigrateGroups(db) }),
	})
}
