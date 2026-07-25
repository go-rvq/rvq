package integration

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// A fixture for the POST-SAVE REFRESH rules: after a successful save only what
// reflects the change is refreshed, and never by reloading the page.
//
//	Article  listing + detailing + editing   — NEW, row DETAIL → EDIT
//	Tag      listing + editing, NO detailing — the row opens EDIT directly
//	Site     singleton + detailing + editing — EDIT must refresh the DETAIL
//	Pref     singleton + editing only        — there is only the EDIT

type RFArticle struct {
	ID    uint
	Title string
	Body  string
}

// the page title follows the Title field, so a save must be visible in <title>
func (a *RFArticle) String() string { return a.Title }

type RFTag struct {
	ID   uint
	Name string
}

type RFSite struct {
	ID    uint
	Name  string
	Motto string
}

type RFPref struct {
	ID    uint
	Theme string
}

const (
	RefreshArticleURI = "articles"
	RefreshTagURI     = "tags"
	RefreshSiteURI    = "site"
	RefreshPrefURI    = "pref"
)

// NewRefreshDB migrates the schema on a fresh in-memory SQLite database.
func NewRefreshDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(nextMemName()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&RFArticle{}, &RFTag{}, &RFSite{}, &RFPref{}); err != nil {
		return nil, err
	}
	return db, nil
}

// NewRefreshApp registers the four models above.
func NewRefreshApp(db *gorm.DB) *presets.Builder {
	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.Permission(perm.New().AllowAll())
	p.DataOperator(gorm2op.DataOperator(db))

	article := p.Model(&RFArticle{}).URIName(RefreshArticleURI)
	article.Listing("ID", "Title")
	article.Detailing("Title", "Body")
	article.Editing("Title", "Body")

	// no detailing: a row opens the edit form
	tag := p.Model(&RFTag{}).URIName(RefreshTagURI)
	tag.Listing("ID", "Name")
	tag.Editing("Name")

	site := p.Model(&RFSite{}).URIName(RefreshSiteURI).Singleton(true)
	site.Detailing("Name", "Motto")
	site.Editing("Name", "Motto")

	pref := p.Model(&RFPref{}).URIName(RefreshPrefURI).Singleton(true)
	pref.Editing("Theme")

	return p
}

// NewRefreshSeededHandler builds NewRefreshApp with one row per model.
func NewRefreshSeededHandler() (http.Handler, error) {
	db, err := NewRefreshDB()
	if err != nil {
		return nil, err
	}
	for _, rec := range []any{
		&RFArticle{ID: 1, Title: "A1", Body: "b1"},
		&RFArticle{ID: 2, Title: "A2", Body: "b2"},
		&RFTag{ID: 1, Name: "T1"},
		&RFSite{ID: 1, Name: "S", Motto: "m"},
		&RFPref{ID: 1, Theme: "dark"},
	} {
		if err = db.Create(rec).Error; err != nil {
			return nil, err
		}
	}
	return NewRefreshApp(db), nil
}
