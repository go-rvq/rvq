package presets

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/qor5/admin/v3/model"
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/perm"
	. "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
)

type DetailingBuilder struct {
	mb                 *ModelBuilder
	actions            []*ActionBuilder
	pageFunc           web.PageFunc
	fetcher            FetchFunc
	tabPanels          []TabComponentFunc
	afterTitleCompFunc ObjectComponentFunc
	pageHandlers       PageHandlers
	verifiers          perm.PermVerifiers

	SectionsBuilder
	RowMenuFields

	EditingRestrictionField[*DetailingBuilder]
	DeletingRestrictionField[*DetailingBuilder]
}

func NewDetailingBuilder(mb *ModelBuilder, sb SectionsBuilder) *DetailingBuilder {
	d := &DetailingBuilder{mb: mb, SectionsBuilder: sb}
	d.RowMenuFields.init(mb)
	d.EditingRestriction = NewObjRestriction(d, func(r *ObjRestriction[*DetailingBuilder]) {
		r.Insert(mb.EditingRestriction)
	})
	d.DeletingRestriction = NewObjRestriction(d, func(r *ObjRestriction[*DetailingBuilder]) {
		r.Insert(mb.DetailingRestriction)
	})
	return d
}

func (mb *ModelBuilder) newDetailing() (r *DetailingBuilder) {
	mb.detailing = NewDetailingBuilder(mb, SectionsBuilder{
		mb:            mb,
		FieldsBuilder: *mb.NewFieldsBuilder(mb.detailFieldBuilders.HasMode(DETAIL)...),
	})

	mb.detailing.FetchFunc(mb.Fetcher)
	return
}

// string / []string / *FieldsSection
func (mb *ModelBuilder) Detailing(vs ...interface{}) (r *DetailingBuilder) {
	r = mb.detailing

	if !mb.hasDetailing {
		if len(vs) == 0 {
			// put audited fields to end

			var end []any
			for _, f := range r.fields {
				if f.audited {
					end = append(end, f.name)
				} else {
					vs = append(vs, f.name)
				}
			}

			vs = append(vs, end...)
		}

		rmb := r.RowMenu()

		rmb.SetRowMenuItem("Delete").ComponentFunc(
			NewDeletingMenuItemBuilder(mb.Info()).
				SetWrapEvent(func(rctx *RecordMenuItemContext, e *web.VueEventTagBuilder) {
					cb := web.DecodeCallback(rctx.Ctx.R.FormValue(ParamPostChangeCallback))
					mode := GetOverlay(rctx.Ctx)
					if mode.Overlayed() {
						cb.AddScript("closer.show = false")
					} else {
						cb.AddScript(web.Plaid().PushState(true).Location(web.Location(nil).URL(mb.modelInfo.ListingHrefCtx(rctx.Ctx))).Go())
					}
					e.ValidQuery(ParamPostChangeCallback, cb.Encode())
				}).
				Build(),
		)
	}

	mb.hasDetailing = true

	if len(vs) == 0 {
		return
	}

	r.Only(vs...)
	return r
}

func (mb *ModelBuilder) DisableDetailing() *ModelBuilder {
	mb.hasDetailing = false
	return mb
}

func (mb *ModelBuilder) SetDetailingBuilder(dt *DetailingBuilder) *ModelBuilder {
	mb.detailing = dt
	return mb
}

func (b DetailingBuilder) Clone() *DetailingBuilder {
	b.EditingRestriction = b.EditingRestriction.Clone(&b)
	b.DeletingRestriction = b.DeletingRestriction.Clone(&b)
	return &b
}

func (b *DetailingBuilder) ModelBuilder() *ModelBuilder {
	return b.mb
}

// string / []string / *FieldsSection
func (b *DetailingBuilder) Only(vs ...interface{}) (r *DetailingBuilder) {
	r = b
	r.FieldsBuilder = *r.FieldsBuilder.Only(vs...)
	return
}

func (b *DetailingBuilder) Prepend(vs ...interface{}) (r *DetailingBuilder) {
	r = b
	r.FieldsBuilder = *r.FieldsBuilder.Prepend(vs...)
	return
}

func (b *DetailingBuilder) Except(vs ...string) (r *DetailingBuilder) {
	r = b
	r.FieldsBuilder = *r.FieldsBuilder.Except(vs...)
	return
}

func (b *DetailingBuilder) PageFunc(pf web.PageFunc) (r *DetailingBuilder) {
	b.pageFunc = pf
	return b
}

func (b *DetailingBuilder) FetchFunc(v FetchFunc) (r *DetailingBuilder) {
	b.fetcher = v
	return b
}

func (b *DetailingBuilder) WrapFetchFunc(w func(in FetchFunc) FetchFunc) (r *DetailingBuilder) {
	b.fetcher = w(b.fetcher)
	return b
}

func (b *DetailingBuilder) GetFetchFunc() FetchFunc {
	return b.fetcher
}

func (b *DetailingBuilder) AfterTitleCompFunc(v ObjectComponentFunc) (r *DetailingBuilder) {
	if v == nil {
		panic("value required")
	}
	b.afterTitleCompFunc = v
	return b
}

func (b *DetailingBuilder) GetPageFunc() web.PageFunc {
	if b.pageFunc != nil {
		return b.pageFunc
	}
	return b.defaultPageFunc
}

func (b *DetailingBuilder) AppendTabsPanelFunc(v TabComponentFunc) (r *DetailingBuilder) {
	b.tabPanels = append(b.tabPanels, v)
	return b
}

func (b *DetailingBuilder) TabsPanelFunc() (r []TabComponentFunc) {
	return b.tabPanels
}

func (b *DetailingBuilder) TabsPanels(vs ...TabComponentFunc) (r *DetailingBuilder) {
	b.tabPanels = vs
	return b
}

func (b *DetailingBuilder) Verifier(vf ...*perm.PermVerifierBuilder) (r *DetailingBuilder) {
	b.verifiers.Add(vf...)
	return b
}

func (b *DetailingBuilder) GetVerifiers() perm.PermVerifiers {
	return b.verifiers
}

func (b *DetailingBuilder) defaultPageFunc(ctx *web.EventContext) (r web.PageResponse, err error) {
	if b.mb.deletingDisabled {
		err = ErrReadRecordNotAllowed
		return
	}

	id := ctx.Param(ParamID)

	var (
		mid  ID
		obj  = b.mb.NewModel()
		msgr = MustGetMessages(ctx.Context())
	)

	if !b.mb.singleton {
		if id == "" {
			err = msgr.ErrEmptyParamID
			return
		}

		r.Body = VContainer(h.Text(id))

		if mid, err = b.mb.ParseRecordID(id); err != nil {
			return
		}
	}

	if b.mb.permissioner.Reader(ctx.R, mid, ParentsModelID(ctx.R)...).Denied() {
		r.Body = h.Div(h.Text(MustGetMessages(ctx.Context()).ErrPermissionDenied.Error()))
		return
	}

	err = b.GetFetchFunc()(obj, mid, ctx)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			if b.mb.singleton {
				return b.mb.editing.CreatingBuilder().DefaultPageFuncMode(true, ctx)
			}
			return b.mb.p.DefaultNotFoundPageFunc(ctx)
		}
		return
	}

	if b.mb.singleton {
		r.PageTitle = b.mb.TTitle(ctx.Context())
	} else {
		r.PageTitle = b.mb.RecordTitle(obj, ctx)
	}

	form := NewFormBuilder(ctx, b.mb, &b.FieldsBuilder, obj)
	form.mode = DETAIL
	f := b.configureForm(form.Build())

	if len(f.MainPortals) > 0 {
		AddPortals(ctx, f.MainPortals...)
	}

	f.MainPortals = nil

	r.Body = f.Component()

	return
}

func (b *DetailingBuilder) BuildPage(vf *perm.PermVerifierBuilder, builder func(ctx *web.EventContext, obj any, mid model.ID, r *web.PageResponse) (err error)) func(ctx *web.EventContext) (r web.PageResponse, err error) {
	if vf == nil {
		vf = perm.PermVerifier()
	}
	return b.mb.BindPageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		var mid ID

		obj := b.mb.NewModel()
		msgr := MustGetMessages(ctx.Context())

		if !b.mb.singleton {
			id := ctx.Param(ParamID)
			if id == "" {
				err = msgr.ErrEmptyParamID
				return
			}

			if mid, err = b.mb.ParseRecordID(id); err != nil {
				return
			}
		}

		v := vf.Build(b.mb.permissioner.Reader(ctx.R, mid, ParentsModelID(ctx.R)...))

		if v.Denied() {
			r.Body = h.Div(h.Text(MustGetMessages(ctx.Context()).ErrPermissionDenied.Error()))
			return
		}

		err = b.GetFetchFunc()(obj, mid, ctx)
		if err != nil {
			if errors.Is(err, ErrRecordNotFound) {
				return b.mb.p.DefaultNotFoundPageFunc(ctx)
			}
			return
		}

		r.PageTitle = b.mb.RecordTitle(obj, ctx)

		err = builder(ctx, obj, mid, &r)
		return
	})
}

func (d *DetailingBuilder) AddPageFunc(vf *perm.PermVerifierBuilder, pth string, handler func(ctx *web.EventContext, obj any, mid model.ID, r *web.PageResponse) (err error), methods ...string) (ph *PageHandler) {
	if vf != nil {
		if !vf.Valid() {
			vf.Path(pth)
		}
		d.verifiers.Add(vf)
	}
	ph = NewPageHandler(pth, d.mb.p.Wrap(d.mb.p.GetDetailLayoutFunc()(d.BuildPage(vf, handler), d.mb.GetLayoutConfig())), methods...)
	d.pageHandlers.Add(ph)
	return
}

func (d *DetailingBuilder) AddRawPageFunc(path string, f web.PageFunc, methods ...string) (ph *PageHandler) {
	ph = NewPageHandler(path, d.mb.p.Wrap(f), methods...)
	d.pageHandlers.Add(ph)
	return
}

func (b *DetailingBuilder) detailingEvent(ctx *web.EventContext) (r web.EventResponse, err error) {
	if b.mb.deletingDisabled {
		err = ErrReadRecordNotAllowed
		return
	}

	id := ctx.Param(ParamID)
	obj := b.mb.NewModel()
	msgr := MustGetMessages(ctx.Context())
	targetPortal := ctx.R.FormValue(ParamTargetPortal)

	if id == "" && !b.mb.singleton {
		err = msgr.ErrEmptyParamID
		return
	}

	var mid ID
	if mid, err = b.mb.ParseRecordID(id); err != nil {
		return
	}

	if b.mb.permissioner.Reader(ctx.R, mid, ParentsModelID(ctx.R)...).Denied() {
		err = perm.PermissionDenied
		return
	}

	if err = b.GetFetchFunc()(obj, mid, ctx); err != nil {
		return
	}

	form := NewFormBuilder(ctx, b.mb, &b.FieldsBuilder, obj)
	form.mode = DETAIL
	f := b.configureForm(form.Build()).Component()

	mode := GetOverlay(ctx)
	if mode.IsDrawer() {
		b.mb.p.Drawer(mode).
			SetValidPortalName(targetPortal).
			SetScrollable(true).
			Respond(&r, f)
	} else {
		b.mb.p.Dialog().
			SetScrollable(true).
			SetTargetPortal(targetPortal).
			Respond(ctx, &r, f)
	}
	return
}

func (b *DetailingBuilder) EventBuilder(builder func(ctx *web.EventContext, obj any, mid model.ID, r *web.EventResponse) (comp h.HTMLComponent, err error)) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id := ctx.Param(ParamID)
		obj := b.mb.NewModel()
		msgr := MustGetMessages(ctx.Context())

		if id == "" {
			err = msgr.ErrEmptyParamID
			return
		}
		var mid ID
		if mid, err = b.mb.ParseRecordID(id); err != nil {
			return
		}

		if b.mb.permissioner.Reader(ctx.R, mid, ParentsModelID(ctx.R)...).Denied() {
			err = perm.PermissionDenied
			return
		}

		if err = b.GetFetchFunc()(obj, mid, ctx); err != nil {
			return
		}

		var comp h.HTMLComponent
		if comp, err = builder(ctx, obj, mid, &r); err != nil {
			return
		}

		if comp == nil {
			return
		}

		targetPortal := ctx.R.FormValue(ParamTargetPortal)

		if targetPortal == "" {
			r.Body = comp
			return
		}

		r.UpdatePortal(targetPortal, comp)

		return
	}

}

func (b *DetailingBuilder) configureForm(f *Form) *Form {
	var (
		ctx        = f.b.ctx
		obj        = f.Obj
		portalName = ctx.R.FormValue(ParamTargetPortal)
	)

	if b.afterTitleCompFunc != nil {
		ctx.WithContextValue(ctxDetailingAfterTitleComponent, b.afterTitleCompFunc(obj, ctx))
	}

	if msg, ok := ctx.Flash.(string); ok {
		f.Notice = VSnackbar(h.Text(msg)).ModelValue(true).Location("top").Color("success")
	}

	f.Portal = portalName

	if b.EditingRestriction.CanObj(obj, ctx) {
		var cb web.Callback
		cb.Decode(ctx.R.FormValue(ParamPostChangeCallback))

		overlayMode := f.b.overlayMode

		if overlayMode.Overlayed() {
			cb.AddScript(web.Plaid().
				URL(ctx.R.RequestURI).
				EventFunc(actions.Detailing).
				StringQuery(ctx.Queries().Encode()).
				Go())
		} else {
			cb.AddScript(web.Plaid().
				URL(ctx.R.RequestURI).
				StringQuery(ctx.Queries().Encode()).
				Go())
		}

		editMode := overlayMode.Up()
		editingPortalID := ctx.UID()
		editPortal := formPortalName + editingPortalID

		onclick := web.Plaid().
			URL(b.mb.Info().ListingHrefCtx(ctx)).
			EventFunc(actions.Edit).
			Query(ParamID, f.b.id).
			ValidQuery(ParamTargetPortal, editPortal).
			ValidQuery(ParamOverlay, editMode.String()).
			ValidQuery(ParamPostChangeCallback, cb.String())

		f.MainPortals = append(f.MainPortals,
			web.Portal().
				Name(editPortal))

		f.PrimaryAction = h.HTMLComponents{
			VBtn("").
				Variant(VariantFlat).
				Color("primary").
				Attr(":disabled", "isFetching").
				Attr(":loading", "isFetching").
				Attr("@click", onclick.Go()).
				Attr("@click.middle",
					fmt.Sprintf(`(e) => e.view.window.open(%q, "_blank")`, b.mb.Info().EditingHrefCtx(ctx, f.b.id))).
				Icon(true).
				Density("comfortable").
				Children(VIcon("mdi-pencil")),
		}
	}

	f.Tabs = b.tabPanels

	if v, ok := GetComponentFromContext(ctx, ctxDetailingAfterTitleComponent); ok {
		f.TopRightActions = append(f.TopRightActions, v)
	}

	f.Title = MustGetMessages(ctx.Context()).DetailingObjectTitle(
		b.mb.TTitle(ctx.Context()),
		b.mb.RecordTitle(obj, ctx))
	sharedPortal := ctx.UID()
	f.MainPortals = append(f.MainPortals, web.Portal().Name(sharedPortal))

	var menus h.HTMLComponents
	b.RowMenu().listingItemFuncs(ctx).
		ForEachRowMenuItemFunc(sharedPortal, func(rctx *RecordMenuItemContext, name string) string {
			name = sharedPortal + "--" + name
			f.MainPortals = append(f.MainPortals, web.Portal().Name(name))
			return name
		}, func(i int, m vx.RowMenuItemFunc) {
			c := m(0, f.Obj, f.b.id, ctx)
			if c != nil {
				menus = append(menus, c)
			}
		})

	actionsMenus, actionsErrors := BuildMenuItemCompomentsOfActions(sharedPortal, ctx, f.b.mb, f.b.id, obj, b.actions...)

	if len(actionsErrors) > 0 {
		f.Body = append(actionsErrors, f.Body)
	}

	if len(actionsMenus) > 0 {
		menus = append(menus, VDivider())
		for _, menuItem := range actionsMenus {
			menus = append(menus, menuItem)
		}
	}

	if len(menus) > 0 {
		f.Menu = append(
			f.Menu,
			VList(menus...).Class("record-menu").
				OpenStrategy("single").
				Class("primary--text").
				Density(DensityCompact),
		)
	}

	return f
}

func (b *DetailingBuilder) doAction(ctx *web.EventContext) (r web.EventResponse, err error) {
	var (
		action *ActionBuilder
		id     string
	)
	if id, action, err = b.parseRequestAction(ctx); err != nil {
		return
	}

	var mid ID
	if !b.mb.singleton {
		if mid, err = b.mb.ParseRecordID(id); err != nil {
			return
		}
	}

	if b.mb.permissioner.ObjectReadActioner(ctx.R, action.PermName(), mid, ParentsModelID(ctx.R)...).Denied() {
		err = perm.PermissionDenied
		return
	}

	_, err = action.Do(b.mb, id, ctx, &r)
	return
}

func (b *DetailingBuilder) formAction(ctx *web.EventContext) (r web.EventResponse, err error) {
	var (
		action *ActionBuilder
		id     string
	)
	if id, action, err = b.parseRequestAction(ctx); err != nil {
		return
	}

	var mid ID
	if !b.mb.singleton {
		if mid, err = b.mb.ParseRecordID(id); err != nil {
			return
		}
	}

	if b.mb.permissioner.ObjectReadActioner(ctx.R, action.PermName(), mid, ParentsModelID(ctx.R)...).Denied() {
		err = perm.PermissionDenied
		return
	}

	err = action.View(b.mb, id, ctx, &r)
	return
}

func (b *DetailingBuilder) Fetch(id string, ctx *web.EventContext) (obj any, err error) {
	obj = b.mb.NewModel()
	var mid ID
	if !b.mb.singleton {
		if mid, err = b.mb.ParseRecordID(id); err != nil {
			return
		}
	}

	if err = b.GetFetchFunc()(obj, mid, ctx); err != nil {
		return
	}
	return
}

func (b *DetailingBuilder) parseRequestAction(ctx *web.EventContext) (id string, action *ActionBuilder, err error) {
	action = getAction(b.actions, ctx.R.FormValue(ParamAction))
	if action == nil {
		err = errors.New("action required")
		return
	}

	if !b.mb.singleton {
		if id = ctx.R.FormValue(ParamID); id == "" {
			err = errors.New("id required")
		}
		return
	}

	var enabled bool
	if enabled, err = action.IsEnabled(id, ctx); err != nil {
		return
	}
	if !enabled {
		err = errors.New("action disabled")
		return
	}
	return
}

// EditDetailField EventFunc: click detail field component edit button
func (b *DetailingBuilder) EditDetailField(ctx *web.EventContext) (r web.EventResponse, err error) {
	key := ctx.Queries().Get(SectionFieldName)

	f := b.Section(key)

	obj := b.mb.NewModel()
	var mid ID
	if mid, err = b.mb.ParseRecordID(ctx.Queries().Get(ParamID)); err != nil {
		return
	}

	err = b.GetFetchFunc()(obj, mid, ctx)
	if err != nil {
		return
	}
	if f.setter != nil {
		f.setter(obj, ctx)
	}

	r.UpdatePortal(f.FieldPortalName(), f.editComponent(obj, &FieldContext{
		ToComponentOptions: &ToComponentOptions{},
		EventContext:       ctx,
		FormKey:            f.name,
		Path:               FieldPath{f.name},
		Name:               f.name,
	}, ctx))
	return r, nil
}

// SaveDetailField EventFunc: click save button
func (b *DetailingBuilder) SaveDetailField(ctx *web.EventContext) (r web.EventResponse, err error) {
	key := ctx.Queries().Get(SectionFieldName)

	f := b.Section(key)

	obj := b.mb.NewModel()

	var mid ID
	if mid, err = b.mb.ParseRecordID(ctx.Queries().Get(ParamID)); err != nil {
		return
	}

	err = b.GetFetchFunc()(obj, mid, ctx)
	if err != nil {
		return
	}
	if f.setter != nil {
		f.setter(obj, ctx)
	}

	err = f.saver(obj, mid, ctx)
	if err != nil {
		ShowMessage(&r, err.Error(), "warning")
		return r, nil
	}

	r.UpdatePortal(
		f.FieldPortalName(),
		f.viewComponent(&FieldContext{
			ToComponentOptions: &ToComponentOptions{},
			EventContext:       ctx,
			Obj:                obj,
			FormKey:            f.name,
			Path:               FieldPath{f.name},
			Name:               f.name,
		}, ctx))
	return r, nil
}

// EditDetailListField Event: click detail list field element edit button
func (b *DetailingBuilder) EditDetailListField(ctx *web.EventContext) (r web.EventResponse, err error) {
	var (
		fieldName          string
		index, deleteIndex int64
	)

	fieldName = ctx.Queries().Get(SectionFieldName)
	f := b.Section(fieldName)

	index, err = strconv.ParseInt(ctx.Queries().Get(f.EditBtnKey()), 10, 64)
	if err != nil {
		return
	}
	deleteIndex = -1
	if ctx.Queries().Get(f.DeleteBtnKey()) != "" {
		deleteIndex, err = strconv.ParseInt(ctx.Queries().Get(f.EditBtnKey()), 10, 64)
		if err != nil {
			return
		}
	}

	var mid ID
	if mid, err = b.mb.ParseRecordID(ctx.Queries().Get(ParamID)); err != nil {
		return
	}
	obj := b.mb.NewModel()
	err = b.GetFetchFunc()(obj, mid, ctx)
	if err != nil {
		return
	}
	if f.setter != nil {
		f.setter(obj, ctx)
	}

	r.UpdatePortal(
		f.FieldPortalName(),
		f.listComponent(&FieldContext{
			ToComponentOptions: &ToComponentOptions{},
			EventContext:       ctx,
			Obj:                obj,
			Mode:               FieldModeStack{DETAIL},
		}, ctx, int(deleteIndex), int(index), -1))

	return
}

// SaveDetailListField Event: click detail list field element Save button
func (b *DetailingBuilder) SaveDetailListField(ctx *web.EventContext) (r web.EventResponse, err error) {
	var (
		fieldName string
		index     int64
	)

	fieldName = ctx.Queries().Get(SectionFieldName)
	f := b.Section(fieldName)

	index, err = strconv.ParseInt(ctx.Queries().Get(f.SaveBtnKey()), 10, 64)
	if err != nil {
		return
	}

	var mid ID
	if mid, err = b.mb.ParseRecordID(ctx.Queries().Get(ParamID)); err != nil {
		return
	}
	obj := b.mb.NewModel()
	err = b.GetFetchFunc()(obj, mid, ctx)
	if err != nil {
		return
	}
	if f.setter != nil {
		f.setter(obj, ctx)
	}

	err = f.saver(obj, mid, ctx)
	if err != nil {
		ShowMessage(&r, err.Error(), "warning")
		return r, nil
	}

	r.UpdatePortal(
		f.FieldPortalName(),
		f.listComponent(&FieldContext{
			ToComponentOptions: &ToComponentOptions{},
			EventContext:       ctx,
			Obj:                obj,
			Mode:               FieldModeStack{DETAIL},
		}, ctx, -1, -1, int(index)))

	return
}

// DeleteDetailListField Event: click detail list field element Delete button
func (b *DetailingBuilder) DeleteDetailListField(ctx *web.EventContext) (r web.EventResponse, err error) {
	var (
		fieldName string
		index     int64
	)

	fieldName = ctx.Queries().Get(SectionFieldName)
	f := b.Section(fieldName)

	index, err = strconv.ParseInt(ctx.Queries().Get(f.DeleteBtnKey()), 10, 64)
	if err != nil {
		return
	}

	var mid ID
	if mid, err = b.mb.ParseRecordID(ctx.Queries().Get(ParamID)); err != nil {
		return
	}

	obj := b.mb.NewModel()
	err = b.GetFetchFunc()(obj, mid, ctx)
	if err != nil {
		return
	}
	if f.setter != nil {
		f.setter(obj, ctx)
	}

	// delete from slice
	var list any
	if list, err = reflectutils.Get(obj, f.name); err != nil {
		return
	}
	listValue := reflect.ValueOf(list)
	if listValue.Kind() != reflect.Slice {
		err = errors.New("field is not a slice")
		return
	}
	newList := reflect.MakeSlice(reflect.TypeOf(list), 0, 0)
	for i := 0; i < listValue.Len(); i++ {
		if i != int(index) {
			newList = reflect.Append(newList, listValue.Index(i))
		}
	}
	if err = reflectutils.Set(obj, f.name, newList.Interface()); err != nil {
		return
	}

	err = f.saver(obj, mid, ctx)
	if err != nil {
		ShowMessage(&r, err.Error(), "warning")
		return r, nil
	}

	r.UpdatePortal(f.FieldPortalName(),
		f.listComponent(&FieldContext{
			ToComponentOptions: &ToComponentOptions{},
			EventContext:       ctx,
			Obj:                obj,
			Mode:               FieldModeStack{DETAIL},
		}, ctx, int(index), -1, -1))

	return
}

// CreateDetailListField Event: click detail list field element Add row button
func (b *DetailingBuilder) CreateDetailListField(ctx *web.EventContext) (r web.EventResponse, err error) {
	fieldName := ctx.Queries().Get(SectionFieldName)
	f := b.Section(fieldName)

	var mid ID
	if mid, err = b.mb.ParseRecordID(ctx.Queries().Get(ParamID)); err != nil {
		return
	}

	obj := b.mb.NewModel()
	err = b.GetFetchFunc()(obj, mid, ctx)
	if err != nil {
		return
	}
	if f.setter != nil {
		f.setter(obj, ctx)
	}

	var list any
	if list, err = reflectutils.Get(obj, f.name); err != nil {
		return
	}

	listLen := 0
	if list != nil {
		listValue := reflect.ValueOf(list)
		if listValue.Kind() != reflect.Slice {
			err = errors.New(fmt.Sprintf("the kind of list field is %s, not slice", listValue.Kind()))
			return
		}
		listLen = listValue.Len()
	}

	if err = reflectutils.Set(obj, f.name+"[]", f.editingFB.model); err != nil {
		return
	}

	if err = f.saver(obj, mid, ctx); err != nil {
		ShowMessage(&r, err.Error(), "warning")
		return r, nil
	}

	r.UpdatePortal(f.FieldPortalName(),
		f.listComponent(&FieldContext{
			ToComponentOptions: &ToComponentOptions{},
			EventContext:       ctx,
			Obj:                obj,
			Mode:               FieldModeStack{DETAIL},
		}, ctx, -1, listLen, -1))

	return
}

func (b *DetailingBuilder) HiddenField(f ...string) *DetailingBuilder {
	b.FieldsBuilder.HiddenField(f...)
	return b
}
