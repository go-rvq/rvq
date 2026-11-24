package media

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MediaFormat is a format of the images of the site: a proportion with a
// name — a portrait among landscape covers. The fields of images that offer
// it are its MediaBoxUsageFormat; a record's image is cropped in the one
// chosen in its box (MediaBox.Format), its sizes the field's in the
// proportion (UsageFormat.Sizes).
type MediaFormat struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// Name is the suffix of the names of its crops (FormatSizeName) and what
	// the boxes keep: not changed once made.
	Name  string `gorm:"size:64;uniqueIndex"`
	Label string
	// Width:Height is its proportion.
	Width  int
	Height int
	// Code says the code declared it (Builder.Formats): made when missing,
	// never deleted.
	Code bool
}

func (MediaFormat) TableName() string { return "media_formats" }

// String is how the format reads: its label and its name.
func (f MediaFormat) String() string {
	if f.Name == "" {
		return f.Label
	}
	return f.Label + " (" + f.Name + ")"
}

// MediaBoxUsageFormat is a format a field of images offers: the field
// (MediaBoxUsage.Key), the format, its place among the field's formats, and
// the width of its largest size in the field — 0: the sizes by the rule of
// UsageFormat.Sizes.
type MediaBoxUsageFormat struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FormatID  uuid.UUID    `gorm:"type:uuid;uniqueIndex:media_box_usage_format,priority:2"`
	Format    *MediaFormat `gorm:"constraint:OnDelete:RESTRICT"`
	// UsageID is the field's key (MediaBoxUsage.Key)
	UsageID  string         `gorm:"column:usage_key;size:255;uniqueIndex:media_box_usage_format,priority:1"`
	Usage    *MediaBoxUsage `gorm:"foreignKey:UsageID;references:Key;constraint:OnDelete:CASCADE"`
	Position int
	Width    int
}

func (MediaBoxUsageFormat) TableName() string { return "media_box_usage_formats" }

// String is how it reads: its field's key and its format's name.
func (l MediaBoxUsageFormat) String() string {
	if l.Format == nil {
		return l.UsageID
	}
	return l.UsageID + " · " + l.Format.Name
}

// formatNameRe is a format's name: the suffix of its sizes' names.
var formatNameRe = regexp.MustCompile(`^[a-z0-9-]+$`)

// UsageFormat is a format as a field offers it: its name, its label, its
// proportion RatioW:RatioH, and the width of its largest size in the field
// (0: by the rule).
type UsageFormat struct {
	Name   string
	Label  string
	RatioW int
	RatioH int
	Width  int
}

// Sizes are the sizes of the field (sizes) in the format, named
// FormatSizeName. With a Width, the largest is Width wide and the others as
// their widths are to the largest's; without, the shorter side of each size
// is the shorter side of the format's — 1920×1080 in 3:4 is 1080×1440, in
// 1:1 1080×1080.
func (f UsageFormat) Sizes(sizes []UsageSize) []UsageSize {
	if f.RatioW <= 0 || f.RatioH <= 0 {
		return nil
	}
	maxW := 0
	for _, s := range sizes {
		maxW = max(maxW, s.Width)
	}
	ratio := float64(f.RatioH) / float64(f.RatioW)
	out := make([]UsageSize, 0, len(sizes))
	for _, s := range sizes {
		var w, h float64
		switch {
		case f.Width > 0 && maxW > 0:
			w = float64(f.Width) * float64(s.Width) / float64(maxW)
			h = w * ratio
		case f.RatioW <= f.RatioH:
			w = float64(min(s.Width, s.Height))
			h = w * ratio
		default:
			h = float64(min(s.Width, s.Height))
			w = h / ratio
		}
		if w < 1 || h < 1 {
			continue
		}
		out = append(out, UsageSize{Name: media_library.FormatSizeName(s.Name, f.Name),
			Width: int(math.Round(w)), Height: int(math.Round(h))})
	}
	return out
}

// UsageFormats are the formats the field offers, in their order.
func (u *MediaBoxUsage) UsageFormats() []UsageFormat { return u.formats }

// FormatSizes are the sizes of the field in the format f.
func (u *MediaBoxUsage) FormatSizes(f UsageFormat) []UsageSize { return f.Sizes(u.EffectiveSizes()) }

// format is the format of the field named name; nil when it offers none.
func (u *MediaBoxUsage) format(name string) *UsageFormat {
	for _, f := range u.formats {
		if f.Name == name {
			return &f
		}
	}
	return nil
}

// AllSizes are the sizes of the field and of each of its formats: what the
// crops of its files are made in (Recrop).
func (u *MediaBoxUsage) AllSizes() []UsageSize {
	l := u.EffectiveSizes()
	for _, f := range u.formats {
		l = append(l, u.FormatSizes(f)...)
	}
	return l
}

// usageFormats are the formats the field key offers, in their order.
func (b *Builder) usageFormats(key string) []UsageFormat {
	var links []MediaBoxUsageFormat
	if b.db.Preload("Format").Where("usage_key = ?", key).Order("position, created_at").Find(&links).Error != nil {
		return nil
	}
	out := make([]UsageFormat, 0, len(links))
	for _, l := range links {
		if l.Format == nil {
			continue
		}
		out = append(out, UsageFormat{Name: l.Format.Name, Label: l.Format.Label,
			RatioW: l.Format.Width, RatioH: l.Format.Height, Width: l.Width})
	}
	return out
}

// Formats declares formats of the code: each made when the site has none of
// its name (the admin's changes kept), and never deleted. Made by
// RegisterUsages.
func (b *Builder) Formats(formats ...MediaFormat) *Builder {
	b.codeFormats = append(b.codeFormats, formats...)
	return b
}

// registerFormats makes the formats of the code the site has not.
func (b *Builder) registerFormats() error {
	for _, f := range b.codeFormats {
		if !formatNameRe.MatchString(f.Name) || f.Width <= 0 || f.Height <= 0 {
			return fmt.Errorf("media: the format %q: a name of [a-z0-9-] and a proportion", f.Name)
		}
		f.Code = true
		f.ID = uuid.Nil
		if err := b.db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "name"}},
			DoUpdates: clause.Assignments(map[string]any{"code": true})}).Create(&f).Error; err != nil {
			return fmt.Errorf("media: the format %s: %w", f.Name, err)
		}
	}
	return nil
}

// UsageFormatsOf are the formats the field key offers, in their order.
func (b *Builder) UsageFormatsOf(key string) []UsageFormat { return b.usageFormats(key) }

// FormatNameValid says name may be a format's: lower-case letters, digits
// and "-".
func FormatNameValid(name string) bool { return formatNameRe.MatchString(name) }

// UsagesOfFormat are the keys of the fields that offer the format id.
func (b *Builder) UsagesOfFormat(id uuid.UUID) (keys []string) {
	b.db.Model(&MediaBoxUsageFormat{}).Where("format_id = ?", id).Order("usage_key").Pluck("usage_key", &keys)
	return
}

// DropBoxFormat takes the format name off the boxes of the records of the
// field key: their images the field's sizes again — the format no longer
// offered there.
func (b *Builder) DropBoxFormat(key, name string) error {
	u := b.usage(key)
	if u == nil {
		return nil
	}
	model := b.usageModel(key)
	if model == nil {
		return nil
	}
	stmt := &gorm.Statement{DB: b.db}
	if err := stmt.Parse(model); err != nil {
		return err
	}
	f := stmt.Schema.LookUpField(u.Field)
	if f == nil || len(stmt.Schema.PrimaryFields) == 0 {
		return nil
	}
	var ids []string
	if err := b.db.Model(&MediaLibraryUse{}).Where("record_table = ? AND field = ?", u.RecordTable, u.Field).
		Distinct().Pluck("record_id", &ids).Error; err != nil {
		return err
	}
	where := func(q *gorm.DB, parts []string) *gorm.DB {
		for i, pf := range stmt.Schema.PrimaryFields {
			q = q.Where(pf.DBName+" = ?", parts[i])
		}
		return q
	}
	for _, id := range ids {
		parts := strings.Split(id, ",")
		if len(parts) != len(stmt.Schema.PrimaryFields) {
			continue
		}
		var raw string
		if err := where(b.db.Table(stmt.Schema.Table), parts).Select(f.DBName).Row().Scan(&raw); err != nil {
			continue // a record gone meanwhile
		}
		var box media_library.MediaBox
		if box.Scan(raw) != nil || box.Format != name {
			continue
		}
		box.Format = ""
		if err := where(b.db.Table(stmt.Schema.Table), parts).Update(f.DBName, box).Error; err != nil {
			return err
		}
	}
	return nil
}

// FormatSizesOf are the sizes of the field key in the format f, the width of
// its largest width (0: by the rule).
func (b *Builder) FormatSizesOf(key string, f *MediaFormat, width int) []UsageSize {
	var u MediaBoxUsage
	if f == nil || b.db.Where("usage_key = ?", key).Take(&u).Error != nil {
		return nil
	}
	return UsageFormat{Name: f.Name, Label: f.Label, RatioW: f.Width, RatioH: f.Height, Width: width}.Sizes(u.EffectiveSizes())
}
