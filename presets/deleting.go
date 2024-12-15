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
	pk := ctx.R.FormValue(ParamID)

	if pk == "" {
		ShowMessage(&r, MustGetMessages(ctx.Context()).ErrEmptyParamID.Error(), "warning")
		return
	}

	var (
		obj = b.mb.NewModel()
		id  ID
	)

	if id, err1 = b.mb.ParseRecordID(pk); err1 != nil {
		ShowMessage(&r, err1.Error(), "warning")
		err1 = nil
		return
	}

	id.SetTo(obj)

	if !b.DeletingRestriction.CanObj(obj, ctx) {
		ShowMessage(&r, perm.PermissionDenied.Error(), "warning")
		return
	}

	if len(pk) > 0 {
		err := b.Deleter(obj, id, ctx)
		if err != nil {
			ShowMessage(&r, err.Error(), "warning")
			return
		}
	}

	ShowMessage(&r, MustGetMessages(ctx.Context()).SuccessfullyDeleted, "")

	web.AppendRunScripts(&r, "closer.show = false")

	if postSaveConfig := ctx.R.URL.Query().Get(ParamPostChangeCallback); postSaveConfig != "" {
		cb := web.DecodeCallback(postSaveConfig)
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

	Dialog(targetPortal).
		Wrap(func(comp *VDialogBuilder) {
			comp.MaxWidth("600px")
		}).
		Respond(ctx, &r, VCard(
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
