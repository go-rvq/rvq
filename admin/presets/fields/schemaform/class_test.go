package schemaform

import (
	"context"
	"net/url"
	"reflect"
	"testing"

	"github.com/gad-lang/gad"
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

// optionsSrc is a layout's options, as a layout declares them: nested classes,
// a choice of classes, and a type only the application knows (Ref — a class
// here, what TypeOf recognizes before it is read as a record).
const optionsSrc = `
class Ref { id str }
[label="Lista"]
class listLayout { comp? enum { index_list, compact } }
[label="Grade"]
class gridLayout { [label="Colunas", hint="de 1 a 6"] columns int = 3 }
class Options {
    [label="Tipo de post"] postType? Ref
    posts class {
        [label="Desenho"] layout? listLayout|gridLayout
    }
    title str = "x"
}
return Options
`

// classSchema is the form of optionsSrc: Ref is drawn as "ref".
func classSchema(t *testing.T) (*Builder, *Schema) {
	t.Helper()
	b := New().Type("ref", TextComponentFunc).
		TypeOf(func(o gad.Object) string {
			if c, ok := o.(*gad.Class); ok && c.Name() == "Ref" {
				return "ref"
			}
			return ""
		})
	vm, ret, err := b.runProgram(optionsSrc)
	if err != nil {
		t.Fatal(err)
	}
	s, err := b.Class(vm, ret.(*gad.Class))
	if err != nil {
		t.Fatal(err)
	}
	return b, s
}

// A class is a form: a class inside it — `posts class { … }` — a record of
// its own, a union of classes a choice of one, each class by its name, a type
// the application names drawn as it says; defaults and metadata kept.
func TestClassSchema(t *testing.T) {
	_, s := classSchema(t)

	if f := s.FieldAt("postType"); f == nil || f.Type != "ref" || !f.Nullable {
		t.Errorf("postType: %+v", f)
	}
	if f := s.FieldAt("posts"); f == nil || f.Type != FormType || f.Schema == nil || f.Schema.Choice {
		t.Errorf("posts is a record: %+v", f)
	}
	layout := s.FieldAt("posts.layout")
	if layout == nil || layout.Schema == nil || !layout.Schema.Choice {
		t.Fatalf("posts.layout is a choice: %+v", layout)
	}
	var names []string
	for _, f := range layout.Schema.Fields {
		names = append(names, f.Name)
	}
	if !reflect.DeepEqual(names, []string{"listLayout", "gridLayout"}) {
		t.Errorf("the choices: %v", names)
	}
	if f := s.FieldAt("posts.layout.gridLayout"); f == nil || f.Meta["label"] != "Grade" {
		t.Errorf("a class's metadata is its choice's: %+v", f)
	}
	if f := s.FieldAt("posts.layout.gridLayout.columns"); f == nil || f.Default != int64(3) && f.Default != 3 {
		t.Errorf("columns: %+v", f)
	}
	if f := s.FieldAt("posts.layout.listLayout.comp"); f == nil || f.Enum == nil {
		t.Errorf("comp is an enum: %+v", f)
	}
	if f := s.FieldAt("title"); f == nil || f.Default != "x" {
		t.Errorf("title: %+v", f)
	}
}

// A choice holds the class chosen and nothing else: filled, only what it
// holds is.
func TestChoiceFilled(t *testing.T) {
	_, s := classSchema(t)
	got := s.Filled(map[string]any{"posts": map[string]any{"layout": map[string]any{"gridLayout": nil}}})
	layout := got.(map[string]any)["posts"].(map[string]any)["layout"].(map[string]any)
	if len(layout) != 1 || layout["gridLayout"] == nil {
		t.Errorf("the class chosen, alone and filled: %v", layout)
	}
	empty := s.Filled(nil).(map[string]any)["posts"].(map[string]any)["layout"].(map[string]any)
	if len(empty) != 0 {
		t.Errorf("nothing chosen, nothing held: %v", empty)
	}

	// a field the value lacks starts as its class's default
	if grid := layout["gridLayout"].(map[string]any); grid["columns"] == nil {
		t.Errorf("the default of columns: %v", grid)
	}
	if got := s.Filled(map[string]any{"title": "y"}).(map[string]any)["title"]; got != "y" {
		t.Errorf("a value given is kept: %v", got)
	}
}

// The form of a choice: a select of the classes, by their labels, and the
// form of the one chosen only while it is; choosing another starts it empty.
// Labels and hints the application does not give come from the metadata.
func TestChoiceComponent(t *testing.T) {
	b, s := classSchema(t)
	comp := b.ComponentFunc(s)(&presets.FieldContext{
		ToComponentOptions: &presets.ToComponentOptions{},
		Name:               "Value",
		FormKey:            "Value",
		Mode:               presets.FieldModeStack{presets.EDIT},
	}, &web.EventContext{})
	out, err := h.Marshal(comp, context.Background())
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	mustContain(t, got,
		`:model-value='Object.keys(form["Value"].posts.layout || {})[0]'`,
		`"title":"Lista","value":"listLayout"`, `"title":"Grade","value":"gridLayout"`,
		`v-if='Object.keys(form["Value"].posts.layout || {})[0] === "gridLayout"'`,
		`form["Value"].posts.layout["gridLayout"].columns`,
		`"gridLayout": {"columns": 0}`,
		"Desenho", "Colunas", "de 1 a 6", "Tipo de post")
}

// The detail of a choice: the class chosen, by its label, and what it holds.
func TestChoiceDisplay(t *testing.T) {
	b, s := classSchema(t)
	comp := b.DetailComponent(&web.EventContext{}, s,
		map[string]any{"posts": map[string]any{"layout": map[string]any{"gridLayout": map[string]any{"columns": 4}}}})
	out, err := h.Marshal(comp, context.Background())
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), "Grade", "Colunas", "4")
	mustNotContain(t, string(out), "Lista")
}

// What the form posts for a choice is the class chosen, alone: the one it
// posted something under.
func TestChoiceDecode(t *testing.T) {
	b, s := classSchema(t)
	values := url.Values{
		"Value.posts.layout.gridLayout.columns": {"4"},
		"Value.title":                           {"y"},
	}
	v, err := b.DecodeForm(&web.EventContext{}, s, values, "Value")
	if err != nil {
		t.Fatal(err)
	}
	out, err := b.EncodeValue(v)
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, out, "gridLayout:", "columns: 4", "title: \"y\"")
	mustNotContain(t, out, "listLayout")

	v, _ = b.DecodeForm(&web.EventContext{}, s, url.Values{"Value.title": {"y"}}, "Value")
	out, _ = b.EncodeValue(v)
	mustNotContain(t, out, "gridLayout", "listLayout")
}
