package presets

import (
	"reflect"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

type (
	PageTitle interface {
		PageTitle() string
	}

	ContextTitler interface {
		ContextTitle(ctx *web.EventContext) string
	}

	ModelFormUnmarshallHandler interface {
		Handler(obj interface{}, ctx *web.EventContext) (err error)
	}

	ModelFormUnmarshallHandlerFunc func(obj interface{}, ctx *web.EventContext) (err error)
	ModelFormUnmarshallHandlers    []ModelFormUnmarshallHandler

	RecordMenuItemContext struct {
		Ctx             *web.EventContext
		RecordIndex     int
		Obj             any
		ID              string
		AddRecordPortal func(name string) (fqn string)
		TempPortal      string
	}

	RecordMenuItemFunc  func(rctx *RecordMenuItemContext) h.HTMLComponent
	RecordMenuItemFuncs []RecordMenuItemFunc

	RecordEncoderFactory[T any] func(ctx *web.EventContext) func(r T) any

	ParentUriPart int
)

func (f RecordEncoderFactory[T]) Any() RecordEncoderFactory[any] {
	return func(ctx *web.EventContext) func(r any) any {
		do := f(ctx)
		return func(r any) any {
			return do(r.(T))
		}
	}
}

func (f RecordEncoderFactory[T]) EncodeSlice(ctx *web.EventContext, slice any) any {
	sliceValue := reflect.ValueOf(slice)

	if sliceValue.Kind() != reflect.Slice {
		sliceValue = reflect.ValueOf([]any{slice})
	}

	newSlice := reflect.MakeSlice(reflect.TypeOf([]any{}), sliceValue.Len(), sliceValue.Len())
	encode := f(ctx)

	for i := 0; i < sliceValue.Len(); i++ {
		newSlice.Index(i).Set(reflect.ValueOf(encode(sliceValue.Index(i).Interface().(T))))
	}

	return newSlice.Interface()
}

func (f ModelFormUnmarshallHandlerFunc) Handler(obj interface{}, ctx *web.EventContext) (err error) {
	return f(obj, ctx)
}

func (h ModelFormUnmarshallHandlers) Handler(obj interface{}, ctx *web.EventContext) (err error) {
	for _, handler := range h {
		if err = handler.Handler(obj, ctx); err != nil {
			return
		}
	}
	return
}

func (h *ModelFormUnmarshallHandlers) Append(handler ...ModelFormUnmarshallHandler) {
	*h = append(*h, handler...)
}

func (h *ModelFormUnmarshallHandlers) AppendFunc(handler ...ModelFormUnmarshallHandlerFunc) {
	for _, f := range handler {
		h.Append(f)
	}
}

func (funcs RecordMenuItemFuncs) ToRowMenuItemFuncs(sharedPortal string, addRowPortal func(rctx *RecordMenuItemContext, name string) string) (r []vx.RowMenuItemFunc) {
	r = make([]vx.RowMenuItemFunc, len(funcs))
	funcs.ForEachRowMenuItemFunc(sharedPortal, addRowPortal, func(i int, itemFunc vx.RowMenuItemFunc) {
		r[i] = itemFunc
	})
	return
}

func (funcs RecordMenuItemFuncs) ForEachRowMenuItemFunc(sharedPortal string, addRowPortal func(rctx *RecordMenuItemContext, name string) string, cb func(i int, f vx.RowMenuItemFunc)) {
	for i, f := range funcs {
		cb(i, func(recordIndex int, obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent {
			rctx := &RecordMenuItemContext{
				RecordIndex: recordIndex,
				Ctx:         ctx,
				Obj:         obj,
				ID:          id,
				TempPortal:  sharedPortal,
			}
			rctx.AddRecordPortal = func(name string) (fqn string) {
				return addRowPortal(rctx, name)
			}
			return f(rctx)
		})
	}
}
