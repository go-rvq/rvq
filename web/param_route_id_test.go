package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Param reads the route's path value before the query. A handler mounted under
// a pattern that names {id} therefore sees the record in the path, never the id
// its own event sent — which is why an event that carries an "id" of its own has
// to read the form directly.
func TestParamPrefersTheRoutePathValue(t *testing.T) {
	mux := http.NewServeMux()

	var fromParam, fromForm string
	mux.HandleFunc("/admin/pages/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := &EventContext{R: r}
		fromParam = ctx.Param("id")
		fromForm = r.FormValue("id")
	})

	// The edit page of record 4, and an event asking for media library file 1.
	req := httptest.NewRequest(http.MethodPost, "/admin/pages/4?id=1", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)

	if fromParam != "4" {
		t.Errorf("Param(\"id\") = %q, want the route's %q", fromParam, "4")
	}
	if fromForm != "1" {
		t.Errorf("FormValue(\"id\") = %q, want the event's %q", fromForm, "1")
	}
	if fromParam == fromForm {
		t.Error("the two agree, so this test no longer pins the collision")
	}
}
