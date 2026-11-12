package role

import (
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// A line removed in the form — the value posted without it — is gone once
// the role is saved; the one kept stays.
func TestRolePermissionRemoved(t *testing.T) {
	mine := Role{ID: uuid.New(), Name: "Mine", Permissions: pq.StringArray{"admin:kept:*", "admin:gone:*"}}
	db, _, h := systemApp(t, func(db *gorm.DB) { db.Create(&mine) })
	event := func(name string, fields [][2]string) string {
		mb := multipartestutils.NewMultipartBuilder().PageURL("/admin/roles").EventFunc(name).
			Query(presets.ParamID, mine.ID.String())
		for _, f := range fields {
			mb = mb.AddField(f[0], f[1])
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, mb.BuildEventFuncRequest())
		return w.Body.String()
	}
	stamp := regexp.MustCompile(`__formSign\\?":\s*\\?"([^"\\]+)`).FindStringSubmatch(event(actions.Edit, nil))
	if stamp == nil {
		t.Fatal("the edit form carries no stamp")
	}
	body := event(actions.Update, [][2]string{
		{"Name", "Mine"}, {presets.RecordStampFormKey, stamp[1]},
		{"Permissions", "admin:kept:*"},
	})
	if got := linesOf(db, mine); got != "admin:kept:*" {
		t.Errorf("the lines after the removal: %q\n%.400s", got, body)
	}
}
