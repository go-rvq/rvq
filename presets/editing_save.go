package presets

import (
	"context"

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

	resetSaveContext := web.WithContextValue(ctx, CtxSaveContext, context.Background())
	defer resetSaveContext()

	var done []func(success bool) error

	defer func() {
		for _, f := range done {
			if e := f(err == nil); err == nil && e != nil {
				err = e
			}
		}
	}()

	if usingB.preSaveCallback != nil {
		if d, e := usingB.preSaveCallback(ctx, obj); e != nil {
			err = e
			return
		} else if d != nil {
			done = append(done, d)
		}
	}

	if usingB.preValidate != nil {
		if err = usingB.preValidate(ctx, obj); err != nil {
			return
		}
	}

	if vErr = usingB.Validators.Validate(obj, FieldModeStack{EDIT}, ctx); vErr.HaveErrors() {
		usingB.UpdateOverlayContent(ctx, r, obj, "", &vErr)
		return &vErr
	}

	usingB.FieldsBuilder.WalkO(usingB.mb.modelInfo, obj, FieldModeStack{EDIT}, ctx, &FieldWalkHandleOptions{
		SkipNestedNil: true,
		Handler: func(field *FieldContext) (s FieldWalkState) {
			vErr.Merge(field.Field.Validators.Validate(field))
			return s
		},
	})

	if vErr.HaveErrors() {
		usingB.UpdateOverlayContent(ctx, r, obj, "", &vErr)
		return &vErr
	}

	if usingB.postValidate != nil {
		if err = usingB.postValidate(ctx, obj); err != nil {
			return
		}
	}

	err1 := usingB.Saver(obj, mid, ctx)
	if err1 != nil {
		usingB.UpdateOverlayContent(ctx, r, obj, "", err1)
		return err1
	}

	if usingB.postSaveCallback != nil {
		if d, e := usingB.postSaveCallback(ctx, obj); e != nil {
			err = e
			return
		} else if d != nil {
			done = append(done, d)
		}
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
