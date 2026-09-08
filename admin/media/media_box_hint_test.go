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
	got := originalSizeHint(ctx, cover)
	for _, want := range []string{"1800×771", "≈7:3"} {
		if !strings.Contains(got, want) {
			t.Errorf("hint = %q, want it to name %q", got, want)
		}
	}
	// Only the largest: the smaller crops are covered by it.
	for _, unwanted := range []string{"900", "500"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("hint = %q names the smaller size %q", got, unwanted)
		}
	}

	// The largest is the one with the most pixels, not the widest.
	mixed := &media_library.MediaBoxConfig{Sizes: map[string]*base.Size{
		"wide": {Width: 1600, Height: 400},
		"tall": {Width: 900, Height: 1200},
	}}
	if got := originalSizeHint(ctx, mixed); !strings.Contains(got, "900×1200") {
		t.Errorf("hint = %q, want it to name 900×1200", got)
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

func TestAspectRatio(t *testing.T) {
	for _, c := range []struct {
		w, h int
		want string
	}{
		{1920, 1080, "16:9"},
		{800, 600, "4:3"},
		{1000, 1000, "1:1"},
		{1500, 1000, "3:2"},
		// Reduces to 600:257, which says nothing: the nearest small ratio is
		// given instead, and it is the proportion sold as 21:9.
		{1800, 771, "≈7:3"},
		{900, 1200, "3:4"},
	} {
		if got := aspectRatio(c.w, c.h); got != c.want {
			t.Errorf("aspectRatio(%d, %d) = %q, want %q", c.w, c.h, got, c.want)
		}
	}
}
