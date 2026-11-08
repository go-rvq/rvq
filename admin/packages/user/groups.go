package user

import (
	"errors"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-rvq/rvq/admin/role"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// The groups of users: a group has members — each with the group's place
// among their groups (the order of their groups) — and roles in its order,
// as a user has. A user's stack is theirs, then their groups' (Stack): the
// lines of the database for them and for each group (a share with a group
// is its members'), the user's roles, then each group's roles, the groups
// in the user's order.

// Group is a group of users.
type Group struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"size:120;not null"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Group) TableName() string { return "user_groups" }

// String is the name of the group (its title in the admin).
func (g *Group) String() string { return g.Name }

func (g *Group) BeforeCreate(*gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

// GroupMember is a user in a group: Position the group's place among the
// user's groups.
type GroupMember struct {
	GroupID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID   uuid.UUID `gorm:"type:uuid;primaryKey;index"`
	Position int       `gorm:"not null;default:0"`
}

func (GroupMember) TableName() string { return "user_group_members" }

// GroupRole is a role of a group: Position its place in the group's order.
type GroupRole struct {
	GroupID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoleID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Position int       `gorm:"not null;default:0"`
}

func (GroupRole) TableName() string { return "user_group_roles" }

// GroupSubject is the subject of the group id in the policies: "group:<id>".
func GroupSubject(id uuid.UUID) string { return "group:" + id.String() }

// AutoMigrateGroups makes the tables of the groups and of the preferences.
func AutoMigrateGroups(db *gorm.DB) error {
	return db.AutoMigrate(&Group{}, &GroupMember{}, &GroupRole{}, &Preference{})
}

// groupGeneration is the version of the groups — their members, their
// roles —: in the version of the users' stacks.
var groupGeneration atomic.Int64

// GroupGeneration is the version of the groups.
func GroupGeneration() int64 { return groupGeneration.Load() }

// MemberHook is told of a user joining (added) or leaving a group, by by.
type MemberHook func(db *gorm.DB, group, user uuid.UUID, added bool, by string) error

var (
	hooksMu     sync.Mutex
	memberHooks []MemberHook
)

// OnMemberChange adds a hook told of the members of the groups changing
// (e.g. the history of the sharing).
func OnMemberChange(h MemberHook) {
	hooksMu.Lock()
	defer hooksMu.Unlock()
	memberHooks = append(memberHooks, h)
}

func memberChanged(db *gorm.DB, group, user uuid.UUID, added bool, by string) error {
	groupGeneration.Add(1)
	hooksMu.Lock()
	hs := slices.Clone(memberHooks)
	hooksMu.Unlock()
	for _, h := range hs {
		if err := h(db, group, user, added, by); err != nil {
			return err
		}
	}
	return nil
}

var ErrGroupName = errors.New("user: a group needs a name")

// CreateGroup makes a group.
func CreateGroup(db *gorm.DB, name, description string) (*Group, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrGroupName
	}
	g := &Group{Name: strings.TrimSpace(name), Description: description}
	return g, db.Create(g).Error
}

// AddMember puts the user in the group, the last of their groups; already
// there, nothing.
func AddMember(db *gorm.DB, group, user uuid.UUID, by string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var n int64
		if err := tx.Model(&GroupMember{}).Where("group_id = ? AND user_id = ?", group, user).Count(&n).Error; err != nil {
			return err
		}
		if n > 0 {
			return nil
		}
		var last struct{ P *int }
		if err := tx.Model(&GroupMember{}).Select("MAX(position) AS p").Where("user_id = ?", user).Scan(&last).Error; err != nil {
			return err
		}
		pos := 0
		if last.P != nil {
			pos = *last.P + 1
		}
		if err := tx.Create(&GroupMember{GroupID: group, UserID: user, Position: pos}).Error; err != nil {
			return err
		}
		return memberChanged(tx, group, user, true, by)
	})
}

// RemoveMember takes the user out of the group: what they had by it, no
// more.
func RemoveMember(db *gorm.DB, group, user uuid.UUID, by string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("group_id = ? AND user_id = ?", group, user).Delete(&GroupMember{})
		if res.Error != nil || res.RowsAffected == 0 {
			return res.Error
		}
		return memberChanged(tx, group, user, false, by)
	})
}

// DeleteGroup deletes the group: its members out of it (each told).
func DeleteGroup(db *gorm.DB, group uuid.UUID, by string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var users []uuid.UUID
		if err := tx.Model(&GroupMember{}).Where("group_id = ?", group).Pluck("user_id", &users).Error; err != nil {
			return err
		}
		for _, u := range users {
			if err := RemoveMember(tx, group, u, by); err != nil {
				return err
			}
		}
		if err := tx.Where("group_id = ?", group).Delete(&GroupRole{}).Error; err != nil {
			return err
		}
		groupGeneration.Add(1)
		return tx.Delete(&Group{}, "id = ?", group).Error
	})
}

// GroupsOf are the groups of the user, in their order.
func GroupsOf(db *gorm.DB, user uuid.UUID) (out []*Group, err error) {
	err = db.Joins("JOIN user_group_members m ON m.group_id = user_groups.id").
		Where("m.user_id = ?", user).Order("m.position, user_groups.name").Find(&out).Error
	return
}

// SetGroupPositions keeps the order of the groups of the user: each its
// place in ids.
func SetGroupPositions(db *gorm.DB, user uuid.UUID, ids []uuid.UUID) error {
	for i, g := range ids {
		if err := db.Model(&GroupMember{}).Where("user_id = ? AND group_id = ?", user, g).Update("position", i).Error; err != nil {
			return err
		}
	}
	groupGeneration.Add(1)
	return nil
}

// MembersOf are the users of the group.
func MembersOf(db *gorm.DB, group uuid.UUID) (ids []uuid.UUID, err error) {
	err = db.Model(&GroupMember{}).Where("group_id = ?", group).Pluck("user_id", &ids).Error
	return
}

// SetGroupRoles gives the group the roles ids, in their order.
func SetGroupRoles(db *gorm.DB, group uuid.UUID, ids []uuid.UUID) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", group).Delete(&GroupRole{}).Error; err != nil {
			return err
		}
		for i, r := range ids {
			if err := tx.Create(&GroupRole{GroupID: group, RoleID: r, Position: i}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err == nil {
		groupGeneration.Add(1)
	}
	return err
}

// GroupRolesOf are the roles of the group, in its order.
func GroupRolesOf(db *gorm.DB, group uuid.UUID) (roles []*role.Role, err error) {
	err = db.Joins("JOIN user_group_roles gr ON gr.role_id = roles.id").
		Where("gr.group_id = ?", group).Order("gr.position, roles.name").Find(&roles).Error
	return
}

// Preference is what a user chose for themselves: their default language
// (the admin's, the mails').
type Preference struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Language  string    `gorm:"size:20"`
	UpdatedAt time.Time
}

func (Preference) TableName() string { return "user_preferences" }

// LanguageOf is the default language of the user ("" none chosen).
func LanguageOf(db *gorm.DB, user uuid.UUID) string {
	var p Preference
	if db.Limit(1).Find(&p, "user_id = ?", user).RowsAffected == 0 {
		return ""
	}
	return p.Language
}

// SetLanguage keeps the default language of the user.
func SetLanguage(db *gorm.DB, user uuid.UUID, lang string) error {
	p := &Preference{UserID: user, Language: lang, UpdatedAt: time.Now()}
	return db.Save(p).Error
}

// Emailer is a user with an e-mail address (optional: a user without it is
// mailed at their account when it is an address).
type Emailer interface {
	GetEmail() string
}

// EmailOf is the e-mail address of the user: GetEmail, else their account
// when it is an address; "" none.
func EmailOf(u User) string {
	if u == nil {
		return ""
	}
	if e, ok := u.(Emailer); ok && e.GetEmail() != "" {
		return e.GetEmail()
	}
	if a := u.GetAccountName(); strings.Contains(a, "@") {
		return a
	}
	return ""
}
