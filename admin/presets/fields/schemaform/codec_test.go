package schemaform

import (
	"net/url"
	"strings"
	"testing"
)

// The value the form edited is carried out by the builder's codec: YAML unless
// the application says otherwise, and the record keeps the schema's order in
// either.
func TestCodec(t *testing.T) {
	s, err := Parse("[]{label str; href}")
	if err != nil {
		t.Fatal(err)
	}
	value := s.Decode(url.Values{
		"V[0].label": {"Facebook"},
		"V[0].href":  {"#"},
	}, "V")

	yamlOut, err := New().EncodeValue(value)
	if err != nil {
		t.Fatal(err)
	}
	if want := "- label: Facebook\n  href: '#'\n"; yamlOut != want {
		t.Errorf("yaml =\n%s\nwant\n%s", yamlOut, want)
	}

	jsonOut, err := New().Codec(JSONEncode, JSONDecode).EncodeValue(value)
	if err != nil {
		t.Fatal(err)
	}
	if want := "[\n  {\n    \"label\": \"Facebook\",\n    \"href\": \"#\"\n  }\n]\n"; jsonOut != want {
		t.Errorf("json =\n%s\nwant\n%s", jsonOut, want)
	}
}

// And read back the same way.
func TestCodecDecodeValue(t *testing.T) {
	for name, tc := range map[string]struct {
		builder *Builder
		stored  string
	}{
		"yaml": {New(), "- a\n- b\n"},
		"json": {New().Codec(JSONEncode, JSONDecode), `["a","b"]`},
	} {
		got, err := tc.builder.DecodeValue(tc.stored)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		items, ok := got.([]any)
		if !ok || len(items) != 2 || items[0] != "a" || items[1] != "b" {
			t.Errorf("%s: %#v", name, got)
		}
	}
}

// An empty text is no value at all, whatever the codec: the caller starts from
// the empty shape instead of from an error.
func TestCodecDecodeEmpty(t *testing.T) {
	for _, b := range []*Builder{New(), New().Codec(JSONEncode, JSONDecode)} {
		for _, stored := range []string{"", "   ", "\n"} {
			got, err := b.DecodeValue(stored)
			if err != nil || got != nil {
				t.Errorf("%q => %#v, %v", stored, got, err)
			}
		}
	}
}

// An application may carry the value any way it likes.
func TestCodecCustom(t *testing.T) {
	b := New().Codec(
		func(value any) (string, error) { return "escrito", nil },
		func(stored string) (any, error) { return strings.ToUpper(stored), nil },
	)

	if out, _ := b.EncodeValue([]any{1}); out != "escrito" {
		t.Errorf("encode = %q", out)
	}
	if out, _ := b.DecodeValue("lido"); out != "LIDO" {
		t.Errorf("decode = %v", out)
	}
}
