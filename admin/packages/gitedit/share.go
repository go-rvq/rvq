package gitedit

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DraftShare is a draft shared by its owner with a user: they reach it whole
// — the page, the IDE, the editor, git, the preview — until it expires or is
// revoked. Revoking keeps it: whoever had a draft, and when, is known.
type DraftShare struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	// Repo is the editor's repository (Builder.Name; "" the site's).
	Repo string `gorm:"column:repo;index;size:64;not null;default:''"`
	// Owner is the key of the owner of the draft.
	Owner string `gorm:"column:owner_key;index;size:128"`
	// User is the key of whom it is shared with.
	User string `gorm:"column:user_key;index;size:128"`
	// CreatedBy is the key of who shared it.
	CreatedBy string `gorm:"size:128"`
	// ExpiresAt, when set, is when it stops: reached no more from then.
	ExpiresAt *time.Time
	// RevokedAt is when it was revoked, and RevokedBy by whom.
	RevokedAt *time.Time
	RevokedBy string `gorm:"size:128"`
}

// TableName is the table of the sharings.
func (DraftShare) TableName() string { return "gitedit_draft_shares" }

// Active says whether it reaches the draft at now: not revoked, not expired.
func (s *DraftShare) Active(now time.Time) bool {
	return s.RevokedAt == nil && (s.ExpiresAt == nil || s.ExpiresAt.After(now))
}

// ErrShareSelf refuses a draft shared with its owner.
var ErrShareSelf = errors.New("gitedit: a draft is not shared with its owner")

// active is the query of the sharings active at now.
func active(db *gorm.DB, now time.Time) *gorm.DB {
	return db.Where("revoked_at IS NULL AND (expires_at IS NULL OR expires_at > ?)", now)
}

// sharedWith says whether the draft of owner is shared with user, now.
func (b *Builder) sharedWith(ctx context.Context, owner, user string) bool {
	if b.DB == nil {
		return false
	}
	var n int64
	active(b.DB.WithContext(ctx).Model(&DraftShare{}), time.Now()).
		Where("repo = ? AND owner_key = ? AND user_key = ?", b.Name, owner, user).Count(&n)
	return n > 0
}

// Share shares the draft of owner with user — by by — until expires (nil:
// until revoked). A sharing of them active is replaced: one each.
func (b *Builder) Share(ctx context.Context, owner, user, by string, expires *time.Time) (*DraftShare, error) {
	if owner == user {
		return nil, ErrShareSelf
	}
	if b.DB == nil {
		return nil, errors.New("gitedit: no database for the sharings")
	}
	s := &DraftShare{ID: uuid.New(), Repo: b.Name, Owner: owner, User: user, CreatedBy: by, ExpiresAt: expires}
	err := b.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		if err := active(tx.Model(&DraftShare{}), now).Where("repo = ? AND owner_key = ? AND user_key = ?", b.Name, owner, user).
			Updates(map[string]any{"revoked_at": now, "revoked_by": by}).Error; err != nil {
			return err
		}
		return tx.Create(s).Error
	})
	return s, err
}

// Revoke revokes the sharing id, by by.
func (b *Builder) Revoke(ctx context.Context, id uuid.UUID, by string) error {
	if b.DB == nil {
		return errors.New("gitedit: no database for the sharings")
	}
	now := time.Now()
	return b.DB.WithContext(ctx).Model(&DraftShare{}).Where("repo = ? AND id = ? AND revoked_at IS NULL", b.Name, id).
		Updates(map[string]any{"revoked_at": now, "revoked_by": by}).Error
}

// SharesOf are the sharings of the draft of owner — active and past —, the
// newest first.
func (b *Builder) SharesOf(ctx context.Context, owner string) (shares []*DraftShare, err error) {
	if b.DB == nil {
		return nil, nil
	}
	err = b.DB.WithContext(ctx).Where("repo = ? AND owner_key = ?", b.Name, owner).Order("created_at DESC").Find(&shares).Error
	return
}

// SharedWith are the sharings active with user: the drafts they reach.
func (b *Builder) SharedWith(ctx context.Context, user string) (shares []*DraftShare, err error) {
	if b.DB == nil {
		return nil, nil
	}
	err = active(b.DB.WithContext(ctx), time.Now()).Where("repo = ? AND user_key = ?", b.Name, user).Order("created_at DESC").
		Find(&shares).Error
	return
}

// share is the sharing id.
func (b *Builder) share(ctx context.Context, id string) (*DraftShare, error) {
	if b.DB == nil {
		return nil, gorm.ErrRecordNotFound
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	s := &DraftShare{}
	return s, b.DB.WithContext(ctx).First(s, "repo = ? AND id = ?", b.Name, uid).Error
}
