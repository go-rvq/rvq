package presets

import (
	"github.com/jinzhu/inflection"
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	"github.com/qor5/x/v3/perm"
)

func (b *EditingBuilder) singletonPageFunc(ctx *web.EventContext) (r web.PageResponse, err error) {
	if b.mb.Info().Verifier().Do(PermUpdate).WithReq(ctx.R).IsAllowed() != nil {
		err = perm.PermissionDenied
		return
	}

	msgr := MustGetMessages(ctx.R)
	title := msgr.EditingObjectTitle(i18n.T(ctx.R, ModelsI18nModuleKey, inflection.Singular(b.mb.label)), "")
	r.PageTitle = title
	r.Body = web.Portal().
		Name(actions.Edit).
		Loader(
			web.GET().
				EventFunc(actions.Edit).
				URL(b.mb.Info().ListingHrefCtx(ctx)).
				Query(ParamTargetPortal, actions.Edit),
		)
	return
}

func (b *EditingBuilder) doSingletonEditing(ctx *web.EventContext) (r web.PageResponse, err error) {
	if !b.CanEditObj(ctx, nil) {
		err = perm.PermissionDenied
		return
	}

	msgr := MustGetMessages(ctx.R)
	title := msgr.EditingObjectTitle(i18n.T(ctx.R, ModelsI18nModuleKey, inflection.Singular(b.mb.label)), "")
	r.PageTitle = title
	obj := b.mb.NewModel()
	err = b.Fetcher(b.mb.NewModel(), "", ctx)
	if err == ErrRecordNotFound {
		if err = b.Saver(obj, "", ctx); err != nil {
			return
		}
		err = b.Fetcher(obj, "", ctx)
	}
	if err != nil {
		return
	}
	r.Body = b.form(obj, ctx).Component()
	return
}
