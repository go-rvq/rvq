package presets

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/mpvl/unique"
	"github.com/qor5/admin/v3/model"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/perm"
	v "github.com/qor5/x/v3/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
)

func (b *FieldBuilder) AutoNested(mb *ModelBuilder, fb *FieldsBuilder) (r *FieldBuilder) {
	switch b.rt.Kind() {
	case reflect.Slice:
		return b.Nested(NestedSlice(mb, fb))
	default:
		return b.Nested(NestedStruct(mb, fb))
	}
}

func (b *FieldBuilder) Nested(n Nested) (r *FieldBuilder) {
	b.nested = n
	n.Build(b)
	return b
}

type NameLabel struct {
	name        string
	label       string
	labelKey    string
	i18nLabel   func(ctx context.Context) string
	hiddenLabel bool
}

func (n *NameLabel) Name() string {
	return n.name
}

func (n *NameLabel) SetName(name string) {
	n.name = name
}

func (n *NameLabel) Label() string {
	return n.label
}

func (n *NameLabel) SetLabel(label string) {
	n.label = label
}

func (n *NameLabel) LabelKey() string {
	return n.labelKey
}

func (n *NameLabel) SetLabelKey(labelKey string) {
	n.labelKey = labelKey
}

func (n *NameLabel) I18nLabel() func(ctx context.Context) string {
	return n.i18nLabel
}

func (n *NameLabel) SetI18nLabel(i18nLabel func(ctx context.Context) string) {
	n.i18nLabel = i18nLabel
}

func (n *NameLabel) HiddenLabel() bool {
	return n.hiddenLabel
}

func (n *NameLabel) SetHiddenLabel(hiddenLabel bool) {
	n.hiddenLabel = hiddenLabel
}

type FieldLabelsFunc func(b *FieldsBuilder, ctx web.ContextValuer) map[string]string

type FieldsBuilder struct {
	builder                    *Builder
	model                      interface{}
	defaults                   *FieldDefaults
	fieldLabelsFromContextFunc FieldLabelsFunc
	fields                     FieldBuilders
	// string / []string / *FieldsSection
	fieldsLayout                 []interface{}
	skipFieldVerifier            func(name string) bool
	FieldToComponentSetup        FieldContextSetups
	BeforeSetObjectFieldsHandler SetObjectFieldsHandlers
	PostSetObjectFieldsHandler   SetObjectFieldsHandlers
	beginComponentFuncs          []func(info *ModelInfo, obj interface{}, mode FieldModeStack, parentFormValueKey string, ctx *web.EventContext) h.HTMLComponent
	hiddenFields                 []string
}

func (b *FieldsBuilder) BeginComponent(f func(info *ModelInfo, obj interface{}, mode FieldModeStack, parentFormValueKey string, ctx *web.EventContext) h.HTMLComponent) *FieldsBuilder {
	b.beginComponentFuncs = append(b.beginComponentFuncs, f)
	return b
}

func (b *FieldsBuilder) HiddenField(f ...string) *FieldsBuilder {
	b.hiddenFields = append(b.hiddenFields, f...)
	unique.Sort(unique.StringSlice{&b.hiddenFields})
	for _, name := range f {
		b.Field(name).ComponentFunc(func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return h.Input("").Type("hidden").
				Attr(web.VField(field.FormKey, field.Value())...)
		})
	}
	return b
}

func (b *FieldsBuilder) SkipFieldVerifier() func(name string) bool {
	return b.skipFieldVerifier
}

func (b *FieldsBuilder) SetSkipFieldVerifier(skipFieldVerifier func(name string) bool) {
	b.skipFieldVerifier = skipFieldVerifier
}

type FieldsSection struct {
	Title string
	Rows  [][]string
}

func NewFieldsBuilder(builder *Builder) *FieldsBuilder {
	return &FieldsBuilder{
		builder: builder,
	}
}

func (b *FieldsBuilder) FieldNames() (r []any) {
	for _, field := range b.fields {
		r = append(r, field.name)
	}
	return
}

func (b *FieldsBuilder) Defaults(v *FieldDefaults) (r *FieldsBuilder) {
	b.defaults = v
	return b
}

func (b *FieldsBuilder) Unmarshal(opts *FieldsSetterOptions, toObj interface{}, info *ModelInfo, removeDeletedAndSort bool, ctx *web.EventContext) (vErr web.ValidationErrors) {
	t := reflect.TypeOf(toObj)
	if t.Kind() != reflect.Ptr {
		panic("toObj must be pointer")
	}

	if opts == nil {
		opts = &FieldsSetterOptions{}
	}

	fromObj := reflect.New(t.Elem()).Interface()
	// don't panic for fields that set in SetterFunc
	_ = ctx.UnmarshalForm(fromObj)

	if err := info.mb.BeforeFormUnmarshallHandlers.Handler(toObj, ctx); err != nil {
		vErr.GlobalError(err.Error())
		return
	}
	// testingutils.PrintlnJson("Unmarshal fromObj", fromObj)

	modifiedIndexes := ContextModifiedIndexesBuilder(ctx).FromHidden(ctx.R)

	vErr = b.SetObjectFields(opts, fromObj, toObj, &FieldContext{
		ToComponentOptions: &ToComponentOptions{
			SkipPermVerify: opts.SkipPermVerify,
		},
		EventContext: ctx,
		Obj:          fromObj,
		ModelInfo:    info,
	}, removeDeletedAndSort, modifiedIndexes, ctx)

	if err := info.mb.PostFormUnmarshallHandlers.Handler(toObj, ctx); err != nil {
		vErr.GlobalError(err.Error())
	}
	return
}

func (b *FieldsBuilder) DoSkipFieldVerifier(field string) bool {
	if b.skipFieldVerifier != nil {
		return b.skipFieldVerifier(field)
	}
	return false
}

func (b *FieldsBuilder) IsAllowed(r *http.Request, info *ModelInfo, obj interface{}, field string, perm ...string) bool {
	if !b.DoSkipFieldVerifier(field) {
		v := info.Permissioner().ReqObjectFielder(r, obj, field)
		for _, perm := range perm {
			if !v.Do(perm).Allowed() {
				return false
			}
		}
	}
	return true
}

type FieldsSetterOptions struct {
	SkipPermVerify bool
}

func (b *FieldsBuilder) SetObjectFields(opts *FieldsSetterOptions, fromObj interface{}, toObj interface{}, parent *FieldContext, removeDeletedAndSort bool, modifiedIndexes *ModifiedIndexesBuilder, ctx *web.EventContext) (vErr web.ValidationErrors) {
	if err := b.BeforeSetObjectFieldsHandler.Handler(fromObj, toObj, parent); err != nil {
		vErr.FieldError(parent.FormKey, err.Error())
		return
	}

	for _, f := range b.fields {
		if f.disabled {
			continue
		}

		info := parent.ModelInfo
		if !opts.SkipPermVerify && info != nil {
			if !b.IsAllowed(ctx.R, info, toObj, f.name, PermCreate, PermUpdate) {
				continue
			}
		}

		if f.nested != nil {
			switch f.rt.Kind() {
			case reflect.Slice:
				vErr.Merge(b.setWithChildFromObjs(opts, fromObj, parent, f, info, modifiedIndexes, toObj, removeDeletedAndSort, ctx))
				continue
			default:
				rt := reflectutils.GetType(toObj, f.name)
				childFromObj := reflectutils.MustGet(fromObj, f.name)
				if childFromObj == nil {
					childFromObj = reflect.New(rt.Elem()).Interface()
				}
				childToObj := reflectutils.MustGet(toObj, f.name)
				if childToObj == nil {
					childToObj = reflect.New(rt.Elem()).Interface()
				}
				if rt.Kind() == reflect.Struct {
					prv := reflect.New(rt)
					prv.Elem().Set(reflect.ValueOf(childToObj))
					childToObj = prv.Interface()
				}
				childInfo := f.nested.Model().Info().ChildOf(info, fromObj)
				pf := f.NewContext(childInfo, ctx, parent, childFromObj)

				vErr.Merge(f.nested.FieldsBuilder().SetObjectFields(opts, childFromObj, childToObj, pf, removeDeletedAndSort, modifiedIndexes, ctx))
				if err := reflectutils.Set(toObj, f.name, childToObj); err != nil {
					panic(err)
				}
				continue
			}
		} else if !opts.SkipPermVerify && info != nil {
			if !b.IsAllowed(ctx.R, info, toObj, f.name, PermCreate, PermUpdate) {
				continue
			}
		}

		fctx := f.NewContext(info, ctx, parent, toObj)

		val, err1 := reflectutils.Get(fromObj, f.name)
		if err1 != nil {
			switch err1.Error() {
			case "no such field",
				"reflect: call of reflect.Value.Interface on zero Value",
				"reflect.Value.Interface: cannot return value obtained from unexported field or method":
				goto set
			default:
			}
		} else {
			fctx.ValueOverride = val
		}

		if err1 = f.ValueFormatters.FormatValue(fctx); err1 != nil {
			vErr.FieldError(f.name, err1.Error())
			continue
		}

		reflectutils.Set(toObj, f.name, fctx.ValueOverride)

	set:
		if f.setterFunc == nil {
			continue
		}

		f.Setup.Setup(fctx)

		err1 = f.setterFunc(toObj, fctx, ctx)
		if err1 != nil {
			vErr.FieldError(fctx.FormKey, err1.Error())
		}
	}

	if vErr.HaveErrors() {
		return
	}

	if err := b.BeforeSetObjectFieldsHandler.Handler(fromObj, toObj, parent); err != nil {
		vErr.FieldError(parent.FormKey, err.Error())
	}
	return
}

func (b *FieldsBuilder) setToObjNilOrDelete(toObj interface{}, formKey string, f *FieldBuilder, modifiedIndexes *ModifiedIndexesBuilder, removeDeletedAndSort bool) {
	if !removeDeletedAndSort {
		if modifiedIndexes.deletedValues != nil && modifiedIndexes.deletedValues[formKey] != nil {
			for _, idx := range modifiedIndexes.deletedValues[formKey] {
				sliceFieldName := fmt.Sprintf("%s[%s]", f.name, idx)
				err := reflectutils.Set(toObj, sliceFieldName, nil)
				if err != nil {
					panic(err)
				}
			}
		}
		return
	}

	childToObjs := reflectutils.MustGet(toObj, f.name)
	if childToObjs == nil {
		return
	}

	t := reflectutils.GetType(toObj, f.name)
	newSlice := reflect.MakeSlice(t, 0, 0)
	modifiedIndexes.SortedForEach(childToObjs, formKey, func(obj interface{}, i int) {
		// remove deleted
		if modifiedIndexes.DeletedContains(formKey, i) {
			return
		}
		newSlice = reflect.Append(newSlice, reflect.ValueOf(obj))
	})

	err := reflectutils.Set(toObj, f.name, newSlice.Interface())
	if err != nil {
		panic(err)
	}

	return
}

func (b *FieldsBuilder) setWithChildFromObjs(
	opts *FieldsSetterOptions,
	fromObj interface{},
	parent *FieldContext,
	f *FieldBuilder,
	info *ModelInfo,
	modifiedIndexes *ModifiedIndexesBuilder,
	toObj interface{},
	removeDeletedAndSort bool,
	ctx *web.EventContext,
) (verr web.ValidationErrors) {
	childFromObjs := reflectutils.MustGet(fromObj, f.name)
	if childFromObjs == nil || reflect.TypeOf(childFromObjs).Kind() != reflect.Slice {
		return
	}

	fieldFormKey := parent.ChildFieldFormKey(f.name)

	i := 0
	reflectutils.ForEach(childFromObjs, func(childFromObj interface{}) {
		defer func() { i++ }()
		if childFromObj == nil {
			return
		}

		sliceFieldName := fmt.Sprintf("%s[%d]", f.name, i)

		childToObj := reflectutils.MustGet(toObj, sliceFieldName)
		if childToObj == nil {
			arrayElementType := reflectutils.GetType(toObj, sliceFieldName)

			if arrayElementType.Kind() == reflect.Ptr {
				arrayElementType = arrayElementType.Elem()
			} else {
				panic(fmt.Sprintf("%s (%T) must be a pointer", sliceFieldName, arrayElementType))
			}

			err := reflectutils.Set(toObj, sliceFieldName, reflect.New(arrayElementType).Interface())
			if err != nil {
				panic(err)
			}
			childToObj = reflectutils.MustGet(toObj, sliceFieldName)
		}
		// if is deleted, do nothing, later, it will be set to nil
		if modifiedIndexes.DeletedContains(fieldFormKey, i) {
			err := reflectutils.Set(toObj, sliceFieldName, nil)
			if err != nil {
				panic(err)
			}
			return
		}

		pf := f.NewContextBuilder(info, ctx, parent, fromObj).Index(i).Build()

		// fmt.Printf("childFromObj %#+v\n", childFromObj)
		// fmt.Printf("childToObj %#+v\n", childToObj)
		// fmt.Printf("fieldContext %#+v\n", pf)
		verr2 := f.nested.FieldsBuilder().SetObjectFields(opts, childFromObj, childToObj, pf, removeDeletedAndSort, modifiedIndexes, ctx)
		if verr2.HaveErrors() {
			verr.Merge(verr2)
		}
	})
	return
}

func (b FieldsBuilder) Clone() *FieldsBuilder {
	b.fields = append([]*FieldBuilder{}, b.fields...)
	return &b
}

func (b *FieldsBuilder) Model(v interface{}) (r *FieldsBuilder) {
	b.model = v
	return b
}

func (b *FieldsBuilder) GetModel() interface{} {
	return b.model
}

func (b *FieldsBuilder) Field(name string) (r *FieldBuilder) {
	r = b.GetField(name)
	if r != nil {
		return
	}

	r = b.appendNewFieldWithName(name)
	return
}

func (b *FieldsBuilder) CloneField(name ...string) (r *FieldsBuilder) {
	for _, name := range name {
		for i, f := range b.fields {
			if f.name == name {
				b.fields[i] = f.Clone()
			}
		}
	}
	return b
}

func (b *FieldsBuilder) FieldLabelsFromContextFunc(f FieldLabelsFunc) (r *FieldsBuilder) {
	b.fieldLabelsFromContextFunc = f
	return b
}

func (b *FieldsBuilder) WrapFieldLabelsFromContextFunc(f func(old FieldLabelsFunc) FieldLabelsFunc) (r *FieldsBuilder) {
	if b.fieldLabelsFromContextFunc != nil {
		b.fieldLabelsFromContextFunc = f(b.fieldLabelsFromContextFunc)
	} else {
		b.fieldLabelsFromContextFunc = f(func(b *FieldsBuilder, ctx web.ContextValuer) map[string]string {
			return nil
		})
	}
	return b
}

func (b *FieldsBuilder) GetFieldLabelsFromContext(ctx web.ContextValuer) (labels map[string]string) {
	if labels = GetFieldLabels(ctx, b); labels == nil {
		if b.fieldLabelsFromContextFunc != nil {
			labels = b.fieldLabelsFromContextFunc(b, ctx)
		}
		WithFieldLabels(ctx, b, labels)
	}
	return
}

func (b *FieldsBuilder) GetFieldOrDefault(name string) (r *FieldBuilder) {
	if r = b.GetField(name); r == nil {
		r = b.NewFieldWithName(name)
	}
	return
}

func (b *FieldsBuilder) GetField(name string) (r *FieldBuilder) {
	for _, f := range b.fields {
		if f.name == name {
			return f
		}
	}
	return
}

func (b *FieldsBuilder) HasField(name string) bool {
	for _, f := range b.fields {
		if f.name == name {
			return true
		}
	}
	return false
}

func (b *FieldsBuilder) getFieldNamesFromLayout() []string {
	var ns []string
	for _, iv := range b.fieldsLayout {
		switch t := iv.(type) {
		case string:
			ns = append(ns, t)
		case []string:
			for _, n := range t {
				ns = append(ns, n)
			}
		case *FieldsSection:
			for _, row := range t.Rows {
				for _, n := range row {
					ns = append(ns, n)
				}
			}
		default:
			panic("unknown fields layout, must be string/[]string/*FieldsSection")
		}
	}
	return ns
}

func (b *FieldsBuilder) Prepend(names ...any) (r *FieldsBuilder) {
	return b.Only(append(names, b.fieldsLayout...)...)
}

func (b *FieldsBuilder) Append(names ...any) (r *FieldsBuilder) {
	return b.Only(append(b.fieldsLayout, names...)...)
}

func (b *FieldsBuilder) Only(vs ...interface{}) (r *FieldsBuilder) {
	if len(vs) == 0 {
		return b
	}

	r = b.Clone()

	r.fieldsLayout = vs
	var (
		newFields []*FieldBuilder
		exists    = make(map[string]any)
	)
	for _, fn := range r.getFieldNamesFromLayout() {
		if _, ok := exists[fn]; !ok {
			exists[fn] = nil
			field := b.GetField(fn)
			if field != nil {
				newFields = append(newFields, field)
			}
		}
	}

	r.fields = newFields
	return
}

func (b *FieldsBuilder) appendNewFieldWithName(name string) (r *FieldBuilder) {
	r = b.NewFieldWithName(name)
	if r != nil {
		b.fields = append(b.fields, r)
	}
	return
}

func (b *FieldsBuilder) NewFieldWithName(name string) (r *FieldBuilder) {
	r = &FieldBuilder{
		mode: ALL,
	}

	r.name = name

	if b.model == nil {
		panic("model must be provided")
	}

	fType := reflectutils.GetType(b.model, name)
	if fType == nil {
		fType = reflect.TypeOf("")
	}
	r.rt = fType

	// if b.defaults == nil {
	// 	panic("field defaults must be provided")
	// }

	// ft := b.defaults.fieldTypeByTypeOrCreate(fType)
	// r.ComponentFunc(ft.compFunc).
	// 	SetterFunc(ft.setterFunc)
	return
}

func (b *FieldsBuilder) appendFieldAfterClone(ob *FieldsBuilder, name string) {
	f := ob.GetField(name)
	if f == nil {
		b.appendNewFieldWithName(name)
	} else {
		b.fields = append(b.fields, f.Clone())
	}
}

func (b *FieldsBuilder) Except(patterns ...string) (r *FieldsBuilder) {
	if len(patterns) == 0 {
		return b
	}

	r = b.Clone()
	r.fields = nil

	for _, f := range b.fields {
		if hasMatched(patterns, f.name) {
			continue
		}
		r.appendFieldAfterClone(b, f.name)
	}
	return
}

func (b *FieldsBuilder) String() (r string) {
	var names []string
	for _, f := range b.fields {
		names = append(names, f.name)
	}
	return fmt.Sprint(names)
}

func (b *FieldsBuilder) CurrentLayout() (layout []interface{}) {
	if b.fieldsLayout == nil {
		layout = make([]interface{}, 0, len(b.fields))
		for _, f := range b.fields {
			layout = append(layout, f.name)
		}
	} else {
		layout = b.fieldsLayout[:]
		layoutFM := make(map[string]struct{})
		for _, fn := range b.getFieldNamesFromLayout() {
			layoutFM[fn] = struct{}{}
		}
		for _, f := range b.fields {
			if _, ok := layoutFM[f.name]; ok {
				continue
			}
			layout = append(layout, f.name)
		}
	}
	return
}

type ToComponentOptions struct {
	FieldVerifier  *perm.Verifier
	SkipPermVerify bool
}

func (b *FieldsBuilder) ToComponent(opts *ToComponentOptions, info *ModelInfo, obj interface{}, mode FieldModeStack, ctx *web.EventContext) h.HTMLComponent {
	return b.toComponentWithModifiedIndexes(opts, info, obj, mode, nil, ctx)
}

func (b *FieldsBuilder) toComponentWithModifiedIndexes(opts *ToComponentOptions, info *ModelInfo, obj interface{}, mode FieldModeStack, parent *FieldContext, ctx *web.EventContext) h.HTMLComponent {
	modifiedIndexes := ContextModifiedIndexesBuilder(ctx)
	return b.toComponentWithFormValueKey(opts, info, obj, mode, parent, modifiedIndexes, ctx)
}

func (b *FieldsBuilder) toComponentWithFormValueKey(opts *ToComponentOptions, info *ModelInfo, obj interface{}, mode FieldModeStack, parent *FieldContext, modifiedIndexes *ModifiedIndexesBuilder, ctx *web.EventContext) h.HTMLComponent {
	var (
		comps     []h.HTMLComponent
		okNames   = make(map[string]any)
		parentKey string
	)

	if parent == nil {
		comps = append(comps, modifiedIndexes.ToFormHidden())
	} else {
		parentKey = parent.FormKey
	}

	vErr, _ := ctx.Flash.(*web.ValidationErrors)
	if vErr == nil {
		vErr = &web.ValidationErrors{}
	}

	// changes mode if not is embedded
	if s := info.Schema(); s != nil && model.HasPrimaryFields(s) {
		if !mode.Dot().Is(LIST, DETAIL) {
			if id, _, _ := info.LookupID(obj); id.IsZero() {
				mode = append(mode, NEW)
			}
		}
	} else {
		mode = append(mode, EDIT)
	}

	for _, f := range b.beginComponentFuncs {
		comps = append(comps, f(info, obj, mode, parentKey, ctx))
	}

	for _, name := range b.hiddenFields {
		fComp := b.fieldToComponentWithFormValueKey(opts, info, obj, mode, parent, ctx, name, vErr)
		if fComp != nil {
			comps = append(comps, fComp)
		}
	}

	layout := b.CurrentLayout()

	for _, iv := range layout {
		var comp h.HTMLComponent
		switch t := iv.(type) {
		case string:
			if _, ok := okNames[t]; ok {
				continue
			}
			okNames[t] = nil

			comp = b.fieldToComponentWithFormValueKey(opts, info, obj, mode, parent, ctx, t, vErr)
		case []string:
			colsComp := make([]h.HTMLComponent, 0, len(t))
			for _, n := range t {
				if _, ok := okNames[n]; ok {
					continue
				}
				okNames[n] = nil
				fComp := b.fieldToComponentWithFormValueKey(opts, info, obj, mode, parent, ctx, n, vErr)
				if fComp == nil {
					continue
				}
				colsComp = append(colsComp, v.VCol(fComp).Class("pr-4"))
			}
			if len(colsComp) > 0 {
				comp = v.VRow(colsComp...).NoGutters(true)
			}
		case *FieldsSection:
			rowsComp := make([]h.HTMLComponent, 0, len(t.Rows))
			for _, row := range t.Rows {
				colsComp := make([]h.HTMLComponent, 0, len(row))
				for _, n := range row {
					if _, ok := okNames[n]; ok {
						continue
					}
					okNames[n] = nil

					fComp := b.fieldToComponentWithFormValueKey(opts, info, obj, mode, parent, ctx, n, vErr)
					if fComp == nil {
						continue
					}
					colsComp = append(colsComp, v.VCol(fComp).Class("pr-4"))
				}
				if len(colsComp) > 0 {
					rowsComp = append(rowsComp, v.VRow(colsComp...).NoGutters(true))
				}
			}
			if len(rowsComp) > 0 {
				var titleComp h.HTMLComponent
				if t.Title != "" {
					titleComp = h.Label(t.Title).Class("v-label theme--light text-caption")
				}
				comp = h.Div(
					titleComp,
					v.VCard(rowsComp...).Variant(v.VariantOutlined).Class("mx-0 mt-1 mb-4 px-4 pb-0 pt-4"),
				)
			}
		default:
			panic("unknown fields layout, must be string/[]string/*FieldsSection")
		}
		if comp == nil {
			continue
		}
		comps = append(comps, comp)
	}

	return h.Components(comps...)
}

func (b *FieldsBuilder) fieldToComponentWithFormValueKey(opts *ToComponentOptions, info *ModelInfo, obj interface{}, mode FieldModeStack, parent *FieldContext, ctx *web.EventContext, name string, vErr *web.ValidationErrors) h.HTMLComponent {
	f := b.GetFieldOrDefault(name)

	if f.disabled || (f.compFunc == nil && f.nested == nil) {
		return nil
	}

	fctx := f.NewContext(info, ctx, parent, obj)
	fctx.Mode = mode
	fctx.Errors = vErr.GetRemoveFieldErrors(fctx.FormKey)

	if !opts.SkipPermVerify && info != nil {
		var (
			fqn  = fctx.Path.NoIndex().Fqn()
			perm = PermGet
			mode = mode.Dot()
		)
		switch mode {
		case EDIT:
			perm = PermUpdate
		case NEW:
			perm = PermCreate
		}

		if !b.IsAllowed(ctx.R, info.Root(), fctx.Root().Obj, fqn, perm) {
			if mode == EDIT && b.IsAllowed(ctx.R, info.Root(), fctx.Root().Obj, fqn, PermGet) {
				fctx.ReadOnly = true
			} else {
				return nil
			}
		}
	}

	b.FieldToComponentSetup.Setup(fctx)
	f.Setup.Setup(fctx)
	f.ToComponentSetup.Setup(fctx)

	if fctx.Disabled || !f.IsEnabled(fctx) {
		return nil
	}

	return f.ToComponent(fctx)
}

type RowFunc func(obj interface{}, path FieldPath, formKey string, content h.HTMLComponent, ctx *web.EventContext) h.HTMLComponent

func defaultRowFunc(obj interface{}, path FieldPath, formKey string, content h.HTMLComponent, ctx *web.EventContext) h.HTMLComponent {
	return content
}

func (b *FieldsBuilder) ToComponentForEach(opts *ToComponentOptions, field *FieldContext, slice interface{}, mode FieldModeStack, ctx *web.EventContext, rowFunc RowFunc) h.HTMLComponent {
	var (
		info            *ModelInfo
		r               []h.HTMLComponent
		modifiedIndexes = ContextModifiedIndexesBuilder(ctx)
	)

	if field != nil {
		info = field.ModelInfo
	} else {
		field = &FieldContext{
			ToComponentOptions: opts,
			FormKey:            "",
		}
	}

	if rowFunc == nil {
		rowFunc = defaultRowFunc
	}

	modifiedIndexes.SortedForEach(slice, field.FormKey, func(obj interface{}, i int) {
		if modifiedIndexes.DeletedContains(field.FormKey, i) {
			return
		}

		parent := *field
		parent.FormKey = fmt.Sprintf("%s[%d]", parent.FormKey, i)
		parent.Path.AppendIndex(i)
		comps := b.toComponentWithFormValueKey(opts, info.ItemOf(slice, i), obj, mode, &parent, modifiedIndexes, ctx)
		r = append(r, rowFunc(obj, parent.Path, parent.FormKey, comps, ctx))
	})

	return h.Components(r...)
}

type ModifiedIndexesBuilder struct {
	deletedValues map[string]map[int]bool
	sortedValues  map[string][]string
}

type deletedIndexBuilderKeyType int

const theDeleteIndexBuilderKey deletedIndexBuilderKeyType = iota

const (
	deletedHiddenNamePrefix = "__Deleted"
	sortedHiddenNamePrefix  = "__Sorted"
)

func ContextModifiedIndexesBuilder(ctx *web.EventContext) (r *ModifiedIndexesBuilder) {
	r, ok := ctx.R.Context().Value(theDeleteIndexBuilderKey).(*ModifiedIndexesBuilder)
	if !ok {
		r = &ModifiedIndexesBuilder{
			deletedValues: make(map[string]map[int]bool),
		}
		ctx.R = ctx.R.WithContext(context.WithValue(ctx.R.Context(), theDeleteIndexBuilderKey, r))
	}
	return r
}

func (b *ModifiedIndexesBuilder) AppendDeleted(sliceFormKey string, index int) (r *ModifiedIndexesBuilder) {
	if v, ok := b.deletedValues[sliceFormKey]; ok {
		v[index] = true
	} else {
		b.deletedValues[sliceFormKey] = map[int]bool{index: true}
	}
	return b
}

func (b *ModifiedIndexesBuilder) SetSorted(sliceFormKey string, indexes []string) (r *ModifiedIndexesBuilder) {
	if b.sortedValues == nil {
		b.sortedValues = make(map[string][]string)
	}
	b.sortedValues[sliceFormKey] = indexes
	return b
}

func (b *ModifiedIndexesBuilder) DeletedContains(sliceFormKey string, index int) (r bool) {
	if b.deletedValues == nil || b.deletedValues[sliceFormKey] == nil {
		return false
	}
	return b.deletedValues[sliceFormKey][index]
}

func (b *ModifiedIndexesBuilder) SortedForEach(slice interface{}, sliceFormKey string, f func(obj interface{}, i int)) {
	sortedIndexes, ok := b.sortedValues[sliceFormKey]
	if !ok {
		sliceVal := reflect.ValueOf(slice)
		for i := 0; i < sliceVal.Len(); i++ {
			obj := sliceVal.Index(i).Interface()
			f(obj, i)
		}
		return
	}

	sliceLen := reflect.ValueOf(slice).Len()
	for j1 := 0; j1 < sliceLen; j1++ {
		if slices.Contains(sortedIndexes, fmt.Sprint(j1)) {
			continue
		}
		sortedIndexes = append(sortedIndexes, fmt.Sprint(j1))
	}

	for _, j := range sortedIndexes {
		obj, err := reflectutils.Get(slice, fmt.Sprintf("[%s]", j))
		if obj == nil || err != nil {
			continue
		}
		j1, _ := strconv.Atoi(j)
		f(obj, j1)
	}
}

func deleteHiddenSliceFormKey(sliceFormKey string) string {
	return deletedHiddenNamePrefix + "." + sliceFormKey
}

func sortedHiddenSliceFormKey(sliceFormKey string) string {
	return sortedHiddenNamePrefix + "." + sliceFormKey
}

func (b *ModifiedIndexesBuilder) FromHidden(req *http.Request) (r *ModifiedIndexesBuilder) {
	if b.deletedValues == nil {
		b.deletedValues = make(map[string]map[int]bool)
	}
	if b.sortedValues == nil {
		b.sortedValues = make(map[string][]string)
	}
	for k, vs := range req.Form {
		if strings.HasPrefix(k, deletedHiddenNamePrefix) {
			m := make(map[int]bool)
			for _, s := range strings.Split(vs[0], ",") {
				if i, err := strconv.Atoi(s); err == nil {
					m[i] = true
				}
			}
			b.deletedValues[k[len(deletedHiddenNamePrefix)+1:]] = m
		}

		if strings.HasPrefix(k, sortedHiddenNamePrefix) {
			b.sortedValues[k[len(sortedHiddenNamePrefix)+1:]] = strings.Split(vs[0], ",")
		}
	}
	return b
}

func (b *ModifiedIndexesBuilder) ToFormHidden() h.HTMLComponent {
	var hidden []h.HTMLComponent
	for sliceFormKey, values := range b.deletedValues {
		var keys []string
		for k := range values {
			keys = append(keys, fmt.Sprint(k))
		}
		hidden = append(hidden, h.Input("").Type("hidden").
			Attr(web.VField(deleteHiddenSliceFormKey(sliceFormKey), strings.Join(keys, ","))...))
	}

	for sliceFormKey, values := range b.sortedValues {
		hidden = append(hidden, h.Input("").Type("hidden").
			Attr(web.VField(sortedHiddenSliceFormKey(sliceFormKey), strings.Join(values, ","))...))
	}
	return h.Components(hidden...)
}
