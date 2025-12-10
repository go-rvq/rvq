package helper

import (
	"fmt"
	"reflect"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type AssociationManagerBuilder struct {
	Parent       *presets.ModelBuilder
	Field        string
	configurer   func(mb *presets.ModelBuilder)
	dataOperator func(do *gorm2op.DataOperatorBuilder) *gorm2op.DataOperatorBuilder
}

func (b *AssociationManagerBuilder) DataOperator() func(do *gorm2op.DataOperatorBuilder) *gorm2op.DataOperatorBuilder {
	return b.dataOperator
}

func (b *AssociationManagerBuilder) SetDataOperator(dataOperator func(do *gorm2op.DataOperatorBuilder) *gorm2op.DataOperatorBuilder) *AssociationManagerBuilder {
	b.dataOperator = dataOperator
	return b
}

func (b *AssociationManagerBuilder) Configurer() func(mb *presets.ModelBuilder) {
	return b.configurer
}

func (b *AssociationManagerBuilder) SetConfigurer(configurer func(mb *presets.ModelBuilder)) *AssociationManagerBuilder {
	b.configurer = configurer
	return b
}

func NewAssociationManagerBuilder(parent *presets.ModelBuilder, field string) *AssociationManagerBuilder {
	return &AssociationManagerBuilder{Parent: parent, Field: field}
}

func (b *AssociationManagerBuilder) Build() {
	var (
		do           = b.Parent.CurrentDataOperator().(*gorm2op.DataOperatorBuilder).Clone()
		assoc        = do.DB().Model(b.Parent.Model()).Association(b.Field)
		parentColumn string
		childField   string
	)

	for _, r := range assoc.Relationship.References {
		if r.OwnPrimaryKey {
			parentColumn = r.ForeignKey.DBName
		} else {
			childField = r.ForeignKey.DBName
		}
	}

	switch assoc.Relationship.Type {
	case schema.HasMany:
		tbName := assoc.Relationship.FieldSchema.Table
		query := fmt.Sprintf("%s.%s = ?", tbName, parentColumn)
		do.
			SetPreparer(func(db *gorm.DB, mode gorm2op.Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
				if id.IsZero() {
					parentsID := presets.ParentsModelID(ctx.R)
					params.Where(query, parentsID.Last().Value())
				}

				return db
			})
	case schema.Many2Many:
		tbName := assoc.Relationship.JoinTable.Table
		query := fmt.Sprintf("EXISTS (SELECT 1 FROM %s rel WHERE rel.%s = ? AND rel.%s = %s.id)",
			tbName,
			parentColumn,
			childField,
			assoc.Relationship.FieldSchema.Table,
		)

		do.
			SetPreparer(func(db *gorm.DB, mode gorm2op.Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
				parentsID := presets.ParentsModelID(ctx.R)
				params.Where(query, parentsID.Last().Value())
				return db
			})
	}

	do.
		SetCreator(func(db *gorm.DB, obj interface{}, ctx *web.EventContext) (err error) {
			parentsID := presets.ParentsModelID(ctx.R)
			parent := b.Parent.NewModel()
			parentsID.Last().SetTo(parent)

			return db.Model(parent).Association(b.Field).Append(obj)
		})

	do.UpdateCallbacks().
		Pre(func(state *gorm2op.CallbackState) (err error) {
			// disable parent column ref updates
			state.DB.Omit(parentColumn)
			return nil
		})

	if b.dataOperator != nil {
		do = b.dataOperator(do)
	}

	modelType := assoc.Relationship.Field.IndirectFieldType.Elem()
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	model := reflect.New(modelType).Interface()
	Child := presets.NewModelBuilder(b.Parent.Builder(), model,
		presets.ModelConfig().
			SetId(b.Field).
			SetUriName(strcase.ToKebab(b.Field)).
			SetDataOperator(do),
	)

	b.Parent.AddChildH(Child, func(mb *presets.ModelBuilder) {
		if b.configurer != nil {
			b.configurer(mb)
		}
	})
}
