package publish

import (
	"cmp"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/admin/v3/utils"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	"github.com/qor5/x/v3/perm"
	v "github.com/qor5/x/v3/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"
)

const (
	PortalSchedulePublishDialog = "publish_PortalSchedulePublishDialog"
	PortalPublishCustomDialog   = "publish_PortalPublishCustomDialog"

	VarCurrentDisplayID = "vars.publish_VarCurrentDisplayID"
)

func duplicateVersionAction(db *gorm.DB, mb *presets.ModelBuilder, _ *Builder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		slug := ctx.Param(presets.ParamID)
		obj := mb.NewModel()
		if err = utils.PrimarySluggerWhere(db, mb.NewModel(), slug).First(obj).Error; err != nil {
			return
		}

		ver, ok := obj.(VersionInterface)
		if !ok {
			err = errInvalidObject
			return
		}

		oldVersion := ver.EmbedVersion().Version
		newVersion, err := ver.CreateVersion(db, slug, mb.NewModel())
		if err != nil {
			return
		}

		if err = reflectutils.Set(obj, "Version", Version{
			Version:       newVersion, // In fact, it is also set in CreateVersion, just in case
			VersionName:   newVersion, // In fact, it is also set in CreateVersion, just in case
			ParentVersion: oldVersion,
		}); err != nil {
			return
		}

		if _, ok := ver.(StatusInterface); ok {
			st := ver.(StatusInterface).EmbedStatus()
			st.Status = StatusDraft
			st.OnlineUrl = ""

			// _, err = reflectutils.Get(obj, "Status")
			// if err == nil {
			// 	if err = reflectutils.Set(obj, "Status", Status{Status: StatusDraft}); err != nil {
			// 		return
			// 	}
			// }
		}
		if _, ok := ver.(ScheduleInterface); ok {
			sched := ver.(ScheduleInterface).EmbedSchedule()

			sched.ScheduledStartAt = nil
			sched.ScheduledEndAt = nil
			sched.ActualStartAt = nil
			sched.ActualEndAt = nil
			// _, err = reflectutils.Get(obj, "Schedule")
			// if err == nil {
			// 	if err = reflectutils.Set(obj, "Schedule", Schedule{}); err != nil {
			// 		return
			// 	}
			// }
		}

		_, err = reflectutils.Get(obj, "CreatedAt")
		if err == nil {
			if err = reflectutils.Set(obj, "CreatedAt", time.Time{}); err != nil {
				return
			}
		}
		_, err = reflectutils.Get(obj, "UpdatedAt")
		if err == nil {
			if err = reflectutils.Set(obj, "UpdatedAt", time.Time{}); err != nil {
				return
			}
		}
		err = nil

		slug = obj.(presets.SlugEncoder).PrimarySlug()
		if err = mb.Editing().Creating().Saver(obj, slug, ctx); err != nil {
			presets.ShowMessage(&r, err.Error(), "error")
			return
		}

		if !mb.HasDetailing() {
			// close dialog and open editing
			web.AppendRunScripts(&r,
				presets.CloseListingDialogVarScript,
				web.Plaid().EventFunc(actions.Edit).Query(presets.ParamID, slug).Go(),
			)
			return
		}
		// close dialog and open detailingDrawer
		web.AppendRunScripts(&r,
			presets.CloseListingDialogVarScript,
			presets.CloseRightDrawerVarScript,
			web.Plaid().EventFunc(actions.Detailing).Query(presets.ParamID, slug).Go(),
		)

		msgr := i18n.MustGetModuleMessages(ctx.R, I18nPublishKey, Messages_en_US).(*Messages)
		presets.ShowMessage(&r, msgr.SuccessfullyCreated, "")

		r.RunScript = web.Plaid().ThenScript(r.RunScript).Go()
		return
	}
}

func selectVersion(pm *presets.ModelBuilder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id := ctx.R.FormValue("select_id")

		if !pm.HasDetailing() {
			// close dialog and open editing
			web.AppendRunScripts(&r,
				presets.CloseListingDialogVarScript,
				web.Plaid().EventFunc(actions.Edit).Query(presets.ParamID, id).Go(),
			)
			return
		}
		// close dialog and open detailingDrawer
		web.AppendRunScripts(&r,
			presets.CloseListingDialogVarScript,
			fmt.Sprintf("if (!!%s && %s != %q) { %s }", VarCurrentDisplayID, VarCurrentDisplayID, id, presets.CloseRightDrawerVarScript+";"+web.Plaid().EventFunc(actions.Detailing).Query(presets.ParamID, id).Go()),
		)
		return
	}
}

func renameVersionDialog(_ *presets.ModelBuilder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		versionName := ctx.R.FormValue("version_name")
		okAction := web.Plaid().
			URL(ctx.R.URL.Path).
			EventFunc(eventRenameVersion).
			Queries(ctx.Queries()).Go()

		r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{
			Name: actions.Dialog.PortalName(),
			Body: web.Scope(
				v.VDialog(
					v.VCard(
						v.VCardTitle(h.Text("Version")),
						v.VCardText(
							v.VTextField().Attr(web.VField("VersionName", versionName)...).Variant(v.FieldVariantUnderlined),
						),
						v.VCardActions(
							v.VSpacer(),
							v.VBtn("Cancel").
								Variant(v.VariantFlat).
								Class("ml-2").
								On("click", "locals.renameVersionDialog = false"),

							v.VBtn("OK").
								Color("primary").
								Variant(v.VariantFlat).
								Theme(v.ThemeDark).
								Attr("@click", "locals.renameVersionDialog = false; "+okAction),
						),
					),
				).MaxWidth("420px").Attr("v-model", "locals.renameVersionDialog"),
			).Init("{renameVersionDialog:true}").VSlot("{locals}"),
		})
		return
	}
}

func renameVersion(mb *presets.ModelBuilder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		if mb.Info().Verifier().Do(presets.PermUpdate).WithReq(ctx.R).IsAllowed() != nil {
			presets.ShowMessage(&r, perm.PermissionDenied.Error(), "warning")
			return
		}

		id := ctx.R.FormValue(presets.ParamID)
		obj := mb.NewModel()
		err = mb.Editing().Fetcher(obj, id, ctx)
		if err != nil {
			return
		}

		name := ctx.R.FormValue("VersionName")
		if err = reflectutils.Set(obj, "Version.VersionName", name); err != nil {
			return
		}

		if err = mb.Editing().Saver(obj, id, ctx); err != nil {
			return
		}

		listQueries := ctx.Queries().Get(presets.ParamListingQueries)
		r.RunScript = web.Plaid().
			URL(ctx.R.URL.Path).
			StringQuery(listQueries).
			Query(presets.ParamPortalID, ctx.R.FormValue(presets.ParamPortalID)).
			EventFunc(actions.UpdateListingDialog).
			Go()
		return
	}
}

func deleteVersionDialog(_ *presets.ModelBuilder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		versionName := ctx.R.FormValue("version_name")

		utilMsgr := i18n.MustGetModuleMessages(ctx.R, utils.I18nUtilsKey, Messages_en_US).(*utils.Messages)
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nPublishKey, Messages_en_US).(*Messages)

		r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{
			Name: presets.DeleteConfirmPortalName,
			Body: utils.DeleteDialog(
				msgr.DeleteVersionConfirmationText(versionName),
				"locals.deleteConfirmation = false;"+web.Plaid().
					URL(ctx.R.URL.Path).
					EventFunc(eventDeleteVersion).
					Queries(ctx.Queries()).Go(),
				utilMsgr),
		})
		return
	}
}

const paramCurrentDisplaySlug = "current_display_id"

func deleteVersion(mb *presets.ModelBuilder, pm *presets.ModelBuilder, db *gorm.DB) web.EventFunc {
	return wrapEventFuncWithShowError(func(ctx *web.EventContext) (web.EventResponse, error) {
		var r web.EventResponse

		if mb.Info().Verifier().Do(presets.PermDelete).WithReq(ctx.R).IsAllowed() != nil {
			return r, perm.PermissionDenied
		}

		slug := ctx.R.FormValue(presets.ParamID)
		if len(slug) <= 0 {
			return r, errors.New("no delete_id")
		}

		if err := mb.Listing().Deleter(mb.NewModel(), slug, ctx); err != nil {
			return r, err
		}

		currentDisplaySlug := ctx.R.FormValue(paramCurrentDisplaySlug)
		if slug == currentDisplaySlug {
			deletedVersion := mb.NewModel().(presets.SlugDecoder).PrimaryColumnValuesBySlug(slug)["version"]

			// find the older version first then find the max version
			version := mb.NewModel()
			db := utils.PrimarySluggerWhere(db, version, slug, "version").Order("version DESC").WithContext(ctx.R.Context())
			err := db.Where("version < ?", deletedVersion).First(version).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return r, err
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err := db.First(version).Error
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return r, err
				}
				if errors.Is(err, gorm.ErrRecordNotFound) {
					r.PushState = web.Location(nil).URL(pm.Info().ListingHref(presets.ParentsModelID(ctx.R)...))
					return r, nil
				}
			}

			currentDisplaySlug = version.(presets.SlugEncoder).PrimarySlug()
			web.AppendRunScripts(&r, fmt.Sprintf("%s = %q", VarCurrentDisplayID, currentDisplaySlug))

			if !pm.HasDetailing() {
				web.AppendRunScripts(&r,
					web.Plaid().EventFunc(actions.Edit).Query(presets.ParamID, currentDisplaySlug).Go(),
				)
			} else {
				web.AppendRunScripts(&r,
					presets.CloseRightDrawerVarScript,
					web.Plaid().EventFunc(actions.Detailing).Query(presets.ParamID, currentDisplaySlug).Go(),
				)
			}
		}

		listQuery, err := url.ParseQuery(ctx.Queries().Get(presets.ParamListingQueries))
		if err != nil {
			return r, err
		}
		if slug == cmp.Or(listQuery.Get("select_id"), listQuery.Get("f_select_id")) {
			listQuery.Set("select_id", currentDisplaySlug)
		}

		web.AppendRunScripts(&r,
			web.Plaid().
				URL(ctx.R.URL.Path).
				Queries(listQuery).
				EventFunc(actions.UpdateListingDialog).
				Query(presets.ParamPortalID, ctx.R.FormValue(presets.ParamPortalID)).
				Go(),
			// web.Plaid().EventFunc(actions.ReloadList).Go(), // TODO: This will reload the dialog list, I don't know how to reload the main list yet.
		)
		return r, nil
	})
}
