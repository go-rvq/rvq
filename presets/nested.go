package presets

type Nested interface {
	FieldWalker
	Model() *ModelBuilder
	FieldsBuilder() *FieldsBuilder
	Build(b *FieldBuilder)
}
