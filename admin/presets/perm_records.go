package presets

import (
	"fmt"

	"github.com/go-rvq/rvq/web"
)

// The records a permission names ("<7>", RecordPermPart): found by their id
// or their title, for a field that writes permissions (vx-permission-input),
// only those the user may read — as the record selector of the menu.

// PermRecord is a record in a permission: its part's id ("7" of "<7>") and
// its title.
type PermRecord struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// permRecordID is the id of obj as a permission names it.
func (mb *ModelBuilder) permRecordID(obj any) (string, bool) {
	id, err := mb.RecordID(obj)
	if err != nil || id.IsZero() {
		return "", false
	}
	v := id.GetValue("ID")
	if v == nil {
		return "", false
	}
	return fmt.Sprint(v), true
}

// PermRecords are the records of mb, of the parents' records, the user of
// ctx may read, whose title or id has keyword: at most limit.
func (mb *ModelBuilder) PermRecords(ctx *web.EventContext, parents IDSlice, keyword string, limit int) (out []PermRecord, err error) {
	objs, _, err := mb.menuRecordSearch(ctx, parents, keyword, nil, 1, int64(limit))
	if err != nil {
		return nil, err
	}
	if keyword != "" && len(objs) == 0 {
		// the id itself
		if obj := mb.menuRecordFind(ctx, parents, keyword); obj != nil {
			objs = append(objs, obj)
		}
	}
	for _, obj := range objs {
		if _, ok := mb.menuRecordReadable(ctx.R, parents, obj); !ok {
			continue
		}
		if id, ok := mb.permRecordID(obj); ok {
			out = append(out, PermRecord{ID: id, Label: mb.RecordTitle(obj, withParents(ctx, parents))})
		}
	}
	return
}

// PermRecordLabel is the title of the record of mb of the id, of the
// parents' records; ok false when there is none the user of ctx may read.
func (mb *ModelBuilder) PermRecordLabel(ctx *web.EventContext, parents IDSlice, id string) (label string, ok bool) {
	obj := mb.menuRecordFind(ctx, parents, id)
	if obj == nil {
		return "", false
	}
	return mb.RecordTitle(obj, withParents(ctx, parents)), true
}
