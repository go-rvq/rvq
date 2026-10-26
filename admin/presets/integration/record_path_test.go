package integration_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	. "github.com/go-rvq/rvq/admin/presets/integration"
)

type rpSetting struct {
	ID   uint
	Name string
}

type rpNote struct {
	ID   string `gorm:"primarykey"`
	Text string
}

// A record path goes to text and back (model and id of each level, root
// first) and to the admin's address of the record: the detail under its
// parents, the singleton's page. A path that names no model, an id the
// model cannot parse, or a singleton with an id are errors.
func TestRecordPath(t *testing.T) {
	db, err := NewNestedModelsDB()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&rpNote{}); err != nil {
		t.Fatal(err)
	}
	p := NewNestedModelsApp(db)
	settings := p.Model(&rpSetting{}, presets.ModelConfig().SetSingleton(true).SetId("settings"))
	notes := p.Model(&rpNote{}, presets.ModelConfig().SetId("notes"))
	notes.Listing("ID", "Text")
	notes.Detailing("Text")
	p.Build()

	surveys := p.GetModel(&NMSurvey{})
	places := p.GetModel(&NMPlace{})
	products := p.GetModel(&NMProduct{})
	id := func(mb *presets.ModelBuilder, v string) model.ID {
		t.Helper()
		r, err := mb.ParseRecordID(v)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	cases := []struct {
		name, path, href string
		mb               *presets.ModelBuilder
		parents          model.IDSlice
		id               model.ID
	}{
		{"root", surveys.MenuID() + ":1", "/admin/surveys/1", surveys, nil, id(surveys, "1")},
		{"child of child", surveys.MenuID() + ":1/" + places.MenuID() + ":2/" + products.MenuID() + ":3",
			"/admin/surveys/1/places/2/products/3", products,
			model.IDSlice{id(surveys, "1"), id(places, "2")}, id(products, "3")},
		{"singleton", "settings", "/admin/settings", settings, nil, model.ID{}},
		{"an id with : and /", "notes:a:b%2Fc", "/admin/notes/a:b%2Fc", notes, nil, id(notes, "a:b/c")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.mb.Info().RecordPath(c.parents, c.id)
			if got != c.path {
				t.Fatalf("path %q, want %q", got, c.path)
			}
			segs, err := p.ParseRecordPath(got)
			if err != nil {
				t.Fatal(err)
			}
			if last := segs[len(segs)-1]; last.Model != c.mb || last.ID.String() != c.id.String() {
				t.Errorf("parsed: %s %q, want %s %q", last.Model.MenuID(), last.ID.String(), c.mb.MenuID(), c.id.String())
			}
			href, err := p.RecordPathHref(got)
			if err != nil {
				t.Fatal(err)
			}
			if href != c.href {
				t.Errorf("href %q, want %q", href, c.href)
			}
		})
	}

	// the href opens the record
	if err = db.Create(&NMSurvey{ID: 1, Name: "S1", Places: []*NMPlace{{ID: 2, Name: "P2",
		Products: []*NMProduct{{ID: 3, Name: "X"}}}}}).Error; err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	p.ServeHTTP(w, httptest.NewRequest("GET", "/admin/surveys/1/places/2/products/3", nil))
	if w.Code != 200 {
		t.Errorf("the detail of the product: %d", w.Code)
	}

	for _, bad := range []string{"", "nothing:1", surveys.MenuID() + ":1/nothing:2",
		surveys.MenuID() + ":x", surveys.MenuID(), "settings:1", places.MenuID() + ":2"} {
		if _, err := p.RecordPathHref(bad); !errors.Is(err, presets.ErrRecordPath) {
			t.Errorf("%q: %v, want ErrRecordPath", bad, err)
		}
	}
}
