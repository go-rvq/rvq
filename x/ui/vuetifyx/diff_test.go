package vuetifyx

import (
	"reflect"
	"testing"
)

// The lines both texts have are paired, every one of them: a block changed
// around a line kept leaves that line kept (the half-match shortcut of
// diffmatchpatch, with a timeout, took it into the change).
func TestLineChangesPairsEveryCommonLine(t *testing.T) {
	old := "@main\n    p one\n    p two\n    p three\n    p four\n    p gone\n    p keep\n    p tail\n"
	cur := "@main\n    p one\n    p 2-3-4 together\n    p keep\n    p inserted\n    p tail\n"
	removed, added, _, _ := LineChanges(old, cur)
	if want := []int{3, 4, 5, 6}; !reflect.DeepEqual(removed, want) {
		t.Errorf("removed %v, want %v", removed, want)
	}
	if want := []int{3, 5}; !reflect.DeepEqual(added, want) {
		t.Errorf("added %v, want %v", added, want)
	}
}
