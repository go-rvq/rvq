package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type LCArticle struct {
	ID    uint
	Name  string
	Cover string
	Intro string
}

func (a *LCArticle) String() string          { return a.Name }
func (a *LCArticle) LinkCardCover() string   { return a.Cover }
func (a *LCArticle) LinkCardSummary() string { return a.Intro }

type LCSecret struct {
	ID   uint
	Name string
}

func (a *LCSecret) String() string { return a.Name }

var lcSeq int64

// The cards of the links of a text (presets_LinkCards): a record — its
// title, its detail's address, its cover and summary —, a listing; the
// application's targets by its resolver; a record not there, a path that
// names nothing, or what the user may not see: no card.
func TestLinkCards(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:linkcards_%d?mode=memory&cache=shared", atomic.AddInt64(&lcSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&LCArticle{}, &LCSecret{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&LCArticle{ID: 1, Name: "Roofs", Cover: "/system/roofs.jpg", Intro: "All about roofs"})
	db.Create(&LCSecret{ID: 1, Name: "Hidden"})

	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.Permission(perm.New().Policies(
		perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything),
		perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(perm.Anything).On("*:lcsecrets:*"),
	).SubjectsFunc(func(*http.Request) []string { return []string{"editor"} }))
	p.DataOperator(gorm2op.DataOperator(db))
	arts := p.Model(&LCArticle{}).URIName("lc-articles")
	arts.Listing("ID", "Name")
	arts.Detailing("Name")
	p.Model(&LCSecret{}).URIName("lc-secrets").Listing("ID", "Name")
	p.LinkCardResolver(func(_ *web.EventContext, target string) *presets.LinkCard {
		if target == "/en-us/about" {
			return &presets.LinkCard{Href: target, Title: "About us"}
		}
		return nil
	})

	art := arts.Info().RecordPath(nil, arts.MustRecordID(&LCArticle{ID: 1}))
	missing := arts.Info().RecordPath(nil, arts.MustRecordID(&LCArticle{ID: 9}))
	listing := arts.Info().ListingPath(nil)
	secretMB := p.GetModel(&LCSecret{})
	secret := secretMB.Info().RecordPath(nil, secretMB.MustRecordID(&LCSecret{ID: 1}))

	q := url.Values{"target": {art, listing, "/en-us/about", missing, secret, "admin:nothing:", "/unknown"}}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/lc-articles?__execute_event__="+presets.LinkCardsEvent, strings.NewReader(q.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.ServeHTTP(w, r)
	var out struct {
		Data map[string]presets.LinkCard `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("%v: %s", err, w.Body.String())
	}
	if c := out.Data[art]; c.Title != "Roofs" || c.Href != arts.Info().DetailingHref("1") || c.Cover != "/system/roofs.jpg" || c.Summary != "All about roofs" {
		t.Errorf("the record's card: %+v", c)
	}
	if c := out.Data[listing]; c.Href != arts.Info().ListingHref() || c.Title == "" {
		t.Errorf("the listing's card: %+v", c)
	}
	if c := out.Data["/en-us/about"]; c.Title != "About us" {
		t.Errorf("the application's card: %+v", c)
	}
	for _, none := range []string{missing, secret, "admin:nothing:", "/unknown"} {
		if c, ok := out.Data[none]; ok {
			t.Errorf("%s: a card %+v", none, c)
		}
	}
}
