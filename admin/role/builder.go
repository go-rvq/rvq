package role

import (
	"fmt"
	"github.com/lib/pq"
	"net/http"
	"strings"
	"sync/atomic"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/admin/userdocs"
	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/ory/ladon"
	"gorm.io/gorm"
)

type Builder struct {
	db        *gorm.DB
	actions   []*DefaultOptionItem
	resources []*DefaultOptionItem
	// editorSubject is the subject that has permission to edit roles
	// empty value means anyone can edit roles
	editorSubject    string
	roleMb           *presets.ModelBuilder
	AfterInstallFunc presets.ModelInstallFunc
	// systemRoles are the roles the application needs (SystemRoles)
	systemRoles []SystemRole
	pb          *presets.Builder
	// permissionsHelp is the URL of the documentation of the permissions,
	// opened by the help button of the form's permissions ("" none)
	permissionsHelp func(ctx *web.EventContext) string
}

// PermissionsHelpURL sets where the help button (?) of the permissions of the
// role's form leads: the documentation of policies and permissions. Without
// it, or when it says "", there is no button.
func (b *Builder) PermissionsHelpURL(f func(ctx *web.EventContext) string) *Builder {
	b.permissionsHelp = f
	return b
}

func New(db *gorm.DB) *Builder {
	uuidkey.MustRegister(db) // its records have UUID keys
	return &Builder{
		db: db,
		actions: []*DefaultOptionItem{
			{Text: "All", Value: "*"},
			{Text: "List", Value: presets.PermList},
			{Text: "Get", Value: presets.PermGet},
			{Text: "Create", Value: presets.PermCreate},
			{Text: "Update", Value: presets.PermUpdate},
			{Text: "Delete", Value: presets.PermDelete},
			{Text: "Delete With Related", Value: presets.PermDeleteWithRelated},
		},
	}
}

func (b *Builder) Actions(vs []*DefaultOptionItem) *Builder {
	b.actions = vs
	return b
}

func (b *Builder) AfterInstall(v presets.ModelInstallFunc) *Builder {
	b.AfterInstallFunc = v
	return b
}

func (b *Builder) Resources(vs []*DefaultOptionItem) *Builder {
	b.resources = vs
	return b
}

func (b *Builder) EditorSubject(v string) *Builder {
	b.editorSubject = v
	return b
}

func (b *Builder) Install(pb *presets.Builder) (err error) {
	ConfigureMessages(pb.I18n())
	b.pb = pb

	if b.editorSubject != "" {
		permB := pb.GetPermission()
		if permB == nil {
			panic("pb does not have a permission builder")
		}
		ctxf := permB.GetContextFunc()
		ssf := permB.GetSubjectsFunc()
		permB.ContextFunc(func(r *http.Request, objs []interface{}) perm.Context {
			c := make(perm.Context)
			if ctxf != nil {
				c = ctxf(r, objs)
			}
			ss := ssf(r)
			hasRoleEditorSubject := false
			for _, s := range ss {
				if s == b.editorSubject {
					hasRoleEditorSubject = true
					break
				}
			}
			c["has_role_editor_subject"] = hasRoleEditorSubject
			return c
		})
		permB.CreatePolicies(
			perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(perm.Anything).On("*:roles:*").Given(perm.Conditions{
				"has_role_editor_subject": &ladon.BooleanCondition{
					BooleanValue: false,
				},
			}),
			perm.PolicyFor(b.editorSubject).WhoAre(perm.Allowed).ToDo(perm.Anything).On("*:roles:*"),
		)
	}

	b.roleMb = pb.Model(&Role{}, presets.ModelConfig().
		SetModuleKey(MessagesKey)).
		MenuIcon("mdi-account-key")

	ed := b.roleMb.Editing(
		"Name",
		"Permissions",
	)

	// its lines, in order (perm.Stack), edited as one value — the tree of the
	// permissions of the admin it touches (userdocs.PermissionInput) —, any
	// part of the admin; stored a line a scope, what is given to a resource
	// by its unique name
	permOpts := func(ctx *web.EventContext) (userdocs.PermissionInputOptions, error) {
		o := userdocs.PermissionInputOptions{URL: b.roleMb.Info().ListingHref(), Event: PermissionInputEvent,
			StringView: true}
		// the help of the permissions: a ? in the bar of the field
		if b.permissionsHelp != nil {
			o.HelpURL, o.HelpTitle = b.permissionsHelp(ctx), GetMessages(ctx.Context()).PermissionsHelp
		}
		return o, nil
	}
	ed.Field("Permissions").ComponentFunc(userdocs.PermissionFieldComponentFunc(pb, permOpts)).
		SetterFunc(userdocs.PermissionFieldSetter)
	b.roleMb.RegisterEventFunc(PermissionInputEvent, userdocs.PermissionInputEventFunc(pb, permOpts))

	ed.Validators.AppendFunc(func(obj interface{}, _ presets.FieldModeStack, ctx *web.EventContext) (errs web.ValidationErrors) {
		r := obj.(*Role)
		var lines pq.StringArray
		for _, l := range r.Permissions {
			if l = strings.TrimSpace(l); l == "" {
				continue
			}
			if perm.CheckLine(l) != nil {
				errs.FieldError("Permissions", fmt.Sprintf(GetMessages(ctx.Context()).ErrLine, l))
			}
			lines = append(lines, l)
		}
		r.Permissions = lines
		return
	})

	ed.SaveFunc(func(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
		if err = gorm2op.DataOperator(b.db.Omit("Policies")).Save(obj, id, ctx); err == nil {
			b.changed()
		}
		return
	})

	// the lines in the detail, one a line, as written
	b.roleMb.Detailing().Field("Permissions").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		r := field.Obj.(*Role)
		var items []h.HTMLComponent
		for _, l := range r.Permissions {
			items = append(items, h.Div(h.Code(l)).Class("text-body-2"))
		}
		return h.Div(h.Div(h.Text(field.Label)).Class("text-caption text-medium-emphasis"),
			h.Div(items...).Attr("data-role-lines", true)).Class("mb-4")
	})

	b.roleMb.Listing().DeleteFunc(func(obj interface{}, id model.ID, cascade bool, ctx *web.EventContext) (err error) {
		err = b.db.Transaction(func(tx *gorm.DB) error {
			// the id as text: a model.ID is no SQL value
			if err := tx.Delete(&perm.DefaultDBPolicy{}, "refer_id = ?", id.String()).Error; err != nil {
				return err
			}
			if err := tx.Delete(&Role{}, "id = ?", id.String()).Error; err != nil {
				return err
			}

			return nil
		})
		if err == nil {
			b.changed()
		}
		return
	})

	b.roleMb.Listing("Name", "SystemKey")
	b.roleMb.Detailing("Name", "SystemKey", "Permissions")
	b.installSystem(b.roleMb)
	// the roles of before the lines: their policies made lines, once
	if err = MigrateLines(b.db); err != nil {
		return err
	}
	// what was given to resources inside the groups of the menu: by their
	// unique names; the placeholders out (perm.CleanLines)
	if err = MigrateUngroup(b.db); err != nil {
		return err
	}
	if len(b.systemRoles) > 0 {
		if err = b.EnsureSystemRoles(); err != nil {
			return err
		}
		b.changed()
	}

	if b.AfterInstallFunc != nil {
		return b.AfterInstallFunc(b.roleMb)
	}

	return nil
}

// generation is the version of the roles — of their lines —: changed by each
// change of a role, it is part of the version of the stacks of the users
// (user.StackSource), which are made again.
var generation atomic.Int64

// Generation is the version of the roles in this process.
func Generation() int64 { return generation.Load() }

// Changed tells the roles changed — by this process, or by another one (a
// listener): the stacks of the users are made again at their next request.
// PermissionInputEvent is the event of the field of the permissions of the
// roles (userdocs.PermissionInputEventFunc).
const PermissionInputEvent = "role_PermissionInput"

func Changed(pb *perm.Builder) {
	generation.Add(1)
	if pb != nil {
		pb.Forget()
	}
}

func (b *Builder) changed() {
	var pb *perm.Builder
	if b.pb != nil {
		pb = b.pb.GetPermission()
	}
	Changed(pb)
}

// MigrateLines makes the permissions of the roles of before the lines —
// policies whose ReferID is the role's key — their lines, once: a role of no
// lines and of policies gets theirs (perm.PolicyLines), in the order they
// were made, and its policies go; the originals of a system role kept in the
// form of before ("effect|actions|resources" a line) are made lines too.
func MigrateLines(db *gorm.DB) error {
	if err := db.AutoMigrate(&Role{}); err != nil {
		return err
	}
	if !db.Migrator().HasTable(&perm.DefaultDBPolicy{}) {
		return nil
	}
	var roles []*Role
	if err := db.Find(&roles).Error; err != nil {
		return err
	}
	for _, r := range roles {
		if err := db.Where("refer_id = ?", r.ID.String()).Order("created_at, id").Find(&r.Policies).Error; err != nil {
			return err
		}
		if len(r.Permissions) > 0 || len(r.Policies) == 0 && !strings.Contains(r.SystemPolicies, "|") {
			continue
		}
		var lines pq.StringArray
		for _, p := range r.Policies {
			lines = append(lines, perm.PolicyLines(p.Effect, p.Actions, p.Resources)...)
		}
		originals := r.SystemPolicies
		if strings.Contains(originals, "|") {
			var ls []string
			for _, k := range strings.Split(originals, "\n") {
				parts := strings.SplitN(k, "|", 3)
				if len(parts) != 3 {
					continue
				}
				ls = append(ls, perm.PolicyLines(parts[0], strings.Split(parts[1], ","), []string{parts[2]})...)
			}
			originals = strings.Join(ls, "\n")
		}
		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(r).UpdateColumns(map[string]any{"permissions": lines, "system_policies": originals}).Error; err != nil {
				return err
			}
			return tx.Where("refer_id = ?", r.ID.String()).Delete(&perm.DefaultDBPolicy{}).Error
		})
		if err != nil {
			return fmt.Errorf("role %s: its lines: %w", r.Name, err)
		}
	}
	return nil
}

// MigrateUngroup cleans the lines of the roles (perm.CleanLines): what they
// give to resources inside the groups of the menu by their unique names
// ("admin:content/:posts:*" is "admin:posts:*") — what is given to a
// resource is stored out of its groups: a resource moved to another group
// keeps its permissions —, and no placeholder ("admin:users:", which decides
// nothing). Only the roles whose lines change; once done, nothing more.
func MigrateUngroup(db *gorm.DB) error {
	var roles []*Role
	if err := db.Select("id", "name", "permissions").Find(&roles).Error; err != nil {
		return err
	}
	for _, r := range roles {
		lines, changed := perm.CleanLines(r.Permissions)
		if !changed {
			continue
		}
		if err := db.Model(r).UpdateColumn("permissions", pq.StringArray(lines)).Error; err != nil {
			return fmt.Errorf("role %s: its lines out of the groups: %w", r.Name, err)
		}
	}
	return nil
}

// Stack is the stack of a user of roles, in their order: the policies of the
// code whose subjects match subjects (perm.Builder.CodeLines) — the code —,
// then the database's (perm.Builder.DBLines: what is granted a record, a
// share), then the lines of each role, of the origin "role: <name>". Not
// resolved.
func Stack(pb *perm.Builder, subjects []string, roles []*Role) *perm.Stack {
	s := &perm.Stack{}
	if pb != nil {
		for _, l := range append(pb.CodeLines(subjects), pb.DBLines(subjects)...) {
			s.AddLine(l)
		}
	}
	for _, r := range roles {
		s.Add("role: "+r.Name, r.Permissions...)
	}
	return s
}
