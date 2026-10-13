package vuetifyx

import (
	"context"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
)

// A resizable drawer says so, with the key its width is kept under, and its
// limits; with no key, resizable alone.
func TestNavigationDrawerResizable(t *testing.T) {
	out, err := h.Marshal(VXNavigationDrawer().Resizable("menu.width").MinWidth(220).MaxWidth(600), context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`resizable='menu.width'`, `:min-width='220'`, `:max-width='600'`} {
		if !strings.Contains(string(out), want) {
			t.Errorf("no %s in %s", want, out)
		}
	}
	out, _ = h.Marshal(VXNavigationDrawer().Resizable(""), context.Background())
	if !strings.Contains(string(out), "resizable") || strings.Contains(string(out), "resizable='") {
		t.Errorf("resizable with no key: %s", out)
	}
}
