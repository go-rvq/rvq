package schemaform

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/gad-lang/gad"
	"github.com/go-rvq/rvq/web"
)

// With ChoiceAsName, a union of classes is the choice of one BY ITS NAME: an
// enum of the classes' names, each by its label (else its name humanized) and
// hint; the value posted is the name.
func TestChoiceAsName(t *testing.T) {
	b := New().ChoiceAsName(true)
	vm, ret, err := b.runProgram(`
[label="Serviços", hint="o que diz um serviço"]
class ServiceOptions { icon? str }
class PortfolioOptions { gradient? str }
class PostTypeConfig {
    [label="Layout das postagens"]
    PostLayout? ServiceOptions | PortfolioOptions
}
return PostTypeConfig`)
	if err != nil {
		t.Fatal(err)
	}
	s, err := b.Class(vm, ret.(*gad.Class))
	if err != nil {
		t.Fatal(err)
	}
	f := s.Fields[0]
	if f.Schema != nil || f.Enum == nil {
		t.Fatalf("not an enum: %+v", f)
	}
	items, err := (&Context{Field: f, Builder: b, Event: &web.EventContext{}}).EnumItems()
	if err != nil {
		t.Fatal(err)
	}
	want := []EnumItem{
		{Name: "ServiceOptions", Label: "Serviços", Hint: "o que diz um serviço"},
		{Name: "PortfolioOptions", Label: "Portfolio Options"},
	}
	if !reflect.DeepEqual(items, want) {
		t.Errorf("items %+v", items)
	}

	value, err := b.DecodeForm(&web.EventContext{}, s, url.Values{"k.PostLayout": {"PortfolioOptions"}}, "k")
	if err != nil {
		t.Fatal(err)
	}
	if got := recordValue(value, "PostLayout"); got != "PortfolioOptions" {
		t.Errorf("decoded %#v", value)
	}
	if _, err := b.DecodeForm(&web.EventContext{}, s, url.Values{"k.PostLayout": {"Other"}}, "k"); err == nil {
		t.Error("a class not in the union taken")
	}

	// without it, the choice of a class with its fields
	s, err = New().Class(vm, ret.(*gad.Class))
	if err != nil {
		t.Fatal(err)
	}
	if f := s.Fields[0]; f.Schema == nil || !f.Schema.Choice {
		t.Errorf("not a choice: %+v", f)
	}
}
