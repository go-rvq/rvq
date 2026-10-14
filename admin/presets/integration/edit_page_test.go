package integration_test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/x/i18n"
)

type EditPageProduct struct {
	ID   uint
	Name string
}

type EditPageSetting struct {
	ID   uint
	Name string
}

// The edit form as a page: `/{listing}/{id}/edit` — the address an edit
// overlay shows while it is open — renders the form of the record of the
// address (its id the path's, not a query's); a singleton with a detailing
// has its form at `/{singleton}/edit`.
func TestEditPage(t *testing.T) {
	db := TestDB
	if err := db.AutoMigrate(&EditPageProduct{}, &EditPageSetting{}); err != nil {
		t.Fatal(err)
	}
	db.Exec("DELETE FROM edit_page_products")
	db.Exec("DELETE FROM edit_page_settings")
	p := &EditPageProduct{Name: "zz-edit-product"}
	db.Create(p)
	db.Create(&EditPageSetting{ID: 1, Name: "zz-edit-setting"})

	b := presets.New(i18n.New()).URIPrefix("/admin")
	b.DataOperator(gorm2op.DataOperator(db))
	b.Model(&EditPageProduct{}).URIName("edit-products").Detailing("Name")
	b.Model(&EditPageSetting{}, presets.ModelConfig().SetSingleton(true).SetId("edit-settings")).
		URIName("edit-settings").Detailing("Name")

	for _, c := range []struct{ path, want string }{
		{fmt.Sprintf("/admin/edit-products/%d/edit", p.ID), "zz-edit-product"},
		{"/admin/edit-settings/edit", "zz-edit-setting"},
	} {
		w := httptest.NewRecorder()
		b.ServeHTTP(w, httptest.NewRequest("GET", c.path, nil))
		body := w.Body.String()
		if w.Code != 200 || !strings.Contains(body, c.want) || strings.Contains(body, ">404<") {
			t.Errorf("%s: %d, the form of the record (%q) not there:\n%.600s", c.path, w.Code, c.want, body)
		}
	}
}
