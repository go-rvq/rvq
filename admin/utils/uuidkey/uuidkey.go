// Package uuidkey gives records UUID keys.
//
// A key is a UUID v7: it starts with the time it was made, so keys made one
// after the other sort — and sit in an index — together, as integers did.
//
// Register installs one gorm callback that fills, before a record is created,
// every primary-key field of type uuid.UUID still zero with a new key. It
// replaces a BeforeCreate on each model, which a model with a BeforeCreate of
// its own would silently shadow (Go promotes only the outermost method), and it
// covers what no model hook sees: a link table's rows, a record created through
// an association.
package uuidkey

import (
	"reflect"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// CallbackName is the name the create callback is registered under.
const CallbackName = "rvq:uuid_key"

var uuidType = reflect.TypeOf(uuid.UUID{})

// New is a new key: a UUID v7, or — should the clock-based generator fail — a
// random (v4) one.
func New() uuid.UUID {
	if id, err := uuid.NewV7(); err == nil {
		return id
	}
	return uuid.New()
}

// Register installs the callback on db. It is idempotent.
func Register(db *gorm.DB) error {
	cb := db.Callback().Create()
	if cb.Get(CallbackName) != nil {
		return nil
	}
	return cb.Before("gorm:create").Register(CallbackName, fill)
}

// fill gives every zero uuid.UUID primary key of the statement's record(s) a
// new key.
func fill(db *gorm.DB) {
	s := db.Statement.Schema
	if s == nil || db.Statement.ReflectValue.Kind() == reflect.Invalid {
		return
	}
	var fields []*schema.Field
	for _, f := range s.PrimaryFields {
		if f.FieldType == uuidType {
			fields = append(fields, f)
		}
	}
	if len(fields) == 0 {
		return
	}

	ctx := db.Statement.Context
	set := func(rv reflect.Value) {
		for _, f := range fields {
			if v, zero := f.ValueOf(ctx, rv); zero || v == uuid.Nil {
				_ = f.Set(ctx, rv, New())
			}
		}
	}

	rv := db.Statement.ReflectValue
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			set(reflect.Indirect(rv.Index(i)))
		}
	case reflect.Struct:
		set(rv)
	}
}

// ShortPath is the two-level path of a key for a file store: the first two
// characters, then the rest — "01/9a8f…" — so no directory holds every
// record's files.
func ShortPath(id uuid.UUID) string {
	s := id.String()
	return s[:2] + "/" + s[2:]
}

// Model is gorm.Model with a UUID key: the same fields, the same columns, for a
// model to embed instead.
type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// MustRegister is Register, for a constructor: it panics if gorm refuses the
// callback. A package whose models have UUID keys calls it on the db it is
// given, so its records get keys whoever opened the connection.
func MustRegister(db *gorm.DB) {
	if db == nil {
		return
	}
	if err := Register(db); err != nil {
		panic(err)
	}
}

// Ptr is id as a nullable key: a pointer to a copy, or nil for uuid.Nil — what
// a nullable foreign-key field holds.
func Ptr(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

// Val is the key a nullable field holds, or uuid.Nil when it holds none.
func Val(p *uuid.UUID) uuid.UUID {
	if p == nil {
		return uuid.Nil
	}
	return *p
}
