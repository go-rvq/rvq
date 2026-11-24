package media

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/fields/schemaform"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// ActionRecrop is the action that makes the crops of a field's files again.
const ActionRecrop = "recrop"

const eventRecrop = "mediaSizes_recrop"

// sizesSchema is the form of the sizes of a field: a table of a name, a
// width and a height.
const sizesSchema = `class Size {
	name str
	[min=1] width int
	[min=1] height int
}

[layout={name: "table", columns: [#name, #width, #height]}]
export type default []Size`

// formatsSchema is the form of the formats of a field: a table of a name, a
// label, a width and a height (the largest size of the format).
const formatsSchema = `class Format {
	name str
	label str
	[min=1] width int
	[min=1] height int
}

[layout={name: "table", columns: [#name, #label, #width, #height]}]
export type default []Format`

// formatNameRe is a format's name: the suffix of its sizes' names.
var formatNameRe = regexp.MustCompile(`^[a-z0-9-]+$`)

// configSizes is the admin of the sizes of the fields of images
// (MediaBoxUsage): their list — the code's or the admin's, how their crops
// are —, the form of a field's sizes, and the action that makes its crops
// again. The fields are the code's (RegisterUsages): none made nor deleted
// here.
func configSizes(pb *presets.Builder, mb *Builder) {
	sf := schemaform.New()
	schema, err := sf.Parse(sizesSchema)
	if err != nil {
		panic(err)
	}
	fschema, err := sf.Parse(formatsSchema)
	if err != nil {
		panic(err)
	}
	m := pb.Model(&MediaBoxUsage{}, presets.ModelConfig().SetModuleKey(I18nMediaLibraryKey)).
		MenuIcon("mdi-crop").
		URIName("media-sizes")
	m.SetCreatingDisabled(true).SetDeletingDisabled(true)

	summary := func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Td(h.Text(sizesText(GetMessages(ctx.Context()), field.Obj.(*MediaBoxUsage))))
	}
	status := func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Td(h.Text(recropStatus(GetMessages(ctx.Context()), field.Obj.(*MediaBoxUsage))))
	}

	names := func(field *presets.FieldContext, ctx *web.EventContext) (string, string) {
		return mb.usageNames(field.Obj.(*MediaBoxUsage), ctx)
	}
	original := func(field *presets.FieldContext, ctx *web.EventContext) string {
		return originalText(GetMessages(ctx.Context()), field.Obj.(*MediaBoxUsage))
	}
	thumb := func(field *presets.FieldContext, ctx *web.EventContext) string {
		return thumbText(GetMessages(ctx.Context()), field.Obj.(*MediaBoxUsage))
	}
	l := m.Listing("Model", "Label", "Sizes", "Original", "Thumb", "RecropState")
	l.Field("Model").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		model, _ := names(field, ctx)
		return h.Td(h.Text(model))
	})
	l.Field("Label").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		_, label := names(field, ctx)
		return h.Td(h.Text(label))
	})
	l.Field("Sizes").ComponentFunc(summary)
	l.Field("Original").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Td(h.Text(original(field, ctx)))
	})
	l.Field("Thumb").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Td(h.Text(thumb(field, ctx)))
	})
	l.Field("RecropState").ComponentFunc(status)

	d := m.Detailing("Model", "Label", "Sizes", "Formats", "Original", "Thumb", "RecropState", "RecropError")
	labeled := func(label, text string) h.HTMLComponent {
		return h.Div(h.Div(h.Text(label)).Class("text-caption text-medium-emphasis"), h.Div(h.Text(text))).Class("mb-4")
	}
	d.Field("Model").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		model, _ := names(field, ctx)
		return labeled(field.Label, model)
	})
	d.Field("Label").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		_, label := names(field, ctx)
		return labeled(field.Label, label)
	})
	d.Field("Sizes").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Div(h.Div(h.Text(field.Label)).Class("text-caption text-medium-emphasis"),
			h.Div(h.Text(sizesText(GetMessages(ctx.Context()), field.Obj.(*MediaBoxUsage))))).Class("mb-4")
	})
	// the formats: each its label, its name and its largest size
	d.Field("Formats").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		var l []string
		for _, f := range field.Obj.(*MediaBoxUsage).UsageFormats() {
			l = append(l, fmt.Sprintf("%s (%s) %d×%d", f.Label, f.Name, f.Width, f.Height))
		}
		return labeled(field.Label, strings.Join(l, ", "))
	})
	d.Field("Original").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return labeled(field.Label, original(field, ctx))
	})
	d.Field("Thumb").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return labeled(field.Label, thumb(field, ctx))
	})
	d.Field("RecropState").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Div(h.Div(h.Text(field.Label)).Class("text-caption text-medium-emphasis"),
			h.Div(h.Text(recropStatus(GetMessages(ctx.Context()), field.Obj.(*MediaBoxUsage))))).Class("mb-4")
	})
	d.Action(ActionRecrop).
		Icon("mdi-crop-rotate").
		SetI18nLabel(func(ctx context.Context) string { return GetMessages(ctx).Recrop }).
		OnClick(func(_ *web.EventContext, id string, _ any) string {
			return web.Plaid().EventFunc(eventRecrop).Query("id", id).Go()
		})
	m.RegisterEventFunc(eventRecrop, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		u := mb.usage(ctx.R.FormValue("id"))
		if u == nil {
			return r, fmt.Errorf("media: no field %q", ctx.R.FormValue("id"))
		}
		if m.Permissioner().ReqObjectActioner(ctx.R, u, presets.ActionPerm(ActionRecrop)).Denied() {
			return r, perm.PermissionDenied
		}
		msgs := GetMessages(ctx.Context())
		if mb.StartRecrop(u.Key) {
			ctx.Flash = msgs.RecropStarted
		} else {
			ctx.Flash = msgs.RecropBusy
		}
		r.Reload = true
		return
	})

	ed := m.Editing("Sizes", "Formats", []string{"OriginalWidth", "OriginalHeight"}, []string{"ThumbWidth", "ThumbHeight"})
	// the formats: the sizes in other proportions, one chosen in a record's
	// image (MediaBox.Format)
	ed.Field("Formats").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := field.Obj.(*MediaBoxUsage)
		value, _ := sf.DecodeValue(u.Formats)
		form := vue.UserComponent(sf.ComponentFunc(fschema)(field, ctx)).Assign("form", field.FormKey, fschema.Filled(value))
		var errs h.HTMLComponents
		for _, e := range field.Errors {
			errs = append(errs, h.Div(h.Text(e)).Class("text-error text-body-2 mb-2").Attr("data-formats-error", true))
		}
		return h.Div(
			errs, form,
			h.Div(h.Text(GetMessages(ctx.Context()).FormatsHint)).Class("text-caption text-medium-emphasis mt-1 mb-3"),
		).Attr("data-media-formats", true)
	}).SetterFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) error {
		value, err := sf.DecodeForm(ctx, fschema, ctx.R.Form, field.FormKey)
		if err != nil {
			return err
		}
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}
		var l []UsageFormat
		if err = json.Unmarshal(b, &l); err != nil {
			return err
		}
		obj.(*MediaBoxUsage).Formats = encodeFormats(l)
		return nil
	})
	// the limits of the original and of the thumbnail: the admin's; empty
	// (0), the code's — said as the hint
	for _, f := range []string{"OriginalWidth", "OriginalHeight", "ThumbWidth", "ThumbHeight"} {
		ed.Field(f).ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			u := field.Obj.(*MediaBoxUsage)
			msgs := GetMessages(ctx.Context())
			value := ""
			if n, _ := field.Value().(int); n > 0 {
				value = fmt.Sprint(n)
			}
			hint := msgs.OriginalOfTheCodeNone
			if u.DefaultOriginalWidth > 0 && u.DefaultOriginalHeight > 0 {
				hint = fmt.Sprintf(msgs.OriginalOfTheCode, u.DefaultOriginalWidth, u.DefaultOriginalHeight)
			}
			if strings.HasPrefix(f, "Thumb") {
				hint = msgs.ThumbOfTheCodeNone
				if u.DefaultThumbWidth > 0 && u.DefaultThumbHeight > 0 {
					hint = fmt.Sprintf(msgs.ThumbOfTheCode, u.DefaultThumbWidth, u.DefaultThumbHeight)
				}
			}
			return v.VTextField().Type("number").Attr("min", 1).Label(field.Label).Suffix("px").
				Hint(hint).PersistentHint(true).Variant(v.FieldVariantUnderlined).
				ErrorMessages(field.Errors...).
				Attr(web.VField(field.FormKey, value)...).Attr("data-original-field", f)
		}).SetterFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) error {
			n := 0
			if raw := strings.TrimSpace(ctx.R.FormValue(field.FormKey)); raw != "" {
				var err error
				if n, err = strconv.Atoi(raw); err != nil {
					return err
				}
			}
			u := obj.(*MediaBoxUsage)
			switch f {
			case "OriginalWidth":
				u.OriginalWidth = n
			case "OriginalHeight":
				u.OriginalHeight = n
			case "ThumbWidth":
				u.ThumbWidth = n
			default:
				u.ThumbHeight = n
			}
			return nil
		})
	}
	ed.Field("Sizes").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := field.Obj.(*MediaBoxUsage)
		stored := u.Sizes
		if decodeSizes(stored) == nil {
			stored = u.DefaultSizes // the code's to start from
		}
		value, _ := sf.DecodeValue(stored)
		form := vue.UserComponent(sf.ComponentFunc(schema)(field, ctx)).Assign("form", field.FormKey, schema.Filled(value))
		// what is wrong with the sizes as a whole (a name of another field's):
		// the table says only what is wrong with each cell
		var errs h.HTMLComponents
		for _, e := range field.Errors {
			errs = append(errs, h.Div(h.Text(e)).Class("text-error text-body-2 mb-2").Attr("data-sizes-error", true))
		}
		return h.Div(errs, form)
	}).SetterFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) error {
		value, err := sf.DecodeForm(ctx, schema, ctx.R.Form, field.FormKey)
		if err != nil {
			return err
		}
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}
		var l []UsageSize
		if err = json.Unmarshal(b, &l); err != nil {
			return err
		}
		u := obj.(*MediaBoxUsage)
		// the code's again: none of the admin's
		if encodeSizes(l) == u.DefaultSizes {
			l = nil
		}
		u.Sizes = encodeSizes(l)
		return nil
	})
	ed.Validators.AppendFunc(func(obj any, _ presets.FieldModeStack, ctx *web.EventContext) (errs web.ValidationErrors) {
		u := obj.(*MediaBoxUsage)
		msgs := GetMessages(ctx.Context())
		seen := map[string]bool{}
		for _, s := range decodeSizes(u.Sizes) {
			switch {
			case strings.TrimSpace(s.Name) == "":
				errs.FieldError("Sizes", msgs.SizeNameRequired)
			case seen[s.Name]:
				errs.FieldError("Sizes", fmt.Sprintf(msgs.SizeNameTwice, s.Name))
			case s.Width <= 0 || s.Height <= 0:
				errs.FieldError("Sizes", fmt.Sprintf(msgs.SizeNotPositive, s.Name))
			}
			seen[s.Name] = true
		}
		formats := map[string]bool{}
		for _, f := range u.UsageFormats() {
			switch {
			case !formatNameRe.MatchString(f.Name):
				errs.FieldError("Formats", fmt.Sprintf(msgs.FormatNameInvalid, f.Name))
			case formats[f.Name]:
				errs.FieldError("Formats", fmt.Sprintf(msgs.FormatTwice, f.Name))
			case f.Width <= 0 || f.Height <= 0:
				errs.FieldError("Formats", fmt.Sprintf(msgs.FormatNotPositive, f.Name))
			}
			formats[f.Name] = true
		}
		if u.OriginalWidth < 0 || u.OriginalHeight < 0 || (u.OriginalWidth > 0) != (u.OriginalHeight > 0) {
			errs.FieldError("OriginalWidth", msgs.OriginalBoth)
		}
		if u.ThumbWidth < 0 || u.ThumbHeight < 0 || (u.ThumbWidth > 0) != (u.ThumbHeight > 0) {
			errs.FieldError("ThumbWidth", msgs.OriginalBoth)
		}
		if c := mb.sizesConflict(u.Key, u.EffectiveSizes()); c != "" {
			errs.FieldError("Sizes", fmt.Sprintf(msgs.SizesConflict, c))
		}
		return
	})
	ed.SaveFunc(func(obj any, _ model.ID, ctx *web.EventContext) error {
		u := obj.(*MediaBoxUsage)
		before := mb.usage(u.Key)
		if err := mb.db.Model(&MediaBoxUsage{}).Where("usage_key = ?", u.Key).
			Updates(map[string]any{"sizes": u.Sizes, "formats": u.Formats, "original_width": u.OriginalWidth, "original_height": u.OriginalHeight,
				"thumb_width": u.ThumbWidth, "thumb_height": u.ThumbHeight}).Error; err != nil {
			return err
		}
		// sizes of other measures: the crops made again in them
		if before == nil || encodeSizes(before.AllSizes()) != encodeSizes(u.AllSizes()) {
			mb.StartRecrop(u.Key)
		}
		return nil
	})
}

// sizesText is the sizes of the field u as a person reads them — "slider
// 1200×480, slider.sm 900×360" —, and whether they are the admin's.
func sizesText(m *Messages, u *MediaBoxUsage) string {
	var parts []string
	for _, s := range u.EffectiveSizes() {
		parts = append(parts, fmt.Sprintf("%s %d×%d", s.Name, s.Width, s.Height))
	}
	text := strings.Join(parts, ", ")
	if decodeSizes(u.Sizes) != nil {
		text += " · " + m.SizesOfTheAdmin
	}
	return text
}

// originalText is the limit of the original of an image of the field u as a
// person reads it — "1920×1080" —, and whether it is the admin's; none, the
// media library's.
func originalText(m *Messages, u *MediaBoxUsage) string {
	o := u.EffectiveOriginalMax()
	if o == nil {
		return ""
	}
	text := fmt.Sprintf("%d×%d", o.Width, o.Height)
	if u.OriginalWidth > 0 && u.OriginalHeight > 0 {
		text += " · " + m.SizesOfTheAdmin
	}
	return text
}

// thumbText is the limit of the thumbnail of an image of the field u as a
// person reads it, and whether it is the admin's; none, empty.
func thumbText(m *Messages, u *MediaBoxUsage) string {
	t := u.EffectiveThumbMax()
	if t == nil {
		return ""
	}
	text := fmt.Sprintf("%d×%d", t.Width, t.Height)
	if u.ThumbWidth > 0 && u.ThumbHeight > 0 {
		text += " · " + m.SizesOfTheAdmin
	}
	return text
}

// recropStatus is how the crops of the field u are: being made (done of
// total), made, or what went wrong.
func recropStatus(m *Messages, u *MediaBoxUsage) string {
	switch u.RecropState {
	case RecropRunning:
		return fmt.Sprintf(m.RecropRunning, u.RecropDone, u.RecropTotal)
	case RecropDone:
		return fmt.Sprintf(m.RecropMade, u.RecropTotal)
	case RecropFailed:
		return m.RecropFailedText
	}
	return ""
}
