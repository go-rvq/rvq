package presets

import (
	"context"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
)

// A message with a body shows it in itself, below its text; a persistent one
// stays until it is closed (no timeout), the others go after their
// duration.
func TestRenderFlashBodyPersistent(t *testing.T) {
	out := h.MustString(RenderFlash(&FlashMessage{Text: "Made.", Persistent: true,
		Body: h.Code("the-code").Attr("data-code", true)}), context.Background())
	for _, want := range []string{"Made.", "data-flash-body", `data-code`, "the-code", `:timeout='-1'`} {
		if !strings.Contains(out, want) {
			t.Errorf("no %s in:\n%s", want, out)
		}
	}
	if strings.Index(out, "Made.") > strings.Index(out, "the-code") {
		t.Error("the body before the text")
	}
	plain := h.MustString(RenderFlash(&FlashMessage{Text: "Saved."}), context.Background())
	if strings.Contains(plain, "data-flash-body") || !strings.Contains(plain, `:timeout='5000'`) {
		t.Errorf("a plain message:\n%s", plain)
	}
}
