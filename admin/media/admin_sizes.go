package media

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/fields/schemaform"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/x/perm"
)

// ActionRecrop is the action that makes the crops of a field's files again.
const ActionRecrop = "recrop"

const eventRecrop = "mediaSizes_recrop"

// sizesSchema is the form of the sizes of a field: a table of a name, a
// width and a height.
const sizesSchema = `[list, layout={name: "table", columns: [#name, #width, #height]}]
class Form {
	name str
	[min=1] width int
	[min=1] height int
}`

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
	l := m.Listing("Model", "Label", "Sizes", "RecropState")
	l.Field("Model").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		model, _ := names(field, ctx)
		return h.Td(h.Text(model))
	})
	l.Field("Label").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		_, label := names(field, ctx)
		return h.Td(h.Text(label))
	})
	l.Field("Sizes").ComponentFunc(summary)
	l.Field("RecropState").ComponentFunc(status)

	d := m.Detailing("Model", "Label", "Sizes", "RecropState", "RecropError")
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

	ed := m.Editing("Sizes")
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
		if c := mb.sizesConflict(u.Key, u.EffectiveSizes()); c != "" {
			errs.FieldError("Sizes", fmt.Sprintf(msgs.SizesConflict, c))
		}
		return
	})
	ed.SaveFunc(func(obj any, _ model.ID, ctx *web.EventContext) error {
		u := obj.(*MediaBoxUsage)
		before := mb.usage(u.Key)
		if err := mb.db.Model(&MediaBoxUsage{}).Where("usage_key = ?", u.Key).
			Update("sizes", u.Sizes).Error; err != nil {
			return err
		}
		// sizes of other measures: the crops made again in them
		if before == nil || encodeSizes(before.EffectiveSizes()) != encodeSizes(u.EffectiveSizes()) {
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
