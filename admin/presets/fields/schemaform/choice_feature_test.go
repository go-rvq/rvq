package schemaform

import (
	"reflect"
	"testing"

	"github.com/gad-lang/gad"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/features"
)

var choiceFeature = features.Register("schemaform.test.choice", "A test feature", "a class of a choice needs it").CompiledIf(true)

// A class that names a feature (`[feature="…"]`) is offered in a choice —
// with its fields, or by its name — only while the feature is on; a feature
// the binary does not know is off.
func TestChoiceFeature(t *testing.T) {
	const src = `
class Default { title? str }
[feature="schemaform.test.choice"]
class Issues { perPage? int }
[feature="schemaform.test.nobody"]
class Ghost { x? str }
class Page { Layout Default | Issues | Ghost }
return Page`
	names := func(asName bool) []string {
		t.Helper()
		b := New().ChoiceAsName(asName)
		vm, ret, err := b.runProgram(src)
		if err != nil {
			t.Fatal(err)
		}
		s, err := b.Class(vm, ret.(*gad.Class))
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		f := s.Fields[0]
		if asName {
			items, err := (&Context{Field: f, Builder: b, Event: &web.EventContext{}}).EnumItems()
			if err != nil {
				t.Fatal(err)
			}
			for _, it := range items {
				out = append(out, it.Name)
			}
			return out
		}
		for _, c := range f.Schema.Fields {
			out = append(out, c.Name)
		}
		return out
	}
	defer features.Set(choiceFeature.Name, false)
	features.Set(choiceFeature.Name, false)
	for _, asName := range []bool{false, true} {
		if got := names(asName); !reflect.DeepEqual(got, []string{"Default"}) {
			t.Errorf("off (by name %v): %v", asName, got)
		}
	}
	features.Set(choiceFeature.Name, true)
	for _, asName := range []bool{false, true} {
		if got := names(asName); !reflect.DeepEqual(got, []string{"Default", "Issues"}) {
			t.Errorf("on (by name %v): %v", asName, got)
		}
	}
}

// A field offered only with a feature on: its own (`[feature="…"] x? T`), or
// the one of the class it is typed by; off, it is not in the schema.
func TestFieldFeature(t *testing.T) {
	const src = `
[feature="schemaform.test.choice"]
class IssuesOptions { perPage? int }
class Default {
	title? str
	issues? IssuesOptions
	[feature="schemaform.test.choice"]
	beforeBody? bool
	[feature="schemaform.test.nobody"]
	ghost? str
}
return Default`
	fields := func() []string {
		t.Helper()
		b := New()
		vm, ret, err := b.runProgram(src)
		if err != nil {
			t.Fatal(err)
		}
		s, err := b.Class(vm, ret.(*gad.Class))
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		for _, f := range s.Fields {
			out = append(out, f.Name)
		}
		return out
	}
	defer features.Set(choiceFeature.Name, false)
	features.Set(choiceFeature.Name, false)
	if got := fields(); !reflect.DeepEqual(got, []string{"title"}) {
		t.Errorf("off: %v", got)
	}
	features.Set(choiceFeature.Name, true)
	if got := fields(); !reflect.DeepEqual(got, []string{"title", "issues", "beforeBody"}) {
		t.Errorf("on: %v", got)
	}
}
