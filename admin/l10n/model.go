package l10n

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/utils/db_utils"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	. "github.com/theplant/htmlgo"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

const (
	FieldLocationLabel = "L10nLocationLabel"
	FieldLocations     = "L10nLocations"
)

func GetModelLocations[T LocaleInterface](b *Builder, m *presets.ModelBuilder, obj T) (_ []*T, err error) {
	db := b.db.Session(&gorm.Session{})
	mid := m.MustRecordID(obj)
	slice := m.NewModelSlice().(*[]*T)

	if err = db_utils.ModelIdWhere(db, nil, mid, "LocaleCode").
		Find(slice).Error; err != nil {
		return
	}

	return *slice, nil
}
func (b *Builder) GetModelLocationsI(m *presets.ModelBuilder, mid model.ID) (records []LocaleInterface, err error) {
	db := b.db.Session(&gorm.Session{})
	slice := m.NewModelSlice()

	if err = db_utils.ModelIdWhere(db, nil, mid, "LocaleCode").
		Find(slice).Error; err != nil {
		return
	}

	reflectutils.ForEach(reflect.ValueOf(slice).Elem().Interface(), func(item interface{}) {
		records = append(records, item.(LocaleInterface))
	})

	return records, nil
}

func (b *Builder) ModelInstall(pb *presets.Builder, m *presets.ModelBuilder) error {
	db := b.db
	obj := m.NewModel()
	_ = obj.(presets.SlugEncoder)
	_ = obj.(presets.SlugDecoder)
	_ = obj.(LocaleInterface)

	for _, fbs := range m.FieldBuilders() {
		fbs.Field("Locale")
		fbs.Field("LocaleCode").SetDisabled(true)
	}

	listing := m.Listing()
	listing.WrapSearchFunc(func(searcher presets.SearchFunc) presets.SearchFunc {
		return func(model interface{}, params *presets.SearchParams, ctx *web.EventContext) (r interface{}, totalCount int, err error) {
			if skip, _ := params.ContextValue(SkipLocaleCode).(bool); !skip {
				if localeCode := params.Query.Get("f_locale_code"); localeCode != "" {
					con := presets.SQLCondition{
						Query: "#TABLE#.locale_code = ?",
						Args:  []interface{}{localeCode},
					}
					params.SQLConditions = append(params.SQLConditions, &con)
				} else if localeCode = params.Query.Get("f_locale_code.in"); localeCode != "" {
					con := presets.SQLCondition{
						Query: "#TABLE#.locale_code IN ?",
						Args:  []any{strings.Split(localeCode, ",")},
					}
					params.SQLConditions = append(params.SQLConditions, &con)
				} else if localeCode = params.Query.Get("f_locale_code.notIn"); localeCode != "" {
					con := presets.SQLCondition{
						Query: "#TABLE#.locale_code NOT IN ?",
						Args:  []any{strings.Split(localeCode, ",")},
					}
					params.SQLConditions = append(params.SQLConditions, &con)
				} else if localeCode := web.GetContexValue(LocaleCode, params.Context, ctx.Context()); localeCode != nil {
					con := presets.SQLCondition{
						Query: "#TABLE#.locale_code = ?",
						Args:  []interface{}{localeCode},
					}
					params.SQLConditions = append(params.SQLConditions, &con)
				}
			}

			return searcher(model, params, ctx)
		}
	})

	setter := func(setter presets.SetterFunc) presets.SetterFunc {
		return func(obj interface{}, ctx *web.EventContext) {
			id := ctx.Param(presets.ParamID)
			if id == "" {
				if localeCode := ctx.R.Context().Value(LocaleCode); localeCode != nil {
					if err := reflectutils.Set(obj, "LocaleCode", localeCode); err != nil {
						return
					}
				}
			}
			if setter != nil {
				setter(obj, ctx)
			}
		}
	}

	ed := m.Editing().
		WrapPostSetterFunc(setter)

	if ed.HasCreatingBuilder() {
		ed.CreatingBuilder().WrapPostSetterFunc(setter)
	}

	m.Listing().WrapDeleteFunc(func(in presets.DeleteFunc) presets.DeleteFunc {
		return func(obj interface{}, id model.ID, cascade bool, ctx *web.EventContext) (err error) {
			if b.disableDeletionForDefaultInternationalizedRecord {
				var (
					countDB = db_utils.ModelIdWhere(
						db.Session(&gorm.Session{}).
							Where(id.Schema.
								FieldByName("LocaleCode").
								QuotedFullDBName()+" NOT IN (?)",
								b.GetDefaultLocaleCode()),
						m.NewModel(),
						id,
						"LocaleCode",
					)

					count int64
				)

				if err = countDB.Count(&count).Error; err != nil {
					return
				}

				if count > 0 {
					return errors.New(MustGetMessages(ctx.Context()).ErrDeleteInternationalizedRecord)
				}
			}

			if err = in(obj, id, cascade, ctx); err != nil {
				return
			}
			locale := id.GetValue("LocaleCode").(string)
			locale = fmt.Sprintf("%s(del:%d)", locale, time.Now().UnixMilli())

			var withoutKeys []string
			if ctx.R.URL.Query().Get("all_versions") == "true" {
				withoutKeys = append(withoutKeys, "Version")
			}

			if err = db_utils.ModelIdWhere(db.Unscoped(), obj, id, withoutKeys...).Update("locale_code", locale).Error; err != nil {
				return
			}
			return
		}
	})

	rmb := m.Listing().RowMenu()
	rmb.RowMenuItem("Localize").ComponentFunc(localizeRowMenuItemFunc(m.Info(), "", url.Values{}))

	locationLabelComp := func(field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
		loc := b.GetLocale(field.Obj.(LocaleInterface).EmbedLocale().LocaleCode)
		if loc == nil {
			return nil
		}
		return Text(loc.label)
	}

	l := m.Listing()
	l.Field(FieldLocationLabel).
		ComponentFunc(presets.ListingFieldComponentFuncWrapper(locationLabelComp)).
		SetI18nLabel(func(ctx context.Context) string {
			return MustGetMessages(ctx).Location
		})

	m.Detailing().Field(FieldLocations).
		ComponentFunc(presets.FieldComponentWrapper(func(field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			locations, _ := b.GetModelLocationsI(m, m.MustRecordID(field.Obj))
			if len(locations) == 0 {
				return nil
			}

			type entry struct {
				Label string
				Code  string
				Obj   LocaleInterface
			}

			entries := make([]*entry, len(locations))

			for i, e := range locations {
				var (
					code  = e.EmbedLocale().LocaleCode
					loc   = b.GetLocale(code)
					label string
				)

				if loc == nil {
					label = code
				} else {
					label = loc.label
				}
				entries[i] = &entry{
					Label: label,
					Code:  code,
					Obj:   e,
				}
			}

			slices.SortFunc(entries, func(a *entry, b *entry) int {
				if a.Label < b.Label {
					return -1
				}
				if a.Label > b.Label {
					return 1
				}
				return 0
			})

			current := field.Obj.(LocaleInterface).EmbedLocale().LocaleCode

			onclick := web.Plaid().URL(m.Info().ListingHrefCtx(ctx))

			if m.HasDetailing() {
				onclick.EventFunc(actions.Detailing)
			} else {
				if !m.CanEditObj(obj, ctx) {
					onclick = nil
				} else {
					onclick.EventFunc(actions.Edit)
				}
			}

			overlayMode := presets.GetOverlay(ctx)
			onclick.Query(presets.ParamOverlay, overlayMode.Up())

			portal := ctx.UID()
			var chips = HTMLComponents{web.Portal().Name(portal)}

			if overlayMode.Overlayed() {
				onclick.Query(presets.ParamTargetPortal, portal)
			}

			for _, e := range entries {
				c := v.VChip(Text(e.Label)).
					Density(v.DensityCompact)

				if e.Code == current {
					c.Color(v.ColorPrimary).Variant(v.VariantFlat)
				}

				if onclick != nil {
					id := e.Obj.(presets.SlugEncoder).PrimarySlug()
					onclick := onclick.Clone().Query(presets.ParamID, id)

					c.Attr("@click",
						onclick.Go()).
						Attr("@click.middle",
							fmt.Sprintf(`(e) => e.view.window.open(%q, "_blank")`, m.Info().DetailingHrefCtx(ctx, id)))
				}

				chips = append(chips, c)
			}
			return chips
		})).
		SetI18nLabel(func(ctx context.Context) string {
			return MustGetMessages(ctx).Localizations
		})

	m.Detailing().
		Action(FieldLocalizedEntries).
		ShowInList().
		Icon("mdi-translate").
		SetI18nLabel(func(ctx context.Context) string {
			return MustGetMessages(ctx).Localizations
		}).
		ComponentFunc(func(id string, ctx *web.EventContext) (HTMLComponent, error) {
			obj := m.NewModel()
			mid := m.MustParseRecordID(id)
			err := m.Fetcher(obj, mid, ctx)
			if err != nil {
				return nil, errors.New("Fetcher object failed: " + err.Error())
			}

			db := db.Session(&gorm.Session{})
			slice := m.NewModelSlice()

			if err = db_utils.ModelIdWhere(db, nil, mid, "LocaleCode").
				Where(mid.Schema.Table()+".locale_code != ?", mid.GetValue("LocaleCode")).
				Find(slice).Error; err != nil {
				return nil, errors.New("Find entries failed: " + err.Error())
			}

			type record struct {
				ID          string
				LocaleCode  string
				LocaleLabel string
				Title       string
			}

			var records []*record

			reflectutils.ForEach(reflect.ValueOf(slice).Elem().Interface(), func(item interface{}) {
				var localeCode = item.(LocaleInterface).EmbedLocale().LocaleCode
				records = append(records, &record{
					ID:          item.(presets.SlugEncoder).PrimarySlug(),
					Title:       m.RecordTitle(item, ctx),
					LocaleCode:  localeCode,
					LocaleLabel: b.GetLocaleLabel(localeCode),
				})
			})

			portalName := "_" + ctx.UID()
			indexUrl := m.Info().ListingHrefCtx(ctx)

			msgs := MustGetMessages(ctx.Context())

			return vue.UserComponent(
				web.Portal().Name(portalName),
				v.VDataTable(
					web.Slot(
						RawHTML(`{{ value }}`),
						v.VBtn("").
							Icon("mdi-eye").
							Attr("@click",
								web.Plaid().
									EventFunc(actions.Detailing).
									URL(indexUrl).
									Query(presets.ParamOverlay, actions.Dialog).
									Query(presets.ParamTargetPortal, portalName).
									Query(presets.ParamID, web.Var(`item.ID`)).Go(),
							).
							Attr("@click.middle", fmt.Sprintf(`(e) => e.view.window.open("%s/"+item.ID, "_blank")`, indexUrl)),
					).Name("item.actions").Scope("{ item, value }"),
				).Items(records).Headers([]any{
					map[string]any{
						"title": msgs.Location,
						"key":   "LocaleLabel",
					},
					map[string]any{
						"title": "",
						"key":   "Title",
					},
					map[string]any{
						"title": msgs.Actions,
						"key":   "actions",
					},
				}),
			), nil
		})

	registerEventFuncs(db, m, b)

	pb.FieldDefaults(presets.LIST).
		FieldType(Locale{}).
		ComponentFunc(localeListFunc(db, b))
	pb.FieldDefaults(presets.WRITE).
		FieldType(Locale{}).
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			value := b.localeValue(field, ctx)
			return Input("").Type("hidden").Attr(web.VField("LocaleCode", value)...)
		}).
		SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
			value := EmbedLocale(obj).LocaleCode
			if !slices.Contains(b.GetSupportLocaleCodesFromRequest(ctx.R), value) {
				return IncorrectLocaleErr
			}

			return nil
		})

	pb.AddWrapHandler(WrapHandlerKey, b.EnsureLocale)
	pb.AddMenuTopItemFunc(MenuTopItemFunc, runSwitchLocaleFunc(b))
	pb.I18n().
		RegisterForModule(language.English, I18nLocalizeKey, Messages_en_US).
		RegisterForModule(language.SimplifiedChinese, I18nLocalizeKey, Messages_zh_CN).
		RegisterForModule(language.Japanese, I18nLocalizeKey, Messages_ja_JP)
	return nil
}
