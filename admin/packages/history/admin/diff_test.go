package admin

import (
	"encoding/json"
	"testing"

	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
)

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
