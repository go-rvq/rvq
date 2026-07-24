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

// A multi-level nested list editor over records with COMPOSITE primary keys —
// the association-object flavor of many-to-many, where the join rows carry
// editable attributes:
//
//	Cart (ID)
//	 └─ Items[]   CartItem   PK (CartID, Sku)          — required Qty
//	     └─ Notes[] CartNote PK (CartID, Sku, Seq)     — required Text
//
// Every level's has-many reconciliation (create/update/delete) keys by the FULL
// composite primary key (gorm2op via gormutils.PKMapKey), and the list editor
// submits every PK column as a hidden field so a row is matched by its whole key.

type CartNote struct {
	CartID uint   `gorm:"primaryKey"`
	Sku    string `gorm:"primaryKey"`
	Seq    uint   `gorm:"primaryKey"`
	Text   string
	Pos    int
}

type CartItem struct {
	CartID uint   `gorm:"primaryKey"`
	Sku    string `gorm:"primaryKey"`
	Qty    int
	Pos    int
	Notes  []*CartNote `gorm:"foreignKey:CartID,Sku;references:CartID,Sku"`
}

type Cart struct {
	ID    uint
	Name  string
	Items []*CartItem `gorm:"foreignKey:CartID;references:ID"`
}

const CartURI = "carts"

// NewCompositeDB migrates the composite-key schema on a fresh in-memory database.
func NewCompositeDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(nextMemName()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&Cart{}, &CartItem{}, &CartNote{}); err != nil {
		return nil, err
	}
	return db, nil
}

// NewCompositeApp builds the nested composite-key list-editor preset app.
func NewCompositeApp(db *gorm.DB) *presets.Builder {
	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.Permission(perm.New().AllowAll())
	p.DataOperator(gorm2op.DataOperator(db))

	mb := p.Model(&Cart{}).URIName(CartURI)

	// leaf: CartNote (all three PK fields hidden so they round-trip; Text required)
	noteMB := presets.NewModelBuilder(mb.Builder(), &CartNote{})
	noteEd := noteMB.Editing(&presets.FieldsSection{Rows: [][]string{{"Text"}}})
	noteEd.HiddenField("CartID")
	noteEd.HiddenField("Sku")
	noteEd.HiddenField("Seq")
	noteEd.Field("Text").SetRequired(true)

	// CartItem (both PK fields hidden; Qty required; nests Notes)
	itemMB := presets.NewModelBuilder(mb.Builder(), &CartItem{})
	itemEd := itemMB.Editing(&presets.FieldsSection{Rows: [][]string{{"Qty"}, {"Notes"}}})
	itemEd.HiddenField("CartID")
	itemEd.HiddenField("Sku")
	itemEd.Field("Qty").SetRequired(true)
	itemEd.Field("Notes").Nested(presets.NestedSlice(noteMB, &noteEd.FieldsBuilder))

	// root Cart
	ed := mb.Editing("Name", "Items")
	ed.Field("Name").SetRequired(true)
	ed.Field("Items").Nested(presets.NestedSlice(itemMB, &itemEd.FieldsBuilder))

	// preload the whole tree on read + reconcile the Items has-many on save.
	mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		b := do.(*gorm2op.DataOperatorBuilder).
			WrapPrepare(func(old gorm2op.Preparer) gorm2op.Preparer {
				return func(db *gorm.DB, mode gorm2op.Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
					db = old(db, mode, obj, id, params, ctx)
					if mode.Is(gorm2op.Fetch, gorm2op.Search) {
						db = db.Preload("Items.Notes")
					}
					return db
				}
			})
		return gorm2op.SaveHasManyAssociation("Items").Build(b)
	})

	return p
}

// NewCompositeSeededHandler builds NewCompositeApp seeded with one cart:
// Cart "c1" → Item (Sku "A", Qty 2) → Note (Seq 1, "hello").
func NewCompositeSeededHandler() (http.Handler, error) {
	db, err := NewCompositeDB()
	if err != nil {
		return nil, err
	}
	cart := &Cart{ID: 1, Name: "c1", Items: []*CartItem{
		{CartID: 1, Sku: "A", Qty: 2, Pos: 0, Notes: []*CartNote{
			{CartID: 1, Sku: "A", Seq: 1, Text: "hello", Pos: 0},
		}},
	}}
	if err = db.Create(cart).Error; err != nil {
		return nil, err
	}
	return NewCompositeApp(db), nil
}
