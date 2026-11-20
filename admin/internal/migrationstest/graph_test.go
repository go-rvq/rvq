// Package migrationstest checks the graph of the migrations of the packages
// of the admin (x/migrate): each registers its own in its init().
package migrationstest

import (
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/dashboard"
	"github.com/go-rvq/rvq/admin/media"
	"github.com/go-rvq/rvq/admin/microsite"
	"github.com/go-rvq/rvq/admin/note"
	"github.com/go-rvq/rvq/admin/packages/accesskeys"
	"github.com/go-rvq/rvq/admin/packages/backup"
	dbtools "github.com/go-rvq/rvq/admin/packages/db-tools"
	"github.com/go-rvq/rvq/admin/packages/gitedit"
	history "github.com/go-rvq/rvq/admin/packages/history/admin"
	mailsender "github.com/go-rvq/rvq/admin/packages/mail-sender"
	"github.com/go-rvq/rvq/admin/packages/notifications"
	orgs "github.com/go-rvq/rvq/admin/packages/orgs/models"
	people "github.com/go-rvq/rvq/admin/packages/people/models"
	"github.com/go-rvq/rvq/admin/packages/shared"
	"github.com/go-rvq/rvq/admin/packages/user"
	validators "github.com/go-rvq/rvq/admin/packages/validators/models"
	"github.com/go-rvq/rvq/admin/pagebuilder"
	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/admin/seo"
	"github.com/go-rvq/rvq/admin/userdocs"
	"github.com/go-rvq/rvq/admin/worker"
	"github.com/go-rvq/rvq/x/lifecycle"
	"github.com/go-rvq/rvq/x/migrate"
)

// Every package's migration registered, the graph whole (no dependency
// missing, no cycle); what depends after its dependency.
func TestGraph(t *testing.T) {
	keys := []lifecycle.Key{
		dashboard.Key, media.Key, microsite.Key, note.Key, accesskeys.Key, backup.Key,
		dbtools.Key, gitedit.Key, history.Key, mailsender.Key, notifications.Key, orgs.Key,
		people.Key, shared.Key, user.Key, validators.Key, pagebuilder.Key, role.Key, seo.Key,
		userdocs.Key, worker.Key,
	}
	for _, k := range keys {
		if !migrate.Graph.Has(k) {
			t.Errorf("not registered: %s", k)
		}
	}
	plan, err := migrate.Graph.Plan()
	if err != nil {
		t.Fatal(err)
	}
	before := func(a, b lifecycle.Key) {
		t.Helper()
		if i, j := strings.Index(plan, ". "+string(a)+"\n"), strings.Index(plan, ". "+string(b)+" "); i < 0 || j < 0 || i > j {
			t.Errorf("%s not before %s:\n%s", a, b, plan)
		}
	}
	before(validators.Key, people.Key)
	before(role.Key, accesskeys.Key)
}
