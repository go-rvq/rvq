package presets

import (
	"errors"
	"reflect"

	"github.com/qor5/admin/v3/model"
	"github.com/qor5/admin/v3/presets/data"
	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type (
	BulkActionConfirmatorSearchParamsHandler      func(sp *data.SearchParams, ctx *web.EventContext) (err error)
	BulkActionConfirmatorConfirmationComponent    func(records any, ctx *web.EventContext) (comp h.HTMLComponent, err error)
	BulkActionConfirmatorConfirmationTableHandler func(records any, table *v.VDataTableBuilder, ctx *web.EventContext) (comp h.HTMLComponent, err error)
	BulkActionConfirmatorHelpComponent            func(ctx *web.EventContext) (comp h.HTMLComponent, err error)
	BulkActionConfirmatorUpdateFunc               func(ctx *web.EventContext, mids model.IDSlice) (err error)

	BulkActionCofirmatorBuilder struct {
		ab                       *BulkActionBuilder
		searchParamsHandler      BulkActionConfirmatorSearchParamsHandler
		confirmationComp         BulkActionConfirmatorConfirmationComponent
		confirmationTableHandler BulkActionConfirmatorConfirmationTableHandler
		helpComp                 BulkActionConfirmatorHelpComponent
		preComp                  BulkActionConfirmatorHelpComponent
		postComp                 BulkActionConfirmatorHelpComponent
		wrapComp                 func(comp h.HTMLComponent, ctx *web.EventContext) h.HTMLComponent
		updator                  BulkActionConfirmatorUpdateFunc
	}
)

func BulkActionConfirmator(ab *BulkActionBuilder) *BulkActionCofirmatorBuilder {
	return &BulkActionCofirmatorBuilder{ab: ab}
}

func (b *BulkActionCofirmatorBuilder) SearchParamsHandler() BulkActionConfirmatorSearchParamsHandler {
	return b.searchParamsHandler
}

func (b *BulkActionCofirmatorBuilder) SetSearchParamsHandler(searcherHandler BulkActionConfirmatorSearchParamsHandler) *BulkActionCofirmatorBuilder {
	b.searchParamsHandler = searcherHandler
	return b
}

func (b *BulkActionCofirmatorBuilder) ConfirmationComp() BulkActionConfirmatorConfirmationComponent {
	return b.confirmationComp
}

func (b *BulkActionCofirmatorBuilder) SetConfirmationComp(confirmationComp BulkActionConfirmatorConfirmationComponent) *BulkActionCofirmatorBuilder {
	b.confirmationComp = confirmationComp
	return b
}

func (b *BulkActionCofirmatorBuilder) ConfirmationTableHandler() BulkActionConfirmatorConfirmationTableHandler {
	return b.confirmationTableHandler
}

func (b *BulkActionCofirmatorBuilder) SetConfirmationTableHandler(confirmationTableHandler BulkActionConfirmatorConfirmationTableHandler) *BulkActionCofirmatorBuilder {
	b.confirmationTableHandler = confirmationTableHandler
	return b
}

func (b *BulkActionCofirmatorBuilder) HelpComp() BulkActionConfirmatorHelpComponent {
	return b.helpComp
}

func (b *BulkActionCofirmatorBuilder) SetHelpComp(helpComp BulkActionConfirmatorHelpComponent) *BulkActionCofirmatorBuilder {
	b.helpComp = helpComp
	return b
}

func (b *BulkActionCofirmatorBuilder) PreComp() BulkActionConfirmatorHelpComponent {
	return b.preComp
}

func (b *BulkActionCofirmatorBuilder) SetPreComp(preComp BulkActionConfirmatorHelpComponent) *BulkActionCofirmatorBuilder {
	b.preComp = preComp
	return b
}

func (b *BulkActionCofirmatorBuilder) PostComp() BulkActionConfirmatorHelpComponent {
	return b.postComp
}

func (b *BulkActionCofirmatorBuilder) SetPostComp(postComp BulkActionConfirmatorHelpComponent) *BulkActionCofirmatorBuilder {
	b.postComp = postComp
	return b
}

func (b *BulkActionCofirmatorBuilder) WrapComp() func(comp h.HTMLComponent, ctx *web.EventContext) h.HTMLComponent {
	return b.wrapComp
}

func (b *BulkActionCofirmatorBuilder) SetWrapComp(wrapComp func(comp h.HTMLComponent, ctx *web.EventContext) h.HTMLComponent) *BulkActionCofirmatorBuilder {
	b.wrapComp = wrapComp
	return b
}

func (b *BulkActionCofirmatorBuilder) Updator() BulkActionConfirmatorUpdateFunc {
	return b.updator
}

func (b *BulkActionCofirmatorBuilder) SetUpdator(updator BulkActionConfirmatorUpdateFunc) *BulkActionCofirmatorBuilder {
	b.updator = updator
	return b
}

func (b *BulkActionCofirmatorBuilder) Component(selectedIds []string, ctx *web.EventContext) (comp h.HTMLComponent, err error) {
	var mids model.IDSlice
	for _, s := range selectedIds {
		if mid, err := b.ab.l.mb.ParseRecordID(s); err == nil {
			mids = append(mids, mid)
		}
	}

	if len(mids) == 0 {
		err = errors.New(MustGetMessages(ctx.Context()).ListingNoRecordToShow)
		return
	}

	var (
		records = b.ab.l.mb.NewModelSlice()
		sp      data.SearchParams
	)

	sp.Page = -1
	sp.WhereModelIDs(mids)

	if records, _, err = b.ab.l.Searcher(records, &sp, ctx); err != nil {
		return
	}

	msgr := MustGetMessages(ctx.Context())

	if reflect.ValueOf(records).Len() == 0 {
		err = errors.New(msgr.BulkActionNoAvailableRecords)
		return
	}

	if b.confirmationComp != nil {
		comp, err = b.confirmationComp(records, ctx)
		return
	}

	table := v.VDataTable(
		web.Slot(h.Text(`{{index+1}}`)).Name("item.index").Scope("{index}"),
	)

	comp = table

	if b.confirmationTableHandler != nil {
		if comp, err = b.confirmationTableHandler(records, table, ctx); err != nil {
			return
		}
	} else {
		table.Items(records)
	}

	var help h.HTMLComponent
	if b.helpComp != nil {
		if help, err = b.helpComp(ctx); err != nil {
			return
		}
	}

	var comps h.HTMLComponents

	if b.preComp != nil {
		var c h.HTMLComponent
		if c, err = b.preComp(ctx); err != nil {
			return
		}
		comps = append(comps, c)
	}

	comps = append(comps,
		h.If(help != nil,
			v.VAlert(help).
				Type(v.ColorInfo).
				Variant(v.VariantTonal).
				Density(v.DensityCompact),
		),
		h.Div(h.RawHTML(msgr.BulkActionConfirmationText(b.ab.RequestTitle(ctx)))).Class("my-3"),
		comp,
	)

	if b.postComp != nil {
		var c h.HTMLComponent
		if c, err = b.postComp(ctx); err != nil {
			return
		}
		comps = append(comps, c)
	}

	comp = comps

	if b.wrapComp != nil {
		comp = b.wrapComp(comp, ctx)
	}

	return
}

func (b *BulkActionCofirmatorBuilder) Update(selectedIds []string, ctx *web.EventContext) (err error) {
	var (
		mids model.IDSlice
		msgr = MustGetMessages(ctx.Context())
	)

	for _, s := range selectedIds {
		if mid, err := b.ab.l.mb.ParseRecordID(s); err == nil {
			mids = append(mids, mid)
		}
	}

	if len(mids) == 0 {
		err = errors.New(msgr.BulkActionNoAvailableRecords)
		return
	}

	if err = b.updator(ctx, mids); err != nil {
		return
	}

	return
}

func (b *BulkActionCofirmatorBuilder) Build() {
	b.ab.
		ComponentFunc(func(selectedIds []string, ctx *web.EventContext) (comp h.HTMLComponent, err error) {
			return b.Component(selectedIds, ctx)
		}).
		UpdateFunc(func(selectedIds []string, ctx *web.EventContext) (err error) {
			return b.Update(selectedIds, ctx)
		})
}
