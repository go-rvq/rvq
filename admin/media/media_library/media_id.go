package media_library

import (
	"bytes"
	"encoding/json"
)

// MediaID is the key of the media library record a MediaBox shows, as text:
// "" when the box is empty.
//
// It reads what boxes saved before keys were UUIDs as well — a number, with 0
// for an empty box — and writes a string.
type MediaID string

func (id MediaID) String() string { return string(id) }

// IsZero reports whether the box points at no record.
func (id MediaID) IsZero() bool { return id == "" || id == "0" }

func (id MediaID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

func (id *MediaID) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	switch {
	case bytes.Equal(b, []byte("null")):
		*id = ""
		return nil
	case len(b) > 0 && b[0] == '"':
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*id = MediaID(s)
	default:
		var n json.Number
		if err := json.Unmarshal(b, &n); err != nil {
			return err
		}
		*id = MediaID(n)
	}
	if *id == "0" {
		*id = ""
	}
	return nil
}
