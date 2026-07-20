package shared

import (
	"context"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/helper/user"
	"github.com/go-rvq/rvq/admin/packages/perms"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

const (
	fieldSubjects = "share_subjects"
	fieldView     = "share_view"
	fieldEdit     = "share_edit"
	fieldDelete   = "share_delete"
)

// Install adds a permissioned "Share" detail action to mb: a dialog listing who
// the record is currently shared with and letting an authorized subject invite
// more users (with view/edit/delete). Sharing sends a pending invite per user
// (notified via the optional Notifier); access is only granted once the invited
// user accepts (see Accept), which writes the share (perm.DefaultDBPolicy.
// SharedID). Mount it only on aggregate resources whose owner controls access
// (e.g. an organization or a project). The action is guarded by the "share" verb.
func Install(mb *presets.ModelBuilder, db *gorm.DB, notifier ...Notifier) {
	if !mb.HasDetailing() {
		return
	}
	if err := AutoMigrateInvites(db); err != nil {
		panic(err)
	}
	var n Notifier
	if len(notifier) > 0 {
		n = notifier[0]
	}
	registerMessages(mb.Builder().I18n())

	reload := func() {
		if pb := mb.Builder().GetPermission(); pb != nil {
			from := time.Now().Add(-time.Minute)
			pb.LoadDBPoliciesToMemory(db, &from)
		}
	}

	resourceOf := func(id string, ctx *web.EventContext) (string, bool) {
		rid, err := mb.ParseRecordID(id)
		if err != nil || rid.IsZero() {
			return "", false
		}
		return perms.RecordResource(mb, rid, presets.ParentsModelID(ctx.R)...) + "*", true
	}

	mb.Detailing().Action(PermShare).
		Icon("mdi-account-multiple-plus").
		SetI18nLabel(func(ctx context.Context) string { return msgs(ctx).ShareTitle }).
		SetVerifier(func(ctx *web.EventContext) *perm.Verifier {
			return mb.Permissioner().ReqListDo(ctx.R, PermShare)
		}).
		ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
			resource, ok := resourceOf(id, ctx)
			if !ok {
				return h.Text(""), nil
			}
			pols, err := ForResource(db, resource)
			if err != nil {
				return nil, err
			}
			return shareBody(msgs(ctx.Context()), pols), nil
		}).
		UpdateFunc(func(id string, ctx *web.EventContext) error {
			resource, ok := resourceOf(id, ctx)
			if !ok {
				return nil
			}
			subjects := parseSubjects(ctx.R.FormValue(fieldSubjects))
			if len(subjects) == 0 {
				return nil
			}
			var actions []string
			if ctx.R.FormValue(fieldView) == "true" {
				actions = append(actions, VerbView...)
			}
			if ctx.R.FormValue(fieldEdit) == "true" {
				actions = append(actions, VerbEdit...)
			}
			if ctx.R.FormValue(fieldDelete) == "true" {
				actions = append(actions, VerbDelete...)
			}
			if len(actions) == 0 {
				actions = VerbView
			}
			var invitedBy string
			if u := user.GetCurrentUser(ctx.R); u != nil {
				invitedBy = u.GetID().String()
			}
			// send pending invites — access is granted on acceptance, not now.
			if _, err := Invite(db, resource, subjects, actions, invitedBy, n); err != nil {
				return err
			}
			reload()
			return nil
		})
}

func shareBody(m *Messages, pols []perm.DefaultDBPolicy) h.HTMLComponent {
	rows := h.HTMLComponents{h.Tr(h.Th(m.SharedWith), h.Th(""))}
	if len(pols) == 0 {
		rows = append(rows, h.Tr(h.Td(h.Text("—")), h.Td(h.Text(m.NoShares))))
	}
	for _, p := range pols {
		rows = append(rows, h.Tr(
			h.Td(h.Text(p.Subject)),
			h.Td(h.Text(strings.Join([]string(p.Actions), ", "))),
		))
	}
	return h.Div(
		v.VTable(h.Tbody(rows...)).Density(v.DensityCompact).Class("mb-4"),
		v.VTextarea().Label(m.ShareWith).Rows(2).Hint(m.ShareWithHint).PersistentHint(true).
			Attr(web.VField(fieldSubjects, "")...),
		v.VCheckbox().Label(m.CanView).Attr(web.VField(fieldView, true)...),
		v.VCheckbox().Label(m.CanEdit).Attr(web.VField(fieldEdit, false)...),
		v.VCheckbox().Label(m.CanDelete).Attr(web.VField(fieldDelete, false)...),
	)
}

// parseSubjects splits a multi-subject input into trimmed, de-duplicated users
// (separated by commas, semicolons or newlines).
func parseSubjects(s string) []string {
	seen := map[string]bool{}
	var out []string
	for _, f := range strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r'
	}) {
		if f = strings.TrimSpace(f); f != "" && !seen[f] {
			seen[f] = true
			out = append(out, f)
		}
	}
	return out
}
