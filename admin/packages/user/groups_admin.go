package user

import (
	"fmt"
	"slices"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GroupsModelID is the id of the groups of users in the admin.
const GroupsModelID = "user_groups"

// Groups puts the groups of users in the admin of p: a group's name, its
// members and its roles in order — as a user's. A member added or taken out
// is told (OnMemberChange: what they reach by the group, at once); the roles
// change the stacks of the members at once.
func (c *Builder) Groups(p *presets.Builder) *presets.ModelBuilder {
	db := c.db
	mb := p.Model(&Group{}, presets.ModelConfig().SetId(GroupsModelID).SetModuleKey(MessagesKey)).
		URIName("user-groups").MenuIcon("mdi-account-group-outline")

	mb.Listing("Name", "Description", "Members").SearchColumns("name", "description")
	ed := mb.Editing("Name", "Description", "Members", "Roles")
	mb.Detailing("Name", "Description", "Members", "Roles", "CreatedAt", "UpdatedAt")

	by := func(ctx *web.EventContext) string {
		if u := GetCurrentUser(ctx.R); u != nil {
			return u.GetID().String()
		}
		return ""
	}
	idOf := func(obj any) uuid.UUID { return obj.(*Group).ID }

	ed.Validators.AppendFunc(func(obj interface{}, _ presets.FieldModeStack, ctx *web.EventContext) (err web.ValidationErrors) {
		if g := obj.(*Group); strings.TrimSpace(g.Name) == "" {
			err.FieldError("Name", GetMessages(ctx.Context()).ErrGroupName)
		}
		return
	})

	// the members: the users of the admin
	mb.Listing().Field("Members").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		ids, _ := MembersOf(db, idOf(field.Obj))
		return h.Td(h.Text(fmt.Sprint(len(ids))))
	})
	mb.Detailing().Field("Members").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return readonly(field, c.userNames(db, idOf(field.Obj)))
	})
	ed.Field("Members").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		var values []string
		if id := idOf(field.Obj); id != uuid.Nil {
			ids, _ := MembersOf(db, id)
			for _, m := range ids {
				values = append(values, m.String())
			}
		}
		return VAutocomplete().Label(field.Label).Multiple(true).Chips(true).ClosableChips(true).
			Items(c.userChoices(db)).ItemTitle("text").ItemValue("value").
			Attr(web.VField(field.FormKey, values)...).ErrorMessages(field.Errors...).
			Attr("data-group-members", true)
	}).SetterFunc(func(any, *presets.FieldContext, *web.EventContext) error { return nil })

	// the roles, in order — the order of the group's in the members' stacks
	mb.Detailing().Field("Roles").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		roles, _ := GroupRolesOf(db, idOf(field.Obj))
		var names []string
		for _, r := range roles {
			names = append(names, r.Name)
		}
		return readonly(field, names)
	})
	ed.Field("Roles").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		values := []string{}
		if id := idOf(field.Obj); id != uuid.Nil {
			roles, _ := GroupRolesOf(db, id)
			for _, r := range roles {
				values = append(values, r.ID.String())
			}
		}
		var all []*role.Role
		db.Order("name").Find(&all)
		items := []DefaultOptionItem{}
		for _, r := range all {
			items = append(items, DefaultOptionItem{Text: r.Name, Value: r.ID.String()})
		}
		field.ValueOverride = values
		return presets.WriteSortableSliceFieldComponentFuncWraper(func(field *presets.FieldContext, itemsModel, deleteFunc string) h.HTMLComponent {
			return VSelect(
				web.Slot(
					VIcon("mdi-delete").Color("error").Attr("@click", fmt.Sprintf(deleteFunc, "index")),
				).Name("append-inner"),
			).Items(items).ItemTitle("text").ItemValue("value").
				Density(DensityCompact).HideDetails(true).Class("mb-2").
				Attr("v-model", "element.__value").Attr("data-group-role-select", true)
		})(field, ctx)
	}).SetterFunc(func(any, *presets.FieldContext, *web.EventContext) error { return nil })

	// saved: the members and the roles as the form says
	ed.WrapSaveFunc(func(in presets.SaveFunc) presets.SaveFunc {
		return func(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
			if err = in(obj, id, ctx); err != nil {
				return
			}
			g := obj.(*Group)
			if err = setMembers(db, g.ID, parseIDs(ctx.FormSliceValues("Members")), by(ctx)); err != nil {
				return
			}
			roles := parseIDs(ctx.FormSliceValues("Roles"))
			if err = SetGroupRoles(db, g.ID, roles); err == nil {
				role.Changed(p.GetPermission())
			}
			return
		}
	})

	// deleted: its members out of it, each told
	mb.Listing().WrapDeleteFunc(func(in presets.DeleteFunc) presets.DeleteFunc {
		return func(obj interface{}, id model.ID, cascade bool, ctx *web.EventContext) error {
			gid, err := uuid.Parse(id.String())
			if err != nil {
				return err
			}
			return DeleteGroup(db, gid, by(ctx))
		}
	})
	return mb
}

// setMembers makes the members of the group want: the others out, the new
// ones in (each told).
func setMembers(db *gorm.DB, group uuid.UUID, want []uuid.UUID, by string) error {
	have, err := MembersOf(db, group)
	if err != nil {
		return err
	}
	for _, u := range have {
		if !slices.Contains(want, u) {
			if err := RemoveMember(db, group, u, by); err != nil {
				return err
			}
		}
	}
	for _, u := range want {
		if !slices.Contains(have, u) {
			if err := AddMember(db, group, u, by); err != nil {
				return err
			}
		}
	}
	return nil
}

// parseIDs are the ids of vs, each once, in order (the invalid left out).
func parseIDs(vs []string) (out []uuid.UUID) {
	for _, v := range vs {
		if id, err := uuid.Parse(v); err == nil && !slices.Contains(out, id) {
			out = append(out, id)
		}
	}
	return
}

// usersTable is the table of the users of c.
func (c *Builder) usersTable(db *gorm.DB) string {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(c.mb.NewModel()); err != nil {
		return "users"
	}
	return stmt.Schema.Table
}

// userChoices are the users to choose: their names (else accounts).
func (c *Builder) userChoices(db *gorm.DB) (items []DefaultOptionItem) {
	var rows []struct{ ID, Name, Account string }
	db.Table(c.usersTable(db)).Select("id, name, account").Where("deleted_at IS NULL").
		Order("name, account").Limit(1000).Scan(&rows)
	for _, r := range rows {
		title := r.Name
		if title == "" {
			title = r.Account
		}
		items = append(items, DefaultOptionItem{Text: title, Value: r.ID})
	}
	return
}

// userNames are the names of the members of the group.
func (c *Builder) userNames(db *gorm.DB, group uuid.UUID) (names []string) {
	ids, _ := MembersOf(db, group)
	if len(ids) == 0 {
		return nil
	}
	var rows []struct{ Name, Account string }
	db.Table(c.usersTable(db)).Select("name, account").Where("id IN ?", ids).Order("name, account").Scan(&rows)
	for _, r := range rows {
		if r.Name != "" {
			names = append(names, r.Name)
		} else {
			names = append(names, r.Account)
		}
	}
	return
}

// readonly is the field as read: its values, comma separated.
func readonly(field *presets.FieldContext, values []string) h.HTMLComponent {
	field.ValueOverride = strings.Join(values, ", ")
	return presets.ReadonlyComponentFunc(field, nil)
}
