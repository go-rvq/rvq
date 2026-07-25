package integration

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/packages/helper"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// A fixture of NESTED MODELS — has-many fields taken as child models with their
// own listing/detailing/editing, two levels deep:
//
//	Survey → Survey.Places → Survey.Places.Products
//
// This is the shape a real application reaches with
// helper.NewNestedSliceBuilder: browsing the root opens the Places listing in a
// dialog, opening a place there opens ITS detail (with the Products listing) in
// another dialog — three listings and two detailings alive at the same time.
//
// It is the stress case for form hosts: every level renders the same kind of
// host, so any state they shared would make one level's button drive another
// level's overlay.

type NMProduct struct {
	ID      uint
	PlaceID uint
	Name    string
	Price   int
}

type NMPlace struct {
	ID       uint
	SurveyID uint
	Name     string
	Products []*NMProduct `gorm:"foreignKey:PlaceID"`
}

type NMSurvey struct {
	ID     uint
	Name   string
	Places []*NMPlace `gorm:"foreignKey:SurveyID"`
}

const NestedModelsRootURI = "surveys"

// NewNestedModelsDB migrates the schema on a fresh in-memory SQLite database.
func NewNestedModelsDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(nextMemName()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&NMSurvey{}, &NMPlace{}, &NMProduct{}); err != nil {
		return nil, err
	}
	return db, nil
}

// NewNestedModelsApp builds the two-level nested-model app: each level is taken
// as a child model with NewNestedSliceBuilder, exactly as an application nests a
// has-many field it wants to browse instead of edit inline.
func NewNestedModelsApp(db *gorm.DB) *presets.Builder {
	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.Permission(perm.New().AllowAll())
	p.DataOperator(gorm2op.DataOperator(db))

	mb := p.Model(&NMSurvey{}).URIName(NestedModelsRootURI)
	mb.Listing("ID", "Name")
	mb.Detailing("Name", "Places")
	mb.Editing("Name")

	// level 1: Survey.Places
	helper.NewNestedSliceBuilder(mb, "Places").
		SetCallback(func(places *presets.ModelBuilder) {
			places.Listing("ID", "Name")
			places.Detailing("Name", "Products")
			places.Editing("Name")

			// level 2: Survey.Places.Products
			helper.NewNestedSliceBuilder(places, "Products").
				SetCallback(func(products *presets.ModelBuilder) {
					products.Listing("ID", "Name", "Price")
					products.Detailing("Name", "Price")
					products.Editing("Name", "Price")
				}).
				Build()
		}).
		Build()

	return p
}

// NewNestedModelsSeededHandler builds NewNestedModelsApp seeded with one survey
// holding two places, each holding two products.
func NewNestedModelsSeededHandler() (http.Handler, error) {
	db, err := NewNestedModelsDB()
	if err != nil {
		return nil, err
	}
	survey := &NMSurvey{ID: 1, Name: "S1", Places: []*NMPlace{
		{ID: 1, SurveyID: 1, Name: "Place A", Products: []*NMProduct{
			{ID: 1, PlaceID: 1, Name: "A-one", Price: 10},
			{ID: 2, PlaceID: 1, Name: "A-two", Price: 20},
		}},
		{ID: 2, SurveyID: 1, Name: "Place B", Products: []*NMProduct{
			{ID: 3, PlaceID: 2, Name: "B-one", Price: 30},
		}},
	}}
	if err = db.Create(survey).Error; err != nil {
		return nil, err
	}
	return NewNestedModelsApp(db), nil
}
