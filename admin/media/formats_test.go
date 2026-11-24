package media

import (
	"testing"

	"github.com/go-rvq/rvq/admin/media/media_library"
)

// A format of a field: its sizes the field's in its proportion — the
// largest its own, the others scaled by their widths —, named with its
// suffix; a box of the format shows their crops, the field's when it has not
// one.
func TestFormatSizesAndURL(t *testing.T) {
	u := &MediaBoxUsage{DefaultSizes: encodeSizes([]UsageSize{
		{Name: "site.thumb", Width: 500, Height: 281},
		{Name: "site.middle", Width: 900, Height: 506},
		{Name: "site.full", Width: 1920, Height: 1080},
	}), Formats: encodeFormats([]UsageFormat{{Name: "retrato", Label: "Retrato 3:4", Width: 1080, Height: 1440}})}
	got := map[string][2]int{}
	for _, s := range u.FormatSizes(*u.format("retrato")) {
		got[s.Name] = [2]int{s.Width, s.Height}
	}
	want := map[string][2]int{"site.full_retrato": {1080, 1440}, "site.middle_retrato": {506, 675}, "site.thumb_retrato": {281, 375}}
	for k, w := range want {
		if got[k] != w {
			t.Errorf("%s = %v, want %v", k, got[k], w)
		}
	}
	if n := len(u.AllSizes()); n != 6 {
		t.Errorf("all the sizes: %d", n)
	}

	box := &media_library.MediaBox{Url: "/system/m/1/file.jpg", Format: "retrato",
		FileSizes: map[string]int{"site.full": 10, "site.full_retrato": 12}}
	if got := box.URL("site.full"); got != "/system/m/1/file.site.full_retrato.jpg" {
		t.Errorf("the format's crop: %s", got)
	}
	if got := box.URL("site.middle"); got != "/system/m/1/file.site.middle.jpg" {
		t.Errorf("no crop of the format: the field's, %s", got)
	}
	box.Format = ""
	if got := box.URL("site.full"); got != "/system/m/1/file.site.full.jpg" {
		t.Errorf("no format: %s", got)
	}
}
