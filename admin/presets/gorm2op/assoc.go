package gorm2op

import (
	"reflect"

	"github.com/go-rvq/rvq/web/zeroer"
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
				newValues = reflect.MakeSlice(reflect.TypeOf(v), 0, 0)
				oldValues = reflect.MakeSlice(reflect.TypeOf(v), 0, 0)
			)
			reflect.ValueOf(state.Obj).Elem().FieldByName(b.field).Set(reflect.MakeSlice(reflect.TypeOf(v), 0, 0))

			pkFieldName := assoc.Relationship.Schema.PrimaryFields[0].Name

			for i := 0; i < items.Len(); i++ {
				item := items.Index(i)

				if zeroer.IsZero(item.Elem().FieldByName(pkFieldName)) {
					newValues = reflect.Append(newValues, item)
				} else {
					oldValues = reflect.Append(oldValues, item)
				}
			}

			if oldValues.Len() > 0 {
				up := b.updator
				if up == nil {
					up = func(db *gorm.DB, r any) error {
						return db.Updates(r).Error
					}
				}

				for l, i := oldValues.Len(), 0; i < l; i++ {
					v := oldValues.Index(i).Interface()
					db := state.SharedDB.Session(&gorm.Session{}).Model(v)
					if err = up(db, v); err != nil {
						return
					}
				}
				if err = assoc.Unscoped().Replace(oldValues.Interface()); err != nil {
					return
				}
			}

			if newValues.Len() > 0 {
				db := state.SharedDB.Session(&gorm.Session{}).Model(state.Obj).Association(b.field)
				if err = db.Append(newValues.Interface()); err != nil {
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
