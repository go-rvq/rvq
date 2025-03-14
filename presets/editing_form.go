package presets

import (
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
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
	var mid ID
	if mid, err = b.mb.ParseRecordID(ctx.Queries().Get(ParamID)); err != nil {
		return
	}

	if !mid.IsZero() || b.mb.singleton {
		if err = b.Fetcher(obj, mid, ctx); err != nil {
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

	targetPortal := ctx.R.FormValue(ParamTargetPortal)
	overlay := actions.OverlayMode(ctx.R.FormValue(ParamOverlay))
	if overlay.IsDrawer() {
		targetPortal = overlay.PortalName()
	}

	f := b.form(obj, ctx)
	f.ScopeDisabled = ctx.R.FormValue(ParamEditFormUnscoped) == "true"

	comp := f.Component()
	mode := GetOverlay(ctx)

	if mode.IsDrawer() {
		b.mb.p.Drawer(mode).
			SetValidPortalName(targetPortal).
			Respond(&r, comp)
	} else if mode.IsDialog() {
		b.mb.p.Dialog().
			SetScrollable(true).
			SetValidWidth(b.mb.rightDrawerWidth).
			SetTargetPortal(targetPortal).
			Respond(ctx, &r, comp)
	} else {
		r.Body = comp
	}

	return
}

func (b *EditingBuilder) SaveBtn(ctx *web.EventContext, id string, edit bool, targetPortal string) h.HTMLComponent {
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

	if targetPortal != "" {
		queries.Del(ParamTargetPortal)
	}

	onClick := web.Plaid().
		EventFunc(event).
		Queries(queries).
		Method("POST").
		ValidQuery(ParamTargetPortal, targetPortal).
		URL(b.mb.Info().ListingHrefCtx(ctx))

	return web.Scope(VBtn("").
		Color("primary").
		Variant(VariantFlat).
		Attr(":disabled", "isFetching").
		Attr(":loading", "isFetching").
		Attr("@click", onClick.Go()).
		Icon(true).
		Density("comfortable").
		Children(VIcon("mdi-content-save"))).Form()
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
			b.mb.TTitle(ctx.Context()),
			b.mb.female,
		)
	} else {
		disableUpdateBtn = !f.b.mb.Info().CanUpdate(f.b.ctx.R, f.Obj)

		editingTitleText := f.b.msgr.EditingObjectTitle(
			b.mb.TTitle(ctx.Context()),
			b.mb.RecordTitle(f.Obj, ctx))
		if b.editingTitleFunc != nil {
			f.Title = b.editingTitleFunc(f.b.obj, editingTitleText, f.b.ctx)
		} else {
			f.Title = editingTitleText
		}
	}

	if !disableUpdateBtn {
		f.PrimaryAction = b.SaveBtn(f.b.ctx, f.b.id, f.b.mode != NEW, f.Portal)
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
