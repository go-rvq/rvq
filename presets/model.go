package presets

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"github.com/qor5/admin/v3/model"
	"github.com/qor5/admin/v3/reflect_utils"
	"github.com/qor5/web/v3"
	"github.com/qor5/web/v3/datafield"
	"github.com/qor5/x/v3/i18n"
	"github.com/qor5/x/v3/perm"
	"github.com/sunfmin/reflectutils"
)

type ModelBuilderOptionFunc func(mb *ModelBuilder)

func ModelWithID(id string) ModelBuilderOptionFunc {
	return func(mb *ModelBuilder) {
		mb.id = id
	}
}

func (f ModelBuilderOptionFunc) Apply(mb *ModelBuilder) {
	f(mb)
}

type ModelBuilderOption interface {
	Apply(mb *ModelBuilder)
}

type ModelBuilder struct {
	parent *ModelBuilder
	model  interface{}

	ModelBuilderConfigAttributes

	p *Builder
	writeFieldBuilders,
	listFieldBuilders,
	detailFieldBuilders FieldBuilders

	modelType           reflect.Type
	menuGroupName       string
	notInMenu           bool
	menuIcon            string
	defaultURLQueryFunc func(*http.Request) url.Values
	fieldLabels         map[string]func(ctx *web.EventContext) string
	fieldHints          map[string]func(ctx *web.EventContext) string
	placeholders        []string
	listing             *ListingBuilder
	detailing           *DetailingBuilder
	editing             *EditingBuilder
	creating            *EditingBuilder
	editingDisabled     bool
	deletingDisabled    bool
	detailingDisabled   bool
	creatingDisabled    bool
	hasDetailing        bool
	rightDrawerWidth    string
	link                string
	layoutConfig        *LayoutConfig
	modelInfo           *ModelInfo
	permissioner        *ModelPermissioner
	plugins             []ModelPlugin
	subRoutesSetup      func(mux *http.ServeMux, uri string)
	listingMenu         Menu
	detailingMenu       Menu

	EventsHub web.EventsHub

	children         []*ModelBuilder
	saveURI          string
	routeSetuper     func(mux *http.ServeMux, uri string)
	itemRouteSetuper []func(mux *http.ServeMux, uri string)

	pages PageHandlers

	verifierModel *ModelBuilder
	verifiers     perm.PermVerifiers

	detailingParentsIDResolver   ParentsModelIDResolver
	BeforeFormUnmarshallHandlers ModelFormUnmarshallHandlers
	PostFormUnmarshallHandlers   ModelFormUnmarshallHandlers

	datafield.DataField[*ModelBuilder]

	ListingRestrictionField[*ModelBuilder]
	CreatingRestrictionField[*ModelBuilder]
	EditingRestrictionField[*ModelBuilder]
	DetailingRestrictionField[*ModelBuilder]
	DeletingRestrictionField[*ModelBuilder]
}

func NewModelBuilder(p *Builder, model interface{}, options ...ModelBuilderOption) (mb *ModelBuilder) {
	mb = datafield.New(&ModelBuilder{
		p:     p,
		model: model,
	})

	mb.modelType = reflect.TypeOf(model)
	if mb.modelType.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("model %#+v must be pointer", model))
	}

	modelstr := mb.modelType.String()
	modelName := modelstr[strings.LastIndex(modelstr, ".")+1:]

	for _, option := range options {
		option.Apply(mb)
	}

	if mb.label == "" {
		mb.label = strcase.ToCamel(modelName)
	}

	if mb.pluralLabel == "" {
		mb.pluralLabel = strcase.ToCamel(inflection.Plural(modelName))
	}

	if mb.id == "" {
		if mb.singleton {
			mb.id = strcase.ToSnake(mb.label)
		} else {
			mb.id = strcase.ToSnake(mb.pluralLabel)
		}
	}

	if mb.uriName == "" {
		mb.uriName = strcase.ToKebab(mb.id)
	}

	mb.modelInfo = &ModelInfo{mb: mb}
	mb.permissioner = &ModelPermissioner{mb: mb}

	mb.EventsHub.Wraper(mb.BindEventFunc)

	setupers := p.ModelSetupFactories.Of(mb)
	setupers.Init()
	_, fields := reflect_utils.UniqueFieldsOfReflectType(mb.modelType)

	fieldBuilders := func(def *FieldDefaults) FieldBuilders {
		fb := def.NewFieldBuilders(fields)
		setupers.InitFields(&fb)
		return def.SetupFields(fb, setupers)
	}

	mb.listFieldBuilders = fieldBuilders(p.listFieldDefaults)
	mb.detailFieldBuilders = fieldBuilders(p.detailFieldDefaults)
	mb.writeFieldBuilders = fieldBuilders(p.writeFieldDefaults)

	mb.ListingRestriction = NewRestriction(mb, func(r *Restriction[*ModelBuilder]) {
		r.Handler(OkHandlerFunc(func(ctx *web.EventContext) (_, _ bool) {
			return r.dot.permissioner.ReqLister(ctx.R).Denied(), true
		}))
	})
	mb.CreatingRestriction = NewRestriction(mb, func(r *Restriction[*ModelBuilder]) {
		r.Handler(OkHandlerFunc(func(ctx *web.EventContext) (_, _ bool) {
			return r.dot.permissioner.ReqCreator(ctx.R).Denied(), true
		}))
	})
	mb.EditingRestriction = NewObjRestriction(mb, func(r *ObjRestriction[*ModelBuilder]) {
		r.ObjHandler(OkObjHandlerFunc(func(obj any, ctx *web.EventContext) (_, _ bool) {
			return r.dot.permissioner.ReqObjectUpdater(ctx.R, obj).Denied(), true
		}))
	})
	mb.DetailingRestriction = NewObjRestriction(mb, func(r *ObjRestriction[*ModelBuilder]) {
		r.ObjHandler(OkObjHandlerFunc(func(obj any, ctx *web.EventContext) (_, _ bool) {
			return r.dot.permissioner.ReqObjectReader(ctx.R, obj).Denied(), true
		}))
	})

	mb.DeletingRestriction = NewObjRestriction(mb, func(r *ObjRestriction[*ModelBuilder]) {
		r.ObjHandler(OkObjHandlerFunc(func(obj any, ctx *web.EventContext) (_, _ bool) {
			return r.dot.permissioner.ReqObjectDeleter(ctx.R, obj).Denied(), true
		}))
	})

	// Be aware the uriName here is still the original struct
	mb.newListing()
	mb.newEditing()
	mb.newDetailing()

	return
}

func (mb *ModelBuilder) String() string {
	var flags []string
	if mb.singleton {
		flags = append(flags, "S")
	} else {
		flags = append(flags, "*"+mb.pluralLabel)
	}
	flags = append(flags, mb.uriName)

	return fmt.Sprintf("%s [%s]", mb.label, strings.Join(flags, ", "))
}

func (mb *ModelBuilder) I18nModuleKeyOrDefault() i18n.ModuleKey {
	if mb.moduleKey == "" {
		return ModelsI18nModuleKey
	}
	return mb.moduleKey
}

func (mb *ModelBuilder) WriteFieldBuilders() FieldBuilders {
	return mb.writeFieldBuilders
}

func (mb *ModelBuilder) ListFieldBuilders() FieldBuilders {
	return mb.listFieldBuilders
}

func (mb *ModelBuilder) DetailFieldBuilders() FieldBuilders {
	return mb.detailFieldBuilders
}

func (mb *ModelBuilder) Female() bool {
	return mb.female
}

func (mb *ModelBuilder) SetFemale(v bool) *ModelBuilder {
	mb.female = v
	return mb
}

func (mb *ModelBuilder) SetPlural(v bool) *ModelBuilder {
	mb.plural = v
	return mb
}

func (mb *ModelBuilder) UriName() string {
	return mb.uriName
}

func (mb *ModelBuilder) SetUriName(uriName string) *ModelBuilder {
	mb.uriName = uriName
	return mb
}

func (mb *ModelBuilder) SetLabel(label string) *ModelBuilder {
	mb.label = label
	return mb
}

func (mb *ModelBuilder) SetPluralLabel(label string) *ModelBuilder {
	mb.pluralLabel = label
	return mb
}

func (mb *ModelBuilder) ModelType() reflect.Type {
	return mb.modelType
}

func (mb *ModelBuilder) VerifierModel() *ModelBuilder {
	return mb.verifierModel
}

func (mb *ModelBuilder) SetVerifierModel(verifier *ModelBuilder) *ModelBuilder {
	mb.verifierModel = verifier
	return mb
}

func (mb *ModelBuilder) GetVerifierModel() *ModelBuilder {
	if mb.verifierModel != nil {
		return mb.verifierModel
	}
	return mb
}

func (mb *ModelBuilder) Depth() (i int) {
	for mb.parent != nil {
		i++
		mb = mb.parent
	}
	return
}

func (mb *ModelBuilder) Parent() *ModelBuilder {
	return mb.parent
}

func (mb *ModelBuilder) Root() *ModelBuilder {
	for mb.parent != nil {
		mb = mb.parent
	}
	return mb
}

func (mb *ModelBuilder) Parents() (parents []*ModelBuilder) {
	for mb.parent != nil {
		parents = append(parents, mb.parent)
		mb = mb.parent
	}
	for i, j := 0, len(parents)-1; i < j; i, j = i+1, j-1 {
		parents[i], parents[j] = parents[j], parents[i]
	}
	return
}

func (mb *ModelBuilder) DataOperator() DataOperator {
	return mb.dataOperator
}

func (mb *ModelBuilder) SetDataOperator(dataOperator DataOperator) *ModelBuilder {
	mb.dataOperator = dataOperator
	return mb
}

func (mb *ModelBuilder) CurrentDataOperator() DataOperator {
	if mb.dataOperator != nil {
		return mb.dataOperator
	}
	return mb.p.dataOperator
}

func (mb *ModelBuilder) Schema() (s Schema) {
	var (
		err error
		do  = mb.CurrentDataOperator()
	)
	if do == nil {
		return
	}
	if s, err = do.Schema(mb.model); err != nil {
		panic(err)
	}
	return s
}

func (mb *ModelBuilder) WithDataOperator(h func(do DataOperator)) bool {
	if do := mb.CurrentDataOperator(); do != nil {
		h(do)
		return true
	}
	return false
}

func (mb *ModelBuilder) UpdateDataOperator(h func(do DataOperator) DataOperator) *ModelBuilder {
	if mb.dataOperator != nil {
		mb.dataOperator = h(mb.dataOperator)
	} else {
		mb.dataOperator = h(mb.p.dataOperator.CloneDataOperator())
	}
	return mb
}

func (mb *ModelBuilder) SetSaveURI(updatingURI string) {
	mb.saveURI = updatingURI
}

func (mb *ModelBuilder) Builder() *Builder {
	return mb.p
}

func (mb *ModelBuilder) HasDetailing() bool {
	return mb.hasDetailing
}

func (mb *ModelBuilder) GetSingleton() bool {
	return mb.singleton
}

func (mb *ModelBuilder) RightDrawerWidth(v string) *ModelBuilder {
	mb.rightDrawerWidth = v
	return mb
}

func (mb *ModelBuilder) Link(v string) *ModelBuilder {
	mb.link = v
	return mb
}

func (mb *ModelBuilder) Searcher(model interface{}, params *SearchParams, ctx *web.EventContext) (r interface{}, totalCount int, err error) {
	return mb.CurrentDataOperator().Search(model, params, ctx)
}

func (mb *ModelBuilder) Fetcher(obj interface{}, id ID, ctx *web.EventContext) (err error) {
	return mb.CurrentDataOperator().Fetch(obj, id, ctx)
}

func (mb *ModelBuilder) TitleFetcher(obj interface{}, id ID, ctx *web.EventContext) (err error) {
	return mb.CurrentDataOperator().FetchTitle(obj, id, ctx)
}

func (mb *ModelBuilder) Saver(obj interface{}, id ID, ctx *web.EventContext) (err error) {
	return mb.CurrentDataOperator().Save(obj, id, ctx)
}

func (mb *ModelBuilder) Creator(obj interface{}, ctx *web.EventContext) (err error) {
	return mb.CurrentDataOperator().Create(obj, ctx)
}

func (mb *ModelBuilder) Deleter(obj interface{}, id ID, cascade bool, ctx *web.EventContext) (err error) {
	return mb.CurrentDataOperator().Delete(obj, id, cascade, ctx)
}

func (mb *ModelBuilder) Model() interface{} {
	return mb.model
}

func (mb *ModelBuilder) NewModel() (r interface{}) {
	return reflect.New(mb.modelType.Elem()).Interface()
}

func (mb *ModelBuilder) NewModelSlice() (r interface{}) {
	return reflect.New(reflect.SliceOf(mb.modelType)).Interface()
}

func (mb *ModelBuilder) NewFieldsBuilder(fields ...*FieldBuilder) *FieldsBuilder {
	return &FieldsBuilder{
		model:    mb.model,
		defaults: mb.p.listFieldDefaults,
		fields:   fields,
	}
}

func (mb *ModelBuilder) Info() (r *ModelInfo) {
	return mb.modelInfo
}

func (mb *ModelBuilder) Permissioner() (r *ModelPermissioner) {
	return mb.modelInfo.Permissioner()
}

func (mb *ModelBuilder) URIName(v string) (r *ModelBuilder) {
	mb.uriName = v
	return mb
}

func (mb *ModelBuilder) DefaultURLQueryFunc(v func(*http.Request) url.Values) (r *ModelBuilder) {
	mb.defaultURLQueryFunc = v
	return mb
}

func (mb *ModelBuilder) InMenu(v bool) (r *ModelBuilder) {
	mb.notInMenu = !v
	return mb
}

func (mb *ModelBuilder) IsInMenu() bool {
	return !mb.notInMenu
}

func (mb *ModelBuilder) MenuIcon(v string) (r *ModelBuilder) {
	mb.menuIcon = v
	return mb
}

func (mb *ModelBuilder) GetMenuIcon() string {
	return mb.menuIcon
}

func (mb *ModelBuilder) Label(v string) (r *ModelBuilder) {
	mb.label = v
	return mb
}

func (mb *ModelBuilder) GetLabel() string {
	return mb.label
}

func (mb *ModelBuilder) Labels(vs ...string) (r *ModelBuilder) {
	if mb.fieldLabels == nil {
		mb.fieldLabels = make(map[string]func(ctx *web.EventContext) string)
	}

	for i := 0; i < len(vs); i = i + 2 {
		name, label := vs[i], vs[i+1]
		mb.fieldLabels[name] = func(*web.EventContext) string {
			return label
		}
	}
	return mb
}

func (mb *ModelBuilder) FieldLabelFunc(name string, f func(ctx *web.EventContext) string) (r *ModelBuilder) {
	if mb.fieldLabels == nil {
		mb.fieldLabels = make(map[string]func(ctx *web.EventContext) string)
	}
	mb.fieldLabels[name] = f
	return mb
}
func (mb *ModelBuilder) Hints(vs ...string) (r *ModelBuilder) {
	if mb.fieldHints == nil {
		mb.fieldHints = make(map[string]func(ctx *web.EventContext) string)
	}

	for i := 0; i < len(vs); i = i + 2 {
		name, label := vs[i], vs[i+1]
		mb.fieldHints[name] = func(*web.EventContext) string {
			return label
		}
	}
	return mb
}

func (mb *ModelBuilder) FieldHintFunc(name string, f func(ctx *web.EventContext) string) (r *ModelBuilder) {
	if mb.fieldHints == nil {
		mb.fieldHints = make(map[string]func(ctx *web.EventContext) string)
	}
	mb.fieldHints[name] = f
	return mb
}

func (mb *ModelBuilder) LayoutConfig(v *LayoutConfig) (r *ModelBuilder) {
	mb.layoutConfig = v
	return mb
}

func (mb *ModelBuilder) Placeholders(vs ...string) (r *ModelBuilder) {
	mb.placeholders = append(mb.placeholders, vs...)
	return mb
}

func (mb *ModelBuilder) Singleton(v bool) (r *ModelBuilder) {
	mb.singleton = v
	return mb
}

func (mb *ModelBuilder) FieldLabel(field *FieldBuilder, ctx *web.EventContext) (r string) {
	if f, _ := mb.fieldLabels[field.name]; f != nil {
		return f(ctx)
	}

	return field.ContextLabel(mb.Info(), ctx.Context())
}

func (mb *ModelBuilder) FieldHint(field *FieldBuilder, ctx *web.EventContext) (r string) {
	if f, _ := mb.fieldHints[field.name]; f != nil {
		return f(ctx)
	}

	return field.ContextHint(mb.Info(), ctx.Context())
}

func (b *ModelBuilder) Verifier(vf ...*perm.PermVerifierBuilder) (r *ModelBuilder) {
	b.verifiers = append(b.verifiers, vf...)
	return b
}

func (b *ModelBuilder) GetVerifiers() perm.PermVerifiers {
	return b.verifiers
}

func (mb *ModelBuilder) SubRoutesSetup(f func(mux *http.ServeMux, baseUri string)) *ModelBuilder {
	mb.subRoutesSetup = f
	return mb
}

func (mb *ModelBuilder) BindEventFunc(f web.EventHandler) web.EventHandler {
	if mb.parent == nil {
		return f
	}
	return web.EventFunc(func(ctx *web.EventContext) (r web.EventResponse, err error) {
		WithModel(ctx, mb)
		if ctx, err = mb.LoadParentsID(ctx); err != nil {
			return
		}
		return f.Handle(ctx)
	})
}

func (mb *ModelBuilder) AddPageFunc(vf *perm.PermVerifierBuilder, path string, f web.PageFunc, methods ...string) (ph *PageHandler) {
	if vf != nil && !vf.Valid() {
		vf.Path(path)
	}
	ph = NewPageHandler(path, mb.p.Wrap(mb.p.layoutFunc(mb.BindVerifiedPageFunc(vf, f), mb.layoutConfig)), methods...)
	mb.pages.Add(ph)
	return
}

func (mb *ModelBuilder) AddRawPageFunc(path string, f web.PageFunc, methods ...string) (ph *PageHandler) {
	ph = NewPageHandler(path, mb.p.Wrap(f), methods...)
	mb.pages.Add(ph)
	return
}

func (mb *ModelBuilder) BindPageFunc(f web.PageFunc) web.PageFunc {
	return func(ctx *web.EventContext) (r web.PageResponse, err error) {
		WithModel(ctx, mb)
		if ctx, err = mb.LoadParentsID(ctx); err != nil {
			return
		}
		return f(ctx)
	}
}

func (mb *ModelBuilder) BindVerifiedPageFunc(vf *perm.PermVerifierBuilder, f web.PageFunc) web.PageFunc {
	if vf != nil {
		mb.Verifier(vf)
		old := f
		f = func(ctx *web.EventContext) (_ web.PageResponse, err error) {
			if vf.Build(mb.permissioner.ReqList(ctx.R)).Denied() {
				err = perm.PermissionDenied
				return
			}
			return old(ctx)
		}
	}
	return mb.BindPageFunc(f)
}

func (mb *ModelBuilder) TPageLabel(ctx context.Context, args ...string) (s string) {
	if mb.singleton {
		return mb.TTitle(ctx, args...)
	}
	return mb.TTitlePlural(ctx, args...)
}

func (mb *ModelBuilder) TTitleAuto(ctx context.Context, args ...string) string {
	if mb.singleton && !mb.plural {
		return mb.TTitle(ctx, args...)
	}
	return mb.TTitlePlural(ctx, args...)
}

func (mb *ModelBuilder) TTitle(ctx context.Context, args ...string) string {
	return i18n.Translate(mb.Translator(), ctx, mb.label, args...)
}

func (mb *ModelBuilder) TTitlePlural(ctx context.Context, args ...string) string {
	return i18n.Translate(mb.Translator(), ctx, mb.pluralLabel, args...)
}

func (mb *ModelBuilder) TTheTitle(ctx context.Context, args ...string) string {
	return MustGetMessages(ctx).TheTitle(mb.female, mb.TTitle(ctx), args...)
}

func (mb *ModelBuilder) DefaultRecordTitle(ctx *web.EventContext, obj any) string {
	if mb.singleton {
		return mb.TTitle(ctx.Context())
	}
	return mb.TTitle(ctx.Context()) + "#" + mb.MustRecordID(obj).String()
}

func (mb *ModelBuilder) RecordTitle(obj any, ctx *web.EventContext) (s string) {
	if mb.singleton {
		return ""
	}
	switch t := obj.(type) {
	case ContextTitler:
		s = t.ContextTitle(ctx)
	case ContextStringer:
		s = t.ContextString(ctx)
	case PageTitle:
		s = t.PageTitle()
	case fmt.Stringer:
		s = t.String()
	}

	if s != "" {
		return
	}

	return mb.DefaultRecordTitle(ctx, obj)
}

func (mb *ModelBuilder) RecordTitleFetch(obj any, ctx *web.EventContext) (_ string, err error) {
	if err = mb.TitleFetcher(obj, mb.MustRecordID(obj), ctx); err != nil {
		return
	}
	return mb.RecordTitle(obj, ctx), nil
}

func (mb *ModelBuilder) LoadParentsID(ctx *web.EventContext) (_ *web.EventContext, err error) {
	var parentsID IDSlice
	if parentsID, err = mb.ParseParentsID(ctx.R); err != nil {
		return
	}
	ctx.WithContextValue(ParentsModelIDKey, parentsID)
	parentsIDPtr := web.GetContextValuer(ctx.R.Context(), ParentsModelIDKey)

	bc := GetOrInitBreadcrumbs(ctx.R)
	parents := mb.Parents()

	if root := mb.Root(); root.menuGroupName != "" {
		bc.Append(&Breadcrumb{
			Label: mb.p.menuGroups.MenuGroup(root.menuGroupName).TTitle(ctx.Context()),
		})
	}

	for i, id := range parentsID {
		bc.Append(&Breadcrumb{
			Label: parents[i].TTitlePlural(ctx.Context()),
			URI:   parents[i].modelInfo.ListingHref(parentsID[:i]...),
		})

		var (
			model = parents[i].NewModel()
			label string
		)

		id.SetTo(model)

		parentsIDPtr.Set(parentsID[:i])
		if label, err = parents[i].RecordTitleFetch(model, ctx); err != nil {
			return
		}

		bc.Append(&Breadcrumb{
			Label: label,
			URI:   parents[i].Info().DetailingHref(id.String(), parentsID[:i]...),
		})
	}

	parentsIDPtr.Set(parentsID)

	if !mb.singleton {
		uri := mb.Info().ListingHref(parentsID...)
		if ctx.R.URL.Path != uri {
			bc.Append(&Breadcrumb{
				URI:   mb.Info().ListingHref(parentsID...),
				Label: mb.TTitlePlural(ctx.Context()),
			})
		}
	}
	return ctx, nil
}

func (mb *ModelBuilder) Children() []*ModelBuilder {
	return mb.children
}

func (mb *ModelBuilder) URI() string {
	dotUri := mb.uriName
	if mb.menuGroupName != "" {
		dotUri = mb.menuGroupName + "/" + dotUri
	}
	if mb.parent != nil {
		if mb.parent.singleton {
			return fmt.Sprintf("%s/%s", mb.parent.URI(), dotUri)
		}
		return fmt.Sprintf("%s/{parent_%d_id}/%s", mb.parent.URI(), mb.parent.Depth(), dotUri)
	}
	return dotUri
}

func (mb *ModelBuilder) SplitedURI() (r []any) {
	if mb.parent != nil {
		r = mb.parent.SplitedURI()

		if !mb.parent.singleton {
			r = append(r, ParentUriPart(mb.parent.Depth()))
		}
	}
	r = append(r, mb.uriName)
	return
}

func (mb *ModelBuilder) ID() string {
	if mb.parent != nil {
		return mb.parent.ID() + "." + mb.uriName
	}
	return mb.uriName
}

func (mb *ModelBuilder) SaveURI() (uri string) {
	if mb.saveURI != "" {
		return mb.saveURI
	}

	uri = mb.URI()

	if mb.editing.mb != mb {
		uri = mb.editing.mb.URI()
	}

	return uri
}

func (mb *ModelBuilder) AddRowMenuItem(label string, componentFunc RecordMenuItemFunc) {
	mb.listing.RowMenu().SetRowMenuItem(label).ComponentFunc(componentFunc)
	mb.detailing.RowMenu().SetRowMenuItem(label).ComponentFunc(componentFunc)
}

func (mb *ModelBuilder) GetChildByID(id string) *ModelBuilder {
	for _, child := range mb.children {
		if child.id == id {
			return child
		}
	}
	return nil
}

func (mb *ModelBuilder) AddChild(child *ModelBuilder) {
	mb.AddChildH(child, nil)
}

func (mb *ModelBuilder) AddChildH(child *ModelBuilder, h ...func(mb *ModelBuilder)) *ModelBuilder {
	child.parent = mb
	mb.children = append(mb.children, child)
	if !child.notInMenu {
		label := child.label
		if !child.singleton {
			label = child.pluralLabel
		}
		mb.AddRowMenuItem(label, childRowMenuItemFunc(child))
	}

	mb.p.ModelConfigurators.ConfigureModel(child)

	for _, h := range h {
		h(child)
	}
	return mb
}

func (mb *ModelBuilder) TakeFieldAsChild(fieldName string, h ...func(mb *ModelBuilder)) *ModelBuilder {
	var (
		rfield, _ = mb.ModelType().Elem().FieldByName(fieldName)
		cfg       = ModelConfig().
				SetId(fieldName).
				SetSingleton(rfield.Type.Kind() != reflect.Slice)

		fieldObject interface{}
	)

	if cfg.singleton {
		// required field is a ptr (*Type)
		fieldObject = reflect.New(rfield.Type.Elem()).Interface()
	} else {
		// require field is slice of ptr ([]*Type)
		fieldObject = reflect.New(rfield.Type.Elem().Elem()).Interface()
		cfg.SetPluralLabel(fieldName)
	}

	cfg.SetModuleKey(mb.moduleKey)

	return mb.AddChildH(NewModelBuilder(
		mb.Builder(),
		fieldObject,
		cfg,
	), h...)
}

func (mb *ModelBuilder) ChildOf(parent *ModelBuilder) *ModelBuilder {
	mb.parent = parent
	return mb
}

func (mb *ModelBuilder) RecordID(obj interface{}) (id ID, err error) {
	s := mb.Schema()
	if s == nil {
		return
	}
	id.Schema = s
	id.Fields = s.PrimaryFields()
	id.Values = make([]any, len(id.Fields))

	for i, field := range id.Fields {
		id.Values[i] = reflectutils.MustGet(obj, field.Name())
	}

	return
}

func (mb *ModelBuilder) MustRecordID(obj interface{}) ID {
	id, err := mb.RecordID(obj)
	if err != nil {
		panic(err)
	}
	return id
}

func (mb *ModelBuilder) ParseRecordIDTo(dst any, v string) (err error) {
	var id ID
	if id, err = mb.ParseRecordID(v); err != nil {
		return
	}
	id.SetTo(dst)
	return
}

func (mb *ModelBuilder) ParseRecordID(v string) (id ID, err error) {
	if v == "" {
		return
	}

	var schema model.Schema
	if schema, err = mb.CurrentDataOperator().Schema(mb.model); err != nil {
		return
	}

	return ParseRecordID(schema, v)
}
func (mb *ModelBuilder) MustParseRecordID(v string) (id ID) {
	var err error
	id, err = mb.ParseRecordID(v)
	if err != nil {
		panic(err)
	}
	return
}

func (mb *ModelBuilder) DetailingParentsIDResolver() ParentsModelIDResolver {
	return mb.detailingParentsIDResolver
}

func (mb *ModelBuilder) SetDetailingParentsIDResolver(detailingParentsIDResolver ParentsModelIDResolver) *ModelBuilder {
	mb.detailingParentsIDResolver = detailingParentsIDResolver
	return mb
}

func (mb *ModelBuilder) GetLayoutConfig() *LayoutConfig {
	return mb.layoutConfig
}

func (mb *ModelBuilder) FieldBuilders() (r []FieldsBuilderInterface) {
	r = append(r, mb.listing, mb.editing)
	if mb.creating != nil {
		r = append(r, mb.creating)
	}
	if mb.detailing != nil {
		r = append(r, mb.detailing)
	}
	return
}

func (mb *ModelBuilder) RouteSetuper(f func(mux *http.ServeMux, uri string)) *ModelBuilder {
	mb.routeSetuper = f
	return mb
}

func (mb *ModelBuilder) WrapRouteSetuper(f func(old func(mux *http.ServeMux, uri string)) func(mux *http.ServeMux, uri string)) *ModelBuilder {
	mb.routeSetuper = f(mb.routeSetuper)
	return mb
}

func (mb *ModelBuilder) ItemRouteSetuper(f func(mux *http.ServeMux, uri string)) *ModelBuilder {
	mb.itemRouteSetuper = append(mb.itemRouteSetuper, f)
	return mb
}

func (mb *ModelBuilder) EditingDisabled() bool {
	return mb.editingDisabled
}

func (mb *ModelBuilder) SetEditingDisabled(v bool) *ModelBuilder {
	mb.editingDisabled = v
	return mb
}

func (mb *ModelBuilder) DeletingDisabled() bool {
	return mb.editingDisabled
}

func (mb *ModelBuilder) SetDeletingDisabled(v bool) *ModelBuilder {
	mb.deletingDisabled = v
	return mb
}

func (mb *ModelBuilder) CreatingDisabled() bool {
	return mb.creatingDisabled
}

func (mb *ModelBuilder) SetCreatingDisabled(v bool) *ModelBuilder {
	mb.creatingDisabled = v
	return mb
}

func (mb *ModelBuilder) DetailingDisabled() bool {
	return mb.detailingDisabled
}

func (mb *ModelBuilder) SetDetailingDisabled(v bool) *ModelBuilder {
	mb.detailingDisabled = v
	return mb
}
