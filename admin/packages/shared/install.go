package shared

import (
	"context"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/packages/perms"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	fieldSubjects = "share_subjects"
	fieldTemplate = "share_template"
	fieldView     = "share_view"
	fieldEdit     = "share_edit"
	fieldDelete   = "share_delete"
)

// ScopeFunc resolves the permission-template scope (typically the organization
// id) for the record being shared, so the dialog can offer that scope's
// templates. Return nil for global templates only.
type ScopeFunc func(ctx *web.EventContext) *uuid.UUID

// Option configures Install.
type Option func(*shareConfig)

type shareConfig struct {
	notifier Notifier
	scope    ScopeFunc
}

// WithNotifier notifies the invited user and the inviter of share events.
func WithNotifier(n Notifier) Option { return func(c *shareConfig) { c.notifier = n } }

// WithTemplates lets the dialog offer named permission templates (from the
// scope's "configuration") instead of raw verb checkboxes. scope resolves the
// template scope (e.g. the organization id) for the record being shared.
func WithTemplates(scope ScopeFunc) Option { return func(c *shareConfig) { c.scope = scope } }

// Install adds a permissioned "Share" detail action to mb: a dialog listing who
// the record is currently shared with and letting an authorized subject invite
// more users. Sharing sends a pending invite per user (notified via a Notifier
// when WithNotifier is set); access is only granted once the invited user
// accepts (see Accept), which writes the share (perm.DefaultDBPolicy.SharedID).
// Verbs come from a permission template (WithTemplates) or from view/edit/delete
// checkboxes. Mount it only on aggregate resources whose owner controls access
// (e.g. an organization or a project). Guarded by the "share" verb.
func Install(mb *presets.ModelBuilder, db *gorm.DB, opts ...Option) {
	if !mb.HasDetailing() {
		return
	}
	cfg := &shareConfig{}
	for _, o := range opts {
		o(cfg)
	}
	n := cfg.notifier
	if err := AutoMigrateInvites(db); err != nil {
		panic(err)
	}
	if cfg.scope != nil {
		if err := AutoMigrateTemplates(db); err != nil {
			panic(err)
		}
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
			var templates []ShareTemplate
			if cfg.scope != nil {
				templates, _ = Templates(db, cfg.scope(ctx))
			}
			return shareBody(msgs(ctx.Context()), pols, templates), nil
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
			// a chosen template wins over the checkboxes
			if tid := ctx.R.FormValue(fieldTemplate); tid != "" {
				if id, e := uuid.Parse(tid); e == nil {
					actions, _ = TemplateActions(db, id)
				}
			}
			if len(actions) == 0 {
				if ctx.R.FormValue(fieldView) == "true" {
					actions = append(actions, VerbView...)
				}
				if ctx.R.FormValue(fieldEdit) == "true" {
					actions = append(actions, VerbEdit...)
				}
				if ctx.R.FormValue(fieldDelete) == "true" {
					actions = append(actions, VerbDelete...)
				}
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

func shareBody(m *Messages, pols []perm.DefaultDBPolicy, templates []ShareTemplate) h.HTMLComponent {
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

	// permission input: a template selector when templates exist, otherwise
	// the raw verb checkboxes.
	var permInput h.HTMLComponent
	if len(templates) > 0 {
		items := make([]map[string]string, len(templates))
		for i, t := range templates {
			items[i] = map[string]string{"title": t.Name, "value": t.ID.String()}
		}
		permInput = v.VSelect().Label(m.Permissions).Items(items).
			ItemTitle("title").ItemValue("value").Clearable(true).
			Attr(web.VField(fieldTemplate, "")...)
	} else {
		permInput = h.Div(
			v.VCheckbox().Label(m.CanView).Attr(web.VField(fieldView, true)...),
			v.VCheckbox().Label(m.CanEdit).Attr(web.VField(fieldEdit, false)...),
			v.VCheckbox().Label(m.CanDelete).Attr(web.VField(fieldDelete, false)...),
		)
	}

	return h.Div(
		v.VTable(h.Tbody(rows...)).Density(v.DensityCompact).Class("mb-4"),
		v.VTextarea().Label(m.ShareWith).Rows(2).Hint(m.ShareWithHint).PersistentHint(true).
			Attr(web.VField(fieldSubjects, "")...),
		permInput,
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
