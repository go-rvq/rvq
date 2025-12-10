package helper

import (
	"reflect"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/web/zeroer"
)

func InlineEdit(m *presets.ModelBuilder, fieldName string, cb ...func(mb *presets.ModelBuilder)) (childModel *presets.ModelBuilder) {
	field, _ := m.ModelType().Elem().FieldByName(fieldName)
	fieldModel := reflect.New(field.Type.Elem()).Interface()

	childModel = presets.NewModelBuilder(m.Builder(), fieldModel)
	childModel.SetModuleKey(m.ModuleKey())
	InlineEditModel(m, childModel, fieldName)
	for _, f := range cb {
		f(childModel)
	}
	return childModel
}

func InlineEditModel(m, childModel *presets.ModelBuilder, fieldName string) {
	var (
		childEditing = childModel.Editing()
		_, hasID     = childModel.ModelType().Elem().FieldByName("ID")
	)

	if hasID {
		childEditing.Field("ID").SetDisabled(false).ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			v := field.Value()
			if zeroer.IsZero(v) {
				return nil
			}
			return vue.UserComponent().Assign("form", field.FormKey, v)
		})
	}

	detailing := m.Detailing()
	childDetailing := childModel.Detailing()
	if hasID {
		childDetailing.Field("ID").SetDisabled(true)
	}

	detailing.Field(fieldName).Nested(presets.NestedStruct(childModel, &childDetailing.FieldsBuilder))

	m.WithEditingBuilders(func(e *presets.EditingBuilder) {
		e.Field(fieldName).Nested(presets.NestedStruct(childModel, &childEditing.FieldsBuilder))
	})
}
