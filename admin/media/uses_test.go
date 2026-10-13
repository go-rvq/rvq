package media

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"

	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// usesPost is a record with media boxes: a box, a pointer to one.
type usesPost struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title string
	Cover media_library.MediaBox  `gorm:"type:text"`
	Thumb *media_library.MediaBox `gorm:"type:text"`
}

// usesDB is a database of the test's own — shared by its connections, its
// foreign keys on —, the uses tracked, with two files.
func usesTestDB(t *testing.T) (*gorm.DB, media_library.MediaLibrary, media_library.MediaLibrary) {
	t.Helper()
	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared&_foreign_keys=on"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	uuidkey.MustRegister(db)
	if err := db.AutoMigrate(&media_library.MediaLibrary{}, &MediaLibraryUse{}, &usesPost{}); err != nil {
		t.Fatal(err)
	}
	if err := TrackUses(db); err != nil {
		t.Fatal(err)
	}
	var a, b media_library.MediaLibrary
	if err := db.Create(&a).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&b).Error; err != nil {
		t.Fatal(err)
	}
	return db, a, b
}

func boxOf(m media_library.MediaLibrary) media_library.MediaBox {
	return media_library.MediaBox{ID: media_library.MediaID(m.ID.String()), Url: "/system/x.png"}
}

// usesText is the uses of db as "table:field=file" (the file as a or b).
func usesText(t *testing.T, db *gorm.DB, a, b media_library.MediaLibrary) string {
	t.Helper()
	var uses []MediaLibraryUse
	if err := db.Order("record_table, record_id, field").Find(&uses).Error; err != nil {
		t.Fatal(err)
	}
	name := map[uuid.UUID]string{a.ID: "a", b.ID: "b"}
	var out []string
	for _, u := range uses {
		out = append(out, u.Field+"="+name[u.MediaLibraryID])
	}
	// the records' keys are random: in no order of theirs
	sort.Strings(out)
	return strings.Join(out, ";")
}

// A use is the file a box of a record shows: set as the record is made,
// changed as its box is, gone as its box is emptied or the record deleted;
// a write that does not touch the boxes leaves them.
func TestUsesTracked(t *testing.T) {
	db, a, b := usesTestDB(t)
	thumb := boxOf(b)
	p := usesPost{ID: uuid.New(), Title: "one", Cover: boxOf(a), Thumb: &thumb}
	if err := db.Create(&p).Error; err != nil {
		t.Fatal(err)
	}
	if got := usesText(t, db, a, b); got != "Cover=a;Thumb=b" {
		t.Fatalf("made: %q", got)
	}
	// a save with the cover changed
	p.Cover = boxOf(b)
	if err := db.Save(&p).Error; err != nil {
		t.Fatal(err)
	}
	if got := usesText(t, db, a, b); got != "Cover=b;Thumb=b" {
		t.Errorf("the cover changed: %q", got)
	}
	// a write of other columns: the uses as they were
	if err := db.Model(&p).Update("title", "two").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&p).Updates(usesPost{Title: "three"}).Error; err != nil {
		t.Fatal(err)
	}
	if got := usesText(t, db, a, b); got != "Cover=b;Thumb=b" {
		t.Errorf("another column written: %q", got)
	}
	// the box by a map of columns, then by a struct
	if err := db.Model(&p).Updates(map[string]any{"cover": boxOf(a)}).Error; err != nil {
		t.Fatal(err)
	}
	if got := usesText(t, db, a, b); got != "Cover=a;Thumb=b" {
		t.Errorf("Updates(map): %q", got)
	}
	thumbA := boxOf(a)
	if err := db.Model(&p).Updates(usesPost{Thumb: &thumbA}).Error; err != nil {
		t.Fatal(err)
	}
	if got := usesText(t, db, a, b); got != "Cover=a;Thumb=a" {
		t.Errorf("Updates(struct): %q", got)
	}
	// emptied
	p.Cover, p.Thumb = media_library.MediaBox{}, nil
	if err := db.Save(&p).Error; err != nil {
		t.Fatal(err)
	}
	if got := usesText(t, db, a, b); got != "" {
		t.Errorf("emptied: %q", got)
	}
	// deleted
	p.Cover = boxOf(a)
	db.Save(&p)
	if err := db.Delete(&p).Error; err != nil {
		t.Fatal(err)
	}
	if got := usesText(t, db, a, b); got != "" {
		t.Errorf("deleted: %q", got)
	}
}

// A box of a key that names no file (of before the keys were UUIDs, a file
// purged) has no use — and the record is saved.
func TestUsesOfNoFile(t *testing.T) {
	db, a, b := usesTestDB(t)
	for _, id := range []media_library.MediaID{"42", media_library.MediaID(uuid.New().String())} {
		p := usesPost{ID: uuid.New(), Cover: media_library.MediaBox{ID: id}}
		if err := db.Create(&p).Error; err != nil {
			t.Fatalf("%s: %v", id, err)
		}
	}
	if got := usesText(t, db, a, b); got != "" {
		t.Errorf("uses of no file: %q", got)
	}
}

// The uses of a file, read; the foreign key keeps a file in use from going.
func TestUsesOfAndRestrict(t *testing.T) {
	db, a, b := usesTestDB(t)
	p := usesPost{ID: uuid.New(), Cover: boxOf(a)}
	db.Create(&p)
	uses, err := UsesOf(db, a.ID)
	if err != nil || len(uses) != 1 || uses[0].RecordID != p.ID.String() || uses[0].Field != "Cover" {
		t.Fatalf("UsesOf: %+v %v", uses, err)
	}
	if uses[0].String() != "uses_posts "+p.ID.String()+" (Cover)" {
		t.Errorf("String: %q", uses[0].String())
	}
	if err := db.Unscoped().Delete(&media_library.MediaLibrary{}, "id = ?", a.ID).Error; err == nil {
		t.Error("a file in use deleted")
	}
	if err := db.Unscoped().Delete(&media_library.MediaLibrary{}, "id = ?", b.ID).Error; err != nil {
		t.Errorf("a file of no use: %v", err)
	}
	msg := fileInUseMessage(Messages_en_US, uses)
	if !strings.Contains(msg, "(1)") || !strings.Contains(msg, p.ID.String()) {
		t.Errorf("the message: %q", msg)
	}
}

// The uses made again from the records: what the tracking did not see
// (written before it, a delete by a condition).
func TestBackfillUses(t *testing.T) {
	db, a, b := usesTestDB(t)
	p1 := usesPost{ID: uuid.New(), Cover: boxOf(a)}
	p2 := usesPost{ID: uuid.New(), Cover: boxOf(b)}
	// written with no hooks: no uses
	db.Session(&gorm.Session{SkipHooks: true}).Create(&p1)
	db.Session(&gorm.Session{SkipHooks: true}).Create(&p2)
	// a stale use of a record gone
	db.Create(&MediaLibraryUse{RecordTable: "uses_posts", RecordID: uuid.NewString(), Field: "Cover", MediaLibraryID: a.ID})
	n, err := BackfillUses(db, &usesPost{}, &media_library.MediaLibrary{})
	if err != nil || n != 2 {
		t.Fatalf("BackfillUses: %d %v", n, err)
	}
	if got := usesText(t, db, a, b); got != "Cover=a;Cover=b" {
		t.Errorf("after the backfill: %q", got)
	}
}

// The admin deletes no file in use — it says where it is used —; a file of no
// use, it deletes.
func TestDeleteFileInUseRefused(t *testing.T) {
	db, a, b := usesTestDB(t)
	p := usesPost{ID: uuid.New(), Cover: boxOf(a)}
	if err := db.Create(&p).Error; err != nil {
		t.Fatal(err)
	}
	pb := presets.New(i18n.New()).URIPrefix("/admin")
	pb.DataOperator(gorm2op.DataOperator(db))
	pb.Permission(perm.New().AllowAll().SubjectsFunc(func(*http.Request) []string { return []string{"x"} }))
	if err := New(db).Install(pb); err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	pb.Build(mux)
	del := func(id uuid.UUID) string {
		req := multipartestutils.NewMultipartBuilder().PageURL("/admin/media-library").EventFunc(doDeleteEvent).
			Query("field", "Cover").Query("id", id.String()).AddField("cfg", "{}").BuildEventFuncRequest()
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		return w.Body.String()
	}
	body := del(a.ID)
	if db.First(&media_library.MediaLibrary{}, "id = ?", a.ID).Error != nil {
		t.Fatalf("a file in use deleted: %.300s", body)
	}
	if !strings.Contains(body, "in use (1)") || !strings.Contains(body, p.ID.String()) {
		t.Errorf("no reason, nor where: %.500s", body)
	}
	del(b.ID)
	if db.First(&media_library.MediaLibrary{}, "id = ?", b.ID).Error == nil {
		t.Error("a file of no use not deleted")
	}
}
