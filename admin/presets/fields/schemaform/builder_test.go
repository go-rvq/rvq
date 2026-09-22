package schemaform

import (
	"context"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// render draws the schema for a presets field named Value and returns the HTML.
func render(t *testing.T, b *Builder, src string) string {
	t.Helper()

	schema, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}

	comp := b.ComponentFunc(schema)(&presets.FieldContext{
		ToComponentOptions: &presets.ToComponentOptions{},
		Name:               "Value",
		FormKey:            "Value",
		Mode:               presets.FieldModeStack{presets.EDIT},
	}, &web.EventContext{})

	out, err := h.Marshal(comp, context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

// A record: each field binds under the form key of the presets field.
func TestComponentFuncBindsEachFieldUnderTheFormKey(t *testing.T) {
	got := render(t, New(), "{label str; href}")

	for _, want := range []string{
		`v-model='form["Value"].label'`,
		`v-model='form["Value"].href'`, // untyped: a text field all the same
	} {
		if !strings.Contains(got, want) {
			t.Errorf("o form não liga %s:\n%s", want, got)
		}
	}
}

// A list: the sorter iterates over the value, and each item binds its own
// fields under the slot's item.
func TestComponentFuncDrawsAListWithTheSorter(t *testing.T) {
	got := render(t, New(), "[]{label str; icon str; href}")

	for _, want := range []string{
		`vx-array-sorter`,
		`v-model='form["Value"]'`,
		// the list as it is edited is the DEFAULT slot: the sorter draws its
		// own rows only while sorting
		`v-slot:default`,
		`v-for='(item, itemIndex) in form["Value"]'`,
		`item-title='label'`, // what a row shows while being sorted
		`v-model='item.label'`,
		`v-model='item.icon'`,
		`v-model='item.href'`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("a lista não traz %s:\n%s", want, got)
		}
	}
}

// A field that is itself an interface draws the form again, one level down —
// and a list of them draws a sorter inside the form.
func TestComponentFuncRecursesIntoANestedForm(t *testing.T) {
	got := render(t, New(), "{title str; sub interface {x str}}")
	for _, want := range []string{
		`v-model='form["Value"].title'`,
		`v-model='form["Value"].sub.x'`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("o form aninhado não traz %s:\n%s", want, got)
		}
	}

	got = render(t, New(), "{title str; subs interface[] {x str}}")
	for _, want := range []string{
		`v-model='form["Value"].subs'`, // the nested sorter, over the field
		`v-model='item.x'`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("a lista aninhada não traz %s:\n%s", want, got)
		}
	}
}

// The type is what picks the component, and an application registers its own.
func TestTypeRegistry(t *testing.T) {
	b := New().Type("color", func(c *Context) h.HTMLComponent {
		return v.VTextField().Attr("data-color", "1").Attr("v-model", c.Value)
	})

	got := render(t, b, "{a str; b color}")
	if !strings.Contains(got, `data-color='1'`) {
		t.Errorf("o componente do type não foi usado:\n%s", got)
	}
	if !strings.Contains(got, `v-model='form["Value"].b'`) {
		t.Errorf("o componente do type não recebeu o valor certo:\n%s", got)
	}

	// replacing the default is allowed: it is a type like any other
	b2 := New().Type(DefaultType, func(c *Context) h.HTMLComponent {
		return h.Div().Attr("data-mine", "1")
	})
	if !strings.Contains(render(t, b2, "{a str}"), `data-mine='1'`) {
		t.Error("o type padrão não foi substituído")
	}
}

// A type nobody registered is not silently drawn as text: the form says which
// field asked for what, and what there is.
func TestUnknownTypeIsReported(t *testing.T) {
	got := render(t, New(), "{a str; b mistério}")

	for _, want := range []string{"b", "mistério", "não tem componente registrado"} {
		if !strings.Contains(got, want) {
			t.Errorf("o erro não diz %q:\n%s", want, got)
		}
	}
	// the rest of the form is still drawn
	if !strings.Contains(got, `v-model='form["Value"].a'`) {
		t.Error("um type desconhecido derrubou o resto do form")
	}
}

func TestTypesListsWhatIsRegistered(t *testing.T) {
	names := New().Type("color", nil).Types()

	var hasForm, hasStr, hasColor bool
	for _, n := range names {
		switch n {
		case FormType:
			hasForm = true
		case DefaultType:
			hasStr = true
		case "color":
			hasColor = true
		}
	}
	if !hasForm || !hasStr || !hasColor {
		t.Errorf("Types() = %v", names)
	}
}

// The types a Builder knows out of the box, each drawn by the component that
// suits it — and a number bound as a NUMBER, since the value is written out
// again and a quoted one would come back a string.
func TestDefaultTypes(t *testing.T) {
	got := render(t, New(), "{a str; n int; u uint; on bool; c color; f float; d decimal; t text; markup html; at time; day date; span duration}")

	for what, want := range map[string]string{
		"str":   `<v-text-field`,
		"int":   `v-model.number='form["Value"].n'`,
		"uint":  `v-model.number='form["Value"].u'`,
		"bool":  `<v-switch`,
		"color": `<v-color-picker`,
		"text":  `<v-textarea`,
		"html":  `<vx-tiptap-editor`,
		// a decimal keeps its text: through a JS number it would come back a float
		"float":   `v-model='form["Value"].f'`,
		"decimal": `v-model='form["Value"].d'`,
		// the gad time namespace: an instant, a day, and a span written as text
		"time":     `<vx-datetimepicker`,
		"date":     `<vx-datepicker`,
		"duration": `v-model='form["Value"].span'`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("%s: falta %s:\n%s", what, want, got)
		}
	}

	// uint does not go below zero
	if !strings.Contains(got, `min='0'`) {
		t.Error("uint sem piso em zero")
	}
	// each one binds its own field
	for _, want := range []string{
		`v-model='form["Value"].a'`,
		`v-model='form["Value"].on'`,
		`v-model='form["Value"].c'`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("falta %s", want)
		}
	}
}

// A field typed with an enum is a select over the enum's values.
func TestEnumFieldIsASelect(t *testing.T) {
	got := render(t, New(), "enum Perm { Read, Write }\ninterface {perm Perm}")

	for _, want := range []string{
		`<v-select`,
		`{"title":"Read","value":"Read"}`,
		`v-model='form["Value"].perm'`,
		`required`, // not optional: it may not be left empty
	} {
		if !strings.Contains(got, want) {
			t.Errorf("o select não traz %s:\n%s", want, got)
		}
	}
}

// Empty is only allowed where the schema said it is: the `?` after the name.
// A required input carries `required`; an optional one does not, and may be
// cleared.
func TestOnlyAnOptionalFieldMayBeEmpty(t *testing.T) {
	optional := render(t, New(), "enum Perm { Read, Write }\ninterface {perm? Perm}")
	if strings.Contains(optional, "required") {
		t.Errorf("o field opcional foi marcado como obrigatório:\n%s", optional)
	}
	if !strings.Contains(optional, `:clearable='true'`) {
		t.Errorf("o field opcional não pode ser limpo:\n%s", optional)
	}

	required := render(t, New(), "enum Perm { Read, Write }\ninterface {perm Perm}")
	if !strings.Contains(required, "required") {
		t.Errorf("o field obrigatório não foi marcado:\n%s", required)
	}
	if strings.Contains(required, `:clearable='true'`) {
		t.Errorf("o field obrigatório pode ser limpo:\n%s", required)
	}

	// and the same for a plain text field
	if strings.Contains(render(t, New(), "{a? str}"), "required") {
		t.Error("texto opcional marcado como obrigatório")
	}
	if !strings.Contains(render(t, New(), "{a str}"), "required") {
		t.Error("texto obrigatório não marcado")
	}
}

// An application may take over an enum by its name, like any other type.
func TestEnumTypeCanBeOverridden(t *testing.T) {
	b := New().Type("Perm", func(c *Context) h.HTMLComponent {
		return h.Div().Attr("data-perm", "1")
	})
	if !strings.Contains(render(t, b, "enum Perm { Read }\ninterface {perm Perm}"), `data-perm='1'`) {
		t.Error("o componente registrado para o enum não foi usado")
	}
}

// The words around a field come from the FieldInfoFunc, asked by PATH: the
// names from the root down. A list adds no name of its own, so a field of the
// records in `links` is `links.href`.
func TestFieldInfoByPath(t *testing.T) {
	var asked []string

	b := New().FieldInfo(func(ctx *web.EventContext, path string) FieldInfo {
		asked = append(asked, path)
		switch path {
		case "title":
			return FieldInfo{Label: "O título", Hint: "aparece no topo"}
		case "links.href":
			return FieldInfo{Label: "Endereço", Help: h.Div(h.Text("comece com https://"))}
		}
		return FieldInfo{}
	})

	got := render(t, b, "{title str; sub interface {note str}; links interface[] {href str}}")

	// every field is asked for, by its own path
	for _, want := range []string{"title", "sub", "sub.note", "links", "links.href"} {
		var found bool
		for _, p := range asked {
			if p == want {
				found = true
			}
		}
		if !found {
			t.Errorf("o path %q não foi perguntado; perguntou: %v", want, asked)
		}
	}

	// the label replaces the humanized name, and the hint shows
	if !strings.Contains(got, `label='O título'`) {
		t.Errorf("o label não veio do FieldInfo:\n%s", got)
	}
	if !strings.Contains(got, `hint='aparece no topo'`) {
		t.Errorf("o hint não veio do FieldInfo:\n%s", got)
	}

	// a field with no answer keeps the humanized name
	if !strings.Contains(got, `label='Note'`) {
		t.Errorf("o fallback do label não é o nome humanizado:\n%s", got)
	}

	// the help hangs from a `?` beside the field
	if !strings.Contains(got, "mdi-help-circle-outline") {
		t.Errorf("o help não trouxe o ícone:\n%s", got)
	}
	if !strings.Contains(got, "comece com https://") {
		t.Errorf("o help não foi desenhado:\n%s", got)
	}
}

// Without the function nothing changes: the label is the name humanized, and
// there is no `?`.
func TestFieldInfoIsOptional(t *testing.T) {
	got := render(t, New(), "{title str}")

	if !strings.Contains(got, `label='Title'`) {
		t.Errorf("label = ?:\n%s", got)
	}
	if strings.Contains(got, "mdi-help-circle-outline") {
		t.Error("apareceu um ? sem help")
	}
}

// The values an enum field offers come from the EnumInfoFunc, asked by the
// field's PATH — one entry per item, in the order they should be offered. The
// path, and not the enum's name, because an enum need not have one.
func TestEnumInfoByPath(t *testing.T) {
	var asked []string

	b := New().EnumInfo(func(ctx *web.EventContext, path string) []EnumItem {
		asked = append(asked, path)
		if path == "perm" {
			return []EnumItem{
				{Name: "Write", Label: "Escrita", Hint: "pode alterar"},
				{Name: "Read"}, // no label: shows its own name
			}
		}
		return nil
	})

	got := render(t, b, "enum Perm { Read, Write }\nenum Cor { Azul }\ninterface {perm Perm; cor Cor}")

	// asked for each enum field, by its path
	for _, want := range []string{"perm", "cor"} {
		var found bool
		for _, p := range asked {
			if p == want {
				found = true
			}
		}
		if !found {
			t.Errorf("o path %q não foi perguntado; perguntou: %v", want, asked)
		}
	}

	// the answer decides the order and the words, and the value is what goes
	// into the record
	if !strings.Contains(got, `"title":"Escrita","value":"Write"`) {
		t.Errorf("o item não veio do EnumInfo:\n%s", got)
	}
	if !strings.Contains(got, `{"title":"Read","value":"Read"}`) {
		t.Errorf("um item sem label devia mostrar o próprio valor:\n%s", got)
	}
	if !strings.Contains(got, `:item-value='"value"'`) || !strings.Contains(got, `:item-title='"title"'`) {
		t.Errorf("o select não separa valor de rótulo:\n%s", got)
	}
	// an item's hint rides as its subtitle, and only then does the select read
	// each item's own props
	if !strings.Contains(got, `"subtitle":"pode alterar"`) {
		t.Errorf("o hint do item não virou subtítulo:\n%s", got)
	}
	if !strings.Contains(got, `:item-props`) {
		t.Errorf("o select não leu os props do item:\n%s", got)
	}

	// a field it did not answer for keeps what its enum declared
	if !strings.Contains(got, `{"title":"Azul","value":"Azul"}`) {
		t.Errorf("o enum declarado não foi usado onde o EnumInfo não respondeu:\n%s", got)
	}
}

// Without the function the field offers what the enum declared, each showing
// its own name.
func TestEnumInfoIsOptional(t *testing.T) {
	got := render(t, New(), "enum Perm { Read, Write }\ninterface {perm Perm}")

	for _, want := range []string{
		`{"title":"Read","value":"Read"}`,
		`{"title":"Write","value":"Write"}`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("falta %s:\n%s", want, got)
		}
	}
}

// What the form STORES is the enum member's NAME, never the number behind it —
// the value is written out as YAML and read back by name, so an enum that
// renumbers its members must not rewrite every record.
func TestEnumStoresTheMemberName(t *testing.T) {
	got := render(t, New(), "enum Perm { Read = 4, Write = 8 }\ninterface {perm Perm}")

	for _, want := range []string{
		`{"title":"Read","value":"Read"}`,
		`{"title":"Write","value":"Write"}`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("falta %s:\n%s", want, got)
		}
	}
	for _, unwanted := range []string{`"value":"4"`, `"value":4`, `"value":"8"`, `"value":8`} {
		if strings.Contains(got, unwanted) {
			t.Errorf("o número do enum entrou no que se grava (%s):\n%s", unwanted, got)
		}
	}
}

// A list of PLAIN VALUES: the sorter iterates over the value and the ITEM is
// the value, so it binds at the item's index — a primitive cannot be written
// back through the slot's `item` variable.
func TestComponentFuncDrawsAListOfValues(t *testing.T) {
	got := render(t, New(), "[]str")

	for _, want := range []string{
		`vx-array-sorter`,
		`v-model='form["Value"]'`,
		`v-slot:default`,
		`v-for='(item, itemIndex) in form["Value"]'`,
		`v-model='form["Value"][itemIndex]'`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("a lista de valores não liga %s:\n%s", want, got)
		}
	}
	// the list carries the words; a label on every row would only repeat them
	if strings.Count(got, `label='Value'`) > 1 {
		t.Errorf("o item da lista repete o label da lista:\n%s", got)
	}
}

// A field that is a list of values is that same list, one level down.
func TestComponentFuncDrawsAFieldListOfValues(t *testing.T) {
	got := render(t, New(), "{title str; tags []str}")

	for _, want := range []string{
		`v-model='form["Value"].title'`,
		`v-model='form["Value"].tags'`,
		`v-model='form["Value"].tags[itemIndex]'`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("o field lista não liga %s:\n%s", want, got)
		}
	}
}

// The sorter sorts; adding and removing an item are the form's own, so a list
// carries a button per row to remove it and one at the end to add another —
// shaped like the ones already there.
func TestComponentFuncListAddsAndRemoves(t *testing.T) {
	got := render(t, New(), "[]{label str; count int; on bool}")

	for _, want := range []string{
		`@click='form["Value"].splice(itemIndex, 1)'`,
		`@click='form["Value"].push({"label": "", "count": 0, "on": false})'`,
		`Adicionar`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("a lista não traz %s:\n%s", want, got)
		}
	}

	// a list of plain values pushes the empty value of its type
	if got := render(t, New(), "[]str"); !strings.Contains(got, `push("")`) {
		t.Errorf("a lista de valores não acrescenta um item vazio:\n%s", got)
	}
	if got := render(t, New(), "[]int"); !strings.Contains(got, `push(0)`) {
		t.Errorf("a lista de int não acrescenta um zero:\n%s", got)
	}
}

// A list the user may not write shows neither button.
func TestComponentFuncReadOnlyListHasNoButtons(t *testing.T) {
	schema, err := Parse("[]{label str}")
	if err != nil {
		t.Fatal(err)
	}
	comp := New().ComponentFunc(schema)(&presets.FieldContext{
		ToComponentOptions: &presets.ToComponentOptions{},
		Name:               "Value",
		FormKey:            "Value",
		Mode:               presets.FieldModeStack{presets.DETAIL},
	}, &web.EventContext{})

	out, err := h.Marshal(comp, context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, unwanted := range []string{"Adicionar", "splice("} {
		if strings.Contains(string(out), unwanted) {
			t.Errorf("uma lista só de leitura traz %q:\n%s", unwanted, out)
		}
	}
}

// While sorting, a row is the sorter's own and shows one field: the one a
// reader would name the row by.
func TestListItemTitleField(t *testing.T) {
	for src, want := range map[string]string{
		"[]{icon str; label str}": "label", // a name beats the first field
		"[]{icon str; title str}": "title",
		"[]{icon str; href str}":  "icon", // no name: the first text field
		"[]{n int; label str}":    "label",
		"[]{n int; u uint}":       "", // nothing to show
	} {
		got := render(t, New(), src)
		if want == "" {
			if strings.Contains(got, "item-title=") {
				t.Errorf("%s: não há field de texto para o título:\n%s", src, got)
			}
			continue
		}
		if !strings.Contains(got, "item-title='"+want+"'") {
			t.Errorf("%s: o título da linha devia ser %q:\n%s", src, want, got)
		}
	}
}
