package role

import (
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// A permission removed in the form — posted as the list editor posts it, its
// ID and __deleted — is gone once the role is saved; the one kept stays.
func TestRolePermissionRemoved(t *testing.T) {
	mine := Role{ID: uuid.New(), Name: "Mine"}
	var kept, gone perm.DefaultDBPolicy
	db, _, h := systemApp(t, func(db *gorm.DB) {
		db.Create(&mine)
		kept = perm.DefaultDBPolicy{ReferID: mine.ID.String(), Subject: "Mine", Effect: perm.Allowed,
			Actions: pq.StringArray{"*"}, Resources: pq.StringArray{"admin:kept:*"}}
		gone = perm.DefaultDBPolicy{ReferID: mine.ID.String(), Subject: "Mine", Effect: perm.Allowed,
			Actions: pq.StringArray{"*"}, Resources: pq.StringArray{"admin:gone:*"}}
		db.Create(&kept)
		db.Create(&gone)
	})
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
		{"Permissions.__present", "1"},
		{"Permissions[0].__pos", "0"}, {"Permissions[0].__index", "0"}, {"Permissions[0].ID", kept.ID.String()},
		{"Permissions[0].Effect", perm.Allowed},
		{"Permissions[0].Resources[0].__value", "admin:kept:*"}, {"Permissions[0].Resources[0].$id", "1"},
		{"Permissions[1].__pos", "1"}, {"Permissions[1].__index", "1"}, {"Permissions[1].ID", gone.ID.String()},
		{"Permissions[1].__deleted", "true"},
		{"Permissions[1].Effect", perm.Allowed},
		{"Permissions[1].Resources[0].__value", "admin:gone:*"}, {"Permissions[1].Resources[0].$id", "1"},
	})
	if got := resourcesOf(db, mine); got != "admin:kept:*=allow/Mine" {
		t.Errorf("the permissions after the removal: %q\n%.400s", got, body)
	}
}
