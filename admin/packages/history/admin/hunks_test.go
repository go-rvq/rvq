package admin

import "testing"

// selectAll returns a selection covering every hunk between current and target.
func selectAll(current, target string) map[int]bool {
	sel := map[int]bool{}
	for _, hk := range fieldHunks(current, target) {
		sel[hk.Index] = true
	}
	return sel
}

func TestApplyHunksAllYieldsTarget(t *testing.T) {
	current := "a b c"
	target := "A b C"
	if got := applyHunks(current, target, selectAll(current, target)); got != target {
		t.Fatalf("all selected = %q, want target %q", got, target)
	}
}

func TestApplyHunksNoneYieldsCurrent(t *testing.T) {
	current := "a b c"
	target := "A b C"
	if got := applyHunks(current, target, map[int]bool{}); got != current {
		t.Fatalf("none selected = %q, want current %q", got, current)
	}
}

func TestFieldHunksTwoRegions(t *testing.T) {
	current := "a b c"
	target := "A b C"
	hunks := fieldHunks(current, target)
	if len(hunks) != 2 {
		t.Fatalf("hunks = %d, want 2 (%+v)", len(hunks), hunks)
	}
	// Selecting only the first hunk reverts just the "a"->"A" region.
	if got := applyHunks(current, target, map[int]bool{0: true}); got != "A b c" {
		t.Fatalf("select {0} = %q, want %q", got, "A b c")
	}
	// Selecting only the second reverts just the "c"->"C" region.
	if got := applyHunks(current, target, map[int]bool{1: true}); got != "a b C" {
		t.Fatalf("select {1} = %q, want %q", got, "a b C")
	}
}

func TestFieldHunksIdenticalHasNone(t *testing.T) {
	if hunks := fieldHunks("same", "same"); len(hunks) != 0 {
		t.Fatalf("identical values should have no hunks, got %d", len(hunks))
	}
}

func TestApplyHunksPureInsertAndDelete(t *testing.T) {
	// current has an extra word; target adds another — one delete, one insert.
	current := "keep remove"
	target := "keep add"
	all := applyHunks(current, target, selectAll(current, target))
	if all != target {
		t.Fatalf("all selected = %q, want %q", all, target)
	}
	none := applyHunks(current, target, map[int]bool{})
	if none != current {
		t.Fatalf("none selected = %q, want %q", none, current)
	}
}
