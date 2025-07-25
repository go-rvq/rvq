package admin

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/activity"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/gorm"
)

func initPermission(b *presets.Builder, db *gorm.DB) {
	perm.Verbose = true

	b.Permission(
		perm.New().Policies(
			perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything),
		).SubjectsFunc(func(r *http.Request) []string {
			u := getCurrentUser(r)
			if u == nil {
				return nil
			}
			return u.GetRoles()
		}).ContextFunc(func(r *http.Request, objs []interface{}) perm.Context {
			c := make(perm.Context)
			for _, obj := range objs {
				switch v := obj.(type) {
				case *activity.ActivityLog:
					u := getCurrentUser(r)
					if u.GetID() == v.GetUserID() {
						c["is_authorized"] = true
					} else {
						c["is_authorized"] = false
					}
				}
			}
			return c
		}).DBPolicy(perm.NewDBPolicy(db)),
	)
}
