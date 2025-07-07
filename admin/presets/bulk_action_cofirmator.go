package presets

import (
	"errors"
	"fmt"
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
	BulkActionConfirmatorUpdateFunc               func(ctx *web.EventContext, r *web.EventResponse, mids model.IDSlice) (err error)

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

	var (
		records    = b.ab.l.mb.NewModelSlice()
		sp         data.SearchParams
		totalCount int
		comps      h.HTMLComponents
	)

	sp.Page = -1

	if len(mids) == 0 {
		if b.ab.allowEmpty {
			sp.MustCount = true
		} else {
			err = errors.New(MustGetMessages(ctx.Context()).ListingNoRecordToShow)
			return
		}
	} else {
		sp.WhereModelIDs(mids)
	}

	if records, totalCount, err = b.ab.l.Searcher(records, &sp, ctx); err != nil {
		return
	}

	msgr := MustGetMessages(ctx.Context())

	if sp.MustCount {
		comps = append(comps, v.VAlert().Type(v.TypeWarning).Text(msgr.ListingSelectedCountNoticeText(totalCount)))
	} else if reflect.ValueOf(records).Len() == 0 {
		err = errors.New(msgr.BulkActionNoAvailableRecords)
		return
	} else {
		if b.confirmationComp != nil {
			comp, err = b.confirmationComp(records, ctx)
			return
		}

		table := v.VDataTable(
			web.Slot(h.Text(`{{index+1}}`)).
				Name("item.index").
				Scope("{index}"),
		).
			PageText(msgr.PaginationPage).
			ItemsPerPageText(msgr.PaginationRowsPerPage).
			NoDataText(msgr.ListingNoRecordToShow)

		comp = table

		if b.confirmationTableHandler != nil {
			if comp, err = b.confirmationTableHandler(records, table, ctx); err != nil {
				return
			}
		} else {
			table.Items(records)
		}
	}

	var help h.HTMLComponent
	if b.helpComp != nil {
		if help, err = b.helpComp(ctx); err != nil {
			return
		}
	}

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
		h.Div(h.RawHTML(msgr.BulkActionConfirmationText(b.ab.RequestTitle(ctx.Context()), fmt.Sprint(totalCount)))).Class("my-3"),
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

func (b *BulkActionCofirmatorBuilder) Update(selectedIds []string, ctx *web.EventContext, r *web.EventResponse) (err error) {
	var (
		mids model.IDSlice
		msgr = MustGetMessages(ctx.Context())
	)

	for _, s := range selectedIds {
		if mid, err := b.ab.l.mb.ParseRecordID(s); err == nil {
			mids = append(mids, mid)
		}
	}

	if len(mids) == 0 && !b.ab.allowEmpty {
		err = errors.New(msgr.BulkActionNoAvailableRecords)
		return
	}

	return b.updator(ctx, r, mids)
}

func (b *BulkActionCofirmatorBuilder) Build() *BulkActionBuilder {
	return b.ab.
		ComponentFunc(func(selectedIds []string, ctx *web.EventContext) (comp h.HTMLComponent, err error) {
			return b.Component(selectedIds, ctx)
		}).
		UpdateFunc(func(selectedIds []string, ctx *web.EventContext, r *web.EventResponse) (err error) {
			return b.Update(selectedIds, ctx, r)
		})
}
