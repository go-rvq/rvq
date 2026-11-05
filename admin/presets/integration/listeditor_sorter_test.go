package integration_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/go-rvq/rvq/admin/presets/integration"
)

// The list editor's sorter (two items or more) gets its items as an array —
// the vx-draggable of the sorter iterates them —, never as a string of JSON.
func TestListEditorSorterItemsArray(t *testing.T) {
	h, err := NewSeededHandler()
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/"+ProductURI+"/1/edit", nil))
	body := strings.NewReplacer(`&`, "&", `<`, "<", `>`, ">", `&#34;`, `"`, `\"`, `"`).Replace(w.Body.String())
	if !strings.Contains(body, `items: [{"index"`) {
		i := strings.Index(body, `"index"`)
		t.Errorf("the sorter's items not an array (%d): %.200s", w.Code, body[max(0, i-60):])
	}
}
