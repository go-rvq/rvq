package perms

import (
	"context"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
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
			return managerBody(msgs(ctx.Context()), mb, pols), nil
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
			rid, _ := mb.ParseRecordID(id)
			for _, subject := range subjects {
				// whole-record grant/revoke
				referID := referBase + ":" + subject
				// a shared policy can only be removed through the Sharing dialog,
				// never here (SPEC §: SharedID grants are managed in Sharing).
				if shared, _ := IsShared(db, referID); shared {
					continue
				}
				var err error
				if len(verbs) == 0 {
					err = Revoke(db, referID)
				} else {
					_, err = Grant(db, referID, subject, resource, verbs)
				}
				if err != nil {
					return err
				}
				// fine-grained per-field grants (restrict to specific fields)
				if err = applyFieldGrants(db, mb, rid, subject, referBase, ctx); err != nil {
					return err
				}
				// per-action grants (record-level detail actions)
				if err = applyActionGrants(db, mb, rid, subject, referBase, RecordResource(mb, rid), ctx); err != nil {
					return err
				}
				// per-page grants (auto-perm listing/detail pages)
				if err = applyPageGrants(db, mb, rid, subject, referBase, ctx); err != nil {
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

// flatField is a field flattened from the (possibly nested) mode field tree.
type flatField struct {
	Label string   // display path, e.g. "ChavesPix › Chave"
	Path  []string // field path, e.g. ["ChavesPix","Chave"]
}

func flatten(nodes []FieldNode, parent []string, label string) (out []flatField) {
	for _, n := range nodes {
		path := append(append([]string{}, parent...), n.Name)
		lbl := n.Name
		if label != "" {
			lbl = label + " › " + n.Name
		}
		if len(n.Children) == 0 {
			out = append(out, flatField{Label: lbl, Path: path})
		} else {
			out = append(out, flatten(n.Children, path, lbl)...)
		}
	}
	return
}

// fieldFormKey encodes a mode+field path into a stable form field name.
func fieldFormKey(mode string, path []string) string {
	return "perms_f_" + mode + "_" + strings.Join(path, "__")
}

// fieldMatrix renders, for a mode, the flattened fields as checkboxes.
func fieldMatrix(mb *presets.ModelBuilder, mode, title string) h.HTMLComponent {
	fields := flatten(FieldNodes(mb, mode), nil, "")
	if len(fields) == 0 {
		return nil
	}
	var boxes h.HTMLComponents
	for _, f := range fields {
		boxes = append(boxes, v.VCheckbox().Label(f.Label).Density(v.DensityCompact).HideDetails(true).
			Attr(web.VField(fieldFormKey(mode, f.Path), false)...))
	}
	return v.VExpansionPanel(
		v.VExpansionPanelTitle(h.Text(title)),
		v.VExpansionPanelText(boxes...),
	)
}

// applyFieldGrants grants/revokes the per-field permissions selected in the
// dialog for one subject: detail-mode fields grant the view verb, edit-mode
// fields the edit verb, on the field's exact resource.
func applyFieldGrants(db *gorm.DB, mb *presets.ModelBuilder, rid model.ID, subject, referBase string, ctx *web.EventContext) error {
	apply := func(mode string, verbs []string) error {
		for _, f := range flatten(FieldNodes(mb, mode), nil, "") {
			res := FieldResource(mb, rid, f.Path...)
			referID := referBase + ":f:" + mode + ":" + strings.Join(f.Path, ".") + ":" + subject
			if ctx.R.FormValue(fieldFormKey(mode, f.Path)) == "true" {
				if _, err := Grant(db, referID, subject, res, verbs); err != nil {
					return err
				}
			} else if err := Revoke(db, referID); err != nil {
				return err
			}
		}
		return nil
	}
	if err := apply(ModeDetail, VerbView); err != nil {
		return err
	}
	return apply(ModeEdit, VerbEdit)
}

// actionFormKey encodes an action verb into a stable form field name.
func actionFormKey(verb string) string { return "perms_a_" + verb }

// actionMatrix renders the record-level actions as checkboxes.
func actionMatrix(mb *presets.ModelBuilder, title string) h.HTMLComponent {
	actions := ActionNodes(mb)
	if len(actions) == 0 {
		return nil
	}
	var boxes h.HTMLComponents
	for _, a := range actions {
		boxes = append(boxes, v.VCheckbox().Label(a.Name).Density(v.DensityCompact).HideDetails(true).
			Attr(web.VField(actionFormKey(a.Verb), false)...))
	}
	return v.VExpansionPanel(
		v.VExpansionPanelTitle(h.Text(title)),
		v.VExpansionPanelText(boxes...),
	)
}

// applyActionGrants grants/revokes the per-action permissions selected in the
// dialog for one subject: a checked action allows its verb on the record.
func applyActionGrants(db *gorm.DB, mb *presets.ModelBuilder, rid model.ID, subject, referBase, record string, ctx *web.EventContext) error {
	for _, a := range ActionNodes(mb) {
		referID := referBase + ":a:" + a.Verb + ":" + subject
		if ctx.R.FormValue(actionFormKey(a.Verb)) == "true" {
			if _, err := Grant(db, referID, subject, record, []string{a.Verb}); err != nil {
				return err
			}
		} else if err := Revoke(db, referID); err != nil {
			return err
		}
	}
	return nil
}

// pageFormKey encodes a page (mode+path) into a stable form field name.
func pageFormKey(mode, pagePath string) string { return "perms_p_" + mode + "_" + pagePath }

// pageMatrix renders the (detail) auto-perm pages as checkboxes.
func pageMatrix(mb *presets.ModelBuilder, title string) h.HTMLComponent {
	pages := PageNodes(mb)
	if len(pages) == 0 {
		return nil
	}
	var boxes h.HTMLComponents
	for _, p := range pages {
		boxes = append(boxes, v.VCheckbox().Label(p.Name).Density(v.DensityCompact).HideDetails(true).
			Attr(web.VField(pageFormKey(p.Mode, p.Path), false)...))
	}
	return v.VExpansionPanel(
		v.VExpansionPanelTitle(h.Text(title)),
		v.VExpansionPanelText(boxes...),
	)
}

// applyPageGrants grants/revokes the per-page permissions selected in the
// dialog for one subject: a checked page allows the page (its path is the perm
// segment) on the record/listing resource.
func applyPageGrants(db *gorm.DB, mb *presets.ModelBuilder, rid model.ID, subject, referBase string, ctx *web.EventContext) error {
	for _, p := range PageNodes(mb) {
		res := PageResource(mb, p.Mode, rid, p.Path)
		referID := referBase + ":p:" + p.Mode + ":" + p.Path + ":" + subject
		if ctx.R.FormValue(pageFormKey(p.Mode, p.Path)) == "true" {
			// the page path is the perm segment, granted with the view verb
			if _, err := Grant(db, referID, subject, res, VerbView); err != nil {
				return err
			}
		} else if err := Revoke(db, referID); err != nil {
			return err
		}
	}
	return nil
}

func managerBody(m *Messages, mb *presets.ModelBuilder, pols []perm.DefaultDBPolicy) h.HTMLComponent {
	rows := h.HTMLComponents{h.Tr(h.Th(m.Subject), h.Th(m.Permissions), h.Th(m.Shared))}
	if len(pols) == 0 {
		rows = append(rows, h.Tr(h.Td(h.Text("—")), h.Td(h.Text(m.NoGrants)), h.Td(h.Text(""))))
	}
	for _, p := range pols {
		// SharedID is shown read-only; shared policies are managed in the Sharing
		// dialog and must not be removed from the permission manager.
		sharedCell := h.Td(h.Text(""))
		if p.SharedID != nil {
			sharedCell = h.Td(
				v.VChip(h.Text(p.SharedID.String())).Size(v.SizeXSmall).Attr("readonly", true).Class("mr-1"),
				h.Span(m.SharedNote).Class("text-caption text-medium-emphasis"),
			)
		}
		rows = append(rows, h.Tr(
			h.Td(h.Text(p.Subject)),
			h.Td(h.Text(strings.Join([]string(p.Actions), ", "))),
			sharedCell,
		))
	}
	// fine-grained per-field matrix (nested fields included), by mode
	var panels h.HTMLComponents
	if p := fieldMatrix(mb, ModeDetail, m.FieldsView); p != nil {
		panels = append(panels, p)
	}
	if p := fieldMatrix(mb, ModeEdit, m.FieldsEdit); p != nil {
		panels = append(panels, p)
	}
	if p := actionMatrix(mb, m.Actions); p != nil {
		panels = append(panels, p)
	}
	if p := pageMatrix(mb, m.Pages); p != nil {
		panels = append(panels, p)
	}

	return h.Div(
		v.VAlert(h.Text(m.Help)).Type("info").Variant(v.VariantTonal).Density(v.DensityComfortable).Class("mb-4"),
		v.VTable(h.Tbody(rows...)).Density(v.DensityCompact).Class("mb-4"),
		v.VTextarea().Label(m.Subject).Rows(2).Hint(m.SubjectHint).PersistentHint(true).
			Attr(web.VField(fieldSubject, "")...),
		v.VCheckbox().Label(m.View).Attr(web.VField(fieldView, false)...),
		v.VCheckbox().Label(m.Edit).Attr(web.VField(fieldEdit, false)...),
		v.VCheckbox().Label(m.Delete).Attr(web.VField(fieldDelete, false)...),
		h.Div(h.Text(m.FieldsHelp)).Class("text-caption text-medium-emphasis mt-2"),
		v.VExpansionPanels(panels...).Class("my-2"),
		h.Div(h.Text(m.RevokeHint)).Class("text-caption text-medium-emphasis"),
	)
}
