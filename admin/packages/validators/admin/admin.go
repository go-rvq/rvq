// Package admin mounts the validators package into an rvq presets admin: a
// Validator CRUD where records are seeded by the application and only the Value
// (GAD algorithm), Doc (markdown) and Messages (per-language error texts)
// overrides are user-editable. Clearing an override and saving resets it to the
// registered default. The detail view renders the documentation (markdown) that
// explains the algorithm.
package admin

import (
	"context"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/helper/comps"
	vmessages "github.com/go-rvq/rvq/admin/packages/validators/messages"
	vmodels "github.com/go-rvq/rvq/admin/packages/validators/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/thirdpart/gorm/datatypes"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

func configureValidator(b *presets.Builder, db *gorm.DB) *presets.ModelBuilder {
	mb := Model(b, &vmodels.Validator{}, presets.ModelWithID(ValidatorModelID)).
		MenuIcon("mdi-check-decagram")

	// records are seeded by the application, not created/removed by users.
	mb.SetCreatingDisabled(true)
	mb.SetDeletingDisabled(true)

	mb.Listing("Name", "Description", "DefaultLang").
		SearchColumns("name", "description").
		OrderBy("name")

	ed := mb.Editing("Value", "Doc", "Messages")

	// Value: the GAD algorithm (big textarea). Empty uses the registered initial.
	ed.Field("Value").
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			m := vmessages.Get(ctx.Context())
			return presets.LongTextFieldComponentFunc(field, ctx).
				Rows(14).AutoGrow(true).Attr("hint", m.ValueHint).Attr("persistent-hint", true)
		})

	// Doc: markdown documentation (big textarea).
	ed.Field("Doc").
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			m := vmessages.Get(ctx.Context())
			return presets.LongTextFieldComponentFunc(field, ctx).
				Rows(10).AutoGrow(true).Attr("hint", m.DocHint).Attr("persistent-hint", true)
		})

	// Messages: {LANG: {KEY: VALUE}} edited with the JSONMap editor (the language
	// is the row key; the value cell holds the inner {key: text} as JSON). A field
	// validator enforces the shape.
	comps.JSONMapField(mb, ed.Field("Messages"), languageKeyComp).
		NameLabelFunc(func(ctx context.Context) string { return vmessages.Get(ctx).MessagesLang }).
		ValueLabelFunc(func(ctx context.Context) string { return vmessages.Get(ctx).MessagesValue }).
		Build()

	ed.Field("Messages").Validator(presets.FieldValidatorFunc(func(field *presets.FieldContext) (verr web.ValidationErrors) {
		val, _ := field.Value().(datatypes.NullJSONMap)
		if err := vmodels.ValidateMessagesFormat(val); err != nil {
			verr.FieldError(field.Name, vmessages.Get(field.EventContext.Context()).ErrMessagesFormat)
		}
		return
	}))

	configureDetailing(mb)
	return mb
}

// languageKeyComp renders the row-key cell (the language code) as a text field.
func languageKeyComp(_ *presets.FieldContext, nameKey string) h.HTMLComponent {
	return v.VTextField().Density("compact").HideDetails(true).Attr("v-model", "item."+nameKey)
}

// configureDetailing renders the detail view: the metadata, the effective
// algorithm and the documentation (markdown) explaining it.
func configureDetailing(mb *presets.ModelBuilder) {
	dt := mb.Detailing("Name", "Description", "DefaultLang", "Algorithm", "Doc")

	dt.Field("Algorithm").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		val := field.Obj.(*vmodels.Validator)
		m := vmessages.Get(ctx.Context())
		return h.Div(
			v.VLabel(h.Text(m.Algorithm)),
			h.Pre(val.EffectiveValue()).
				Style("background:rgba(0,0,0,.05);padding:.75rem;border-radius:8px;overflow:auto;font-family:ui-monospace,monospace"),
		).Class("mb-4")
	})

	dt.Field("Doc").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		val := field.Obj.(*vmodels.Validator)
		m := vmessages.Get(ctx.Context())
		doc := val.EffectiveDoc()
		if doc == "" {
			return nil
		}
		return h.Div(
			v.VLabel(h.Text(m.Doc)),
			h.Div(h.RawHTML(mdToHTML(doc))).Class("markdown-body"),
		).Class("mb-4")
	})
}
