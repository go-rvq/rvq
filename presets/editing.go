package presets

import (
	"github.com/jinzhu/inflection"
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/perm"
	. "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type EditingBuilder struct {
	mb               *ModelBuilder
	Fetcher          FetchFunc
	Setter           SetterFunc
	Saver            SaveFunc
	Validators       Validators
	tabPanels        []TabComponentFunc
	hiddenFuncs      []ObjectComponentFunc
	sidePanel        ObjectComponentFunc
	actionsFunc      ObjectComponentFunc
	editingTitleFunc EditingTitleComponentFunc
	onChangeAction   OnChangeActionFunc
	pageFunc         web.PageFunc
	editionDisabled  OkHandled
	FieldsBuilder
}

func NewEditingBuilder(mb *ModelBuilder, fieldsBuilder FieldsBuilder) *EditingBuilder {
	eb := &EditingBuilder{mb: mb, FieldsBuilder: fieldsBuilder}
	return eb
}

func (mb *ModelBuilder) newEditing() (r *EditingBuilder) {
	mb.editing = NewEditingBuilder(mb, *mb.NewFieldsBuilder(mb.writeFieldBuilders.HasMode(LIST)...))

	mb.editing.FetchFunc(mb.Fetcher)
	mb.editing.SaveFunc(mb.Saver)
	return
}

// string / []string / *FieldsSection
func (mb *ModelBuilder) Editing(vs ...interface{}) (r *EditingBuilder) {
	r = mb.editing
	if len(vs) == 0 {
		return
	}

	r.Only(vs...)
	return r
}

func (mb *ModelBuilder) SetEditingBuilder(b *EditingBuilder) {
	mb.editing = b
}

func (b *EditingBuilder) ModelBuilder() *ModelBuilder {
	return b.mb
}

// string / []string / *FieldsSection
func (b *EditingBuilder) Only(vs ...interface{}) (r *EditingBuilder) {
	r = b
	r.FieldsBuilder = *r.FieldsBuilder.Only(vs...)
	return
}

func (b *EditingBuilder) Except(vs ...string) (r *EditingBuilder) {
	r = b
	r.FieldsBuilder = *r.FieldsBuilder.Except(vs...)
	return
}

func (b *EditingBuilder) FetchFunc(v FetchFunc) (r *EditingBuilder) {
	b.Fetcher = v
	return b
}

func (b *EditingBuilder) WrapFetchFunc(w func(in FetchFunc) FetchFunc) (r *EditingBuilder) {
	b.Fetcher = w(b.Fetcher)
	return b
}

func (b *EditingBuilder) SaveFunc(v SaveFunc) (r *EditingBuilder) {
	b.Saver = v
	return b
}

func (b *EditingBuilder) WrapSaveFunc(w func(in SaveFunc) SaveFunc) (r *EditingBuilder) {
	b.Saver = w(b.Saver)
	return b
}

func (b *EditingBuilder) SetterFunc(v SetterFunc) (r *EditingBuilder) {
	b.Setter = v
	return b
}

func (b *EditingBuilder) OnChangeActionFunc(v OnChangeActionFunc) (r *EditingBuilder) {
	b.onChangeAction = v
	return b
}

func (b *EditingBuilder) WrapSetterFunc(w func(in SetterFunc) SetterFunc) (r *EditingBuilder) {
	b.Setter = w(b.Setter)
	return b
}

func (b *EditingBuilder) AppendTabsPanelFunc(v TabComponentFunc) (r *EditingBuilder) {
	b.tabPanels = append(b.tabPanels, v)
	return b
}

func (b *EditingBuilder) TabsPanels(vs ...TabComponentFunc) (r *EditingBuilder) {
	b.tabPanels = vs
	return b
}

func (b *EditingBuilder) SidePanelFunc(v ObjectComponentFunc) (r *EditingBuilder) {
	b.sidePanel = v
	return b
}

func (b *EditingBuilder) AppendHiddenFunc(v ObjectComponentFunc) (r *EditingBuilder) {
	b.hiddenFuncs = append(b.hiddenFuncs, v)
	return b
}

func (b *EditingBuilder) ActionsFunc(v ObjectComponentFunc) (r *EditingBuilder) {
	b.actionsFunc = v
	return b
}

func (b *EditingBuilder) EditingTitleFunc(v EditingTitleComponentFunc) (r *EditingBuilder) {
	b.editingTitleFunc = v
	return b
}

func (b *EditingBuilder) FetchAndUnmarshal(id string, removeDeletedAndSort bool, ctx *web.EventContext) (obj interface{}, vErr web.ValidationErrors) {
	obj = b.mb.NewModel()
	if len(id) > 0 || b.mb.singleton {
		err1 := b.Fetcher(obj, id, ctx)
		if err1 != nil {
			if !(err1 == ErrRecordNotFound && b.mb.singleton) {
				vErr.GlobalError(err1.Error())
				// b.UpdateOverlayContent(ctx, &r, obj, "", err1)
				return
			}
		}
	}

	vErr = b.RunSetterFunc(ctx, removeDeletedAndSort, obj)
	return
}

func (b *EditingBuilder) RunSetterFunc(ctx *web.EventContext, removeDeletedAndSort bool, toObj interface{}) (vErr web.ValidationErrors) {
	if b.Setter != nil {
		b.Setter(toObj, ctx)
	}

	vErr = b.Unmarshal(toObj, b.mb.Info(), removeDeletedAndSort, ctx)

	return
}

func (mb *EditingBuilder) EditionDisabled() OkHandled {
	return mb.editionDisabled
}

func (mb *EditingBuilder) SetEditionDisabled(editDisabled OkHandled) *EditingBuilder {
	mb.editionDisabled = editDisabled
	return mb
}

func (mb *EditingBuilder) CanEdit(ctx *web.EventContext) bool {
	if mb.editionDisabled == nil {
		return !CallOkHandled(mb.mb.editionDisabled, ctx)
	}
	return !CallOkHandled(mb.editionDisabled, ctx)
}

func (mb *EditingBuilder) CanEditObj(ctx *web.EventContext, obj interface{}) bool {
	if !mb.CanEdit(ctx) {
		return false
	}

	return mb.mb.Info().CanUpdate(ctx.R, obj)
}

func (b *EditingBuilder) UpdateOverlayContent(
	ctx *web.EventContext,
	r *web.EventResponse,
	obj interface{},
	successMessage string,
	err error,
) {
	ctx.Flash = err

	if err != nil {
		if _, ok := err.(*web.ValidationErrors); !ok {
			vErr := &web.ValidationErrors{}
			vErr.GlobalError(err.Error())
			ctx.Flash = vErr
		}
	}

	if ctx.Flash == nil {
		ctx.Flash = successMessage
	}

	f := b.form(obj, ctx)
	if f.b.overlayMode.IsDrawer() {
		f.Portal = f.b.overlayMode.PortalName()
	}

	f.ScopeDisabled = ctx.R.FormValue(ParamEditFormUnscoped) == "true"

	f.Respond(r)
}

func (b *EditingBuilder) defaultPageFunc(ctx *web.EventContext) (r web.PageResponse, err error) {
	var (
		id  = ctx.R.PathValue(ParamID)
		obj = b.mb.NewModel()
	)

	if len(id) > 0 || b.mb.singleton {
		if err = b.Fetcher(obj, id, ctx); err != nil {
			if err == ErrRecordNotFound {
				if b.mb.singleton {
					err = nil
				} else {
					return b.mb.p.DefaultNotFoundPageFunc(ctx)
				}
			} else {
				return
			}
		}
	} else {
		return b.mb.p.DefaultNotFoundPageFunc(ctx)
	}

	r.Body = VContainer(h.Text(id))

	msgr := MustGetMessages(ctx.R)
	r.PageTitle = msgr.DetailingObjectTitle(inflection.Singular(b.mb.label), b.mb.RecordTitle(obj, ctx))

	if b.mb.Info().Verifier().Do(PermGet).ObjectOn(obj).WithReq(ctx.R).IsAllowed() != nil {
		r.Body = h.Div(h.Text(perm.PermissionDenied.Error()))
		return
	}

	portalName := actions.Edit

	overlay := actions.OverlayMode(ctx.R.FormValue(ParamOverlay))
	if overlay != "" {
		portalName = overlay.PortalName()
	}

	WithScope(ctx, web.Scope())

	// set portal to edit btn
	ctx.R.Form.Set(ParamTargetPortal, portalName)
	EditFormUnscoped(ctx, true)

	f := b.form(obj, ctx)
	comp := f.Component()

	r.Body = VContainer(web.Portal(comp).
		Name(portalName))

	// /ctx.WithContextValue(CtxActionsComponent, h.HTMLComponents{
	// .EditBtn(ctx, id, true),
	// })

	return
}

func (b *EditingBuilder) GetPageFunc() web.PageFunc {
	if b.pageFunc != nil {
		return b.pageFunc
	}
	return b.defaultPageFunc
}

func (b *EditingBuilder) SetPageFunc(pageFunc web.PageFunc) *EditingBuilder {
	b.pageFunc = pageFunc
	return b
}
