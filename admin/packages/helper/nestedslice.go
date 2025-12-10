package helper

import (
	"fmt"
	"reflect"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/sunfmin/reflectutils"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type NestedSliceBuilderInfo struct {
	JoinTable       *schema.Schema
	DeleteQuery     string
	LinkInsertQuery string
	Target          *presets.ModelBuilder
}

func (i *NestedSliceBuilderInfo) Delete(db *gorm.DB, parentID, id any) error {
	return db.Exec(i.DeleteQuery, parentID, id).Error
}

func (i *NestedSliceBuilderInfo) Deleter() func(db *gorm.DB, obj interface{}, id model.ID, cascade bool, ctx *web.EventContext) (err error) {
	return func(db *gorm.DB, obj interface{}, id model.ID, cascade bool, ctx *web.EventContext) (err error) {
		parentID := presets.ParentsModelID(ctx.R).Last().Value()
		return i.Delete(db, parentID, id)
	}
}

type NestedSliceItemsSearcher func(field *presets.FieldContext, b *web.VueEventTagBuilder) *web.VueEventTagBuilder

const FieldContextNestedSliceItemsSearcherKey = "NestedSliceItemsSearcher"

type NestedSliceBuilder struct {
	baseModel           *presets.ModelBuilder
	linkModel           *presets.ModelBuilder
	fieldModel          *presets.ModelBuilder
	targetModel         *presets.ModelBuilder
	fieldName           string
	inline              bool
	callback            func(SubGrupos *presets.ModelBuilder)
	callbackInfo        func(info *NestedSliceBuilderInfo, SubGrupos *presets.ModelBuilder)
	configureSelector   ModelSelectorConfiguror
	itemsSearcher       ModelSelectorItemsSearcher
	wrapItemsSearcher   NestedSliceItemsSearcher
	done                func(b *NestedSliceBuilder)
	recordEncodeFactory *RecordEncodeFactory
}

func NewNestedSliceBuilder(baseModel *presets.ModelBuilder, fieldName string) *NestedSliceBuilder {
	return &NestedSliceBuilder{baseModel: baseModel, fieldName: fieldName}
}

func (b *NestedSliceBuilder) Inline(inline bool) *NestedSliceBuilder {
	b.inline = inline
	return b
}

func (b *NestedSliceBuilder) SetRecordEncodeFactory(recordEncodeFactory *RecordEncodeFactory) *NestedSliceBuilder {
	b.recordEncodeFactory = recordEncodeFactory
	return b
}

func (b *NestedSliceBuilder) LinkModel() *presets.ModelBuilder {
	return b.linkModel
}

func (b *NestedSliceBuilder) SetLinkModel(linkModel *presets.ModelBuilder) *NestedSliceBuilder {
	b.linkModel = linkModel
	return b
}

func (b *NestedSliceBuilder) FieldModel() *presets.ModelBuilder {
	return b.fieldModel
}

func (b *NestedSliceBuilder) TargetModel() *presets.ModelBuilder {
	return b.targetModel
}

func (b *NestedSliceBuilder) Callback() func(FieldModel *presets.ModelBuilder) {
	return b.callback
}

func (b *NestedSliceBuilder) SetCallback(cb func(FieldModel *presets.ModelBuilder)) *NestedSliceBuilder {
	b.callback = cb
	return b
}

func (b *NestedSliceBuilder) CallbackInfo() func(info *NestedSliceBuilderInfo, FieldModel *presets.ModelBuilder) {
	return b.callbackInfo
}

func (b *NestedSliceBuilder) SetCallbackInfo(callbackInfo func(info *NestedSliceBuilderInfo, FieldModel *presets.ModelBuilder)) *NestedSliceBuilder {
	b.callbackInfo = callbackInfo
	return b
}

func (b *NestedSliceBuilder) ConfigureSelector() func(input *vx.VXAdvancedSelectBuilder) {
	return b.configureSelector
}

func (b *NestedSliceBuilder) SetConfigureSelector(configureSelector ModelSelectorConfiguror) *NestedSliceBuilder {
	b.configureSelector = configureSelector
	return b
}

func (b *NestedSliceBuilder) ItemsSearcher() ModelSelectorItemsSearcher {
	return b.itemsSearcher
}

func (b *NestedSliceBuilder) SetItemsSearcher(itemsSearcher ModelSelectorItemsSearcher) *NestedSliceBuilder {
	b.itemsSearcher = itemsSearcher
	return b
}

func (b *NestedSliceBuilder) WrapItemsSearcher() NestedSliceItemsSearcher {
	return b.wrapItemsSearcher
}

func (b *NestedSliceBuilder) SetWrapItemsSearcher(wrapItemsSearcher NestedSliceItemsSearcher) *NestedSliceBuilder {
	b.wrapItemsSearcher = wrapItemsSearcher
	return b
}

func (b *NestedSliceBuilder) OnDone() func(b *NestedSliceBuilder) {
	return b.done
}

func (b *NestedSliceBuilder) SetOnDone(f func(b *NestedSliceBuilder)) *NestedSliceBuilder {
	b.done = f
	return b
}

func (b *NestedSliceBuilder) Build() *NestedSliceBuilder {
	var (
		db       = b.baseModel.Builder().GetDataOperator().(*gorm2op.DataOperatorBuilder).DB().Session(&gorm.Session{})
		relation = db.Model(b.baseModel.Model()).Association(b.fieldName).Relationship
		// baseTable    = relation.Schema.Table
		relatedTable             = relation.FieldSchema.Table
		filterQuery, insertQuery string
		deleteQuery, linkQuery   string

		info       = &NestedSliceBuilderInfo{}
		itemsModel *presets.ModelBuilder
	)

	switch relation.Type {
	case schema.Many2Many:
		joinTable := relation.JoinTable
		filterQuery = fmt.Sprintf("EXISTS (SELECT 1 FROM %s rel WHERE rel.%s = ? AND rel.%s = %s.id)",
			joinTable.Name, joinTable.DBNames[0], joinTable.DBNames[1], relatedTable)
		insertQuery = fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (?, ?)",
			joinTable.Table, joinTable.DBNames[0], joinTable.DBNames[1])
		deleteQuery = fmt.Sprintf("DELETE FROM %s WHERE %s = ? AND %s = ?",
			joinTable.Table, joinTable.DBNames[0], joinTable.DBNames[1])
		linkQuery = fmt.Sprintf(m2mInsertQuery, relatedTable, joinTable.Table, joinTable.DBNames[0], joinTable.DBNames[1])
		info.JoinTable = joinTable
		info.DeleteQuery = deleteQuery
		info.LinkInsertQuery = linkQuery
	}

	b.baseModel.TakeFieldAsChild(b.fieldName, func(FieldModel *presets.ModelBuilder) {
		if len(FieldModel.Schema().PrimaryFields()) > 1 {
			panic("NestSlice doesn't supports ModelBuilder with many primary fields")
		}

		b.fieldModel = FieldModel

		switch relation.Type {
		case schema.HasMany:
			FieldModel.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
				return do.(*gorm2op.DataOperatorBuilder).
					WrapPrepare(func(old gorm2op.Preparer) gorm2op.Preparer {
						return func(db *gorm.DB, mode gorm2op.Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
							if !mode.Is(gorm2op.Fetch, gorm2op.FetchTitle) {
								parentID := presets.ParentsModelID(ctx.R).Last().Value()
								params.Where(relation.References[0].ForeignKey.DBName+" = ?", parentID)
							}
							return old(db, mode, obj, id, params, ctx)
						}
					})
			})
		case schema.Many2Many:
			FieldModel.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
				return do.(*gorm2op.DataOperatorBuilder).
					WrapPrepare(func(old gorm2op.Preparer) gorm2op.Preparer {
						return func(db *gorm.DB, mode gorm2op.Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
							if !mode.Is(gorm2op.Fetch, gorm2op.FetchTitle) {
								parentID := presets.ParentsModelID(ctx.R).Last().Value()
								params.Where(filterQuery, parentID)
							}
							return old(db, mode, obj, id, params, ctx)
						}
					}).
					SetCreator(func(db *gorm.DB, obj interface{}, ctx *web.EventContext) (err error) {
						return db.Transaction(func(db *gorm.DB) (err error) {
							if err = db.Create(obj).Error; err != nil {
								return
							}
							parentID := presets.ParentsModelID(ctx.R).Last().Value()
							itemID := FieldModel.MustRecordID(obj).Value()
							return db.Exec(insertQuery, parentID, itemID).Error
						})
					})
			})

			FieldModel.UpdateDataOperator(func(dataOperator presets.DataOperator) presets.DataOperator {
				return dataOperator.(*gorm2op.DataOperatorBuilder).SetDeleter(info.Deleter())
			})

			targetModel := b.linkModel
			if targetModel == nil {
				targetModel = FieldModel.Builder().GetModel(FieldModel.Model())
			}

			info.Target = targetModel
			b.targetModel = targetModel

			if b.configureSelector == nil {
				b.configureSelector, _ = targetModel.GetData(ModelSelectorConfigurorKey).(ModelSelectorConfiguror)
			}

			if b.recordEncodeFactory == nil {
				b.recordEncodeFactory, _ = b.targetModel.GetData(ModelSelectorEncoderKey).(*RecordEncodeFactory)
			}

			FieldModel.SetDetailingBuilder(targetModel.
				Detailing(),
			)

			FieldModel.MenuIcon(targetModel.GetMenuIcon())

			itemsModel = presets.NewModelBuilder(
				FieldModel.Builder(),
				&NestedSliceItems{},
				presets.ModelConfig().
					SetDataOperator(FieldModel.DataOperator()).
					SetUriName(FieldModel.UriName()).
					SetId(b.fieldName).
					SetLabel(targetModel.GetLabel())).
				MenuIcon(targetModel.GetMenuIcon()).
				ChildOf(FieldModel.Parent()).
				SetVerifierModel(FieldModel.Parent())

			editing := itemsModel.Editing()
			itemsSearcher := b.itemsSearcher

			searchTag := web.Plaid().
				EventFunc(actions.ListData)

			if b.recordEncodeFactory != nil {
				b.configureSelector = b.recordEncodeFactory.Configure(targetModel.Listing(), searchTag, b.configureSelector)
			}

			if itemsSearcher == nil {
				if c, _ := targetModel.GetData(ModelSelectorItemsSearcherKey).(ModelSelectorItemsSearcher); c != nil {
					itemsSearcher = c
				} else {
					itemsSearcher = func(field *presets.FieldContext, tagBuilder *web.VueEventTagBuilder) *web.VueEventTagBuilder {
						return tagBuilder
					}
				}
			}

			if b.wrapItemsSearcher == nil {
				b.wrapItemsSearcher = func(field *presets.FieldContext, b *web.VueEventTagBuilder) *web.VueEventTagBuilder {
					return b
				}
			}

			editing.Field("Items").
				SetHiddenLabel(true).
				ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
					uri := targetModel.Info().ListingHref()
					searcher := itemsSearcher(field, searchTag.Clone().URL(uri))
					searcher = b.wrapItemsSearcher(field, searcher)

					if v, _ := field.ContextValue(FieldContextNestedSliceItemsSearcherKey).(NestedSliceItemsSearcher); v != nil {
						searcher = v(field, searcher)
					}

					var values []any

					selector := vx.VXSelectMany().
						Label(field.Label).
						ItemValue("ID").
						ItemText("Text").
						ItemsSearcher(searcher)

					if b.configureSelector != nil {
						b.configureSelector(selector)
					}

					val := field.Value()
					if val != nil {
						if b.recordEncodeFactory != nil {
							selector.Items(b.recordEncodeFactory.REF.EncodeSlice(ctx, val))
						} else {
							selector.Items(val)
						}

						reflectutils.ForEach(val, func(v any) {
							values = append(values, targetModel.MustRecordID(v).String())
						})
					}

					return vue.FormField(selector).
						Value(field.FormKey, values).
						Bind()
				})

			editing.SetSkipFieldVerifier(func(name string) bool {
				return true
			})

			save := func(obj interface{}, ctx *web.EventContext) (err error) {
				p := obj.(*NestedSliceItems)
				var ids = make(map[string]any)
				for _, item := range p.Items {
					if item != "" {
						ids[item] = nil
					}
				}

				if len(ids) == 0 {
					return
				}

				var idSlice []any

				for key := range ids {
					idSlice = append(idSlice, FieldModel.MustParseRecordID(key).Value())
				}

				parentID := presets.ParentsModelID(ctx.R).Last().Value()
				err = db.Session(&gorm.Session{}).Exec(info.LinkInsertQuery, parentID, idSlice).Error
				return
			}

			editing.WrapSaveFunc(func(in presets.SaveFunc) presets.SaveFunc {
				return func(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
					return save(obj, ctx)
				}
			})

			editing.CreatingBuilder().WrapCreateFunc(func(in presets.CreateFunc) presets.CreateFunc {
				return save
			})

			FieldModel.SetCreatingBuilder(editing)
			FieldModel.EventsHub.Merge(&itemsModel.EventsHub)
		}

		if b.callback != nil {
			b.callback(FieldModel)
		}

		if b.callbackInfo != nil {
			b.callbackInfo(info, FieldModel)
		}

		switch relation.Type {
		case schema.HasMany:
			editing := FieldModel.Editing()
			creating := editing.CreatingBuilder()

			for _, ref := range relation.References {
				creating.Field(ref.ForeignKey.Name).SetDisabled(true)
				editing.Field(ref.ForeignKey.Name).SetDisabled(true)
			}

			if creating.ModelBuilder() == FieldModel {
				FieldModel.BeforeFormUnmarshallHandlers.AppendFunc(func(obj interface{}, ctx *web.EventContext) (err error) {
					parent := presets.ParentsModelID(ctx.R).Last()
					fieldsNames := make([]string, len(parent.Fields))
					for i := range fieldsNames {
						fieldsNames[i] = relation.References[i].ForeignKey.Name
					}
					parent.Related(FieldModel.Schema(), fieldsNames...).SetTo(obj)
					return
				})
			}
		}
	})

	if b.inline && itemsModel != nil {
		switch relation.Type {
		case schema.Many2Many:
			f := b.baseModel.Editing().Field(b.fieldName)
			f.ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
				f := itemsModel.Editing().Field("Items")
				return f.GetCompFunc()(field, ctx)
			})
			f.SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
				var ids = make(map[string]any)
				if v := ctx.R.Form[field.FormKey]; len(v) > 0 {
					for _, s := range v {
						if s != "" {
							ids[s] = nil
						}
					}
				} else {
					for _, key := range ctx.FormSliceKeys(field.FormKey) {
						if v, _ := ctx.R.MultipartForm.Value[key.Key]; len(v) > 0 {
							if s := v[0]; s != "" {
								ids[v[0]] = nil
							}
						}
					}
				}

				var records = reflect.ValueOf(b.targetModel.NewModelSlice()).Elem()

				for key := range ids {
					r := b.targetModel.NewModel()
					b.targetModel.MustParseRecordID(key).SetTo(r)
					records = reflect.Append(records, reflect.ValueOf(r))
				}

				reflectutils.Set(obj, field.Name, records.Interface())
				return
			})
			b.baseModel.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
				return do.(*gorm2op.DataOperatorBuilder).WithModeSplitCallbacks(gorm2op.Write, func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
					cb.Post(func(state *gorm2op.CallbackState) (err error) {
						v, _ := reflectutils.Get(state.Obj, b.fieldName)
						return state.DB.Association(b.fieldName).Replace(v)
					})
				})
			})
			b.baseModel.Detailing().Field(b.fieldName).ComponentFunc(presets.FieldComponentWrapper(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
				items := field.Value()
				if items != nil {
					var labels []string
					reflectutils.ForEach(items, func(i interface{}) {
						labels = append(labels, b.targetModel.RecordTitle(i, ctx))
					})
					if len(labels) > 0 {
						var comp = make(h.HTMLComponents, len(labels))
						for i, label := range labels {
							comp[i] = v.VChip(h.Text(label))
						}
						return v.VChipGroup(comp...)
					}
				}
				return nil
			}))
		}
	}

	if b.done != nil {
		b.done(b)
	}

	return b
}

func NestedSlice(baseModel *presets.ModelBuilder, fieldName string, cb ...func(mb *presets.ModelBuilder)) {
	b := NewNestedSliceBuilder(baseModel, fieldName)
	for _, b.callback = range cb {
	}
	b.Build()
}

type NestedSliceItems struct {
	Items []string
}

const m2mInsertQuery = `with 
data as (
	select id::BIGINT as f_id, ?::BIGINT as p_id FROM %[1]s WHERE id IN ?
) 
, data_ok as (
	select f_id, p_ID from data d where not exists (
	select 1 
	from %[2]s fp 
	where fp.%[3]s = p_id and fp.%[4]s = f_id)
)
insert into %[2]s (%[3]s, %[4]s) select p_id, f_id from data_ok;
`
