package integration

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Item is the has-many child rendered by the list editor.
type Item struct {
	ID        uint
	ProductID uint
	Label     string
	Pos       int
}

// Product is the parent record edited through the presets Editing form.
type Product struct {
	ID    uint
	Name  string
	Items []*Item `gorm:"foreignKey:ProductID"`
}

// The model URI and the nested field the browser posts under.
const (
	ProductURI = "products"
	ItemsField = "Items"
)

var dbSeq int64

// nextMemName returns a unique shared-cache in-memory SQLite DSN, so each opened
// database is isolated from the others while its connection pool stays pointed at
// a single in-memory database.
func nextMemName() string {
	return fmt.Sprintf("file:presetstest_%d?mode=memory&cache=shared", atomic.AddInt64(&dbSeq, 1))
}

// NewDB opens a fresh, isolated in-memory SQLite database with the schema
// migrated.
func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(nextMemName()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&Product{}, &Item{}); err != nil {
		return nil, err
	}
	return db, nil
}

// NewApp builds the preset app: a Product model whose Items has-many field is
// rendered by the list-editor table builder, with Name required and the
// has-many association persisted by gorm2op.SaveHasManyAssociation. This is the
// shared fixture exercised by both the Go list-editor tests and the bun/TypeScript
// integration tests.
func NewApp(db *gorm.DB) *presets.Builder {
	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.Permission(perm.New().AllowAll())
	p.DataOperator(gorm2op.DataOperator(db))

	mb := p.Model(&Product{}).URIName(ProductURI)

	// A detailing, so the tests can exercise its Edit button — the form host: the
	// button only turns presets.DetailingEditScope.show on.
	mb.Detailing("Name")

	// The child model + its editing fields (the columns the table renders).
	// Label is required so a removed row with an empty Label would fail validation
	// unless the deleted row is skipped.
	itemMB := presets.NewModelBuilder(mb.Builder(), &Item{})
	ied := itemMB.Editing(&presets.FieldsSection{
		Rows: [][]string{{"Label"}},
	})
	ied.HiddenField("ID")
	ied.Field("Label").SetRequired(true)

	ed := mb.Editing("Name", ItemsField)
	ed.Field("Name").SetRequired(true)

	// Persist the has-many association (deleted children are removed) and preload
	// Items on every read (Search + Fetch) so both the listing and the edit form
	// show the children.
	mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		b := do.(*gorm2op.DataOperatorBuilder).
			WrapPrepare(func(old gorm2op.Preparer) gorm2op.Preparer {
				return func(db *gorm.DB, mode gorm2op.Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
					db = old(db, mode, obj, id, params, ctx)
					if mode.Is(gorm2op.Fetch, gorm2op.Search) {
						db = db.Preload(ItemsField)
					}
					return db
				}
			})
		return gorm2op.SaveHasManyAssociation(ItemsField).Build(b)
	})

	ed.Field(ItemsField).
		Nested(presets.NestedSlice(itemMB, &ied.FieldsBuilder).
			SetComponentBuilder(presets.NewListEditorTableBuilder().
				LayoutFunc(func(fb *presets.FieldsBuilder) presets.FieldsLayout {
					return presets.FieldsLayout{"Label"}
				})))

	return p
}

// NewSeededHandler builds NewApp on a fresh in-memory SQLite database seeded with
// one product (id=1, "P1") holding two items (A, B). It is the entry point for
// out-of-process integration servers (e.g. the bun TypeScript tests) that serve
// the very same app the Go tests exercise.
func NewSeededHandler() (http.Handler, error) {
	db, err := NewDB()
	if err != nil {
		return nil, err
	}
	if err = db.Create(&Product{ID: 1, Name: "P1", Items: []*Item{
		{ID: 1, ProductID: 1, Label: "A", Pos: 0},
		{ID: 2, ProductID: 1, Label: "B", Pos: 1},
	}}).Error; err != nil {
		return nil, err
	}
	return NewApp(db), nil
}
