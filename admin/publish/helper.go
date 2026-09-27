package publish

import (
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

func wrapEventFuncWithShowError(f web.EventFunc) web.EventFunc {
	return func(ctx *web.EventContext) (web.EventResponse, error) {
		r, err := f(ctx)
		if err != nil {
			presets.ShowMessage(&r, err.Error(), "error")
		}
		return r, nil
	}
}
