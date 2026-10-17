package schemaform

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// show renders the stored value (YAML) the way the schema describes it: the
// detail when compact is false, the listing cell when it is true.
func show(t *testing.T, b *Builder, src, stored string, compact bool) string {
	t.Helper()

	schema, err := b.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	value, err := b.DecodeValue(stored)
	if err != nil {
		t.Fatal(err)
	}

	var comp h.HTMLComponent
	if compact {
		comp = b.ListComponent(&web.EventContext{}, schema, value)
	} else {
		comp = b.DetailComponent(&web.EventContext{}, schema, value)
	}
	out, err := h.Marshal(comp, context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func mustContain(t *testing.T, got string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("want %q in:\n%s", w, got)
		}
	}
}

func mustNotContain(t *testing.T, got string, unwanted ...string) {
	t.Helper()
	for _, w := range unwanted {
		if strings.Contains(got, w) {
			t.Errorf("did not want %q in:\n%s", w, got)
		}
	}
}

// A read view binds nothing: it shows the value itself.
func TestDisplayBindsNothing(t *testing.T) {
	got := show(t, New(), "export class default { title str }", "title: Hello", false)
	mustContain(t, got, "Title", "Hello")
	mustNotContain(t, got, "v-model", "v-text-field", "<input")
}

// A record in a detail: each field under its label, an empty one as a dash.
func TestDisplayRecordDetail(t *testing.T) {
	got := show(t, New(), "export class default { title str; note? str }", "title: Hello", false)
	mustContain(t, got, "Title", "Hello", "Note", "—")
}

// A record in a cell: on one line, and an empty field is left out.
func TestDisplayRecordCompact(t *testing.T) {
	got := show(t, New(), "export class default { title str; note? str; count int }", "title: Hello\ncount: 3", true)
	mustContain(t, got, "Title: ", "Hello", "Count: ", "3")
	mustNotContain(t, got, "Note", "<div")
}

// A list of plain values: a bulleted list in a detail, commas in a cell.
func TestDisplayValues(t *testing.T) {
	stored := "- Woburn, MA\n- Newton, MA\n"

	got := show(t, New(), "export type default []str", stored, false)
	mustContain(t, got, "<li>Woburn, MA</li>", "<li>Newton, MA</li>")

	got = show(t, New(), "export type default []str", stored, true)
	mustContain(t, got, "Woburn, MA, Newton, MA")
	mustNotContain(t, got, "<li>")
}

// The default list of records: every record a card, a space between them.
func TestDisplayRecordsDetail(t *testing.T) {
	got := show(t, New(), "class Item { label str; href }\nexport type default []Item", "- label: A\n  href: /a\n- label: B\n  href: /b\n", false)
	mustContain(t, got, "A", "/a", "B", "/b", "Href")
	if n := strings.Count(got, "<v-card "); n != 2 {
		t.Errorf("want a card per record, got %d:\n%s", n, got)
	}
	if strings.Count(got, "class='mt-3'") != 1 {
		t.Errorf("a space between the cards, not before the first:\n%s", got)
	}
}

// The table layout: the columns it names, in its order, a row per record.
func TestDisplayTableDetail(t *testing.T) {
	src := `class Item { label str; href; hidden? str }
[layout={name: "table", columns: [#href, #label]}]
export type default []Item`
	got := show(t, New(), src, "- label: A\n  href: /a\n  hidden: secret\n", false)
	mustContain(t, got, "<v-table", "<th", "Href", "Label", "<td", "/a")
	mustNotContain(t, got, "Hidden", "secret")
	if strings.Index(got, "Href") > strings.Index(got, "Label") {
		t.Errorf("the columns keep the order the schema names:\n%s", got)
	}
}

// A list of records in a cell: the title of each, and how many are left out.
func TestDisplayRecordsCompact(t *testing.T) {
	var stored strings.Builder
	for i := 1; i <= ListMaxItems+2; i++ {
		fmt.Fprintf(&stored, "- label: Item %d\n  href: /%d\n", i, i)
	}
	got := show(t, New(), "class Item { href; label str }\nexport type default []Item", stored.String(), true)
	mustContain(t, got, "Item 1", fmt.Sprintf("Item %d", ListMaxItems), "+2")
	mustNotContain(t, got, fmt.Sprintf("Item %d", ListMaxItems+1), "/1", "<v-table")
}

// An enum shows the label it is offered by, not its name.
func TestDisplayEnumShowsTheLabel(t *testing.T) {
	b := New().EnumInfo(func(_ *web.EventContext, path string) []EnumItem {
		return []EnumItem{{Name: "Read", Label: "Leitura"}, {Name: "Write", Label: "Escrita"}}
	})
	got := show(t, b, "enum Perm { Read, Write }\nexport class default { perm Perm }", "perm: Write", false)
	mustContain(t, got, "Escrita")
	mustNotContain(t, got, ">Write<")
}

// A list the application gives (EnumItemsFunc) shows the label as well.
func TestDisplayEnumItemsShowsTheLabel(t *testing.T) {
	b := New().EnumItems(func(_ *web.EventContext, path string, _ *Field) ([]EnumItem, error) {
		if path == "locale" {
			return []EnumItem{{Name: "pt-BR", Label: "Português"}}, nil
		}
		return nil, nil
	})
	got := show(t, b, "export class default { locale str }", "locale: pt-BR", true)
	mustContain(t, got, "Português")
}

func TestDisplayBool(t *testing.T) {
	got := show(t, New(), "export class default { on bool; off bool }", "on: true\noff: false", false)
	mustContain(t, got, "mdi-check", "mdi-minus")
}

// A color is a swatch of it; a value that is not a color never reaches the
// style attribute.
func TestDisplayColor(t *testing.T) {
	got := show(t, New(), "export class default { c color }", "c: '#ff0000'", false)
	mustContain(t, got, "background:#ff0000", "#ff0000")

	got = show(t, New(), "export class default { c color }", `c: "red;background:url(x)"`, false)
	mustContain(t, got, "background:transparent")
	mustNotContain(t, got, "url(x)\"", "background:red;")
}

// html is sanitized in a detail, and only its text in a cell.
func TestDisplayHTML(t *testing.T) {
	stored := `body: "<p>Hi <b>there</b><script>alert(1)</script></p>"`

	got := show(t, New(), "export class default { body html }", stored, false)
	mustContain(t, got, "<b>there</b>")
	mustNotContain(t, got, "<script", "alert(1)")

	got = show(t, New(), "export class default { body html }", stored, true)
	mustContain(t, got, "Hi there")
	mustNotContain(t, got, "<b>", "<p>")
}

// text keeps its line breaks in a detail.
func TestDisplayLongText(t *testing.T) {
	got := show(t, New(), "export class default { t text }", "t: |\n  one\n  two\n", false)
	mustContain(t, got, "white-space:pre-wrap", "one\ntwo")
}

// The labels come from FieldInfoFunc, as in the editor — in the table header
// too, with the hint behind it.
func TestDisplayUsesFieldInfo(t *testing.T) {
	b := New().FieldInfo(func(_ *web.EventContext, path string) FieldInfo {
		if path == "href" {
			return FieldInfo{Label: "Endereço", Hint: "para onde vai"}
		}
		return FieldInfo{}
	})
	got := show(t, b, "export class default { href }", "href: /a", false)
	mustContain(t, got, "Endereço")

	got = show(t, b, `class Item { href }
[layout="table"]
export type default []Item`, "- href: /a\n", false)
	mustContain(t, got, "Endereço", `title='para onde vai'`)
}

// A form inside the form is shown the same way, one level down.
func TestDisplayNestedForm(t *testing.T) {
	got := show(t, New(), "class Link { href str }\nexport class default { title str; links []Link }",
		"title: T\nlinks:\n  - href: /x\n  - href: /y\n", false)
	mustContain(t, got, "Links", "/x", "/y")
}

// An empty list: a dash in a detail, nothing in a cell.
func TestDisplayEmptyList(t *testing.T) {
	mustContain(t, show(t, New(), "export type default []str", "", false), "—")
	if got := show(t, New(), "export type default []str", "", true); strings.TrimSpace(got) != "" {
		t.Errorf("an empty list in a cell shows nothing, got %q", got)
	}
}

// A type an application registers a display for is shown by it.
func TestDisplayRegistry(t *testing.T) {
	b := New().Display("str", func(c *Context) h.HTMLComponent {
		return h.Span(fmt.Sprint(c.Data)).Attr("data-mine", "1")
	})
	mustContain(t, show(t, b, "export class default { a str }", "a: x", false), `data-mine='1'`)

	if _, ok := New().DisplayFunc(FormType); !ok {
		t.Error("FormType always has a display")
	}
}

// A Record (a value read back from the form) shows like the stored map.
func TestDisplayRecordValue(t *testing.T) {
	schema, err := Parse("export class default { title str }")
	if err != nil {
		t.Fatal(err)
	}
	out, err := h.Marshal(New().DetailComponent(&web.EventContext{}, schema,
		Record{{Name: "title", Value: "From the form"}}), context.Background())
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), "From the form")
}

// A grid shows each record as a card of its labelled fields, columns per
// row; a cell keeps to the titles.
func TestDisplayGrid(t *testing.T) {
	src := `class Item { label str; href }
[layout={name: "grid", columns: 2}]
export type default []Item`
	stored := "- label: A\n  href: /a\n- label: B\n  href: /b\n"
	got := show(t, New(), src, stored, false)
	mustContain(t, got, "<v-row", `:md='6'`, "<v-card", "Label", "/a", "/b")
	if strings.Count(got, "<v-card ") != 2 {
		t.Errorf("one card per record:\n%s", got)
	}
	mustNotContain(t, show(t, New(), src, stored, true), "<v-card", "<v-row")
}

// A record inside the record is set in, behind a line, in the detail as in the
// form; a cell keeps to one line.
func TestDisplayNestedRecordIsSetIn(t *testing.T) {
	src := "class Address { street str }\nexport class default { name str; address Address }"
	got := show(t, New(), src, "name: A\naddress:\n  street: Rua 1\n", false)
	mustContain(t, got, `class='`+nestedRecordClass+`'`, "Rua 1")
	if strings.Count(got, nestedRecordClass) != 1 {
		t.Errorf("only the inner record is set in:\n%s", got)
	}
	mustNotContain(t, show(t, New(), src, "name: A\naddress:\n  street: Rua 1\n", true), nestedRecordClass)
}

// A field by its path, shown as the detail shows it — by the display of its
// type —; a path of no field, nothing.
func TestFieldDisplay(t *testing.T) {
	b := New()
	schema, err := b.Parse("class Address { city str; zip? str }\nexport class default { title str; done bool; address Address }")
	if err != nil {
		t.Fatal(err)
	}
	value, err := b.DecodeValue("title: Hello\ndone: true\naddress:\n  city: Boston\n")
	if err != nil {
		t.Fatal(err)
	}
	show, portal := b.FieldDisplay(&web.EventContext{}, schema, value)
	if portal == nil {
		t.Error("no portal")
	}
	html := func(path string) string {
		comp := show(path)
		if comp == nil {
			return "<nil>"
		}
		out, err := h.Marshal(comp, context.Background())
		if err != nil {
			t.Fatal(err)
		}
		return string(out)
	}
	mustContain(t, html("title"), "Hello")
	mustNotContain(t, html("title"), "Title", "v-model")
	mustContain(t, html("address.city"), "Boston")
	mustContain(t, html("address.zip"), "—")
	// a bool by its display (BoolDisplayFunc): its icon, not its text
	mustContain(t, html("done"), "mdi-check")
	for _, p := range []string{"nothing", "title.x", "address.nothing"} {
		if got := html(p); got != "<nil>" {
			t.Errorf("%s: %s", p, got)
		}
	}
}

// A field's warnings — its record's $warnings of it — under its value, in
// every read-only view: the detail, a field by its path, a read-only form.
func TestDisplayWarnings(t *testing.T) {
	b := New()
	schema, err := b.Parse("export class default { title str; zip str }")
	if err != nil {
		t.Fatal(err)
	}
	value := map[string]any{"title": "Hello", "zip": "00000",
		WarningsField: []any{map[string]any{"fieldName": "zip", "message": "not checked"}}}
	html := func(comp h.HTMLComponent) string {
		out, err := h.Marshal(comp, context.Background())
		if err != nil {
			t.Fatal(err)
		}
		return string(out)
	}
	detail := html(b.DetailComponent(&web.EventContext{}, schema, value))
	mustContain(t, detail, "not checked", WarningClass)
	if strings.Count(detail, "data-schemaform-warning") != 1 {
		t.Errorf("one warning, of zip:\n%s", detail)
	}
	show, _ := b.FieldDisplay(&web.EventContext{}, schema, value)
	mustContain(t, html(show("zip")), "not checked")
	mustNotContain(t, html(show("title")), "not checked")
	form := html(b.ComponentFunc(schema)(&presets.FieldContext{ToComponentOptions: &presets.ToComponentOptions{},
		Name: "Value", FormKey: "Value", Mode: presets.FieldModeStack{presets.DETAIL}, ValueOverride: value}, &web.EventContext{}))
	mustContain(t, form, "not checked")
}

// A prefix and a suffix (`[prefix="R$ "]`, `[suffix=" kg"]`): in the input of
// a text or a number, and with the value shown; a number shown as the
// language writes it, with its decimal places (`[places=N]`, else its
// step's, else its own) — places a decimal's input's step.
func TestPrefixSuffix(t *testing.T) {
	b := New()
	src := `export class default { [prefix="R$ ", places=2] value decimal; [suffix=" kg", step=0.5] weight float; [prefix="#"] code str; count int }`
	schema, err := b.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	form := render(t, b, src)
	mustContain(t, form, `prefix='R$ '`, `suffix=' kg'`, `prefix='#'`, `step='0.01'`, `step='0.5'`)

	value := map[string]any{"value": "10", "weight": 72.5, "code": "A1", "count": 1234}
	out, err := h.Marshal(b.DetailComponent(&web.EventContext{}, schema, value), context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// no request: English
	mustContain(t, string(out), "R$ 10.00", "72.5 kg", "#A1", "1,234")

	for _, c := range []struct {
		lang     string
		v        any
		decimals int
		want     string
	}{
		{"pt-BR", "10", 2, "10,00"},
		{"pt-BR", "1234.5", 2, "1.234,50"},
		{"pt-BR", 1234.567, -1, "1.234,567"},
		{"en", "1234.5", 0, "1,235"},
		{"pt-BR", "x", 2, "x"},
	} {
		if got := FormatNumber(language.MustParse(c.lang), c.v, c.decimals); got != c.want {
			t.Errorf("%s %v %d: %q, want %q", c.lang, c.v, c.decimals, got, c.want)
		}
	}
}

// A date, a time, both: as the language writes them (the short forms of
// i18n): with no request, English; with pt-BR's messages, 10/03/2026 08:05.
func TestDateTimeDisplay(t *testing.T) {
	b := New()
	schema, err := b.Parse(`export class default {
	day date
	at time
	hour time
	[prefix="em "] other date
	bad date
}`)
	if err != nil {
		t.Fatal(err)
	}
	value := map[string]any{"day": "2026-03-10", "at": "2026-03-10T08:05:00", "hour": "08:05",
		"other": "2026-12-25", "bad": "not a date"}
	show := func(ev *web.EventContext) string {
		out, err := h.Marshal(b.DetailComponent(ev, schema, value), context.Background())
		if err != nil {
			t.Fatal(err)
		}
		return string(out)
	}
	mustContain(t, show(&web.EventContext{}), "2026-03-10", "2026-03-10 08:05:00", "08:05:00", "em 2026-12-25", "not a date")

	// pt-BR: its messages on the request
	ib := i18n.New().SupportLanguages(language.BrazilianPortuguese).
		RegisterForModule(language.BrazilianPortuguese, i18n.DefaultKey, i18n.Default_pt_BR)
	var r *http.Request
	ib.EnsureLanguage(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) { r = req })).
		ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	mustContain(t, show(&web.EventContext{R: r}), "10/03/2026", "10/03/2026 08:05", ">08:05<", "em 25/12/2026")
}

// An int of places is fixed point — the integer is the value × 10^places,
// money in cents —: its form holds the value as shown (Filled: 1050 is
// "10.50"), read back the integer (Decoder: "10.5", "10,50" are 1050; more
// places, refused); its limits, of the integer, shown divided; shown as the
// value: 1050 of places 2 is "R$ 10.50".
func TestFixedPointInt(t *testing.T) {
	b := New()
	src := `export class default {
	[prefix="R$ ", places=2, min=100, max=500000] price int
	[places=1] score uint
	count int
}`
	form := render(t, b, src)
	mustContain(t, form, `v-model='form["Value"].price'`,
		`step='0.01'`, `min='1'`, `max='5000'`, `step='0.1'`,
		"Math.round(Number(v) * 100) >= 100", "the smallest value is 1", "the largest value is 5000", `\d{0,2}`)
	schema, err := b.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	filled := schema.Filled(map[string]any{"price": 1050, "score": 47, "count": 3}).(map[string]any)
	if filled["price"] != "10.50" || filled["score"] != "4.7" || filled["count"] != 3 {
		t.Errorf("filled: %v", filled)
	}
	for in, want := range map[string]any{"10.5": int64(1050), "10,50": int64(1050), "7": int64(700)} {
		rec, issues := (&Decoder{Schema: schema}).Decode(map[string]any{"price": in, "score": "4.7", "count": "3"})
		if len(issues) != 0 || rec[0].Value != want || rec[1].Value != uint64(47) {
			t.Errorf("%s: %v %v", in, rec, issues)
		}
	}
	for _, in := range []string{"10.501", "x", "0.5"} {
		if _, issues := (&Decoder{Schema: schema}).Decode(map[string]any{"price": in, "score": "1", "count": "3"}); len(issues) == 0 {
			t.Errorf("%s: no issue", in)
		}
	}
	out, err := h.Marshal(b.DetailComponent(&web.EventContext{}, schema,
		map[string]any{"price": 1050, "score": 47, "count": 1234}), context.Background())
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), "R$ 10.50", "4.7", "1,234")
	if got := DisplayText(context.Background(), schema.Fields[0], 123456); got != "R$ 1,234.56" {
		t.Errorf("DisplayText: %q", got)
	}
}
