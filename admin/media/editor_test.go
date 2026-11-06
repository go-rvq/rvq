package media

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/media/oss"
	"github.com/go-rvq/rvq/admin/media/storage"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func pngOf(w, h int, c color.Color) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	var b bytes.Buffer
	_ = png.Encode(&b, img)
	return b.Bytes()
}

type part struct {
	name, file string
	data       []byte
}

// post is a request of a multipart form: the fields, then the files.
func post(t *testing.T, fields map[string]string, files ...part) *web.EventContext {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for k, v := range fields {
		_ = w.WriteField(k, v)
	}
	for _, f := range files {
		fw, _ := w.CreateFormFile(f.name, f.file)
		_, _ = fw.Write(f.data)
	}
	_ = w.Close()
	r := httptest.NewRequest("POST", "/admin/x", &body)
	r.Header.Set("Content-Type", w.FormDataContentType())
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		t.Fatal(err)
	}
	return &web.EventContext{R: r}
}

func editorBuilder(t *testing.T) (*Builder, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	pb := presets.New(i18n.New()).URIPrefix("/admin")
	pb.DataOperator(gorm2op.DataOperator(db))
	b := New(db)
	if err := b.Install(pb); err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&media_library.MediaLibrary{}); err != nil {
		t.Fatal(err)
	}
	old := oss.Storage
	oss.Storage = storage.NewFileSystem(t.TempDir())
	t.Cleanup(func() { oss.Storage = old })
	return b, db
}

func uploadOne(t *testing.T, b *Builder, cfg string, files ...part) media_library.MediaBox {
	t.Helper()
	resp, err := uploadMediaBoxes(b)(post(t, map[string]string{"cfg": cfg}, files...))
	if err != nil {
		t.Fatal(err)
	}
	boxes, _ := resp.Data.([]media_library.MediaBox)
	if len(boxes) != 1 {
		t.Fatalf("uploaded: %+v", boxes)
	}
	return boxes[0]
}

// exists: the file of url is in the storage.
func exists(url string) bool {
	_, err := oss.Storage.Get(url)
	return err == nil
}

// The editor saves the image it made: the result is the file (its size), the
// original kept apart and the edit beside it; reopened, on the original with
// the edit; "Back to the original", the original again.
func TestImageEditSaveReopenRestore(t *testing.T) {
	b, db := editorBuilder(t)
	box := uploadOne(t, b, "", part{"NewFiles", "foto.png", pngOf(80, 60, color.RGBA{200, 0, 0, 255})})
	id := box.ID.String()

	open := func() map[string]any {
		t.Helper()
		r, err := editImage(b)(post(t, map[string]string{"id": id, "thumb": "default", "field": "Foto"}))
		if err != nil {
			t.Fatal(err)
		}
		d, _ := r.Data.(map[string]any)
		return d
	}
	d := open()
	if d["edited"] != false || filepath.Ext(d["name"].(string)) != ".png" {
		t.Fatalf("first opening: %+v", d)
	}

	values, _ := json.Marshal(box)
	edit := `{"version":1,"mode":"perspectiva","perspectiva":[{"x":1,"y":1},{"x":70,"y":2},{"x":72,"y":50},{"x":0,"y":55}]}`
	_, err := saveImageEdit(b)(post(t, map[string]string{
		"id": id, "thumb": "default", "field": "Foto", "Foto.Values": string(values), "Edit": edit,
	}, part{"Result", "foto.png", pngOf(40, 30, color.RGBA{0, 0, 200, 255})}))
	if err != nil {
		t.Fatal(err)
	}
	var m media_library.MediaLibrary
	db.First(&m, "id = ?", id)
	if m.File.Width != 40 || m.File.Height != 30 {
		t.Errorf("the result is the file: %d×%d", m.File.Width, m.File.Height)
	}
	if m.File.Edit == nil || m.File.Edit.Source == "" || !exists(m.File.Edit.Source) {
		t.Fatalf("the original kept apart: %+v", m.File.Edit)
	}
	if string(m.File.Edit.State) != edit {
		t.Errorf("the edit: %s", m.File.Edit.State)
	}

	d = open()
	if d["edited"] != true || d["edit"].(map[string]any)["mode"] != "perspectiva" {
		t.Errorf("reopened on the edit: %+v", d)
	}
	if src := d["src"].(string); src[:len(m.File.Edit.Source)] != m.File.Edit.Source {
		t.Errorf("reopened on the original: %s", src)
	}

	// edited again: the original is still the first one
	source := m.File.Edit.Source
	_, err = saveImageEdit(b)(post(t, map[string]string{
		"id": id, "thumb": "default", "field": "Foto", "Foto.Values": string(values), "Edit": `{"version":1}`,
	}, part{"Result", "foto.png", pngOf(20, 20, color.White)}))
	if err != nil {
		t.Fatal(err)
	}
	m = media_library.MediaLibrary{}
	db.First(&m, "id = ?", id)
	if m.File.Edit.Source != source || m.File.Width != 20 {
		t.Errorf("edited again: %+v %d", m.File.Edit, m.File.Width)
	}

	_, err = saveImageEdit(b)(post(t, map[string]string{
		"id": id, "thumb": "default", "field": "Foto", "Foto.Values": string(values), "Restore": "1",
	}))
	if err != nil {
		t.Fatal(err)
	}
	m = media_library.MediaLibrary{}
	db.First(&m, "id = ?", id)
	if m.File.Edit != nil || m.File.Width != 80 || m.File.Height != 60 {
		t.Errorf("back to the original: %+v %d×%d", m.File.Edit, m.File.Width, m.File.Height)
	}
}

// The fixed crop of a size: the size's CropOption, over the image edited.
func TestImageEditCropOptionOfSize(t *testing.T) {
	b, db := editorBuilder(t)
	cfg := `{"Sizes":{"thumb":{"Width":20,"Height":10}}}`
	box := uploadOne(t, b, cfg, part{"NewFiles", "foto.png", pngOf(80, 60, color.Black)})
	id := box.ID.String()
	values, _ := json.Marshal(box)
	_, err := saveImageEdit(b)(post(t, map[string]string{
		"id": id, "thumb": "thumb", "field": "Foto", "Foto.Values": string(values), "cfg": cfg,
		"CropOption": `{"x":10,"y":5,"width":40,"height":20}`,
	}))
	if err != nil {
		t.Fatal(err)
	}
	var m media_library.MediaLibrary
	db.First(&m, "id = ?", id)
	co := m.File.CropOptions["thumb"]
	if co == nil || co.X != 10 || co.Width != 40 || m.File.Edit != nil {
		t.Fatalf("the size's crop: %+v, edit %+v", co, m.File.Edit)
	}
	r, _ := editImage(b)(post(t, map[string]string{"id": id, "thumb": "thumb", "field": "Foto", "cfg": cfg}))
	d := r.Data.(map[string]any)
	if d["aspect"].(float64) != 2 || d["edit"].(map[string]any)["fixo"] == nil {
		t.Errorf("reopened with the size's crop: %+v", d)
	}
}

// An upload edited in the browser: by default only the result; with
// UploadEdit, its original and the edit kept too.
func TestUploadWithEdit(t *testing.T) {
	b, db := editorBuilder(t)
	result := part{"NewFiles", "cheque.png", pngOf(40, 20, color.White)}
	original := part{"NewOriginal0", "cheque.jpg", pngOf(90, 70, color.Black)}
	edit := `{"version":1,"mode":"perspectiva"}`

	box := uploadOne(t, b, "", result, original)
	var m media_library.MediaLibrary
	db.First(&m, "id = ?", box.ID.String())
	if m.File.Edit != nil || m.File.Width != 40 {
		t.Errorf("by default only the result: %+v", m.File.Edit)
	}

	ctx := post(t, map[string]string{"cfg": `{"UploadEdit":true}`, "NewEdit0": edit}, result, original)
	resp, err := uploadMediaBoxes(b)(ctx)
	if err != nil {
		t.Fatal(err)
	}
	box = resp.Data.([]media_library.MediaBox)[0]
	m = media_library.MediaLibrary{}
	db.First(&m, "id = ?", box.ID.String())
	if m.File.Edit == nil || m.File.Edit.SourceName != "cheque.jpg" || string(m.File.Edit.State) != edit || !exists(m.File.Edit.Source) {
		t.Fatalf("UploadEdit: the original and the edit kept: %+v", m.File.Edit)
	}
	if filepath.Ext(m.File.Edit.Source) != ".jpg" || m.File.Width != 40 {
		t.Errorf("the original apart, the result the file: %s %d", m.File.Edit.Source, m.File.Width)
	}
}
