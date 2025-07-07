package presets

import (
	"path"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/web/v3/vue"
	"github.com/qor5/x/v3/i18n"
	. "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
)

func (b *Builder) DefaultLayout(in web.PageFunc, cfg *LayoutConfig) (out web.PageFunc) {
	return func(ctx *web.EventContext) (pr web.PageResponse, err error) {
		b.InjectAssets(ctx)

		// call CreateMenus before in(ctx) to fill the menuGroupName for modelBuilders first
		var (
			titleDisabled      bool
			breadCrumbDisabled bool
			innerPr            web.PageResponse
		)

		innerPr, err = in(ctx)

		if ctx.W.Writed() || innerPr.RedirectURL != "" {
			return innerPr, nil
		}

		// b.RunSwitchLanguageFunc(ctx)

		toolbar := h.Div(
			b.RunBrandFunc(ctx),
			h.Div(
				VBtn("").
					Icon("mdi-menu").
					Density(DensityCompact).
					Variant(VariantText).
					Attr("@click", "vars.navDrawer = !vars.navDrawer").Density(DensityCompact),
			).Class(H75),
		).Class("d-flex ga-1 v-card-text")

		menu := b.CreateMenus(ctx)
		if err != nil {
			title := MustGetMessages(ctx.Context()).Error
			innerPr.PageTitle = title
			innerPr.Body = VAlert(h.Text(err.Error())).Icon("$error").Color("error").Title(title)
			err = nil
			breadCrumbDisabled = true
			titleDisabled = true
		}

		var profile h.HTMLComponent
		if b.profileFunc != nil {
			profile = b.profileFunc(ctx)
		}

		// showNotificationCenter := cfg == nil || !cfg.NotificationCenterInvisible
		// var notifier h.HTMLComponent
		// if b.notificationCountFunc != nil && b.notificationContentFunc != nil {
		//	notifier = web.Portal().Name(NotificationCenterPortalName).Loader(web.GET().EventFunc(actions.NotificationCenter))
		// }
		// ctx.R = ctx.R.WithContext(context.WithValue(ctx.R.Context(), ctxNotifyCenter, notifier))

		// showSearchBox := cfg == nil || !cfg.SearchBoxInvisible

		// _ := i18n.MustGetModuleMessages(ctx.R, CoreI18nModuleKey, Messages_en_US).(*Messages)

		pr.PageTitle = innerPr.PageTitle
		pr.Actions = innerPr.Actions
		pr.Menu = innerPr.Menu

		var breadcrumbs []h.HTMLComponent

		if !breadCrumbDisabled {
			if pr.PageTitle != "" {
				bc := GetOrInitBreadcrumbs(ctx.R)
				home := &Breadcrumb{
					Label: i18n.T(ctx.Context(), ModelsI18nModuleKey, b.brandTitle),
					URI:   b.prefix + "/",
				}
				if path.Clean(ctx.R.URL.Path) != path.Clean(home.URI) {
					msgr := MustGetMessages(ctx.Context())
					bc.Prepend(home)
					bc.Append(&Breadcrumb{Label: pr.PageTitle})
					breadcrumbs = append(breadcrumbs, bc.Component(msgr.YouAreHere))
				}
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
			web.Portal(flash).Name(FlashPortalName),
		)

		var menuCloser h.HTMLComponent

		if menu := append(pr.Menu, GetMenuComponent(ctx)...); len(menu) > 0 {
			menuCloser = VBtn("").
				Variant(VariantFlat).
				Icon(true).
				Density(DensityComfortable).
				Children(
					VIcon("mdi-menu"),
				).Attr("@click.menu", "vars.contentPageMenu = !vars.contentPageMenu")

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

		msg := MustGetMessages(ctx.Context())
		actions := h.SimplifyItems(append(pr.Actions, GetActionsComponent(ctx)...))

		hasHeader := (!titleDisabled && innerPr.PageTitle != "") || len(actions) > 0 || menuCloser != nil

		pr.Body = VApp(
			portals,

			// App(true).
			// Fixed(true),
			// ClippedLeft(true),

			VSnackbar(
				h.Text(msg.CopiedToClipboard),
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

			vx.VXNavigationDrawer(
				// b.RunBrandProfileSwitchLanguageDisplayFunc(b.RunBrandFunc(ctx), profile, b.RunSwitchLanguageFunc(ctx), ctx),
				// b.RunBrandFunc(ctx),
				// profile,
				menu,
				// ).Class("ma-2").
				// 	Style("height: calc(100% - 20px); border: 1px solid grey"),
			).
				SlotTop(toolbar).
				SlotBottom(profile).
				Width("320").
				Attr("v-model", "vars.navDrawer").
				VariantMenu().
				ContainerProps(`{permanent:true, floating:true, elevation:0, class:"border-e"}`),
			VMain(
				scoped(
					h.If(hasHeader,
						VAppBar(
							h.Div(
								VProgressLinear().
									Attr(":active", "isFetching").
									Class("ml-4").
									Attr("style", "position: fixed; z-index: 99;").
									Indeterminate(true).
									Height(2).
									Color(b.progressBarColor),
								VAppBarNavIcon().
									Density("compact").
									Class("mr-2").
									Attr("v-if", "!vars.navDrawer").
									On("click.stop", "vars.navDrawer = !vars.navDrawer"),

								h.If(!titleDisabled && innerPr.PageTitle != "", h.Div(
									VToolbarTitle(innerPr.PageTitle), // Class("text-h6 font-weight-regular"),
								).Class("mr-auto")),
								actions,
								menuCloser,
							).Class("d-flex align-center mx-2 border-b w-100").Style("height: 48px"),
						).Elevation(0),
					),
					h.If(len(breadcrumbs) > 0, VContainer(breadcrumbs...).Fluid(true).Style("padding-top:0;padding-bottom:0")),
					innerPr.Body,
				),
			).Class("v-main__page_content"),
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
navDrawer: true,
contentPageMenu: false
}`)

		return
	}
}
