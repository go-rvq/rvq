package presets

import (
	"fmt"
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
			mb:        b.mb,
			Fetcher:   b.Fetcher,
			Setter:    b.Setter,
			Saver:     b.Saver,
			Deleter:   b.Deleter,
			Validator: b.Validator,
		}
	}

	b.mb.creating.FieldsBuilder = *b.FieldsBuilder.Clone()
	r = b.mb.creating
	if len(vs) == 0 {
		for _, f := range b.fields {
			vs = append(vs, f.name)
		}
	}

	r.FieldsBuilder = *b.FieldsBuilder.Only(vs...)

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
	var (
		creating = b.CreatingBuilder()
		obj      = b.mb.NewModel()
	)

	if b.mb.Info().Verifier().Do(PermCreate).ObjectOn(obj).WithReq(ctx.R).IsAllowed() != nil {
		b.UpdateOverlayContent(ctx, r, obj, "", perm.PermissionDenied)
		return perm.PermissionDenied
	}

	if vErr := b.RunSetterFunc(ctx, false, obj); vErr.HaveErrors() {
		creating.UpdateOverlayContent(ctx, r, obj, "", &vErr)
		return &vErr
	}

	if creating.Validator != nil {
		if vErr := creating.Validator(obj, ctx); vErr.HaveErrors() {
			creating.UpdateOverlayContent(ctx, r, obj, "", &vErr)
			return &vErr
		}
	}

	err1 := creating.Saver(obj, "", ctx)
	if err1 != nil {
		creating.UpdateOverlayContent(ctx, r, obj, "", err1)
		return err1
	}

	overlayType := ctx.R.FormValue(ParamOverlay)
	script := CloseRightDrawerVarScript
	if overlayType == actions.Dialog {
		script = closeDialogVarScript
	}
	if silent {
		script = ""
	}

	afterUpdateScript := ctx.R.FormValue(ParamOverlayAfterUpdateScript)
	if afterUpdateScript != "" {
		web.AppendRunScripts(r, script, strings.NewReplacer(".go()",
			fmt.Sprintf(".query(%s, %s).go()",
				h.JSONString(ParamOverlayUpdateID),
				h.JSONString(vx.ObjectID(obj)),
			)).Replace(afterUpdateScript),
		)

		return
	}

	if isInDialogFromQuery(ctx) {
		/*web.AppendRunScripts(r,
			web.Plaid().
				URL(ctx.R.RequestURI).
				EventFunc(actions.UpdateListingDialog).
				StringQuery(ctx.R.URL.Query().Get(ParamListingQueries)).
				Go(),
		)*/

		if postSaveConfig := ctx.R.URL.Query().Get(PostSaveCallback); postSaveConfig != "" {
			cb := web.DecodeCallback(postSaveConfig)
			if len(cb.Scripts) > 0 {
				web.AppendRunScripts(r, cb.Script())
			}
			r.ReloadPortals = append(r.ReloadPortals, cb.ReloadPortals...)
		}
	} else {
		r.PushState = web.Location(nil)
	}
	web.AppendRunScripts(r, script)
	return
}

func (b *EditingBuilder) formNew(ctx *web.EventContext) (r web.EventResponse, err error) {
	if b.mb.Info().Verifier().Do(PermCreate).WithReq(ctx.R).IsAllowed() != nil {
		ShowMessage(&r, perm.PermissionDenied.Error(), "warning")
		return
	}

	obj := b.mb.NewModel()
	body := b.CreatingBuilder().editFormFor(obj, ctx, "")

	if ctx.R.FormValue(ParamOverlay) != "" {
		b.mb.p.overlay(ctx, &r, body, b.mb.rightDrawerWidth)
	} else {
		r.Body = body
	}
	return
}

func (b *EditingBuilder) CreatingBuilder() (c *EditingBuilder) {
	c = b
	if b.mb.creating != nil {
		c = b.mb.creating
	}
	return c
}
