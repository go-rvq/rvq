package shared

import (
	"errors"
	"time"

	"github.com/go-rvq/rvq/admin/packages/perms"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ShareStatus is the lifecycle of a share invite.
type ShareStatus string

const (
	StatusPending  ShareStatus = "pending"
	StatusAccepted ShareStatus = "accepted"
	StatusRejected ShareStatus = "rejected"
)

// ShareInvite is a pending/answered invitation to share a record with a user.
// Access is only granted (the perm rules are only written) once the invited
// user accepts — see Accept.
type ShareInvite struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	SharedID  uuid.UUID `gorm:"type:uuid;index"`
	Resource  string
	Subject   string         `gorm:"index"` // the invited user
	Actions   pq.StringArray `gorm:"type:text[]"`
	InvitedBy string
	Status    ShareStatus `gorm:"size:12;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ShareInvite) TableName() string { return "shared_invites" }

func (i *ShareInvite) BeforeCreate(*gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	if i.Status == "" {
		i.Status = StatusPending
	}
	return nil
}

// AutoMigrateInvites creates the share-invite table.
func AutoMigrateInvites(db *gorm.DB) error { return db.AutoMigrate(&ShareInvite{}) }

var (
	ErrInviteNotPending = errors.New("the share invite is not pending")
	ErrInviteNotYours   = errors.New("the share invite is addressed to another user")
)

// Notifier is the app's hook to surface share events in its UI (e.g. the rvq
// notification indicator). A nil Notifier means "do not notify".
type Notifier interface {
	// ShareInvited notifies the invited user (inv.Subject) of a pending invite,
	// and typically the inviter too.
	ShareInvited(db *gorm.DB, inv *ShareInvite) error
	// ShareAccepted notifies the inviter that inv.Subject accepted.
	ShareAccepted(db *gorm.DB, inv *ShareInvite) error
}

// Invite creates one pending invite per subject (all sharing one new SharedID)
// and notifies them. No permission is granted yet — the invited user must
// Accept. invitedBy is the inviting user's identifier.
func Invite(db *gorm.DB, resource string, subjects, actions []string, invitedBy string, n Notifier) ([]*ShareInvite, error) {
	sharedID := uuid.New()
	var invites []*ShareInvite
	for _, s := range subjects {
		inv := &ShareInvite{
			SharedID:  sharedID,
			Resource:  resource,
			Subject:   s,
			Actions:   actions,
			InvitedBy: invitedBy,
			Status:    StatusPending,
		}
		if err := db.Create(inv).Error; err != nil {
			return nil, err
		}
		if n != nil {
			if err := n.ShareInvited(db, inv); err != nil {
				return nil, err
			}
		}
		invites = append(invites, inv)
	}
	return invites, nil
}

// Accept grants the invited subject the actual access (GrantShared, stamped with
// the invite's SharedID) and marks the invite accepted, then notifies the
// inviter. Only a pending invite can be accepted.
func Accept(db *gorm.DB, inviteID uuid.UUID, n Notifier) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var inv ShareInvite
		if err := tx.First(&inv, "id = ?", inviteID).Error; err != nil {
			return err
		}
		if inv.Status != StatusPending {
			return ErrInviteNotPending
		}
		if _, err := perms.GrantShared(tx, ReferID(inv.SharedID, inv.Subject), inv.Subject, inv.Resource, []string(inv.Actions), inv.SharedID); err != nil {
			return err
		}
		if err := tx.Model(&inv).Update("status", StatusAccepted).Error; err != nil {
			return err
		}
		if n != nil {
			return n.ShareAccepted(tx, &inv)
		}
		return nil
	})
}

// AcceptFor accepts an invite only when it is addressed to subject, so a user
// can never accept someone else's invite. Use it from the UI with the current
// user as subject.
func AcceptFor(db *gorm.DB, inviteID uuid.UUID, subject string, n Notifier) error {
	var inv ShareInvite
	if err := db.Select("subject").First(&inv, "id = ?", inviteID).Error; err != nil {
		return err
	}
	if inv.Subject != subject {
		return ErrInviteNotYours
	}
	return Accept(db, inviteID, n)
}

// RejectFor rejects an invite only when it is addressed to subject.
func RejectFor(db *gorm.DB, inviteID uuid.UUID, subject string) error {
	var inv ShareInvite
	if err := db.Select("subject").First(&inv, "id = ?", inviteID).Error; err != nil {
		return err
	}
	if inv.Subject != subject {
		return ErrInviteNotYours
	}
	return Reject(db, inviteID)
}

// Reject marks a pending invite rejected (no access granted).
func Reject(db *gorm.DB, inviteID uuid.UUID) error {
	res := db.Model(&ShareInvite{}).
		Where("id = ? AND status = ?", inviteID, StatusPending).
		Update("status", StatusRejected)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrInviteNotPending
	}
	return nil
}

// PendingFor lists the pending invites addressed to subject (for the UI to show
// and let the user accept/reject).
func PendingFor(db *gorm.DB, subject string) ([]ShareInvite, error) {
	var out []ShareInvite
	err := db.Where("subject = ? AND status = ?", subject, StatusPending).
		Order("created_at desc").Find(&out).Error
	return out, err
}
