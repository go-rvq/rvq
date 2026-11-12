package perms

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/gorm"
)

type restoreThing struct {
	ID        uint
	Name      string
	DeletedAt gorm.DeletedAt
}

func (r *restoreThing) String() string { return r.Name }

// linesSource is the stack of a request: its header X-Lines, by ";".
type linesSource struct{}

func (linesSource) StackOf(r *http.Request) (string, string, func() (*perm.Stack, error)) {
	lines := r.Header.Get("X-Lines")
	if lines == "" {
		return "", "", nil
	}
	// one subject (a "<7>" in it would be a pattern of the policies), its
	// version the lines: made again when they change
	return "user:test", lines, func() (*perm.Stack, error) {
		return (&perm.Stack{}).Add("role: R", strings.Split(lines, ";")...), nil
	}
}

// restoreApp is an admin of things of a trash (AutoTrash): one, two and
// three deleted.
func restoreApp(t *testing.T) (*gorm.DB, *presets.ModelBuilder) {
	db := sharedTestDB(t)
	if err := db.AutoMigrate(&restoreThing{}); err != nil {
		t.Fatal(err)
	}
	for i, n := range []string{"one", "two", "three"} {
		db.Create(&restoreThing{ID: uint(i + 1), Name: n})
	}
	db.Delete(&restoreThing{}, []uint{1, 2, 3})
	p := presets.New(i18n.New())
	p.DataOperator(gorm2op.DataOperator(db))
	p.Permission(perm.New().Stacks(linesSource{}))
	AutoTrash(p, db)
	mb := p.Model(&restoreThing{}, presets.ModelWithID("things"))
	p.Build(http.NewServeMux())
	return db, mb
}

func ctxOf(lines ...string) *web.EventContext {
	r := httptest.NewRequest("POST", "/admin/things", nil)
	r.Header.Set("X-Lines", strings.Join(lines, ";"))
	return &web.EventContext{R: r}
}

func deletedIDs(db *gorm.DB) (ids []uint) {
	db.Unscoped().Model(&restoreThing{}).Where("deleted_at IS NOT NULL").Order("id").Pluck("id", &ids)
	return
}

// The permissions of the trash in the tree: browsing it of the listing
// (@trash), restoring of the record (@restore) — not an action !restore.
func TestTrashPermissionsTree(t *testing.T) {
	_, mb := restoreApp(t)
	var list, record []string
	mb.Builder().BuildPermissions().Tree().Walk(func(_ []*presets.PermNode, n *presets.PermNode) {
		var names []string
		for _, a := range n.Actions {
			names = append(names, a.Name)
		}
		switch n.Unique {
		case "admin:things:":
			list = names
		case "admin:things:<*>:":
			record = names
		}
	})
	if !slices.Contains(list, PermTrash) || slices.Contains(list, "!restore") {
		t.Errorf("the listing's: %q", list)
	}
	if !slices.Contains(record, PermRestore) {
		t.Errorf("the record's: %q", record)
	}
}

// Restoring: who may browse the trash, each record of its @restore — one
// not, none restored.
func TestRestoreRecords(t *testing.T) {
	db, mb := restoreApp(t)
	restore := func(ctx *web.EventContext, ids ...string) error {
		return restoreRecords(ctx, mb, db, false, ids)
	}
	if err := restore(ctxOf("admin:things:<*>:@restore"), "1"); !errors.Is(err, perm.PermissionDenied) {
		t.Errorf("no trash: %v", err)
	}
	if err := restore(ctxOf("admin:things:@trash"), "1"); err == nil || !strings.Contains(err.Error(), "one") {
		t.Errorf("no @restore: %v", err)
	}
	// one of two not: none restored, it named
	err := restore(ctxOf("admin:things:{@trash, <1>:@restore}"), "1", "2")
	if err == nil || !strings.Contains(err.Error(), "two") || strings.Contains(err.Error(), "one") {
		t.Errorf("one of two: %v", err)
	}
	if got := deletedIDs(db); !slices.Equal(got, []uint{1, 2, 3}) {
		t.Errorf("restored though refused: deleted %v", got)
	}
	if err := restore(ctxOf("admin:things:{@trash, <1>:@restore}"), "1"); err != nil {
		t.Fatal(err)
	}
	if err := restore(ctxOf("admin:things:{@trash, <*>:@restore}"), "2"); err != nil {
		t.Fatal(err)
	}
	if got := deletedIDs(db); !slices.Equal(got, []uint{3}) {
		t.Errorf("deleted %v, want [3]", got)
	}
	// who may delete a record (shared) may restore it
	if !slices.Contains(VerbDelete, PermRestore) {
		t.Error("the delete of a record shared without its restore")
	}
}
