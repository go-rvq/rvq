package gorm2op

import (
	"reflect"

	"github.com/qor5/web/v3/zeroer"
	"gorm.io/gorm"
)

type SaveHasManyAssociationBuilder struct {
	field     string
	pre, post []Callback
	updator   func(db *gorm.DB, r any) error
}

func (b *SaveHasManyAssociationBuilder) Updator(f func(db *gorm.DB, r any) (err error)) *SaveHasManyAssociationBuilder {
	b.updator = f
	return b
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

		if v, ok := state.GetOk(b.field); ok {
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
				up := b.updator
				if up == nil {
					up = func(db *gorm.DB, r any) error {
						return db.Updates(r).Error
					}
				}
				for _, value := range oldValues {
					db := state.SharedDB.Session(&gorm.Session{}).Model(value)
					if err = up(db, value); err != nil {
						return
					}
				}
			}
			if len(newValues) > 0 {
				if err = assoc.Append(newValues...); err != nil {
					return
				}
			}
			err = assoc.Unscoped().Replace(all...)
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
