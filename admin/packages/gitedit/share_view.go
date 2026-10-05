package gitedit

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

const dateFormat = "2006-01-02 15:04"

// whoseDraft says, on the page of another's draft, whose it is and how the
// user reaches it: shared with them (until when), or every draft (!drafts).
func (b *Builder) whoseDraft(ctx *web.EventContext, m *Messages, owner string, own bool) h.HTMLComponent {
	if own {
		return nil
	}
	r := ctx.R
	actor, _ := b.Identity(r)
	how := m.ViaDrafts
	if shares, _ := b.SharedWith(r.Context(), actor); len(shares) > 0 {
		for _, s := range shares {
			if s.Owner == owner {
				how = m.ViaShare
				if s.ExpiresAt != nil {
					how += " · " + fmt.Sprintf(m.ShareUntil, s.ExpiresAt.Local().Format(dateFormat))
				}
				break
			}
		}
	}
	return v.VAlert(
		h.Strong(fmt.Sprintf(m.DraftOf, b.user(r, owner).Label())),
		h.Text(" — "+how),
	).Type("info").Variant(v.VariantTonal).Density(v.DensityCompact).Class("mb-4").
		Attr("data-draft-owner", owner)
}

// shareState is how a sharing stands at now: active (until when), expired,
// revoked (when, by whom).
func (b *Builder) shareState(ctx *web.EventContext, m *Messages, s *DraftShare, now time.Time) (string, bool) {
	switch {
	case s.RevokedAt != nil:
		return fmt.Sprintf(m.ShareRevokedOn, s.RevokedAt.Local().Format(dateFormat),
			b.user(ctx.R, s.RevokedBy).Label()), false
	case s.ExpiresAt != nil && !s.ExpiresAt.After(now):
		return fmt.Sprintf(m.ShareExpiredOn, s.ExpiresAt.Local().Format(dateFormat)), false
	case s.ExpiresAt != nil:
		return fmt.Sprintf(m.ShareUntil, s.ExpiresAt.Local().Format(dateFormat)), true
	}
	return m.ShareForever, true
}

// sharesView is the tab of the sharings of the draft of owner: whom it is
// shared with — active, and the past ones —; sharing it (its owner, !share)
// and revoking (its owner, or !drafts).
func (b *Builder) sharesView(ctx *web.EventContext, m *Messages, portal, owner string, own bool) h.HTMLComponent {
	r := ctx.R
	shares, _ := b.SharesOf(r.Context(), owner)
	canShare := own && b.allowedAction(r, ActionShare)
	canRevoke := canShare || b.allowedAction(r, ActionDrafts)
	now := time.Now()

	var rows, past []h.HTMLComponent
	for _, s := range shares {
		state, active := b.shareState(ctx, m, s, now)
		var revoke h.HTMLComponent
		if active && canRevoke {
			revoke = v.VBtn(m.RevokeAction).PrependIcon("mdi-account-cancel-outline").Variant(v.VariantText).
				Size(v.SizeSmall).Color("error").Attr("data-revoke", s.ID.String()).
				Attr("@click", b.actionOnClick(ctx, ActionRevoke, portal, "id", s.ID.String()))
		}
		row := h.Tr(
			h.Td(h.Text(b.user(r, s.User).Label())),
			h.Td(h.Text(s.CreatedAt.Local().Format(dateFormat)+" · "+b.user(r, s.CreatedBy).Label())),
			h.Td(h.Text(state)),
			h.Td(revoke).Class("text-right"),
		).Attr("data-share", s.ID.String())
		if active {
			rows = append(rows, row)
		} else {
			past = append(past, row)
		}
	}
	table := func(rows []h.HTMLComponent) h.HTMLComponent {
		return v.VTable(
			h.Thead(h.Tr(h.Th(m.ShareUser), h.Th(m.ShareSince), h.Th(m.ShareState), h.Th(""))),
			h.Tbody(rows...),
		).Density(v.DensityCompact)
	}

	var share h.HTMLComponent
	if canShare && b.Users != nil {
		share = v.VBtn(m.ShareAction).PrependIcon("mdi-account-multiple-plus-outline").Variant(v.VariantTonal).
			Size(v.SizeSmall).Color("primary").Class("mb-3").Attr("data-share-action", true).
			Attr("@click", b.actionOnClick(ctx, ActionShare, portal))
	}
	var active h.HTMLComponent = h.P(h.Text(m.SharesNone)).Class("text-medium-emphasis")
	if len(rows) > 0 {
		active = table(rows)
	}
	var history h.HTMLComponent
	if len(past) > 0 {
		history = h.Div(h.H4(m.SharesHistory).Class("mt-6 mb-2"), table(past))
	}
	return v.VCard(v.VCardText(
		h.P(h.Text(m.SharesHint)).Class("text-body-2 mb-3"),
		share,
		active,
		history,
	)).Variant(v.VariantOutlined)
}

// draftsView is the tab of the drafts the user reaches: the ones shared with
// them; and, with !drafts, every draft — its owner, how it stands, its last
// commit, whom it is shared with —, each with its page, editor, preview and
// git.
func (b *Builder) draftsView(ctx *web.EventContext, m *Messages, all bool) h.HTMLComponent {
	r := ctx.R
	c := r.Context()
	actor, _ := b.Identity(r)
	link := func(text, href, icon string) *v.VBtnBuilder {
		return v.VBtn(text).PrependIcon(icon).Variant(v.VariantText).Size(v.SizeSmall).Href(href)
	}
	links := func(key string) h.HTMLComponent {
		page := b.UserURL(key)
		if key == actor {
			page = b.page.Page().FullPath()
		}
		items := []h.HTMLComponent{
			link(m.OpenDraft, page, "mdi-source-branch"),
			link(m.OpenIDE, page+"/editor", "mdi-open-in-new").Attr("target", "_blank"),
		}
		if b.previewPrefix != "" {
			p := b.previewPrefix + "/u/" + key
			if key == actor {
				p = b.PreviewURL
			}
			items = append(items, link(m.Preview, p, "mdi-eye-outline").Attr("target", "_blank"))
		}
		return h.Div(items...).Class("d-flex flex-wrap")
	}

	shared, _ := b.SharedWith(c, actor)
	var mine []h.HTMLComponent
	for _, s := range shared {
		until := m.ShareForever
		if s.ExpiresAt != nil {
			until = fmt.Sprintf(m.ShareUntil, s.ExpiresAt.Local().Format(dateFormat))
		}
		mine = append(mine, h.Tr(
			h.Td(h.Text(b.user(ctx.R, s.Owner).Label())),
			h.Td(h.Text(until)),
			h.Td(links(s.Owner)),
		).Attr("data-shared-draft", s.Owner))
	}
	var sharedComp h.HTMLComponent = h.P(h.Text(m.SharedWithMeNone)).Class("text-medium-emphasis")
	if len(mine) > 0 {
		sharedComp = v.VTable(
			h.Thead(h.Tr(h.Th(m.DraftOwner), h.Th(m.ShareState), h.Th(""))),
			h.Tbody(mine...),
		).Density(v.DensityCompact)
	}
	out := []h.HTMLComponent{h.H4(m.SharedWithMe).Class("mb-2"), sharedComp}

	if all {
		out = append(out, h.H4(m.AllDrafts).Class("mt-6 mb-2"), b.allDrafts(ctx, m, links))
	}
	return v.VCard(v.VCardText(out...)).Variant(v.VariantOutlined)
}

// allDrafts is every draft: its owner, how it stands, its last commit, whom
// it is shared with (active), and its links.
func (b *Builder) allDrafts(ctx *web.EventContext, m *Messages, links func(key string) h.HTMLComponent) h.HTMLComponent {
	c := ctx.Context()
	entries, _ := os.ReadDir(b.Repo.DraftsDir)
	var keys []string
	for _, e := range entries {
		if _, err := os.Stat(filepath.Join(b.Repo.DraftsDir, e.Name(), ".git")); e.IsDir() && err == nil {
			keys = append(keys, e.Name())
		}
	}
	sort.Strings(keys)
	now := time.Now()
	var rows []h.HTMLComponent
	for _, key := range keys {
		d, err := b.Repo.Draft(c, key)
		if err != nil {
			continue
		}
		var state []string
		if changes, err := d.Status(c); err == nil && len(changes) > 0 {
			state = append(state, fmt.Sprintf(m.ChangesCount, len(changes)))
		}
		if s, err := d.Sync(c); err == nil {
			if s.Ahead > 0 {
				state = append(state, fmt.Sprintf(m.Ahead, s.Ahead))
			}
			if s.Behind > 0 {
				state = append(state, fmt.Sprintf(m.Behind, s.Behind))
			}
		}
		if len(state) == 0 {
			state = append(state, m.UpToDate)
		}
		last := ""
		if log, err := d.Log(c, 1); err == nil && len(log) > 0 {
			last = log[0].Subject + " · " + log[0].Author + " · " + log[0].Date.Format(dateFormat)
		}
		var with []string
		shares, _ := b.SharesOf(c, key)
		for _, s := range shares {
			if s.Active(now) {
				with = append(with, b.user(ctx.R, s.User).Label())
			}
		}
		rows = append(rows, h.Tr(
			h.Td(h.Text(b.user(ctx.R, key).Label())),
			h.Td(h.Text(strings.Join(state, " · "))),
			h.Td(h.Text(last)).Class("text-caption"),
			h.Td(h.Text(strings.Join(with, ", "))),
			h.Td(links(key)),
		).Attr("data-draft", key))
	}
	if len(rows) == 0 {
		return h.P(h.Text(m.NoDrafts)).Class("text-medium-emphasis")
	}
	return v.VTable(
		h.Thead(h.Tr(h.Th(m.DraftOwner), h.Th(m.ShareState), h.Th(m.LastCommit), h.Th(m.SharedWithCol), h.Th(""))),
		h.Tbody(rows...),
	).Density(v.DensityCompact)
}
