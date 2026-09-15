package presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"

	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

// formScope responds with a self-opening FormHost: it owns the form's `closer`
// and its `form` scope, and loads the real create/edit event (innerEvent) into
// its inner portal as soon as it mounts (Show(true)).
//
// This is the entry point for callers that cannot host the form themselves (a
// row menu, an action button elsewhere, …). Pages that CAN host it — the
// detailing (edit) and the listing (create) — embed a FormHost directly and just
// flip `<scope>.show`, saving this round-trip.
//
// Either way the semantics are the same: the host owns the closer, so turning it
// off destroys the form and its scope completely, and turning it on loads a fresh
// one; the form's own re-renders (list editor, validation) never recreate the
// scope.
func (b *EditingBuilder) formScope(ctx *web.EventContext, scope, innerEvent string) (r web.EventResponse, err error) {
	innerPortal := ctx.UID()

	queries := ctx.Queries()
	queries.Set(ParamTargetPortal, innerPortal)

	comp := FormHost(scope, innerPortal,
		web.Plaid().
			URL(b.mb.Info().ListingHrefCtx(ctx)).
			EventFunc(innerEvent).
			Queries(queries)).
		Show(true).
		Component()

	if target := ctx.R.FormValue(ParamTargetPortal); target != "" {
		r.UpdatePortal(target, comp)
	} else {
		r.Body = comp
	}
	return
}

func (b *EditingBuilder) formEditScope(ctx *web.EventContext) (web.EventResponse, error) {
	return b.formScope(ctx, DetailingEditScope, actions.Edit)
}

func (b *EditingBuilder) formNewScope(ctx *web.EventContext) (web.EventResponse, error) {
	return b.formScope(ctx, ListingNewScope, actions.New)
}

func (b *EditingBuilder) respondFormEdit(ctx *web.EventContext, obj any) (r web.EventResponse, err error) {
	targetPortal := ctx.R.FormValue(ParamTargetPortal)
	overlay := actions.OverlayMode(ctx.R.FormValue(ParamOverlay))
	if overlay.IsDrawer() && (targetPortal == "" || CloserIsGlobal(ctx)) {
		// see editing_new.go: the layout's portal, whenever the caller's closer is
		// reachable from there
		targetPortal = overlay.PortalName()
	}

	// The `form` scope is established once by the EditForm/NewForm wrapper; the
	// form here renders inside it (never creating its own scope).
	f := b.form(obj, ctx)
	comp := f.Component()
	mode := GetOverlay(ctx)

	if mode.IsDrawer() {
		b.mb.p.Drawer(mode).
			SetScrollable(true).
			SetValidPortalName(targetPortal).
			SetCloserProvided(CloserProvided(ctx)).
			SetCloserRef(CloserRef(ctx)).
			Respond(&r, comp)
	} else if mode.IsDialog() {
		b.mb.p.Dialog().
			SetScrollable(true).
			SetValidWidth(b.mb.rightDrawerWidth).
			SetTargetPortal(targetPortal).
			Respond(ctx, &r, comp)
	} else if mode == actions.Content {
		// Content renders into a portal the caller already has on the page: the
		// one it asked for, or the layout's global content portal when it did
		// not (see Builder.contentDrawer, which decides the same way).
		portal := targetPortal
		if portal == "" {
			portal = actions.RightDrawer.PortalName()
		}
		r.UpdatePortal(portal, comp)
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

	return b.respondFormEdit(ctx, obj)
}

func (b *EditingBuilder) SaveBtn(ctx *web.EventContext, id string, edit bool, targetPortal string) h.HTMLComponent {
	var (
		queries = ctx.Queries()
		event   = actions.Create
	)

	if id != "" {
		queries.Set(ParamID, id)
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
		Attr("data-event", event).
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

	if b.topRightActionsFunc != nil {
		if c := b.topRightActionsFunc(f.Obj, f.b.ctx); c != nil {
			f.TopRightActions = append(f.TopRightActions, c)
		}
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

	// The record as it was when this form was rendered, signed. The update
	// refuses to run if the stored record has moved since (see
	// EditingBuilder.VerifyRecordStamp). A creation has nothing to compare to.
	if f.b.mode != NEW {
		if stamp := b.recordStampField(f.b.obj); stamp != nil {
			hiddenComps = append(hiddenComps, stamp)
		}
	}

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
