package admin

import (
	"encoding/json"
	"testing"
)

func TestReferenceIDFromObjectAndScalar(t *testing.T) {
	// belongs-to snapshot: the embedded record — id comes from its ID field.
	m := map[string]json.RawMessage{"Author": json.RawMessage(`{"ID":42,"Name":"Ana"}`)}
	if got := referenceID(m, "Author"); got != "42" {
		t.Fatalf("referenceID(object) = %q, want 42", got)
	}
	// foreign-key snapshot: the id itself.
	m2 := map[string]json.RawMessage{"AuthorID": json.RawMessage(`7`)}
	if got := referenceID(m2, "AuthorID"); got != "7" {
		t.Fatalf("referenceID(scalar) = %q, want 7", got)
	}
}

func TestRefElemLabel(t *testing.T) {
	id, label := refElemLabel(json.RawMessage(`{"ID":3,"Name":"Go"}`))
	if id != "3" || label != "Go" {
		t.Fatalf("refElemLabel = (%q,%q), want (3,Go)", id, label)
	}
	// no title field → label falls back to #id.
	id, label = refElemLabel(json.RawMessage(`{"ID":9}`))
	if id != "9" || label != "#9" {
		t.Fatalf("refElemLabel(no title) = (%q,%q), want (9,#9)", id, label)
	}
}

func TestRefItemsDecode(t *testing.T) {
	items := refItems(json.RawMessage(`[{"ID":1,"Name":"Go"},{"ID":2,"Name":"Rust"}]`))
	if len(items) != 2 || items[0].label != "Go" || items[0].id != "1" || items[1].label != "Rust" {
		t.Fatalf("refItems = %+v, want Go#1, Rust#2", items)
	}
	if refItems(json.RawMessage(`{"a":1}`)) != nil {
		t.Fatal("refItems of a non-array should be nil")
	}
}

// The to-many diff is over the whole list, by id: it names which item was
// removed (in OLD, not NEW) and which was added (in NEW, not OLD).
func TestManyRefChangesAddedRemoved(t *testing.T) {
	old := json.RawMessage(`[{"ID":1,"Name":"Go"},{"ID":2,"Name":"Rust"}]`)
	neu := json.RawMessage(`[{"ID":1,"Name":"Go"},{"ID":3,"Name":"Java"}]`)
	_, _, removed, added := manyRefChanges(old, neu)
	if !removed["2"] || len(removed) != 1 {
		t.Fatalf("removed = %v, want {2} (Rust)", removed)
	}
	if !added["3"] || len(added) != 1 {
		t.Fatalf("added = %v, want {3} (Java)", added)
	}
	// The renderer produces both sides.
	oldC, newC, _, handled := manyReferenceDiff(old, neu)
	if !handled || oldC == nil || newC == nil {
		t.Fatalf("manyReferenceDiff did not produce both sides")
	}
}

// lineDiffSides drives the Prism handler's per-line marks: which lines are
// removed on OLD and added on NEW.
func TestLineDiffSides(t *testing.T) {
	old := "a\nb\nc"
	neu := "a\nB\nc"
	removed, added := lineDiffSides(old, neu)
	if len(removed) != 1 || removed[0] != 2 {
		t.Fatalf("removed = %v, want [2]", removed)
	}
	if len(added) != 1 || added[0] != 2 {
		t.Fatalf("added = %v, want [2]", added)
	}
	// A pure insertion at the end marks only the new line.
	removed, added = lineDiffSides("a\nb", "a\nb\nc")
	if len(removed) != 0 {
		t.Fatalf("removed = %v, want none", removed)
	}
	if len(added) != 1 || added[0] != 3 {
		t.Fatalf("added = %v, want [3]", added)
	}
}

func TestIsJSONArray(t *testing.T) {
	if !isJSONArray(json.RawMessage(`  [1,2]`)) {
		t.Fatal("array not detected")
	}
	if isJSONArray(json.RawMessage(`{"a":1}`)) {
		t.Fatal("object misdetected as array")
	}
}
