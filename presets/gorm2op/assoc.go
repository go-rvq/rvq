package gorm2op

import (
	"reflect"

	"github.com/qor5/x/v3/zeroer"
	"gorm.io/gorm"
)

type SaveHasManyAssociationBuilder struct {
	field     string
	pre, post []Callback
}

func SaveHasManyAssociation(field string) *SaveHasManyAssociationBuilder {
	return &SaveHasManyAssociationBuilder{field: field}
}

func (b *SaveHasManyAssociationBuilder) Pre(f ...Callback) *SaveHasManyAssociationBuilder {
	b.pre = append(b.pre, f...)
	return b
}

func (b *SaveHasManyAssociationBuilder) Post(f ...Callback) *SaveHasManyAssociationBuilder {
	b.post = append(b.post, f...)
	return b
}

func (b *SaveHasManyAssociationBuilder) Build(ob *DataOperatorBuilder) *DataOperatorBuilder {
	pre := func(state *CallbackState) (err error) {
		value := reflect.ValueOf(state.Obj).Elem().FieldByName(b.field)
		if value.Len() == 0 {
			return
		}
		valueI := value.Interface()
		state.Set(b.field, valueI)
		value.Set(reflect.Zero(value.Type()))
		return
	}
	post := func(state *CallbackState) (err error) {
		assoc := state.SharedDB.Session(&gorm.Session{}).Model(state.Obj).Association(b.field)

		if v, ok := state.Get(b.field); ok {
			var (
				items     = reflect.ValueOf(v)
				all       = make([]any, items.Len())
				newValues []any
				oldValues []any
			)

			pkFieldName := assoc.Relationship.Schema.PrimaryFields[0].Name

			for i := 0; i < items.Len(); i++ {
				item := items.Index(i)
				itemObject := item.Interface()

				if zeroer.IsZero(item.Elem().FieldByName(pkFieldName)) {
					newValues = append(newValues, itemObject)
				} else {
					oldValues = append(oldValues, itemObject)
				}
				all[i] = itemObject
			}

			if len(oldValues) > 0 {
				for _, value := range oldValues {
					db := state.SharedDB.Session(&gorm.Session{}).Model(value)
					if err = db.Updates(value).Error; err != nil {
						return
					}
				}
			}
			if len(newValues) > 0 {
				if err = assoc.Append(newValues); err != nil {
					return
				}
			}
			err = assoc.Unscoped().Replace(all...)
		}
		return
	}

	return ob.
		PreCreate(b.pre...).
		PreCreate(pre).
		PostCreate(post).
		PostCreate(b.post...).
		PreUpdate(b.pre...).
		PreUpdate(pre).
		PostUpdate(post).
		PostUpdate(b.post...)
}
