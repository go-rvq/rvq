package schemaform

import (
	"net/url"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/web"
)

// countries is a list that is not fixed — it comes from the request, not from
// the schema — which is what EnumItemsFunc is for.
func countries(_ *web.EventContext, path string, _ *Field) []EnumItem {
	if path != "country" && path != "addresses.country" {
		return nil
	}
	return []EnumItem{{Name: "BR", Label: "Brasil"}, {Name: "PT"}}
}

// A field the schema left untyped becomes a select as soon as a list answers
// for it, and the select offers exactly that list.
func TestEnumItemsDrawsASelect(t *testing.T) {
	b := New().EnumItems(countries)
	got := render(t, b, "{country; name str}")

	for _, want := range []string{
		`v-select`,
		`v-model='form["Value"].country'`,
		`Brasil`,
		`PT`, // no label of its own: it shows its name
	} {
		if !strings.Contains(got, want) {
			t.Errorf("o select não traz %s:\n%s", want, got)
		}
	}
	// the field the list says nothing about stays what it was
	if !strings.Contains(got, `v-model='form["Value"].name'`) {
		t.Errorf("o outro field mudou:\n%s", got)
	}
}

// Inside a list, the path is the one the form asks by — a list adds no level.
func TestEnumItemsInsideAList(t *testing.T) {
	b := New().EnumItems(countries)
	got := render(t, b, "{addresses: []{country; street str}}")

	if !strings.Contains(got, `v-model='item.country'`) || !strings.Contains(got, "Brasil") {
		t.Errorf("o select dentro da lista não foi montado:\n%s", got)
	}
}

// A type an application registered keeps its component: registering a type is
// how an application takes a field over.
func TestEnumItemsDoNotStealARegisteredType(t *testing.T) {
	b := New().EnumItems(func(_ *web.EventContext, _ string, _ *Field) []EnumItem {
		return []EnumItem{{Name: "x"}}
	})
	if got := render(t, b, "{n int}"); !strings.Contains(got, "v-text-field") {
		t.Errorf("o int deixou de ser int:\n%s", got)
	}
}

// Saving: a value the list did not offer is refused, and the value still comes
// back whole so the form can be shown again.
func TestDecodeFormChecksTheList(t *testing.T) {
	b := New().EnumItems(countries)
	s, err := Parse("{country; name str}")
	if err != nil {
		t.Fatal(err)
	}

	value, err := b.DecodeForm(&web.EventContext{}, s, url.Values{
		"V.country": {"BR"},
		"V.name":    {"Ada"},
	}, "V")
	if err != nil {
		t.Fatalf("um valor da lista foi recusado: %v", err)
	}
	if rec, ok := value.(Record); !ok || rec[0].Value != "BR" {
		t.Errorf("value = %#v", value)
	}

	value, err = b.DecodeForm(&web.EventContext{}, s, url.Values{
		"V.country": {"ZZ"},
		"V.name":    {"Ada"},
	}, "V")
	if err == nil {
		t.Fatal("um valor fora da lista passou")
	}
	for _, want := range []string{"country", `"ZZ"`, "BR", "PT"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("o erro não diz %s: %v", want, err)
		}
	}
	if rec, ok := value.(Record); !ok || len(rec) != 2 {
		t.Errorf("o valor não voltou inteiro: %#v", value)
	}
}

// An enum the schema DID declare is checked the same way, with no function at
// all.
func TestDecodeFormChecksADeclaredEnum(t *testing.T) {
	s, err := Parse("enum Perm { Read, Write }\ninterface { perm Perm; opt? Perm }")
	if err != nil {
		t.Fatal(err)
	}
	b := New()

	if _, err := b.DecodeForm(&web.EventContext{}, s, url.Values{"V.perm": {"Read"}}, "V"); err != nil {
		t.Fatalf("um valor do enum foi recusado: %v", err)
	}

	_, err = b.DecodeForm(&web.EventContext{}, s, url.Values{"V.perm": {"Delete"}}, "V")
	if err == nil || !strings.Contains(err.Error(), "Delete") {
		t.Fatalf("err = %v", err)
	}

	// required and empty is reported as empty; optional and empty is fine
	_, err = b.DecodeForm(&web.EventContext{}, s, url.Values{"V.perm": {""}}, "V")
	if err == nil || !strings.Contains(err.Error(), "escolha um valor") {
		t.Fatalf("err = %v", err)
	}
	if strings.Contains(err.Error(), "opt") {
		t.Errorf("o field opcional vazio foi cobrado: %v", err)
	}
}

// Schema.Decode reads and does not check — checking is the builder's, which is
// the one that knows what the select offered.
func TestSchemaDecodeDoesNotCheck(t *testing.T) {
	s, _ := Parse("enum Perm { Read }\ninterface { perm Perm }")
	if rec, ok := s.Decode(url.Values{"V.perm": {"Delete"}}, "V").(Record); !ok || rec[0].Value != "Delete" {
		t.Errorf("value = %#v", s.Decode(url.Values{"V.perm": {"Delete"}}, "V"))
	}
}
