package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// UUID is a UUID of a record — a key, a foreign key — whose zero is none:
// it goes to the database as NULL, and a NULL comes back as the zero
// (IsZero), never as 00000000-0000-0000-0000-000000000000, which no record
// has. Not a pointer: a foreign key that is optional is a UUID that may be
// zero.
type UUID uuid.UUID

// Nil is the zero UUID: none.
var Nil UUID

// NewUUID is a new UUID: a v7 (it starts with the time it was made, so keys
// made one after the other sort together), or a random one (v4) should the
// clock-based generator fail.
func NewUUID() UUID {
	if id, err := uuid.NewV7(); err == nil {
		return UUID(id)
	}
	return UUID(uuid.New())
}

// ParseUUID is the UUID of s; "" is the zero.
func ParseUUID(s string) (UUID, error) {
	if s == "" {
		return Nil, nil
	}
	id, err := uuid.Parse(s)
	return UUID(id), err
}

// MustParseUUID is ParseUUID, panicking on a bad s.
func MustParseUUID(s string) UUID {
	id, err := ParseUUID(s)
	if err != nil {
		panic(err)
	}
	return id
}

// Parse is ParseUUID: what parses a record's id from its text (presets).
func (UUID) Parse(s string) (UUID, error) { return ParseUUID(s) }

// IsZero says it is none (zeroer.Zeroer).
func (u UUID) IsZero() bool { return u == Nil }

// UUID is it as the uuid package's.
func (u UUID) UUID() uuid.UUID { return uuid.UUID(u) }

// String is its text; "" for the zero.
func (u UUID) String() string {
	if u.IsZero() {
		return ""
	}
	return uuid.UUID(u).String()
}

// Value is it for the database: NULL for the zero.
func (u UUID) Value() (driver.Value, error) {
	if u.IsZero() {
		return nil, nil
	}
	return uuid.UUID(u).String(), nil
}

// Scan reads it from the database: NULL is the zero.
func (u *UUID) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*u = Nil
		return nil
	case string:
		if v == "" {
			*u = Nil
			return nil
		}
	case []byte:
		if len(v) == 0 {
			*u = Nil
			return nil
		}
	}
	var id uuid.UUID
	if err := id.Scan(src); err != nil {
		return fmt.Errorf("model.UUID: %w", err)
	}
	*u = UUID(id)
	return nil
}

// GormDataType is its type in the database.
func (UUID) GormDataType() string { return "uuid" }

// MarshalText is its text; "" for the zero.
func (u UUID) MarshalText() ([]byte, error) { return []byte(u.String()), nil }

// UnmarshalText reads its text; "" is the zero.
func (u *UUID) UnmarshalText(b []byte) error {
	id, err := ParseUUID(string(bytes.TrimSpace(b)))
	if err != nil {
		return err
	}
	*u = id
	return nil
}

// MarshalJSON is its text, or null for the zero.
func (u UUID) MarshalJSON() ([]byte, error) {
	if u.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(u.String())
}

// UnmarshalJSON reads its text; null and "" are the zero.
func (u *UUID) UnmarshalJSON(b []byte) error {
	if bytes.Equal(bytes.TrimSpace(b), []byte("null")) {
		*u = Nil
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return u.UnmarshalText([]byte(s))
}
