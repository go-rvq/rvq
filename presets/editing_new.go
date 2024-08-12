package presets

import (
	"strings"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/perm"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
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
			Saver:      b.Saver,
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

	b.mb.creating.FieldsBuilder = fb
	return r
}

func (b *EditingBuilder) defaultCreate(ctx *web.EventContext) (r web.EventResponse, err error) {
	uErr := b.doCreate(ctx, &r, false)
	if uErr == nil {
		msgr := MustGetMessages(ctx.R)
		ShowMessage(&r, msgr.SuccessfullyCreated, "")
	}
	return r, nil
}

func (b *EditingBuilder) doCreate(
	ctx *web.EventContext,
	r *web.EventResponse,
	// will not close drawer/dialog
	silent bool,
) (err error) {
	b = b.CreatingBuilder()
	obj := b.mb.NewModel()

	if b.mb.Info().Verifier().Do(PermCreate).ObjectOn(obj).WithReq(ctx.R).IsAllowed() != nil {
		b.UpdateOverlayContent(ctx, r, obj, "", perm.PermissionDenied)
		return perm.PermissionDenied
	}

	if vErr := b.RunSetterFunc(ctx, false, obj); vErr.HaveErrors() {
		b.UpdateOverlayContent(ctx, r, obj, "", &vErr)
		return &vErr
	}
	if vErr := b.Validators.Validate(obj, FieldModeStack{NEW}, ctx); vErr.HaveErrors() {
		b.UpdateOverlayContent(ctx, r, obj, "", &vErr)
		return &vErr
	}

	{
		var vErr web.ValidationErrors
		b.FieldsBuilder.Walk(b.mb.modelInfo, obj, FieldModeStack{NEW}, ctx, func(field *FieldContext) (s FieldWalkState) {
			vErr.Merge(field.Field.Validators.Validate(field))
			return s
		})

		if vErr.HaveErrors() {
			b.UpdateOverlayContent(ctx, r, obj, "", &vErr)
			return &vErr
		}
	}

	err1 := b.Saver(obj, "", ctx)
	if err1 != nil {
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
	script := "closer.show = false"

	if overlay != "" {
		web.AppendRunScripts(r, script)
	} else {
		r.PushState = web.Location(nil)
	}
	return
}

func (b *EditingBuilder) formNew(ctx *web.EventContext) (r web.EventResponse, err error) {
	b = b.CreatingBuilder()

	if b.mb.Info().Verifier().Do(PermCreate).WithReq(ctx.R).IsAllowed() != nil {
		ShowMessage(&r, perm.PermissionDenied.Error(), "warning")
		return
	}

	obj := b.mb.NewModel()
	f := b.form(obj, ctx)
	if f.b.overlayMode.IsDrawer() {
		f.Portal = f.b.overlayMode.PortalName()
	}
	f.Respond(&r)
	return
}

func (b *EditingBuilder) CreatingBuilder() (c *EditingBuilder) {
	c = b
	if b.mb.creating != nil {
		c = b.mb.creating
	}
	return c
}
