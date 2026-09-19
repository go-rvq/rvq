package admin

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/gorm"

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
	// A fresh session with an explicit ORDER (First would otherwise order by the
	// Revision schema's table "revisions", which is not the FROM — the model maps
	// to <table>_revisions via .Table): both avoid a leaked scope and the wrong
	// table qualifier.
	if err := h.db.Session(&gorm.Session{}).Table(h.table).
		Where("record_key = ? AND hash = ?", recordKey, hash).
		Order("created_at").Take(&rev).Error; err != nil {
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

// ChangedFieldsLabel formats a revision's changed fields as a nested label:
// "Field1, Field2, Field3 [ Sub1, Sub2 [ sub ] ]". Top-level fields are in the
// model's versioned order; a structured field that only changed in some of its
// sub-fields is shown with those sub-fields (recursively) in brackets, computed
// by comparing this revision's value against the parent's. The first revision
// (no parent) lists the top-level fields only (everything is new).
func (h *ModelHistory) ChangedFieldsLabel(rev *histmodels.Revision, ctx *web.EventContext) string {
	changed := map[string]bool{}
	for _, f := range rev.ChangedFields.Data {
		changed[f] = true
	}
	m, _ := fieldMap(rev)

	var parentMap map[string]json.RawMessage
	if len(rev.Parent) > 0 {
		if p, err := h.Revision(rev.RecordKey, rev.Parent); err == nil {
			parentMap, _ = fieldMap(p)
		}
	}

	var parts []string
	for _, f := range h.resolved {
		if !changed[f] {
			continue
		}
		if parentMap == nil {
			parts = append(parts, fieldLabelOf(h.mb, f, ctx)) // root: top-level labels only
			continue
		}
		parts = append(parts, formatChangedField(h.mb, f, parentMap[f], m[f], ctx))
	}
	return strings.Join(parts, ", ")
}

// formatChangedField renders one changed field (label-translated via mb's field
// translator), descending into a structured value to show only the sub-fields
// that differ: "Name" for a leaf change, "Name [ sub1, sub2 ]" when it is an
// object whose sub-keys changed. Sub-fields are translated with the child model
// builder of the field (mb.GetChildByID(name)) when there is one.
func formatChangedField(mb *presets.ModelBuilder, name string, oldRaw, newRaw json.RawMessage, ctx *web.EventContext) string {
	label := fieldLabelOf(mb, name, ctx)
	child := nestedModelOf(mb, name)
	var oldObj, newObj map[string]json.RawMessage
	if json.Unmarshal(oldRaw, &oldObj) == nil && json.Unmarshal(newRaw, &newObj) == nil {
		keys := sortedUnionKeys(oldObj, newObj)
		var subs []string
		for _, k := range keys {
			if rawChanged(oldObj[k], newObj[k]) {
				subs = append(subs, formatChangedField(child, k, oldObj[k], newObj[k], ctx))
			}
		}
		if len(subs) > 0 {
			return label + " [ " + strings.Join(subs, ", ") + " ]"
		}
	}
	return label
}

// fieldLabelOf translates a field name through mb's field translator (keys
// <Model><Field>), falling back to the raw name. A nil mb or ctx yields the name.
func fieldLabelOf(mb *presets.ModelBuilder, name string, ctx *web.EventContext) string {
	if mb == nil || ctx == nil {
		return name
	}
	if s := i18n.Translate(mb.FieldTranslator(), ctx.Context(), name); s != "" && s != name {
		return s
	}
	return name
}

// nestedModelOf returns the model builder that describes a field's sub-fields —
// the field's Nested (NestedStruct / NestedSlice) Model — so sub-field labels
// translate with the right model. Looks in the Editing then Detailing fields
// builders. nil when the field is not nested.
func nestedModelOf(mb *presets.ModelBuilder, name string) *presets.ModelBuilder {
	if mb == nil {
		return nil
	}
	fbs := []*presets.FieldsBuilder{&mb.Editing().FieldsBuilder}
	if mb.HasDetailing() {
		fbs = append(fbs, &mb.Detailing().FieldsBuilder)
	}
	for _, fb := range fbs {
		if fb == nil {
			continue
		}
		if f := fb.GetField(name); f != nil {
			if n := f.GetNested(); n != nil {
				if cm := n.Model(); cm != nil {
					return cm
				}
			}
		}
	}
	return nil
}

// fieldLabelPath translates a (possibly nested/indexed) field path to a readable
// label, descending child model builders for each segment: "PageOptions.Layout"
// → "Opções da Página / Layout"; "Tags[0].Name" keeps the index on the segment.
func (h *ModelHistory) fieldLabelPath(path string, ctx *web.EventContext) string {
	mb := h.mb
	parts := strings.Split(path, ".")
	out := make([]string, len(parts))
	for i, p := range parts {
		name, idx := p, ""
		if j := strings.IndexByte(p, '['); j >= 0 {
			name, idx = p[:j], p[j:]
		}
		out[i] = fieldLabelOf(mb, name, ctx) + idx
		mb = nestedModelOf(mb, name)
	}
	return strings.Join(out, " / ")
}

// sortedUnionKeys is the sorted union of two maps' keys.
func sortedUnionKeys(a, b map[string]json.RawMessage) []string {
	set := map[string]bool{}
	for k := range a {
		set[k] = true
	}
	for k := range b {
		set[k] = true
	}
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// rawChanged reports whether two raw JSON values differ (compared compacted, so
// formatting never counts as a change).
func rawChanged(a, b json.RawMessage) bool {
	return !bytes.Equal(compactJSON(a), compactJSON(b))
}

func compactJSON(r json.RawMessage) []byte {
	if len(r) == 0 {
		return nil
	}
	var buf bytes.Buffer
	if json.Compact(&buf, r) != nil {
		return r
	}
	return buf.Bytes()
}

// topField reduces a (possibly nested/indexed) path to its top-level versioned
// field name ("PageOptions.Layout" → "PageOptions", "Tags[0].Name" → "Tags").
func topField(path string) string {
	if i := strings.IndexAny(path, ".["); i >= 0 {
		return path[:i]
	}
	return path
}

// revisionChanged reports whether a revision recorded a change to the given
// (possibly nested) field, from its ChangedFields list (matched on the
// top-level field). A revision with no list (legacy) is treated as changed so it
// is never hidden.
func revisionChanged(rev *histmodels.Revision, field string) bool {
	cf := rev.ChangedFields.Data
	if cf == nil {
		return true
	}
	top := topField(field)
	for _, c := range cf {
		if c == top {
			return true
		}
	}
	return false
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
