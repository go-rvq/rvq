package perms

import (
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/gorm"
)

type trashThing struct {
	ID        uint
	Name      string
	DeletedAt gorm.DeletedAt
}

// A second SetupTrash of a model (its own, after AutoTrash's) adds no second
// tab nor restore action: it only completes the options of the first.
func TestSetupTrashOnce(t *testing.T) {
	db := sharedTestDB(t)
	if err := db.AutoMigrate(&trashThing{}); err != nil {
		t.Fatal(err)
	}
	p := presets.New(i18n.New())
	p.DataOperator(gorm2op.DataOperator(db))
	AutoTrash(p, db)
	mb := p.Model(&trashThing{})
	if trashSetups[mb] == nil {
		t.Fatal("AutoTrash did not set the trash up")
	}

	restored := false
	SetupTrash(mb, db, TrashOptions{TableName: "trash_things",
		OnRestore: func(*web.EventContext, []string) error { restored = true; return nil }})

	n := 0
	for _, a := range mb.Listing().GetBulkActions() {
		if a.Name() == "restore" {
			n++
		}
	}
	if n != 1 {
		t.Fatalf("%d restore actions, want 1", n)
	}
	o := trashSetups[mb]
	if o.OnRestore == nil {
		t.Fatal("the OnRestore of the second call was lost")
	}
	_ = o.OnRestore(nil, nil)
	if !restored {
		t.Fatal("the OnRestore of the second call did not run")
	}
}
