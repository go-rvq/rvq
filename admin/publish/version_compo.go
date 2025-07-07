package publish

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web/vue"
	"gorm.io/gorm"

	"github.com/go-rvq/rvq/admin/note"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/utils"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

type VersionComponentConfig struct {
	// If you want to use custom publish dialog, you can update the portal named PublishCustomDialogPortalName
	PublishEvent   func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) string
	UnPublishEvent func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) string
	RePublishEvent func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) string
	Top            bool
}

func DefaultVersionComponentFunc(b *presets.ModelBuilder, cfg ...VersionComponentConfig) presets.FieldComponentFunc {
	var config VersionComponentConfig
	if len(cfg) > 0 {
		config = cfg[0]
	}
	return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		var (
			obj            = field.Obj
			version        VersionInterface
			status         StatusInterface
			primarySlugger presets.SlugEncoder
			ok             bool
			versionSwitch  *v.VChipBuilder
			publishBtn     h.HTMLComponent
		)
		msgr := GetMessages(ctx.Context())
		utilsMsgr := i18n.MustGetModuleMessages(ctx.Context(), utils.I18nUtilsKey, utils.Messages_en_US).(*utils.Messages)

		primarySlugger, ok = obj.(presets.SlugEncoder)
		if !ok {
			panic("obj should be SlugEncoder")
		}

		var (
			tempPortal = field.Name + "Portal"
			div        = h.Div(web.Portal().Name(tempPortal)).Class("w-100 d-inline-flex")
		)

		if !config.Top {
			div.Class("pb-4")
		}

		if version, ok = obj.(VersionInterface); ok {
			versionSwitch = v.VChip(
				h.Text(version.EmbedVersion().VersionName),
			).Label(true).Variant(v.VariantOutlined).
				Attr("style", "height:40px;").
				On("click", web.Plaid().EventFunc(actions.OpenListingDialog).
					URL(b.Info().PresetsPrefix()+"/"+field.ModelInfo.URI()+"-version-list-dialog").
					Query("select_id", primarySlugger.PrimarySlug()).
					Query(presets.ParamTargetPortal, tempPortal).
					BeforeScript(fmt.Sprintf("%s ||= ''", VarCurrentDisplayID)).
					ThenScript(fmt.Sprintf("%s = %q", VarCurrentDisplayID, primarySlugger.PrimarySlug())).
					Go()).
				Class(v.W100)
			if status, ok = obj.(StatusInterface); ok {
				versionSwitch.AppendChild(LiveChipsFormBuilder.Status(status.EmbedStatus().Status, msgr).Class("mx-2"))
			}
			versionSwitch.AppendChild(v.VSpacer())
			versionSwitch.AppendIcon("mdi-chevron-down")

			div.AppendChildren(versionSwitch)
			div.AppendChildren(v.VBtn(msgr.Duplicate).PrependIcon("mdi-file-document-multiple").
				Height(40).Class("ml-2").Variant(v.VariantOutlined).
				Attr("@click", fmt.Sprintf(`locals.action="%s";locals.commonConfirmDialog = true`, EventDuplicateVersion)))
		}

		if status, ok = obj.(StatusInterface); ok {
			gotoPublishedURL := h.A(v.VIcon("mdi-open-in-new")).
				Attr("v-if", `locals.`+FieldOnlineUrl).
				Attr(":href", `locals.`+FieldOnlineUrl).
				Attr("target", `_blank`).
				Class("v-btn v-btn--icon v-theme--light v-btn--density-comfortable v-btn--size-default v-btn--variant-flat")

			switch status.EmbedStatus().Status {
			case StatusDraft, StatusOffline:
				publishEvent := fmt.Sprintf(`locals.action="%s";locals.commonConfirmDialog = true`, EventPublish)
				if config.PublishEvent != nil {
					publishEvent = config.PublishEvent(obj, field, ctx)
				}
				publishBtn = h.Div(
					v.VBtn(msgr.Publish).Attr("@click", publishEvent).Rounded("0").
						Class("rounded-s ml-2").Variant(v.VariantFlat).Color(v.ColorPrimary).Height(40),
					gotoPublishedURL,
				)
			case StatusOnline:
				unPublishEvent := fmt.Sprintf(`locals.action="%s";locals.commonConfirmDialog = true`, EventUnpublish)
				if config.UnPublishEvent != nil {
					unPublishEvent = config.UnPublishEvent(obj, field, ctx)
				}
				rePublishEvent := fmt.Sprintf(`locals.action="%s";locals.commonConfirmDialog = true`, EventRepublish)
				if config.RePublishEvent != nil {
					rePublishEvent = config.RePublishEvent(obj, field, ctx)
				}
				publishBtn = h.Div(
					v.VBtn(msgr.Unpublish).Attr("@click", unPublishEvent).
						Class("ml-2").Variant(v.VariantFlat).Color(v.ColorError).Height(40),
					v.VBtn(msgr.Republish).Attr("@click", rePublishEvent).
						Class("ml-2").Variant(v.VariantFlat).Color(v.ColorPrimary).Height(40),
					gotoPublishedURL,
				).Class("d-inline-flex")
			}
			div.AppendChildren(publishBtn)

			// Publish/UnPublish/Republish ConfirmDialog
			div.AppendChildren(
				utils.ConfirmDialog(msgr.Areyousure, web.Plaid().EventFunc(web.Var("locals.action")).
					Query(presets.ParamID, primarySlugger.PrimarySlug()).Go(),
					utilsMsgr),
			)

			// Publish/UnPublish/Republish CustomDialog
			if config.UnPublishEvent != nil || config.RePublishEvent != nil || config.PublishEvent != nil {
				div.AppendChildren(web.Portal().Name(PortalPublishCustomDialog))
			}
		}

		if _, ok = obj.(ScheduleInterface); ok {
			var scheduleBtn h.HTMLComponent
			clickEvent := web.POST().
				EventFunc(eventSchedulePublishDialog).
				Query(presets.ParamOverlay, actions.Dialog).
				Query(presets.ParamID, primarySlugger.PrimarySlug()).
				URL(fmt.Sprintf("%s/%s", b.Info().PresetsPrefix(), b.Info().URI())).Go()
			if config.Top {
				scheduleBtn = v.VAutocomplete().PrependInnerIcon("mdi-alarm").Density(v.DensityCompact).
					Variant(v.FieldVariantSoloFilled).ModelValue("Schedule Publish Time").
					BgColor(v.ColorPrimaryLighten2).Readonly(true).
					Width(600).HideDetails(true).Attr("@click", clickEvent).Class("ml-2 text-caption")
			} else {
				scheduleBtn = v.VBtn("").Children(v.VIcon("mdi-alarm").Size(v.SizeXLarge)).Rounded("0").Class("ml-1 rounded-e").
					Variant(v.VariantFlat).Color(v.ColorPrimary).Height(40).Attr("@click", clickEvent)
			}
			div.AppendChildren(scheduleBtn)
			// SchedulePublishDialog
			div.AppendChildren(web.Portal().Name(PortalSchedulePublishDialog))
		}

		var onlineUrl string
		if status != nil {
			onlineUrl = status.EmbedStatus().OnlineUrl
		}
		if pu, _ := obj.(PublicUrlInterface); pu != nil {
			onlineUrl = pu.GetPublicUrl(b, ctx)
		}
		return vue.UserComponent(div).Scope("locals", vue.Var(`{action: "", commonConfirmDialog: false, `+FieldOnlineUrl+`: `+strconv.Quote(onlineUrl)+` }`))
	}
}

func DefaultVersionBar(db *gorm.DB) presets.ObjectComponentFunc {
	return func(obj interface{}, ctx *web.EventContext) h.HTMLComponent {
		msgr := GetMessages(ctx.Context())
		res := h.Div().Class("d-inline-flex align-center")

		slugEncoderIf := obj.(presets.SlugEncoder)
		slugDncoderIf := obj.(presets.SlugDecoder)
		mp := slugDncoderIf.PrimaryColumnValuesBySlug(slugEncoderIf.PrimarySlug())

		currentObj := reflect.New(reflect.TypeOf(obj).Elem()).Interface()
		err := db.Where("id = ?", mp["id"]).Where("status = ?", StatusOnline).First(&currentObj).Error
		if err != nil {
			return res
		}
		versionIf := currentObj.(VersionInterface)
		currentVersionStr := fmt.Sprintf("%s: %s", msgr.OnlineVersion, versionIf.EmbedVersion().VersionName)
		res.AppendChildren(v.VChip(h.Span(currentVersionStr)).Density(v.DensityCompact).Color(v.ColorSuccess))

		if _, ok := currentObj.(ScheduleInterface); !ok {
			return res
		}

		nextObj := reflect.New(reflect.TypeOf(obj).Elem()).Interface()
		flagTime := db.NowFunc()
		count := int64(0)
		err = db.Model(nextObj).Where("id = ?", mp["id"]).Where("scheduled_start_at >= ?", flagTime).Count(&count).Error
		if err != nil {
			return res
		}

		if count == 0 {
			return res
		}

		err = db.Where("id = ?", mp["id"]).Where("scheduled_start_at >= ?", flagTime).Order("scheduled_start_at ASC").First(&nextObj).Error
		if err != nil {
			return res
		}
		res.AppendChildren(
			h.Div(
				h.Div().Class(fmt.Sprintf(`w-100 bg-%s`, v.ColorSuccessLighten2)).Style("height:4px"),
				v.VIcon("mdi-circle").Size(v.SizeXSmall).Color(v.ColorSuccess).Attr("style", "position:absolute;left:0;right:0;margin-left:auto;margin-right:auto"),
			).Class("h-100 d-flex align-center").Style("position:relative;width:40px"),
		)
		versionIf = nextObj.(VersionInterface)
		// TODO use nextVersion I18n
		nextText := fmt.Sprintf("%s: %s", msgr.OnlineVersion, versionIf.EmbedVersion().VersionName)
		res.AppendChildren(v.VChip(h.Span(nextText)).Density(v.DensityCompact).Color(v.ColorSecondary))
		if count >= 2 {
			res.AppendChildren(
				h.Div(
					h.Div().Class(fmt.Sprintf(`w-100 bg-%s`, v.ColorSecondaryLighten1)).Style("height:4px"),
				).Class("h-100 d-flex align-center").Style("width:40px"),
				h.Div(
					h.Text(fmt.Sprintf(`+%v`, count)),
				).Class(fmt.Sprintf(`text-caption bg-%s`, v.ColorSecondaryLighten1)),
			)
		}
		return res
	}
}

func configureVersionListDialog(db *gorm.DB, b *presets.Builder, pm *presets.ModelBuilder) {
	// actually, VersionListDialog is a listing
	// use this URL : URLName-version-list-dialog
	mb := b.Model(pm.NewModel()).
		URIName(pm.Info().URI() + "-version-list-dialog").
		InMenu(false)

	registerEventFuncsForVersion(mb, pm, db)

	// TODO: i18n
	lb := mb.Listing("Version", "State", "StartAt", "EndAt", "Notes", "Option").
		DialogWidth("900").
		Title("Version List").
		SearchColumns("version", "version_name").
		PerPage(10).
		WrapSearchFunc(func(in presets.SearchFunc) presets.SearchFunc {
			return func(model interface{}, params *presets.SearchParams, ctx *web.EventContext) (r interface{}, totalCount int, err error) {
				id := ctx.R.FormValue("select_id")
				if id == "" {
					id = ctx.R.FormValue("f_select_id")
				}
				if id != "" {
					cs := mb.NewModel().(presets.SlugDecoder).PrimaryColumnValuesBySlug(id)
					con := presets.SQLCondition{
						Query: "id = ?",
						Args:  []interface{}{cs["id"]},
					}
					params.SQLConditions = append(params.SQLConditions, &con)
				}
				params.OrderBy = "created_at DESC"

				return in(model, params, ctx)
			}
		})
	lb.CellWrapperFunc(func(cell h.MutableAttrHTMLComponent, id string, obj interface{}, dataTableID string, ctx *web.EventContext) h.HTMLComponent {
		return cell
	})
	lb.Field("Version").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		obj := field.Obj
		versionName := obj.(VersionInterface).EmbedVersion().VersionName
		p := obj.(presets.SlugEncoder)
		id := ctx.R.FormValue("select_id")
		if id == "" {
			id = ctx.R.FormValue("f_select_id")
		}

		queries := ctx.Queries()
		queries.Set("select_id", p.PrimarySlug())
		onChange := web.Plaid().
			URL(ctx.R.URL.Path).
			Queries(queries).
			EventFunc(actions.UpdateListingDialog).
			Query(presets.ParamPortalID, ctx.R.FormValue(presets.ParamPortalID)).
			Go()

		return h.Td().Children(
			h.Div().Class("d-inline-flex align-center").Children(
				v.VRadio().ModelValue(p.PrimarySlug()).TrueValue(id).Attr("@change", onChange),
				h.Text(versionName),
			),
		)
	})
	lb.Field("State").ComponentFunc(StatusListFunc(&LiveChipsListBuilder))
	lb.Field("StartAt").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		p := field.Obj.(ScheduleInterface)

		return h.Td(
			h.Text(ScheduleTimeString(p.EmbedSchedule().ScheduledStartAt)),
		)
	}).Label("Start at")
	lb.Field("EndAt").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		p := field.Obj.(ScheduleInterface)
		return h.Td(
			h.Text(ScheduleTimeString(p.EmbedSchedule().ScheduledEndAt)),
		)
	}).Label("End at")

	lb.Field("Notes").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		p := field.Obj.(presets.SlugEncoder)
		rt := pm.Info().Label()
		ri := p.PrimarySlug()
		userID, _ := note.GetUserData(ctx)
		count := note.GetUnreadNotesCount(db, userID, rt, ri)

		return h.Td(
			h.If(count > 0,
				v.VBadge().Content(count).Color("red"),
			).Else(
				h.Text(""),
			),
		)
	}).Label("Unread Notes")

	lb.Field("Option").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		msgr := GetMessages(ctx.Context())
		pmsgr := presets.MustGetMessages(ctx.Context())

		obj := field.Obj
		id := obj.(presets.SlugEncoder).PrimarySlug()
		versionName := obj.(VersionInterface).EmbedVersion().VersionName
		status := obj.(StatusInterface).EmbedStatus().Status
		disable := status == StatusOnline || status == StatusOffline
		deniedUpdate := mb.EditingDisabled() || mb.Permissioner().ReqObjectUpdater(ctx.R, obj).Denied()
		deniedDelete := mb.DeletingDisabled() || mb.Permissioner().ReqObjectDeleter(ctx.R, obj).Denied()
		return h.Td().Children(
			v.VBtn(msgr.Rename).Disabled(disable || deniedUpdate).PrependIcon("mdi-rename-box").Size(v.SizeXSmall).Color(v.ColorPrimary).Variant(v.VariantText).
				On("click", web.Plaid().
					URL(ctx.R.URL.Path).
					EventFunc(eventRenameVersionDialog).
					Query(presets.ParamListingQueries, ctx.Queries().Encode()).
					Query(presets.ParamOverlay, actions.Dialog).
					Query(presets.ParamID, id).
					Query("version_name", versionName).
					Go(),
				),
			v.VBtn(pmsgr.Delete).Disabled(disable || deniedDelete).PrependIcon("mdi-delete").Size(v.SizeXSmall).Color(v.ColorPrimary).Variant(v.VariantText).
				On("click", web.Plaid().
					URL(ctx.R.URL.Path).
					EventFunc(eventDeleteVersionDialog).
					Query(presets.ParamListingQueries, ctx.Queries().Encode()).
					Query(presets.ParamOverlay, actions.Dialog).
					Query(presets.ParamID, id).
					Query("version_name", versionName).
					Query(paramCurrentDisplaySlug, web.Var(VarCurrentDisplayID)).
					Go(),
				),
		)
	})
	lb.NewButtonFunc(func(ctx *web.EventContext) h.HTMLComponent { return nil })
	lb.FooterAction("Cancel").ButtonCompFunc(func(ctx *web.EventContext) h.HTMLComponent {
		return v.VBtn("Cancel").Variant(v.VariantElevated).Attr("@click", "vars.presetsListingDialog=false")
	})
	lb.FooterAction("Save").ButtonCompFunc(func(ctx *web.EventContext) h.HTMLComponent {
		id := ctx.R.FormValue("select_id")
		if id == "" {
			id = ctx.R.FormValue("f_select_id")
		}

		return v.VBtn("Save").Disabled(id == "").Variant(v.VariantElevated).Color(v.ColorSecondary).Attr("@click", web.Plaid().
			Query("select_id", id).
			URL(pm.Info().PresetsPrefix()+"/"+pm.Info().URI()).
			EventFunc(eventSelectVersion).
			Go())
	})
	lb.RowMenu().Empty()

	lb.FilterDataFunc(func(ctx *web.EventContext) vx.FilterData {
		return []*vx.FilterItem{
			{
				Key:          "all",
				Invisible:    true,
				SQLCondition: ``,
			},
			{
				Key:          "online_versions",
				Invisible:    true,
				SQLCondition: `status = 'online'`,
			},
			{
				Key:          "named_versions",
				Invisible:    true,
				SQLCondition: `version <> version_name`,
			},
		}
	})

	lb.FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		msgr := GetMessages(ctx.Context())
		id := ctx.R.FormValue("select_id")
		if id == "" {
			id = ctx.R.FormValue("f_select_id")
		}
		return []*presets.FilterTab{
			{
				Label: msgr.FilterTabAllVersions,
				ID:    "all",
				Query: url.Values{"all": []string{"1"}, "select_id": []string{id}},
			},
			{
				Label: msgr.FilterTabOnlineVersion,
				ID:    "online_versions",
				Query: url.Values{"online_versions": []string{"1"}, "select_id": []string{id}},
			},
			{
				Label: msgr.FilterTabNamedVersions,
				ID:    "named_versions",
				Query: url.Values{"named_versions": []string{"1"}, "select_id": []string{id}},
			},
		}
	})
}
