package publish

import (
	"fmt"

	"github.com/qor5/admin/v3/model"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
)

func Publish(mb *presets.ModelBuilder, publisher *Builder, actionName string, ctx *web.EventContext, mid model.ID) (obj any, err error) {
	obj = mb.NewModel()
	err = mb.Editing().Fetcher(obj, mid, ctx)
	if err != nil {
		return
	}
	reqCtx := publisher.WithContextValues(ctx.R.Context())
	err = publisher.Publish(mb, obj, reqCtx)
	if err != nil {
		return
	}

	if publisher.ab != nil {
		if _, exist := publisher.ab.GetModelBuilder(obj); exist {
			publisher.ab.AddCustomizedRecord(actionName, false, ctx.R.Context(), obj)
		}
	}

	return
}

func publishAction(mb *presets.ModelBuilder, publisher *Builder, actionName string) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		var mid model.ID
		if mid, err = mb.ParseRecordID(ctx.Param(presets.ParamID)); err != nil {
			return
		}

		var obj any
		if obj, err = Publish(mb, publisher, actionName, ctx, mid); err != nil {
			return
		}

		if status, ok := obj.(StatusInterface); ok {
			web.AppendRunScripts(&r, fmt.Sprintf("locals.%s = %q", FieldOnlineUrl, status.EmbedStatus().OnlineUrl))
		}

		if script := ctx.R.FormValue(ParamScriptAfterPublish); script != "" {
			web.AppendRunScripts(&r, script)
		} else {
			presets.ShowMessage(&r, "success", "")
			r.Reload = true
		}
		return
	}
}

func Unpublish(mb *presets.ModelBuilder, publisher *Builder, actionName string, ctx *web.EventContext, mid model.ID) (obj any, err error) {
	obj = mb.NewModel()

	if err = mb.Editing().Fetcher(obj, mid, ctx); err != nil {
		return
	}

	if obj.(StatusInterface).EmbedStatus().Status != StatusOnline {
		return
	}

	reqCtx := publisher.WithContextValues(ctx.R.Context())
	err = publisher.UnPublish(mb, obj, reqCtx)
	if err != nil {
		return
	}
	if publisher.ab != nil {
		if _, exist := publisher.ab.GetModelBuilder(obj); exist {
			publisher.ab.AddCustomizedRecord(actionName, false, ctx.R.Context(), obj)
		}
	}
	return
}

func unpublishAction(mb *presets.ModelBuilder, publisher *Builder, actionName string) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		var mid model.ID
		if mid, err = mb.ParseRecordID(ctx.Param(presets.ParamID)); err != nil {
			return
		}

		if _, err = Unpublish(mb, publisher, actionName, ctx, mid); err != nil {
			return
		}
		presets.ShowMessage(&r, "success", "")
		r.Reload = true
		return
	}
}
