package gitedit

import (
	"fmt"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// gitPerm is a permission a repository asks: what it is for, its permission
// ("!git", "@get"), and whether the user of the request has it.
type gitPerm struct {
	label, perm string
	ok          bool
}

// hasGit says whether the request reaches the files by git: a URL of them
// told (GitURL), and !git.
func (b *Builder) hasGit(ctx *web.EventContext) bool {
	return b.GitURL != nil && b.allowedAction(ctx.R, ActionGit)
}

// gitView is the tab Git: the repositories by git, each its own section —
// the draft (the user's, or the other's the page is of), and the site —:
// its URL, what a push does, how to clone it, the permissions it asks, each
// had or not, and the hook that signs the commits.
func (b *Builder) gitView(ctx *web.EventContext, m *Messages) h.HTMLComponent {
	r := ctx.R
	owner, own := b.ownerKey(r)
	// the resource of the page, out of the groups: its permissions in full
	base := ""
	if ver := b.page.Page().ActionVerifier(r, presets.PermGet); ver != nil {
		if base = ver.PreferredResource(); base == "" {
			base = ver.Resource()
		}
	}
	getOK := true
	if ver := b.page.Page().ActionVerifier(r, presets.PermGet); ver != nil {
		getOK = ver.Allowed()
	}
	perm := func(label, action string) gitPerm {
		return gitPerm{label, presets.ActionPerm(action), b.allowedAction(r, action)}
	}
	files := []gitPerm{
		perm(m.CreateAction, ActionCreate), perm(m.EditAction, ActionEdit), perm(m.RenameAction, ActionRename),
		perm(m.MoveAction, ActionMove), perm(m.DeleteAction, ActionDelete),
	}
	draftPerms := append([]gitPerm{{m.GitPermGet, presets.PermGet, getOK}, perm(m.GitPermPush, ActionGit)}, files...)
	sitePerms := append(append([]gitPerm{}, draftPerms...), perm(m.PublishAction, ActionPublish))

	title := m.GitDraftTitle
	if !own {
		title = fmt.Sprintf(m.GitOtherDraftTitle, b.user(r, owner).Label())
	}
	return h.Div(
		b.gitSection(ctx, m, "draft", title, m.GitDraftWhat, b.GitURL(r, owner), base, draftPerms),
		b.gitSection(ctx, m, "site", m.GitSiteTitle, m.GitSiteWhat, b.GitURL(r, ""), base, sitePerms),
	)
}

// gitSection is the section of a repository by git.
func (b *Builder) gitSection(ctx *web.EventContext, m *Messages, kind, title, what, url, base string, perms []gitPerm) h.HTMLComponent {
	r := ctx.R
	copyField := func(label, value, attr string) h.HTMLComponent {
		return v.VTextField().Label(label).ModelValue(value).Readonly(true).
			Variant(v.FieldVariantOutlined).Density(v.DensityCompact).HideDetails(true).Class("mb-3").
			Attr("append-inner-icon", "mdi-content-copy").Attr(attr, value).
			Attr("@click:append-inner", fmt.Sprintf("navigator.clipboard.writeText(%q)", value))
	}
	var rows []h.HTMLComponent
	refused := false
	for _, p := range perms {
		icon, color := "mdi-check", "success"
		if !p.ok {
			icon, color, refused = "mdi-close", "error", true
		}
		rows = append(rows, h.Tr(
			h.Td(v.VIcon(icon).Color(color).Size(v.SizeSmall)).Style("width: 32px"),
			h.Td(h.Text(p.label)),
			h.Td(h.Code(base+p.perm)),
		).Attr("data-git-perm", p.perm).Attr("data-git-perm-ok", strconv.FormatBool(p.ok)))
	}
	var missing h.HTMLComponent
	if refused {
		missing = h.P(h.Text(m.GitPermsMissing)).Class("text-caption text-error mb-2")
	}
	hook := fmt.Sprintf("curl -fsSL -u %s %s -o .git/hooks/commit-msg && chmod +x .git/hooks/commit-msg",
		shellQuote(b.actor(r).Login), absURL(r, b.hookPath()))
	var key h.HTMLComponent
	if b.HelpURL != "" {
		key = h.P(h.Text(m.GitKeyHint+" "), h.A(h.Text(m.GitKeyHelp)).Href(strings.TrimSuffix(b.HelpURL, "/")+"/08-git").
			Attr("target", "_blank")).Class("text-caption mt-2")
	} else {
		key = h.P(h.Text(m.GitKeyHint)).Class("text-caption mt-2")
	}
	return v.VCard(
		v.VCardTitle(h.Text(title)),
		v.VCardSubtitle(h.Text(what)),
		v.VCardText(
			copyField(m.GitURLLabel, url, "data-git-url"),
			copyField(m.GitCloneLabel, "git clone "+url, "data-git-clone"),
			h.H4(m.GitPermsTitle).Class("mt-2 mb-1"),
			missing,
			v.VTable(h.Tbody(rows...)).Density(v.DensityCompact).Class("mb-3"),
			h.P(h.Text(m.GitHookHint)).Class("text-caption mb-2"),
			copyField(m.GitHook, hook, "data-git-hook"),
			key,
		),
	).Variant(v.VariantOutlined).Class("mb-4").Attr("data-git-repo", kind)
}
