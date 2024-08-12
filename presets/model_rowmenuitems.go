package presets

import (
	"fmt"
	"net/url"
	_ "unsafe"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	. "github.com/qor5/x/v3/ui/vuetify"
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

		msgr := MustGetMessages(ctx.R)
		if mi.mb.Info().Verifier().Do(PermUpdate).ObjectOn(obj).WithReq(ctx.R).IsAllowed() != nil {
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
	label := mb.label
	if !mb.singleton {
		label = mb.pluralLabel
	}

	return func(ctx *RecordMenuItemContext) h.HTMLComponent {
		mi := mb.modelInfo

		if mi.mb.Info().Verifier().Do(PermList).ObjectOn(ctx.Obj).WithReq(ctx.Ctx.R).IsAllowed() != nil {
			return nil
		}

		title := i18n.PT(ctx.Ctx.R, ModelsI18nModuleKey, mb.parent.label, humanizeString(label))
		uri := mi.ListingHref(append(ParentsModelID(ctx.Ctx.R), ID{Value: ctx.ID})...)
		return VListItem(
			web.Slot(
				VIcon(mi.mb.menuIcon),
			).Name("prepend"),
			VListItemTitle(h.Text(title)),
		).Attr("@click", web.Plaid().
			EventFunc(actions.OpenListingDialog).
			Query(ParamTargetPortal, ctx.TempPortal).
			URL(uri).
			Go()).
			Attr("@click.middle",
				fmt.Sprintf("(e) => e.view.window.open(%q, '_blank')", uri))
	}
}
