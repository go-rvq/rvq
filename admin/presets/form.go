package presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
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
	pre, post   []ModeObjectComponentFunc
}

func NewFormBuilder(ctx *web.EventContext, mb *ModelBuilder, fb *FieldsBuilder, obj interface{}) *FormBuilder {
	rawID := mb.MustRecordID(obj)

	f := &FormBuilder{
		id:          vx.ObjectID(obj),
		fb:          fb,
		ctx:         ctx,
		msgr:        MustGetMessages(ctx.Context()),
		mb:          mb,
		obj:         obj,
		overlayMode: GetOverlay(ctx),
	}
	if rawID.IsZero() {
		f.mode = NEW
	} else {
		f.mode = EDIT
	}
	return f
}

func (f *FormBuilder) Mode() FieldMode {
	return f.mode
}

func (f *FormBuilder) SetPre(v []ModeObjectComponentFunc) *FormBuilder {
	f.pre = v
	return f
}

func (f *FormBuilder) SetPost(v []ModeObjectComponentFunc) *FormBuilder {
	f.post = v
	return f
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
			).Slot("{ locals }").LocalsInit(`{ show: true }`)
		}
	}

	form.Body = f.fb.ToComponentFull(f.pre, f.post, &ToComponentOptions{}, f.mb.Info(), f.obj, FieldModeStack{f.mode}, ctx)

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
	Wrap          func(h h.HTMLComponent) h.HTMLComponent
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

	overlay := f.b.overlayMode

	cb.PreBody = append(cb.PreBody, f.Notice)

	if f.b.ctx.R.FormValue(ParamActionsDisabled) != "true" {
		cb.PrimaryAction = f.PrimaryAction
		cb.BottomActions = f.Actions
	}

	if overlay.Overlayed() {
		if !f.ScopeDisabled {
			cb.Scope = web.Scope().Form().Locals().Vars()
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
	f.RespondToPortal(f.Portal, r)
}

func (f *Form) RespondToPortal(portal string, r *web.EventResponse) {
	comp := f.Component()
	oldWrap := f.Wrap

	f.Wrap = func(h h.HTMLComponent) h.HTMLComponent {
		if oldWrap != nil {
			h = oldWrap(h)
		}
		return web.Scope(h).FormInit("{}")
	}

	switch f.b.overlayMode {
	case actions.Dialog:
		f.MB.p.Dialog().
			SetTargetPortal(portal).
			SetContentPortalName(f.Portal+"Content").
			SetScrollable(true).
			RootWrap(f.Wrap).
			Respond(f.b.ctx, r, comp)
	default:
		if f.Portal != "" && f.b.overlayMode.IsDrawer() {
			d := f.MB.p.Drawer(f.b.overlayMode).
				SetPortalName(portal).
				SetScrollable(true).
				SetValidWidth(f.MB.rightDrawerWidth).
				RootWrap(f.Wrap)

			d.Respond(r, comp)
		} else if portal != "" {
			r.UpdatePortal(portal, comp)
		} else {
			r.Body = comp
		}
	}
}
