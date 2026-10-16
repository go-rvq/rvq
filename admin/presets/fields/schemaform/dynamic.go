package schemaform

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gad-lang/gad"
)

// The metadata of a field drawn by what is answered: each a value, or a
// function `(this, parents) => …` — this the record the field is of, an
// instance of its class with the values answered so far (typed); parents
// the records above it, from the nearest, each `{fieldName, self}`: the
// field of the parent the record is, and the parent's instance.
//
//		[disabled=(this) => !this.state] city str
//		[if=(this, parents) => parents[0].self.kind == "house"] rooms int
//
//	  - if: the field is there only when true (Vue's v-if);
//	  - show: it is shown only when true — hidden, it stays in the page (CSS
//	    display);
//	  - disabled: it is drawn, not edited;
//	  - readonly: it is drawn, its value not changed;
//	  - class: the CSS classes of its box (a str, an array, a dict of a class
//	    and when it holds).
//
// A field not there, hidden or disabled is not asked: not required, its
// value dropped. A field read-only is not required.
const (
	MetaIf       = "if"
	MetaShow     = "show"
	MetaDisabled = "disabled"
	MetaReadonly = "readonly"
	MetaClass    = "class"
	// MetaOptionsError is what a field does when the values of its
	// `[options=(this, parents) => …]` could not be had (a request that
	// failed): OptionsErrorRefuse (the default), the value refused;
	// OptionsErrorWarn, the value kept as given and a warning of it put in
	// the record (WarningsField).
	MetaOptionsError   = "optionsError"
	OptionsErrorRefuse = "error"
	OptionsErrorWarn   = "warn"
	// WarningsField is the field of a record holding its warnings: a list of
	// {field, message}.
	WarningsField = "$warnings"
)

// MetaChanged is what is done as a field's value changes, in the page:
// `[changed=(this, parents) => …]` sets fields of this — its record —, at
// once or once a Request it answers is resolved (its thens set them); the
// fields it set, the page takes (Decoder.Changed). The server takes the
// values as posted: it checks them as any.
//
//	[changed=(this, parents) => len(this.zip) >= 5 ? Request("https://api.zippopotam.us/us/" + this.zip[:5]).then(func(r) {
//	    this.state = r.places[0].state
//	    this.city = r.places[0]["place name"]
//	}) : nil]
//	zip str
const MetaChanged = "changed"

// MetaRound is a float's or a decimal's own rounding: `[round=(this, parents,
// value; places=0) => …]`, called with the value read — places its
// `[places=N]`, 0 without —, what it answers kept as the value (a number, or
// its text).
//
//	[places=2, round=(this, parents, value; places=0) => float(int(value * 100 + 0.5)) / 100] price decimal
const MetaRound = "round"

// MetaRoundDisplay is a float's or a decimal's rounding as it is shown —
// the detail, a read-only form, DisplayText —, the value kept as it is:
// `[roundDisplay=(this, parents, value; places=0) => …]`, as MetaRound.
//
//	[places=2, roundDisplay=up] price decimal
const MetaRoundDisplay = "roundDisplay"

// CondMetas are the metadata of a field drawn by what is answered.
var CondMetas = []string{MetaIf, MetaShow, MetaDisabled, MetaReadonly, MetaClass}

// Warning is a value kept although it could not be checked (OptionsErrorWarn):
// the field's path, and why.
type Warning struct {
	Path    string `json:"path"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FieldState is how a field is drawn, by what is answered: its if, show,
// disabled, readonly and class; of a field of `[options=(this, parents) =>
// …]`, its options — Items when the function answered them, none when it
// answered nil (NoOptions), Request (its key) when they are to be fetched
// (Decoder.Options).
type FieldState struct {
	If        bool       `json:"if"`
	Show      bool       `json:"show"`
	Disabled  bool       `json:"disabled"`
	Readonly  bool       `json:"readonly"`
	Class     any        `json:"class,omitempty"`
	Dynamic   bool       `json:"dynamic,omitempty"`
	Items     []EnumItem `json:"items,omitempty"`
	NoOptions bool       `json:"noOptions,omitempty"`
	Request   string     `json:"request,omitempty"`
	Error     string     `json:"error,omitempty"`
}

// Asked says the field is asked: there, shown and not disabled.
func (s FieldState) Asked() bool { return s.If && s.Show && !s.Disabled }

// IsDynamic says the field is drawn by what is answered: a metadata of
// CondMetas, or options of a function.
func (f *Field) IsDynamic() bool { return len(f.Cond) > 0 || f.OptionsFunc != nil }

// CondStatic is the field's metadata of CondMetas a function holds, by name —
// what the page asks the state of (Decoder.State) —; a static one, its
// value.
func (f *Field) CondStatic(key string) (any, bool) {
	v, ok := f.Cond[key]
	if !ok || isCallable(v) {
		return nil, false
	}
	return metaValue(v), true
}

// dynamic says some field of the schema, or of a record of it, is drawn by
// what is answered.
func (s *Schema) dynamic() bool {
	if s == nil {
		return false
	}
	for _, f := range s.Fields {
		if f.IsDynamic() || (f.Schema != nil && f.Schema != s && f.Schema.dynamic()) {
			return true
		}
		if f.Schema != nil && f.Schema.Item != nil && f.Schema.Item.Schema.dynamic() {
			return true
		}
	}
	return false
}

// scoped says the reading of the schema needs the records as this and
// parents: a field drawn by what is answered (dynamic), or one of its own
// rounding (MetaRound).
func (s *Schema) scoped() bool {
	return s.dynamic() || s.rounds()
}

// rounds says some field of the schema, or of a record of it, has its own
// rounding (MetaRound).
func (s *Schema) rounds() bool {
	if s == nil {
		return false
	}
	for _, f := range s.Fields {
		if f.RoundFunc != nil || (f.Schema != nil && f.Schema != s && f.Schema.rounds()) {
			return true
		}
		if f.Schema != nil && f.Schema.Item != nil && f.Schema.Item.Schema.rounds() {
			return true
		}
	}
	return false
}

// isCallable says v is a function a schema wrote (not a type).
func isCallable(v gad.Object) bool {
	if _, isType := v.(gad.ObjectType); isType {
		return false
	}
	if _, isClass := v.(*gad.Class); isClass {
		return false
	}
	_, ok := v.(gad.CallerObject)
	return ok
}

// condMeta reads a metadata of CondMetas, or `[optionsError=…]`, of f.
func condMeta(f *Field, key string, v gad.Object) error {
	if key == MetaOptionsError {
		switch s := v.ToString(); s {
		case OptionsErrorRefuse, OptionsErrorWarn:
			f.OptionsError = s
			return nil
		default:
			return fmt.Errorf("%s=%q: %q or %q", MetaOptionsError, s, OptionsErrorRefuse, OptionsErrorWarn)
		}
	}
	if f.Cond == nil {
		f.Cond = map[string]gad.Object{}
	}
	f.Cond[key] = v
	return nil
}

// scope is a record being read: its class, its values typed, the field of
// its parent it is, its parent.
type scope struct {
	class  *gad.Class
	fields []*Field
	typed  map[string]any
	name   string
	parent *scope
	this   gad.Object
}

// child is the scope of the record of the field name of s — its fields,
// of class —, typed as typed.
func (s *scope) child(name string, class *gad.Class, fields []*Field, typed any) *scope {
	return &scope{class: class, fields: fields, typed: asMap(typed), name: name, parent: s}
}

func asMap(v any) map[string]any {
	switch x := v.(type) {
	case Record:
		return x.Map()
	case map[string]any:
		return x
	}
	return nil
}

// This is the record as its class's instance — its values typed set —; of
// no class (or one its values do not fit), a dict of them.
func (s *scope) This(d *Decoder) gad.Object {
	if s.this != nil {
		return s.this
	}
	values := gad.Dict{}
	for k, v := range s.typed {
		if v != nil && !strings.HasPrefix(k, "$") {
			values[k] = typedFieldValue(fieldNamed(s.fields, k), v)
		}
	}
	s.this = values
	if s.class != nil && d.VM != nil {
		if obj, err := gad.NewInvoker(d.VM, s.class).Invoke(gad.Args{}, nil); err == nil {
			if inst, ok := obj.(*gad.ClassInstance); ok && inst.Init(d.VM, values) == nil {
				s.this = inst
			}
		}
	}
	return s.this
}

// Parents are the records above s, from the nearest: {fieldName, self}.
func (s *scope) Parents(d *Decoder) gad.Array {
	var out gad.Array
	for c := s; c != nil && c.parent != nil; c = c.parent {
		out = append(out, gad.Dict{"fieldName": gad.Str(c.name), "self": c.parent.This(d)})
	}
	if out == nil {
		out = gad.Array{}
	}
	return out
}

// evalMeta is the metadata v of f: a function called with this and parents
// of the record being read; a value, itself.
func (d *Decoder) evalMeta(f *Field, key string, v gad.Object) (gad.Object, error) {
	if !isCallable(v) {
		return v, nil
	}
	sc := d.scope
	if sc == nil {
		sc = &scope{}
	}
	ret, err := d.call(v, sc.This(d), sc.Parents(d))
	if err != nil {
		return nil, fmt.Errorf("%s: [%s]: %w", f.Name, key, err)
	}
	return ret, nil
}

// state is how f is drawn in the record being read (d.scope).
func (d *Decoder) state(f *Field) FieldState {
	st := FieldState{If: true, Show: true}
	for _, key := range CondMetas {
		v, ok := f.Cond[key]
		if !ok {
			continue
		}
		ret, err := d.evalMeta(f, key, v)
		if err != nil {
			log.Printf("ERROR: form: %v", err)
			st.Error = err.Error()
			continue
		}
		switch key {
		case MetaIf:
			st.If = !ret.IsFalsy()
		case MetaShow:
			st.Show = !ret.IsFalsy()
		case MetaDisabled:
			st.Disabled = !ret.IsFalsy()
		case MetaReadonly:
			st.Readonly = !ret.IsFalsy()
		case MetaClass:
			st.Class = metaValue(ret)
		}
	}
	return st
}

// optionsOf is what the `[options=(this, parents) => …]` of f answers in
// the record being read: nil (no options), the values, or a request.
func (d *Decoder) optionsOf(f *Field) (gad.Object, error) {
	ret, err := d.evalMeta(f, MetaOptions, f.OptionsFunc)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		ret = gad.Nil
	}
	return ret, nil
}

// itemsOf are the values ret — what options answered, resolved — gives.
func itemsOf(ret gad.Object) ([]EnumItem, error) {
	if ret == nil || ret == gad.Nil {
		return []EnumItem{}, nil
	}
	if a, ok := ret.(gad.Array); ok && len(a) == 0 {
		return []EnumItem{}, nil
	}
	e, err := optionsEnum(ret)
	if err != nil {
		return nil, err
	}
	return e.Items, nil
}

// dynamicItems are the values f may hold, by its options' function: the
// page's (PageOptions), when it says them; else what the function answers —
// a request resolved (not by the page's check: NoFetch, any value then) —.
// Values that could not be had: an error, or — `[optionsError="warn"]` — a
// warning, and any value.
func (d *Decoder) dynamicItems(path string, f *Field) ([]EnumItem, error) {
	if d.PageOptions != nil {
		if items, ok := d.PageOptions[path]; ok {
			return items, nil
		}
	}
	ret, err := d.optionsOf(f)
	if err == nil {
		if req, ok := ret.(*Request); ok {
			if d.NoFetch {
				return nil, nil
			}
			ret, err = req.Resolve(d.Ctx, d.VM)
		}
	}
	var items []EnumItem
	if err == nil {
		items, err = itemsOf(ret)
	}
	if err != nil {
		if f.OptionsError == OptionsErrorWarn {
			log.Printf("WARNING: form: the options of %s: %v", path, err)
			d.warnings = append(d.warnings, Warning{Path: path, Field: f.Name,
				Message: "the options could not be checked (" + err.Error() + "): the value was kept as given"})
			return nil, nil
		}
		return nil, err
	}
	return items, nil
}

// Warnings are the values kept although they could not be checked, of the
// last read.
func (d *Decoder) Warnings() []Warning { return d.warnings }

// itemScope makes the scope the record i of the list f, of class; restore
// undoes it.
func (d *Decoder) itemScope(f *Field, class *gad.Class, fields []*Field, i int) (restore func()) {
	outer := d.scope
	if outer == nil {
		return func() {}
	}
	items, _ := outer.typed[f.Name].([]any)
	var typed any
	if i < len(items) {
		typed = items[i]
	}
	d.scope = outer.child(f.Name+"["+strconv.Itoa(i)+"]", class, fields, typed)
	return func() { d.scope = outer }
}

// warningsOf are the warnings of the fields of the record at path.
func (d *Decoder) warningsOf(path string) []any {
	var out []any
	for _, w := range d.warnings {
		parent := ""
		if i := strings.LastIndex(w.Path, "."); i >= 0 {
			parent = w.Path[:i]
		}
		if parent == path {
			out = append(out, map[string]any{"fieldName": w.Field, "message": w.Message})
		}
	}
	return out
}

// rootScope is the scope of the form's record: its class, data typed.
func (d *Decoder) rootScope(data map[string]any) *scope {
	var (
		class  *gad.Class
		fields []*Field
	)
	if d.Schema != nil {
		class, fields = d.Schema.recClass, d.Schema.Fields
	}
	return &scope{class: class, fields: fields, typed: d.Typed(data).Map()}
}

// State is how each field drawn by what is answered (IsDynamic) is, by
// data — the values the page holds, as Unflatten makes them —: by its path.
// Options of a request are not fetched: the page asks them (Options).
func (d *Decoder) State(data map[string]any) map[string]FieldState {
	d.prepare()
	defer d.done()
	out := map[string]FieldState{}
	d.scope = d.rootScope(data)
	d.walk(d.Schema.Fields, "", func(f *Field, path string) {
		if !f.IsDynamic() {
			return
		}
		st := FieldState{If: true, Show: true}
		if len(f.Cond) > 0 {
			st = d.state(f)
		}
		if f.OptionsFunc != nil {
			st.Dynamic = true
			ret, err := d.optionsOf(f)
			switch {
			case err != nil:
				st.Error = err.Error()
			case ret == gad.Nil:
				st.NoOptions = true
			default:
				if req, ok := ret.(*Request); ok {
					st.Request = req.key()
				} else if items, err := itemsOf(ret); err != nil {
					st.Error = err.Error()
				} else {
					st.Items = items
				}
			}
		}
		out[path] = st
	})
	return out
}

// Options are the values of the field at path by its options' function —
// a request fetched — for data (as State).
func (d *Decoder) Options(path string, data map[string]any) ([]EnumItem, error) {
	d.prepare()
	defer d.done()
	var (
		items []EnumItem
		err   = fmt.Errorf("form: no field %q of options of a function", path)
	)
	d.scope = d.rootScope(data)
	d.walk(d.Schema.Fields, "", func(f *Field, p string) {
		if p != path || f.OptionsFunc == nil {
			return
		}
		var ret gad.Object
		if ret, err = d.optionsOf(f); err != nil {
			return
		}
		if req, ok := ret.(*Request); ok {
			if ret, err = req.Resolve(d.Ctx, d.VM); err != nil {
				return
			}
		}
		items, err = itemsOf(ret)
	})
	return items, err
}

// Changed calls the `[changed=…]` of the field at path, by data (as
// State) — a request it answers fetched —: the fields of its record it set
// (this's), by name, and their values; nil when it set none.
func (d *Decoder) Changed(path string, data map[string]any) (map[string]any, error) {
	d.prepare()
	defer d.done()
	var (
		out map[string]any
		err = fmt.Errorf("form: no field %q of a changed", path)
	)
	d.scope = d.rootScope(data)
	d.walk(d.Schema.Fields, "", func(f *Field, p string) {
		if p != path || f.ChangedFunc == nil {
			return
		}
		sc := d.scope
		this := sc.This(d)
		before := d.recordValues(sc, this)
		var ret gad.Object
		if ret, err = d.evalMeta(f, MetaChanged, f.ChangedFunc); err != nil {
			return
		}
		if req, ok := ret.(*Request); ok {
			if _, err = req.Resolve(d.Ctx, d.VM); err != nil {
				return
			}
		}
		for name, v := range d.recordValues(sc, this) {
			if fmt.Sprint(v) != fmt.Sprint(before[name]) {
				if out == nil {
					out = map[string]any{}
				}
				// an int of places (fixed point): as the form holds it
				if f := fieldNamed(sc.fields, name); f != nil && v != nil {
					if p := fixedPlaces(f); p > 0 {
						v = FixedShown(v, p)
					}
				}
				out[name] = v
			}
		}
	})
	return out, err
}

// recordValues are the values of the fields of the record of sc — of a
// value, not of a record —, as this holds them.
func (d *Decoder) recordValues(sc *scope, this gad.Object) map[string]any {
	getter, _ := this.(gad.IndexGetter)
	out := map[string]any{}
	for _, f := range sc.fields {
		if f.Schema != nil || strings.HasPrefix(f.Name, "$") || getter == nil {
			continue
		}
		v, err := getter.IndexGet(d.VM, gad.Str(f.Name))
		if err != nil || v == nil {
			continue
		}
		out[f.Name] = metaValue(v)
	}
	return out
}

// walk calls fn with each field drawn by what is answered of fields — of
// the records in them too —, at its path, d.scope its record's.
func (d *Decoder) walk(fields []*Field, path string, fn func(f *Field, path string)) {
	for _, f := range fields {
		p := join(path, f.Name)
		if (f.IsDynamic() || f.ChangedFunc != nil) && f.Name != WarningsField {
			fn(f, p)
		}
		switch {
		case f.Schema != nil && f.Schema.Slice && f.Schema.Item == nil:
			d.walkList(f, f.Schema.Fields, f.Schema.recClass, p, fn)
		case f.Schema != nil && f.Schema.Slice && f.Schema.Item.Schema != nil:
			d.walkList(f, f.Schema.Item.Schema.Fields, f.Schema.Item.Schema.recClass, p, fn)
		case f.Schema != nil && !f.Schema.Slice && !f.Schema.Choice:
			outer := d.scope
			d.scope = outer.child(f.Name, f.Schema.recClass, f.Schema.Fields, outer.typed[f.Name])
			d.walk(f.Schema.Fields, p, fn)
			d.scope = outer
		}
	}
}

func (d *Decoder) walkList(f *Field, fields []*Field, class *gad.Class, path string, fn func(f *Field, path string)) {
	outer := d.scope
	items, _ := outer.typed[f.Name].([]any)
	for i, it := range items {
		d.scope = outer.child(f.Name+"["+strconv.Itoa(i)+"]", class, fields, it)
		d.walk(fields, path+"["+strconv.Itoa(i)+"]", fn)
	}
	d.scope = outer
}
