package schemaform

import (
	"errors"
	"fmt"
	"github.com/go-rvq/rvq/admin/presets"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gad-lang/gad"
)

// FormName is the name of the interface that IS the form. A schema written in
// parts declares others — the interfaces and enums it uses — and names this
// one; a schema written in one piece needs no name, and is given this one.
const FormName = "Form"

// RunTimeout bounds how long a schema may run. A schema is a declaration and
// runs in microseconds; one that does not return by then is not a schema.
var RunTimeout = time.Second

// Parse reads a schema with the types of a new Builder. See Builder.Parse.
func Parse(src string) (*Schema, error) { return New().Parse(src) }

// Parse RUNS the schema and reads the form off what it returns.
//
// The schema is gad: it is compiled and run in a VM of its own, and the
// interface named Form is what the run returns. Everything the form needs is
// then read from that interface through gad's own reflection (IndexGetter) —
// its fields and their resolved types, its array depth, its `[k=v, …]`
// metadata — so a schema can use whatever gad can declare: interfaces by name,
// enums, typed arrays, metadata on the interface and on its fields.
//
// The type names a field may use are the builder's: every type registered on
// it is in scope while the schema runs, so `color` or `html` resolve even though
// gad has no such types.
func (b *Builder) Parse(src string) (*Schema, error) {
	program, err := formProgram(src)
	if err != nil {
		return nil, err
	}

	key := b.typesKey() + "\x00" + program
	if s, ok := parsed.Load(key); ok {
		return s.(*Schema), nil
	}

	vm, form, err := b.run(program)
	if err != nil {
		return nil, err
	}

	r := &reader{vm: vm, enums: map[string]*Enum{}, visiting: map[*gad.Interface]bool{}}
	s, err := r.schema(form)
	if err != nil {
		return nil, err
	}
	if err := s.readLayout(); err != nil {
		return nil, err
	}
	s.shareEnums(r.enums)

	parsed.Store(key, s)
	return s, nil
}

// parsed caches what a schema read into, by the program and the types it ran
// with: a schema is read on every draw and every save, and it is the same
// declaration each time. What is cached is never changed afterwards.
var parsed sync.Map

// typesKey names the set of types a schema runs with, for the cache.
func (b *Builder) typesKey() string {
	names := make([]string, 0, len(b.types))
	for name := range b.types {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}

var (
	// anonIface is an interface written without a name where a declaration
	// starts — at the start of a line, after a `;` or a `}`, or after the `]`
	// that closes its metadata block: `interface { … }`, `interface[] { … }`,
	// `[layout="table"] interface []{ … }`.
	anonIface = regexp.MustCompile(`(?m)(^|[;}\]])([ \t]*)interface([ \t]*)([\[{])`)
	// formIface is the interface named Form, where a declaration starts.
	formIface = regexp.MustCompile(`(?m)(^|[;}\]])[ \t]*interface[ \t]+` + FormName + `\b`)
)

// formProgram is the schema written as a program that returns its form.
//
// A schema in one piece may leave out the keyword and the name — `{…}`,
// `[]{…}`, `[]str` — or leave out only the name (`interface []{…}`, and an
// enum declared before it); either way it is the form, and is named so. A
// schema in parts names its form itself.
func formProgram(src string) (string, error) {
	trimmed := strings.TrimSpace(src)
	if trimmed == "" {
		return "", errors.New("schemaform: the schema is empty")
	}

	switch {
	case strings.HasPrefix(trimmed, "{"), strings.HasPrefix(trimmed, "[]"):
		trimmed = "interface " + FormName + " " + trimmed
	case formIface.MatchString(trimmed):
	default:
		switch n := len(anonIface.FindAllStringIndex(trimmed, -1)); n {
		case 1:
			trimmed = anonIface.ReplaceAllString(trimmed, "${1}${2}interface "+FormName+" ${4}")
		case 0:
			return "", fmt.Errorf("schemaform: the form is the interface %q, and there is none", FormName)
		default:
			return "", fmt.Errorf(
				"schemaform: a schema written in parts names its form %q; there are %d interfaces without a name",
				FormName, n)
		}
	}

	// Whatever named the form, an interface left without a name beside it is
	// one nothing can refer to: a mistake written down.
	if anonIface.MatchString(trimmed) {
		return "", fmt.Errorf(
			"schemaform: the form is the interface %q; an interface beside it must have a name too", FormName)
	}

	return trimmed + "\nreturn " + FormName + "\n", nil
}

// run compiles and runs the program, with the builder's types in scope, and
// returns the VM — the field types are symbols only it resolves — and the form.
func (b *Builder) run(program string) (*gad.VM, *gad.Interface, error) {
	vm, ret, err := b.runProgram(program)
	if err != nil {
		return nil, nil, err
	}
	form, ok := ret.(*gad.Interface)
	if !ok {
		return nil, nil, fmt.Errorf("schemaform: %q is not an interface, it is %s", FormName, ret.Type().Name())
	}
	return vm, form, nil
}

// runProgram compiles and runs the program, with the builder's types in scope,
// and returns the VM and what the program returned.
func (b *Builder) runProgram(program string) (*gad.VM, gad.Object, error) {
	builtins := b.Builtins()
	st := gad.NewSymbolTable(builtins.NameSet)
	res, err := gad.Compile(st, []byte(program), gad.CompileOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("schemaform: %w", err)
	}

	vm := gad.NewVM(builtins.Build(), res.Bytecode)
	timer := time.AfterFunc(RunTimeout, vm.Abort)
	defer timer.Stop()

	ret, err := vm.Run()
	if err != nil {
		return nil, nil, fmt.Errorf("schemaform: %w", err)
	}
	return vm, ret, nil
}

// Builtins are gad's, with the builder's types in scope: what a program whose
// classes are read as forms (Class) runs with — an application adds its own
// names to them.
func (b *Builder) Builtins() *gad.Builtins {
	builtins := gad.NewBuiltins()
	for name := range b.types {
		// A name gad already has as a TYPE stays gad's — `str`, `int`, `bool`
		// are the same thing to both. Any other is declared as a marker: `color`
		// is not a gad type, and `time` is gad's time NAMESPACE, not a type.
		if t, ok := builtins.NameSet[name]; ok {
			if _, isType := builtins.Objects[t].(gad.ObjectType); isType {
				continue
			}
		}
		builtins.Set(name, &typeMarker{name: name})
	}
	return builtins
}

// TypeMarker is the object a program sees for the type name — a type gad has
// not (`color`, `contact_form.email`): a field typed by it is of that type.
// An application puts it where its programs find it (a module's member).
func TypeMarker(name string) gad.Object { return &typeMarker{name: name} }

// Class reads a gad class as a form: a record of its fields, the way an
// interface's are read, with the default each one has. A field typed by a
// class — `a class { … }`, or one by name — is a record of its own, a group
// of the form; a type only the application knows is named by TypeOf. vm is
// the VM that ran the code the class comes from.
func (b *Builder) Class(vm *gad.VM, c *gad.Class) (*Schema, error) {
	r := &reader{vm: vm, enums: map[string]*Enum{}, visiting: map[*gad.Interface]bool{}, typeOf: b.typeOf,
		choiceAsName: b.choiceAsName, className: b.className}
	s, err := r.class(c)
	if err != nil {
		return nil, err
	}
	s.shareEnums(r.enums)
	return s, nil
}

// class reads the fields of a gad class — what a layout's Config declares —
// the way an interface's are read, with the default each one has.
func (r *reader) class(c *gad.Class) (*Schema, error) {
	if r.visitingClass[c] {
		return nil, fmt.Errorf("schemaform: the class %q contains itself", c.Name())
	}
	if r.visitingClass == nil {
		r.visitingClass = map[*gad.Class]bool{}
	}
	r.visitingClass[c] = true
	defer delete(r.visitingClass, c)

	s := &Schema{}
	// the fields of its parents (`*Parent`) first, in their order: a class
	// extends them
	for _, p := range c.RawParents() {
		ps, err := r.class(p.Type)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", p.Type.Name(), err)
		}
		// (their validations, before its own)
		s.Validations = append(s.Validations, ps.Validations...)
		owner := p.Type.Name()
		if r.className != nil {
			if n := r.className(p.Type); n != "" {
				owner = n
			}
		}
		for _, f := range ps.Fields {
			f.Owner = owner
		}
		appendRows(s, ps.Fields)
	}
	own, err := classValidations(c)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", c.Name(), err)
	}
	s.Validations = append(s.Validations, own...)
	for _, cf := range c.RawFields() {
		types := make(gad.Array, len(cf.Types))
		for i, t := range cf.Types {
			types[i] = t
		}
		f := &Field{Name: cf.Name, Meta: metaOf(cf.Meta), Nullable: cf.Nullable}
		if err := r.typeInto(f, types); err != nil {
			return nil, fmt.Errorf("%s: %w", f.Name, err)
		}
		if err := options(f, cf.Meta); err != nil {
			return nil, fmt.Errorf("%s: %w", f.Name, err)
		}
		if cf.Value != nil && cf.Value != gad.Nil {
			f.Default = metaValue(cf.Value)
		}
		if err := r.appendField(s, f); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// choice reads a union of classes (`listLayout|gridLayout`) into a choice of
// one of them (Schema.Choice): a field per class, named as the class is, its
// schema the class's record, its metadata the class's. nil when a type of the
// union is not a class.
func (r *reader) choice(types gad.Array) (*Schema, error) {
	s := &Schema{Choice: true}
	for _, t := range types {
		c, ok := t.(*gad.Class)
		if !ok {
			return nil, nil
		}
		sub, err := r.class(c)
		if err != nil {
			return nil, err
		}
		s.Fields = append(s.Fields, &Field{
			Name: c.Name(), Type: FormType, Schema: sub, Nullable: true, Meta: metaOf(c.Meta),
		})
	}
	return s, nil
}

// classNames reads a union of classes as an enum of their names, each labelled
// by its class's metadata — its label, else its name humanized — and hinted
// by its hint. nil when a type of the union is not a class.
func classNames(types gad.Array) *Enum {
	e := &Enum{}
	for _, t := range types {
		c, ok := t.(*gad.Class)
		if !ok {
			return nil
		}
		it := EnumItem{Name: c.Name(), Label: presets.HumanizeString(c.Name())}
		meta := metaOf(c.Meta)
		if l, _ := meta["label"].(string); l != "" {
			it.Label = l
		}
		if hint, _ := meta["hint"].(string); hint != "" {
			it.Hint = hint
		}
		e.Names = append(e.Names, c.Name())
		e.Items = append(e.Items, it)
	}
	e.Name = strings.Join(e.Names, "|")
	return e
}

// typeMarker is a type the builder knows and gad does not — `color`, `html` —
// declared while the schema runs so a field may be typed with it. It is only
// ever reflected upon: its name is the component the field is drawn with.
type typeMarker struct{ name string }

var _ gad.Object = (*typeMarker)(nil)

func (m *typeMarker) Type() gad.ObjectType { return gad.TStr }
func (m *typeMarker) ToString() string     { return m.name }
func (m *typeMarker) IsFalsy() bool        { return false }

// AssignTo and CanAssign make it a type a class field may declare
// (`m text`): the value is the form's to check (DecodeForm), not gad's.
func (m *typeMarker) AssignTo(_ *gad.VM, obj gad.Object, _ gad.TypeAssigner) (gad.Object, error) {
	return obj, nil
}
func (m *typeMarker) CanAssign(gad.Object) (bool, error) { return true, nil }
func (m *typeMarker) Equal(right gad.Object) bool {
	r, ok := right.(*typeMarker)
	return ok && r.name == m.name
}

// reader reads a form off the interface a schema returned, through gad's
// reflection, with the VM that resolves the symbols the types are.
type reader struct {
	vm *gad.VM
	// enums are the enums the fields hold, by name.
	enums map[string]*Enum
	// visiting are the interfaces being read, down from the form: one met again
	// on the way down contains itself. visitingClass, the classes.
	visiting      map[*gad.Interface]bool
	visitingClass map[*gad.Class]bool
	// typeOf names the application's own types (Builder.TypeOf).
	typeOf TypeOfFunc
	// choiceAsName reads a union of classes as the choice of a class's name
	// (Builder.ChoiceAsName).
	choiceAsName bool
	// className names the classes a class extends (Builder.ClassName).
	className ClassNameFunc
}

// index is obj[key] through the object's own reflection.
func (r *reader) index(obj gad.Object, key string) (gad.Object, error) {
	ig, ok := obj.(gad.IndexGetter)
	if !ok {
		return nil, fmt.Errorf("schemaform: %s has no %q", obj.Type().Name(), key)
	}
	return ig.IndexGet(r.vm, gad.Str(key))
}

// array is obj[key] as an array.
func (r *reader) array(obj gad.Object, key string) (gad.Array, error) {
	v, err := r.index(obj, key)
	if err != nil {
		return nil, err
	}
	if v == gad.Nil {
		return nil, nil
	}
	arr, ok := v.(gad.Array)
	if !ok {
		return nil, fmt.Errorf("schemaform: %q is a %s, not an array", key, v.Type().Name())
	}
	return arr, nil
}

// depth is obj[key] as an int.
func (r *reader) depth(obj gad.Object) (int, error) {
	v, err := r.index(obj, "@depth")
	if err != nil {
		return 0, err
	}
	n, ok := v.(gad.Int)
	if !ok {
		return 0, fmt.Errorf("schemaform: @depth is a %s", v.Type().Name())
	}
	return int(n), nil
}

// meta is obj's `[k=v, …]` block.
func (r *reader) meta(obj gad.Object) (Meta, error) {
	v, err := r.index(obj, "@meta")
	if err != nil {
		return nil, err
	}
	return metaOf(v), nil
}

// schema reads one interface: a record, a list of records (`[]{…}`), or a list
// of plain values (`[]str`, an interface whose element is a type) — one schema
// per `[]`, so a list of lists nests.
func (r *reader) schema(iface *gad.Interface) (*Schema, error) {
	if r.visiting[iface] {
		return nil, fmt.Errorf("schemaform: the interface %q contains itself", iface.IName)
	}
	r.visiting[iface] = true
	defer delete(r.visiting, iface)

	depth, err := r.depth(iface)
	if err != nil {
		return nil, err
	}
	meta, err := r.meta(iface)
	if err != nil {
		return nil, err
	}
	elem, err := r.array(iface, "@elem")
	if err != nil {
		return nil, err
	}

	var s *Schema
	if len(elem) > 1 {
		return nil, fmt.Errorf("schemaform: a list of values holds ONE type, and %q declares %d", iface.IName, len(elem))
	}
	if len(elem) == 1 {
		// `interface Form []str` — the element is a type, and the item IS the
		// value.
		item := &Field{}
		if err := r.typeInto(item, elem); err != nil {
			return nil, err
		}
		s = &Schema{Slice: true, Item: item}
	} else {
		s = &Schema{Slice: depth > 0}
		// the getters (`get path str`) first: fields shown, never edited
		props, err := r.array(iface, "props")
		if err != nil {
			return nil, err
		}
		for _, o := range props {
			f, err := r.getter(o)
			if err != nil {
				return nil, err
			}
			if f != nil {
				if err := r.appendField(s, f); err != nil {
					return nil, err
				}
			}
		}
		fields, err := r.array(iface, "fields")
		if err != nil {
			return nil, err
		}
		for _, o := range fields {
			f, err := r.field(o)
			if err != nil {
				return nil, err
			}
			if err := r.appendField(s, f); err != nil {
				return nil, err
			}
		}
	}

	for i := 1; i < depth; i++ {
		s = &Schema{Slice: true, Item: &Field{Type: FormType, Schema: s}}
	}
	s.Meta = meta
	return s, nil
}

// field reads one field of an interface. Only fields are a form: a method, a
// getter, a setter or a prop describes behaviour, and there is nothing to edit
// in it.
func (r *reader) field(o gad.Object) (*Field, error) {
	name, err := r.index(o, "name")
	if err != nil {
		return nil, err
	}
	types, err := r.array(o, "types")
	if err != nil {
		return nil, err
	}
	meta, err := r.meta(o)
	if err != nil {
		return nil, err
	}

	f := &Field{Name: name.ToString(), Meta: meta}
	// `?` after the name has no reflection key of its own; it is the field's.
	if gf, ok := o.(*gad.InterfaceField); ok {
		f.Nullable = gf.Nullable
	}
	if err := r.typeInto(f, types); err != nil {
		return nil, fmt.Errorf("%s: %w", f.Name, err)
	}
	rawMeta, err := r.index(o, "@meta")
	if err != nil {
		return nil, err
	}
	if err := options(f, rawMeta); err != nil {
		return nil, fmt.Errorf("%s: %w", f.Name, err)
	}
	return f, nil
}

// MetaOptions is the metadata of a field that gives the values it may hold,
// with their labels — a select, whatever the field's type:
//
//	[options=(;opt1="option 1", opt2="option 2")] value str
//	[options=[["opt1", "option 1"], ["opt2", "option 2"]]] value str
//	[options=["opt1", "opt2"]] value str
//
// in the order written: a key-value array (the key the value, the value its
// label), or an array of pairs [value, label] — or of values alone, each its
// own label. The value is what the record holds, as text.
const MetaOptions = "options"

// MetaOpen is the bound of a range (Range[T]) that may be left blank:
// `[open="to"]` — from a date on, with no end —, `[open="from"]`.
const MetaOpen = "open"

// The limits of a value (Field.Min…): `[min=0, max=10, step=0.5]` of a
// number; `[min="2026-01-01", max="2026-12-31"]` of a date or a time, as
// text; `[minlength=3, maxlength=200]` of a text.
const (
	MetaMin       = "min"
	MetaMax       = "max"
	MetaStep      = "step"
	MetaMinLength = "minlength"
	MetaMaxLength = "maxlength"
	MetaPattern   = "pattern"
)

// NumberTypes and DateTypes are the types that take [min, max, step];
// TextTypes, [minlength, maxlength].
var (
	NumberTypes = map[string]bool{"int": true, "uint": true, "float": true, "decimal": true}
	DateTypes   = map[string]bool{"date": true, "time": true, "calendarDate": true, "calendarTime": true}
	TextTypes   = map[string]bool{DefaultType: true, "text": true}
)

// limits reads a field's [min, max, step, minlength, maxlength] — a range's
// into its bounds.
func limits(f *Field, key string, v gad.Object) error {
	target := f
	if f.Type == RangeType && f.Range != nil {
		target = f.Range
	}
	if f.Schema != nil && f.Schema.Slice {
		// a list: [min, max] are how many items it holds
		if key != MetaMin && key != MetaMax {
			return fmt.Errorf("%s: a list has [min, max], how many items it holds; not %s", key, key)
		}
		n, ok := v.(gad.Int)
		if !ok || n < 0 {
			return fmt.Errorf("%s: of a list, want a whole number of items, 0 or more, got %s", key, v.ToString())
		}
		if key == MetaMin {
			f.MinItems = int(n)
		} else {
			f.MaxItems = int(n)
		}
		if f.MaxItems > 0 && f.MinItems > f.MaxItems {
			return fmt.Errorf("%s=%d is more than %s=%d", MetaMin, f.MinItems, MetaMax, f.MaxItems)
		}
		return nil
	}
	if target.Enum != nil || target.Schema != nil {
		return fmt.Errorf("%s: only a number, a date, a time or a text has it", key)
	}
	t := target.Type
	switch key {
	case MetaPattern:
		if t != DefaultType {
			return fmt.Errorf("%s: only a line of text (str) has it, not %s", key, t)
		}
		str, ok := v.(gad.Str)
		if !ok || str == "" {
			return fmt.Errorf("%s: want a regular expression, as text", key)
		}
		if _, err := PatternRegexp(string(str)); err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
		target.Pattern = string(str)
		return nil
	case MetaPlaceholder:
		if !TextTypes[t] && !NumberTypes[t] {
			return fmt.Errorf("%s: only a text or a number has it, not %s", key, t)
		}
		return nil
	case MetaMinLength, MetaMaxLength:
		if !TextTypes[t] {
			return fmt.Errorf("%s: only a text (str, text) has it, not %s", key, t)
		}
		n, ok := v.(gad.Int)
		if !ok || n < 0 {
			return fmt.Errorf("%s: want a whole number, 0 or more, got %s", key, v.ToString())
		}
		if key == MetaMinLength {
			target.MinLength = int(n)
		} else {
			target.MaxLength = int(n)
		}
		if target.MaxLength > 0 && target.MinLength > target.MaxLength {
			return fmt.Errorf("%s=%d is more than %s=%d", MetaMinLength, target.MinLength, MetaMaxLength, target.MaxLength)
		}
		return nil
	}
	var s string
	switch {
	case NumberTypes[t] || key == MetaStep && DateTypes[t]:
		n, ok := number(v)
		if !ok {
			return fmt.Errorf("%s: want a number, got %s", key, v.ToString())
		}
		if key == MetaStep && n <= 0 {
			return fmt.Errorf("%s: want a number more than 0, got %s", key, v.ToString())
		}
		if (t == "int" || t == "uint") && n != float64(int64(n)) {
			return fmt.Errorf("%s=%s: %s is a whole number", key, v.ToString(), t)
		}
		s = strconv.FormatFloat(n, 'f', -1, 64)
	case DateTypes[t]:
		str, ok := v.(gad.Str)
		if !ok {
			return fmt.Errorf("%s: a date or a time is given as text (%s=\"2026-01-30\"), got %s", key, key, v.ToString())
		}
		s = string(str)
	default:
		return fmt.Errorf("%s: only a number, a date or a time has it, not %s", key, t)
	}
	switch key {
	case MetaMin:
		target.Min = s
	case MetaMax:
		target.Max = s
	default:
		target.Step = s
	}
	if target.Min != "" && target.Max != "" && LimitAfter(t, target.Min, target.Max) {
		return fmt.Errorf("%s=%s is after %s=%s", MetaMin, target.Min, MetaMax, target.Max)
	}
	return nil
}

// PatternRegexp is the regular expression of a [pattern=…]: matching the
// whole value, as HTML5 does. What Go and the browser read differently is
// refused: a group of flags or a named one (`(?…`, but `(?:…)`), lookaround, a
// backreference, `\A`, `\z`, `\Q…\E`.
func PatternRegexp(p string) (*regexp.Regexp, error) {
	if strings.Contains(strings.ReplaceAll(p, "(?:", ""), "(?") {
		return nil, fmt.Errorf("%q: no (?…) but (?:…) — flags, named groups, lookaround —: the browser reads them differently", p)
	}
	for _, esc := range []string{`\A`, `\z`, `\Q`, `\E`} {
		if strings.Contains(p, esc) {
			return nil, fmt.Errorf("%q: no %s, the browser has none (the pattern matches the whole value already)", p, esc)
		}
	}
	re, err := regexp.Compile("^(?:" + p + ")$")
	if err != nil {
		return nil, fmt.Errorf("%q is no regular expression: %w", p, err)
	}
	return re, nil
}

// LimitAfter reports whether the limit a comes after b, of the type t: numbers
// by value, dates and times (ISO texts) by their text.
func LimitAfter(t, a, b string) bool {
	if NumberTypes[t] {
		x, _ := strconv.ParseFloat(a, 64)
		y, _ := strconv.ParseFloat(b, 64)
		return x > y
	}
	return a > b
}

func number(v gad.Object) (float64, bool) {
	switch n := v.(type) {
	case gad.Int:
		return float64(n), true
	case gad.Uint:
		return float64(n), true
	case gad.Float:
		return float64(n), true
	case gad.Decimal:
		f, err := strconv.ParseFloat(n.ToString(), 64)
		return f, err == nil
	}
	return 0, false
}

// MetaValidation is the metadata of a field, or of a class, that checks its
// value: a function, or an array of them —
//
//	[validation=cep] zip str
//	[validation=[notWeekend, notHoliday]] day date
//	[validation=func(r) { … }] class Address { … }
//
// — accumulated: a class's after its parents', a field of a class's before
// its class's. What calls them, and with what, is the application's.
const MetaValidation = "validation"

// validations reads a `[validation=…]`: a function, or an array of them.
func validations(v gad.Object) ([]gad.Object, error) {
	var out []gad.Object
	add := func(o gad.Object) error {
		if _, ok := o.(gad.CallerObject); !ok {
			return fmt.Errorf("%s: want a function, or an array of them, got %s", MetaValidation, o.Type().Name())
		}
		out = append(out, o)
		return nil
	}
	if arr, ok := v.(gad.Array); ok {
		for _, o := range arr {
			if err := add(o); err != nil {
				return nil, err
			}
		}
		return out, nil
	}
	return out, add(v)
}

// classValidations are the `[validation=…]` of the class c's metadata.
func classValidations(c *gad.Class) ([]gad.Object, error) {
	for _, kv := range c.Meta {
		if kv.K.ToString() == MetaValidation {
			return validations(kv.V)
		}
	}
	return nil, nil
}

// The types of a file and of an image: their value, a file sent.
const (
	FileType  = "file"
	ImageType = "image"
)

// MetaAccept and MetaMaxSize are the types a file or an image takes and its
// largest size (Field.Accept, Field.MaxSize).
const (
	MetaAccept  = "accept"
	MetaMaxSize = "maxSize"
)

// fileMeta reads a file's or an image's [accept=…] and [maxSize=…] — a
// list's, its items'.
func fileMeta(f *Field, key string, v gad.Object) error {
	target := f
	if f.Schema != nil && f.Schema.Slice && f.Schema.Item != nil {
		target = f.Schema.Item
	}
	if target.Type != FileType && target.Type != ImageType {
		return fmt.Errorf("%s: only a file or an image (or a list of them) has it, not %s", key, target.Type)
	}
	if key == MetaAccept {
		str, ok := v.(gad.Str)
		if !ok || strings.TrimSpace(string(str)) == "" {
			return fmt.Errorf("%s: want the types it takes, as text: \"application/pdf,.docx\"", key)
		}
		for _, a := range strings.Split(string(str), ",") {
			a = strings.TrimSpace(a)
			if a == "" || !strings.HasPrefix(a, ".") && !strings.Contains(a, "/") {
				return fmt.Errorf("%s: %q is neither a media type (\"image/png\", \"image/*\") nor an extension (\".pdf\")", key, a)
			}
			if target.Type == ImageType && !strings.HasPrefix(a, ".") && !strings.HasPrefix(a, "image/") {
				return fmt.Errorf("%s: an image takes images only, not %q", key, a)
			}
		}
		target.Accept = string(str)
		return nil
	}
	n, err := ParseSize(v)
	if err != nil {
		return fmt.Errorf("%s: %w", key, err)
	}
	target.MaxSize = n
	return nil
}

// ParseSize reads a size: a number of bytes, or a text of one with its
// unit — "500KB", "5MB", "1GB" (of 1024).
func ParseSize(v gad.Object) (int64, error) {
	switch n := v.(type) {
	case gad.Int:
		if n > 0 {
			return int64(n), nil
		}
	case gad.Uint:
		if n > 0 {
			return int64(n), nil
		}
	case gad.Str:
		s := strings.ToUpper(strings.TrimSpace(string(n)))
		mult := int64(1)
		for _, u := range []struct {
			suffix string
			mult   int64
		}{{"GB", 1 << 30}, {"MB", 1 << 20}, {"KB", 1 << 10}, {"B", 1}} {
			if strings.HasSuffix(s, u.suffix) {
				s, mult = strings.TrimSpace(strings.TrimSuffix(s, u.suffix)), u.mult
				break
			}
		}
		if f, err := strconv.ParseFloat(s, 64); err == nil && f > 0 {
			return int64(f * float64(mult)), nil
		}
	}
	return 0, fmt.Errorf("want a size: a number of bytes, or \"500KB\", \"5MB\"; got %s", v.ToString())
}

// The bounds of a range: its value's keys.
const (
	RangeFrom = "from"
	RangeTo   = "to"
)

// options reads the field's `[options=…]` into the enum of the values it may
// hold — the item's, for a list of plain values (each item a select).
func options(f *Field, rawMeta gad.Object) error {
	kva, _ := rawMeta.(gad.KeyValueArray)
	for _, kv := range kva {
		if kv.K.ToString() == MetaOpen {
			// a bound of a range that may be left blank
			if f.Type != RangeType {
				return fmt.Errorf("%s: only a range (Range[T]) has a bound left open", MetaOpen)
			}
			switch open := kv.V.ToString(); open {
			case RangeFrom, RangeTo:
				f.RangeOpen = open
			default:
				return fmt.Errorf("%s=%q: the bound left open is %q or %q", MetaOpen, open, RangeFrom, RangeTo)
			}
			continue
		}
		switch k := kv.K.ToString(); k {
		case MetaAccept, MetaMaxSize:
			if err := fileMeta(f, kv.K.ToString(), kv.V); err != nil {
				return err
			}
			continue
		case MetaValidation:
			fns, err := validations(kv.V)
			if err != nil {
				return err
			}
			f.Validations = append(f.Validations, fns...)
			continue
		case MetaMin, MetaMax, MetaStep, MetaMinLength, MetaMaxLength, MetaPattern, MetaPlaceholder:
			if err := limits(f, k, kv.V); err != nil {
				return err
			}
			continue
		}
		if kv.K.ToString() != MetaOptions {
			continue
		}
		e, err := optionsEnum(kv.V)
		if err != nil {
			return fmt.Errorf("%s: %w", MetaOptions, err)
		}
		target := f
		if f.Schema != nil && f.Schema.Slice && f.Schema.Item != nil && f.Schema.Item.Schema == nil {
			target = f.Schema.Item
		} else if f.Type == RangeType && f.Range != nil {
			// a range of choices: its bounds are of them
			target = f.Range
		} else if f.Schema != nil {
			return fmt.Errorf("%s: a record has no options; only a value has", MetaOptions)
		}
		target.Enum = e
	}
	// (choices, whatever the order of the meta: no limits)
	for _, l := range []*Field{f, f.Range} {
		if l != nil && l.Enum != nil && (l.Min != "" || l.Max != "" || l.Step != "" || l.MinLength > 0 || l.MaxLength > 0 || l.Pattern != "") {
			return fmt.Errorf("%s: only a number, a date, a time or a text has limits, not a choice", f.Name)
		}
	}
	return nil
}

// optionsEnum reads an `[options=…]` value into an enum, in its order.
func optionsEnum(v gad.Object) (*Enum, error) {
	e := &Enum{Options: true}
	add := func(value, label gad.Object) error {
		if value == gad.Nil {
			return fmt.Errorf("an option with no value")
		}
		it := EnumItem{Name: value.ToString()}
		if label != nil && label != gad.Nil {
			it.Label = label.ToString()
		}
		for _, n := range e.Names {
			if n == it.Name {
				return fmt.Errorf("the value %q twice", it.Name)
			}
		}
		e.Names = append(e.Names, it.Name)
		e.Items = append(e.Items, it)
		return nil
	}
	switch t := v.(type) {
	case gad.KeyValueArray:
		for _, kv := range t {
			if err := add(kv.K, kv.V); err != nil {
				return nil, err
			}
		}
	case gad.Array:
		for i, item := range t {
			var err error
			switch p := item.(type) {
			case gad.Array:
				if len(p) != 2 {
					return nil, fmt.Errorf("[%d]: a pair is [value, label], got %d items", i, len(p))
				}
				err = add(p[0], p[1])
			case *gad.KeyValue:
				err = add(p.K, p.V)
			default:
				err = add(p, nil)
			}
			if err != nil {
				return nil, fmt.Errorf("[%d]: %w", i, err)
			}
		}
	default:
		return nil, fmt.Errorf("want a key-value array or an array of [value, label], got %s", v.Type().Name())
	}
	if len(e.Names) == 0 {
		return nil, fmt.Errorf("no option")
	}
	return e, nil
}

// getter reads a property of an interface that is a getter alone — `get name
// Type` — as a read-only field of the type it returns; nil for a property with
// a setter, which describes behaviour.
func (r *reader) getter(o gad.Object) (*Field, error) {
	p, ok := o.(*gad.InterfaceProp)
	if !ok || p.Getter == nil || len(p.Setters) > 0 {
		return nil, nil
	}
	meta, err := r.meta(o)
	if err != nil {
		return nil, err
	}
	f := &Field{Name: p.Name, Meta: meta, ReadOnly: true}
	// the getter returns one value, a typedIdent: its types, resolved
	var types gad.Array
	if len(p.Getter.Return) > 0 {
		if types, err = r.array(p.Getter.Return[0], "types"); err != nil {
			return nil, fmt.Errorf("%s: %w", f.Name, err)
		}
	}
	if err := r.typeInto(f, types); err != nil {
		return nil, fmt.Errorf("%s: %w", f.Name, err)
	}
	return f, nil
}

// typeInto reads what a field's types are into it: the component it is drawn
// with, and — for an interface, an array or an enum — the rest of what it
// holds.
func (r *reader) typeInto(f *Field, types gad.Array) error {
	switch len(types) {
	case 0:
		f.Type = DefaultType
		return nil
	case 1:
	default:
		// a union of classes is a choice of one of them: by its name, or with
		// its fields
		if r.choiceAsName {
			if e := classNames(types); e != nil {
				f.Type, f.Enum = e.Name, e
				return nil
			}
		}
		if choice, err := r.choice(types); err != nil || choice != nil {
			if err != nil {
				return err
			}
			f.Type, f.Schema = FormType, choice
			return nil
		}
		// Several types describe no single input: the field is reported as the
		// union it is written as, which no component is registered for.
		names := make([]string, len(types))
		for i, t := range types {
			names[i] = typeName(t)
		}
		f.Type = strings.Join(names, "|")
		return nil
	}

	if r.typeOf != nil {
		if name := r.typeOf(types[0]); name != "" {
			f.Type = name
			return nil
		}
	}

	switch t := types[0].(type) {
	case *gad.Class:
		// a record of its own: `a class { … }`, or a class by name
		sub, err := r.class(t)
		if err != nil {
			return err
		}
		f.Type, f.Schema = FormType, sub
	case *gad.Interface:
		sub, err := r.schema(t)
		if err != nil {
			return err
		}
		f.Type, f.Schema = FormType, sub
	case *gad.ArrayType:
		sub, err := r.arrayType(t)
		if err != nil {
			return err
		}
		if sub == nil {
			f.Type = t.Name() // several element types: no single input
			return nil
		}
		f.Type, f.Schema = FormType, sub
	case *gad.Enum:
		e, err := r.enum(t)
		if err != nil {
			return err
		}
		f.Type, f.Enum = e.Name, e
	case *gad.RangeOfType:
		// Range[T]: from and to, of T — a union of numbers, a number
		bound := &Field{}
		if u, ok := t.Elem.(*gad.TypeUnion); ok {
			if numericUnion(u) {
				bound.Type = "float"
			} else if err := r.typeInto(bound, u.Types); err != nil {
				return err
			}
		} else if err := r.typeInto(bound, []gad.Object{t.Elem}); err != nil {
			return err
		}
		f.Type, f.Range = RangeType, bound
	default:
		f.Type = typeName(t)
	}
	return nil
}

// arrayType reads `tags []str` into a list of plain values, one schema per
// `[]`. Several element types (`[]<int|str>`) describe no single input: nil.
func (r *reader) arrayType(t *gad.ArrayType) (*Schema, error) {
	depth, err := r.depth(t)
	if err != nil {
		return nil, err
	}
	elem, err := r.array(t, "@elem")
	if err != nil {
		return nil, err
	}
	if len(elem) != 1 {
		return nil, nil
	}

	item := &Field{}
	if err := r.typeInto(item, elem); err != nil {
		return nil, err
	}
	s := &Schema{Slice: true, Item: item}
	for i := 1; i < depth; i++ {
		s = &Schema{Slice: true, Item: &Field{Type: FormType, Schema: s}}
	}
	return s, nil
}

// enum reads an enum — declared beside the interface or inline — into the
// values a field may hold, by NAME: the name is what the form stores.
func (r *reader) enum(e *gad.Enum) (*Enum, error) {
	if known, ok := r.enums[e.EnumName]; ok {
		return known, nil
	}
	names, err := r.array(e, "@names")
	if err != nil {
		return nil, err
	}
	out := &Enum{Name: e.EnumName}
	if m := e.Module; m != nil && m.Name != "" && m.InitCompiledFunc == nil && !m.IsMain() {
		// an enum of a module of gad's own, declared in Go (time.Months) — not
		// one of the code's, which runs (InitCompiledFunc) —: its words, the
		// module's
		out.Module = e.Module.Name
	}
	for _, n := range names {
		out.Names = append(out.Names, n.ToString())
	}
	r.enums[out.Name] = out
	return out, nil
}

// numericUnion reports whether every type of u is a number's: a range of
// them is drawn as numbers.
func numericUnion(u *gad.TypeUnion) bool {
	for _, t := range u.Types {
		switch typeName(t) {
		case "int", "uint", "float", "decimal":
		default:
			return false
		}
	}
	return len(u.Types) > 0
}

// typeName is the name a type is drawn by. `any` is what gad makes of a field
// the schema left untyped, and an untyped field is plain text.
func typeName(t gad.Object) string {
	switch t := t.(type) {
	case *typeMarker:
		return t.name
	case gad.ObjectType:
		if name := t.Name(); name != "any" {
			return name
		}
		return DefaultType
	}
	return t.ToString()
}

// metaOf reads a `[k=v, …]` block into Go values.
func metaOf(o gad.Object) Meta {
	kva, ok := o.(gad.KeyValueArray)
	if !ok || len(kva) == 0 {
		return nil
	}
	m := make(Meta, len(kva))
	for _, kv := range kva {
		m[kv.K.ToString()] = metaValue(kv.V)
	}
	return m
}

func metaValue(o gad.Object) any {
	switch v := o.(type) {
	case gad.Str:
		return string(v)
	case gad.RawStr:
		return string(v)
	case gad.Int:
		return int64(v)
	case gad.Uint:
		return int64(v)
	case gad.Float:
		return float64(v)
	case gad.Bool:
		return bool(v)
	case gad.Flag:
		return bool(v)
	case gad.Array:
		out := make([]any, len(v))
		for i, e := range v {
			out[i] = metaValue(e)
		}
		return out
	case gad.KeyValueArray:
		return metaOf(v)
	case gad.Dict:
		out := make(Meta, len(v))
		for k, e := range v {
			out[k] = metaValue(e)
		}
		return out
	}
	if o == gad.Nil {
		return nil
	}
	return o.ToString()
}

// readLayout reads how the schema is drawn from its metadata — `layout`, a
// layout's name or `{name: …, <its config>}` — and refuses what cannot be
// drawn: an unknown layout, a list layout of something that is not a list of
// records, a config the layout's Config does not declare or its checks refuse.
func (s *Schema) readLayout() error {
	name, cfg, err := layoutMeta(s.Meta)
	if err != nil {
		return err
	}
	return s.setLayout(name, cfg)
}

// layoutMeta is the layout the metadata names, and the config it gives it.
func layoutMeta(meta Meta) (name string, cfg Meta, err error) {
	for _, k := range []string{"columns", "num_columns"} {
		if _, ok := meta[k]; ok {
			return "", nil, fmt.Errorf(`schemaform: %q is the config of a layout: write it in layout={name: "…", %s: …}`, k, k)
		}
	}
	switch v := meta["layout"].(type) {
	case nil:
	case string:
		name = v
	case Meta:
		cfg = make(Meta, len(v))
		for k, x := range v {
			if k == "name" {
				name, _ = x.(string)
				continue
			}
			cfg[k] = x
		}
	default:
		return "", nil, fmt.Errorf("schemaform: layout is a name or {name: …, …}, not %T", v)
	}
	return name, cfg, nil
}

// setLayout gives the schema the layout name, told cfg.
func (s *Schema) setLayout(name string, cfg Meta) error {
	if name == "" {
		name = LayoutForm
	}
	l := LookupLayout(name)
	if l == nil {
		return fmt.Errorf("schemaform: layout %q does not exist (there is %q)", name, strings.Join(LayoutNames(), `", "`))
	}
	if l.Records && (!s.Slice || s.Item != nil) {
		return fmt.Errorf("schemaform: layout %q draws a list of records, and this is not one", name)
	}
	config, err := l.read(s, cfg)
	if err != nil {
		return fmt.Errorf("schemaform: layout %q: %w", name, err)
	}
	s.Layout, s.LayoutConfig = name, config
	return nil
}

// shareEnums hands the enums the schema's fields hold to every schema in it.
func (s *Schema) shareEnums(enums map[string]*Enum) {
	s.Enums = enums
	for _, f := range append(append([]*Field{}, s.Fields...), s.Item) {
		if f != nil && f.Schema != nil {
			f.Schema.shareEnums(enums)
		}
	}
}
