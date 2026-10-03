package publish

import (
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

// notDeleted refuses a deleted record (presets.ModelBuilder.IsDeleted): in
// the trash, it is not published, unpublished, scheduled nor renamed.
func notDeleted(mb *presets.ModelBuilder, obj any) error {
	if mb.IsDeleted(obj) {
		return presets.ErrUpdateRecordNotAllowed
	}
	return nil
}

func wrapEventFuncWithShowError(f web.EventFunc) web.EventFunc {
	return func(ctx *web.EventContext) (web.EventResponse, error) {
		r, err := f(ctx)
		if err != nil {
			presets.ShowMessage(&r, err.Error(), "error")
		}
		return r, nil
	}
}
