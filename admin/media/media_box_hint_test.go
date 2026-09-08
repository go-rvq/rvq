package media

import (
	"context"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/web"
)

func TestOriginalSizeHint(t *testing.T) {
	// A context with no messages registered falls back to Messages_en_US,
	// which is what this asserts against.
	ctx := context.Background()

	// The sizes hermon configures for a cover.
	cover := &media_library.MediaBoxConfig{Sizes: map[string]*base.Size{
		"site.thumb":  {Width: 500, Height: 214},
		"site.middle": {Width: 900, Height: 385},
		"site.full":   {Width: 1800, Height: 771},
	}}
	if got := originalSizeHint(ctx, cover); !strings.Contains(got, "1800×771") {
		t.Errorf("hint = %q, want it to name 1800×771", got)
	}

	// Widest and tallest in different entries: the original has to cover both.
	mixed := &media_library.MediaBoxConfig{Sizes: map[string]*base.Size{
		"wide": {Width: 1600, Height: 400},
		"tall": {Width: 600, Height: 1200},
	}}
	if got := originalSizeHint(ctx, mixed); !strings.Contains(got, "1600×1200") {
		t.Errorf("hint = %q, want it to name 1600×1200", got)
	}

	for name, cfg := range map[string]*media_library.MediaBoxConfig{
		"nil config":  nil,
		"no sizes":    {},
		"empty sizes": {Sizes: map[string]*base.Size{}},
		"zeroed size": {Sizes: map[string]*base.Size{"x": {}}},
	} {
		if got := originalSizeHint(ctx, cfg); got != "" {
			t.Errorf("%s: hint = %q, want none", name, got)
		}
	}
}

// The hint has to reach the rendered field: the media box used to drop it,
// passing only the label.
func TestMediaBoxRendersTheHint(t *testing.T) {
	// Write reaches for the EventContext to build its portals.
	ctx := web.ContextWithEventContext(context.Background(),
		&web.EventContext{R: httptest.NewRequest(http.MethodGet, "/", nil)})

	out, err := h.MarshallString(
		QMediaBox(nil).
			FieldName("Cover").
			Value(&media_library.MediaBox{}).
			Label("Imagem de Destaque").
			Hint("Original recomendada: pelo menos 1800×771 px.").
			Config(&media_library.MediaBoxConfig{}),
		ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out, "Imagem de Destaque") {
		t.Error("the label is missing")
	}
	if !strings.Contains(out, "Original recomendada: pelo menos 1800×771 px.") {
		t.Errorf("the hint is missing from:\n%s", out)
	}
}
