package presets

import (
	"fmt"
	"net/url"
	_ "unsafe"

	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

func EditDeleteRowMenuItemFuncs(mi *ModelInfo, url string, editExtraParams url.Values) []RecordMenuItemFunc {
	return []RecordMenuItemFunc{
		editRowMenuItemFunc(mi, url, editExtraParams),
		NewDeletingMenuItemBuilder(mi).SetUrl(url).SetUrlValues(editExtraParams).Build(),
	}
}

func editRowMenuItemFunc(mi *ModelInfo, url string, editExtraParams url.Values) RecordMenuItemFunc {
	return func(rctx *RecordMenuItemContext) h.HTMLComponent {
		var (
			ctx = rctx.Ctx
			obj = rctx.Obj
			id  = rctx.ID
		)

		msgr := MustGetMessages(ctx.Context())
		if mi.mb.permissioner.ReqObjectUpdater(ctx.R, obj).Denied() {
			return nil
		}

		onclick := web.Plaid().
			EventFunc(actions.Edit).
			Queries(editExtraParams).
			Query(ParamID, id).
			URL(url)
		if IsInDialog(ctx) {
			onclick.URL(ctx.R.RequestURI).
				Query(ParamOverlay, actions.Dialog).
				Query(ParamListingQueries, ctx.Queries().Encode())
		}
		return VListItem(
			web.Slot(
				VIcon("mdi-pencil"),
			).Name("prepend"),

			VListItemTitle(h.Text(msgr.Edit)),
		).Attr("@click", onclick.Go())
	}
}

func childRowMenuItemFunc(mb *ModelBuilder) RecordMenuItemFunc {
	return func(ctx *RecordMenuItemContext) h.HTMLComponent {
		if mb.notInMenu {
			return nil
		}

		var (
			r         = ctx.Ctx.R
			mi        = mb.modelInfo
			parentsID = append(ParentsModelID(ctx.Ctx.R), mb.parent.MustParseRecordID(ctx.ID))
		)

		if mb.singleton {
			if mb.hasDetailing {
				if mb.permissioner.Reader(r, ID{}, parentsID...).Denied() {
					return nil
				}
			} else {
				if mb.permissioner.Updater(r, ID{}, parentsID...).Denied() {
					return nil
				}
			}
		} else if mb.permissioner.Lister(ctx.Ctx.R, parentsID...).Denied() {
			return nil
		}

		var (
			event string
			title = mb.TTitleAuto(ctx.Ctx.Context())
			uri   = mi.ListingHref(parentsID...)
		)

		if mb.singleton {
			if mb.hasDetailing {
				event = actions.Detailing
			} else {
				event = actions.Edit
			}
		} else {
			event = actions.OpenListingDialog
		}

		return VListItem(
			web.Slot(
				VIcon(mi.mb.menuIcon),
			).Name("prepend"),
			VListItemTitle(h.Text(title)),
		).Attr("@click", web.Plaid().
			EventFunc(event).
			Query(ParamTargetPortal, ctx.TempPortal).
			Query(ParamOverlay, actions.Dialog).
			URL(uri).
			Go()).
			Attr("@click.middle",
				fmt.Sprintf(`(e) => e.view.window.open(%q, "_blank")`, uri))
	}
}
