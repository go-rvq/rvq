package publish_test

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/go-rvq/rvq/admin/docs/docsrc/examples/examples_admin"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/theplant/testenv"
	"gorm.io/gorm"
)

var (
	DB             *gorm.DB
	SQLDB          *sql.DB
	PresetsBuilder *presets.Builder
)

func TestMain(m *testing.M) {
	env, err := testenv.New().DBEnable(true).SetUp()
	if err != nil {
		panic(err)
	}
	defer env.TearDown()

	DB = env.DB
	SQLDB, err = DB.DB()
	if err != nil {
		panic(err)
	}
	PresetsBuilder = presets.New().DataOperator(gorm2op.DataOperator(DB)).URIPrefix("/samples/publish-example")
	examples_admin.PublishExample(PresetsBuilder, DB)

	m.Run()
}

type Flow struct {
	db *gorm.DB
	h  http.Handler
}
