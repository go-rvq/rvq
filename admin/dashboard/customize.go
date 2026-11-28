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
func (b *Panel) customizeTools(ctx *web.EventContext, msgs *Messages, prefs *Prefs, layout []Section) h.HTMLComponent {
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
func (b *Panel) sectionDialog(msgs *Messages) h.HTMLComponent {
	save := b.event(EventPref).Query("op", "section").
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
func (b *Panel) addDialog(ctx *web.EventContext, msgs *Messages, prefs *Prefs, layout []Section) h.HTMLComponent {
	var groups h.HTMLComponents
	anyRemoved := false
	for _, g := range b.sortedGroups() {
		var boxes h.HTMLComponents
		var shows []string
		original, _ := b.sectionHead(ctx, Section{ID: g.ID})
		original = cmp.Or(original, msgs.Untitled)
		for _, w := range b.all() {
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
	add := b.event(EventPref).Query("op", web.Var(`locals.addTarget ? "add-to" : "add"`)).
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
.dashboard [data-dashboard-col].dashboard-dragging, .dashboard [data-dashboard-group].dashboard-dragging { opacity: .4; }
.dashboard .dashboard-section-handle { cursor: grab; }
.dashboard .dashboard-col { position: relative; }
.dashboard .dashboard-resize { position: absolute; top: 8px; bottom: 8px; right: 0; width: 8px; cursor: ew-resize; border-radius: 4px; touch-action: none; }
.dashboard .dashboard-resize:hover, .dashboard .dashboard-resizing .dashboard-resize { background: rgba(var(--v-theme-primary), .35); }
.dashboard .dashboard-disabled { border-color: rgb(var(--v-theme-warning)) !important; border-style: dashed !important;
  background: repeating-linear-gradient(-45deg, rgba(var(--v-theme-warning), .07) 0 10px, transparent 10px 20px) !important; }
`

// dashScript is the dragging of the customize mode (window.__rvqDash, the
// events of the template call it). A widget dragged over a section — of any
// group — goes before the widget under the pointer (or last). A section
// dragged by its handle goes before or after the section under the pointer
// among its siblings only: never into another (its parent kept). Dropped,
// the layout is kept — the sections and their widgets in the page's order —
// in the panel of the page (data-dashboard).
func dashScript() string {
	return `({window}) => {
	const w = window, d = w.document;
	let dragged = null, sec = null;
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
	// sibling is the section under the pointer that is a sibling of the one
	// dragged (the one it is in, up to that level), null none
	const sibling = (e) => {
		let t = e.target.closest && e.target.closest("[data-dashboard-group]");
		while (t && t.parentElement !== sec.parentElement) t = t.parentElement.closest("[data-dashboard-group]");
		return t === sec ? null : t;
	};
	const clear = () => d.querySelectorAll(".dashboard-drop, .dashboard-dragging").forEach((x) => x.classList.remove("dashboard-drop", "dashboard-dragging"));
	const dash = w.__rvqDash = {
		start(e) {
			const col = e.target.closest && e.target.closest("[data-dashboard-col]");
			if (!col || sec) return;
			dragged = col;
			e.dataTransfer.effectAllowed = "move";
			e.dataTransfer.setData("text/plain", col.getAttribute("data-dashboard-col"));
			col.classList.add("dashboard-dragging");
		},
		end() {
			clear();
			dragged = sec = null;
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
			const el = dragged;
			dash.end();
			dash.save(el);
		},
		// a section, by its handle
		secStart(e) {
			e.stopPropagation();
			const box = e.target.closest("[data-dashboard-group]");
			if (!box) return;
			sec = box;
			e.dataTransfer.effectAllowed = "move";
			e.dataTransfer.setData("text/plain", box.getAttribute("data-dashboard-group"));
			e.dataTransfer.setDragImage(box, 16, 16);
			box.classList.add("dashboard-dragging");
		},
		secOver(e) {
			if (!sec || e.__rvqSec) return;
			e.__rvqSec = true;
			const t = sibling(e);
			if (!t) return;
			e.preventDefault();
			const r = t.getBoundingClientRect();
			const at = e.clientY < r.top + r.height / 2 ? t : t.nextSibling;
			if (at !== sec && at !== sec.nextSibling) t.parentElement.insertBefore(sec, at);
		},
		secDrop(e) {
			if (!sec || e.__rvqSec) return;
			e.__rvqSec = true;
			e.preventDefault();
			const el = sec;
			dash.end();
			dash.save(el);
		},
		// a widget's width, by its right edge: in columns of the grid of 12
		// of its row, MinSize the least; let go, kept
		resizeStart(e) {
			e.preventDefault();
			e.stopPropagation();
			const col = e.target.closest("[data-dashboard-col]");
			if (!col) return;
			const unit = col.parentElement.getBoundingClientRect().width / 12;
			const from = +col.getAttribute("data-dashboard-size"), x0 = e.clientX;
			let n = from;
			col.setAttribute("draggable", "false");
			col.classList.add("dashboard-resizing");
			const move = (ev) => {
				n = Math.max(` + strconv.Itoa(int(MinSize)) + `, Math.min(12, from + Math.round((ev.clientX - x0) / unit)));
				col.style.flex = "0 0 " + (n / 12 * 100) + "%";
				col.style.maxWidth = (n / 12 * 100) + "%";
			};
			const up = () => {
				w.removeEventListener("pointermove", move);
				w.removeEventListener("pointerup", up);
				col.classList.remove("dashboard-resizing");
				col.setAttribute("draggable", "true");
				if (n === from) return;
				col.setAttribute("data-dashboard-size", n);
				const root = col.closest("[data-dashboard]") || d;
				plaid().eventFunc(` + strconv.Quote(EventPref) + `).query("panel", root.getAttribute("data-dashboard") || "home").
					query("op", "size").query("arg", col.getAttribute("data-dashboard-col")).query("size", String(n)).go();
			};
			w.addEventListener("pointermove", move);
			w.addEventListener("pointerup", up);
		},
		layout(root) {
			return [...root.querySelectorAll("[data-dashboard-row]")].map((r) => ({
				id: r.getAttribute("data-dashboard-row"),
				widgets: [...r.querySelectorAll(":scope > [data-dashboard-col]")].map((c) => c.getAttribute("data-dashboard-col")),
			}));
		},
		save(el) {
			const root = el.closest("[data-dashboard]") || d;
			plaid().eventFunc(` + strconv.Quote(EventPref) + `).query("panel", root.getAttribute("data-dashboard") || "home").
				query("op", "layout").query("arg", JSON.stringify(dash.layout(root))).go();
		},
	};
}`
}
