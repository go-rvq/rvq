package schemaform

import (
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/web"
)

// The schema is a gad program: it runs, and the form is read off the interface
// named Form through gad's reflection — metadata included, symbols as names.
func TestParseReadsMetadata(t *testing.T) {
	s, err := Parse(`[layout="table", columns=[#label, #icon], extra=2]
interface Form []{
	label str
	[width=120]
	icon str
	color? str
}`)
	if err != nil {
		t.Fatal(err)
	}

	if s.Layout != LayoutTable {
		t.Errorf("layout = %q, want %q", s.Layout, LayoutTable)
	}
	if got := strings.Join(s.Columns, ","); got != "label,icon" {
		t.Errorf("columns = %q", got)
	}
	if s.Meta["extra"] != int64(2) {
		t.Errorf("meta = %#v", s.Meta)
	}
	// a field's own block comes along
	if s.Fields[1].Meta["width"] != int64(120) {
		t.Errorf("icon.meta = %#v", s.Fields[1].Meta)
	}
	if !s.Fields[2].Nullable {
		t.Error("color? devia aceitar vazio")
	}
}

// With no layout, a list is drawn as forms, as it always was.
func TestParseDefaultLayout(t *testing.T) {
	s, err := Parse("[]{a str}")
	if err != nil {
		t.Fatal(err)
	}
	if s.Layout != LayoutForm || len(s.Columns) != 0 {
		t.Errorf("layout = %q, columns = %v", s.Layout, s.Columns)
	}
}

// What cannot be drawn is refused where the schema is read.
func TestParseRefusesWhatCannotBeDrawn(t *testing.T) {
	for name, tc := range map[string]struct{ src, want string }{
		"layout que não existe": {`[layout="grid"] interface Form []{a str}`, `"grid" does not exist`},
		"tabela de um registro": {`[layout="table"] interface Form {a str}`, "is not one"},
		"tabela de valores":     {`[layout="table"] interface Form []str`, "is not one"},
		"coluna sem field":      {`[layout="table", columns=[#a, #nope]] interface Form []{a str}`, `"nope" names no field`},
		"columns não é lista":   {`[layout="table", columns=#a] interface Form []{a str}`, "want a list"},
	} {
		_, err := Parse(tc.src)
		if err == nil {
			t.Errorf("%s: passou", name)
			continue
		}
		if !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: err = %v", name, err)
		}
	}
}

// The builder's types are in scope while the schema runs: a type gad does not
// have (`color`), and one the application registered, resolve — and `time`,
// which is gad's time NAMESPACE, is the builder's time type here.
func TestParseUsesTheBuilderTypes(t *testing.T) {
	b := New().Type("slug", TextComponentFunc)
	s, err := b.Parse("{a color; b slug; c time; d html; e int}")
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range []string{"color", "slug", "time", "html", "int"} {
		if got := s.Fields[i].Type; got != want {
			t.Errorf("%s = %q, want %q", s.Fields[i].Name, got, want)
		}
	}

	// a type only another builder registered is not in scope here
	if _, err := New().Parse("{b slug}"); err == nil {
		t.Error("um type que este builder não registrou passou")
	}
}

// A schema that does not return is not a schema: the run is cut short.
func TestParseRunIsBounded(t *testing.T) {
	defer func(old time.Duration) { RunTimeout = old }(RunTimeout)
	RunTimeout = 50 * time.Millisecond

	done := make(chan error, 1)
	go func() {
		_, err := Parse("interface Form { a str }\nfor { }")
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Error("um schema que não termina passou")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("o schema não foi interrompido")
	}
}

// A schema is read on every draw and every save; reading it again is the cache.
func TestParseIsCached(t *testing.T) {
	a, err := Parse("{cached str}")
	if err != nil {
		t.Fatal(err)
	}
	b, err := Parse("{cached str}")
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Error("o mesmo schema foi lido duas vezes")
	}
}

// A list of records laid out as a TABLE: one row per record, one column per
// field Columns names, in that order; the header carries the label, the cell
// only the input, on one line.
func TestTableLayout(t *testing.T) {
	b := New().FieldInfo(func(_ *web.EventContext, path string) FieldInfo {
		return FieldInfo{Label: map[string]string{"label": "Texto", "icon": "Ícone"}[path],
			Hint: map[string]string{"icon": "classe Font Awesome"}[path]}
	})
	got := render(t, b, `[layout="table", columns=[#icon, #label]]
interface Form []{label str; icon str; hidden str}`)

	for _, want := range []string{
		`<v-table`,
		`v-for='(item, itemIndex) in form["Value"]'`,
		`title='classe Font Awesome'`, // the hint, behind the header
		`v-model='item.icon'`,
		`v-model='item.label'`,
		`hide-details`,
		`@click='form["Value"].splice(itemIndex, 1)'`,
		`Adicionar`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("a tabela não traz %s:\n%s", want, got)
		}
	}

	// the columns come in the order Columns names them
	if i, j := strings.Index(got, ">Ícone<"), strings.Index(got, ">Texto<"); i < 0 || j < 0 || i > j {
		t.Errorf("colunas fora de ordem:\n%s", got)
	}
	// a field left out is not drawn — it keeps its value in the record
	if strings.Contains(got, "item.hidden") {
		t.Errorf("um field fora das colunas foi desenhado:\n%s", got)
	}
	// and a cell carries no label of its own: the header has it
	if strings.Contains(got, `label='Texto'`) || strings.Contains(got, `label='Ícone'`) {
		t.Errorf("a célula repete o label do cabeçalho:\n%s", got)
	}
}

// Without columns, a table shows every field, in the schema's order.
func TestTableLayoutAllColumns(t *testing.T) {
	got := render(t, New(), `[layout="table"] interface Form []{a str; b int}`)
	if !strings.Contains(got, "<v-table") {
		t.Fatalf("não é uma tabela:\n%s", got)
	}
	if i, j := strings.Index(got, "item.a"), strings.Index(got, "item.b"); i < 0 || j < 0 || i > j {
		t.Errorf("colunas = todos os fields, na ordem do schema:\n%s", got)
	}
}
