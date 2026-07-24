package gorm2op

import (
	"net/url"
	"reflect"

	"github.com/go-rvq/rvq/admin/presets"
	gormutils "github.com/go-rvq/rvq/thirdpart/gorm/utils"
	"github.com/go-rvq/rvq/web"
	"gorm.io/gorm"
)

type SaveHasManyAssociationBuilder struct {
	field        string
	fieldFormKey func(fieldName string, state *CallbackState) string
	pre, post    []Callback
	updator      func(db *gorm.DB, r any) error
}

func (b *SaveHasManyAssociationBuilder) Updator(f func(db *gorm.DB, r any) (err error)) *SaveHasManyAssociationBuilder {
	b.updator = f
	return b
}

func SaveHasManyAssociation(field string, fieldFormKey ...func(fieldName string, state *CallbackState) string) *SaveHasManyAssociationBuilder {
	b := &SaveHasManyAssociationBuilder{field: field}
	for _, b.fieldFormKey = range fieldFormKey {
	}
	if b.fieldFormKey == nil {
		b.fieldFormKey = func(fieldName string, state *CallbackState) string {
			return fieldName
		}
	}
	return b
}

func (b *SaveHasManyAssociationBuilder) Pre(f ...Callback) *SaveHasManyAssociationBuilder {
	b.pre = append(b.pre, f...)
	return b
}

func (b *SaveHasManyAssociationBuilder) Post(f ...Callback) *SaveHasManyAssociationBuilder {
	b.post = append(b.post, f...)
	return b
}

// assocPresentKey namespaces the "was the field in the form" flag stored on the
// callback state between the pre and post callbacks.
type assocPresentKey string

func (b *SaveHasManyAssociationBuilder) Build(ob *DataOperatorBuilder) *DataOperatorBuilder {
	// pre detaches the submitted slice from the object (so the base create/update
	// does not cascade-save it) and remembers whether the list was present in the
	// form, so post can act on the user's intent.
	pre := func(state *CallbackState) (err error) {
		fieldFormKey := b.fieldFormKey(b.field, state)
		present := presets.ListEditorInitialized(formValues(state.Ctx), fieldFormKey)
		state.Set(assocPresentKey(b.field), present)

		value := reflect.ValueOf(state.Obj).Elem().FieldByName(b.field)
		if present {
			state.Set(b.field, value.Interface())
		}
		value.Set(reflect.Zero(value.Type()))
		return
	}

	// post reconciles the persisted children with the submitted list:
	//   - items flagged __deleted (or missing from the list) are removed;
	//   - items flagged __new are created;
	//   - the remaining (existing) items are updated.
	// Classification is delegated to presets.PartitionListEditorItems, which uses
	// only the __deleted/__new form flags (never the primary key) so the list
	// editor works for element types without an ID. When the list was present but
	// empty every child is removed. When an AssociationAuditor is available in the
	// context (the parent model is audited), each create/update/delete is logged.
	post := func(state *CallbackState) (err error) {
		present, _ := state.Get(assocPresentKey(b.field)).(bool)
		if !present {
			return
		}

		v, ok := state.GetOk(b.field)
		if !ok {
			return
		}

		var (
			db           = state.SharedDB.Session(&gorm.Session{})
			assoc        = db.Model(state.Obj).Association(b.field)
			fieldFormKey = b.fieldFormKey(b.field, state)
		)
		if assoc.Error != nil {
			return assoc.Error
		}
		if assoc.Relationship == nil || len(assoc.Relationship.FieldSchema.PrimaryFields) == 0 {
			return
		}

		var (
			items     = reflect.ValueOf(v)
			sliceType = items.Type()
			// support composite primary keys: every primary field participates in
			// the key used to match/skip rows, not just the first.
			pkNames = gormutils.PrimaryFieldNames(assoc.Relationship.FieldSchema)
		)

		// part classifies the submitted items by their __deleted/__new flags, in
		// __pos order. A nil part means the list editor was not initialized for this
		// field, so its children are left untouched.
		part := presets.PartitionListEditorItems(formValues(state.Ctx), fieldFormKey, items)
		if part == nil {
			return
		}
		kept := part.KeptSlice(sliceType)

		auditor := AssociationAuditorFromContext(state.Ctx.R.Context())

		// oldByID snapshots the persisted children in a single query. It powers
		// both the update diff and the deletion log, and is only fetched when an
		// auditor is present (otherwise the reads are pure overhead). The primary
		// key is only consulted here, for auditing — never for classification.
		var oldByID map[any]any
		if auditor != nil {
			if oldByID, err = b.snapshotExisting(db, state.Obj, sliceType, pkNames); err != nil {
				return
			}
			// keptIDs: primary keys of the existing children we keep, used to
			// detect which persisted rows are being removed.
			keptIDs := make(map[any]bool, len(part.Kept))
			for _, item := range part.Kept {
				if !gormutils.PKAllZero(item, pkNames) {
					keptIDs[gormutils.PKMapKey(item, pkNames)] = true
				}
			}
			for id, old := range oldByID {
				if !keptIDs[id] {
					if err = auditor.LogDeleted(db, old); err != nil {
						return
					}
				}
			}
		}

		// persist field changes on existing children.
		up := b.updator
		if up == nil {
			up = func(db *gorm.DB, r any) error {
				return db.Updates(r).Error
			}
		}
		for _, item := range part.Others {
			// an "existing" row must have a primary key to update; without one there
			// is nothing to target (gorm would raise "WHERE conditions required").
			// For a composite key every primary field must be set. Leave a keyless
			// row for Replace to create — a safety net for rows that reach Others
			// without a PK (e.g. a client that omitted __new).
			if gormutils.PKAllZero(item, pkNames) {
				continue
			}
			rec := item.Interface()
			if err = up(db.Session(&gorm.Session{}).Model(rec), rec); err != nil {
				return
			}
			if auditor != nil {
				if old, okOld := oldByID[gormutils.PKMapKey(item, pkNames)]; okOld {
					if err = auditor.LogUpdated(db, old, rec); err != nil {
						return
					}
				}
			}
		}

		// a single Replace creates the new children, keeps the updated ones and
		// deletes every persisted row not in kept (removed and empty-list cases).
		if err = db.Model(state.Obj).Association(b.field).Unscoped().Replace(kept.Interface()); err != nil {
			return
		}

		if auditor != nil {
			for _, item := range part.New {
				if err = auditor.LogCreated(db, item.Interface()); err != nil {
					return
				}
			}
		}
		return
	}

	return ob.
		CreateCallbacks().
		Pre(b.pre...).Pre(pre).
		Post(b.post...).Post(post).
		Dot().
		UpdateCallbacks().
		Pre(b.pre...).Pre(pre).
		Post(b.post...).Post(post).
		Dot()
}

// snapshotExisting loads the persisted children of the association into a map
// keyed by their (possibly composite) primary key, in one query.
func (b *SaveHasManyAssociationBuilder) snapshotExisting(db *gorm.DB, obj any, sliceType reflect.Type, pkNames []string) (map[any]any, error) {
	existingPtr := reflect.New(sliceType).Interface()
	if err := db.Model(obj).Association(b.field).Find(existingPtr); err != nil {
		return nil, err
	}
	existing := reflect.ValueOf(existingPtr).Elem()
	out := make(map[any]any, existing.Len())
	for i := 0; i < existing.Len(); i++ {
		row := existing.Index(i)
		out[gormutils.PKMapKey(row, pkNames)] = row.Interface()
	}
	return out, nil
}

// formValues returns the (reordered) submitted form values. UnmarshalForm
// reindexes MultipartForm.Value by __pos, so item indexes there line up with the
// decoded slice order. It delegates to presets.ListEditorFormValues so the
// persistence layer reads from the very same source as the list-editor render.
func formValues(ctx *web.EventContext) url.Values {
	return presets.ListEditorFormValues(ctx)
}
