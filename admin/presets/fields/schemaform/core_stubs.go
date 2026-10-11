//go:build formscore

package schemaform

// The build of the schema alone (tag formscore): what reads a form's schema
// and checks its values — Parse, Class, the groups —, with none of what
// draws it (presets, the components). What the Builder holds of the drawing
// is here with no body: its types known by their names only.

// ComponentFunc is, in this build, a type known by its name: nothing draws.
type ComponentFunc func(c *Context) any

// Context is, in this build, the field alone.
type Context struct {
	Field *Field
}

// The functions of the words and the values of a field: none in this build.
type (
	FieldInfoFunc func()
	EnumInfoFunc  func()
	EnumItemsFunc func()
	MapsKeyFunc   func()
)

// New is a Builder that knows the types of the drawing build, by name.
func New() *Builder {
	types := map[string]ComponentFunc{}
	for _, name := range []string{DefaultType, "int", "uint", "bool", "color", "float", "decimal", "text", "html",
		"time", "date", "calendarDate", "calendarTime", RangeType, "duration", FileType, ImageType, AddressType} {
		types[name] = nil
	}
	return &Builder{types: types, displays: map[string]ComponentFunc{}}
}

var TextComponentFunc ComponentFunc
