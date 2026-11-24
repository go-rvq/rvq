package media

import (
	"testing"

	"github.com/google/uuid"

	"github.com/go-rvq/rvq/admin/media/media_library"
)

// A format in a field: each of the field's sizes in its proportion — the
// shorter side of the size the shorter of the format's, or, with a width of
// its own, the largest that wide and the others as their widths are —, named
// with its suffix; a box of the format shows their crops, the field's when it
// has not one.
func TestFormatSizesAndURL(t *testing.T) {
	sizes := []UsageSize{
		{Name: "site.thumb", Width: 500, Height: 281},
		{Name: "site.middle", Width: 900, Height: 506},
		{Name: "site.full", Width: 1920, Height: 1080},
	}
	got := func(f UsageFormat) map[string][2]int {
		m := map[string][2]int{}
		for _, s := range f.Sizes(sizes) {
			m[s.Name] = [2]int{s.Width, s.Height}
		}
		return m
	}
	for _, c := range []struct {
		f    UsageFormat
		want map[string][2]int
	}{
		{UsageFormat{Name: "retrato", RatioW: 3, RatioH: 4},
			map[string][2]int{"site.full_retrato": {1080, 1440}, "site.middle_retrato": {506, 675}, "site.thumb_retrato": {281, 375}}},
		{UsageFormat{Name: "quadrado", RatioW: 1, RatioH: 1},
			map[string][2]int{"site.full_quadrado": {1080, 1080}, "site.thumb_quadrado": {281, 281}}},
		{UsageFormat{Name: "largo", RatioW: 21, RatioH: 9},
			map[string][2]int{"site.full_largo": {2520, 1080}}},
		{UsageFormat{Name: "retrato", RatioW: 3, RatioH: 4, Width: 1200},
			map[string][2]int{"site.full_retrato": {1200, 1600}, "site.middle_retrato": {563, 750}}},
	} {
		g := got(c.f)
		for k, w := range c.want {
			if g[k] != w {
				t.Errorf("%s %d:%d (%d): %s = %v, want %v", c.f.Name, c.f.RatioW, c.f.RatioH, c.f.Width, k, g[k], w)
			}
		}
	}

	u := &MediaBoxUsage{DefaultSizes: encodeSizes(sizes), formats: []UsageFormat{{Name: "retrato", RatioW: 3, RatioH: 4}}}
	if n := len(u.AllSizes()); n != 6 {
		t.Errorf("all the sizes: %d", n)
	}
	if u.format("retrato") == nil || u.format("quadrado") != nil {
		t.Error("the formats of the field: only the ones it offers")
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

// The formats of the code are made once (the admin's changes kept); a field
// offers the ones linked to it, in their order; a format taken off a field
// takes it off the boxes of its records, the boxes of other fields kept.
func TestFormatsOfFields(t *testing.T) {
	db, a, _ := usesTestDB(t)
	_, b := usesSizesApp(t, db)
	b.Formats(MediaFormat{Name: "retrato", Label: "Retrato 3:4", Width: 3, Height: 4},
		MediaFormat{Name: "quadrado", Label: "Quadrado 1:1", Width: 1, Height: 1})
	if err := b.registerFormats(); err != nil {
		t.Fatal(err)
	}
	// the admin's label kept when the code's are made again
	db.Model(&MediaFormat{}).Where("name = ?", "retrato").Update("label", "Vertical")
	if err := b.registerFormats(); err != nil {
		t.Fatal(err)
	}
	var formats []MediaFormat
	db.Order("name").Find(&formats)
	if len(formats) != 2 || formats[1].Label != "Vertical" || !formats[1].Code {
		t.Fatalf("the formats of the code: %+v", formats)
	}
	quadrado, retrato := formats[0], formats[1]

	key := UsageKey("uses_posts", "Cover")
	db.Create(&MediaBoxUsageFormat{UsageID: key, FormatID: retrato.ID, Position: 2})
	db.Create(&MediaBoxUsageFormat{UsageID: key, FormatID: quadrado.ID, Position: 1, Width: 600})
	u := b.usage(key)
	if u == nil || len(u.UsageFormats()) != 2 || u.UsageFormats()[0].Name != "quadrado" {
		t.Fatalf("the formats of the field, in their order: %+v", u)
	}
	// the field's size thumb 100×100: 600 wide in its own width
	if s := u.FormatSizes(*u.format("quadrado")); len(s) != 1 || s[0].Width != 600 || s[0].Height != 600 {
		t.Errorf("a width of its own: %+v", s)
	}
	if u.format("largo") != nil {
		t.Error("a format the field does not offer")
	}

	// two records of the format, one of no format
	box := boxOf(a)
	box.Format = "retrato"
	p1, p2, p3 := usesPost{ID: uuid.New(), Cover: box}, usesPost{ID: uuid.New(), Cover: box}, usesPost{ID: uuid.New(), Cover: boxOf(a)}
	for _, p := range []*usesPost{&p1, &p2, &p3} {
		if err := db.Create(p).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := b.DropBoxFormat(key, "retrato"); err != nil {
		t.Fatal(err)
	}
	var posts []usesPost
	db.Find(&posts)
	for _, p := range posts {
		if p.Cover.Format != "" || string(p.Cover.ID) != a.ID.String() {
			t.Errorf("the box after the format is taken off: %+v", p.Cover)
		}
	}
}
