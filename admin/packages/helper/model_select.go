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
	"github.com/go-rvq/rvq/web/zeroer"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/sunfmin/reflectutils"
	"gorm.io/gorm"
)

type ModelSelectorConfiguror func(input *vx.VXAdvancedSelectBuilder)

func (f ModelSelectorConfiguror) Wrap(configuror ModelSelectorConfiguror) ModelSelectorConfiguror {
	if f == nil {
		return configuror
	}
	return func(input *vx.VXAdvancedSelectBuilder) {
		f(input)
		configuror(input)
	}
}

type key string

var ModelSelectorConfigurorKey key = "ModelSelectorBuilderConfiguror"
var ModelSelectorItemsSearcherKey key = "ModelSelectorItemsSearcher"
var ModelSelectorEncoderKey key = "ModelSelectorEncoder"

type ModelSelectorItemsSearcher func(field *presets.FieldContext, tagBuilder *web.VueEventTagBuilder) *web.VueEventTagBuilder

type RecordEncodeFactory struct {
	Name    string
	TextKey string
	REF     presets.RecordEncoderFactory[any]
}

func (ref *RecordEncodeFactory) Configure(l *presets.ListingBuilder, searcher *web.VueEventTagBuilder, configuror ModelSelectorConfiguror) ModelSelectorConfiguror {
	l.SetRecordEncoderFactory(ref.Name, ref.REF)
	searcher.Query(presets.ParamListingEncoder, ref.Name)
	return configuror.Wrap(func(input *vx.VXAdvancedSelectBuilder) {
		input.ItemText(ref.TextKey)
	})
}

func SetDetaultItemsSearch(mb *presets.ModelBuilder, s ModelSelectorItemsSearcher) {
	mb.SetData(ModelSelectorItemsSearcherKey, s)
}

func SetToTitleRecordEncoderFactory[T any](l *presets.ListingBuilder, ttf ...presets.ToTitleFactory[T]) {
	SetRecordEncoderFactory(
		l.ModelBuilder(),
		&RecordEncodeFactory{
			Name:    presets.ToTitleRecordEncoderName,
			TextKey: "Title",
			REF:     presets.ToTitleRecordEncoderFactory(l, ttf...),
		})
}

func SetToAnyTitleRecordEncoderFactory(l *presets.ListingBuilder, ttf ...presets.ToTitleFactory[any]) {
	SetRecordEncoderFactory(
		l.ModelBuilder(),
		&RecordEncodeFactory{
			Name:    presets.ToTitleRecordEncoderName,
			TextKey: "Title",
			REF:     presets.ToTitleRecordEncoderFactory(l, ttf...),
		})
}

func SetRecordEncoderFactory(mb *presets.ModelBuilder, enc *RecordEncodeFactory) {
	mb.SetData(ModelSelectorEncoderKey, enc)
}

type ModelSelectorBuilder struct {
	Model               *presets.ModelBuilder
	Field               string
	foreignModel        *presets.ModelBuilder
	parents             presets.ParentsModelIDResolver
	listingDisabled     bool
	detailingDisabled   bool
	creatingDisabled    bool
	itemsSearcher       ModelSelectorItemsSearcher
	configureSelector   func(fctx *presets.FieldContext, input *vx.VXAdvancedSelectBuilder)
	recordEncodeFactory *RecordEncodeFactory
	many                bool
}

func NewModelSelectorBuilder(model *presets.ModelBuilder, field string) *ModelSelectorBuilder {
	return &ModelSelectorBuilder{Model: model, Field: field}
}

func (b *ModelSelectorBuilder) ForeignModel() *presets.ModelBuilder {
	return b.foreignModel
}

func (b *ModelSelectorBuilder) SetForeignModel(foreignModel *presets.ModelBuilder) *ModelSelectorBuilder {
	b.foreignModel = foreignModel
	return b
}

func (b *ModelSelectorBuilder) ListingDisabled() bool {
	return b.listingDisabled
}

func (b *ModelSelectorBuilder) SetListingDisabled(v bool) *ModelSelectorBuilder {
	b.listingDisabled = v
	return b
}

func (b *ModelSelectorBuilder) DetailingDisabled() bool {
	return b.detailingDisabled
}

func (b *ModelSelectorBuilder) SetDetailingDisabled(v bool) *ModelSelectorBuilder {
	b.detailingDisabled = v
	return b
}

func (b *ModelSelectorBuilder) CreatingDisabled() bool {
	return b.creatingDisabled
}

func (b *ModelSelectorBuilder) SetCreatingDisabled(v bool) *ModelSelectorBuilder {
	b.creatingDisabled = v
	return b
}

func (b *ModelSelectorBuilder) Parents() presets.ParentsModelIDResolver {
	return b.parents
}

func (b *ModelSelectorBuilder) SetParents(parents presets.ParentsModelIDResolver) *ModelSelectorBuilder {
	b.parents = parents
	return b
}

func (b *ModelSelectorBuilder) ItemsSearcher() ModelSelectorItemsSearcher {
	return b.itemsSearcher
}

func (b *ModelSelectorBuilder) SetItemsSearcher(f ModelSelectorItemsSearcher) *ModelSelectorBuilder {
	b.itemsSearcher = f
	return b
}

func (b *ModelSelectorBuilder) ConfigureSelector() func(fctx *presets.FieldContext, input *vx.VXAdvancedSelectBuilder) {
	return b.configureSelector
}

func (b *ModelSelectorBuilder) SetConfigureSelector(configureSelector func(fctx *presets.FieldContext, input *vx.VXAdvancedSelectBuilder)) *ModelSelectorBuilder {
	b.configureSelector = configureSelector
	return b
}

func (b *ModelSelectorBuilder) RecordEncodeFactory() *RecordEncodeFactory {
	return b.recordEncodeFactory
}

func (b *ModelSelectorBuilder) SetRecordEncodeFactory(recordEncodeFactory *RecordEncodeFactory) *ModelSelectorBuilder {
	b.recordEncodeFactory = recordEncodeFactory
	return b
}

func (b *ModelSelectorBuilder) SetMany(v bool) *ModelSelectorBuilder {
	b.many = v
	return b
}

func (b *ModelSelectorBuilder) Many() bool {
	return b.many
}

func (b *ModelSelectorBuilder) Build() *ModelSelectorBuilder {
	if b.foreignModel == nil {
		field, _ := b.Model.ModelType().Elem().FieldByName(b.Field)
		b.foreignModel = b.Model.Builder().GetModel(field.Type)
	}

	b.Model.UpdateDataOperator(func(dataOperator presets.DataOperator) presets.DataOperator {
		return dataOperator.(*gorm2op.DataOperatorBuilder).WrapPrepare(func(old gorm2op.Preparer) gorm2op.Preparer {
			return func(db *gorm.DB, mode gorm2op.Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
				return old(db, mode, obj, id, params, ctx).Joins(b.Field)
			}
		})
	})

	itemsSearcher := b.itemsSearcher

	if itemsSearcher == nil {
		if c, _ := b.foreignModel.GetData(ModelSelectorItemsSearcherKey).(ModelSelectorItemsSearcher); c != nil {
			itemsSearcher = c
		} else {
			itemsSearcher = func(_ *presets.FieldContext, tagBuilder *web.VueEventTagBuilder) *web.VueEventTagBuilder {
				return tagBuilder
			}
		}
	}

	configureSelector := b.configureSelector
	if configureSelector == nil {
		configureSelector = func(fctx *presets.FieldContext, input *vx.VXAdvancedSelectBuilder) {}
	}

	if b.recordEncodeFactory == nil {
		b.recordEncodeFactory, _ = b.foreignModel.GetData(ModelSelectorEncoderKey).(*RecordEncodeFactory)
	}

	var formKeySufix string
	if reflectutils.GetType(b.Model.Model(), b.Field+"ID") != nil {
		formKeySufix = "ID"
	}

	configuror, _ := b.foreignModel.GetData(ModelSelectorConfigurorKey).(ModelSelectorConfiguror)

	searcher := web.GET().
		EventFunc(actions.ListData).Clone()

	if b.recordEncodeFactory != nil {
		configuror = b.recordEncodeFactory.Configure(b.foreignModel.Listing(), searcher, configuror)
	}

	setter := func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
		if field.ReadOnly {
			return
		}
		if b.many {
			var (
				i    int
				keys []string
			)
			fv := ctx.R.PostForm[field.FormKey]
			if len(fv) > 0 {
				for _, v := range fv {
					if v != "" {
						keys = append(keys, v)
					}
				}
			} else {
				for {
					key := fmt.Sprintf("%s[%d]", field.FormKey, i)
					fv = ctx.R.PostForm[key]
					if len(fv) == 0 {
						break
					}
					if fv[0] != "" {
						keys = append(keys, fv[0])
					}
					i++
				}
			}

			field := reflect.ValueOf(obj).Elem().FieldByName(field.Name)
			field.Set(reflect.Zero(field.Type()))

			for _, key := range keys {
				id, _ := b.foreignModel.ParseRecordID(key)
				if !id.IsZero() {
					entry := b.foreignModel.NewModel()
					id.SetTo(entry)
					field.Set(reflect.Append(field, reflect.ValueOf(entry)))
				}
			}
		} else {
			if v := ctx.R.Form[field.FormKey]; len(v) > 0 {
				var id presets.ID
				if v[0] == "" {
					f := reflect.ValueOf(obj).Elem().FieldByName(field.Name)
					f.Set(reflect.Zero(f.Type()))
					if formKeySufix != "" {
						f := reflect.ValueOf(obj).Elem().FieldByName(field.Name + formKeySufix)
						f.Set(reflect.Zero(f.Type()))
					}
				} else {
					if id, err = b.foreignModel.ParseRecordID(v[0]); err != nil {
						return
					}
					if !id.IsZero() {
						s, _ := b.foreignModel.CurrentDataOperator().Schema(obj)
						id.Related(s, b.Field+formKeySufix).SetTo(obj)
					}
				}
			}
		}
		return
	}

	comp := func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		if field.ReadOnly {
			if comp := b.ReadonlyComponent(field, ctx); comp == nil {
				return nil
			} else {
				return vx.VXReadonlyField(comp).
					Label(field.Label)
			}
		} else if mode := field.Mode.Dot(); mode.Is(presets.NEW) {
			if b.creatingDisabled {
				return nil
			}
		}

		searcher := searcher.Clone().
			URL(b.foreignModel.Info().ListingHref())

		val := field.Value()
		var (
			assign  any
			loadVal = func() {
				if b.many {
					var ids []string
					reflectutils.ForEach(val, func(item any) {
						id, _ := b.foreignModel.RecordID(item)
						ids = append(ids, id.String())
					})
					assign = ids
				} else {
					id, _ := b.foreignModel.RecordID(val)
					assign = id.String()
				}
			}
		)

		if val != nil {
			loadVal()
		} else {
			assign = ""
			if b.many {

			} else if val2 := reflectutils.MustGet(field.Obj, field.Name+formKeySufix); val2 != nil {
				if id := b.foreignModel.MustParseRecordID(fmt.Sprint(val2)); !id.IsZero() {
					assign = id.String()
					val = b.foreignModel.NewModel()
					id.SetTo(val)
					b.foreignModel.Fetcher(val, id, ctx)
					reflectutils.Set(field.Obj, field.Name, val)
				}
			}
		}

		selector := vx.VXSelectOne().
			Label(field.Label).
			Attr(web.VField(field.FormKey, assign)...).
			ItemValue("ID").
			ItemText("Text").
			ItemsSearcher(itemsSearcher(field, searcher)).
			Attr(":model-value", assign).
			Many(b.many)

		configureSelector(field, selector)

		if configuror != nil {
			configuror(selector)
		}

		if val == nil {
			if b.many {
			} else if formKeySufix != "" {
				id := reflectutils.MustGet(field.Obj, b.Field+formKeySufix)
				if !zeroer.IsZero(id) {
					val = b.foreignModel.NewModel()
					if err := b.foreignModel.CurrentDataOperator().Fetch(val, b.foreignModel.MustParseRecordID(fmt.Sprint(id)), field.EventContext); err != nil {
						field.Errors = append(field.Errors, fmt.Sprint("Fetch record: "+err.Error()))
					} else {
						reflectutils.Set(field.Obj, b.Field, val)
					}
				}
			}
		}

		if val != nil {
			if b.recordEncodeFactory != nil {
				selector.Items(b.recordEncodeFactory.REF.EncodeSlice(ctx, val))
			} else {
				selector.Items([]any{val})
			}
		}

		selector.ErrorMessages(field.Errors...)

		return selector
	}

	b.Model.WithEditingBuilders(func(e *presets.EditingBuilder) {
		e.Except(b.Field + "ID").
			Field(b.Field).
			ComponentFunc(comp).
			SetterFunc(setter)
	})

	if !b.detailingDisabled {
		d := b.Model.Detailing()
		d.Field(b.Field).
			ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) (comp h.HTMLComponent) {
				if comp = b.ReadonlyComponent(field, ctx); comp == nil {
					return
				}
				return vx.VXReadonlyField(comp).
					Label(field.Label)
			})

		d.Except(b.Field + "ID")
	}

	if !b.listingDisabled {
		l := b.Model.Listing()
		l.Field(b.Field).
			ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) (comp h.HTMLComponent) {
				return h.Td(b.ReadonlyComponent(field, ctx))
			})
		l.Except(b.Field + "ID")
	}

	return b
}

func (b *ModelSelectorBuilder) ReadonlyComponent(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	var (
		value = field.RawValue()
		texts []string
	)

	if value == nil {
		return nil
	}

	if b.many {
		reflectutils.ForEach(value, func(item any) {
			texts = append(texts, b.RecordToString(ctx, item, b.foreignModel.MustRecordID(item)))
		})
	} else {
		var id model.ID
		id = b.foreignModel.MustRecordID(value)

		if id.IsZero() {
			v := reflectutils.MustGet(field.Obj, b.Field+"ID")
			if id.IsZero() {
				return nil
			}
			id = b.foreignModel.MustParseRecordID(fmt.Sprint(v))
		}

		texts = []string{b.RecordToString(ctx, value, id)}
	}

	if len(texts) == 0 {
		return nil
	}

	return b.ReadonlyComponentOfRecord(value, texts, ctx, !field.Mode.Dot().Has(presets.LIST))
}

func (b *ModelSelectorBuilder) RecordToString(ctx *web.EventContext, record any, id model.ID) string {
	if id.IsZero() {
		return fmt.Sprintf("%s %v", b.foreignModel.TTitle(ctx.Context()), id)
	}
	return b.foreignModel.RecordTitle(record, ctx)
}

func (b *ModelSelectorBuilder) ReadonlyComponentOfRecord(record any, text []string, ctx *web.EventContext, link bool) h.HTMLComponent {
	var (
		om = presets.GetOverlay(ctx)

		each = func(record any, text string) h.HTMLComponent {
			var (
				id           = b.foreignModel.MustRecordID(record)
				detailingUri string
			)

			if !link {
				return h.Text(text)
			}

			if b.foreignModel.HasDetailing() {
				detailingUri = b.foreignModel.Info().ListingHref(id)
			} else {
				detailingUri = b.foreignModel.Info().EditingHref(id)
			}

			comp := v.VChip(h.Text(text))

			if !om.IsDialog() {
				var (
					uri     = b.foreignModel.Info().ListingHref()
					onclick = web.Plaid().URL(uri).
						Query(presets.ParamID, id.String())
				)

				if b.foreignModel.HasDetailing() {
					onclick.EventFunc(actions.Detailing)
				} else {
					onclick.EventFunc(actions.Edit)
				}

				if om.IsDrawer() {
					onclick.Query(presets.ParamOverlay, actions.Dialog)
				} else {
					onclick.Query(presets.ParamOverlay, actions.RightDrawer)
				}

				comp.SetAttr("@click.self.prevent", onclick.Go())
				comp.SetAttr("@click.middle",
					fmt.Sprintf(`(e) => e.view.window.open(%q, "_blank")`, detailingUri))
			}
			return comp
		}
	)

	if b.many {
		var (
			comp    h.HTMLComponents
			records = reflect.ValueOf(record)
		)
		for i, s := range text {
			comp = append(comp, each(records.Index(i).Interface(), s))
		}

		if len(text) > 0 {
			return v.VChipGroup(comp...)
		}

		return comp
	}
	return each(record, text[0])
}

func ModelSelect(model *presets.ModelBuilder, field string) {
	NewModelSelectorBuilder(model, field).Build()
}
