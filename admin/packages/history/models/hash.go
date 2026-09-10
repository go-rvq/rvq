package models

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
)

// Hash is a content hash: stored as bytea, but addressed as hex. Its String()
// is the hex form, so it serializes into a clean URL id (ID.String uses %v,
// which honors Stringer), and Parse() decodes that id back — the round-trip
// presets' record-id routing relies on (see presets.ParseRecordID, which tries
// a type's Parse(string) (T, error) before sql.Scanner).
type Hash []byte

// String is the hex form of the hash.
func (h Hash) String() string { return hex.EncodeToString(h) }

// Parse decodes a hex id back into a Hash.
func (Hash) Parse(s string) (Hash, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("history: bad hash %q: %w", s, err)
	}
	return Hash(b), nil
}

// Value stores the hash as bytea.
func (h Hash) Value() (driver.Value, error) {
	if len(h) == 0 {
		return nil, nil
	}
	return []byte(h), nil
}

// Scan reads the hash from a bytea (or text) column.
func (h *Hash) Scan(v any) error {
	switch b := v.(type) {
	case nil:
		*h = nil
	case []byte:
		*h = append(Hash(nil), b...)
	case string:
		*h = append(Hash(nil), b...)
	default:
		return fmt.Errorf("history: cannot scan %T into Hash", v)
	}
	return nil
}
