package schemaform

import (
	"strings"
	"testing"
)

// A file and an image: what they take ([accept=…]; an image, images), their
// largest size ([maxSize=…]); a list of them, its items' — and [min, max],
// how many items a list holds. Refused where they mean nothing.
func TestFileFields(t *testing.T) {
	s, err := Parse(`interface Form {
		[accept="application/pdf,.docx", maxSize="5MB"] doc file
		photo? image
		[accept="image/png,image/jpeg", maxSize=200000, min=1, max=3] photos []image
		[min=2] files? []file
		[options=["a", "b", "c"], max=2] picks? []str
	}`)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, f := range s.Fields {
		l := f
		if f.Schema != nil && f.Schema.Item != nil {
			l = f.Schema.Item
		}
		got[f.Name] = strings.Join([]string{l.Type, l.AcceptOf(), FormatSize(l.MaxSize), itoa(f.MinItems), itoa(f.MaxItems)}, "|")
	}
	for name, want := range map[string]string{
		"doc":    "file|application/pdf,.docx|5.0 MB|0|0",
		"photo":  "image|image/*|0 B|0|0",
		"photos": "image|image/png,image/jpeg|195.3 KB|1|3",
		"files":  "file||0 B|2|0",
		"picks":  "str||0 B|0|2",
	} {
		if got[name] != want {
			t.Errorf("%s: %q, want %q", name, got[name], want)
		}
	}
	for src, want := range map[string]string{
		`interface Form { [accept=".pdf"] a str }`:              "only a file or an image",
		`interface Form { [accept="pdf"] a file }`:              "neither a media type",
		`interface Form { [accept="application/pdf"] a image }`: "takes images only",
		`interface Form { [maxSize="big"] a file }`:             "want a size",
		`interface Form { [min=3, max=1] a []file }`:            "is more than",
		`interface Form { [min=1.5] a []file }`:                 "whole number of items",
		`interface Form { [minlength=1] a []file }`:             "a list has [min, max]",
	} {
		if _, err := Parse(src); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v (want %q)", src, err, want)
		}
	}
}
