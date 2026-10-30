package dashboard

import (
	"cmp"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// The customize mode of the dashboard: its tools — a new section, the
// widgets removed added back —, their dialogs (the locals of the page's
// scope: Page), and the dragging of the widgets (dashScript).

// customizeLocals are the locals of the dialogs of the customize mode.
const customizeLocals = `{secOpen: false, secId: "", secCustom: true, secTitle: "", secIcon: "", secParent: "", secParents: [], addOpen: false, add: [], addTarget: "", addTitle: ""}`

// customizeTools are the tools of the customize mode, its dialogs and its
// script.
func (b *Builder) customizeTools(ctx *web.EventContext, msgs *Messages, prefs *Prefs, layout []Section) h.HTMLComponent {
	return h.Div(
		h.Tag("component").Attr(":is", "'style'").Children(h.RawHTML(dashStyle)),
		h.Tag("user-component").Attr(":setup", "[("+dashScript()+")]"),
		h.Div(
			v.VBtn(msgs.NewSection).PrependIcon("mdi-plus").Variant(v.VariantTonal).Size(v.SizeSmall).
				Attr("@click", `locals.secId = ""; locals.secCustom = true; locals.secTitle = ""; locals.secIcon = ""; locals.secParent = ""; locals.secParents = `+
					h.JSONString(b.parentChoices(ctx, msgs, layout, ""))+`; locals.secOpen = true`).
				Attr("data-dashboard-new-section", true),
			v.VBtn(msgs.AddWidgets).PrependIcon("mdi-view-grid-plus-outline").Variant(v.VariantTonal).Size(v.SizeSmall).
				Attr("@click", `locals.add = []; locals.addTarget = ""; locals.addOpen = true`).Attr("data-dashboard-add", true),
			h.Div(h.Text(msgs.DragHint)).Class("text-caption text-medium-emphasis"),
		).Class("d-flex flex-wrap align-center ga-2 mt-4"),
		b.sectionDialog(msgs),
		b.addDialog(ctx, msgs, prefs, layout),
	)
}

// sectionDialog makes or changes a section of the user's own: its title and
// icon.
func (b *Builder) sectionDialog(msgs *Messages) h.HTMLComponent {
	save := web.Plaid().EventFunc(EventPref).Query("op", "section").
		Query("arg", web.Var("locals.secId")).Query("title", web.Var("locals.secTitle")).Query("icon", web.Var("locals.secIcon")).
		Query("parent", web.Var("locals.secParent")).Go()
	return v.VDialog(
		v.VCard(
			v.VCardTitle(
				h.Span(msgs.EditSection).Attr("v-if", "locals.secId"),
				h.Span(msgs.NewSection).Attr("v-else", true),
			),
			v.VCardText(
				// a group's section: its title and icon the group's
				h.Div(
					v.VTextField().Label(msgs.SectionTitle).Attr("v-model", "locals.secTitle").Autofocus(true).
						Variant(v.VariantOutlined).Density(v.DensityCompact).Attr("data-dashboard-section-title", true),
					v.VTextField().Label(msgs.SectionIcon).Attr("v-model", "locals.secIcon").Placeholder("mdi-…").
						Attr(":prepend-inner-icon", "/^mdi-[a-z0-9-]+$/.test(locals.secIcon || '') ? locals.secIcon : 'mdi-shape-outline'").
						Variant(v.VariantOutlined).Density(v.DensityCompact).Class("mt-2").Attr("data-dashboard-section-icon", true),
				).Attr("v-if", "locals.secCustom"),
				// the section it is in: a tree
				v.VSelect().Label(msgs.ParentSection).Attr(":items", "locals.secParents").ItemTitle("title").ItemValue("value").
					Attr("v-model", "locals.secParent").Variant(v.VariantOutlined).Density(v.DensityCompact).Class("mt-2").
					PrependInnerIcon("mdi-file-tree-outline").Attr("data-dashboard-section-parent", true),
			),
			v.VCardActions(
				v.VSpacer(),
				v.VBtn(msgs.Cancel).Variant(v.VariantText).Attr("@click", "locals.secOpen = false"),
				v.VBtn(msgs.Save).Color("primary").Variant(v.VariantFlat).Attr(":disabled", "locals.secCustom && !(locals.secTitle || '').trim()").
					Attr("@click", "locals.secOpen = false; "+save).Attr("data-dashboard-section-save", true),
			),
		),
	).Attr("v-model", "locals.secOpen").MaxWidth("480").Attr("data-dashboard-section-dialog", true)
}

// addDialog adds widgets, in their original sections: with no target (Add
// widgets) the ones removed, each back to its original section; with a
// target (a section's +: locals.addTarget) also the others not in it — a
// disabled one marked, an enabled one saying the section it is in and its
// original —, all moved into the target, enabled.
func (b *Builder) addDialog(ctx *web.EventContext, msgs *Messages, prefs *Prefs, layout []Section) h.HTMLComponent {
	var groups h.HTMLComponents
	anyRemoved := false
	for _, g := range b.sortedGroups() {
		var boxes h.HTMLComponents
		var shows []string
		original, _ := b.sectionHead(ctx, Section{ID: g.ID})
		original = cmp.Or(original, msgs.Untitled)
		for _, w := range b.widgets {
			if w.Group != g.ID || !w.visible(ctx.R) {
				continue
			}
			cur := sectionOf(layout, w.ID)
			removed := cur == nil
			anyRemoved = anyRemoved || removed
			// shown: a removed one always; another, with a target not its
			// section
			show := "true"
			if !removed {
				show = "!!locals.addTarget && locals.addTarget !== " + h.JSONString(cur.ID)
			}
			shows = append(shows, "("+show+")")
			var sub h.HTMLComponent
			switch {
			case removed:
				sub = v.VChip(h.Text(msgs.Removed)).Size(v.SizeXSmall).Variant(v.VariantTonal)
			case prefs.disabled(w.ID):
				sub = v.VChip(h.Text(msgs.Disabled)).Size(v.SizeXSmall).Color("warning").Variant(v.VariantFlat).
					PrependIcon("mdi-toggle-switch-off-outline")
			default:
				sub = h.Span(strings.NewReplacer("{section}", b.sectionPath(ctx, msgs, layout, cur.ID), "{original}", original).
					Replace(msgs.InSection)).Class("text-caption text-medium-emphasis")
			}
			boxes = append(boxes, h.Div(
				v.VCheckbox().Label(w.Title(ctx.Context())).Attr(":value", h.JSONString(w.ID)).Attr("v-model", "locals.add").
					Density(v.DensityCompact).HideDetails(true).Attr("data-dashboard-removed", w.ID),
				h.Div(sub).Class("ms-10 mt-n2 mb-1"),
			).Attr("v-if", show).Attr("data-dashboard-add-item", w.ID).
				AttrIf("data-dashboard-add-disabled", w.ID, !removed && prefs.disabled(w.ID)))
		}
		if len(boxes) == 0 {
			continue
		}
		_, icon := b.sectionHead(ctx, Section{ID: g.ID})
		groups = append(groups, h.Div(
			h.Div(
				h.If(icon != "", v.VIcon(icon).Size(v.SizeSmall).Class("me-2 text-medium-emphasis")),
				h.Text(original),
			).Class("d-flex align-center text-subtitle-2 font-weight-bold mt-3"),
			h.Div(boxes...),
		).Attr("v-if", strings.Join(shows, " || ")).Attr("data-dashboard-removed-group", g.ID))
	}
	add := web.Plaid().EventFunc(EventPref).Query("op", web.Var(`locals.addTarget ? "add-to" : "add"`)).
		Query("section", web.Var("locals.addTarget")).Query("arg", web.Var("locals.add.join(',')")).Go()
	return v.VDialog(
		v.VCard(
			v.VCardTitle(
				h.Span(msgs.AddWidgets).Attr("v-if", "!locals.addTarget"),
				h.Span(msgs.AddInto+" ").Attr("v-if", "locals.addTarget"),
				h.Span("{{ locals.addTitle }}").Class("font-weight-bold").Attr("v-if", "locals.addTarget"),
			),
			v.VCardText(
				h.Div(groups...),
				h.If(!anyRemoved, h.Div(h.Text(msgs.NoRemoved)).Class("text-medium-emphasis py-4").
					Attr("v-if", "!locals.addTarget").Attr("data-dashboard-none", true)),
			).Style("max-height: 70vh; overflow-y: auto"),
			v.VCardActions(
				v.VSpacer(),
				v.VBtn(msgs.Cancel).Variant(v.VariantText).Attr("@click", "locals.addOpen = false"),
				v.VBtn(msgs.AddSelected).Color("primary").Variant(v.VariantFlat).Attr(":disabled", "!locals.add.length").
					Attr("@click", "locals.addOpen = false; "+add).Attr("data-dashboard-add-selected", true),
			),
		),
	).Attr("v-model", "locals.addOpen").MaxWidth("600").Attr("data-dashboard-add-dialog", true)
}

// dashStyle is the look of the customize mode: a place to drop in, the
// widget dragged, a widget disabled (striped, in the warning color).
const dashStyle = `
.dashboard .dashboard-row { min-height: 96px; border: 1px dashed rgba(var(--v-border-color), var(--v-border-opacity)); border-radius: 8px; padding: 4px; }
.dashboard .dashboard-row.dashboard-drop { border-color: rgb(var(--v-theme-primary)); background: rgba(var(--v-theme-primary), .05); }
.dashboard [data-dashboard-col] { cursor: grab; }
.dashboard [data-dashboard-col].dashboard-dragging { opacity: .4; }
.dashboard .dashboard-disabled { border-color: rgb(var(--v-theme-warning)) !important; border-style: dashed !important;
  background: repeating-linear-gradient(-45deg, rgba(var(--v-theme-warning), .07) 0 10px, transparent 10px 20px) !important; }
`

// dashScript is the dragging of the widgets (window.__rvqDash, the events
// of the template call it): a widget dragged over a section goes before the
// widget under the pointer (or last), and dropped the layout is kept — the
// sections and their widgets as the page has them.
func dashScript() string {
	return `({window}) => {
	const w = window, d = w.document;
	let dragged = null;
	const place = (row, e) => {
		for (const c of row.querySelectorAll(":scope > [data-dashboard-col]")) {
			if (c === dragged) continue;
			const r = c.getBoundingClientRect();
			if (e.clientY < r.top || (e.clientY <= r.bottom && e.clientX < r.left + r.width / 2)) return c;
		}
		return null;
	};
	// put places the widget dragged in row, before the widget under the
	// pointer (or last)
	const put = (row, e) => {
		const before = place(row, e);
		if (before) { if (before !== dragged.nextSibling) row.insertBefore(dragged, before); }
		else if (row.lastElementChild !== dragged) row.appendChild(dragged);
	};
	const dash = w.__rvqDash = {
		start(e) {
			const col = e.target.closest && e.target.closest("[data-dashboard-col]");
			if (!col) return;
			dragged = col;
			e.dataTransfer.effectAllowed = "move";
			e.dataTransfer.setData("text/plain", col.getAttribute("data-dashboard-col"));
			col.classList.add("dashboard-dragging");
		},
		end() {
			if (dragged) dragged.classList.remove("dashboard-dragging");
			d.querySelectorAll(".dashboard-drop").forEach((x) => x.classList.remove("dashboard-drop"));
		},
		over(e) {
			if (!dragged) return;
			e.preventDefault();
			e.currentTarget.classList.add("dashboard-drop");
			put(e.currentTarget, e);
		},
		leave(e) {
			if (!e.currentTarget.contains(e.relatedTarget)) e.currentTarget.classList.remove("dashboard-drop");
		},
		drop(e) {
			if (!dragged) return;
			e.preventDefault();
			// where it was dropped (a quick drag may have been over another
			// section last)
			put(e.currentTarget, e);
			dash.end();
			dragged = null;
			dash.save();
		},
		layout() {
			return [...d.querySelectorAll("[data-dashboard-row]")].map((r) => ({
				id: r.getAttribute("data-dashboard-row"),
				widgets: [...r.querySelectorAll(":scope > [data-dashboard-col]")].map((c) => c.getAttribute("data-dashboard-col")),
			}));
		},
		save() {
			plaid().eventFunc(` + strconv.Quote(EventPref) + `).query("op", "layout").query("arg", JSON.stringify(dash.layout())).go();
		},
	};
}`
}
