package worker

import (
	"context"
	"github.com/tnclong/go-que/pg"

	"github.com/go-rvq/rvq/x/lifecycle"
	"github.com/go-rvq/rvq/x/migrate"
	"gorm.io/gorm"
)

// its migrations, in the graph of every package's (x/migrate): `app migrate`
// and the serve run them; New still runs them too, for the applications
// that do not run the graph yet
func init() {
	migrate.Graph.Add(lifecycle.Node{
		Key:   Key,
		Start: migrate.Func(func(_ context.Context, db *gorm.DB) error { return Migrate(db) }),
	})
}

// Migrate makes the tables of the jobs, off the old "qor_" prefix, and the
// queue's.
func Migrate(db *gorm.DB) error {
	if err := renameLegacy(db); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Job{}, &JobInstance{}, &JobLog{}, &GoQueError{}); err != nil {
		return err
	}
	// the queue's table (goque_jobs): its library makes it
	rdb, err := db.DB()
	if err != nil {
		return err
	}
	_, err = pg.New(rdb)
	return err
}
