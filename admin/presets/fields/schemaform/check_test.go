package schemaform

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

// checkForm is a class of each kind of check the admin does on save.
const checkForm = `
[ordered]
class Address {
	street str
	city str
}
class Link {
	[pattern="https?://.+"] href str
}
[ordered]
class Form {
	name str
	nick? str
	[min=1, max=10] rooms int
	[minlength=3, maxlength=5] code? str
	[pattern="[A-Z]{2}"] uf? str
	[mask="#####-###"] zip? str
	[validation=/\d{3}/] three? str
	[validation=func(v) => v == "bad" ? "not that one" : nil] word? str
	[min=1, max=10] stay? Range[int]
	[open="to"] when? Range[date]
	billing? Address
	links? []Link
}
return Form`

// What is wrong, field by field, by the keys the form posts.
func checkErrors(t *testing.T, s *Schema, values url.Values) map[string][]string {
	t.Helper()
	_, err := New().DecodeForm(&web.EventContext{}, s, values, "V")
	if err == nil {
		return nil
	}
	var fe *FieldErrors
	if !errors.As(err, &fe) {
		t.Fatalf("not by field: %v", err)
	}
	return fe.Fields
}

// The admin checks a form posted on save, field by field: a required one
// left empty, its type, its limits, its pattern and mask, its
// `[validation=…]` (a regular expression, a function), a range (its order,
// its bounds' limits); an optional card left blank asks nothing; the fields
// of a list by their index.
func TestDecodeFormChecks(t *testing.T) {
	s, err := readClass(checkForm)
	if err != nil {
		t.Fatal(err)
	}
	ok := url.Values{"V.name": {"Ada"}, "V.rooms": {"3"}}
	if got := checkErrors(t, s, ok); len(got) != 0 {
		t.Fatalf("a form right: %v", got)
	}
	for _, c := range []struct {
		name string
		set  url.Values
		want []string
	}{
		{"required left empty", url.Values{"V.name": {""}}, []string{"V.name"}},
		{"not a number", url.Values{"V.rooms": {"x"}}, []string{"V.rooms"}},
		{"over the max", url.Values{"V.rooms": {"11"}}, []string{"V.rooms"}},
		{"too short", url.Values{"V.code": {"ab"}}, []string{"V.code"}},
		{"too long", url.Values{"V.code": {"abcdef"}}, []string{"V.code"}},
		{"off the pattern", url.Values{"V.uf": {"mg"}}, []string{"V.uf"}},
		{"off the mask", url.Values{"V.zip": {"12345678"}}, []string{"V.zip"}},
		{"off the regular expression", url.Values{"V.three": {"12"}}, []string{"V.three"}},
		{"refused by the function", url.Values{"V.word": {"bad"}}, []string{"V.word"}},
		{"a range out of order", url.Values{"V.stay.from": {"7"}, "V.stay.to": {"2"}}, []string{"V.stay"}},
		{"a range's bound off its limits", url.Values{"V.stay.from": {"0"}, "V.stay.to": {"2"}}, []string{"V.stay.from"}},
		{"a range's bound left empty", url.Values{"V.stay.from": {"2"}}, []string{"V.stay.to"}},
		{"a date range out of order", url.Values{"V.when.from": {"2026-05-02"}, "V.when.to": {"2026-05-01"}}, []string{"V.when"}},
		{"a card begun: its required asked", url.Values{"V.billing.street": {"Main"}}, []string{"V.billing.city"}},
		{"a list's field, by its index", url.Values{"V.links[0].href": {"https://a.b"}, "V.links[1].href": {"nope"}}, []string{"V.links[1].href"}},
	} {
		values := url.Values{}
		for k, v := range ok {
			values[k] = v
		}
		for k, v := range c.set {
			values[k] = v
		}
		got := checkErrors(t, s, values)
		if len(got) != len(c.want) {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
			continue
		}
		for _, k := range c.want {
			if len(got[k]) == 0 {
				t.Errorf("%s: nothing under %s (%v)", c.name, k, got)
			}
		}
	}
	// within, all of them: nothing said
	all := url.Values{"V.name": {"Ada"}, "V.rooms": {"10"}, "V.code": {"abc"}, "V.uf": {"MG"}, "V.zip": {"12345-678"},
		"V.three": {"123"}, "V.word": {"good"}, "V.stay.from": {"2"}, "V.stay.to": {"2"}, "V.when.from": {"2026-05-01"},
		"V.billing.street": {"Main"}, "V.billing.city": {"Town"}, "V.links[0].href": {"http://x.y"}}
	if got := checkErrors(t, s, all); len(got) != 0 {
		t.Errorf("all within: %v", got)
	}
}

// A range is read back as its from and its to (it was read as nothing).
func TestDecodeFormRange(t *testing.T) {
	s, err := readClass("class Form { r? Range[int]; d? Range[date] }; return Form")
	if err != nil {
		t.Fatal(err)
	}
	v, err := New().DecodeForm(&web.EventContext{}, s, url.Values{"V.r.from": {"1"}, "V.r.to": {"5"}}, "V")
	if err != nil {
		t.Fatal(err)
	}
	rec := v.(Record)
	r, ok := rec[0].Value.(Record)
	if !ok || len(r) != 2 || r[0].Value != int64(1) || r[1].Value != int64(5) {
		t.Errorf("the range: %#v", rec[0].Value)
	}
	if rec[1].Value != nil {
		t.Errorf("an optional range left blank: %#v", rec[1].Value)
	}
}

// What the server found wrong on save is said under each input: the errors
// of the field's keys (presets' SubErrors) in the scope, each input binding
// its own key — of a list's item, by its index —, gone once it changes.
func TestErrorsUnderEachInput(t *testing.T) {
	schema, err := New().Parse(`interface Link { href str }
	interface Form {
		name str
		links []Link
		r? Range[int]
	}`)
	if err != nil {
		t.Fatal(err)
	}
	comp := New().ComponentFunc(schema)(&presets.FieldContext{
		ToComponentOptions: &presets.ToComponentOptions{},
		Name:               "Value",
		FormKey:            "Value",
		Mode:               presets.FieldModeStack{presets.EDIT},
		SubErrors:          map[string][]string{"Value.name": {"fill this field"}, "Value.links[1].href": {"bad"}},
	}, &web.EventContext{})
	out, err := h.Marshal(comp, context.Background())
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	for _, want := range []string{
		`vars.sfErr_Value = {\"Value.links[1].href\":[\"bad\"],\"Value.name\":[\"fill this field\"]}`,
		`:error-messages='(vars.sfErr_Value || {})["Value" + ".name"]'`,
		`(vars.sfErr_Value || {})["Value" + ".links" + "[" + itemIndex + "]" + ".href"]`,
		`delete vars.sfErr_Value["Value" + ".name"]`,
		`data-range-error`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("no %s in:\n%s", want, got)
		}
	}
}
