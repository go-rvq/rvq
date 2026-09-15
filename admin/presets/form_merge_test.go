package presets

import (
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"

	"github.com/go-rvq/rvq/web"
)

func TestFormCarriesField(t *testing.T) {
	ctx := &web.EventContext{R: &http.Request{
		Form: url.Values{
			"Title":            {"Hi"},
			"TitleWithSlug":    {""}, // present but empty
			"Cover.Values":     {"[]"},
			"Tags[0]":          {"a"},
			"Config.NotListed": {"true"},
		},
		MultipartForm: &multipart.Form{
			Value: map[string][]string{"Body": {"x"}},
			File:  map[string][]*multipart.FileHeader{"Attachment.file": nil},
		},
	}}

	cases := map[string]bool{
		"Title":         true,  // exact key
		"TitleWithSlug": true,  // present, empty value still counts
		"Cover":         true,  // structured sub-key Cover.Values
		"Tags":          true,  // indexed sub-key Tags[0]
		"Config":        true,  // nested sub-key Config.NotListed
		"Body":          true,  // multipart value
		"Attachment":    true,  // multipart file sub-key
		"TypeID":        false, // absent entirely
		"Config.TypeID": false, // absent nested field
		"":              true,  // no key to check
	}
	for key, want := range cases {
		if got := formCarriesField(ctx, key); got != want {
			t.Errorf("formCarriesField(%q) = %v, want %v", key, got, want)
		}
	}
}
