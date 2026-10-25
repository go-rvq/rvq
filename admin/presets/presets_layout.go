package presets

import (
	_ "embed"
	"errors"
	"fmt"
	"path"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/printer"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/x/account"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

//go:embed assets/fix-layout.css
var fixLayoutCSS []byte

// EmbeddedWindowName is the name of a window — an iframe's, name="…" — the
// admin runs in embedded in another page (a site's dialog): there it starts
// with the side menu closed, the room is the page's; the button of the
// toolbar opens it. The name stays the window's while the admin navigates in
// it.
const EmbeddedWindowName = "rvq-embedded"

func BuildBreadcrumbsComponent(b *Builder, bc *BreadcrumbsBuilder, ctx *web.EventContext, pageTitle string) (breadcrumb h.HTMLComponent) {
	home := &Breadcrumb{
		Label: i18n.T(ctx.Context(), ModelsI18nModuleKey, b.brandTitle),
		URI:   b.URIPrefixCtx(ctx.R.Context()) + "/",
	}

	if path.Clean(ctx.R.URL.Path) != path.Clean(b.prefix+"/") {
		msgr := MustGetMessages(ctx.Context())
		bc.Prepend(home)
		if len(pageTitle) > 0 {
			bc.Append(&Breadcrumb{Label: pageTitle})
		}

		return bc.Component(msgr.YouAreHere)
	}
	return
}

func (b *Builder) DefaultLayout(in web.PageFunc, cfg *LayoutConfig) (out web.PageFunc) {
	return func(ctx *web.EventContext) (pr web.PageResponse, err error) {
		defer func() {
			lang := ctx.Injector.GetHTMLLang()

			if len(lang) == 0 {
				// set current locale fallback
				lang = i18n.DynaFromContext(ctx.Context()).GetLanguage()
				ctx.Injector.HTMLLang(lang)
			}

			b.appendLocaleScripts(ctx)
		}()

		b.InjectAssets(ctx)

		// call CreateMenus before in(ctx) to fill the menuGroupName for modelBuilders first
		var (
			titleDisabled      bool
			breadCrumbDisabled bool
			innerPr            web.PageResponse
		)

		innerPr, err = in(ctx)

		// another account of the browser chosen (its link: SwitchAccountQuery)
		// whose user may not see this page: its home, with a notice
		if err != nil && errors.Is(err, perm.PermissionDenied) && ctx.R.URL.Query().Get(SwitchAccountQuery) == "1" {
			pr.RedirectURL = b.URIPrefixCtx(ctx.R.Context()) + "/?" + NoAccessQuery + "=1"
			return pr, nil
		}

		if ctx.W.Writed() || innerPr.RedirectURL != "" {
			return innerPr, nil
		}

		// b.RunSwitchLanguageFunc(ctx)

		m := MustGetMessages(ctx.Context())

		menu := b.CreateMenus(ctx)
		if err != nil {
			innerPr.PageTitle = m.Error
			innerPr.Body = VAlert(h.Text(err.Error())).Icon("$error").Color("error").Title(m.Error)
			err = nil
			breadCrumbDisabled = true
			titleDisabled = true
		}

		var profile h.HTMLComponent
		if b.profileFunc != nil {
			profile = b.profileFunc(ctx)
		}

		// the top bar: opens and closes the menu, the brand, the user (their
		// button, their notifications: the ProfileFunc); the progress of the
		// requests at its foot
		topBar := VAppBar(
			h.Div(
				VBtn("").
					Icon("mdi-menu").
					Variant(VariantText).
					Attr("title", m.MenuToggle).
					Attr("aria-label", m.MenuToggle).
					Attr("@click.stop", "vars.navDrawer = !vars.navDrawer"),
				h.Div(b.RunBrandFunc(ctx)).Class("rvq-app-bar__brand"),
				VSpacer(),
				h.Div(b.topBarItemsOf(ctx), profile).Class("rvq-app-bar__user d-flex align-center ga-1"),
				// the theme kept in the browser, taken again as the page opens
				// (the toggle of the menu of the user mounts only with it open)
				h.Div(h.Tag("vx-theme-toggle").Attr("storage-key", UserMenuThemeStorageKey)).Style("display: none"),
			).Class("d-flex align-center ga-2 px-2 w-100"),
			VProgressLinear().
				Attr(":active", "isFetching").
				Absolute(true).
				Location(LocationBottom).
				Indeterminate(true).
				Height(2).
				Color(b.progressBarColor),
		).
			Elevation(0).
			Height(56).
			Class("rvq-app-bar border-b noprint")

		// a page of the account N of a browser holding several (/admin/u/N/…)
		// stays in it: its requests, its links (x/account)
		if n := account.IndexOf(ctx.R); n > 0 {
			ctx.Injector.HeadHTMLComponent("rvq-account", accountScript(b.prefix, n), true)
		}

		var helpBar h.HTMLComponent
		if b.helpBarFunc != nil {
			helpBar = b.helpBarFunc(ctx)
		}

		// showNotificationCenter := cfg == nil || !cfg.NotificationCenterInvisible
		// var notifier h.HTMLComponent
		// if b.notificationCountFunc != nil && b.notificationContentFunc != nil {
		//	notifier = web.Portal().Name(NotificationCenterPortalName).Loader(web.GET().EventFunc(actions.NotificationCenter))
		// }
		// ctx.R = ctx.R.WithContext(context.WithValue(ctx.R.Context(), ctxNotifyCenter, notifier))

		// showSearchBox := cfg == nil || !cfg.SearchBoxInvisible

		// _ := i18n.MustGetModuleMessages(ctx.R, CoreI18nModuleKey, Messages_en_US).(*Messages)

		// the page another account was switched to, which it may not see:
		// this one (its home) instead, said
		if ctx.R.URL.Query().Get(NoAccessQuery) == "1" {
			innerPr.Body = h.HTMLComponents{
				VAlert(h.Text(m.SwitchedNoAccess)).Type("warning").Variant(VariantTonal).
					Closable(true).Density(DensityCompact).Class("mx-4 mt-3").Attr("data-no-access", true),
				innerPr.Body,
			}
		}

		pr.PageTitle = innerPr.PageTitle
		pr.Actions = innerPr.Actions
		pr.Menu = innerPr.Menu

		var breadcrumbs h.HTMLComponent

		if !breadCrumbDisabled {
			if pr.PageTitle != "" {
				breadcrumbs = BuildBreadcrumbsComponent(b, GetOrInitBreadcrumbs(ctx.R), ctx, pr.PageTitle)
			}
		}

		var flash h.HTMLComponent
		if ctx.Flash != nil {
			flash = RenderFlash(ctx.Flash)
			ctx.Flash = nil
		}

		portals := append(GetPortals(ctx),
			web.Portal().Name(actions.RightDrawer.PortalName()),
			web.Portal().Name(actions.LeftDrawer.PortalName()),
			web.Portal().Name(actions.TopDrawer.PortalName()),
			web.Portal().Name(actions.BottomDrawer.PortalName()),
			web.Portal().Name(actions.StartDrawer.PortalName()),
			web.Portal().Name(actions.EndDrawer.PortalName()),
			web.Portal().Name(actions.Dialog.PortalName()),
			web.Portal().Name(DeleteConfirmPortalName),
			web.Portal().Name(DefaultConfirmDialogPortalName),
			web.Portal().Name(ListingDialogPortalName),
			// Where the login shows up when a session dies under an open page —
			// at the layout root, so it covers everything (see LoginPortalName).
			web.Portal().Name(LoginPortalName),
			web.Portal(flash).Name(FlashPortalName),
		)

		var menuCloser h.HTMLComponent

		if menu := append(pr.Menu, GetMenuComponent(ctx)...); len(menu) > 0 {
			menuCloser = PageBarIconButton("mdi-dots-vertical", m.MoreActions, false).
				Attr("@click.menu", "vars.contentPageMenu = !vars.contentPageMenu")

			innerPr.Body = h.HTMLComponents{
				vx.VXNavigationDrawer(menu).
					// Attr("@input", "plaidForm.dirty && vars.presetsRightDrawer == false && !confirm('You have unsaved changes on this form. If you close it, you will lose all unsaved changes. Are you sure you want to close it?') ? vars.presetsRightDrawer = true: vars.presetsRightDrawer = $event"). // remove because drawer plaidForm has to be reset when UpdateOverlayContent
					Attr("v-model", "vars.contentPageMenu").
					Location(LocationRight).
					Temporary(true).
					Density(DensityCompact).
					VariantMenu(),
				// Fixed(true).
				innerPr.Body,
			}
		}

		scoped := func(comp ...h.HTMLComponent) h.HTMLComponent {
			scope := GetScope(ctx)
			if scope != nil {
				return scope.Children(comp...)
			}
			if len(comp) == 1 {
				return comp[0]
			}
			return h.HTMLComponents(comp)
		}

		actions := h.SimplifyItems(append(pr.Actions, GetActionsComponent(ctx)...))

		if !b.printButtonDisabled && (cfg == nil || !cfg.PrintInvisible) {
			actions = append(actions, VMenu(
				web.Slot(
					VTooltip(
						web.Slot(
							PageBarIconButton("mdi-printer", "", false).
								Class("noprint").
								Attr("v-bind", "Vue.mergeProps(menu, tooltip)"),
						).
							Name("activator").
							Scope("{props:tooltip}"),
					).Text(m.PrinterOptions.Title).Location(LocationBottom),
				).
					Name("activator").
					Scope("{props:menu}"),
				VList(
					VListItem().Title(m.PrinterOptions.WithHeaders).Attr("@click", `vars.print(vars, ".rvq-page")`),
					VListItem().Title(m.PrinterOptions.WithoutHeaders).Attr("@click", `vars.print(vars, ".rvq-page-content")`),
				).Density(DensityCompact),
			))
		}

		for _, action := range actions {
			if t, _ := action.(h.TagGetter); t != nil {
				t.GetHTMLTagBuilder().Class("noprint")
			}
		}

		hasHeader := (!titleDisabled && innerPr.PageTitle != "") || len(actions) > 0 || menuCloser != nil

		pr.Body = VApp(
			portals,

			// App(true).
			// Fixed(true),
			// ClippedLeft(true),

			VSnackbar(
				h.Text(m.CopiedToClipboard),
			).
				Attr("v-model", "copiedToClipboard.value").
				Attr("color", "info").
				Location(LocationTop),

			VSnackbar(
				h.Text("{{vars.presetsMessage.message||''}}"),
				h.Div().Attr("v-if", "vars.presetsMessage.htmlMessage").Attr("v-html", "vars.presetsMessage.htmlMessage"),
				web.Slot(
					VBtn("").
						Icon("mdi-close").
						Attr("@click", `vars.presetsMessage.show = false`),
				).Name("actions"),
			).
				ZIndex(1000000).
				Attr("v-model", "vars.presetsMessage.show").
				Attr(":color", "vars.presetsMessage.color").
				Attr("v-if", "vars.presetsMessage").
				Timeout(3000).
				Location(LocationTop),

			// the top bar, the first of the layout: over the whole width
			topBar,
			// the help bar: at the right edge, the first of the layout's drawers
			// there — the drawers of the page open beside it
			helpBar,
			vx.VXNavigationDrawer(
				// Configurable items before and after the model menu (the language
				// selector is a default pre-menu item).
				// the filter and what follows the menu: with the menu whole
				h.Div(b.preMenu(ctx)).Attr("v-show", "!vars.navRail || vars.navMobile"),
				menu,
				h.Div(b.postMenu(ctx)).Attr("v-show", "!vars.navRail || vars.navMobile"),
				// ).Class("ma-2").
				// 	Style("height: calc(100% - 20px); border: 1px solid grey"),
			).
				// its foot: its icons alone, or whole again (menu_rail.go)
				SlotBottom(navRailToggle(m)).
				Width("320").
				// its right edge drags its width, kept in the browser
				Resizable("rvq.presets.navDrawer.width").
				MinWidth(220).
				MaxWidth(600).
				Attr("v-model", "vars.navDrawer").
				VariantMenu().
				// on a narrow screen (vars.navMobile) the menu goes over the page —
				// closed by a click out of it —, not beside it, which would leave
				// the page no room
				ContainerProps(`{permanent: !vars.navMobile, temporary: vars.navMobile, floating:true, elevation:0, class:"border-e",
					rail: !!vars.navRail && !vars.navMobile, railWidth: `+fmt.Sprint(navRailWidth)+`}`),
			VMain(
				scoped(
					h.If(hasHeader,
						VAppBar(
							h.Div(
								// the title; the actions at the end, with it or without
								// (a page not found has none)
								h.Div(h.If(!titleDisabled && innerPr.PageTitle != "",
									VToolbarTitle(innerPr.PageTitle),
								)).Class("mr-auto"),
								actions,
								menuCloser,
							).Class("d-flex align-center mx-2 border-b w-100 rvq-page-header").Style("height: 48px"),
						).Elevation(0).Class("rvq-page__toolbar"),
					),
					h.If(breadcrumbs != nil, VContainer(breadcrumbs).Class("rvq-page-breadcrumbs").Fluid(true).Style("padding-top:0;padding-bottom:0")),
					h.Div(innerPr.Body).Class("rvq-page-content"),
				),
			).Class("v-main__page_content").Class("rvq-page"),
			vx.VXDialog(
				h.Form(
					h.Input("options").Type("hidden").Value(h.JSONString(map[string]any{
						"PrintButton": strings.ReplaceAll(printButton, "{{title}}", m.PrinterOptions.Print),
					})),
					h.Input("body").Type("hidden").Attr(":value", "vars.printer.body"),
				).
					ID("rvq-printer__form").
					Attr(":action", fmt.Sprintf("window.location.pathname+%q", printer.PrintPath)).
					Method("post").
					Target("rvq-printer__frame"),
				h.Iframe().
					Name("rvq-printer__frame").
					Style("height: 100%; border:none"),
			).
				Closable(true).
				Expandable(true).
				Title(m.PrinterOptions.Preview).
				Density(DensityCompact).
				Attr("height", "90%").
				Attr("v-model", "vars.printer.show").
				Attr("@open", `window.setTimeout(() => window.document.getElementById("rvq-printer__form").submit(), 10)`),
		).Attr("id", "vt-app")

		pr.Body = innerPr.Wrapers.Wrap(pr.Body)

		pr.Body = vue.UserComponent(pr.Body).
			AssignMany("vars", `{
	presetsRightDrawer: false, 
	presetsLeftDrawer: false,
	presetsTopDrawer: false,
	presetsBottomDrawer: false,
	presetsStartDrawer: false,
	presetsEndDrawer: false,  
	presetsDialog: false, 
	presetsListingDialog: false,
	// a narrow screen (Vuetify's md breakpoint): the menu starts hidden
	navMobile: window.innerWidth < 960,
	navDrawer: window.name !== "`+EmbeddedWindowName+`" && window.innerWidth >= 960,
	// the main menu as its icons alone (kept in the browser: menu_rail.go)
	navRail: (() => { try { return window.localStorage.getItem("`+NavRailStorageKey+`") === "1" } catch (e) { return false } })(),
	contentPageMenu: false,
	menuFilter: "",
	printer: {show:false},
	print: (vars, selector) => {
		const d = window.document;
		let e = d.querySelector(selector),
			body = [];
	
		body[body.length] = e.outerHTML;
		d.querySelectorAll("link").forEach((e) => {
			if (e.getAttribute("rel") === "stylesheet") {
				body[body.length] = e.outerHTML
			}
		})
	
		d.querySelectorAll("style").forEach((e) => {
			body[body.length] = e.outerHTML
		})
	body[body.length] = "<style>.rvq-page {padding-left:0} .rvq-page__toolbar{left:0 !important;width:100%!important} </style>"
		vars.printer = {show: true, body: body.join("\n") };
	}
}`)
		return
	}
}

const printButton = `<div style="position: fixed; right: 20px; top: 20px; z-index:100000000">
  <button class="v-btn v-btn--elevated v-btn--icon v-theme--light bg-primary v-btn--density-default v-btn--size-default v-btn--variant-elevated"
          title="{{title}}"
          type="button" onclick="window.print()"><span class="v-btn__overlay"></span><span class="v-btn__underlay"></span><span
    class="v-btn__content" data-no-activator=""><i
    aria-hidden="true" class="mdi-printer mdi v-icon notranslate v-theme--light v-icon--size-default"></i></span></button>
</div>`
