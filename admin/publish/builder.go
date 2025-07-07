package publish

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/iancoleman/strcase"
	"github.com/qor5/admin/v3/activity"
	"github.com/qor5/admin/v3/media/storage"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/utils"
	utils2 "github.com/qor5/admin/v3/utils/db_utils"
	"github.com/qor5/web/v3"
	"github.com/sunfmin/reflectutils"
	"github.com/theplant/htmlgo"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Builder struct {
	db      *gorm.DB
	storage storage.Storage
	// models            []*presets.ModelBuilder
	ab                *activity.Builder
	ctxValueProviders []ContextValueFunc
	afterInstallFuncs []func()
}

type ContextValueFunc func(ctx context.Context) context.Context

func New(db *gorm.DB, storage storage.Storage) *Builder {
	return &Builder{
		db:      db,
		storage: storage,
	}
}

func (b *Builder) Activity(v *activity.Builder) (r *Builder) {
	b.ab = v
	return b
}

func (b *Builder) AfterInstall(f func()) *Builder {
	b.afterInstallFuncs = append(b.afterInstallFuncs, f)
	return b
}

func (b *Builder) ModelInstall(pb *presets.Builder, m *presets.ModelBuilder) error {
	obj := m.NewModel()
	_ = obj.(presets.SlugEncoder)
	_ = obj.(presets.SlugDecoder)

	var withVersion bool

	if model, ok := obj.(VersionInterface); ok {
		withVersion = true
		if schedulePublishModel, ok := model.(ScheduleInterface); ok {
			VersionPublishModels[m.Info().URI()] = &Model{
				Record:  reflect.ValueOf(schedulePublishModel).Elem().Interface(),
				Builder: m,
			}
		}

		b.configVersionAndPublish(pb, m)
	} else {
		if schedulePublishModel, ok := obj.(ScheduleInterface); ok {
			NonVersionPublishModels[m.Info().URI()] = &Model{
				Record:  reflect.ValueOf(schedulePublishModel).Elem().Interface(),
				Builder: m,
			}
		}
	}

	if model, ok := obj.(ListInterface); ok {
		if schedulePublishModel, ok := model.(ScheduleInterface); ok {
			ListPublishModels[m.Info().URI()] = &Model{
				Record:  reflect.ValueOf(schedulePublishModel).Elem().Interface(),
				Builder: m,
			}
		}
	}

	if _, ok := obj.(StatusInterface); ok {
		if !withVersion {
			b.configVersionAndPublish(pb, m)
		}
		detailFields := m.Detailing().GetSections()
		for _, detailField := range detailFields {
			wrapper := func(in presets.ObjectBoolFunc) presets.ObjectBoolFunc {
				return func(obj interface{}, ctx *web.EventContext) bool {
					return in(obj, ctx) && EmbedStatus(obj).Status == StatusDraft
				}
			}
			detailField.WrapComponentEditBtnFunc(wrapper)
			detailField.WrapComponentHoverFunc(wrapper)
		}
	}

	registerEventFuncsForResource(m, b)
	return nil
}

func (b *Builder) configVersionAndPublish(pb *presets.Builder, m *presets.ModelBuilder) {
	ed := m.Editing().HiddenField(FieldOnlineUrl)
	creating := ed.Creating().HiddenField(FieldOnlineUrl).Except(VersionsPublishBar, FieldStatus)
	detailing := m.Detailing()

	detailing.HiddenField(FieldOnlineUrl).Field(VersionsPublishBar).ComponentFunc(DefaultVersionComponentFunc(m))
	listing := m.Listing()

	listing.WrapDeleteFunc(func(in presets.DeleteFunc) presets.DeleteFunc {
		return func(obj interface{}, id presets.ID, cascade bool, ctx *web.EventContext) (err error) {
			if obj, err = UnPublish.Execute(m, b, ActivityUnPublish, ctx, id); err != nil {
				return
			}
			return in(obj, id, cascade, ctx)
		}
	})

	if _, ok := m.Model().(VersionInterface); ok {
		panic("test deletion only by version")

		listing.WrapSearchFunc(makeSearchFunc(m, b.db))
		listing.RowMenu().RowMenuItem("Delete").ComponentFunc(func(rctx *presets.RecordMenuItemContext) htmlgo.HTMLComponent {
			// DeleteRowMenu should be disabled when using the version interface
			return nil
		})

		setter := makeSetVersionSetterFunc(b.db)
		ed.WrapPostSetterFunc(setter)
		creating.WrapPostSetterFunc(setter)
		configureVersionListDialog(b.db, pb, m)
	}

	listing.Field(ListingFieldDraftCount).ComponentFunc(DraftCountComponentFunc(b.db))
	listing.Field(ListingFieldLive).ComponentFunc(LiveComponentFunc(b.db, &LiveChipsListBuilder))
	detailing.Field(ListingFieldLive).ComponentFunc(LiveComponentFunc(b.db, &LiveChipsFormBuilder))
}

func makeSearchFunc(mb *presets.ModelBuilder, db *gorm.DB) func(searcher presets.SearchFunc) presets.SearchFunc {
	return func(searcher presets.SearchFunc) presets.SearchFunc {
		return func(model interface{}, params *presets.SearchParams, ctx *web.EventContext) (r interface{}, totalCount int, err error) {
			stmt := &gorm.Statement{DB: db}
			stmt.Parse(model)
			tn := stmt.Schema.Table

			var pks []string
			condition := ""
			for _, f := range stmt.Schema.Fields {
				if f.Name == "DeletedAt" {
					condition = "WHERE deleted_at IS NULL"
				}
			}
			for _, f := range stmt.Schema.PrimaryFields {
				if f.Name != "Version" {
					pks = append(pks, f.DBName)
				}
			}
			pkc := strings.Join(pks, ",")

			sql := fmt.Sprintf(`
			(%s, version) IN (
				SELECT %s, version
				FROM (
					SELECT %s, version,
						ROW_NUMBER() OVER (PARTITION BY %s ORDER BY CASE WHEN status = '%s' THEN 0 ELSE 1 END, version DESC) as rn
					FROM %s %s
				) subquery
				WHERE subquery.rn = 1
			)`, pkc, pkc, pkc, pkc, StatusOnline, tn, condition)

			if _, ok := mb.NewModel().(ScheduleInterface); ok {
				// Also need to get the most recent planned to publish
				sql = fmt.Sprintf(`
				(%s, version) IN (
					SELECT %s, version
					FROM (
						SELECT %s, version,
							ROW_NUMBER() OVER (
								PARTITION BY %s 
								ORDER BY 
									CASE WHEN status = '%s' THEN 0 ELSE 1 END,
									CASE 
										WHEN scheduled_start_at >= now() THEN scheduled_start_at
										ELSE NULL 
									END,
									version DESC
							) as rn
						FROM %s %s
					) subquery
					WHERE subquery.rn = 1
				)`, pkc, pkc, pkc, pkc, StatusOnline, tn, condition)
			}

			con := presets.SQLCondition{
				Query: sql,
			}
			params.SQLConditions = append(params.SQLConditions, &con)

			return searcher(model, params, ctx)
		}
	}
}

func makeSetVersionSetterFunc(db *gorm.DB) func(presets.SetterFunc) presets.SetterFunc {
	return func(in presets.SetterFunc) presets.SetterFunc {
		return func(obj interface{}, ctx *web.EventContext) {
			if ctx.Param(presets.ParamID) == "" {
				version := fmt.Sprintf("%s-v01", db.NowFunc().Format("2006-01-02"))
				if err := reflectutils.Set(obj, "Version.Version", version); err != nil {
					return
				}
				if err := reflectutils.Set(obj, "Version.VersionName", version); err != nil {
					return
				}
			}
			if in != nil {
				in(obj, ctx)
			}
		}
	}
}

func (b *Builder) Install(pb *presets.Builder) error {
	pb.FieldDefaults(presets.LIST).
		FieldType(Status{}).
		ComponentFunc(StatusListFunc(&LiveChipsListBuilder))

	pb.I18n().
		RegisterForModule(language.English, I18nPublishKey, Messages_en_US).
		RegisterForModule(language.SimplifiedChinese, I18nPublishKey, Messages_zh_CN).
		RegisterForModule(language.Japanese, I18nPublishKey, Messages_ja_JP)

	utils.Install(pb)
	for _, f := range b.afterInstallFuncs {
		f()
	}
	return nil
}

func (b *Builder) ContextValueFuncs(vs ...ContextValueFunc) *Builder {
	b.ctxValueProviders = append(b.ctxValueProviders, vs...)
	return b
}

func (b *Builder) WithContextValues(ctx context.Context) context.Context {
	for _, v := range b.ctxValueProviders {
		ctx = v(ctx)
	}
	return ctx
}

// 幂等
func (b *Builder) Publish(mb *presets.ModelBuilder, record interface{}, ctx context.Context) (err error) {
	err = utils2.Transact(b.db, func(tx *gorm.DB) (err error) {
		if cb, ok := mb.GetData(ModelPublishCallbackKey).(ModelPublishCallback); ok {
			var done func(err error) error
			if done, err = cb(tx, ctx, record); err != nil {
				return
			}
			if done != nil {
				defer func() {
					err = done(err)
				}()
			}
		}
		// publish content
		if r, ok := record.(PublishInterface); ok {
			var objs []*PublishAction
			objs, err = r.GetPublishActions(mb, b.db, ctx, b.storage)
			if err != nil {
				return
			}
			if err = UploadOrDelete(objs, b.storage); err != nil {
				return
			}
		}

		// update status
		if r, ok := record.(StatusInterface); ok {
			now := b.db.NowFunc()
			if version, ok := record.(VersionInterface); ok {
				var modelSchema *schema.Schema
				modelSchema, err = schema.Parse(record, &sync.Map{}, b.db.NamingStrategy)
				if err != nil {
					return
				}
				scope := setPrimaryKeysConditionWithoutVersion(b.db.Model(reflect.New(modelSchema.ModelType).Interface()), record, modelSchema).Where("version <> ? AND status = ?", version.EmbedVersion().Version, StatusOnline)

				oldVersionUpdateMap := make(map[string]interface{})
				if _, ok := record.(ScheduleInterface); ok {
					oldVersionUpdateMap["scheduled_end_at"] = nil
					oldVersionUpdateMap["actual_end_at"] = &now
				}
				if _, ok := record.(ListInterface); ok {
					oldVersionUpdateMap["list_deleted"] = true
				}
				oldVersionUpdateMap["status"] = StatusOffline
				if err = scope.Updates(oldVersionUpdateMap).Error; err != nil {
					return
				}
			}
			updateMap := make(map[string]interface{})

			if r, ok := record.(ScheduleInterface); ok {
				r.EmbedSchedule().ActualStartAt = &now
				r.EmbedSchedule().ScheduledStartAt = nil
				updateMap["scheduled_start_at"] = r.EmbedSchedule().ScheduledStartAt
				updateMap["actual_start_at"] = r.EmbedSchedule().ActualStartAt
			}
			if _, ok := record.(ListInterface); ok {
				updateMap["list_updated"] = true
			}
			updateMap["status"] = StatusOnline
			updateMap["online_url"] = r.EmbedStatus().OnlineUrl
			if err = b.db.Model(record).Updates(updateMap).Error; err != nil {
				return
			}
		}

		// publish callback
		if r, ok := record.(AfterPublishInterface); ok {
			if err = r.AfterPublish(b.db, b.storage, ctx); err != nil {
				return
			}
		}
		return
	})
	return
}

func (b *Builder) UnPublish(mb *presets.ModelBuilder, record interface{}, ctx context.Context) (err error) {
	err = utils2.Transact(b.db, func(tx *gorm.DB) (err error) {
		if cb, ok := mb.GetData(ModelUnpublishCallbackKey).(ModelUnpublishCallback); ok {
			var done func(err error) error
			if done, err = cb(tx, ctx, record); err != nil {
				return
			}
			if done != nil {
				defer func() {
					err = done(err)
				}()
			}
		}

		// unpublish content
		if r, ok := record.(UnPublishInterface); ok {
			var objs []*PublishAction
			objs, err = r.GetUnPublishActions(mb, b.db, ctx, b.storage)
			if err != nil {
				return
			}
			if err = UploadOrDelete(objs, b.storage); err != nil {
				return
			}
		}

		// update status
		if _, ok := record.(StatusInterface); ok {
			updateMap := make(map[string]interface{})
			if r, ok := record.(ScheduleInterface); ok {
				now := b.db.NowFunc()
				r.EmbedSchedule().ActualEndAt = &now
				r.EmbedSchedule().ScheduledEndAt = nil
				updateMap["scheduled_end_at"] = r.EmbedSchedule().ScheduledEndAt
				updateMap["actual_end_at"] = r.EmbedSchedule().ActualEndAt
			}
			if _, ok := record.(ListInterface); ok {
				updateMap["list_deleted"] = true
			}
			updateMap["status"] = StatusOffline
			if err = b.db.Model(record).Updates(updateMap).Error; err != nil {
				return
			}
		}

		// unpublish callback
		if r, ok := record.(AfterUnPublishInterface); ok {
			if err = r.AfterUnPublish(b.db, b.storage, ctx); err != nil {
				return
			}
		}
		return
	})
	return
}

func UploadOrDelete(objs []*PublishAction, Storage storage.Storage) (err error) {
	for _, obj := range objs {
		if obj.IsDelete {
			fmt.Printf("deleting %s \n", obj.Url)
			err = Storage.Delete(obj.Url)
		} else {
			fmt.Printf("uploading %s \n", obj.Url)
			_, err = Storage.Put(obj.Url, strings.NewReader(obj.Content))
		}
		if err != nil {
			return
		}
	}
	return nil
}

func setPrimaryKeysConditionWithoutVersion(db *gorm.DB, record interface{}, s *schema.Schema) *gorm.DB {
	querys := []string{}
	args := []interface{}{}
	for _, p := range s.PrimaryFields {
		if p.Name == "Version" {
			continue
		}
		val, _ := p.ValueOf(db.Statement.Context, reflect.ValueOf(record))
		querys = append(querys, fmt.Sprintf("%s = ?", strcase.ToSnake(p.Name)))
		args = append(args, val)
	}
	return db.Where(strings.Join(querys, " AND "), args...)
}

func setPrimaryKeysConditionWithoutFields(db *gorm.DB, record interface{}, s *schema.Schema, ignoreFields ...string) *gorm.DB {
	querys := []string{}
	args := []interface{}{}
	for _, p := range s.PrimaryFields {
		if slices.Contains(ignoreFields, p.Name) {
			continue
		}
		val, _ := p.ValueOf(db.Statement.Context, reflect.ValueOf(record))
		querys = append(querys, fmt.Sprintf("%s = ?", strcase.ToSnake(p.Name)))
		args = append(args, val)
	}
	return db.Where(strings.Join(querys, " AND "), args...)
}
