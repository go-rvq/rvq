package gorm2op

import (
	"errors"
	"net/url"
	"reflect"

	"github.com/go-rvq/rvq/admin/presets"
	gormutils "github.com/go-rvq/rvq/thirdpart/gorm/utils"
	"github.com/go-rvq/rvq/web"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
		auditor := AssociationAuditorFromContext(state.Ctx.R.Context())

		// oldByID snapshots the persisted children of this record in a single
		// query: an existing item of the form must be one of them — its key
		// comes from the form, and a key of a child of another record would
		// update (and Replace re-parent) that child —; it powers the update diff
		// and the deletion log too.
		var oldByID map[any]any
		if oldByID, err = b.snapshotExisting(db, state.Obj, sliceType, pkNames); err != nil {
			return
		}
		// a child of a has-many (has-one) is of one record only: a key of the
		// form that is not of one of its children is of another record's — the
		// update would change it and Replace re-parent it. (The rows of a
		// many-to-many are shared: choosing one by its key is what the list
		// does.)
		if t := assoc.Relationship.Type; t == schema.HasMany || t == schema.HasOne {
			for _, item := range part.Others {
				if !gormutils.PKAllZero(item, pkNames) {
					if _, ok := oldByID[gormutils.PKMapKey(item, pkNames)]; !ok {
						return ErrNotAChild
					}
				}
			}
			for _, item := range part.New {
				var taken bool
				if taken, err = keyOfAnotherRecord(db, assoc.Relationship.FieldSchema, item, oldByID, pkNames); err != nil {
					return
				}
				if taken {
					return ErrNotAChild
				}
			}
		}
		// the nested lists of an item are saved with it (gorm upserts them, the
		// foreign key included): a nested row of another record would be
		// re-parented
		for _, items := range [][]reflect.Value{part.Others, part.New} {
			for _, item := range items {
				if err = nestedOfAnotherRecord(db, assoc.Relationship.FieldSchema, item); err != nil {
					return
				}
			}
		}
		kept := part.KeptSlice(sliceType)

		if auditor != nil {
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
			// a row of a many-to-many not linked to this record yet is shared:
			// it is linked (Replace), its columns are not the form's to change
			if _, linked := oldByID[gormutils.PKMapKey(item, pkNames)]; !linked {
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

// ErrNotAChild: an item of the list of a form (a has-many) whose
// key is not of a child of the record saved: refused.
var ErrNotAChild = errors.New("gorm2op: an item of the list is not a child of the record saved")

// pkWhere is the condition of the primary key of rv (of schema sch), every
// field — a zero one too: it may be a value of the key —; set tells whether
// any field is set (none: a row to create).
func pkWhere(db *gorm.DB, sch *schema.Schema, rv reflect.Value) (where map[string]any, set bool) {
	where = make(map[string]any, len(sch.PrimaryFields))
	for _, f := range sch.PrimaryFields {
		v, zero := f.ValueOf(db.Statement.Context, rv)
		where[f.DBName] = v
		set = set || !zero
	}
	return
}

// keyOfAnotherRecord tells whether the new item carries the primary key of a
// persisted row that is not one of oldByID (the children of the
// record saved): creating it would upsert that row.
func keyOfAnotherRecord(db *gorm.DB, sch *schema.Schema, item reflect.Value, oldByID map[any]any, pkNames []string) (bool, error) {
	rv := indirectValue(item)
	if !rv.IsValid() || rv.Kind() != reflect.Struct {
		return false, nil
	}
	where, set := pkWhere(db, sch, rv)
	if !set {
		return false, nil
	}
	if _, ok := oldByID[gormutils.PKMapKey(item, pkNames)]; ok {
		return false, nil
	}
	var n int64
	err := db.Session(&gorm.Session{NewDB: true}).Unscoped().Table(sch.Table).Where(where).Count(&n).Error
	return n > 0, err
}

// nestedOfAnotherRecord walks the has-many and has-one lists of the item (of
// schema sch), deep: a nested row whose primary key (a field of it set) is of
// a persisted row must be a child of its parent item already — ErrNotAChild
// otherwise (a new parent has no children yet). A row of a key not persisted
// is created.
func nestedOfAnotherRecord(db *gorm.DB, sch *schema.Schema, item reflect.Value) error {
	rv := indirectValue(item)
	if !rv.IsValid() || rv.Kind() != reflect.Struct {
		return nil
	}
	// the schema of the item's own type: the relations of a relationship's
	// FieldSchema may hold the owner's
	if sch == nil || sch.ModelType != rv.Type() {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(reflect.New(rv.Type()).Interface()); err != nil {
			return err
		}
		sch = stmt.Schema
	}
	ctx := db.Statement.Context
	for _, rel := range sch.Relationships.Relations {
		if rel.Field.Schema != sch && rel.Field.Schema != nil && rel.Field.Schema.ModelType != rv.Type() {
			continue
		}
		if (rel.Type != schema.HasMany && rel.Type != schema.HasOne) || len(rel.FieldSchema.PrimaryFields) == 0 {
			continue
		}
		fv := indirectValue(rel.Field.ReflectValueOf(ctx, rv))
		if !fv.IsValid() {
			continue
		}
		var children []reflect.Value
		if fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array {
			for i := 0; i < fv.Len(); i++ {
				children = append(children, fv.Index(i))
			}
		} else {
			children = append(children, fv)
		}
		for _, child := range children {
			cv := indirectValue(child)
			if !cv.IsValid() || cv.Kind() != reflect.Struct {
				continue
			}
			if pk, set := pkWhere(db, rel.FieldSchema, cv); set {
				of := make(map[string]any, len(pk)+len(rel.References))
				for k, v := range pk {
					of[k] = v
				}
				for _, ref := range rel.References {
					switch {
					case ref.OwnPrimaryKey:
						v, _ := ref.PrimaryKey.ValueOf(ctx, rv)
						of[ref.ForeignKey.DBName] = v
					case ref.PrimaryValue != "":
						of[ref.ForeignKey.DBName] = ref.PrimaryValue
					}
				}
				var persisted, ofParent int64
				q := func() *gorm.DB {
					return db.Session(&gorm.Session{NewDB: true}).Unscoped().Table(rel.FieldSchema.Table)
				}
				if err := q().Where(pk).Count(&persisted).Error; err != nil {
					return err
				}
				if persisted > 0 {
					if err := q().Where(of).Count(&ofParent).Error; err != nil {
						return err
					}
					if ofParent == 0 {
						return ErrNotAChild
					}
				}
			}
			if err := nestedOfAnotherRecord(db, rel.FieldSchema, cv); err != nil {
				return err
			}
		}
	}
	return nil
}

// indirectValue dereferences pointers and interfaces: the invalid Value for a
// nil one.
func indirectValue(v reflect.Value) reflect.Value {
	for v.IsValid() && (v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}
