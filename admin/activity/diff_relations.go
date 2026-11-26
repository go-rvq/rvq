package activity

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// The relations of a record in its diff. The record of before (findOld) is
// read by its columns only, the one saved carries what its form set — its
// belongs-to, its lists —: compared as they are, every edit would log them
// as changed ("" → the project, "" → the items). So:
//
//   - a belongs-to or a has-one is not diffed: its foreign key, a column, is;
//   - a has-many or a many-to-many the record saved carries is loaded in the
//     one of before (loadOldRelations) and compared by the key of each
//     element: what came, what left and, of a has-many, what changed in one
//     kept — unless its elements are audited on their own (a model of the
//     activity), whose own logs say it;
//   - a relation of an element of a list is not diffed.

// relationsOf are the gorm relations of the type of obj, by field name; nil
// when it is not a model.
func relationsOf(db *gorm.DB, obj any) map[string]*schema.Relationship {
	if db == nil || obj == nil {
		return nil
	}
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(obj); err != nil {
		return nil
	}
	return stmt.Schema.Relationships.Relations
}

// loadOldRelations loads in old the has-many and many-to-many lists that now
// carries (not nil), from the database: findOld read its columns only.
func loadOldRelations(db *gorm.DB, old, now any, rels map[string]*schema.Relationship) {
	ov, nv := reflect.Indirect(reflect.ValueOf(old)), reflect.Indirect(reflect.ValueOf(now))
	for name, rel := range rels {
		if rel.Type != schema.HasMany && rel.Type != schema.Many2Many {
			continue
		}
		nf, of := nv.FieldByName(name), ov.FieldByName(name)
		if !nf.IsValid() || !of.IsValid() || nf.Kind() != reflect.Slice || nf.IsNil() || !of.CanAddr() {
			continue
		}
		// not loaded: the list of before unknown, not diffed (nil)
		if err := db.Session(&gorm.Session{NewDB: true}).Model(old).Association(name).Find(of.Addr().Interface()); err != nil {
			of.Set(reflect.Zero(of.Type()))
		}
	}
}

// relationDiff diffs the relation rel, the field name of the record, of
// before (old) and now: false, not a relation.
func (db *DiffBuilder) relationDiff(name string, old, now reflect.Value) bool {
	rel := db.rels[name]
	if rel == nil {
		return false
	}
	switch rel.Type {
	case schema.HasMany, schema.Many2Many:
	default:
		return true
	}
	if old.Kind() != reflect.Slice || now.Kind() != reflect.Slice || old.IsNil() || now.IsNil() {
		// one not loaded: nothing known of a change
		return true
	}
	if rel.Type == schema.HasMany && db.mb != nil && db.mb.activity != nil {
		if _, ok := db.mb.activity.GetModelBuilder(reflect.New(rel.FieldSchema.ModelType).Interface()); ok {
			return true
		}
	}
	pk := rel.FieldSchema.PrioritizedPrimaryField
	key := func(v reflect.Value) string {
		if pk == nil {
			return ""
		}
		k, zero := pk.ValueOf(db.ctx, reflect.Indirect(v))
		if zero {
			return ""
		}
		return fmt.Sprint(k)
	}
	olds := map[string]reflect.Value{}
	for i := 0; i < old.Len(); i++ {
		if k := key(old.Index(i)); k != "" {
			olds[k] = old.Index(i)
		}
	}
	kept := map[string]bool{}
	for i := 0; i < now.Len(); i++ {
		e := now.Index(i)
		k := key(e)
		o, ok := olds[k]
		if k == "" || !ok {
			db.diffs = append(db.diffs, Diff{Field: name, Now: label(e, k)})
			continue
		}
		kept[k] = true
		if rel.Type == schema.HasMany {
			// what changed in it, its relations not
			prefix := formatFieldByDot(name, label(e, k))
			sub := &DiffBuilder{mb: db.mb, ctx: db.ctx, root: prefix, skip: rel.FieldSchema.Relationships.Relations}
			if err := sub.diffLoop(o, e, prefix); err == nil {
				db.diffs = append(db.diffs, sub.diffs...)
			}
		}
	}
	for i := 0; i < old.Len(); i++ {
		if k := key(old.Index(i)); k != "" && !kept[k] {
			db.diffs = append(db.diffs, Diff{Field: name, Old: label(old.Index(i), k)})
		}
	}
	return true
}

// label is how an element of a list is told in the history: its String, or
// its key; else its value, never its address.
func label(v reflect.Value, key string) string {
	if v.CanInterface() {
		if s, ok := v.Interface().(fmt.Stringer); ok && !(v.Kind() == reflect.Ptr && v.IsNil()) {
			if t := strings.TrimSpace(s.String()); t != "" {
				return t
			}
		}
	}
	if key != "" {
		return key
	}
	return valueString(v)
}

// valueString is the value of v (of what it points to), never an address.
func valueString(v reflect.Value) string {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}
	if !v.CanInterface() {
		return ""
	}
	return fmt.Sprintf("%+v", v.Interface())
}
