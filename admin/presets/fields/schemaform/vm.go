package schemaform

import (
	"errors"
	"fmt"
	"github.com/go-rvq/rvq/admin/presets"
	"regexp"
	"sort"
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

// Class reads a gad class as a form: a record of its fields, the way an
// interface's are read, with the default each one has. A field typed by a
// class — `a class { … }`, or one by name — is a record of its own, a group
// of the form; a type only the application knows is named by TypeOf. vm is
// the VM that ran the code the class comes from.
func (b *Builder) Class(vm *gad.VM, c *gad.Class) (*Schema, error) {
	r := &reader{vm: vm, enums: map[string]*Enum{}, visiting: map[*gad.Interface]bool{}, typeOf: b.typeOf, choiceAsName: b.choiceAsName}
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
	for _, cf := range c.RawFields() {
		types := make(gad.Array, len(cf.Types))
		for i, t := range cf.Types {
			types[i] = t
		}
		f := &Field{Name: cf.Name, Meta: metaOf(cf.Meta), Nullable: cf.Nullable}
		if err := r.typeInto(f, types); err != nil {
			return nil, fmt.Errorf("%s: %w", f.Name, err)
		}
		if cf.Value != nil && cf.Value != gad.Nil {
			f.Default = metaValue(cf.Value)
		}
		s.Fields = append(s.Fields, f)
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
				s.Fields = append(s.Fields, f)
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
			s.Fields = append(s.Fields, f)
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
	return f, nil
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
	for _, n := range names {
		out.Names = append(out.Names, n.ToString())
	}
	r.enums[out.Name] = out
	return out, nil
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
