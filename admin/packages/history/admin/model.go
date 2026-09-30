package admin

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/go-rvq/rvq/admin/activity"
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/thirdpart/gorm/datatypes"
	"github.com/go-rvq/rvq/web"
	"github.com/google/uuid"
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
	extraFields  []string
	allFields    bool
	wholeFields  map[string]bool
	noRevert     map[string]bool
	htmlFields   map[string]bool
	refFields    map[string]bool
	fieldDiffers map[string]FieldDiffFunc
	fieldCodecs  map[string]FieldContentCodec
	fieldLangs   map[string]string
	fetcher      presets.FetchFunc

	table      string
	resolved   []string
	schemaRefs map[string]bool // lazily-detected relation fields (nil until computed)
}

// FieldContentCodec adapts a non-string field (e.g. a JSON map) to the text-based
// partial-revert flow: ToText turns the field's stored JSON into the editable
// text form (e.g. YAML), and Apply sets the field back from patched text (e.g.
// parse YAML into the map and assign it). Register one with FieldContentHandler.
type FieldContentCodec struct {
	ToText func(rawJSON string) string
	Apply  func(obj any, text string) error
}

// New starts a per-model history activation on db.
func New(db *gorm.DB) *ModelHistory {
	return &ModelHistory{
		db:           db,
		wholeFields:  map[string]bool{},
		noRevert:     map[string]bool{},
		htmlFields:   map[string]bool{},
		refFields:    map[string]bool{},
		fieldDiffers: map[string]FieldDiffFunc{},
		fieldCodecs:  map[string]FieldContentCodec{},
		fieldLangs:   map[string]string{},
	}
}

// FieldContentHandler registers a text codec (and a Prism language) for a field,
// so a non-string field (e.g. a JSON map) can take part in the text-based partial
// revert: it is diffed/edited as text (language), and patched text is applied
// back through the codec. Registering a codec is what lets such a field accept
// partial revert.
func (h *ModelHistory) FieldContentHandler(field, language string, codec FieldContentCodec) *ModelHistory {
	h.fieldCodecs[field] = codec
	if language != "" {
		h.fieldLangs[field] = language
	}
	return h
}

// fieldLanguage is the Prism language registered for a field (empty = plain).
func (h *ModelHistory) fieldLanguage(field string) string { return h.fieldLangs[field] }

func (h *ModelHistory) Model(mb *presets.ModelBuilder) *ModelHistory { h.mb = mb; return h }

// Fields sets the versioned fields explicitly. With none, the model's EDIT
// fields are used.
func (h *ModelHistory) Fields(names ...string) *ModelHistory { h.fields = names; return h }

// AllFields versions every struct field of the model.
func (h *ModelHistory) AllFields() *ModelHistory { h.allFields = true; return h }

// ExtraFields versions additional fields beyond the resolved set (EDIT fields by
// default) — for values edited outside the model's own form, e.g. a SEO Setting
// managed by a nested model, so changing them still records a revision.
func (h *ModelHistory) ExtraFields(names ...string) *ModelHistory {
	h.extraFields = append(h.extraFields, names...)
	return h
}

// WholeFields marks fields that only accept whole-field (not partial content)
// revert.
func (h *ModelHistory) WholeFields(names ...string) *ModelHistory {
	for _, n := range names {
		h.wholeFields[n] = true
	}
	return h
}

// NoRevertFields marks fields a revert never restores. They are versioned and
// compared like any other — what they said before is part of the record's
// history, and often what explains the rest of it — but the value they have now
// stays: a revert of the whole record touches every other field and leaves
// these alone, and a revert of one of them is refused.
//
// It is for a field that is not the editor's to choose: one the application
// keeps in step with something else (a configuration file, another record), so
// restoring an old value would only put it out of step until the next sync.
func (h *ModelHistory) NoRevertFields(names ...string) *ModelHistory {
	for _, n := range names {
		h.noRevert[n] = true
	}
	return h
}

// Revertible reports whether a revert may restore this field.
func (h *ModelHistory) Revertible(field string) bool { return !h.noRevert[field] }

// revertible is the fields of names a revert may restore, in order.
func (h *ModelHistory) revertible(names []string) []string {
	out := make([]string, 0, len(names))
	for _, n := range names {
		if h.Revertible(n) {
			out = append(out, n)
		}
	}
	return out
}

// HTMLFields marks fields whose value is HTML (e.g. Body), so the diff UI uses
// an HTML diff for them instead of a plain-text one.
func (h *ModelHistory) HTMLFields(names ...string) *ModelHistory {
	for _, n := range names {
		h.htmlFields[n] = true
	}
	return h
}

// ReferenceFields marks fields that select another model (a ModelSelector /
// foreign-key field, e.g. an author). Their diff compares the whole value (never
// partial) and shows the referenced record rendered — its title with the id
// beside it. Relation fields are auto-detected from the gorm schema; this is the
// manual override for cases the schema does not surface.
func (h *ModelHistory) ReferenceFields(names ...string) *ModelHistory {
	for _, n := range names {
		h.refFields[n] = true
	}
	return h
}

// Fetcher overrides how a record is (re)loaded before snapshotting it, with the
// same signature as ModelBuilder's fetch (presets.FetchFunc). By default the
// model's own Editing fetcher is used, which applies its configured preloads so
// associations are present in the snapshot; set this to control loading (e.g. to
// add preloads the model's fetcher does not).
func (h *ModelHistory) Fetcher(fn presets.FetchFunc) *ModelHistory {
	h.fetcher = fn
	return h
}

// Table is the model's revisions table name (<model table>_revisions).
func (h *ModelHistory) Table() string { return h.table }

// IsHTML reports whether a field was declared as HTML (see HTMLFields).
func (h *ModelHistory) IsHTML(field string) bool { return h.htmlFields[field] }

// Fields returns the resolved versioned field names.
func (h *ModelHistory) ResolvedFields() []string { return h.resolved }

// AcceptsPartial reports whether a field accepts partial (content-level) revert.
// A field with a registered content codec always does (it is patched as text —
// e.g. a JSON map edited as YAML). Otherwise a reference/relation or whole-only
// field never does: a hunk-level patch would corrupt it.
func (h *ModelHistory) AcceptsPartial(field string) bool {
	if !h.Revertible(field) {
		return false
	}
	if _, ok := h.fieldCodecs[field]; ok {
		return true
	}
	return !h.wholeFields[field] && !h.isReferenceField(field)
}

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

	// Migrate legacy rows: fill ChangedFields on revisions saved before the
	// column existed. Idempotent and cheap when there is nothing to do; a failure
	// must not block boot.
	if err := h.Backfill(); err != nil {
		log.Printf("history: backfill %s: %v", h.table, err)
	}

	// Seed the initial revision for records that predate history (a database
	// populated before revisions were configured). Idempotent; must not block boot.
	if err := h.SeedInitialRevisions(); err != nil {
		log.Printf("history: seed %s: %v", h.table, err)
	}

	h.mb.Editing().WrapSaveFunc(func(in presets.SaveFunc) presets.SaveFunc {
		return func(obj interface{}, id model.ID, ctx *web.EventContext) error {
			if err := in(obj, id, ctx); err != nil {
				return err
			}
			return h.capture(obj, ctx)
		}
	})
	// Create uses a separate Creator (not the Saver), so wrap it too — otherwise
	// the first save of a new record records no initial revision.
	h.mb.Editing().WrapCreateFunc(func(in presets.CreateFunc) presets.CreateFunc {
		return func(obj interface{}, ctx *web.EventContext) error {
			if err := in(obj, ctx); err != nil {
				return err
			}
			return h.capture(obj, ctx)
		}
	})
	h.installDeletion()
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
// struct field (AllFields), else the model's EDIT fields; then any ExtraFields
// are appended (deduplicated) — so a field edited outside the form (e.g. a SEO
// Setting managed by a nested model) can still be versioned.
func (h *ModelHistory) resolveFields() []string {
	var names []string
	switch {
	case len(h.fields) > 0:
		names = append(names, h.fields...)
	case h.allFields:
		stmt := &gorm.Statement{DB: h.db}
		if err := stmt.Parse(h.mb.NewModel()); err != nil {
			panic(err)
		}
		for _, f := range stmt.Schema.Fields {
			names = append(names, f.Name)
		}
	default:
		for _, n := range h.mb.Editing().FieldNames() {
			if s, ok := n.(string); ok {
				names = append(names, s)
			}
		}
	}
	for _, e := range h.extraFields {
		if !slicesContains(names, e) {
			names = append(names, e)
		}
	}
	return names
}

func slicesContains(ss []string, v string) bool {
	for _, s := range ss {
		if s == v {
			return true
		}
	}
	return false
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

// Backfill fills ChangedFields on revisions that predate the column (stored
// NULL): for each such record it walks the chain oldest→newest and records, per
// revision, the fields that changed from the previous one (all on the first).
// Idempotent — it only reads records that still have a NULL row and only writes
// the rows still missing the list, so re-running (every boot) is a no-op once
// done.
// SeedInitialRevisions creates the first revision (a mirror of the current
// record) for every record of the model that has no revision yet. This is the
// case of configuring history on an already-populated database: existing records
// predate the revisions table, so without a seed they would have no baseline to
// compare or revert to. Idempotent: a record that already has any revision is
// skipped, so it is a no-op on every boot after the first. Best-effort per
// record — one failure is logged, the rest proceed.
func (h *ModelHistory) SeedInitialRevisions() error {
	slicePtr := h.mb.NewModelSlice() // *[]Model
	if err := h.db.Session(&gorm.Session{}).Find(slicePtr).Error; err != nil {
		return err
	}
	sliceVal := reflect.ValueOf(slicePtr).Elem() // []Model

	// A synthetic context so the fetcher (which applies the model's preloads) can
	// load each record's associations for a complete snapshot.
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	ctx := &web.EventContext{R: req}
	fetch := h.fetcher
	if fetch == nil {
		fetch = h.mb.Editing().Fetcher
	}

	var firstErr error
	for i := 0; i < sliceVal.Len(); i++ {
		// The addressable pointer to the element, so MustRecordID and the fetcher
		// see the concrete record.
		obj := sliceVal.Index(i).Addr().Interface()

		id := h.mb.MustRecordID(obj)
		if id.IsZero() {
			continue
		}
		recordKey := id.String()

		var n int64
		if err := h.db.Table(h.table).Where("record_key = ?", recordKey).Count(&n).Error; err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if n > 0 {
			continue // already has history
		}

		// Reload with associations for a faithful mirror.
		if fetch != nil {
			fresh := h.mb.NewModel()
			if fetch(fresh, id, ctx) == nil {
				obj = fresh
			}
		}
		if _, _, err := h.createRevision(obj, uuid.Nil, "", histmodels.Origin{}); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (h *ModelHistory) Backfill() error {
	var keys []string
	if err := h.db.Table(h.table).
		Where("changed_fields IS NULL").
		Distinct().
		Pluck("record_key", &keys).Error; err != nil {
		return err
	}
	for _, key := range keys {
		revs, err := h.Chain(key) // oldest → newest
		if err != nil {
			return err
		}
		var prev map[string]json.RawMessage
		for i := range revs {
			snap, err := fieldMap(&revs[i])
			if err != nil {
				return err
			}
			var changed []string
			if prev == nil {
				changed = append(changed, h.resolved...)
			} else {
				for _, f := range h.resolved {
					if fieldValue(snap, f) != fieldValue(prev, f) {
						changed = append(changed, f)
					}
				}
			}
			if revs[i].ChangedFields.Data == nil {
				if err := h.db.Table(h.table).
					Where("hash = ? AND record_key = ?", []byte(revs[i].Hash), key).
					Update("changed_fields", datatypes.NewJSONType(changed)).Error; err != nil {
					return err
				}
			}
			prev = snap
		}
	}
	return nil
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
	// Reload through the fetcher (the configured one, else the model's Editing
	// fetcher) so associations — e.g. PageOptions and its Galleries — are present
	// in the snapshot. The just-saved obj may carry them nil, which would store a
	// null value and later show a false "changed" diff. Best-effort: keep obj on
	// failure.
	fetch := h.fetcher
	if fetch == nil {
		fetch = h.mb.Editing().Fetcher
	}
	if fetch != nil {
		if id := h.mb.MustRecordID(obj); !id.IsZero() {
			fresh := h.mb.NewModel()
			if fetch(fresh, id, ctx) == nil {
				obj = fresh
			}
		}
	}

	var creatorID uuid.UUID
	var creator string
	var origin histmodels.Origin
	if ctx != nil && ctx.R != nil {
		if u := user.GetCurrentUser(ctx.R); u != nil {
			creatorID = u.GetID()
			creator = u.GetName()
		}
		origin = originOf(ctx.R)
	}

	hash, created, err := h.createRevision(obj, creatorID, creator, origin)
	if err != nil || !created {
		return err
	}
	// Tell the activity log (if any is being written for this save) to reference
	// this revision instead of duplicating its diff.
	if ctx != nil && ctx.R != nil {
		ctx.R = activity.WithRevisionRef(ctx.R, h.table, hash)
	}
	return nil
}

// Record records the revision of obj — a record of the model changed outside
// the admin's forms: by a site, a job, a command — with its author (creatorID
// and creator; uuid.Nil and "" for none) and where they were: the request r
// they made it by (SetOriginFunc; nil: nowhere known). obj must be loaded
// whole (the associations of the versioned fields too): its snapshot is what
// the revision keeps. Nothing when it did not change since its last revision
// (the same snapshot); the hash is the revision's, and whether it was created.
func (h *ModelHistory) Record(obj any, creatorID uuid.UUID, creator string, r *http.Request) (histmodels.Hash, bool, error) {
	return h.createRevision(obj, creatorID, creator, originOf(r))
}

// RecordDeletion records the deletion of obj — deleted outside the admin's
// forms: by a site, a job, a command —: a revision of its own (EventDeleted),
// its snapshot obj as it was — loaded whole, as for Record —, with who deleted
// it (creatorID, creator) and from where (the request r; nil: nowhere known).
func (h *ModelHistory) RecordDeletion(obj any, creatorID uuid.UUID, creator string, r *http.Request) (histmodels.Hash, error) {
	hash, _, err := h.createEventRevision(obj, histmodels.EventDeleted, creatorID, creator, originOf(r))
	return hash, err
}

// installDeletion records, in a revision of its own, the deletion of a record
// by the admin — the model's data operator (gorm2op) —: the record as it was
// before, who deleted it and from where.
func (h *ModelHistory) installDeletion() {
	if _, ok := h.mb.CurrentDataOperator().(*gorm2op.DataOperatorBuilder); !ok {
		return
	}
	type deletingKey struct{}
	h.mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		return do.(*gorm2op.DataOperatorBuilder).WithDeleteCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
			cb.Pre(func(state *gorm2op.CallbackState) error {
				// the record as it is, before it is deleted
				obj := state.Obj
				if id := h.mb.MustRecordID(obj); !id.IsZero() && state.Ctx != nil {
					fresh := h.mb.NewModel()
					if h.mb.Editing().Fetcher(fresh, id, state.Ctx) == nil {
						obj = fresh
					}
				}
				state.Set(deletingKey{}, obj)
				return nil
			})
			cb.Post(func(state *gorm2op.CallbackState) error {
				obj := state.Get(deletingKey{})
				if obj == nil {
					return nil
				}
				var (
					creatorID uuid.UUID
					creator   string
					origin    histmodels.Origin
				)
				if state.Ctx != nil && state.Ctx.R != nil {
					if u := user.GetCurrentUser(state.Ctx.R); u != nil {
						creatorID, creator = u.GetID(), u.GetName()
					}
					origin = originOf(state.Ctx.R)
				}
				_, _, err := h.createEventRevision(obj, histmodels.EventDeleted, creatorID, creator, origin)
				return err
			})
		})
	})
}

// createRevision snapshots obj (already loaded, with its associations) and
// creates its revision unless an identical one already exists (dedup). It needs
// no request context, so it also serves the boot-time seed. Returns the hash and
// whether a new revision was created.
func (h *ModelHistory) createRevision(obj interface{}, creatorID uuid.UUID, creator string, origin histmodels.Origin) (histmodels.Hash, bool, error) {
	return h.createEventRevision(obj, "", creatorID, creator, origin)
}

// createEventRevision is createRevision of an event (Revision.Event): ""
// for an edit — deduplicated by its snapshot —, or what else happened to the
// record, as its deletion — always a revision of its own, its snapshot the
// record as it was, nothing changed in it.
func (h *ModelHistory) createEventRevision(obj interface{}, event string, creatorID uuid.UUID, creator string, origin histmodels.Origin) (histmodels.Hash, bool, error) {
	snap, data, err := h.snapshot(obj)
	if err != nil {
		return nil, false, err
	}
	recordKey := h.mb.MustRecordID(obj).String()
	// Fold the record key into the hash so it is unique per record (sole PK,
	// clean nested route) yet still collapses an unchanged save. An event's
	// folds the event and when it happened too: its own revision.
	prefix := recordKey + "\x00"
	if event != "" {
		prefix += event + "\x00" + strconv.FormatInt(time.Now().UnixNano(), 10) + "\x00"
	}
	sum := sha256.Sum256(append([]byte(prefix), data...))
	hash := histmodels.Hash(sum[:])

	if event == "" {
		var exists int64
		if err = h.db.Table(h.table).
			Where("hash = ? AND record_key = ?", hash, recordKey).
			Count(&exists).Error; err != nil {
			return nil, false, err
		}
		if exists > 0 {
			return hash, false, nil
		}
	}

	parent, err := h.latestHash(recordKey)
	if err != nil {
		return nil, false, err
	}

	// The fields whose value changed from the parent (all of them on the first
	// revision). Kept alongside the full snapshot so the history is queryable by
	// field without walking the chain or re-diffing.
	var changed []string
	if event == "" {
		if changed, err = h.changedFields(recordKey, parent, snap); err != nil {
			return nil, false, err
		}
	}

	rev := histmodels.Revision{
		Hash:          hash,
		RecordKey:     recordKey,
		Parent:        parent,
		Fields:        datatypes.NewJSONType(snap),
		ChangedFields: datatypes.NewJSONType(changed),
		CreatedAt:     time.Now(),
		CreatorID:     creatorID,
		Creator:       creator,
		Origin:        origin,
		Event:         event,
	}
	if err = h.db.Table(h.table).Create(&rev).Error; err != nil {
		return nil, false, err
	}
	return hash, true, nil
}
