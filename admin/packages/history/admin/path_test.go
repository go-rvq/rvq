package admin

import (
	"encoding/json"
	"testing"
)

func TestFieldValuePaths(t *testing.T) {
	m := map[string]json.RawMessage{
		"Title":       json.RawMessage(`"Hello"`),
		"PageOptions": json.RawMessage(`{"Layout":"grid","Order":3}`),
		"Tags":        json.RawMessage(`[{"Name":"a"},{"Name":"b"}]`),
	}
	cases := map[string]string{
		"Title":              "Hello",
		"PageOptions.Layout": "grid",
		"PageOptions.Order":  "3",
		"Tags[0].Name":       "a",
		"Tags[1].Name":       "b",
		"Tags[2].Name":       "", // out of range
		"PageOptions.Nope":   "", // missing
		"Nope":               "", // missing top
	}
	for path, want := range cases {
		if got := fieldValue(m, path); got != want {
			t.Errorf("fieldValue(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestParsePath(t *testing.T) {
	got := parsePath("a[0].b.c[1]")
	want := []pathToken{
		{key: "a"}, {idx: 0, isIdx: true}, {key: "b"}, {key: "c"}, {idx: 1, isIdx: true},
	}
	if len(got) != len(want) {
		t.Fatalf("parsePath len = %d, want %d (%+v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("token %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}
