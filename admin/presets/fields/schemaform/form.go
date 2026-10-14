package schemaform

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gad-lang/gad"
)

// The one logic of a form, wherever it is drawn: the values a form posts
// read back by its schema — each field typed and checked (left empty when
// required, its type, its limits, its mask and pattern, one of its values,
// a range in order, a file by its field, a list by how many) and its
// `[validation=…]` called —, what is wrong said field by field (Issues, by
// path). The site and the admin read with it, each its own way of drawing
// and keeping; the browser too (forms/wasm, the build formscore).

// Issue is what is wrong with the value at Path ("a.b", a list's record
// "links[1].href"; "" the form whole): a code (CodeRequired…), and the words
// of a `[validation=…]` that said why (Message).
type Issue struct {
	Path    string `json:"path"`
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	// Field is the field of the value; Value, the value as posted; Err, the
	// error of the application's (Items) the value could not be checked by:
	// what a message is said with.
	Field *Field `json:"-"`
	Value string `json:"-"`
	Err   error  `json:"-"`
}

func (e Issue) Error() string { return e.Path + ": " + e.Code }

// More codes of an Issue (CodeRequired… in check.go): not one of its values;
// a file of a type its field does not take, too large, no image; a list of
// fewer, of more items than it takes; a `[validation=…]` that said why.
const (
	CodeChoice     = "choice"
	CodeFileType   = "file_type"
	CodeFileSize   = "file_size"
	CodeImage      = "image"
	CodeItemsMin   = "items_min"
	CodeItemsMax   = "items_max"
	CodeValidation = "validation"
	// CodeItems: the values the field may hold could not be fetched (Err);
	// a value never checked is not one that passed.
	CodeItems = "items"
)

// The longest a value is when its field says none ([maxlength=…]): a line,
// a text.
const (
	MaxLine = 1000
	MaxText = 20000
)

// TypeReader reads a value of a type of the application's (Decoder.Types) —
// raw, as posted, at path —: the value, or what is wrong with it (a code);
// given only a value not empty (Blank).
type TypeReader func(f *Field, path string, raw any) (any, string)

// FileSource is what the files of a form are (Decoder.Files): the file the
// value raw of the file field f names — checked by f: its size, its type —,
// as the value holds it, or what is wrong with it (a code). raw is never
// empty.
type FileSource interface {
	File(f *Field, raw any) (any, string)
}

// Decoder reads the values a form posts back by its schema, and checks them:
// the one logic of a form (see Issue).
type Decoder struct {
	Schema *Schema
	// VM calls the `[validation=…]` functions (the schema's own: Builder.Class
	// keeps it, when nil); Ctx bounds them.
	VM  *gad.VM
	Ctx context.Context
	// Types read the values of the application's types (a phone, an e-mail);
	// one of the builder's own (DefaultType…) is read as such.
	Types map[string]TypeReader
	// Files are the files posted; nil, a file field is refused (invalid).
	Files FileSource
	// Items are the values a field may hold — path the field's in the schema,
	// no index (`links.href`) —: a list only the request knows
	// (EnumItemsFunc); nil, those of its enum.
	Items func(path string, f *Field) ([]EnumItem, error)
	// NoValidations: the `[validation=…]` are not called (a value read only to
	// be typed: Typed).
	NoValidations bool
	// PageOptions are the values of the fields of options of a function the
	// page holds, by path (a check of the page's: Check); NoFetch, a request
	// of options is not fetched — any value taken —: the page's check.
	PageOptions map[string][]EnumItem
	NoFetch     bool

	issues   []Issue
	warnings []Warning
	// plain: the fields read whatever is answered (no metadata functions
	// called): what Typed reads; scope, the record being read.
	plain  bool
	scope  *scope
	cancel context.CancelFunc
}

// prepare sets what a read needs: the VM, a bound of its time; done undoes
// it.
func (d *Decoder) prepare() {
	d.issues, d.warnings = nil, nil
	if d.VM == nil && d.Schema != nil {
		d.VM = d.Schema.vm
	}
	if d.Ctx == nil {
		d.Ctx, d.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	}
}

func (d *Decoder) done() {
	if d.cancel != nil {
		d.cancel()
		d.Ctx, d.cancel = nil, nil
	}
	d.scope = nil
}

// Decode reads data — the form's values, a record nested as JSON holds it
// (Nest: as the admin posts them) — by the schema: the record, typed, and
// what is wrong with it; the form's own `[validation=…]` last.
func (d *Decoder) Decode(data map[string]any) (Record, []Issue) {
	return d.decode(d.Schema.Fields, d.Schema.Validations, data)
}

// Value is Decode of a schema of any shape — a record, a list (of records,
// of values), a choice of classes —: raw as JSON holds it (a list, a map of
// its indexes).
func (d *Decoder) Value(raw any) (any, []Issue) {
	s := d.Schema
	if !s.Slice && !s.Choice {
		m, _ := raw.(map[string]any)
		return d.Decode(m)
	}
	d.issues = nil
	if d.VM == nil {
		d.VM = s.vm
	}
	if d.Ctx == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		d.Ctx = ctx
		defer func() { d.Ctx = nil }()
	}
	return d.field(&Field{Schema: s, Nullable: true}, raw, ""), d.issues
}

// DecodeFields is Decode of some fields of the form alone (a step): the
// form's own validations not called.
func (d *Decoder) DecodeFields(fields []*Field, data map[string]any) (Record, []Issue) {
	return d.decode(fields, nil, data)
}

func (d *Decoder) decode(fields []*Field, validations []gad.Object, data map[string]any) (Record, []Issue) {
	d.prepare()
	defer d.done()
	// the fields drawn by what is answered: by the record typed
	if !d.plain && d.Schema.dynamic() {
		d.scope = d.rootScope(data)
	}
	rec := d.record(fields, data, "")
	d.recordValidations(validations, "", rec)
	return rec, d.issues
}

// Typed is data by the schema, typed (a number an int, a date a date…): what
// is wrong not said.
func (d *Decoder) Typed(data map[string]any) Record {
	quiet := *d
	quiet.NoValidations, quiet.plain = true, true
	quiet.Ctx, quiet.cancel, quiet.scope = nil, nil, nil
	rec, _ := quiet.Decode(data)
	return rec
}

func (d *Decoder) fail(path, code string) {
	d.issues = append(d.issues, Issue{Path: path, Code: code})
}

// failOf says what is wrong with the value raw of f at path.
func (d *Decoder) failOf(f *Field, path, code string, raw any) {
	d.issues = append(d.issues, Issue{Path: path, Code: code, Field: f, Value: Text(raw)})
}

// failed says something is wrong at path (under: or under it).
func (d *Decoder) failed(path string, under bool) bool {
	for _, e := range d.issues {
		if e.Path == path || under && (path == "" || strings.HasPrefix(e.Path, path+".") || strings.HasPrefix(e.Path, path+"[")) {
			return true
		}
	}
	return false
}

func (d *Decoder) record(fields []*Field, data map[string]any, path string) Record {
	rec := make(Record, 0, len(fields))
	warnings := -1
	for _, f := range fields {
		p := join(path, f.Name)
		if f.Name == WarningsField {
			// the form's own (WithWarnings): set once the record is read
			warnings = len(rec)
			rec = append(rec, RecordField{Name: f.Name})
			continue
		}
		read := f
		if !d.plain && d.scope != nil && len(f.Cond) > 0 {
			st := d.state(f)
			if !st.Asked() {
				// not there, hidden or disabled: not asked
				rec = append(rec, RecordField{Name: f.Name})
				continue
			}
			if st.Readonly {
				c := *f
				c.Nullable = true
				read = &c
			}
		}
		rec = append(rec, RecordField{Name: f.Name, Value: d.field(read, data[f.Name], p)})
	}
	if warnings >= 0 {
		if w := d.warningsOf(path); len(w) > 0 {
			rec[warnings].Value = w
		}
	}
	if !d.NoValidations {
		for i, f := range fields {
			if f.Schema == nil || f.Schema.Slice {
				d.fieldValidations(f, join(path, f.Name), rec[i].Value, rec)
			}
		}
	}
	return rec
}

// field reads the value raw of f, at path.
func (d *Decoder) field(f *Field, raw any, path string) any {
	switch {
	case f.ReadOnly:
		// shown, never edited: the application keeps what it holds
		// (Schema.KeepReadOnly)
		return raw
	case f.Type == RangeType:
		return d.rangeOf(f, raw, path)
	case IsFileType(f.Type):
		return d.file(f, raw, path)
	case f.Schema != nil && f.Schema.Slice:
		return d.list(f, raw, path)
	case f.Schema != nil && f.Schema.Choice:
		return d.choice(f, raw, path)
	case f.Schema != nil:
		sub, _ := raw.(map[string]any)
		// an optional card left blank — a form posts its inputs, empty —
		// is not there
		if f.Nullable && Blank(sub) {
			return nil
		}
		outer := d.scope
		if outer != nil {
			d.scope = outer.child(f.Name, f.Schema.recClass, f.Schema.Fields, outer.typed[f.Name])
		}
		sr := d.record(f.Schema.Fields, sub, path)
		d.scope = outer
		d.recordValidations(append(append([]gad.Object{}, f.Validations...), f.Schema.Validations...), path, sr)
		return sr
	}
	return d.value(f, raw, path)
}

// IsFileType says t is the type of a file field: a file, an image.
func IsFileType(t string) bool { return t == FileType || t == ImageType }

func (d *Decoder) file(f *Field, raw any, path string) any {
	if Blank(raw) {
		if !f.Nullable {
			d.fail(path, CodeRequired)
		}
		return nil
	}
	if d.Files == nil {
		d.fail(path, CodeInvalid)
		return nil
	}
	v, code := d.Files.File(f, raw)
	if code != "" {
		d.fail(path, code)
		return nil
	}
	return v
}

// choice reads a choice of classes (Schema.Choice): an object of the one
// key of the class chosen.
func (d *Decoder) choice(f *Field, raw any, path string) any {
	m, _ := raw.(map[string]any)
	for _, cf := range f.Schema.Fields {
		sub, ok := m[cf.Name]
		if !ok {
			continue
		}
		return Record{{Name: cf.Name, Value: d.field(cf, sub, join(path, cf.Name))}}
	}
	if !f.Nullable {
		d.fail(path, CodeRequired)
	}
	return nil
}

// list reads a list: of plain values, each of its item's type (an empty
// item dropped); of records, each one at its index (`links[1]`); then how
// many ([min=…, max=…]).
func (d *Decoder) list(f *Field, raw any, path string) any {
	items := asList(raw)
	if len(items) == 0 {
		if !f.Nullable {
			d.fail(path, CodeRequired)
		}
		return []any{}
	}
	out := make([]any, 0, len(items))
	item := f.Schema.Item
	switch {
	case item == nil:
		// a list of records of its own fields (a schema of a list)
		for i, x := range items {
			sub, _ := x.(map[string]any)
			p := path + "[" + strconv.Itoa(i) + "]"
			restore := d.itemScope(f, f.Schema.recClass, f.Schema.Fields, i)
			sr := d.record(f.Schema.Fields, sub, p)
			restore()
			d.recordValidations(f.Schema.Validations, p, sr)
			out = append(out, sr)
		}
	case item.Schema != nil:
		// of records: each at its index, its record's validations its own
		for i, x := range items {
			sub, _ := x.(map[string]any)
			p := path + "[" + strconv.Itoa(i) + "]"
			restore := d.itemScope(f, item.Schema.recClass, item.Schema.Fields, i)
			sr := d.record(item.Schema.Fields, sub, p)
			restore()
			d.recordValidations(item.Schema.Validations, p, sr)
			out = append(out, sr)
		}
	default:
		it := *item
		it.Nullable = true // an empty item is dropped, not an error
		for _, x := range items {
			var v any
			if IsFileType(it.Type) {
				v = d.file(&it, x, path)
			} else {
				v = d.value(&it, x, path)
			}
			if v != nil && v != "" {
				out = append(out, v)
			}
		}
	}
	switch {
	case len(out) == 0 && !f.Nullable:
		d.fail(path, CodeRequired)
	case f.MinItems > 0 && len(out) < f.MinItems:
		d.fail(path, CodeItemsMin)
	case f.MaxItems > 0 && len(out) > f.MaxItems:
		d.fail(path, CodeItemsMax)
	}
	return out
}

// asList is raw as a list: a list, or the values of a map by their index
// (what a form posts as `key[0]`, `key[1]`).
func asList(raw any) []any {
	switch x := raw.(type) {
	case []any:
		return x
	case map[string]any:
		var out []any
		for i := 0; ; i++ {
			v, ok := x[strconv.Itoa(i)]
			if !ok {
				return out
			}
			out = append(out, v)
		}
	}
	return nil
}

// rangeOf reads a range (Range[T]): its from and its to, each of the
// bounds' type and limits — both asked once it is begun, but the one left
// open ([open=…]) —, the from not after the to. An optional range left blank
// is none.
func (d *Decoder) rangeOf(f *Field, raw any, path string) any {
	m, _ := raw.(map[string]any)
	if f.Nullable && Blank(m) {
		return nil
	}
	bound := Field{Type: DefaultType}
	if f.Range != nil {
		bound = *f.Range
	}
	out := make(Record, 0, 2)
	for _, name := range []string{RangeFrom, RangeTo} {
		b := bound
		b.Name, b.Meta = name, nil
		b.Nullable = f.RangeOpen == name
		out = append(out, RecordField{Name: name, Value: d.value(&b, m[name], path+"."+name)})
	}
	from, to := out[0].Value, out[1].Value
	if from != nil && to != nil && !d.failed(path, true) && After(bound.Enum, from, to) {
		d.fail(path, CodeRange)
	}
	return out
}

// value reads one value of f: typed — a type of the application's by its
// reader (Types) —, checked by its limits; on an error, what was posted.
func (d *Decoder) value(f *Field, raw any, path string) any {
	fail := func(code string) any {
		d.failOf(f, path, code, raw)
		if Text(raw) == "" {
			return zeroOf(f)
		}
		return Text(raw)
	}
	_, must := f.Meta["required"]
	if read, ok := d.Types[f.Type]; ok {
		if Blank(raw) {
			if !f.Nullable {
				return fail(CodeRequired)
			}
			return nil
		}
		v, code := read(f, path, raw)
		if code != "" {
			return fail(code)
		}
		if _, isText := raw.(string); isText {
			if code := CheckLimits(f, Text(raw)); code != "" {
				return fail(code)
			}
		}
		return v
	}
	switch {
	case f.Type == "bool" && f.Nullable:
		// optional: yes, no, or not said (nil); with [required], yes or no
		if b, ok := raw.(bool); ok {
			return b
		}
		switch Text(raw) {
		case "true", "1", "on":
			return true
		case "false", "0", "off":
			return false
		case "":
			if must {
				return fail(CodeRequired)
			}
			return nil
		}
		return fail(CodeChoice)
	case f.Type == "bool":
		// a box that must be checked (accepting terms) says so: [required];
		// else false when not checked
		v := false
		switch x := raw.(type) {
		case bool:
			v = x
		default:
			s := Text(raw)
			v = s == "true" || s == "1" || s == "on"
		}
		if must && !v {
			return fail(CodeRequired)
		}
		return v
	case f.Type == AddressType:
		return d.address(f, raw, path)
	}
	s := Text(raw)
	if s == "" {
		if !f.Nullable {
			return fail(CodeRequired)
		}
		return nil
	}
	if items, err := d.items(path, f); err != nil {
		d.issues = append(d.issues, Issue{Path: path, Code: CodeItems, Field: f, Value: s, Err: err})
		return s
	} else if items != nil {
		if !slices.ContainsFunc(items, func(it EnumItem) bool { return it.Name == s }) {
			return fail(CodeChoice)
		}
		// one of its values: typed as its field (an option of an int, an int)
	}
	var v any = s
	switch f.Type {
	case "int":
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fail(CodeInvalid)
		}
		v = n
	case "uint":
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return fail(CodeInvalid)
		}
		v = n
	case "float", "decimal":
		n, err := strconv.ParseFloat(s, 64)
		if err != nil || math.IsNaN(n) || math.IsInf(n, 0) {
			return fail(CodeInvalid)
		}
		v = n
	case "date", "calendarDate":
		if _, err := time.Parse("2006-01-02", s); err != nil {
			return fail(CodeInvalid)
		}
	case "calendarTime":
		// what datetime-local gives: 2026-01-30T14:30 (its seconds, maybe)
		if _, err := time.Parse("2006-01-02T15:04", s); err != nil {
			if _, err := time.Parse("2006-01-02T15:04:05", s); err != nil {
				return fail(CodeInvalid)
			}
		}
	case "time":
		if _, err := time.Parse("15:04", s); err != nil {
			if _, err := time.Parse("15:04:05", s); err != nil {
				return fail(CodeInvalid)
			}
		}
	case "color":
		if !colorRe.MatchString(s) {
			return fail(CodeInvalid)
		}
	}
	if f.MaxLength == 0 && !NumberTypes[f.Type] {
		max := MaxLine
		if f.Type == "text" || f.Type == "html" {
			max = MaxText
		}
		if len(s) > max {
			return fail(CodeTooLong)
		}
	}
	if code := CheckLimits(f, s); code != "" {
		return fail(code)
	}
	return v
}

var colorRe = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

// items are the values f may hold: the application's (Items), else its
// enum's; nil when it holds any.
func (d *Decoder) items(path string, f *Field) ([]EnumItem, error) {
	if f.OptionsFunc != nil {
		if d.plain || d.scope == nil {
			return nil, nil
		}
		return d.dynamicItems(path, f)
	}
	if d.Items != nil {
		items, err := d.Items(SchemaPath(path), f)
		if err != nil || items != nil {
			return items, err
		}
	}
	if f.Enum != nil {
		if f.Enum.Items != nil {
			return f.Enum.Items, nil
		}
		items := make([]EnumItem, len(f.Enum.Names))
		for i, n := range f.Enum.Names {
			items[i] = EnumItem{Name: n}
		}
		return items, nil
	}
	return nil, nil
}

var indexRe = regexp.MustCompile(`\[\d+\]`)

// SchemaPath is a value's path as its field's in the schema: no index
// (`links[1].href` → `links.href`).
func SchemaPath(path string) string { return indexRe.ReplaceAllString(path, "") }

// zeroOf is the value of f's type that is none: what a required number left
// empty is kept as, the value keeping its shape.
func zeroOf(f *Field) any {
	switch f.Type {
	case "int":
		return int64(0)
	case "uint":
		return uint64(0)
	case "float", "decimal":
		return 0.0
	}
	return ""
}

// address reads an address (AddressType): its words and parts, its
// coordinates; typed with no place found, its words alone.
func (d *Decoder) address(f *Field, raw any, path string) any {
	m, _ := raw.(map[string]any)
	if m == nil && Text(raw) != "" {
		m = map[string]any{"formatted": Text(raw)}
	}
	if m == nil || Text(m["formatted"]) == "" {
		if !f.Nullable {
			d.fail(path, CodeRequired)
		}
		return nil
	}
	rec := make(Record, 0, len(AddressFields))
	for _, name := range AddressFields {
		var v any
		switch name {
		case "lat", "lng":
			if n, err := strconv.ParseFloat(Text(m[name]), 64); err == nil && !math.IsNaN(n) && !math.IsInf(n, 0) {
				v = n
			}
		default:
			s := Text(m[name])
			if len(s) > MaxLine {
				d.fail(path, CodeTooLong)
				return nil
			}
			v = s
		}
		rec = append(rec, RecordField{Name: name, Value: v})
	}
	return rec
}

// fieldValidations run the `[validation=…]` of the field at path over its
// value (nil: left blank, not checked) of the record rec: the first that
// refuses it says why. A regular expression is the whole value matching it.
func (d *Decoder) fieldValidations(f *Field, path string, value any, rec Record) {
	if d.NoValidations || len(f.Validations) == 0 || value == nil || d.failed(path, false) {
		return
	}
	val, r := GadValue(value), GadValue(rec)
	for _, fn := range f.Validations {
		if re, ok := fn.(*gad.Regexp); ok {
			if !(*regexp.Regexp)(re).MatchString(Text(value)) {
				d.fail(path, CodePattern)
				return
			}
			continue
		}
		if d.result(path, false, func() (gad.Object, error) { return d.call(fn, val, r) }) {
			return
		}
	}
}

// recordValidations run the validations of the record rec at path — a
// field's of a class, then its class's — once its fields have nothing
// wrong: one says what is wrong with it, a str, or with its fields, a dict
// of their names.
func (d *Decoder) recordValidations(fns []gad.Object, path string, rec Record) {
	if d.NoValidations || len(fns) == 0 || rec == nil || d.failed(path, true) {
		return
	}
	r := GadValue(rec)
	for _, fn := range fns {
		if d.result(path, true, func() (gad.Object, error) { return d.call(fn, r) }) {
			return
		}
	}
}

// call calls fn with args — as many as it takes: `func(value)` is called
// with the value alone; of a function of several methods, the one of the
// args (else of the value alone). Without a VM, nothing is called.
func (d *Decoder) call(fn gad.Object, args ...gad.Object) (gad.Object, error) {
	if d.VM == nil {
		return gad.Nil, nil
	}
	var a gad.Array = args
	if spec, ok := fn.(interface{ GetFuncSpec() *gad.FuncSpec }); ok {
		co, _ := spec.GetFuncSpec().CallerMethodWithValidationCheckOfArgs(gad.Args{a})
		for n := len(a) - 1; co == nil && n >= 1; n-- {
			a = a[:n]
			co, _ = spec.GetFuncSpec().CallerMethodWithValidationCheckOfArgs(gad.Args{a})
		}
		if co != nil {
			fn = co
		}
	}
	for {
		cm, ok := fn.(*gad.CallerMethod)
		if !ok || cm.CallerObject == nil {
			break
		}
		fn = cm.CallerObject
	}
	if cf, ok := fn.(*gad.CompiledFunction); ok {
		items := cf.Params.Items
		variadic := len(items) > 0 && items[len(items)-1].Var
		if !variadic && len(items) < len(a) {
			a = a[:len(items)]
		}
	}
	return gad.NewInvoker(d.VM, fn).WithContext(d.Ctx).Invoke(gad.Args{a}, nil)
}

// result reads what a validation answered: nil or true, nothing wrong;
// false, the value invalid; a str, why; a dict (of a record: fields), why,
// field by field. A function that fails is the form's error: logged, the
// value invalid. It reports whether something was wrong.
func (d *Decoder) result(path string, record bool, call func() (gad.Object, error)) bool {
	ret, err := call()
	if err != nil {
		log.Printf("ERROR: form: [validation] of %s: %v", orRoot(path), err)
		d.fail(path, CodeInvalid)
		return true
	}
	switch r := ret.(type) {
	case nil:
		return false
	case gad.Bool:
		if r {
			return false
		}
		d.fail(path, CodeInvalid)
		return true
	case gad.Str:
		d.issues = append(d.issues, Issue{Path: path, Code: CodeValidation, Message: string(r)})
		return true
	case gad.Dict:
		if record && len(r) > 0 {
			for k, m := range r {
				d.issues = append(d.issues, Issue{Path: join(path, k), Code: CodeValidation, Message: m.ToString()})
			}
			return true
		}
		if len(r) == 0 {
			return false
		}
	}
	if ret == gad.Nil {
		return false
	}
	log.Printf("ERROR: form: [validation] of %s answered %s: nil, a str or (of a class) a dict", orRoot(path), ret.Type().Name())
	d.fail(path, CodeInvalid)
	return true
}

func orRoot(path string) string {
	if path == "" {
		return "the form"
	}
	return path
}

// GadValuer is a value of the application's (a file sent) as gad holds it:
// what a `[validation=…]` is given.
type GadValuer interface {
	GadValue() gad.Object
}

// GadValue is a value read as gad holds it: a record a dict of its fields.
func GadValue(v any) gad.Object {
	switch x := v.(type) {
	case nil:
		return gad.Nil
	case GadValuer:
		return x.GadValue()
	case Record:
		dict := make(gad.Dict, len(x))
		for _, f := range x {
			dict[f.Name] = GadValue(f.Value)
		}
		return dict
	case map[string]any:
		dict := make(gad.Dict, len(x))
		for k, e := range x {
			dict[k] = GadValue(e)
		}
		return dict
	case []any:
		arr := make(gad.Array, len(x))
		for i, e := range x {
			arr[i] = GadValue(e)
		}
		return arr
	case string:
		return gad.Str(x)
	case int64:
		return gad.Int(x)
	case float64:
		return gad.Float(x)
	case bool:
		return gad.Bool(x)
	}
	if o, err := gad.ToObject(v); err == nil {
		return o
	}
	return gad.Str(fmt.Sprint(v))
}

// Text is a value posted as text, trimmed.
func Text(raw any) string {
	switch v := raw.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(v)
	case bool:
		return strconv.FormatBool(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return strings.TrimSpace(fmt.Sprint(raw))
}

// Blank says a value posted holds nothing: no value but empty ones (a box
// not checked, an empty list).
func Blank(v any) bool {
	switch x := v.(type) {
	case nil:
		return true
	case map[string]any:
		for _, e := range x {
			if !Blank(e) {
				return false
			}
		}
		return true
	case []any:
		return len(x) == 0
	case bool:
		return !x
	}
	return Text(v) == ""
}

// Nest is the values of a form posted flat (`key.a.b`, `key.list[0].x`, as
// the admin posts them), under key, as a record nested — what a Decoder
// reads. A list is a map of its indexes ("0", "1"…), read as a list.
func Nest(values url.Values, key string) map[string]any {
	out := map[string]any{}
	for k, vs := range values {
		var rest string
		switch {
		case key == "":
			rest = k
		case strings.HasPrefix(k, key+"."):
			rest = k[len(key)+1:]
		case strings.HasPrefix(k, key+"["):
			// a list's: `key[0].x`
			rest = k[len(key)+1:]
		default:
			continue
		}
		var v any
		if len(vs) == 1 {
			v = vs[0]
		} else {
			list := make([]any, len(vs))
			for i, s := range vs {
				list[i] = s
			}
			v = list
		}
		setNested(out, splitKey(rest), v)
	}
	return out
}

// splitKey is a posted key by its parts: `a.b[0].c` → a, b, 0, c.
func splitKey(k string) []string {
	k = strings.ReplaceAll(k, "[", ".")
	k = strings.ReplaceAll(k, "]", "")
	return strings.Split(k, ".")
}

func setNested(m map[string]any, parts []string, v any) {
	for _, p := range parts[:len(parts)-1] {
		sub, ok := m[p].(map[string]any)
		if !ok {
			sub = map[string]any{}
			m[p] = sub
		}
		m = sub
	}
	m[parts[len(parts)-1]] = v
}

// join is the path of name under path.
func join(path, name string) string {
	if path == "" {
		return name
	}
	return path + "." + name
}

// Map is the record's fields by their names (a record within one stays a
// record).
func (r Record) Map() map[string]any {
	out := make(map[string]any, len(r))
	for _, f := range r {
		out[f.Name] = f.Value
	}
	return out
}

// ReadOrder reads an order (Builder.Order): its items of a field and a
// direction (ASC, DESC) — `key[0].field`, `key[0].dir` —, an item of no
// field or no direction dropped.
func ReadOrder(_ *Field, _ string, raw any) (any, string) {
	out := []any{}
	for _, it := range asList(raw) {
		m, _ := it.(map[string]any)
		field, dir := Text(m["field"]), strings.ToUpper(Text(m["dir"]))
		if field == "" || (dir != "ASC" && dir != "DESC") {
			continue
		}
		out = append(out, map[string]any{"field": field, "dir": dir})
	}
	return out, ""
}
