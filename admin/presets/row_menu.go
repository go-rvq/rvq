package presets

import (
	"fmt"
	"strconv"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/iancoleman/strcase"
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

func (b *RowMenuBuilder) SetRowMenuItem(name string) *RowMenuItemBuilder {
	return b.rowMenuItem(true, name)
}

func (b *RowMenuBuilder) RowMenuItem(name string) *RowMenuItemBuilder {
	return b.rowMenuItem(false, name)
}

func (b *RowMenuBuilder) rowMenuItem(set bool, name string) *RowMenuItemBuilder {
	if v, ok := b.items[strcase.ToSnake(name)]; ok {
		if set {
			panic("duplicated item " + strconv.Quote(name))
		}
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
		var mid ID
		if mid, err = b.mb.ParseRecordID(ctx.R.FormValue(ParamID)); err != nil {
			return
		}
		if ib.permAction != "" {
			if b.mb.permissioner.Actioner(ctx.R, ib.permAction, mid, ParentsModelID(ctx.R)...).Denied() {
				err = perm.PermissionDenied
				return
			}

			obj := b.mb.NewModel()
			err = b.mb.editing.Fetcher(obj, mid, ctx)
			if err != nil {
				return r, err
			}
		}
		if ib.clickF == nil {
			return r, nil
		}
		return ib.clickF(ctx, mid.String())
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
			id  = rctx.ID
			mid = b.rmb.mb.MustRecordID(rctx.Obj)
		)

		if b.permAction != "" && b.rmb.mb.permissioner.Actioner(ctx.R, b.permAction, mid, ParentsModelID(ctx.R)...).Denied() {
			return nil
		}

		return VListItem(
			web.Slot(
				VIcon(b.icon),
			).Name("prepend"),

			VListItemTitle(h.Text(b.rmb.mb.TFormat(ctx.Context(), "%sRowMenuItem%s", b.name))),
		).Attr("@click", web.Plaid().
			EventFunc(b.eventID).
			Query(ParamID, id).
			Go())
	}
}
