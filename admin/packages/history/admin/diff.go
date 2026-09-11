package admin

import (
	"encoding/json"
	"strconv"
	"strings"

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
// pathToken is one step of a field path: a struct-field key or an array index.
type pathToken struct {
	key   string
	idx   int
	isIdx bool
}

// parsePath breaks a field path into tokens, supporting nested struct fields and
// array indices: "PageOptions.Layout", "Galleries[0].Title", "a[0][1].b".
func parsePath(path string) (toks []pathToken) {
	for _, seg := range strings.Split(path, ".") {
		for seg != "" {
			i := strings.IndexByte(seg, '[')
			if i < 0 {
				toks = append(toks, pathToken{key: seg})
				break
			}
			if i > 0 {
				toks = append(toks, pathToken{key: seg[:i]})
			}
			j := strings.IndexByte(seg, ']')
			if j < 0 {
				return
			}
			n, err := strconv.Atoi(seg[i+1 : j])
			if err != nil {
				return
			}
			toks = append(toks, pathToken{idx: n, isIdx: true})
			seg = seg[j+1:]
		}
	}
	return
}

// navigateRaw follows a field path into a snapshot map, descending struct fields
// (objects) and array indices. It returns the raw JSON at the path.
func navigateRaw(m map[string]json.RawMessage, path string) (json.RawMessage, bool) {
	toks := parsePath(path)
	if len(toks) == 0 || toks[0].isIdx {
		return nil, false
	}
	raw, ok := m[toks[0].key]
	if !ok {
		return nil, false
	}
	for _, t := range toks[1:] {
		if t.isIdx {
			var arr []json.RawMessage
			if json.Unmarshal(raw, &arr) != nil || t.idx < 0 || t.idx >= len(arr) {
				return nil, false
			}
			raw = arr[t.idx]
			continue
		}
		var obj map[string]json.RawMessage
		if json.Unmarshal(raw, &obj) != nil {
			return nil, false
		}
		if raw, ok = obj[t.key]; !ok {
			return nil, false
		}
	}
	return raw, true
}

// fieldValue returns a field's value as a display string, following a path
// (nested fields and array indices — "PageOptions.Layout", "Tags[0].Name"). A
// JSON string is unquoted; anything else is its raw JSON.
func fieldValue(m map[string]json.RawMessage, field string) string {
	raw, ok := navigateRaw(m, field)
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
