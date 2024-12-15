package l10n

import (
	"context"
	"reflect"
	"slices"

	"github.com/qor5/admin/v3/activity"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/admin/v3/utils/db_utils"
	"github.com/qor5/web/v3"
	"github.com/qor5/web/v3/vue"
	. "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"
)

const (
	Localize   = "l10n_LocalizeEvent"
	DoLocalize = "l10n_DoLocalizeEvent"

	FromID      = "l10n_DoLocalize_FromID"
	FromVersion = "l10n_DoLocalize_FromVersion"
	FromLocale  = "l10n_DoLocalize_FromLocale"

	LocalizeFrom = "Localize From"
	LocalizeTo   = "Localize To"
)

func registerEventFuncs(db *gorm.DB, mb *presets.ModelBuilder, lb *Builder, ab *activity.Builder) {
	mb.RegisterEventHandler(Localize, localizeToConfirmation(db, lb, mb))
	mb.RegisterEventHandler(DoLocalize, doLocalizeTo(db, mb, lb, ab))
}

type SelectLocale struct {
	Label string
	Code  string
}

func localizeToConfirmation(db *gorm.DB, lb *Builder, mb *presets.ModelBuilder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		var (
			presetsMsgr = presets.MustGetMessages(ctx.Context())
			msgr        = MustGetMessages(ctx.Context())
			paramID     = ctx.Param(presets.ParamID)
			mid         = mb.MustParseRecordID(paramID)
			id          = mid.GetValue("ID")
			fromLocale  = mid.GetValue("LocaleCode").(string)
			obj         = mb.NewModelSlice()
		)

		if err = db.Session(&gorm.Session{}).
			Distinct("locale_code").
			Where("id = ? AND locale_code <> ?", id, fromLocale).
			Find(obj).Error; err != nil {
			return
		}

		if err != nil {
			return
		}
		var (
			vo           = reflect.ValueOf(obj).Elem()
			existLocales []string
			chips        []h.HTMLComponent
		)

		for i := 0; i < vo.Len(); i++ {
			code := vo.Index(i).Elem().FieldByName("LocaleCode").String()
			existLocales = append(existLocales, code)
			chips = append(chips, VChip(h.Text(lb.GetLocale(code).Label())))
		}

		var (
			toLocales     = lb.GetSupportLocaleCodesFromRequest(ctx.R)
			selectLocales []SelectLocale
		)

		for _, locale := range toLocales {
			if locale == fromLocale {
				continue
			}
			if !slices.Contains(existLocales, locale) || vo.Len() == 0 {
				selectLocales = append(selectLocales, SelectLocale{Label: MustGetTranslation(ctx.Context(), lb.GetLocaleLabel(locale)), Code: locale})
			}
		}

		attr := web.VField("LocalizeTo", vue.Var("[]"))

		var cb = &presets.ContentComponentBuilder{
			Context: ctx,
			Title:   msgr.Localize,
			Body: h.HTMLComponents{
				VContainer(
					VRow(
						VCol(
							vx.VXReadonlyField().
								Label(msgr.LocalizeFrom).
								Value(lb.GetLocaleLabel(fromLocale)),
						).Class("px-0 py-0"),
					),
					h.If(len(existLocales) > 0,
						VRow(
							VCol(
								h.Label(msgr.CurrentLocalizations).Class("v-label theme--light text-caption"),
								VChipGroup(chips...),
							).Class("px-0 py-0"),
						),
					),
					VRow(
						VCol(
							VSelect().
								Attr(attr...).
								Variant(FieldVariantUnderlined).
								Label(msgr.LocalizeTo).
								Multiple(true).
								Chips(true).
								Items(selectLocales).
								ItemTitle("Label").
								HideDetails(true).
								ItemValue("Code"),
						).Class("px-0 py-0"),
					)),
			},
			BottomActions: h.HTMLComponents{
				VSpacer(),
				VBtn(presetsMsgr.OK).
					Color("primary").
					Variant(VariantFlat).
					Theme(ThemeDark).
					Attr("@click", web.Plaid().
						EventFunc(DoLocalize).
						Query(presets.ParamID, paramID).
						Query("localize_from", fromLocale).
						URL(ctx.R.URL.Path).
						Go()),
			},

			Overlay: &presets.ContentComponentBuilderOverlay{
				Mode: actions.Dialog,
			},
		}

		mb.Builder().Dialog().
			SetScrollable(true).
			SetTargetPortal(actions.Dialog.PortalName()).
			Respond(ctx, &r, cb.BuildOverlay())

		return
	}
}

type doLocalize struct {
	LocalizeTo []string
}

func doLocalizeTo(db *gorm.DB, mb *presets.ModelBuilder, lb *Builder, ab *activity.Builder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {

		var (
			mid         = mb.MustParseRecordID(ctx.Param(presets.ParamID))
			fromID      = mid.GetValue("ID")
			fromVersion = mid.GetValue("Version")
			fromLocale  = mid.GetValue("LocaleCode")
			toForm      doLocalize
			to          = make(map[string]interface{})
		)

		ctx.UnmarshalForm(&toForm)

		for _, v := range toForm.LocalizeTo {
			for _, lc := range lb.GetSupportLocaleCodes() {
				if v == lc {
					to[v] = struct{}{}
					break
				}
			}
		}

		if len(to) == 0 {
			web.AppendRunScripts(&r, "vars.localizeConfirmation = false")
			return
		}

		fromObj := mb.NewModel()

		if err = db_utils.ModelIdWhere(db, mb.NewModel(), mid).First(fromObj).Error; err != nil {
			return
		}

		var toObjs []interface{}
		defer func(fromObj interface{}) {
			if ab == nil {
				return
			}
			if _, ok := ab.GetModelBuilder(fromObj); !ok {
				return
			}
			if len(toObjs) > 0 {
				if err = ab.AddCustomizedRecord(LocalizeFrom, false, ctx.R.Context(), fromObj); err != nil {
					return
				}
				for _, toObj := range toObjs {
					if err = ab.AddCustomizedRecord(LocalizeTo, false, ctx.R.Context(), toObj); err != nil {
						return
					}
				}
			}
		}(reflect.Indirect(reflect.ValueOf(fromObj)).Interface())
		me := mb.Editing()

		for toLocale := range to {
			toObj := mb.NewModel()
			mid.SetTo(toObj)

			if err = reflectutils.Set(toObj, "LocaleCode", toLocale); err != nil {
				return
			}

			me.SetObjectFields(fromObj, toObj, &presets.FieldContext{
				Obj:       fromObj,
				ModelInfo: mb.Info(),
			}, false, presets.ContextModifiedIndexesBuilder(ctx).FromHidden(ctx.R), ctx)

			if vErr := me.Validators.Validate(toObj, presets.FieldModeStack{presets.EDIT}, ctx); vErr.HaveErrors() {
				presets.ShowMessage(&r, vErr.Error(), "error")
				return
			}

			newContext := context.WithValue(ctx.R.Context(), FromID, fromID)
			newContext = context.WithValue(newContext, FromVersion, fromVersion)
			newContext = context.WithValue(newContext, FromLocale, fromLocale)
			ctx.R = ctx.R.WithContext(newContext)

			var done func() error

			if cb, _ := mb.GetData(LocalizeOptions).(*ModelLocalizeOptions); cb != nil {
				if done, err = cb.LocalizeCallback(ctx, fromObj, toObj); err != nil {
					return
				}
			}

			if err = me.CreatingBuilder().Creator(toObj, ctx); err != nil {
				return
			}

			if done != nil {
				if err = done(); err != nil {
					return
				}
			}

			toObjs = append(toObjs, toObj)
		}

		presets.ShowMessage(&r, MustGetTranslation(ctx.Context(), "SuccessfullyLocalized"), "")

		// refresh current page
		r.Reload = true
		return
	}
}
