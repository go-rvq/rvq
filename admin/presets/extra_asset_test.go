package presets

import (
	"strings"
	"testing"
)

// TestExtraAssetVersion: the address of an extra asset changes with its
// content (no browser keeps an old one).
func TestExtraAssetVersion(t *testing.T) {
	b := &Builder{}
	b.ExtraAsset("/x.css", "text/css", "a{}")
	ea := b.extraAssets[0]
	first := b.extraRef(ea)
	if !strings.HasPrefix(first, b.extraFullPath(ea)+"?v=") {
		t.Fatal(first)
	}
	b.ExtraAsset("/x.css", "text/css", "b{}")
	if again := b.extraRef(ea); again == first {
		t.Error("same address for another content", again)
	}
}
