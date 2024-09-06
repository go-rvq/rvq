package presets

import (
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	"github.com/qor5/x/v3/perm"
	h "github.com/theplant/htmlgo"

	. "github.com/qor5/x/v3/ui/vuetify"
)

func (b *EditingBuilder) formEdit(ctx *web.EventContext) (r web.EventResponse, err error) {
	if b.mb.Info().Verifier().Do(PermGet).WithReq(ctx.R).IsAllowed() != nil {
		ShowMessage(&r, perm.PermissionDenied.Error(), "warning")
		return
	}

	obj := b.mb.NewModel()
	id := ctx.R.Form.Get(ParamID)

	if len(id) > 0 || b.mb.singleton {
		if err = b.Fetcher(obj, id, ctx); err != nil {
			if err == ErrRecordNotFound && b.mb.singleton {
				err = nil
			} else {
				return
			}
		}
	} else {
		err = ErrRecordNotFound
		return
	}

	comp := b.form(obj, ctx).Component()
	mode := GetOverlay(ctx)
	targetPortal := ctx.R.FormValue(ParamTargetPortal)

	if mode.IsDrawer() {
		b.mb.p.Drawer(mode).
			SetValidPortalName(targetPortal).
			Respond(&r, comp)
	} else if mode.IsDialog() {
		b.mb.p.Dialog().
			SetScrollable(true).
			SetValidWidth(b.mb.rightDrawerWidth).
			SetTargetPortal(targetPortal).
			Respond(&r, comp)
	} else {
		r.Body = VContainer(comp)
	}

	return
}

func (b *EditingBuilder) EditBtn(ctx *web.EventContext, id string, edit bool, targetPortal string) *VBtnBuilder {
	var (
		queries = ctx.Queries()
		event   = actions.Create
	)

	if id != "" {
		queries.Set(ParamID, id)
	}

	if GetEditFormUnscoped(ctx) {
		queries.Set(ParamEditFormUnscoped, "true")
	}

	if edit {
		event = actions.Update
	}

	onClick := web.Plaid().
		EventFunc(event).
		Queries(queries).
		ValidQuery(ParamTargetPortal, targetPortal).
		URL(b.mb.Info().ListingHrefCtx(ctx))

	return VBtn("").
		Color("primary").
		Variant(VariantFlat).
		Attr(":disabled", "isFetching").
		Attr(":loading", "isFetching").
		Attr("@click", onClick.Go()).
		Icon(true).
		Density("comfortable").
		Children(VIcon("mdi-content-save"))
}

func (b *EditingBuilder) ConfigureForm(f *Form) *Form {
	var (
		disableUpdateBtn bool
		ctx              = f.b.ctx
		portalName       = ctx.R.FormValue(ParamTargetPortal)
	)

	f.Portal = portalName

	if f.b.mode == NEW {
		f.Title = f.b.msgr.CreatingObjectTitle(
			i18n.T(f.b.ctx.R, ModelsI18nModuleKey, b.mb.label),
			b.mb.female,
		)
	} else {
		disableUpdateBtn = !f.b.mb.Info().CanUpdate(f.b.ctx.R, f.Obj)

		editingTitleText := f.b.msgr.EditingObjectTitle(
			i18n.T(f.b.ctx.R, ModelsI18nModuleKey, b.mb.label),
			b.mb.RecordTitle(f.Obj, ctx))
		if b.editingTitleFunc != nil {
			f.Title = b.editingTitleFunc(f.b.obj, editingTitleText, f.b.ctx)
		} else {
			f.Title = editingTitleText
		}
	}

	if !disableUpdateBtn {
		f.PrimaryAction = b.EditBtn(f.b.ctx, f.b.id, f.b.mode != NEW, f.Portal)
	}

	if b.actionsFunc != nil {
		actions := b.actionsFunc(f.Obj, f.b.ctx)
		if comps, ok := actions.(h.HTMLComponents); ok {
			f.Actions = comps
		} else {
			f.Actions = append(f.Actions, actions)
		}
	}

	var hiddenComps []h.HTMLComponent
	for _, hf := range b.hiddenFuncs {
		hiddenComps = append(hiddenComps, hf(f.b.obj, f.b.ctx))
	}

	if len(hiddenComps) > 0 {
		f.Body = h.Components(
			h.Components(hiddenComps...),
			f.Body,
		)
	}

	return f
}

func (b *EditingBuilder) form(obj interface{}, ctx *web.EventContext) *Form {
	return b.ConfigureForm(NewFormBuilder(ctx, b.mb, &b.FieldsBuilder, obj).Build())
}
