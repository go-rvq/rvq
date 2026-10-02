package userdocs

import (
	"regexp"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gad-lang/gad"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func testSource(readme string) *Source {
	return &Source{
		Package: "example.com/pkg",
		FS: fstest.MapFS{
			"en/things/README.md":                {Data: []byte(readme)},
			"en/things/actions/do.md":            {Data: []byte("# Do\n\nIt does.\n")},
			"en/things/images/list.png":          {Data: []byte("png")},
			"en/things/children/parts/README.md": {Data: []byte("# Parts\n")},
		},
	}
}

func TestSourceFiles(t *testing.T) {
	src := testSource("# Things\n")
	lang, err := src.Lang()
	if err != nil || lang != "en" {
		t.Fatalf("lang %q, %v", lang, err)
	}
	files, err := src.Files()
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, f := range files {
		paths = append(paths, f.Path)
	}
	if strings.Join(paths, ",") != "things/README.md,things/actions/do.md,things/children/parts/README.md" {
		t.Errorf("files %q", paths)
	}
	two := &Source{FS: fstest.MapFS{"en/a.md": {}, "pt-BR/a.md": {}}}
	if lang, err := two.Lang(); err != nil || lang != "en" {
		t.Errorf("the source of en and pt-BR: %q, %v", lang, err)
	}
	if two.LangOf("pt-BR") != "pt-BR" || two.LangOf("pt") != "pt-BR" || two.LangOf("es") != "" {
		t.Error("LangOf")
	}
	noEn := &Source{FS: fstest.MapFS{"es/a.md": {}, "pt-BR/a.md": {}}}
	if _, err := noEn.Lang(); err == nil {
		t.Error("two languages and no en were accepted")
	}
}

// The pictures of a language: its own, else the source's.
func TestSourceAsset(t *testing.T) {
	src := &Source{FS: fstest.MapFS{
		"en/a/images/x.png": {}, "en/a/images/y.png": {},
		"pt-BR/a/images/x.png": {},
	}}
	for _, c := range []struct{ locale, file, want string }{
		{"pt-BR", "a/images/x.png", "pt-BR"},
		{"pt-BR", "a/images/y.png", "en"},
		{"en", "a/images/x.png", "en"},
		{"", "a/images/x.png", "en"},
	} {
		if got, err := src.Asset(c.locale, c.file); err != nil || got != c.want {
			t.Errorf("Asset(%s, %s) = %s, %v; want %s", c.locale, c.file, got, err, c.want)
		}
	}
}

// A package written in en and in pt-BR: pt-BR's documents as they are; one
// pt-BR lacks, translated from en's.
func TestSyncWritten(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	src := testSource("# Things\n")
	src.FS.(fstest.MapFS)["pt-BR/things/README.md"] = &fstest.MapFile{Data: []byte("# Coisas\n")}
	src.FS.(fstest.MapFS)["pt-BR/things/only.md"] = &fstest.MapFile{Data: []byte("# Só aqui\n")}
	jobs, err := Sync(db, []*Source{src}, []string{"en", "pt-BR"})
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 || jobs[0].Locale != "pt-BR" || len(jobs[0].Files) != 2 {
		t.Fatalf("jobs %+v", jobs)
	}
	Translate(db, &fakeTranslator{}, jobs[0])
	pt, _ := value(t, db, "pt-BR")
	if pt["things/README.md"] != "# Coisas\n" || pt["things/only.md"] != "# Só aqui\n" {
		t.Errorf("pt-BR's own %q", pt)
	}
	if !strings.Contains(pt["things/actions/do.md"], "[en>pt-br]") {
		t.Errorf("not translated %q", pt["things/actions/do.md"])
	}
	en, _ := value(t, db, "en")
	if _, ok := en["things/only.md"]; ok || en["things/README.md"] != "# Things\n" {
		t.Errorf("en %q", en)
	}
	// nothing changed: nothing to translate
	if jobs, _ = Sync(db, []*Source{src}, []string{"en", "pt-BR"}); len(jobs) != 0 {
		t.Errorf("a sync with nothing new left %+v", jobs)
	}
}

func TestDocFile(t *testing.T) {
	for node, file := range map[string]string{
		"things":                "things/README.md",
		"things/actions/do":     "things/actions/do.md",
		"things/children/parts": "things/children/parts/README.md",
		"groups/admin":          "groups/admin/README.md",
		"pages/db-tools":        "pages/db-tools/README.md",
		"posts/pages/report":    "posts/pages/report.md",
		"pages/pages/report":    "pages/pages/report.md",
		"posts/actions":         "posts/actions/README.md",
	} {
		if got := DocFile(node); got != file {
			t.Errorf("DocFile(%s) = %s, want %s", node, got, file)
		}
		if got := docNode(file); got != node {
			t.Errorf("docNode(%s) = %s, want %s", file, got, node)
		}
	}
}

func TestMask(t *testing.T) {
	src := "# Title\n\nOpen {%= admin.model(\"x\").link %} and `the code`.\n\n" +
		"- an item [a link](other.md)\n\n```\nfn()\n```\n![pic](images/a.png)\n"
	masked, masks := Mask(src)
	for _, keep := range []string{"{%=", "`the code`", "fn()", "(other.md)", "(images/a.png)", "# "} {
		if strings.Contains(masked, keep) {
			t.Errorf("%q not masked:\n%s", keep, masked)
		}
	}
	if !strings.Contains(masked, "Title") || !strings.Contains(masked, "an item") {
		t.Errorf("the words were masked:\n%s", masked)
	}
	if got := Unmask(masked, masks); got != src {
		t.Errorf("unmasked\n%s\nwant\n%s", got, src)
	}
}

func TestRenderTemplate(t *testing.T) {
	out, err := RenderTemplate("Open {%= admin.name %}{% if admin.on begin %}!{% end %}",
		gad.Dict{"admin": gad.Dict{"name": gad.Str("Languages"), "on": gad.True}})
	if err != nil || out != "Open Languages!" {
		t.Errorf("%q, %v", out, err)
	}
	if _, err := RenderTemplate("{%= nope( %}", gad.Dict{"admin": gad.Dict{}}); err == nil {
		t.Error("a bad template was accepted")
	}
}

type fakeTranslator struct{ calls int }

func (f *fakeTranslator) Enabled() bool                 { return true }
func (f *fakeTranslator) Code(l string) (string, error) { return strings.ToLower(l), nil }
func (f *fakeTranslator) Translate(src, dst string, texts []string) ([]string, error) {
	f.calls++
	// the mark after the tokens that open a line (a heading's #, masked), as
	// a translation is: in place
	lead := regexp.MustCompile(`(?m)^((?:⟦\d+⟧)*)(\S)`)
	out := make([]string, len(texts))
	for i, t := range texts {
		out[i] = lead.ReplaceAllString(t, "${1}["+src+">"+dst+"]${2}")
	}
	return out, nil
}

func value(t *testing.T, db *gorm.DB, locale string) (val, ini map[string]string) {
	t.Helper()
	var row UserDocPackage
	if err := db.Take(&row, "locale_code = ? AND id = ?", locale, "example.com/pkg").Error; err != nil {
		t.Fatal(err)
	}
	v, _ := ParseFiles(row.Value)
	i, _ := ParseFiles(row.InitialValue)
	return filesByPath(v), filesByPath(i)
}

// The documents on boot: the source's language as they are, the others
// translated; an edit kept; a source changed translated again, alone.
func TestSync(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	src := testSource("# Things\n\nOpen {%= admin.model(\"things\").link %}.\n")
	tr := &fakeTranslator{}
	jobs, err := Sync(db, []*Source{src}, []string{"en", "pt-BR"})
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 || jobs[0].Locale != "pt-BR" || len(jobs[0].Files) != 3 {
		t.Fatalf("jobs %+v", jobs)
	}
	en, _ := value(t, db, "en")
	if en["things/README.md"] != "# Things\n\nOpen {%= admin.model(\"things\").link %}.\n" {
		t.Errorf("en %q", en["things/README.md"])
	}
	if err := Translate(db, tr, jobs[0]); err != nil {
		t.Fatal(err)
	}
	pt, ini := value(t, db, "pt-BR")
	if got := pt["things/README.md"]; !strings.Contains(got, "[en>pt-br]") || !strings.Contains(got, `{%= admin.model("things").link %}`) ||
		!strings.HasPrefix(got, "# ") || got != ini["things/README.md"] {
		t.Errorf("pt-BR %q", got)
	}

	// nothing changed: nothing to translate
	if jobs, _ = Sync(db, []*Source{src}, []string{"en", "pt-BR"}); len(jobs) != 0 {
		t.Errorf("a sync with nothing new left %+v", jobs)
	}

	// an edit of pt-BR; a change of the source's action
	var row UserDocPackage
	db.Take(&row, "locale_code = ? AND id = ?", "pt-BR", "example.com/pkg")
	files, _ := ParseFiles(row.Value)
	for i := range files {
		if files[i].Path == "things/README.md" {
			files[i].Content = "# Coisas editadas\n"
		}
	}
	edited, _ := FormatFiles(files)
	db.Model(&row).Update("value", edited)
	src.FS.(fstest.MapFS)["en/things/actions/do.md"] = &fstest.MapFile{Data: []byte("# Do\n\nIt does more.\n")}

	jobs, err = Sync(db, []*Source{src}, []string{"en", "pt-BR"})
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 || len(jobs[0].Files) != 1 || jobs[0].Files[0].Path != "things/actions/do.md" {
		t.Fatalf("jobs %+v", jobs)
	}
	Translate(db, tr, jobs[0])
	pt, _ = value(t, db, "pt-BR")
	if pt["things/README.md"] != "# Coisas editadas\n" {
		t.Errorf("the edit was lost: %q", pt["things/README.md"])
	}
	if !strings.Contains(pt["things/actions/do.md"], "It does more") {
		t.Errorf("the change was not translated: %q", pt["things/actions/do.md"])
	}

	// a package gone
	if _, err := Sync(db, nil, []string{"en"}); err != nil {
		t.Fatal(err)
	}
	var gone UserDocPackage
	db.Take(&gone, "locale_code = ? AND id = ?", "en", "example.com/pkg")
	if !gone.Unused {
		t.Error("a package gone is not unused")
	}
}

// A package written in pt-BR alone: en is translated from it.
func TestSyncFromPortuguese(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	src := &Source{Package: "example.com/pkg", FS: fstest.MapFS{"pt-BR/things/README.md": {Data: []byte("# Coisas\n")}}}
	jobs, err := Sync(db, []*Source{src}, []string{"en", "pt-BR"})
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 || jobs[0].Locale != "en" || jobs[0].SourceLocale != "pt-BR" {
		t.Fatalf("jobs %+v", jobs)
	}
	Translate(db, &fakeTranslator{}, jobs[0])
	en, _ := value(t, db, "en")
	if !strings.HasPrefix(en["things/README.md"], "# [pt-br>en]Coisas") {
		t.Errorf("en %q", en["things/README.md"])
	}
}
