package schemaform

import (
	"context"
	"fmt"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
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
	got := show(t, New(), "{title str}", "title: Hello", false)
	mustContain(t, got, "Title", "Hello")
	mustNotContain(t, got, "v-model", "v-text-field", "<input")
}

// A record in a detail: each field under its label, an empty one as a dash.
func TestDisplayRecordDetail(t *testing.T) {
	got := show(t, New(), "{title str; note? str}", "title: Hello", false)
	mustContain(t, got, "Title", "Hello", "Note", "—")
}

// A record in a cell: on one line, and an empty field is left out.
func TestDisplayRecordCompact(t *testing.T) {
	got := show(t, New(), "{title str; note? str; count int}", "title: Hello\ncount: 3", true)
	mustContain(t, got, "Title: ", "Hello", "Count: ", "3")
	mustNotContain(t, got, "Note", "<div")
}

// A list of plain values: a bulleted list in a detail, commas in a cell.
func TestDisplayValues(t *testing.T) {
	stored := "- Woburn, MA\n- Newton, MA\n"

	got := show(t, New(), "[]str", stored, false)
	mustContain(t, got, "<li>Woburn, MA</li>", "<li>Newton, MA</li>")

	got = show(t, New(), "[]str", stored, true)
	mustContain(t, got, "Woburn, MA, Newton, MA")
	mustNotContain(t, got, "<li>")
}

// A list of records laid out as forms: every record, with a line between them.
func TestDisplayRecordsDetail(t *testing.T) {
	got := show(t, New(), "[]{label str; href}", "- label: A\n  href: /a\n- label: B\n  href: /b\n", false)
	mustContain(t, got, "A", "/a", "B", "/b", "Href")
	if n := strings.Count(got, "<v-divider"); n != 1 {
		t.Errorf("want 1 divider between 2 records, got %d:\n%s", n, got)
	}
}

// The table layout: the columns it names, in its order, a row per record.
func TestDisplayTableDetail(t *testing.T) {
	src := `[layout="table", columns=[#href, #label]]
interface Form []{label str; href; hidden? str}`
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
	got := show(t, New(), "[]{href; label str}", stored.String(), true)
	mustContain(t, got, "Item 1", fmt.Sprintf("Item %d", ListMaxItems), "+2")
	mustNotContain(t, got, fmt.Sprintf("Item %d", ListMaxItems+1), "/1", "<v-table")
}

// An enum shows the label it is offered by, not its name.
func TestDisplayEnumShowsTheLabel(t *testing.T) {
	b := New().EnumInfo(func(_ *web.EventContext, path string) []EnumItem {
		return []EnumItem{{Name: "Read", Label: "Leitura"}, {Name: "Write", Label: "Escrita"}}
	})
	got := show(t, b, "enum Perm { Read, Write }\ninterface {perm Perm}", "perm: Write", false)
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
	got := show(t, b, "{locale str}", "locale: pt-BR", true)
	mustContain(t, got, "Português")
}

func TestDisplayBool(t *testing.T) {
	got := show(t, New(), "{on bool; off bool}", "on: true\noff: false", false)
	mustContain(t, got, "mdi-check", "mdi-minus")
}

// A color is a swatch of it; a value that is not a color never reaches the
// style attribute.
func TestDisplayColor(t *testing.T) {
	got := show(t, New(), "{c color}", "c: '#ff0000'", false)
	mustContain(t, got, "background:#ff0000", "#ff0000")

	got = show(t, New(), "{c color}", `c: "red;background:url(x)"`, false)
	mustContain(t, got, "background:transparent")
	mustNotContain(t, got, "url(x)\"", "background:red;")
}

// html is sanitized in a detail, and only its text in a cell.
func TestDisplayHTML(t *testing.T) {
	stored := `body: "<p>Hi <b>there</b><script>alert(1)</script></p>"`

	got := show(t, New(), "{body html}", stored, false)
	mustContain(t, got, "<b>there</b>")
	mustNotContain(t, got, "<script", "alert(1)")

	got = show(t, New(), "{body html}", stored, true)
	mustContain(t, got, "Hi there")
	mustNotContain(t, got, "<b>", "<p>")
}

// text keeps its line breaks in a detail.
func TestDisplayLongText(t *testing.T) {
	got := show(t, New(), "{t text}", "t: |\n  one\n  two\n", false)
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
	got := show(t, b, "{href}", "href: /a", false)
	mustContain(t, got, "Endereço")

	got = show(t, b, `[layout="table"]
interface Form []{href}`, "- href: /a\n", false)
	mustContain(t, got, "Endereço", `title='para onde vai'`)
}

// A form inside the form is shown the same way, one level down.
func TestDisplayNestedForm(t *testing.T) {
	got := show(t, New(), "{title str; links interface[] {href str}}",
		"title: T\nlinks:\n  - href: /x\n  - href: /y\n", false)
	mustContain(t, got, "Links", "/x", "/y")
}

// An empty list: a dash in a detail, nothing in a cell.
func TestDisplayEmptyList(t *testing.T) {
	mustContain(t, show(t, New(), "[]str", "", false), "—")
	if got := show(t, New(), "[]str", "", true); strings.TrimSpace(got) != "" {
		t.Errorf("an empty list in a cell shows nothing, got %q", got)
	}
}

// A type an application registers a display for is shown by it.
func TestDisplayRegistry(t *testing.T) {
	b := New().Display("str", func(c *Context) h.HTMLComponent {
		return h.Span(fmt.Sprint(c.Data)).Attr("data-mine", "1")
	})
	mustContain(t, show(t, b, "{a str}", "a: x", false), `data-mine='1'`)

	if _, ok := New().DisplayFunc(FormType); !ok {
		t.Error("FormType always has a display")
	}
}

// A Record (a value read back from the form) shows like the stored map.
func TestDisplayRecordValue(t *testing.T) {
	schema, err := Parse("{title str}")
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
