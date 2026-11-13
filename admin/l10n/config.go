package l10n

import (
	"fmt"
	"net/url"
	"reflect"
	"slices"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/utils/bcp47"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	"gorm.io/gorm"
)

const (
	WrapHandlerKey  = "l10nWrapHandlerKey"
	MenuTopItemFunc = "l10nMenuTopItemFunc"
	SlugLocaleCode  = "locale_code"
	BuilderKey      = "l10nBuilderKey"
)

// localeListFunc is the cell of the languages of a record in a listing: the
// flag of each (its code, of a language of no region), its name in a
// tooltip, the record's own first and marked; a click opens the detail of
// the record in that language, the middle button in a new tab.
func localeListFunc(db *gorm.DB, lb *Builder) presets.FieldComponentFunc {
	return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		obj := field.Obj
		id, err := reflectutils.Get(obj, "ID")
		if err != nil {
			return nil
		}
		current := EmbedLocale(obj).LocaleCode

		objs := reflect.New(reflect.SliceOf(reflect.TypeOf(obj).Elem())).Interface()
		err = db.Distinct("locale_code").Where("id = ? AND locale_code <> ?", id, current).Find(objs).Error
		if err != nil {
			return nil
		}
		vo := reflect.ValueOf(objs).Elem()
		var existLocales []string
		for i := 0; i < vo.Len(); i++ {
			existLocales = append(existLocales, vo.Index(i).FieldByName("LocaleCode").String())
		}

		var others []string
		for _, locale := range lb.GetSupportLocaleCodesFromRequest(ctx.R) {
			if slices.Contains(existLocales, locale) {
				others = append(others, locale)
			}
		}

		flags := []h.HTMLComponent{lb.flagLink(field, ctx, current, true)}
		for _, locale := range others {
			flags = append(flags, lb.flagLink(field, ctx, locale, false))
		}
		return h.Td(h.Div(flags...).Class("d-flex align-center ga-1"))
	}
}

// flagLink is the flag of the record of field in the language code (else its
// code), its name in a tooltip; current, marked. A click opens its detail —
// in the listing's overlay —, the middle button in a new tab.
func (lb *Builder) flagLink(field *presets.FieldContext, ctx *web.EventContext, code string, current bool) h.HTMLComponent {
	text := bcp47.Flag(code)
	if text == "" {
		text = code
	}
	flag := h.Span(text).Attr("v-bind", "props").Attr("data-locale", code).
		Class("px-1 rounded text-body-1").Style("cursor:pointer;line-height:1.6")
	if current {
		flag.Class("border border-primary border-opacity-100")
	} else {
		flag.Style("cursor:pointer;line-height:1.6;opacity:.55")
	}
	if mi := field.ModelInfo; mi != nil && mi.Builder() != nil {
		mb := mi.Builder()
		if mid, err := mb.RecordID(field.Obj); err == nil {
			o := mb.NewModel()
			mid.Clone().WithField("LocaleCode", code).SetTo(o)
			if enc, ok := o.(presets.SlugEncoder); ok {
				slug := enc.PrimarySlug()
				open := presets.GetItemFormHosts(ctx).OpenDetailExpr(slug)
				if open == "" {
					open = web.Plaid().EventFunc(actions.Detailing).URL(mi.ListingHrefCtx(ctx)).
						Query(presets.ParamID, slug).Query(presets.ParamOverlay, actions.Dialog).Go()
				}
				flag.Attr("@click.stop", open).
					Attr("@click.middle.stop", fmt.Sprintf(`(e) => e.view.window.open(%q, "_blank")`, mi.DetailingHrefCtx(ctx, slug)))
			}
		}
	}
	return VTooltip(
		web.Slot(flag).Name("activator").Scope("{ props }"),
	).Text(MustGetTranslation(ctx.Context(), lb.GetLocaleLabel(code))).Location(LocationTop)
}

func runSwitchLocaleFunc(lb *Builder) func(ctx *web.EventContext) (r h.HTMLComponent) {
	return func(ctx *web.EventContext) (r h.HTMLComponent) {
		supportLocales := lb.GetSupportLocaleCodesFromRequest(ctx.R)

		if len(lb.GetSupportLocaleCodes()) <= 1 || len(supportLocales) == 0 {
			return nil
		}

		localeQueryName := lb.queryName

		if len(supportLocales) == 1 {
			return h.Template().Children(
				h.Div(
					VList(
						VListItem(
							web.Slot(
								// icon was language
								VIcon("mdi-translate").Size(SizeSmall).
									Class("ml-1").
									Attr("style", "margin-right: 16px"),
							).Name("prepend"),
							VListItemTitle(
								h.Div(h.Text(fmt.Sprintf("%s%s %s",
									MustGetTranslation(ctx.Context(), "Location"),
									MustGetTranslation(ctx.Context(), "Colon"),
									MustGetTranslation(ctx.Context(),
										lb.GetLocaleLabel(supportLocales[0]))))).Role("button"),
							),
						).Class("pa-0").Density(DensityCompact),
					).Class("pa-0 ma-n4 mt-n6"),
				).Attr("@click", web.Plaid().Query(localeQueryName, supportLocales[0]).Go()),
			)
		}

		locale := lb.GetCorrectLocaleCode(ctx.R)

		var locales []h.HTMLComponent
		for _, contry := range supportLocales {
			locales = append(locales,
				h.Div(
					VListItem(
						VListItemTitle(
							h.Div(h.Text(MustGetTranslation(ctx.Context(), lb.GetLocaleLabel(contry)))),
						),
					).Attr("@click", web.Plaid().Query(localeQueryName, contry).Go()),
				),
			)
		}

		return VMenu().Children(
			h.Template().Attr("v-slot:activator", "{props}").Children(
				h.Div(
					VList(
						VListItem(
							web.Slot(
								// icon was language
								VIcon("mdi-translate").Size(SizeSmall).Class("ml-1"),
							).Name("prepend"),
							VListItemTitle(
								h.Text(fmt.Sprintf("%s%s %s",
									MustGetTranslation(ctx.Context(), "Location"),
									MustGetTranslation(ctx.Context(), "Colon"),
									MustGetTranslation(ctx.Context(), lb.GetLocaleLabel(locale)))),
							),
							web.Slot(
								// icon was arrow_drop_down
								VIcon("mdi-arrow-down").Class("mr-1"),
							).Name("append"),
						).Class("pa-0").Density(DensityCompact),
					).Class("pa-0 ma-n4 mt-n6"),
				).Attr("v-bind", "props"),
			),

			VList(
				locales...,
			).Density(DensityCompact),
		)
	}
}

func localizeRowMenuItemFunc(mi *presets.ModelInfo, url string, editExtraParams url.Values) presets.RecordMenuItemFunc {
	return func(rctx *presets.RecordMenuItemContext) h.HTMLComponent {
		var (
			ctx = rctx.Ctx
			obj = rctx.Obj
			id  = rctx.ID
		)
		if mi.Permissioner().ReqObjectUpdater(ctx.R, obj).Denied() {
			return nil
		}

		return VListItem(
			web.Slot(
				// icon was language
				VIcon("mdi-translate"),
			).Name("prepend"),
			VListItemTitle(h.Text(MustGetTranslation(ctx.Context(), "Localize"))),
		).Attr("@click", web.Plaid().
			EventFunc(Localize).
			Queries(editExtraParams).
			Query(presets.ParamID, id).
			URL(url).
			Go())
	}
}
