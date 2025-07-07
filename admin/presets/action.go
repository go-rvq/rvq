package presets

import (
	"context"
	"fmt"
	"net/url"

	"github.com/iancoleman/strcase"
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	"github.com/qor5/x/v3/perm"
	. "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type ActionLinkHandler func(baseModel *ModelBuilder, ctx *web.EventContext, r *web.EventResponse, q url.Values, id string)
type ActionOnClickHandler func(ctx *web.EventContext, id string, obj any) string

type ActionType uint8

const (
	ActionTypeListItem ActionType = iota
	ActionTypeList
	ActionTypeDetailing
	ActionTypePage
)

type ListingDoActionOptions struct {
	ListReloadDisabled bool
	Flash              any
}

type PageDoActionOptions struct {
	ReloadPage         bool
	AutoReloadDisabled bool
}

type ActionBuilder struct {
	NameLabel

	p                           *Builder
	onClick                     ActionOnClickHandler
	buttonCompFunc              ButtonComponentFunc
	updateFunc                  ActionUpdateFunc
	compFunc                    ActionComponentFunc
	buildCompFunc               ActionBuildComponentFunc
	wrapButtonFunc              func(old ButtonComponentFunc) ButtonComponentFunc
	fullComponent               bool
	dialogWidth                 string
	buttonColor                 string
	icon                        string
	doBtnLabel                  func(ctx *web.EventContext) string
	enabledFunc                 ActionEnabledFunc
	enabledObjFunc              ActionEnabledObjFunc
	linkHandler                 ActionLinkHandler
	showInList                  bool
	typ                         ActionType
	db                          *DetailingBuilder
	primaryActionPositionBottom bool
	verifier                    func(ctx *web.EventContext) *perm.Verifier
}

func Action(p *Builder, name string) *ActionBuilder {
	return &ActionBuilder{p: p, NameLabel: NameLabel{name: name}}
}

func getAction(actions []*ActionBuilder, name string) *ActionBuilder {
	for _, f := range actions {
		if f.name == name {
			return f
		}
	}
	return nil
}

func (b *ActionBuilder) Builder() *Builder {
	return b.p
}

func (b *ActionBuilder) String() string {
	return b.name
}

func (b *ActionBuilder) Type() ActionType {
	return b.typ
}

func (b *ActionBuilder) Label(v string) (r *ActionBuilder) {
	b.label = v
	return b
}

func (b *ActionBuilder) Icon(v string) *ActionBuilder {
	b.icon = v
	return b
}

func (b *ActionBuilder) SetI18nLabel(i18nLabel func(ctx context.Context) string) *ActionBuilder {
	b.NameLabel.SetI18nLabel(i18nLabel)
	return b
}

// ButtonCompFunc defines the components of the button area.
func (b *ActionBuilder) ButtonCompFunc(v ButtonComponentFunc) (r *ActionBuilder) {
	b.buttonCompFunc = v
	return b
}

func (b *ActionBuilder) WrapButtonFunc(v func(old ButtonComponentFunc) ButtonComponentFunc) *ActionBuilder {
	b.wrapButtonFunc = v
	return b
}

func (b *ActionBuilder) OnClick(f ActionOnClickHandler) *ActionBuilder {
	b.onClick = f
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

func (b *ActionBuilder) BuildComponentFunc(f ActionBuildComponentFunc) (r *ActionBuilder) {
	b.buildCompFunc = f
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

func (b *ActionBuilder) ShowInList() *ActionBuilder {
	b.showInList = true
	return b
}

func (b *ActionBuilder) SetShowInList(v bool) *ActionBuilder {
	b.showInList = v
	return b
}

func (b *ActionBuilder) SetVerifier(v func(ctx *web.EventContext) *perm.Verifier) *ActionBuilder {
	b.verifier = v
	return b
}

func (b *ActionBuilder) PrimaryActionPositionBottom(v bool) *ActionBuilder {
	b.primaryActionPositionBottom = v
	return b
}

func (b *ActionBuilder) IsEnabledObj(obj any, id string, ctx *web.EventContext) (ok bool, err error) {
	if b.enabledObjFunc == nil {
		return true, nil
	}
	return b.enabledObjFunc(obj, id, ctx)
}

func (b *ActionBuilder) RequestTitle(mb *ModelBuilder, ctx context.Context) (label string) {
	label = b.labelKey
	if label == "" {
		if b.i18nLabel != nil {
			return b.i18nLabel(ctx)
		}

		label = b.name
	}

	return i18n.Translate(mb.ActionTranslator(), ctx, label)
}

func (b *ActionBuilder) LinkHandler() ActionLinkHandler {
	return b.linkHandler
}

func (b *ActionBuilder) SetLinkHandler(linkHandler ActionLinkHandler) *ActionBuilder {
	b.linkHandler = linkHandler
	return b
}

func (b *ActionBuilder) form(baseModel *ModelBuilder, id string, ctx *web.EventContext) (comp h.HTMLComponent, err error) {
	if comp, err = b.Form(baseModel, id, actions.OverlayMode(ctx.Param(ParamOverlay)), ctx); err != nil {
		return
	}
	comp = web.Scope(comp).FormInit()
	return
}

func (b *ActionBuilder) View(baseModel *ModelBuilder, id string, ctx *web.EventContext, r *web.EventResponse) (err error) {
	if b.verifier != nil {
		if b.verifier(ctx).Denied() {
			return perm.PermissionDenied
		}
	}

	if b.linkHandler != nil {
		q := ctx.R.URL.Query()
		q.Del("__execute_event__")
		q.Del("actionOpen")

		b.linkHandler(baseModel, ctx, r, q, id)
		return
	}

	var comp h.HTMLComponent
	if comp, err = b.form(baseModel, id, ctx); err != nil {
		return
	}

	p := b.p
	if p == nil {
		p = baseModel.p
	}
	p.Dialog(b.dialogWidth).
		SetValidTargetPortalName(ctx.R.FormValue(ParamTargetPortal)).
		Respond(ctx, r, comp)
	return
}

func (b *ActionBuilder) Do(baseModel *ModelBuilder, id string, ctx *web.EventContext, r *web.EventResponse) (success bool, err error) {
	if b.verifier != nil {
		if b.verifier(ctx).Denied() {
			return false, perm.PermissionDenied
		}
	}

	if b.enabledFunc != nil {
		var ok bool
		if ok, err = b.IsEnabled(id, ctx); err != nil {
			return
		}
		if !ok {
			err = ErrActionNotAllowed
			return
		}
	}
	if b.enabledObjFunc != nil {
		var obj any
		if obj, err = b.db.Fetch(id, ctx); err != nil {
			return
		}

		var ok bool
		if ok, err = b.IsEnabledObj(obj, id, ctx); err != nil {
			return
		} else if !ok {
			err = ErrActionNotAllowed
			return
		}
	}

	switch b.typ {
	case ActionTypePage:
		ctx.WithData(&PageDoActionOptions{})
	case ActionTypeList:
		ctx.WithData(&ListingDoActionOptions{})
	}

	if err = b.updateFunc(id, ctx); err != nil {
		if verr, _ := err.(*web.ValidationErrors); verr != nil {
			if verr.HaveGlobalErrors() {
				ctx.Flash = verr
			}
			err = nil
		} else {
			return
		}
	}

	if ctx.W.Writed() {
		return
	}

	if ctx.Flash != nil {
		switch ctx.Flash.(type) {
		case string, *FlashMessage, FlashMessages:
			goto done
		default:
			var comp h.HTMLComponent
			if comp, err = b.form(baseModel, id, ctx); err != nil {
				return
			}
			p := b.p
			if p == nil {
				p = baseModel.p
			}
			p.DialogPortal(ctx.Param(ParamTargetPortal), b.dialogWidth).Respond(ctx, r, comp)
		}
		return
	} else {
		ctx.Flash = MustGetMessages(ctx.Context()).SuccessfullyExecutedAction
	}
done:
	success = true

	switch b.typ {
	case ActionTypeDetailing:
		if !IsInDialog(ctx) {
			r.PushState = web.Location(url.Values{})
		}
		r.RunScript = "closer.show = false"
	case ActionTypePage:
		data := ctx.Data().(*PageDoActionOptions)
		if data.ReloadPage {
			WithEventHandlerWrapperNoFlash(ctx)
			r.RunScript = "window.location.href = window.location.href"
		} else {
			if !data.AutoReloadDisabled && !IsInDialog(ctx) {
				r.PushState = web.Location(url.Values{})
			}
			r.RunScript = "closer.show = false"
		}
	case ActionTypeList, ActionTypeListItem:
		data := ctx.Data().(*ListingDoActionOptions)

		if isInDialogFromQuery(ctx) {
			r.AppendRunScript("closer.show = false")
		}

		if !data.ListReloadDisabled {
			r.AppendRunScript("presetsListing.loader.go()")
		}
	}

	if ctx.Flash != nil {
		ShowMessage(r, ctx.Flash)
	}

	return
}

func (b *ActionBuilder) Form(mb *ModelBuilder, id string, overlay actions.OverlayMode, ctx *web.EventContext) (comp h.HTMLComponent, err error) {
	if b.fullComponent {
		return b.compFunc(id, ctx)
	}

	var body h.HTMLComponent

	if b.compFunc != nil {
		if body, err = b.compFunc(id, ctx); err != nil {
			return
		}
	}

	cb := ContentComponentBuilder{
		Overlay: &ContentComponentBuilderOverlay{
			Mode:  overlay,
			Width: b.dialogWidth,
		},
		Title: b.RequestTitle(mb, ctx.Context()),
		Body:  body,
	}

	if b.updateFunc != nil {
		click := web.Plaid().
			MergeQuery(true).
			URL(ctx.R.RequestURI)

		switch b.typ {
		case ActionTypeList:
			click.EventFunc(actions.DoListingAction)
		case ActionTypeListItem, ActionTypePage, ActionTypeDetailing:
			click.EventFunc(actions.DoAction)
		}

		btn := VBtn("").
			Color("primary").
			Title(b.DoBtnLabel(ctx)).
			Variant(VariantFlat).
			Attr(":disabled", "isFetching").
			Attr(":loading", "isFetching").
			Icon(true).
			Density("comfortable").
			Children(VIcon("mdi-chevron-right")).
			Attr("@click", click.Go())

		if b.primaryActionPositionBottom {
			cb.BottomActions = append(cb.BottomActions, btn)
		} else {
			cb.PrimaryAction = btn
		}
	}

	if b.buildCompFunc != nil {
		if err = b.buildCompFunc(id, ctx, &cb); err != nil {
			return
		}
	}

	if ctx.Flash != nil {
		cb.Notice(ctx.Flash)
		ctx.Flash = nil
	}

	return cb.BuildOverlay(), nil
}

func (b *ActionBuilder) PermName() string {
	switch b.typ {
	case ActionTypeList, ActionTypeDetailing:
		return strcase.ToSnake(b.name)
	default:
		return "list_item:" + strcase.ToSnake(b.name)
	}
}

func (b *ActionBuilder) BuildButton(defaultBtnBuilder ButtonComponentFunc, onclick *web.VueEventTagBuilder, id string, obj any, ctx *web.EventContext) h.HTMLComponent {
	var (
		btnf = b.buttonCompFunc
		oc   = &OnClick{
			Builder: onclick,
		}
	)

	if btnf == nil {
		btnf = defaultBtnBuilder
	}

	if b.onClick != nil {
		oc.Raw = b.onClick(ctx, id, obj)
	}

	if b.wrapButtonFunc != nil {
		btnf = b.wrapButtonFunc(btnf)
	}

	return btnf(ctx, oc)
}

func BuildMenuItemCompomentsOfActions(sharedPortal string, ctx *web.EventContext, mb *ModelBuilder, id string, obj any, actionBuilders ...*ActionBuilder) (items []*VListItemBuilder, errors h.HTMLComponents) {
	for _, action := range actionBuilders {
		if mb.permissioner.ReqObjectActioner(ctx.R, obj, action.PermName()).Denied() {
			continue
		}

		if ok, err := action.IsEnabledObj(obj, id, ctx); err != nil {
			errors = append(errors, VAlert(h.Text(fmt.Sprintf("Action %q: IsEnabledObj: %v", action.RequestTitle(mb, ctx.Context()), err))).Color("error"))
		} else if !ok {
			continue
		}

		var (
			onclick = web.Plaid().
				EventFunc(actions.Action).
				Query(ParamID, id).
				Query(ParamAction, action.name).
				Query(ParamTargetPortal, sharedPortal).
				Query(ParamOverlay, actions.Dialog).
				URL(mb.Info().ListingHrefCtx(ctx))

			btn = action.BuildButton(func(ctx *web.EventContext, onclick *OnClick) h.HTMLComponent {
				return VListItem(
					h.If(
						action.icon != "", web.Slot(
							VIcon(action.icon),
						).Name("prepend"),
					),
					VListItemTitle(h.Text(action.RequestTitle(mb, ctx.Context()))),
				).Attr("@click", onclick.String())
			}, onclick, id, obj, ctx)
		)

		if li, _ := btn.(*VListItemBuilder); li != nil {
			items = append(items, li)
		} else {
			items = append(items, VListItem(btn))
		}

	}

	return
}

func (b *ListingBuilder) Action(name string) (r *ActionBuilder) {
	if r = getAction(b.actions, name); r != nil {
		return
	}

	r = &ActionBuilder{
		db:  b.mb.detailing,
		typ: ActionTypeList,
	}
	r.name = name
	b.actions = append(b.actions, r)
	return
}

func (b *ListingBuilder) ItemAction(name string) (r *ActionBuilder) {
	if r = getAction(b.itemActions, name); r != nil {
		return
	}

	r = &ActionBuilder{
		db:  b.mb.detailing,
		typ: ActionTypeListItem,
	}
	r.name = name
	b.itemActions = append(b.itemActions, r)
	return
}

func (b *ListingBuilder) AppendItemAction(actions ...*ActionBuilder) *ListingBuilder {
	b.itemActions = append(b.itemActions, actions...)
	return b
}

func (b *DetailingBuilder) Action(name string) (r *ActionBuilder) {
	if r = getAction(b.actions, name); r != nil {
		return r
	}
	r = &ActionBuilder{
		typ: ActionTypeDetailing,
		db:  b,
	}
	r.name = name
	b.actions = append(b.actions, r)
	return
}
