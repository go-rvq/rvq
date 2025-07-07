package presets

import (
	"fmt"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/web/v3/vue"
	"github.com/qor5/x/v3/perm"
	. "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
)

func (b *ListingBuilder) DeleteFunc(v DeleteFunc) (r *ListingBuilder) {
	b.Deleter = v
	return b
}

func (b *ListingBuilder) WrapDeleteFunc(w func(in DeleteFunc) DeleteFunc) (r *ListingBuilder) {
	b.Deleter = w(b.Deleter)
	return b
}

func (b *ListingBuilder) doDelete(ctx *web.EventContext) (r web.EventResponse, err error) {
	pk := ctx.R.FormValue(ParamID)

	if pk == "" {
		ShowMessage(&r, MustGetMessages(ctx.Context()).ErrEmptyParamID.Error(), "warning")
		return
	}

	var (
		obj = b.mb.NewModel()
		id  ID
	)

	if id, err = b.mb.ParseRecordID(pk); err != nil {
		return
	}

	id.SetTo(obj)

	if !b.DeletingRestriction.CanObj(obj, ctx) {
		err = perm.PermissionDenied
		return
	}

	if len(pk) > 0 {
		if err = b.Deleter(obj, id, ctx.R.FormValue("cascade") == "true", ctx); err != nil {
			return
		}
	}

	ShowMessage(&r, MustGetMessages(ctx.Context()).SuccessfullyDeleted)

	web.AppendRunScripts(&r, "closer.show = false")

	if postDeleteConfig := ctx.R.URL.Query().Get(ParamPostDeleteCallback); postDeleteConfig != "" {
		cb := web.DecodeCallback(postDeleteConfig)
		web.AppendRunScripts(&r, web.ApplyChangeEvent(cb.Script(), web.Deleted, pk))
		r.ReloadPortals = append(r.ReloadPortals, cb.ReloadPortals...)
	}
	return
}

func (b *ListingBuilder) deleteConfirmation(ctx *web.EventContext) (r web.EventResponse, err error) {
	msgr := MustGetMessages(ctx.Context())
	id := ctx.R.FormValue(ParamID)
	targetPortal := ctx.R.FormValue(ParamTargetPortal)

	var (
		obj           = b.mb.NewModel()
		modelTitle    = b.mb.TTitle(ctx.Context())
		theModelTitle = b.mb.TTheTitle(ctx.Context())
		title         string
		ido           ID
	)

	if ido, err = b.mb.ParseRecordID(id); err != nil {
		return
	}

	ido.SetTo(obj)

	if title, err = b.mb.RecordTitleFetch(obj, ctx); err != nil {
		return
	}

	b.mb.p.Dialog().
		SetTargetPortal(targetPortal).
		Respond(ctx, &r, vue.UserComponent(
			vx.VXDialog().
				Density(DensityCompact).
				Style("max-width: 500px").
				SlotBody(
					VAlert(
						h.RawHTML(msgr.DeleteConfirmationHtml(modelTitle, theModelTitle, title)),
					).Type(TypeWarning).
						Variant(VariantTonal),
					h.If(false, VSwitch().
						Attr("v-model", []byte("cascade.value")).
						Color(ColorWarning).
						Label("Excluir relacionados").
						Hint("Exclui automaticamente todas entidades relacionadas").
						PersistentHint(true),
					),
				).
				ToolbarProps(fmt.Sprintf(`{color:%q}`, ColorError)).
				SlotBottom(h.HTMLComponents{
					VSpacer(),
					VBtn(msgr.Cancel).
						Variant(VariantFlat).
						Class("ml-2").
						On("click", "closer.show = false"),
					VBtn(msgr.Delete).
						Color(ColorError).
						Variant(VariantFlat).
						Theme(ThemeDark).
						Attr("@click", web.Plaid().
							EventFunc(actions.DoDelete).
							Queries(ctx.Queries()).
							Query("cascade", web.Var("cascade.value")).
							URL(ctx.R.URL.Path).
							Go()),
				}).
				Title(msgr.Delete),
		).ScopeVar("cascade", "{value: false}"))
	return
}
