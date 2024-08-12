package presets

type NestedFieldBuilder struct {
	mb *ModelBuilder
	*FieldsBuilder
}

func NewNestedFieldBuilder(mb *ModelBuilder, fieldsBuilder *FieldsBuilder) *NestedFieldBuilder {
	return &NestedFieldBuilder{mb: mb, FieldsBuilder: fieldsBuilder}
}
