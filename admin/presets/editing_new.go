package presets

import (
	"context"
	"strings"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/perm"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
)

func (mb *ModelBuilder) SetCreatingBuilder(b *EditingBuilder) {
	mb.creating = b
}

func (b *EditingBuilder) Creating(vs ...interface{}) (r *EditingBuilder) {
	if b.mb.creating != nil && len(vs) == 0 {
		return b.mb.creating
	}

	if b.mb.creating == nil {
		b.mb.creating = &EditingBuilder{
			mb:         b.mb,
			Fetcher:    b.Fetcher,
			Setter:     b.Setter,
			Creator:    b.Creator,
			Validators: append(Validators{}, b.Validators...),
		}
	}

	fb := *b.FieldsBuilder.Clone()

	r = b.mb.creating
	if len(vs) == 0 {
		for _, f := range b.fields {
			vs = append(vs, f.name)
		}
	}

	if len(vs) == 0 {
		fb.fields = b.fields.HasMode(NEW)
		b.fields = b.fields.HasMode(EDIT)
	} else {
		fb.fields = append(FieldBuilders{}, b.fields...)
		fb = *fb.Only(vs...)
	}

	r.FieldsBuilder = fb
	return r
}

func (b *EditingBuilder) defaultCreate(ctx *web.EventContext) (r web.EventResponse, err error) {
	uErr := b.doCreate(ctx, &r, false)
	if uErr == nil {
		msgr := MustGetMessages(ctx.Context())
		ShowMessage(&r, msgr.SuccessfullyCreated, "success")
	}
	return r, nil
}

func (b *EditingBuilder) doCreate(
	ctx *web.EventContext,
	r *web.EventResponse,
	// will not close drawer/dialog
	silent bool,
) (err error) {
	if b.mb.creatingDisabled {
		err = ErrCreateRecordNotAllowed
		return
	}

	b = b.CreatingBuilder()
	obj := b.mb.NewModel()

	if b.mb.permissioner.ReqCreator(ctx.R).Denied() {
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

	if vErr := b.RunSetterFunc(nil, ctx, false, obj); vErr.HaveErrors() {
		b.UpdateOverlayContent(ctx, r, obj, "", &vErr)
		return &vErr
	}

	if b.preSaveCallback != nil {
		if d, e := b.preSaveCallback(ctx, obj); e != nil {
			err = e
			return
		} else if d != nil {
			done = append(done, d)
		}
	}

	if b.preValidate != nil {
		if err = b.preValidate(ctx, obj); err != nil {
			return
		}
	}

	if vErr := b.Validators.Validate(obj, FieldModeStack{NEW}, ctx); vErr.HaveErrors() {
		b.UpdateOverlayContent(ctx, r, obj, "", &vErr)
		return &vErr
	}

	{
		var vErr web.ValidationErrors
		b.FieldsBuilder.WalkOptions(b.mb.modelInfo, obj, FieldModeStack{NEW}, ctx, &FieldWalkHandleOptions{
			SkipNestedNil: true,
			Handler: func(field *FieldContext) (s FieldWalkState) {
				vErr.Merge(field.Field.Validators.Validate(field))
				return s
			},
		})

		if vErr.HaveErrors() {
			b.UpdateOverlayContent(ctx, r, obj, "", &vErr)
			return &vErr
		}
	}

	restore := RemoveEmptySliceItems(obj, ContextModifiedIndexesBuilder(ctx))
	err1 := b.Creator(obj, ctx)
	if err1 != nil {
		restore()
		b.UpdateOverlayContent(ctx, r, obj, "", err1)
		return err1
	}

	id := vx.ObjectID(obj)

	if postSaveConfig := ctx.R.URL.Query().Get(ParamPostChangeCallback); postSaveConfig != "" {
		cb := web.DecodeCallback(postSaveConfig)
		if len(cb.Scripts) > 0 {
			web.AppendRunScripts(r, strings.ReplaceAll(cb.Script(), "$ID$", id))
		}
		r.ReloadPortals = append(r.ReloadPortals, cb.ReloadPortals...)
	}

	overlay := actions.OverlayMode(ctx.R.FormValue(ParamOverlay))
	if overlay.Overlayed() {
		web.AppendRunScripts(r, "closer.show = false")
	} else {
		r.PushState = web.Location(nil)
	}
	return
}

func (b *EditingBuilder) formNew(ctx *web.EventContext) (r web.EventResponse, err error) {
	if b.mb.creatingDisabled {
		err = ErrCreateRecordNotAllowed
		return
	}

	b = b.CreatingBuilder()

	if b.mb.permissioner.ReqCreator(ctx.R).Denied() {
		err = perm.PermissionDenied
		return
	}

	obj := b.mb.NewModel()

	if b.New != nil {
		if err = b.New(ctx, obj); err != nil {
			return
		}
	}

	respondTargetPortal := ctx.R.FormValue(ParamTargetPortal)
	overlay := actions.OverlayMode(ctx.R.FormValue(ParamOverlay))
	if overlay.IsDrawer() {
		respondTargetPortal = overlay.PortalName()
	}
	targetPortal := respondTargetPortal
	ctx.R.Form.Set(ParamTargetPortal, targetPortal)

	f := b.form(obj, ctx)
	f.Portal = targetPortal
	f.Wrap = func(c h.HTMLComponent) h.HTMLComponent {
		return web.Scope(web.Portal(c).
			Name(targetPortal)).FormInit()
	}

	f.RespondToPortal(respondTargetPortal, &r)
	return
}

func (b *EditingBuilder) HasCreatingBuilder() bool {
	return b.mb.creating != nil
}

func (b *EditingBuilder) CreatingBuilder() (c *EditingBuilder) {
	if b.mb.creating != nil {
		return b.mb.creating
	}
	return b
}
