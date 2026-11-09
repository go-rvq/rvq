package shared

import (
	"context"

	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UsersOf are the users of the admin (the package user) to the sharing:
// found by the builder, their e-mail (user.EmailOf), their default language
// (user.LanguageOf), their groups in their order.
func UsersOf(b *user.Builder, db *gorm.DB) Users { return adminUsers{b, db} }

type adminUsers struct {
	b  *user.Builder
	db *gorm.DB
}

func (a adminUsers) find(ctx context.Context, id string) user.User {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil
	}
	u, err := a.b.Find(a.db.WithContext(ctx), uid)
	if err != nil {
		return nil
	}
	return u
}

func (a adminUsers) Exists(ctx context.Context, id string) bool { return a.find(ctx, id) != nil }

func (a adminUsers) Name(ctx context.Context, id string) string {
	if u := a.find(ctx, id); u != nil {
		return u.GetName()
	}
	return id
}

func (a adminUsers) Email(ctx context.Context, id string) string {
	return user.EmailOf(a.find(ctx, id))
}

func (a adminUsers) Language(ctx context.Context, id string) string {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ""
	}
	return user.LanguageOf(a.db.WithContext(ctx), uid)
}

func (a adminUsers) Groups(ctx context.Context, id string) (out []string) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil
	}
	gs, _ := user.GroupsOf(a.db.WithContext(ctx), uid)
	for _, g := range gs {
		out = append(out, g.ID.String())
	}
	return
}

func (a adminUsers) group(ctx context.Context, id string) *user.Group {
	var g user.Group
	if a.db.WithContext(ctx).Limit(1).Find(&g, "id = ?", id).RowsAffected == 0 {
		return nil
	}
	return &g
}

func (a adminUsers) Members(ctx context.Context, id string) (out []string) {
	gid, err := uuid.Parse(id)
	if err != nil {
		return nil
	}
	ids, _ := user.MembersOf(a.db.WithContext(ctx), gid)
	for _, u := range ids {
		out = append(out, u.String())
	}
	return
}

func (a adminUsers) GroupExists(ctx context.Context, id string) bool { return a.group(ctx, id) != nil }

func (a adminUsers) GroupName(ctx context.Context, id string) string {
	if g := a.group(ctx, id); g != nil {
		return g.Name
	}
	return id
}

// ListenGroups puts in the history of the records shared with a group its
// members joining and leaving (user.OnMemberChange).
func (s *Service) ListenGroups() {
	user.OnMemberChange(func(db *gorm.DB, group, u uuid.UUID, added bool, by string) error {
		return s.MemberChangedIn(db, group.String(), u.String(), added, by, Origin{})
	})
}
