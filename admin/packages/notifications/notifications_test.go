package notifications

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	return db
}

// Told to each user once; unread counted, the newest first; read — one's
// own only —, all read.
func TestNotifications(t *testing.T) {
	db := testDB(t)
	ana, bia := uuid.New(), uuid.New()
	if err := Notify(db, Notification{Title: "A", Level: LevelWarning, Source: "x"}, ana, bia, ana, uuid.Nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Millisecond)
	if err := Notify(db, Notification{Title: "B"}, ana); err != nil {
		t.Fatal(err)
	}
	if n, _ := Unread(db, ana); n != 2 {
		t.Fatalf("ana's unread: %d", n)
	}
	if n, _ := Unread(db, bia); n != 1 {
		t.Fatalf("bia's unread: %d", n)
	}
	ns, _ := Latest(db, ana, 10)
	if len(ns) != 2 || ns[0].Title != "B" || ns[1].Level != LevelWarning || ns[0].Level != LevelInfo {
		t.Fatalf("ana's latest: %+v", ns)
	}
	bs, _ := Latest(db, bia, 10)
	// bia's id marked by ana: not read
	if err := MarkRead(db, ana, bs[0].ID); err != nil {
		t.Fatal(err)
	}
	if n, _ := Unread(db, bia); n != 1 {
		t.Error("another's read")
	}
	if err := MarkRead(db, ana, ns[0].ID); err != nil {
		t.Fatal(err)
	}
	if n, _ := Unread(db, ana); n != 1 {
		t.Errorf("one read: %d", n)
	}
	if err := MarkAllRead(db, ana); err != nil {
		t.Fatal(err)
	}
	if n, _ := Unread(db, ana); n != 0 {
		t.Errorf("all read: %d", n)
	}
}

// Every word in each language.
func TestMessagesFilled(t *testing.T) {
	for _, m := range []*Messages{Messages_en_US, Messages_pt_BR} {
		if m.BellTitle == "" || m.MarkAllRead == "" || m.SeeAll == "" || m.Notifications == "" || m.Open == "" {
			t.Errorf("%+v", m)
		}
	}
}
