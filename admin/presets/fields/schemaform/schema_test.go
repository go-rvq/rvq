package schemaform

import (
	"strings"
	"testing"
)

// The schema PAC carries: a list of records, one field of it untyped.
func TestParseSliceOfRecords(t *testing.T) {
	s, err := Parse("[]{label str; icon str; href}")
	if err != nil {
		t.Fatal(err)
	}

	if !s.Slice {
		t.Error("[]{…} descreve uma lista")
	}
	if got, want := len(s.Fields), 3; got != want {
		t.Fatalf("fields = %d, want %d", got, want)
	}

	for i, want := range []struct {
		name, typ string
	}{
		{"label", "str"},
		{"icon", "str"},
		{"href", DefaultType}, // untyped is plain text
	} {
		if got := s.Fields[i]; got.Name != want.name || got.Type != want.typ {
			t.Errorf("field %d = %q %q, want %q %q", i, got.Name, got.Type, want.name, want.typ)
		}
	}
}

// A form, not a list.
func TestParseRecord(t *testing.T) {
	s, err := Parse("{title str; count int}")
	if err != nil {
		t.Fatal(err)
	}
	if s.Slice {
		t.Error("{…} descreve um registro, não uma lista")
	}
	if got, want := s.Fields[1].Type, "int"; got != want {
		t.Errorf("type = %q, want %q", got, want)
	}
}

// A field that is itself an interface is a form again — however it is spelled.
func TestParseNestedInterface(t *testing.T) {
	cases := map[string]struct {
		src   string
		slice bool
	}{
		"long form":        {"{a str; sub interface {x str}}", false},
		"long form slice":  {"{a str; sub interface[] {x str}}", true},
		"short form":       {"{a str; sub: {x str}}", false},
		"short form slice": {"{a str; sub: []{x str}}", true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			s, err := Parse(c.src)
			if err != nil {
				t.Fatal(err)
			}
			sub := s.Fields[1]
			if sub.Name != "sub" {
				t.Fatalf("name = %q", sub.Name)
			}
			if sub.Type != FormType {
				t.Errorf("type = %q, want %q", sub.Type, FormType)
			}
			if sub.Schema == nil {
				t.Fatal("o field não trouxe o schema dele")
			}
			if sub.Schema.Slice != c.slice {
				t.Errorf("slice = %v, want %v", sub.Schema.Slice, c.slice)
			}
			if got, want := sub.Schema.Fields[0].Name, "x"; got != want {
				t.Errorf("nested field = %q, want %q", got, want)
			}
		})
	}
}

// The `?` after a name marks the field nullable.
func TestParseNullable(t *testing.T) {
	s, err := Parse("{a str; b? str}")
	if err != nil {
		t.Fatal(err)
	}
	if s.Fields[0].Nullable {
		t.Error("a não é nullable")
	}
	if !s.Fields[1].Nullable {
		t.Error("b? é nullable")
	}
}

// Only fields are a form: an accessor describes behaviour, and there is nothing
// to edit in it.
//
// A METHOD is not covered here: gad reads the return type of `e() str` as one
// more body item, so it would arrive as a field named "str". A form schema has
// no reason to declare methods, and this says what the walk does with what it
// is given.
func TestParseKeepsOnlyFields(t *testing.T) {
	s, err := Parse("{a str; get b; set c; prop d str}")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(s.Fields), 1; got != want {
		var names []string
		for _, f := range s.Fields {
			names = append(names, f.Name)
		}
		t.Fatalf("fields = %d (%v), want %d", got, names, want)
	}
	if s.Fields[0].Name != "a" {
		t.Errorf("field = %q, want a", s.Fields[0].Name)
	}
}

func TestParseRejectsWhatIsNotASchema(t *testing.T) {
	for name, src := range map[string]string{
		"empty":       "",
		"not a form":  "42",
		"broken":      "{a str",
		"two of them": "{a str} interface {b str}",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := Parse(src); err == nil {
				t.Error("devia recusar")
			} else if !strings.HasPrefix(err.Error(), "schemaform:") {
				t.Errorf("erro sem contexto: %v", err)
			}
		})
	}
}

// An enum declared beside the interface is a closed set of values, and a field
// typed with it carries it.
func TestParseEnum(t *testing.T) {
	s, err := Parse("enum Perm { Read, Write }\ninterface {perm Perm; other str}")
	if err != nil {
		t.Fatal(err)
	}

	e, ok := s.Enums["Perm"]
	if !ok {
		t.Fatalf("o enum não foi declarado: %v", s.Enums)
	}
	if got, want := strings.Join(e.Names, ","), "Read,Write"; got != want {
		t.Errorf("valores = %q, want %q", got, want)
	}

	if s.Fields[0].Enum != e {
		t.Error("o field não recebeu o enum do type dele")
	}
	if s.Fields[1].Enum != nil {
		t.Error("um field de outro type não tem enum")
	}
}

// An enum declared once serves the forms inside the form too.
func TestParseEnumReachesNestedForms(t *testing.T) {
	s, err := Parse("enum Perm { Read, Write }\ninterface {sub interface {perm Perm}}")
	if err != nil {
		t.Fatal(err)
	}
	if s.Fields[0].Schema.Fields[0].Enum == nil {
		t.Error("o form aninhado não enxergou o enum")
	}
}

// Empty is only allowed where the schema said it is: the `?` after the name.
func TestFieldRequired(t *testing.T) {
	s, err := Parse("{btnColor? color; label str}")
	if err != nil {
		t.Fatal(err)
	}
	if s.Fields[0].Required() {
		t.Error("btnColor? é opcional")
	}
	if !s.Fields[1].Required() {
		t.Error("um field sem ? é obrigatório")
	}
}
