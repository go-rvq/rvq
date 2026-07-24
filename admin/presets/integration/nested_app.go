package integration

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// A four-level nested list-editor fixture: L0 (the edited root) → L1 → L2 → L3.
// Every level has a required Name (to exercise validation at any depth) and a
// has-many slice to the next level rendered by the list-editor table builder.

type NL3 struct {
	ID   uint
	L2ID uint
	Name string
	Pos  int
}

type NL2 struct {
	ID   uint
	L1ID uint
	Name string
	Pos  int
	L3s  []*NL3 `gorm:"foreignKey:L2ID"`
}

type NL1 struct {
	ID   uint
	L0ID uint
	Name string
	Pos  int
	L2s  []*NL2 `gorm:"foreignKey:L1ID"`
}

type NL0 struct {
	ID   uint
	Name string
	L1s  []*NL1 `gorm:"foreignKey:L0ID"`
}

const NestedRootURI = "l0s"

// NewNestedDB migrates the 4-level schema on a fresh in-memory SQLite database.
func NewNestedDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(nextMemName()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&NL0{}, &NL1{}, &NL2{}, &NL3{}); err != nil {
		return nil, err
	}
	return db, nil
}

// nestedEditing wires a required Name + an optional nested has-many field into a
// child model's editing. childField / childFB nest the next level (empty for the
// leaf). It uses the default (card) list-editor builder — not the table builder —
// because the card renders each item's full body, so a sub-list-editor nests
// inside it (the table builder only renders flat leaf columns).
func nestedEditing(itemMB *presets.ModelBuilder, childField string, childMB *presets.ModelBuilder, childFB *presets.FieldsBuilder) *presets.EditingBuilder {
	rows := [][]string{{"Name"}}
	if childField != "" {
		rows = append(rows, []string{childField})
	}
	ed := itemMB.Editing(&presets.FieldsSection{Rows: rows})
	ed.HiddenField("ID")
	ed.Field("Name").SetRequired(true)
	if childField != "" {
		ed.Field(childField).Nested(presets.NestedSlice(childMB, childFB))
	}
	return ed
}

// NewNestedApp builds the 4-level nested list-editor preset app.
func NewNestedApp(db *gorm.DB) *presets.Builder {
	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.Permission(perm.New().AllowAll())
	p.DataOperator(gorm2op.DataOperator(db))

	mb := p.Model(&NL0{}).URIName(NestedRootURI)

	// leaf up: L3, L2, L1
	l3MB := presets.NewModelBuilder(mb.Builder(), &NL3{})
	l3ed := nestedEditing(l3MB, "", nil, nil)

	l2MB := presets.NewModelBuilder(mb.Builder(), &NL2{})
	l2ed := nestedEditing(l2MB, "L3s", l3MB, &l3ed.FieldsBuilder)

	l1MB := presets.NewModelBuilder(mb.Builder(), &NL1{})
	l1ed := nestedEditing(l1MB, "L2s", l2MB, &l2ed.FieldsBuilder)

	// root L0
	ed := mb.Editing("Name", "L1s")
	ed.Field("Name").SetRequired(true)
	ed.Field("L1s").Nested(presets.NestedSlice(l1MB, &l1ed.FieldsBuilder))

	// preload the whole tree on read so the edit form shows every level.
	mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		b := do.(*gorm2op.DataOperatorBuilder).
			WrapPrepare(func(old gorm2op.Preparer) gorm2op.Preparer {
				return func(db *gorm.DB, mode gorm2op.Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
					db = old(db, mode, obj, id, params, ctx)
					if mode.Is(gorm2op.Fetch, gorm2op.Search) {
						db = db.Preload("L1s.L2s.L3s")
					}
					return db
				}
			})
		return gorm2op.SaveHasManyAssociation("L1s").Build(b)
	})

	return p
}

// NewNestedSeededHandler builds NewNestedApp seeded with one L0 tree:
// L0 "root" → L1 "a" → L2 "a1" → L3 "a1x".
func NewNestedSeededHandler() (http.Handler, error) {
	db, err := NewNestedDB()
	if err != nil {
		return nil, err
	}
	root := &NL0{ID: 1, Name: "root", L1s: []*NL1{
		{ID: 1, L0ID: 1, Name: "a", Pos: 0, L2s: []*NL2{
			{ID: 1, L1ID: 1, Name: "a1", Pos: 0, L3s: []*NL3{
				{ID: 1, L2ID: 1, Name: "a1x", Pos: 0},
			}},
		}},
	}}
	if err = db.Create(root).Error; err != nil {
		return nil, err
	}
	return NewNestedApp(db), nil
}
