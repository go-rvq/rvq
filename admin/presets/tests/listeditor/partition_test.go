package listeditor

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	. "github.com/go-rvq/rvq/admin/presets/integration"
)

// labels extracts each item's Label for order-sensitive assertions.
func labels(vs []reflect.Value) []string {
	out := make([]string, 0, len(vs))
	for _, v := range vs {
		out = append(out, reflect.Indirect(v).FieldByName("Label").String())
	}
	return out
}

func eq(t *testing.T, name string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s: got %v, want %v", name, got, want)
		return
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("%s: got %v, want %v", name, got, want)
			return
		}
	}
}

// TestPartitionListEditorItems exercises the helper directly: classification is
// by the __deleted/__new flags only (never the primary key), and every partition
// is returned in __pos order. Item B has a zero PK yet is NOT flagged __new, so
// it must land in Others (not New) — proving classification ignores the ID.
func TestPartitionListEditorItems(t *testing.T) {
	items := []*Item{
		{ID: 1, Label: "A"}, // index 0, __pos 2 -> existing (Others)
		{ID: 0, Label: "B"}, // index 1, __pos 0 -> zero PK but NOT __new (Others)
		{ID: 2, Label: "C"}, // index 2, __pos 3 -> __deleted
		{ID: 0, Label: "D"}, // index 3, __pos 1 -> __new
	}
	vals := url.Values{
		"Items.__present":    {"1"}, // the list editor was initialized
		"Items[0].__pos":     {"2"},
		"Items[1].__pos":     {"0"},
		"Items[2].__pos":     {"3"},
		"Items[2].__deleted": {"true"},
		"Items[3].__pos":     {"1"},
		"Items[3].__new":     {"true"},
	}

	p := presets.PartitionListEditorItems(vals, "Items", reflect.ValueOf(items))
	if p == nil {
		t.Fatal("expected a partition, got nil (present marker set)")
	}

	// order by __pos: B(0), D(1), A(2), C(3)
	eq(t, "Deleted", labels(p.Deleted), []string{"C"})
	eq(t, "New", labels(p.New), []string{"D"})
	eq(t, "Others", labels(p.Others), []string{"B", "A"})  // zero-PK B included, __pos order
	eq(t, "Kept", labels(p.Kept), []string{"B", "D", "A"}) // New+Others, __pos order

	// KeptSlice returns a typed []*Item in the same order.
	ks := p.KeptSlice(reflect.TypeOf(items)).Interface().([]*Item)
	eq(t, "KeptSlice", labels(reflectValues(ks)), []string{"B", "D", "A"})
}

// TestPartitionListEditorItems_NotInitialized: without the presence marker the
// list editor was not initialized, so the helper returns nil (the field was not
// submitted — leave its children untouched), even if item keys are present.
func TestPartitionListEditorItems_NotInitialized(t *testing.T) {
	items := []*Item{{ID: 1, Label: "A"}}
	vals := url.Values{
		"Items[0].ID":    {"1"},
		"Items[0].Label": {"A"},
		// no "Items.__present"
	}
	if p := presets.PartitionListEditorItems(vals, "Items", reflect.ValueOf(items)); p != nil {
		t.Fatalf("expected nil (not initialized), got %+v", p)
	}
}

func reflectValues(items []*Item) []reflect.Value {
	out := make([]reflect.Value, len(items))
	for i := range items {
		out[i] = reflect.ValueOf(items[i])
	}
	return out
}
