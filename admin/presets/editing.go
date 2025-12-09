package presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/jinzhu/inflection"
)

type EditingBuilder struct {
	mb               *ModelBuilder
	Fetcher          FetchFunc
	Setter           SetterFunc
	PostSetter       SetterFunc
	Saver            SaveFunc
	Creator          CreateFunc
	New              func(ctx *web.EventContext, obj any) (err error)
	preValidate      func(ctx *web.EventContext, obj any) (err error)
	postValidate     func(ctx *web.EventContext, obj any) (err error)
	preSaveCallback  SaveCallbackFunc
	postSaveCallback SaveCallbackFunc
	Validators       Validators
	tabPanels        []TabComponentFunc
	hiddenFuncs      []ObjectComponentFunc
	sidePanel        ObjectComponentFunc
	actionsFunc      ObjectComponentFunc
	editingTitleFunc EditingTitleComponentFunc
	onChangeAction   OnChangeActionFunc
	pageFunc         web.PageFunc
	FieldsBuilder
	preComponents  []ModeObjectComponentFunc
	postComponents []ModeObjectComponentFunc
	maxPostSize    int64

	EditingRestrictionField[*EditingBuilder]
}

func NewEditingBuilder(mb *ModelBuilder, fieldsBuilder FieldsBuilder) *EditingBuilder {
	e := &EditingBuilder{
		mb:            mb,
		FieldsBuilder: fieldsBuilder,
		maxPostSize:   web.DefaulMaxPostSize,
	}
	e.EditingRestriction = NewObjRestriction(e, func(r *ObjRestriction[*EditingBuilder]) {
		r.Insert(mb.EditingRestriction)
	})
	return e
}

func (mb *ModelBuilder) newEditing() (r *EditingBuilder) {
	mb.editing = NewEditingBuilder(mb, *mb.NewFieldsBuilder(mb.writeFieldBuilders.HasMode(LIST)...))

	mb.editing.FetchFunc(mb.Fetcher)
	mb.editing.SaveFunc(mb.Saver)
	mb.editing.CreateFunc(mb.Creator)
	return
}

func (mb *ModelBuilder) WithEditingBuilders(do func(e *EditingBuilder)) *ModelBuilder {
	if mb.creating != nil {
		do(mb.creating)
	}
	do(mb.editing)
	return mb
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

func (b *EditingBuilder) MaxPostSize(size int64) *EditingBuilder {
	if size <= 0 {
		size = web.DefaulMaxPostSize
	}
	b.maxPostSize = size
	return b
}

func (b *EditingBuilder) GetMaxPostSize() int64 {
	return b.maxPostSize
}

// string / []string / *FieldsSection
func (b *EditingBuilder) Only(vs ...interface{}) *EditingBuilder {
	b.FieldsBuilder = *b.FieldsBuilder.Only(vs...)
	return b
}

// string / []string / *FieldsSection
func (b *EditingBuilder) Prepend(vs ...interface{}) *EditingBuilder {
	b.FieldsBuilder = *b.FieldsBuilder.Prepend(vs...)
	return b
}

// string / []string / *FieldsSection
func (b *EditingBuilder) Append(vs ...interface{}) *EditingBuilder {
	b.FieldsBuilder = *b.FieldsBuilder.Append(vs...)
	return b
}

func (b *EditingBuilder) Except(vs ...string) *EditingBuilder {
	b.FieldsBuilder = *b.FieldsBuilder.Except(vs...)
	return b
}

func (b *EditingBuilder) FetchFunc(v FetchFunc) *EditingBuilder {
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

func (b *EditingBuilder) CreateFunc(v CreateFunc) (r *EditingBuilder) {
	b.Creator = v
	return b
}

func (b *EditingBuilder) WrapCreateFunc(w func(in CreateFunc) CreateFunc) (r *EditingBuilder) {
	b.Creator = w(b.Creator)
	return b
}

func (b *EditingBuilder) SetterFunc(v SetterFunc) (r *EditingBuilder) {
	b.Setter = v
	return b
}

func (b *EditingBuilder) WrapSetterFunc(w func(in SetterFunc) SetterFunc) (r *EditingBuilder) {
	b.Setter = w(b.Setter)
	return b
}

func (b *EditingBuilder) PostSetterFunc(v SetterFunc) (r *EditingBuilder) {
	b.PostSetter = v
	return b
}

func (b *EditingBuilder) WrapPostSetterFunc(w func(in SetterFunc) SetterFunc) (r *EditingBuilder) {
	if b.PostSetter == nil {
		b.PostSetter = func(any, *web.EventContext) {}
	}
	b.PostSetter = w(b.PostSetter)
	return b
}

func (b *EditingBuilder) PreValidate(f func(ctx *web.EventContext, obj any) (err error)) *EditingBuilder {
	b.preValidate = f
	return b
}

func (b *EditingBuilder) WrapPreValidate(f func(old func(ctx *web.EventContext, obj any) (err error)) func(ctx *web.EventContext, obj any) error) *EditingBuilder {
	return b.PreValidate(f(b.preValidate))
}

func (b *EditingBuilder) PostValidate(f func(ctx *web.EventContext, obj any) (err error)) *EditingBuilder {
	b.postValidate = f
	return b
}

func (b *EditingBuilder) WrapPostValidate(f func(old func(ctx *web.EventContext, obj any) (err error)) func(ctx *web.EventContext, obj any) error) *EditingBuilder {
	return b.PostValidate(f(b.preValidate))
}

func (b *EditingBuilder) PreSaveCallback(f SaveCallbackFunc) *EditingBuilder {
	b.preSaveCallback = f
	return b
}

func (b *EditingBuilder) WrapPreSaveCallback(f func(old SaveCallbackFunc) SaveCallbackFunc) (r *EditingBuilder) {
	return b.PreSaveCallback(f(b.preSaveCallback))
}

func (b *EditingBuilder) PostSaveCallback(f SaveCallbackFunc) *EditingBuilder {
	b.postSaveCallback = f
	return b
}

func (b *EditingBuilder) WrapPostSaveCallback(f func(old SaveCallbackFunc) SaveCallbackFunc) (r *EditingBuilder) {
	return b.PostSaveCallback(f(b.postSaveCallback))
}

func (b *EditingBuilder) WrapNew(f func(ctx *web.EventContext, obj any) (err error)) *EditingBuilder {
	if b.New == nil {
		b.New = f
	} else {
		old := b.New
		b.New = func(ctx *web.EventContext, obj any) (err error) {
			if err = old(ctx, obj); err != nil {
				return
			}
			return f(ctx, obj)
		}
	}
	return b
}

func (b *EditingBuilder) OnChangeActionFunc(v OnChangeActionFunc) (r *EditingBuilder) {
	b.onChangeAction = v
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

func (b *EditingBuilder) FetchAndUnmarshal(opts *FieldsSetterOptions, id ID, removeDeletedAndSort bool, ctx *web.EventContext) (obj interface{}, vErr web.ValidationErrors) {
	obj = b.mb.NewModel()
	if !id.IsZero() || b.mb.singleton {
		err1 := b.Fetcher(obj, id, ctx)
		if err1 != nil {
			if !(err1 == ErrRecordNotFound && b.mb.singleton) {
				vErr.GlobalError(err1.Error())
				// b.UpdateOverlayContent(ctx, &r, obj, "", err1)
				return
			}
		}
	}

	vErr = b.RunSetterFunc(opts, ctx, removeDeletedAndSort, obj)
	return
}

func (b *EditingBuilder) RunSetterFunc(opts *FieldsSetterOptions, ctx *web.EventContext, removeDeletedAndSort bool, toObj interface{}) (vErr web.ValidationErrors) {
	if b.Setter != nil {
		b.Setter(toObj, ctx)
	}

	if err := ctx.R.ParseMultipartForm(b.maxPostSize); err != nil {
		vErr.GlobalError(err.Error())
		return
	}

	vErr = b.Unmarshal(opts, toObj, b.mb.Info(), removeDeletedAndSort, ctx)

	if b.PostSetter != nil {
		b.PostSetter(toObj, ctx)
	}

	return
}

func (b *EditingBuilder) PreComponent(f ModeObjectComponentFunc) *EditingBuilder {
	b.preComponents = append(b.preComponents, f)
	return b
}

func (b *EditingBuilder) PostComponent(f ModeObjectComponentFunc) *EditingBuilder {
	b.postComponents = append(b.postComponents, f)
	return b
}

func (b *EditingBuilder) ToComponent(opts *ToComponentOptions, obj interface{}, mode FieldModeStack, ctx *web.EventContext) h.HTMLComponent {
	var (
		comp h.HTMLComponents
		add  = func(c h.HTMLComponent) {
			if comps, ok := c.(h.HTMLComponents); ok {
				comp = append(comp, comps...)
			} else {
				comp = append(comp, c)
			}
		}
	)

	for _, f := range b.preComponents {
		add(f(mode, obj, ctx))
	}

	add(b.FieldsBuilder.ToComponent(opts, b.mb.Info(), obj, mode, ctx))

	for _, f := range b.postComponents {
		add(f(mode, obj, ctx))
	}
	return comp
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
	f.ScopeDisabled = ctx.R.FormValue(ParamEditFormUnscoped) == "true"

	f.Respond(r)
}

func (b *EditingBuilder) defaultPageFunc(ctx *web.EventContext) (r web.PageResponse, err error) {
	return b.DefaultPageFuncMode(false, ctx)
}

func (b *EditingBuilder) DefaultPageFuncMode(creating bool, ctx *web.EventContext) (r web.PageResponse, err error) {
	if b.mb.editingDisabled {
		err = ErrUpdateRecordNotAllowed
		return
	}

	var (
		msgr = MustGetMessages(ctx.Context())
		obj  = b.mb.NewModel()
		mid  ID
	)

	if creating {
		if b.mb.permissioner.Creator(ctx.R, ParentsModelID(ctx.R)...).Denied() {
			r.Body = h.Div(h.Text(perm.PermissionDenied.Error()))
			return
		}
		r.PageTitle = msgr.CreatingObjectTitle(inflection.Singular(b.mb.label), b.mb.female)
	} else {
		if !b.mb.singleton {
			if mid, err = b.mb.ParseRecordID(ctx.Queries().Get(ParamID)); err != nil {
				return
			}
		}

		if b.mb.permissioner.Updater(ctx.R, mid, ParentsModelID(ctx.R)...).Denied() {
			r.Body = h.Div(h.Text(perm.PermissionDenied.Error()))
			return
		}

		if !mid.IsZero() || b.mb.singleton {
			if err = b.Fetcher(obj, mid, ctx); err != nil {
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

		r.PageTitle = msgr.DetailingObjectTitle(inflection.Singular(b.mb.label), b.mb.RecordTitle(obj, ctx))
	}

	portalName := actions.Edit

	overlay := actions.OverlayMode(ctx.R.FormValue(ParamOverlay))
	if overlay != "" {
		portalName = overlay.PortalName()
	}

	WithScope(ctx, web.Scope())

	// set portal to edit btn
	ctx.R.Form.Set(ParamTargetPortal, portalName)

	if !creating {
		EditFormUnscoped(ctx, true)
	}

	f := b.form(obj, ctx)
	comp := f.Component()
	r.PageTitle = f.Title

	r.Body = VContainer(web.Portal(comp).
		Name(portalName)).Fluid(true)

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

func (b *EditingBuilder) HiddenField(f ...string) *EditingBuilder {
	b.FieldsBuilder.HiddenField(f...)
	return b
}
