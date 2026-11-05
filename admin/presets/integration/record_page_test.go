package integration_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	. "github.com/go-rvq/rvq/admin/presets/integration"
	"github.com/go-rvq/rvq/web"
)

type rpgNote struct {
	ID   string `gorm:"primarykey"`
	Text string
}

// A page of the detail of a record (RecordLayout) fetches its record: the
// detail's fetch.
func TestRecordPageFetches(t *testing.T) {
	db, err := NewNestedModelsDB()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&rpgNote{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&rpgNote{ID: "n1", Text: "zz-the-text"})
	p := NewNestedModelsApp(db)
	notes := p.Model(&rpgNote{}, presets.ModelConfig().SetId("rpg_notes"))
	notes.Detailing().PagesRegistrator().New(presets.HttpPage("/text")).
		RecordLayout(func(ctx *web.EventContext, _ model.ID, _ []model.ID, fetch func() (any, error)) (r web.PageResponse, err error) {
			obj, err := fetch()
			if err != nil || obj == nil {
				return r, err
			}
			r.Body = h.Div(h.Text(obj.(*rpgNote).Text))
			return
		}).Build()
	p.Build()
	w := httptest.NewRecorder()
	p.ServeHTTP(w, httptest.NewRequest("GET", notes.Info().DetailingHref("n1")+"/text", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "zz-the-text") {
		t.Errorf("the page of the record: %d %.300s", w.Code, w.Body.String())
	}
}
