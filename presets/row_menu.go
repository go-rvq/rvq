package presets

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	. "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type RowMenuFields struct {
	rowMenu *RowMenuBuilder
}

type RowMenuBuilder struct {
	mb              *ModelBuilder
	listings        []string
	defaultListings []string
	items           map[string]*RowMenuItemBuilder
}

func (b *RowMenuFields) init(mb *ModelBuilder) *RowMenuBuilder {
	b.rowMenu = &RowMenuBuilder{
		mb:    mb,
		items: make(map[string]*RowMenuItemBuilder),
	}
	return b.rowMenu
}

func (b *RowMenuFields) RowMenu(listings ...string) *RowMenuBuilder {
	rmb := b.rowMenu
	if len(listings) == 0 {
		return rmb
	}
	rmb.listings = listings
	for _, li := range rmb.listings {
		rmb.RowMenuItem(li)
	}

	return rmb
}

func (b *RowMenuBuilder) Empty() {
	b.listings = nil
	b.items = make(map[string]*RowMenuItemBuilder)
}

func (b *RowMenuBuilder) listingItemFuncs(ctx *web.EventContext) (fs RecordMenuItemFuncs) {
	listings := b.defaultListings
	if len(b.listings) > 0 {
		listings = b.listings
	}
	for _, li := range listings {
		if ib, ok := b.items[strcase.ToSnake(li)]; ok {
			comp := ib.getComponentFunc(ctx)
			if comp != nil {
				fs = append(fs, comp)
			}
		}
	}
	return fs
}

type RowMenuItemBuilder struct {
	rmb        *RowMenuBuilder
	name       string
	icon       string
	clickF     RowMenuItemClickFunc
	compF      RecordMenuItemFunc
	permAction string
	eventID    string
}

func (b *RowMenuBuilder) RowMenuItem(name string) *RowMenuItemBuilder {
	if v, ok := b.items[strcase.ToSnake(name)]; ok {
		return v
	}

	ib := &RowMenuItemBuilder{
		rmb:     b,
		name:    name,
		eventID: fmt.Sprintf("%s_rowMenuItemFunc_%s", b.mb.uriName, name),
	}
	b.items[strcase.ToSnake(name)] = ib
	b.defaultListings = append(b.defaultListings, name)

	b.mb.RegisterEventFunc(ib.eventID, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id := ctx.R.FormValue(ParamID)
		if ib.permAction != "" {
			obj := b.mb.NewModel()
			err = b.mb.editing.Fetcher(obj, id, ctx)
			if err != nil {
				return r, err
			}
			err = b.mb.Info().Verifier().Do(ib.permAction).ObjectOn(obj).WithReq(ctx.R).IsAllowed()
			if err != nil {
				return r, err
			}
		}
		if ib.clickF == nil {
			return r, nil
		}
		return ib.clickF(ctx, id)
	})

	return ib
}

func (b *RowMenuItemBuilder) Icon(v string) *RowMenuItemBuilder {
	b.icon = v
	return b
}

type RowMenuItemClickFunc func(ctx *web.EventContext, id string) (r web.EventResponse, err error)

func (b *RowMenuItemBuilder) OnClick(v RowMenuItemClickFunc) *RowMenuItemBuilder {
	b.clickF = v
	return b
}

func (b *RowMenuItemBuilder) ComponentFunc(v RecordMenuItemFunc) *RowMenuItemBuilder {
	b.compF = v
	return b
}

func (b *RowMenuItemBuilder) PermAction(v string) *RowMenuItemBuilder {
	b.permAction = v
	return b
}

func (b *RowMenuItemBuilder) getComponentFunc(_ *web.EventContext) RecordMenuItemFunc {
	if b.compF != nil {
		return b.compF
	}

	return func(rctx *RecordMenuItemContext) h.HTMLComponent {
		var (
			ctx = rctx.Ctx
			obj = rctx.Obj
			id  = rctx.ID
		)

		if b.permAction != "" && b.rmb.mb.Info().Verifier().Do(b.permAction).ObjectOn(obj).WithReq(ctx.R).IsAllowed() != nil {
			return nil
		}
		return VListItem(
			web.Slot(
				VIcon(b.icon),
			).Name("prepend"),

			VListItemTitle(h.Text(i18n.PT(ctx.R, ModelsI18nModuleKey, strcase.ToCamel(b.rmb.mb.label+" RowMenuItem"), b.name))),
		).Attr("@click", web.Plaid().
			EventFunc(b.eventID).
			Query(ParamID, id).
			Go())
	}
}
