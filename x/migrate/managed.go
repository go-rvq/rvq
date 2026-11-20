package migrate

import (
	"context"

	"github.com/go-rvq/rvq/x/lifecycle"
	"gorm.io/gorm"
)

// managedPlugin marks a database whose migrations the graph runs (Manage): a
// plugin of its gorm configuration, whose map of plugins every session of it
// shares.
type managedPlugin struct{}

const managedName = "rvq:migrate:managed"

func (managedPlugin) Name() string                 { return managedName }
func (managedPlugin) Initialize(db *gorm.DB) error { return nil }

// Manage says the migrations of db are the graph's — `app migrate`, the
// node Key of the command that serves —: building an application on it
// migrates nothing (Auto).
func Manage(db *gorm.DB) {
	if !Managed(db) {
		_ = db.Use(managedPlugin{})
	}
}

// Managed tells whether the migrations of db are the graph's (Manage).
func Managed(db *gorm.DB) bool {
	_, ok := db.Config.Plugins[managedName]
	return ok
}

// Auto runs f, a migration, as an application is built — unless the
// migrations of db are the graph's (Manage), where f is a node of it. The
// constructors of the packages migrate through it, for the applications that
// do not run the graph.
func Auto(db *gorm.DB, f func(db *gorm.DB) error) error {
	if Managed(db) {
		return nil
	}
	return f(db)
}

// Offer runs f, a migration of key — of an application: its order the order
// it is built in —, as the application is built, unless the graph runs the
// migrations of db (Manage): then f is offered to ms (presets.Builder.
// Migrations), after the one offered before it (lifecycle.Nodes.Chain), and
// runs there. A key offered already is not again.
//
// after are the keys it runs after besides (a package's migration whose
// tables it writes: rvq/media's uses of a user saved), optional.
func Offer(ms *lifecycle.Nodes, db *gorm.DB, key lifecycle.Key, f func(db *gorm.DB) error, after ...lifecycle.Key) error {
	if !Managed(db) {
		return f(db)
	}
	if !ms.Has(key) {
		ms.Chain(lifecycle.Node{Key: key, After: after, Optional: after,
			Start: Func(func(_ context.Context, db *gorm.DB) error { return f(db) })})
	}
	return nil
}
