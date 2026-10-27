package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	. "github.com/go-rvq/rvq/admin/presets/integration"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/multipartestutils"
)

type amNote struct {
	ID   string `gorm:"primarykey"`
	Text string
}

// An action of the menu of a detail that is a page — done in its own
// dialog — shows the page again once done, as a button of the toolbar does:
// its click says so (ParamToolbarAction "page"), and its response goes to
// the page again and closes its dialog. A detail in an overlay does not say
// it.
func TestActionMenuOfDetailPage(t *testing.T) {
	db, err := NewNestedModelsDB()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&amNote{}); err != nil {
		t.Fatal(err)
	}
	p := NewNestedModelsApp(db)
	done := ""
	notes := p.Model(&amNote{}, presets.ModelConfig().SetId("am-notes"))
	notes.Listing("ID", "Text")
	notes.Detailing("Text").Action("Resend").UpdateFunc(func(id string, ctx *web.EventContext) error {
		done = id
		return nil
	})
	p.Build()
	db.Create(&amNote{ID: "n1", Text: "N"})

	serve := func(r *http.Request) string {
		t.Helper()
		w := httptest.NewRecorder()
		p.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("%s: %d %s", r.URL, w.Code, w.Body.String())
		}
		return w.Body.String()
	}
	if page := serve(httptest.NewRequest("GET", "/admin/am-notes/n1", nil)); !strings.Contains(page, presets.ParamToolbarAction) {
		t.Error("the menu of the detail page does not say its action is of a page")
	}
	drawer := serve(multipartestutils.NewMultipartBuilder().PageURL("/admin/am-notes").EventFunc(actions.Detailing).
		Query(presets.ParamID, "n1").Query(presets.ParamOverlay, string(actions.RightDrawer)).BuildEventFuncRequest())
	if strings.Contains(drawer, presets.ParamToolbarAction) {
		t.Error("the menu of a detail in a drawer says its action is of a page")
	}

	got := serve(multipartestutils.NewMultipartBuilder().PageURL("/admin/am-notes").EventFunc(actions.DoAction).
		Query(presets.ParamAction, "Resend").Query(presets.ParamID, "n1").
		Query(presets.ParamOverlay, string(actions.Dialog)).Query(presets.ParamToolbarAction, "page").BuildEventFuncRequest())
	if done != "n1" || strings.Contains(got, `"pushState":null`) || !strings.Contains(got, "closer.show = false") {
		t.Errorf("done %q: %s", done, got)
	}
}
