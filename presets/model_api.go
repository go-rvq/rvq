package presets

import (
	"github.com/qor5/web/v3"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
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
)

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
