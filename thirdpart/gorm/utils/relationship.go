package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-rvq/rvq/web/zeroer"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// FieldValuer reads a value by (Go) field name. model.ID satisfies it, so a
// record id can bind a parent filter without this package importing that type.
type FieldValuer interface {
	GetValue(fieldName string) any
}

// ForeignKeyFields resolves, from the gorm relationship of a belongs-to field,
// the owner's foreign-key field names in related-primary-key order. It supports
// composite foreign keys (two or more columns). Returns nil when the relationship
// cannot be resolved.
func ForeignKeyFields(db *gorm.DB, model any, field string) []string {
	rel := db.Model(model).Association(field).Relationship
	if rel == nil || len(rel.References) == 0 {
		return nil
	}
	byRelatedPK := make(map[string]string, len(rel.References))
	for _, ref := range rel.References {
		byRelatedPK[ref.PrimaryKey.Name] = ref.ForeignKey.Name
	}
	out := make([]string, 0, len(rel.FieldSchema.PrimaryFields))
	for _, pf := range rel.FieldSchema.PrimaryFields {
		if fk, ok := byRelatedPK[pf.Name]; ok {
			out = append(out, fk)
		}
	}
	return out
}

// HasManyParentFilter builds the WHERE matching a has-many child by its parent
// ("<fk> = ? [AND ...]") and the owner primary-key field names in placeholder
// order. Every reference participates, so composite foreign keys (two or more
// fields) work.
func HasManyParentFilter(rel *schema.Relationship) (sql string, ownerFields []string) {
	conds := make([]string, 0, len(rel.References))
	for _, ref := range rel.References {
		conds = append(conds, ref.ForeignKey.DBName+" = ?")
		ownerFields = append(ownerFields, ref.PrimaryKey.Name)
	}
	return strings.Join(conds, " AND "), ownerFields
}

// M2MParentFilter builds an EXISTS(...) correlated subquery matching the related
// rows linked to a parent through the join table, and the owner primary-key field
// names in placeholder order. Owner-side join columns are bound to the parent id;
// related-side columns are correlated to the related table (aliased "j" for the
// join). Handles single and composite keys of any type — no fixed SQL cast.
func M2MParentFilter(rel *schema.Relationship) (sql string, ownerFields []string) {
	conds := make([]string, 0, len(rel.References))
	for _, ref := range rel.References {
		if ref.OwnPrimaryKey {
			conds = append(conds, fmt.Sprintf("j.%s = ?", ref.ForeignKey.DBName))
			ownerFields = append(ownerFields, ref.PrimaryKey.Name)
		} else {
			conds = append(conds, fmt.Sprintf("j.%s = %s.%s", ref.ForeignKey.DBName, rel.FieldSchema.Table, ref.PrimaryKey.DBName))
		}
	}
	return fmt.Sprintf("EXISTS (SELECT 1 FROM %s j WHERE %s)", rel.JoinTable.Table, strings.Join(conds, " AND ")), ownerFields
}

// ParentFilterArgs extracts the owner primary-key values (in fieldNames order)
// from a parent id, to bind the placeholders of a parent filter. Supports single
// and composite keys.
func ParentFilterArgs(id FieldValuer, fieldNames []string) []any {
	args := make([]any, len(fieldNames))
	for i, name := range fieldNames {
		args[i] = id.GetValue(name)
	}
	return args
}

// PrimaryFieldNames returns the Go field names of every primary-key field of a
// schema, so composite keys are handled and not just the first field.
func PrimaryFieldNames(s *schema.Schema) []string {
	names := make([]string, len(s.PrimaryFields))
	for i, f := range s.PrimaryFields {
		names[i] = f.Name
	}
	return names
}

// PKAllZero reports whether every primary-key field of item is zero, i.e. the row
// has no usable key to target an update (a new / keyless row).
func PKAllZero(item reflect.Value, pkNames []string) bool {
	for _, name := range pkNames {
		if !zeroer.IsZero(reflect.Indirect(item).FieldByName(name)) {
			return false
		}
	}
	return true
}

// PKMapKey builds a comparable map key from item's primary-key field(s),
// supporting composite keys unambiguously.
func PKMapKey(item reflect.Value, pkNames []string) any {
	if len(pkNames) == 1 {
		return NormalizeKey(reflect.Indirect(item).FieldByName(pkNames[0]))
	}
	parts := make([]any, len(pkNames))
	for i, name := range pkNames {
		parts[i] = NormalizeKey(reflect.Indirect(item).FieldByName(name))
	}
	return fmt.Sprintf("%v", parts)
}

// NormalizeKey makes a primary-key value comparable as a map key regardless of
// its concrete integer/uint/string type.
func NormalizeKey(v reflect.Value) any {
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
