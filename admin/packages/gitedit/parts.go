package gitedit

import (
	"context"
	"errors"
	"net/http"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// The parts of the editor that are pages under its page — each an item of
// the menu under it —: the files of the draft (the IDE), its sharings, and
// the drafts the user reaches. Of another's draft (…/u/{user}), its files
// and sharings are under its page too, out of the menu.
const (
	filesSub  = "/files"
	sharesSub = "/shares"
	draftsSub = "/drafts"
)

// partsPortal is the portal of the actions of a part (their dialogs).
const partsPortal = "gitEditAction"

// installParts makes the pages of the parts — their actions the page's —,
// to be built.
func (b *Builder) installParts(p *presets.Builder) (pages []*presets.PageBuilder) {
	title := func(f func(m *Messages) string) func(ctx context.Context) string {
		return func(ctx context.Context) string { return f(GetMessages(ctx)) }
	}
	part := func(pg *presets.PageBuilder, sub string, words func(m *Messages) string, f partFunc) *presets.PageBuilder {
		part := pg.SubPage(sub).Title(title(words)).Layout(b.partLayout(f))
		b.setupActions(p, part)
		pages = append(pages, part)
		return part
	}
	files := func(m *Messages) string { return m.Files }
	shares := func(m *Messages) string { return m.SharesTab }
	drafts := func(m *Messages) string { return m.DraftsTab }

	part(b.page, filesSub, files, b.filesPart).InMenu(true).MenuIcon("mdi-file-code-outline")
	part(b.userPage, filesSub, files, b.filesPart)
	if b.DB != nil {
		part(b.page, sharesSub, shares, b.sharesPart).InMenu(true).MenuIcon("mdi-account-multiple-outline")
		part(b.userPage, sharesSub, shares, b.sharesPart)
	}
	part(b.page, draftsSub, drafts, b.draftsPart).InMenu(true).MenuIcon("mdi-folder-account-outline").
		// the drafts shared with the user (DB), or every draft (!drafts)
		MenuVisible(func(r *http.Request) bool { return b.DB != nil || b.allowedAction(r, ActionDrafts) })
	return
}

// part is what a part shows: its body, and the buttons of the page's header
// (beside its print button).
type part struct {
	body    h.HTMLComponent
	actions h.HTMLComponents
}

// partFunc is a part, of the draft of owner (own: the user's).
type partFunc func(ctx *web.EventContext, m *Messages, owner string, own bool) (part, error)

// partLayout is the page of a part: the page's permission, the draft's
// (mayReach); whose draft it is, when another's.
func (b *Builder) partLayout(f partFunc) web.PageFunc {
	return func(ctx *web.EventContext) (r web.PageResponse, err error) {
		m := GetMessages(ctx.Context())
		if ver := b.page.Page().ActionVerifier(ctx.R, presets.PermGet); ver != nil && ver.Denied() {
			return r, perm.PermissionDenied
		}
		owner, own := b.ownerKey(ctx.R)
		if !b.mayReach(ctx.R, owner) {
			return r, perm.PermissionDenied
		}
		pt, err := f(ctx, m, owner, own)
		if errors.Is(err, ErrEmpty) {
			// no commit yet: its first files, by git or the default ones
			r.Body = b.emptyView(ctx, m)
			return r, nil
		}
		if err != nil {
			return
		}
		r.Body = h.Div(
			web.Portal().Name(partsPortal),
			b.whoseDraft(ctx, m, owner, own),
			pt.body,
		).Class("pa-4")
		r.Actions = pt.actions
		return
	}
}

// filesPart are the files of the draft: the IDE, in the page or over the
// whole of it (the window) — expanded by a button of the page's header, back
// by one of the IDE's header (vars.ideFull).
func (b *Builder) filesPart(ctx *web.EventContext, m *Messages, _ string, _ bool) (part, error) {
	if _, _, err := b.Draft(ctx.R); err != nil {
		return part{}, err
	}
	expand := v.VBtn("").Icon("mdi-arrow-expand").Density(v.DensityComfortable).
		Attr("title", m.IdeExpand).Attr("@click", "vars.ideFull = true").Attr("data-ide-full", true)
	restore := "[{id: 'restore', icon: 'mdi-arrow-collapse', title: " + h.JSONString(m.IdeCollapse) +
		", onClick: () => { vars.ideFull = false }}]"
	return part{
		body: h.Div(
			b.ideTag(ctx.R).Attr(":height", "vars.ideFull ? '100vh' : 'calc(100vh - 160px)'").
				Attr(":header-actions", "vars.ideFull ? "+restore+" : []"),
		).Attr(web.VAssign("vars", `{ideFull: false}`)...).
			Attr(":style", "vars.ideFull ? "+
				"'position: fixed; inset: 0; z-index: 2000; background: rgb(var(--v-theme-surface))' : ''"),
		actions: h.HTMLComponents{expand},
	}, nil
}

// sharesPart are the sharings of the draft.
func (b *Builder) sharesPart(ctx *web.EventContext, m *Messages, owner string, own bool) (part, error) {
	return part{body: b.sharesView(ctx, m, partsPortal, owner, own)}, nil
}

// draftsPart are the drafts the user reaches.
func (b *Builder) draftsPart(ctx *web.EventContext, m *Messages, _ string, _ bool) (part, error) {
	return part{body: b.draftsView(ctx, m, b.allowedAction(ctx.R, ActionDrafts))}, nil
}

// partsNav are the links to the pages of another's draft — the draft, its
// files, its sharings —: out of the menu, which has the user's own.
func (b *Builder) partsNav(ctx *web.EventContext, m *Messages, owner string) h.HTMLComponent {
	base := b.UserURL(owner)
	link := func(text, href, icon string) h.HTMLComponent {
		return v.VBtn(text).PrependIcon(icon).Variant(v.VariantText).Size(v.SizeSmall).
			Attr("@click", web.Plaid().PushStateURL(href).Go()).Attr("data-part-link", href)
	}
	items := []h.HTMLComponent{
		link(m.DraftOther, base, "mdi-source-branch"),
		link(m.Files, base+filesSub, "mdi-file-code-outline"),
	}
	if b.DB != nil {
		items = append(items, link(m.SharesTab, base+sharesSub, "mdi-account-multiple-outline"))
	}
	return h.Div(items...).Class("d-flex flex-wrap mt-1")
}
