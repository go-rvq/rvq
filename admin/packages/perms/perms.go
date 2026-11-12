// Package perms provides reusable, resource-agnostic permission tooling for
// rvq presets admins: a per-record permission manager (a shortcut to the perm
// API, backed by perm.DefaultDBPolicy) and a soft-delete trash view gated by a
// 'trash' listing permission. Both work with any ModelBuilder.
package perms

import (
	"strings"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission verbs granted per record. They map to the presets action names
// checked by the verifier.
var (
	VerbView = []string{presets.PermList, presets.PermGet}
	VerbEdit = []string{presets.PermUpdate}
	// who may delete a record may bring it back from the trash
	VerbDelete = []string{presets.PermDelete, presets.PermRestore}
)

// PermTrash is the listing permission verb that gates the trash view and the
// restore of its records (presets.PermTrash); PermRestore, the record's verb
// that brings it back (presets.PermRestore).
const (
	PermTrash   = presets.PermTrash
	PermRestore = presets.PermRestore
)

// PermManageAction is the name of the detail action of the permission
// manager; PermManage, the perm verb that guards it — the action's own
// (presets.ActionPerm of its name), so the permission tree grants the verb
// that is checked.
const PermManageAction = "manage_permissions"

// PermManage is the perm verb that guards the permission manager itself.
var PermManage = presets.ActionPerm(PermManageAction)

// RecordResource returns the exact permission resource string of a record, as
// computed by the presets permissioner (the module "presets" is already seeded
// into Resource(), so it is not added again). For a nested resource, pass the
// parent ids (the same the request carries) so the resource matches the runtime
// check.
func RecordResource(mb *presets.ModelBuilder, id model.ID, parentID ...model.ID) string {
	return uniqueResource(mb.Permissioner().Verifier(id, parentID...))
}

// uniqueResource is the resource of v by the unique name of what it is of —
// which decides first, and does not change when a model moves in the menu
// (perm.Verifier.Prefer) —, else its resource.
func uniqueResource(v *perm.Verifier) string {
	if r := v.PreferredResource(); r != "" {
		return r
	}
	return v.Resource()
}

// ListResource returns the permission resource of the whole listing. For a
// nested resource, pass the parent ids.
func ListResource(mb *presets.ModelBuilder, parentID ...model.ID) string {
	return uniqueResource(mb.Permissioner().ListVerifier(parentID...))
}

// Grant creates or updates a DB policy granting subject the actions on
// resource, keyed by referID (so it can be updated or revoked later).
func Grant(db *gorm.DB, referID, subject, resource string, actions []string) (*perm.DefaultDBPolicy, error) {
	return keep(db, referID, subject, resource, actions, perm.Allowed)
}

// DenyShared is like GrantShared, denying: the subject may none of actions
// on resource — whatever else allows them.
func DenyShared(db *gorm.DB, referID, subject, resource string, actions []string, sharedID uuid.UUID) (*perm.DefaultDBPolicy, error) {
	p, err := keep(db, referID, subject, resource, actions, perm.Denied)
	if err != nil {
		return nil, err
	}
	p.SharedID = model.UUID(sharedID)
	return p, db.Model(p).Update("shared_id", sharedID).Error
}

// keep creates or updates the policy of referID: effect for subject on
// resource.
func keep(db *gorm.DB, referID, subject, resource string, actions []string, effect string) (*perm.DefaultDBPolicy, error) {
	uuidkey.MustRegister(db) // policies have UUID keys
	// the resource asked ends in its verb ("…:<7>:@edit"): a resource that
	// ends a part holds the verbs after it
	if strings.HasSuffix(resource, ":") {
		resource += "*"
	}
	p := &perm.DefaultDBPolicy{}
	err := db.Where("refer_id = ?", referID).First(p).Error
	switch err {
	case nil:
		p.Subject = subject
		p.Effect = effect
		p.Actions = actions
		p.Resources = []string{resource}
		return p, db.Save(p).Error
	case gorm.ErrRecordNotFound:
		p = &perm.DefaultDBPolicy{
			ReferID:   referID,
			Subject:   subject,
			Effect:    effect,
			Actions:   actions,
			Resources: []string{resource},
		}
		return p, db.Create(p).Error
	default:
		return nil, err
	}
}

// GrantShared is like Grant but stamps the policy with sharedID, marking it as
// part of a resource share (organization/project). On update the SharedID is
// (re)applied so a policy adopted by a share is tracked.
func GrantShared(db *gorm.DB, referID, subject, resource string, actions []string, sharedID uuid.UUID) (*perm.DefaultDBPolicy, error) {
	p, err := Grant(db, referID, subject, resource, actions)
	if err != nil {
		return nil, err
	}
	p.SharedID = model.UUID(sharedID)
	return p, db.Model(p).Update("shared_id", sharedID).Error
}

// Revoke soft-deletes the policy identified by referID.
func Revoke(db *gorm.DB, referID string) error {
	return db.Where("refer_id = ?", referID).Delete(&perm.DefaultDBPolicy{}).Error
}

// RevokeShared soft-deletes every policy belonging to the given share.
func RevokeShared(db *gorm.DB, sharedID uuid.UUID) error {
	return db.Where("shared_id = ?", sharedID).Delete(&perm.DefaultDBPolicy{}).Error
}

// ListShared returns the active policies belonging to the given share, ordered
// by subject.
func ListShared(db *gorm.DB, sharedID uuid.UUID) (out []perm.DefaultDBPolicy, err error) {
	err = db.Where("shared_id = ?", sharedID).Order("subject").Find(&out).Error
	return
}

// IsShared reports whether the policy identified by referID carries a SharedID
// (and therefore may only be removed through the sharing UI, not the permission
// manager).
func IsShared(db *gorm.DB, referID string) (bool, error) {
	var p perm.DefaultDBPolicy
	if err := db.Select("shared_id").Where("refer_id = ?", referID).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return !p.SharedID.IsZero(), nil
}

// List returns the active policies whose resources contain resource.
func List(db *gorm.DB, resource string) (out []perm.DefaultDBPolicy, err error) {
	err = db.Where("? = ANY(resources)", resource).Order("subject").Find(&out).Error
	return
}
