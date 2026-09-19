package seo

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

// VariableItem is one entry in the SEO "+ Variable" menu: a label shown to the
// user and the expression inserted at the cursor (without braces — the field's
// Template wraps it, e.g. into `{expr}`).
type VariableItem struct {
	Label string
	Expr  string
}

// VariableGroup groups variable items under a template global. Name is the
// global's identifier (Site, Page, Post, …), used as the group's stable value;
// Label is what the user sees, and Description a short subtitle explaining the
// global (like the help overlay). When Label is empty it falls back to Name.
type VariableGroup struct {
	Name        string
	Label       string
	Description string
	Items       []VariableItem
}

// RegisterVariableGroups sets the provider of the "+ Variable" menu groups. The
// app supplies them (with translated field labels) because the template globals
// are the app's; without it the editor shows no variable menu.
func (b *Builder) RegisterVariableGroups(f func(ctx *web.EventContext) []VariableGroup) *Builder {
	b.variableGroups = f
	return b
}

// variablesMenu renders the "+ Variable" multi-level menu from the registered
// groups. Each leaf inserts its expression through the send-variables component
// (its Template wraps it into the `{expr}` form). Returns nil when no groups are
// registered.
func (b *Builder) variablesMenu(ctx *web.EventContext) h.HTMLComponent {
	var groups []VariableGroup
	if b.variableGroups != nil {
		groups = b.variableGroups(ctx)
	}
	return b.renderVariablesMenu(ctx, groups)
}

// variablesMenuWithVars is variablesMenu plus a group for the setting's own
// custom variables, so a field can insert `{Name}` for a declared variable (and
// a descendant sees the inherited ones once merged — here it shows the setting's
// own).
func (b *Builder) variablesMenuWithVars(ctx *web.EventContext, setting *Setting) h.HTMLComponent {
	var groups []VariableGroup
	if b.variableGroups != nil {
		groups = b.variableGroups(ctx)
	}
	if g := customVarsMenuGroup(setting); g != nil {
		msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nSeoKey, Messages_en_US).(*Messages)
		if msgr.CustomVarsGroup != "" {
			g.Label = msgr.CustomVarsGroup
		}
		groups = append(groups, *g)
	}
	return b.renderVariablesMenu(ctx, groups)
}

func (b *Builder) renderVariablesMenu(ctx *web.EventContext, groups []VariableGroup) h.HTMLComponent {
	if len(groups) == 0 {
		return nil
	}
	msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nSeoKey, Messages_en_US).(*Messages)

	var groupComps []h.HTMLComponent
	for _, g := range groups {
		label := g.Label
		if label == "" {
			label = g.Name
		}
		var items []h.HTMLComponent
		for _, it := range g.Items {
			items = append(items,
				VListItem().Title(it.Label).
					Attr("@click", fmt.Sprintf("$refs.seo.addTags(%q)", it.Expr)),
			)
		}
		activator := VListItem().Attr("v-bind", "props").Title(label)
		if g.Description != "" {
			activator.Subtitle(g.Description)
		}
		groupComps = append(groupComps,
			VListGroup(
				web.Slot(activator).Name("activator").Scope("{ props }"),
				h.Components(items...),
			).Value(g.Name),
		)
	}

	return VMenu(
		web.Slot(
			VBtn(msgr.AddVariable).PrependIcon("mdi-plus").
				Variant(VariantOutlined).Size(SizeSmall).Attr("v-bind", "props"),
		).Name("activator").Scope("{ props }"),
		VList(groupComps...).Density(DensityCompact),
	)
}
