package examples_admin

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/admin/worker"
	"gorm.io/gorm"
)

func ActionWorkerExample(b *presets.Builder, db *gorm.DB) http.Handler {
	if err := db.AutoMigrate(&ExampleResource{}); err != nil {
		panic(err)
	}

	b.DataOperator(gorm2op.DataOperator(db))

	mb := b.Model(&ExampleResource{})
	mb.Listing().ActionsAsMenu(true)

	wb := worker.NewWithQueue(db, Que)
	b.Use(wb)
	addActionJobs(mb, wb)
	wb.Listen()
	return b
}
