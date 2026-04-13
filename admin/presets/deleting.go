package presets

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
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
		obj     = b.mb.NewModel()
		id      ID
		related = ctx.R.FormValue("cascade") == "true"
	)

	if id, err = b.mb.ParseRecordID(pk); err != nil {
		return
	}

	id.SetTo(obj)

	if related {
		if !b.DeletingWithRelatedRestriction.CanObj(obj, ctx) {
			err = perm.PermissionDenied
			return
		}
	} else if !b.DeletingRestriction.CanObj(obj, ctx) {
		err = perm.PermissionDenied
		return
	}

	if len(pk) > 0 {
		if err = b.Deleter(obj, id, related, ctx); err != nil {
			return
		}
	}

	ShowMessage(&r, MustGetMessages(ctx.Context()).SuccessfullyDeleted)

	web.AppendRunScripts(&r, "closer.show = false; presetsListing?.loader?.go()")

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
		relatedComp   h.HTMLComponent
	)

	if ido, err = b.mb.ParseRecordID(id); err != nil {
		return
	}

	ido.SetTo(obj)

	if b.relatedDeletionConfig != nil {
		rdCtx := &RelatedDeletionContext{
			Builder: b,
			Obj:     obj,
			Context: ctx,
			FormKey: "cascade.value",
		}

		if err = b.relatedDeletionConfig(rdCtx); err != nil {
			return
		}

		if rdCtx.Enabled {
			if !b.DeletingWithRelatedRestriction.CanObj(obj, ctx) {
				err = perm.PermissionDenied
				return
			}

			hint := rdCtx.Description

			if hint == "" {
				hint = msgr.DeleteRelatedHint
			}

			relatedComps := h.HTMLComponents{
				VSwitch().
					Attr("v-model", []byte(rdCtx.FormKey)).
					Color(ColorWarning).
					Label(msgr.DeleteRelated).
					Hint(hint).
					PersistentHint(true),
			}

			if !rdCtx.ShowRelatedDisabled && b.showRelatedItensForDeletionActionFunc != nil {
				portalID := ctx.UID()
				relatedComps = append(relatedComps,
					web.Portal().Name(portalID),
					VBtn(msgr.ShowRelatedItemsTitle).
						Color(ColorInfo).
						Size(SizeSmall).
						Variant(VariantText).
						Attr("@click", web.Plaid().
							EventFunc(actions.ShowRelatedItensForDeletion).
							Query(ParamTargetPortal, portalID).
							Query(ParamID, id).
							URL(ctx.R.URL.Path).
							Go()))
			}

			relatedComp = VCard(VCardText(relatedComps))

			if rdCtx.WrapComponent != nil {
				relatedComp = rdCtx.WrapComponent(relatedComp)
			}
		} else if !b.DeletingRestriction.CanObj(obj, ctx) {
			err = perm.PermissionDenied
			return
		}
	} else if !b.DeletingRestriction.CanObj(obj, ctx) {
		err = perm.PermissionDenied
		return
	}

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
					relatedComp,
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
