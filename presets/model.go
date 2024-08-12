package presets

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"github.com/qor5/admin/v3/reflect_utils"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	"github.com/sunfmin/reflectutils"
)

type ModelBuilderOptionFunc func(mb *ModelBuilder)

func (f ModelBuilderOptionFunc) Apply(mb *ModelBuilder) {
	f(mb)
}

type ModelBuilderOption interface {
	Apply(mb *ModelBuilder)
}

func WithDataOperator(do DataOperator) ModelBuilderOptionFunc {
	return func(mb *ModelBuilder) {
		mb.dataOperator = do
	}
}

type ModelBuilder struct {
	p *Builder
	writeFieldBuilders,
	listFieldBuilders,
	detailFieldBuilders FieldBuilders

	parent              *ModelBuilder
	model               interface{}
	primaryField        string
	modelType           reflect.Type
	menuGroupName       string
	notInMenu           bool
	menuIcon            string
	defaultURLQueryFunc func(*http.Request) url.Values
	fieldLabels         []string
	placeholders        []string
	listing             *ListingBuilder
	detailing           *DetailingBuilder
	editing             *EditingBuilder
	creating            *EditingBuilder
	hasDetailing        bool
	female              bool
	rightDrawerWidth    string
	link                string
	layoutConfig        *LayoutConfig
	modelInfo           *ModelInfo
	plugins             []ModelPlugin
	subRoutesSetup      func(mux *http.ServeMux, uri string)
	listingMenu         Menu
	detailingMenu       Menu

	EventsHub web.EventsHub

	children []*ModelBuilder
	saveURI  string

	creationDisabled OkHandled
	editionDisabled  OkHandled
	verifierModel    *ModelBuilder

	detailingParentsIDResolver ParentsModelIDResolver
	ModelBuilderConfigAttributes
	BeforeFormUnmarshallHandlers ModelFormUnmarshallHandlers
	PostFormUnmarshallHandlers   ModelFormUnmarshallHandlers
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

func NewModelBuilder(p *Builder, model interface{}, options ...ModelBuilderOption) (mb *ModelBuilder) {
	mb = &ModelBuilder{p: p, model: model, primaryField: "ID"}
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

	if mb.uriName == "" {
		if mb.singleton {
			mb.uriName = strcase.ToKebab(mb.label)
		} else {
			mb.uriName = strcase.ToKebab(mb.pluralLabel)
		}
	}

	mb.modelInfo = &ModelInfo{mb: mb}
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

	// Be aware the uriName here is still the original struct
	mb.newListing()
	mb.newDetailing()
	mb.newEditing()
	mb.editing.Creating()

	return
}

func (mb *ModelBuilder) Female() bool {
	return mb.female
}

func (mb *ModelBuilder) SetFemale(female bool) *ModelBuilder {
	mb.female = female
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

func (mb *ModelBuilder) WithDataOperator(h func(dataOperator DataOperator)) bool {
	if do := mb.CurrentDataOperator(); do != nil {
		h(do)
		return true
	}
	return false
}

func (mb *ModelBuilder) UpdateDataOperator(h func(dataOperator DataOperator) DataOperator) *ModelBuilder {
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

func (mb *ModelBuilder) Fetcher(obj interface{}, id string, ctx *web.EventContext) (err error) {
	return mb.CurrentDataOperator().Fetch(obj, id, ctx)
}

func (mb *ModelBuilder) TitleFetcher(obj interface{}, id string, ctx *web.EventContext) (err error) {
	return mb.CurrentDataOperator().FetchTitle(obj, id, ctx)
}

func (mb *ModelBuilder) Saver(obj interface{}, id string, ctx *web.EventContext) (err error) {
	return mb.CurrentDataOperator().Save(obj, id, ctx)
}

func (mb *ModelBuilder) Deleter(obj interface{}, id string, ctx *web.EventContext) (err error) {
	return mb.CurrentDataOperator().Delete(obj, id, ctx)
}

func (mb *ModelBuilder) RegisterEventFunc(eventFuncId string, ef web.EventFunc) (key string) {
	return mb.EventsHub.RegisterEventFunc(eventFuncId, ef)
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

func (mb *ModelBuilder) URIName(v string) (r *ModelBuilder) {
	mb.uriName = v
	return mb
}

func (mb *ModelBuilder) DefaultURLQueryFunc(v func(*http.Request) url.Values) (r *ModelBuilder) {
	mb.defaultURLQueryFunc = v
	return mb
}

func (mb *ModelBuilder) PrimaryField(v string) (r *ModelBuilder) {
	mb.primaryField = v
	return mb
}

func (mb *ModelBuilder) InMenu(v bool) (r *ModelBuilder) {
	mb.notInMenu = !v
	return mb
}

func (mb *ModelBuilder) MenuIcon(v string) (r *ModelBuilder) {
	mb.menuIcon = v
	return mb
}

func (mb *ModelBuilder) Label(v string) (r *ModelBuilder) {
	mb.label = v
	return mb
}

func (mb *ModelBuilder) GetLabel() string {
	return mb.label
}

func (mb *ModelBuilder) Labels(vs ...string) (r *ModelBuilder) {
	mb.fieldLabels = append(mb.fieldLabels, vs...)
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

func (mb *ModelBuilder) GetComponentFuncField(obj interface{}, field *FieldBuilder) (r *FieldContext) {
	r = &FieldContext{
		Obj:       obj,
		ModelInfo: mb.Info(),
		Name:      field.name,
		Label:     mb.getLabel(field.NameLabel),
	}
	return
}

func (mb *ModelBuilder) getLabel(field NameLabel) (r string) {
	if len(field.label) > 0 {
		return field.label
	}

	for i := 0; i < len(mb.fieldLabels)-1; i = i + 2 {
		if mb.fieldLabels[i] == field.name {
			return mb.fieldLabels[i+1]
		}
	}

	return humanizeString(field.name)
}

func (mb *ModelBuilder) SubRoutesSetup(f func(mux *http.ServeMux, baseUri string)) *ModelBuilder {
	mb.subRoutesSetup = f
	return mb
}

func (mb *ModelBuilder) BindEventFunc(f web.EventFunc) web.EventFunc {
	if mb.parent == nil {
		return f
	}
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		if ctx, err = mb.LoadParentsID(ctx); err != nil {
			return
		}
		if r, err = f(ctx); err != nil {
			ShowMessage(&r, err.Error(), "error")
			err = nil
		}
		return
	}
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

func (mb *ModelBuilder) TTitle(r *http.Request) string {
	return i18n.T(r, ModelsI18nModuleKey, mb.label)
}

func (mb *ModelBuilder) TTheTitle(r *http.Request) string {
	return MustGetMessages(r).TheTitle(mb.female, mb.TTitle(r))
}

func (mb *ModelBuilder) DefaultRecordTitle(ctx *web.EventContext, obj any) string {
	idStr := mb.ID(obj).String()
	return mb.TTitle(ctx.R) + "#" + idStr
}

func (mb *ModelBuilder) RecordTitle(obj any, ctx *web.EventContext) string {
	if pt, ok := obj.(PageTitle); ok {
		return pt.PageTitle()
	} else if s, ok := obj.(fmt.Stringer); ok {
		return s.String()
	}

	return mb.DefaultRecordTitle(ctx, obj)
}

func (mb *ModelBuilder) RecordTitleFetch(obj any, ctx *web.EventContext) (_ string, err error) {
	if err = mb.TitleFetcher(obj, mb.ID(obj).String(), ctx); err != nil {
		return
	}
	return mb.RecordTitle(obj, ctx), nil
}

func (mb *ModelBuilder) LoadParentsID(ctx *web.EventContext) (_ *web.EventContext, err error) {
	var parentsID IDSlice
	if parentsID, err = mb.ParseParentsID(ctx.R); err != nil {
		return
	}
	bc := GetOrInitBreadcrumbs(ctx.R)
	parents := mb.Parents()
	for i, id := range parentsID {
		bc.Append(&Breadcrumb{
			Label: i18n.T(ctx.R, ModelsI18nModuleKey, id.Model.pluralLabel),
			URI:   parents[i].modelInfo.ListingHref(parentsID[:i]...),
		})

		var (
			model = id.Model.NewModel()
			label string
		)

		id.SetTo(model)

		if label, err = id.Model.RecordTitleFetch(model, ctx); err != nil {
			return
		}

		bc.Append(&Breadcrumb{
			Label: label,
			URI:   id.Model.Info().DetailingHref(id.Value, parentsID[:i]...),
		})
	}

	if !mb.singleton && ctx.Param("id") != "" {
		bc.Append(&Breadcrumb{
			URI:   mb.Info().ListingHref(parentsID...),
			Label: i18n.T(ctx.R, ModelsI18nModuleKey, mb.pluralLabel),
		})
	}
	return ctx.WithContextValue(ParentsModelIDKey, parentsID), nil
}

func (mb *ModelBuilder) Children() []*ModelBuilder {
	return mb.children
}

func (mb *ModelBuilder) URI() string {
	if mb.parent != nil {
		if mb.parent.singleton {
			return fmt.Sprintf("%s/%s", mb.parent.URI(), mb.uriName)
		}
		return fmt.Sprintf("%s/{parent_%d_id}/%s", mb.parent.URI(), mb.parent.Depth(), mb.uriName)
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
	mb.listing.RowMenu().RowMenuItem(label).ComponentFunc(componentFunc)
	mb.detailing.RowMenu().RowMenuItem(label).ComponentFunc(componentFunc)
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
			label = inflection.Plural(child.label)
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

func (mb *ModelBuilder) ID(obj interface{}) ID {
	if obj == nil {
		return ID{}
	}
	return ID{Model: mb, Value: reflectutils.MustGet(obj, mb.primaryField)}
}

func (mb *ModelBuilder) ParseID(v string) (id ID, err error) {
	id.Model = mb
	field, _ := mb.modelType.Elem().FieldByName("ID")
	switch field.Type.Kind() {
	case reflect.String:
		id.Value = v
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var bsize int
		switch field.Type.Kind() {
		case reflect.Uint8:
			bsize = 8
		case reflect.Uint16:
			bsize = 16
		case reflect.Uint32:
			bsize = 32
		case reflect.Uint64:
			bsize = 64
		}
		if id.Value, err = strconv.ParseUint(v, 10, bsize); err != nil {
			return
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var bsize int
		switch field.Type.Kind() {
		case reflect.Int8:
			bsize = 8
		case reflect.Int16:
			bsize = 16
		case reflect.Int32:
			bsize = 32
		case reflect.Int64:
			bsize = 64
		}

		if id.Value, err = strconv.ParseInt(v, 10, bsize); err != nil {
			return
		}
	default:
		err = errors.New(fmt.Sprintf("Unsupported type: %v", field.Type))
		return
	}

	id.Value = reflect.ValueOf(id.Value).Convert(field.Type).Interface()
	return
}

func (mb *ModelBuilder) CreationDisabled() OkHandled {
	return mb.creationDisabled
}

func (mb *ModelBuilder) SetCreationDisabled(creationDisabled OkHandled) {
	mb.creationDisabled = creationDisabled
}

func (mb *ModelBuilder) EditionDisabled() OkHandled {
	return mb.editionDisabled
}

func (mb *ModelBuilder) SetEditionDisabled(editDisabled OkHandled) {
	mb.editionDisabled = editDisabled
}

func (mb *ModelBuilder) DetailingParentsIDResolver() ParentsModelIDResolver {
	return mb.detailingParentsIDResolver
}

func (mb *ModelBuilder) SetDetailingParentsIDResolver(detailingParentsIDResolver ParentsModelIDResolver) {
	mb.detailingParentsIDResolver = detailingParentsIDResolver
}
