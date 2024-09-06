package presets

import (
	"fmt"

	"github.com/qor5/web/v3"
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

func (n *NestedSliceBuilder) Walk(fctx *FieldContext, handle FieldWalkHandle) (s FieldWalkState) {
	var (
		i     int
		slice = fctx.RawValue()
	)
	if slice != nil {
		reflectutils.ForEach(slice, func(v interface{}) {
			defer func() { i++ }()
			fieldInfo := n.mb.Info().ChildOf(fctx.ModelInfo, fctx.Obj).ItemOf(slice, i)
			n.fb.walk(fieldInfo, v, fctx.Mode, fmt.Sprintf("%s[%d]", fctx.FormKey, i), fctx.EventContext, handle)
		})
	}
	return
}
