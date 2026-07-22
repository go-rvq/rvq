package web

import (
	"net/url"
	"testing"
)

type reorderItem struct {
	ID    uint
	Valor string
}
type reorderDoc struct {
	Parcelas []*reorderItem
}

// TestReorderSlicesByPos: the decoded slice follows __pos (not the array index),
// items without __pos go to the end, and the form keys are reindexed to match.
func TestReorderSlicesByPos(t *testing.T) {
	// submitted order [0,1,2] but __pos says the wanted order is 1,2,0;
	// a 4th item has no __pos → must land at the end.
	values := url.Values{
		"Parcelas[0].ID":        {"10"},
		"Parcelas[0].Valor":     {"a"},
		"Parcelas[0].__pos":     {"2"},
		"Parcelas[1].ID":        {"11"},
		"Parcelas[1].Valor":     {"b"},
		"Parcelas[1].__pos":     {"0"},
		"Parcelas[1].__deleted": {"true"},
		"Parcelas[2].ID":        {"12"},
		"Parcelas[2].Valor":     {"c"},
		"Parcelas[2].__pos":     {"1"},
		"Parcelas[3].ID":        {"13"},
		"Parcelas[3].Valor":     {"d"}, // no __pos → end
	}
	var d reorderDoc
	dec := (&EventContext{}).UnmarshalFormValues(values, &d)
	if dec != nil {
		t.Fatal(dec)
	}

	// expected order by __pos: b(11,pos0), c(12,pos1), a(10,pos2), then d(13,no pos)
	want := []struct {
		id  uint
		val string
	}{{11, "b"}, {12, "c"}, {10, "a"}, {13, "d"}}
	if len(d.Parcelas) != len(want) {
		t.Fatalf("len = %d, want %d", len(d.Parcelas), len(want))
	}
	for i, w := range want {
		if d.Parcelas[i].ID != w.id || d.Parcelas[i].Valor != w.val {
			t.Errorf("[%d] = {ID:%d Valor:%q}, want {ID:%d Valor:%q}", i, d.Parcelas[i].ID, d.Parcelas[i].Valor, w.id, w.val)
		}
	}

	// the deleted item (ID 11) moved to position 0; its __deleted key must have
	// been reindexed to Parcelas[0].__deleted.
	if values.Get("Parcelas[0].__deleted") != "true" {
		t.Errorf("reindex: expected Parcelas[0].__deleted=true, got form=%v", values)
	}
	if values.Get("Parcelas[1].__deleted") != "" {
		t.Errorf("reindex: stale Parcelas[1].__deleted should be gone")
	}
}

// TestReorderNoPosNoop: without any __pos, the order is unchanged.
func TestReorderNoPosNoop(t *testing.T) {
	values := url.Values{
		"Parcelas[0].ID": {"1"},
		"Parcelas[1].ID": {"2"},
	}
	var d reorderDoc
	_ = (&EventContext{}).UnmarshalFormValues(values, &d)
	if len(d.Parcelas) != 2 || d.Parcelas[0].ID != 1 || d.Parcelas[1].ID != 2 {
		t.Fatalf("order changed unexpectedly: %+v", d.Parcelas)
	}
}
