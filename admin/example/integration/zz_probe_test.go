package integration_test

import (
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/admin/example/admin"
	. "github.com/go-rvq/rvq/web/multipartestutils"
)

func TestZZProbe(t *testing.T) {
	h := admin.TestHandler(TestDB)
	dbr, _ := TestDB.DB()
	pageBuilderContainerTestData.TruncatePut(dbr)

	for _, u := range []string{
		"/page_builder/pages/editors/10_2024-05-21-v01_International",
		"/page_builder/pages/editors/10_2024-05-21-v01_International?__execute_event__=page_builder_AddContainerEvent&modelName=BrandGrid",
		"/page_builder/pages/editors",
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, NewMultipartBuilder().PageURL(u).BuildEventFuncRequest())
		body := w.Body.String()
		if len(body) > 60 {
			body = body[:60]
		}
		t.Logf("%-100s -> %d %q", u, w.Code, body)
	}
}
