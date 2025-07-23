package presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"

	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

func (b *EditingBuilder) respondFormEdit(ctx *web.EventContext, obj any, initForm ...bool) (r web.EventResponse, err error) {
	targetPortal := ctx.R.FormValue(ParamTargetPortal)
	overlay := actions.OverlayMode(ctx.R.FormValue(ParamOverlay))
	if overlay.IsDrawer() && targetPortal == "" {
		targetPortal = overlay.PortalName()
	}

	f := b.form(obj, ctx)
	f.ScopeDisabled = ctx.R.FormValue(ParamEditFormUnscoped) == "true"

	comp := f.Component()
	mode := GetOverlay(ctx)

	for _, v := range initForm {
		if v {
			comp = web.Scope(comp).FormInit()
		}
	}

	if mode.IsDrawer() {
		b.mb.p.Drawer(mode).
			SetScrollable(true).
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

func (b *EditingBuilder) formEdit(ctx *web.EventContext) (r web.EventResponse, err error) {
	if b.mb.editingDisabled {
		err = ErrUpdateRecordNotAllowed
		return
	}

	obj := b.mb.NewModel()
	var mid ID
	if mid, err = b.mb.ParseRecordID(ctx.Queries().Get(ParamID)); err != nil {
		return
	}

	if b.mb.permissioner.Updater(ctx.R, mid, ParentsModelID(ctx.R)...).Denied() {
		err = perm.PermissionDenied
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

	return b.respondFormEdit(ctx, obj, true)
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
			b.mb.TTitleAuto(ctx.Context()),
			b.mb.female,
		)
	} else {
		disableUpdateBtn = f.b.mb.permissioner.ReqObjectUpdater(f.b.ctx.R, f.Obj).Denied()
		var editingTitleText string
		if b.mb.singleton {
			editingTitleText = f.b.msgr.EditingTitle(b.mb.TTitleAuto(ctx.Context()))
		} else {
			editingTitleText = f.b.msgr.EditingObjectTitle(
				b.mb.TTitle(ctx.Context()),
				b.mb.RecordTitle(f.Obj, ctx))
		}
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
	return b.ConfigureForm(NewFormBuilder(ctx, b.mb, &b.FieldsBuilder, obj).SetPre(b.preComponents).SetPost(b.postComponents).Build())
}
