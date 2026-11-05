package media

import (
	"mime/multipart"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

type setterObj struct {
	Midia media_library.MediaBox
}

// The media box of an item of a list editor reads the values as decoded —
// reindexed by the order of the items (__pos) —, not r.Form, which keeps the
// indexes as posted: a row moved keeps its media.
func TestMediaBoxSetterReadsDecodedValues(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)
	// as posted: the media of the row at [1]
	r.Form = url.Values{"Itens[1].Midia.Values": {`{"ID":"2","FileName":"b.pdf"}`}}
	// as decoded: the row moved to [0] (the one at [0] removed)
	r.MultipartForm = &multipart.Form{Value: map[string][]string{
		"Itens[0].Midia.Values": {`{"ID":"2","FileName":"b.pdf"}`},
	}}
	obj := &setterObj{}
	field := &presets.FieldContext{Name: "Midia", FormKey: "Itens[0].Midia"}
	if err := MediaBoxSetterFunc(nil)(obj, field, &web.EventContext{R: r}); err != nil {
		t.Fatal(err)
	}
	if obj.Midia.FileName != "b.pdf" {
		t.Errorf("the media of the row moved: %+v", obj.Midia)
	}
}

type captionObj struct {
	Midia   media_library.MediaBox
	Legenda string
}

// With a DescriptionField, the media box's description is the record's field
// beside it — its form key a sibling of the box's —, out of the box.
func TestMediaBoxDescriptionField(t *testing.T) {
	cfg := &media_library.MediaBoxConfig{DescriptionField: "Legenda"}
	if k := MediaBoxDescriptionKey("Itens[0].Imagens[1].Midia", cfg); k != "Itens[0].Imagens[1].Legenda" {
		t.Errorf("the key: %q", k)
	}
	if k := MediaBoxDescriptionKey("Midia", cfg); k != "Legenda" {
		t.Errorf("the key at the root: %q", k)
	}
	if k := MediaBoxDescriptionKey("Midia", nil); k != "Midia.Description" {
		t.Errorf("the box's own: %q", k)
	}

	r := httptest.NewRequest("POST", "/", nil)
	r.MultipartForm = &multipart.Form{Value: map[string][]string{
		"Itens[0].Midia.Values": {`{"ID":"2","FileName":"a.png"}`},
		"Itens[0].Legenda":      {"a fachada"},
	}}
	r.Form = url.Values(r.MultipartForm.Value)
	obj := &captionObj{}
	field := &presets.FieldContext{Name: "Midia", FormKey: "Itens[0].Midia", Obj: obj}
	field.SetContextValue(MediaBoxConfig, cfg)
	if err := MediaBoxSetterFunc(nil)(obj, field, &web.EventContext{R: r}); err != nil {
		t.Fatal(err)
	}
	if obj.Legenda != "a fachada" || obj.Midia.Description != "" || obj.Midia.FileName != "a.png" {
		t.Errorf("saved: %+v", obj)
	}
	if d := MediaBoxDescription(field, cfg); d != "a fachada" {
		t.Errorf("shown: %q", d)
	}
}
