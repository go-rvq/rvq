package presets

import (
	"fmt"
	"net/url"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/web/v3/datafield"
	"github.com/qor5/x/v3/i18n"
	. "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type ActionFormBuilderHandler[T any] func(ctx *ActionFormContext[T]) (err error)

type ActionFormContext[T any] struct {
	Context *web.EventContext
	ID      string
	Obj     T
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
		actx := NewActionUpdateContext[T](ctx, id)
		for _, cb := range c {
			if err = cb(actx); err != nil {
				ctx.Flash = err
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

type ActionLinkHandler func(baseModel *ModelBuilder, ctx *web.EventContext, r *web.EventResponse, q url.Values, id string)

type ActionBuilder struct {
	NameLabel

	buttonCompFunc ComponentFunc
	updateFunc     ActionUpdateFunc
	compFunc       ActionComponentFunc
	fullComponent  bool
	dialogWidth    string
	buttonColor    string
	icon           string
	doBtnLabel     func(ctx *web.EventContext) string
	enabledFunc    ActionEnabledFunc
	enabledObjFunc ActionEnabledObjFunc
	linkHandler    ActionLinkHandler
}

func getAction(actions []*ActionBuilder, name string) *ActionBuilder {
	for _, f := range actions {
		if f.name == name {
			return f
		}
	}
	return nil
}

func (b *ActionBuilder) Label(v string) (r *ActionBuilder) {
	b.label = v
	return b
}

func (b *ActionBuilder) Icon(v string) *ActionBuilder {
	b.icon = v
	return b
}

func (b *ActionBuilder) SetI18nLabel(i18nLabel func(ctx web.ContextValuer) string) *ActionBuilder {
	b.NameLabel.SetI18nLabel(i18nLabel)
	return b
}

// ButtonCompFunc defines the components of the button area.
func (b *ActionBuilder) ButtonCompFunc(v ComponentFunc) (r *ActionBuilder) {
	b.buttonCompFunc = v
	return b
}

// UpdateFunc defines event when the button is clicked.
func (b *ActionBuilder) UpdateFunc(v ActionUpdateFunc) (r *ActionBuilder) {
	b.updateFunc = v
	return b
}

func (b *ActionBuilder) WrapUpdateFunc(f func(old ActionUpdateFunc) ActionUpdateFunc) *ActionBuilder {
	b.updateFunc = f(b.updateFunc)
	return b
}

// ComponentFunc defines the components in dialog of button click.
func (b *ActionBuilder) ComponentFunc(v ActionComponentFunc) (r *ActionBuilder) {
	b.compFunc = v
	return b
}

func (b *ActionBuilder) SetDialogWidth(v string) (r *ActionBuilder) {
	b.dialogWidth = v
	return b
}

// SetButtonColor defines the color of default button if buttonCompFunc is not defined.
func (b *ActionBuilder) SetButtonColor(v string) (r *ActionBuilder) {
	b.buttonColor = v
	return b
}

func (b *ActionBuilder) SetFullComponent(v bool) (r *ActionBuilder) {
	b.fullComponent = v
	return b
}

func (b *ActionBuilder) SetDoBtnLabel(f func(ctx *web.EventContext) string) *ActionBuilder {
	b.doBtnLabel = f
	return b
}

func (b *ActionBuilder) DoBtnLabel(ctx *web.EventContext) string {
	if b.doBtnLabel != nil {
		return b.doBtnLabel(ctx)
	}
	return MustGetMessages(ctx.Context()).Update
}

func (b *ActionBuilder) SetEnabled(f ActionEnabledFunc) *ActionBuilder {
	b.enabledFunc = f
	return b
}

func (b *ActionBuilder) GetEnabled() ActionEnabledFunc {
	return b.enabledFunc
}

func (b *ActionBuilder) WrapEnabled(f func(old ActionEnabledFunc) ActionEnabledFunc) *ActionBuilder {
	return b.SetEnabled(f(b.enabledFunc))
}

func (b *ActionBuilder) IsEnabled(id string, ctx *web.EventContext) (ok bool, err error) {
	if b.enabledFunc == nil {
		return true, nil
	}
	return b.enabledFunc(id, ctx)
}

func (b *ActionBuilder) GetEnabledObj() ActionEnabledObjFunc {
	return b.enabledObjFunc
}

func (b *ActionBuilder) SetEnabledObj(f ActionEnabledObjFunc) *ActionBuilder {
	b.enabledObjFunc = f
	return b
}

func (b *ActionBuilder) WrapEnabledObj(f func(old ActionEnabledObjFunc) ActionEnabledObjFunc) *ActionBuilder {
	return b.SetEnabledObj(f(b.enabledObjFunc))
}

func (b *ActionBuilder) IsEnabledObj(obj any, id string, ctx *web.EventContext) (ok bool, err error) {
	if b.enabledObjFunc == nil {
		return true, nil
	}
	return b.enabledObjFunc(obj, id, ctx)
}

func (b *ActionBuilder) RequestTitle(mb *ModelBuilder, ctx web.ContextValuer) (label string) {
	label = b.labelKey
	if label == "" {
		if b.i18nLabel != nil {
			return b.i18nLabel(ctx)
		}

		label = b.name
	}

	return i18n.Translate(mb.ActionTranslator(), ctx.Context(), label)
}

func (b *ActionBuilder) LinkHandler() ActionLinkHandler {
	return b.linkHandler
}

func (b *ActionBuilder) SetLinkHandler(linkHandler ActionLinkHandler) *ActionBuilder {
	b.linkHandler = linkHandler
	return b
}

func (b *ActionBuilder) form(baseModel *ModelBuilder, id string, ctx *web.EventContext) h.HTMLComponent {
	comp := b.Form(baseModel, id, actions.OverlayMode(ctx.Param(ParamOverlay)), ctx)
	return web.Scope(comp).FormInit()
}

func (b *ActionBuilder) View(baseModel *ModelBuilder, id string, ctx *web.EventContext, r *web.EventResponse) (err error) {
	if b.linkHandler != nil {
		q := ctx.R.URL.Query()
		q.Del("__execute_event__")
		q.Del("actionOpen")

		b.linkHandler(baseModel, ctx, r, q, id)
		return
	}

	baseModel.p.dialog(ctx, r, b.form(baseModel, id, ctx), b.dialogWidth)
	return
}

func (b *ActionBuilder) Do(baseModel *ModelBuilder, id string, ctx *web.EventContext, r *web.EventResponse) (err error) {
	if err = b.updateFunc(id, ctx); err != nil || ctx.Flash != nil {
		if ctx.Flash == nil {
			ctx.Flash = err
		}

		baseModel.p.dialog(ctx, r, b.form(baseModel, id, ctx), b.dialogWidth)
		return
	}

	r.PushState = web.Location(url.Values{})
	r.RunScript = "closer.show = false"
	GetFlashMessages(ctx).RespondTo(r)
	return nil
}

func (b *ActionBuilder) Form(mb *ModelBuilder, id string, overlay actions.OverlayMode, ctx *web.EventContext) h.HTMLComponent {
	if b.fullComponent {
		return b.compFunc(id, ctx)
	}

	var body h.HTMLComponent

	if b.compFunc != nil {
		body = b.compFunc(id, ctx)
	}

	cb := ContentComponentBuilder{
		Overlay: &ContentComponentBuilderOverlay{
			Mode:       overlay,
			Scrollable: true,
		},
		Title: b.RequestTitle(mb, ctx),
		Body:  body,
	}

	if b.updateFunc != nil {
		cb.PrimaryAction = VBtn("").
			Color("primary").
			Title(b.DoBtnLabel(ctx)).
			Variant(VariantFlat).
			Attr(":disabled", "isFetching").
			Attr(":loading", "isFetching").
			Icon(true).
			Density("comfortable").
			Children(VIcon("mdi-chevron-right")).
			Attr("@click", web.Plaid().
				EventFunc(actions.DoAction).
				MergeQuery(true).
				URL(ctx.R.RequestURI).
				Go())
	}

	if ctx.Flash != nil {
		cb.Notice(ctx.Flash)
	}

	WithRespondDialogHandlers(ctx, func(d *DialogBuilder) {
		d.SetScrollable(true)
	})

	return cb.BuildOverlay()
}

func BuildMenuItemCompomentsOfActions(sharedPortal string, ctx *web.EventContext, mb *ModelBuilder, id string, obj any, actionBuilders ...*ActionBuilder) (items []*VListItemBuilder, errors h.HTMLComponents) {
	for _, action := range actionBuilders {
		if mb.Info().Verifier().SnakeDo(PermActions, action.name).WithReq(ctx.R).IsAllowed() != nil {
			continue
		}

		if ok, err := action.IsEnabledObj(obj, id, ctx); err != nil {
			errors = append(errors, VAlert(h.Text(fmt.Sprintf("Action %q: IsEnabledObj: %v", action.RequestTitle(mb, ctx), err))).Color("error"))
		} else if !ok {
			continue
		}

		if action.buttonCompFunc != nil {
			items = append(items, VListItem(action.buttonCompFunc(ctx)))
		} else {
			buttonColor := action.buttonColor
			if buttonColor == "" {
				buttonColor = ColorPrimary
			}

			onclick := web.Plaid().
				EventFunc(actions.Action).
				Query(ParamID, id).
				Query(ParamAction, action.name).
				Query(ParamTargetPortal, sharedPortal).
				Query(ParamOverlay, actions.Dialog).
				URL(ctx.R.RequestURI)
			items = append(items, VListItem(
				h.If(
					action.icon != "", web.Slot(
						VIcon(action.icon),
					).Name("prepend"),
				),
				VListItemTitle(h.Text(action.RequestTitle(mb, ctx))),
			).Attr("@click", onclick.Go()))
		}
	}

	return
}

func ActionForm[T any](action *ActionBuilder, eb *EditingBuilder, handler ActionFormBuilderHandler[T]) *ActionBuilder {
	if action.compFunc == nil {
		action.compFunc = func(id string, ctx *web.EventContext) h.HTMLComponent {
			obj := ctx.ContextValue(ctxActionFormObject)
			if obj == nil {
				obj = eb.mb.NewModel()
			}
			return eb.ToComponent(obj, FieldModeStack{NEW}, ctx)
		}
	}
	handlers := ActionFormBuilderHandlers[T]{
		func(ctx *ActionFormContext[T]) (err error) {
			ctx.Obj = eb.mb.NewModel().(T)
			ctx.Context.WithContextValue(ctxActionFormObject, ctx.Obj)

			verr := eb.RunSetterFunc(ctx.Context, true, ctx.Obj)
			if verr.HaveErrors() {
				goto done
			}

			if verr = eb.Validators.Validate(ctx.Obj, FieldModeStack{NEW}, ctx.Context); verr.HaveErrors() {
				goto done
			}

			eb.FieldsBuilder.Walk(eb.mb.modelInfo, ctx.Obj, FieldModeStack{NEW}, ctx.Context, func(field *FieldContext) (s FieldWalkState) {
				verr.Merge(field.Field.Validators.Validate(field))
				return s
			})
		done:
			if verr.HaveErrors() {
				err = &verr
			}
			return
		},
		handler,
	}
	return action.WrapUpdateFunc(func(old ActionUpdateFunc) ActionUpdateFunc {
		if old != nil {
			handlers.Append(func(ctx *ActionFormContext[T]) (err error) {
				return old(ctx.ID, ctx.Context)
			})
		}
		return handlers.Handler()
	})
}

func (b *ListingBuilder) Action(name string) (r *ActionBuilder) {
	builder := getAction(b.actions, name)
	if builder != nil {
		return builder
	}

	r = &ActionBuilder{}
	r.name = name
	b.actions = append(b.actions, r)
	return
}

func (b *ListingBuilder) ItemAction(actions ...*ActionBuilder) *ListingBuilder {
	b.itemActions = append(b.itemActions, actions...)
	return b
}

func (b *DetailingBuilder) Action(name string) (r *ActionBuilder) {
	builder := getAction(b.actions, name)
	if builder != nil {
		return builder
	}

	r = &ActionBuilder{}
	r.name = name
	b.actions = append(b.actions, r)
	return
}
