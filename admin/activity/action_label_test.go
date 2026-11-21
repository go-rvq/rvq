package activity

import (
	"context"
	"testing"
)

// An action of another package is shown by its label registered, in the
// language of the ctx; one not registered, by its name; this package's, by
// its messages.
func TestActionLabel(t *testing.T) {
	type langKey struct{}
	RegisterActionLabel("TestBlock", func(ctx context.Context) string {
		if ctx.Value(langKey{}) == "pt-BR" {
			return "Bloqueio"
		}
		return "Block"
	})
	ctx := context.Background()
	if got := ActionLabel(context.WithValue(ctx, langKey{}, "pt-BR"), "TestBlock"); got != "Bloqueio" {
		t.Errorf("registered, pt-BR: %q", got)
	}
	if got := ActionLabel(ctx, "TestBlock"); got != "Block" {
		t.Errorf("registered: %q", got)
	}
	if got := ActionLabel(ctx, "TestNone"); got != "TestNone" {
		t.Errorf("not registered: %q", got)
	}
	if got := ActionLabel(ctx, ActivityCreate); got != Messages_en_US.ActionCreate {
		t.Errorf("this package's: %q", got)
	}
}
