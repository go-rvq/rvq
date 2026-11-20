package activity

import (
	"context"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/lifecycle"
	"github.com/go-rvq/rvq/x/migrate"
	"gorm.io/gorm"
)

// Key is the node of the log of the activities in the graphs of the
// commands (x/lifecycle).
const Key lifecycle.Key = "rvq/activity"

// offerMigration offers p the migration of the log — the application's
// model (New) —, which New still runs too, for the applications that do not
// run the graph yet.
func (ab *Builder) offerMigration(p *presets.Builder) {
	if p.Migrations().Has(Key) {
		return
	}
	logModel := ab.logModel
	p.Migrations().Chain(lifecycle.Node{
		Key: Key,
		Start: migrate.Func(func(_ context.Context, db *gorm.DB) error {
			return db.AutoMigrate(logModel)
		}),
	})
}
