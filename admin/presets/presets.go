package presets

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/utils/bcp47"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/datafield"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/iancoleman/strcase"
	"go.uber.org/zap"
	"golang.org/x/text/language"
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
	preMenuItems                          []*SideMenuItem
	postMenuItems                         []*SideMenuItem
	preMenuOrder                          []string
	postMenuOrder                         []string
	notificationCountFunc                 func(ctx *web.EventContext) int
	notificationContentFunc               ComponentFunc
	brandTitle                            string
	vuetifyOptions                        string
	progressBarColor                      string
	rightDrawerWidth                      string
	printButtonDisabled                   bool
	writeFieldDefaults                    *FieldDefaults
	listFieldDefaults                     *FieldDefaults
	detailFieldDefaults                   *FieldDefaults
	extraAssets                           []*extraAsset
	assetFunc                             AssetFunc
	menuTree                              *MenuGroupBuilder
	menuItems                             map[string]*MenuItem
	wrapHandlers                          map[string]func(in http.Handler) (out http.Handler)
	plugins                               []Plugin
	ModelConfigurators                    ModelConfigurators
	ModelSetupFactories                   ModelSetupFactories
	skipNotFoundHandler                   func(r *http.Request) bool
	muxSetup                              []func(prefix string, r *http.ServeMux)
	permissions                           *PermMenu
	verifiers                             perm.PermVerifiers
	pagesRegistrator                      *PagesRegistrator
	formSigner                            FormSigner
	recordUserFinder                      RecordUserFinder

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
		// signs the form values the user must not change (the record stamp, see
		// record_stamp.go). Random key: SetFormSigner to survive a restart or to
		// share it between instances.
		formSigner: NewHMACFormSigner(nil),
		skipNotFoundHandler: func(r *http.Request) bool {
			return false
		},
	})
	r.GetWebBuilder().RegisterEventHandler(EventOpenConfirmDialog, web.EventFunc(r.openConfirmDialog))
	r.layoutFunc = r.DefaultLayout
	r.detailLayoutFunc = r.DefaultLayout
	r.pagesRegistrator = NewPagesRegistrator(
		r,
		func() string {
			return r.prefix
		},
		func(pf web.PageFunc, do DoPageBuilder) http.Handler {
			return r.Wrap(pf, do)
		},
		RequestPermVerifierFunc(func(*http.Request) *perm.Verifier {
			return r.verifier
		}),
	).LayoutFunc(func(config *LayoutConfig, f func(ctx *web.EventContext) (r web.PageResponse, err error)) web.PageFunc {
		if config == nil {
			config = r.homePageLayoutConfig
		}
		return r.layoutFunc(f, config)
	})

	// Default pre-menu items (first in the side menu): the language selector
	// (invisible with a single language) and, below it, a menu filter.
	r.AddPreMenuItem(
		&SideMenuItem{
			Name:    PreMenuItemLanguageSwitch,
			Enabled: true,
			Handler: r.RunSwitchLanguageFunc,
		},
		&SideMenuItem{
			Name:    PreMenuItemMenuFilter,
			Enabled: true,
			Handler: r.RunMenuFilterFunc,
		},
	)
	return r
}

func (b *Builder) PagesRegistrator() *PagesRegistrator {
	return b.pagesRegistrator
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
	// Only a model that is in the menu takes a key there — see ModelNotInMenu.
	if !r.notInMenu {
		if _, err := b.RegisterMenuItem(MenuItemModel, r.id, r); err != nil {
			panic(err)
		}
	}
	return r
}

func (b *Builder) GetModelByID(id string) *ModelBuilder {
	for _, mb := range b.models {
		if found := findModelByID(mb, id); found != nil {
			return found
		}
	}
	return nil
}

// findModelByID returns mb or the first of its (recursive) children whose id
// matches, so nested resources (AddChild) are resolvable by id.
func findModelByID(mb *ModelBuilder, id string) *ModelBuilder {
	if mb.id == id {
		return mb
	}
	for _, child := range mb.children {
		if found := findModelByID(child, id); found != nil {
			return found
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
		if found := findModelByType(model, t); found != nil {
			return found
		}
	}
	return nil
}

// findModelByType returns mb or the first of its (recursive) children whose
// modelType matches t. This lets GetModel resolve nested resources (added via
// AddChild), so cross-model references such as field selectors keep working when
// resources are mounted under a parent (e.g. under an organization).
func findModelByType(mb *ModelBuilder, t reflect.Type) *ModelBuilder {
	if mb.modelType == t {
		return mb
	}
	for _, child := range mb.children {
		if found := findModelByType(child, t); found != nil {
			return found
		}
	}
	return nil
}

func (b *Builder) DataOperator(v DataOperator) (r *Builder) {
	b.dataOperator = v
	return b
}

// FormSigner signs the form values the user must not be able to change — today
// the edit form's record stamp (see record_stamp.go).
func (b *Builder) FormSigner() FormSigner {
	return b.formSigner
}

// SetFormSigner replaces the signer. The default one is created with a random
// key at startup, so forms rendered before a restart are no longer accepted and
// instances do not accept each other's — pass a signer with YOUR key
// (NewHMACFormSigner) when either matters.
func (b *Builder) SetFormSigner(v FormSigner) (r *Builder) {
	b.formSigner = v
	return b
}

// RecordUserFinder loads the user an `UpdatedByID` points at, to name whoever
// changed a record while somebody had its form open.
func (b *Builder) RecordUserFinder() RecordUserFinder {
	return b.recordUserFinder
}

func (b *Builder) SetRecordUserFinder(v RecordUserFinder) (r *Builder) {
	b.recordUserFinder = v
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

// MenuTree is the root of the menu: every model, page and group hangs from it,
// directly or through a group. The root is a sentinel — it has no entry of its
// own and never renders as a group.
func (b *Builder) MenuTree() *MenuGroupBuilder {
	if b.menuTree == nil {
		b.menuTree = &MenuGroupBuilder{b: b}
		b.menuItems = map[string]*MenuItem{}
	}
	return b.menuTree
}

// MenuItems is the key registry: one entry per key, whether or not its value
// has arrived.
func (b *Builder) MenuItems() map[string]*MenuItem {
	b.MenuTree()
	return b.menuItems
}

// MenuItemOf is the item a key names, creating the placeholder when the key has
// only been referenced so far. The placeholder goes to the root, in
// registration order, until something moves it.
//
// Moving is always done through the key: whether the item already carries its
// value or is still waiting for one, there is a single item per key, and that
// item is what moves.
func (b *Builder) MenuItemOf(ref MenuRef) *MenuItem {
	root := b.MenuTree()
	key := ref.Key()
	if it := b.menuItems[key]; it != nil {
		return it
	}
	it := &MenuItem{Type: ref.Type, Name: ref.Name}
	b.menuItems[key] = it
	it.parent = root
	root.items = append(root.items, it)
	return it
}

// RegisterMenuItem gives a key its value.
//
// A key is unique. Registering a value for a key that already has one is a
// mistake — two models with the same id, a page registered twice — and is
// reported instead of one silently replacing the other. A key that exists
// without a value is a place someone reserved: registering fills it where it
// stands, keeping the position the reservation gave it.
func (b *Builder) RegisterMenuItem(typ MenuItemType, name string, value any) (*MenuItem, error) {
	if value == nil {
		return nil, fmt.Errorf("presets: register menu item %s: no value", menuKey(typ, name))
	}

	it := b.MenuItemOf(MenuRef{Type: typ, Name: name})
	if it.Value != nil {
		if it.Value == value {
			return it, nil
		}
		return nil, fmt.Errorf("presets: menu item %s is already registered", it.Key())
	}
	it.Value = value
	return it, nil
}

// MoveMenuItem puts the item a key names inside dst, at index when given.
//
// The key need not be registered: an unregistered key is placed as a
// reservation and waits there for its value, so a group may name what does not
// exist yet and the order of configuration stops mattering. Calling it again
// for the same key moves the item — the last call wins, and an item is in
// exactly one place at a time.
func (b *Builder) MoveMenuItem(ref MenuRef, dst *MenuGroupBuilder, index ...int) error {
	root := b.MenuTree()
	if dst == nil {
		dst = root
	}

	it := b.MenuItemOf(ref)

	// A group cannot be moved into itself or into one of its own descendants:
	// the tree would stop being one.
	if g := it.Group(); g != nil && g.contains(dst) {
		return fmt.Errorf("presets: moving menu item %s into %q would make a cycle", it.Key(), dst.Path())
	}

	if p := it.parent; p != nil {
		if i := p.indexOf(it.Key()); i >= 0 {
			p.items = append(p.items[:i], p.items[i+1:]...)
		}
	}

	at := len(dst.items)
	if len(index) > 0 && index[0] >= 0 && index[0] < at {
		at = index[0]
	}
	dst.items = append(dst.items, nil)
	copy(dst.items[at+1:], dst.items[at:])
	dst.items[at] = it
	it.parent = dst
	return nil
}

// menuGroupOf is the group a key sits in, or nil when it sits at the root. The
// root sentinel is not a group anyone is in.
func (b *Builder) menuGroupOf(ref MenuRef) *MenuGroupBuilder {
	it := b.MenuItems()[ref.Key()]
	if it == nil {
		return nil
	}
	if p := it.Parent(); p != nil && p.item != nil {
		return p
	}
	return nil
}

// MenuGroup finds or registers the group of that name. A group named by another
// group before this call already exists as a reservation, and this fills it.
func (b *Builder) MenuGroup(name string) *MenuGroupBuilder {
	b.MenuTree()
	if it := b.menuItems[menuKey(MenuItemGroup, name)]; it != nil {
		if g := it.Group(); g != nil {
			return g
		}
	}

	g := &MenuGroupBuilder{name: name, b: b}
	it, err := b.RegisterMenuItem(MenuItemGroup, name, g)
	if err != nil {
		panic(err)
	}
	g.item = it
	return g
}

// MenuOrder places items at the root, in the order given. Each is a MenuRef, a
// *MenuGroupBuilder or — in the older form — a name.
//
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
	root := b.MenuTree()
	for _, item := range items {
		if err := b.MoveMenuItem(menuRefOf(item), root); err != nil {
			panic(err)
		}
	}
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

// menuNode is one node of the data-driven side menu (a VTreeview item): a group
// (has Children) or a leaf (Value is the target path). Props carries the item's
// component props (the prepend icon, and href on leaves).
type menuNode struct {
	Title    string         `json:"title"`
	Value    string         `json:"value"`
	Props    map[string]any `json:"props,omitempty"`
	Children []*menuNode    `json:"children,omitempty"`
}

func (b *Builder) CreateMenus(ctx *web.EventContext) (r h.HTMLComponent) {
	var (
		mMap = make(map[string]*ModelBuilder)
		pMap = make(map[string]*HttpPageBuilder)
	)

	for _, m := range b.models {
		if !m.notInMenu {
			mMap[m.id] = m
		}
	}

	for _, page := range b.pagesRegistrator.httpPages {
		if !page.notInMenu {
			pMap[page.path] = page
		}
	}

	var (
		openedGroups = []string{} // groups to auto-expand: the active item's ancestor chain
		activated    = []string{} // the active leaf's value (its path), for the highlight
		inOrderMap   = make(map[string]struct{})
	)

	// pageNode / modelNode build a leaf node (or nil when hidden/denied), record
	// it as placed, and flag the active one.
	pageNode := func(p *HttpPageBuilder) *menuNode {
		if p == nil || p.notInMenu || (p.verififer != nil && p.Verifier(ctx.R).Denied()) {
			return nil
		}
		inOrderMap[p.path] = struct{}{}
		if p.isMenuItemActive(ctx) {
			activated = []string{p.fullPath}
		}
		props := map[string]any{"href": p.fullPath}
		if p.menuIcon != "" {
			props["prependIcon"] = p.menuIcon
		}
		return &menuNode{Title: p.TTitle(ctx.Context()), Value: p.fullPath, Props: props}
	}
	modelNode := func(m *ModelBuilder) *menuNode {
		if m == nil || m.notInMenu || m.permissioner.ReqLister(ctx.R).Denied() {
			return nil
		}
		inOrderMap[m.id] = struct{}{}
		href := m.Info().ListingHref(ParentsModelID(ctx.R)...)
		if m.link != "" {
			href = m.link
		}
		if m.defaultURLQueryFunc != nil {
			href = fmt.Sprintf("%s?%s", href, m.defaultURLQueryFunc(ctx.R).Encode())
		}
		if m.isMenuItemActive(ctx) {
			activated = []string{href}
		}
		icon := m.menuIcon
		if icon == "" {
			icon = defaultMenuIcon(m.label)
		}
		return &menuNode{
			Title: m.TPageLabel(ctx.Context()),
			Value: href,
			Props: map[string]any{"href": href, "prependIcon": icon},
		}
	}

	// groupNode builds a group node (and its nested sub-groups) recursively. It
	// returns nil when the group has no visible child. active reports whether the
	// group (or a descendant) holds the active item, so each ancestor adds itself
	// to openedGroups and the whole chain auto-expands.
	var (
		groupNode func(v *MenuGroupBuilder) (node *menuNode, active bool)
		itemNodes func(items []*MenuItem) (nodes []*menuNode, active bool)
	)
	groupNode = func(v *MenuGroupBuilder) (node *menuNode, active bool) {
		groupIcon := v.icon
		if groupIcon == "" {
			groupIcon = defaultMenuIcon(v.name)
		}
		title := v.TTitle(ctx.Context())

		children, active := itemNodes(v.items)
		if len(children) == 0 {
			return nil, false
		}
		if active {
			// opened is keyed by the item value, which for a group is its key.
			openedGroups = append(openedGroups, menuKey(MenuItemGroup, v.name))
		}
		return &menuNode{
			Title:    title,
			Value:    menuKey(MenuItemGroup, v.name),
			Props:    map[string]any{"prependIcon": groupIcon},
			Children: children,
		}, active
	}

	// itemNodes renders a run of items in order. An item still waiting for its
	// value renders nothing: it holds a place, and there is nothing to show
	// until whoever owns it registers it.
	itemNodes = func(items []*MenuItem) (nodes []*menuNode, active bool) {
		for _, it := range items {
			var (
				node        *menuNode
				childActive bool
			)
			switch value := it.Value.(type) {
			case *ModelBuilder:
				node = modelNode(value)
			case *HttpPageBuilder:
				node = pageNode(value)
			case *MenuGroupBuilder:
				node, childActive = groupNode(value)
			}
			if node == nil {
				continue
			}
			nodes = append(nodes, node)
			if childActive || (len(activated) > 0 && activated[0] == node.Value) {
				active = true
			}
		}
		return
	}

	// The root's items, in order. Anything never moved sits here in the order it
	// was registered, so a model or page nobody placed still shows up.
	nodes, _ := itemNodes(b.MenuTree().items)

	// The side menu is a VTreeview: it filters (search, bound to vars.menuFilter)
	// and auto-expands matches natively, and a leaf's value is its path — so
	// activating one navigates there (SPA push-state). openedGroups pre-expands the
	// active item's ancestor chain; activated pre-highlights it. The prepend slot
	// draws each item's icon (item-props does not carry it to the row).
	tree := VTreeview(
		h.Template(
			VIcon("").Attr(":icon", "item.props.prependIcon").
				Attr("v-if", "item.props && item.props.prependIcon").
				Size(SizeSmall),
		).Attr("v-slot:prepend", "{ item }"),
	).
		Items(nodes).
		ItemTitle("title").
		ItemValue("value").
		ItemChildren("children").
		ItemProps(true).
		Activatable(true).
		ActiveStrategy("single").
		OpenStrategy("multiple").
		Density(DensityCompact).
		Slim(true).
		Class("main-menu primary--text").
		Attr("v-model:opened", "locals.opened").
		Attr("v-model:activated", "locals.activated").
		Attr(":search", "vars.menuFilter").
		// Navigate when a leaf (value is a path) is activated; group values are
		// "group:<name>" and are ignored.
		Attr("@update:activated", `(v) => { const p = Array.isArray(v) ? v[v.length-1] : v; if (p && String(p).charAt(0) === '/') { plaid().vars(vars).pushStateURL(String(p)).go(); } }`)

	r = web.Scope(tree).Slot("{ locals }").LocalsInit(
		fmt.Sprintf(`{ opened: %s}`, h.JSONString(openedGroups)),
		fmt.Sprintf(`{ activated: %s}`, h.JSONString(activated)),
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
							h.Div(h.Text(fmt.Sprintf("%s%s %s", msgr.Language, msgr.Colon, bcp47.FlagLabel(supportLanguages[0].String())))).Role("button"),
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

	_, mi := language.MatchStrings(matcher, lang, accept)

	// Each item's value is the full URL that switches to that language (same page,
	// with ?lang=…). Selecting it navigates there — a full page reload, so the
	// whole admin (the side menu included) re-renders in the chosen language, and
	// the i18n middleware persists it in the "lang" cookie. This mirrors the login
	// page's language select, which is the proven idiom.
	type langItem struct {
		Label string `json:"Label"`
		Value string `json:"Value"`
	}
	var items []langItem
	var currentURL string
	for i, tag := range supportLanguages {
		u, _ := url.Parse(ctx.R.RequestURI)
		qs := u.Query()
		// Drop the ajax event params from the current URL — the request may be a
		// Plaid reload (…?__execute_event__=__reload__), and the language link must
		// land on the plain page, not re-fire that event.
		qs.Del("__execute_event__")
		qs.Del("__reload__")
		qs.Set(queryName, tag.String())
		u.RawQuery = qs.Encode()
		if i == mi {
			currentURL = u.String()
		}
		// Label as "🏳 name (code)" (the flag, then the name and the code), so the admin-language selector reads the same way; the items
		// are already limited to the registered i18n languages.
		items = append(items, langItem{Label: bcp47.FlagLabel(tag.String()), Value: u.String()})
	}

	return VAutocomplete().
		Label(msgr.Language).
		Items(items).
		ItemTitle("Label").
		ItemValue("Value").
		ModelValue(currentURL).
		// The emitted value is the selected item's Value — the URL that switches to
		// that language. Navigating to it is a full page reload, so the whole admin
		// (side menu included) re-renders in the chosen language; the i18n
		// middleware persists it in the "lang" cookie.
		Attr("@update:model-value", "$event && (window.location.href = $event)").
		PrependInnerIcon("mdi-translate").
		Density(DensityCompact).
		Variant(VariantOutlined).
		HideDetails(true).
		// mt-2 gives the floating outlined label room at the top so it is not
		// clipped by the drawer's edge.
		Class("mx-3 mt-2")
}

// RunMenuFilterFunc renders the side-menu text filter (a pre-menu item). It binds
// to the global vars.menuFilter, which the VTreeview menu reads as its search —
// filtering and auto-expanding matches natively.
func (b *Builder) RunMenuFilterFunc(ctx *web.EventContext) (r h.HTMLComponent) {
	msgr := MustGetMessages(ctx.Context())
	return VTextField().
		Attr("v-model", "vars.menuFilter").
		Attr("placeholder", msgr.Search).
		PrependInnerIcon("mdi-magnify").
		Clearable(true).
		Density(DensityCompact).
		Variant(VariantOutlined).
		HideDetails(true).
		Class("mx-3 mt-2 mb-1")
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

	// LoginPortalName is where the login shows up when the session dies while a
	// page is open. It sits at the layout root so the dialog covers everything,
	// and the page underneath keeps whatever the user had typed.
	LoginPortalName = "presets_LoginPortalName"

	// LoginDoneURI is the tiny page the login lands on inside the dialog: it
	// tells the page around it that the session is back, and that is the whole
	// content — nothing heavy is loaded inside the frame.
	LoginDoneURI = "/login-done"

	// LoginDialogURI answers with the dialog itself. It is what a page asks for
	// after being told its session is gone (web.EventResponse.LoginURI), so it
	// must be reachable WITHOUT a session — it is whitelisted in the login
	// middleware.
	LoginDialogURI = "/login-dialog"
)

// LoginDialogVar is the layout variable that shows the login dialog.
const LoginDialogVar = "vars.presetsLoginDialog"

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
		b.contentDrawer(ctx, r, comp)
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

// contentDrawer puts comp in the target portal (the right drawer's by default),
// which is already open: its width is the portal's.
func (b *Builder) contentDrawer(ctx *web.EventContext, r *web.EventResponse, comp h.HTMLComponent) {
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
		defer func() {
			lang := ctx.Injector.GetHTMLLang()
			if len(lang) == 0 {
				lang = i18n.DynaFromContext(ctx.Context()).GetLanguage()
				ctx.Injector.HTMLLang(lang)
			}
		}()

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
				innerPr.Body,
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
	r = web.ParseRequest(r)

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
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}
