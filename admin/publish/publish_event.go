package publish

import (
	"fmt"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

type Executor interface {
	Execute(mb *presets.ModelBuilder, publisher *Builder, actionName string, ctx *web.EventContext, mid model.ID) (obj any, err error)
	Title(msgr *Messages) string
	ActivityName() string
	Accept(status string) bool
}

type PublishFlag uint8

const (
	FlagPublishRepublish PublishFlag = iota + 1
	FlagPublishOrRepublish
)

type PublishOptions struct {
	Flag PublishFlag
}

type PublishOption func(*PublishOptions)

func PublishWithFlag(flag PublishFlag) PublishOption {
	return func(o *PublishOptions) {
		o.Flag = flag
	}
}

type PublishExecutor struct {
	opts PublishOptions
}

func (p PublishExecutor) Option(opt ...PublishOption) PublishExecutor {
	for _, o := range opt {
		o(&p.opts)
	}
	return p
}

func (p PublishExecutor) Execute(mb *presets.ModelBuilder, publisher *Builder, actionName string, ctx *web.EventContext, mid model.ID) (obj any, err error) {
	obj = mb.NewModel()
	err = mb.Editing().Fetcher(obj, mid, ctx)
	if err != nil {
		return
	}
	err = p.Do(mb, publisher, actionName, ctx, obj)
	return
}

func (p PublishExecutor) Accept(status string) (ok bool) {
	switch status {
	case StatusDraft, StatusOffline:
		if p.opts.Flag == FlagPublishRepublish {
			return
		}
	case StatusOnline:
		switch p.opts.Flag {
		case FlagPublishOrRepublish:
		default:
			return
		}
	}
	return true
}

func (p PublishExecutor) Do(mb *presets.ModelBuilder, publisher *Builder, actionName string, ctx *web.EventContext, obj any) (err error) {
	reqCtx := publisher.WithContextValues(ctx.R.Context())
	if status, ok := obj.(StatusInterface); ok {
		switch status.EmbedStatus().Status {
		case StatusDraft, StatusOffline:
			if p.opts.Flag == FlagPublishRepublish {
				return
			}
		case StatusOnline:
			switch p.opts.Flag {
			case FlagPublishRepublish, FlagPublishOrRepublish:
				if err = UnPublish.Do(mb, publisher, actionName, ctx, obj); err != nil {
					return
				}
			default:
				return
			}
		}
	}

	if err = publisher.Publish(mb, obj, reqCtx); err != nil {
		return
	}

	if publisher.ab != nil {
		if _, exist := publisher.ab.GetModelBuilder(obj); exist {
			publisher.ab.AddCustomizedRecord(actionName, false, ctx.R.Context(), obj)
		}
	}

	return
}

func (p PublishExecutor) ActivityName() string {
	switch p.opts.Flag {
	case FlagPublishRepublish:
		return ActivityRepublish
	case FlagPublishOrRepublish:
		return ActivityPublishOrRepublish
	default:
		return ActivityPublish
	}
}

func (p PublishExecutor) Title(msgr *Messages) string {
	switch p.opts.Flag {
	case FlagPublishRepublish:
		return msgr.Republish
	case FlagPublishOrRepublish:
		return msgr.PublishOrRepublish
	default:
		return msgr.Publish
	}
}

type UnpublishExecutor struct {
}

func (e UnpublishExecutor) Execute(mb *presets.ModelBuilder, publisher *Builder, actionName string, ctx *web.EventContext, mid model.ID) (obj any, err error) {
	obj = mb.NewModel()
	if err = mb.Editing().Fetcher(obj, mid, ctx); err != nil {
		return
	}
	err = e.Do(mb, publisher, actionName, ctx, obj)
	return
}

func (UnpublishExecutor) Accept(status string) (ok bool) {
	return status == StatusOnline
}

func (e UnpublishExecutor) Do(mb *presets.ModelBuilder, publisher *Builder, actionName string, ctx *web.EventContext, obj any) (err error) {
	if !e.Accept(obj.(StatusInterface).EmbedStatus().Status) {
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

func (p UnpublishExecutor) Title(msgr *Messages) string {
	return msgr.Unpublish
}

func (p UnpublishExecutor) ActivityName() string {
	return ActivityUnPublish
}

var (
	Publish            PublishExecutor
	PublishOrRepublish = Publish.Option(PublishWithFlag(FlagPublishOrRepublish))
	RePublish          = Publish.Option(PublishWithFlag(FlagPublishRepublish))
	UnPublish          UnpublishExecutor
)

func ExecutorFromActivityName(name string) Executor {
	switch name {
	case ActivityPublish:
		return Publish
	case ActivityPublishOrRepublish:
		return PublishOrRepublish
	case ActivityRepublish:
		return RePublish
	case ActivityUnPublish:
		return UnPublish
	default:
		return nil
	}
}

func publishAction(mb *presets.ModelBuilder, publisher *Builder, actionName string) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		var mid model.ID
		if mid, err = mb.ParseRecordID(ctx.Param(presets.ParamID)); err != nil {
			return
		}

		var e = Publish
		switch actionName {
		case ActivityPublishOrRepublish:
			e = PublishOrRepublish
		case ActivityRepublish:
			e = RePublish
		}

		var obj any
		if obj, err = e.Execute(mb, publisher, actionName, ctx, mid); err != nil {
			return
		}

		if status, ok := obj.(StatusInterface); ok {
			web.AppendRunScripts(&r, fmt.Sprintf("locals.%s = %q", FieldOnlineUrl, status.EmbedStatus().OnlineUrl))
		}

		if script := ctx.R.FormValue(ParamScriptAfterPublish); script != "" {
			web.AppendRunScripts(&r, script)
		} else {
			var (
				msgr = GetMessages(ctx.Context())
				msg  = msgr.SuccessfullyPublished
			)

			switch actionName {
			case ActivityPublishOrRepublish:
				msg = msgr.SuccessfullyPublishedOrRepublished
			case ActivityRepublish:
				msg = msgr.SuccessfullyRepublished
			}

			presets.ShowMessage(&r, msg, "")
			r.Reload = true
		}
		return
	}
}

func unpublishAction(mb *presets.ModelBuilder, publisher *Builder, actionName string) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		var mid model.ID
		if mid, err = mb.ParseRecordID(ctx.Param(presets.ParamID)); err != nil {
			return
		}

		if _, err = UnPublish.Execute(mb, publisher, actionName, ctx, mid); err != nil {
			return
		}
		presets.ShowMessage(&r, GetMessages(ctx.Context()).SuccessfullyUnpublished, "")
		r.Reload = true
		return
	}
}
