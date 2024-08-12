package presets

import (
	"path"

	"github.com/qor5/web/v3"
)

type OkHandled func(ctx *web.EventContext) (ok bool, handled bool)

func DefaultOkHandled(ctx *web.EventContext) (ok, handled bool) {
	return
}

func WrapOkHandled(old OkHandled, do func(h OkHandled) OkHandled) OkHandled {
	if old == nil {
		old = DefaultOkHandled
	}
	return do(old)
}

func CallOkHandled(h OkHandled, ctx *web.EventContext) (ok bool) {
	if h != nil {
		ok, _ = h(ctx)
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
