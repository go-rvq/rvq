package media

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// MediaLibraryUse is a use of a file of the media library: the field of a
// record that shows it (a media_library.MediaBox) — its table, the record's
// key, the field —, and the file, by a foreign key that keeps a file in use
// from being deleted (ON DELETE RESTRICT; the admin refuses a deletion of
// one too, soft or not: deleteFileField).
//
// The uses are kept as the records are written (TrackUses: gorm's callbacks
// of create, update and delete): a box set, changed or emptied is its use
// set, changed or gone; a record deleted, its uses gone. What gorm writes with
// no record to read — Delete(&T{}, "condition"), a raw query — is out of it:
// BackfillUses makes them again from the records.
type MediaLibraryUse struct {
	RecordTable    string                      `gorm:"primaryKey;size:128"`
	RecordID       string                      `gorm:"primaryKey;size:255"`
	Field          string                      `gorm:"primaryKey;size:128"`
	MediaLibraryID uuid.UUID                   `gorm:"type:uuid;not null;index"`
	MediaLibrary   *media_library.MediaLibrary `gorm:"foreignKey:MediaLibraryID;constraint:OnDelete:RESTRICT"`
	UpdatedAt      time.Time
}

func (MediaLibraryUse) TableName() string { return "media_library_uses" }

// String is the use as a person reads it: "posts 6f1c… (Cover)".
func (u MediaLibraryUse) String() string {
	return u.RecordTable + " " + u.RecordID + " (" + u.Field + ")"
}

var mediaBoxType = reflect.TypeOf(media_library.MediaBox{})

// boxFields are the fields of s that hold a MediaBox (or a pointer to one).
var boxFieldsCache sync.Map // *schema.Schema → []*schema.Field

func boxFields(s *schema.Schema) []*schema.Field {
	if v, ok := boxFieldsCache.Load(s); ok {
		return v.([]*schema.Field)
	}
	var out []*schema.Field
	for _, f := range s.Fields {
		t := f.FieldType
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t == mediaBoxType && f.DBName != "" {
			out = append(out, f)
		}
	}
	boxFieldsCache.Store(s, out)
	return out
}

var tracked sync.Map // *gorm.Config → true

// TrackUses keeps the uses of the files of the media library (MediaLibraryUse)
// as db writes records: after a create, an update (a save, Updates, Update)
// and a delete. Once per db (its gorm.Config); the table must exist.
func TrackUses(db *gorm.DB) error {
	if _, done := tracked.LoadOrStore(db.Config, true); done {
		return nil
	}
	cb := db.Callback()
	if err := cb.Create().After("gorm:create").Register("media:uses", afterWrite); err != nil {
		return err
	}
	if err := cb.Update().After("gorm:update").Register("media:uses", afterWrite); err != nil {
		return err
	}
	return cb.Delete().After("gorm:delete").Register("media:uses", afterDelete)
}

// statementOf is the statement of tx when it is one to look at: no error, a
// model of a schema with media boxes (not the uses themselves).
func statementOf(tx *gorm.DB) (*gorm.Statement, []*schema.Field) {
	stmt := tx.Statement
	if tx.Error != nil || stmt == nil || stmt.Schema == nil || stmt.SkipHooks {
		return nil, nil
	}
	if stmt.Schema.Table == (MediaLibraryUse{}).TableName() {
		return nil, nil
	}
	fields := boxFields(stmt.Schema)
	if len(fields) == 0 {
		return nil, nil
	}
	return stmt, fields
}

// records are the records of the statement's model: a struct, or a slice or
// array of them (a create of many).
func records(rv reflect.Value) []reflect.Value {
	rv = reflect.Indirect(rv)
	switch rv.Kind() {
	case reflect.Struct:
		return []reflect.Value{rv}
	case reflect.Slice, reflect.Array:
		out := make([]reflect.Value, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			if r := reflect.Indirect(rv.Index(i)); r.Kind() == reflect.Struct {
				out = append(out, r)
			}
		}
		return out
	}
	return nil
}

// recordID is the key of the record rv as text — the values of a composite
// one joined by ","; "" when it has none yet.
func recordID(ctx context.Context, s *schema.Schema, rv reflect.Value) string {
	if len(s.PrimaryFields) == 0 {
		return ""
	}
	parts := make([]string, 0, len(s.PrimaryFields))
	for _, pf := range s.PrimaryFields {
		v, zero := pf.ValueOf(ctx, rv)
		if zero {
			return ""
		}
		parts = append(parts, fmt.Sprint(v))
	}
	return strings.Join(parts, ",")
}

// selected says the statement writes the field f: not omitted, and selected
// when it selects some.
func selected(stmt *gorm.Statement, f *schema.Field) bool {
	named := func(list []string) bool {
		return slices.Contains(list, f.Name) || slices.Contains(list, f.DBName)
	}
	if named(stmt.Omits) {
		return false
	}
	return len(stmt.Selects) == 0 || slices.Contains(stmt.Selects, "*") || named(stmt.Selects)
}

// asBox is v — a value of a MediaBox field, or of a map of columns — as a box.
func asBox(v any) (media_library.MediaBox, bool) {
	switch b := v.(type) {
	case media_library.MediaBox:
		return b, true
	case *media_library.MediaBox:
		if b == nil {
			return media_library.MediaBox{}, true
		}
		return *b, true
	case nil:
		return media_library.MediaBox{}, true
	}
	return media_library.MediaBox{}, false
}

// afterWrite keeps the uses of what a create or an update wrote: of a map of
// columns (Updates(map)), the boxes in it; of a struct other than the model
// (Updates(struct)), its boxes not empty (gorm writes no zero field from
// one); of the model itself (Create, Save), every box it writes — an empty
// one its use gone.
func afterWrite(tx *gorm.DB) {
	stmt, fields := statementOf(tx)
	if stmt == nil {
		return
	}
	ctx := stmt.Context
	recs := records(stmt.ReflectValue)
	switch dest := stmt.Dest.(type) {
	case map[string]any:
		for _, f := range fields {
			v, ok := dest[f.DBName]
			if !ok {
				v, ok = dest[f.Name]
			}
			if !ok || !selected(stmt, f) {
				continue
			}
			if box, ok := asBox(v); ok {
				for _, rec := range recs {
					setUse(tx, stmt.Schema.Table, recordID(ctx, stmt.Schema, rec), f.Name, box)
				}
			}
		}
		return
	}
	other := stmt.Dest != nil && stmt.Model != nil && !samePointer(stmt.Dest, stmt.Model)
	var destRec reflect.Value
	if other {
		dr := records(reflect.ValueOf(stmt.Dest))
		if len(dr) != 1 {
			return
		}
		destRec = dr[0]
	}
	for _, rec := range recs {
		id := recordID(ctx, stmt.Schema, rec)
		if id == "" {
			continue
		}
		for _, f := range fields {
			if !selected(stmt, f) {
				continue
			}
			from := rec
			if other {
				from = destRec
			}
			v, zero := f.ValueOf(ctx, from)
			if other && zero && len(stmt.Selects) == 0 {
				continue // Updates(struct): a zero field not written
			}
			if box, ok := asBox(v); ok {
				setUse(tx, stmt.Schema.Table, id, f.Name, box)
			}
		}
	}
}

// afterDelete forgets the uses of the records deleted (soft deleted too).
func afterDelete(tx *gorm.DB) {
	stmt, _ := statementOf(tx)
	if stmt == nil {
		return
	}
	for _, rec := range records(stmt.ReflectValue) {
		if id := recordID(stmt.Context, stmt.Schema, rec); id != "" {
			if err := usesDB(tx).Where("record_table = ? AND record_id = ?", stmt.Schema.Table, id).
				Delete(&MediaLibraryUse{}).Error; err != nil {
				_ = tx.AddError(err)
			}
		}
	}
}

func samePointer(a, b any) bool {
	va, vb := reflect.ValueOf(a), reflect.ValueOf(b)
	return va.Kind() == reflect.Ptr && vb.Kind() == reflect.Ptr && va.Pointer() == vb.Pointer()
}

// usesDB is tx — its transaction — for the uses: no hooks (no use of a use).
func usesDB(tx *gorm.DB) *gorm.DB {
	return tx.Session(&gorm.Session{NewDB: true, SkipHooks: true})
}

// setUse is the use of the field of the record: the box's file, or none — an
// empty box, a key that names no file (a box of before the keys were UUIDs, a
// file purged) —. It reports whether a use was written.
func setUse(tx *gorm.DB, table, recordID, field string, box media_library.MediaBox) bool {
	if recordID == "" {
		return false
	}
	db := usesDB(tx)
	where := db.Where("record_table = ? AND record_id = ? AND field = ?", table, recordID, field)
	file, err := uuid.Parse(string(box.ID))
	if box.ID.IsZero() || err != nil || !fileExists(db, file) {
		if err := where.Delete(&MediaLibraryUse{}).Error; err != nil {
			_ = tx.AddError(err)
		}
		return false
	}
	err = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "record_table"}, {Name: "record_id"}, {Name: "field"}},
		DoUpdates: clause.AssignmentColumns([]string{"media_library_id", "updated_at"}),
	}).Create(&MediaLibraryUse{RecordTable: table, RecordID: recordID, Field: field,
		MediaLibraryID: file, UpdatedAt: time.Now()}).Error
	if err != nil {
		_ = tx.AddError(err)
		return false
	}
	return true
}

// fileExists says the media library has the file (deleted softly too: the
// key is still there for the foreign key).
func fileExists(db *gorm.DB, id uuid.UUID) bool {
	var n int64
	db.Model(&media_library.MediaLibrary{}).Unscoped().Where("id = ?", id).Count(&n)
	return n > 0
}

// UsesOf are the uses of the file id, by table, record and field.
func UsesOf(db *gorm.DB, id uuid.UUID) ([]MediaLibraryUse, error) {
	var uses []MediaLibraryUse
	err := db.Where("media_library_id = ?", id).
		Order("record_table, record_id, field").Find(&uses).Error
	return uses, err
}

// BackfillUses makes the uses again from the records of models — what
// TrackUses does not see (a delete by a condition, a raw query, what was
// written before it) —: of each model with media boxes, each record's. It
// answers how many uses it wrote.
func BackfillUses(db *gorm.DB, models ...any) (n int, err error) {
	for _, m := range models {
		stmt := &gorm.Statement{DB: db}
		if err = stmt.Parse(m); err != nil {
			return n, fmt.Errorf("media: the uses of %T: %w", m, err)
		}
		fields := boxFields(stmt.Schema)
		if len(fields) == 0 {
			continue
		}
		table := stmt.Schema.Table
		// what the table knew of the uses is made again: what is gone, gone
		if err = usesDB(db).Where("record_table = ?", table).Delete(&MediaLibraryUse{}).Error; err != nil {
			return n, err
		}
		typ := reflect.Indirect(reflect.ValueOf(m)).Type()
		batch := reflect.New(reflect.SliceOf(typ))
		res := db.Model(m).FindInBatches(batch.Interface(), 500, func(tx *gorm.DB, _ int) error {
			for _, rec := range records(batch) {
				id := recordID(db.Statement.Context, stmt.Schema, rec)
				if id == "" {
					continue
				}
				for _, f := range fields {
					v, _ := f.ValueOf(context.Background(), rec)
					box, ok := asBox(v)
					if !ok || box.ID.IsZero() {
						continue
					}
					if setUse(db, table, id, f.Name, box) {
						n++
					}
				}
			}
			return nil
		})
		if res.Error != nil {
			return n, fmt.Errorf("media: the uses of %s: %w", table, res.Error)
		}
	}
	return n, nil
}

// fileInUseMessage says a file is in use, and where: its first uses.
func fileInUseMessage(m *Messages, uses []MediaLibraryUse) string {
	const shown = 5
	var b strings.Builder
	fmt.Fprintf(&b, m.FileInUse, len(uses))
	for i, u := range uses {
		if i == shown {
			b.WriteString(" …")
			break
		}
		b.WriteString(" ")
		b.WriteString(u.String())
		if i < len(uses)-1 && i < shown-1 {
			b.WriteString(";")
		}
	}
	return b.String()
}
