package web

import (
	"context"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
)

// A scope whose locals are renamed ({ locals: drop }) gives no "locals" of
// its own: an outer scope's locals stay seen inside it.
func TestScopeRenamedSlot(t *testing.T) {
	out := h.MustString(Scope(h.Div()).Slot("{ locals: drop }").LocalsInit("{over: false}"), context.Background())
	if !strings.Contains(out, "v-slot='{ locals: drop }'") {
		t.Fatalf("the slot: %s", out)
	}
	out = h.MustString(Scope(h.Div()).LocalsInit("{a: 1}").Slot("{ form }"), context.Background())
	if !strings.Contains(out, "v-slot='{ form, locals }'") {
		t.Fatalf("not renamed: %s", out)
	}
}
