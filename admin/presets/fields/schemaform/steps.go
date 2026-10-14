package schemaform

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gad-lang/gad"
	"github.com/go-rvq/rvq/web/str_utils"
)

// The metadata of a form of steps, before its class: `[steps] class Form
// { … }` — each field of the class a step: a field alone (a value, a class of
// fields under its name: `a.name`), or a group `{ … }` (its fields the
// record's own, no name added) —; `[steps, autosave]`, each step saved as it
// is passed (the application keeping it).
const (
	MetaSteps    = "steps"
	MetaAutosave = "autosave"
	// MetaVertical: the steps of the class down the page (`[steps,
	// vertical]`), as a group's.
	MetaVertical = "vertical"
	// MetaNext is the step after a step, by what was answered: a
	// `[next=(r) => …]` before it, or the method `next(r)` of its class; r
	// the answer so far, a dict. It answers the step — its class field,
	// `Form.b`, `Form.$2` (a group) —, or nil: the one after it in the
	// class.
	MetaNext = "next"
)

// StepTimeout bounds a step's next.
var StepTimeout = 5 * time.Second

// Step is a step of a form of steps: a field of the class, or a group.
type Step struct {
	// Key is the field's name — of a group, `$N` —, what `next` answers.
	Key string
	// Label is its `[label=…]`, else its name humanized (of a group, its
	// number).
	Label string
	// Group: a `{ … }` of the class, its fields the record's own.
	Group bool
	// Fields are its fields: the field, or the fields of the group.
	Fields []*Field
	// next is its `[next=…]`; class, the class of a field of one, whose
	// method `next` is the step after it.
	next  gad.Object
	class *gad.Class
}

// Flow is the steps of a form of steps, in order.
type Flow struct {
	Steps []*Step
	// Autosave: `[autosave]` — each step saved as it is passed.
	Autosave bool
	// Vertical: `[steps, vertical]` — the steps down the page, each over its
	// panel; across, else.
	Vertical bool
	vm       *gad.VM
}

// Flow is the form's steps; nil when it is not of steps (`[steps]` before
// its class) or was not read of a class (Builder.Class).
func (s *Schema) Flow() *Flow {
	if s == nil || s.class == nil || !ClassFlag(s.class, MetaSteps) {
		return nil
	}
	fl := &Flow{Autosave: ClassFlag(s.class, MetaAutosave), Vertical: ClassFlag(s.class, MetaVertical), vm: s.vm}
	for _, f := range s.Fields {
		if f.Row != 0 && f.Group != nil {
			if n := len(fl.Steps); n > 0 && fl.Steps[n-1].Group && fl.Steps[n-1].Key == f.Group.Name {
				fl.Steps[n-1].Fields = append(fl.Steps[n-1].Fields, f)
				continue
			}
			step := &Step{Key: f.Group.Name, Group: true, Fields: []*Field{f}}
			step.Label = f.Group.Meta.String("label")
			if step.Label == "" {
				step.Label = strconv.Itoa(len(fl.Steps) + 1)
			}
			fl.Steps = append(fl.Steps, step)
			continue
		}
		step := &Step{Key: f.Name, Fields: []*Field{f}, Label: f.Meta.String("label")}
		if step.Label == "" {
			step.Label = str_utils.HumanizeString(f.Name)
		}
		fl.Steps = append(fl.Steps, step)
	}
	// what the class says of each: its [next=…], the class of a field
	for _, step := range fl.Steps {
		cf, err := s.class.IndexGet(s.vm, gad.Str(step.Key))
		if err != nil {
			continue
		}
		field, ok := cf.(*gad.ClassField)
		if !ok {
			continue
		}
		for _, kv := range field.Meta {
			if kv.K.ToString() == MetaNext {
				step.next = kv.V
			}
		}
		if len(field.Types) == 1 {
			step.class, _ = field.Types[0].(*gad.Class)
		}
	}
	return fl
}

// ClassFlag says the class's metadata says key, and not false.
func ClassFlag(c *gad.Class, key string) bool {
	for _, kv := range c.Meta {
		if kv.K.ToString() == key {
			return !kv.V.IsFalsy()
		}
	}
	return false
}

// Index is the position of the step key; -1 when there is none.
func (fl *Flow) Index(key string) int {
	for i, s := range fl.Steps {
		if s.Key == key {
			return i
		}
	}
	return -1
}

// Step is the step key; nil when there is none.
func (fl *Flow) Step(key string) *Step {
	if i := fl.Index(key); i >= 0 {
		return fl.Steps[i]
	}
	return nil
}

// Next is the step after key, by what was answered (data, the answer so
// far, typed: Decoder.Typed): its [next=…], else the method next of its
// class, else the one after it in the class; "" after the last.
func (fl *Flow) Next(key string, data map[string]any) (string, error) {
	i := fl.Index(key)
	if i < 0 {
		return "", fmt.Errorf("form: no step %q", key)
	}
	step := fl.Steps[i]
	ordered := ""
	if i+1 < len(fl.Steps) {
		ordered = fl.Steps[i+1].Key
	}
	fn := step.next
	if fn == nil && step.class != nil {
		inst, err := fl.instance(step, data)
		if err != nil {
			return "", err
		}
		if inst != nil {
			if m := inst.GetMethod(MetaNext); m != nil {
				fn = m
			}
		}
	}
	if fn == nil {
		return ordered, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), StepTimeout)
	defer cancel()
	ret, err := gad.NewInvoker(fl.vm, fn).WithContext(ctx).Invoke(gad.Args{gad.Array{GadValue(data)}}, nil)
	if err != nil {
		return "", fmt.Errorf("form: the next of %s: %w", key, err)
	}
	var next string
	switch r := ret.(type) {
	case nil:
		return ordered, nil
	case *gad.ClassField:
		next = r.Name
	case gad.Str:
		next = string(r)
	default:
		if ret == gad.Nil {
			return ordered, nil
		}
		return "", fmt.Errorf("form: the next of %s answered %s: a step (Form.b, Form.$2) or nil", key, ret.Type().Name())
	}
	if fl.Index(next) < 0 {
		return "", fmt.Errorf("form: the next of %s answered %q: no such step", key, next)
	}
	return next, nil
}

// instance is the record of the step of a class, as an instance of it: what
// its method next is of.
func (fl *Flow) instance(step *Step, data map[string]any) (*gad.ClassInstance, error) {
	obj, err := gad.NewInvoker(fl.vm, step.class).Invoke(gad.Args{}, nil)
	if err != nil {
		return nil, fmt.Errorf("form: the step %s: %w", step.Key, err)
	}
	inst, ok := obj.(*gad.ClassInstance)
	if !ok {
		return nil, nil
	}
	fields := step.Fields
	if !step.Group && len(fields) == 1 && fields[0].Schema != nil {
		fields = fields[0].Schema.Fields
	}
	if d, ok := TypedGadValue(fields, data[step.Key]).(gad.Dict); ok {
		if err := inst.Init(fl.vm, d); err != nil {
			return nil, stepError(step, err)
		}
	}
	return inst, nil
}

// Path is the way an answer goes through the steps — typed, its values
// (Decoder.Typed) —: from the first, each next its step's Next, until the
// last. A way that comes back to a step is the form's error.
func (fl *Flow) Path(typed map[string]any) ([]string, error) {
	if len(fl.Steps) == 0 {
		return nil, nil
	}
	seen := map[string]bool{}
	var path []string
	for key := fl.Steps[0].Key; key != ""; {
		if seen[key] {
			return path, fmt.Errorf("form: the steps go back to %s", key)
		}
		seen[key] = true
		path = append(path, key)
		next, err := fl.Next(key, typed)
		if err != nil {
			return path, err
		}
		key = next
	}
	return path, nil
}

// DecodeStep reads the fields of the step key alone — what a step posts as
// it is passed —; the form's own validations wait for the last step.
func (d *Decoder) DecodeStep(fl *Flow, key string, data map[string]any) (Record, []Issue) {
	step := fl.Step(key)
	if step == nil {
		return nil, []Issue{{Code: CodeInvalid}}
	}
	return d.DecodeFields(step.Fields, data)
}

// DecodeFlow is Decode of a form of steps: the fields of the steps its way
// goes through (Path) — a step it skips is not asked —, then the form's own
// validations.
func (d *Decoder) DecodeFlow(fl *Flow, data map[string]any) (Record, []Issue) {
	path, err := fl.Path(d.Typed(data).Map())
	if err != nil {
		return nil, []Issue{issueOf(err)}
	}
	var fields []*Field
	for _, key := range path {
		fields = append(fields, fl.Step(key).Fields...)
	}
	return d.decode(fields, d.Schema.Validations, data)
}
