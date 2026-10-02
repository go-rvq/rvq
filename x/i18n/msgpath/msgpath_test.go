package msgpath

import (
	"reflect"
	"strings"
	"testing"
)

type html string

type Base struct {
	Save string `i18n:"hint='Button that saves.'"`
}

type Sub struct {
	Submit string `i18n:"label='Submit button'"`
}

type Option struct {
	Name  string `yaml:"name"`
	Value int    `yaml:"value"`
}

type msgs struct {
	ModuleDescription string `i18n:"hint='The description of the module.'"`
	Base
	Title    string `i18n:"label='Page title', hint='The title of the page.'"`
	Deleted  string `i18n:"hint='Shown once records were deleted.', fields=(;'%d'='how many', '%[2]s'='the model\\'s name')"`
	Body     html   `i18n:"type=html"`
	Form     Sub
	Ptr      *Sub
	Help     [][2]string `i18n:"hint='Sections: a title and its text.'"`
	Options  []Option    `i18n:"type=form, schema='x'"`
	Count    int
	Format   func(int) string
	internal string
}

func newMsgs() *msgs {
	return &msgs{
		ModuleDescription: "The test messages.",
		Base:              Base{Save: "Save"},
		Title:             "Title",
		Deleted:           "%d deleted",
		Body:              "<b>x</b>",
		Form:              Sub{Submit: "Go"},
		Ptr:               &Sub{Submit: "Ptr"},
		Help:              [][2]string{{"A", "a"}, {"B", "b"}},
		Options:           []Option{{"one", 1}},
		Format:            func(int) string { return "" },
		internal:          "not a text",
	}
}

func TestParseTag(t *testing.T) {
	sf, _ := reflect.TypeOf(msgs{}).FieldByName("Deleted")
	tag, err := ParseTag(sf.Tag)
	if err != nil {
		t.Fatal(err)
	}
	want := Tag{Type: Text, Hint: "Shown once records were deleted.",
		Fields: []Field{{"%d", "how many"}, {"%[2]s", "the model's name"}}}
	if !reflect.DeepEqual(tag, want) {
		t.Errorf("got %+v, want %+v", tag, want)
	}
	for _, bad := range []string{
		`i18n:"type=pdf"`, `i18n:"color='red'"`, `i18n:"label=x"`, `i18n:"fields='x'"`, `i18n:"label='x"`,
	} {
		if _, err := ParseTag(reflect.StructTag(bad)); err == nil {
			t.Errorf("%s was accepted", bad)
		}
	}
	if tag, err := ParseTag(""); err != nil || tag.Type != Text {
		t.Errorf("no tag: %+v %v", tag, err)
	}
}

func TestWalk(t *testing.T) {
	entries, skipped, err := Walk(newMsgs())
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, e := range entries {
		got = append(got, e.Path+"="+strings.TrimSpace(e.Value)+":"+string(e.Type))
	}
	want := []string{
		"ModuleDescription=The test messages.:text", "Save=Save:text", "Title=Title:text", "Deleted=%d deleted:text", "Body=<b>x</b>:html",
		"Form.Submit=Go:text", "Ptr.Submit=Ptr:text",
		"Help[0][0]=A:text", "Help[0][1]=a:text", "Help[1][0]=B:text", "Help[1][1]=b:text",
		"Options=- name: one\n  value: 1:form",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got\n%q\nwant\n%q", got, want)
	}
	var sk []string
	for _, s := range skipped {
		sk = append(sk, s.Path)
	}
	if !reflect.DeepEqual(sk, []string{"Count", "Format"}) {
		t.Errorf("skipped %q", sk)
	}
	if entries[1].Label() != "Save" || entries[2].Label() != "Page title" || entries[5].Label() != "Submit button" {
		t.Errorf("labels %q %q %q", entries[1].Label(), entries[2].Label(), entries[5].Label())
	}
	if _, _, err := Walk("x"); err == nil {
		t.Error("a string was walked")
	}
}

func TestGetSet(t *testing.T) {
	m := newMsgs()
	for path, want := range map[string]string{
		"Save": "Save", "Form.Submit": "Go", "Ptr.Submit": "Ptr", "Help[1][0]": "B", "Body": "<b>x</b>",
	} {
		if got, err := Get(m, path); err != nil || got != want {
			t.Errorf("Get(%s) = %q, %v", path, got, err)
		}
	}
	for path, v := range map[string]string{
		"Save": "Salvar", "Form.Submit": "Ir", "Help[1][0]": "Bê", "Help[3][1]": "new", "Body": "<i>y</i>",
		"Options": "- name: two\n  value: 2\n",
	} {
		if err := Set(m, path, v); err != nil {
			t.Fatalf("Set(%s): %v", path, err)
		}
	}
	if m.Save != "Salvar" || m.Form.Submit != "Ir" || m.Help[1][0] != "Bê" || len(m.Help) != 4 || m.Help[3][1] != "new" ||
		m.Body != "<i>y</i>" || !reflect.DeepEqual(m.Options, []Option{{"two", 2}}) {
		t.Errorf("after Set: %+v", m)
	}
	for _, bad := range []string{"Nope", "Count", "Help[x]", "Title.X", "Help[0", "", "internal", "Help[0][5]"} {
		if err := Set(m, bad, "v"); err == nil {
			t.Errorf("Set(%q) was accepted", bad)
		}
	}
	if _, err := Get(m, "Help[9][0]"); err == nil {
		t.Error("Get past the end was accepted")
	}
	if err := Set(*m, "Save", "x"); err == nil {
		t.Error("Set on a value was accepted")
	}
}

func TestClone(t *testing.T) {
	m := newMsgs()
	c := Clone(m)
	if err := Set(c, "Help[0][0]", "changed"); err != nil {
		t.Fatal(err)
	}
	Set(c, "Ptr.Submit", "changed")
	Set(c, "Save", "changed")
	Set(c, "Options", "[]")
	if m.Help[0][0] != "A" || m.Ptr.Submit != "Ptr" || m.Save != "Save" || len(m.Options) != 1 {
		t.Errorf("the original changed: %+v", m)
	}
	if c.Format == nil {
		t.Error("the func was lost")
	}
}

func TestDescription(t *testing.T) {
	if got := Description(newMsgs()); got != "The test messages." {
		t.Errorf("Description %q", got)
	}
	if got := Description(&Sub{}); got != "" {
		t.Errorf("no ModuleDescription: %q", got)
	}
	if got := Description("x"); got != "" {
		t.Errorf("a string: %q", got)
	}
}
