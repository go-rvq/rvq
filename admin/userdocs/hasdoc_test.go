package userdocs

import (
	"testing"
	"testing/fstest"
)

// HasDoc: a node whose document some source has (its own, or the general
// one of an action); one none has.
func TestHasDoc(t *testing.T) {
	b := &Builder{sources: []*Source{
		{Package: "a", FS: fstest.MapFS{"en/dashboard/widgets/w1/README.md": {Data: []byte("# W1")}}},
		{Package: "b", FS: fstest.MapFS{"pt-BR/actions/restore.md": {Data: []byte("# Restore")}}},
	}}
	for node, want := range map[string]bool{
		"dashboard/widgets/w1":      true,
		"dashboard/widgets/w2":      false,
		"posts/actions/restore":     true,
		"posts/actions/nonexisting": false,
	} {
		if got := b.HasDoc(node); got != want {
			t.Errorf("HasDoc(%q) = %v, want %v", node, got, want)
		}
	}
}
