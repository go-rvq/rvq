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

// The description for accessibility of a media box, unless its config
// disables it (DisableDescription).
func TestMediaBoxDisableDescription(t *testing.T) {
	ectx := &web.EventContext{R: httptest.NewRequest(http.MethodGet, "/", nil)}
	ctx := web.ContextWithEventContext(context.Background(), ectx)
	ectx.R = ectx.R.WithContext(ctx)
	box := &media_library.MediaBox{ID: "1", Url: "/system/1/file.png", FileName: "file.png"}
	render := func(cfg *media_library.MediaBoxConfig) string {
		out, err := h.MarshallString(mediaBoxThumbnails(ectx, box, "Foto", cfg, false, false), ctx)
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	if out := render(&media_library.MediaBoxConfig{}); !strings.Contains(out, `form["Foto.Description"]`) {
		t.Error("the description is not offered")
	}
	if out := render(&media_library.MediaBoxConfig{DisableDescription: true}); strings.Contains(out, `form["Foto.Description"]`) {
		t.Error("the description is offered, disabled")
	}
}

// What a media box holds: an image's dimensions and the file's size (the
// original's for the file itself), a file's size only.
func TestMediaBoxInfo(t *testing.T) {
	img := &media_library.MediaBox{FileName: "a.png", Width: 800, Height: 600, FileSizes: map[string]int{"original": 120500}}
	if got := MediaBoxInfo(img, ""); got != "800 × 600 · 120kB" {
		t.Errorf("image: %q", got)
	}
	pdf := &media_library.MediaBox{FileName: "a.pdf", FileSizes: map[string]int{"original": 2048}}
	if got := MediaBoxInfo(pdf, ""); got != "2kB" {
		t.Errorf("file: %q", got)
	}
	if got := MediaBoxInfo(&media_library.MediaBox{FileName: "b.txt"}, ""); got != "" {
		t.Errorf("unknown: %q", got)
	}
}

// A box of a list of images (Card) is drawn as a card in every event —
// the image in the lightbox of the list, its description beside it; none,
// only the button that chooses it.
func TestMediaBoxCard(t *testing.T) {
	ectx := &web.EventContext{R: httptest.NewRequest(http.MethodGet, "/", nil)}
	ctx := web.ContextWithEventContext(context.Background(), ectx)
	ectx.R = ectx.R.WithContext(ctx)
	cfg := &media_library.MediaBoxConfig{Card: true, DescriptionField: "Legenda"}
	render := func(box *media_library.MediaBox) string {
		out, err := h.MarshallString(mediaBoxThumbnails(ectx, box, "Itens[1].Imagens[2].Midia", cfg, false, false), ctx)
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	out := render(&media_library.MediaBox{ID: "1", Url: "/system/1/file.png", FileName: "file.png", Width: 80, Height: 60})
	for _, want := range []string{`data-image-card`, `data-lightbox='Itens[1].Imagens'`, `data-lightbox-src='/system/1/file.png'`,
		`form["Itens[1].Imagens[2].Legenda"]`, `form["Itens[1].Imagens[2].Midia.Values"]`, `80 × 60`} {
		if !strings.Contains(out, want) {
			t.Errorf("no %s in %s", want, out)
		}
	}
	if strings.Contains(out, "mdi-upload-outline") {
		t.Error("the box's own buttons in a card")
	}
	if out := render(&media_library.MediaBox{}); !strings.Contains(out, "data-image-card-empty") || strings.Contains(out, "data-lightbox") {
		t.Errorf("an empty card: %s", out)
	}
}

// An empty media box is a card to drop a file on — uploaded and chosen at
// once (UploadMediaBoxesEvent, then chooseFileEvent) —, its click the chooser;
// one with a file, its buttons.
func TestMediaBoxEmptyIsADrop(t *testing.T) {
	ectx := &web.EventContext{R: httptest.NewRequest(http.MethodGet, "/", nil)}
	ctx := web.ContextWithEventContext(context.Background(), ectx)
	ectx.R = ectx.R.WithContext(ctx)
	render := func(box *media_library.MediaBox, cfg *media_library.MediaBoxConfig) string {
		out, err := h.MarshallString(mediaBoxThumbnails(ectx, box, "Foto", cfg, false, false), ctx)
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	out := render(&media_library.MediaBox{}, &media_library.MediaBoxConfig{AllowType: "image"})
	for _, want := range []string{`data-media-box-drop='Foto'`, `data-file-drop`, `@drop.prevent`, UploadMediaBoxesEvent, chooseFileEvent, openFileChooserEvent, "Drop the image here", `accept='image/*'`} {
		if !strings.Contains(out, want) {
			t.Errorf("an empty box: no %s", want)
		}
	}
	if strings.Contains(out, `multiple`) {
		t.Error("an empty box takes one file")
	}
	// with a file: a card as an image's of a list — the image alone in the
	// lightbox, its actions over it (edit, replace, delete), the crops of its
	// sizes —, no drop
	out = render(&media_library.MediaBox{ID: "1", Url: "/system/1/file.png", FileName: "file.png", Width: 80, Height: 60},
		&media_library.MediaBoxConfig{Sizes: map[string]*base.Size{"thumb": {Width: 40, Height: 30}}})
	for _, want := range []string{`data-media-card='Foto'`, `data-image-card`, `data-lightbox-src='/system/1/file.png'`, `data-media-edit`, `data-media-replace`,
		`data-media-delete`, `data-image-card-size='thumb'`, `80 × 60`, editImageEvent, deleteFileEvent} {
		if !strings.Contains(out, want) {
			t.Errorf("a box with a file: no %s", want)
		}
	}
	if strings.Contains(out, "data-media-box-drop") || strings.Contains(out, "data-lightbox='") {
		t.Error("a box with a file: a drop, or a group of the lightbox")
	}
	// read only: no actions
	out2, _ := h.MarshallString(mediaBoxThumbnails(ectx, &media_library.MediaBox{ID: "1", Url: "/system/1/file.png", FileName: "file.png"}, "Foto", &media_library.MediaBoxConfig{}, false, true), ctx)
	if strings.Contains(out2, "data-image-card-actions") {
		t.Error("read only: actions")
	}
	// a file not an image: its icon, no lightbox nor editor
	out3, _ := h.MarshallString(mediaBoxThumbnails(ectx, &media_library.MediaBox{ID: "2", Url: "/system/2/nota.pdf", FileName: "nota.pdf"}, "Foto", &media_library.MediaBoxConfig{}, false, false), ctx)
	if strings.Contains(out3, "data-lightbox-src") || strings.Contains(out3, "data-media-edit") || !strings.Contains(out3, "data-media-delete") {
		t.Error("a file: no lightbox, no editor, its delete")
	}
}
