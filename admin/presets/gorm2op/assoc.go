package gorm2op

import (
	"net/url"
	"reflect"
	"strconv"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/zeroer"
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
		present := sliceSubmittedInForm(state.Ctx, fieldFormKey)
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
	//   - items with a zero primary key are created;
	//   - the remaining items are updated.
	// When the list was present but empty every child is removed. When an
	// AssociationAuditor is available in the context (the parent model is
	// audited), each create/update/delete is logged.
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
			pkName    = assoc.Relationship.FieldSchema.PrimaryFields[0].Name

			// kept: the desired final child set (new + updated), passed to Replace.
			kept = reflect.MakeSlice(sliceType, 0, items.Len())
			// updates: existing children whose fields must be persisted.
			updates []any
			// newItems: children to be created (zero pk); logged after Replace
			// assigns their ids.
			newItems []any
			// keptIDs: primary keys of the existing children we keep, used to
			// detect which persisted rows are being removed.
			keptIDs = map[any]bool{}
		)

		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			if itemFlag(state.Ctx, fieldFormKey, i, presets.ListEditorDeletedField) {
				// deleted rows post only {ID, __deleted}; drop them from the kept
				// set so Replace removes them.
				continue
			}
			kept = reflect.Append(kept, item)

			pkVal := reflect.Indirect(item).FieldByName(pkName)
			if zeroer.IsZero(pkVal) || itemFlag(state.Ctx, fieldFormKey, i, presets.ListEditorNewField) {
				newItems = append(newItems, item.Interface())
			} else {
				keptIDs[normalizeKey(pkVal)] = true
				updates = append(updates, item.Interface())
			}
		}

		auditor := AssociationAuditorFromContext(state.Ctx.R.Context())

		// oldByID snapshots the persisted children in a single query. It powers
		// both the update diff and the deletion log, and is only fetched when an
		// auditor is present (otherwise the reads are pure overhead).
		var oldByID map[any]any
		if auditor != nil {
			if oldByID, err = b.snapshotExisting(db, state.Obj, sliceType, pkName); err != nil {
				return
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
		for _, rec := range updates {
			if err = up(db.Session(&gorm.Session{}).Model(rec), rec); err != nil {
				return
			}
			if auditor != nil {
				if old, okOld := oldByID[normalizeKey(reflect.Indirect(reflect.ValueOf(rec)).FieldByName(pkName))]; okOld {
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
			for _, rec := range newItems {
				if err = auditor.LogCreated(db, rec); err != nil {
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
// keyed by their primary key, in one query.
func (b *SaveHasManyAssociationBuilder) snapshotExisting(db *gorm.DB, obj any, sliceType reflect.Type, pkName string) (map[any]any, error) {
	existingPtr := reflect.New(sliceType).Interface()
	if err := db.Model(obj).Association(b.field).Find(existingPtr); err != nil {
		return nil, err
	}
	existing := reflect.ValueOf(existingPtr).Elem()
	out := make(map[any]any, existing.Len())
	for i := 0; i < existing.Len(); i++ {
		row := existing.Index(i)
		out[normalizeKey(reflect.Indirect(row).FieldByName(pkName))] = row.Interface()
	}
	return out, nil
}

// formValues returns the (reordered) submitted form values. UnmarshalForm
// reindexes MultipartForm.Value by __pos, so item indexes there line up with the
// decoded slice order.
func formValues(ctx *web.EventContext) url.Values {
	if ctx != nil && ctx.R != nil {
		if ctx.R.MultipartForm != nil && ctx.R.MultipartForm.Value != nil {
			return ctx.R.MultipartForm.Value
		}
		if ctx.R.Form != nil {
			return ctx.R.Form
		}
	}
	return nil
}

// sliceSubmittedInForm reports whether the list field was part of the submitted
// form, even when it is empty. It honours the explicit presence marker and also
// falls back to detecting any item key, so it works with older forms.
func sliceSubmittedInForm(ctx *web.EventContext, fieldFormKey string) bool {
	values := formValues(ctx)
	if values == nil {
		return false
	}
	if _, ok := values[fieldFormKey+"."+presets.ListEditorPresentField]; ok {
		return true
	}
	prefix := fieldFormKey + "["
	for k := range values {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// itemFlag reports whether the boolean metadata field (e.g. __deleted, __new) is
// truthy for item i of the list.
func itemFlag(ctx *web.EventContext, fieldFormKey string, i int, field string) bool {
	values := formValues(ctx)
	if values == nil {
		return false
	}
	key := fieldFormKey + "[" + strconv.Itoa(i) + "]." + field
	switch values.Get(key) {
	case "true", "1", "on":
		return true
	}
	return false
}

// normalizeKey makes primary key values comparable as map keys regardless of
// their concrete integer/uint/string type.
func normalizeKey(v reflect.Value) any {
	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint())
	case reflect.String:
		return v.String()
	default:
		return v.Interface()
	}
}
