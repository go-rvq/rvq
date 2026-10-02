package schemaform

import (
	"net/url"
	"strings"
	"testing"
)

// A field declared `get name Type` is read-only: drawn so, and a post does not
// change it — KeepReadOnly puts the stored value back.
func TestGetterIsReadOnly(t *testing.T) {
	s, err := Parse(`[]{get path str; content text}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Fields) != 2 || s.Fields[0].Name != "path" || !s.Fields[0].ReadOnly || s.Fields[0].Type != DefaultType ||
		s.Fields[1].Name != "content" || s.Fields[1].ReadOnly || s.Fields[1].Type != "text" {
		for _, f := range s.Fields {
			t.Logf("%+v", f)
		}
		t.Fatal("not the fields of the schema")
	}

	posted := s.Decode(url.Values{
		"V[0].path": {"hacked.md"}, "V[0].content": {"new text"},
		"V[1].path": {"b.md"}, "V[1].content": {"b"},
	}, "V")
	stored := []any{
		map[string]any{"path": "a/README.md", "content": "old"},
		map[string]any{"path": "b.md", "content": "b"},
	}
	kept := s.KeepReadOnly(posted, stored).([]any)
	if p, _ := fieldOf(kept[0], "path"); p != "a/README.md" {
		t.Errorf("the post changed a read-only field: %v", p)
	}
	if c, _ := fieldOf(kept[0], "content"); c != "new text" {
		t.Errorf("the editable field was not taken: %v", c)
	}
}

// A getter is drawn read-only, the other fields editable.
func TestGetterComponent(t *testing.T) {
	got := render(t, New(), `{get path str; content text}`)
	i := strings.Index(got, `form[&#34;Value&#34;].path`)
	if i < 0 {
		i = strings.Index(got, ".path")
	}
	if i < 0 {
		t.Fatalf("no path field:\n%s", got)
	}
	start := strings.LastIndex(got[:i], "<v-text-field")
	if start < 0 || !strings.Contains(got[start:i+40], "readonly") {
		t.Errorf("the getter is not read-only:\n%s", got[max(0, start):min(len(got), i+80)])
	}
	if !strings.Contains(got, "<v-textarea") {
		t.Errorf("the content is not a textarea:\n%s", got)
	}
}
