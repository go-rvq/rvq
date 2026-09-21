package schemaform

import (
	"net/url"
	"testing"

	"gopkg.in/yaml.v3"
)

// What the browser posts for a list of records, and what it must come back as:
// the same YAML the value was written in, fields in the order the schema
// declares them.
func TestDecodeListOfRecords(t *testing.T) {
	s, err := Parse("[]{label str; icon str; href}")
	if err != nil {
		t.Fatal(err)
	}

	got, err := yaml.Marshal(s.Decode(url.Values{
		"V[0].label": {"Facebook"},
		"V[0].icon":  {"fa-facebook-f"},
		"V[0].href":  {"#"},
		"V[1].label": {"Instagram"},
		"V[1].icon":  {"fa-instagram"},
		"V[1].href":  {"https://instagram.com/x"},
	}, "V"))
	if err != nil {
		t.Fatal(err)
	}

	want := `- label: Facebook
  icon: fa-facebook-f
  href: '#'
- label: Instagram
  icon: fa-instagram
  href: https://instagram.com/x
`
	if string(got) != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

// A list of plain values: the item IS the value, so the form posts `V[0]`.
func TestDecodeListOfValues(t *testing.T) {
	s, err := Parse("[]str")
	if err != nil {
		t.Fatal(err)
	}

	got, err := yaml.Marshal(s.Decode(url.Values{
		"V[0]": {"Woburn, MA"},
		"V[1]": {"Newton, MA"},
	}, "V"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "- Woburn, MA\n- Newton, MA\n"; string(got) != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

// An empty list is an empty list, not nil — the form simply has no index.
func TestDecodeEmptyList(t *testing.T) {
	s, _ := Parse("[]str")
	got, err := yaml.Marshal(s.Decode(url.Values{}, "V"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "[]\n"; string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Each type comes back as what it IS, so the YAML is typed and not all strings.
func TestDecodeTypes(t *testing.T) {
	s, err := Parse("{num int; u uint; f float; active bool; hidden bool; s str}")
	if err != nil {
		t.Fatal(err)
	}

	got, err := yaml.Marshal(s.Decode(url.Values{
		"V.num":    {"-3"},
		"V.u":      {"7"},
		"V.f":      {"1.5"},
		"V.active": {"true"},
		"V.s":      {"texto"},
	}, "V"))
	if err != nil {
		t.Fatal(err)
	}

	// a switch that is off does not post at all, and comes back false
	want := `num: -3
u: 7
f: 1.5
active: true
hidden: false
s: texto
`
	if string(got) != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

// A field that accepts nil and was left empty is nil; a required one keeps the
// empty value of its type, so the shape of the record never changes.
func TestDecodeEmptyFields(t *testing.T) {
	s, err := Parse("{a? str; b str; c? int}")
	if err != nil {
		t.Fatal(err)
	}

	got, err := yaml.Marshal(s.Decode(url.Values{"V.a": {""}, "V.b": {""}}, "V"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "a: null\nb: \"\"\nc: null\n"; string(got) != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

// A form inside the form, and a list inside a record, come back nested.
func TestDecodeNested(t *testing.T) {
	s, err := Parse("{title str; sub: {x str}; tags []str}")
	if err != nil {
		t.Fatal(err)
	}

	got, err := yaml.Marshal(s.Decode(url.Values{
		"V.title":   {"t"},
		"V.sub.x":   {"dentro"},
		"V.tags[0]": {"a"},
		"V.tags[1]": {"b"},
	}, "V"))
	if err != nil {
		t.Fatal(err)
	}

	want := `title: t
sub:
    x: dentro
tags:
    - a
    - b
`
	if string(got) != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
