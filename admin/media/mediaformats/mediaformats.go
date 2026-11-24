// Package mediaformats is the admin of the formats of the images
// (media.MediaFormat) and of the formats the fields of images offer
// (media.MediaBoxUsageFormat): two models, siblings of the sizes of the
// images (media-sizes), a format of a field choosing its field and its format
// by their selectors (helper.ModelSelectorBuilder).
//
// It is a package of its own because the selectors' package (helper) imports
// the media package: media does not import it back.
package mediaformats

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/media"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/packages/helper"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// The URIs of the models.
const (
	FormatsURIName      = "media-formats"
	UsageFormatsURIName = "media-usage-formats"
)

// Install registers the formats of the images and the formats of the fields
// in pb, beside the sizes of the images (media.Builder.Install first: the
// selector of a field lists them).
func Install(pb *presets.Builder, mb *media.Builder) error {
	usages := pb.GetModel(&media.MediaBoxUsage{})
	if usages == nil {
		return errors.New("mediaformats: the sizes of the images (media) are not installed")
	}
	db := mb.DB()
	cfg := func() *presets.ModelBuilderConfig {
		return presets.ModelConfig().SetModuleKey(media.I18nMediaLibraryKey)
	}
	formats := pb.Model(&media.MediaFormat{}, cfg()).URIName(FormatsURIName).MenuIcon("mdi-aspect-ratio")
	links := pb.Model(&media.MediaBoxUsageFormat{}, cfg()).URIName(UsageFormatsURIName).MenuIcon("mdi-crop-portrait")
	configureFormats(formats, mb, db)
	configureLinks(links, usages, mb, db)
	return nil
}

func msgs(ctx context.Context) *media.Messages { return media.GetMessages(ctx) }

func td(s string) h.HTMLComponent { return h.Td(h.Text(s)) }

func labeled(label, text string) h.HTMLComponent {
	return h.Div(h.Div(h.Text(label)).Class("text-caption text-medium-emphasis"), h.Div(h.Text(text))).Class("mb-4")
}

// usageTitle is how the field key reads: "<model> · <field>".
func usageTitle(mb *media.Builder, db *gorm.DB, key string, ctx *web.EventContext) string {
	var u media.MediaBoxUsage
	if db.Where("usage_key = ?", key).Take(&u).Error != nil {
		return key
	}
	model, field := mb.UsageNames(&u, ctx)
	return model + " · " + field
}

func configureFormats(m *presets.ModelBuilder, mb *media.Builder, db *gorm.DB) {
	ratio := func(f *media.MediaFormat) string { return fmt.Sprintf("%d:%d", f.Width, f.Height) }
	usagesText := func(f *media.MediaFormat, ctx *web.EventContext) string {
		var l []string
		for _, key := range mb.UsagesOfFormat(f.ID) {
			l = append(l, usageTitle(mb, db, key, ctx))
		}
		return strings.Join(l, ", ")
	}
	helper.SetToTitleRecordEncoderFactory(m.Listing(), func(ctx *web.EventContext) func(r *media.MediaFormat) string {
		return func(r *media.MediaFormat) string { return fmt.Sprintf("%s (%s)", r.Label, r.Name) }
	})

	l := m.Listing("Label", "Name", "Ratio", "Usages", "Code").OrderBy("name")
	l.Field("Ratio").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return td(ratio(field.Obj.(*media.MediaFormat)))
	})
	l.Field("Usages").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return td(usagesText(field.Obj.(*media.MediaFormat), ctx))
	})

	d := m.Detailing("Label", "Name", "Ratio", "Usages", "Code")
	d.Field("Ratio").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return labeled(field.Label, ratio(field.Obj.(*media.MediaFormat)))
	})
	d.Field("Usages").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return labeled(field.Label, usagesText(field.Obj.(*media.MediaFormat), ctx))
	})

	ed := m.Editing("Label", "Name", []string{"Width", "Height"})
	for _, f := range []string{"Label", "Name", "Width", "Height"} {
		ed.Field(f).Required(true)
	}
	// the name: what the crops and the boxes keep — read only once made
	ed.Field("Name").WrapComponentFunc(func(old presets.FieldComponentFunc) presets.FieldComponentFunc {
		return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			if f := field.Obj.(*media.MediaFormat); f.ID != uuid.Nil {
				return labeled(field.Label, f.Name)
			}
			return old(field, ctx)
		}
	})
	ed.Validators.AppendFunc(func(obj any, _ presets.FieldModeStack, ctx *web.EventContext) (errs web.ValidationErrors) {
		f := obj.(*media.MediaFormat)
		m := msgs(ctx.Context())
		if f.ID == uuid.Nil {
			f.Name = strings.TrimSpace(f.Name)
			var n int64
			db.Model(&media.MediaFormat{}).Where("name = ?", f.Name).Count(&n)
			switch {
			case !media.FormatNameValid(f.Name):
				errs.FieldError("Name", fmt.Sprintf(m.FormatNameInvalid, f.Name))
			case n > 0:
				errs.FieldError("Name", fmt.Sprintf(m.FormatNameTaken, f.Name))
			}
		}
		if f.Width <= 0 || f.Height <= 0 {
			errs.FieldError("Width", m.FormatNotPositive)
		}
		return
	})
	ed.WrapSaveFunc(func(in presets.SaveFunc) presets.SaveFunc {
		return func(obj any, id model.ID, ctx *web.EventContext) error {
			f := obj.(*media.MediaFormat)
			var before media.MediaFormat
			if f.ID != uuid.Nil && db.Take(&before, "id = ?", f.ID).Error == nil {
				// the name kept, and the code's mark
				f.Name, f.Code = before.Name, before.Code
			}
			if err := in(obj, id, ctx); err != nil {
				return err
			}
			// another proportion: the crops of the fields that offer it
			if before.ID != uuid.Nil && (before.Width != f.Width || before.Height != f.Height) {
				for _, key := range mb.UsagesOfFormat(f.ID) {
					mb.StartRecrop(key)
				}
			}
			return nil
		}
	})
	l.WrapDeleteFunc(func(in presets.DeleteFunc) presets.DeleteFunc {
		return func(obj any, id model.ID, cascade bool, ctx *web.EventContext) error {
			var f media.MediaFormat
			if err := db.Take(&f, "id = ?", id.String()).Error; err != nil {
				return err
			}
			m := msgs(ctx.Context())
			if f.Code {
				return errors.New(m.FormatOfCode)
			}
			if keys := mb.UsagesOfFormat(f.ID); len(keys) > 0 {
				var l []string
				for _, key := range keys {
					l = append(l, usageTitle(mb, db, key, ctx))
				}
				return fmt.Errorf(m.FormatInUse, strings.Join(l, ", "))
			}
			return in(obj, id, cascade, ctx)
		}
	})
}

func configureLinks(m *presets.ModelBuilder, usages *presets.ModelBuilder, mb *media.Builder, db *gorm.DB) {
	// a field in its selector: "<model> · <field>"
	helper.SetToTitleRecordEncoderFactory(usages.Listing(), func(ctx *web.EventContext) func(r *media.MediaBoxUsage) string {
		return func(r *media.MediaBoxUsage) string {
			model, field := mb.UsageNames(r, ctx)
			return model + " · " + field
		}
	})
	format := func(l *media.MediaBoxUsageFormat) *media.MediaFormat {
		if l.Format != nil {
			return l.Format
		}
		var f media.MediaFormat
		if db.Take(&f, "id = ?", l.FormatID).Error != nil {
			return nil
		}
		return &f
	}
	sizesText := func(l *media.MediaBoxUsageFormat) string {
		f := format(l)
		if f == nil {
			return ""
		}
		var parts []string
		for _, s := range mb.FormatSizesOf(l.UsageID, f, l.Width) {
			parts = append(parts, fmt.Sprintf("%s %d×%d", s.Name, s.Width, s.Height))
		}
		return strings.Join(parts, ", ")
	}
	text := func(field *presets.FieldContext, ctx *web.EventContext) string {
		l := field.Obj.(*media.MediaBoxUsageFormat)
		switch field.Name {
		case "Usage":
			return usageTitle(mb, db, l.UsageID, ctx)
		case "Format":
			if f := format(l); f != nil {
				return fmt.Sprintf("%s (%s)", f.Label, f.Name)
			}
		case "Width":
			if l.Width > 0 {
				return strconv.Itoa(l.Width)
			}
		case "Sizes":
			return sizesText(l)
		}
		return ""
	}

	ed := m.Editing("Usage", "Format", []string{"Position", "Width"})
	ed.Field("Usage").Required(true)
	ed.Field("Format").Required(true)
	helper.NewModelSelectorBuilder(m, "Usage").SetForeignModel(usages).
		SetConfigureSelector(func(_ *presets.FieldContext, input *vx.VXAdvancedSelectBuilder) { input.ItemText("Title") }).
		Build()
	helper.NewModelSelectorBuilder(m, "Format").
		SetConfigureSelector(func(_ *presets.FieldContext, input *vx.VXAdvancedSelectBuilder) { input.ItemText("Title") }).
		Build()

	// the listing and the detail: the texts of the field and of the format
	// (after the selectors, which draw their own)
	ls := m.Listing("Usage", "Format", "Position", "Width", "Sizes").OrderBy("usage_key, position, created_at")
	d := m.Detailing("Usage", "Format", "Position", "Width", "Sizes")
	for _, f := range []string{"Usage", "Format", "Width", "Sizes"} {
		ls.Field(f).ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return td(text(field, ctx))
		})
		d.Field(f).ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return labeled(field.Label, text(field, ctx))
		})
	}

	ed.Validators.AppendFunc(func(obj any, _ presets.FieldModeStack, ctx *web.EventContext) (errs web.ValidationErrors) {
		l := obj.(*media.MediaBoxUsageFormat)
		m := msgs(ctx.Context())
		if l.UsageID == "" {
			errs.FieldError("Usage", m.FormatUsageRequired)
		}
		if l.FormatID == uuid.Nil {
			errs.FieldError("Format", m.FormatRequired)
		}
		if l.Width < 0 {
			errs.FieldError("Width", m.FormatWidthNegative)
		}
		if l.UsageID != "" && l.FormatID != uuid.Nil {
			var n int64
			db.Model(&media.MediaBoxUsageFormat{}).Where("usage_key = ? AND format_id = ? AND id <> ?",
				l.UsageID, l.FormatID, l.ID).Count(&n)
			if n > 0 {
				name := ""
				if f := format(l); f != nil {
					name = f.Name
				}
				errs.FieldError("Format", fmt.Sprintf(m.FormatTwice, name))
			}
		}
		return
	})
	// saved: the crops of its field in it; moved to another field, or to
	// another format, the boxes of the field before take the old one off
	ed.WrapSaveFunc(func(in presets.SaveFunc) presets.SaveFunc {
		return func(obj any, id model.ID, ctx *web.EventContext) error {
			l := obj.(*media.MediaBoxUsageFormat)
			var before media.MediaBoxUsageFormat
			if l.ID != uuid.Nil {
				db.Preload("Format").Take(&before, "id = ?", l.ID)
			}
			if err := in(obj, id, ctx); err != nil {
				return err
			}
			if before.ID != uuid.Nil && before.Format != nil &&
				(before.UsageID != l.UsageID || before.FormatID != l.FormatID) {
				if err := mb.DropBoxFormat(before.UsageID, before.Format.Name); err != nil {
					return err
				}
			}
			mb.StartRecrop(l.UsageID)
			return nil
		}
	})
	// deleted: the boxes of its field its sizes again
	ls.WrapDeleteFunc(func(in presets.DeleteFunc) presets.DeleteFunc {
		return func(obj any, id model.ID, cascade bool, ctx *web.EventContext) error {
			var link media.MediaBoxUsageFormat
			if err := db.Preload("Format").Take(&link, "id = ?", id.String()).Error; err != nil {
				return err
			}
			if err := in(obj, id, cascade, ctx); err != nil {
				return err
			}
			if link.Format != nil {
				return mb.DropBoxFormat(link.UsageID, link.Format.Name)
			}
			return nil
		}
	})
}
