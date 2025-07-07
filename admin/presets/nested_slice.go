package presets

import (
	"fmt"
	"reflect"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/zeroer"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
)

type NestedSliceBuilder struct {
	mb                     *ModelBuilder
	fb                     *FieldsBuilder
	displayFieldInSorter   string
	addListItemRowEvent    string
	removeListItemRowEvent string
	sortListItemsEvent     string
}

func NestedSlice(mb *ModelBuilder, fb *FieldsBuilder) *NestedSliceBuilder {
	return &NestedSliceBuilder{mb: mb, fb: fb}
}

func (n *NestedSliceBuilder) Model() *ModelBuilder {
	return n.mb
}

func (n *NestedSliceBuilder) FieldsBuilder() *FieldsBuilder {
	return n.fb
}

func (n *NestedSliceBuilder) DisplayFieldInSorter() string {
	return n.displayFieldInSorter
}

func (n *NestedSliceBuilder) SetDisplayFieldInSorter(displayFieldInSorter string) *NestedSliceBuilder {
	n.displayFieldInSorter = displayFieldInSorter
	return n
}

func (n *NestedSliceBuilder) AddListItemRowEvent() string {
	return n.addListItemRowEvent
}

func (n *NestedSliceBuilder) SetAddListItemRowEvent(addListItemRowEvent string) *NestedSliceBuilder {
	n.addListItemRowEvent = addListItemRowEvent
	return n
}

func (n *NestedSliceBuilder) RemoveListItemRowEvent() string {
	return n.removeListItemRowEvent
}

func (n *NestedSliceBuilder) SetRemoveListItemRowEvent(removeListItemRowEvent string) *NestedSliceBuilder {
	n.removeListItemRowEvent = removeListItemRowEvent
	return n
}

func (n *NestedSliceBuilder) SortListItemsEvent() string {
	return n.sortListItemsEvent
}

func (n *NestedSliceBuilder) SetSortListItemsEvent(sortListItemsEvent string) *NestedSliceBuilder {
	n.sortListItemsEvent = sortListItemsEvent
	return n
}

func (n *NestedSliceBuilder) Build(b *FieldBuilder) {
	b.ComponentFunc(func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return NewListEditor(field).Value(field.Value()).
			DisplayFieldInSorter(n.displayFieldInSorter).
			AddListItemRowEvent(n.addListItemRowEvent).
			RemoveListItemRowEvent(n.removeListItemRowEvent).
			SortListItemsEvent(n.sortListItemsEvent).Component(ctx)
	})
}

func (n *NestedSliceBuilder) Walk(fctx *FieldContext, opts *FieldWalkHandleOptions) (s FieldWalkState) {

	if fctx.Obj == nil {
		if opts.SkipNestedNil {
			return
		}
	}

	var (
		i     int
		slice = fctx.RawValue()
	)

	if slice == nil {
		if opts.InitializeSlices {
			st := reflectutils.GetType(fctx.Obj, fctx.Name)
			s := reflect.MakeSlice(st, 1, 1)
			s.Index(0).Set(reflect.ValueOf(n.mb.Model()))
			slice = s.Interface()
		} else {
			return
		}
	}

	reflectutils.ForEach(slice, func(v interface{}) {
		defer func() { i++ }()
		if opts.SkipNestedNil && zeroer.IsNil(v) {
			return
		}
		fieldInfo := n.mb.Info().ChildOf(fctx.ModelInfo, fctx.Obj).ItemOf(slice, i)
		n.fb.walk(fieldInfo, v, fctx.Mode, append(fctx.Path, FieldPathIndex(i)), fmt.Sprintf("%s[%d]", fctx.FormKey, i), fctx.EventContext, opts)
	})
	return
}
