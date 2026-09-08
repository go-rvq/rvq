package db_tools

import (
	"os"
	"sync"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Backup is scanned into by a raw query, and Scan parses the destination as a
// gorm schema. A func field without `gorm:"-"` fails that parse with
// "unsupported data type", which is how the nightly backup broke.
func TestBackupParsesAsGormSchema(t *testing.T) {
	s, err := schema.Parse(&Backup{}, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatalf("gorm cannot parse Backup: %v", err)
	}

	for _, name := range []string{"OpenFunc", "DetailFunc"} {
		if f := s.LookUpField(name); f != nil && f.DBName != "" {
			t.Errorf("%s was mapped to column %q; it is behaviour, not data", name, f.DBName)
		}
	}

	// The columns the query actually fills must still be there.
	for _, name := range []string{"CreatedAt", "DbName", "Message", "Size", "Auto"} {
		if f := s.LookUpField(name); f == nil || f.DBName == "" {
			t.Errorf("%s is not mapped to a column", name)
		}
	}
}

// The failure came from Scan, which parses the destination before it copies
// anything. This reproduces that call shape against a real driver.
func TestBackupIsScannable(t *testing.T) {
	dsn := os.Getenv("DB_TOOLS_TEST_DSN")
	if dsn == "" {
		t.Skip("set DB_TOOLS_TEST_DSN to run against a database")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	var bkp Backup
	err = db.Raw(`SELECT now() AS created_at, 'app' AS db_name,
	                     'teste' AS message, 42::bigint AS size, true AS auto`).
		Scan(&bkp).Error
	if err != nil {
		t.Fatalf("Scan into Backup failed: %v", err)
	}

	if bkp.DbName != "app" || bkp.Message != "teste" || bkp.Size != 42 || !bkp.Auto {
		t.Errorf("scanned %+v", bkp)
	}
}
