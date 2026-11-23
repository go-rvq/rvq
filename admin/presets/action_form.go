package presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/datafield"
	"reflect"
)

type ActionFormBuilderHandler[T any] func(ctx *ActionFormContext[T]) (err error)

type ActionFormContext[T any] struct {
	Context *web.EventContext
	ID      string
	Obj     any
	Form    T
	Err     error
	datafield.DataField[*ActionFormContext[T]]
}

func NewActionUpdateContext[T any](context *web.EventContext, id string) *ActionFormContext[T] {
	return datafield.New(&ActionFormContext[T]{
		Context: context,
		ID:      id,
	})
}

type ActionFormBuilderHandlers[T any] []ActionFormBuilderHandler[T]

func (c ActionFormBuilderHandlers[T]) Handler() ActionUpdateFunc {
	return func(id string, ctx *web.EventContext) (err error) {
		actx := GetActionFormContext[T](ctx)
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

func (c *ActionFormBuilderHandlers[T]) Append(handlers ...ActionFormBuilderHandler[T]) {
	*c = append(*c, handlers...)
}

func (c *ActionFormBuilderHandlers[T]) Prepend(handlers ...ActionFormBuilderHandler[T]) {
	*c = append(handlers, *c...)
}

type ActionFormBuilder[T any] struct {
	action      *ActionBuilder
	eb          *EditingBuilder
	fetch       bool
	handlers    ActionFormBuilderHandlers[T]
	initForm    func(fctx *ActionFormContext[T]) error
	refresh     func(fctx *ActionFormContext[T]) error
	decodedForm func(ctx *ActionFormContext[T]) (err error)
}

func (b *ActionFormBuilder[T]) Handler(h ActionFormBuilderHandler[T]) *ActionFormBuilder[T] {
	b.handlers.Append(h)
	return b
}

func (b *ActionFormBuilder[T]) Fetch(v bool) *ActionFormBuilder[T] {
	b.fetch = v
	return b
}

func (b *ActionFormBuilder[T]) InitForm(f func(fctx *ActionFormContext[T]) error) *ActionFormBuilder[T] {
	b.initForm = f
	return b
}

// Refresh changes the form decoded for it to be drawn again
// (ActionBuilder.FormRefresh): a choice that fills it, a total computed
// again. The queries of the click are in fctx.Context.R.
func (b *ActionFormBuilder[T]) Refresh(f func(fctx *ActionFormContext[T]) error) *ActionFormBuilder[T] {
	b.refresh = f
	return b
}

func (b *ActionFormBuilder[T]) DecodedForm(f func(fctx *ActionFormContext[T]) error) *ActionFormBuilder[T] {
	b.decodedForm = f
	return b
}

func (b *ActionFormBuilder[T]) FetchObject(fctx *ActionFormContext[T]) (err error) {
	var obj any
	if obj, err = b.action.db.Fetch(fctx.ID, fctx.Context); err != nil {
		return
	}
	fctx.Obj = obj
	return nil
}

func (b *ActionFormBuilder[T]) RenderForm(fctx *ActionFormContext[T]) (_ h.HTMLComponent, err error) {
	// its list editors: their + sends the action again (go-rvq/rvq#17)
	fctx.Context.WithContextValue(ctxListEditorActionForm, &listEditorActionForm{
		addRow: func(formKey string) *web.VueEventTagBuilder {
			return web.Plaid().
				MergeQuery(true).
				URL(actionRequestURI(fctx.Context.R)).
				EventFunc(b.action.doEvent()).
				Query(ParamActionFormAddRow, formKey)
		},
	})
	return b.eb.ToComponent(&ToComponentOptions{
		SkipPermVerify: true,
	}, fctx.Form, FieldModeStack{NEW}, fctx.Context), nil
}

func (b *ActionFormBuilder[T]) CreateContext(id string, ctx *web.EventContext) (fctx *ActionFormContext[T], err error) {
	fctx = NewActionUpdateContext[T](ctx, id)
	fctx.Form = b.eb.mb.NewModel().(T)

	if b.fetch {
		if err = b.FetchObject(fctx); err != nil {
			return
		}
	}

	if b.initForm != nil {
		if err = b.initForm(fctx); err != nil {
			return
		}
	}
	return
}

func (b *ActionFormBuilder[T]) GetOrCreateContext(id string, ctx *web.EventContext) (fctx *ActionFormContext[T], err error) {
	fctx = GetActionFormContext[T](ctx)
	if fctx == nil {
		if fctx, err = b.CreateContext(id, ctx); err != nil {
			return
		}
		WithActionFormContext(ctx, fctx)
	}
	return fctx, nil
}

func (b *ActionFormBuilder[T]) decodeForm(ctx *ActionFormContext[T]) (err error) {
	verr := b.eb.RunSetterFunc(&FieldsSetterOptions{
		SkipPermVerify: true,
	}, ctx.Context, true, ctx.Form)
	if verr.HaveErrors() {
		goto done
	}

	// its lists: the rows removed (__deleted) out, in their order
	dropListEditorDeleted(&b.eb.FieldsBuilder, ctx.Context, ctx.Form)

	if verr = b.eb.Validators.Validate(ctx.Form, FieldModeStack{NEW}, ctx.Context); verr.HaveErrors() {
		goto done
	}

	b.eb.FieldsBuilder.Walk(b.eb.mb.modelInfo, ctx.Form, FieldModeStack{NEW}, ctx.Context, func(field *FieldContext) (s FieldWalkState) {
		if field.Field.IsEnabled(field) {
			verr.Merge(field.Field.Validate(field))
		}
		return s
	})

	if b.decodedForm != nil {
		if err = b.decodedForm(ctx); err != nil {
			return
		}
	}

done:
	if verr.HaveErrors() {
		ctx.Err = &verr
	}
	return
}

func (b *ActionFormBuilder[T]) Build() *ActionBuilder {
	b.handlers.Prepend(b.decodeForm)

	// the + of a list of the form, a refresh: its values decoded (not
	// validated), the row added, the Refresh (go-rvq/rvq#17)
	b.action.formRedraw = func(id string, ctx *web.EventContext) (err error) {
		var fctx *ActionFormContext[T]
		if fctx, err = b.GetOrCreateContext(id, ctx); err != nil {
			return
		}
		b.eb.RunSetterFunc(&FieldsSetterOptions{SkipPermVerify: true}, ctx, false, fctx.Form)
		if formKey := ctx.R.FormValue(ParamActionFormAddRow); formKey != "" {
			if err = appendListEditorRows(ctx, fctx.Form, formKey); err != nil {
				return
			}
		}
		if ctx.R.FormValue(ParamActionFormRefresh) != "" && b.refresh != nil {
			err = b.refresh(fctx)
		}
		return
	}

	if b.fetch {
		b.handlers.Prepend(b.FetchObject)
	}

	if b.action.compFunc == nil {
		b.action.compFunc = func(_ string, ctx *web.EventContext) (_ h.HTMLComponent, err error) {
			return b.RenderForm(GetActionFormContext[T](ctx))
		}
	}

	old := b.action.compFunc
	b.action.compFunc = func(id string, ctx *web.EventContext) (_ h.HTMLComponent, err error) {
		if _, err = b.GetOrCreateContext(id, ctx); err != nil {
			return
		}
		return old(id, ctx)
	}

	return b.action.WrapUpdateFunc(func(old ActionUpdateFunc) ActionUpdateFunc {
		if old != nil {
			b.handlers.Append(func(ctx *ActionFormContext[T]) (err error) {
				return old(ctx.ID, ctx.Context)
			})
		}

		return func(id string, ctx *web.EventContext) (err error) {
			var fctx *ActionFormContext[T]
			if fctx, err = b.GetOrCreateContext(id, ctx); err != nil {
				return
			}
			for _, cb := range b.handlers {
				if err = cb(fctx); err != nil {
					return
				} else if fctx.Err != nil {
					ctx.Flash = fctx.Err
					return
				}
			}
			return
		}
	})
}

func ActionForm[T any](action *ActionBuilder, eb *EditingBuilder, h ...ActionFormBuilderHandler[T]) *ActionFormBuilder[T] {
	return &ActionFormBuilder[T]{
		action:   action,
		eb:       eb,
		handlers: h,
	}
}

// dropListEditorDeleted takes out of each list (NestedSlice) of fb in obj the
// rows its list editor flagged __deleted, the others in their __pos order
// (PartitionListEditorItems): what a form of an action decodes is what was
// left on it — an editing form leaves that to the save of the association.
func dropListEditorDeleted(fb *FieldsBuilder, ctx *web.EventContext, obj any) {
	vals := ListEditorFormValues(ctx)
	rv := reflect.ValueOf(obj)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return
	}
	for _, f := range fb.fields {
		if _, ok := f.nested.(*NestedSliceBuilder); !ok {
			continue
		}
		fv := rv.FieldByName(f.name)
		if !fv.IsValid() || fv.Kind() != reflect.Slice || !fv.CanSet() {
			continue
		}
		if part := PartitionListEditorItems(vals, f.name, fv); part != nil {
			fv.Set(part.KeptSlice(fv.Type()))
		}
	}
}
