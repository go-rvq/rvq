package media

import (
	"image"
	_ "image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/media/oss"
	"github.com/go-rvq/rvq/admin/media/storage"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestCenteredCrop(t *testing.T) {
	// a wide image, a square: its height whole, centered across
	if c := centeredCrop(400, 200, 1, 1); *c != (base.CropOption{X: 100, Y: 0, Width: 200, Height: 200}) {
		t.Errorf("square of a wide image: %+v", c)
	}
	// a tall image, 2.5:1: its width whole, centered down
	if c := centeredCrop(250, 400, 5, 2); *c != (base.CropOption{X: 0, Y: 150, Width: 250, Height: 100}) {
		t.Errorf("2.5:1 of a tall image: %+v", c)
	}
	if centeredCrop(0, 10, 1, 1) != nil {
		t.Error("a crop of an image of no width")
	}
	if !sameRatio(&base.CropOption{Width: 300, Height: 100}, 1200, 400) || sameRatio(&base.CropOption{Width: 300, Height: 100}, 1200, 480) {
		t.Error("sameRatio")
	}
}

// The sizes of a field: the admin's when it has some, else the code's — what
// the code gave besides the measures kept —; the server reads them by the
// use, not as the browser sends them.
func TestUsageSizes(t *testing.T) {
	db, _, _ := usesTestDB(t)
	if err := db.AutoMigrate(&MediaBoxUsage{}); err != nil {
		t.Fatal(err)
	}
	b := &Builder{db: db}
	db.Create(&MediaBoxUsage{Key: "uses_posts.Cover", RecordTable: "uses_posts", Field: "Cover",
		DefaultSizes: encodeSizes([]UsageSize{{Name: "thumb", Width: 100, Height: 100}})})

	// the browser's sizes, not trusted: the code's
	cfg := b.cfgOf(`{"Usage":"uses_posts.Cover","Sizes":{"thumb":{"Width":1,"Height":1,"Cols":4}}}`)
	if s := cfg.Sizes["thumb"]; s.Width != 100 || s.Height != 100 || s.Cols != 4 {
		t.Errorf("the code's: %+v", s)
	}
	// the admin's
	db.Model(&MediaBoxUsage{}).Where("usage_key = ?", "uses_posts.Cover").
		Update("sizes", encodeSizes([]UsageSize{{Name: "thumb", Width: 200, Height: 80}}))
	if s := b.cfgOf(`{"Usage":"uses_posts.Cover"}`).Sizes["thumb"]; s.Width != 200 || s.Height != 80 {
		t.Errorf("the admin's: %+v", s)
	}
	// no use: as sent
	if s := b.cfgOf(`{"Sizes":{"x":{"Width":3,"Height":4}}}`).Sizes["x"]; s.Width != 3 {
		t.Errorf("no use: %+v", s)
	}
	// a name another field has in other measures
	db.Create(&MediaBoxUsage{Key: "pages.Cover", DefaultSizes: encodeSizes([]UsageSize{{Name: "thumb", Width: 50, Height: 50}})})
	if b.sizesConflict("uses_posts.Cover", []UsageSize{{Name: "thumb", Width: 200, Height: 80}}) == "" {
		t.Error("no conflict of a name of other measures")
	}
	if c := b.sizesConflict("uses_posts.Cover", []UsageSize{{Name: "thumb", Width: 50, Height: 50}}); c != "" {
		t.Errorf("a conflict of the same measures: %s", c)
	}
}

// A file chosen for a field: a size it has in other measures (the field's
// changed) is made again, not kept.
func TestMergeNewSizes(t *testing.T) {
	var m media_library.MediaLibrary
	m.File.Sizes = map[string]*base.Size{"slider": {Width: 1200, Height: 400}, "thumb": {Width: 50, Height: 50}}
	sizes, crop := mergeNewSizes(&m, &media_library.MediaBoxConfig{Sizes: map[string]*base.Size{
		"slider": {Width: 1200, Height: 480}, "thumb": {Width: 50, Height: 50}}})
	if !crop || sizes["slider"].Height != 480 {
		t.Errorf("the size changed: %v %+v", crop, sizes["slider"])
	}
	if _, crop = mergeNewSizes(&m, &media_library.MediaBoxConfig{Sizes: map[string]*base.Size{"thumb": {Width: 50, Height: 50}}}); crop {
		t.Error("made again with the same measures")
	}
}

// usesSizesApp is an admin of usesPost — its Cover of the size thumb 100×100
// — with the media library, its fields of images registered.
func usesSizesApp(t *testing.T, db *gorm.DB) (*presets.Builder, *Builder) {
	t.Helper()
	pb := presets.New(i18n.New()).URIPrefix("/admin")
	pb.DataOperator(gorm2op.DataOperator(db))
	pb.Permission(perm.New().AllowAll().SubjectsFunc(func(*http.Request) []string { return []string{"x"} }))
	mb := New(db)
	if err := mb.Install(pb); err != nil {
		t.Fatal(err)
	}
	pm := pb.Model(&usesPost{})
	pm.Editing("Title", "Cover").Field("Cover").WithContextValue(MediaBoxConfig, &media_library.MediaBoxConfig{
		AllowType: "image", Sizes: map[string]*base.Size{"thumb": {Width: 100, Height: 100}}})
	if err := mb.RegisterUsages(pb); err != nil {
		t.Fatal(err)
	}
	return pb, mb
}

// The fields of images of the models registered with the code's sizes; the
// admin's kept as the application starts again.
func TestRegisterUsages(t *testing.T) {
	db, _, _ := usesTestDB(t)
	pb, mb := usesSizesApp(t, db)
	u := mb.usage("uses_posts.Cover")
	if u == nil || u.DefaultSizes != `[{"name":"thumb","width":100,"height":100}]` || u.Sizes != "" {
		t.Fatalf("registered: %+v", u)
	}
	db.Model(&MediaBoxUsage{}).Where("usage_key = ?", u.Key).Update("sizes", `[{"name":"thumb","width":60,"height":60}]`)
	if err := mb.RegisterUsages(pb); err != nil {
		t.Fatal(err)
	}
	if u = mb.usage(u.Key); u.Sizes != `[{"name":"thumb","width":60,"height":60}]` {
		t.Errorf("the admin's lost: %+v", u)
	}
}

// uploaded is a file of the media library: the test's image (357×340),
// stored in the test's storage.
func uploaded(t *testing.T, db *gorm.DB, sizes map[string]*base.Size) media_library.MediaLibrary {
	t.Helper()
	dir := t.TempDir()
	was := oss.Storage
	oss.Storage = storage.NewFileSystem(dir)
	t.Cleanup(func() { oss.Storage = was })
	png, err := os.ReadFile("integration/testfile.png")
	if err != nil {
		t.Fatal(err)
	}
	m := media_library.MediaLibrary{SelectedType: media_library.ALLOW_TYPE_IMAGE}
	if err = m.File.Scan(multipartestutils.CreateMultipartFileHeader("test.png", png)); err != nil {
		t.Fatal(err)
	}
	m.File.Sizes = sizes
	if err = base.SaveUploadAndCropImage(&base.Config{}, db, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

// imageSize is the width and height of the crop name of the file m.
func imageSize(t *testing.T, m media_library.MediaLibrary, name string) (int, int) {
	t.Helper()
	f, err := os.Open(filepath.Join(oss.Storage.(*storage.FileSystem).Base, m.File.URL(name)))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	c, _, err := image.DecodeConfig(f)
	if err != nil {
		t.Fatal(err)
	}
	return c.Width, c.Height
}

// The sizes of a field changed: the crops of its files made again in them —
// of the new ratio, centered —, and the box of each record told; its state,
// done.
func TestRecrop(t *testing.T) {
	db, _, _ := usesTestDB(t)
	_, mb := usesSizesApp(t, db)
	m := uploaded(t, db, map[string]*base.Size{"thumb": {Width: 100, Height: 100}})
	if w, h := imageSize(t, m, "thumb"); w != 100 || h != 100 {
		t.Fatalf("the first crop: %dx%d", w, h)
	}
	p := usesPost{ID: uuid.New(), Cover: media_library.MediaBox{ID: media_library.MediaID(m.ID.String()),
		Url: m.File.Url, FileSizes: m.File.FileSizes}}
	if err := db.Create(&p).Error; err != nil {
		t.Fatal(err)
	}

	db.Model(&MediaBoxUsage{}).Where("usage_key = ?", "uses_posts.Cover").
		Update("sizes", `[{"name":"thumb","width":120,"height":48}]`)
	if err := mb.Recrop("uses_posts.Cover"); err != nil {
		t.Fatal(err)
	}
	var got media_library.MediaLibrary
	db.First(&got, "id = ?", m.ID)
	if w, h := imageSize(t, got, "thumb"); w != 120 || h != 48 {
		t.Errorf("the crop made again: %dx%d", w, h)
	}
	// 2.5:1 centered in the 357×340 image: its width whole
	if c := got.File.CropOptions["thumb"]; c == nil || c.X != 0 || c.Width != 357 || c.Height != 143 || c.Y != 99 {
		t.Errorf("the area: %+v", c)
	}
	var post usesPost
	db.First(&post, "id = ?", p.ID)
	if post.Cover.FileSizes["thumb"] != got.File.FileSizes["thumb"] || post.Cover.FileSizes["thumb"] == 0 {
		t.Errorf("the box not told: %v, the file's %v", post.Cover.FileSizes, got.File.FileSizes)
	}
	u := mb.usage("uses_posts.Cover")
	if u.RecropState != RecropDone || u.RecropDone != 1 || u.RecropTotal != 1 || u.AppliedSizes != u.Sizes {
		t.Errorf("the state: %+v", u)
	}
	// being made: not twice
	db.Model(&MediaBoxUsage{}).Where("usage_key = ?", u.Key).Update("recrop_state", RecropRunning)
	if mb.StartRecrop(u.Key) {
		t.Error("a recrop started twice")
	}
}

// The admin of the sizes: the fields listed, the code's sizes; a field's
// sizes saved — the admin's —, its crops made again; sizes of a name another
// field has in other measures, refused.
func TestSizesAdmin(t *testing.T) {
	db, _, _ := usesTestDB(t)
	pb, mb := usesSizesApp(t, db)
	db.Create(&MediaBoxUsage{Key: "pages.Cover", RecordTable: "pages", Field: "Cover",
		DefaultSizes: encodeSizes([]UsageSize{{Name: "wide", Width: 300, Height: 100}})})
	mux := http.NewServeMux()
	pb.Build(mux)
	serve := func(req *http.Request) string {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		return w.Body.String()
	}
	// the listing's rows, as the page asks for them
	list := serve(multipartestutils.NewMultipartBuilder().PageURL("/admin/media-sizes").
		EventFunc(actions.ReloadList).BuildEventFuncRequest())
	// the field and its sizes, the row opening its record (by its key)
	if !strings.Contains(list, "Uses Posts") || !strings.Contains(list, "thumb 100×100") ||
		!strings.Contains(list, `$presetsItemDetailing.id = \"uses_posts.Cover\"`) {
		t.Fatalf("the listing: %.900s", list)
	}
	event := func(name string, fields [][2]string) string {
		b := multipartestutils.NewMultipartBuilder().PageURL("/admin/media-sizes").EventFunc(name).
			Query(presets.ParamID, "uses_posts.Cover")
		for _, f := range fields {
			b = b.AddField(f[0], f[1])
		}
		return serve(b.BuildEventFuncRequest())
	}
	edit := event(actions.Edit, nil)
	stamp := regexp.MustCompile(`__formSign\\?":\s*\\?"([^"\\]+)`).FindStringSubmatch(edit)
	if stamp == nil {
		t.Fatalf("the edit form carries no stamp: %.800s", edit)
	}
	save := func(name string, w, h int) string {
		return event(actions.Update, [][2]string{{presets.RecordStampFormKey, stamp[1]},
			{"Sizes[0].name", name}, {"Sizes[0].width", strconv.Itoa(w)}, {"Sizes[0].height", strconv.Itoa(h)}})
	}
	// of a name pages.Cover has in other measures: refused
	body := save("wide", 200, 80)
	if u := mb.usage("uses_posts.Cover"); u.Sizes != "" || !strings.Contains(body, "pages.Cover") {
		t.Errorf("a conflict saved: %+v %.400s", u, body)
	}
	body = save("thumb", 120, 48)
	u := mb.usage("uses_posts.Cover")
	if u.Sizes != `[{"name":"thumb","width":120,"height":48}]` {
		t.Fatalf("the admin's not saved: %+v\n%.600s", u, body)
	}
	// its crops being made again (no file: done at once)
	for i := 0; i < 50 && mb.usage(u.Key).RecropState == RecropRunning; i++ {
		time.Sleep(20 * time.Millisecond)
	}
	if got := mb.usage(u.Key); got.RecropState != RecropDone || got.AppliedSizes != got.Sizes {
		t.Errorf("the crops after a save: %+v", got)
	}
}

// The crop of an image by hand, in the admin: its cropper opens for a file of
// a UUID key, and the area chosen makes the crop again.
func TestCropByHand(t *testing.T) {
	db, _, _ := usesTestDB(t)
	pb, _ := usesSizesApp(t, db)
	m := uploaded(t, db, map[string]*base.Size{"thumb": {Width: 100, Height: 100}})
	mux := http.NewServeMux()
	pb.Build(mux)
	event := func(name string, fields [][2]string) string {
		b := multipartestutils.NewMultipartBuilder().PageURL("/admin/uses-posts").EventFunc(name).
			Query("field", "Cover").Query("id", m.ID.String()).Query("thumb", "thumb").
			AddField("cfg", `{"Usage":"uses_posts.Cover"}`)
		for _, f := range fields {
			b = b.AddField(f[0], f[1])
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, b.BuildEventFuncRequest())
		return w.Body.String()
	}
	if body := event(loadImageCropperEvent, nil); !strings.Contains(body, "cropper") || !strings.Contains(body, m.ID.String()) {
		t.Fatalf("the cropper did not open: %.400s", body)
	}
	box := `{"ID":"` + m.ID.String() + `","Url":"` + m.File.Url + `"}`
	event(cropImageEvent, [][2]string{{"Cover.Values", box}, {"CropOption", `{"x":10,"y":20,"width":200,"height":200}`}})
	var got media_library.MediaLibrary
	db.First(&got, "id = ?", m.ID)
	if c := got.File.CropOptions["thumb"]; c == nil || c.X != 10 || c.Y != 20 || c.Width != 200 {
		t.Errorf("the area chosen: %+v", c)
	}
	if w, h := imageSize(t, got, "thumb"); w != 100 || h != 100 {
		t.Errorf("the crop: %dx%d", w, h)
	}
}
