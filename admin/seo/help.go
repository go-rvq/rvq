package seo

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// SEO help overlay.
//
// A "?" button sits in the SEO editing form's actions, next to Save. It opens an
// overlay — a dialog or a navigation drawer, chosen by the URL the same way the
// detailing Edit button chooses (see DialogBuilder.Respond) — that explains how
// SEO works.
//
// The overlay body is a reusable frame, common to every context (what the SEO
// fields are, how the {{variables}} tags work, and that an empty field inherits
// down the SEO chain), with a context-specific section mounted inside it: the
// Global/Page/Post explanation and how the merge reaches that context.
//
// The frame and its text live here (multi-language via the seo i18n Messages).
// The context+merge section is supplied by the application through
// RegisterHelpSection, because the chain of contexts (and any templating or
// fallbacks an app adds on top) is the application's — e.g. the site-side
// composition in hermon-cms. A context with no registered section still shows the
// frame.
const (
	// eventSEOHelp opens the help overlay.
	eventSEOHelp = "seo_help"
	// ParamSEOHelpContext carries the help context key HelpButton passes.
	ParamSEOHelpContext = "seo_help_context"

	// HelpContextGlobal is the context key for the Global SEO editor, which the
	// seo package wires itself. Applications define their own keys (e.g. "page",
	// "post") for the models they mount, passing them to AddEditingHelpButton and
	// RegisterHelpSection.
	HelpContextGlobal = "global"
)

// HelpSectionFunc renders the context-specific part of the help overlay — the
// Global/Page/Post explanation and how the merge reaches this context — mounted
// inside the common frame. It returns nil to contribute nothing.
type HelpSectionFunc func(ctx *web.EventContext) h.HTMLComponent

// RegisterHelpSection registers the help body for a context key. Returns the
// builder for chaining.
func (b *Builder) RegisterHelpSection(key string, f HelpSectionFunc) *Builder {
	if b.helpSections == nil {
		b.helpSections = map[string]HelpSectionFunc{}
	}
	b.helpSections[key] = f
	return b
}

// registerHelpEvent wires the help overlay event. Called from Install.
func (b *Builder) registerHelpEvent(pb *presets.Builder) {
	pb.GetWebBuilder().RegisterEventFunc(eventSEOHelp, b.seoHelpDialog(pb))
}

// AddEditingHelpButton adds the "?" help button to an editing form's top-right
// action area, next to Save, for the given help context key. Use it on any
// editing builder that edits a Setting — the Global SEO model, or an app-mounted
// per-record SEO (Page/Post).
func (b *Builder) AddEditingHelpButton(eb *presets.EditingBuilder, contextKey string) {
	eb.TopRightActionsFunc(func(_ interface{}, ctx *web.EventContext) h.HTMLComponent {
		return b.HelpButton(ctx, contextKey)
	})
}

// HelpButton is the "?" icon button (with a tooltip) that opens the help overlay
// for contextKey. Exposed so an application can place it itself when it does not
// use AddEditingHelpButton.
func (b *Builder) HelpButton(ctx *web.EventContext, contextKey string) h.HTMLComponent {
	msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nSeoKey, Messages_en_US).(*Messages)
	return VBtn("").
		Variant(VariantText).
		Density("comfortable").
		Icon(true).
		Attr("data-event", eventSEOHelp).
		Attr("@click", web.Plaid().
			EventFunc(eventSEOHelp).
			Query(ParamSEOHelpContext, contextKey).
			Go()).
		Children(
			VIcon("mdi-help-circle-outline"),
			VTooltip(h.Text(msgr.HelpTooltip)).Attr("activator", "parent"),
		)
}

// seoHelpDialog is the event handler: it assembles the common frame with the
// registered context section mounted inside and responds with an overlay
// (dialog/drawer per the URL).
func (b *Builder) seoHelpDialog(pb *presets.Builder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nSeoKey, Messages_en_US).(*Messages)

		var contextComp h.HTMLComponent
		if b.helpSections != nil {
			if f := b.helpSections[ctx.R.FormValue(ParamSEOHelpContext)]; f != nil {
				contextComp = f(ctx)
			}
		}

		intro := HelpBlock(msgr.HelpIntroTitle, msgr.HelpIntro)

		var body h.HTMLComponent
		if contextComp != nil {
			// The context section (Global/Page/Post + merge) is mounted inside the
			// frame that repeats for every context, and owns everything after the
			// intro. It renders Vuetify components (VAlert, VTable) and single-brace
			// `{ expr }` template syntax — which Vue does not interpolate — so it
			// must compile normally. v-pre here would leave those components
			// unstyled (a raw <v-alert> with a bare left border).
			body = h.Div(intro, contextComp).Class("pa-4")
		} else {
			// The generic frame is plain HTML that documents the built-in
			// {{variables}} literally: v-pre keeps Vue from interpolating `{{…}}`
			// (which even throws "illegal character U+2026"). No Vuetify component
			// lives here, so v-pre costs nothing.
			body = h.Div(
				intro,
				HelpBlock(msgr.HelpVariablesTitle, msgr.HelpVariables),
				HelpBlock(msgr.HelpInheritanceTitle, msgr.HelpInheritance),
			).Class("pa-4").Attr("v-pre", true)
		}

		pb.Dialog("760").Respond(ctx, &r,
			vx.VXDialog().
				Title(msgr.HelpTitle).
				SlotBody(body).
				Closable(true),
		)
		return
	}
}

// HelpBlock renders one titled help block; body is HTML, so a message may use
// simple markup (<code>, <ul>, <b>). Exposed so a context section can build its
// own blocks with the same look as the frame.
func HelpBlock(title, body string) h.HTMLComponent {
	if title == "" && body == "" {
		return h.Div()
	}
	return h.Div(
		h.H3(title).Class("text-subtitle-1 font-weight-bold mt-4 mb-1"),
		h.Div(h.RawHTML(body)).Class("text-body-2"),
	)
}
