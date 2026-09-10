package admin

import (
	"encoding/json"

	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// FieldChange is one field's difference between two revisions.
type FieldChange struct {
	Field string
	Old   string
	New   string
}

// Revision loads one revision of a record by its hash.
func (h *ModelHistory) Revision(recordKey string, hash []byte) (*histmodels.Revision, error) {
	var rev histmodels.Revision
	if err := h.db.Table(h.table).
		Where("record_key = ? AND hash = ?", recordKey, hash).
		First(&rev).Error; err != nil {
		return nil, err
	}
	return &rev, nil
}

// fieldMap returns a revision's field→raw-JSON snapshot.
func fieldMap(rev *histmodels.Revision) (map[string]json.RawMessage, error) {
	if rev.Fields.Data == nil {
		return map[string]json.RawMessage{}, nil
	}
	return rev.Fields.Data, nil
}

// fieldValue returns a field's value as a display string: a JSON string is
// unquoted (so HTML/text reads naturally), anything else is its raw JSON.
func fieldValue(m map[string]json.RawMessage, field string) string {
	raw, ok := m[field]
	if !ok {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	return string(raw)
}

// Diff compares every versioned field between two revisions of a record,
// returning only the fields that differ.
func (h *ModelHistory) Diff(recordKey string, aHash, bHash []byte) ([]FieldChange, error) {
	am, bm, err := h.twoFieldMaps(recordKey, aHash, bHash)
	if err != nil {
		return nil, err
	}
	var changes []FieldChange
	for _, f := range h.resolved {
		ov, nv := fieldValue(am, f), fieldValue(bm, f)
		if ov != nv {
			changes = append(changes, FieldChange{Field: f, Old: ov, New: nv})
		}
	}
	return changes, nil
}

// FieldDiff returns one field's old and new values across two revisions.
func (h *ModelHistory) FieldDiff(recordKey, field string, aHash, bHash []byte) (old, now string, err error) {
	am, bm, err := h.twoFieldMaps(recordKey, aHash, bHash)
	if err != nil {
		return "", "", err
	}
	return fieldValue(am, field), fieldValue(bm, field), nil
}

func (h *ModelHistory) twoFieldMaps(recordKey string, aHash, bHash []byte) (a, b map[string]json.RawMessage, err error) {
	ra, err := h.Revision(recordKey, aHash)
	if err != nil {
		return
	}
	rb, err := h.Revision(recordKey, bHash)
	if err != nil {
		return
	}
	if a, err = fieldMap(ra); err != nil {
		return
	}
	b, err = fieldMap(rb)
	return
}

// HTMLDiff renders an inline (ins/del) HTML diff of two field values — used for
// HTML/text fields such as Body.
func HTMLDiff(old, now string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(old, now, false)
	dmp.DiffCleanupSemantic(diffs)
	return dmp.DiffPrettyHtml(diffs)
}
