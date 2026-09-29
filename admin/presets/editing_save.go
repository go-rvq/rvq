package presets

import (
	"context"

	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
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

	obj, vErr := usingB.FetchAndUnmarshal(nil, mid, true, ctx)
	if vErr.HaveErrors() {
		usingB.UpdateOverlayContent(ctx, r, obj, "", &vErr)
		return &vErr
	}

	if !b.CanEditObj(obj, ctx) {
		return perm.PermissionDenied
	}

	// Before ANY validation: what the form says is only worth validating if it
	// was written over the record as it stands now — this compares the record as
	// it was when the form was rendered against the record today.
	if err = usingB.VerifyRecordStamp(mid, ctx); err != nil {
		// re-render with the message: the form then shows the stored values and
		// carries a fresh stamp, so a deliberate re-save is possible.
		usingB.UpdateOverlayContent(ctx, r, obj, "", err)
		return
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

	usingB.FieldsBuilder.WalkOptions(usingB.mb.modelInfo, obj, FieldModeStack{EDIT}, ctx, &FieldWalkHandleOptions{
		SkipNestedNil:         true,
		SkipListEditorDeleted: true,
		Handler: func(field *FieldContext) (s FieldWalkState) {
			vErr.Merge(field.Field.Validate(field))
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

	restore := RemoveEmptySliceItems(obj, ContextModifiedIndexesBuilder(ctx))
	err1 := usingB.Saver(obj, mid, ctx)
	if err1 != nil {
		restore()
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
	script := `closer.show = false;`

	// Refresh what the change affects. The form does not know its opener: it just
	// calls back into the closer, which IS the host's state — and the host was
	// rendered by whoever must refresh (see FormHostBuilder.OnSave).
	web.AppendRunScripts(r, PostSaveScript(mid.String()))

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
	// a save that has something to say — ctx.Flash — keeps it
	if err = b.doUpdate(ctx, &r, false); err == nil && ctx.Flash == nil {
		ctx.Flash = MustGetMessages(ctx.Context()).SuccessfullyUpdated
	}
	return
}

func (b *EditingBuilder) SaveOverlayContent(
	ctx *web.EventContext,
	r *web.EventResponse,
) (err error) {
	return b.doUpdate(ctx, r, true)
}
