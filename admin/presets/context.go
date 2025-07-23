package presets

import (
	"context"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/utils"
	"github.com/go-rvq/rvq/web"
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
	ModelIDKey
	ParentsModelIDKey
	parentsRecordKey
	CtxEventHandlerWrapperNoFlash
	CtxSkipAutoBreadcrumb
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

func WithActionsComponent(ctx web.ContextValuer, comp ...h.HTMLComponent) {
	ctx.WithContextValue(CtxActionsComponent, h.HTMLComponents(comp))
}

func GetActionsComponent(ctx web.ContextValuer) h.HTMLComponents {
	v, _ := ctx.ContextValue(CtxActionsComponent).(h.HTMLComponents)
	return v
}

func WithMenuComponent(ctx web.ContextValuer, comp ...h.HTMLComponent) {
	ctx.WithContextValue(CtxMenuComponent, h.HTMLComponents(comp))
}

func GetMenuComponent(ctx web.ContextValuer) h.HTMLComponents {
	v, _ := ctx.ContextValue(CtxMenuComponent).(h.HTMLComponents)
	return v
}

func GetComponentFromContext(ctx web.ContextValuer, key presetsCtx) (h.HTMLComponent, bool) {
	v, ok := ctx.ContextValue(key).(h.HTMLComponent)
	return v, ok
}

func GetPortals(ctx web.ContextValuer) h.HTMLComponents {
	v, _ := ctx.ContextValue(CtxPortals).(h.HTMLComponents)
	return v
}

func WithPortals(ctx web.ContextValuer, portal ...h.HTMLComponent) {
	vlr := web.GetContextValuer(ctx.Context(), CtxPortals)
	if vlr == nil {
		ctx.WithContextValue(CtxPortals, h.HTMLComponents(portal))
	} else {
		vlr.Set(h.HTMLComponents(portal))
	}
}

func AddPortals(ctx web.ContextValuer, portal ...h.HTMLComponent) {
	vlr := web.GetContextValuer(ctx.Context(), CtxPortals)
	if vlr == nil {
		ctx.WithContextValue(CtxPortals, h.HTMLComponents(portal))
	} else {
		vlr.Set(append(vlr.Get().(h.HTMLComponents), portal...))
	}
}

func GetModel(ctx web.ContextValuer) (v *ModelBuilder) {
	v, _ = ctx.ContextValue(ctxModel).(*ModelBuilder)
	return
}

func WithModel(ctx web.ContextValuer, model *ModelBuilder) {
	ctx.WithContextValue(ctxModel, model)
}

func WithParentsRecord(ctx web.ContextValuer, parents []any) {
	ctx.WithContextValue(parentsRecordKey, parents)
}

func GetParentsRecord(ctx web.ContextValuer) utils.Anies {
	v, _ := ctx.ContextValue(parentsRecordKey).([]any)
	return v
}

func WithModelAndLoadBreadcrumbsAndParents(ctx *web.EventContext, mb *ModelBuilder) (err error) {
	if err = mb.LoadParentsID(ctx); err != nil {
		return
	}
	if !IsSkipAutoBreadcrumb(ctx) {
		var records []any
		if records, err = mb.LoadBreadCrumbs(ctx); err != nil {
			return
		}
		WithParentsRecord(ctx, records)
	}
	return
}

func WithScope(ctx web.ContextValuer, scope *web.ScopeBuilder) {
	ctx.WithContextValue(ctxScope, scope)
}

func GetScope(ctx web.ContextValuer) (scope *web.ScopeBuilder) {
	scope, _ = ctx.ContextValue(ctxScope).(*web.ScopeBuilder)
	return
}

func EditFormUnscoped(ctx web.ContextValuer, v bool) {
	ctx.WithContextValue(ctxEditFormUnscoped, v)
}

func GetEditFormUnscoped(ctx web.ContextValuer) (ok bool) {
	ok, _ = ctx.ContextValue(ctxEditFormUnscoped).(bool)
	return
}

func WithRespondDialogHandlers(ctx web.ContextValuer, f ...func(d *DialogBuilder)) {
	ctx.WithContextValue(CtxRespondDialogHandlers, append(GetRespondDialogHandlers(ctx), f...))
}

func GetRespondDialogHandlers(ctx web.ContextValuer) (handlers []func(d *DialogBuilder)) {
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

func GetSaveContext(ctx web.ContextValuer) (v context.Context) {
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

func WithEventHandlerWrapperNoFlash(ctx web.ContextValuer) {
	ctx.WithContextValue(CtxEventHandlerWrapperNoFlash, true)
}

func IsSkipAutoBreadcrumb(ctx web.ContextValuer) (ok bool) {
	ok, _ = ctx.ContextValue(CtxSkipAutoBreadcrumb).(bool)
	return
}

func WithSkipAutoBreadcrumb(ctx web.ContextValuer) (ok bool) {
	ctx.WithContextValue(CtxSkipAutoBreadcrumb, true)
	return
}
