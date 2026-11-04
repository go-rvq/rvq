package listeditor

import (
	"net/http"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	. "github.com/go-rvq/rvq/web/multipartestutils"
)

// The form on its own page (/new, /{id}/edit, no overlay) lives in the
// portal of the page (actions.Edit): adding a row answers into it — it
// replaced the body of the page, the header and the save gone.
func TestAddRow_PageModeUpdatesTheFormPortal(t *testing.T) {
	app, err := twoPersistedItems()
	if err != nil {
		t.Fatal(err)
	}
	RunCase(t, TestCase{
		Name: "add row on the page of the form",
		ReqFunc: func() *http.Request {
			return signed(t, app, NewMultipartBuilder().
				PageURL("/admin/products").
				EventFunc(actions.AddRowEvent).
				Query(presets.ParamID, "1").
				Query(presets.ParamAddRowFormKey, "Items").
				AddField("Items.__present", "1").
				AddField("Items[0].ID", "1").
				AddField("Items[0].Label", "A").
				BuildEventFuncRequest())
		},
		EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
			if er.Body != "" {
				t.Fatalf("the body of the page replaced")
			}
			if len(er.UpdatePortals) != 1 || er.UpdatePortals[0].Name != actions.Edit {
				t.Fatalf("portals %+v, want the form's (%s)", er.UpdatePortals, actions.Edit)
			}
		},
	}, app)
}
