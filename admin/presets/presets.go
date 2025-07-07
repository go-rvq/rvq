package presets

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/datafield"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type Builder struct {
	prefix                                string
	models                                []*ModelBuilder
	handler                               http.Handler
	builder                               *web.Builder
	i18nBuilder                           *i18n.Builder
	logger                                *zap.Logger
	permissionBuilder                     *perm.Builder
	verifier                              *perm.Verifier
	layoutFunc                            func(in web.PageFunc, cfg *LayoutConfig) (out web.PageFunc)
	detailLayoutFunc                      func(in web.PageFunc, cfg *LayoutConfig) (out web.PageFunc)
	dataOperator                          DataOperator
	messagesFunc                          MessagesFunc
	homePageFunc                          web.PageFunc
	notFoundFunc                          web.PageFunc
	homePageLayoutConfig                  *LayoutConfig
	notFoundPageLayoutConfig              *LayoutConfig
	brandFunc                             ComponentFunc
	profileFunc                           ComponentFunc
	switchLanguageFunc                    ComponentFunc
	brandProfileSwitchLanguageDisplayFunc func(brand, profile, switchLanguage h.HTMLComponent) h.HTMLComponent
	menuTopItems                          map[string]ComponentFunc
	notificationCountFunc                 func(ctx *web.EventContext) int
	notificationContentFunc               ComponentFunc
	brandTitle                            string
	vuetifyOptions                        string
	progressBarColor                      string
	rightDrawerWidth                      string
	writeFieldDefaults                    *FieldDefaults
	listFieldDefaults                     *FieldDefaults
	detailFieldDefaults                   *FieldDefaults
	extraAssets                           []*extraAsset
	assetFunc                             AssetFunc
	menuGroups                            MenuGroups
	menuOrder                             []interface{}
	wrapHandlers                          map[string]func(in http.Handler) (out http.Handler)
	plugins                               []Plugin
	ModelConfigurators                    ModelConfigurators
	ModelSetupFactories                   ModelSetupFactories
	skipNotFoundHandler                   func(r *http.Request) bool
	muxSetup                              []func(prefix string, r *http.ServeMux)
	permissions                           *PermMenu
	pageHandlers                          PageHandlers
	pages                                 []*PageBuilder
	verifiers                             perm.PermVerifiers

	datafield.DataField[*Builder]
}

type AssetFunc func(ctx *web.EventContext)

type extraAsset struct {
	path        string
	contentType string
	body        web.ComponentsPack
	refTag      string
}

const (
	CoreI18nModuleKey   i18n.ModuleKey = "CoreI18nModuleKey"
	ModelsI18nModuleKey i18n.ModuleKey = "ModelsI18nModuleKey"
)

func New(i18nB *i18n.Builder) *Builder {
	l, _ := zap.NewDevelopment()
	r := datafield.New(&Builder{
		logger:  l,
		builder: web.New(),
		i18nBuilder: i18nB.
			RegisterForModule(language.English, CoreI18nModuleKey, Messages_en_US).
			RegisterForModule(language.BrazilianPortuguese, CoreI18nModuleKey, Messages_pt_BR),
		writeFieldDefaults:   NewFieldDefaults(WRITE),
		listFieldDefaults:    NewFieldDefaults(LIST),
		detailFieldDefaults:  NewFieldDefaults(DETAIL),
		progressBarColor:     "amber",
		menuTopItems:         make(map[string]ComponentFunc),
		brandTitle:           "Admin",
		rightDrawerWidth:     "600",
		verifier:             perm.NewVerifier(PermModule, nil),
		homePageLayoutConfig: &LayoutConfig{SearchBoxInvisible: true},
		notFoundPageLayoutConfig: &LayoutConfig{
			SearchBoxInvisible:          true,
			NotificationCenterInvisible: true,
		},
		wrapHandlers:        make(map[string]func(in http.Handler) (out http.Handler)),
		ModelSetupFactories: DefaultModelSetupFactories,
		skipNotFoundHandler: func(r *http.Request) bool {
			return false
		},
	})
	r.GetWebBuilder().RegisterEventHandler(EventOpenConfirmDialog, web.EventFunc(r.openConfirmDialog))
	r.layoutFunc = r.DefaultLayout
	r.detailLayoutFunc = r.DefaultLayout
	return r
}

func (b *Builder) I18n() (r *i18n.Builder) {
	return b.i18nBuilder
}

func (b *Builder) Permission(v *perm.Builder) (r *Builder) {
	b.permissionBuilder = v
	b.verifier = perm.NewVerifier(PermModule, v)
	return b
}

func (b *Builder) GetPermission() (r *perm.Builder) {
	return b.permissionBuilder
}

func (b *Builder) URIPrefix(v string) (r *Builder) {
	b.prefix = strings.TrimRight(v, "/")
	return b
}

func (b *Builder) GetURIPrefix() string {
	return b.prefix
}

func (b *Builder) LayoutFunc(v func(in web.PageFunc, cfg *LayoutConfig) (out web.PageFunc)) (r *Builder) {
	b.layoutFunc = v
	return b
}

func (b *Builder) GetLayoutFunc() func(in web.PageFunc, cfg *LayoutConfig) (out web.PageFunc) {
	return b.layoutFunc
}

func (b *Builder) DetailLayoutFunc(v func(in web.PageFunc, cfg *LayoutConfig) (out web.PageFunc)) (r *Builder) {
	b.detailLayoutFunc = v
	return b
}

func (b *Builder) GetDetailLayoutFunc() func(in web.PageFunc, cfg *LayoutConfig) (out web.PageFunc) {
	return b.detailLayoutFunc
}

func (b *Builder) HomePageLayoutConfig(v *LayoutConfig) (r *Builder) {
	b.homePageLayoutConfig = v
	return b
}

func (b *Builder) NotFoundPageLayoutConfig(v *LayoutConfig) (r *Builder) {
	b.notFoundPageLayoutConfig = v
	return b
}

func (b *Builder) Builder(v *web.Builder) (r *Builder) {
	b.builder = v
	return b
}

func (b *Builder) GetWebBuilder() (r *web.Builder) {
	return b.builder
}

func (b *Builder) Logger(v *zap.Logger) (r *Builder) {
	b.logger = v
	return b
}

func (b *Builder) MessagesFunc(v MessagesFunc) (r *Builder) {
	b.messagesFunc = v
	return b
}

func (b *Builder) HomePageFunc(v web.PageFunc) (r *Builder) {
	b.homePageFunc = v
	return b
}

func (b *Builder) NotFoundFunc(v web.PageFunc) (r *Builder) {
	b.notFoundFunc = v
	return b
}

func (b *Builder) BrandFunc(v ComponentFunc) (r *Builder) {
	b.brandFunc = v
	return b
}

func (b *Builder) ProfileFunc(v ComponentFunc) (r *Builder) {
	b.profileFunc = v
	return b
}

func (b *Builder) GetProfileFunc() ComponentFunc {
	return b.profileFunc
}

func (b *Builder) SwitchLanguageFunc(v ComponentFunc) (r *Builder) {
	b.switchLanguageFunc = v
	return b
}

func (b *Builder) BrandProfileSwitchLanguageDisplayFuncFunc(f func(brand, profile, switchLanguage h.HTMLComponent) h.HTMLComponent) (r *Builder) {
	b.brandProfileSwitchLanguageDisplayFunc = f
	return b
}

func (b *Builder) NotificationFunc(contentFunc ComponentFunc, countFunc func(ctx *web.EventContext) int) (r *Builder) {
	b.notificationCountFunc = countFunc
	b.notificationContentFunc = contentFunc
	b.GetWebBuilder().RegisterEventHandler(actions.NotificationCenter, web.EventFunc(b.notificationCenter))
	return b
}

func (b *Builder) BrandTitle(v string) (r *Builder) {
	b.brandTitle = v
	return b
}

func (b *Builder) GetBrandTitle() string {
	return b.brandTitle
}

func (b *Builder) VuetifyOptions(v string) (r *Builder) {
	b.vuetifyOptions = v
	return b
}

func (b *Builder) RightDrawerWidth(v string) (r *Builder) {
	b.rightDrawerWidth = v
	return b
}

func (b *Builder) ProgressBarColor(v string) (r *Builder) {
	b.progressBarColor = v
	return b
}

func (b *Builder) GetProgressBarColor() string {
	return b.progressBarColor
}

func (b *Builder) AssetFunc(v AssetFunc) (r *Builder) {
	b.assetFunc = v
	return b
}

func (b *Builder) ExtraAsset(path string, contentType string, body web.ComponentsPack, refTag ...string) (r *Builder) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	var theOne *extraAsset
	for _, ea := range b.extraAssets {
		if ea.path == path {
			theOne = ea
			break
		}
	}

	if theOne == nil {
		theOne = &extraAsset{path: path, contentType: contentType, body: body}
		b.extraAssets = append(b.extraAssets, theOne)
	} else {
		theOne.contentType = contentType
		theOne.body = body
	}

	if len(refTag) > 0 {
		theOne.refTag = refTag[0]
	}

	return b
}

func (b *Builder) FieldDefaults(v FieldMode) (r *FieldDefaults) {
	if v == WRITE {
		return b.writeFieldDefaults
	}

	if v == LIST {
		return b.listFieldDefaults
	}

	if v == DETAIL {
		return b.detailFieldDefaults
	}

	return r
}

func (b *Builder) NewFieldsBuilder(v FieldMode) (r *FieldsBuilder) {
	r = NewFieldsBuilder(b).Defaults(b.FieldDefaults(v))
	return
}

func (b *Builder) Model(v interface{}, opts ...ModelBuilderOption) (r *ModelBuilder) {
	r = NewModelBuilder(b, v, opts...)
	b.ModelConfigurators.ConfigureModel(r)
	b.models = append(b.models, r)
	return r
}

func (b *Builder) GetModelByID(id string) *ModelBuilder {
	for _, mb := range b.models {
		if mb.id == id {
			return mb
		}
	}
	return nil
}

func (b *Builder) GetModel(typ any) *ModelBuilder {
	var t reflect.Type
	switch tp := typ.(type) {
	case reflect.Type:
		if tp.Kind() == reflect.Struct {
			tp = reflect.PointerTo(tp)
		}
		t = tp
	default:
		t = reflect.TypeOf(typ)
	}

	for _, model := range b.models {
		if model.modelType == t {
			return model
		}
	}
	return nil
}

func (b *Builder) DataOperator(v DataOperator) (r *Builder) {
	b.dataOperator = v
	return b
}

func (b *Builder) GetDataOperator() DataOperator {
	return b.dataOperator
}

func modelNames(ms []*ModelBuilder) (r []string) {
	for _, m := range ms {
		r = append(r, m.uriName)
	}
	return
}

func (b *Builder) MenuGroup(name string) *MenuGroupBuilder {
	mgb := b.menuGroups.MenuGroup(name)
	if !b.isMenuGroupInOrder(mgb) {
		b.menuOrder = append(b.menuOrder, mgb)
	}
	return mgb
}

func (b *Builder) isMenuGroupInOrder(mgb *MenuGroupBuilder) bool {
	for _, v := range b.menuOrder {
		if v == mgb {
			return true
		}
	}
	return false
}

func (b *Builder) removeMenuGroupInOrder(mgb *MenuGroupBuilder) {
	for i, om := range b.menuOrder {
		if om == mgb {
			b.menuOrder = append(b.menuOrder[:i], b.menuOrder[i+1:]...)
			break
		}
	}
}

// item can be Slug name, model name, *MenuGroupBuilder
// the underlying logic is using Slug name,
// so if the Slug name is customized, item must be the Slug name
// example:
// b.MenuOrder(
//
//	b.MenuGroup("Product Management").SubItems(
//		"products",
//		"Variant",
//	),
//	"customized-uri",
//
// )
func (b *Builder) MenuOrder(items ...interface{}) {
	for _, item := range items {
		switch v := item.(type) {
		case string:
			b.menuOrder = append(b.menuOrder, v)
		case *MenuGroupBuilder:
			if b.isMenuGroupInOrder(v) {
				b.removeMenuGroupInOrder(v)
			}
			b.menuOrder = append(b.menuOrder, v)
			for _, item := range v.subMenuItems {
				if item[0] == '/' {
					if p := b.GetPage(item); p != nil {
						p.menuGroup = v.name
					}
				} else if mb := b.GetModelByID(item); mb != nil {
					mb.menuGroupName = v.name
				}
			}
		default:
			panic(fmt.Sprintf("unknown menu order item type: %T\n", item))
		}
	}
}

type defaultMenuIconRE struct {
	re   *regexp.Regexp
	icon string
}

var defaultMenuIconREs = []defaultMenuIconRE{
	// user
	{re: regexp.MustCompile(`\busers?|members?\b`), icon: "mdi-account"},
	// store
	{re: regexp.MustCompile(`\bstores?\b`), icon: "mdi-store"},
	// order
	{re: regexp.MustCompile(`\borders?\b`), icon: "mdi-cart"},
	// product
	{re: regexp.MustCompile(`\bproducts?\b`), icon: "mdi-format-list-bulleted"},
	// post
	{re: regexp.MustCompile(`\bposts?|articles?\b`), icon: "mdi-note"},
	// web
	{re: regexp.MustCompile(`\bweb|site\b`), icon: "mdi-web"},
	// seo
	{re: regexp.MustCompile(`\bseo\b`), icon: "mdi-search-web"},
	// i18n
	{re: regexp.MustCompile(`\bi18n|translations?\b`), icon: "mdi-translate"},
	// chart
	{re: regexp.MustCompile(`\banalytics?|charts?|statistics?\b`), icon: "mdi-google-analytics"},
	// dashboard
	{re: regexp.MustCompile(`\bdashboard\b`), icon: "mdi-view-dashboard"},
	// setting
	{re: regexp.MustCompile(`\bsettings?\b`), icon: "mdi-cog"},
}

func defaultMenuIcon(mLabel string) string {
	ws := strings.Join(strings.Split(strcase.ToSnake(mLabel), "_"), " ")
	for _, v := range defaultMenuIconREs {
		if v.re.MatchString(ws) {
			return v.icon
		}
	}

	return "mdi-alert-octagon-outline"
}

const (
	menuFontWeight    = "500"
	subMenuFontWeight = "400"
)

func (b *Builder) CreateMenus(ctx *web.EventContext) (r h.HTMLComponent) {
	var (
		mMap = make(map[string]*ModelBuilder)
		pMap = make(map[string]*PageBuilder)
	)

	for _, m := range b.models {
		if !m.notInMenu {
			mMap[m.id] = m
		}
	}

	for _, page := range b.pages {
		if !page.notInMenu {
			pMap[page.path] = page
		}
	}

	var (
		activeMenuItem string
		selection      string
		menus          []h.HTMLComponent
		inOrderMap     = make(map[string]struct{})
	)

	for _, om := range b.menuOrder {
		switch v := om.(type) {
		case *MenuGroupBuilder:
			disabled := false
			groupIcon := v.icon
			if groupIcon == "" {
				groupIcon = defaultMenuIcon(v.name)
			}

			var title string
			if v.title != nil {
				title = v.TTitle(ctx.Context())
			} else {
				title = i18n.T(ctx.Context(), ModelsI18nModuleKey, v.name)
			}

			subMenus := []h.HTMLComponent{
				h.Template(
					VListItem(
						web.Slot(
							VIcon(groupIcon),
						).Name("prepend"),
						VListItemTitle().Attr("style", fmt.Sprintf("white-space: normal; font-weight: %s;font-size: 14px;", menuFontWeight)),
						// VListItemTitle(h.Text(i18n.T(ctx.R, ModelsI18nModuleKey, v.name))).
					).Attr("v-bind", "props").
						Title(title).
						Class("rounded-lg"),
					// Value(i18n.T(ctx.R, ModelsI18nModuleKey, v.name)),
				).Attr("v-slot:activator", "{ props }"),
			}

			subCount := 0
			for _, subOm := range v.subMenuItems {
				if subOm[0] == '/' {
					p := pMap[subOm]
					if p == nil || p.notInMenu || (p.verififer != nil && p.Verifier(ctx.R).Denied()) {
						continue
					}
					subMenus = append(subMenus, p.menuItem(ctx, true))
					subCount++
					inOrderMap[p.path] = struct{}{}
					if p.isMenuItemActive(ctx) {
						// activeMenuItem = m.label
						activeMenuItem = v.name
						selection = p.path
					}
				} else {
					m, _ := mMap[subOm]
					if m == nil {
						continue
					}
					if m.notInMenu {
						continue
					}
					if m.permissioner.ReqLister(ctx.R).Denied() {
						continue
					}
					subMenus = append(subMenus, m.menuItem(ctx, true))
					subCount++
					inOrderMap[m.id] = struct{}{}
					if m.isMenuItemActive(ctx) {
						// activeMenuItem = m.label
						activeMenuItem = v.name
						selection = m.label
					}
				}
			}
			if subCount == 0 {
				continue
			}
			if disabled {
				continue
			}

			menus = append(menus,
				VListGroup(subMenus...).Value(v.name),
			)
		case string:
			if v[0] == '/' {
				p := pMap[v]
				if p == nil || p.notInMenu || (p.verififer != nil && p.Verifier(ctx.R).Denied()) {
					continue
				}

				menuItem := p.menuItem(ctx, false)
				menus = append(menus, menuItem)
				inOrderMap[p.path] = struct{}{}

				if p.isMenuItemActive(ctx) {
					selection = p.path
				}
			} else {
				m, ok := mMap[v]
				if !ok {
					m = mMap[inflection.Plural(strcase.ToKebab(v))]
				}
				if m == nil {
					continue
				}
				if m.permissioner.ReqLister(ctx.R).Denied() {
					continue
				}

				if m.notInMenu {
					continue
				}
				menuItem := m.menuItem(ctx, false)
				menus = append(menus, menuItem)
				inOrderMap[m.id] = struct{}{}

				if m.isMenuItemActive(ctx) {
					selection = m.label
				}
			}
		}
	}

	for _, m := range b.models {
		_, ok := inOrderMap[m.id]
		if ok {
			continue
		}

		if m.permissioner.ReqLister(ctx.R).Denied() {
			continue
		}

		if m.notInMenu {
			continue
		}

		if m.isMenuItemActive(ctx) {
			selection = m.label
		}
		menus = append(menus, m.menuItem(ctx, false))
	}

	for _, p := range b.pages {
		_, ok := inOrderMap[p.path]
		if ok {
			continue
		}

		if p == nil || p.notInMenu || (p.verififer != nil && p.Verifier(ctx.R).Denied()) {
			continue
		}

		if p.isMenuItemActive(ctx) {
			selection = p.path
		}

		menus = append(menus, p.menuItem(ctx, false))
	}

	r = web.Scope(
		VList(menus...).Class("main-menu").
			OpenStrategy("single").
			Class("primary--text").
			Density(DensityCompact).
			Attr("v-model:opened", "locals.menuOpened").
			Attr("v-model:selected", "locals.selection"),
		// .Attr("v-model:selected", h.JSONString([]string{"Pages"})),
	).Slot("{ locals }").LocalsInit(
		fmt.Sprintf(`{ menuOpened:  ["%s"]}`, activeMenuItem),
		fmt.Sprintf(`{ selection:  ["%s"]}`, selection),
	)
	return
}

func (b *Builder) RunBrandFunc(ctx *web.EventContext) (r h.HTMLComponent) {
	if b.brandFunc != nil {
		return b.brandFunc(ctx)
	}
	return h.H1(i18n.T(ctx.Context(), ModelsI18nModuleKey, b.brandTitle)).Class("text-h6")
}

func (b *Builder) RunSwitchLanguageFunc(ctx *web.EventContext) (r h.HTMLComponent) {
	if b.switchLanguageFunc != nil {
		return b.switchLanguageFunc(ctx)
	}

	supportLanguages := b.I18n().GetSupportLanguagesFromRequest(ctx.R)

	if len(b.I18n().GetSupportLanguages()) <= 1 || len(supportLanguages) == 0 {
		return nil
	}
	queryName := b.I18n().GetQueryName()
	msgr := MustGetMessages(ctx.Context())
	if len(supportLanguages) == 1 {
		return h.Template().Children(
			h.Div(
				VList(
					VListItem(
						web.Slot(
							VIcon("mdi-widget-translate").Size(SizeSmall).Class("mr-4 ml-1"),
						).Name("prepend"),
						VListItemTitle(
							h.Div(h.Text(fmt.Sprintf("%s%s %s", msgr.Language, msgr.Colon, display.Self.Name(supportLanguages[0])))).Role("button"),
						),
					).Class("pa-0").Density(DensityCompact),
				).Class("pa-0 ma-n4 mt-n6"),
			).Attr("@click", web.Plaid().Query(queryName, supportLanguages[0].String()).Go()),
		)
	}

	matcher := language.NewMatcher(supportLanguages)

	lang := ctx.R.FormValue(queryName)
	if lang == "" {
		lang = b.i18nBuilder.GetCurrentLangFromCookie(ctx.R)
	}

	accept := ctx.R.Header.Get("Accept-Language")

	var displayLanguage language.Tag
	_, i := language.MatchStrings(matcher, lang, accept)
	displayLanguage = supportLanguages[i]

	var languages []h.HTMLComponent
	for _, tag := range supportLanguages {
		languages = append(languages,
			h.Div(
				VListItem(
					VListItemTitle(
						h.Div(h.Text(display.Self.Name(tag))),
					),
				).Attr("@click", web.Plaid().Query(queryName, tag.String()).Go()),
			),
		)
	}

	oldIcon := VMenu().Children(
		h.Template().Attr("v-slot:activator", "{isActive, props}").Children(
			h.Div(
				VList(
					VListItem(
						VListItemTitle(
							h.Text(fmt.Sprintf("%s%s %s", msgr.Language, msgr.Colon, display.Self.Name(displayLanguage))),
						).Class("text-subtitle-2 font-weight-regular"),
						web.Slot(
							VIcon("mdi-translate-variant").Size(SizeSmall).Class(""),
						).Name("append"),
					).Class("pa-0").Density(DensityCompact),
				).Class("pa-0 ma-n4 mt-n6"),
			).Attr("v-bind", "props"),
		),
		VList(
			languages...,
		).Density(DensityCompact),
	)
	_ = oldIcon
	return VMenu().Children(
		h.Template().Attr("v-slot:activator", "{isActive, props}").Children(
			VRow(
				VCol(
					VIcon("mdi-translate")).Cols(1),
				VCol(VIcon("mdi-menu-down")).Cols(1),
			).Attr("v-bind", "props"),
		),
		VList(
			languages...,
		).Density(DensityCompact),
	)
}

func (b *Builder) AddMenuTopItemFunc(key string, v ComponentFunc) (r *Builder) {
	b.menuTopItems[key] = v
	return b
}

func (b *Builder) RunBrandProfileSwitchLanguageDisplayFunc(brand, profile, switchLanguage h.HTMLComponent, ctx *web.EventContext) (r h.HTMLComponent) {
	if b.brandProfileSwitchLanguageDisplayFunc != nil {
		return b.brandProfileSwitchLanguageDisplayFunc(brand, profile, switchLanguage)
	}

	var items []h.HTMLComponent
	items = append(items,
		h.If(brand != nil,
			VListItem(
				VCardText(brand),
			),
		),
		h.If(profile != nil,
			VListItem(
				VCardText(profile),
			),
		),
		h.If(switchLanguage != nil,
			VListItem(
				VCardText(switchLanguage),
			).Density(DensityCompact),
		),
	)
	for _, v := range b.menuTopItems {
		items = append(items,
			h.If(v(ctx) != nil,
				VListItem(
					VCardText(v(ctx)),
				),
			))
	}

	return h.Div(
		items...,
	)
}

const (
	NotificationCenterPortalName   = "notification-center"
	DefaultConfirmDialogPortalName = "presets_ConfirmDialogPortalName"
	ListingDialogPortalName        = "presets_ListingDialogPortalName"
	FormPortalName                 = "presets_FormPortalName"
	FlashPortalName                = "flash"
)

const (
	CloseRightDrawerVarScript   = "vars.presetsRightDrawer = false"
	closeDialogVarScript        = "vars.presetsDialog = false"
	CloseListingDialogVarScript = "vars.presetsListingDialog = false"
)

func (b *Builder) Overlay(ctx *web.EventContext, r *web.EventResponse, comp h.HTMLComponent, width string) {
	overlayType := actions.OverlayMode(ctx.Param(ParamOverlay))

	if overlayType == actions.Dialog {
		b.dialog(ctx, r, comp, width)
		return
	} else if overlayType == actions.Content {
		b.contentDrawer(ctx, r, comp, width)
		return
	}
	b.rightDrawer(r, comp, width)
}

func (b *Builder) rightDrawer(r *web.EventResponse, comp h.HTMLComponent, width string) {
	if width == "" {
		width = b.rightDrawerWidth
	}
	r.UpdatePortal(
		actions.RightDrawer.PortalName(),
		VNavigationDrawer(
			web.GlobalEvents().Attr("@keyup.esc", "vars.presetsRightDrawer = false"),
			web.Portal(comp).Name(actions.RightDrawer.ContentPortalName()),
		).
			// Attr("@input", "plaidForm.dirty && vars.presetsRightDrawer == false && !confirm('You have unsaved changes on this form. If you close it, you will lose all unsaved changes. Are you sure you want to close it?') ? vars.presetsRightDrawer = true: vars.presetsRightDrawer = $event"). // remove because drawer plaidForm has to be reset when UpdateOverlayContent
			Class("v-navigation-drawer--temporary").
			Attr("v-model", "vars.presetsRightDrawer").
			Location(LocationRight).
			Temporary(true).
			// Fixed(true).
			Width(width).
			Attr(":height", `"100%"`),
		// Temporary(true),
		// HideOverlay(true).
		// Floating(true).

	)
	r.RunScript = "setTimeout(function(){ vars.presetsRightDrawer = true }, 100)"
}

func (b *Builder) contentDrawer(ctx *web.EventContext, r *web.EventResponse, comp h.HTMLComponent, width string) {
	if width == "" {
		width = b.rightDrawerWidth
	}
	portalName := ctx.Param(ParamTargetPortal)
	p := actions.RightDrawer.PortalName()
	if portalName != "" {
		p = portalName
	}
	r.UpdatePortal(p, comp)
}

// 				Attr("@input", "alert(plaidForm.dirty) && !confirm('You have unsaved changes on this form. If you close it, you will lose all unsaved changes. Are you sure you want to close it?') ? vars.presetsDialog = true : vars.presetsDialog = $event").

type LayoutConfig struct {
	SearchBoxInvisible          bool
	NotificationCenterInvisible bool
}

func (b *Builder) notificationCenter(ctx *web.EventContext) (er web.EventResponse, err error) {
	total := b.notificationCountFunc(ctx)
	content := b.notificationContentFunc(ctx)
	icon := VIcon("mdi-bell-outline").Size(20).Color("grey-darken-1")
	er.Body = VMenu().Children(
		h.Template().Attr("v-slot:activator", "{ props }").Children(
			VBtn("").Icon(true).Children(
				h.If(total > 0,
					VBadge(
						icon,
					).Content(total).Floating(true).Color("red"),
				).Else(icon),
			).Attr("v-bind", "props").
				Density(DensityCompact).
				Variant(VariantText),
			// .Class("ml-1")
		),
		VCard(content),
	)
	return
}

const (
	ConfirmDialogConfirmEvent     = "presets_ConfirmDialog_ConfirmEvent"
	ConfirmDialogPromptText       = "presets_ConfirmDialog_PromptText"
	ConfirmDialogDialogPortalName = "presets_ConfirmDialog_DialogPortalName"
)

// for pages outside the default presets layout
func (b *Builder) PlainLayout(in web.PageFunc) (out web.PageFunc) {
	return func(ctx *web.EventContext) (pr web.PageResponse, err error) {
		b.InjectAssets(ctx)

		var innerPr web.PageResponse
		innerPr, err = in(ctx)
		if err == perm.PermissionDenied {
			pr.Body = h.Text(MustGetMessages(ctx.Context()).ErrPermissionDenied.Error())
			return pr, nil
		}
		if err != nil {
			panic(err)
		}

		pr.PageTitle = fmt.Sprintf("%s - %s", innerPr.PageTitle, i18n.T(ctx.Context(), ModelsI18nModuleKey, b.brandTitle))
		pr.Body = VApp(
			web.Portal().Name(actions.Dialog.PortalName()),
			web.Portal().Name(DeleteConfirmPortalName),
			web.Portal().Name(DefaultConfirmDialogPortalName),

			VProgressLinear().
				Attr(":active", "isFetching").
				Attr("style", "position: fixed; z-index: 99").
				Indeterminate(true).
				Height(2).
				Color(b.progressBarColor),
			h.Template(
				VSnackbar(h.Text("{{vars.presetsMessage.message}}")).
					Attr("v-model", "vars.presetsMessage.show").
					Attr(":color", "vars.presetsMessage.color").
					Timeout(2000).
					Location(LocationTop).
					ZIndex(1000000),
			).Attr("v-if", "vars.presetsMessage"),
			VMain(
				innerPr.Body.(h.HTMLComponent),
			),
		).
			Attr("id", "vt-app").
			Attr(web.VAssign("vars", `{presetsDialog: false, presetsMessage: {show: false, color: "success", 
message: ""}}`)...)

		return
	}
}

func (b *Builder) FormatHtmlValue(v string) string {
	return strings.Replace(strings.Replace(v, "{{prefix}}", b.prefix, -1), "{{exe_mtime}}", web.ExeMTime, -1)
}

func (b *Builder) InjectAssets(ctx *web.EventContext) {
	ctx.Injector.HeadHTML(b.FormatHtmlValue(`
			<link rel="stylesheet" href="{{prefix}}/assets/main.css?{{exe_mtime}}" async>
			<link rel="stylesheet" href="{{prefix}}/vuetify/assets/index.css?{{exe_mtime}}" async>
			<script src="{{prefix}}/assets/vue.js?{{exe_mtime}}"></script>
			<style>
				[v-cloak] {
					display: none;
				}
				.vx-list-item--active {
					position: relative;
				}
				.vx-list-item--active:after {
					opacity: .12;
					background-color: currentColor;
					bottom: 0;
					content: "";
					left: 0;
					pointer-events: none;
					position: absolute;
					right: 0;
					top: 0;
					transition: .3s cubic-bezier(.25,.8,.5,1);
					line-height: 0;
				}
				.vx-list-item--active:hover {
					background-color: inherit !important;
				}
			</style>
		`))

	b.InjectExtraAssets(ctx)

	ctx.Injector.TailHTML(b.FormatHtmlValue(`
			<script src="{{prefix}}/assets/main.js?{{exe_mtime}}"></script>
			`))

	if b.assetFunc != nil {
		b.assetFunc(ctx)
	}
}

func (b *Builder) InjectExtraAssets(ctx *web.EventContext) {
	for _, ea := range b.extraAssets {
		if len(ea.refTag) > 0 {
			ctx.Injector.HeadHTML(ea.refTag)
			continue
		}

		if strings.HasSuffix(ea.path, "css") {
			ctx.Injector.HeadHTML(fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\">", b.extraFullPath(ea)))
			continue
		}

		ctx.Injector.HeadHTML(fmt.Sprintf("<script src=\"%s\"></script>", b.extraFullPath(ea)))
	}
}

func (b *Builder) defaultHomePageFunc(ctx *web.EventContext) (r web.PageResponse, err error) {
	r.Body = h.Div().Text("home")
	return
}

func (b *Builder) getHomePageFunc() web.PageFunc {
	if b.homePageFunc != nil {
		return b.homePageFunc
	}
	return b.defaultHomePageFunc
}

func (b *Builder) DefaultNotFoundPageFunc(ctx *web.EventContext) (r web.PageResponse, err error) {
	msgr := MustGetMessages(ctx.Context())
	r.Body = h.Div(
		h.H1("404").Class("mb-2"),
		h.Text(msgr.NotFoundPageNotice),
	).Class("text-center mt-8")
	return
}

func (b *Builder) getNotFoundPageFunc() web.PageFunc {
	pf := b.DefaultNotFoundPageFunc
	if b.notFoundFunc != nil {
		pf = b.notFoundFunc
	}
	return pf
}
func (b *Builder) SkipNotFoundHandlerFunc(f func(r *http.Request) bool) *Builder {
	b.skipNotFoundHandler = f
	return b
}

func (b *Builder) extraFullPath(ea *extraAsset) string {
	return b.prefix + "/extra" + ea.path
}

func (b *Builder) Build(mux ...*http.ServeMux) {
	var (
		mx    *http.ServeMux
		names = make(map[string]any)
		mns   = modelNames(b.models)
	)

	for _, mx = range mux {
	}

	if mx == nil {
		mx = http.NewServeMux()
	}

	for _, mn := range mns {
		if names[mn] != nil {
			panic(fmt.Sprintf("Duplicated model name %q", mn))
		}
		names[mn] = nil
	}

	b.SetupRoutes(mx)
	b.permissions = b.BuildPermissions()
}

func (b *Builder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	web.ParseRequest(r)

	if b.handler == nil {
		b.Build()
	}
	redirectSlashes(b.handler).ServeHTTP(w, r)
}

func (b *Builder) Models() []*ModelBuilder {
	return b.models
}

func (b *Builder) Permissions() *PermMenu {
	return b.permissions
}

func (b *Builder) Verifier(vf ...*perm.PermVerifierBuilder) (r *Builder) {
	b.verifiers = append(b.verifiers, vf...)
	return b
}

func (b *Builder) BindVerifiedPageFunc(vf *perm.PermVerifierBuilder, f web.PageFunc) web.PageFunc {
	if vf != nil {
		b.Verifier(vf)
		old := f
		f = func(ctx *web.EventContext) (_ web.PageResponse, err error) {
			if vf.Build(b.verifier.Spawn().WithReq(ctx.R).Do(PermFromRequest(ctx.R))).Denied() {
				err = perm.PermissionDenied
				return
			}
			return old(ctx)
		}
	}
	return f
}

func redirectSlashes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > 1 && path[len(path)-1] == '/' {
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path[:len(path)-1], r.URL.RawQuery)
			} else {
				path = path[:len(path)-1]
			}
			redirectURL := fmt.Sprintf("//%s%s", r.Host, path)
			http.Redirect(w, r, redirectURL, 301)
			return
		}
		next.ServeHTTP(w, r)
	})
}
