package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleJoinTable is the table of the roles of the users; its column position
// is the order of a user's roles — the order of the user's stack.
const RoleJoinTable = "user_role_join"

// EnsureRolePositions gives the table of the roles of the users its column
// position (the order of a user's roles), when it lacks it.
func EnsureRolePositions(db *gorm.DB) error {
	m := db.Migrator()
	if !m.HasTable(RoleJoinTable) || m.HasColumn(RoleJoinTable, "position") {
		return nil
	}
	return db.Exec(`ALTER TABLE ` + RoleJoinTable + ` ADD COLUMN position integer NOT NULL DEFAULT 0`).Error
}

// OrderedRoles are the roles of the user id, in the user's order (then by
// name).
func OrderedRoles(db *gorm.DB, id uuid.UUID) (roles []*role.Role, err error) {
	err = db.Joins("JOIN "+RoleJoinTable+" j ON j.role_id = roles.id").
		Where("j.user_id = ?", id).Order("j.position, roles.name").Find(&roles).Error
	return
}

// SetRolePositions keeps the order of the roles of the user id: each its
// place in ids.
func SetRolePositions(db *gorm.DB, id uuid.UUID, ids []uuid.UUID) error {
	for i, rid := range ids {
		if err := db.Table(RoleJoinTable).Where("user_id = ? AND role_id = ?", id, rid).
			Update("position", i).Error; err != nil {
			return err
		}
	}
	return nil
}

// SubjectOf is the subject of the stack of the user id: "user:<id>".
func SubjectOf(id uuid.UUID) string { return "user:" + id.String() }

// Stack is the stack of the permissions of the user u, resolved: the
// policies of the code for the user's subjects (UserRoles, and the user's
// own) first, then the lines of the user's roles, in the user's order (the
// roles u carries: WithRoles loads them so).
//
// The groups of the user are in it after them: the lines of the database
// for each group (group:<id>), and each group's roles after the user's, the
// groups in the user's order.
func (b *Builder) Stack(pb *perm.Builder, u User) *perm.Stack {
	subjects := append(b.UserRoles(u), SubjectOf(u.GetID()))
	var roles []*role.Role
	if u.GetAccountName() != b.lb.GetInitialUserAccount() {
		roles = u.GetRoles()
	}
	if b.db != nil && !u.Anonymous() {
		groups, _ := GroupsOf(b.db, u.GetID())
		for _, g := range groups {
			subjects = append(subjects, GroupSubject(g.ID))
			grs, _ := GroupRolesOf(b.db, g.ID)
			roles = append(roles, grs...)
		}
	}
	return role.Stack(pb, subjects, roles).Resolve()
}

// NewUser is a new record of the users' model.
func (b *Builder) NewUser() User { return b.mb.NewModel().(User) }

// Find is the user id (nil when none).
func (b *Builder) Find(db *gorm.DB, id uuid.UUID) (User, error) {
	u := b.mb.NewModel().(User)
	if err := db.Where("id = ?", id).First(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// FindByAccount is the user of the account (nil when none).
func (b *Builder) FindByAccount(db *gorm.DB, account string) (User, error) {
	u := b.mb.NewModel().(User)
	if err := db.Where("account = ?", account).First(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// StackSource is the source of the stacks of the users (perm.Builder.Stacks):
// the user of a request — none for the anonymous one, whose requests the
// policies decide —, of the version of its roles, in its order, and of the
// roles (role.Generation).
func (b *Builder) StackSource(pb *perm.Builder) perm.StackSource {
	return stackSource{b: b, pb: pb}
}

type stackSource struct {
	b  *Builder
	pb *perm.Builder
}

func (s stackSource) StackOf(r *http.Request) (string, string, func() (*perm.Stack, error)) {
	u := GetCurrentUser(r)
	if u == nil || u.Anonymous() {
		return "", "", nil
	}
	ids := make([]string, 0, len(u.GetRoles()))
	for _, ro := range u.GetRoles() {
		ids = append(ids, ro.ID.String())
	}
	version := strconv.FormatInt(role.Generation(), 10) + "|" + strconv.FormatInt(GroupGeneration(), 10) +
		"|" + strings.Join(ids, ",")
	return SubjectOf(u.GetID()), version, func() (*perm.Stack, error) {
		return s.b.Stack(s.pb, u), nil
	}
}

// installPermissionsPage puts the page of the permissions of a user in its
// detail: its stack — Tudo, every line in its place, where it comes from, its
// state (active, replaced by which, a removal) —, and the active lines —
// Ativas, a list, or a tree.
func (b *Builder) installPermissionsPage() {
	pb := b.mb.Builder().GetPermission()
	page := b.mb.Detailing().PagesRegistrator().New(
		presets.HttpPage("/permissions").
			MenuIcon("mdi-shield-key-outline").
			TitleFunc(func(ctx context.Context) string { return GetMessages(ctx).Permissions })).
		RecordLayout(func(ctx *web.EventContext, mid model.ID, _ []model.ID, fetch func() (any, error)) (r web.PageResponse, err error) {
			obj, err := fetch()
			if err != nil {
				return
			}
			u := obj.(User)
			roles, err := OrderedRoles(b.db, u.GetID())
			if err != nil {
				return
			}
			u.SetRoles(roles)
			m := GetMessages(ctx.Context())
			r.PageTitle = m.Permissions
			r.Body = b.permissionsBody(m, b.Stack(pb, u), u)
			return
		})
	page.Build()
}

// permissionsBody is the stack s of the user u: its tabs.
func (b *Builder) permissionsBody(m *Messages, s *perm.Stack, u User) h.HTMLComponent {
	var note h.HTMLComponent
	if u.GetAccountName() == b.lb.GetInitialUserAccount() {
		note = v.VAlert(h.Text(m.PermissionsFixed)).Type("info").Density(v.DensityCompact).Class("mb-4")
	}

	// Tudo: every line
	var rows []h.HTMLComponent
	for _, l := range s.Lines {
		state := m.PermissionsActiveS
		class := "text-success"
		switch {
		case l.ReplacedBy > 0:
			state, class = fmt.Sprintf(m.PermissionsReplaced, l.ReplacedBy), "text-medium-emphasis"
		case l.Kind == perm.LineReset:
			state, class = m.PermissionsReset, "text-medium-emphasis"
		case len(l.ExceptBy) > 0:
			idx := make([]string, len(l.ExceptBy))
			for i, e := range l.ExceptBy {
				idx[i] = "#" + strconv.Itoa(e)
			}
			state = fmt.Sprintf(m.PermissionsExcept, strings.Join(idx, ", "))
		}
		rows = append(rows, h.Tr(
			h.Td(h.Text(strconv.Itoa(l.Index))),
			h.Td(h.Code(l.Text())),
			h.Td(h.Text(l.Origin)),
			h.Td(h.Text(state)).Class(class),
		).Attr("data-line", strconv.Itoa(l.Index)).Attr("data-replaced-by", strconv.Itoa(l.ReplacedBy)))
	}
	all := h.HTMLComponent(h.P(h.Text(m.PermissionsNone)).Class("text-medium-emphasis"))
	if len(rows) > 0 {
		all = v.VTable(
			h.Thead(h.Tr(h.Th(m.PermissionsIndex), h.Th(m.PermissionsLine), h.Th(m.PermissionsOrigin), h.Th(m.PermissionsState))),
			h.Tbody(rows...),
		).Density(v.DensityCompact).Attr("data-permissions-all", true)
	}

	// Ativas: the lines that count — an allow is a line of no prefix —, a
	// list or a tree
	var activeLines []string
	var items, excepts []h.HTMLComponent
	for _, l := range s.Active() {
		text := l.Pattern
		if l.Kind == perm.LineDeny {
			text = l.Text()
		}
		activeLines = append(activeLines, text)
		item := h.Div(h.Code(text)).Class("text-body-2")
		if len(l.Except) > 0 {
			// what it no longer counts for: the later lines narrower
			item.AppendChildren(h.Span(" "+fmt.Sprintf(m.PermissionsExceptOf, strings.Join(l.Except, ", "))).
				Class("text-medium-emphasis").Attr("data-except", true))
			excepts = append(excepts, h.Div(h.Code(text),
				h.Span(" "+fmt.Sprintf(m.PermissionsExceptOf, strings.Join(l.Except, ", "))).Class("text-medium-emphasis")).
				Class("text-body-2"))
		}
		items = append(items, item)
	}
	active := h.HTMLComponent(h.P(h.Text(m.PermissionsNone)).Class("text-medium-emphasis"))
	if len(items) > 0 {
		active = h.Div(
			v.VSwitch().Label(m.PermissionsTree).Attr("v-model", "locals.tree").Density(v.DensityCompact).
				Color("primary").HideDetails(true).Attr("data-permissions-tree-switch", true),
			h.Div(items...).Attr("v-if", "!locals.tree").Attr("data-permissions-active", true),
			h.Pre(perm.Tree(activeLines)).Attr("v-if", "locals.tree").
				Class("text-body-2 pa-2 bg-grey-lighten-4 rounded").Attr("data-permissions-tree", true),
			// the tree has the lines alone: their exceptions under it
			h.Div(excepts...).Attr("v-if", "locals.tree").Class("mt-2"),
		)
	}

	return h.Div(
		note,
		web.Scope(
			v.VTabs(
				v.VTab(h.Text(m.PermissionsAll)).Value("all"),
				v.VTab(h.Text(m.PermissionsActive)).Value("active"),
			).Attr("v-model", "locals.tab").Color("primary").Class("mb-4"),
			v.VTabsWindow(
				v.VTabsWindowItem(all).Value("all"),
				v.VTabsWindowItem(active).Value("active"),
			).Attr("v-model", "locals.tab"),
		).Slot("{ locals }").LocalsInit(`{tab: "all", tree: false}`),
	).Class("pa-4").Attr("data-user-permissions", true)
}
