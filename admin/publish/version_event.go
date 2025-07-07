package publish

import (
	"cmp"
	"errors"
	"fmt"
	"net/url"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/utils"
	utils2 "github.com/go-rvq/rvq/admin/utils/db_utils"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	"gorm.io/gorm"
)

const (
	PortalSchedulePublishDialog = "publish_PortalSchedulePublishDialog"
	PortalPublishCustomDialog   = "publish_PortalPublishCustomDialog"

	VarCurrentDisplayID = "vars.publish_VarCurrentDisplayID"
)

func duplicateVersionAction(mb *presets.ModelBuilder, pb *Builder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		mid := mb.MustParseRecordID(ctx.Param(presets.ParamID))
		obj := mb.NewModel()
		if err = utils2.ModelIdWhere(pb.db, mb.NewModel(), mid).First(obj).Error; err != nil {
			return
		}

		ver, ok := obj.(VersionInterface)
		if !ok {
			err = errInvalidObject
			return
		}

		oldVersion := *ver.EmbedVersion()
		newVersion, err := ver.CreateVersion(pb.db, mid, mb.NewModel())
		if err != nil {
			return
		}

		if err = reflectutils.Set(obj, "Version", Version{
			Version:       newVersion, // In fact, it is also set in CreateVersion, just in case
			VersionName:   newVersion, // In fact, it is also set in CreateVersion, just in case
			ParentVersion: oldVersion.Version,
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

		if err = mb.Editing().Creating().Creator(obj, ctx); err != nil {
			presets.ShowMessage(&r, err.Error(), "error")
			return
		}

		if !mb.HasDetailing() {
			// close dialog and open editing
			web.AppendRunScripts(&r,
				presets.CloseListingDialogVarScript,
				web.Plaid().EventFunc(actions.Edit).Query(presets.ParamID, mid.String()).Go(),
			)
			return
		}
		// close dialog and open detailingDrawer
		web.AppendRunScripts(&r,
			presets.CloseListingDialogVarScript,
			presets.CloseRightDrawerVarScript,
			web.Plaid().EventFunc(actions.Detailing).Query(presets.ParamID, mid.String()).Go(),
		)

		msgr := GetMessages(ctx.Context())
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

		r.UpdatePortal(actions.Dialog.PortalName(),
			web.Scope(
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
			).LocalsInit("{renameVersionDialog:true}").Slot("{locals}"),
		)
		return
	}
}

func renameVersion(mb *presets.ModelBuilder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		mid := mb.MustParseRecordID(ctx.R.FormValue(presets.ParamID))

		if mb.Permissioner().Updater(ctx.R, mid, presets.ParentsModelID(ctx.R)...).Denied() {
			err = perm.PermissionDenied
			return
		}

		obj := mb.NewModel()
		err = mb.Editing().Fetcher(obj, mid, ctx)
		if err != nil {
			return
		}

		name := ctx.R.FormValue("VersionName")
		if err = reflectutils.Set(obj, "Version.VersionName", name); err != nil {
			return
		}

		if err = mb.Editing().Saver(obj, mid, ctx); err != nil {
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

		utilMsgr := i18n.MustGetModuleMessages(ctx.Context(), utils.I18nUtilsKey, Messages_en_US).(*utils.Messages)
		msgr := GetMessages(ctx.Context())

		r.UpdatePortal(presets.DeleteConfirmPortalName,
			utils.DeleteDialog(
				msgr.DeleteVersionConfirmationText(versionName),
				"locals.deleteConfirmation = false;"+web.Plaid().
					URL(ctx.R.URL.Path).
					EventFunc(eventDeleteVersion).
					Queries(ctx.Queries()).Go(),
				utilMsgr),
		)
		return
	}
}

const paramCurrentDisplaySlug = "current_display_id"

func deleteVersion(mb *presets.ModelBuilder, pm *presets.ModelBuilder, db *gorm.DB) web.EventFunc {
	return wrapEventFuncWithShowError(func(ctx *web.EventContext) (web.EventResponse, error) {
		var r web.EventResponse

		mid := mb.MustParseRecordID(ctx.R.FormValue(presets.ParamID))
		if mid.IsZero() {
			return r, errors.New("no delete_id")
		}

		if mb.Permissioner().Deleter(ctx.R, mid, presets.ParentsModelID(ctx.R)...).Denied() {
			return r, perm.PermissionDenied
		}

		if err := mb.Listing().Deleter(mb.NewModel(), mid, false, ctx); err != nil {
			return r, err
		}

		currentDisplaySlug := ctx.R.FormValue(paramCurrentDisplaySlug)
		if mid.String() == currentDisplaySlug {
			deletedVersion := mid.GetValue("Version").(string)

			// find the older version first then find the max version
			version := mb.NewModel()
			db := utils2.ModelIdWhere(db, version, mid, "Version").Order(mid.Schema.Table() + ".version DESC").WithContext(ctx.R.Context())
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
		if mid.String() == cmp.Or(listQuery.Get("select_id"), listQuery.Get("f_select_id")) {
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
