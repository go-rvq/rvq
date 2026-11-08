package shared

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gad-lang/gad"
	mail_sender "github.com/go-rvq/rvq/admin/packages/mail-sender"
	"github.com/go-rvq/rvq/x/mailtpl"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Resource is a model whose records are shared: its table, the column of
// its id, of its owner ("" none: the owner is its parent's) and its parent
// (the records it inherits its access from).
type Resource struct {
	// Name is the resource's (the model's id).
	Name string
	// Table and IDColumn ("id") are where its records are.
	Table, IDColumn string
	// OwnerColumn is the column of the owner's id; empty: the owner is the
	// parent's.
	OwnerColumn string
	// Parent is the resource it inherits from and the column of its id.
	Parent       string
	ParentColumn string
	// Label is the name of a record (in the mails); nil: its id.
	Label func(ctx context.Context, db *gorm.DB, recordID string) string
}

// Users are the users of the application, to the sharing.
type Users interface {
	// Exists says whether the user id is one (a sharing is with an existing
	// user only).
	Exists(ctx context.Context, id string) bool
	// Name is the user's name; Email their address (the owner confirms by
	// it); Language their default language (the mails' language).
	Name(ctx context.Context, id string) string
	Email(ctx context.Context, id string) string
	Language(ctx context.Context, id string) string
	// Groups are the ids of the groups of the user, in the user's order.
	Groups(ctx context.Context, id string) []string
	// GroupExists says whether the group id is one; GroupName its name.
	GroupExists(ctx context.Context, id string) bool
	GroupName(ctx context.Context, id string) string
}

// Options of a Service.
type Options struct {
	DB    *gorm.DB
	Users Users
	// Sender sends the mails (the owner's confirmation, the invite).
	Sender mail_sender.Sender
	// ConfirmURL is the link of the page where the owner confirms a request
	// of token.
	ConfirmURL func(token string) string
	// TTL is how long the owner has to confirm (48h).
	TTL time.Duration
	// Now is the clock (time.Now).
	Now func() time.Time
}

// Service is the sharing of the records of the resources registered.
type Service struct {
	o         Options
	mu        sync.RWMutex
	resources map[string]*Resource
	// lastAt keeps the times of the events increasing: their order
	lastAt time.Time
	hooks  map[string][]AccessHook
}

// AccessHook is told of an access given or changed (role its role) or taken
// away (removed), in the transaction tx — e.g. the policies of the
// permissions of the admin kept with it.
type AccessHook func(tx *gorm.DB, a *Access, role *ShareTemplate, removed bool) error

// OnAccess adds a hook of the accesses of the resource.
func (s *Service) OnAccess(resource string, h AccessHook) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.hooks == nil {
		s.hooks = map[string][]AccessHook{}
	}
	s.hooks[resource] = append(s.hooks[resource], h)
}

// accessChanged tells the hooks of the resource of a.
func (s *Service) accessChanged(tx *gorm.DB, a *Access, role *ShareTemplate, removed bool) error {
	s.mu.RLock()
	hs := slices.Clone(s.hooks[a.Resource])
	s.mu.RUnlock()
	for _, h := range hs {
		if err := h(tx, a, role, removed); err != nil {
			return err
		}
	}
	return nil
}

var (
	ErrUnknownResource = errors.New("shared: the resource is not registered")
	ErrNoRecord        = errors.New("shared: no such record")
	ErrNoOwner         = errors.New("shared: the record has no owner to confirm")
	ErrNoOwnerEmail    = errors.New("shared: the owner has no e-mail to confirm by")
	ErrNoSuchSubject   = errors.New("shared: no such user or group")
	ErrRole            = errors.New("shared: the role is not of the resource")
	ErrAlreadyShared   = errors.New("shared: the record is shared with the subject already: change their role")
	ErrNotShared       = errors.New("shared: the record is not shared with the subject")
	ErrNotPending      = errors.New("shared: it is not pending")
	ErrExpired         = errors.New("shared: the confirmation expired")
	ErrNotOwner        = errors.New("shared: only the owner confirms")
	ErrNotYours        = errors.New("shared: the invite is another's")
	ErrOwnerIsUser     = errors.New("shared: the new owner must be another user")
)

// New is a service on the database of o, its tables made.
func New(o Options) (*Service, error) {
	if o.DB == nil || o.Users == nil {
		return nil, errors.New("shared: a database and the users are needed")
	}
	if o.TTL == 0 {
		o.TTL = 48 * time.Hour
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	if err := AutoMigrate(o.DB); err != nil {
		return nil, err
	}
	return &Service{o: o, resources: map[string]*Resource{}}, nil
}

// Register makes a resource shared.
func (s *Service) Register(r Resource) {
	if r.IDColumn == "" {
		r.IDColumn = "id"
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resources[r.Name] = &r
}

func (s *Service) resource(name string) (*Resource, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if r := s.resources[name]; r != nil {
		return r, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrUnknownResource, name)
}

// --- the record, its parents, its owner ---

type level struct {
	res *Resource
	id  string
}

// chain is the record and the ones above it, the record first.
func (s *Service) chain(ctx context.Context, resource, id string) ([]level, error) {
	var out []level
	for resource != "" {
		r, err := s.resource(resource)
		if err != nil {
			return nil, err
		}
		out = append(out, level{r, id})
		if r.Parent == "" {
			break
		}
		var parents []*string
		if err := s.o.DB.WithContext(ctx).Table(r.Table).Where(r.IDColumn+" = ?", id).
			Pluck(r.ParentColumn, &parents).Error; err != nil {
			return nil, err
		}
		if len(parents) == 0 || parents[0] == nil || *parents[0] == "" {
			break
		}
		resource, id = r.Parent, *parents[0]
	}
	return out, nil
}

// ownerAt is the owner of the record of l ("" none).
func (s *Service) ownerAt(ctx context.Context, l level) (string, error) {
	return ownerIn(s.o.DB.WithContext(ctx), l)
}

// ownerIn is the owner of the record of l in db ("" none).
func ownerIn(db *gorm.DB, l level) (string, error) {
	if l.res.OwnerColumn == "" {
		return "", nil
	}
	var owners []*string
	if err := db.Table(l.res.Table).Where(l.res.IDColumn+" = ?", l.id).
		Pluck(l.res.OwnerColumn, &owners).Error; err != nil {
		return "", err
	}
	if len(owners) == 0 {
		return "", ErrNoRecord
	}
	if owners[0] == nil {
		return "", nil
	}
	return *owners[0], nil
}

// Owner is the owner of the record: its own, else the nearest above it.
func (s *Service) Owner(ctx context.Context, resource, id string) (string, error) {
	ch, err := s.chain(ctx, resource, id)
	if err != nil {
		return "", err
	}
	if len(ch) > 0 {
		if _, err := s.ownerAt(ctx, ch[0]); errors.Is(err, ErrNoRecord) {
			return "", err
		}
	}
	for _, l := range ch {
		o, err := s.ownerAt(ctx, l)
		if err != nil {
			return "", err
		}
		if o != "" {
			return o, nil
		}
	}
	return "", ErrNoOwner
}

// --- the role of a user ---

// Effective is the role of a user on a record: Owner (every action), or the
// role of the access that decides — the most specific record's, of the user
// themselves first, then of their groups in their order —; Via is the
// record of that access and Subject its subject. No Role, no access.
type Effective struct {
	Owner    bool
	Role     *ShareTemplate
	Resource string
	RecordID string
	Subject  string
}

// Allows says whether the role gives action.
func (e *Effective) Allows(action string) bool {
	return e != nil && (e.Owner || e.Role.Allows(action))
}

// Visible says whether the user sees the record.
func (e *Effective) Visible() bool { return e != nil && (e.Owner || (e.Role != nil && !e.Role.None)) }

// subjects are the user's subjects, in their order: the user, their groups.
func (s *Service) subjects(ctx context.Context, user string) []string {
	out := []string{UserSubject(user)}
	for _, g := range s.o.Users.Groups(ctx, user) {
		out = append(out, GroupSubject(g))
	}
	return out
}

// RoleOf is the role of user on the record.
func (s *Service) RoleOf(ctx context.Context, user, resource, id string) (*Effective, error) {
	ch, err := s.chain(ctx, resource, id)
	if err != nil {
		return nil, err
	}
	for _, l := range ch {
		if o, err := s.ownerAt(ctx, l); err != nil {
			return nil, err
		} else if o != "" && o == user {
			return &Effective{Owner: true, Resource: l.res.Name, RecordID: l.id}, nil
		}
	}
	subs := s.subjects(ctx, user)
	for _, l := range ch {
		var acc []Access
		if err := s.o.DB.WithContext(ctx).Preload("Role").
			Where("resource = ? AND record_id = ? AND subject IN ?", l.res.Name, l.id, subs).Find(&acc).Error; err != nil {
			return nil, err
		}
		for _, sub := range subs {
			for _, a := range acc {
				if a.Subject == sub {
					return &Effective{Role: a.Role, Resource: l.res.Name, RecordID: l.id, Subject: sub}, nil
				}
			}
		}
	}
	return &Effective{}, nil
}

// Can says whether user may do action on the record.
func (s *Service) Can(ctx context.Context, user, resource, id, action string) bool {
	e, err := s.RoleOf(ctx, user, resource, id)
	return err == nil && e.Allows(action)
}

// AccessibleIDs are the records of the resource user sees: owned, shared
// with them or their groups, or under a record they see — the most
// specific access deciding (None takes a record out).
func (s *Service) AccessibleIDs(ctx context.Context, user, resource string) ([]string, error) {
	r, err := s.resource(resource)
	if err != nil {
		return nil, err
	}
	db := s.o.DB.WithContext(ctx)
	subs := s.subjects(ctx, user)
	// the accesses of the records of the resource: of each, the one deciding
	var acc []Access
	if err := db.Preload("Role").Where("resource = ? AND subject IN ?", r.Name, subs).Find(&acc).Error; err != nil {
		return nil, err
	}
	decide := map[string]*Access{}
	for _, sub := range subs {
		for i := range acc {
			if a := &acc[i]; a.Subject == sub && decide[a.RecordID] == nil {
				decide[a.RecordID] = a
			}
		}
	}
	seen := map[string]bool{}
	var out []string
	add := func(id string) {
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	// owned
	if r.OwnerColumn != "" {
		var ids []string
		if err := db.Table(r.Table).Where(r.OwnerColumn+" = ?", user).Pluck(r.IDColumn, &ids).Error; err != nil {
			return nil, err
		}
		for _, id := range ids {
			add(id)
		}
	}
	// shared directly
	for id, a := range decide {
		if a.Role != nil && !a.Role.None {
			add(id)
		}
	}
	// under a record seen: the parent's, but the ones taken out (None);
	// under a record owned, every one
	if r.Parent != "" {
		parents, err := s.AccessibleIDs(ctx, user, r.Parent)
		if err != nil {
			return nil, err
		}
		owned, err := s.ownedIDs(ctx, user, r.Parent)
		if err != nil {
			return nil, err
		}
		if len(parents) > 0 {
			var rows []struct{ ID, Parent string }
			if err := db.Table(r.Table).Select(r.IDColumn+" AS id, "+r.ParentColumn+" AS parent").
				Where(r.ParentColumn+" IN ?", parents).Scan(&rows).Error; err != nil {
				return nil, err
			}
			for _, row := range rows {
				if a := decide[row.ID]; a != nil && a.Role != nil && a.Role.None && !slices.Contains(owned, row.Parent) {
					continue
				}
				add(row.ID)
			}
		}
	}
	return out, nil
}

// ownedIDs are the records of the resource owned by user, the ones under
// them owned too.
func (s *Service) ownedIDs(ctx context.Context, user, resource string) ([]string, error) {
	r, err := s.resource(resource)
	if err != nil {
		return nil, err
	}
	var ids []string
	if r.OwnerColumn != "" {
		if err := s.o.DB.WithContext(ctx).Table(r.Table).Where(r.OwnerColumn+" = ?", user).
			Pluck(r.IDColumn, &ids).Error; err != nil {
			return nil, err
		}
	}
	if r.Parent != "" {
		parents, err := s.ownedIDs(ctx, user, r.Parent)
		if err != nil {
			return nil, err
		}
		if len(parents) > 0 {
			var more []string
			if err := s.o.DB.WithContext(ctx).Table(r.Table).Where(r.ParentColumn+" IN ?", parents).
				Pluck(r.IDColumn, &more).Error; err != nil {
				return nil, err
			}
			ids = append(ids, more...)
		}
	}
	return ids, nil
}

// Scope is db limited to the records of the resource user sees.
func (s *Service) Scope(ctx context.Context, db *gorm.DB, user, resource string) (*gorm.DB, error) {
	r, err := s.resource(resource)
	if err != nil {
		return nil, err
	}
	ids, err := s.AccessibleIDs(ctx, user, resource)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return db.Where("1 = 0"), nil
	}
	return db.Where(r.Table+"."+r.IDColumn+" IN ?", ids), nil
}

// --- who has access ---

// Subjects are the accesses of the record, and the ones of the records
// above it (inherited), the record's first.
func (s *Service) Subjects(ctx context.Context, resource, id string) ([]Access, error) {
	ch, err := s.chain(ctx, resource, id)
	if err != nil {
		return nil, err
	}
	var out []Access
	for _, l := range ch {
		var acc []Access
		if err := s.o.DB.WithContext(ctx).Preload("Role").Where("resource = ? AND record_id = ?", l.res.Name, l.id).
			Order("created_at").Find(&acc).Error; err != nil {
			return nil, err
		}
		out = append(out, acc...)
	}
	return out, nil
}

// History is the history of the record, the oldest first.
func (s *Service) History(ctx context.Context, resource, id string) ([]Event, error) {
	var out []Event
	return out, s.o.DB.WithContext(ctx).Where("resource = ? AND record_id = ?", resource, id).
		Order("at, id").Find(&out).Error
}

// HistoryOf is the history of the subject (a user, a group), the newest
// first.
func (s *Service) HistoryOf(ctx context.Context, subject string) ([]Event, error) {
	var out []Event
	return out, s.o.DB.WithContext(ctx).Where("subject = ? OR actor = ?", subject, strings.TrimPrefix(subject, SubjectUser)).
		Order("at DESC").Find(&out).Error
}

// --- the flows ---

// Origin is where an action came from (the history's IP and browser).
type Origin struct{ IP, UserAgent string }

func (s *Service) event(tx *gorm.DB, e Event, o Origin) error {
	s.mu.Lock()
	at := s.o.Now()
	if !at.After(s.lastAt) {
		at = s.lastAt.Add(time.Microsecond)
	}
	s.lastAt = at
	s.mu.Unlock()
	e.At = at
	e.IP, e.UserAgent = o.IP, o.UserAgent
	return tx.Create(&e).Error
}

func rolePtr(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

func (s *Service) role(tx *gorm.DB, resource string, id uuid.UUID) (*ShareTemplate, error) {
	var t ShareTemplate
	if err := tx.First(&t, "id = ?", id).Error; err != nil {
		return nil, ErrRole
	}
	if t.Resource != "" && t.Resource != resource {
		return nil, ErrRole
	}
	return &t, nil
}

// subjectExists says whether the subject is a user or a group there is.
func (s *Service) subjectExists(ctx context.Context, subject string) bool {
	switch {
	case strings.HasPrefix(subject, SubjectUser):
		return s.o.Users.Exists(ctx, strings.TrimPrefix(subject, SubjectUser))
	case strings.HasPrefix(subject, SubjectGroup):
		return s.o.Users.GroupExists(ctx, strings.TrimPrefix(subject, SubjectGroup))
	}
	return false
}

// subjectName is the name of the subject.
func (s *Service) subjectName(ctx context.Context, subject string) string {
	if g, ok := strings.CutPrefix(subject, SubjectGroup); ok {
		return s.o.Users.GroupName(ctx, g)
	}
	return s.o.Users.Name(ctx, strings.TrimPrefix(subject, SubjectUser))
}

// newToken is a token and its hash.
func newToken() (token, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}
	token = base64.RawURLEncoding.EncodeToString(b)
	return token, tokenHash(token), nil
}

func tokenHash(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// Share asks the owner of the record — by e-mail — to share it with the
// subject (a user or a group) with the role, by by; the role None, at once
// (nil Request).
func (s *Service) Share(ctx context.Context, resource, id, subject string, roleID uuid.UUID, by string, o Origin) (*Request, error) {
	if !s.subjectExists(ctx, subject) {
		return nil, ErrNoSuchSubject
	}
	var n int64
	if err := s.o.DB.WithContext(ctx).Model(&Access{}).
		Where("resource = ? AND record_id = ? AND subject = ?", resource, id, subject).Count(&n).Error; err != nil {
		return nil, err
	}
	if n > 0 {
		return nil, ErrAlreadyShared
	}
	// None gives nothing — it takes the record out of what is shared above
	// it —: at once, no owner to ask, no invite to accept (one invited to
	// lose a record would decline it)
	role, err := s.role(s.o.DB.WithContext(ctx), resource, roleID)
	if err != nil {
		return nil, err
	}
	if role.None {
		return nil, s.o.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return s.grant(tx, resource, id, subject, roleID, by, o)
		})
	}
	return s.request(ctx, RequestShare, resource, id, subject, roleID, by, o)
}

// ChangeRole changes the role of the subject on the record: to one with an
// action the old one has not, asking the owner (a Request); to one with
// fewer, at once (nil Request).
func (s *Service) ChangeRole(ctx context.Context, resource, id, subject string, roleID uuid.UUID, by string, o Origin) (*Request, error) {
	db := s.o.DB.WithContext(ctx)
	var a Access
	if err := db.Preload("Role").Where("resource = ? AND record_id = ? AND subject = ?", resource, id, subject).
		First(&a).Error; err != nil {
		return nil, ErrNotShared
	}
	to, err := s.role(db, resource, roleID)
	if err != nil {
		return nil, err
	}
	if More(a.Role, to) {
		return s.request(ctx, RequestRoleChange, resource, id, subject, roleID, by, o)
	}
	return nil, db.Transaction(func(tx *gorm.DB) error { return s.setRole(tx, &a, to, by, o) })
}

func (s *Service) setRole(tx *gorm.DB, a *Access, to *ShareTemplate, by string, o Origin) error {
	before, actsBefore := a.RoleID, a.Role.Actions
	if err := tx.Model(&Access{}).Where("id = ?", a.ID).Updates(map[string]any{"role_id": to.ID, "updated_at": s.o.Now()}).Error; err != nil {
		return err
	}
	a.RoleID = to.ID
	if err := s.accessChanged(tx, a, to, false); err != nil {
		return err
	}
	return s.event(tx, Event{Resource: a.Resource, RecordID: a.RecordID, Subject: a.Subject, Event: EventRoleChanged,
		RoleBefore: rolePtr(before), RoleAfter: rolePtr(to.ID), ActionsBefore: actsBefore, ActionsAfter: to.Actions, By: by}, o)
}

// Revoke takes the subject's access to the record away, at once.
func (s *Service) Revoke(ctx context.Context, resource, id, subject, by string, o Origin) error {
	return s.o.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var a Access
		if err := tx.Preload("Role").Where("resource = ? AND record_id = ? AND subject = ?", resource, id, subject).
			First(&a).Error; err != nil {
			return ErrNotShared
		}
		if err := tx.Delete(&a).Error; err != nil {
			return err
		}
		if err := s.accessChanged(tx, &a, a.Role, true); err != nil {
			return err
		}
		return s.event(tx, Event{Resource: resource, RecordID: id, Subject: subject, Event: EventRevoked,
			RoleBefore: rolePtr(a.RoleID), ActionsBefore: a.Role.Actions, By: by}, o)
	})
}

// TransferOwner asks the owner — by e-mail — to make the user the record's
// owner; the former owner keeps the role keep.
func (s *Service) TransferOwner(ctx context.Context, resource, id, user string, keep uuid.UUID, by string, o Origin) (*Request, error) {
	owner, err := s.Owner(ctx, resource, id)
	if err != nil {
		return nil, err
	}
	if user == owner || !s.o.Users.Exists(ctx, user) {
		return nil, ErrOwnerIsUser
	}
	return s.request(ctx, RequestOwnerTransfer, resource, id, UserSubject(user), keep, by, o)
}

// request makes the request and mails the owner the link to confirm it.
func (s *Service) request(ctx context.Context, kind, resource, id, subject string, roleID uuid.UUID, by string, o Origin) (*Request, error) {
	owner, err := s.Owner(ctx, resource, id)
	if err != nil {
		return nil, err
	}
	email := s.o.Users.Email(ctx, owner)
	if email == "" {
		return nil, ErrNoOwnerEmail
	}
	db := s.o.DB.WithContext(ctx)
	var role *ShareTemplate
	if roleID != uuid.Nil {
		if role, err = s.role(db, resource, roleID); err != nil {
			return nil, err
		}
	}
	token, hash, err := newToken()
	if err != nil {
		return nil, err
	}
	req := &Request{Kind: kind, Resource: resource, RecordID: id, Subject: subject, RoleID: rolePtr(roleID),
		Owner: owner, By: by, TokenHash: hash, ExpiresAt: s.o.Now().Add(s.o.TTL), State: StatePending, CreatedAt: s.o.Now()}
	ev := map[string]string{RequestShare: EventConfirmationRequested, RequestRoleChange: EventRoleChangeRequested,
		RequestOwnerTransfer: EventOwnerTransferRequest}[kind]
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(req).Error; err != nil {
			return err
		}
		e := Event{Resource: resource, RecordID: id, Subject: subject, Event: ev, By: by, Note: req.ID.String()}
		if role != nil {
			e.RoleAfter, e.ActionsAfter = &role.ID, role.Actions
		}
		return s.event(tx, e, o)
	})
	if err != nil {
		return nil, err
	}
	data := s.mailData(ctx, req, role)
	if s.o.ConfirmURL != nil {
		data["ConfirmURL"] = gad.Str(s.o.ConfirmURL(token))
	}
	if err := s.mail(ctx, "shared/confirm", owner, email, data); err != nil {
		return req, err
	}
	return req, nil
}

// mailData are the data of the mails of the request.
func (s *Service) mailData(ctx context.Context, req *Request, role *ShareTemplate) gad.Dict {
	label := req.RecordID
	if r, err := s.resource(req.Resource); err == nil && r.Label != nil {
		label = r.Label(ctx, s.o.DB.WithContext(ctx), req.RecordID)
	}
	d := gad.Dict{
		"Kind":       gad.Str(req.Kind),
		"Resource":   gad.Str(req.Resource),
		"Record":     gad.Str(label),
		"By":         gad.Str(s.o.Users.Name(ctx, req.By)),
		"Subject":    gad.Str(s.subjectName(ctx, req.Subject)),
		"Group":      gad.Bool(strings.HasPrefix(req.Subject, SubjectGroup)),
		"Role":       gad.Str(""),
		"Expires":    gad.Str(req.ExpiresAt.Format("2006-01-02 15:04 MST")),
		"ConfirmURL": gad.Str(""),
	}
	if role != nil {
		d["Role"] = gad.Str(role.Name)
	}
	return d
}

// mail sends the template name to user (at email), in their language.
func (s *Service) mail(ctx context.Context, name, user, email string, data gad.Dict) error {
	if s.o.Sender == nil {
		return errors.New("shared: no sender of mails")
	}
	lang := s.o.Users.Language(ctx, user)
	m, err := mailtpl.Render(name, lang, data)
	if err != nil {
		return err
	}
	return s.o.Sender.Send(mail_sender.Message().To(email).Subject(m.Subject).Body(m.HTML).
		Kind("shared").Locale(m.Language).User(user))
}

// RequestOf is the request of token.
func (s *Service) RequestOf(ctx context.Context, token string) (*Request, error) {
	var r Request
	if err := s.o.DB.WithContext(ctx).First(&r, "token_hash = ?", tokenHash(token)).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// pending is the request of token pending, by its owner; expired, marked
// so.
func (s *Service) pending(tx *gorm.DB, token, by string, o Origin) (*Request, error) {
	var r Request
	if err := tx.First(&r, "token_hash = ?", tokenHash(token)).Error; err != nil {
		return nil, ErrNotPending
	}
	if r.State != StatePending {
		return nil, ErrNotPending
	}
	if by != r.Owner {
		return nil, ErrNotOwner
	}
	if s.o.Now().After(r.ExpiresAt) {
		now := s.o.Now()
		tx.Model(&r).Updates(map[string]any{"state": StateExpired, "answered_at": &now})
		_ = s.event(tx, Event{Resource: r.Resource, RecordID: r.RecordID, Subject: r.Subject, Event: EventConfirmationExpired,
			By: by, Note: r.ID.String()}, o)
		return nil, ErrExpired
	}
	return &r, nil
}

// Confirm is the owner (by) confirming the request of token: a sharing with
// a user, an invite to them; with a group, the access at once; a role
// change, at once; a new owner, an invite to them.
func (s *Service) Confirm(ctx context.Context, token, by string, o Origin) error {
	var expired bool
	err := s.o.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		r, err := s.pending(tx, token, by, o)
		if errors.Is(err, ErrExpired) {
			expired = true
			return nil
		}
		if err != nil {
			return err
		}
		now := s.o.Now()
		if err := tx.Model(r).Updates(map[string]any{"state": StateConfirmed, "answered_at": &now}).Error; err != nil {
			return err
		}
		if err := s.event(tx, Event{Resource: r.Resource, RecordID: r.RecordID, Subject: r.Subject, Event: EventOwnerConfirmed,
			RoleAfter: r.RoleID, By: by, Note: r.ID.String()}, o); err != nil {
			return err
		}
		switch {
		case r.Kind == RequestRoleChange:
			var a Access
			if err := tx.Preload("Role").Where("resource = ? AND record_id = ? AND subject = ?", r.Resource, r.RecordID, r.Subject).
				First(&a).Error; err != nil {
				return ErrNotShared
			}
			to, err := s.role(tx, r.Resource, *r.RoleID)
			if err != nil {
				return err
			}
			return s.setRole(tx, &a, to, by, o)
		case r.Kind == RequestShare && (strings.HasPrefix(r.Subject, SubjectGroup) || s.isNone(tx, r.Resource, r.RoleID)):
			// a group's, at once for its members; None, no invite to accept
			return s.grant(tx, r.Resource, r.RecordID, r.Subject, *r.RoleID, r.By, o)
		default:
			inv := &AccessInvite{RequestID: r.ID, Kind: r.Kind, Resource: r.Resource, RecordID: r.RecordID, Subject: r.Subject,
				RoleID: r.RoleID, By: r.By, State: StatePending, CreatedAt: now}
			if err := tx.Create(inv).Error; err != nil {
				return err
			}
			return s.event(tx, Event{Resource: r.Resource, RecordID: r.RecordID, Subject: r.Subject, Event: EventInvited,
				RoleAfter: r.RoleID, By: r.By, Note: inv.ID.String()}, o)
		}
	})
	if err == nil && expired {
		return ErrExpired
	}
	return err
}

// Reject is the owner (by) refusing the request of token.
func (s *Service) Reject(ctx context.Context, token, by string, o Origin) error {
	var expired bool
	err := s.o.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		r, err := s.pending(tx, token, by, o)
		if errors.Is(err, ErrExpired) {
			expired = true
			return nil
		}
		if err != nil {
			return err
		}
		now := s.o.Now()
		if err := tx.Model(r).Updates(map[string]any{"state": StateRejected, "answered_at": &now}).Error; err != nil {
			return err
		}
		return s.event(tx, Event{Resource: r.Resource, RecordID: r.RecordID, Subject: r.Subject, Event: EventOwnerRejected,
			By: by, Note: r.ID.String()}, o)
	})
	if err == nil && expired {
		return ErrExpired
	}
	return err
}

// isNone says the role id is None.
func (s *Service) isNone(tx *gorm.DB, resource string, id *uuid.UUID) bool {
	if id == nil {
		return false
	}
	role, err := s.role(tx, resource, *id)
	return err == nil && role.None
}

// grant gives the subject the role on the record.
func (s *Service) grant(tx *gorm.DB, resource, id, subject string, roleID uuid.UUID, by string, o Origin) error {
	role, err := s.role(tx, resource, roleID)
	if err != nil {
		return err
	}
	a := &Access{Resource: resource, RecordID: id, Subject: subject, RoleID: roleID, GrantedBy: by,
		CreatedAt: s.o.Now(), UpdatedAt: s.o.Now()}
	if err := tx.Create(a).Error; err != nil {
		return err
	}
	if err := s.accessChanged(tx, a, role, false); err != nil {
		return err
	}
	return s.event(tx, Event{Resource: resource, RecordID: id, Subject: subject, Event: EventGranted,
		RoleAfter: &role.ID, ActionsAfter: role.Actions, By: by}, o)
}

// Invites are the invites pending of the user.
func (s *Service) Invites(ctx context.Context, user string) ([]AccessInvite, error) {
	var out []AccessInvite
	return out, s.o.DB.WithContext(ctx).Where("subject = ? AND state = ?", UserSubject(user), StatePending).
		Order("created_at DESC").Find(&out).Error
}

// Accept is the user accepting the invite: the access; of a new owner, the
// record theirs, the former owner with the role kept.
func (s *Service) Accept(ctx context.Context, inviteID uuid.UUID, user string, o Origin) error {
	return s.o.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		inv, err := s.invite(tx, inviteID, user)
		if err != nil {
			return err
		}
		now := s.o.Now()
		if err := tx.Model(inv).Updates(map[string]any{"state": StateAccepted, "answered_at": &now}).Error; err != nil {
			return err
		}
		if err := s.event(tx, Event{Resource: inv.Resource, RecordID: inv.RecordID, Subject: inv.Subject, Event: EventAccepted,
			By: user, Note: inv.ID.String()}, o); err != nil {
			return err
		}
		if inv.Kind != RequestOwnerTransfer {
			return s.grant(tx, inv.Resource, inv.RecordID, inv.Subject, *inv.RoleID, inv.By, o)
		}
		return s.transfer(tx, inv, user, o)
	})
}

// transfer makes user the owner of the record of inv.
func (s *Service) transfer(tx *gorm.DB, inv *AccessInvite, user string, o Origin) error {
	r, err := s.resource(inv.Resource)
	if err != nil {
		return err
	}
	if r.OwnerColumn == "" {
		return ErrNoOwner
	}
	l := level{r, inv.RecordID}
	former, err := ownerIn(tx, l)
	if err != nil {
		return err
	}
	if err := tx.Table(r.Table).Where(r.IDColumn+" = ?", inv.RecordID).Update(r.OwnerColumn, user).Error; err != nil {
		return err
	}
	// the new owner's access, if any, is no more: they own it
	var mine Access
	if tx.Preload("Role").Where("resource = ? AND record_id = ? AND subject = ?", inv.Resource, inv.RecordID, UserSubject(user)).
		Limit(1).Find(&mine).RowsAffected > 0 {
		if err := tx.Delete(&mine).Error; err != nil {
			return err
		}
		if err := s.accessChanged(tx, &mine, mine.Role, true); err != nil {
			return err
		}
	}
	if err := s.event(tx, Event{Resource: inv.Resource, RecordID: inv.RecordID, Subject: UserSubject(user),
		Event: EventOwnerTransferred, By: user, Note: "from " + former}, o); err != nil {
		return err
	}
	if former != "" && inv.RoleID != nil {
		return s.grant(tx, inv.Resource, inv.RecordID, UserSubject(former), *inv.RoleID, user, o)
	}
	return nil
}

// Decline is the user refusing the invite.
func (s *Service) Decline(ctx context.Context, inviteID uuid.UUID, user string, o Origin) error {
	return s.o.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		inv, err := s.invite(tx, inviteID, user)
		if err != nil {
			return err
		}
		now := s.o.Now()
		if err := tx.Model(inv).Updates(map[string]any{"state": StateRejected, "answered_at": &now}).Error; err != nil {
			return err
		}
		return s.event(tx, Event{Resource: inv.Resource, RecordID: inv.RecordID, Subject: inv.Subject, Event: EventRejected,
			By: user, Note: inv.ID.String()}, o)
	})
}

func (s *Service) invite(tx *gorm.DB, id uuid.UUID, user string) (*AccessInvite, error) {
	var inv AccessInvite
	if err := tx.First(&inv, "id = ?", id).Error; err != nil {
		return nil, ErrNotPending
	}
	if inv.Subject != UserSubject(user) {
		return nil, ErrNotYours
	}
	if inv.State != StatePending {
		return nil, ErrNotPending
	}
	return &inv, nil
}

// MemberChanged records, in the history of every record shared with the
// group, that the user joined (added) or left it — their access, the
// group's, at once.
func (s *Service) MemberChanged(ctx context.Context, group, user string, added bool, by string, o Origin) error {
	var acc []Access
	db := s.o.DB.WithContext(ctx)
	if err := db.Where("subject = ?", GroupSubject(group)).Find(&acc).Error; err != nil {
		return err
	}
	ev := EventMemberRemoved
	if added {
		ev = EventMemberAdded
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for _, a := range acc {
			if err := s.event(tx, Event{Resource: a.Resource, RecordID: a.RecordID, Subject: GroupSubject(group), Event: ev,
				By: by, Note: UserSubject(user)}, o); err != nil {
				return err
			}
		}
		return nil
	})
}
