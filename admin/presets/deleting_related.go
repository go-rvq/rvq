package presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

type RelatedDeletionContext struct {
	Builder             *ListingBuilder
	Context             *web.EventContext
	FormKey             string
	Obj                 any
	Enabled             bool
	Description         string
	WrapComponent       func(comp h.HTMLComponent) h.HTMLComponent
	ShowRelatedDisabled bool
}

type ShowRelatedItensForDeletionContext struct {
	Response *web.EventResponse
	Context  *web.EventContext
	Obj      any
	ID       model.ID
	Dialog   *vx.VXDialogBuilder
}

type RelatedDeletionFunc func(ctx *RelatedDeletionContext) (err error)

type ShowRelatedItensFormDeletionFunc func(ctx *ShowRelatedItensForDeletionContext) (body h.HTMLComponent, err error)

func (b *ListingBuilder) showRelatedItensForDeletion(ctx *web.EventContext) (r web.EventResponse, err error) {
	targetPortal := ctx.R.FormValue(ParamTargetPortal)

	var (
		id   = ctx.R.FormValue(ParamID)
		msgr = MustGetMessages(ctx.Context())
		obj  = b.mb.NewModel()
		ido  ID
		body h.HTMLComponent
	)

	if ido, err = b.mb.ParseRecordID(id); err != nil {
		return
	}

	ido.SetTo(obj)

	srCtx := &ShowRelatedItensForDeletionContext{
		Response: &r,
		Context:  ctx,
		Obj:      obj,
		ID:       ido,
		Dialog: vx.VXDialog().
			Density(DensityCompact).
			SlotBody(body).
			Title(msgr.RelatedItemsForDeletionActionTitle),
	}

	if body, err = b.showRelatedItensForDeletionActionFunc(srCtx); err != nil {
		return
	}

	if body != nil {
		srCtx.Dialog.SetSlotBody(body)
	}

	b.mb.p.DialogPortal(targetPortal).
		Respond(ctx, &r, srCtx.Dialog)
	return
}

func (b *ListingBuilder) ShowRelatedItensForDeletionActionFunc(f ShowRelatedItensFormDeletionFunc) *ListingBuilder {
	b.showRelatedItensForDeletionActionFunc = f
	return b
}
