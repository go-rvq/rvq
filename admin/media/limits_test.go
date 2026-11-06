package media

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/media/oss"
	"github.com/go-rvq/rvq/admin/media/storage"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
)

// The size an image fits a limit in: its proportions kept, the limit turned
// with the image (a portrait one fits it on its side).
func TestFitSize(t *testing.T) {
	fullHD := base.Size{Width: 1920, Height: 1080}
	for _, c := range []struct{ w, h, ww, wh int }{
		{1280, 720, 1280, 720},   // fits: as it is
		{4000, 3000, 1440, 1080}, // landscape
		{4000, 1200, 1920, 576},  // wide: the width limits it
		{3000, 4000, 1080, 1440}, // portrait: the limit turned
		{3000, 3000, 1080, 1080}, // square: the limit as it is
	} {
		if s := base.FitSize(c.w, c.h, fullHD); s.Width != c.ww || s.Height != c.wh {
			t.Errorf("%d×%d: %d×%d, want %d×%d", c.w, c.h, s.Width, s.Height, c.ww, c.wh)
		}
	}
	// a side of 0: no limit on it
	if s := base.NewSize(5000, 100).RescaleToMax(2000, 0); s.Width != 2000 || s.Height != 40 {
		t.Errorf("no limit of the height: %+v", s)
	}
}

// The limits of a box: its own; with no sizes, Full HD and 320×320; with
// sizes, none (the media library's).
func TestBoxLimits(t *testing.T) {
	var none *media_library.MediaBoxConfig
	if o, th := none.OriginalLimit(), none.ThumbLimit(); *o != media_library.DefaultOriginalMax || *th != media_library.DefaultThumbMax {
		t.Errorf("no config: %+v %+v", o, th)
	}
	sized := &media_library.MediaBoxConfig{Sizes: map[string]*base.Size{"x": {Width: 10, Height: 10}}}
	if sized.OriginalLimit() != nil || sized.ThumbLimit() != nil {
		t.Error("a box of sizes has limits of its own")
	}
	own := &media_library.MediaBoxConfig{Sizes: sized.Sizes, OriginalMax: &base.Size{Width: 800, Height: 600}, ThumbMax: &base.Size{Width: 90, Height: 90}}
	if o, th := own.OriginalLimit(), own.ThumbLimit(); o.Width != 800 || th.Width != 90 {
		t.Errorf("its own: %+v %+v", o, th)
	}
}

func gradPNG(t *testing.T, w, h int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y += 7 {
		for x := 0; x < w; x += 7 {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 90, 255})
		}
	}
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// An upload to a box: the original made smaller to its limit when larger
// (the browser makes it so before; the server checks it); its thumbnail the
// browser's when right, else made by the server.
func TestUploadLimits(t *testing.T) {
	db, _, _ := usesTestDB(t)
	pb := presets.New(i18n.New()).URIPrefix("/admin")
	pb.DataOperator(gorm2op.DataOperator(db))
	b := New(db)
	if err := b.Install(pb); err != nil {
		t.Fatal(err)
	}
	old := oss.Storage
	oss.Storage = storage.NewFileSystem(t.TempDir())
	t.Cleanup(func() { oss.Storage = old })

	upload := func(cfg string, file, thumb []byte) media_library.MediaLibrary {
		t.Helper()
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		fw, _ := w.CreateFormFile("NewFiles", "foto.png")
		fw.Write(file)
		if thumb != nil {
			tw, _ := w.CreateFormFile("NewThumb0", "foto.png")
			tw.Write(thumb)
		}
		w.WriteField("cfg", cfg)
		w.Close()
		r := httptest.NewRequest("POST", "/admin/x?__execute_event__="+UploadMediaBoxesEvent, &body)
		r.Header.Set("Content-Type", w.FormDataContentType())
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			t.Fatal(err)
		}
		resp, err := uploadMediaBoxes(b)(&web.EventContext{R: r})
		if err != nil {
			t.Fatal(err)
		}
		boxes, _ := resp.Data.([]media_library.MediaBox)
		if len(boxes) != 1 {
			t.Fatalf("uploaded: %+v", boxes)
		}
		var m media_library.MediaLibrary
		if err := db.First(&m, "id = ?", boxes[0].ID.String()).Error; err != nil {
			t.Fatal(err)
		}
		if boxes[0].Width != m.File.Width || boxes[0].Height != m.File.Height {
			t.Errorf("the box's measures: %+v, the file's %d×%d", boxes[0], m.File.Width, m.File.Height)
		}
		return m
	}
	measures := func(m media_library.MediaLibrary, name string) string {
		t.Helper()
		w, h := imageSize(t, m, name)
		return fmt.Sprintf("%d×%d", w, h)
	}

	// no sizes: Full HD, a thumbnail of 320×320 made here
	m := upload(`{"AllowType":"image"}`, gradPNG(t, 3000, 2000), nil)
	if got := measures(m, "original"); got != "1620×1080" || m.File.Width != 1620 {
		t.Errorf("the original: %s (%d×%d)", got, m.File.Width, m.File.Height)
	}
	if got := measures(m, media_library.ThumbSizeName); got != "320×213" || m.File.FileSizes[media_library.ThumbSizeName] == 0 {
		t.Errorf("the thumbnail: %s %v", got, m.File.FileSizes)
	}
	// a portrait one: the limit on its side
	if got := measures(upload(`{"AllowType":"image"}`, gradPNG(t, 2000, 3000), nil), "original"); got != "1080×1620" {
		t.Errorf("portrait: %s", got)
	}
	// the box's own limit
	if got := measures(upload(`{"AllowType":"image","OriginalMax":{"Width":800,"Height":600}}`, gradPNG(t, 3000, 2000), nil), "original"); got != "800×533" {
		t.Errorf("its own limit: %s", got)
	}
	// with sizes: as it is (the media library's limit, none), no thumbnail
	m = upload(`{"AllowType":"image","Sizes":{"card":{"Width":100,"Height":100}}}`, gradPNG(t, 3000, 2000), nil)
	if got := measures(m, "original"); got != "3000×2000" || m.File.FileSizes[media_library.ThumbSizeName] != 0 {
		t.Errorf("with sizes: %s %v", got, m.File.FileSizes)
	}

	// the browser's thumbnail, right: kept as sent
	thumb := gradPNG(t, 320, 213)
	m = upload(`{"AllowType":"image"}`, gradPNG(t, 1620, 1080), thumb)
	if m.File.FileSizes[media_library.ThumbSizeName] != len(thumb) || measures(m, media_library.ThumbSizeName) != "320×213" {
		t.Errorf("the browser's not kept: %v", m.File.FileSizes)
	}
	// not of the image's measures: made here
	bad := gradPNG(t, 100, 100)
	m = upload(`{"AllowType":"image"}`, gradPNG(t, 1620, 1080), bad)
	if m.File.FileSizes[media_library.ThumbSizeName] == len(bad) || measures(m, media_library.ThumbSizeName) != "320×213" {
		t.Errorf("a wrong one kept: %v %s", m.File.FileSizes, measures(m, media_library.ThumbSizeName))
	}
}

// A box of no sizes registered too, with the code's limits; the admin's —
// trusted, not the browser's — given to its config, and saved in its form.
func TestUsageLimits(t *testing.T) {
	db, _, _ := usesTestDB(t)
	pb := presets.New(i18n.New()).URIPrefix("/admin")
	pb.DataOperator(gorm2op.DataOperator(db))
	pb.Permission(perm.New().AllowAll().SubjectsFunc(func(*http.Request) []string { return []string{"x"} }))
	mb := New(db)
	if err := mb.Install(pb); err != nil {
		t.Fatal(err)
	}
	pm := pb.Model(&usesPost{})
	ed := pm.Editing("Title", "Cover", "Thumb")
	ed.Field("Cover").WithContextValue(MediaBoxConfig, &media_library.MediaBoxConfig{
		AllowType: "image", Sizes: map[string]*base.Size{"card": {Width: 100, Height: 100}}})
	if err := mb.RegisterUsages(pb); err != nil {
		t.Fatal(err)
	}
	cover, thumb := mb.usage("uses_posts.Cover"), mb.usage("uses_posts.Thumb")
	if cover == nil || cover.DefaultOriginalWidth != 0 || cover.DefaultThumbWidth != 0 {
		t.Errorf("a box of sizes: %+v", cover)
	}
	if thumb == nil || thumb.DefaultOriginalWidth != 1920 || thumb.DefaultOriginalHeight != 1080 ||
		thumb.DefaultThumbWidth != 320 || thumb.DefaultThumbHeight != 320 {
		t.Fatalf("a box of no sizes: %+v", thumb)
	}
	// the browser's limit not trusted: the code's
	cfg := mb.cfgOf(`{"Usage":"uses_posts.Thumb","OriginalMax":{"Width":9000,"Height":9000}}`)
	if cfg.OriginalMax.Width != 1920 || cfg.ThumbMax.Width != 320 {
		t.Errorf("the code's: %+v %+v", cfg.OriginalMax, cfg.ThumbMax)
	}

	// the admin's, saved in its form
	mux := http.NewServeMux()
	pb.Build(mux)
	event := func(name string, fields [][2]string) string {
		b := multipartestutils.NewMultipartBuilder().PageURL("/admin/media-sizes").EventFunc(name).
			Query(presets.ParamID, "uses_posts.Thumb")
		for _, f := range fields {
			b = b.AddField(f[0], f[1])
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, b.BuildEventFuncRequest())
		return w.Body.String()
	}
	edit := event(actions.Edit, nil)
	if !strings.Contains(edit, "data-original-field") || !strings.Contains(edit, "320×320") {
		t.Errorf("the form of the limits: %.600s", edit)
	}
	stamp := regexp.MustCompile(`__formSign\\?":\s*\\?"([^"\\]+)`).FindStringSubmatch(edit)
	if stamp == nil {
		t.Fatalf("no stamp: %.600s", edit)
	}
	save := func(ow, oh, tw, th string) string {
		return event(actions.Update, [][2]string{{presets.RecordStampFormKey, stamp[1]},
			{"OriginalWidth", ow}, {"OriginalHeight", oh}, {"ThumbWidth", tw}, {"ThumbHeight", th}})
	}
	// a width with no height: refused
	save("1000", "800", "200", "")
	if u := mb.usage("uses_posts.Thumb"); u.OriginalWidth != 0 || u.ThumbWidth != 0 {
		t.Errorf("half a limit saved: %+v", u)
	}
	body := save("1000", "800", "200", "150")
	u := mb.usage("uses_posts.Thumb")
	if u.OriginalWidth != 1000 || u.OriginalHeight != 800 || u.ThumbWidth != 200 || u.ThumbHeight != 150 {
		t.Fatalf("the admin's not saved: %+v\n%.600s", u, body)
	}
	cfg = mb.cfgOf(`{"Usage":"uses_posts.Thumb"}`)
	if cfg.OriginalMax.Width != 1000 || cfg.ThumbMax.Height != 150 {
		t.Errorf("the admin's: %+v %+v", cfg.OriginalMax, cfg.ThumbMax)
	}
}

// The script of an upload: the browser fits the images and makes their
// thumbnails (fitUploads), posted as NewThumb<i>.
func TestUploadScriptLimits(t *testing.T) {
	s := uploadThroughEditor(&media_library.MediaBoxConfig{}, "files", uploadFormScript(&media_library.MediaBoxConfig{}, "fd"), false, "")
	for _, want := range []string{"ed.fitUploads(list, {width: 1920, height: 1080}, {width: 320, height: 320})", "'NewThumb' + i"} {
		if !strings.Contains(s, want) {
			t.Errorf("no %q:\n%s", want, s)
		}
	}
	sized := &media_library.MediaBoxConfig{Sizes: map[string]*base.Size{"x": {Width: 10, Height: 10}}}
	if s = uploadThroughEditor(sized, "files", "", false, ""); !strings.Contains(s, "ed.fitUploads(list, null, null)") {
		t.Errorf("a box of sizes:\n%s", s)
	}
}
