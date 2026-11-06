package vuetifyx

import (
	"context"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
)

// A truncated text: cut on one line within its width, its title the whole
// text; its click opens it (when cut) without reaching the row.
func TestTruncatedText(t *testing.T) {
	got := h.MustString(TruncatedText("Pagamento a fornecedor de materiais", "200px"), context.TODO())
	for _, want := range []string{"max-width: 200px", "title='Pagamento a fornecedor de materiais'", "textOverflow: &#39;ellipsis&#39;",
		"$event.stopPropagation()", "scrollWidth", "{open: false}"} {
		if !strings.Contains(got, want) && !strings.Contains(got, strings.ReplaceAll(want, "&#39;", "'")) {
			t.Errorf("no %q:\n%s", want, got)
		}
	}
	if got := h.MustString(TruncatedText("x", ""), context.TODO()); !strings.Contains(got, "max-width: 280px") {
		t.Errorf("the default width:\n%s", got)
	}
}

// A truncated cell: the room left in its row, its text cut at its width.
func TestTruncatedCell(t *testing.T) {
	got := h.MustString(TruncatedCell("Uma descrição longa"), context.TODO())
	for _, want := range []string{"<td", "width: 100%; max-width: 0", "max-width: 100%", "data-truncated-text"} {
		if !strings.Contains(got, want) {
			t.Errorf("no %q:\n%s", want, got)
		}
	}
}
