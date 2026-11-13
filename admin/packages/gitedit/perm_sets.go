package gitedit

import (
	"context"
	"strings"

	"github.com/go-rvq/rvq/admin/presets"
)

// installPermSets gives the page its sets of permissions (presets.PermSet),
// of the actions it has: "admin:/site-files:%publishers".
//
//	readers      see the files, their history, the preview
//	authors      readers, and work on one's draft and save it
//	publishers   authors, and publish the draft
//	developers   readers, and git and WebDAV (a clone asks to read them all)
//	sharers      share the draft, by a public link too
//	reviewers    see every user's draft
//	maintainers  start the repository, restore the default files
func (b *Builder) installPermSets() {
	page := b.page.Page()
	act := func(names ...string) string {
		var out []string
		for _, n := range names {
			out = append(out, strings.TrimPrefix(presets.ActionPerm(n), "!"))
		}
		return "!{" + strings.Join(out, ", ") + "}"
	}
	words := func(f func(m *Messages) (string, string)) (func(ctx context.Context) string, func(ctx context.Context) string) {
		return func(ctx context.Context) string { t, _ := f(b.msgs(ctx)); return t },
			func(ctx context.Context) string { _, d := f(b.msgs(ctx)); return d }
	}
	readers := "@get"
	if b.PreviewPath != "" {
		readers = "{@get, " + presets.ActionPerm(ActionPreview) + "}"
	}
	add := func(name, body string, f func(m *Messages) (string, string)) {
		t, d := words(f)
		page.PermSet(name, body, t, d)
	}
	add("readers", readers, func(m *Messages) (string, string) { return m.SetReaders, m.SetReaders_Desc })
	add("authors", "{%readers, "+act(ActionCreate, ActionEdit, ActionRename, ActionMove, ActionDelete, ActionImport,
		ActionCommit, ActionUpdate, ActionDiscard, ActionReset)+"}",
		func(m *Messages) (string, string) { return m.SetAuthors, m.SetAuthors_Desc })
	add("publishers", "{%authors, "+presets.ActionPerm(ActionPublish)+"}",
		func(m *Messages) (string, string) { return m.SetPublishers, m.SetPublishers_Desc })
	add("developers", "{%readers, "+presets.ActionPerm(ActionGit)+"}",
		func(m *Messages) (string, string) { return m.SetDevelopers, m.SetDevelopers_Desc })
	sharers := presets.ActionPerm(ActionShare)
	if b.PublicPreviewPath != "" {
		sharers = act(ActionShare, ActionPublicPreview)
	}
	add("sharers", sharers, func(m *Messages) (string, string) { return m.SetSharers, m.SetSharers_Desc })
	add("reviewers", presets.ActionPerm(ActionDrafts), func(m *Messages) (string, string) { return m.SetReviewers, m.SetReviewers_Desc })
	if b.Defaults != nil {
		add("maintainers", act(ActionInitDefaults, ActionResetDefaults),
			func(m *Messages) (string, string) { return m.SetMaintainers, m.SetMaintainers_Desc })
	}
}
