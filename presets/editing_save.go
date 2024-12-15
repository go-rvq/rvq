package presets

import (
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/perm"
)

func (b *EditingBuilder) doUpdate(
	ctx *web.EventContext,
	r *web.EventResponse,
	// will not close drawer/dialog
	silent bool,
) (err error) {
	var mid ID
	if mid, err = b.mb.ParseRecordID(ctx.Queries().Get(ParamID)); err != nil {
		return
	}
	if mid.IsZero() && !b.mb.singleton {
		return MustGetMessages(ctx.Context()).ErrEmptyParamID
	}

	usingB := b

	obj, vErr := usingB.FetchAndUnmarshal(mid, true, ctx)
	if vErr.HaveErrors() {
		usingB.UpdateOverlayContent(ctx, r, obj, "", &vErr)
		return &vErr
	}

	if !b.CanEditObj(obj, ctx) {
		b.UpdateOverlayContent(ctx, r, obj, "", perm.PermissionDenied)
		return perm.PermissionDenied
	}

	if vErr = usingB.Validators.Validate(obj, FieldModeStack{EDIT}, ctx); vErr.HaveErrors() {
		usingB.UpdateOverlayContent(ctx, r, obj, "", &vErr)
		return &vErr
	}

	usingB.FieldsBuilder.Walk(usingB.mb.modelInfo, obj, FieldModeStack{EDIT}, ctx, func(field *FieldContext) (s FieldWalkState) {
		vErr.Merge(field.Field.Validators.Validate(field))
		return s
	})

	if vErr.HaveErrors() {
		usingB.UpdateOverlayContent(ctx, r, obj, "", &vErr)
		return &vErr
	}

	err1 := usingB.Saver(obj, mid, ctx)
	if err1 != nil {
		usingB.UpdateOverlayContent(ctx, r, obj, "", err1)
		return err1
	}

	overlay := actions.OverlayMode(ctx.R.FormValue(ParamOverlay))
	script := "closer.show = false"

	if postSaveConfig := ctx.R.URL.Query().Get(ParamPostChangeCallback); postSaveConfig != "" {
		cb := web.DecodeCallback(postSaveConfig)
		web.AppendRunScripts(r, web.ApplyChangeEvent(cb.Script(), web.Updated, mid.String()))
		r.ReloadPortals = append(r.ReloadPortals, cb.ReloadPortals...)
	}

	if overlay != "" {
		web.AppendRunScripts(r, script)
	} else {
		r.PushState = web.Location(nil)
	}
	return
}

func (b *EditingBuilder) defaultUpdate(ctx *web.EventContext) (r web.EventResponse, err error) {
	uErr := b.doUpdate(ctx, &r, false)
	if uErr == nil {
		msgr := MustGetMessages(ctx.Context())
		ShowMessage(&r, msgr.SuccessfullyUpdated, "")
	}
	return r, nil
}

func (b *EditingBuilder) SaveOverlayContent(
	ctx *web.EventContext,
	r *web.EventResponse,
) (err error) {
	return b.doUpdate(ctx, r, true)
}
