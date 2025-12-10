package nested

import (
	"reflect"

	"github.com/go-rvq/rvq/admin/presets"
)

type FieldMode = presets.FieldMode

type Nested struct {
	mb *presets.ModelBuilder
}

type FieldBuilderHandle func(*presets.FieldsBuilder) *presets.FieldsBuilder

func New(mb *presets.ModelBuilder) *Nested {
	return &Nested{mb: mb}
}

func (n *Nested) Field(name string) *NestedFieldBuilder {
	nestedModel := reflect.ValueOf(n.mb.Model()).Elem().FieldByName(name).Interface()

	if rt := reflect.TypeOf(nestedModel); rt.Kind() != reflect.Ptr {
		nestedModel = reflect.New(rt).Interface()
	}

	fieldMb := presets.NewModelBuilder(n.mb.Builder(), nestedModel)
	fieldMb.SetModuleKey(n.mb.ModuleKey())
	return &NestedFieldBuilder{n.mb, fieldMb, name}
}

type NestedFieldBuilder struct {
	baseModel  *presets.ModelBuilder
	fieldModel *presets.ModelBuilder
	field      string
}

func (b *NestedFieldBuilder) Editing(cb ...FieldBuilderHandle) {
	fb := &b.fieldModel.Editing().FieldsBuilder
	b.baseModel.Editing().Field(b.field).AutoNested(b.fieldModel, fb)
	for _, cb := range cb {
		cb(fb)
	}
}

func (b *NestedFieldBuilder) Detailing(cb ...FieldBuilderHandle) {
	fb := &b.fieldModel.Detailing().FieldsBuilder
	b.baseModel.Detailing().Field(b.field).AutoNested(b.fieldModel, fb)
	for _, cb := range cb {
		cb(fb)
	}
}

func (b *NestedFieldBuilder) Listing(cb ...FieldBuilderHandle) {
	fb := &b.fieldModel.Listing().FieldsBuilder
	b.baseModel.Listing().Field(b.field).AutoNested(b.fieldModel, fb)
	for _, cb := range cb {
		cb(fb)
	}
}
