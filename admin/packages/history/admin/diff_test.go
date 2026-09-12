package admin

import (
	"encoding/json"
	"testing"

	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
)

func TestFormatChangedField(t *testing.T) {
	mustJSON := func(v any) json.RawMessage {
		b, _ := json.Marshal(v)
		return b
	}

	// nil mb/ctx → raw names (no i18n in the test).

	// Leaf change: just the name.
	if got := formatChangedField(nil, "Title", mustJSON("a"), mustJSON("b"), nil); got != "Title" {
		t.Fatalf("leaf = %q, want Title", got)
	}

	// Structured change: only the changed sub-keys, sorted, in brackets.
	old := mustJSON(map[string]any{"Layout": "x", "Config": 1, "Galleries": []int{1}})
	neu := mustJSON(map[string]any{"Layout": "y", "Config": 1, "Galleries": []int{1, 2}})
	got := formatChangedField(nil, "PageOptions", old, neu, nil)
	want := "PageOptions [ Galleries, Layout ]"
	if got != want {
		t.Fatalf("structured = %q, want %q", got, want)
	}

	// Nested structured change recurses.
	old2 := mustJSON(map[string]any{"Seo": map[string]any{"Title": "a", "Desc": "d"}})
	neu2 := mustJSON(map[string]any{"Seo": map[string]any{"Title": "b", "Desc": "d"}})
	got2 := formatChangedField(nil, "Meta", old2, neu2, nil)
	want2 := "Meta [ Seo [ Title ] ]"
	if got2 != want2 {
		t.Fatalf("nested = %q, want %q", got2, want2)
	}
}

func rawMap(t *testing.T, m map[string]any) map[string]json.RawMessage {
	t.Helper()
	out := map[string]json.RawMessage{}
	for k, v := range m {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("marshal %s: %v", k, err)
		}
		out[k] = b
	}
	return out
}

func TestFieldValueUnquotesString(t *testing.T) {
	m := rawMap(t, map[string]any{"Body": "<p>hi</p>"})
	if got, want := fieldValue(m, "Body"), "<p>hi</p>"; got != want {
		t.Fatalf("fieldValue string = %q, want %q", got, want)
	}
}

func TestFieldValueRawForNonString(t *testing.T) {
	m := rawMap(t, map[string]any{"Position": 42})
	if got, want := fieldValue(m, "Position"), "42"; got != want {
		t.Fatalf("fieldValue number = %q, want %q", got, want)
	}
}

func TestFieldValueMissing(t *testing.T) {
	if got := fieldValue(map[string]json.RawMessage{}, "Nope"); got != "" {
		t.Fatalf("missing field = %q, want empty", got)
	}
}

func TestFieldMapEmptyRevision(t *testing.T) {
	m, err := fieldMap(&histmodels.Revision{})
	if err != nil {
		t.Fatalf("fieldMap: %v", err)
	}
	if len(m) != 0 {
		t.Fatalf("empty revision fieldMap len = %d, want 0", len(m))
	}
}
