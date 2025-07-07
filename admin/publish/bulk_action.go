package publish

import (
	"errors"
	"fmt"
	"reflect"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/reflect_utils"
	"github.com/go-rvq/rvq/admin/utils/db_utils"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

type BulkAction struct {
	b                     *Builder
	m                     *presets.ModelBuilder
	Finder                func(db *gorm.DB, e Executor, selectedIds model.IDSlice) (records any, err error)
	ConfirmationComponent func(ctx *web.EventContext, records any) (comp h.HTMLComponent, err error)
	ConfirmationTable     func(records any, table *v.VDataTableBuilder, ctx *web.EventContext) (comp h.HTMLComponent, err error)
}

func (b *Builder) BulkAction(m *presets.ModelBuilder) *BulkAction {
	return &BulkAction{b: b, m: m}
}

func (ba *BulkAction) component(e Executor, selectedIds []string, ctx *web.EventContext) (comp h.HTMLComponent, err error) {
	var mids model.IDSlice
	for _, s := range selectedIds {
		if mid, err := ba.m.ParseRecordID(s); err == nil {
			mids = append(mids, mid)
		}
	}

	if len(mids) == 0 {
		err = errors.New(presets.MustGetMessages(ctx.Context()).ListingNoRecordToShow)
		return
	}

	var records any

	{
		db := ba.b.db.Session(&gorm.Session{})

		if ba.Finder != nil {
			if records, err = ba.Finder(db, e, mids); err != nil {
				return
			}
		} else {
			records = ba.m.NewModelSlice()
			db = db_utils.ModelIdsWhere(db, mids)
			if err = db.Find(records).Error; err != nil {
				return
			}

			records = reflect.ValueOf(records).Elem().Interface()

			records = reflect_utils.Filter(records, func(item any) bool {
				if s, ok := item.(StatusInterface); ok {
					return e.Accept(s.EmbedStatus().Status)
				}
				return true
			})
		}
	}

	msgr := GetMessages(ctx.Context())

	if reflect.ValueOf(records).Len() == 0 {
		err = errors.New(msgr.BulkActionNoRecordsText(e.ActivityName()))
		return
	}

	if ba.ConfirmationComponent != nil {
		comp, err = ba.ConfirmationComponent(ctx, records)
		return
	}

	table := v.VDataTable(
		web.Slot(h.Text(`{{index+1}}`)).Name("item.index").Scope("{index}"),
	)

	comp = table

	if ba.ConfirmationTable != nil {
		if comp, err = ba.ConfirmationTable(records, table, ctx); err != nil {
			return
		}
	} else {
		table.Items(records)
	}

	return h.HTMLComponents{
		v.VAlert(h.RawHTML(msgr.ActivityHelp(e.ActivityName()))).
			Type(v.ColorInfo).
			Variant(v.VariantTonal).
			Density(v.DensityCompact),
		h.Div(h.RawHTML(msgr.BulkActionConfirmationText(e.Title(msgr)))).Class("my-3"),
		comp,
	}, nil
}

func (ba *BulkAction) Update(e Executor, selectedIds []string, ctx *web.EventContext) (err error) {
	var (
		mids model.IDSlice
		msgr = presets.MustGetMessages(ctx.Context())
	)

	for _, s := range selectedIds {
		if mid, err := ba.m.ParseRecordID(s); err == nil {
			mids = append(mids, mid)
		}
	}

	if len(mids) == 0 {
		err = errors.New(msgr.ListingNoRecordToShow)
		return
	}

	for _, mid := range mids {
		if _, err = e.Execute(ba.m, ba.b, "Bulk"+e.ActivityName(), ctx, mid); err != nil {
			return
		}
	}

	return
}

func (ba *BulkAction) Build() *presets.BulkActionBuilder {
	const activityParam = "publish__activity"
	getExecutor := func(ctx *web.EventContext) (e Executor, err error) {
		activity := ctx.R.FormValue(activityParam)

		if e = ExecutorFromActivityName(activity); e == nil {
			err = fmt.Errorf("publish activity %q not found", activity)
		}

		return
	}

	l := ba.m.Listing()
	b := l.BulkAction("Publisher")
	return b.
		ButtonCompFunc(func(ctx *web.EventContext, title func() string, onclick *web.VueEventTagBuilder) h.HTMLComponent {
			var (
				msgr = GetMessages(ctx.Context())
				item = func(e Executor) h.HTMLComponent {
					activity := e.ActivityName()
					return v.VListItem(
						web.Slot(
							v.VBtn("",
								v.VIcon("mdi-information-outline"),
								v.VTooltip(h.Text(msgr.ActivityHelp(activity))).
									Attr("activator", "parent"),
							).
								Variant(v.VariantText).
								Size(v.SizeSmall).
								Color(v.ColorInfo).
								Attr("@click", `(e) => e.stopPropagation()`),
						).Name("append"),
					).Title(msgr.ActivityTitle(activity)).
						Attr("@click", b.DefaulButtonOnClick(ctx).Query(activityParam, activity).Go())
				}
			)

			return v.VBtn(
				"",
				h.Text(msgr.Publication),
				v.VMenu(
					v.VList(
						item(Publish),
						item(PublishOrRepublish),
						item(RePublish),
						item(UnPublish),
					),
				).
					Attr("activator", "parent")).
				PrependIcon("mdi-web").
				AppendIcon("mdi-menu-down")
		}).
		SetComponentHandler(func(cb *presets.ContentComponentBuilder, ctx *web.EventContext) {
			e, _ := getExecutor(ctx)
			if e != nil {
				cb.Title = GetMessages(ctx.Context()).ActivityTitle(e.ActivityName())
			}
		}).
		ComponentFunc(func(selectedIds []string, ctx *web.EventContext) (comp h.HTMLComponent, err error) {
			var e Executor
			if e, err = getExecutor(ctx); err != nil {
				return
			}
			return ba.component(e, selectedIds, ctx)
		}).
		UpdateFunc(func(selectedIds []string, ctx *web.EventContext, r *web.EventResponse) (err error) {
			var e Executor
			if e, err = getExecutor(ctx); err != nil {
				return
			}

			return ba.Update(e, selectedIds, ctx)
		})
}
