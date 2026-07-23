package presets

import (
	"net/url"
	"reflect"
	"sort"
	"strconv"

	"github.com/go-rvq/rvq/web"
)

// ListEditorFormValues returns the submitted form values a list editor should
// read per-item metadata from: the multipart values (which ReorderSlicesByPos
// reindexes to match the decoded slice order) when present, otherwise the parsed
// form. Both the render and the persistence layer read from here so their
// per-item indexes stay aligned with the (reordered) decoded slice.
func ListEditorFormValues(ctx *web.EventContext) url.Values {
	if ctx != nil && ctx.R != nil {
		if ctx.R.MultipartForm != nil && ctx.R.MultipartForm.Value != nil {
			return ctx.R.MultipartForm.Value
		}
		if ctx.R.Form != nil {
			return ctx.R.Form
		}
	}
	return nil
}

// ListEditorItemsPartition groups the items of a submitted list-editor slice by
// their per-item metadata flags.
//
// Classification is by form flags ONLY (__deleted / __new) and never by primary
// key, so the list editor works for element types that have no ID (or whose ID
// is assigned before insert). The list editor emits __new for browser-created
// rows, which is the single source of truth for "this is a creation".
type ListEditorItemsPartition struct {
	// Deleted holds the items flagged __deleted (to be removed).
	Deleted []reflect.Value
	// New holds the non-deleted items flagged __new (to be created).
	New []reflect.Value
	// Others holds the non-deleted, non-new items (existing rows, to be updated).
	Others []reflect.Value
	// Kept holds every non-deleted item — New and Others — in __pos order.
	// It is the desired final set (e.g. for gorm's Association.Replace).
	Kept []reflect.Value
}

// KeptSlice returns Kept as a freshly allocated slice value of sliceType (the
// concrete slice type of the field), ready to hand to code that needs a typed
// slice instead of the []reflect.Value view.
func (p ListEditorItemsPartition) KeptSlice(sliceType reflect.Type) reflect.Value {
	out := reflect.MakeSlice(sliceType, 0, len(p.Kept))
	for _, v := range p.Kept {
		out = reflect.Append(out, v)
	}
	return out
}

// ListEditorInitialized reports whether the list editor for fieldFormKey was
// rendered/submitted, i.e. it seeded its presence marker (`<fieldFormKey>` +
// "." + ListEditorPresentField == "1"). When false the field was not part of the
// form and its children must be left untouched.
func ListEditorInitialized(values url.Values, fieldFormKey string) bool {
	return values != nil && values.Get(fieldFormKey+"."+ListEditorPresentField) == "1"
}

// PartitionListEditorItems classifies items — the slice reflect.Value of the
// list field fieldFormKey — using the submitted form values. See
// ListEditorItemsPartition for the classification rules. Every partition is
// returned in the user-intended __pos order.
//
// It returns nil when the list editor was not initialized for fieldFormKey (its
// presence marker is absent), so callers can distinguish "the field was not
// submitted — leave it alone" from "submitted, possibly empty".
func PartitionListEditorItems(values url.Values, fieldFormKey string, items reflect.Value) *ListEditorItemsPartition {
	if !ListEditorInitialized(values, fieldFormKey) {
		return nil
	}

	p := &ListEditorItemsPartition{}

	// order the item indexes by their submitted __pos (stable; a missing __pos
	// falls back to the natural index), so all partitions come out ordered.
	order := make([]int, items.Len())
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(a, b int) bool {
		return listEditorItemPos(values, fieldFormKey, order[a]) < listEditorItemPos(values, fieldFormKey, order[b])
	})

	for _, i := range order {
		item := items.Index(i)
		// stable __index form keys are sparse, so the decoded slice can carry nil
		// (or zero) holes where indexes are unused — skip them.
		if item.Kind() == reflect.Ptr && item.IsNil() {
			continue
		}
		if ListEditorItemFlag(values, fieldFormKey, i, ListEditorDeletedField) {
			p.Deleted = append(p.Deleted, item)
			continue
		}
		p.Kept = append(p.Kept, item)
		if ListEditorItemFlag(values, fieldFormKey, i, ListEditorNewField) {
			p.New = append(p.New, item)
		} else {
			p.Others = append(p.Others, item)
		}
	}
	return p
}

// listEditorItemPos returns item i's submitted __pos, or i when it is absent or
// unparseable, so ordering is stable for older/partial forms.
func listEditorItemPos(values url.Values, fieldFormKey string, i int) int {
	if values != nil {
		if v := values.Get(fieldFormKey + "[" + strconv.Itoa(i) + "]." + ListEditorPositionField); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				return n
			}
		}
	}
	return i
}

// ListEditorItemFlag reports whether the boolean metadata field (e.g. __deleted,
// __new) is truthy for item i of the list fieldFormKey in the submitted values.
func ListEditorItemFlag(values url.Values, fieldFormKey string, i int, field string) bool {
	if values == nil {
		return false
	}
	switch values.Get(fieldFormKey + "[" + strconv.Itoa(i) + "]." + field) {
	case "true", "1", "on":
		return true
	}
	return false
}
