package l10n

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// A localized record that keeps its deleted rows.
type softLocaleRecord struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	DeletedAt gorm.DeletedAt
	Locale
}

func (r *softLocaleRecord) PrimarySlug() string {
	return fmt.Sprintf("%d_%s", r.ID, r.LocaleCode)
}

func (r *softLocaleRecord) PrimaryColumnValuesBySlug(slug string) map[string]string {
	id, code, _ := strings.Cut(slug, "_")
	return map[string]string{"ID": id, "LocaleCode": code}
}

func softLocaleApp(t *testing.T) (*gorm.DB, *presets.ModelBuilder) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&softLocaleRecord{}); err != nil {
		t.Fatal(err)
	}
	pb := presets.New(i18n.New()).DataOperator(gorm2op.DataOperator(db))
	m := pb.Model(&softLocaleRecord{})

	b := New(db).DefaultLocaleCode("en-US")
	b.RegisterLocale("en-US", "en-us", "English")
	b.RegisterLocale("pt-BR", "pt-br", "Português do Brasil")
	if err := b.ModelInstall(pb, m); err != nil {
		t.Fatal(err)
	}
	return db, m
}

// Deleting keeps the row's locale: a locale_code that is no locale broke every
// foreign key to the locales (and every key that includes it).
func TestDeleteKeepsTheLocale(t *testing.T) {
	db, m := softLocaleApp(t)
	rec := &softLocaleRecord{ID: 1, Title: "x", Locale: Locale{LocaleCode: "pt-BR"}}
	if err := db.Create(rec).Error; err != nil {
		t.Fatal(err)
	}

	id, err := m.ParseRecordID("1_pt-BR")
	if err != nil {
		t.Fatal(err)
	}
	ctx := &web.EventContext{R: httptest.NewRequest("POST", "/", nil)}
	if err := m.Listing().Deleter(&softLocaleRecord{}, id, false, ctx); err != nil {
		t.Fatalf("delete: %v", err)
	}

	var got softLocaleRecord
	if err := db.Unscoped().Where("id = 1").First(&got).Error; err != nil {
		t.Fatal(err)
	}
	if got.LocaleCode != "pt-BR" {
		t.Errorf("locale_code = %q, want it kept: pt-BR", got.LocaleCode)
	}
	if !got.DeletedAt.Valid {
		t.Error("the record was not deleted (softly)")
	}
}

// Localizing to a locale whose row was deleted softly frees its key first;
// the other locales, and a live row, are left alone.
func TestPurgeSoftDeletedLocale(t *testing.T) {
	db, m := softLocaleApp(t)
	db.Create(&softLocaleRecord{ID: 1, Title: "en", Locale: Locale{LocaleCode: "en-US"}})
	db.Create(&softLocaleRecord{ID: 1, Title: "pt", Locale: Locale{LocaleCode: "pt-BR"}})
	db.Where("id = 1 AND locale_code = ?", "pt-BR").Delete(&softLocaleRecord{})

	id, _ := m.ParseRecordID("1_en-US")
	count := func(locale string) (n int64) {
		db.Unscoped().Model(&softLocaleRecord{}).Where("id = 1 AND locale_code = ?", locale).Count(&n)
		return
	}

	if err := purgeSoftDeletedLocale(db, m, id, "pt-BR"); err != nil {
		t.Fatal(err)
	}
	if n := count("pt-BR"); n != 0 {
		t.Errorf("the deleted pt-BR row is still there")
	}
	// a live row is never purged
	if err := purgeSoftDeletedLocale(db, m, id, "en-US"); err != nil {
		t.Fatal(err)
	}
	if n := count("en-US"); n != 1 {
		t.Errorf("the live en-US row was purged")
	}
}
