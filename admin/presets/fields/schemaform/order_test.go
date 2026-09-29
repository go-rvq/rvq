package schemaform

import (
	"context"
	"net/url"
	"reflect"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

func orderBuilder(t *testing.T) (*Builder, *Schema) {
	t.Helper()
	b := New().Order("order", func(*Context) []OrderOption {
		return []OrderOption{{Value: "CreatedAt", Label: "Criado em"}, {Value: "Title", Label: "Título"}}
	})
	s, err := b.Parse(`{[label="Ordem", hint="a prioridade"] sort? order}`)
	if err != nil {
		t.Fatal(err)
	}
	return b, s
}

// The form of an order: every field offered by its label, the ordering ones
// first; ↑/↓ to move, ASC / DESC / — to order or not. Nothing to write.
func TestOrderComponent(t *testing.T) {
	b, s := orderBuilder(t)
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
	mustContain(t, got, "Ordem", "a prioridade",
		`{"value":"CreatedAt","label":"Criado em"}`, `{"value":"Title","label":"Título"}`,
		"mdi-arrow-up", "mdi-arrow-down", `text='ASC'`, `text='DESC'`,
		`form["Value"].sort = cur.filter(x => x.field !== f)`,
		`@click='() => {`, `@update:model-value='(d) => {`)
}

// What the form posts: the ordering items, in their order; one not ordering, or
// of a direction that is neither, left out.
func TestDecodeOrder(t *testing.T) {
	b, s := orderBuilder(t)
	got, err := b.DecodeForm(&web.EventContext{}, s, url.Values{
		"Value.sort[0].field": {"Title"}, "Value.sort[0].dir": {"asc"},
		"Value.sort[1].field": {"Slug"}, "Value.sort[1].dir": {""},
		"Value.sort[2].field": {"CreatedAt"}, "Value.sort[2].dir": {"DESC"},
	}, "Value")
	if err != nil {
		t.Fatal(err)
	}
	want := []any{map[string]any{"field": "Title", "dir": "ASC"}, map[string]any{"field": "CreatedAt", "dir": "DESC"}}
	rec, _ := got.(Record)
	if len(rec) != 1 || rec[0].Name != "sort" || !reflect.DeepEqual(rec[0].Value, want) {
		t.Errorf("got %#v", got)
	}
}

// The detail: the ordering fields by their labels, ↑ ascending, ↓ descending.
func TestOrderDisplay(t *testing.T) {
	b, s := orderBuilder(t)
	comp := b.DetailComponent(&web.EventContext{}, s, map[string]any{"sort": []any{
		map[string]any{"field": "Title", "dir": "ASC"}, map[string]any{"field": "CreatedAt", "dir": "DESC"},
	}})
	out, err := h.Marshal(comp, context.Background())
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), "Título ↑, Criado em ↓")
}
