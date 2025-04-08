package presets

import (
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	. "github.com/qor5/x/v3/ui/vuetify"
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

			menu    = b.CreateMenus(ctx)
			toolbar = VContainer(
				VRow(
					VCol(b.RunBrandFunc(ctx)).Cols(8),
					VCol(
						b.RunSwitchLanguageFunc(ctx),
						// VBtn("").Children(
						//	languageSwitchIcon,
						//	VIcon("mdi-menu-down"),
						// ).Attr("variant", "plain").
						//	Attr("icon", ""),
					).Cols(2),

					VCol(
						VAppBarNavIcon().Attr("icon", "mdi-menu").
							Class("text-grey-darken-1").
							Attr("@click", "vars.navDrawer = !vars.navDrawer").Density(DensityCompact),
					).Cols(2),
				).Attr("align", "center").Attr("justify", "center"),
			)
		)

		innerPr, err = in(ctx)
		if ctx.W.Writed() {
			return
		}
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
			profile = VAppBar(
				b.profileFunc(ctx),
			).Location("bottom").Class("border-t-sm border-b-0").Elevation(0)
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

		var breadcrumbs []h.HTMLComponent

		if !breadCrumbDisabled {
			if pr.PageTitle != "" {
				bc := GetOrInitBreadcrumbs(ctx.R)
				home := &Breadcrumb{
					Label: i18n.T(ctx.Context(), ModelsI18nModuleKey, b.brandTitle),
					URI:   b.prefix + "/",
				}
				if pr.PageTitle != home.Label {
					msgr := MustGetMessages(ctx.Context())
					bc.Prepend(home)
					bc.Append(&Breadcrumb{Label: pr.PageTitle})
					breadcrumbs = append(breadcrumbs, bc.Component(msgr.YouAreHere))
				}
			}
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
			web.Portal().Name(FlashPortalName),
		)

		var menuCloser h.HTMLComponent

		if menu := GetMenuComponent(ctx); menu != nil {
			menuCloser = VBtn("").
				Variant(VariantFlat).
				Icon(true).
				Density(DensityComfortable).
				Children(
					VIcon("mdi-menu"),
				).Attr("@click.menu", "vars.contentPageMenu = !vars.contentPageMenu")

			innerPr.Body = h.HTMLComponents{
				VNavigationDrawer(menu).
					// Attr("@input", "plaidForm.dirty && vars.presetsRightDrawer == false && !confirm('You have unsaved changes on this form. If you close it, you will lose all unsaved changes. Are you sure you want to close it?') ? vars.presetsRightDrawer = true: vars.presetsRightDrawer = $event"). // remove because drawer plaidForm has to be reset when UpdateOverlayContent
					Attr("v-model", "vars.contentPageMenu").
					Location(LocationRight).
					// Fixed(true).
					Attr(":height", `"100%"`),
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

			h.Template(
				VSnackbar(
					h.Text("{{vars.presetsMessage.message}}"),
					web.Slot(
						VBtn("").
							Icon("mdi-close").
							Attr("@click", `vars.presetsMessage.show = false`),
					).Name("actions"),
				).
					Attr("v-model", "vars.presetsMessage.show").
					Attr(":color", "vars.presetsMessage.color").
					Timeout(3000).
					Location(LocationTop),
			).Attr("v-if", "vars.presetsMessage"),

			VNavigationDrawer(
				// b.RunBrandProfileSwitchLanguageDisplayFunc(b.RunBrandFunc(ctx), profile, b.RunSwitchLanguageFunc(ctx), ctx),
				// b.RunBrandFunc(ctx),
				// profile,
				VLayout(
					VMain(
						toolbar,
						VCard(
							menu,
						).Variant(VariantText),
					),
					// VDivider(),
					profile,
				).Class("ma-2 border-sm rounded-lg elevation-0"),
				// ).Class("ma-2").
				// 	Style("height: calc(100% - 20px); border: 1px solid grey"),
			).
				Width(320).
				// App(true).
				// Clipped(true).
				// Fixed(true).
				Attr("v-model", "vars.navDrawer").
				// Attr("style", "border-right: 1px solid grey ").
				Permanent(true).
				Floating(true).
				Elevation(0),

			VMain(
				scoped(VAppBar(
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

						h.If(!titleDisabled, h.Div(
							VToolbarTitle(innerPr.PageTitle), // Class("text-h6 font-weight-regular"),
						).Class("mr-auto")),
						GetActionsComponent(ctx),
						menuCloser,
					).Class("d-flex align-center mx-2 border-b w-100").Style("height: 48px"),
				).
					Elevation(0),
					h.If(len(breadcrumbs) > 0, VContainer(breadcrumbs...).Fluid(true).Style("padding-top:0;padding-bottom:0")),
					innerPr.Body),
			).Class("v-main__page_content"),
		).Attr("id", "vt-app").
			Attr(web.VAssign("vars", `{
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
}`)...)

		return
	}
}
