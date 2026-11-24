package media

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MediaBoxUsage is a field of images of a model — a use of MediaBox: its
// key, "<table>.<Field>", is the record_table and field of its uses
// (MediaLibraryUse) — and the sizes its crops are made in: the code's
// (DefaultSizes, written each time the application starts: RegisterUsages),
// or the admin's (Sizes, when it has some). The server reads them by the key
// (MediaBoxConfig.Usage), not as the browser sends them.
//
// When the sizes of a field change — the admin's saved, or the code's when
// the admin has none — the crops of its files are made again in them, in
// the background (Recrop): RecropState, RecropDone of RecropTotal,
// RecropError.
type MediaBoxUsage struct {
	Key         string `gorm:"column:usage_key;primaryKey;size:255"`
	RecordTable string `gorm:"size:128;index"`
	Field       string `gorm:"size:128"`
	// Model is the model's name in the admin, the field's label: what a
	// person reads of it.
	Model string
	Label string
	// DefaultSizes are the code's, Sizes the admin's (none: the code's):
	// JSON of []UsageSize.
	DefaultSizes string `gorm:"type:text"`
	Sizes        string `gorm:"type:text"`
	// AppliedSizes are the sizes the crops were last made in (Recrop); none
	// when they never were here.
	AppliedSizes string `gorm:"type:text"`
	// Formats are the admin's formats of the field (JSON of []UsageFormat):
	// its sizes in other proportions, one chosen in a record's box
	// (MediaBox.Format) — a portrait cover among landscape ones.
	Formats string `gorm:"type:text"`
	// DefaultOriginalWidth × DefaultOriginalHeight is the code's limit of
	// the original of an image (MediaBoxConfig.OriginalLimit; 0: the media
	// library's), OriginalWidth × OriginalHeight the admin's (0: the code's):
	// a larger image is made smaller to fit it — in the browser before it is
	// sent, and checked by the server.
	DefaultOriginalWidth  int
	DefaultOriginalHeight int
	OriginalWidth         int
	OriginalHeight        int
	// DefaultThumbWidth × DefaultThumbHeight is the code's limit of the
	// thumbnail of an image (MediaBoxConfig.ThumbLimit; 0: none), ThumbWidth
	// × ThumbHeight the admin's (0: the code's).
	DefaultThumbWidth  int
	DefaultThumbHeight int
	ThumbWidth         int
	ThumbHeight        int
	RecropState        string `gorm:"size:16"`
	RecropTotal        int
	RecropDone         int
	RecropError        string `gorm:"type:text"`
	RecropAt           *time.Time
	UpdatedAt          time.Time
}

func (MediaBoxUsage) TableName() string { return "media_box_usages" }

// PrimarySlug is the record's key in the admin (its key is no "ID");
// PrimaryColumnValuesBySlug reads it back whole — a table's name has "_",
// which a slug of parts is split by.
func (u MediaBoxUsage) PrimarySlug() string { return u.Key }

func (MediaBoxUsage) PrimaryColumnValuesBySlug(slug string) map[string]string {
	return map[string]string{"Key": slug}
}

// The states of a recrop.
const (
	RecropRunning = "running"
	RecropDone    = "done"
	RecropFailed  = "error"
)

// UsageSize is a size of the crops of a field: its name (what the code and
// the templates ask the image for: MediaBox.URL(name)), its width and height.
type UsageSize struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// UsageKey is the key of the field of the table: "<table>.<Field>".
func UsageKey(table, field string) string { return table + "." + field }

// sizesList is sizes as a list, by their names.
func sizesList(sizes map[string]*base.Size) []UsageSize {
	out := make([]UsageSize, 0, len(sizes))
	for name, s := range sizes {
		if s != nil {
			out = append(out, UsageSize{Name: name, Width: s.Width, Height: s.Height})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func encodeSizes(l []UsageSize) string {
	if len(l) == 0 {
		return ""
	}
	b, _ := json.Marshal(l)
	return string(b)
}

func decodeSizes(s string) []UsageSize {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	var l []UsageSize
	if json.Unmarshal([]byte(s), &l) != nil {
		return nil
	}
	return l
}

// EffectiveSizes are the sizes the field's crops are made in: the admin's
// when it has some, else the code's.
func (u *MediaBoxUsage) EffectiveSizes() []UsageSize {
	if l := decodeSizes(u.Sizes); len(l) > 0 {
		return l
	}
	return decodeSizes(u.DefaultSizes)
}

// EffectiveOriginalMax is the limit of the original of an image of the field:
// the admin's when it has one, else the code's; nil, none of them.
func (u *MediaBoxUsage) EffectiveOriginalMax() *base.Size {
	if u.OriginalWidth > 0 && u.OriginalHeight > 0 {
		return &base.Size{Width: u.OriginalWidth, Height: u.OriginalHeight}
	}
	if u.DefaultOriginalWidth > 0 && u.DefaultOriginalHeight > 0 {
		return &base.Size{Width: u.DefaultOriginalWidth, Height: u.DefaultOriginalHeight}
	}
	return nil
}

// EffectiveThumbMax is the limit of the thumbnail of an image of the field:
// the admin's when it has one, else the code's; nil, none of them.
func (u *MediaBoxUsage) EffectiveThumbMax() *base.Size {
	if u.ThumbWidth > 0 && u.ThumbHeight > 0 {
		return &base.Size{Width: u.ThumbWidth, Height: u.ThumbHeight}
	}
	if u.DefaultThumbWidth > 0 && u.DefaultThumbHeight > 0 {
		return &base.Size{Width: u.DefaultThumbWidth, Height: u.DefaultThumbHeight}
	}
	return nil
}

// sizesOf is l as the sizes of a MediaBoxConfig — what code gave a size of
// the same name besides its width and height (Padding, Sm, Cols) kept.
func sizesOf(l []UsageSize, code map[string]*base.Size) map[string]*base.Size {
	out := make(map[string]*base.Size, len(l))
	for _, s := range l {
		size := &base.Size{Width: s.Width, Height: s.Height}
		if c := code[s.Name]; c != nil {
			size.Padding, size.Sm, size.Cols = c.Padding, c.Sm, c.Cols
		}
		out[s.Name] = size
	}
	return out
}

// usage is the use of the field key; nil when there is none.
func (b *Builder) usage(key string) *MediaBoxUsage {
	if key == "" {
		return nil
	}
	var u MediaBoxUsage
	if b.db.Where("usage_key = ?", key).Take(&u).Error != nil {
		return nil
	}
	return &u
}

// withUsage is cfg with the sizes and the limits of the original and of the
// thumbnail of its use
// (cfg.Usage): the ones the server trusts.
func (b *Builder) withUsage(cfg *media_library.MediaBoxConfig) *media_library.MediaBoxConfig {
	if cfg == nil || cfg.Usage == "" {
		return cfg
	}
	u := b.usage(cfg.Usage)
	if u == nil {
		return cfg
	}
	c := *cfg
	if l := u.EffectiveSizes(); len(l) > 0 {
		c.Sizes = sizesOf(l, cfg.Sizes)
	}
	// its formats, the one of the box: its sizes the box's
	c.Formats = nil
	for _, f := range u.UsageFormats() {
		c.Formats = append(c.Formats, media_library.FormatOption{Name: f.Name, Label: f.Label})
	}
	if f := u.format(c.Format); f != nil {
		if l := u.FormatSizes(*f); len(l) > 0 {
			c.Sizes = sizesOf(l, nil)
		}
	} else {
		c.Format = ""
	}
	c.OriginalMax = u.EffectiveOriginalMax()
	c.ThumbMax = u.EffectiveThumbMax()
	return &c
}

// cfgOf is the config a browser sent back (v, its JSON) with the sizes of its
// use read here.
func (b *Builder) cfgOf(v string) *media_library.MediaBoxConfig {
	return b.withUsage(stringToCfg(v))
}

// fieldUsage is the config of the field of the context as it is drawn: its
// use (the record's table, the field's name) — of a field of a record, not
// of one of a list or a form in it —, and its sizes.
func (b *Builder) fieldUsage(field *presets.FieldContext, cfg *media_library.MediaBoxConfig) *media_library.MediaBoxConfig {
	if cfg == nil || !imagesOf(cfg) || field == nil || field.Parent != nil || field.Obj == nil {
		return cfg
	}
	stmt := &gorm.Statement{DB: b.db}
	if stmt.Parse(field.Obj) != nil || stmt.Schema == nil {
		return cfg
	}
	c := *cfg
	c.Usage = UsageKey(stmt.Schema.Table, field.Name)
	return b.withUsage(&c)
}

// RegisterUsages writes the uses of MediaBox of the models of pb — each
// field of images of their editing (its MediaBoxConfig, sizes or not) —,
// their sizes and the limits of their original and of their thumbnail the
// code's (DefaultSizes, DefaultOriginal*, DefaultThumb*); the admin's kept. A field whose sizes now are not
// the ones its crops were made in (the code's changed, and the admin has none)
// has them made again (Recrop). Call it once the models are all registered.
func (b *Builder) RegisterUsages(pb *presets.Builder) error {
	for _, mb := range pb.Models() {
		eb := mb.Editing()
		if eb == nil {
			continue
		}
		stmt := &gorm.Statement{DB: b.db}
		if err := stmt.Parse(mb.NewModel()); err != nil || stmt.Schema == nil {
			continue
		}
		for _, f := range boxFields(stmt.Schema) {
			fb := eb.GetField(f.Name)
			if fb == nil {
				continue
			}
			// a box with no config: one of no sizes
			cfg := &media_library.MediaBoxConfig{}
			if ctx := *fb.ContextPtr(); ctx != nil {
				if c, _ := ctx.Value(MediaBoxConfig).(*media_library.MediaBoxConfig); c != nil {
					cfg = c
				}
			}
			if !imagesOf(cfg) {
				continue
			}
			var ow, oh, tw, th int
			if o := cfg.OriginalLimit(); o != nil {
				ow, oh = o.Width, o.Height
			}
			if t := cfg.ThumbLimit(); t != nil {
				tw, th = t.Width, t.Height
			}
			u := MediaBoxUsage{
				Key:          UsageKey(stmt.Schema.Table, f.Name),
				RecordTable:  stmt.Schema.Table,
				Field:        f.Name,
				Model:        mb.Info().URIName(),
				Label:        f.Name,
				DefaultSizes: encodeSizes(sizesList(cfg.Sizes)),
				UpdatedAt:    time.Now(),

				DefaultOriginalWidth:  ow,
				DefaultOriginalHeight: oh,
				DefaultThumbWidth:     tw,
				DefaultThumbHeight:    th,
			}
			b.usageModels.Store(u.Key, mb)
			if err := b.db.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "usage_key"}},
				DoUpdates: clause.AssignmentColumns([]string{"record_table", "field", "model", "label", "default_sizes",
					"default_original_width", "default_original_height", "default_thumb_width", "default_thumb_height", "updated_at"}),
			}).Create(&u).Error; err != nil {
				return fmt.Errorf("media: the sizes of %s: %w", u.Key, err)
			}
			if got := b.usage(u.Key); got != nil && got.AppliedSizes != "" &&
				got.AppliedSizes != encodeSizes(got.EffectiveSizes()) {
				b.StartRecrop(u.Key)
			}
		}
	}
	return nil
}

// sizesConflict is why the sizes l of the field key may not be: a name a size
// of another field has in other measures — the crops of a file used in both
// would be made in one, then the other.
func (b *Builder) sizesConflict(key string, l []UsageSize) string {
	var others []MediaBoxUsage
	b.db.Where("usage_key <> ?", key).Find(&others)
	for _, s := range l {
		for _, o := range others {
			for _, os := range o.EffectiveSizes() {
				if os.Name == s.Name && (os.Width != s.Width || os.Height != s.Height) {
					return fmt.Sprintf("%s: %s %dx%d", o.Key, os.Name, os.Width, os.Height)
				}
			}
		}
	}
	return ""
}

// centeredCrop is the largest area of the ratio w:h centered in an image of
// width × height.
func centeredCrop(width, height, w, h int) *base.CropOption {
	if width <= 0 || height <= 0 || w <= 0 || h <= 0 {
		return nil
	}
	r := float64(w) / float64(h)
	cw, ch := float64(width), float64(height)
	if cw/ch > r {
		cw = ch * r
	} else {
		ch = cw / r
	}
	return &base.CropOption{
		X: int(math.Round((float64(width) - cw) / 2)), Y: int(math.Round((float64(height) - ch) / 2)),
		Width: int(math.Round(cw)), Height: int(math.Round(ch)),
	}
}

// cropOfRatio is a copy of the area of crops whose ratio is w:h — of the one
// first by name, so the choice is always the same —; nil when none is.
func cropOfRatio(crops map[string]*base.CropOption, w, h int) *base.CropOption {
	names := make([]string, 0, len(crops))
	for name, c := range crops {
		if sameRatio(c, w, h) {
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return nil
	}
	sort.Strings(names)
	c := *crops[names[0]]
	return &c
}

// sameRatio says the area keeps the ratio w:h (to 1%).
func sameRatio(c *base.CropOption, w, h int) bool {
	if c == nil || c.Height == 0 || h == 0 {
		return false
	}
	a, b := float64(c.Width)/float64(c.Height), float64(w)/float64(h)
	return math.Abs(a-b)/b < 0.01
}

// usageModelBuilder is the model of the field key in the admin
// (RegisterUsages); nil when none.
func (b *Builder) usageModelBuilder(key string) *presets.ModelBuilder {
	m, _ := b.usageModels.Load(key)
	mb, _ := m.(*presets.ModelBuilder)
	return mb
}

// usageModel is a record of the model of the field key; nil when none.
func (b *Builder) usageModel(key string) any {
	if mb := b.usageModelBuilder(key); mb != nil {
		return mb.NewModel()
	}
	return nil
}

// usageNames are the names the admin gives the model and the field of u —
// what it has, as registered, when the admin has no model of them.
func (b *Builder) usageNames(u *MediaBoxUsage, ctx *web.EventContext) (model, field string) {
	model, field = u.Model, u.Label
	mb := b.usageModelBuilder(u.Key)
	if mb == nil {
		return
	}
	model = mb.TTitlePlural(ctx.Context())
	if fb := mb.Editing().GetField(u.Field); fb != nil {
		if l := mb.FieldLabel(fb, ctx); l != "" {
			field = l
		}
	}
	return
}

// uploadConfig is the config an image uploaded to a box of cfg is saved by:
// the media library's, its original's limit the box's (OriginalLimit).
func (b *Builder) uploadConfig(cfg *media_library.MediaBoxConfig) *base.Config {
	c := b.Config
	if s := cfg.OriginalLimit(); s != nil {
		c.OriginalMaxWidth, c.OriginalMaxHeight = s.Width, s.Height
	}
	return &c
}

// imagesOf says the box of cfg takes images (its type an image's, or any).
func imagesOf(cfg *media_library.MediaBoxConfig) bool {
	return cfg.AllowType == "" || cfg.AllowType == media_library.ALLOW_TYPE_IMAGE
}

// UsageFormat is a format of a field of images: its name (the suffix of its
// sizes: FormatSizeName), its label, and the largest of its sizes — the
// field's others in its proportion, scaled by their widths (FormatSizes).
type UsageFormat struct {
	Name   string `json:"name"`
	Label  string `json:"label"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func encodeFormats(l []UsageFormat) string {
	if len(l) == 0 {
		return ""
	}
	b, _ := json.Marshal(l)
	return string(b)
}

func decodeFormats(s string) []UsageFormat {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	var l []UsageFormat
	if json.Unmarshal([]byte(s), &l) != nil {
		return nil
	}
	return l
}

// UsageFormats are the formats of the field.
func (u *MediaBoxUsage) UsageFormats() []UsageFormat { return decodeFormats(u.Formats) }

// FormatSizes are the sizes of the field in the format f: each of its sizes
// (EffectiveSizes) in the proportion of f, the largest f's own size, the
// others scaled as their widths are to its; named FormatSizeName.
func (u *MediaBoxUsage) FormatSizes(f UsageFormat) []UsageSize {
	sizes := u.EffectiveSizes()
	maxW := 0
	for _, s := range sizes {
		maxW = max(maxW, s.Width)
	}
	if maxW == 0 || f.Width <= 0 || f.Height <= 0 {
		return nil
	}
	out := make([]UsageSize, 0, len(sizes))
	for _, s := range sizes {
		out = append(out, UsageSize{
			Name:   media_library.FormatSizeName(s.Name, f.Name),
			Width:  int(math.Round(float64(f.Width) * float64(s.Width) / float64(maxW))),
			Height: int(math.Round(float64(f.Height) * float64(s.Width) / float64(maxW))),
		})
	}
	return out
}

// format is the format of the field named name; nil when it has none.
func (u *MediaBoxUsage) format(name string) *UsageFormat {
	for _, f := range u.UsageFormats() {
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
	for _, f := range u.UsageFormats() {
		l = append(l, u.FormatSizes(f)...)
	}
	return l
}
