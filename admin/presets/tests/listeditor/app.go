// Package listeditor holds a self-contained preset app and an integration test
// that exercises the list-editor deleted-state behaviour end to end, validating
// the rendered UI (not just the persisted data).
//
// # The domain
//
// A Product (the parent, edited via the presets Editing form) owns a has-many
// slice of Item children rendered by the list-editor table component. Name is a
// required field on the Product, used to force a validation failure so the form
// re-renders with the children still in their edited state.
//
//	Product ─┬─ Name  (required)
//	         └─ Items [] ── Label   (rendered by NewListEditorTableBuilder)
//
// # The reactive form model (what the browser posts)
//
// The list editor keeps each child's state on a single flat reactive `form`
// object. Each row's form key is bracketed by its stable __index (NOT its slice
// position), so the key is preserved when the list is sorted:
//
//	Items.__present             -> "1"     (the list editor was initialized)
//	Items[__index].ID                      (hidden PK; empty for a new row)
//	Items[__index].Label                   (the real field)
//	Items[__index].__index      -> __index (stable per-item id, seeded at init)
//	Items[__index].__pos        -> n       (display / submit order)
//	Items[__index].__deleted    -> "true"  (client-side removal toggle)
//	Items[__index].__new        -> "true"  (browser-created row)
//
// Ordering vs identity:
//   - __pos drives the display and persistence order; the server reorders the
//     decoded slice by __pos (ReorderSlicesByPos) and the helper
//     presets.PartitionListEditorItems returns the items in __pos order.
//   - __index is the stable identity: the form key never moves when a row is
//     sorted, so a removed row cannot leak its metadata onto whatever later
//     occupies its old array position. New rows get a fresh __index.
//
// Classification for persistence uses the flags ONLY (__deleted / __new), never
// the primary key — so the list editor works for element types without an ID.
//
// Removing a row is a pure client-side toggle (`form["Items[i].__deleted"]=true`,
// see ListEditorItemContext.DeleteExpr); it is only sent to the server on the
// next submit. On a submit that fails validation the server must re-render each
// removed row as deleted AND re-seed its __deleted flag, so the removal is not
// silently lost — that is exactly what this app lets the test assert against the
// produced HTML.
//
// # The flow under test (see the *_test.go files for the assertions)
//
//  1. Two persisted items; remove item #1 (client toggles __deleted); clear the
//     required Name; SAVE. Expected: the Name error renders AND item #1 stays
//     rendered as deleted, with its __deleted re-seeded by the component Setup.
//  2. Add a brand-new item (flagged __new, no ID); remove it; clear Name; SAVE.
//     Expected: the new item also stays rendered as deleted (its submitted
//     __deleted honoured unconditionally, because a new item only ever exists in
//     the posted form).
//  3. A successful save with item #1 deleted actually removes it from the
//     database while keeping the others; a __new row (no ID) is created.
//  4. Items submitted reordered by __pos keep their form key from their __index
//     (index_test.go), and the partition helper classifies by flags in __pos
//     order (partition_test.go).
package listeditor

import (
	"fmt"
	"sync/atomic"

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

// The model URI and the nested field/form key the browser posts under.
const (
	ProductURI   = "products"
	ItemsField   = "Items"
	itemsFormKey = ItemsField
)

var dbSeq int64

// newDB opens a fresh, isolated in-memory SQLite database with the schema
// migrated. A unique shared-cache name per call keeps the gorm connection pool
// pointing at a single in-memory database while isolating each test from the
// others.
func newDB() (*gorm.DB, error) {
	name := fmt.Sprintf("file:listeditor_%d?mode=memory&cache=shared", atomic.AddInt64(&dbSeq, 1))
	db, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&Product{}, &Item{}); err != nil {
		return nil, err
	}
	return db, nil
}

// newApp builds the preset app: a Product model whose Items has-many field is
// rendered by the list-editor table builder, with Name required and the
// has-many association persisted by gorm2op.SaveHasManyAssociation.
func newApp(db *gorm.DB) *presets.Builder {
	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.Permission(perm.New().AllowAll())
	p.DataOperator(gorm2op.DataOperator(db))

	mb := p.Model(&Product{}).URIName(ProductURI)

	// The child model + its editing fields (the columns the table renders).
	itemMB := presets.NewModelBuilder(mb.Builder(), &Item{})
	ied := itemMB.Editing(&presets.FieldsSection{
		Rows: [][]string{{"Label"}},
	})
	ied.HiddenField("ID")

	ed := mb.Editing("Name", ItemsField)
	ed.Field("Name").SetRequired(true)

	// Persist the has-many association (deleted children are removed) and make
	// the parent Fetcher preload Items so opening the edit form shows them.
	mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		b := do.(*gorm2op.DataOperatorBuilder).
			SetFinder(func(qdb *gorm.DB, obj interface{}, ctx *web.EventContext) (any, error) {
				return obj, qdb.Preload("Items").First(obj).Error
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
