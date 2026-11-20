package login

import (
	"context"
	"testing"

	"github.com/go-rvq/rvq/x/lifecycle"
)

// The renewal as a node: scheduled by the rule given when it starts, stopped
// when it closes; a rule that does not parse logged, the command going on.
func TestSecureKeyRenewalNode(t *testing.T) {
	b := newSecureKeyBuilder(t)
	asked := 0
	n := b.SecureKeyRenewalNode(func() string { asked++; return "@daily" })
	if n.Key != SecureKeyKey || asked != 0 {
		t.Fatalf("the node: %+v; the rule asked: %d", n, asked)
	}
	var g lifecycle.Graph
	g.Add(n)
	err := g.Run(context.Background(), func(context.Context) error {
		b.secureKeyCronMu.Lock()
		defer b.secureKeyCronMu.Unlock()
		if b.secureKeyCron == nil {
			t.Error("not scheduled while running")
		}
		return nil
	})
	if err != nil || asked != 1 {
		t.Fatalf("run: %v; the rule asked: %d", err, asked)
	}
	if b.secureKeyCron != nil {
		t.Error("still scheduled after the node closed")
	}

	var bad lifecycle.Graph
	bad.Add(b.SecureKeyRenewalNode(func() string { return "not a rule" }))
	if err := bad.Run(context.Background(), nil); err != nil {
		t.Errorf("a rule that does not parse stopped the command: %v", err)
	}
}
