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
			return
		}
	} else {
		err = ErrRecordNotFound
		return
	}

	comp := b.form(obj, ctx).FullComponent()
	mode := GetOverlay(ctx)
	targetPortal := ctx.R.FormValue(ParamTargetPortal)

	if mode.IsDrawer() {
		b.mb.p.Drawer(mode).
			SetValidPortalName(targetPortal).
			Respond(&r, comp)
	} else if mode.IsDialog() {
		b.mb.p.Dialog().
			SetValidWidth(b.mb.rightDrawerWidth).
			SetTargetPortal(targetPortal).
			Respond(&r, comp)
	} else {
		r.Body = comp
	}

	return
}

func (b *EditingBuilder) ConfigureForm(f *Form) *Form {
	var (
		event            = actions.Update
		disableUpdateBtn bool
		ctx              = f.b.ctx
		portalName       = ctx.R.FormValue(ParamTargetPortal)
	)

	f.Portal = portalName

	if f.b.mode == NEW {
		event = actions.Create
		f.Title = h.Text(f.b.msgr.CreatingObjectTitle(
			i18n.T(f.b.ctx.R, ModelsI18nModuleKey, b.mb.label),
			b.mb.female,
		))
	} else {
		disableUpdateBtn = !f.b.mb.Info().CanUpdate(f.b.ctx.R, f.Obj)

		editingTitleText := f.b.msgr.EditingObjectTitle(
			i18n.T(f.b.ctx.R, ModelsI18nModuleKey, b.mb.label),
			getPageTitle(f.b.obj, f.b.id))
		if b.editingTitleFunc != nil {
			f.Title = b.editingTitleFunc(f.b.obj, editingTitleText, f.b.ctx)
		} else {
			f.Title = h.Text(editingTitleText)
		}
	}

	queries := f.b.ctx.Queries()

	if f.b.id != "" {
		queries.Set(ParamID, f.b.id)
	}

	uri := b.mb.Info().ListingHrefCtx(f.b.ctx)
	onClick := web.Plaid().
		EventFunc(event).
		Queries(queries).
		URL(uri)

	if !disableUpdateBtn {
		btn := VBtn("").
			Color("primary").
			Variant(VariantFlat).
			Attr(":disabled", "isFetching").
			Attr(":loading", "isFetching")

		// if overlayType == "" && portalName != "" {
		// não está funcionando,
		//	var cb web.Callback
		//	cb.ReloadPortals = append(cb.ReloadPortals, portalName)
		//		onClick.Query(ParamPostChangeCallback, cb.Encode())
		//	}

		btn.Attr("@click", onClick.Go())
		btn.Icon(true).
			Density("comfortable").
			Children(VIcon("mdi-content-save"))
		f.PrimaryAction = btn
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

func (b *EditingBuilder) formBuilder(obj interface{}, ctx *web.EventContext) (fb *FormBuilder) {
	return NewFormBuilder(ctx, b.mb, &b.FieldsBuilder, obj)
}

func (b *EditingBuilder) form(obj interface{}, ctx *web.EventContext) *Form {
	return b.ConfigureForm(b.formBuilder(obj, ctx).Build())
}
