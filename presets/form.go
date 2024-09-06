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
		form.Portal = overlayType.PortalName()
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
	Title        string
	Notice,
	Body h.HTMLComponent
	Tabs []TabComponentFunc

	TopLeftActions  h.HTMLComponents
	TopRightActions h.HTMLComponents

	MainPortals   h.HTMLComponents
	PrimaryAction h.HTMLComponent
	Menu          h.HTMLComponents
	Actions       h.HTMLComponents
	ScopeDisabled bool
}

func (f *Form) Component() (comp h.HTMLComponent) {
	var cb = &ContentComponentBuilder{
		Obj:     f.Obj,
		Context: f.b.ctx,
		Tabs:    f.Tabs,
		Title:   f.Title,
		Menu:    f.Menu,
		Body:    f.Body,

		TopLeftActions:  f.TopLeftActions,
		TopRightActions: f.TopRightActions,
		MainPortals:     f.MainPortals,

		Overlay: &ContentComponentBuilderOverlay{
			Mode: f.b.overlayMode,
		},
	}

	var (
		overlay = f.b.overlayMode
	)

	cb.PreBody = append(cb.PreBody, f.Notice)

	if f.b.ctx.R.FormValue(ParamActionsDisabled) != "true" {
		cb.PrimaryAction = f.PrimaryAction
		cb.BottomActions = f.Actions
	}

	if overlay.Overlayed() {
		if !f.ScopeDisabled {
			cb.Scope = web.Scope().Form().Locals().Vars().Closes()
		}
		return cb.BuildOverlay()
	}
	if !f.ScopeDisabled {
		scope := GetScope(f.b.ctx)
		if scope == nil {
			scope = web.Scope()
			cb.Scope = scope
		}
		scope.Form().Locals().Vars()
	}
	return cb.BuildPage()
}

func (f *Form) Respond(r *web.EventResponse) {
	comp := f.Component()

	switch f.b.overlayMode {
	case actions.Dialog:
		f.MB.p.Dialog().
			SetTargetPortal(f.Portal).
			SetContentPortalName(f.Portal+"Content").
			SetScrollable(true).
			Respond(r, comp)
	default:
		if f.Portal != "" && f.b.overlayMode.IsDrawer() {
			d := f.MB.p.Drawer(f.b.overlayMode).
				SetPortalName(f.Portal).
				SetValidWidth(f.MB.rightDrawerWidth)

			d.Respond(r, comp)
		} else if f.Portal != "" {
			r.UpdatePortal(f.Portal, comp)
		} else {
			r.Body = comp
		}
	}
}
