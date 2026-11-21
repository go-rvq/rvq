package presets_test

import (
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type statusThing struct {
	ID     uint
	Status string
}

// The restriction of a deletion is asked of the record as it is — read —,
// not of one made of its id alone: a rule of its fields ("a draft only")
// lets the draft be deleted and keeps the other.
func TestDeletingRestrictionOfTheRecord(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&statusThing{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&[]statusThing{{ID: 1, Status: "draft"}, {ID: 2, Status: "done"}})
	b := presets.New(i18n.New()).URIPrefix("/admin").DataOperator(gorm2op.DataOperator(db))
	mb := b.Model(&statusThing{}).URIName("things")
	mb.DeletingRestriction.ObjHandler(presets.OkObjHandlerFunc(func(obj any, ctx *web.EventContext) (bool, bool) {
		draft := obj.(*statusThing).Status == "draft"
		return !draft, !draft
	}))
	b.Build()
	del := func(id string) {
		r := multipartestutils.NewMultipartBuilder().PageURL("/admin/things?id=" + id).
			EventFunc(actions.DoDelete).BuildEventFuncRequest()
		r = r.WithContext(web.WithSkipFormSign(r.Context()))
		b.ServeHTTP(httptest.NewRecorder(), r)
	}
	del("1")
	del("2")
	var left []statusThing
	db.Order("id").Find(&left)
	if len(left) != 1 || left[0].ID != 2 {
		t.Errorf("left: %+v; want the one done only", left)
	}
}
