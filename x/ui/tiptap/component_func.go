package tiptap

import (
	"context"
	"reflect"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/media"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/fields/tiptap"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/admin/reflect_utils"
	"github.com/go-rvq/rvq/utils"
	"github.com/go-rvq/rvq/utils/context_utils"
	"github.com/go-rvq/rvq/web"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

type contextKey string

const (
	tipTapContextKey contextKey = "tipTap"
	MediaBuilderKey  contextKey = "tiptap:mediaBuilder"

	dbKey contextKey = "tiptap:db"
)

func ComponentFunc(mb *presets.ModelBuilder, mode presets.FieldMode, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	t := context_utils.OrContextValue(field.Field.ContextPtr(), tipTapContextKey, func() (t *tiptap.Builder) {
		fieldIndex := reflect_utils.AllFieldsOf(mb.Model()).Get(field.Name).Index
		mb.WithDataOperator(func(do presets.DataOperator) {
			do.(*gorm2op.DataOperatorBuilder).
				WithWriteCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
					cb.Pre(func(state *gorm2op.CallbackState) (err error) {
						rv := reflect.ValueOf(state.Obj).Elem()
						rField := rv.FieldByIndex(fieldIndex)

						v := rField.String()
						if len(v) > 0 {
							var r *tiptap.StoreImagesResult
							if r, err = t.StoreImages(
								context.WithValue(context.Background(), dbKey, state.SharedDB),
								utils.UnsafeStringToBytes(v)); err != nil {
								return
							}
							if r.Changed() {
								v = string(r.NewValue)
								rField.SetString(v)
							}
						}
						return
					})
				})
		})

		mediaBuilder, _ := field.Field.GetData(MediaBuilderKey).(*media.Builder)
		if mediaBuilder == nil {
			mediaBuilder, _ = mb.GetData(MediaBuilderKey).(*media.Builder)
			if mediaBuilder == nil {
				mediaBuilder, _ = mb.Builder().GetData(MediaBuilderKey).(*media.Builder)
			}
		}

		t = tiptap.New().
			Model(mb)

		if mediaBuilder != nil {
			t.Store(tiptap.MediaLibraryStorer(mediaBuilder, dbKey))
		}

		t.Fields(field.Name).
			WrapEditor(func(ctx *presets.FieldContext, comp *vx.VXTipTapEditorBuilder) {
				comp.Template(true)
			}).
			Build(mode)
		return
	})

	return t.ComponentFunc(field.Mode.Dot())(field, ctx)
}
