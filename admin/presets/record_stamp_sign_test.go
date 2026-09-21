package presets

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

// stampedApp registers the guarded model SignForm signs for: the stamp is found
// through the model's own editing builder.
func stampedApp() *Builder {
	b := New(i18n.New())
	b.Model(&stampedModel{}).URIName("stamped")
	return b
}

// multipartRequest is the kind of request a test builds by hand: a couple of
// fields and a file.
func multipartRequest(t *testing.T) *http.Request {
	t.Helper()

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	if err := w.WriteField("Title", "a título"); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteField("Tags", "one"); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteField("Tags", "two"); err != nil {
		t.Fatal(err)
	}
	fw, err := w.CreateFormFile("Cover", "cover.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = fw.Write([]byte("\x89PNG not really")); err != nil {
		t.Fatal(err)
	}
	if err = w.Close(); err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("POST", "/posts?__execute_event__=presets_Update&id=1", &body)
	r.Header.Set("Content-Type", w.FormDataContentType())
	return r
}

// The signed request carries the stamp of the record, and the form it already
// had arrives whole — every field, the repeated one, and the file.
func TestSignFormKeepsTheFormAndAddsTheStamp(t *testing.T) {
	b := stampedApp()
	at := time.Date(2026, 7, 27, 10, 30, 0, 0, time.UTC)

	signed, err := b.SignForm(multipartRequest(t), &stampedModel{ID: 1, UpdatedAt: at})
	if err != nil {
		t.Fatal(err)
	}

	if err = signed.ParseMultipartForm(1 << 20); err != nil {
		t.Fatal(err)
	}

	if got, want := signed.FormValue("Title"), "a título"; got != want {
		t.Errorf("Title = %q, want %q", got, want)
	}
	if got := signed.MultipartForm.Value["Tags"]; len(got) != 2 || got[0] != "one" || got[1] != "two" {
		t.Errorf("Tags = %v, want [one two]", got)
	}

	fh := signed.MultipartForm.File["Cover"]
	if len(fh) != 1 {
		t.Fatalf("o arquivo não sobreviveu: %v", signed.MultipartForm.File)
	}
	if fh[0].Filename != "cover.png" {
		t.Errorf("filename = %q, want cover.png", fh[0].Filename)
	}
	f, err := fh[0].Open()
	if err != nil {
		t.Fatal(err)
	}
	content, _ := io.ReadAll(f)
	if string(content) != "\x89PNG not really" {
		t.Errorf("conteúdo do arquivo = %q", content)
	}

	// the stamp itself: signed, and unsigning gives the record's instant
	got, err := b.FormSigner().Unsign(signed.FormValue(RecordStampFormKey))
	if err != nil {
		t.Fatalf("o valor não veio assinado: %v", err)
	}
	if want := RecordStampValue(at); got != want {
		t.Errorf("stamp = %q, want %q", got, want)
	}
}

// The verifier is the reader of this field: a request the helper signed passes
// it, and the same request without the stamp does not.
func TestSignFormSatisfiesVerifyRecordStamp(t *testing.T) {
	at := time.Date(2026, 7, 27, 10, 30, 0, 0, time.UTC)

	b := New(i18n.New())
	mb := b.Model(&stampedModel{}).URIName("posts")
	mb.Editing().FetchFunc(func(obj any, id ID, ctx *web.EventContext) error {
		obj.(*stampedModel).ID = 1
		obj.(*stampedModel).UpdatedAt = at
		return nil
	})

	// the fetcher above ignores the id, and parsing one would need a data
	// operator this test has no reason to provide
	var mid ID

	verify := func(r *http.Request) error {
		ctx := &web.EventContext{R: r}
		return mb.Editing().VerifyRecordStamp(mid, ctx)
	}

	unsigned := multipartRequest(t)
	if err := verify(unsigned); !errors.Is(err, ErrRecordStampMissing) {
		t.Errorf("sem o stamp: err = %v, want ErrRecordStampMissing", err)
	}

	signed, err := b.SignForm(multipartRequest(t), &stampedModel{ID: 1, UpdatedAt: at})
	if err != nil {
		t.Fatal(err)
	}
	if err := verify(signed); err != nil {
		t.Errorf("o request assinado foi recusado: %v", err)
	}

	// signing an out-of-date record is what the guard is for
	stale, err := b.SignForm(multipartRequest(t), &stampedModel{ID: 1, UpdatedAt: at.Add(-time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	if err := verify(stale); !errors.Is(err, ErrRecordChanged) {
		t.Errorf("stamp velho: err = %v, want ErrRecordChanged", err)
	}
}

// A model with no UpdatedAt is stamped all the same: the stamp is the hash of
// the fields its form edits.
func TestSignFormStampsAModelWithNoUpdatedAt(t *testing.T) {
	b, mb := stateApp(t)
	obj := baseStateModel()

	signed, err := b.SignForm(multipartRequest(t), obj)
	if err != nil {
		t.Fatal(err)
	}
	if err = signed.ParseMultipartForm(1 << 20); err != nil {
		t.Fatal(err)
	}

	got, err := b.FormSigner().Unsign(signed.FormValue(RecordStampFormKey))
	if err != nil {
		t.Fatalf("o valor não veio assinado: %v", err)
	}
	if want := mb.Editing().RecordStateHash(obj); got != want {
		t.Errorf("stamp = %q, want o hash do estado %q", got, want)
	}
}

// The caller's request is not touched, and stays readable.
func TestSignFormDoesNotTouchTheOriginal(t *testing.T) {
	b := stampedApp()
	r := multipartRequest(t)

	if _, err := b.SignForm(r, &stampedModel{ID: 1, UpdatedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		t.Fatalf("o body do original não sobreviveu: %v", err)
	}
	if got := r.FormValue(RecordStampFormKey); got != "" {
		t.Errorf("o original ganhou o stamp: %q", got)
	}
	if got, want := r.FormValue("Title"), "a título"; got != want {
		t.Errorf("Title do original = %q, want %q", got, want)
	}
}

func TestSignFormURLEncoded(t *testing.T) {
	b := stampedApp()
	at := time.Date(2026, 7, 27, 10, 30, 0, 0, time.UTC)

	r := httptest.NewRequest("POST", "/posts", strings.NewReader("Title=a+t%C3%ADtulo&Tags=one&Tags=two"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	signed, err := b.SignForm(r, &stampedModel{ID: 1, UpdatedAt: at})
	if err != nil {
		t.Fatal(err)
	}
	if err = signed.ParseForm(); err != nil {
		t.Fatal(err)
	}
	if got, want := signed.PostFormValue("Title"), "a título"; got != want {
		t.Errorf("Title = %q, want %q", got, want)
	}
	if got := signed.PostForm["Tags"]; len(got) != 2 {
		t.Errorf("Tags = %v, want dois valores", got)
	}
	got, err := b.FormSigner().Unsign(signed.PostFormValue(RecordStampFormKey))
	if err != nil {
		t.Fatal(err)
	}
	if want := RecordStampValue(at); got != want {
		t.Errorf("stamp = %q, want %q", got, want)
	}
}

// A body that is not a form has nowhere to put the field, and saying so beats
// returning a request that would be refused later.
func TestSignFormRefusesABodyThatIsNotAForm(t *testing.T) {
	b := stampedApp()

	r := httptest.NewRequest("POST", "/posts", strings.NewReader(`{"Title":"x"}`))
	r.Header.Set("Content-Type", "application/json")

	if _, err := b.SignForm(r, &stampedModel{ID: 1, UpdatedAt: time.Now()}); err == nil {
		t.Error("um body JSON devia dar erro")
	}
}
