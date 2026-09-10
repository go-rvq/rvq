package admin

import (
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
)

// Chain returns a record's revisions oldest → newest.
func (h *ModelHistory) Chain(recordKey string) ([]histmodels.Revision, error) {
	var revs []histmodels.Revision
	err := h.db.Table(h.table).
		Where("record_key = ?", recordKey).
		Order("created_at ASC").
		Find(&revs).Error
	return revs, err
}

// FieldRevision is a revision as it left one field.
type FieldRevision struct {
	Revision histmodels.Revision
	Value    string
}

// FieldHistory returns the revisions in which a given field's value changed —
// the timeline of just that field (e.g. every edit of Post.Body). Consecutive
// revisions that left the field untouched are collapsed.
func (h *ModelHistory) FieldHistory(recordKey, field string) ([]FieldRevision, error) {
	revs, err := h.Chain(recordKey)
	if err != nil {
		return nil, err
	}
	var out []FieldRevision
	var prev string
	first := true
	for i := range revs {
		m, err := fieldMap(&revs[i])
		if err != nil {
			return nil, err
		}
		v := fieldValue(m, field)
		if first || v != prev {
			out = append(out, FieldRevision{Revision: revs[i], Value: v})
			prev = v
			first = false
		}
	}
	return out, nil
}
