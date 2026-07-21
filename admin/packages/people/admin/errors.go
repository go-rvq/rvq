package admin

import (
	"context"
	"errors"

	"github.com/go-rvq/rvq/admin/packages/people/messages"
	"github.com/go-rvq/rvq/admin/packages/people/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

// errText maps each model validation sentinel to its translated message.
func errText(m *messages.Messages, err error) string {
	switch {
	case errors.Is(err, models.ErrInvalidDocumentType):
		return m.ErrInvalidDocumentType
	}
	return ""
}

// translateErr replaces a known model validation error with its translated
// message for the request language, leaving other errors untouched.
func translateErr(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if s := errText(messages.Get(ctx), err); s != "" {
		return errors.New(s)
	}
	return err
}

// wrapSaveErrors translates the model validation errors of a resource's save
// path to the request language, so the admin shows them localized.
func wrapSaveErrors(mb *presets.ModelBuilder) {
	ed := mb.Editing()
	ed.WrapSaveFunc(func(in presets.SaveFunc) presets.SaveFunc {
		return func(obj interface{}, id presets.ID, ctx *web.EventContext) error {
			return translateErr(ctx.Context(), in(obj, id, ctx))
		}
	})
	if ed.HasCreatingBuilder() {
		ed.CreatingBuilder().WrapSaveFunc(func(in presets.SaveFunc) presets.SaveFunc {
			return func(obj interface{}, id presets.ID, ctx *web.EventContext) error {
				return translateErr(ctx.Context(), in(obj, id, ctx))
			}
		})
	}
}
