package login_session

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/place"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbSeq int64

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:ls_%d?mode=memory&cache=shared", atomic.AddInt64(&dbSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&LoginSession{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func accessRequest(ip string) *http.Request {
	r := httptest.NewRequest("GET", "/admin/site-files.git/info/refs", nil)
	r.RemoteAddr = ip + ":5000"
	r.Header.Set("User-Agent", "git/2.43.0")
	return r
}

// The accesses of the WebDAV and the git: one row a device, IP and way in,
// in a window — its requests counted, written at most once a minute —, its
// place; another IP, another row; the old ones purged.
func TestRecordAccess(t *testing.T) {
	db := testDB(t)
	m := NewManager(login.New(i18n.New()))
	m.PlaceFunc = func(r *http.Request) place.Place { return place.Place{City: "Boston", Country: "US"} }
	user := uuid.New()

	for i := 0; i < 3; i++ {
		if err := m.RecordAccess(db, accessRequest("203.0.113.9"), Access{UserID: user, Kind: KindGit, Auth: login.AuthAccessKey}); err != nil {
			t.Fatal(err)
		}
	}
	var rows []LoginSession
	db.Find(&rows)
	if len(rows) != 1 || rows[0].Requests != 1 || rows[0].Place.String() != "Boston, US" || rows[0].Kind != KindGit ||
		rows[0].Auth != login.AuthAccessKey || rows[0].LastAccessAt == nil {
		t.Fatalf("one row, written once in the minute: %+v", rows)
	}
	// a minute later: the requests counted in between, written
	m.mu.Lock()
	for _, p := range m.accesses {
		p.written = p.written.Add(-2 * accessWriteEvery)
	}
	m.mu.Unlock()
	if err := m.RecordAccess(db, accessRequest("203.0.113.9"), Access{UserID: user, Kind: KindGit, Auth: login.AuthAccessKey}); err != nil {
		t.Fatal(err)
	}
	db.First(&rows[0], "id = ?", rows[0].ID)
	if rows[0].Requests != 4 {
		t.Errorf("the requests of the window: %d", rows[0].Requests)
	}
	// another process (a restart): the window found in the database
	m2 := NewManager(login.New(i18n.New()))
	if err := m2.RecordAccess(db, accessRequest("203.0.113.9"), Access{UserID: user, Kind: KindGit, Auth: login.AuthAccessKey}); err != nil {
		t.Fatal(err)
	}
	var n int64
	db.Model(&LoginSession{}).Count(&n)
	if n != 1 {
		t.Errorf("a restart made another row: %d", n)
	}
	// another IP, another way in: rows of their own
	_ = m.RecordAccess(db, accessRequest("198.51.100.7"), Access{UserID: user, Kind: KindGit, Auth: login.AuthAccessKey})
	_ = m.RecordAccess(db, accessRequest("203.0.113.9"), Access{UserID: user, Kind: KindWebDAV, Auth: login.AuthPassword})
	db.Model(&LoginSession{}).Count(&n)
	if n != 3 {
		t.Errorf("rows of their own: %d", n)
	}

	// purged: the rows older than the retention
	old := time.Now().Add(-100 * 24 * time.Hour)
	db.Model(&LoginSession{}).Where("ip = ?", "198.51.100.7").
		Updates(map[string]any{"last_access_at": old, "expired_at": old, "created_at": old})
	if purged, err := m.Purge(db, 90*24*time.Hour); err != nil || purged != 1 {
		t.Errorf("purged: %d %v", purged, err)
	}
	db.Model(&LoginSession{}).Count(&n)
	if n != 2 {
		t.Errorf("after the purge: %d", n)
	}
}
