package perms

import (
	"context"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

// manager dialog form fields
const (
	fieldSubject = "perms_subject"
	fieldView    = "perms_view"
	fieldEdit    = "perms_edit"
	fieldDelete  = "perms_delete"
)

// InstallManager adds a permissioned "Manage Permissions" detail action to mb.
// It is a shortcut to the perm API: the dialog lists the record's grants and
// lets an authorized subject grant/revoke a set of roles/users' verbs
// (view/edit/delete) over the record, writing perm.DefaultDBPolicy rows (record
// resource + glob wildcard covering fields, actions and pages) and reloading
// the verifier. The action is itself guarded by the PermManage verb.
func InstallManager(mb *presets.ModelBuilder, db *gorm.DB) {
	if !mb.HasDetailing() {
		return
	}
	registerMessages(mb.Builder())

	reload := func() {
		if pb := mb.Builder().GetPermission(); pb != nil {
			from := time.Now().Add(-time.Minute)
			pb.LoadDBPoliciesToMemory(db, &from)
		}
	}

	scope := func(id string) (referBase, resource string, ok bool) {
		rid, err := mb.ParseRecordID(id)
		if err != nil || rid.IsZero() {
			return "", "", false
		}
		res := RecordResource(mb, rid)
		return "perms:" + res, res + "*", true
	}

	mb.Detailing().Action(PermManage).
		Icon("mdi-shield-key").
		SetI18nLabel(func(ctx context.Context) string { return msgs(ctx).ManageTitle }).
		SetVerifier(func(ctx *web.EventContext) *perm.Verifier {
			return mb.Permissioner().ReqListDo(ctx.R, PermManage)
		}).
		ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
			_, resource, ok := scope(id)
			if !ok {
				return h.Text(""), nil
			}
			pols, err := List(db, resource)
			if err != nil {
				return nil, err
			}
			return managerBody(msgs(ctx.Context()), pols), nil
		}).
		UpdateFunc(func(id string, ctx *web.EventContext) error {
			referBase, resource, ok := scope(id)
			if !ok {
				return nil
			}
			subjects := parseSubjects(ctx.R.FormValue(fieldSubject))
			if len(subjects) == 0 {
				return nil
			}
			var verbs []string
			if ctx.R.FormValue(fieldView) == "true" {
				verbs = append(verbs, VerbView...)
			}
			if ctx.R.FormValue(fieldEdit) == "true" {
				verbs = append(verbs, VerbEdit...)
			}
			if ctx.R.FormValue(fieldDelete) == "true" {
				verbs = append(verbs, VerbDelete...)
			}
			for _, subject := range subjects {
				referID := referBase + ":" + subject
				var err error
				if len(verbs) == 0 {
					err = Revoke(db, referID)
				} else {
					_, err = Grant(db, referID, subject, resource, verbs)
				}
				if err != nil {
					return err
				}
			}
			reload()
			return nil
		})
}

// parseSubjects splits a multi-subject input into trimmed, de-duplicated
// roles/users (separated by commas, semicolons or newlines).
func parseSubjects(s string) []string {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r'
	})
	seen := map[string]bool{}
	var out []string
	for _, f := range fields {
		if f = strings.TrimSpace(f); f != "" && !seen[f] {
			seen[f] = true
			out = append(out, f)
		}
	}
	return out
}

func managerBody(m *Messages, pols []perm.DefaultDBPolicy) h.HTMLComponent {
	rows := h.HTMLComponents{h.Tr(h.Th(m.Subject), h.Th(m.Permissions))}
	if len(pols) == 0 {
		rows = append(rows, h.Tr(h.Td(h.Text("—")), h.Td(h.Text(m.NoGrants))))
	}
	for _, p := range pols {
		rows = append(rows, h.Tr(
			h.Td(h.Text(p.Subject)),
			h.Td(h.Text(strings.Join([]string(p.Actions), ", "))),
		))
	}
	return h.Div(
		v.VAlert(h.Text(m.Help)).Type("info").Variant(v.VariantTonal).Density(v.DensityComfortable).Class("mb-4"),
		v.VTable(h.Tbody(rows...)).Density(v.DensityCompact).Class("mb-4"),
		v.VTextarea().Label(m.Subject).Rows(2).Hint(m.SubjectHint).PersistentHint(true).
			Attr(web.VField(fieldSubject, "")...),
		v.VCheckbox().Label(m.View).Attr(web.VField(fieldView, false)...),
		v.VCheckbox().Label(m.Edit).Attr(web.VField(fieldEdit, false)...),
		v.VCheckbox().Label(m.Delete).Attr(web.VField(fieldDelete, false)...),
		h.Div(h.Text(m.RevokeHint)).Class("text-caption text-medium-emphasis"),
	)
}
