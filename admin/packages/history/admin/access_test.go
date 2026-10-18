package admin

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRecorderAccessHours(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&AccessHour{}); err != nil {
		t.Fatal(err)
	}
	r := NewRecorder(db)
	for range 3 {
		r.Hit("posts_revisions", "a_en")
	}
	r.Hit("pages_revisions", "b_en")
	r.Flush()
	r.Hit("posts_revisions", "a_en")
	r.Flush()

	var rows []AccessHour
	if err := db.Order("table_name").Find(&rows).Error; err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].RecordKey != "b_en" || rows[0].N != 1 || rows[1].N != 4 {
		t.Fatalf("rows %+v", rows)
	}
	if rows[1].Hour.Minute() != 0 || rows[1].Hour.Second() != 0 {
		t.Errorf("hour not truncated: %s", rows[1].Hour)
	}
}
