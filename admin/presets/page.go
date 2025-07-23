package presets

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
)

type PageBuilder struct {
	r       *PagesRegistrator
	page    *HttpPageBuilder
	events  web.EventsHub
	actions []*ActionBuilder
	layout  *LayoutConfig
	pf      web.PageFunc
	handler func(f web.PageFunc) http.Handler
}

func Page(r *PagesRegistrator, page *HttpPageBuilder) *PageBuilder {
	return &PageBuilder{
		r:    r,
		page: page,
	}
}

func (b *PageBuilder) Page() *HttpPageBuilder {
	return b.page
}

func (b *PageBuilder) Title(f func(ctx context.Context) string) *PageBuilder {
	b.page.titleFunc = f
	return b
}

func (b *PageBuilder) SubTitle(f func(ctx context.Context) string) *PageBuilder {
	b.page.subTitleFunc = f
	return b
}

func (b *PageBuilder) InMenu(v bool) *PageBuilder {
	b.page.InMenu(v)
	return b
}

func (b *PageBuilder) MenuIcon(v string) *PageBuilder {
	b.page.menuIcon = v
	return b
}

func (b *PageBuilder) EventHandler(eventFuncId string, ef web.EventHandler) *PageBuilder {
	b.events.RegisterEventHandler(eventFuncId, ef)
	return b
}

func (b *PageBuilder) EventFunc(eventFuncId string, ef web.EventFunc) *PageBuilder {
	return b.EventHandler(eventFuncId, ef)
}

func (b *PageBuilder) PrivateEventFunc(eventFuncId string, ef web.EventFunc, vf ...func(v *perm.Verifier) *perm.Verifier) *PageBuilder {
	return b.EventFunc(eventFuncId, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		v := b.page.Verifier(ctx.R)
		for _, f := range vf {
			v = f(v)
		}
		if v.Denied() {
			err = perm.PermissionDenied
			return
		}
		return ef(ctx)
	})
}

func (b *PageBuilder) Raw(pf web.PageFunc) *PageBuilder {
	b.pf = pf
	b.handler = func(f web.PageFunc) http.Handler {
		return b.r.Wrap(f, b.buildPage)
	}
	return b
}

func (b *PageBuilder) Layout(pf web.PageFunc, layout ...*LayoutConfig) *PageBuilder {
	var lc *LayoutConfig
	if len(layout) > 0 {
		lc = layout[0]
	}
	b.pf = pf
	b.handler = func(f web.PageFunc) http.Handler {
		return b.r.Wrap(b.r.layoutFunc(lc, func(ctx *web.EventContext) (r web.PageResponse, err error) {
			if r, err = f(ctx); err == nil {
				if r.PageTitle == "" {
					r.PageTitle = b.page.TTitle(ctx.Context())
				}
			}
			return
		}), b.buildPage)
	}
	return b
}

func (b *PageBuilder) RecordIDLayout(pf func(ctx *web.EventContext, mid model.ID, parents []model.ID) (r web.PageResponse, err error), layout ...*LayoutConfig) *PageBuilder {
	var lc *LayoutConfig
	if len(layout) > 0 {
		lc = layout[0]
	}
	b.pf = func(ctx *web.EventContext) (r web.PageResponse, err error) {
		return pf(ctx, ctx.ContextValue(ModelIDKey).(model.ID), ctx.ContextValue(ParentsModelIDKey).(model.IDSlice))
	}
	b.handler = func(f web.PageFunc) http.Handler {
		return b.r.Wrap(b.r.recordIdLayoutFunc(lc, func(ctx *web.EventContext) (r web.PageResponse, err error) {
			if r, err = f(ctx); err == nil {
				if r.PageTitle == "" {
					r.PageTitle = b.page.TTitle(ctx.Context())
				}
			}
			return
		}), b.buildPage)
	}
	return b
}

func (b *PageBuilder) RecordLayout(f func(ctx *web.EventContext, mid model.ID, parents []model.ID, fetch func() (any, error)) (r web.PageResponse, err error), layout ...*LayoutConfig) *PageBuilder {
	return b.Layout(b.r.Record(f), layout...)
}

func (b *PageBuilder) parseRequestAction(ctx *web.EventContext) (id string, action *ActionBuilder, err error) {
	action = getAction(b.actions, ctx.R.FormValue(ParamAction))
	if action == nil {
		err = errors.New("action required")
		return
	}

	var enabled bool
	if enabled, err = action.IsEnabled(id, ctx); err != nil {
		return
	}

	if !enabled {
		err = errors.New("action disabled")
		return
	}
	return
}

func (b *PageBuilder) actionForm(ctx *web.EventContext) (r web.EventResponse, err error) {
	var action *ActionBuilder
	if _, action, err = b.parseRequestAction(ctx); err != nil {
		return
	}
	err = action.View(nil, "", ctx, &r)
	return
}

func (b *PageBuilder) doAction(ctx *web.EventContext) (r web.EventResponse, err error) {
	var action *ActionBuilder
	if _, action, err = b.parseRequestAction(ctx); err != nil {
		return
	}

	_, err = action.Do(nil, "", ctx, &r)
	return
}

func (b *PageBuilder) buildPage(p *web.PageBuilder) {
	b.EventFunc(actions.Action, b.actionForm)
	b.EventFunc(actions.DoAction, b.doAction)
	p.MergeHub((&b.events).Wrap(WrapEventHandler))
}

func (b *PageBuilder) Action(name string) *ActionBuilder {
	action := Action(b.r.b, name)
	action.typ = ActionTypePage
	b.actions = append(b.actions, action)
	return action
}

func (b *PageBuilder) PrivateAction(name string, vf ...func(v *perm.Verifier) *perm.Verifier) *ActionBuilder {
	return b.Action(name).SetVerifier(func(ctx *web.EventContext) *perm.Verifier {
		v := b.page.Verifier(ctx.R)
		for _, f := range vf {
			v = f(v)
		}
		return v
	})
}

func (b *PageBuilder) Private(md ...RequestPermVerifierMiddleware) *PageBuilder {
	v := b.r.baseVerifier
	if len(md) > 0 {
		for _, md := range md {
			v = md(v)
		}
	}
	b.page.AutoPerm().BaseVerifier(v)
	b.page.PostBuild(func(*PageHandler) {
		b.page.BaseVerifier(v)
	})
	return b
}

func (b *PageBuilder) Build() {
	if b.pf != nil {
		b.page.HandlerFromPageFunc(b.handler, b.pf)
	}
}
