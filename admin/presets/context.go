package presets

import (
	"context"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	h "github.com/theplant/htmlgo"
)

type presetsCtx int

const (
	ctxInDialog presetsCtx = iota
	CtxActionsComponent
	CtxMenuComponent
	CtxPortals
	ctxDetailingAfterTitleComponent
	ctxModel
	ctxScope
	ctxEditFormUnscoped
	ctxFlashMessages
	ctxActionFormContext
	ctxBulkActionFormContext
	CtxRespondDialogHandlers
	ctxFieldLabels
	CtxSaveContext
	ParentsModelIDKey
	CtxEventHandlerWrapperNoFlash
)

func IsInDialog(ctx *web.EventContext) bool {
	return actions.OverlayMode(ctx.R.FormValue(ParamOverlay)).IsDialog()
}

func IsInDrawer(ctx *web.EventContext) bool {
	return actions.OverlayMode(ctx.R.FormValue(ParamOverlay)).IsDrawer()
}

func OverlayMode(ctx *web.EventContext) actions.OverlayMode {
	return actions.OverlayMode(ctx.R.FormValue(ParamOverlay))
}

func WithActionsComponent(ctx *web.EventContext, comp ...h.HTMLComponent) {
	ctx.WithContextValue(CtxActionsComponent, h.HTMLComponents(comp))
}

func GetActionsComponent(ctx *web.EventContext) h.HTMLComponents {
	v, _ := ctx.ContextValue(CtxActionsComponent).(h.HTMLComponents)
	return v
}

func WithMenuComponent(ctx *web.EventContext, comp ...h.HTMLComponent) {
	ctx.WithContextValue(CtxMenuComponent, h.HTMLComponents(comp))
}

func GetMenuComponent(ctx *web.EventContext) h.HTMLComponents {
	v, _ := ctx.ContextValue(CtxMenuComponent).(h.HTMLComponents)
	return v
}

func GetComponentFromContext(ctx *web.EventContext, key presetsCtx) (h.HTMLComponent, bool) {
	v, ok := ctx.ContextValue(key).(h.HTMLComponent)
	return v, ok
}

func GetPortals(ctx *web.EventContext) h.HTMLComponents {
	v, _ := ctx.ContextValue(CtxPortals).(h.HTMLComponents)
	return v
}

func WithPortals(ctx *web.EventContext, portal ...h.HTMLComponent) {
	vlr := web.GetContextValuer(ctx.R.Context(), CtxPortals)
	if vlr == nil {
		ctx.WithContextValue(CtxPortals, h.HTMLComponents(portal))
	} else {
		vlr.Set(h.HTMLComponents(portal))
	}
}

func AddPortals(ctx *web.EventContext, portal ...h.HTMLComponent) {
	vlr := web.GetContextValuer(ctx.R.Context(), CtxPortals)
	if vlr == nil {
		ctx.WithContextValue(CtxPortals, h.HTMLComponents(portal))
	} else {
		vlr.Set(append(vlr.Get().(h.HTMLComponents), portal...))
	}
}

func GetModel(ctx *web.EventContext) (v *ModelBuilder) {
	v, _ = ctx.ContextValue(ctxModel).(*ModelBuilder)
	return
}

func WithModel(ctx *web.EventContext, model *ModelBuilder) {
	ctx.WithContextValue(ctxModel, model)
}

func WithScope(ctx *web.EventContext, scope *web.ScopeBuilder) {
	ctx.WithContextValue(ctxScope, scope)
}

func GetScope(ctx *web.EventContext) (scope *web.ScopeBuilder) {
	scope, _ = ctx.ContextValue(ctxScope).(*web.ScopeBuilder)
	return
}

func EditFormUnscoped(ctx *web.EventContext, v bool) {
	ctx.WithContextValue(ctxEditFormUnscoped, v)
}

func GetEditFormUnscoped(ctx *web.EventContext) (ok bool) {
	ok, _ = ctx.ContextValue(ctxEditFormUnscoped).(bool)
	return
}

func WithRespondDialogHandlers(ctx *web.EventContext, f ...func(d *DialogBuilder)) {
	ctx.WithContextValue(CtxRespondDialogHandlers, append(GetRespondDialogHandlers(ctx), f...))
}

func GetRespondDialogHandlers(ctx *web.EventContext) (handlers []func(d *DialogBuilder)) {
	handlers, _ = ctx.ContextValue(CtxRespondDialogHandlers).([]func(d *DialogBuilder))
	return
}

func WithFieldLabels(ctx web.ContextValuer, fb *FieldsBuilder, labels map[string]string) {
	v, _ := ctx.ContextValue(ctxFieldLabels).(map[*FieldsBuilder]map[string]string)
	if v != nil {
		v[fb] = labels
	} else {
		ctx.WithContextValue(ctxFieldLabels, map[*FieldsBuilder]map[string]string{fb: labels})
	}
}

func GetFieldLabels(ctx web.ContextValuer, fb *FieldsBuilder) map[string]string {
	v, _ := ctx.ContextValue(ctxFieldLabels).(map[*FieldsBuilder]map[string]string)
	if v != nil {
		return v[fb]
	}
	return nil
}

func GetSaveContext(ctx *web.EventContext) (v context.Context) {
	v, _ = ctx.ContextValue(CtxSaveContext).(context.Context)
	return
}

func GetActionFormContext[T any](ctx web.ContextValuer) *ActionFormContext[T] {
	v, _ := ctx.ContextValue(ctxActionFormContext).(*ActionFormContext[T])
	return v
}

func WithActionFormContext[T any](ctx web.ContextValuer, v *ActionFormContext[T]) {
	ctx.WithContextValue(ctxActionFormContext, v)
}

func GetBulkActionFormContext[T any](ctx web.ContextValuer) *BulkActionFormContext[T] {
	v, _ := ctx.ContextValue(ctxBulkActionFormContext).(*BulkActionFormContext[T])
	return v
}

func WithBulkActionFormContext[T any](ctx web.ContextValuer, v *BulkActionFormContext[T]) {
	ctx.WithContextValue(ctxBulkActionFormContext, v)
}

func WithEventHandlerWrapperNoFlash(ctx *web.EventContext) {
	ctx.WithContextValue(CtxEventHandlerWrapperNoFlash, true)
}
