package presets

import (
	"context"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
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

// TestLazyExtraAsset: a lazy extra asset is served (its address, its
// version) and not in the pages' head; ExtraAssetRef finds it.
func TestLazyExtraAsset(t *testing.T) {
	b := &Builder{prefix: "/admin"}
	b.ExtraAsset("/a.js", "text/javascript", "a()")
	b.LazyExtraAsset("/big.js", "text/javascript", "big()")
	ref := b.ExtraAssetRef("big.js")
	if !strings.HasPrefix(ref, "/admin/!assets/extra/big.js?v=") {
		t.Fatal(ref)
	}
	if b.ExtraAssetRef("/none.js") != "" {
		t.Error("an asset not registered has no address")
	}
	ctx := &web.EventContext{Injector: &web.PageInjector{}}
	b.InjectExtraAssets(ctx)
	head := h.MustString(ctx.Injector.GetHeadHTMLComponent(), context.Background())
	if !strings.Contains(head, "/extra/a.js") || strings.Contains(head, "big.js") {
		t.Errorf("the head: %s", head)
	}
}
