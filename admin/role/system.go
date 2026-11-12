package role

import (
	"context"
	"errors"
	"fmt"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// SystemRole is a role the application needs: made on boot when missing (or
// the role of its name adopted), never deleted nor renamed; its permissions
// changed as anyone's, and reset to Policies (its originals) by the action
// ActionResetPermissions. A Fixed role's permissions are the application's
// code's (Administrador: all): none of the database's.
type SystemRole struct {
	// Key is what the application knows it by ("editor"): kept in the role.
	Key string
	// Name is the role's (the subject of its permissions).
	Name string
	// Description is what the role is for, in the language of ctx.
	Description func(ctx context.Context) string
	// Policies are its original permissions (made lines: SystemPolicy.Lines).
	Policies []SystemPolicy
	// Lines are original permissions too, as lines (perm.Stack), after the
	// Policies'.
	Lines []string
	// Fixed says its permissions are the code's: none editable.
	Fixed bool
	// Formerly are originals of its former versions: a role that still has
	// one as it was loses it on boot (an original changed — narrowed, say —
	// reaches the roles made before).
	Formerly []SystemPolicy
}

// Lines are the lines of a policy: a resource each, "deny " before a
// deny's (perm.PolicyLines).
func (p SystemPolicy) Lines() []string {
	effect := p.Effect
	if effect == "" {
		effect = perm.Allowed
	}
	return perm.PolicyLines(effect, p.Actions, p.Resources)
}

// SystemPolicy is a permission of a SystemRole: allows (Effect "allow",
// the default) or denies the actions of the resources — by the unique names
// of the models and the pages ("admin:posts:*", "admin:/site-files:!git"),
// never by the groups of the menu, which change.
type SystemPolicy struct {
	Effect    string
	Actions   []string
	Resources []string
}

// Allow is a SystemPolicy that allows anything of resources.
func Allow(resources ...string) SystemPolicy {
	return SystemPolicy{Effect: perm.Allowed, Actions: []string{"*"}, Resources: resources}
}

// ActionResetPermissions is the action that resets a system role's
// permissions to its originals.
const ActionResetPermissions = "reset_permissions"

// The refusals on system roles.
var (
	ErrSystemRoleDelete = errors.New("a role of the system is not deleted")
	ErrSystemRoleRename = errors.New("a role of the system is not renamed")
	ErrSystemRoleFixed  = errors.New("the permissions of this role are the application's: they are not changed")
)

// SystemRoles registers the roles the application needs (SystemRole).
func (b *Builder) SystemRoles(roles ...SystemRole) *Builder {
	b.systemRoles = append(b.systemRoles, roles...)
	return b
}

// GetSystemRoles are the roles the application registered.
func (b *Builder) GetSystemRoles() []SystemRole { return b.systemRoles }

// systemRole is the SystemRole of key; nil when there is none.
func (b *Builder) systemRole(key string) *SystemRole {
	for i := range b.systemRoles {
		if b.systemRoles[i].Key == key {
			return &b.systemRoles[i]
		}
	}
	return nil
}

// lines are the originals of sr: its policies' lines, then its lines.
func (sr *SystemRole) lines() []string {
	return append(policiesLines(sr.Policies), sr.Lines...)
}

func policiesLines(ps []SystemPolicy) (out []string) {
	for _, p := range ps {
		out = append(out, p.Lines()...)
	}
	return
}

// EnsureSystemRoles makes the roles the application needs: a role of its key
// kept; else the role of its name adopted — its permissions kept, the
// originals it lacks added: nothing it could do lost —; else made, with its
// originals.
func (b *Builder) EnsureSystemRoles() error {
	for _, sr := range b.systemRoles {
		err := b.db.Transaction(func(tx *gorm.DB) error {
			var r Role
			err := tx.Where("system_key = ?", sr.Key).First(&r).Error
			if err == nil {
				return b.syncOriginals(tx, &r, &sr)
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if err := tx.Where("name = ?", sr.Name).First(&r).Error; err == nil {
				if err := tx.Model(&r).Update("system_key", sr.Key).Error; err != nil {
					return err
				}
				if sr.Fixed {
					return nil
				}
				return b.syncOriginals(tx, &r, &sr)
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			r = Role{ID: uuid.New(), Name: sr.Name, SystemKey: sr.Key}
			if !sr.Fixed {
				r.SystemPolicies = originalsText(sr.lines())
				r.Permissions = sr.lines()
			}
			return tx.Omit("Policies").Create(&r).Error
		})
		if err != nil {
			return fmt.Errorf("role %s: %w", sr.Name, err)
		}
	}
	return nil
}

// syncOriginals brings r to the originals of sr as they are now, by what
// changed since it was last given them (Role.SystemPolicies): an original
// gained is added at the end — unless it has it —; one lost is taken, when r
// has it as it was; what was changed by hand stays, in its place. A role
// given none yet (made before the originals were kept) loses those of its
// former versions (Formerly) it has as they were, and is given the originals
// it lacks. The originals are kept as given.
func (b *Builder) syncOriginals(tx *gorm.DB, r *Role, sr *SystemRole) error {
	if sr.Fixed {
		return nil
	}
	now := sr.lines()
	nowText := originalsText(now)
	if r.SystemPolicies == nowText {
		return nil
	}
	current := map[string]bool{}
	for _, l := range now {
		current[l] = true
	}
	// what it was given — else, its former versions
	gone := map[string]bool{}
	if r.SystemPolicies != "" {
		for _, k := range strings.Split(r.SystemPolicies, "\n") {
			if k != "" && !current[k] {
				gone[k] = true
			}
		}
	} else {
		for _, k := range policiesLines(sr.Formerly) {
			if !current[k] {
				gone[k] = true
			}
		}
	}
	given := map[string]bool{}
	for _, k := range strings.Split(r.SystemPolicies, "\n") {
		given[k] = true
	}
	var lines pq.StringArray
	seen := map[string]bool{}
	for _, l := range r.Permissions {
		if gone[l] {
			continue
		}
		lines = append(lines, l)
		seen[l] = true
	}
	for _, l := range now {
		// an original it was given and has not: taken by hand — not again
		if seen[l] || r.SystemPolicies != "" && given[l] {
			continue
		}
		lines = append(lines, l)
		seen[l] = true
	}
	r.Permissions = lines
	r.SystemPolicies = nowText
	return tx.Model(r).Select("Permissions", "SystemPolicies").Updates(r).Error
}

// originalsText is the originals of a role as Role.SystemPolicies keeps them.
func originalsText(lines []string) string {
	return strings.Join(lines, "\n")
}

// ResetPermissions sets the permissions of the system role r to its
// originals.
func (b *Builder) ResetPermissions(r *Role) error {
	sr := b.systemRole(r.SystemKey)
	if sr == nil || sr.Fixed {
		return ErrSystemRoleFixed
	}
	r.Permissions = sr.lines()
	r.SystemPolicies = originalsText(r.Permissions)
	err := b.db.Model(r).Select("Permissions", "SystemPolicies").Updates(r).Error
	if err == nil {
		b.changed()
	}
	return err
}

// installSystem guards the system roles on the model: not deleted, not
// renamed, a fixed one's permissions not changed; their mark in the listing
// and the detail; the action that resets their permissions.
func (b *Builder) installSystem(mb *presets.ModelBuilder) {
	ed := mb.Editing()
	ed.Validators.AppendFunc(func(obj any, _ presets.FieldModeStack, ctx *web.EventContext) (errs web.ValidationErrors) {
		r := obj.(*Role)
		if r.ID == uuid.Nil {
			return
		}
		var old Role
		if b.db.First(&old, "id = ?", r.ID).Error != nil || old.SystemKey == "" {
			return
		}
		r.SystemKey = old.SystemKey
		m := GetMessages(ctx.Context())
		if r.Name != old.Name {
			errs.FieldError("Name", m.ErrSystemRoleRename)
		}
		if sr := b.systemRole(old.SystemKey); sr != nil && sr.Fixed && len(r.Permissions) > 0 {
			errs.FieldError("Permissions", m.ErrSystemRoleFixed)
		}
		return
	})

	mb.Listing().WrapDeleteFunc(func(in presets.DeleteFunc) presets.DeleteFunc {
		return func(obj any, id model.ID, cascade bool, ctx *web.EventContext) error {
			var r Role
			if err := b.db.First(&r, "id = ?", id.String()).Error; err == nil && r.SystemKey != "" {
				return errors.New(GetMessages(ctx.Context()).ErrSystemRoleDelete)
			}
			return in(obj, id, cascade, ctx)
		}
	})

	// the mark: a role of the system, and what it is for
	system := func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		r := field.Obj.(*Role)
		m := GetMessages(ctx.Context())
		if r.SystemKey == "" {
			// an element, empty: the listing makes its cell of it (a text is
			// none, and the row lost a cell)
			return h.Div()
		}
		text := m.SystemRole
		if sr := b.systemRole(r.SystemKey); sr != nil && sr.Description != nil {
			text += " — " + sr.Description(ctx.Context())
		}
		return h.Div(h.Text(text)).Class("text-body-2").Attr("data-system-role", r.SystemKey)
	}
	mb.Listing().Field("SystemKey").ComponentFunc(system)
	mb.Detailing().Field("SystemKey").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Div(h.Div(h.Text(field.Label)).Class("text-caption text-medium-emphasis"), system(field, ctx)).Class("mb-4")
	})

	d := mb.Detailing()
	d.Action(ActionResetPermissions).
		Icon("mdi-restore").
		SetEnabledObj(func(obj any, _ string, _ *web.EventContext) (bool, error) {
			r := obj.(*Role)
			sr := b.systemRole(r.SystemKey)
			return sr != nil && !sr.Fixed, nil
		}).
		SetI18nLabel(func(ctx context.Context) string { return GetMessages(ctx).ResetPermissions }).
		OnClick(func(_ *web.EventContext, id string, _ any) string {
			return web.Plaid().EventFunc(eventResetPermissions).Query("id", id).Go()
		})
	mb.RegisterEventFunc(eventResetPermissions, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		var role Role
		if err = b.db.First(&role, "id = ?", ctx.R.FormValue("id")).Error; err != nil {
			return
		}
		if mb.Permissioner().ReqObjectActioner(ctx.R, &role, presets.ActionPerm(ActionResetPermissions)).Denied() {
			return r, perm.PermissionDenied
		}
		if err = b.ResetPermissions(&role); err != nil {
			return
		}
		ctx.Flash = GetMessages(ctx.Context()).PermissionsReset
		r.Reload = true
		return
	})
}

const eventResetPermissions = "role_resetPermissions"
