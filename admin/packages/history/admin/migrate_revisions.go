package admin

import (
	"context"

	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/lifecycle"
	"github.com/go-rvq/rvq/x/migrate"
	"gorm.io/gorm"
)

// RevisionsKey is the node of the migration of the revisions of table in
// the graphs of the commands (x/lifecycle).
func RevisionsKey(table string) lifecycle.Key {
	return lifecycle.Key("rvq/history/revisions/" + table)
}

// offerRevisionsMigration offers p the migration of the table of the
// revisions of h's model — and its backfill and seed (backfillAndSeed) —,
// which Build still runs too, unless the graph runs the migrations
// (migrate.Manage).
func offerRevisionsMigration(p *presets.Builder, h *ModelHistory) {
	key := RevisionsKey(h.table)
	if p == nil || p.Migrations().Has(key) {
		return
	}
	table := h.table
	p.Migrations().Chain(lifecycle.Node{
		Key: key,
		Start: migrate.Func(func(_ context.Context, db *gorm.DB) error {
			if err := db.Table(table).AutoMigrate(&histmodels.Revision{}); err != nil {
				return err
			}
			h.backfillAndSeed()
			return nil
		}),
	})
}
