package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/fields/markdown"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MDNote struct {
	ID   uint
	Body string
}

var mdSeq int64

// stored is a picture the Storer was given.
type stored struct {
	format, alt, id, field string
	size                   int
}

// mdApp is the app of a note (1) whose Body is a MarkdownBuilder's, its
// pictures given to a Storer that keeps what it is told; extra denies.
func mdApp(t *testing.T, extra ...*perm.PolicyBuilder) (http.Handler, *[]stored) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:markdown_%d?mode=memory&cache=shared", atomic.AddInt64(&mdSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&MDNote{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&MDNote{ID: 1, Body: "**bold**[^a]\n\n[^a]: the note"})

	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.Permission(perm.New().Policies(append([]*perm.PolicyBuilder{
		perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything),
	}, extra...)...).SubjectsFunc(func(*http.Request) []string { return []string{"editor"} }))
	p.DataOperator(gorm2op.DataOperator(db))
	mb := p.Model(&MDNote{}).URIName("md-notes")
	mb.Listing("ID", "Body")
	mb.Editing("Body")
	mb.Detailing("Body")

	var got []stored
	markdown.New().Model(mb).Fields("Body").
		Store(markdown.StoreFunc(func(ctx context.Context, format string, data []byte, alt string) (string, error) {
			sc := markdown.StoreContextOf(ctx)
			got = append(got, stored{format: format, alt: alt, id: sc.ID.String(), field: sc.Field, size: len(data)})
			return "/system/pictures/" + alt + "." + format, nil
		})).
		Build(presets.FORM | presets.LIST)
	return p, &got
}

// mdFire posts the event of the builder (kind: preview, upload) with fields
// and a file (none: nil).
func mdFire(t *testing.T, h http.Handler, kind string, fields map[string]string, file []byte) map[string]string {
	t.Helper()
	b := multipartestutils.NewMultipartBuilder().PageURL("/admin/md-notes").EventFunc("markdown_" + kind + "_Body")
	for k, v := range fields {
		b = b.AddField(k, v)
	}
	if file != nil {
		b = b.AddReader("file", "My picture.png", bytes.NewReader(file))
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, b.BuildEventFuncRequest())
	var r struct {
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil {
		t.Fatalf("%s: %v: %s", kind, err, w.Body.String())
	}
	return r.Data
}

// mdPage is the page answered, its portal's content unescaped (JSON).
func mdPage(w *httptest.ResponseRecorder) string {
	return strings.NewReplacer(`\u003c`, "<", `\u003e`, ">", `\u0026`, "&", `\"`, `"`).Replace(w.Body.String())
}

var pngBytes = append([]byte("\x89PNG\r\n\x1a\n"), make([]byte, 32)...)

// The MarkdownBuilder: its preview rendered by the server (GFM, footnotes,
// sanitized); a picture stored by its Storer — told the record and the field
// —, its address given; not a picture, or one the user may not put in the
// record, refused; the form with the editor, its preview and its pictures
// by the events; the detail rendered; the listing as text.
func TestMarkdownBuilder(t *testing.T) {
	h, got := mdApp(t)

	out := mdFire(t, h, "preview", map[string]string{"body": "**a** <script>x()</script>\n\n| x |\n|---|\n| 1 |\n\ntext[^n]\n\n[^n]: note"}, nil)
	for _, want := range []string{"<strong>a</strong>", "<table>", `class="footnotes"`} {
		if !strings.Contains(out["html"], want) {
			t.Errorf("preview: no %s in %s", want, out["html"])
		}
	}
	if strings.Contains(out["html"], "<script") {
		t.Errorf("preview not sanitized: %s", out["html"])
	}

	out = mdFire(t, h, "upload", map[string]string{"id": "1"}, pngBytes)
	if out["url"] != "/system/pictures/My picture.png" {
		t.Fatalf("upload: %v", out)
	}
	if len(*got) != 1 || (*got)[0] != (stored{format: "png", alt: "My picture", id: "1", field: "Body", size: len(pngBytes)}) {
		t.Errorf("stored: %+v", *got)
	}
	// a new record's (no id)
	if out = mdFire(t, h, "upload", nil, pngBytes); out["url"] == "" || (*got)[1].id != "" {
		t.Errorf("upload of a new record: %v %+v", out, *got)
	}
	if out = mdFire(t, h, "upload", map[string]string{"id": "1"}, []byte("plain text, not a picture")); out["url"] != "" || out["message"] == "" {
		t.Errorf("not a picture: %v", out)
	}
	if len(*got) != 2 {
		t.Errorf("stored what it refused: %+v", *got)
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/md-notes/1/edit", nil))
	page := mdPage(w)
	for _, want := range []string{"<vx-markdown-editor", "markdown_preview_Body", "markdown_upload_Body"} {
		if !strings.Contains(page, want) {
			t.Errorf("edit page: no %s", want)
		}
	}

	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/md-notes/1", nil))
	if body := mdPage(w); !strings.Contains(body, "<strong>bold</strong>") || !strings.Contains(body, "data-markdown-detail") {
		t.Errorf("detail not rendered")
	}

	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/md-notes", nil))
	if body := mdPage(w); !strings.Contains(body, ">bold[1] [1] the note</td>") {
		i := strings.Index(body, "bold")
		t.Errorf("listing not as text: %q", body[max(0, i-50):min(len(body), i+60)])
	}
}

// Who may not edit the record may not put a picture in it; who may not
// create, in a new one.
func TestMarkdownBuilderPictureDenied(t *testing.T) {
	h, got := mdApp(t, perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(presets.PermUpdate, presets.PermCreate).On(perm.Anything))
	for _, fields := range []map[string]string{{"id": "1"}, nil} {
		if out := mdFire(t, h, "upload", fields, pngBytes); out["url"] != "" || out["message"] == "" {
			t.Errorf("%v: denied, yet %v", fields, out)
		}
	}
	if len(*got) != 0 {
		t.Errorf("stored: %+v", *got)
	}
}
