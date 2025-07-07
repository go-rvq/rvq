package presets

import (
	"path"

	"github.com/qor5/web/v3"
)

type OkHandler interface {
	Handle(ctx *web.EventContext) (ok bool, handled bool)
}

type OkHandlerFunc func(ctx *web.EventContext) (ok, handled bool)

func (f OkHandlerFunc) Handle(ctx *web.EventContext) (ok, handled bool) {
	return f(ctx)
}

type OkObjHandler interface {
	Handle(obj any, ctx *web.EventContext) (ok, handled bool)
}

type OkObjHandlerFunc func(obj any, ctx *web.EventContext) (ok, handled bool)

func (f OkObjHandlerFunc) Handle(obj any, ctx *web.EventContext) (ok, handled bool) {
	return f(obj, ctx)
}

func OkObjHandlerFuncT[T any](f func(obj T, ctx *web.EventContext) (ok, handled bool)) OkObjHandlerFunc {
	return func(obj any, ctx *web.EventContext) (ok, handled bool) {
		return f(obj.(T), ctx)
	}
}

type OkHandlers []OkHandler

func (s OkHandlers) Handle(ctx *web.EventContext) (ok, handled bool) {
	var handled_ bool
	for _, h := range s {
		if h != nil {
			if ok, handled_ = h.Handle(ctx); ok {
				handled = true
				return
			} else if handled_ {
				handled = true
			}
		}
	}
	return
}

type OkObjHandlers []OkObjHandler

func (s OkObjHandlers) Handle(obj any, ctx *web.EventContext) (ok, handled bool) {
	var handled_ bool
	for _, h := range s {
		if h != nil {
			if ok, handled_ = h.Handle(obj, ctx); ok {
				handled = true
				return
			} else if handled_ {
				handled = true
			}
		}
	}
	return
}

var DefaultOkHandler OkHandlerFunc = func(ctx *web.EventContext) (ok, handled bool) {
	return
}

var DefaultOkObjHandler OkObjHandlerFunc = func(obj any, ctx *web.EventContext) (ok, handled bool) {
	return
}

func WrapOkHandled(old OkHandler, do func(h OkHandler) OkHandler) OkHandler {
	if old == nil {
		old = DefaultOkHandler
	}
	return do(old)
}

func WrapOkObjHandled(old OkObjHandler, do func(h OkObjHandler) OkObjHandler) OkObjHandler {
	if old == nil {
		old = DefaultOkObjHandler
	}
	return do(old)
}

func CallOkHandler(h OkHandler, ctx *web.EventContext) (ok bool) {
	if h != nil {
		ok, _ = h.Handle(ctx)
	}
	return
}

func CallOkObjHandlers(h OkObjHandler, obj any, ctx *web.EventContext) (ok bool) {
	if h != nil {
		ok, _ = h.Handle(obj, ctx)
	}
	return
}

const PortalPathKey = "portal_path"

func AddPortalPath(ctx *web.EventContext, name string) (rollback func()) {
	var s []string

	if v := ctx.ContextValue(PortalPathKey); v != nil {
		s = v.([]string)
	}

	ctx.WithContextValue(PortalPathKey, append(s, name))
	return func() {
		ctx.WithContextValue(PortalPathKey, s)
	}
}

func GetPortalPath(ctx *web.EventContext, sub string) string {
	if v := ctx.ContextValue(PortalPathKey); v != nil {
		return path.Join(append(v.([]string), sub)...)
	}
	return sub
}
