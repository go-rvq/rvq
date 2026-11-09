package filecards

import (
	"context"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
)

// An image a thumbnail in the lightbox (vuetifyx's); another file, its icon, opened; each
// its attribute and its download; none, nothing.
func TestCards(t *testing.T) {
	if Cards(nil) != nil {
		t.Fatal("no files, no cards")
	}
	out := h.MustString(Cards([]File{
		{Key: "a", Name: "foto.png", Size: 2048, Type: "image/png", URL: "/f/a"},
		{Key: "b", Name: "pauta.pdf", Size: 10, Type: "application/pdf", URL: "/f/b"},
		{Key: "c", Name: "x.svg", Type: "image/svg+xml", URL: "/f/c"},
	}, "data-answer-file"), context.Background())
	for _, want := range []string{"data-answer-file='a'", "data-answer-file='b'", "mdi-file-pdf-box", "2 KB", "data-lightbox-src='/f/a'", "data-file-download"} {
		if !strings.Contains(out, want) {
			t.Errorf("no %q", want)
		}
	}
	// an SVG is no image of the lightbox (it may run): its icon
	if strings.Contains(out, "data-lightbox-src='/f/c'") {
		t.Error("the SVG in the lightbox")
	}
	if Icon("planilha.xlsx", "") != "mdi-file-excel-outline" || Icon("a", "") != "mdi-file-outline" {
		t.Error("Icon")
	}
}
