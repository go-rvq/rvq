package presets_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type tabbedThing struct {
	ID         uint
	Name, Note string
}

// A form in tabs drawn again — a row added, an error — opens the tab it was
// on (the form key FieldsTabFormKey, set as the tab changes), not the first;
// none in the form, the first.
func TestFieldsTabsKeptByTheForm(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&tabbedThing{}); err != nil {
		t.Fatal(err)
	}
	b := presets.New(i18n.New()).URIPrefix("/admin").DataOperator(gorm2op.DataOperator(db))
	b.Model(&tabbedThing{}).URIName("things").Editing(presets.NewFieldsTabs("abas",
		&presets.FieldsGroup{Name: "a", Items: presets.FieldsLayout{"Name"}},
		&presets.FieldsGroup{Name: "b", Items: presets.FieldsLayout{"Note"}},
	))
	b.Build()
	form := func(tab string) string {
		mb := multipartestutils.NewMultipartBuilder().PageURL("/admin/things").EventFunc(actions.New)
		if tab != "" {
			mb.AddField(presets.FieldsTabFormKey("abas"), tab)
		}
		w := httptest.NewRecorder()
		b.ServeHTTP(w, mb.BuildEventFuncRequest())
		return w.Body.String()
	}
	if body := form("b"); !strings.Contains(body, "{tab: 1}") {
		t.Errorf("the tab of the form not open: %.300s", body)
	}
	body := form("")
	if !strings.Contains(body, "{tab: 0}") {
		t.Errorf("no tab in the form, not the first: %.300s", body)
	}
	// the tab chosen goes with the form
	if !strings.Contains(body, presets.FieldsTabFormKey("abas")) {
		t.Error("the change of tab does not set the form")
	}
}
