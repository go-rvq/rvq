package presets

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/web/zeroer"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/sunfmin/reflectutils"
)

type ListEditorItemType uint8

func (t ListEditorItemType) String() string {
	switch t {
	case ListEditorItemTypeNew:
		return "new"
	case ListEditorItemTypeEdit:
		return "edit"
	case ListEditorItemTypeDelete:
		return "delete"
	default:
		return "-"
	}
}

func (t ListEditorItemType) FormKeyValue(itemFormKey string) (string, string) {
	return itemFormKey + ".__type", t.String()
}

const (
	ListEditorItemTypeNew = iota
	ListEditorItemTypeEdit
	ListEditorItemTypeDelete
)

// Per-item form metadata emitted by the list editor. Each item posts these
// alongside its real fields, keyed as `<itemFormKey>.<field>` (e.g.
// `Parcelas[2].__deleted`). They replace the old global ModifiedIndexesBuilder:
// the server reads each item's own state instead of a shared hidden index.
//
//   - ListEditorPositionField: the item's desired position (used to reorder the
//     decoded slice by intent instead of by array index).
//   - ListEditorNewField:      "true" for items created in the browser.
//   - ListEditorDeletedField:  "true" for items the user removed; a non-new
//     deleted item is posted as `{ID, __deleted:true}` so it can be reverted.
//   - ListEditorPresentField:  always "1"; lets the server tell "the list was in
//     the form but empty" (delete all) from "the field was absent" (leave as is).
//   - ListEditorPurgedField:   client-only "true" for a deleted item the user
//     dismissed for good — the removed placeholder is hidden and can no longer be
//     reverted. The item stays __deleted, so the server still deletes it.
const (
	ListEditorPositionField = web.PosFieldSuffix
	ListEditorNewField      = "__new"
	ListEditorDeletedField  = "__deleted"
	ListEditorPresentField  = "__present"
	ListEditorPurgedField   = "__purged"
)

const (
	ListEditorAddRowParamJsonItems    = "presets_listEditorAddRowJsonItems"
	ListEditorAddRowParamCurrentItems = "presets_listEditorCurrentItems"
)

type ListEditorItemContext struct {
	FieldContext *FieldContext
	ItemVar      string
	RemoveEvent  string
	Item         any
	Path         FieldPath
	// ItemKey internal key of item in list
	ItemKey             int
	ItemFormKey         string
	ItemPositionFormKey string
	Deleted             bool
	New                 bool
	Body                func() h.HTMLComponent
	ItemType            ListEditorItemType
	// Label is the display label of the item (used by the deleted placeholder).
	Label string
}

// deletedKey is the flat form key holding this item's deletion flag,
// e.g. `Parcelas[2].__deleted`.
func (c *ListEditorItemContext) deletedKey() string {
	return c.ItemFormKey + "." + ListEditorDeletedField
}

// DeleteExpr is the @click expression that marks the item as deleted. Deletion
// is toggled on the reactive form (client-side) instead of round-tripping to the
// server, so the item's inputs stay in the DOM and it can be reverted.
func (c *ListEditorItemContext) DeleteExpr() string {
	return fmt.Sprintf("form[%q] = true", c.deletedKey())
}

// RevertExpr is the @click expression that restores a deleted item.
func (c *ListEditorItemContext) RevertExpr() string {
	return fmt.Sprintf("form[%q] = false", c.deletedKey())
}

// DeletedCond is the truthy Vue condition for this item being deleted.
func (c *ListEditorItemContext) DeletedCond() string {
	return fmt.Sprintf("form[%q]", c.deletedKey())
}

// purgedKey is the flat form key holding this item's "removed for good" flag.
func (c *ListEditorItemContext) purgedKey() string {
	return c.ItemFormKey + "." + ListEditorPurgedField
}

// PurgeExpr is the @click expression that dismisses a deleted item for good: it
// stays deleted (so the server removes it) but its placeholder is hidden and it
// can no longer be reverted.
func (c *ListEditorItemContext) PurgeExpr() string {
	return fmt.Sprintf("form[%q] = true", c.purgedKey())
}

// PurgedCond is the truthy Vue condition for this item having been purged.
func (c *ListEditorItemContext) PurgedCond() string {
	return fmt.Sprintf("form[%q]", c.purgedKey())
}

// DeletedVisibleCond is the condition for showing the deleted placeholder: the
// item is deleted and has not been purged.
func (c *ListEditorItemContext) DeletedVisibleCond() string {
	return c.DeletedCond() + " && !" + c.PurgedCond()
}

type ListEditorContainerContext struct {
	FieldContext *FieldContext
}

type ListEditorComponentBuilder interface {
	Container(ctx *ListEditorContainerContext, items h.HTMLComponents) h.HTMLComponent
	Item(ctx *ListEditorItemContext) h.HTMLComponent
	// DeletedItem renders a persisted item that the user removed. It keeps the
	// item's place in the list (so order is preserved) and offers a way to
	// revert the deletion. It is shown, reactively, in place of Item while the
	// item is flagged deleted.
	DeletedItem(ctx *ListEditorItemContext) h.HTMLComponent
}

type ListEditorListBuilder struct {
}

func (b *ListEditorListBuilder) Container(ctx *ListEditorContainerContext, items h.HTMLComponents) h.HTMLComponent {
	return items
}

func (b *ListEditorListBuilder) Item(ctx *ListEditorItemContext) h.HTMLComponent {
	return VCard(
		h.If(!ctx.FieldContext.ReadOnly,
			VToolbar(
				web.Slot(VBtn("").
					Color("error").
					Variant(VariantText).
					Density(DensityCompact).
					Icon("mdi-delete").
					// mark the item deleted in the reactive scope instead of
					// round-tripping to the server, so it can be reverted and its
					// data is not lost.
					Attr("@click", ctx.DeleteExpr()),
				).Name("append"),
			).Density(DensityCompact).AutoHeight(true),
		),
		VCardText(ctx.Body()),
	).Variant(VariantOutlined).Attr("v-show", "!"+ctx.DeletedCond())
}

func (b *ListEditorListBuilder) DeletedItem(ctx *ListEditorItemContext) h.HTMLComponent {
	msgr := MustGetMessages(ctx.FieldContext.EventContext.Context())
	label := ctx.Label
	if label == "" {
		label = msgr.ListEditorDeletedItem
	}
	return VCard(
		VCardText(
			h.Div(
				VIcon("mdi-delete-outline").Class("mr-2").Color("error"),
				h.Span(label).Class("text-decoration-line-through text-medium-emphasis"),
				VSpacer(),
				h.If(!ctx.FieldContext.ReadOnly,
					h.Div(
						VBtn(msgr.ListEditorRevertDeletion).
							Variant(VariantText).
							Color("primary").
							Density(DensityCompact).
							PrependIcon("mdi-undo-variant").
							Attr("@click", ctx.RevertExpr()),
						VBtn(msgr.ListEditorRemoveItem).
							Variant(VariantText).
							Color("error").
							Density(DensityCompact).
							PrependIcon("mdi-delete-forever").
							// dismiss the removed item for good: stays deleted (server
							// removes it) but the placeholder is hidden and unrevertible.
							Attr("@click", ctx.PurgeExpr()),
					).Class("d-flex ga-2"),
				),
			).Class("d-flex align-center"),
		),
	).Variant(VariantTonal).Class("mb-0").Attr("v-show", ctx.DeletedVisibleCond())
}

type ListSorter struct {
	Items []ListSorterItem `json:"items"`
}

type ListSorterItem struct {
	Index int    `json:"index"`
	Label string `json:"label"`
}

type ListEditorItemAppendActionEventContext struct {
	*web.EventContext
	Model     *ModelBuilder
	items     any
	FieldName string
}

func (c *ListEditorItemAppendActionEventContext) Items() any {
	if c.items == nil {
		obj := c.Model.NewModel()
		formKey := c.R.FormValue(ParamAddRowFormKey)

		if err := c.EventContext.UnmarshalFormValues(c.R.PostForm, obj); err != nil {
			panic(err)
		}

		if items, err := reflectutils.Get(obj, formKey); err == nil {
			c.items = items
			return items
		}
	}
	return c.items
}

func (c *ListEditorItemAppendActionEventContext) AddItems(v any) {
	if rv := reflect.ValueOf(v); rv.Type().Kind() != reflect.Slice {
		panic("ListEditorItemAppendActionEventContext AddItems must be a slice")
	} else if rv.Len() > 0 {
		c.Resp.AppendRunScript(fmt.Sprintf("$listEditorAddItems(%s)", h.JSONString(v)))
	}
}

type ListEditorBuilder struct {
	fieldContext           *FieldContext
	value                  interface{}
	displayFieldInSorter   string
	addListItemRowEvent    string
	removeListItemRowEvent string
	sortListItemsEvent     string
	ComponentBuilder       ListEditorComponentBuilder
	addItemEvent           string
}

func NewListEditor(v *FieldContext) *ListEditorBuilder {
	return &ListEditorBuilder{
		fieldContext:           v,
		addListItemRowEvent:    actions.AddRowEvent,
		removeListItemRowEvent: actions.RemoveRowEvent,
		sortListItemsEvent:     actions.SortEvent,
		ComponentBuilder:       &ListEditorListBuilder{},
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

func (b *ListEditorBuilder) AddItemEvent(event string) *ListEditorBuilder {
	b.addItemEvent = event
	return b
}

func (b *ListEditorBuilder) RegisterItemEvent(mb *ModelBuilder, eventName string, handler func(ctx *ListEditorItemAppendActionEventContext) error) *ListEditorBuilder {
	b.addItemEvent = eventName
	mb.RegisterEventFunc(eventName, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		ctx.Resp = &r
		err = handler(&ListEditorItemAppendActionEventContext{
			EventContext: ctx,
			Model:        mb,
		})
		return
	})

	return b
}

// SetComponentBuilder swaps the component builder (e.g. a table renderer). A nil
// value keeps the current one.
func (b *ListEditorBuilder) SetComponentBuilder(v ListEditorComponentBuilder) (r *ListEditorBuilder) {
	if v != nil {
		b.ComponentBuilder = v
	}
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

// itemLabel is the display label of an item, from the configured sorter field
// when set, otherwise a positional fallback.
func (b *ListEditorBuilder) itemLabel(obj any, i int) string {
	if b.displayFieldInSorter != "" {
		return fmt.Sprint(reflectutils.MustGet(obj, b.displayFieldInSorter))
	}
	return fmt.Sprintf("Item %d", i+1)
}

func (b *ListEditorBuilder) BuildComponent(ctx *web.EventContext) h.HTMLComponent {
	var (
		form h.HTMLComponent
		// itemsState feeds the drag sorter's list (index + label per item).
		itemsState []map[string]any
		// itemFormKeys is the ordered list of each rendered item's form key,
		// used by Setup to seed per-item metadata (__pos, __deleted default).
		itemFormKeys []string

		msgr       = MustGetMessages(ctx.Context())
		formKey    = b.fieldContext.FormKey
		tempPortal = ctx.UID()
	)

	if b.value != nil {
		var (
			i     int
			items = b.fieldContext.Nested.FieldsBuilder().
				ToComponentForEach(&ToComponentOptions{}, b.fieldContext, b.value, b.fieldContext.Mode, ctx, func(obj interface{}, path FieldPath, itemFormKey string, content h.HTMLComponent, ctx *web.EventContext) h.HTMLComponent {
					if zeroer.IsNil(obj) {
						return nil
					}

					idx := i
					i++

					label := b.itemLabel(obj, idx)
					itemsState = append(itemsState, map[string]any{"index": idx, "label": label})
					itemFormKeys = append(itemFormKeys, itemFormKey)

					itemCtx := &ListEditorItemContext{
						FieldContext:        b.fieldContext,
						ItemVar:             "",
						RemoveEvent:         b.removeListItemRowEvent,
						Item:                obj,
						Path:                path,
						ItemKey:             idx,
						ItemFormKey:         itemFormKey,
						ItemPositionFormKey: itemFormKey + "." + ListEditorPositionField,
						Label:               label,
						Deleted:             false,
						New:                 false,
						Body: func() h.HTMLComponent {
							return content
						},
					}

					// render both views; each builder toggles its own visibility via
					// ctx.DeletedCond() (Item hidden when deleted, DeletedItem shown),
					// so the deleted placeholder takes the item's place and list order
					// is preserved. No wrapping element is added here, so table rows
					// (<tr>) stay valid direct children of the container.
					return h.Components(
						b.ComponentBuilder.Item(itemCtx),
						b.ComponentBuilder.DeletedItem(itemCtx),
					)
				})
		)

		if i > 0 {
			form = b.ComponentBuilder.Container(&ListEditorContainerContext{
				FieldContext: b.fieldContext,
			}, items)
		}
	}

	var (
		sorter         h.HTMLComponent
		isSortStart    = ctx.R.FormValue(ParamIsStartSort) == "1" && ctx.R.FormValue(ParamSortSectionFormKey) == formKey
		haveSorterIcon = true
	)

	// the drag sorter reuses the reactive `items` state (built above) as its
	// v-model, so deletion state and ordering share one source of truth instead
	// of the global ModifiedIndexesBuilder.
	if len(itemsState) < 2 {
		haveSorterIcon = false
	}

	if haveSorterIcon && isSortStart {
		sorter = VCard(
			VList(
				h.Tag("vx-draggable").Attr("v-model", "items", "handle", ".handle", "animation", "300", "item-key", "index").Children(
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

	var (
		url          = b.fieldContext.ModelInfo.ListingHref(ParentsModelID(ctx.R)...)
		addItemClick string
		addItemTb    = web.Plaid().
				URL(url).
				EventFunc(b.addListItemRowEvent).
				Queries(web.Query(ctx.Queries()).
					Set(ParamAddRowFormKey, b.fieldContext.FormKey).
					URLValues())
		scopeSetup strings.Builder
	)

	fmt.Fprintf(&scopeSetup, `({ scope }) => {
			const formKey = "%s";
			// expose the ambient (root) form as the slot-scoped form binding so the
			// item inputs and the delete/revert expressions bind to the real parent
			// form -- NOT a fresh object (which would reset the field).
			scope.form = form
			// presence marker: always posted so the server can tell an empty list
			// that was submitted (delete all children) from a field that was
			// absent (leave children untouched).
			form[formKey + "." + %q] = "1"
			// seed each item's position on the reactive form so submit order is
			// preserved (the server reorders the decoded slice by __pos). The real
			// field values and the __deleted flag are bound directly by the item
			// components; nothing here overwrites them.
			const keys = %s
			keys && keys.forEach((k, i) => { form[k + %q] = i })
		`,
		formKey,
		ListEditorPresentField,
		h.JSONString(itemFormKeys),
		"."+ListEditorPositionField,
	)

	if len(b.addItemEvent) > 0 {
		addItemClick = "openItemsSelector"
		fmt.Fprintf(&scopeSetup, `
			scope.$listEditorAddItems = (items) => {
				if (!items) return;
				((typeof items) !== "string" && (items = JSON.stringify(items)))
				%s.formData({%s: items}).go()
			} 

			scope.openItemsSelector = () => {
				const items = {},
					prefix = formKey + "["

				Object.entries(form).forEach(([key, value]) => {
					if (key.startsWith(prefix)) items[key] = value
				});

				%s.form(items).scope({$listEditorAddItems: scope.$listEditorAddItems}).go()
			}`,
			addItemTb,
			ListEditorAddRowParamJsonItems,
			web.Plaid().
				EventFunc(b.addItemEvent).
				URL(url).
				Query(ParamTargetPortal, tempPortal).
				FormData(map[string]any{
					ParamAddRowFormKey: formKey,
				}).
				SkipFiles(true),
		)
	} else {
		addItemClick = addItemTb.
			Go()
	}

	scopeSetup.WriteByte('}')

	return h.Div(
		vue.UserComponent(
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
										FieldValue(ParamSortResultFormKey, web.Var("JSON.stringify(items)")).
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
						Attr("data-list-editor-add", formKey).
						Attr("@click", addItemClick),
				),
			).Attr("v-show", h.JSONString(!isSortStart)).
				Class("mt-1 mb-4"),
		).Scope("items", jsonOrEmptyArray(itemsState)).
			Scope("$listEditorAddItems").
			Scope("openItemsSelector").
			Scope("form").
			Setup(scopeSetup.String()),
		web.Portal().Name(tempPortal),
	)
}

// jsonOrEmptyArray marshals v, guaranteeing a JS array literal (never "null")
// so a Vue v-model bound to it stays iterable.
func jsonOrEmptyArray(v []map[string]any) string {
	if len(v) == 0 {
		return "[]"
	}
	return h.JSONString(v)
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

		if jsonItems := ctx.R.FormValue(ListEditorAddRowParamJsonItems); len(jsonItems) > 0 {
			// jsonItems is a JSON array of serialized items; decode them into the
			// slice element type and append each to obj's slice.
			itemsPtr := reflect.New(reflect.SliceOf(t))
			if err = json.Unmarshal([]byte(jsonItems), itemsPtr.Interface()); err != nil {
				return
			}
			items := itemsPtr.Elem()
			for i := 0; i < items.Len(); i++ {
				if err = reflectutils.Set(obj, formKey+"[]", items.Index(i).Interface()); err != nil {
					return
				}
			}
		} else {
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			newVal := reflect.New(t).Interface()
			if err = reflectutils.Set(obj, formKey+"[]", newVal); err != nil {
				return
			}
		}

		// prevent to create new form scope
		ctx.R.Form.Set(ParamEditFormUnscoped, "true")

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

		// prevent to create new form scope
		ctx.R.Form.Set(ParamEditFormUnscoped, "true")
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

		// prevent to create new form scope
		ctx.R.Form.Set(ParamEditFormUnscoped, "true")
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
