package shared

import (
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// The sharing of records of any model (Service): a subject — a user or a
// group — given a role on a record (Access), inherited by the records under
// it (Resource.Parent) — the most specific wins, a role None takes a record
// out —; every grant and every change of the owner confirmed by the owner by
// e-mail (Request), a user's then accepted by them (Invite); all of it in
// the history (Event), kept forever.

// Subjects of an access: a user or a group.
const (
	SubjectUser  = "user:"
	SubjectGroup = "group:"
)

// UserSubject is the subject of the user id.
func UserSubject(id string) string { return SubjectUser + id }

// GroupSubject is the subject of the group id.
func GroupSubject(id string) string { return SubjectGroup + id }

// Access is a subject's role on a record of a resource.
type Access struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Resource  string    `gorm:"size:120;not null;uniqueIndex:uidx_shared_access"`
	RecordID  string    `gorm:"size:64;not null;uniqueIndex:uidx_shared_access"`
	Subject   string    `gorm:"size:128;not null;uniqueIndex:uidx_shared_access;index"`
	RoleID    uuid.UUID `gorm:"type:uuid;not null"`
	Role      *ShareTemplate
	GrantedBy string `gorm:"size:128"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Access) TableName() string { return "shared_access" }

func (a *Access) BeforeCreate(*gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// The events of the history.
const (
	EventConfirmationRequested = "confirmation_requested"
	EventOwnerConfirmed        = "owner_confirmed"
	EventOwnerRejected         = "owner_rejected"
	EventConfirmationExpired   = "confirmation_expired"
	EventInvited               = "invited"
	EventAccepted              = "accepted"
	EventRejected              = "rejected"
	EventGranted               = "granted"
	EventRoleChangeRequested   = "role_change_requested"
	EventRoleChanged           = "role_changed"
	EventRevoked               = "revoked"
	EventOwnerTransferRequest  = "owner_transfer_requested"
	EventOwnerTransferred      = "owner_transferred"
	EventMemberAdded           = "member_added"
	EventMemberRemoved         = "member_removed"
)

// Event is a line of the history of a record's sharing: what happened, to
// whom, the role and the actions before and after, by whom, when, from
// where.
type Event struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Resource      string         `gorm:"size:120;not null;index:idx_shared_events_record"`
	RecordID      string         `gorm:"size:64;not null;index:idx_shared_events_record"`
	Subject       string         `gorm:"size:128;index"`
	Event         string         `gorm:"size:40;not null"`
	RoleBefore    *uuid.UUID     `gorm:"type:uuid"`
	RoleAfter     *uuid.UUID     `gorm:"type:uuid"`
	ActionsBefore pq.StringArray `gorm:"type:text[]"`
	ActionsAfter  pq.StringArray `gorm:"type:text[]"`
	By            string         `gorm:"column:actor;size:128;index"`
	At            time.Time      `gorm:"not null;index"`
	IP            string         `gorm:"size:64"`
	UserAgent     string         `gorm:"size:400"`
	Note          string
}

func (Event) TableName() string { return "shared_events" }

func (e *Event) BeforeCreate(*gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// The kinds of a request to the owner.
const (
	RequestShare         = "share"
	RequestRoleChange    = "role_change"
	RequestOwnerTransfer = "owner_transfer"
)

// The states of a request and of an invite.
const (
	StatePending   = "pending"
	StateConfirmed = "confirmed"
	StateAccepted  = "accepted"
	StateRejected  = "rejected"
	StateExpired   = "expired"
)

// Request is what waits for the owner's confirmation by e-mail: a new
// sharing, a role with more actions, a new owner. Its token (sent by
// e-mail) is kept as a hash.
type Request struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Kind     string    `gorm:"size:20;not null"`
	Resource string    `gorm:"size:120;not null;index:idx_shared_requests_record"`
	RecordID string    `gorm:"size:64;not null;index:idx_shared_requests_record"`
	Subject  string    `gorm:"size:128;not null"`
	// RoleID is the role asked (a sharing, a change); of a new owner, the
	// role the former owner keeps.
	RoleID     *uuid.UUID `gorm:"type:uuid"`
	Owner      string     `gorm:"size:128;not null"`
	By         string     `gorm:"column:requested_by;size:128"`
	TokenHash  string     `gorm:"size:64;uniqueIndex"`
	ExpiresAt  time.Time
	State      string `gorm:"size:12;not null;index"`
	CreatedAt  time.Time
	AnsweredAt *time.Time
}

func (Request) TableName() string { return "shared_requests" }

func (r *Request) BeforeCreate(*gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// AccessInvite is a sharing confirmed by the owner, waiting for the user invited
// (or the new owner) to accept it.
type AccessInvite struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	RequestID  uuid.UUID  `gorm:"type:uuid;not null"`
	Kind       string     `gorm:"size:20;not null"`
	Resource   string     `gorm:"size:120;not null"`
	RecordID   string     `gorm:"size:64;not null"`
	Subject    string     `gorm:"size:128;not null;index"`
	RoleID     *uuid.UUID `gorm:"type:uuid"`
	By         string     `gorm:"column:invited_by;size:128"`
	State      string     `gorm:"size:12;not null;index"`
	CreatedAt  time.Time
	AnsweredAt *time.Time
}

func (AccessInvite) TableName() string { return "shared_access_invites" }

func (i *AccessInvite) BeforeCreate(*gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// AutoMigrate makes the tables of the sharing.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&ShareTemplate{}, &Access{}, &Event{}, &Request{}, &AccessInvite{})
}

// RoleNone is the name of the role of no access: of a record under one
// shared, it takes it out.
const RoleNone = "None"

// Allows says whether the role gives action.
func (t *ShareTemplate) Allows(action string) bool {
	return t != nil && !t.None && slices.Contains([]string(t.Actions), action)
}

// More says whether the role to gives an action the role from does not
// (from nil: any role gives more).
func More(from, to *ShareTemplate) bool {
	if to == nil || to.None {
		return false
	}
	if from == nil || from.None {
		return true
	}
	for _, a := range to.Actions {
		if !slices.Contains([]string(from.Actions), a) {
			return true
		}
	}
	return false
}

// RoleSeed is a role a resource offers.
type RoleSeed struct {
	Name    string
	Actions []string
}

// SeedRoles gives the resource of scope its roles (and None) when it has
// none (idempotent).
func SeedRoles(db *gorm.DB, scope uuid.UUID, resource string, roles ...RoleSeed) error {
	var n int64
	q := db.Model(&ShareTemplate{}).Where("resource = ?", resource)
	if scope != uuid.Nil {
		q = q.Where("scope = ?", scope)
	} else {
		q = q.Where("scope IS NULL")
	}
	if err := q.Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	for _, r := range append(roles, RoleSeed{Name: RoleNone}) {
		t := &ShareTemplate{Name: r.Name, Actions: r.Actions, Resource: resource, None: r.Name == RoleNone}
		if scope != uuid.Nil {
			t.Scope = scopeOf(scope)
		}
		if err := db.Create(t).Error; err != nil {
			return err
		}
	}
	return nil
}

// RolesOf are the roles of the resource in scope (its own; when it has none,
// the ones of any resource), None last.
func RolesOf(db *gorm.DB, scope uuid.UUID, resource string) ([]ShareTemplate, error) {
	q := db.Model(&ShareTemplate{}).Where("resource = ?", resource)
	var n int64
	if err := q.Session(&gorm.Session{}).Count(&n).Error; err != nil {
		return nil, err
	}
	if n == 0 {
		q = db.Model(&ShareTemplate{}).Where("resource = '' OR resource IS NULL")
	}
	if scope != uuid.Nil {
		q = q.Where("scope = ? OR scope IS NULL", scope)
	} else {
		q = q.Where("scope IS NULL")
	}
	var out []ShareTemplate
	if err := q.Order("none, name").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}
