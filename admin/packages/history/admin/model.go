package admin

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/go-rvq/rvq/admin/activity"
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
	"github.com/go-rvq/rvq/admin/packages/user"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/thirdpart/gorm/datatypes"
	"github.com/go-rvq/rvq/web"
	"github.com/sunfmin/reflectutils"
	"gorm.io/gorm"
)

// ModelHistory activates revisions for one model. Fluent:
//
//	history.New(db).Model(mb).Build()                 // versions the EDIT fields
//	history.New(db).Model(mb).Fields("Body").Build()  // only these fields
//	history.New(db).Model(mb).AllFields().Build()     // every struct field
//
// WholeFields marks fields that reject partial (content-level) revert — a field
// whose value would be corrupted by a hunk-level patch (foreign keys, JSON,
// relations). By default such fields are inferred; WholeFields adds overrides.
type ModelHistory struct {
	db           *gorm.DB
	mb           *presets.ModelBuilder
	fields       []string
	allFields    bool
	wholeFields  map[string]bool
	htmlFields   map[string]bool
	fieldDiffers map[string]FieldDiffFunc

	table    string
	resolved []string
}

// New starts a per-model history activation on db.
func New(db *gorm.DB) *ModelHistory {
	return &ModelHistory{
		db:           db,
		wholeFields:  map[string]bool{},
		htmlFields:   map[string]bool{},
		fieldDiffers: map[string]FieldDiffFunc{},
	}
}

func (h *ModelHistory) Model(mb *presets.ModelBuilder) *ModelHistory { h.mb = mb; return h }

// Fields sets the versioned fields explicitly. With none, the model's EDIT
// fields are used.
func (h *ModelHistory) Fields(names ...string) *ModelHistory { h.fields = names; return h }

// AllFields versions every struct field of the model.
func (h *ModelHistory) AllFields() *ModelHistory { h.allFields = true; return h }

// WholeFields marks fields that only accept whole-field (not partial content)
// revert.
func (h *ModelHistory) WholeFields(names ...string) *ModelHistory {
	for _, n := range names {
		h.wholeFields[n] = true
	}
	return h
}

// HTMLFields marks fields whose value is HTML (e.g. Body), so the diff UI uses
// an HTML diff for them instead of a plain-text one.
func (h *ModelHistory) HTMLFields(names ...string) *ModelHistory {
	for _, n := range names {
		h.htmlFields[n] = true
	}
	return h
}

// Table is the model's revisions table name (<model table>_revisions).
func (h *ModelHistory) Table() string { return h.table }

// IsHTML reports whether a field was declared as HTML (see HTMLFields).
func (h *ModelHistory) IsHTML(field string) bool { return h.htmlFields[field] }

// Fields returns the resolved versioned field names.
func (h *ModelHistory) ResolvedFields() []string { return h.resolved }

// AcceptsPartial reports whether a field accepts partial (content-level) revert.
func (h *ModelHistory) AcceptsPartial(field string) bool { return !h.wholeFields[field] }

// Build resolves the table and versioned fields, migrates the revisions table,
// and wraps the model's save to capture a revision.
func (h *ModelHistory) Build() *ModelHistory {
	if h.mb == nil {
		panic("history: Model(mb) is required before Build()")
	}
	h.table = revisionTable(h.db, h.mb.NewModel())
	if err := h.db.Table(h.table).AutoMigrate(&histmodels.Revision{}); err != nil {
		panic(err)
	}
	h.resolved = h.resolveFields()
	revisionsByTable[h.table] = h

	h.mb.Editing().WrapSaveFunc(func(in presets.SaveFunc) presets.SaveFunc {
		return func(obj interface{}, id model.ID, ctx *web.EventContext) error {
			if err := in(obj, id, ctx); err != nil {
				return err
			}
			return h.capture(obj, ctx)
		}
	})
	h.installPublishTag()
	h.installChild()
	return h
}

// revisionTable derives "<model table>_revisions" from the gorm schema.
func revisionTable(db *gorm.DB, m any) string {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(m); err != nil {
		panic(err)
	}
	return stmt.Schema.Table + RevisionsSuffix
}

// resolveFields picks the versioned field names: explicit Fields, else every
// struct field (AllFields), else the model's EDIT fields.
func (h *ModelHistory) resolveFields() []string {
	if len(h.fields) > 0 {
		return h.fields
	}
	if h.allFields {
		stmt := &gorm.Statement{DB: h.db}
		if err := stmt.Parse(h.mb.NewModel()); err != nil {
			panic(err)
		}
		names := make([]string, 0, len(stmt.Schema.Fields))
		for _, f := range stmt.Schema.Fields {
			names = append(names, f.Name)
		}
		return names
	}
	var names []string
	for _, n := range h.mb.Editing().FieldNames() {
		if s, ok := n.(string); ok {
			names = append(names, s)
		}
	}
	return names
}

// snapshot captures the versioned fields of obj as a field→raw-JSON map, plus
// its canonical bytes (json.Marshal sorts map keys, so the bytes are
// deterministic for the same values — the basis of the content hash).
func (h *ModelHistory) snapshot(obj interface{}) (map[string]json.RawMessage, []byte, error) {
	snap := make(map[string]json.RawMessage, len(h.resolved))
	for _, f := range h.resolved {
		v, err := reflectutils.Get(obj, f)
		if err != nil {
			// virtual/non-struct field (e.g. an action column) — skip it.
			continue
		}
		raw, err := json.Marshal(v)
		if err != nil {
			return nil, nil, err
		}
		snap[f] = raw
	}
	data, err := json.Marshal(snap)
	return snap, data, err
}

// changedFields lists the versioned fields of snap whose value differs from the
// parent revision. With no parent (the first revision) every versioned field is
// "changed" — the record's initial state.
func (h *ModelHistory) changedFields(recordKey string, parent histmodels.Hash, snap map[string]json.RawMessage) ([]string, error) {
	if parent == nil {
		out := make([]string, len(h.resolved))
		copy(out, h.resolved)
		return out, nil
	}
	prev, err := h.Revision(recordKey, parent)
	if err != nil {
		return nil, err
	}
	pm, err := fieldMap(prev)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, f := range h.resolved {
		if fieldValue(snap, f) != fieldValue(pm, f) {
			out = append(out, f)
		}
	}
	return out, nil
}

// capture writes a revision of obj if the versioned fields changed. The hash is
// the pure content hash; a revision with the same (hash, record) already stored
// means nothing changed, so none is created.
func (h *ModelHistory) capture(obj interface{}, ctx *web.EventContext) error {
	snap, data, err := h.snapshot(obj)
	if err != nil {
		return err
	}
	recordKey := h.mb.MustRecordID(obj).String()
	// Fold the record key into the hash so it is unique per record (sole PK,
	// clean nested route) yet still collapses an unchanged save.
	sum := sha256.Sum256(append([]byte(recordKey+"\x00"), data...))
	hash := sum[:]

	var exists int64
	if err = h.db.Table(h.table).
		Where("hash = ? AND record_key = ?", hash, recordKey).
		Count(&exists).Error; err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}

	parent, err := h.latestHash(recordKey)
	if err != nil {
		return err
	}

	// The fields whose value changed from the parent (all of them on the first
	// revision). Kept alongside the full snapshot so the history is queryable by
	// field without walking the chain or re-diffing.
	changed, err := h.changedFields(recordKey, parent, snap)
	if err != nil {
		return err
	}

	rev := histmodels.Revision{
		Hash:          hash,
		RecordKey:     recordKey,
		Parent:        parent,
		Fields:        datatypes.NewJSONType(snap),
		ChangedFields: datatypes.NewJSONType(changed),
		CreatedAt:     time.Now(),
	}
	if u := user.GetCurrentUser(ctx.R); u != nil {
		rev.CreatorID = u.GetID()
		rev.Creator = u.GetName()
	}
	if err = h.db.Table(h.table).Create(&rev).Error; err != nil {
		return err
	}
	// Tell the activity log (if any is being written for this save) to reference
	// this revision instead of duplicating its diff.
	if ctx != nil && ctx.R != nil {
		ctx.R = activity.WithRevisionRef(ctx.R, h.table, rev.Hash)
	}
	return nil
}
