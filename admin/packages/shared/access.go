package shared

import (
	"slices"
	"strings"
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
	// NoParentRead: of a record under another, the subject does not read its
	// parent and the parent's supporting records (by default they do: the
	// parent seen read only, its records that are not shared by themselves
	// — its other records shared so stay out)
	NoParentRead bool `gorm:"not null;default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	RoleID    *uuid.UUID `gorm:"type:uuid"`
	Owner     string     `gorm:"size:128;not null"`
	By        string     `gorm:"column:requested_by;size:128"`
	TokenHash string     `gorm:"size:64;uniqueIndex"`
	// NoParentRead is the Access's, asked
	NoParentRead bool `gorm:"not null;default:false"`
	ExpiresAt    time.Time
	State        string `gorm:"size:12;not null;index"`
	CreatedAt    time.Time
	AnsweredAt   *time.Time
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
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	RequestID uuid.UUID  `gorm:"type:uuid;not null"`
	Kind      string     `gorm:"size:20;not null"`
	Resource  string     `gorm:"size:120;not null"`
	RecordID  string     `gorm:"size:64;not null"`
	Subject   string     `gorm:"size:128;not null;index"`
	RoleID    *uuid.UUID `gorm:"type:uuid"`
	By        string     `gorm:"column:invited_by;size:128"`
	// NoParentRead is the Access's, confirmed
	NoParentRead bool   `gorm:"not null;default:false"`
	State        string `gorm:"size:12;not null;index"`
	CreatedAt    time.Time
	AnsweredAt   *time.Time
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

// Allows says whether the role gives action on the record (and all under
// it): "*" any.
func (t *ShareTemplate) Allows(action string) bool {
	if t == nil || t.None {
		return false
	}
	for _, r := range t.RuleList() {
		if r.On == RuleAll && (slices.Contains(r.Do, action) || slices.Contains(r.Do, "*")) {
			return true
		}
	}
	return false
}

// RuleAll is the resource of a rule of the record and all under it.
const RuleAll = "*"

// Rule is a part of a role: its actions (Do) on the resources under the
// record shared — On, after the record's: "*" (RuleAll) the record and all
// under it, "movimentacoes:*" the movements under it, "pagamentos:<*>:!efetivar"
// … (perm patterns).
type Rule struct {
	On string
	Do []string
}

func (r Rule) encode() string { return r.On + "|" + strings.Join(r.Do, ",") }

func decodeRule(s string) (Rule, bool) {
	on, do, ok := strings.Cut(s, "|")
	if !ok || on == "" || do == "" {
		return Rule{}, false
	}
	return Rule{On: on, Do: strings.Split(do, ",")}, true
}

// RuleList are the rules of the role: its Actions (on RuleAll), then its
// Rules.
func (t *ShareTemplate) RuleList() (out []Rule) {
	if t == nil {
		return nil
	}
	if len(t.Actions) > 0 {
		out = append(out, Rule{On: RuleAll, Do: []string(t.Actions)})
	}
	for _, l := range t.Rules {
		if r, ok := decodeRule(l); ok {
			out = append(out, r)
		}
	}
	return
}

// More says whether the role to gives an action — on a resource — the role
// from does not (from nil: any role gives more).
func More(from, to *ShareTemplate) bool {
	if to == nil || to.None {
		return false
	}
	if from == nil || from.None {
		return true
	}
	has := map[string]bool{}
	for _, r := range from.RuleList() {
		for _, a := range r.Do {
			has[r.On+" "+a] = true
		}
	}
	for _, r := range to.RuleList() {
		for _, a := range r.Do {
			if !has[r.On+" "+a] && !has[RuleAll+" "+a] && !has[r.On+" *"] && !has[RuleAll+" *"] {
				return true
			}
		}
	}
	return false
}

// RoleSeed is a role a resource offers: its Actions on the record and all
// under it, its Rules on resources under it.
type RoleSeed struct {
	Name    string
	Actions []string
	Rules   []Rule
}

func (r RoleSeed) rules() (out pq.StringArray) {
	for _, x := range r.Rules {
		out = append(out, x.encode())
	}
	return
}

// SeedRoles gives the resource of scope its roles (and None): the ones
// missing made, the ones there — by name — given the actions and the rules
// of the seed (the roles are the code's); roles of other names kept.
func SeedRoles(db *gorm.DB, scope uuid.UUID, resource string, roles ...RoleSeed) error {
	for _, r := range append(roles, RoleSeed{Name: RoleNone}) {
		var t ShareTemplate
		q := rolesQuery(db, scope, resource).Where("name = ?", r.Name)
		found := q.Limit(1).Find(&t).RowsAffected > 0
		if q.Error != nil {
			return q.Error
		}
		if found {
			if err := db.Model(&t).Updates(map[string]any{"actions": pq.StringArray(r.Actions), "rules": r.rules(),
				"none": r.Name == RoleNone}).Error; err != nil {
				return err
			}
			continue
		}
		t = ShareTemplate{Name: r.Name, Actions: r.Actions, Rules: r.rules(), Resource: resource, None: r.Name == RoleNone}
		if scope != uuid.Nil {
			t.Scope = scopeOf(scope)
		}
		if err := db.Create(&t).Error; err != nil {
			return err
		}
	}
	return nil
}

// RenameRoles renames the roles of the resource of scope (from → to) — the
// accesses keep them —, unless one of the new name is there.
func RenameRoles(db *gorm.DB, scope uuid.UUID, resource string, renames map[string]string) error {
	for from, to := range renames {
		var n int64
		if err := rolesQuery(db, scope, resource).Where("name = ?", to).Count(&n).Error; err != nil {
			return err
		}
		if n > 0 {
			continue
		}
		if err := rolesQuery(db, scope, resource).Where("name = ?", from).Update("name", to).Error; err != nil {
			return err
		}
	}
	return nil
}

// RemoveRoles removes the roles of the names of the resource of scope that
// no access, request or invite holds (one held stays).
func RemoveRoles(db *gorm.DB, scope uuid.UUID, resource string, names ...string) error {
	for _, name := range names {
		var t ShareTemplate
		if rolesQuery(db, scope, resource).Where("name = ?", name).Limit(1).Find(&t).RowsAffected == 0 {
			continue
		}
		var n, m, k int64
		db.Model(&Access{}).Where("role_id = ?", t.ID).Count(&n)
		db.Model(&Request{}).Where("role_id = ? AND state = ?", t.ID, StatePending).Count(&m)
		db.Model(&AccessInvite{}).Where("role_id = ? AND state = ?", t.ID, StatePending).Count(&k)
		if n+m+k > 0 {
			continue
		}
		if err := db.Delete(&t).Error; err != nil {
			return err
		}
	}
	return nil
}

// rolesQuery are the roles of the resource of scope.
func rolesQuery(db *gorm.DB, scope uuid.UUID, resource string) *gorm.DB {
	q := db.Model(&ShareTemplate{}).Where("resource = ?", resource)
	if scope != uuid.Nil {
		return q.Where("scope = ?", scope)
	}
	return q.Where("scope IS NULL")
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
