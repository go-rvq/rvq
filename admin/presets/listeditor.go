package presets

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/web/v3/zeroer"
	. "github.com/qor5/x/v3/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
)

type ListEditorBuilder struct {
	fieldContext           *FieldContext
	value                  interface{}
	displayFieldInSorter   string
	addListItemRowEvent    string
	removeListItemRowEvent string
	sortListItemsEvent     string
}

type ListSorter struct {
	Items []ListSorterItem `json:"items"`
}

type ListSorterItem struct {
	Index int    `json:"index"`
	Label string `json:"label"`
}

func NewListEditor(v *FieldContext) *ListEditorBuilder {
	return &ListEditorBuilder{
		fieldContext:           v,
		addListItemRowEvent:    actions.AddRowEvent,
		removeListItemRowEvent: actions.RemoveRowEvent,
		sortListItemsEvent:     actions.SortEvent,
	}
}

func (b *ListEditorBuilder) Value(v interface{}) (r *ListEditorBuilder) {
	if v == nil {
		return b
	}
	if reflect.TypeOf(v).Kind() != reflect.Slice {
		panic("value must be slice")
	}
	b.value = v
	return b
}

func (b *ListEditorBuilder) DisplayFieldInSorter(v string) (r *ListEditorBuilder) {
	b.displayFieldInSorter = v
	return b
}

func (b *ListEditorBuilder) AddListItemRowEvent(v string) (r *ListEditorBuilder) {
	if v == "" {
		return b
	}
	b.addListItemRowEvent = v
	return b
}

func (b *ListEditorBuilder) RemoveListItemRowEvent(v string) (r *ListEditorBuilder) {
	if v == "" {
		return b
	}
	b.removeListItemRowEvent = v
	return b
}

func (b *ListEditorBuilder) SortListItemsEvent(v string) (r *ListEditorBuilder) {
	if v == "" {
		return b
	}
	b.sortListItemsEvent = v
	return b
}

func (b *ListEditorBuilder) Component(ctx *web.EventContext) h.HTMLComponent {
	msgr := MustGetMessages(ctx.Context())
	formKey := b.fieldContext.FormKey
	var form h.HTMLComponent
	if b.value != nil {
		form = b.fieldContext.Nested.FieldsBuilder().
			ToComponentForEach(&ToComponentOptions{}, b.fieldContext, b.value, b.fieldContext.Mode, ctx, func(obj interface{}, path FieldPath, formKey string, content h.HTMLComponent, ctx *web.EventContext) h.HTMLComponent {
				if zeroer.IsNil(obj) {
					return nil
				}
				return VCard(
					h.If(!b.fieldContext.ReadOnly,
						VToolbar(
							web.Slot(VBtn("").
								Color("error").
								Variant(VariantText).
								Density(DensityCompact).
								Icon("mdi-delete").
								Attr("@click", web.Plaid().
									URL(b.fieldContext.ModelInfo.ListingHref(ParentsModelID(ctx.R)...)).
									EventFunc(b.removeListItemRowEvent).
									Queries(ctx.Queries()).
									Query(ParamRemoveRowFormKey, formKey).
									Go()),
							).Name("append"),
						).Density(DensityCompact).AutoHeight(true),
					),
					VCardText(content),
				).Variant(VariantOutlined)
			})
	}

	isSortStart := ctx.R.FormValue(ParamIsStartSort) == "1" && ctx.R.FormValue(ParamSortSectionFormKey) == formKey
	haveSorterIcon := true
	var sorter h.HTMLComponent
	var sorterData ListSorter
	if b.value != nil {
		deletedIndexes := ContextModifiedIndexesBuilder(ctx)

		deletedIndexes.SortedForEach(b.value, formKey, func(obj interface{}, i int) {
			if deletedIndexes.DeletedContains(b.fieldContext.FormKey, i) {
				return
			}
			label := ""
			if b.displayFieldInSorter != "" {
				label = fmt.Sprint(reflectutils.MustGet(obj, b.displayFieldInSorter))
			} else {
				label = fmt.Sprintf("Item %d", i)
			}
			sorterData.Items = append(sorterData.Items, ListSorterItem{Label: label, Index: i})
		})
	}
	if len(sorterData.Items) < 2 {
		haveSorterIcon = false
	}
	if haveSorterIcon && isSortStart {
		sorter = VCard(
			VList(
				h.Tag("vx-draggable").Attr("v-model", "locals.items", "handle", ".handle", "animation", "300", "item-key", "index").Children(
					h.Template().Attr("#item", " { element } ").Children(
						VListItem(
							web.Slot(
								VIcon("mdi-drag").Class("handle mx-2 cursor-grab"),
							).Name("prepend"),
							VListItemTitle(h.Text("{{element.label}}")),
							VDivider(),
						),
					),
				),
			).Class("pa-0")).Variant(VariantOutlined).Class("mx-0 mt-1 mb-4")
	}
	return h.Div(
		web.Scope(
			h.If(!b.fieldContext.ReadOnly,
				h.Div(
					h.Label(b.fieldContext.Label).Class("v-label theme--light text-caption"),
					VSpacer(),
					h.If(haveSorterIcon,
						h.If(!isSortStart,
							VBtn("").
								Variant(VariantText).
								Icon("mdi-sort-variant").
								Class("mt-n4").
								Attr("@click",
									web.Plaid().
										URL(b.fieldContext.ModelInfo.ListingHref(ParentsModelID(ctx.R)...)).
										EventFunc(b.sortListItemsEvent).
										Queries(ctx.Queries()).
										Query(ParamID, ctx.R.FormValue(ParamID)).
										Query(ParamOverlay, ctx.R.FormValue(ParamOverlay)).
										Query(ParamSortSectionFormKey, b.fieldContext.FormKey).
										Query(ParamIsStartSort, "1").
										Go(),
								),
						).Else(
							VBtn("").
								Variant(VariantText).
								Icon("mdi-check").
								Class("mt-n4").
								Attr("@click",
									web.Plaid().
										URL(b.fieldContext.ModelInfo.ListingHref(ParentsModelID(ctx.R)...)).
										EventFunc(b.sortListItemsEvent).
										Queries(ctx.Queries()).
										Query(ParamID, ctx.R.FormValue(ParamID)).
										Query(ParamOverlay, ctx.R.FormValue(ParamOverlay)).
										Query(ParamSortSectionFormKey, b.fieldContext.FormKey).
										FieldValue(ParamSortResultFormKey, web.Var("JSON.stringify(locals.items)")).
										Query(ParamIsStartSort, "0").
										Go(),
								),
						),
					),
				).Class("d-flex align-end"),
			),
			sorter,
			h.Div(
				form,
				h.If(!b.fieldContext.ReadOnly,
					VBtn(msgr.AddRow).
						Variant(VariantText).
						Color("primary").
						Attr("@click", web.Plaid().
							URL(b.fieldContext.ModelInfo.ListingHref(ParentsModelID(ctx.R)...)).
							EventFunc(b.addListItemRowEvent).
							Queries(web.Query(ctx.Queries()).
								Set(ParamAddRowFormKey, b.fieldContext.FormKey).
								URLValues()).
							Go(),
						),
				),
			).Attr("v-show", h.JSONString(!isSortStart)).
				Class("mt-1 mb-4"),
		).LocalsInit(h.JSONString(sorterData)).Slot("{ locals }"),
	)
}

func addListItemRow(mb *ModelBuilder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		var mid ID
		if mid, err = mb.ParseRecordID(ctx.R.FormValue(ParamID)); err != nil {
			return
		}

		me := mb.Editing()
		if mid.IsZero() {
			me = me.CreatingBuilder()
		}
		obj, _ := me.FetchAndUnmarshal(nil, mid, false, ctx)
		formKey := ctx.R.FormValue(ParamAddRowFormKey)
		t := reflectutils.GetType(obj, formKey+"[0]")
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		newVal := reflect.New(t).Interface()
		if err = reflectutils.Set(obj, formKey+"[]", newVal); err != nil {
			return
		}

		return me.respondFormEdit(ctx, obj)
	}
}

func removeListItemRow(mb *ModelBuilder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		me := mb.Editing()
		var mid ID
		if mid, err = mb.ParseRecordID(ctx.R.FormValue(ParamID)); err != nil {
			return
		}
		if mid.IsZero() {
			me = me.CreatingBuilder()
		}
		formKey := ctx.R.FormValue(ParamRemoveRowFormKey)
		lb := strings.LastIndex(formKey, "[")
		sliceField := formKey[0:lb]
		strIndex := formKey[lb+1 : strings.LastIndex(formKey, "]")]

		var index int
		index, err = strconv.Atoi(strIndex)
		if err != nil {
			return
		}

		obj, _ := me.FetchAndUnmarshal(nil, mid, false, ctx)

		ContextModifiedIndexesBuilder(ctx).AppendDeleted(sliceField, index)

		return me.respondFormEdit(ctx, obj)
	}
}

func sortListItems(mb *ModelBuilder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		me := mb.Editing()
		var mid ID
		if mid, err = mb.ParseRecordID(ctx.R.FormValue(ParamID)); err != nil {
			return
		}
		obj, _ := me.FetchAndUnmarshal(nil, mid, false, ctx)
		sortSectionFormKey := ctx.R.FormValue(ParamSortSectionFormKey)
		mib := ContextModifiedIndexesBuilder(ctx)

		isStartSort := ctx.R.FormValue(ParamIsStartSort)
		if isStartSort != "1" {
			sortResult := ctx.R.FormValue(ParamSortResultFormKey)

			var result []ListSorterItem
			err = json.Unmarshal([]byte(sortResult), &result)
			if err != nil {
				return
			}
			var indexes []string
			for _, i := range result {
				indexes = append(indexes, fmt.Sprint(i.Index))
			}
			mib.SetSorted(sortSectionFormKey, indexes)
		}

		return me.respondFormEdit(ctx, obj)
	}
}

func RemoveEmptySliceItems(obj any, mib *ModifiedIndexesBuilder) func() {
	type State struct {
		key   string
		value any
	}

	var old []State

	for k, m := range mib.deletedValues {
		slice := reflectutils.MustGet(obj, k)
		sliceV := reflect.ValueOf(slice)
		if sliceV.Len() == 0 {
			continue
		}

		old = append(old, State{
			key:   k,
			value: slice,
		})

		length := sliceV.Len()
		newLength := length - len(m)
		newSlice := reflect.MakeSlice(sliceV.Type(), newLength, newLength)

		for i, j := 0, 0; i < length; i++ {
			if _, ok := m[i]; ok {
				continue
			}
			newSlice.Index(j).Set(sliceV.Index(i))
			j++
		}

		fmt.Println(newSlice.Interface())
		if err := reflectutils.Set(obj, k, newSlice.Interface()); err != nil {
			panic(err)
		}
	}

	return func() {
		for _, state := range old {
			if err := reflectutils.Set(obj, state.key, state.value); err != nil {
				panic(err)
			}
		}
	}
}
