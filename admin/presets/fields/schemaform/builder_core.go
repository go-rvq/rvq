package schemaform

import (
	"github.com/gad-lang/gad"
)

// EnumItem is one value of an enum as the form offers it: the value that goes
// into the record, and what the reader sees.
type EnumItem struct {
	// Name is what goes INTO the record — the enum member's name, never the
	// number behind it.
	Name string
	// Label is what the reader sees; empty shows the Name.
	Label string
	// Hint is the line under the item in the list, when it needs one.
	Hint string
}

// Builder maps a TYPE to the component that edits it, and draws a schema with
// them.
//
// A new Builder knows one type, DefaultType ("str"): a text field, which is
// also what an untyped field is. An application adds its own — and may replace
// that one:
//
//	schemaform.New().
//	    Type("int", myNumberField).
//	    Type("color", myColorPicker)
type Builder struct {
	types     map[string]ComponentFunc
	displays  map[string]ComponentFunc
	readers   map[string]TypeReader
	info      FieldInfoFunc
	enumInfo  EnumInfoFunc
	enumItems EnumItemsFunc
	encode    EncodeFunc
	decode    DecodeFunc
	typeOf    TypeOfFunc
	className ClassNameFunc
	// choiceAsName reads a union of classes as the choice of a class's NAME
	// (ChoiceAsName).
	choiceAsName bool
	// mapsKey is the Google Maps key of the addresses (MapsKey).
	mapsKey MapsKeyFunc
}

// TypeOfFunc says what a field's type is, for a type the schema cannot name by
// itself: an application's own — a model, say — given as the gad object the
// field is typed with. It answers the name of the component the field is
// drawn with, or "" to leave the type to the schema.
type TypeOfFunc func(t gad.Object) string

// TypeOf sets the function that names the types only the application knows
// (TypeOfFunc). It is asked first, for every field of a single type.
// ClassNameFunc names a class a form extends (`*Parent`) — Field.Owner of
// the fields it gives —: "" leaves its own name.
type ClassNameFunc func(c *gad.Class) string

// ClassName names the classes a form extends, for an application whose
// classes are not told apart by their names: forms of its own, each a class
// Form, extended by others (`*forms.quote`: "forms.quote").
func (b *Builder) ClassName(f ClassNameFunc) *Builder {
	b.className = f
	return b
}

func (b *Builder) TypeOf(f TypeOfFunc) *Builder {
	b.typeOf = f
	return b
}

// ChoiceAsName reads a union of classes (`layout? ServiceOptions |
// PortfolioOptions`) as the choice of one class BY ITS NAME — a select of the
// classes, each by its label and hint (the class's metadata; its name
// humanized when it has no label), the value saved its name: "ServiceOptions".
// Without it, the union is a choice of one class WITH ITS FIELDS, saved as
// `{"ServiceOptions": {…}}` (Schema.Choice).
func (b *Builder) ChoiceAsName(v bool) *Builder {
	b.choiceAsName = v
	return b
}

// withLabels shows an item's own name when nothing else was given for it.
func withLabels(items []EnumItem) []EnumItem {
	for i := range items {
		if items[i].Label == "" {
			items[i].Label = items[i].Name
		}
	}
	return items
}

// TypeReader registers how the value of the type name is read back from what
// the form posted (see the type TypeReader): the one logic of a form — the
// admin's, the site's, the browser's — reads it so.
func (b *Builder) TypeReader(name string, f TypeReader) *Builder {
	if b.readers == nil {
		b.readers = map[string]TypeReader{}
	}
	b.readers[name] = f
	return b
}

// Decoder is the form's reader of schema (Decoder) with the types this
// builder reads (TypeReader).
func (b *Builder) Decoder(schema *Schema) *Decoder {
	return &Decoder{Schema: schema, Types: b.readers}
}

// Types are the registered type names, FormType included.
func (b *Builder) Types() (names []string) {
	names = append(names, FormType)
	for name := range b.types {
		names = append(names, name)
	}
	return
}

// The field metadata that names a field where the application gives it no
// words (FieldInfoFunc): `[label="Título", hint="…", help="…",
// placeholder="…"]`.
const (
	MetaLabel       = "label"
	MetaHint        = "hint"
	MetaHelp        = "help"
	MetaPlaceholder = "placeholder"
)

// Type registers (or replaces) the component that draws a field of this type.
func (b *Builder) Type(name string, f ComponentFunc) *Builder {
	b.types[name] = f
	return b
}

// AddressType is the type of an address (`where address`): typed and found by
// Google Places — its suggestions as it is typed —, the place on a map below
// it (vuetifyx's <vx-address-field>). Its value is a record of AddressFields:
// schema.org's PostalAddress, with the coordinates of the place and its id at
// Google.
const AddressType = "address"

// AddressFields are the fields of an address, in their order: the address in
// words (formatted) and its parts, the coordinates (lat, lng: numbers) and the
// place's id at Google (placeId). Typed with no place found, only formatted.
var AddressFields = []string{"formatted", "street", "number", "complement", "neighborhood", "city", "region",
	"postalCode", "country", "countryCode", "lat", "lng", "placeId"}
