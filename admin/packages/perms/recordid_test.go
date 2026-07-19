package perms

import (
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// uuidRec exercises a uuid.UUID primary key (built-in support).
type uuidRec struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string
}

// code is a custom id type exposing Parse(string) (code, error).
type code string

func (code) Parse(s string) (code, error) { return code(s), nil }

type codeRec struct {
	ID   code `gorm:"primaryKey"`
	Name string
}

// fbID is a custom id type with a Valuer (so gorm treats it as a scalar
// column) but no Scanner and no Parse method, so it is parsed only through
// IdParserFallback.
type fbID struct{ V string }

func (f fbID) Value() (driver.Value, error) { return f.V, nil }

type fbRec struct {
	ID   fbID `gorm:"primaryKey;type:text"`
	Name string
}

func TestRecordIDParsing(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	b := presets.New(i18n.New()).DataOperator(gorm2op.DataOperator(db))

	t.Run("uuid", func(t *testing.T) {
		mb := b.Model(&uuidRec{})
		want := uuid.New()
		var obj uuidRec
		if err := mb.ParseRecordIDTo(&obj, want.String()); err != nil {
			t.Fatal(err)
		}
		if obj.ID != want {
			t.Errorf("got %v, want %v", obj.ID, want)
		}
	})

	t.Run("parse method", func(t *testing.T) {
		mb := b.Model(&codeRec{})
		var obj codeRec
		if err := mb.ParseRecordIDTo(&obj, "ABCD"); err != nil {
			t.Fatal(err)
		}
		if obj.ID != code("ABCD") {
			t.Errorf("got %q, want ABCD", obj.ID)
		}
	})

	t.Run("fallback register", func(t *testing.T) {
		presets.IdParserFallback[reflect.TypeOf(fbID{})] =
			func(s string) (any, error) { return fbID{V: s}, nil }
		defer delete(presets.IdParserFallback, reflect.TypeOf(fbID{}))

		mb := b.Model(&fbRec{})
		var obj fbRec
		if err := mb.ParseRecordIDTo(&obj, "xyz"); err != nil {
			t.Fatal(err)
		}
		if obj.ID.V != "xyz" {
			t.Errorf("got %+v, want {xyz}", obj.ID)
		}
	})
}
