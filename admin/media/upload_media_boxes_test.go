package media

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
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

// A drop of files (UploadMediaBoxesEvent) uploads each to the media library
// and answers their media boxes, in their order — those of the field's type
// only (cfg).
func TestUploadMediaBoxes(t *testing.T) {
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

	png, err := os.ReadFile("integration/testfile.png")
	if err != nil {
		t.Fatal(err)
	}
	upload := func(cfg string) []media_library.MediaBox {
		t.Helper()
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		for _, f := range []struct {
			name string
			data []byte
		}{{"nota.txt", []byte("uma nota")}, {"foto.png", png}} {
			fw, _ := w.CreateFormFile("NewFiles", f.name)
			fw.Write(f.data)
		}
		if cfg != "" {
			w.WriteField("cfg", cfg)
		}
		w.Close()
		r := httptest.NewRequest("POST", "/admin/x?__execute_event__="+UploadMediaBoxesEvent, &body)
		r.Header.Set("Content-Type", w.FormDataContentType())
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		resp, err := uploadMediaBoxes(b)(&web.EventContext{R: r})
		if err != nil {
			t.Fatal(err)
		}
		boxes, _ := resp.Data.([]media_library.MediaBox)
		return boxes
	}

	boxes := upload("")
	if len(boxes) != 2 || boxes[0].FileName != "nota.txt" || boxes[1].FileName != "foto.png" {
		t.Fatalf("any type, both in order: %+v", boxes)
	}
	for _, mb := range boxes {
		if mb.ID.String() == "" || mb.Url == "" {
			t.Errorf("not uploaded: %+v", mb)
		}
	}
	if boxes[1].Width == 0 {
		t.Errorf("the image's size: %+v", boxes[1])
	}

	boxes = upload(`{"AllowType":"image"}`)
	if len(boxes) != 1 || boxes[0].FileName != "foto.png" {
		t.Fatalf("an image field takes the image only: %+v", boxes)
	}
	var n int64
	db.Model(&media_library.MediaLibrary{}).Count(&n)
	if n != 3 {
		t.Errorf("records: %d, want 3", n)
	}
}
