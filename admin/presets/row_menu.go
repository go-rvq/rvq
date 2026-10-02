package presets

import (
	"context"
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
	title      func(ctx context.Context) string
	child      *ModelBuilder
}

// Items are the items of the menu, in the order of the menu.
func (b *RowMenuBuilder) Items() (items []*RowMenuItemBuilder) {
	listings := b.defaultListings
	if len(b.listings) > 0 {
		listings = b.listings
	}
	for _, li := range listings {
		if ib, ok := b.items[strcase.ToSnake(li)]; ok {
			items = append(items, ib)
		}
	}
	return
}

// Name is the name of the item.
func (b *RowMenuItemBuilder) Name() string { return b.name }

// Title sets the title of the item, for whoever names it out of the menu (the
// documentation): its component shows its own.
func (b *RowMenuItemBuilder) Title(f func(ctx context.Context) string) *RowMenuItemBuilder {
	b.title = f
	return b
}

// TTitle is the title of the item: Title's, else its name.
func (b *RowMenuItemBuilder) TTitle(ctx context.Context) string {
	if b.title != nil {
		return b.title(ctx)
	}
	return HumanizeString(b.name)
}

// deleteTitle is the title of the item that deletes a record.
func deleteTitle(ctx context.Context) string { return MustGetMessages(ctx).Delete }

// Child is the nested model the item opens, nil when it is not one.
func (b *RowMenuItemBuilder) Child() *ModelBuilder { return b.child }

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
