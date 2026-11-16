package db_tools

import (
	"testing"

	"github.com/go-rvq/rvq/x/i18n"
	db_tools "github.com/go-rvq/rvq/x/packages/db-tools"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// A configuration without a retention — none yet, or empty — is given the
// default (one a day, for 30 days); one with a retention keeps it; 0, none
// given.
func TestEnsureRetention(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	b := New(db, i18n.New(), nil)
	if err := b.ensureRetention(); err != nil {
		t.Fatal(err)
	}
	per, _ := b.GetPersistence()
	if !per.Enabled || per.Days != 30 {
		t.Fatalf("none yet: %+v", per)
	}
	cfg, _ := b.GetConfig()
	cfg.Persistence.Data = &db_tools.Persistence{Enabled: true, Weeks: 2}
	db.Save(cfg)
	if err := b.ensureRetention(); err != nil {
		t.Fatal(err)
	}
	if per, _ := b.GetPersistence(); per.Weeks != 2 || per.Days != 0 {
		t.Errorf("one set: %+v", per)
	}
	cfg.Persistence.Data = &db_tools.Persistence{}
	db.Save(cfg)
	b.DefaultRetention(7)
	b.ensureRetention()
	if per, _ := b.GetPersistence(); per.Days != 7 {
		t.Errorf("empty: %+v", per)
	}
	var n int64
	db.Model(&DbBackupConfig{}).Count(&n)
	if n != 1 {
		t.Errorf("%d configurations", n)
	}
}
