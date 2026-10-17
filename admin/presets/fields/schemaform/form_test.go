package schemaform

import (
	"net/url"
	"sort"
	"strings"
	"testing"
)

// issuePaths are the paths of what is wrong, sorted.
func issuePaths(issues []Issue) []string {
	var out []string
	for _, e := range issues {
		out = append(out, e.Path)
	}
	sort.Strings(out)
	return out
}

// The one logic of a form: the Decoder reads what the admin posts (Nest) the
// way it checks it on save — each check by the path of the field.
func TestDecoderChecks(t *testing.T) {
	s, err := readClass(checkForm)
	if err != nil {
		t.Fatal(err)
	}
	decode := func(values url.Values) []Issue {
		_, issues := (&Decoder{Schema: s}).Decode(Nest(values, "V"))
		return issues
	}
	ok := url.Values{"V.name": {"Ada"}, "V.rooms": {"3"}}
	if got := decode(ok); len(got) != 0 {
		t.Fatalf("a form right: %v", got)
	}
	for _, c := range []struct {
		name string
		set  url.Values
		want string
		code string
	}{
		{"required left empty", url.Values{"V.name": {""}}, "name", CodeRequired},
		{"not a number", url.Values{"V.rooms": {"x"}}, "rooms", CodeInvalid},
		{"over the max", url.Values{"V.rooms": {"11"}}, "rooms", CodeMax},
		{"too short", url.Values{"V.code": {"ab"}}, "code", CodeTooShort},
		{"too long", url.Values{"V.code": {"abcdef"}}, "code", CodeTooLong},
		{"off the pattern", url.Values{"V.uf": {"mg"}}, "uf", CodePattern},
		{"off the mask", url.Values{"V.zip": {"12345678"}}, "zip", ""},
		{"off the regular expression", url.Values{"V.three": {"12"}}, "three", CodePattern},
		{"refused by the function", url.Values{"V.word": {"bad"}}, "word", CodeValidation},
		{"a range out of order", url.Values{"V.stay.from": {"7"}, "V.stay.to": {"2"}}, "stay", CodeRange},
		{"a range's bound off its limits", url.Values{"V.stay.from": {"0"}, "V.stay.to": {"2"}}, "stay.from", CodeMin},
		{"a range's bound left empty", url.Values{"V.stay.from": {"2"}}, "stay.to", CodeRequired},
		{"a date range out of order", url.Values{"V.when.from": {"2026-05-02"}, "V.when.to": {"2026-05-01"}}, "when", CodeRange},
		{"a card begun: its required asked", url.Values{"V.billing.street": {"Main"}}, "billing.city", CodeRequired},
		{"a list's field, by its index", url.Values{"V.links[0].href": {"https://a.b"}, "V.links[1].href": {"nope"}}, "links[1].href", CodePattern},
	} {
		values := url.Values{}
		for k, v := range ok {
			values[k] = v
		}
		for k, v := range c.set {
			values[k] = v
		}
		got := decode(values)
		if len(got) != 1 || got[0].Path != c.want || c.code != "" && got[0].Code != c.code {
			t.Errorf("%s: %v, want %s %s", c.name, got, c.want, c.code)
		}
	}
	all := url.Values{"V.name": {"Ada"}, "V.rooms": {"10"}, "V.code": {"abc"}, "V.uf": {"MG"}, "V.zip": {"12345-678"},
		"V.three": {"123"}, "V.word": {"good"}, "V.stay.from": {"2"}, "V.stay.to": {"2"}, "V.when.from": {"2026-05-01"},
		"V.billing.street": {"Main"}, "V.billing.city": {"Town"}, "V.links[0].href": {"http://x.y"}}
	rec, issues := (&Decoder{Schema: s}).Decode(Nest(all, "V"))
	if len(issues) != 0 {
		t.Fatalf("all within: %v", issues)
	}
	if v := rec[2].Value; v != int64(10) {
		t.Errorf("rooms typed: %#v", v)
	}
	// the function's words, as it said them
	_, issues = (&Decoder{Schema: s}).Decode(map[string]any{"name": "Ada", "rooms": "3", "word": "bad"})
	if len(issues) != 1 || issues[0].Message != "not that one" {
		t.Errorf("the function's words: %v", issues)
	}
}

// A group's fields are posted where its class's are (gad [autoGroupKeys],
// set by the reader): flat; the record `::: Form` built of them either way.
func TestDecoderGroups(t *testing.T) {
	s, err := readClass(`export class default {
		name str
		{
			phone str
			city? str
		}
	
}`)
	if err != nil {
		t.Fatal(err)
	}
	_, issues := (&Decoder{Schema: s}).Decode(map[string]any{"name": "Ada"})
	if p := issuePaths(issues); len(p) != 1 || p[0] != "phone" {
		t.Errorf("a group's required field: %v", p)
	}
}

// Of the application: a type by its reader; the values of a field the
// request knows (Items); a list by how many; a field's validations off
// (Typed).
func TestDecoderExtension(t *testing.T) {
	s, err := readClass(`export class default {
		[validation=func(v) => "never"] tel str
		pick? str
		[min=2] tags? []str
	
}`)
	if err != nil {
		t.Fatal(err)
	}
	s.Fields[0].Type = "phone"
	d := &Decoder{
		Schema: s,
		Types: map[string]TypeReader{"phone": func(f *Field, _ string, raw any) (any, string) {
			digits := strings.Map(func(r rune) rune {
				if r >= '0' && r <= '9' {
					return r
				}
				return -1
			}, Text(raw))
			if len(digits) < 10 {
				return nil, CodeInvalid
			}
			return digits, ""
		}},
		Items: func(path string, f *Field) ([]EnumItem, error) {
			if path == "pick" {
				return []EnumItem{{Name: "a"}, {Name: "b"}}, nil
			}
			return nil, nil
		},
	}
	_, issues := d.Decode(map[string]any{"tel": "12", "pick": "c", "tags": []any{"x"}})
	want := []string{"pick", "tags", "tel"}
	if got := issuePaths(issues); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("issues: %v", issues)
	}
	rec := d.Typed(map[string]any{"tel": "(781) 555-0100", "pick": "a", "tags": []any{"x", "", "y"}})
	if rec[0].Value != "7815550100" || rec[1].Value != "a" || len(rec[2].Value.([]any)) != 2 {
		t.Errorf("typed: %#v", rec)
	}
	if _, issues := d.Decode(map[string]any{"tel": "(781) 555-0100"}); len(issues) != 1 || issues[0].Code != CodeValidation {
		t.Errorf("the field's validation after its type: %v", issues)
	}
}

// Nest: the keys the admin posts as a record nested.
func TestNest(t *testing.T) {
	got := Nest(url.Values{"V.a": {"1"}, "V.b.c": {"2"}, "V.l[1].x": {"3"}, "W.z": {"no"}, "V.m": {"4", "5"}}, "V")
	want := map[string]any{"a": "1", "b": map[string]any{"c": "2"}, "l": map[string]any{"1": map[string]any{"x": "3"}}, "m": []any{"4", "5"}}
	if !equalAny(got, want) {
		t.Errorf("Nest: %#v", got)
	}
}

func equalAny(a, b any) bool {
	switch x := a.(type) {
	case map[string]any:
		y, ok := b.(map[string]any)
		if !ok || len(x) != len(y) {
			return false
		}
		for k := range x {
			if !equalAny(x[k], y[k]) {
				return false
			}
		}
		return true
	case []any:
		y, ok := b.([]any)
		if !ok || len(x) != len(y) {
			return false
		}
		for i := range x {
			if !equalAny(x[i], y[i]) {
				return false
			}
		}
		return true
	}
	return a == b
}
