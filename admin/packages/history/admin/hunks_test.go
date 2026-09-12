package admin

import (
	"strings"
	"testing"
)

// selectAll returns a selection covering every hunk between current and target.
func selectAll(current, target string, htmlMode bool) map[int]bool {
	sel := map[int]bool{}
	for _, hk := range fieldHunks(current, target, htmlMode) {
		sel[hk.Index] = true
	}
	return sel
}

func TestApplyHunksAllYieldsTarget(t *testing.T) {
	current := "a b c"
	target := "A b C"
	if got := applyHunks(current, target, selectAll(current, target, false), false); got != target {
		t.Fatalf("all selected = %q, want target %q", got, target)
	}
}

func TestApplyHunksNoneYieldsCurrent(t *testing.T) {
	current := "a b c"
	target := "A b C"
	if got := applyHunks(current, target, map[int]bool{}, false); got != current {
		t.Fatalf("none selected = %q, want current %q", got, current)
	}
}

func TestFieldHunksTwoRegions(t *testing.T) {
	current := "a b c"
	target := "A b C"
	hunks := fieldHunks(current, target, false)
	if len(hunks) != 2 {
		t.Fatalf("hunks = %d, want 2 (%+v)", len(hunks), hunks)
	}
	// Selecting only the first hunk reverts just the "a"->"A" region.
	if got := applyHunks(current, target, map[int]bool{0: true}, false); got != "A b c" {
		t.Fatalf("select {0} = %q, want %q", got, "A b c")
	}
	// Selecting only the second reverts just the "c"->"C" region.
	if got := applyHunks(current, target, map[int]bool{1: true}, false); got != "a b C" {
		t.Fatalf("select {1} = %q, want %q", got, "a b C")
	}
}

func TestFieldHunksIdenticalHasNone(t *testing.T) {
	if hunks := fieldHunks("same", "same", false); len(hunks) != 0 {
		t.Fatalf("identical values should have no hunks, got %d", len(hunks))
	}
}

func TestApplyHunksPureInsertAndDelete(t *testing.T) {
	// current has an extra word; target adds another — one delete, one insert.
	current := "keep remove"
	target := "keep add"
	all := applyHunks(current, target, selectAll(current, target, false), false)
	if all != target {
		t.Fatalf("all selected = %q, want %q", all, target)
	}
	none := applyHunks(current, target, map[int]bool{}, false)
	if none != current {
		t.Fatalf("none selected = %q, want %q", none, current)
	}
}

// In HTML mode a hunk must never split a tag: the change is expressed with whole
// tags/words, so a hunk's text is renderable HTML rather than a `</p><p>` shard.
func TestFieldHunksHTMLKeepsTagsWhole(t *testing.T) {
	current := "<p>hello world</p>"
	target := "<p>hello there</p>"
	hunks := fieldHunks(current, target, true)
	if len(hunks) != 1 {
		t.Fatalf("hunks = %d, want 1 (%+v)", len(hunks), hunks)
	}
	// Only the changed word is in the hunk — no tag fragments.
	if strings.Contains(hunks[0].Old, "<") || strings.Contains(hunks[0].New, "<") {
		t.Fatalf("hunk contains tag fragments: old=%q new=%q", hunks[0].Old, hunks[0].New)
	}
	if strings.TrimSpace(hunks[0].Old) != "world" || strings.TrimSpace(hunks[0].New) != "there" {
		t.Fatalf("hunk = old %q new %q, want old %q new %q", hunks[0].Old, hunks[0].New, "world", "there")
	}
}

func TestApplyHunksHTMLRoundTrip(t *testing.T) {
	current := "<p>a</p><p>keep</p>"
	target := "<p>B</p><p>keep</p>"
	if got := applyHunks(current, target, selectAll(current, target, true), true); got != target {
		t.Fatalf("all selected = %q, want target %q", got, target)
	}
	if got := applyHunks(current, target, map[int]bool{}, true); got != current {
		t.Fatalf("none selected = %q, want current %q", got, current)
	}
}
