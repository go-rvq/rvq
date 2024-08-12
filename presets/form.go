package presets

import (
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	. "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
)

type FormConfigure interface {
	ConfigureForm(f *Form)
}

type FormBuilder struct {
	id          string
	fb          *FieldsBuilder
	ctx         *web.EventContext
	msgr        *Messages
	mode        FieldMode
	mb          *ModelBuilder
	obj         interface{}
	portalName  string
	overlayMode actions.OverlayMode
}

func NewFormBuilder(ctx *web.EventContext, mb *ModelBuilder, fb *FieldsBuilder, obj interface{}) *FormBuilder {
	f := &FormBuilder{
		id:          vx.ObjectID(obj),
		fb:          fb,
		ctx:         ctx,
		msgr:        MustGetMessages(ctx.R),
		mb:          mb,
		obj:         obj,
		overlayMode: GetOverlay(ctx),
	}
	if f.id == "" {
		f.mode = NEW
	} else {
		f.mode = EDIT
	}
	return f
}

func (f *FormBuilder) Mode() FieldMode {
	return f.mode
}

func (f *FormBuilder) SetMode(mode FieldMode) *FormBuilder {
	f.mode = mode
	return f
}

func (f *FormBuilder) PortalName() string {
	return f.portalName
}

func (f *FormBuilder) SetPortalName(portalName string) {
	f.portalName = portalName
}

func (f *FormBuilder) Build() (form *Form) {
	ctx := f.ctx

	form = &Form{
		b:      f,
		Obj:    f.obj,
		MB:     f.mb,
		Portal: f.portalName,
	}

	overlayType := f.overlayMode

	if f.portalName == "" {
		form.Portal = actions.OverlayMode(overlayType).PortalName()
	}

	{
		var (
			text  string
			color string
		)

		if msg, ok := ctx.Flash.(string); ok {
			if len(msg) > 0 {
				text = msg
				color = "success"
			}
		}
		vErr, ok := ctx.Flash.(*web.ValidationErrors)
		if ok {
			gErr := vErr.GetGlobalError()
			if len(gErr) > 0 {
				text = gErr
				color = "error"
			}
		}
		if text != "" {
			form.Notice = web.Scope(
				VSnackbar(
					h.Text(text),
					web.Slot(
						VBtn("").Variant("text").
							Attr("@click", "locals.show = false").
							Children(VIcon("mdi-close")),
					).Name("actions"),
				).Location("top").
					Timeout(-1).
					Color(color).
					Attr("v-model", "locals.show"),
			).VSlot("{ locals }").Init(`{ show: true }`)
		}
	}

	form.Body = f.fb.ToComponent(f.mb.Info(), f.obj, FieldModeStack{f.mode}, ctx)

	return
}

type Form struct {
	b            *FormBuilder
	Obj          interface{}
	MB           *ModelBuilder
	AutoSave     bool
	Portal       string
	PortalLoader *web.VueEventTagBuilder
	Title,
	Notice,
	Body h.HTMLComponent
	Tabs      []TabComponentFunc
	SidePanel ObjectComponentFunc

	MainPortals   h.HTMLComponents
	PrimaryAction h.HTMLComponent
	Menu          h.HTMLComponents
	Actions       h.HTMLComponents
}

func (f *Form) Component() (comp h.HTMLComponent) {
	comp = defaultToPage(commonPageConfig{
		main:      web.Scope(f.Body),
		tabPanels: f.Tabs,
		sidePanel: f.SidePanel,
	}, f.Obj, f.b.ctx)

	var (
		overlay        = f.b.overlayMode
		primaryActions []h.HTMLComponent
	)

	if f.b.ctx.R.FormValue(ParamActionsDisabled) != "true" {
		if f.PrimaryAction != nil {
			primaryActions = append(primaryActions, f.PrimaryAction)
		}

		actionsElements := f.Actions

		if !overlay.Overlayed() && f.PrimaryAction != nil {
			actionsElements = append(actionsElements, primaryActions...)
			primaryActions = nil
		}

		if len(actionsElements) > 0 {
			comp = h.Components(
				VCardText(comp),
				// h.If(!f.AutoSave, VCardActions(actionButtons)),
				VCardActions(append([]h.HTMLComponent{VSpacer()}, actionsElements...)...),
			)
		}

		if overlay.Overlayed() {
			primaryActions = append(primaryActions, h.If(!f.AutoSave,
				VBtn("").
					Variant(VariantFlat).
					Icon(true).
					Children(
						VIcon("mdi-close"),
					).Attr("@click.stop", "closer.show = false"),
			))
		}
	}

	if !overlay.Overlayed() {
		return
	}
	// b.RowMenu().listingItemFuncs(ctx)...

	appBarComps := append([]h.HTMLComponent{
		VToolbarTitle("").Class("pl-2").
			Children(f.Title),
		VSpacer(),
	}, primaryActions...)

	comp = VSheet(
		VCard(comp).Variant(VariantFlat),
	).Class("pa-2")

	if len(f.Menu) > 0 {
		appBarComps = append(h.HTMLComponents{
			VBtn("").
				Variant(VariantFlat).
				Icon(true).
				Density("compact").
				Children(
					VIcon("mdi-menu"),
				).Attr("@click.menu", "locals.menu = !locals.menu"),
		}, appBarComps...)

		comp = h.HTMLComponents{
			VNavigationDrawer(f.Menu).
				// Attr("@input", "plaidForm.dirty && vars.presetsRightDrawer == false && !confirm('You have unsaved changes on this form. If you close it, you will lose all unsaved changes. Are you sure you want to close it?') ? vars.presetsRightDrawer = true: vars.presetsRightDrawer = $event"). // remove because drawer plaidForm has to be reset when UpdateOverlayContent
				Class("v-navigation-drawer--temporary").
				Attr("v-model", "locals.menu").
				Location(LocationLeft).
				Temporary(true).
				// Fixed(true).
				Attr(":height", `"100%"`),
			comp,
		}
	}

	comp = web.Scope(
		f.Notice,
		VLayout(
			h.If(!f.MB.singleton,
				VAppBar(appBarComps...).Color("white").Elevation(0),
			),
			VMain(
				comp,
			),
		).Attr(":height", `"100%"`)).VSlot("{ form, locals, vars, closer }").Init("{menu: false}")
	return
}

func (f *Form) FullComponent() (comp h.HTMLComponent) {
	return append(f.MainPortals, f.Component())
}

func (f *Form) Respond(r *web.EventResponse) {
	comp := f.FullComponent()

	switch f.b.overlayMode {
	case actions.Dialog:
		f.MB.p.Dialog().
			SetTargetPortal(f.Portal).
			SetContentPortalName(f.Portal+"Content").
			Respond(r, comp)
	default:
		if f.Portal != "" && f.b.overlayMode.IsDrawer() {
			d := f.MB.p.Drawer(f.b.overlayMode).
				SetPortalName(f.Portal).
				SetValidWidth(f.MB.rightDrawerWidth)

			d.Respond(r, comp)
		} else {
			r.Body = comp
		}
	}
}
