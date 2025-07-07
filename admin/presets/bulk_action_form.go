package presets

import (
	"reflect"

	"github.com/qor5/web/v3"
	"github.com/qor5/web/v3/datafield"
	h "github.com/theplant/htmlgo"
)

type BulkActionFormBuilderHandler[T any] func(ctx *BulkActionFormContext[T]) (err error)

type BulkActionFormContext[T any] struct {
	Context    *web.EventContext
	Response   *web.EventResponse
	SelectedID []string
	Form       T
	Err        error
	datafield.DataField[*BulkActionFormContext[T]]
}

func NewBulkActionFormContext[T any](context *web.EventContext, r *web.EventResponse, selectedID []string) *BulkActionFormContext[T] {
	var model T
	return datafield.New(&BulkActionFormContext[T]{
		Context:    context,
		Response:   r,
		SelectedID: selectedID,
		Form:       reflect.New(reflect.TypeOf(model).Elem()).Interface().(T),
	})
}

type BulkActionFormBuilder[T any] struct {
	model            T
	action           *BulkActionBuilder
	eb               *EditingBuilder
	handlers         BulkActionFormBuilderHandlers[T]
	insertFormBefore bool
}

type BulkActionFormBuilderHandlers[T any] []BulkActionFormBuilderHandler[T]

func (c BulkActionFormBuilderHandlers[T]) Handler() BulkActionUpdateFunc {
	return func(selectedID []string, ctx *web.EventContext, r *web.EventResponse) (err error) {
		actx := NewBulkActionFormContext[T](ctx, r, selectedID)
		ctx.WithContextValue(ctxBulkActionFormContext, actx)
		for _, cb := range c {
			if err = cb(actx); err != nil {
				return
			} else if actx.Err != nil {
				ctx.Flash = actx.Err
				return
			}
		}
		return
	}
}

func (c *BulkActionFormBuilderHandlers[T]) Append(handlers ...BulkActionFormBuilderHandler[T]) {
	*c = append(*c, handlers...)
}

func (c *BulkActionFormBuilderHandlers[T]) Prepend(handlers ...BulkActionFormBuilderHandler[T]) {
	*c = append(handlers, *c...)
}

func (b *BulkActionFormBuilder[T]) Handler(h BulkActionFormBuilderHandler[T]) *BulkActionFormBuilder[T] {
	b.handlers.Append(h)
	return b
}

func (b *BulkActionFormBuilder[T]) SetInsertFormBefore(v bool) *BulkActionFormBuilder[T] {
	b.insertFormBefore = v
	return b
}

func (b *BulkActionFormBuilder[T]) componentFunc(_ []string, ctx *web.EventContext) (_ h.HTMLComponent, err error) {
	formCtx := GetBulkActionFormContext[T](ctx)
	return b.eb.ToComponent(&ToComponentOptions{
		SkipPermVerify: true,
	}, formCtx.Form, FieldModeStack{NEW}, ctx), nil
}

func (b *BulkActionFormBuilder[T]) wrapComponentFunc(old BulkActionComponentFunc) BulkActionComponentFunc {
	var do = old
	if do == nil {
		do = b.componentFunc
	} else {
		do = func(selectedIds []string, ctx *web.EventContext) (comp h.HTMLComponent, err error) {
			var oldComp, form h.HTMLComponent

			if b.insertFormBefore {
				if form, err = b.componentFunc(selectedIds, ctx); err != nil {
					return
				}

				if oldComp, err = old(selectedIds, ctx); err != nil {
					return
				}

				return h.HTMLComponents{form, oldComp}, nil
			}

			if oldComp, err = old(selectedIds, ctx); err != nil {
				return
			}

			if form, err = b.componentFunc(selectedIds, ctx); err != nil {
				return
			}

			return h.HTMLComponents{oldComp, form}, nil
		}
	}
	return func(selectedIds []string, ctx *web.EventContext) (comp h.HTMLComponent, err error) {
		formCtx := GetBulkActionFormContext[T](ctx)
		if formCtx == nil {
			formCtx = NewBulkActionFormContext[T](ctx, nil, selectedIds)
			ctx.WithContextValue(ctxBulkActionFormContext, formCtx)
		}
		return do(selectedIds, ctx)
	}
}

func (b *BulkActionFormBuilder[T]) wrapUpdateFunc(old BulkActionUpdateFunc) BulkActionUpdateFunc {
	if old != nil {
		b.handlers.Append(func(ctx *BulkActionFormContext[T]) (err error) {
			return old(ctx.SelectedID, ctx.Context, ctx.Response)
		})
	}
	return b.handlers.Handler()
}

func (b *BulkActionFormBuilder[T]) defaultHandler(ctx *BulkActionFormContext[T]) (err error) {
	verr := b.eb.RunSetterFunc(&FieldsSetterOptions{
		SkipPermVerify: true,
	}, ctx.Context, true, ctx.Form)
	if verr.HaveErrors() {
		goto done
	}

	if verr = b.eb.Validators.Validate(ctx.Form, FieldModeStack{NEW}, ctx.Context); verr.HaveErrors() {
		goto done
	}

	b.eb.FieldsBuilder.Walk(b.eb.mb.modelInfo, ctx.Form, FieldModeStack{NEW}, ctx.Context, func(field *FieldContext) (s FieldWalkState) {
		if field.Field.IsEnabled(field) {
			verr.Merge(field.Field.Validators.Validate(field))
		}
		return s
	})
done:
	if verr.HaveErrors() {
		ctx.Err = &verr
	}
	return
}

func (b *BulkActionFormBuilder[T]) Build() *BulkActionBuilder {
	b.handlers.Prepend(b.defaultHandler)
	b.action.WrapComponentFunc(b.wrapComponentFunc)
	return b.action.WrapUpdateFunc(b.wrapUpdateFunc)
}

func BulkActionForm[T any](action *BulkActionBuilder, eb *EditingBuilder, h ...BulkActionFormBuilderHandler[T]) *BulkActionFormBuilder[T] {
	return &BulkActionFormBuilder[T]{
		action:   action,
		eb:       eb,
		handlers: h,
	}
}
