package presets_test

import (
	"errors"
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

// keyedThing gets its key from a hook as it is created (gorm runs
// BeforeSave before BeforeCreate: the check is after it), and is not saved
// without a name.
type keyedThing struct {
	ID   string
	Name string
}

const keyedThingKey = "key-from-the-hook"

func (k *keyedThing) BeforeCreate(*gorm.DB) error {
	k.ID = keyedThingKey
	return nil
}

func (k *keyedThing) AfterCreate(*gorm.DB) error {
	if k.Name == "" {
		return errors.New("no name")
	}
	return nil
}

// A creation that fails after a hook gave the record its key: the form
// shown again — in a drawer — is still a new one's — it has no key, the one given was rolled
// back with the record —, and a second save creates it.
func TestCreateFailedKeepsNoKey(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&keyedThing{}); err != nil {
		t.Fatal(err)
	}
	b := presets.New(i18n.New()).URIPrefix("/admin").DataOperator(gorm2op.DataOperator(db))
	b.Model(&keyedThing{}).URIName("things").Editing("Name")
	b.Build()
	create := func(name string) string {
		r := multipartestutils.NewMultipartBuilder().PageURL("/admin/things").EventFunc(actions.Create).
			AddField("Name", name).AddField(presets.ParamOverlay, string(actions.RightDrawer)).BuildEventFuncRequest()
		w := httptest.NewRecorder()
		b.ServeHTTP(w, r)
		return w.Body.String()
	}
	if body := create(""); !strings.Contains(body, "no name") || strings.Contains(body, keyedThingKey) {
		t.Errorf("the form after the failure: the error %v, the key of the hook %v",
			strings.Contains(body, "no name"), strings.Contains(body, keyedThingKey))
	}
	create("Thing")
	var n int64
	db.Model(&keyedThing{}).Where("name = ?", "Thing").Count(&n)
	if n != 1 {
		t.Errorf("created after the failure: %d", n)
	}
}
