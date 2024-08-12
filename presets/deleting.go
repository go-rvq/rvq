package presets

import (
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/perm"
	. "github.com/qor5/x/v3/ui/vuetify"
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

func (b *ListingBuilder) doDelete(ctx *web.EventContext) (r web.EventResponse, err1 error) {
	if !b.mb.Info().CanDelete(ctx.R) {
		ShowMessage(&r, perm.PermissionDenied.Error(), "warning")
		return
	}

	id := ctx.R.FormValue(ParamID)
	obj := b.mb.NewModel()
	if len(id) > 0 {
		err := b.Deleter(obj, id, ctx)
		if err != nil {
			ShowMessage(&r, err.Error(), "warning")
			return
		}
	}

	ShowMessage(&r, MustGetMessages(ctx.R).SuccessfullyDeleted, "")

	web.AppendRunScripts(&r, "closer.show = false")

	if postSaveConfig := ctx.R.URL.Query().Get(ParamPostChangeCallback); postSaveConfig != "" {
		cb := web.DecodeCallback(postSaveConfig)
		web.AppendRunScripts(&r, web.ApplyChangeEvent(cb.Script(), web.Deleted, id))
		r.ReloadPortals = append(r.ReloadPortals, cb.ReloadPortals...)
	}
	return
}

func (b *ListingBuilder) deleteConfirmation(ctx *web.EventContext) (r web.EventResponse, err error) {
	msgr := MustGetMessages(ctx.R)
	id := ctx.R.FormValue(ParamID)
	targetPortal := ctx.R.FormValue(ParamTargetPortal)

	var (
		obj           = b.mb.NewModel()
		modelTitle    = b.mb.TTitle(ctx.R)
		theModelTitle = b.mb.TTheTitle(ctx.R)
		title         string
		ido           ID
	)

	if ido, err = b.mb.ParseID(id); err != nil {
		return
	}

	ido.SetTo(obj)

	if title, err = b.mb.RecordTitleFetch(obj, ctx); err != nil {
		return
	}

	Dialog(targetPortal).
		Wrap(func(comp *VDialogBuilder) {
			comp.MaxWidth("600px")
		}).
		Respond(&r, VCard(
			VCardTitle(h.Text(msgr.DeleteConfirmationText(modelTitle, theModelTitle, title))),
			VCardActions(
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
						URL(ctx.R.URL.Path).
						Go()),
			),
		))
	return
}
